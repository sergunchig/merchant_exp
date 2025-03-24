package storage

import (
	"errors"
	"fmt"
	"io"
	"merchant_exp/internal/entity"
	"merchant_exp/pkg/logger"
	"os"

	"github.com/tealeg/xlsx/v3"
)

type ExcelReaderObj struct {
	log *logger.AppLogger
}

func (er ExcelReaderObj) Read(file string) ([]entity.Offer, error) {

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
func (er ExcelReaderObj) ReadAsync(file string) (<-chan entity.Offer, error) {
	out := make(chan entity.Offer, 10)

	wb, err := xlsx.OpenFile(file)
	if err != nil {
		return nil, fmt.Errorf("file can't open %w", err)
	}
	sheet := wb.Sheets[0]

	go func() {
		defer close(out)
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

			out <- entity.Offer{
				OfferId:   id,
				Name:      r.GetCell(1).Value,
				Price:     price,
				Available: r.GetCell(3).Bool(),
			}
		}
	}()

	return out, nil
}

func SaveFile(in io.Reader, fileName string) error {
	copyname := fmt.Sprintf("./storage/%s", fileName)
	newFile, err := os.Create(copyname)
	defer newFile.Close()
	if err != nil {
		return fmt.Errorf("can't create file %w", err)
	}
	_, err = io.Copy(newFile, in)
	if err != nil {
		return fmt.Errorf("Can't copy file %w", err)
	}
	return nil
}
