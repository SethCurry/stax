// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// ArtistsColumns holds the columns for the "artists" table.
	ArtistsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString, Size: 255},
	}
	// ArtistsTable holds the schema information for the "artists" table.
	ArtistsTable = &schema.Table{
		Name:       "artists",
		Columns:    ArtistsColumns,
		PrimaryKey: []*schema.Column{ArtistsColumns[0]},
	}
	// CardsColumns holds the columns for the "cards" table.
	CardsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString, Size: 255},
		{Name: "oracle_id", Type: field.TypeString},
		{Name: "color_identity", Type: field.TypeUint8},
	}
	// CardsTable holds the schema information for the "cards" table.
	CardsTable = &schema.Table{
		Name:       "cards",
		Columns:    CardsColumns,
		PrimaryKey: []*schema.Column{CardsColumns[0]},
	}
	// CardFacesColumns holds the columns for the "card_faces" table.
	CardFacesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString, Size: 255},
		{Name: "flavor_text", Type: field.TypeString, Size: 1000},
		{Name: "oracle_text", Type: field.TypeString, Size: 1000},
		{Name: "language", Type: field.TypeString},
		{Name: "cmc", Type: field.TypeFloat32},
		{Name: "power", Type: field.TypeString},
		{Name: "toughness", Type: field.TypeString},
		{Name: "loyalty", Type: field.TypeString},
		{Name: "mana_cost", Type: field.TypeString},
		{Name: "type_line", Type: field.TypeString},
		{Name: "colors", Type: field.TypeString},
		{Name: "card_face_card", Type: field.TypeInt, Nullable: true},
	}
	// CardFacesTable holds the schema information for the "card_faces" table.
	CardFacesTable = &schema.Table{
		Name:       "card_faces",
		Columns:    CardFacesColumns,
		PrimaryKey: []*schema.Column{CardFacesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "card_faces_cards_card",
				Columns:    []*schema.Column{CardFacesColumns[12]},
				RefColumns: []*schema.Column{CardsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// PrintingsColumns holds the columns for the "printings" table.
	PrintingsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "rarity", Type: field.TypeEnum, Enums: []string{"common", "uncommon", "rare", "mythic", "special", "bonus"}},
		{Name: "printing_artist", Type: field.TypeInt, Nullable: true},
		{Name: "printing_set", Type: field.TypeInt, Nullable: true},
		{Name: "printing_card_face", Type: field.TypeInt, Nullable: true},
	}
	// PrintingsTable holds the schema information for the "printings" table.
	PrintingsTable = &schema.Table{
		Name:       "printings",
		Columns:    PrintingsColumns,
		PrimaryKey: []*schema.Column{PrintingsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "printings_artists_artist",
				Columns:    []*schema.Column{PrintingsColumns[2]},
				RefColumns: []*schema.Column{ArtistsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "printings_sets_set",
				Columns:    []*schema.Column{PrintingsColumns[3]},
				RefColumns: []*schema.Column{SetsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "printings_card_faces_card_face",
				Columns:    []*schema.Column{PrintingsColumns[4]},
				RefColumns: []*schema.Column{CardFacesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// PrintingImagesColumns holds the columns for the "printing_images" table.
	PrintingImagesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "url", Type: field.TypeString},
		{Name: "image_type", Type: field.TypeEnum, Enums: []string{"small", "normal", "large", "png", "art_crop", "border_crop"}},
		{Name: "printing_image_printing", Type: field.TypeInt, Nullable: true},
	}
	// PrintingImagesTable holds the schema information for the "printing_images" table.
	PrintingImagesTable = &schema.Table{
		Name:       "printing_images",
		Columns:    PrintingImagesColumns,
		PrimaryKey: []*schema.Column{PrintingImagesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "printing_images_printings_printing",
				Columns:    []*schema.Column{PrintingImagesColumns[3]},
				RefColumns: []*schema.Column{PrintingsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// RulingsColumns holds the columns for the "rulings" table.
	RulingsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "text", Type: field.TypeString},
		{Name: "date", Type: field.TypeTime},
		{Name: "ruling_card", Type: field.TypeInt, Nullable: true},
	}
	// RulingsTable holds the schema information for the "rulings" table.
	RulingsTable = &schema.Table{
		Name:       "rulings",
		Columns:    RulingsColumns,
		PrimaryKey: []*schema.Column{RulingsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "rulings_cards_card",
				Columns:    []*schema.Column{RulingsColumns[3]},
				RefColumns: []*schema.Column{CardsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// SetsColumns holds the columns for the "sets" table.
	SetsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString, Size: 255},
		{Name: "code", Type: field.TypeString, Size: 255},
	}
	// SetsTable holds the schema information for the "sets" table.
	SetsTable = &schema.Table{
		Name:       "sets",
		Columns:    SetsColumns,
		PrimaryKey: []*schema.Column{SetsColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		ArtistsTable,
		CardsTable,
		CardFacesTable,
		PrintingsTable,
		PrintingImagesTable,
		RulingsTable,
		SetsTable,
	}
)

func init() {
	CardFacesTable.ForeignKeys[0].RefTable = CardsTable
	PrintingsTable.ForeignKeys[0].RefTable = ArtistsTable
	PrintingsTable.ForeignKeys[1].RefTable = SetsTable
	PrintingsTable.ForeignKeys[2].RefTable = CardFacesTable
	PrintingImagesTable.ForeignKeys[0].RefTable = PrintingsTable
	RulingsTable.ForeignKeys[0].RefTable = CardsTable
}
