package server

import (
	"fmt"
	"html/template"
	"log"
	"merchant_exp/internal/storage"
	"net/http"
)

func (s *Server) HomeHandler(rw http.ResponseWriter, r *http.Request) {
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
func (s *Server) UploadHandler(rw http.ResponseWriter, r *http.Request) {
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
