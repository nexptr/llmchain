package prompts

import (
	"bytes"
	"io"
	"strings"
	"text/template"
)

type TemplateConfig struct {
	Completion string `json:"completion" yaml:"completion"`
	Chat       string `json:"chat" yaml:"chat"`
	Edit       string `json:"edit" yaml:"edit"`
}

type Template struct {
	tmpl    *template.Template
	funcMap template.FuncMap
}

// NewTemplate 创建初始化渲染类，输入参数为当前模版库路径
func NewTemplate(path string) *Template {

	p := &Template{funcMap: template.FuncMap{
		// "CustomPortContent": CustomPortContent,
	}}

	if path == `` {
		path = "prompts/*.tmpl"
	}

	p.tmpl = template.Must(template.New("").Funcs(p.funcMap).ParseGlob(path))

	return p
}

func (p *Template) Execute(wr io.Writer, name string, data any) error {

	// if !strings.HasSuffix(name, `.tmpl`) {
	// 	name = name + `.tmpl`
	// }

	return p.tmpl.ExecuteTemplate(wr, name, data)
}

func (p *Template) Render(name string, data any) (string, error) {

	if !strings.HasSuffix(name, `.tmpl`) {
		name = name + `.tmpl`
	}

	var buf bytes.Buffer

	err := p.tmpl.ExecuteTemplate(&buf, name, data)
	return buf.String(), err
}
