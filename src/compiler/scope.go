package compiler

import (
	"slices"
)

type scope struct {
	name       name
	returnType interface{} //this is optional field for Type Structure representing return type

	variables map[name]variable
	functions map[name]function

	usingScopes []*scope

	parent   *scope
	children map[name]*scope
}

func newScope(n name) *scope {
	return &scope{
		name:        n,
		returnType:  nil,
		variables:   make(map[name]variable),
		functions:   make(map[name]function),
		usingScopes: make([]*scope, 0),
		parent:      nil,
		children:    make(map[name]*scope, 0)}
}

func (s *scope) GetVariable(name name) (variable, bool) {
	if variable, ok := s.getLocalVariable(name); ok {
		return variable, true
	}

	for _, scope := range s.usingScopes {
		if variable, ok := scope.getLocalVariable(name); ok {
			return variable, true
		}
	}

	if s.parent != nil {
		return s.parent.GetVariable(name)
	}
	return variable{}, false
}

func (s *scope) AddVariable(name string, addr address, t Type) bool {
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

func (s *scope) AddFunction(name name, label string, t Type, signature []Type) bool {
	if _, ok := s.functions[name]; ok {
		return false
	}

	s.functions[name] = function{
		Name:      name,
		Label:     label,
		Type:      t,
		Signature: slices.Clone(signature)}
	return true
}

func (s *scope) UsingScope(scope *scope) {
	s.usingScopes = append(s.usingScopes, scope)
}

func (s *scope) GetScope(name name) (*scope, bool) {
	if child, ok := s.children[name]; ok {
		return child, true
	}

	for _, scope := range s.usingScopes {
		if scope.name == name {
			return scope, true
		}

		if _scope, ok := scope.GetScope(name); ok {
			return _scope, true
		}
	}

	if s.parent != nil {
		if s.parent.name == name {
			return s.parent, true
		}
		return s.parent.GetScope(name)
	}

	return nil, false
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
