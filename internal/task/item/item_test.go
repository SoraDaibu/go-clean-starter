package item

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/SoraDaibu/go-clean-starter/domain"
	"github.com/SoraDaibu/go-clean-starter/internal/logger"
)

type importRepository struct {
	domain.ItemRepository
	created  []uint
	failType uint
}

func (r *importRepository) CreateItem(_ context.Context, item *domain.Item) (*domain.Item, error) {
	if item.TypeID() == r.failType {
		return nil, errors.New("storage unavailable")
	}
	r.created = append(r.created, item.TypeID())
	return item, nil
}

type importTransaction struct{}

func (importTransaction) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestImportItemsReportsRowFailures(t *testing.T) {
	for _, tt := range []struct {
		name        string
		dryRun      bool
		first       string
		failType    uint
		wantError   bool
		wantCreated int
	}{
		{"valid", false, "1,one,description\n", 0, false, 2},
		{"invalid row", false, "bad,one,description\n", 0, true, 1},
		{"empty type", false, ",one,description\n", 0, true, 1},
		{"negative type", false, "-1,one,description\n", 0, true, 1},
		{"invalid dry run", true, "bad,one,description\n", 0, true, 0},
		{"storage failure", false, "1,one,description\n", 1, true, 1},
		{"valid dry run", true, "1,one,description\n", 0, false, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for name, rows := range map[string]string{"a.csv": tt.first, "b.csv": "2,two,description\n"} {
				if err := os.WriteFile(filepath.Join(dir, name), []byte("type_id,name,description\n"+rows), 0600); err != nil {
					t.Fatal(err)
				}
			}
			repo := &importRepository{failType: tt.failType}
			u := NewItemTaskUsecase(importTransaction{}, repo, logger.New("error", "test"))
			err := u.ImportItems(context.Background(), dir, tt.dryRun)
			if (err != nil) != tt.wantError {
				t.Fatalf("ImportItems() error = %v, wantError %v", err, tt.wantError)
			}
			if tt.wantError && !strings.Contains(err.Error(), "1 row error") {
				t.Fatalf("missing error count: %v", err)
			}
			if len(repo.created) != tt.wantCreated {
				t.Fatalf("created %v, want %d successful rows", repo.created, tt.wantCreated)
			}
		})
	}
}
