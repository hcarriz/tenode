package tenode_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/hcarriz/tenode"
	"github.com/hcarriz/tenode/tenode_testdata/rendered"
	"github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

func TestTemplToNode(t *testing.T) {
	type args struct {
		ctx      context.Context
		required templ.Component
		input    []templ.Component
	}
	tests := []struct {
		name string
		args args
		want gomponents.Node
	}{
		{
			name: "text",
			args: args{
				ctx:      context.Background(),
				required: rendered.Text("Hello, World!"),
			},
			want: gomponents.Text("Hello, World!"),
		},
		{
			name: "muliple",
			args: args{
				ctx:      context.Background(),
				required: rendered.Text("Hello"),
				input: []templ.Component{
					rendered.Text(", "),
					rendered.Text("World!"),
				},
			},
			want: gomponents.Text("Hello, World!"),
		},
		{
			name: "nested",
			args: args{
				ctx:      context.Background(),
				required: rendered.Basic("Hello, World!"),
			},
			want: html.Article(html.P(gomponents.Text("Hello, World!"))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tenode.TemplToNode(tt.args.ctx, tt.args.required, tt.args.input...)

			input := bytes.NewBuffer(nil)
			result := bytes.NewBuffer(nil)

			if err := got.Render(result); err != nil {
				t.Errorf("unable to render gomponents: %s", err.Error())
				return
			}

			for _, x := range append([]templ.Component{tt.args.required}, tt.args.input...) {
				if err := x.Render(context.Background(), input); err != nil {
					t.Errorf("unable to render templ: %s", err.Error())
					return
				}

			}

			if input.String() != result.String() {
				t.Errorf("TemplToNode() = %v, want %v", result.String(), input.String())
			}

		})
	}
}

func TestNodeToTempl(t *testing.T) {
	type args struct {
		required gomponents.Node
		optional []gomponents.Node
	}
	tests := []struct {
		name string
		args args
		want templ.Component
	}{
		{
			name: "text",
			args: args{
				required: gomponents.Text("Hello, World!"),
			},
			want: rendered.Text("Hello, World!"),
		},
		{
			name: "basic",
			args: args{
				required: html.Article(html.P(gomponents.Text("Hello, World!"))),
			},
			want: rendered.Basic("Hello, World!"),
		},
		{
			name: "basic - extra",
			args: args{
				required: html.Article(html.P(gomponents.Text("Hello, World!"))),
				optional: []gomponents.Node{
					html.Article(html.P(gomponents.Text("Hello, World!"))),
				},
			},
			want: rendered.Basic("Hello, World!", 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tenode.NodeToTempl(tt.args.required, tt.args.optional...)

			input := bytes.NewBuffer(nil)
			result := bytes.NewBuffer(nil)

			if err := got.Render(context.Background(), result); err != nil {
				t.Errorf("unable to render templ: %s", err.Error())
				return
			}

			for _, x := range append([]gomponents.Node{tt.args.required}, tt.args.optional...) {
				if err := x.Render(input); err != nil {
					t.Errorf("unable to render gomponents: %s", err.Error())
					return
				}

			}

			if input.String() != result.String() {
				t.Errorf("NodeToTempl() = %v, want %v", result.String(), input.String())
			}
		})
	}
}
