package main

import (
	"go/parser"
	"go/token"
	"sort"
	"strings"
)

// Imports returns list of packages imported by
// all sources found in dir.
func Imports(dir string) ([]string, error) {
	fset := token.NewFileSet()

	// Find all packages in current dir
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return nil, err
	}

	// Iterate over each package, each file
	// and add imported packages to map
	imports := make(map[string]struct{})
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, impt := range file.Imports {
				path := strings.Trim(impt.Path.Value, `"`)
				imports[path] = struct{}{}
			}
		}
	}

	// Convert map to slice and sort
	var ret []string
	for name, _ := range imports {
		ret = append(ret, name)
	}
	sort.Strings(ret)
	return ret, nil
}
