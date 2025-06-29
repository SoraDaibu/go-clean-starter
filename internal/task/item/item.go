package item

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/SoraDaibu/go-clean-starter/domain"
	"github.com/rs/zerolog/log"
)

type ImportResult struct {
	FilePath     string
	ItemsCreated int
	ItemsSkipped int
	Errors       []string
}

func (r *ImportResult) addError(err error) {
	r.Errors = append(r.Errors, err.Error())
}

func (u *itemTaskUsecase) ImportItems(ctx context.Context, sourceDir string, dryRun bool) error {
	log.Info().Str("source_dir", sourceDir).Bool("dry_run", dryRun).Msg("Starting item import")

	// Read all CSV files in the source directory
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read source directory")
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	totalResults := []*ImportResult{}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".csv" {
			continue
		}

		filePath := filepath.Join(sourceDir, file.Name())
		result, err := u.importCSVFile(ctx, filePath, dryRun)
		if err != nil {
			log.Error().Err(err).Str("file", filePath).Msg("Failed to import CSV file")
			return fmt.Errorf("failed to import file %s: %w", filePath, err)
		}

		result.FilePath = filePath
		totalResults = append(totalResults, result)

		log.Info().
			Str("file", filePath).
			Int("created", result.ItemsCreated).
			Int("skipped", result.ItemsSkipped).
			Int("errors", len(result.Errors)).
			Msg("Import completed")
	}

	// Summary log
	totalCreated := 0
	totalSkipped := 0
	totalErrors := 0
	for _, result := range totalResults {
		totalCreated += result.ItemsCreated
		totalSkipped += result.ItemsSkipped
		totalErrors += len(result.Errors)
	}

	log.Info().
		Int("files_processed", len(totalResults)).
		Int("total_created", totalCreated).
		Int("total_skipped", totalSkipped).
		Int("total_errors", totalErrors).
		Msg("Import summary")

	return nil
}

func (u *itemTaskUsecase) importCSVFile(ctx context.Context, filePath string, dryRun bool) (*ImportResult, error) {
	result := &ImportResult{}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return result, nil
	}

	// Skip header row (1)
	for i, record := range records[1:] {
		if len(record) < 3 {
			err := fmt.Errorf("invalid CSV format at line %d: expected 3 columns (type_id,name,description), got %d", i+2, len(record))
			log.Error().Err(err).Msg("invalid CSV format")
			result.addError(err)
			continue
		}

		// Parse typeID if provided
		var typeID uint
		if record[0] == "" {
			err := fmt.Errorf("empty type_id for item %s at line %d", record[1], i+2)
			log.Error().Err(err).Msg("empty type_id")
			result.addError(err)
			continue
		}

		typeIDInt, err := strconv.Atoi(record[0])
		if err != nil {
			err := fmt.Errorf("invalid type_id '%s' for item %s at line %d: %w", record[0], record[1], i+2, err)
			log.Error().Err(err).Msg("failed to convert type_id to int")
			result.addError(err)
			continue
		}

		if typeIDInt < 0 {
			err := fmt.Errorf("negative type_id '%d' for item %s at line %d", typeIDInt, record[1], i+2)
			log.Error().Err(err).Msg("invalid type_id")
			result.addError(err)
			continue
		}

		typeID = uint(typeIDInt)

		// Create domain item
		item := domain.NewItem(typeID)

		if dryRun {
			log.Info().
				Str("id", item.ID().String()).
				Interface("type_id", item.TypeID()).
				Msg("DRY RUN: Would create item")
			result.ItemsCreated++
			continue
		}

		// Create item in database with transaction
		err = u.Tx.Do(ctx, func(ctx context.Context) error {
			_, err := u.ItemRepo.CreateItem(ctx, item)
			return err
		})

		if err != nil {
			err := fmt.Errorf("failed to create item at line %d: %w", i+2, err)
			result.addError(err)
			continue
		}

		result.ItemsCreated++
		log.Debug().
			Str("id", item.ID().String()).
			Interface("type_id", item.TypeID()).
			Msg("Item created successfully")
	}

	return result, nil
}
