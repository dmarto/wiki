package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
	"github.com/russross/blackfriday"

	"github.com/namsral/flag"
)

type Config struct {
	data string
	bind string
}

func main() {
	var config Config

	flag.StringVar(&config.data, "data", "./", "path to data")
	flag.StringVar(&config.bind, "bind", "0.0.0.0:8000", "[addr]:<port> to bind to")
	flag.Parse()

	ServerInit(config)
}

type Page struct {
	Title string
	HTML  template.HTML
}

func LoadWiki(title string, dir string) (*Page, error) {
	filename := path.Join(dir, title + ".md")

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

	server := &Server{
		view:   template.Must(template.New("base").Parse(baseTmpl)),
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


var baseTmpl = `
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0, shrink-to-fit=no" />
		<meta http-equiv="x-ua-compatible" content="ie=edge" />

		<title>{{.Title}}</title>

		<style type="text/css">
			@import url("https://unpkg.com/spectre.css/dist/spectre.min.css");
			@import url("https://unpkg.com/spectre.css/dist/spectre-exp.min.css");
			@import url("https://unpkg.com/spectre.css/dist/spectre-icons.min.css");

			@import url("https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/2.8.0/github-markdown.min.css");

			.bg-white {
				background: #ffffff;
			}

			code {
				color: #5764c6;
			}

			ul {
				list-style: url;
			}
		</style>
	</head>
<body class="bg-white">
	<section class="container p-2">
		<div class="columns">
			<div class="column col-6 col-mx-auto">
				<header class="navbar">
					<section class="navbar-section">
						<ul class="breadcrumb">
							<li class="breadcrumb-item">
								<a href="/">Home</a>
							</li>
							<li class="breadcrumb-item">
								<a href="#">{{.Title}}</a>
							</li>
						</ul>
					</section>
					<section class="navbar-section">
						<div class="btn-group btn-group-block">
							<button class="btn btn-primary"><i class="icon icon-mail"></i></button>
							<button class="btn btn-primary"><i class="icon icon-download"></i></button>
							<button class="btn btn-primary"><i class="icon icon-edit"></i></button>
							<button class="btn btn-primary"><i class="icon icon-delete"></i></button>
						</div>
					</section>
				</header>
			</div>
		</div>
		<div class="divider pb-2"></div>
		<div class="columns">
			<div class="column col-6 col-xl-8 col-lg-10 col-md-12 col-mx-auto">
				<div class="markdown-body">
				{{.HTML}}
				</div>
			</div>
		</div>
	</section>
</body>
</html>
`
