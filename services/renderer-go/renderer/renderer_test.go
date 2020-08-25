package renderer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Render_StandAloneURL(t *testing.T) {
	src := `foo https://google.com/ bar`
	html, err := Render(context.Background(), src)
	assert.NoError(t, err)
	assert.Equal(t, "<p>foo <a href=\"https://google.com/\">https://google.com/</a> bar</p>\n", html)
}

// TODO Test abnormal case

func Test_Render_Heading(t *testing.T) {
	src := "# h1\nfoo\n## h2\nhoge"
	html, err := Render(context.Background(), src)
	assert.NoError(t, err)
	assert.Equal(t, "<h1>h1</h1>\n<p>foo</p>\n<h2>h2</h2>\n<p>hoge</p>\n", html)
}

func Test_Render_Link(t *testing.T) {
	src := "[namachan10777](https://namachan10777.dev)"
	html, err := Render(context.Background(), src)
	assert.NoError(t, err)
	assert.Equal(t, "<p><a href=\"https://namachan10777.dev\">namachan10777</a></p>\n", html)
}

func Test_Render_List(t *testing.T) {
	src := "* li1\n* li2\n* li3"
	html, err := Render(context.Background(), src)
	assert.NoError(t, err)
	assert.Equal(t, "<ul>\n<li>li1</li>\n<li>li2</li>\n<li>li3</li>\n</ul>\n", html)
}
