.PHONY : all clean install test

all : r~

clean :
	rm -vf *~
	rm -vfR testdata/
	go clean

test :
	./mktest.sh

#
# Concrete Target
#

r~ : *.go
	go build

install :
	go fmt
	go install
