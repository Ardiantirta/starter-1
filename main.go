package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ardiantirta/starter-1/common"
	"github.com/ardiantirta/starter-1/models"
)

func init() {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func main() {
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")
	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("sslmode", "disable")
	connStr := fmt.Sprintf("%s?%s", connection, val.Encode())

	dbConn, err := gorm.Open("postgres", connStr)
	if err != nil {
		logrus.Error(err)
	}

	if err = dbConn.DB().Ping(); err != nil {
		logrus.Error(err)
	}
	fmt.Println("ping from db")

	defer func() {
		if err = dbConn.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	dbConn.Debug().AutoMigrate(
		&models.Todo{},
	)

	r := mux.NewRouter()

	r.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := common.Message(true, "Welcome to todo service api")
		common.Response(w, response)
		return
	}))).Methods(http.MethodGet)

	// set route /pong

	// set NotFoundHandler

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})

	fmt.Println("server run on port ", viper.GetString("server.address"))
	logrus.Fatal(http.ListenAndServe(viper.GetString("server.address"), handlers.CORS(headersOk, originsOk, methodsOk)(r)))
}
