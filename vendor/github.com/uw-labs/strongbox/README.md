![Strongbox](strongbox-logo.png)

Encryption for git users.

Strongbox makes it easy to encrypt and decrypt files stored in git, with minimal
divergence from a typical git workflow.  Once installed, strongbox enables
normal use of commands such as `git diff` etc. and all of the files that should
be encrypted in the repository remain decrypted on your working copy.

It supports use of different keys per directory if wanted. It can cover as many
or as few files as you wish based on
[.gitattributes](https://www.git-scm.com/docs/gitattributes)


<!-- vim-markdown-toc GFM -->

* [Installation](#installation)
* [Usage](#usage)
* [Existing project](#existing-project)
* [Verification](#verification)
* [Key rotation](#key-rotation)
* [Security](#security)
* [Testing](#testing)
* [Known issues](#known-issues)
	* [Clone file ordering](#clone-file-ordering)
		* [Workaround](#workaround)

<!-- vim-markdown-toc -->

## Installation

You can obtain a binary from https://github.com/uw-labs/strongbox/releases

Alternatively, assuming you have a working [Go](https://golang.org) installation, you can
install via `go get github.com/uw-labs/strongbox`

## Usage

1. As a one time action, install the plugin by running `strongbox -git-config`.
   This will edit global git config to enable strongbox filter and diff
   configuration.

2. In each repository you want to use strongbox, create `.gitattributes` file
   containing the patterns to be managed by strongbox.

   For example:

   ```
   secrets/* filter=strongbox diff=strongbox
   ```

3. Generate a key to use for the encryption, for example:
   ```
   strongbox -gen-key my-key
   ```
   This will add a new key to your `.strongbox_keyring`. By default, the
   keyring is created in the `$HOME` directory, but this location can be changed
   by setting the `$STRONGBOX_HOME` environmental variable.

4. Include a `.strongbox-keyid` file in your repository containing public key
   you want to use (typically by copying a public key from
   `$HOME/.strongbox_keyring` )  This can be in the same directory as the
   protected resource(s) or any parent directory.   When searching for
   `.strongbox-keyid` for a given resource, strongbox will recurse up the
   directory structure until it finds the file.  This allows using different
   keys for different subdirectories within a repository.

## Existing project

Strongbox uses [clean and smudge
filters](https://git-scm.com/book/en/v2/Customizing-Git-Git-Attributes#filters_a)
to encrypt and decrypt files.

If you are cloning a project that uses strongbox, you will need to place the
key into your keyring file prior to cloning (checkout). Otherwise that filter
will fail and not decrypt files on checkout.

If you already have the project locally and added the keys, you can remove and
checkout the files to force the filter:
```
rm <files> && git checkout -- <files>
```

## Verification

Following a `git add`, you can verify the file is encrypted in the index:

```
git show :/path/to/file
```

Verify a file is encrypted in the commit:

```
git show HEAD:/path/to/file
```

What you should see is a Strongbox encrypted resource, and this is what would
be pushed to the remote.

Compare an entire branch (as it would appear on the remote) to master:

```
git diff-index -p master
```

## Key rotation

To rotate keys, update the `.strongbox-keyid` with the new key id, then `touch`
all files/directories covered by `.gitattributes`. All affected files should now
show up as "changed".

## Security

Strongbox uses SIV-AES as defined in rfc5297 in order to achieve authenticated
deterministic encryption.

## Testing

Run integration tests:

```
make test
```

## Known issues

### Clone file ordering

Given a `.strongbox-keyid` in the root of the repository and an encrypted file
in the same directory,*and* alphabetically it comes before the key-id file.

Git checks out files alphanumerically, so if the strongboxed file is being
checked out before the `.strongbox-keyid` is present on disk, strongbox will
fail to find the decryption key.

Order of files being cloned is dictated by the index.

#### Workaround

Clone repository, let the descryption fail. Delete encrypted files and do `git
checkout` on the deleted files.
