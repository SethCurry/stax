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
	txn, err := db.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to create initial transaction")
	}

	artistCache := make(map[string]int)
	setCache := make(map[string]int)
	cardCache := make(map[string]int)
	cardFaceCache := make(map[string][]cachedCardFace)

	numCards, err := txn.Card.Query().Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to query for current card count: %w", err)
	}

	isFresh := numCards == 0

	for {
		card, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				err = txn.Commit()
				if err != nil {
					logger.Error("failed to perform last commit", zap.Error(err))
				}
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
			txn,
			card,
			artistCache,
			setCache,
			cardCache,
			cardFaceCache,
			isFresh)
		if err != nil {
			return err
		}
	}
	return txn.Commit()
}

type cachedCardFace struct {
	Name string
	ID   int
}

func scryfallCardIngestor(
	ctx context.Context,
	logger *zap.Logger,
	db *oracledb.Tx,
	row *scryfall.Card,
	artistCache map[string]int,
	setCache map[string]int,
	cardCache map[string]int,
	cardFacesCache map[string][]cachedCardFace,
	isFresh bool,
) error {
	if row.OracleID == "" {
		logger.Debug("ignoring card because it is missing an oracle ID")
		return nil
	}

	cardID := 0

	if gotCardID, ok := cardCache[row.OracleID]; ok {
		cardID = gotCardID
	} else {
		gotCard, err := getOrCreateCard(ctx, logger, db, row, isFresh)
		if err != nil {
			return fmt.Errorf("failed to get or create card: %w", err)
		}

		cardCache[row.OracleID] = gotCard.ID
		cardID = gotCard.ID
	}

	artistID := 0

	if row.Artist != "" {
		if foundArtistID, ok := artistCache[row.Artist]; ok {
			artistID = foundArtistID
		} else {
			cardArtist, err := getOrCreateCardArtist(ctx, logger, db, row.Artist, isFresh)
			if err != nil {
				return fmt.Errorf("failed to get or create card artist: %w", err)
			}
			artistCache[row.Artist] = cardArtist.ID
			artistID = cardArtist.ID
		}
	}

	setID := 0

	if gotSetID, ok := setCache[row.SetCode]; ok {
		setID = gotSetID
	} else {
		cardSet, err := getOrCreateSet(ctx, logger, db, row.SetName, row.SetCode, isFresh)
		if err != nil {
			return fmt.Errorf("failed to get or create card set: %w", err)
		}

		setCache[row.SetCode] = cardSet.ID
		setID = cardSet.ID
	}

	cardFaceID := 0
	foundCardFace := false

	facesForID, ok := cardFacesCache[row.OracleID]
	if ok {
		for _, v := range facesForID {
			if v.Name == row.Name {
				foundCardFace = true
				cardFaceID = v.ID
			}
		}
	}

	if !foundCardFace {
		cardFace, err := getOrCreateCardFace(ctx, logger, db, row, cardID, isFresh)
		if err != nil {
			return fmt.Errorf("failed to get or create card face: %w", err)
		}

		cardFaceID = cardFace.ID

		cardFacesCache[row.OracleID] = append(facesForID, cachedCardFace{
			Name: cardFace.Name,
			ID:   cardFace.ID,
		})
	}

	cardPrinting, err := getOrCreatePrinting(ctx, logger, db, printing.Rarity(row.Rarity), artistID, setID, cardFaceID, isFresh)
	if err != nil {
		return fmt.Errorf("failed to get or create card printing: %w", err)
	}

	err = createPrintingImagesIfNotExist(ctx, logger, db, row, cardPrinting, isFresh)
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
	db *oracledb.Tx,
	row *scryfall.Card,
	cardPrinting *oracledb.Printing,
	isFresh bool,
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
				cardPrinting,
				isFresh,
			)
		}
	}

	return nil
}

// createSinglePrintingImage creates a single printing image in the database
// if there is not an existing record with the same image type, URL, and printing.
func createSinglePrintingImage(
	ctx context.Context,
	logger *zap.Logger,
	db *oracledb.Tx,
	imageURL string,
	imageType printingimage.ImageType,
	cardPrinting *oracledb.Printing,
	isFresh bool,
) error {
	logger = logger.With(
		zap.String("image_type", string(imageType)),
		zap.String("image_url", imageURL))

	if !isFresh {
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
	}

	_, err := db.PrintingImage.Create().
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
	db *oracledb.Tx, // OracleDB client to interact with the database
	rarity printing.Rarity, // The rarity of the card face
	gotArtistID int, // Pointer to an artist entity (optional)
	gotSetID int, // Set associated with the card face
	gotCardFace int, // The card face we are dealing with
	isFresh bool,
) (*oracledb.Printing, error) {
	// Logger is updated with additional contextual information about the printing
	logger = logger.With(
		zap.String("rarity", string(rarity)),
		zap.Int("set_id", gotSetID),
		zap.Bool("has_artist", gotArtistID != 0), // Checks if artist is present or not
	)

	// Predicate for the artist based on whether it exists or not
	var artistPred predicate.Printing
	if gotArtistID != 0 {
		artistPred = printing.HasArtistWith(artist.IDEQ(gotArtistID))
	} else {
		artistPred = printing.Not(printing.HasArtist())
	}

	if !isFresh {
		// Query to find an existing printing with the given parameters
		existingPrinting, err := db.Printing.Query().Where(
			printing.RarityEQ(rarity),               // Rarity of the card face
			artistPred,                              // Artist predicate based on whether an artist is present or not
			printing.HasSetWith(set.IDEQ(gotSetID)), // Set association of the printing
			printing.HasCardFaceWith(cardface.IDEQ(gotCardFace))).Only(ctx)
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
	}

	// If no previous error occurred and the printing does not exist in the database, a new printing is created with given parameters
	newPrintingQuery := db.Printing.Create().SetSetID(gotSetID).SetCardFaceID(gotCardFace).SetRarity(rarity)
	if gotArtistID != 0 { // If an artist is present, it is set in the new printing query
		newPrintingQuery = newPrintingQuery.SetArtistID(gotArtistID)
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
	db *oracledb.Tx,
	row *scryfall.Card,
	gotCardID int,
	isFresh bool,
) (*oracledb.CardFace, error) {
	logger = logger.With(zap.String("card_face_name", row.Name))

	if !isFresh {
		existingCardFace, err := db.CardFace.Query().Where(
			cardface.NameEQ(row.Name),
			cardface.HasCardWith(card.IDEQ(gotCardID)),
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
		SetCardID(gotCardID).
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
	db *oracledb.Tx,
	row *scryfall.Card,
	isFresh bool,
) (*oracledb.Card, error) {
	logger = logger.With(
		zap.String("card_name", row.Name),
		zap.String("oracle_id", row.OracleID),
		zap.Strings("colors", row.Colors),
	)

	if !isFresh {
		existingCard, err := db.Card.Query().Where(card.NameEQ(row.Name)).Only(ctx)
		if err == nil {
			logger.Debug("card already exists")
			return existingCard, nil
		}

		if !oracledb.IsNotFound(err) {
			logger.Error("failed to query for existing card", zap.Error(err))
			return nil, fmt.Errorf("failed to query for existing card: %w", err)
		}
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
	db *oracledb.Tx,
	setName string,
	setCode string,
	isFresh bool,
) (*oracledb.Set, error) {
	logger = logger.With(zap.String("set_name", setName), zap.String("set_code", setCode))

	if !isFresh {
		existingSet, err := db.Set.Query().Where(set.NameEQ(setName)).Only(ctx)
		if err == nil {
			logger.Debug("card already exists")
			return existingSet, nil
		}

		if !oracledb.IsNotFound(err) {
			logger.Error("failed to query for existing set", zap.Error(err))
			return nil, fmt.Errorf("failed to query for existing set: %w", err)
		}
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
	db *oracledb.Tx,
	artistName string,
	isFresh bool,
) (*oracledb.Artist, error) {
	logger = logger.With(zap.String("artist_name", artistName))

	if !isFresh {
		existingArtist, err := db.Artist.Query().Where(artist.NameEQ(artistName)).Only(ctx)
		if err == nil {
			logger.Debug("artist already exists")
			return existingArtist, nil
		}

		if !oracledb.IsNotFound(err) {
			logger.Error("failed to query for existing artist", zap.Error(err))
			return nil, fmt.Errorf("failed to query for existing artist: %w", err)
		}
	}

	newArtist, err := db.Artist.Create().SetName(artistName).Save(ctx)
	if err != nil {
		logger.Error("failed to create new artist", zap.Error(err))
		return nil, fmt.Errorf("failed to create new artist: %w", err)
	}

	logger.Info("created new artist")

	return newArtist, nil
}
