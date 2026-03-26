#!/bin/sh

echo "Installing gh-get"
cd /workspace/gh-get
gh extension install .

test() {
	local expected_path="$1"
	shift

	echo "gh get $*"
	gh get "$@"

	if [[ -d "$expected_path" ]]; then
		echo "SUCCESS"
		echo ""
		rm -rf "$expected_path"
	else
		echo "FAILURE"
		exit 1
	fi
}

test "/root/github/britter/gh-get" --fork=false britter/gh-get
test "/root/github/britter/gh-get" --fork=false https://github.com/britter/gh-get
test "/root/github/britter/gh-get" --fork=false https://github.com/britter/gh-get.git

export GH_GET_FOLDER=src
test "/root/src/britter/gh-get" --fork=false britter/gh-get
unset GH_GET_FOLDER

export GH_GET_ROOT=/repositories
test "/repositories/britter/gh-get" --fork=false britter/gh-get
unset GH_GET_ROOT

