package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type keyRing interface {
	Load() error
	Save() error
	AddKey(name string, keyID []byte, key []byte)
	Key(keyID []byte) ([]byte, error)
}

type fileKeyRing struct {
	fileName   string
	KeyEntries []keyEntry
}

type keyEntry struct {
	Description string `yaml:"description"`
	KeyID       string `yaml:"key-id"`
	Key         string `yaml:"key"`
}

func (kr *fileKeyRing) AddKey(desc string, keyID []byte, key []byte) {
	kr.KeyEntries = append(kr.KeyEntries, keyEntry{
		Description: desc,
		KeyID:       string(encode(keyID[:])),
		Key:         string(encode(key[:])),
	})
}

func (kr *fileKeyRing) Key(keyID []byte) ([]byte, error) {
	b64 := string(encode(keyID[:]))

	for _, ke := range kr.KeyEntries {
		if ke.KeyID == b64 {
			dec, err := decode([]byte(ke.Key))
			if err != nil {
				return []byte{}, err
			}
			if len(dec) != 32 {
				return []byte{}, fmt.Errorf("unexpected length of key: %d", len(dec))
			}
			return dec, nil
		}
	}

	return []byte{}, errKeyNotFound
}

func (kr *fileKeyRing) Load() error {

	bytes, err := ioutil.ReadFile(kr.fileName)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, kr)
	return err
}

func (kr *fileKeyRing) Save() error {
	ser, err := yaml.Marshal(kr)
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Dir(kr.fileName)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0700)
		if err != nil {
			return fmt.Errorf("error creating strongbox home folder: %s", err)
		}
	}

	return ioutil.WriteFile(kr.fileName, ser, 0600)
}
