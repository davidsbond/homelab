package main

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jacobsa/crypto/siv"
)

const version = "v0.3.1"

var (
	keyLoader     = keyPair
	kr            keyRing
	prefix        = []byte("# STRONGBOX ENCRYPTED RESOURCE ;")
	defaultPrefix = "# STRONGBOX ENCRYPTED RESOURCE ; See https://github.com/uw-labs/strongbox\n# key-id: %s\n"

	// Match lines *not* starting with `#`
	// this should match ciphertext without the strongbox prefix
	prefixStripRegex = regexp.MustCompile(`(?m)^[^#]+$`)

	keyIDRegex = regexp.MustCompile(`key-id: (\w+)`)

	errKeyNotFound            = errors.New("key not found")
	errKeyIDMissingFromHeader = errors.New("strongbox header doesn't contain key-id")

	// flags
	flagGitConfig = flag.Bool("git-config", false, "Configure git for strongbox use")
	flagGenKey    = flag.String("gen-key", "", "Generate a new key and add it to your strongbox keyring")
	flagDecrypt   = flag.Bool("decrypt", false, "Decrypt single resource")
	flagKey       = flag.String("key", "", "Private key to use to decrypt")

	flagClean  = flag.String("clean", "", "intended to be called internally by git")
	flagSmudge = flag.String("smudge", "", "intended to be called internally by git")
	flagDiff   = flag.String("diff", "", "intended to be called internally by git")

	flagVersion = flag.Bool("version", false, "Strongbox version")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n\n")
	fmt.Fprintf(os.Stderr, "\tstrongbox -git-config\n")
	fmt.Fprintf(os.Stderr, "\tstrongbox -gen-key key-name\n")
	fmt.Fprintf(os.Stderr, "\tstrongbox -decrypt\n")
	fmt.Fprintf(os.Stderr, "\tstrongbox -key\n")
	fmt.Fprintf(os.Stderr, "\tstrongbox -version\n")
	os.Exit(2)
}

func main() {
	log.SetPrefix("strongbox: ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Usage = usage
	flag.Parse()

	// Set up keyring file name
	home := deriveHome()

	kr = &fileKeyRing{fileName: filepath.Join(home, ".strongbox_keyring")}

	if *flagVersion || (flag.NArg() == 1 && flag.Arg(0) == "version") {
		fmt.Println(version)
		return
	}

	if *flagGitConfig {
		gitConfig()
		return
	}

	if *flagGenKey != "" {
		genKey(*flagGenKey)
		return
	}

	if *flagDecrypt {
		if *flagKey == "" {
			log.Fatalf("Must provide a key when using -decrypt")
		}
		decryptCLI()
		return
	}

	if *flagClean != "" {
		clean(os.Stdin, os.Stdout, *flagClean)
		return
	}
	if *flagSmudge != "" {
		smudge(os.Stdin, os.Stdout, *flagSmudge)
		return
	}
	if *flagDiff != "" {
		diff(*flagDiff)
		return
	}
}

func deriveHome() string {
	// try explicitly set STRONGBOX_HOME
	if home := os.Getenv("STRONGBOX_HOME"); home != "" {
		return home
	}
	// Try user.Current which works in most cases, but may not work with CGO disabled.
	u, err := user.Current()
	if err == nil && u.HomeDir != "" {
		return u.HomeDir
	}
	// try HOME env var
	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	log.Fatal("Could not call os/user.Current() or find $STRONGBOX_HOME or $HOME. Please recompile with CGO enabled or set $STRONGBOX_HOME or $HOME")
	// not reached
	return ""
}

func decryptCLI() {
	var fn string
	if flag.Arg(0) == "" {
		// no file passed, try to read stdin
		fn = "/dev/stdin"
	} else {
		fn = flag.Arg(0)
	}
	fb, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatalf("Unable to read file to decrypt %v", err)
	}
	dk, err := decode([]byte(*flagKey))
	if err != nil {
		log.Fatalf("Unable to decode private key %v", err)
	}
	out, err := decrypt(fb, dk)
	if err != nil {
		log.Fatalf("Unable to decrypt %v", err)
	}
	fmt.Printf("%s", out)
}

func gitConfig() {
	args := [][]string{
		{"config", "--global", "--replace-all", "filter.strongbox.clean", "strongbox -clean %f"},
		{"config", "--global", "--replace-all", "filter.strongbox.smudge", "strongbox -smudge %f"},
		{"config", "--global", "--replace-all", "filter.strongbox.required", "true"},

		{"config", "--global", "--replace-all", "diff.strongbox.textconv", "strongbox -diff"},
	}
	for _, command := range args {
		cmd := exec.Command("git", command...)
		if out, err := cmd.CombinedOutput(); err != nil {
			log.Fatal(string(out))
		}
	}
	log.Println("git global configuration updated successfully")
}

func genKey(desc string) {
	err := kr.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		log.Fatal(err)
	}

	keyID := sha256.Sum256(key)

	kr.AddKey(desc, keyID[:], key)

	err = kr.Save()
	if err != nil {
		log.Fatal(err)
	}
}

func diff(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	_, err = io.Copy(os.Stdout, f)
	if err != nil {
		log.Fatal(err)
	}
}

func clean(r io.Reader, w io.Writer, filename string) {
	// Read the file, fail on error
	in, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	// Check the file is plaintext, if its an encrypted strongbox file, copy as is, and exit 0
	if bytes.HasPrefix(in, prefix) {
		_, err = io.Copy(w, bytes.NewReader(in))
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	// File is plaintext and needs to be encrypted, get the key, fail on error
	keyID, key, err := keyLoader(filename)
	if err != nil {
		log.Fatal(err)
	}
	// encrypt the file, fail on error
	out, err := encrypt(in, key, keyID)
	if err != nil {
		log.Fatal(err)
	}
	// write out encrypted file, fail on error
	_, err = io.Copy(w, bytes.NewReader(out))
	if err != nil {
		log.Fatal(err)
	}
}

// Called by git on `git checkout`
func smudge(r io.Reader, w io.Writer, filename string) {
	in, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	// file is a non-strongbox file, copy as is and exit
	if !bytes.HasPrefix(in, prefix) {
		_, err = io.Copy(w, bytes.NewReader(in))
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// try to get the key using the header, failing that, try to get the
	// key using the filename
	var key []byte
	key, err = keyFromHeader(in)
	if err != nil {
		_, key, err = keyLoader(filename)
		if err != nil {
			// don't log error if its keyNotFound
			switch err {
			case errKeyNotFound:
			default:
				log.Println(err)
			}
			// Couldn't load the key, just copy as is and return
			if _, err = io.Copy(w, bytes.NewReader(in)); err != nil {
				log.Println(err)
			}
			return
		}
	}

	out, err := decrypt(in, key)
	if err != nil {
		log.Println(err)
		out = in
	}
	if _, err := io.Copy(w, bytes.NewReader(out)); err != nil {
		log.Println(err)
	}
}

// keyFromHeader looks through the file content, trying to get key-id value,
// and look up the key in the keyring
func keyFromHeader(in []byte) ([]byte, error) {
	match := keyIDRegex.FindStringSubmatch(string(in))
	if len(match) != 2 {
		return []byte{}, errKeyIDMissingFromHeader
	}
	decodedKeyID, _ := decode([]byte(match[1]))
	key, err := kr.Key(decodedKeyID)
	//log.Printf("DEBUG: found key %s %e", encode(key), err)
	if err != nil {
		return []byte{}, err
	}
	return key, nil
}

func encrypt(b []byte, key, keyID []byte) ([]byte, error) {
	b = compress(b)
	out, err := siv.Encrypt(nil, key, b, nil)
	if err != nil {
		return nil, err
	}
	var buf []byte
	p := fmt.Sprintf(defaultPrefix, encode(keyID))
	buf = append(buf, []byte(p)...)

	b64 := encode(out)
	for len(b64) > 0 {
		l := 76
		if len(b64) < 76 {
			l = len(b64)
		}
		buf = append(buf, b64[0:l]...)
		buf = append(buf, '\n')
		b64 = b64[l:]
	}
	return buf, nil
}

func compress(b []byte) []byte {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(b)
	if err != nil {
		log.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func decompress(b []byte) []byte {
	zr, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	b, err = ioutil.ReadAll(zr)
	if err != nil {
		log.Fatal(err)
	}
	if err := zr.Close(); err != nil {
		log.Fatal(err)
	}
	return b
}

func encode(decoded []byte) []byte {
	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(decoded)))
	base64.StdEncoding.Encode(b64, decoded)
	return b64
}

func decode(encoded []byte) ([]byte, error) {
	decoded := make([]byte, len(encoded))
	i, err := base64.StdEncoding.Decode(decoded, encoded)
	if err != nil {
		return nil, err
	}
	return decoded[0:i], nil
}

func decrypt(enc []byte, priv []byte) ([]byte, error) {
	// strip the prefix (both single line v0.1 and multiline v0.2)
	ciphertext := prefixStripRegex.Find(enc)
	if ciphertext == nil {
		return nil, errors.New("Couldn't split strongbox prefix and ciphertext")
	}
	b64decoded, err := decode(ciphertext)
	if err != nil {
		return nil, err
	}
	decrypted, err := siv.Decrypt(priv, b64decoded, nil)
	if err != nil {
		return nil, err
	}
	decrypted = decompress(decrypted)
	return decrypted, nil
}

// keyPair returns public, private and error
func keyPair(filename string) ([]byte, []byte, error) {
	keyID, err := findKey(filename)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	err = kr.Load()
	if err != nil {
		return []byte{}, []byte{}, err
	}

	key, err := kr.Key(keyID)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	return keyID, key, nil
}

func findKey(filename string) ([]byte, error) {
	path := filepath.Dir(filename)
	for {
		if fi, err := os.Stat(path); err == nil && fi.IsDir() {
			keyFilename := filepath.Join(path, ".strongbox-keyid")
			if keyFile, err := os.Stat(keyFilename); err == nil && !keyFile.IsDir() {
				return readKey(keyFilename)
			}
		}
		if path == "." {
			break
		}
		path = filepath.Dir(path)
	}
	return []byte{}, fmt.Errorf("failed to find key id for file %s", filename)
}

func readKey(filename string) ([]byte, error) {
	fp, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	b64 := strings.TrimSpace(string(fp))
	b, err := decode([]byte(b64))
	if err != nil {
		return []byte{}, err
	}
	if len(b) != 32 {
		return []byte{}, fmt.Errorf("unexpected key length %d", len(b))
	}
	return b, nil
}
