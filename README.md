# gh-get

gh-get is a [GitHub CLI](https://cli.github.com/) based re-implementation of [ghq](https://github.com/x-motemen/ghq).

## Installing the extension

`gh extension install britter/gh-get`

## Usage

`gh get OWNER/REPO`

This will clone the repository identified by OWNER/REPO into `$HOME/github/$OWNER/$REPO`.

Example:

`gh get britter/gh-get` will clone this respository into `~/github/britter/gh-get`

## Configuration

There are two environment variables that control the location gh-get clones repositories to:

- `GH_GET_FOLDER`: The name of the folder inside the user home that repositories are cloned into. The defaults to `github` which means repositories are cloned into `$HOME/github`. If you prefer a different folder name or path inside your user home, configure this variable.
- `GH_GET_ROOT`: The full path to to the folder that repositories are cloned into. This defaults to `$HOME/$GH_GET_FOLDER`. Change this is you want repositories to be cloned to a location _outside_ your user home. 

## Building

### Using Nix

This project uses the [nix package manager](https://nixos.org) to define a development environment that has everything needed to build gh-get.
In order to enter the dev shell execute `nix develop`.
Inside the dev shell you can build gh-get by running `go build`.

### Without Nix

Make sure you have a Go toolchain installed in your path that matches the version defined in [go.mod](go.mod).
Build gh-get by running `go build`.

## License

Code is under the [Apache Licence v2](https://www.apache.org/licenses/LICENSE-2.0.txt).
