package renderer

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"regexp"

	"github.com/yuin/goldmark"
)

var standAloneUrlRE = regexp.MustCompile(`[^(\(.*\)\[)]https?://[^\s]+[^\]]`)
var exactUrlRE = regexp.MustCompile(`https?://[^\s]+`)
var gamingRE = regexp.MustCompile(`\+\+.+\+\+`)
var gamingREInner = regexp.MustCompile(`[^(\+\+)].+[^(\+\+)]`)
var linkTmpl = template.Must(template.New("link").Parse(`<a href="{{.}}">{{.}}</a>`))

// Render は受け取った文書を HTML に変換する
func Render(ctx context.Context, src string) (string, error) {
	// TODO: これはサンプル実装 (URL の自動リンク) です
	standAloneUrlConverted := standAloneUrlRE.ReplaceAllStringFunc(src, func(standAloneUrl string) string {
		return exactUrlRE.ReplaceAllStringFunc(standAloneUrl, func(exactUrl string) string {
			return fmt.Sprintf("[%s](%s)", exactUrl, exactUrl)
		})
	})
	var htmlBuf bytes.Buffer
	if err := goldmark.Convert([]byte(standAloneUrlConverted), &htmlBuf); err != nil {
		panic(err)
	}
	fullConverted := gamingRE.ReplaceAllStringFunc(htmlBuf.String(), func(inner string) string {
		return fmt.Sprintf("<span class=\"gaming\">%s</span>", inner[2:len(inner)-2])
	})
	return fullConverted, nil
}
