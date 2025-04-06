# crittomane

Quick, hacky command to help me save private obsidian notes in github.

## Install

```shell
git clone git@github.com:giorgiovilardo/crittomane.git
# optional, if you use docker and want to cache dep download
go mod download
go build -ldflags="-s -w" -trimpath -o ./cm main.go
chmod 0755 cm
sudo mv cm /usr/local/bin
```

## Usage

⚠️⚠️⚠️ WARNING ⚠️⚠️⚠️: ALL COMMANDS TRUNCATE/OVERWRITE EXISTING DATA!

⚠️⚠️⚠️ WARNING ⚠️⚠️⚠️: `crittomane` SHOULD BE SAFE AS LONG AS YOUR PASSWORD IS GOOD; STILL DON'T USE IT FOR REALLY PRIVATE
THINGS!

`crittomane` spawns from the very basic need of having a bit of privacy for my
[open-notes](https://github.com/giorgiovilardo/open-notes) repository.
I mainly use markdown to jot down stuff, but there are some things I keep secluded in a `private`
directory inside the main repo just because they're simply not ready for public consumption.
Rather than having a public and a private repo or rely on some obsidian extension that does the same but with less control,
I made this super small utility which has 2 modes of usage:

* `e,` that encrypts, if present, a `private` directory child of the directory in which you run the command and stores everything inside a `private.ctm` file;
* `d,` that decrypts, if present, a `private.ctm` file, unpacking everything into a `private` directory.

This allows me to commit the `private.ctm` file inside the repo, knowing it's reasonably secure as long as I use a decent password.

In any case, never store anything truly sensitive into this file.

### Encrypt the `./private` directory into `private.ctm`

```shell
# yourPassword is optional, as if not there it will be parsed by stdin read
cm e yourPassword
```

### Decrypt `private.ctm` into `./private`

```shell
# yourPassword is optional, as if not there it will be parsed by stdin read
cm d yourPassword
```

## Roadmap

Mostly accepting paths instead of hardcoding everything. And tests.

But it will probably stay like this as it works for my usecase :)
