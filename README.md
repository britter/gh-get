# gh-get

A GitHub CLI based copy of https://github.com/x-motemen/ghq

## Installing the extension

`gh extension install britter/gh-get`

## Usage

`gh get OWNER/REPO`

This will clone the repository identified by OWNER/REPO into `$HOME/github/$OWNER/$REPO`.

Example:

`gh get britter/gh-get` will clone this respository into `~/workspace/britter/gh-get`

## Configuration

There are two environment variables that control the location gh-get clones repositories to:

- `GH_GET_FOLDER`: The name of the folder inside the user home that repositories are cloned into. The defaults to `github` which means repositories are cloned into `$HOME/github`. If you prefer a different folder name or path inside your user home, configure this variable.
- `GH_GET_ROOT`: The full path to to the folder that repositories are cloned into. This defaults to `$HOME/$GH_GET_FOLDER`. Change this is you want repositories to be cloned to a location _outside_ your user home. 

## License

Code is under the [Apache Licence v2](https://www.apache.org/licenses/LICENSE-2.0.txt).
