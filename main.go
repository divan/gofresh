package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	update = flag.Bool("update", false, "Update all packages")
	dryRun = flag.Bool("dry-run", false, "Dry run")
	GOPATH = os.Getenv("GOPATH")
)

func main() {
	flag.Parse()

	imports, err := Imports(".")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d imports, checking for updates...\n", len(imports))

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

	for _, name := range imports {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			pkg := NewPackage(name, GOPATH)
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
			fmt.Printf("Updating %s...\n", green(pkg.Name))
			if !*dryRun {
				err = pkg.Update()
				if err != nil {
					fmt.Printf("%s: %s\n", red(pkg.Name), redBold(err.Error()))
					return
				}
			}
		}

		// TODO: check again?
		outdated = Packages{}
	}

	hasUpdate := len(outdated) > 0
	if !hasUpdate {
		fmt.Println("Everything is up to date.")
		return
	} else {
		for _, pkg := range outdated {
			fmt.Println(pkg)
		}
		fmt.Printf(green("---\nYou have %d packages out of date\n", len(outdated)))
		fmt.Println("To update all packages automatically, run", bold("gofresh -update"))
	}
}
