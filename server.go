package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"github.com/julienschmidt/httprouter"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

type Page struct {
	Title string
	HTML  template.HTML
}

func LoadWiki(title string, datadir string) (*Page, error) {
	filename := path.Join(datadir, title+".txt")

	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	unsafe := blackfriday.MarkdownCommon(body)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	return &Page{
		Title: title,
		HTML:  template.HTML(html),
	}, nil
}


type Server struct {
	config    Config
	router    *httprouter.Router
	templates *Templates
}

func (s *Server) render(name string, w http.ResponseWriter, ctx interface{}) {
	buf, err := s.templates.Exec(name, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ServerInit(config Config) *Server {
	server := &Server{
		config:    config,
		router:    httprouter.New(),
		templates: NewTemplates("base"),
	}

	tmplView := template.New("view")
	template.Must(tmplView.ParseFiles("templates/view.html", "templates/base.html"))

	server.templates.Add("view", tmplView)

	server.router.GET("/", server.RootHandler())
	server.router.GET("/w/:title", server.WikiHandler())

	// TODO 404

	http.ListenAndServe(server.config.bind, server.router)

	return server
}

// Root Route Handler
func (s *Server) RootHandler() httprouter.Handle {
	// TODO, get list of all wikis
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		s.render("view", w, &Page{
			Title: "title",
			HTML:  template.HTML("<h1>asd</h1>"),
		})
	}
}

// Wiki Route Handler
func (s *Server) WikiHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		title := p.ByName("title")

		page, err := LoadWiki(title, s.config.data)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		s.render("view", w, page)
	}
}
