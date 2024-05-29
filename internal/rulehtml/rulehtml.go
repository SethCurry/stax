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

func manaSymbolToClass(symbol string) (string, error) {
	switch symbol {
	case "t", "T":
		return "ms-tap ms-cost ms-shadow", nil
	case "q", "Q":
		return "ms-untap ms-cost ms-shadow", nil
	case "rn", "rN":
		return "ms-saga", nil
	case "rn1", "rN1":
		return "ms-saga ms-saga-1", nil
	case "rn2", "rN2":
		return "ms-saga ms-saga-2", nil
	case "rn3", "rN3":
		return "ms-saga ms-saga-3", nil
	case "rn4", "rN4":
		return "ms-saga ms-saga-4", nil
	}
	return "ms-cost ms-shadow ms-" + strings.Replace(strings.ToLower(symbol), "/", "", -1), nil
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
