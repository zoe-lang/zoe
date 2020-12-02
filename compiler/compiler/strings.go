package zoe

import "github.com/philpearl/intern"

type InternedString int

func SaveInternedString(str string) InternedString {
	val := InternedIds.Save(str)
	return InternedString(val)
}

func GetInternedString(id InternedString) string {
	return InternedIds.Get(int(id))
}

var InternedIds *intern.Intern = intern.New(1024)
var instrNext = InternedIds.Save("next") // used in some nodes rewrite with fake strings that are interned.
