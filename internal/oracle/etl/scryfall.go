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
	"github.com/SethCurry/stax/internal/oracle/oracledb/printing"
	"github.com/SethCurry/stax/internal/oracle/oracledb/printingimage"
	"github.com/SethCurry/stax/internal/oracle/oracledb/set"
	"github.com/SethCurry/stax/pkg/scryfall"
	"go.uber.org/zap"
)

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

		err = scryfallCardIngestor(ctx, logger.With(zap.String("card_name", card.Name)), db, card)
		if err != nil {
			return err
		}
	}
	return nil
}

func scryfallCardIngestor(ctx context.Context, logger *zap.Logger, db *oracledb.Client, row *scryfall.Card) error {
	gotCard, err := getOrCreateCard(ctx, logger, db, row)
	if err != nil {
		return fmt.Errorf("failed to get or create card: %w", err)
	}

	cardArtist, err := getOrCreateCardArtist(ctx, logger, db, row.Artist)
	if err != nil {
		return fmt.Errorf("failed to get or create card artist: %w", err)
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

	err = createPrintingImages(ctx, logger, db, row, cardPrinting)
	if err != nil {
		return fmt.Errorf("failed to create printing images: %w", err)
	}

	return nil
}

func createPrintingImages(ctx context.Context, logger *zap.Logger, db *oracledb.Client, row *scryfall.Card, cardPrinting *oracledb.Printing) error {
	var err error

	if row.ImageURIs.Small != "" {
		err = createSinglePrintingImage(ctx, logger, db, row.ImageURIs.Small, printingimage.ImageTypeSmall, cardPrinting)
		if err != nil {
			return err
		}
	}

	if row.ImageURIs.Normal != "" {
		err = createSinglePrintingImage(ctx, logger, db, row.ImageURIs.Normal, printingimage.ImageTypeNormal, cardPrinting)
		if err != nil {
			return err
		}
	}

	if row.ImageURIs.Large != "" {
		err = createSinglePrintingImage(ctx, logger, db, row.ImageURIs.Large, printingimage.ImageTypeLarge, cardPrinting)
		if err != nil {
			return err
		}
	}

	if row.ImageURIs.PNG != "" {
		err = createSinglePrintingImage(ctx, logger, db, row.ImageURIs.PNG, printingimage.ImageTypePng, cardPrinting)
		if err != nil {
			return err
		}
	}

	if row.ImageURIs.ArtCrop != "" {
		err = createSinglePrintingImage(ctx, logger, db, row.ImageURIs.ArtCrop, printingimage.ImageTypeArtCrop, cardPrinting)
		if err != nil {
			return err
		}
	}

	if row.ImageURIs.BorderCrop != "" {
		err = createSinglePrintingImage(ctx, logger, db, row.ImageURIs.BorderCrop, printingimage.ImageTypeBorderCrop, cardPrinting)
		if err != nil {
			return err
		}
	}

	return nil
}

func createSinglePrintingImage(ctx context.Context, logger *zap.Logger, db *oracledb.Client, imageURL string, imageType printingimage.ImageType, cardPrinting *oracledb.Printing) error {
	_, err := db.PrintingImage.Query().Where(printingimage.ImageTypeEQ(imageType), printingimage.URLEQ(imageURL), printingimage.HasPrintingWith(printing.IDEQ(cardPrinting.ID))).Only(ctx)
	if err == nil {
		return nil
	}

	if !oracledb.IsNotFound(err) {
		return fmt.Errorf("failed to query for existing printing image: %w", err)
	}

	_, err = db.PrintingImage.Create().SetURL(imageURL).SetImageType(imageType).SetPrinting(cardPrinting).Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create new printing image: %w", err)
	}

	return nil
}

func getOrCreatePrinting(ctx context.Context, logger *zap.Logger, db *oracledb.Client, rarity printing.Rarity, gotArtist *oracledb.Artist, gotSet *oracledb.Set, gotCardFace *oracledb.CardFace) (*oracledb.Printing, error) {
	existingPrinting, err := db.Printing.Query().Where(printing.RarityEQ(rarity), printing.HasArtistWith(artist.IDEQ(gotArtist.ID)), printing.HasSetWith(set.IDEQ(gotSet.ID)), printing.HasCardFaceWith(cardface.IDEQ(gotCardFace.ID))).Only(ctx)
	if err == nil {
		return existingPrinting, nil
	}

	if !oracledb.IsNotFound(err) {
		return nil, fmt.Errorf("failed to query for existing printing: %w", err)
	}

	newPrinting, err := db.Printing.Create().SetArtist(gotArtist).SetSet(gotSet).SetCardFace(gotCardFace).SetRarity(rarity).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create new printing: %w", err)
	}

	return newPrinting, nil
}

func getOrCreateCardFace(ctx context.Context, logger *zap.Logger, db *oracledb.Client, row *scryfall.Card, gotCard *oracledb.Card) (*oracledb.CardFace, error) {
	existingCardFace, err := db.CardFace.Query().Where(cardface.NameEQ(row.Name), cardface.HasCardWith(card.IDEQ(gotCard.ID))).Only(ctx)
	if err == nil {
		return existingCardFace, nil
	}

	if !oracledb.IsNotFound(err) {
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
		return nil, fmt.Errorf("failed to create new card face: %w", err)
	}

	return newCardFace, nil
}

func getOrCreateCard(ctx context.Context, logger *zap.Logger, db *oracledb.Client, row *scryfall.Card) (*oracledb.Card, error) {
	existingCard, err := db.Card.Query().Where(card.NameEQ(row.Name)).Only(ctx)
	if err == nil {
		return existingCard, nil
	}

	if !oracledb.IsNotFound(err) {
		return nil, fmt.Errorf("failed to query for existing card: %w", err)
	}

	newCard, err := db.Card.Create().SetName(row.Name).SetOracleID(row.OracleID).SetColorIdentity(0).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create new card: %w", err)
	}

	return newCard, nil
}

func getOrCreateSet(ctx context.Context, logger *zap.Logger, db *oracledb.Client, setName string, setCode string) (*oracledb.Set, error) {
	existingSet, err := db.Set.Query().Where(set.NameEQ(setName)).Only(ctx)
	if err == nil {
		return existingSet, nil
	}

	if !oracledb.IsNotFound(err) {
		return nil, fmt.Errorf("failed to query for existing set: %w", err)
	}

	newSet, err := db.Set.Create().SetName(setName).SetCode(setCode).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create new set: %w", err)
	}

	return newSet, nil
}

func getOrCreateCardArtist(ctx context.Context, logger *zap.Logger, db *oracledb.Client, artistName string) (*oracledb.Artist, error) {
	existingArtist, err := db.Artist.Query().Where(artist.NameEQ(artistName)).Only(ctx)
	if err == nil {
		return existingArtist, nil
	}

	if !oracledb.IsNotFound(err) {
		return nil, fmt.Errorf("failed to query for existing artist: %w", err)
	}

	newArtist, err := db.Artist.Create().SetName(artistName).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create new artist: %w", err)
	}

	return newArtist, nil
}
