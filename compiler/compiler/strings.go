package zoe

import (
	"github.com/philpearl/intern"
)

type InternedString int
type Name = InternedString

func (i InternedString) GetText() string {
	return GetInternedString(i)
}

func SaveInternedString(str string) InternedString {
	val := InternedIds.Save(str)
	// log.Print("saved ", str, " -> ", val)
	return InternedString(val)
}

func GetInternedString(id InternedString) string {
	// log.Print("getting ", id, " -> ", InternedIds.Get(int(id)))
	return InternedIds.Get(int(id))
}

var InternedIds *intern.Intern = intern.New(1024)
var instrNext = InternedIds.Save("next") // used in some nodes rewrite with fake strings that are interned.
