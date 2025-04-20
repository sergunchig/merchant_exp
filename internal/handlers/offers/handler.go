//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

type offerLogger interface {
	Error(msg string)
}

type Handler struct {
	log offerLogger
}

func New(log offerLogger) *Handler {
	return &Handler{
		log: log,
	}
}

func (s *Handler) HomeHandler(rw http.ResponseWriter, r *http.Request) {
	html := `<html>
	<body>
	<div>
		<form action="/uploadandimport" method="post" enctype="multipart/form-data">
			Excel file: <input type="file" name="file">
			<input type="submit" value="Import">
		</form>
	</div>
	</body>
</html>`

	tmpl := template.Must(template.New("loadpage").Parse(html))

	err := tmpl.Execute(rw, nil)

	if err != nil {
		s.log.Error(fmt.Errorf("cant execute template 'home' %w", err).Error())
	}
}

// func (h *Handler) UploadAndImportHandler(rw http.ResponseWriter, r *http.Request) {

// 	offers, err := h.reader.Read(file)

// 	if err != nil {
// 		h.log.Error(err.Error())
// 		rw.Write([]byte(err.Error()))
// 		return
// 	}

// 	err = h.offers.CreateOffers(r.Context(), offers)
// 	if err != nil {
// 		h.log.Error(err.Error())
// 		rw.Write([]byte(err.Error()))
// 		return
// 	}
// 	rw.Write([]byte("Offers import is successfully"))
// }

// todo viewmodel
// func (h *Handler) GetOffers(rw http.ResponseWriter, r *http.Request) {

// 	offers, err := h.offers.GetOffers(r.Context())
// 	if err != nil {
// 		h.log.Error(err.Error())
// 		rw.WriteHeader(http.StatusInternalServerError)
// 		rw.Write([]byte("status 500"))
// 		return
// 	}
// 	data, err := json.Marshal(offers)
// 	if err != nil {
// 		h.log.Error(fmt.Errorf("json marshal error: %w", err).Error())
// 		rw.WriteHeader(http.StatusInternalServerError)
// 		rw.Write([]byte("status 500"))
// 		return
// 	}
// 	rw.Header().Set("Content-Type", "application/json")
// 	rw.WriteHeader(http.StatusOK)
// 	rw.Write(data)
// }
