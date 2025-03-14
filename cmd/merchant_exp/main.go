package main

import (
	"context"
	"fmt"
	"log"
	"merchant_exp/cmd/merchant_exp/models"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Application struct {
	ErrorLog *log.Logger
	Dao      *models.Dao
}

func initApplication() *Application {
	dbconn := os.Getenv("DBCONNECTION")
	fmt.Println("connection string: ", dbconn)
	dao, err := models.NewDao(dbconn)
	if err != nil {
		panic("app can't connect to db")
	}

	app := &Application{
		ErrorLog: log.New(os.Stdout, "Error", log.Ldate|log.Ltime),
		Dao:      dao,
	}

	return app
}

func init() {
	fmt.Println("init  project")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("can't load environment")
		panic("can't load environment")
	}
}
func main() {
	fmt.Println("start project...")
	app := initApplication()
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.HomeHandler)

	http.ListenAndServe(":8080", mux)

}

func (app *Application) HomeHandler(rw http.ResponseWriter, r *http.Request) {

	//todo проработать передачу файла
	path := "C:\\Users\\Serjio\\Documents\\merchant\\first.xlsx"

	ch, err := models.ReadOffersPipe(path)
	if err != nil {
		rw.Write([]byte(err.Error()))
	}
	err = app.Dao.CreateOffersPipe(context.Background(), ch)
	var resStr string
	if err != nil {
		resStr = "Ok"
	} else {
		resStr = err.Error()
	}

	rw.Write([]byte(resStr))
}
