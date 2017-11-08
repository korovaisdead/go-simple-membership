package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	c "github.com/korovaisdead/go-simple-membership/lib/config"
	"github.com/korovaisdead/go-simple-membership/lib/controllers"
	"github.com/korovaisdead/go-simple-membership/lib/storage"
	"net/http"
	"os"
)

func main() {
	config := c.Build()
	fmt.Print(config)
	if err := storage.BuildRedisClient(); err != nil {
		fmt.Errorf(err.Error())
		panic(err)
	}

	var router *mux.Router
	router = mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/register", controllers.RegisterHandler).Methods(http.MethodPost)
	router.HandleFunc("/authenticate", controllers.LoginHandler).Methods(http.MethodPost)

	if err := http.ListenAndServe(":"+config.WebPort, handlers.LoggingHandler(os.Stdout, router)); err != nil {
		panic(err)
	}

	select {}
}
