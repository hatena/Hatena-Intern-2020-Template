package web

//go:generate go-assets-builder --package=web --output=./templates_gen.go --strip-prefix="/templates/" --variable=templateAssets ../templates

import (
	"html/template"
	"io"
	"io/ioutil"

	echo "github.com/labstack/echo/v4"
)

type renderer struct {
	templates map[string]*template.Template
}

func newRenderer() (*renderer, error) {
	templates := make(map[string]*template.Template)

	bs, err := ioutil.ReadAll(templateAssets.Files["wrapper.html"])
	templateAssets.Files["wrapper.html"].Close()
	if err != nil {
		return nil, err
	}
	main := template.Must(template.New("main").Parse(string(bs))).Funcs(template.FuncMap{
		"unescapedHTML": func(html string) template.HTML {
			return template.HTML(html)
		},
	})

	for fileName, file := range templateAssets.Files {
		bs, err := ioutil.ReadAll(file)
		file.Close()
		if err != nil {
			return nil, err
		}
		main := template.Must(main.Clone())
		templates[fileName] = template.Must(main.New(fileName).Parse(string(bs)))
	}

	return &renderer{templates}, nil
}

func (r *renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if data == nil {
		data = map[string]interface{}{}
	}
	data.(map[string]interface{})["CsrfToken"] = c.Get("csrf")
	return r.templates[name].ExecuteTemplate(w, "main", data)
}
