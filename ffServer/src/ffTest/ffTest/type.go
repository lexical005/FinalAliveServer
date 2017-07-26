package main

type it interface {
	Type() string
}
type c1 struct {
}

func (c *c1) Type() string {
	return "c1"
}

type c2 struct {
	*c1
}

func (c *c2) Type() string {
	return "c2"
}
