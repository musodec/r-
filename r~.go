// r~ is a simple utility to remove regular files ending in ~ from the current directory.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	ExitOK = iota
	ExitUsage
	ExitFlagErr // To match flag package usage.
	ExitIOErr
	versionStr = "v0.0.0"
	verboseDoc = "Give more verbose output (i.e. set verbosity=2)"
	verbageDoc = "Control verbosity level: 0=error-only, 1=normal, 2=verbose, 3=debug"
	helpDoc    = "Print brief usage info and exit"
	versionDoc = "Print version and exit"
	recurseDoc = "Recurse through all subdirectories"
	contDoc    = "Continue in case of IO error trying to remove a file"
	emptyDoc   = ""
	suffix     = "~"
	globSuffix = "*" + suffix
)

// Flag vars
var (
	verbose, help, version, recursive, continueOnErr bool
	verbosity                                        int
)

func init() {
	flag.BoolVar(&verbose, "v", false, verboseDoc)
	flag.BoolVar(&verbose, "verbose", false, emptyDoc)
	flag.BoolVar(&help, "h", false, helpDoc)
	flag.BoolVar(&help, "help", false, emptyDoc)
	flag.BoolVar(&version, "V", false, versionDoc)
	flag.BoolVar(&version, "Version", false, emptyDoc)
	flag.BoolVar(&recursive, "r", false, recurseDoc)
	flag.BoolVar(&recursive, "recursive", false, emptyDoc)
	flag.BoolVar(&continueOnErr, "c", false, contDoc)
	flag.BoolVar(&continueOnErr, "continue-on-error", false, emptyDoc)
	flag.IntVar(&verbosity, "y", 1, emptyDoc)
	flag.IntVar(&verbosity, "verbosity", 1, verbageDoc)
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	if verbosity >= 3 {
		log.Printf("NFlags = %d; NArgs = %d", flag.NFlag(), flag.NArg())
	}
	var exitCode = ExitOK

	// Sanity checks on command line:
	if flag.NArg() > 1 {
		log.Printf("%s takes at most one argument.\nIf options are supplied, they must precede the argument.\n",
			os.Args[0])
		flag.Usage()
		os.Exit(ExitUsage)
	}
	if flag.NFlag() > 1 && (help || version) {
		log.Println("Options --help and --version should be supplied alone.")
		exitCode = ExitUsage
	}
	if verbose && verbositySet() && verbosity != 2 {
		log.Println("If both --verbose and --verbosity=n are supplied, then n must equal 2.")
		flag.Usage()
		os.Exit(ExitUsage)
	}
	if verbosity < 0 || verbosity > 3 {
		log.Println("Option --verbosity=n must be in range 0 <= n <= 3.")
		flag.Usage()
		os.Exit(ExitUsage)
	}
	if verbosity >= 3 {
		log.Printf("Verbosity level = %d\n", verbosity)
	}

	// Process the help|version special cases:
	switch {
	case help:
		flag.Usage()
		os.Exit(exitCode)
	case version:
		fmt.Println(versionStr)
		os.Exit(exitCode)
	}

	var dir string
	var err error

	// Check the directory argument:
	switch flag.NArg() {
	case 0:
		dir, err = os.Getwd()
		if err != nil {
			log.Printf("Getwd: %v\n", err.Error())
			os.Exit(ExitIOErr)
		}
		if verbosity >= 2 {
			fmt.Printf("Working Directory: %s\n", dir)
		}
	case 1:
		p := flag.Arg(0)
		if verbosity >= 3 {
			log.Printf("Locating directory %q...\n", p)
		}
		if filepath.IsAbs(p) {
			dir = filepath.Clean(p)
		} else {
			b, err := os.Getwd()
			if err != nil {
				log.Printf("Getwd: %v\n", err.Error())
				os.Exit(ExitIOErr)
			}
			dir = filepath.Join(b, p)
		}
		fi, err := os.Stat(dir)
		if err != nil {
			log.Print(err)
			os.Exit(ExitIOErr)
		}
		if !fi.IsDir() {
			log.Printf("%q not a directory", dir)
			os.Exit(ExitIOErr)
		}
		if verbosity >= 2 {
			fmt.Printf("Working relative to directory %s\n", dir)
		}
	}

	if recursive {
		log.Println("RECURSIVE WIP")
		filepath.Walk(dir, walkFunc)
		os.Exit(ExitFlagErr)
	}
	exitCode = rDir(dir)
	os.Exit(exitCode)
}

func rDir(dir string) int {
	exitCode := ExitOK
	g, err := filepath.Glob(filepath.Join(dir, globSuffix))
	if err == filepath.ErrBadPattern {
		log.Printf("Glob ErrBadPattern: %v" + err.Error())
		os.Exit(ExitIOErr)
	}
	for _, p := range g {
		f, err := os.Stat(p)
		if err != nil {
			log.Printf("Stat: %v\n", err.Error())
			if continueOnErr {
				exitCode = ExitIOErr
			} else {
				os.Exit(ExitIOErr)
			}
		}
		if !f.Mode().IsRegular() {
			if verbosity >= 2 {
				log.Printf("%s is not a regular file - skipping\n",
					filepath.Base(p))
			}
			continue
		}
		err = os.Remove(p)
		if err != nil {
			log.Printf("Failed to remove %q: %s\n", p, err.Error())
			if continueOnErr {
				exitCode = ExitIOErr
			} else {
				os.Exit(ExitIOErr)
			}
		} else if verbosity >= 1 {
			fmt.Printf("Removed %s\n", filepath.Base(p))
		}
	}
	return exitCode
}

var abort bool

func walkFunc(path string, info os.FileInfo, err error) error {
	exitCode := ExitOK
	if abort {
		return filepath.SkipDir
	}
	if err != nil {
		log.Printf("Walk: %v\n", err)
		if !continueOnErr {
			abort = true
			return filepath.SkipDir
		}
	}
	if info.Mode().IsRegular() && strings.HasSuffix(path, suffix) {
		// >>> DUP TODO: Combine
		err = os.Remove(path)
		if err != nil {
			log.Printf("Failed to remove %q: %s\n", path, err.Error())
			if continueOnErr {
				exitCode = ExitIOErr
			} else {
				os.Exit(ExitIOErr)
			}
		} else if verbosity >= 1 {
			fmt.Printf("Removed %s\n", path)
		}
		// <<<
	} else if info.IsDir() && verbosity >= 2 {
		fmt.Printf("Entering %s\n", path)
	}
	_ = exitCode // return/globalize exitCode
	return nil
}

func verbositySet() bool {
	set := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "verbosity" {
			set = true
		}
	})
	return set
}
