// Code generated by ent, DO NOT EDIT.

package bones

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/SethCurry/stax/internal/bones/card"
)

// Card is the model entity for the Card schema.
type Card struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// OracleID holds the value of the "oracle_id" field.
	OracleID string `json:"oracle_id,omitempty"`
	// ColorIdentity holds the value of the "color_identity" field.
	ColorIdentity uint8 `json:"color_identity,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the CardQuery when eager-loading is set.
	Edges        CardEdges `json:"edges"`
	selectValues sql.SelectValues
}

// CardEdges holds the relations/edges for other nodes in the graph.
type CardEdges struct {
	// Faces holds the value of the faces edge.
	Faces []*CardFace `json:"faces,omitempty"`
	// Rulings holds the value of the rulings edge.
	Rulings []*Ruling `json:"rulings,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// FacesOrErr returns the Faces value or an error if the edge
// was not loaded in eager-loading.
func (e CardEdges) FacesOrErr() ([]*CardFace, error) {
	if e.loadedTypes[0] {
		return e.Faces, nil
	}
	return nil, &NotLoadedError{edge: "faces"}
}

// RulingsOrErr returns the Rulings value or an error if the edge
// was not loaded in eager-loading.
func (e CardEdges) RulingsOrErr() ([]*Ruling, error) {
	if e.loadedTypes[1] {
		return e.Rulings, nil
	}
	return nil, &NotLoadedError{edge: "rulings"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Card) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case card.FieldID, card.FieldColorIdentity:
			values[i] = new(sql.NullInt64)
		case card.FieldName, card.FieldOracleID:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Card fields.
func (c *Card) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case card.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			c.ID = int(value.Int64)
		case card.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				c.Name = value.String
			}
		case card.FieldOracleID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field oracle_id", values[i])
			} else if value.Valid {
				c.OracleID = value.String
			}
		case card.FieldColorIdentity:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field color_identity", values[i])
			} else if value.Valid {
				c.ColorIdentity = uint8(value.Int64)
			}
		default:
			c.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Card.
// This includes values selected through modifiers, order, etc.
func (c *Card) Value(name string) (ent.Value, error) {
	return c.selectValues.Get(name)
}

// QueryFaces queries the "faces" edge of the Card entity.
func (c *Card) QueryFaces() *CardFaceQuery {
	return NewCardClient(c.config).QueryFaces(c)
}

// QueryRulings queries the "rulings" edge of the Card entity.
func (c *Card) QueryRulings() *RulingQuery {
	return NewCardClient(c.config).QueryRulings(c)
}

// Update returns a builder for updating this Card.
// Note that you need to call Card.Unwrap() before calling this method if this Card
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Card) Update() *CardUpdateOne {
	return NewCardClient(c.config).UpdateOne(c)
}

// Unwrap unwraps the Card entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Card) Unwrap() *Card {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("bones: Card is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Card) String() string {
	var builder strings.Builder
	builder.WriteString("Card(")
	builder.WriteString(fmt.Sprintf("id=%v, ", c.ID))
	builder.WriteString("name=")
	builder.WriteString(c.Name)
	builder.WriteString(", ")
	builder.WriteString("oracle_id=")
	builder.WriteString(c.OracleID)
	builder.WriteString(", ")
	builder.WriteString("color_identity=")
	builder.WriteString(fmt.Sprintf("%v", c.ColorIdentity))
	builder.WriteByte(')')
	return builder.String()
}

// Cards is a parsable slice of Card.
type Cards []*Card
