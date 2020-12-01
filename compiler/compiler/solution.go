package zoe

import (
	"errors"
	"path/filepath"
)

type Solution struct {
	Files map[int]*File
}

func NewSolution() *Solution {
	return &Solution{
		Files: make(map[int]*File),
	}
}

func (s *Solution) AddFile(uri string, contents string, version int) (*File, error) {
	u := InternedIds.Save(uri)
	f, err := NewFileFromContents(uri, []byte(contents))
	f.Version = version

	if err != nil {
		return nil, err
	}

	s.Files[u] = f
	f.Parse()
	return f, nil
}

// GetFile returns a file that was stored in the solution.
func (s *Solution) GetFile(uri string) (*File, bool) {

	return nil, false
}

// URI represents a file name
type URI string

// Resolve a path and return the absolute path name after link resolution
// It also internalises the name to get an int id that will later be used
// as part of the absolute identifiers
//
// from is the file path the resolve was started from. It must be an absolute path,
// or the empty string, in which case local resolves will fail.
func Resolve(from string, asked string) (string, error) {
	// Resolving is impacted by a couple of things

	// There are three ways a package could be imported.

	// The first one is to look in the standard library ; such imports do not start by
	// . or / and do not contain a '.' in their name.
	// The compiler will thus go look in its known standard directory path as a base
	// to resolve `asked`.

	// If the asked path is not positioned and contains a '.', then it is a third-party
	// module and will be looked for in the module cache directory.

	// Thus, zoe will look in its configured standard library path as the real base.

	// If a zoe.toml file is found, we look for overrides for package resolution
	// How to make them convienent ? (more so than node_modules links for instance ?)

	path := asked

	// Now that we have the final path computed, resolve all its symlinks
	filepath.EvalSymlinks(path)
	return "", errors.New("module '" + asked + "' could not be found")
}
