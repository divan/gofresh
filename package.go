package main

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/vcs"
)

type Package struct {
	Name string
	Dir  string
	Repo *vcs.RepoRoot

	Commits Commits
}

type Packages []*Package

var emojiRune = '✅'

func NewPackage(name, gopath string) *Package {
	dir := filepath.Join(gopath, "src", name)
	repo, err := vcs.RepoRootForImportPath(name, false)
	if err != nil {
		return nil
	}
	return &Package{
		Name: name,
		Dir:  dir,

		Repo: repo,
	}
}

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
	cmd := fmt.Sprintf("go get -u %s", p.Name)
	_, err := Run(cmd, p.Dir)
	return err
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
