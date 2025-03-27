package offer

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/sergunchig/merchant_exp.git/internal/entity"
	"github.com/sergunchig/merchant_exp.git/internal/storage"
)

type repoOffers interface {
	CreateOffers(ctx context.Context, offers []entity.Offer) error
}
type offersReader interface {
	Read(file string) ([]entity.Offer, error)
}
type offerLogger interface {
	Error(msg string)
}

type Handler struct {
	offers repoOffers
	reader offersReader
	log    offerLogger
}

func New(repo repoOffers, reader offersReader, log offerLogger) *Handler {
	return &Handler{
		offers: repo,
		reader: reader,
		log:    log,
	}
}

func (s *Handler) HomeHandler(rw http.ResponseWriter, r *http.Request) {
	html := `<html>
	<body>
	<div>
		<form action="/UploadAndImport" method="post" enctype="multipart/form-data">
			Excel file: <input type="file" name="my_file">
			<input type="submit" value="Import">
		</form>
	</div>
	</body>
</html>`

	tmpl := template.Must(template.New("loadpage").Parse(html))

	err := tmpl.Execute(rw, nil)
	if err != nil {

	}
}
func (h *Handler) UploadAndImportHandler(rw http.ResponseWriter, r *http.Request) {
	uploadData, _, err := r.FormFile("my_file")
	if err != nil {
		log.Println("cant parse file", err)
		http.Error(rw, "request error", http.StatusInternalServerError)
		return
	}
	defer uploadData.Close()

	file := "./storage/excelfile.xlsx"
	err = storage.SaveFile(uploadData, file)
	if err != nil {
		h.log.Error(err.Error())
		rw.Write([]byte(err.Error()))
		return
	}
	offers, err := h.reader.Read(file)

	if err != nil {
		h.log.Error(err.Error())
		rw.Write([]byte(err.Error()))
		return
	}

	err = h.offers.CreateOffers(r.Context(), offers)
	if err != nil {
		h.log.Error(err.Error())
		rw.Write([]byte(err.Error()))
		return
	}
	rw.Write([]byte("Offers import is successfully"))
}
