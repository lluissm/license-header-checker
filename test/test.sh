#!/bin/bash

# Test config
CMD="$1"
PROJECT_DIR=sample-project
CMD_ARGS="${PROJECT_DIR}/licenses/current-license.txt ${PROJECT_DIR} java js cpp go"

# Test results
errors=0

# Colors
RED='\033[0;31m'
GREEN='\033[32m'
NC='\033[0m'

remove_ansi_color() {
	if [[ "$(uname)" == "Darwin" ]]; then
		echo $(echo $1 | sed $'s,\x1b\\[[0-9;]*[a-zA-Z],,g')
	else
		echo $(echo $1 | sed -r "s/\x1B\[([0-9]{1,3}(;[0-9]{1,2})?)?[mGK]//g")
	fi
}

extract_sample_project() {
	tar -xzf $PROJECT_DIR.tar.gz
}

delete_sample_project() {
	rm -rf $PROJECT_DIR
}

on_success() {
	echo -e "${GREEN}OK${NC}"
}

on_failure() {
	((errors++))
	echo -e "${RED}Error${NC}"
}

run_test() {
	flags=$1
	expected=$3

	# extract sample project
	delete_sample_project
	extract_sample_project

	# print test case
	echo -e "\n$2"

	# execute license-header-checker and remove colors
	output=$(remove_ansi_color "$($CMD $flags $CMD_ARGS)")

	# verify result
	if [[ "$output" =~ "$expected" ]]; then
		on_success
	else
		on_failure
	fi
}

# Test cases

# version
flags="-version"
test_case="Testing version..."
expected_output="version:"
run_test "${flags}" "${test_case}" "$expected_output"

# read only
flags=''
test_case='Testing read only...'
expected_output="1 licenses ok, 0 licenses replaced, 0 licenses added \
[!] 2 files had no license but were not changed as the -a (add) option was not supplied. \
[!] 2 files had a different license but were not changed as the -r (replace) option was not supplied."
run_test "$flags" "$test_case" "$expected_output"

# add
flags='-a'
test_case='Testing with -a flag...'
expected_output="1 licenses ok, 0 licenses replaced, 2 licenses added \
[!] 2 files had a different license but were not changed as the -r (replace) option was not supplied."
run_test "$flags" "$test_case" "$expected_output"

# add and replace
flags='-a -r'
test_case='Testing with -a and -r flags...'
expected_output="1 licenses ok, 2 licenses replaced, 2 licenses added"
run_test "$flags" "$test_case" "$expected_output"

# add and replace (with ignore)
flags='-a -r -i src/other'
test_case='Testing with -a and -r and -i flags...'
expected_output="1 licenses ok, 2 licenses replaced, 1 licenses added"
run_test "$flags" "$test_case" "$expected_output"

# add and replace with ignore (and verbose)
flags='-a -r -v -i src/other'
test_case='Testing with -a and -r and -i and -v flags...'
expected_output="\
files: license_ok: - sample-project/src/file-with-license.js license_replaced: - sample-project/test/file-with-old-license.go - sample-project/src/file-with-old-license.cpp license_added: - sample-project/src/file-without-license.java \
options: project_path: sample-project ignore_paths: - src/other extensions: - .java - .js - .cpp - .go flags: - add - replace - verbose license_header: %ssample-project/licenses/current-license.txt \
totals: license_ok: 1 files license_replaced: 2 files license_added: 1 files elapsed_time: 0ms"
run_test "$flags" "$test_case" "$expected_output"

delete_sample_project

if (($errors > 0)); then
	exit 1
fi

exit 0
