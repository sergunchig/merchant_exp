package offer

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sergunchig/merchant_exp.git/internal/entity"
	mock "github.com/stretchr/testify/mock"
)

type excelReaderMock struct {
	mock.Mock
}

func (er excelReaderMock) Read(file string) ([]entity.Offer, error) {

	offers := make([]entity.Offer, 0, 0)
	return offers, nil
}

type repoOffersMock struct {
	mock.Mock
}

func (r repoOffersMock) CreateOffers(ctx context.Context, offers []entity.Offer) error {
	return nil
}

type appLoggerMock struct {
	mock.Mock
}

func (l appLoggerMock) Error(msg string) {
	fmt.Println(msg)
}

func TestUploadAndImportHandler(t *testing.T) {
	ctx := context.Background()
	offers := make([]entity.Offer, 0, 0)

	reader := new(excelReaderMock)
	reader.On("Read", "file").Return(offers, nil)

	repo := new(repoOffersMock)
	repo.On("CreateOffers", ctx, offers).Return(nil)

	log := new(appLoggerMock)

	handler := New(repo, reader, log)

	req, _ := http.NewRequest(http.MethodPost, "/UploadAndImport", nil)
	rr := httptest.NewRecorder()

	handler.UploadAndImportHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("staus %d", status)
	}
}
