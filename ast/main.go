package ast

import (
	"fmt"

	"github.com/jacobsimpson/mp3tag/metadata"
)

type Expression interface {
	String() string
}

type Equal struct {
	LHS metadata.Name
	RHS string
}

func (e *Equal) String() string {
	return fmt.Sprintf("%s=%q", e.LHS, e.RHS)
}

type Has struct {
	LHS metadata.Name
	RHS string
}

func (h *Has) String() string {
	return fmt.Sprintf("%s:%q", h.LHS, h.RHS)
}
