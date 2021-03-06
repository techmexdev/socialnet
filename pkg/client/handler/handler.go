package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

// Server responds to http requests.
type Server struct {
	router *mux.Router
}

type route struct {
	method  string
	path    string
	handler http.HandlerFunc
}

type handler struct {
	client        *http.Client
	template      *template.Template
	serverAddress string
}

// New creates a server that responds to http requests.
func New(serverAddress string) *Server {
	server := &Server{router: mux.NewRouter()}

	t := template.Must(template.ParseFiles(genTemplatePaths()...))

	h := handler{template: t, client: http.DefaultClient, serverAddress: serverAddress}
	rr := []route{
		{method: "GET", path: "/", handler: h.Home},
		{method: "GET", path: "/settings", handler: h.Settings},
		{method: "GET", path: "/user/{username}", handler: h.Profile},
		{method: "GET", path: "/users", handler: h.Users},
		{method: "GET", path: "/feed", handler: h.Feed},
	}

	server.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	for _, r := range rr {
		server.router.HandleFunc(r.path, r.handler).Methods(r.method)
	}

	return server
}

// ServeHTTP handles responding to http requests.
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		return
	}

	s.router.ServeHTTP(w, r)
}

func genTemplatePaths() []string {
	var sharedPaths, notSharedPaths []string

	sharedFilenames := []string{"header", "head", "footer"}
	for _, fn := range sharedFilenames {
		sharedPaths = append(sharedPaths, genTemplatePath(fn, true))
	}

	notSharedFilenames := []string{"index", "profile", "feed", "settings", "users"}
	for _, fn := range notSharedFilenames {
		notSharedPaths = append(notSharedPaths, genTemplatePath(fn, false))
	}

	return append(sharedPaths, notSharedPaths...)
}

func genTemplatePath(filename string, shared bool) string {
	if shared {
		return "./static/_shared/" + filename + ".html"
	}
	return "./static/" + filename + "/" + filename + ".html"
}

func (h *handler) proxy(w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(h.serverAddress)
	if err != nil {
		serverError(w, fmt.Errorf("error parsing %s: %s", h.serverAddress, err))
	}
	httputil.NewSingleHostReverseProxy(url).ServeHTTP(w, r)
}
