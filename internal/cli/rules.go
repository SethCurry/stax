package cli

import (
	"bytes"
	"fmt"

	"github.com/SethCurry/stax/internal/rulehtml"
	"github.com/SethCurry/stax/pkg/ruleparser"
	"go.uber.org/zap"
)

type RulesCmd struct {
	Judge RulesJudgeCmd `cmd:"" help:"Generate SQL for updating the Judge rules website."`
}

type RulesJudgeCmd struct {
	RulesFile string `arg:"rules-file" help:"The .txt rules file to parse rules from."`
}

func (cmd *RulesJudgeCmd) Run(ctx *Context) error {
	logger := ctx.Logger

	logger.Info("parsing rules", zap.String("rules_file", cmd.RulesFile))

	parsedRules, err := ruleparser.ParseFile(cmd.RulesFile)
	if err != nil {
		return fmt.Errorf("failed to parse rules from file %q: %w", cmd.RulesFile, err)
	}

	var buf bytes.Buffer

	err = rulehtml.GenerateTemplate(parsedRules, &buf)
	if err != nil {
		logger.Error("failed to generate template for rules", zap.String("rules_file", cmd.RulesFile), zap.Error(err))
	}

	fmt.Println(buf.String())

	return nil
}
