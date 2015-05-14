package main

import "fmt"

// ShowMax limits number of commits to show.
const ShowMax = 3

type Commits []string

func (commits Commits) String() string {
	if len(commits) == 0 {
		return ""
	}

	count := len(commits)

	isLimited := count > ShowMax

	limit := ShowMax
	if !isLimited {
		limit = count
	}

	var out string
	for _, commit := range commits[:limit] {
		str := cyan(fmt.Sprintf("    %s\n", commit))
		out = fmt.Sprintf("%s%s", out, str)
	}

	if isLimited {
		more := yellow(fmt.Sprintf("and %d more...", count-ShowMax))
		out = fmt.Sprintf("%s%s", out, more)
	}

	return out
}
