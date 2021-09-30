// r~ is a simple utility to remove regular files ending in ~ from the current directory.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
)

// Flag vars
var (
	verbose, help, version, recursive, continueOnErr bool
	verbosity                                        int
)

func init() {
	flag.BoolVar(&verbose, "v", false, verboseDoc)
	flag.BoolVar(&verbose, "verbose", false, verboseDoc)
	flag.BoolVar(&help, "h", false, helpDoc)
	flag.BoolVar(&help, "help", false, helpDoc)
	flag.BoolVar(&version, "V", false, versionDoc)
	flag.BoolVar(&version, "Version", false, versionDoc)
	flag.BoolVar(&recursive, "r", false, recurseDoc)
	flag.BoolVar(&recursive, "recursive", false, recurseDoc)
	flag.BoolVar(&continueOnErr, "c", false, contDoc)
	flag.BoolVar(&continueOnErr, "continue-on-error", false, contDoc)
	flag.IntVar(&verbosity, "y", 1, verbageDoc)
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
		log.Printf("%s takes at most one argument.\nIf flags are supplied, they must precede the argument.\n",
			os.Args[0])
		flag.Usage()
		os.Exit(ExitUsage)
	}
	if flag.NFlag() > 1 && (help || version) {
		log.Println("Flags --help and --version should be supplied alone.")
		exitCode = ExitUsage
	}
	if verbose && verbositySet() && verbosity != 2 {
		log.Println("If both --verbose and --verbosity=n are supplied, then n must equal 2.")
		flag.Usage()
		os.Exit(ExitUsage)
	}
	if verbosity < 0 || verbosity > 3 {
		log.Println("Flag --verbosity=n must be in range 0 <= n <= 3.")
		flag.Usage()
		os.Exit(ExitUsage)
	}
	if recursive {
		log.Println("Not yet implemented")
		os.Exit(ExitFlagErr)
	}

	if verbosity >= 3 {
		log.Printf("Verbosity level = %d\n", verbosity)
	}

	switch {
	case help:
		flag.Usage()
		os.Exit(ExitUsage)
	case version:
		fmt.Println(versionStr)
		os.Exit(ExitUsage)
	}

	switch flag.NArg() {
	case 0:
		d, err := os.Getwd()
		if err != nil {
			log.Printf("Getwd: %v\n", err.Error())
			os.Exit(ExitIOErr)
		}
		if verbosity >= 2 {
			fmt.Printf("Working Directory: %s\n", d)
		}
	case 1:
		d := flag.Arg(1)
		if filepath.IsAbs(d) {
			d = filepath.Clean(d)
			fi, err := os.Stat(d)
			if err != nil {
				log.Print(err)
				os.Exit(ExitIOErr)
			}
			if !fi.IsDir() {
				log.Printf("%q not a directory", d)
				os.Exit(ExitIOErr)
			}
		}
		// default:
		// os.Exit(ExitUsage) // This should never happen
	}
	g, err := filepath.Glob("*~")
	if err == filepath.ErrBadPattern {
		panic("ErrBadPattern: " + err.Error()) // This should never happen.
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
				log.Printf("%s is not a regular file - skipping\n", p)
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
			fmt.Printf("Removed %s\n", p)
		}
	}
	os.Exit(exitCode)
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
