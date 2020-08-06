package grpc

import (
	"context"
	"testing"

	pb "github.com/hatena/Hatena-Intern-2020/services/renderer-go/pb/renderer"
	"github.com/stretchr/testify/assert"
)

func Test_Server_Render(t *testing.T) {
	s := NewServer()
	src := `foo https://google.com/ bar`
	reply, err := s.Render(context.Background(), &pb.RenderRequest{Src: src})
	assert.NoError(t, err)
	assert.Equal(t, `foo <a href="https://google.com/">https://google.com/</a> bar`, reply.Html)
}
