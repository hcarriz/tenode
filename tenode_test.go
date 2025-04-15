package tenode_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/hcarriz/tenode"
	rendered "github.com/hcarriz/tenode/internal/tenode_testdata"
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func TestTemplToNode(t *testing.T) {
	tests := []struct {
		name string
		args templ.Component
		want gomponents.Node
	}{
		{
			name: "text",
			args: rendered.Text("Hello, World!"),
			want: gomponents.Text("Hello, World!"),
		},
		{
			name: "multiple",
			args: templ.Join(rendered.Text("Hello"),

				rendered.Text(", "),
				rendered.Text("World!"),
			),
			want: gomponents.Text("Hello, World!"),
		},
		{
			name: "nested",
			args: rendered.Basic("Hello, World!"),
			want: html.Article(html.P(gomponents.Text("Hello, World!"))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tenode.ToGomponent(tt.args)

			input := bytes.NewBuffer(nil)
			result := bytes.NewBuffer(nil)

			if err := got.Render(result); err != nil {
				t.Errorf("unable to render gomponents: %s", err.Error())
				return
			}

			if err := tt.args.Render(context.Background(), input); err != nil {
				t.Errorf("unable to render templ: %s", err.Error())
				return
			}

			if input.String() != result.String() {
				t.Errorf("TemplToNode() = %v, want %v", result.String(), input.String())
			}

		})
	}
}

func TestNodeToTempl(t *testing.T) {
	tests := []struct {
		name string
		args gomponents.Node
		want templ.Component
	}{
		{
			name: "text",
			args: gomponents.Text("Hello, World!"),
			want: rendered.Text("Hello, World!"),
		},
		{
			name: "basic",
			args: html.Article(html.P(gomponents.Text("Hello, World!"))),
			want: rendered.Basic("Hello, World!"),
		},
		{
			name: "basic - extra",
			args: gomponents.Group{
				html.Article(html.P(gomponents.Text("Hello, World!"))),
				html.Article(html.P(gomponents.Text("Hello, World!"))),
			},
			want: rendered.Basic("Hello, World!", 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tenode.ToTempl(tt.args)

			input := bytes.NewBuffer(nil)
			result := bytes.NewBuffer(nil)

			if err := got.Render(context.Background(), result); err != nil {
				t.Errorf("unable to render templ: %s", err.Error())
				return
			}

			if err := tt.args.Render(input); err != nil {
				t.Errorf("unable to render gomponents: %s", err.Error())
				return
			}

			if input.String() != result.String() {
				t.Errorf("NodeToTempl() = %v, want %v", result.String(), input.String())
			}
		})
	}
}
