// r~ is a simple utility to remove regular files ending in ~ from the current directory.
package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

var verbose bool

func init() {
	flag.BoolVar(&verbose, "v", false, "Give more verbose output")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	log.Printf("NArgs = %d", flag.NArg())

	switch flag.NArg() {
	case 0:
		d, err := os.Getwd()
		if err != nil {
			log.Fatal("Getwd: " + err.Error())
		}
		if verbose {
			log.Printf("Working Directory: %s\n", d)
		}
	case 1:
		d := flag.Arg(1)
		if filepath.IsAbs(d) {
			d = filepath.Clean(d)
			fi, err := os.Stat(d)
			if err != nil {
				log.Fatal(err)
			}
			if !fi.IsDir() {
				log.Fatalf("%q not a directory", d)
				log.Printf("Absolute dir %q", d)
			}
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
	g, err := filepath.Glob("*~")
	if err == filepath.ErrBadPattern {
		log.Fatalf("ErrBadPattern: %s\n", err.Error())
	}
	for _, p := range g {
		f, err := os.Stat(p)
		if err != nil {
			log.Fatal("Stat: " + err.Error())
		}
		if !f.Mode().IsRegular() {
			if verbose {
				log.Printf("%s is not a regular file - skipping\n", p)
			}
			continue
		}
		err = os.Remove(p)
		if err != nil {
			log.Printf("Failed to remove %q: %s\n", p, err.Error())
		} else {
			log.Printf("Removed %s\n", p)
		}
	}
}
