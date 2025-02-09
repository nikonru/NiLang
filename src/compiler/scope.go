package compiler

type scope struct {
	name       name
	returnType interface{} //this is optional field for string representing return type

	variables map[name]variable
	functions map[name]function

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
		usingScopes: make([]name, 0),
		parent:      nil,
		children:    make(map[name]*scope, 0)}
}

func (s *scope) GetVariable(name name) (variable, bool) {
	if variable, ok := s.getLocalVariable(name); ok {
		return variable, true
	}

	for _, scope := range s.usingScopes {
		if child, ok := s.children[scope]; ok {
			if variable, ok := child.getLocalVariable(name); ok {
				return variable, true
			}
		}
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

func (s *scope) AddScope(scope *scope) bool {
	if _, ok := s.children[scope.name]; ok {
		return false
	}

	s.children[scope.name] = scope
	return true
}

func (s *scope) UsingScope(name string) {
	s.usingScopes = append(s.usingScopes, name)
}

func (s *scope) GetParent() *scope {
	return s.parent
}

func (s *scope) SetParent(scope *scope) {
	s.parent = scope
}

func (s *scope) getLocalVariable(name name) (variable, bool) {
	if variable, ok := s.variables[name]; ok {
		return variable, true
	}

	return variable{}, false
}
