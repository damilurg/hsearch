package background

import (
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	resp "github.com/nicklaw5/go-respond"
)

type Api struct {
	router *httprouter.Router
	st     Storage
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (m *Manager) httpApi() {
	srv := &Api{
		router: httprouter.New(),
		st:     m.st,
	}
	srv.router.NotFound = http.HandlerFunc(srv.NotFound)
	srv.router.MethodNotAllowed = http.HandlerFunc(srv.MethodNotAllowed)

	srv.router.GET("/ping", srv.ping)

	//srv.router.GET("/v1/offer", srv.offerList)
	//srv.router.POST("/v1/offers/:id", srv.sayAnswer)

	if strings.HasPrefix(m.cnf.HTTPBind, ":") {
		m.cnf.HTTPBind = "0.0.0.0" + m.cnf.HTTPBind
	}
	log.Printf("[api] Start service and listening on http://%s\n", m.cnf.HTTPBind)
	log.Fatal(http.ListenAndServe(m.cnf.HTTPBind, srv.router))
}

func (api *Api) NotFound(w http.ResponseWriter, _ *http.Request) {
	resp.NewResponse(w).NotFound(&Response{
		Success: false,
		Message: http.StatusText(http.StatusNotFound),
	})
	return
}

func (api *Api) MethodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	resp.NewResponse(w).MethodNotAllowed(&Response{
		Success: false,
		Message: http.StatusText(http.StatusMethodNotAllowed),
	})
	return
}

func (api *Api) ping(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	resp.NewResponse(w).Ok(&Response{
		Success: true,
		Message: "pong",
	})
	return
}

func (api *Api) offerList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return
}
