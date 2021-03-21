package parser

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

func Parse(query string) (Expression, error) {
	return &Equal{
		LHS: "title",
		RHS: "A Thousand Years",
	}, nil
}
