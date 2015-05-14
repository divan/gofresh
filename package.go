package main

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/vcs"
)

// Package represents single Go package/repo in Gopath.
type Package struct {
	Name string
	Dir  string
	Repo *vcs.RepoRoot

	Commits Commits
}

// Packages is an obvious type, but I prefer to have golint happy.
type Packages []*Package

var emojiRune = '✅'

// NewPackage returns new package.
func NewPackage(name, gopath string) *Package {
	dir := filepath.Join(gopath, "src", name)
	repo, err := vcs.RepoRootForImportPath(name, false)
	if err != nil {
		// it's ok, silently discard errors here
		return nil
	}
	return &Package{
		Name: name,
		Dir:  dir,

		Repo: repo,
	}
}

// Refresh updates package info about new commits.
//
// It typically require internet connection to check
// remote side.
func (p *Package) Refresh() error {
	var vcs VCS
	switch p.Repo.VCS.Name {
	case "Git":
		vcs = NewGit(p.Dir)
	case "Mercurial":
		vcs = NewHg(p.Dir)
	default:
		return fmt.Errorf("unknown VCS")
	}

	if err := vcs.Update(); err != nil {
		return err
	}

	p.Commits = vcs.Commits()
	return nil
}

// IsOutdated returns true if package has updates on remote.
func (p *Package) IsOutdated() bool {
	return len(p.Commits) > 0
}

// String implements Stringer for Package.
func (p *Package) String() string {
	count := len(p.Commits)
	out := fmt.Sprintf("%s [%c %d]\n", green(p.Name), emojiRune, count)
	out = fmt.Sprintf("%s%s", out, p.Commits)
	return out
}

// UpdateCmd returns command used to update package.
func (p Package) UpdateCmd() string {
	return fmt.Sprintf("go get -u %s", p.Name)
}

// Update updates package to the latest revision.
func (p *Package) Update() error {
	cmd := p.UpdateCmd()
	_, err := Run(cmd, p.Dir)
	return err
}

// Outdated filters only outdated packages.
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
