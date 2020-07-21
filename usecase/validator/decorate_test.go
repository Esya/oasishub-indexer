package validator_test

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"testing"

	"github.com/figment-networks/oasishub-indexer/config"
	mock "github.com/figment-networks/oasishub-indexer/mock/usecase/validator"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/usecase/validator"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

var (
	// Files keep track of files to cleanup
	Files       []string
	errDbCreate = errors.New("errDbCreate")
	errDbFind   = errors.New("errDbFind")
)

func TestDecorate_Handle(t *testing.T) {
	defer cleanUp(t)

	tests := []struct {
		description string
		fileName    string
		data        [][]string
		dbCreateErr error
		dbFindErr   error
		expectErr   error
	}{
		{"should error if no file provided",
			"",
			[][]string{},
			nil,
			nil,
			validator.ErrMissingFile},
		{"should error if file is missing headers",
			"no_headers.csv",
			[][]string{row("test1"), row("test2")},
			nil,
			nil,
			validator.ErrInValidFile},
		{"should update valid csv",
			"case_1.csv",
			[][]string{headers(), row("test1"), row("test2")},
			nil,
			nil,
			nil},
		{"should error if db errors on create method",
			"case_2.csv",
			[][]string{headers(), row("test1"), row("test2")},
			errDbCreate,
			nil,
			errDbCreate},
		{"should error if unexpected db error on find method",
			"case_3.csv",
			[][]string{headers(), row("test1"), row("test2")},
			nil,
			errDbFind,
			errDbFind},
		{"should not update record if it doesn't already exist in db ",
			"case_4.csv",
			[][]string{headers(), row("test1"), row("test2")},
			nil,
			store.ErrNotFound,
			nil},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cfg := &config.Config{}
			ctrl := gomock.NewController(t)

			if tt.fileName != "" {
				createFile(tt.fileName, tt.data, t)
			}

			dbMock := mock.NewMockValidatorAggStore(ctrl)
			for i, row := range tt.data {
				if i == 0 {
					// skip first row of headers
					continue
				}
				val := &model.ValidatorAgg{RecentAddress: row[1]}

				dbMock.EXPECT().FindBy("recent_address", row[1]).Return(val, tt.dbFindErr).Times(1)
				if tt.dbFindErr == store.ErrNotFound {
					// don't expect CreateOrUpdate to be called for this val
					continue
				} else if tt.dbFindErr != nil {
					// don't expect rest of values in file to be called
					break
				}

				val.LogoURL = row[2]
				val.EntityName = row[0]
				dbMock.EXPECT().CreateOrUpdate(val).Return(tt.dbCreateErr).Times(1)
				if tt.dbCreateErr != nil {
					// don't expect rest of values in file to be called
					break
				}
			}

			uc := validator.NewDecorateUseCase(cfg, dbMock)

			ctx := context.Background()
			err := uc.Execute(ctx, tt.fileName)
			if err != tt.expectErr {
				t.Errorf("unexpected response, want: %+v, got: %+v", tt.expectErr, err)
			}
		})
	}

}

func headers() []string {
	return []string{"amber entities", "entity id (new format)", "logo link"}
}

func row(prefix string) []string {
	return []string{
		fmt.Sprintf("%s-name", prefix),
		fmt.Sprintf("%s-addr", prefix),
		fmt.Sprintf("%s-logo", prefix),
	}
}

func createFile(fileName string, data [][]string, t *testing.T) {
	f, e := os.Create(fileName)
	if e != nil {
		t.Fatal(e)
	}

	writer := csv.NewWriter(f)

	e = writer.WriteAll(data)
	if e != nil {
		t.Fatal(e)
	}

	Files = append(Files, fileName)
}

func cleanUp(t *testing.T) {
	for _, path := range Files {
		if err := os.RemoveAll(path); err != nil {
			t.Error("could not remove file", err)
		}
	}
}
