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
