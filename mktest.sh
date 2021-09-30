#!/usr/bin/bash -e

TESTDIR=testdata

rm -fR $TESTDIR
mkdir $TESTDIR
cd $TESTDIR
touch reg~ .hid~ ok -- --flag-like~
mkfifo ff ff1~
mkdir --verbose dir dir1~
