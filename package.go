package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Package struct {
	Name string
	Dir  string

	Commits Commits
}

type Packages []*Package

var emojiRune = '✅'

func NewPackage(name, gopath string) *Package {
	dir := filepath.Join(gopath, "src", name)
	return &Package{
		Name: name,
		Dir:  dir,
	}
}

func (p *Package) Refresh() error {
	git := NewGit(p.Dir)
	if err := git.Update(); err != nil {
		return err
	}

	p.Commits = git.Commits()
	return nil
}

func (p *Package) IsOutdated() bool {
	return len(p.Commits) > 0
}

func (p *Package) String() string {
	count := len(p.Commits)
	out := fmt.Sprintf("%s [%c %d]\n", green(p.Name), emojiRune, count)
	out = fmt.Sprintf("%s%s", out, p.Commits)
	return out
}

func (p *Package) Update() error {
	fmt.Println("Updating", p.Name, "...")
	fmt.Sprintf("go get -u %s", p.Name)
	//_, err := Run(cmd, ".")
	return nil
}

func (pkgs Packages) Outdated() Packages {
	var outdated Packages
	for _, pkg := range pkgs {
		if pkg.IsOutdated() {
			outdated = append(outdated, pkg)
		}
	}
	return outdated
}

func init() {
	// No legacy is so rich as honesty. (:
	if os.Getenv("TRUTH_MODE") != "" {
		emojiRune = '🐞'
	}
}
