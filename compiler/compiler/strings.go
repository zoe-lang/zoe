package zoe

import "github.com/philpearl/intern"

var internedIds *intern.Intern = intern.New(1024)
var instrNext = internedIds.Save("next") // used in some nodes rewrite with fake strings that are interned.
