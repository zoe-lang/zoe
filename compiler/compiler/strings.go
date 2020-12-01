package zoe

import "github.com/philpearl/intern"

var InternedIds *intern.Intern = intern.New(1024)
var instrNext = InternedIds.Save("next") // used in some nodes rewrite with fake strings that are interned.
