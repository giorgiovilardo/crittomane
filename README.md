# crittomane

Quick, hacky command to help me save private obsidian notes in github.

## Install

```shell
git clone git@github.com:giorgiovilardo/crittomane.git
go mod download
go build -ldflags="-s -w" -trimpath -o ./cm main.go
chmod 0755 cm
sudo mv cm /usr/local/bin
```

## Usage

⚠️⚠️⚠️ WARNING ⚠️⚠️⚠️: ALL COMMANDS TRUNCATE/OVERWRITE EXISTING DATA!

### Encrypt the `./private` directory into `private.ctm`

```shell
cm e yourPassword
```

### Decrypt `private.ctm` into `./private`

```shell
cm d yourPassword
```

## Roadmap

Mostly accepting paths instead of hardcoding everything. And tests.

But it will probably stay like this as it works for my usecase :)
