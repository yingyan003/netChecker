package sorter

import (
	"sort"
)

func (s *Sorter) Sort(elems []IResource) {
	s.elems = elems
	sort.Sort(s)
}

func (s *Sorter) Len() int {
	return len(s.elems)
}

func (s *Sorter) Swap(i, j int) {
	s.elems[i], s.elems[j] = s.elems[j], s.elems[i]
}

func (s *Sorter) Less(i, j int) bool {
	p, q := s.elems[i], s.elems[j]

	var k int
	for k = 0; k < len(s.Lessf)-1; k++ {
		Lessf := s.Lessf[k]
		switch {
		case Lessf(p, q):
			return true
		case Lessf(q, p):
			return false
		}
	}

	return s.Lessf[k](p, q)
}
