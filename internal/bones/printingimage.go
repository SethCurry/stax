// Code generated by ent, DO NOT EDIT.

package bones

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/SethCurry/stax/internal/bones/printing"
	"github.com/SethCurry/stax/internal/bones/printingimage"
)

// PrintingImage is the model entity for the PrintingImage schema.
type PrintingImage struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// URL holds the value of the "url" field.
	URL string `json:"url,omitempty"`
	// ImageType holds the value of the "image_type" field.
	ImageType printingimage.ImageType `json:"image_type,omitempty"`
	// LocalPath holds the value of the "local_path" field.
	LocalPath string `json:"local_path,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PrintingImageQuery when eager-loading is set.
	Edges                   PrintingImageEdges `json:"edges"`
	printing_image_printing *int
	selectValues            sql.SelectValues
}

// PrintingImageEdges holds the relations/edges for other nodes in the graph.
type PrintingImageEdges struct {
	// Printing holds the value of the printing edge.
	Printing *Printing `json:"printing,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// PrintingOrErr returns the Printing value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PrintingImageEdges) PrintingOrErr() (*Printing, error) {
	if e.Printing != nil {
		return e.Printing, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: printing.Label}
	}
	return nil, &NotLoadedError{edge: "printing"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*PrintingImage) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case printingimage.FieldID:
			values[i] = new(sql.NullInt64)
		case printingimage.FieldURL, printingimage.FieldImageType, printingimage.FieldLocalPath:
			values[i] = new(sql.NullString)
		case printingimage.ForeignKeys[0]: // printing_image_printing
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the PrintingImage fields.
func (pi *PrintingImage) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case printingimage.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			pi.ID = int(value.Int64)
		case printingimage.FieldURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field url", values[i])
			} else if value.Valid {
				pi.URL = value.String
			}
		case printingimage.FieldImageType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field image_type", values[i])
			} else if value.Valid {
				pi.ImageType = printingimage.ImageType(value.String)
			}
		case printingimage.FieldLocalPath:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field local_path", values[i])
			} else if value.Valid {
				pi.LocalPath = value.String
			}
		case printingimage.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field printing_image_printing", value)
			} else if value.Valid {
				pi.printing_image_printing = new(int)
				*pi.printing_image_printing = int(value.Int64)
			}
		default:
			pi.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the PrintingImage.
// This includes values selected through modifiers, order, etc.
func (pi *PrintingImage) Value(name string) (ent.Value, error) {
	return pi.selectValues.Get(name)
}

// QueryPrinting queries the "printing" edge of the PrintingImage entity.
func (pi *PrintingImage) QueryPrinting() *PrintingQuery {
	return NewPrintingImageClient(pi.config).QueryPrinting(pi)
}

// Update returns a builder for updating this PrintingImage.
// Note that you need to call PrintingImage.Unwrap() before calling this method if this PrintingImage
// was returned from a transaction, and the transaction was committed or rolled back.
func (pi *PrintingImage) Update() *PrintingImageUpdateOne {
	return NewPrintingImageClient(pi.config).UpdateOne(pi)
}

// Unwrap unwraps the PrintingImage entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pi *PrintingImage) Unwrap() *PrintingImage {
	_tx, ok := pi.config.driver.(*txDriver)
	if !ok {
		panic("bones: PrintingImage is not a transactional entity")
	}
	pi.config.driver = _tx.drv
	return pi
}

// String implements the fmt.Stringer.
func (pi *PrintingImage) String() string {
	var builder strings.Builder
	builder.WriteString("PrintingImage(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pi.ID))
	builder.WriteString("url=")
	builder.WriteString(pi.URL)
	builder.WriteString(", ")
	builder.WriteString("image_type=")
	builder.WriteString(fmt.Sprintf("%v", pi.ImageType))
	builder.WriteString(", ")
	builder.WriteString("local_path=")
	builder.WriteString(pi.LocalPath)
	builder.WriteByte(')')
	return builder.String()
}

// PrintingImages is a parsable slice of PrintingImage.
type PrintingImages []*PrintingImage
