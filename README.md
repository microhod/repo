# Repo

`repo` is a cli application to organise scm repositories in a structured hierarchy.

For example, `github.com/microhod/repo.git` would be stored at `~/src/github.com/microhod/repo`.

## Usage

Run `repo --help` to see full CLI usage.

```
NAME:
   repo - A cli application to organise scm repositories in a structured heirachy

USAGE:
   repo [global options] command [command options] [arguments...]

COMMANDS:
   clone     clone a repo
   organise  organise all repos under the current path into a structured heirachy
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

## Install

### Compile from Source

* Install [golang](https://golang.org/doc/install) `1.22` or later
* Run `go install github.com/microhod/repo@latest` to download and install the binary (this will install to `~/go/bin`)

## Config

Configuration will be auto-generated at first startup and stored at `~/.config/repo/config.json`.

The default configuration is as below:

```json
{
    "remote": {
        "default": {
            "prefix": "ssh://git@github.com"
        }
    },
    "local": {
        "root": "~/src"
    }
}
```

## Todo

- [ ] More complete docs
- [ ] Tests!!!
- [ ] `--dry-run` option for `organise`
- [ ] `--verbose` option for all commands
- [ ] add `profile` command generate terminal profile for utlity commands
