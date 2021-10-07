.PHONY : all clean install test

all : r~

clean :
	rm -vf *~
	rm -vfR testdir/
	go clean

test :
	./dotest.sh

#
# Concrete Target
#

r~ : *.go
	go build

install :
	go fmt
	go install
