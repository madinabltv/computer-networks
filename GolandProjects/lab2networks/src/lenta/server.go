package main

import (
	"github.com/mgutz/logxi/v1"
	"html/template"
	"net/http"
)

const INDEX_HTML = `
    <!doctype html>
    <html lang="ru">
        <head>
            <meta charset="utf-8">
            <title>Последние новости с kruzhok.org/news</title>
			<style>
				a:hover {
					text-decoration: underline;
				}
			</style>
        </head>
        <body>
                {{range .}}
                    <a href="https://kruzhok.org/news{{.Ref}}" style="color: #333 !important;font-size: 24px;font-weight: 700;text-decoration: none;line-height: 34px;" target = "_blank">
                    	{{.Title}}
                    </a>
                    <br/>
                {{end}}
		 <a href="https://kruzhok.org/news">
			<h1 style="color: #333 !important">Заголовки новостей</h1></a>
        </body>
    </html>
    `

var indexHtml = template.Must(template.New("index").Parse(INDEX_HTML))

func serveClient(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	log.Info("got request", "Method", request.Method, "Path", path)
	if path != "/" && path != "/index.html" {
		log.Error("invalid path", "Path", path)
		response.WriteHeader(http.StatusNotFound)
	} else if err := indexHtml.Execute(response, downloadNews()); err != nil {
		log.Error("HTML creation failed", "error", err)
	} else {
		log.Info("response sent to client successfully")
	}
}

func main() {
	http.HandleFunc("/", serveClient)
	log.Info("starting listener")
	log.Error("listener failed", "error", http.ListenAndServe("127.0.0.1:6060", nil))
}
