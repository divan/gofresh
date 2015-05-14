package main

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
	_, err := Run("git fetch origin", git.Dir)
	return err
}

// Commits returns new commits in master branch.
func (git *Git) Commits() []string {
	out, err := Run("git log HEAD..origin/master --oneline", git.Dir)
	if err != nil {
		return nil
	}
	return out
}
