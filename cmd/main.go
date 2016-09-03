package main

import (
	"net/http"
	"github.com/ken5scal/meander"
	"encoding/json"
	"strings"
	"strconv"
	"os"
	"bytes"
	"text/template"
	"fmt"
)

func main() {
	meander.APIKey = os.Getenv("SP_GOOGLE_PLACE_API_KEY")
	http.HandleFunc("/journeys", withCORS(func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, meander.Journeys)
	}))

	http.HandleFunc("/recommendations",
		withCORS(withLog(func(w http.ResponseWriter, r *http.Request) {
			q := &meander.Query{
				Journey: strings.Split(r.URL.Query().Get("journey"), "|"),
			}
			q.Lat, _ = strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
			q.Lng, _ = strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
			q.Radius, _ = strconv.Atoi(r.URL.Query().Get("radius"))
			q.CostRangeStr = r.URL.Query().Get("cost")
			places := q.Run()
			respond(w, r, places)
		})))

	http.ListenAndServe(":8080", http.DefaultServeMux)
}

func respond(w http.ResponseWriter, r *http.Request, data []interface{}) error {
	publicData := make([]interface{}, len(data))
	for i, d := range data {
		publicData[i] = meander.Public(d)
	}
	return json.NewEncoder(w).Encode(publicData)
}

func withCORS(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}

func withLog(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bufbody := new(bytes.Buffer)
		bufbody.ReadFrom(r.Body)
		body := bufbody.String()

		line := LineOfLog{
			r.RemoteAddr,
			r.Header.Get("Content-Type"),
			r.Header.Get("Referer"),
			r.URL.Path,
			r.URL.RawQuery,
			r.Method, body,
		}
		tmpl, err := template.New("line").Parse(TemplateOfLog)
		if err != nil {
			panic(err)
		}

		bufline := new(bytes.Buffer)
		err = tmpl.Execute(bufline, line)
		if err != nil {
			panic(err)
		}

		fmt.Printf(bufline.String())
		f(w, r)
	}
}

type LineOfLog struct {
	RemoteAddr  string
	ContentType string
	Referer     string
	Path        string
	Query       string
	Method      string
	Body        string
}

var TemplateOfLog = `
Remote address:   {{.RemoteAddr}}
Content-Type:     {{.ContentType}}
HTTP method:      {{.Method}}
Referer:		  {{.Referer}}

path:
{{.Path}}

query string:
{{.Query}}

body:
{{.Body}}
`