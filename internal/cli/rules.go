package cli

import (
	"fmt"
	"os"

	"github.com/SethCurry/stax/internal/rulehtml"
	"github.com/SethCurry/stax/pkg/ruleparser"
	"go.uber.org/zap"
)

type RulesCmd struct {
	HTML RulesHTMLCmd `cmd:"" help:"Generate stand-alone HTML documentation of the rules."`
}

type RulesHTMLCmd struct {
	RulesFile string `arg:"rules-file" help:"The .txt rules file to parse rules from."`
	Output    string `type:"path" optional:"output" default:"rules.html" help:"The file to write the output to."`
}

func (cmd *RulesHTMLCmd) Run(ctx *Context) error {
	logger := ctx.Logger

	logger.Info("parsing rules", zap.String("rules_file", cmd.RulesFile))

	parsedRules, err := ruleparser.ParseFile(cmd.RulesFile)
	if err != nil {
		return fmt.Errorf("failed to parse rules from file %q: %w", cmd.RulesFile, err)
	}

	outFd, err := os.Create(cmd.Output)
	if err != nil {
		logger.Fatal("failed to open output file", zap.String("path", cmd.Output), zap.Error(err))
	}
	defer outFd.Close()

	err = rulehtml.GenerateTemplate(parsedRules, outFd)
	if err != nil {
		logger.Fatal("failed to generate templated HTML", zap.Error(err))
	}

	return nil
}
