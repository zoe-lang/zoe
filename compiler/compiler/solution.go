package zoe

import (
	"errors"
	"log"
	"path/filepath"
	"sync"
)

type Solution struct {
	Files map[string]*File
	// The dependencies between the files, where the key is the file being depended upon
	// and the value the files linked to it.
	// Deps map[string][]string
}

func NewSolution() *Solution {
	return &Solution{
		Files: make(map[string]*File),
		// Deps:  make(map[string][]string),
	}
}

type FileCond struct {
	rw *sync.RWMutex
}

func newCondLock() *FileCond {
	return &FileCond{rw: &sync.RWMutex{}}
}

func (fc *FileCond) Acquire() {
	fc.rw.Lock()
}

func (fc *FileCond) SignalDone() {
	fc.rw.Unlock()
}

func (fc *FileCond) Wait() {
	fc.rw.RLock()
	fc.rw.RUnlock()
}

func (s *Solution) AddFile(uri string, contents string, version int) (*File, error) {
	// u := InternedIds.Save(uri)
	f, err := NewFileFromContents(uri, []byte(contents))
	f.Version = version

	s.Files[uri] = f
	if err != nil {
		log.Print(err)
		return nil, err
	}

	f.Parse()
	f.RootNode.Check()
	return f, nil
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

	// 1.
	// The first one is to look in the standard library ; such imports do not start by
	// . or / and do not contain a '.' in their name.
	// The compiler will thus go look in its known standard directory path as a base
	// to resolve `asked`.

	// 2.
	// If the asked path is not positioned and contains a '.', then it is a third-party
	// module and will be looked for in the module cache directory.
	// There will be a simple module manager that will deal with versioning and cloning
	// the package from their git (or other CVS ?) repositories.
	// This approach requires a zoe.toml somewhere that tells the compiler which version
	// of the module it should look for exactly, or if it actually was overriden with
	// a local path for local development changes.

	// 3.
	// If the asked path starts by './', then it is a local import. The name is computed
	// simply by joining asked to from.
	// The use of ../ is disallowed.

	// 4.
	// If the asked path starts by '/', then it is a "pseudo" absolute import, where
	// the real base is the first directory containing a "zoe.toml" file.
	// There is no absolute import to the filesystem, as allowing so would probably
	// create a world of pain.

	path := asked

	// Now that we have the final path computed, resolve all its symlinks
	path, _ = filepath.EvalSymlinks(path)
	return "", errors.New("module '" + asked + "' could not be found")
}
