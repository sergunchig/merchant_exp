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
func (r repoOffersMock) GetOffers(ctx context.Context) ([]entity.Offer, error) {
	args := r.Called(ctx)
	return args.Get(0).([]entity.Offer), args.Error(1)
}

func generateOffers() []entity.Offer {
	return []entity.Offer{
		entity.Offer{OfferId: 1, Name: "cat", Price: 10, Available: true},
		entity.Offer{OfferId: 2, Name: "dog", Price: 13, Available: false},
	}
}

type appLoggerMock struct {
	mock.Mock
}

func (l appLoggerMock) Error(msg string) {
	fmt.Println(msg)
}

func TestGetOffersHandler(t *testing.T) {
	ctx := context.Background()
	//offers := generateOffers()

	reader := new(excelReaderMock)

	repo := new(repoOffersMock)
	repo.On("GetOffers", ctx).Return([]entity.Offer{entity.Offer{OfferId: 5, Name: "badger", Price: 100, Available: true}}, nil)

	log := new(appLoggerMock)

	handler := New(repo, reader, log)

	req, _ := http.NewRequest(http.MethodGet, "/getoffers", nil)
	rr := httptest.NewRecorder()

	handler.GetOffers(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("staus %d", status)
	}
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
