package zoe

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
