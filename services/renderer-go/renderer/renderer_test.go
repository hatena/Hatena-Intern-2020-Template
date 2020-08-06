package renderer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Render(t *testing.T) {
	src := `foo https://google.com/ bar`
	html, err := Render(context.Background(), src)
	assert.NoError(t, err)
	assert.Equal(t, `foo <a href="https://google.com/">https://google.com/</a> bar`, html)
}
