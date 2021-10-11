#!/usr/bin/bash

TESTDIR=testdir
FAIL=0

populate_test_dir () {
    touch reg~ .hid~ spacy\ ~ ok -- --flag-like~
    mkfifo ff ff1~
    mkdir dir dir1~
}

mktestdir () {
    rm -fR $TESTDIR
    mkdir $TESTDIR
    cd $TESTDIR
    populate_test_dir
    cd ..
}

mktestdirs () {
    mktestdir
    cd $TESTDIR
    mkdir --parents x/y/z
    cd x
    populate_test_dir
    cd y/z
    populate_test_dir
    cd ../../../..
}

# Check for supplied argument
if [[ $# -gt 1 ]]; then
    echo 'Too many parameters supplied' >&2
    exit 1
elif [[ $# -eq 1 ]]; then
    case $1 in
	'mktestdir' )
	    mktestdir
	    echo 'Made testdir'
	    exit 0
	    ;;
	'mktestdirs' )
	    mktestdirs
	    echo 'Made testdirs'
	    exit 0
	    ;;
	* )
	    echo 'Unrecognized parameter supplied' >&2
	    exit 1
    esac
fi

# Non-recursive case
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

# Recursive case
for (( N = 0; N <= 3; N++ )); do
    mktestdirs
    cd $TESTDIR
    diff <(ls -AR --color=no) ../testdata/lsAR0
    if [[ $? -ne 0 ]]; then
	echo mktestdirs FAILURE >&2
	exit 1
    fi
    r~ --recursive --verbosity=$N >/tmp/rv$N 2>&1
    diff /tmp/rv$N ../testdata/rv$N
    if [[ $? -ne 0 ]]; then
	FAIL=1
	echo recursive log FAILURE at verbosity=$N
    fi
    diff <(ls -AR --color=no) ../testdata/lsAR1
    if [[ $? -ne 0 ]]; then
	FAIL=1
	echo recursive log FAILURE at verbosity=$N
    fi
    cd ..
done

rm /tmp/rv[0-3]
rm -fR $TESTDIR

# Report
if [[ $FAIL -eq 0 ]]; then
    echo All tests PASS
else
    echo Test FAILURE >&2
    exit 1
fi
