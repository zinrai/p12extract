package main

import "errors"

type onceString struct {
	value string
	set   bool
}

func (s *onceString) String() string { return s.value }

func (s *onceString) Set(v string) error {
	if s.set {
		return errors.New("given more than once")
	}
	s.value, s.set = v, true
	return nil
}
