// r~ is a simple utility to remove regular files ending in ~ from the current directory.
package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	log.SetFlags(0)
	d, err := os.Getwd()
	if err != nil {
		log.Fatal("Getwd: " + err.Error())
	}
	log.Printf("Working Directory: %s\n", d)
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
			log.Printf("%s is not a regular file - skipping\n", p)
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
