package renderer

import (
	"bytes"
	"context"
	"html/template"
	"regexp"
)

var standAloneUrlRE = regexp.MustCompile(`[^(\(.*\)\[)]https?://[^\s]+[^\]]`)
var exactUrlRE = regexp.MustCompile(`https?://[^\s]+`)
var linkTmpl = template.Must(template.New("link").Parse(`<a href="{{.}}">{{.}}</a>`))

// Render は受け取った文書を HTML に変換する
func Render(ctx context.Context, src string) (string, error) {
	// TODO: これはサンプル実装 (URL の自動リンク) です
	html := standAloneUrlRE.ReplaceAllStringFunc(src, func(standAloneUrl string) string {
		return exactUrlRE.ReplaceAllStringFunc(standAloneUrl, func(exactUrl string) string {
			var w bytes.Buffer
			err := linkTmpl.ExecuteTemplate(&w, "link", exactUrl)
			if err != nil {
				return exactUrl
			}
			return w.String()
		})
	})
	return html, nil
}
