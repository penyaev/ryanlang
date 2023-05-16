package compiler

type symbolScope int

const (
	symbolScopeLocal symbolScope = iota
	symbolScopeForeign
	symbolScopeBuiltin
)

func (s symbolScope) String() string {
	switch s {
	case symbolScopeLocal:
		return "local"
	case symbolScopeForeign:
		return "foreign"
	case symbolScopeBuiltin:
		return "builtin"
	default:
		panic("unknown scope")
	}
}

type Symbol struct {
	name  string
	id    int
	scope symbolScope
}
type symbols struct {
	scopeId    int
	storage    map[string]*Symbol
	linkedRoot *symbols
	foreign    []*Symbol
	parent     *symbols
	locals     *int
}

func newSymbols() *symbols {
	var locals int
	s := &symbols{
		storage: map[string]*Symbol{},
		locals:  &locals,
	}
	s.linkedRoot = s
	return s
}

func (s *symbols) push() *symbols {
	next := newSymbols()
	next.parent = s
	next.scopeId = s.scopeId + 1
	return next
}
func (s *symbols) pushLinked() *symbols {
	next := newSymbols()
	next.parent = s
	next.locals = s.locals
	next.scopeId = s.scopeId // preserve scopeId
	next.linkedRoot = s.linkedRoot
	return next
}
func (s *symbols) pop() *symbols {
	return s.parent
}
func (s *symbols) getWithScope(key string, scope symbolScope) *Symbol {
	ret, ok := s.storage[key]
	if !ok || ret.scope != scope {
		ret = nil
	}
	if ret == nil && s.parent != nil {
		ret = s.parent.getWithScope(key, scope)
	}
	return ret
}
func (s *symbols) get(key string) (*Symbol, int) {
	var scopeId = s.scopeId
	ret, ok := s.storage[key]
	if !ok && s.parent != nil {
		ret, scopeId = s.parent.get(key)
		if ret == nil {
			return ret, scopeId
		}
		if (ret.scope == symbolScopeForeign || ret.scope == symbolScopeLocal) && (scopeId != s.scopeId) {
			ret = s.createForeign(key, ret)
			scopeId = s.scopeId
		}
	}
	return ret, scopeId
}
func (s *symbols) getLocal(key string) *Symbol {
	return s.storage[key]
}
func (s *symbols) hasLocal(key string) bool {
	sym, ok := s.storage[key]
	return ok && sym.scope == symbolScopeLocal
}
func (s *symbols) hasForeign(key string) bool {
	sym, ok := s.storage[key]
	return ok && sym.scope == symbolScopeForeign
}
func (s *symbols) createLocal(key string) *Symbol {
	ret := &Symbol{
		name:  key,
		id:    *s.linkedRoot.locals,
		scope: symbolScopeLocal,
	}
	s.set(key, ret)
	*s.linkedRoot.locals++
	return ret
}
func (s *symbols) importLocal(sym *Symbol) {
	s.set(sym.name, sym)
}
func (s *symbols) createForeign(key string, original *Symbol) *Symbol {
	ret := &Symbol{
		name:  key,
		id:    len(s.linkedRoot.foreign),
		scope: symbolScopeForeign,
	}
	s.linkedRoot.foreign = append(s.linkedRoot.foreign, original)
	s.set(key, ret)
	return ret
}
func (s *symbols) createBuiltin(key string, id int) *Symbol {
	ret := &Symbol{
		name:  key,
		id:    id,
		scope: symbolScopeBuiltin,
	}
	s.set(key, ret)
	return ret
}
func (s *symbols) getLocalOrCreate(key string) *Symbol {
	ret := s.getLocal(key)
	if ret == nil {
		ret = s.createLocal(key)
	}
	return ret
}
func (s *symbols) set(key string, symbol *Symbol) {
	s.storage[key] = symbol
}
