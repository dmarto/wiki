package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
	"github.com/russross/blackfriday"
)

type Page struct {
	Title string
	HTML  template.HTML
}

func LoadWiki(title string, datadir string) (*Page, error) {
	filename := path.Join(datadir, title + ".md")

	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	html := blackfriday.MarkdownCommon(body)

	return &Page{
		Title: title,
		HTML:  template.HTML(html),
	}, nil
}

type Server struct {
	view    *template.Template
	config  Config
	router  *httprouter.Router
}

func ServerInit(config Config) *Server {

	tmpl := template.New("base")
	tmpl, _ = template.ParseFiles("templates/base.html")

	server := &Server{
		view:   tmpl,
		config: config,
		router: httprouter.New(),
	}

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
		s.view.Execute(w, &Page{
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

		s.view.Execute(w, page)
	}
}
