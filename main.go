package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	update = flag.Bool("update", false, "Update all packages")
	dryRun = flag.Bool("dry-run", false, "Dry run")
	expand = flag.Bool("expand", false, "Expand list of commits")
)

func main() {
	flag.Usage = Usage
	flag.Parse()

	var packages []string
	// In case package name(s) were specified, check only them
	if len(flag.Args()) != 0 {
		packages = flag.Args()
		fmt.Printf("Checking %d packages for updates...\n", len(packages))
	} else {
		// otherwise, find imports for current package and
		// subpackages
		var err error
		packages, err = Imports(".")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Found %d imports, checking for updates...\n", len(packages))
	}

	var (
		wg   sync.WaitGroup
		pkgs Packages
		ch   = make(chan *Package)
	)

	go func() {
		for pkg := range ch {
			pkgs = append(pkgs, pkg)
		}
	}()

	gopath := GOPATH()
	for _, name := range packages {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			pkg := NewPackage(name, gopath)
			if pkg == nil {
				return
			}
			err := pkg.Refresh()
			if err != nil {
				fmt.Printf("%s: %s\n", red(name), redBold(err.Error()))
				return
			}

			ch <- pkg

		}(name)
	}
	wg.Wait()
	close(ch)

	outdated := pkgs.Outdated()

	// Update, if requested
	if *update {
		for _, pkg := range outdated {
			cmdline := strings.Join(pkg.UpdateCmd(), " ")
			fmt.Println(green(cmdline))
			if !*dryRun {
				err := pkg.Update()
				if err != nil {
					fmt.Printf("%s: %s\n", red(pkg.Name), redBold(err.Error()))
					return
				}
			}
		}

		// TODO: check again?
		outdated = Packages{}
	}

	upToDate := len(outdated) == 0
	if upToDate {
		fmt.Println("Everything is up to date.")
		return
	}

	for _, pkg := range outdated {
		fmt.Print(pkg)
	}
	fmt.Printf(green("---\nYou have %d packages out of date\n", len(outdated)))
	fmt.Println("To update all packages automatically, run", bold("gofresh -update"))
}

func Usage() {
	fmt.Fprintf(os.Stderr, "gofresh [-options]\n")
	fmt.Fprintf(os.Stderr, "gofresh [-options] [package(s)]\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

// GOPATH returns GOPATH to be used for package update checks.
//
// In case there are many dirs in GOPATH, only the first will be used.
// TODO: add multiple dirs support? someone use it w/o vendoring tools?
func GOPATH() string {
	path := os.Getenv("GOPATH")
	colonDelim := func(r rune) bool {
		return r == ':'
	}
	fields := strings.FieldsFunc(path, colonDelim)
	return fields[0]
}
