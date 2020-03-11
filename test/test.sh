#!/bin/bash

CMD="$1"
RED='\033[0;31m'
GREEN='\033[32m'
NC='\033[0m'
errors=0

remove_ansi_color() {
	if [[ "$(uname)" == "Darwin" ]]; then
		output=$(echo $1 | sed $'s,\x1b\\[[0-9;]*[a-zA-Z],,g')
		eval "$2=\"${output}\""
	else
		output=$(echo $1 | sed -r "s/[[:cntrl:]]\[[0-9]{1,3}m//g")
		eval "$2=\"${output}\""
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
	PROJECT_DIR=sample-project
	ARGS="${PROJECT_DIR}/licenses/current-license.txt ${PROJECT_DIR} java js cpp go"
	delete_sample_project
	extract_sample_project
	echo -e "\n$3"
	remove_ansi_color "$($1 $2 $ARGS)" res
	eval "$4=\"${res}\""
}

run_test $CMD '-version' 'Testing version...' result

if [[ $result =~ $'version:' ]]; then
	on_success
else
	on_failure
fi

run_test $CMD '' 'Testing read only...' result

if [[ $result =~ $'1 licenses ok, 0 licenses replaced, 0 licenses added' && \
	$result =~ $'2 files had no license' && \
	$result =~ $'2 files had a different license' ]]; then
	on_success
else
	on_failure
fi

run_test $CMD '-a' 'Testing with -a flag...' result

if [[ $result =~ $'1 licenses ok, 0 licenses replaced, 2 licenses added' && \
	$result =~ $'2 files had a different license' ]]; then
	on_success
else
	on_failure
fi

run_test $CMD '-a -r' 'Testing with -a and -r flags...' result

if [[ $result =~ $'1 licenses ok, 2 licenses replaced, 2 licenses added' ]]; then
	on_success
else
	on_failure
fi

run_test $CMD '-a -r -i src/other' 'Testing with -a and -r and -i flags...' result

if [[ $result =~ $'1 licenses ok, 2 licenses replaced, 1 licenses added' ]]; then
	on_success
else
	on_failure
fi

run_test $CMD '-a -r -v -i src/other' 'Testing with -a and -r and -i and -v flags...' result

if [[ $result =~ $'license_ok: 1 files' && \
	$result =~ $'license_replaced: 2 files' && \
	$result =~ $'license_added: 1 files' && \
	$result =~ $'elapsed_time:' ]]; then
	on_success
else
	on_failure
fi

delete_sample_project

if (($errors > 0)); then
	exit 1
fi
exit 0
