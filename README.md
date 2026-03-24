# gh-get

gh-get is a [GitHub CLI](https://cli.github.com/) based re-implementation of [ghq](https://github.com/x-motemen/ghq).

## Installing the extension

`gh extension install britter/gh-get`

## Usage

`gh get [--fork] OWNER/REPO`

This will clone the repository identified by OWNER/REPO into `$HOME/github/$OWNER/$REPO`.

Example:

`gh get britter/gh-get` will clone this respository into `~/github/britter/gh-get`

> [!TIP]
> gh-get also accepts full GitHub URLs like `https://github.com/britter/gh-get` as well as URLs pointing to branches, tags, commits, files, pull requests, or issues — anything you can copy from your browser.

### Forking

When you don't have write access to a repository, gh-get will ask whether you want to fork it first. If you answer yes, the fork is created under your account and cloned into `$HOME/github/$YOUR_USERNAME/$REPO`.

You can also pass `--fork` to skip the prompt and always fork:

`gh get --fork OWNER/REPO`

If the repository does not allow forking, the original is cloned and a warning is printed.

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

## Testing

The project contains two sorts of tests:

- Unit tests written in go testing individual functions. You can execute these by running `go test ./...`
- Integration tests running inside a docker container. You can execute these by running `cd integration-tests && sudo GH_TOKEN=<personal access token> docker compose up --build --abort-on-container-exit`.

To get a personal access token, you can either look the one currently used by your GitHub CLI installation up in `~/.config/gh/hosts.yaml` or generate a new one in your [GitHub settings](https://github.com/settings/tokens).

> [!IMPORTANT]
> Before submitting pull requests make sure to run all tests and add new test cases for the functionality you've added.

## License

Code is under the [Apache Licence v2](https://www.apache.org/licenses/LICENSE-2.0.txt).
