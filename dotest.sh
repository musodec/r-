#!/usr/bin/bash

TESTDIR=testdir
FAIL=0

mktestdir () {
    rm -fR $TESTDIR
    mkdir $TESTDIR
    cd $TESTDIR
    touch reg~ .hid~ spacy\ ~ ok -- --flag-like~
    mkfifo ff ff1~
    mkdir dir dir1~
    cd ..
}


for (( N = 0; N <= 3; N++ )); do
    mktestdir
    cd $TESTDIR
    diff <(ls -A --color=no) ../testdata/lsA0
    if [[ $? -ne 0 ]]; then
	echo mktestdir FAILURE >&2
	exit 1
    fi
    r~ --verbosity=$N >/tmp/v$N 2>&1
    diff /tmp/v$N ../testdata/v$N
    if [[ $? -ne 0 ]]; then
	FAIL=1
	echo log FAILURE at verbosity=$N
    fi
    diff <(ls -A --color=no) ../testdata/lsA1
    if [[ $? -ne 0 ]]; then
	FAIL=1
	echo r~ FAILURE at verbosity=$N
    fi
    cd ..
done


rm /tmp/v[0-3]
rm -fR $TESTDIR


if [[ $FAIL -eq 0 ]]; then
    echo All tests PASS
else
    echo Test FAILURE >&2
    exit 1
fi
