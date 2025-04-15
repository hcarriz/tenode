package tenode

import (
	"context"
	"io"
)

type TemplInterface interface {
	Render(context.Context, io.Writer) error
}

type GomponentInterface interface {
	Render(io.Writer) error
}

type Gomponent struct {
	templ TemplInterface
}

var _ GomponentInterface = (*Gomponent)(nil)

func (n Gomponent) Render(w io.Writer) error {
	return n.templ.Render(context.Background(), w)
}

var _ TemplInterface = (*Templ)(nil)

type Templ struct {
	gomponent GomponentInterface
}

func (t Templ) Render(_ context.Context, w io.Writer) error {
	return t.gomponent.Render(w)
}

// ToGomponent allows for templ components to be used as a gomponent.
func ToGomponent(component TemplInterface) Gomponent {
	return Gomponent{
		templ: component,
	}
}

// ToTempl allows for gomponents to be used in a templ component.
func ToTempl(component GomponentInterface) Templ {
	return Templ{
		gomponent: component,
	}
}
