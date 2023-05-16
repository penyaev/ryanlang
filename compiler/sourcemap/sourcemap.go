package sourcemap

import "ryanlang/lexer"

type Entry struct {
	from     int
	to       int
	Loc      *lexer.Location
	Comment  string
	parent   *Entry
	children []*Entry
}

func (e *Entry) search(p int) *Entry {
	if p < e.from || p >= e.to {
		return nil
	}

	for _, entry := range e.children {
		if entry.from <= p && entry.to > p {
			return entry.search(p)
		}
	}
	return e
}

func (e *Entry) shift(offset int) {
	e.from += offset
	e.to += offset
	for _, entry := range e.children {
		entry.shift(offset)
	}
}

func (e *Entry) Location() *lexer.Location {
	cur := e
	for cur.Loc == nil && cur.parent != nil {
		cur = cur.parent
	}
	return cur.Loc
}

type SourceMap struct {
	root *Entry
	cur  *Entry
}

func New() *SourceMap {
	root := &Entry{}
	return &SourceMap{
		root: root,
		cur:  root,
	}
}

func (s *SourceMap) Push(p int, loc *lexer.Location, comment string) {
	n := &Entry{
		from:     p,
		to:       0,
		Loc:      loc,
		Comment:  comment,
		parent:   s.cur,
		children: nil,
	}
	s.cur.children = append(s.cur.children, n)
	s.cur = n
}

func (s *SourceMap) Pop(p int) {
	s.cur.to = p
	s.cur = s.cur.parent
}

func (s *SourceMap) PopAll(p int) {
	for s.cur.parent != nil {
		s.Pop(p)
	}
	s.cur.to = p
}

func (s *SourceMap) Search(p int) *Entry {
	return s.root.search(p)
}

func (s *SourceMap) Merge(other *SourceMap, at int) {
	other.root.shift(at)
	s.cur.children = append(s.cur.children, other.root)
	other.root.parent = s.cur
}
