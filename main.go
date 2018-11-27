package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path"

	"github.com/bmatcuk/doublestar"
	"github.com/julienschmidt/httprouter"
	"github.com/namsral/flag"
	"github.com/shurcooL/github_flavored_markdown"
)

type Config struct {
	data string
	bind string
}

func main() {
	var config Config

	flag.StringVar(&config.data, "data", "./", "path to data")
	flag.StringVar(&config.bind, "bind", "0.0.0.0:9000", "[addr]:<port> to bind to")
	flag.Parse()

	exec.Command("xdg-open", "http://" + config.bind).Run()

	ServerInit(config)
}

type Page struct {
	Title string
	HTML  template.HTML
}

func LoadWiki(file string, dir string) (*Page, error) {
	filename := path.Join(dir, file)

	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	html := github_flavored_markdown.Markdown(body)

	return &Page{
		Title: filename,
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
	server.router.GET("/w/*file", server.WikiHandler())

	http.ListenAndServe(server.config.bind, server.router)

	return server
}

func (s *Server) RootHandler() httprouter.Handle {

	pages, _ := doublestar.Glob("./**/*.md")

	list := ""

	for _, page := range pages {
		list += "<h2><a href='w/" + page + "'>" + page + "</a></h2>"
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		s.view.Execute(w, &Page{
			Title: "Index",
			HTML:  template.HTML(list),
		})
	}
}

func (s *Server) WikiHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		file := p.ByName("file")

		page, err := LoadWiki(file, s.config.data)
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
							<li class="breadcrumb-item"><a href="/">Home</a></li>
							<li class="breadcrumb-item"><a href="#">{{.Title}}</a></li>
						</ul>
					</section>
				</header>
			</div>
		</div>
		<div class="divider pb-2"></div>
		<div class="columns">
			<div class="column col-6 col-xl-8 col-lg-10 col-md-12 col-mx-auto">
				<div class="markdown-body">{{.HTML}}</div>
			</div>
		</div>
	</section>
</body>
</html>
`
