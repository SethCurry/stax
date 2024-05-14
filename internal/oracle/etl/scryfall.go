package etl

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/SethCurry/stax/internal/oracle/oracledb"
	"github.com/SethCurry/stax/internal/oracle/oracledb/artist"
	"github.com/SethCurry/stax/internal/oracle/oracledb/card"
	"github.com/SethCurry/stax/internal/oracle/oracledb/cardface"
	"github.com/SethCurry/stax/internal/oracle/oracledb/predicate"
	"github.com/SethCurry/stax/internal/oracle/oracledb/printing"
	"github.com/SethCurry/stax/internal/oracle/oracledb/printingimage"
	"github.com/SethCurry/stax/internal/oracle/oracledb/set"
	"github.com/SethCurry/stax/pkg/scryfall"
	"go.uber.org/zap"
)

// ScryfallCards reads all of the cards from the provided scryfall.BulkReader and creates
// all of the SQL records implied by that object (set, artists, etc).
func ScryfallCards(ctx context.Context, logger *zap.Logger, db *oracledb.Client, reader *scryfall.BulkReader[scryfall.Card]) error {
	for {
		card, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				logger.Debug("got EOF")
				return nil
			}

			return err
		}

		if card == nil {
			logger.Debug("card is nil, breaking")
			break
		}

		err = scryfallCardIngestor(
			ctx,
			logger.With(zap.String("card_name", card.Name)),
			db,
			card)
		if err != nil {
			return err
		}
	}
	return nil
}

func scryfallCardIngestor(
	ctx context.Context,
	logger *zap.Logger,
	db *oracledb.Client,
	row *scryfall.Card,
) error {
	if row.OracleID == "" {
		logger.Debug("ignoring card because it is missing an oracle ID")
		return nil
	}

	gotCard, err := getOrCreateCard(ctx, logger, db, row)
	if err != nil {
		return fmt.Errorf("failed to get or create card: %w", err)
	}

	var cardArtist *oracledb.Artist

	if row.Artist != "" {
		cardArtist, err = getOrCreateCardArtist(ctx, logger, db, row.Artist)
		if err != nil {
			return fmt.Errorf("failed to get or create card artist: %w", err)
		}
	}

	cardSet, err := getOrCreateSet(ctx, logger, db, row.SetName, row.SetCode)
	if err != nil {
		return fmt.Errorf("failed to get or create card set: %w", err)
	}

	cardFace, err := getOrCreateCardFace(ctx, logger, db, row, gotCard)
	if err != nil {
		return fmt.Errorf("failed to get or create card face: %w", err)
	}

	cardPrinting, err := getOrCreatePrinting(ctx, logger, db, printing.Rarity(row.Rarity), cardArtist, cardSet, cardFace)
	if err != nil {
		return fmt.Errorf("failed to get or create card printing: %w", err)
	}

	err = createPrintingImagesIfNotExist(ctx, logger, db, row, cardPrinting)
	if err != nil {
		return fmt.Errorf("failed to create printing images: %w", err)
	}

	return nil
}

// createPrintingImagesIfNotExist creates all printing images for a card if they do not already exist.
// It will ignore any images that are not present in the row.
func createPrintingImagesIfNotExist(
	ctx context.Context,
	logger *zap.Logger,
	db *oracledb.Client,
	row *scryfall.Card,
	cardPrinting *oracledb.Printing,
) error {
	imageURIs := []struct {
		uri   string
		type_ printingimage.ImageType
	}{
		{row.ImageURIs.Small, printingimage.ImageTypeSmall},
		{row.ImageURIs.Normal, printingimage.ImageTypeNormal},
		{row.ImageURIs.Large, printingimage.ImageTypeLarge},
		{row.ImageURIs.PNG, printingimage.ImageTypePng},
		{row.ImageURIs.ArtCrop, printingimage.ImageTypeArtCrop},
		{row.ImageURIs.BorderCrop, printingimage.ImageTypeBorderCrop},
	}

	for _, imageURI := range imageURIs {
		if imageURI.uri != "" {
			_ = createSinglePrintingImage(
				ctx,
				logger,
				db,
				imageURI.uri,
				imageURI.type_,
				cardPrinting)
		}
	}

	return nil
}

// createSinglePrintingImage creates a single printing image in the database
// if there is not an existing record with the same image type, URL, and printing.
func createSinglePrintingImage(
	ctx context.Context,
	logger *zap.Logger,
	db *oracledb.Client,
	imageURL string,
	imageType printingimage.ImageType,
	cardPrinting *oracledb.Printing,
) error {
	logger = logger.With(
		zap.String("image_type", string(imageType)),
		zap.String("image_url", imageURL))

	_, err := db.PrintingImage.Query().
		Where(
			printingimage.ImageTypeEQ(imageType),
			printingimage.URLEQ(imageURL),
			printingimage.HasPrintingWith(printing.IDEQ(cardPrinting.ID))).
		Only(ctx)
	if err == nil {
		logger.Debug("printing image already exists")
		return nil
	}

	if !oracledb.IsNotFound(err) {
		logger.Error("failed to query for existing printing image", zap.Error(err))
		return fmt.Errorf("failed to query for existing printing image: %w", err)
	}

	_, err = db.PrintingImage.Create().
		SetURL(imageURL).
		SetImageType(imageType).
		SetPrinting(cardPrinting).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create new printing image: %w", err)
	}

	return nil
}

// getOrCreatePrinting checks if a printing matching the provided parameters already exists.
// If a matching printing is found, the matching printing is returned.  If no matching
// printing could be found, one will be created and returned.
func getOrCreatePrinting(
	ctx context.Context,
	logger *zap.Logger, // Zap logger for logging purposes
	db *oracledb.Client, // OracleDB client to interact with the database
	rarity printing.Rarity, // The rarity of the card face
	gotArtist *oracledb.Artist, // Pointer to an artist entity (optional)
	gotSet *oracledb.Set, // Set associated with the card face
	gotCardFace *oracledb.CardFace, // The card face we are dealing with
) (*oracledb.Printing, error) {
	// Logger is updated with additional contextual information about the printing
	logger = logger.With(
		zap.String("rarity", string(rarity)),
		zap.String("set_name", gotSet.Name),
		zap.String("card_face_name", gotCardFace.Name),
		zap.Bool("has_artist", gotArtist != nil), // Checks if artist is present or not
	)

	// Predicate for the artist based on whether it exists or not
	var artistPred predicate.Printing
	if gotArtist != nil {
		artistPred = printing.HasArtistWith(artist.IDEQ(gotArtist.ID))
	} else {
		artistPred = printing.Not(printing.HasArtist())
	}

	// Query to find an existing printing with the given parameters
	existingPrinting, err := db.Printing.Query().Where(
		printing.RarityEQ(rarity),                // Rarity of the card face
		artistPred,                               // Artist predicate based on whether an artist is present or not
		printing.HasSetWith(set.IDEQ(gotSet.ID)), // Set association of the printing
		printing.HasCardFaceWith(cardface.IDEQ(gotCardFace.ID))).Only(ctx)
	if err == nil {
		// If existingPrinting is found, a debug log is made and existing printing is returned
		logger.Debug("printing already exists")
		return existingPrinting, nil
	}

	// Error handling if the printing does not exist in the database
	if !oracledb.IsNotFound(err) {
		logger.Error("failed to query for existing printing", zap.Error(err))
		return nil, fmt.Errorf("failed to query for existing printing: %w", err)
	}

	// If no previous error occurred and the printing does not exist in the database, a new printing is created with given parameters
	newPrintingQuery := db.Printing.Create().SetSet(gotSet).SetCardFace(gotCardFace).SetRarity(rarity)
	if gotArtist != nil { // If an artist is present, it is set in the new printing query
		newPrintingQuery = newPrintingQuery.SetArtist(gotArtist)
	}

	// The newly created printing is saved into the database and returned along with any errors that might have occurred during saving
	newPrinting, err := newPrintingQuery.Save(ctx)
	if err != nil {
		logger.Error("failed to create new printing", zap.Error(err))
		return nil, fmt.Errorf("failed to create new printing: %w", err)
	}

	logger.Info("created new printing")

	return newPrinting, nil // Return of the newly created printing and no error
}

func getOrCreateCardFace(
	ctx context.Context,
	logger *zap.Logger,
	db *oracledb.Client,
	row *scryfall.Card,
	gotCard *oracledb.Card,
) (*oracledb.CardFace, error) {
	logger = logger.With(zap.String("card_face_name", row.Name))

	existingCardFace, err := db.CardFace.Query().Where(
		cardface.NameEQ(row.Name),
		cardface.HasCardWith(card.IDEQ(gotCard.ID)),
	).
		Only(ctx)
	if err == nil {
		logger.Debug("card face already exists")
		return existingCardFace, nil
	}

	if !oracledb.IsNotFound(err) {
		logger.Error("failed to query for existing card face", zap.Error(err))
		return nil, fmt.Errorf("failed to query for existing card face: %w", err)
	}

	newCardFace, err := db.CardFace.Create().
		SetName(row.Name).
		SetFlavorText(row.FlavorText).
		SetOracleText(row.OracleText).
		SetLanguage(row.Language).
		SetCmc(row.CMC).
		SetPower(row.Power).
		SetToughness(row.Toughness).
		SetLoyalty(row.Loyalty).
		SetManaCost(row.ManaCost).
		SetTypeLine(row.TypeLine).
		SetColors(strings.Join(row.Colors, "")).
		SetCard(gotCard).
		Save(ctx)
	if err != nil {
		logger.Error("failed to create new card face", zap.Error(err))
		return nil, fmt.Errorf("failed to create new card face: %w", err)
	}

	logger.Info("created new card face")

	return newCardFace, nil
}

func getOrCreateCard(
	ctx context.Context,
	logger *zap.Logger,
	db *oracledb.Client,
	row *scryfall.Card,
) (*oracledb.Card, error) {
	logger = logger.With(
		zap.String("card_name", row.Name),
		zap.String("oracle_id", row.OracleID),
		zap.Strings("colors", row.Colors),
	)

	existingCard, err := db.Card.Query().Where(card.NameEQ(row.Name)).Only(ctx)
	if err == nil {
		logger.Debug("card already exists")
		return existingCard, nil
	}

	if !oracledb.IsNotFound(err) {
		logger.Error("failed to query for existing card", zap.Error(err))
		return nil, fmt.Errorf("failed to query for existing card: %w", err)
	}

	// TODO write function to convert row.Colors to a bitfield
	newCard, err := db.Card.Create().
		SetName(row.Name).
		SetOracleID(row.OracleID).
		SetColorIdentity(0).
		Save(ctx)
	if err != nil {
		logger.Error("failed to create new card", zap.Error(err))
		return nil, fmt.Errorf("failed to create new card: %w", err)
	}

	logger.Info("created new card")

	return newCard, nil
}

func getOrCreateSet(
	ctx context.Context,
	logger *zap.Logger,
	db *oracledb.Client,
	setName string,
	setCode string,
) (*oracledb.Set, error) {
	logger = logger.With(zap.String("set_name", setName), zap.String("set_code", setCode))

	existingSet, err := db.Set.Query().Where(set.NameEQ(setName)).Only(ctx)
	if err == nil {
		logger.Debug("card already exists")
		return existingSet, nil
	}

	if !oracledb.IsNotFound(err) {
		logger.Error("failed to query for existing set", zap.Error(err))
		return nil, fmt.Errorf("failed to query for existing set: %w", err)
	}

	newSet, err := db.Set.Create().SetName(setName).SetCode(setCode).Save(ctx)
	if err != nil {
		logger.Error("failed to create new set", zap.Error(err))
		return nil, fmt.Errorf("failed to create new set: %w", err)
	}

	logger.Info("created new set")

	return newSet, nil
}

// getOrCreateCardArtist checks whether a card artist exists in the database and creates
// it if not.
func getOrCreateCardArtist(
	ctx context.Context,
	logger *zap.Logger,
	db *oracledb.Client,
	artistName string,
) (*oracledb.Artist, error) {
	logger = logger.With(zap.String("artist_name", artistName))

	existingArtist, err := db.Artist.Query().Where(artist.NameEQ(artistName)).Only(ctx)
	if err == nil {
		logger.Debug("artist already exists")
		return existingArtist, nil
	}

	if !oracledb.IsNotFound(err) {
		logger.Error("failed to query for existing artist", zap.Error(err))
		return nil, fmt.Errorf("failed to query for existing artist: %w", err)
	}

	newArtist, err := db.Artist.Create().SetName(artistName).Save(ctx)
	if err != nil {
		logger.Error("failed to create new artist", zap.Error(err))
		return nil, fmt.Errorf("failed to create new artist: %w", err)
	}

	logger.Info("created new artist")

	return newArtist, nil
}