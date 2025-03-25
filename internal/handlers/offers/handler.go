package offer

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/sergunchig/merchant_exp.git/internal/entity"
	"github.com/sergunchig/merchant_exp.git/internal/storage"
	"github.com/sergunchig/merchant_exp.git/pkg/logger"
)

type repoOffers interface {
	GetOffers(ctx context.Context) ([]entity.Offer, error)
}

type Handler struct {
	offers repoOffers
	log    *logger.AppLogger
}

func New(repo repoOffers, log *logger.AppLogger) *Handler {
	return &Handler{
		offers: repo,
		log:    log,
	}
}

func (s *Handler) HomeHandler(rw http.ResponseWriter, r *http.Request) {
	html := `<html>
	<body>
	<div>
		<form action="/upload" method="post" enctype="multipart/form-data">
			Excel file: <input type="file" name="my_file">
			<input type="submit" value="Upload">
		</form>
	</div>
	</body>
</html>`

	tmpl := template.Must(template.New("loadpage").Parse(html))

	err := tmpl.Execute(rw, nil)

	if err != nil {

	}
}
func (h *Handler) UploadHandler(rw http.ResponseWriter, r *http.Request) {
	uploadData, _, err := r.FormFile("my_file")
	if err != nil {
		log.Println("cant parse file", err)
		http.Error(rw, "request error", http.StatusInternalServerError)
		return
	}
	defer uploadData.Close()
	fmt.Println("upload")
	err = storage.SaveFile(uploadData, "excelfile.xlsx")
	if err != nil {
		fmt.Println(err)
	}
}
