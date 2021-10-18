// r~ is a simple utility to remove regular files ending in ~ from the current directory.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ExitCode int

const (
	ExitOK ExitCode = iota
	ExitUsage
	ExitFlagErr // To match flag package usage.
	ExitIOErr
)

var exitCode ExitCode

func setExitCode(c ExitCode) { exitCode = c }

func getExitCode() ExitCode { return exitCode }

func exitWithCode(c ...ExitCode) {
	switch len(c) {
	case 0:
		// Do nothing.
	case 1:
		exitCode = c[0]
	default:
		panic("Internal error: exitWithCode case fallthrough")
	}
	os.Exit(int(getExitCode()))
}

const (
	versionStr = "v0.0.3"
	verboseDoc = "Give more verbose output (i.e. set verbosity=2)"
	verbageDoc = "Control verbosity level: 0=error-only, 1=normal, 2=verbose, 3=debug"
	helpDoc    = "Print brief usage info and exit"
	versionDoc = "Print version and exit"
	recurseDoc = "Recurse through all subdirectories"
	contDoc    = "Continue in case of IO error (or flag parsing errors if supplied earlier)"
	interDoc   = "Prompt for each deletion"
	emptyDoc   = ""
	suffix     = "~"
	globSuffix = "*" + suffix
	contErrPat = `^--?c(ontinue-on-error)?$`
	yesOrNoPat = `(?i:^[[:blank:]]*y(?:es)?\W)|(?:^\s*$)`
	usagePre   = `  r~ [options...] [--] [<dir>]
  r~ --help
  r~ --Version

<dir> defaults to the current working directory, and relative paths
will be resolved with respect to the current working directory.

Options take one or two hyphens:
`
	usagePost = `Exit value:
  0 OK
  1 Usage error
  2 Flag parsing error
  3 IO error
`
)

// Flag vars
var (
	verbose, help, version, recursive, continueOnErr, interactive bool
	verbosity                                                     int
	contErrRe, yesOrNoRe                                          *regexp.Regexp
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
	flag.BoolVar(&interactive, "i", false, interDoc)
	flag.BoolVar(&interactive, "interactive", false, emptyDoc)
	flag.IntVar(&verbosity, "y", 1, verbageDoc)
	flag.IntVar(&verbosity, "verbosity", 1, emptyDoc)
	// setExitCode(ExitOK) // implicitly
	contErrRe = regexp.MustCompile(contErrPat)
	yesOrNoRe = regexp.MustCompile(yesOrNoPat)
}

func main() {
	var err error
	if contOnErrArg() {
		flag.CommandLine.Init(
			filepath.Base(os.Args[0])+" --continue-on-error",
			flag.ContinueOnError)
		err = flag.CommandLine.Parse(os.Args[1:])
		if err != nil {
			setExitCode(ExitFlagErr)
		}
	} else {
		flag.Parse()
	}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage of %s:\n%s", os.Args[0], usagePre)
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(),
			"\n%s", usagePost)
	}

	log.SetFlags(0)
	if verbosity >= 3 {
		log.Printf("NFlags = %d; NArgs = %d", flag.NFlag(), flag.NArg())
		log.Printf("Version = %s\nverbosity = %d\nrecursive = %t\ncontinue-on-error = %t\n",
			versionStr, verbosity, recursive, continueOnErr)
	}

	// Sanity checks on command line:
	if flag.NArg() > 1 {
		log.Printf("%s takes at most one argument.\nIf options are supplied, they must precede the argument.\n",
			os.Args[0])
		flag.Usage()
		exitWithCode(ExitUsage)
	}
	if flag.NFlag() > 1 && (help || version) {
		log.Println("Options --help and --version should be supplied alone.")
		setExitCode(ExitUsage)
	}
	if verbose && verbositySet() && verbosity != 2 {
		log.Println("If both --verbose and --verbosity=n are supplied, then n must equal 2.")
		flag.Usage()
		exitWithCode(ExitUsage)
	} else if verbose {
		verbosity = 2
	}
	if verbosity < 0 || verbosity > 3 {
		log.Println("Option --verbosity=n must be in range 0 <= n <= 3.")
		flag.Usage()
		exitWithCode(ExitUsage)
	}

	// Process the help|version special cases:
	switch {
	case help:
		if getExitCode() == ExitOK {
			flag.CommandLine.SetOutput(os.Stdout)
		}
		flag.Usage()
		exitWithCode()
	case version:
		fmt.Printf("%s %s\n", filepath.Base(os.Args[0]), versionStr)
		exitWithCode()
	}

	var dir string

	// Check the directory argument:
	switch flag.NArg() {
	case 0:
		dir, err = os.Getwd()
		if err != nil {
			log.Printf("Getwd: %v\n", err.Error())
			exitWithCode(ExitIOErr)
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
				exitWithCode(ExitIOErr)
			}
			dir = filepath.Join(b, p)
		}
		fi, err := os.Stat(dir)
		if err != nil {
			log.Print(err)
			exitWithCode(ExitIOErr)
		}
		if !fi.IsDir() {
			log.Printf("%q not a directory", dir)
			exitWithCode(ExitIOErr)
		}
		if verbosity >= 2 {
			fmt.Printf("Working relative to directory %s\n", dir)
		}
	}

	if recursive {
		err = filepath.Walk(dir, walkFunc)
		if err != nil {
			log.Printf("Walk: %v\n", err)
			if c := getExitCode(); c == ExitOK {
				setExitCode(ExitIOErr)
			}
		}
	} else {
		rDir(dir)
	}
	exitWithCode()
}

func rDir(dir string) {
	g, err := filepath.Glob(filepath.Join(dir, globSuffix))
	if err == filepath.ErrBadPattern {
		log.Printf("Glob ErrBadPattern: %v" + err.Error())
		exitWithCode(ExitIOErr)
	}
	for _, p := range g {
		f, err := os.Stat(p)
		if err != nil {
			log.Printf("Stat: %v\n", err.Error())
			if continueOnErr {
				setExitCode(ExitIOErr)
			} else {
				exitWithCode(ExitIOErr)
			}
		}
		if !f.Mode().IsRegular() {
			if verbosity >= 2 {
				log.Printf("%s is not a regular file - skipping\n",
					filepath.Base(p))
			}
			continue
		}
		rm(p, true)
	}
}

func rm(path string, base bool) {
	maybeBase := path
	if base {
		maybeBase = filepath.Base(path)
	}
	if interactive {
		fmt.Printf("Remove %s? [Y/n] ", maybeBase)
		r := bufio.NewReader(os.Stdin)
		s, err := r.ReadString('\n')
		if err != nil {
			log.Printf("ReadAll: %v\n", err)
			if continueOnErr {
				setExitCode(ExitIOErr)
			} else {
				exitWithCode(ExitIOErr)
			}
		}
		if verbosity >= 3 {
			log.Printf("Read: %q\n", s)
		}
		if fi, _ := os.Stdin.Stat(); fi.Mode()&os.ModeCharDevice == 0 {
			fmt.Println()
		}
		if !yesOrNoRe.MatchString(s) {
			return
		}
	}
	err := os.Remove(path)
	if err != nil {
		log.Printf("Failed to remove %q: %s\n", path, err.Error())
		if continueOnErr {
			setExitCode(ExitIOErr)
		} else {
			exitWithCode(ExitIOErr)
		}
	} else if verbosity >= 1 {
		fmt.Printf("Removed %s\n", maybeBase)
	}
}

var abort bool

func walkFunc(path string, info os.FileInfo, err error) error {
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
		rm(path, false)
	} else if verbosity >= 2 && info.IsDir() {
		fmt.Printf("Entering %s\n", path)
	} else if verbosity >= 2 && !info.Mode().IsRegular() {
		log.Printf("%s is not a regular file - skipping\n",
			filepath.Base(path))
	}
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

func contOnErrArg() bool {
	if len(os.Args) <= 1 {
		return false
	}
	for _, a := range os.Args[1:] {
		if contErrRe.MatchString(a) {
			return true
		}
	}
	return false
}
