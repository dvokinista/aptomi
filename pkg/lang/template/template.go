package template

import (
	"bytes"
	"fmt"
	"github.com/Aptomi/aptomi/pkg/errors"
	"strings"
	t "text/template"
)

// Template struct contains text template string as well as its compiled version
type Template struct {
	templateStr      string
	templateCompiled *t.Template
}

// NewTemplate compiles a text template and returns the result in Template struct
// Parameter templateStr must follow syntax defined by text/template
func NewTemplate(templateStr string) (*Template, error) {
	templateCompiled, err := t.New("").Parse(templateStr)
	if err != nil {
		return nil, fmt.Errorf("unable to compile template '%s': %s", templateStr, err)
	}
	return &Template{
		templateStr:      templateStr,
		templateCompiled: templateCompiled,
	}, nil
}

// Evaluate evaluates a compiled text template given a set named parameters
func (template *Template) Evaluate(params *Parameters) (string, error) {
	// Evaluate
	var doc bytes.Buffer

	// Multiple executions of the same template can execute safely in parallel
	err := template.templateCompiled.Execute(&doc, params.params)
	if err != nil {
		return "", errors.NewErrorWithDetails(
			fmt.Sprintf("Unable to evaluate template '%s': %s", template.templateStr, err),
			errors.Details{
				"template": template.templateStr,
				"params":   params,
			},
		)
	}

	// Convert result to bool
	result := doc.String()
	if strings.Contains(result, "<no value>") {
		return "", errors.NewErrorWithDetails(
			fmt.Sprintf("Unable to evaluate template '%s': <no value>", template.templateStr),
			errors.Details{
				"template": template.templateStr,
				"params":   params,
			},
		)
	}

	return result, nil
}
