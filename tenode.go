package tenode

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/maragudk/gomponents"
)

// TemplToNode allows for templ components to be used as a gomponent.
func TemplToNode(ctx context.Context, required templ.Component, input ...templ.Component) gomponents.Node {

	components := append([]templ.Component{required}, input...)

	group := gomponents.Map(components, func(in templ.Component) gomponents.Node {
		return gomponents.NodeFunc(func(w io.Writer) error {
			return in.Render(ctx, w)
		})
	})

	return gomponents.NodeFunc(func(w io.Writer) error {
		for _, x := range group {
			if err := x.Render(w); err != nil {
				return err
			}
		}
		return nil
	})
}

// NodeToTempl allows for gomponents to be used in a templ component.
func NodeToTempl(required gomponents.Node, optional ...gomponents.Node) templ.Component {

	list := append([]gomponents.Node{required}, optional...)

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		for _, x := range list {
			if err := x.Render(w); err != nil {
				return err
			}
		}
		return nil
	})
}
