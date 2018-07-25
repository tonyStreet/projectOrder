package server

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tonyStreet/projectOrder/handler"
	"log"
	"net/http"
	"strings"
	"time"
)

func router() *mux.Router {
	// router
	r := mux.NewRouter()
	r.HandleFunc("/order", handler.CreateOrder).Methods(http.MethodPost)
	r.HandleFunc("/order/{id:[0-9]+}", handler.TakeOrderHandler).Methods(http.MethodPut)
	r.HandleFunc("/orders", handler.ListOrderHandler).Queries("page", "{page:[0-9]+}", "limit", "{limit:[0-9]+}").Methods(http.MethodGet)

	log.Print("url routes:")
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		qt, err := route.GetQueriesTemplates()
		p, err := route.GetPathRegexp()
		qr, err := route.GetQueriesRegexp()
		m, err := route.GetMethods()
		if err != nil {
			log.Fatal("Cannot establish routes")
		}
		log.Println("\t", strings.Join(m, ","), strings.Join(qt, ","), strings.Join(qr, ","), t, p)
		return nil
	})

	return r
}

func server() *http.Server {

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "authorization", "content-type", "X-Auth"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 233 * time.Second,
		Handler:      handlers.CORS(originsOk, headersOk, methodsOk)(router()),
	}

	return srv
}

func Dispatch() error {
	return server().ListenAndServe()
}
