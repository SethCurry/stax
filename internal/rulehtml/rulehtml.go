package rulehtml

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/SethCurry/stax/pkg/ruleparser"
)

//go:embed rules.tmpl
var rootTemplate string

//go:embed mana.css
var manaCss string

func manaSymbolToClass(symbol string) (string, error) {
	switch symbol {
	case "t":
		return "ms-tap", nil
	}
	return "ms-" + strings.Replace(strings.ToLower(symbol), "/", "", -1), nil
}

func getElementID(elementName string) (string, error) {
	return strings.Replace(elementName, ".", "_", -1), nil
}

func GenerateTemplate(parsedRules *ruleparser.Rules, toWriter io.Writer) error {
	parsedTemplate, err := template.New("rules.tmpl").Funcs(template.FuncMap{
		"ManaClass": manaSymbolToClass,
		"ElementID": getElementID,
	}).Parse(rootTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	err = parsedTemplate.Execute(toWriter, parsedRules)
	if err != nil {
		return fmt.Errorf("failed to execute rules template: %w", err)
	}

	return nil
}
