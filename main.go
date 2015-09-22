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
	force  = flag.Bool("f", false, "Use force while updating packages")
	dryRun = flag.Bool("dry-run", false, "Dry run")
	expand = flag.Bool("expand", false, "Expand list of commits")
)

func main() {
	flag.Usage = Usage
	flag.Parse()

	var packages []string

	// In case package name(s) were specified, check only them
	byName := len(flag.Args()) != 0
	if byName {
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
		wg     sync.WaitGroup
		pkgs   Packages
		ch     = make(chan *Package)
		failed bool
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
			pkg, err := NewPackage(name, gopath)
			if err != nil {
				// There always will be error, when processing imports from
				// source, like 'fmt', 'net/http', etc.
				// But for explicitly specified packages by name, we should
				// show user an error.
				if byName {
					failed = true
					fmt.Printf("%s: %s\n", red(name), redBold(err.Error()))
				}
				return
			}
			err = pkg.Refresh()
			if err != nil {
				failed = true
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
			cmdline := strings.Join(pkg.UpdateCmd(*force), " ")
			fmt.Println(green(cmdline))
			if !*dryRun {
				err := pkg.Update(*force)
				if err != nil {
					fmt.Printf("%s: %s\n", red(pkg.Name), redBold(err.Error()))
					failed = true
					continue
				}
			}
		}

		// TODO: check again?
		outdated = Packages{}
	}

	upToDate := len(outdated) == 0
	if upToDate && !failed {
		fmt.Println("Everything is up to date.")
		return
	} else if upToDate && failed {
		fmt.Println("There were some errors, check incomplete or wrong usage.")
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
