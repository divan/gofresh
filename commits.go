package main

import "fmt"

// MaxCommits limits number of commits to show.
const MaxCommits = 3

type Commits []string

func (commits Commits) String() string {
	if len(commits) == 0 {
		return ""
	}

	count := len(commits)

	Max := MaxCommits
	if *expand {
		Max = count
	}
	isLimited := count > Max

	limit := Max
	if !isLimited {
		limit = count
	}

	var out string
	for _, commit := range commits[:limit] {
		str := cyan(fmt.Sprintf("    %s\n", commit))
		out = fmt.Sprintf("%s%s", out, str)
	}

	if isLimited {
		more := yellow(fmt.Sprintf("and %d more...\n", count-Max))
		out = fmt.Sprintf("%s%s", out, more)
	}

	return out
}
