r~
==

r~ is a simple utility to remove regular files ending in __~__ from
directories.

Such files are automatically generated by `emacs` or by the backup
options to everyday Gnu utilities such as `cp`, `mv`, `ln`, and
`install`.  This provides a convenient safety net should you ever need
to backtrack or compare some local changes.  But then the time comes
when you want to clean a directory of any such backup files that may
have crept in.

Of course, you can always try something like `rm *~` or `find . -type
f -name '*~' -exec rm -v '{}' \;` from a shell, or `~ x yes` in dired.
So why `r~`?

* `r~` is shorter and easier to type (and you may end up typing it a
  lot).  I had a colleague who once deleted the entire filesystem from
  a live Unix server by `rm -fR * ~`.
* `r~` distinguishes between regular files (which you probably do want
  to remove) and special files (which you probably don't).
* `r~` provides options to control common use cases.

Usage
-----

`r~ --help` gives:

```
Usage of r~:
  r~ [options...] [--] [<dir>]
  r~ --help
  r~ --Version

<dir> defaults to the current working directory, and relative paths
will be resolved with respect to the current working directory.

Short options may not be combined.  Long options take one or two hyphens:
  -V	Print version and exit
  -Version

  -h	Print brief usage info and exit
  -help

  -i	Prompt for each deletion
  -interactive

  -k	Keep going in case of IO error (or flag parsing errors if supplied before the error occurs)
  -keep-going

  -r	Recurse through all subdirectories
  -recursive

  -v	Give more verbose output (i.e. set verbosity=2)
  -verbose

  -verbosity int
	 (default 1)
  -y int
	Control verbosity level: 0=error-only, 1=normal, 2=verbose, 3=debug (default 1)

Exit value:
  0 OK
  1 Usage error
  2 Flag parsing error
  3 IO error
```

Requirements
------------

Golang installation to build.  (Tested on go version 1.15 so far.)

Bash for testing.

Make optional.

Installation
------------

`make` or `go build` to build in the r~ working directory.

`make install` or `go install`, from the r~ directory, to install as
per GOBIN environment variable.

Other `make` targets supported: `test`, `testdir`, `clean`.

Fun Facts
---------

`make; ./r~` builds itself... only to remove itself.  `r~` was used
many times during development of itself, to keep itself clean.

r~ can be neat as a step in `make clean` type scripts, or prior to `git
add` and the like.  `r~` also serves neatly as a git hook.
