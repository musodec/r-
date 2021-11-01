.PHONY : all clean install test testdir

all : r~

clean :
	rm -vf *~ testdata/*~
	rm -vfR testdir/
	go clean

test :
	go test
	./dotest.sh

testdir :
	./dotest.sh mktestdir

#
# Concrete Target
#

r~ : *.go
	go build

install :
	go fmt
	go install
