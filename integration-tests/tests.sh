!#/bin/sh

echo "Installing gh-get"
cd /workspace/gh-get
gh extension install .

test() {
	local repository_definition="$1"
	local expected_path="$2"

	echo "gh get $repository_definition"
	gh get "$repository_definition"

	if [[ -d "$expected_path" ]]; then
		echo "SUCCESS"
		echo ""
		rm -rf "$expected_path"
	else
		echo "FAILURE"
		exit 1
	fi
}

test "britter/gh-get" "/root/github/britter/gh-get"
test "https://github.com/britter/gh-get" "/root/github/britter/gh-get"
test "https://github.com/britter/gh-get.git" "/root/github/britter/gh-get"

export GH_GET_FOLDER=src
test "britter/gh-get" "/root/src/britter/gh-get"
unset GH_GET_FOLDER

export GH_GET_ROOT=/repositories
test "britter/gh-get" "/repositories/britter/gh-get"

