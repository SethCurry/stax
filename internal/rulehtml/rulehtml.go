package rulehtml

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"

	"github.com/SethCurry/stax/pkg/ruleparser"
)

// go:embed rules.tmpl
var ruleTemplate string

func GenerateTemplate(parsedRules *ruleparser.Rules, toWriter io.Writer) error {
	parsedTemplate, err := template.New("rules.tmpl").Parse(ruleTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	err = parsedTemplate.Execute(toWriter, parsedRules)
	if err != nil {
		return fmt.Errorf("failed to execute rules template: %w", err)
	}

	return nil
}
