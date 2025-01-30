package compiler

type scope struct {
	name name

	variables map[name]variable
	functions map[name]function
	types     map[name]name

	parent   *scope
	children []*scope
}

func newScope(n name) *scope {
	return &scope{
		name:      n,
		variables: make(map[name]variable),
		functions: make(map[name]function),
		types:     make(map[name]name),
		parent:    nil,
		children:  make([]*scope, 0)}
}

func (s *scope) GetVariable(name name) (variable, bool) {
	if variable, ok := s.variables[name]; ok {
		return variable, true
	}
	if s.parent != nil {
		return s.parent.GetVariable(name)
	}
	return variable{}, false
}

func (s *scope) AddVariable(name string, addr address, t name) bool {
	if _, ok := s.variables[name]; ok {
		return false
	}

	s.variables[name] = variable{Name: name, Addr: addr, Type: t}
	return true
}
