//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package excelReader

import (
	"context"
	"errors"
	"fmt"

	"github.com/tealeg/xlsx/v3"

	"github.com/sergunchig/merchant_exp.git/internal/entity"
	"github.com/sergunchig/merchant_exp.git/pkg/logger"
)

type ExcelReader struct {
	log *logger.AppLogger
}

func New(log *logger.AppLogger) *ExcelReader {
	return &ExcelReader{
		log: log,
	}
}

func (er ExcelReader) Read(file string) ([]entity.Offer, error) {
	wb, err := xlsx.OpenFile(file)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("cant read file %s", file), err)
	}
	sheet := wb.Sheets[0]
	offerlist := make([]entity.Offer, 0, sheet.MaxRow)

	for i := 1; i < sheet.MaxRow; i++ {
		r, err := sheet.Row(i)
		if err != nil {
			continue
		}
		id, err := r.GetCell(0).Int()
		if err != nil {
			continue
		}
		price, err := r.GetCell(2).Float()
		if err != nil {
			continue
		}

		offerlist = append(offerlist, entity.Offer{
			OfferId:   id,
			Name:      r.GetCell(1).Value,
			Price:     price,
			Available: r.GetCell(3).Bool(),
		})
	}
	if len(offerlist) != 0 {
		return offerlist, nil
	} else {
		return nil, fmt.Errorf(fmt.Sprintf("error read file %s", file), errors.New(""))
	}
}
func (er ExcelReader) ReadAsync(ctx context.Context, done func(), file string) (<-chan entity.Offer, error) {

	out := make(chan entity.Offer, 10)

	wb, err := xlsx.OpenFile(file)
	if err != nil {
		return nil, fmt.Errorf("error open excel file: %w", err)
	}
	sheet := wb.Sheets[0]

	go func() {
		for i := 1; i < sheet.MaxRow; i++ {
			r, err := sheet.Row(i)
			if err != nil {
				done()
				return
			}
			id, err := r.GetCell(0).Int()
			if err != nil {
				done()
				return
			}
			price, err := r.GetCell(2).Float()
			if err != nil {
				done()
				return
			}

			out <- entity.Offer{
				OfferId:   id,
				Name:      r.GetCell(1).Value,
				Price:     price,
				Available: r.GetCell(3).Bool(),
			}
		}
		close(out)
	}()
	return out, nil
}
