package main

import ()

// VCS represents Version Control System.
type VCS interface {
	Update() error
	Commits() []string
}

// Git implements VCS interface for Git.
type Git struct {
	Dir string
}

// NewGit creates new Git object.
func NewGit(dir string) *Git {
	return &Git{
		Dir: dir,
	}
}

// Update updates info from the remote.
func (git *Git) Update() error {
	_, err := Run(git.Dir, "git", "fetch", "origin")
	return err
}

// Commits returns new commits in master branch.
func (git *Git) Commits() []string {
	out, err := Run(git.Dir, "git", "log", "HEAD..origin/master", "--oneline")
	if err != nil {
		return nil
	}
	return out
}

// Hg implements VCS interface for Mercurial.
type Hg struct {
	Dir string
}

// NewHg creates new Hg object.
func NewHg(dir string) *Hg {
	return &Hg{
		Dir: dir,
	}
}

// Update updates info from the remote.
func (hg *Hg) Update() error {
	return nil
}

// Commits returns new commits in master branch.
func (hg *Hg) Commits() []string {
	out, err := Run(hg.Dir, "hg", "incoming", "-n", "-q", "--template", "{node|short} {desc|strip|firstline}\n")
	if err != nil {
		return nil
	}
	return out
}
