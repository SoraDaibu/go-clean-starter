package item

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/SoraDaibu/go-clean-starter/domain"
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
	u.Logger.Info("Starting item import", "source_dir", sourceDir, "dry_run", dryRun)

	// Read all CSV files in the source directory
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		u.Logger.Error("Failed to read source directory", "err", err)
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
			u.Logger.Error("Failed to import CSV file", "err", err, "file", filePath)
			return fmt.Errorf("failed to import file %s: %w", filePath, err)
		}

		result.FilePath = filePath
		totalResults = append(totalResults, result)

		u.Logger.Info("Import completed",
			"file", filePath,
			"created", result.ItemsCreated,
			"skipped", result.ItemsSkipped,
			"errors", len(result.Errors),
		)
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

	u.Logger.Info("Import summary",
		"files_processed", len(totalResults),
		"total_created", totalCreated,
		"total_skipped", totalSkipped,
		"total_errors", totalErrors,
	)

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
			u.Logger.Error("invalid CSV format", "err", err)
			result.addError(err)
			continue
		}

		// Parse typeID if provided
		var typeID uint
		if record[0] == "" {
			err := fmt.Errorf("empty type_id for item %s at line %d", record[1], i+2)
			u.Logger.Error("empty type_id", "err", err)
			result.addError(err)
			continue
		}

		typeIDInt, err := strconv.Atoi(record[0])
		if err != nil {
			err := fmt.Errorf("invalid type_id '%s' for item %s at line %d: %w", record[0], record[1], i+2, err)
			u.Logger.Error("failed to convert type_id to int", "err", err)
			result.addError(err)
			continue
		}

		if typeIDInt < 0 {
			err := fmt.Errorf("negative type_id '%d' for item %s at line %d", typeIDInt, record[1], i+2)
			u.Logger.Error("invalid type_id", "err", err)
			result.addError(err)
			continue
		}

		typeID = uint(typeIDInt)

		// Create domain item
		item := domain.NewItem(typeID)

		if dryRun {
			u.Logger.Info("DRY RUN: Would create item", "id", item.ID().String(), "type_id", item.TypeID())
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
		u.Logger.Debug("Item created successfully", "id", item.ID().String(), "type_id", item.TypeID())
	}

	return result, nil
}
