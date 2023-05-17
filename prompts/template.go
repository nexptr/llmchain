package prompts

import (
	"bytes"
	"fmt"
	"text/template"
)

type H = map[string]string

type Template struct {
	tmpl      *template.Template
	variables []string
}

func PromptTemplate(prompt string, variables ...string) *Template {
	t, err := New(prompt, variables...)
	if err != nil {
		fmt.Println(err)
	}
	return t
}

func New(prompt string, variables ...string) (*Template, error) {

	tmpl, err := template.New(``).Parse(prompt)

	if err != nil {
		return nil, fmt.Errorf(`paser tmpl failed: %v`, err)
	}

	t := &Template{
		tmpl:      tmpl,
		variables: variables,
	}

	if len(t.variables) == 0 {
		t.variables = []string{`input`}
	}

	return t, nil
}

func (t *Template) Render(vars H) (string, error) {

	//todo verify vars
	for _, v := range t.variables {
		if _, ok := vars[v]; !ok {
			return "", fmt.Errorf(`field: '%s' not set`, v)
		}
	}

	var buf bytes.Buffer

	err := t.tmpl.Execute(&buf, vars)

	return buf.String(), err

}
