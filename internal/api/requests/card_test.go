package requests

import (
	"errors"
	"testing"

	"entgo.io/ent/dialect/sql"
	"github.com/SethCurry/stax/internal/oracle/oracledb/card"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	tables := []struct {
		c   CardByName
		err error
	}{
		{CardByName{"Exact", "Fuzzy"}, errors.New("exact and fuzzy cannot be used at the same time")},
		{CardByName{"", ""}, errors.New("either fuzzy or exact must be specified")},
		{CardByName{"Exact", ""}, nil},
		{CardByName{"", "Fuzzy"}, nil},
	}
	for _, table := range tables {
		err := table.c.Validate()
		if err != nil && err.Error() != table.err.Error() {
			t.Errorf("Validation of %v failed, expected: '%s', got:  '%s'", table.c, table.err, err)
		}
	}
}

func TestAddToSQL(t *testing.T) {
	// Assuming you have a function card.NameEQ and card.NameContainsFold which return appropriate conditions
	tables := []struct {
		name         string
		c            CardByName
		expected     string
		expectedArgs []any
	}{
		{"fuzzy search", CardByName{"", "Fuzzy"}, "SELECT * FROM `cards` WHERE LOWER(`cards`.`name`) LIKE ?", []interface{}{"Fuzzy"}}, // You need to provide the expected SQL based on your actual implementation of card.NameEQ and card.NameContainsFold
		{"exact search", CardByName{"Exact", ""}, "SELECT * FROM `cards` WHERE `cards`.`name` = ?", []interface{}{"Exact"}},
		{"invalid search", CardByName{"", ""}, "SELECT * FROM `cards`", []interface{}{}},
	}
	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			pred := table.c.ToPredicate()

			query := sql.Select("*").From(sql.Table(card.Table))
			pred(query)

			asString, _ := query.Query()

			assert.Equal(t, table.expected, asString)
		})
	}
}
