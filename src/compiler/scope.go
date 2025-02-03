package compiler

type scope struct {
	name       name
	returnType interface{} //this is optional field for string representing return type

	variables map[name]variable
	functions map[name]function
	types     map[name]name

	usingScopes []name

	parent   *scope
	children map[name]*scope
}

func newScope(n name) *scope {
	return &scope{
		name:        n,
		returnType:  nil,
		variables:   make(map[name]variable),
		functions:   make(map[name]function),
		types:       make(map[name]name),
		usingScopes: make([]name, 0),
		parent:      nil,
		children:    make(map[name]*scope, 0)}
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

func (s *scope) UsingScope(name string) {
	s.usingScopes = append(s.usingScopes, name)
}
