package server

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"

	"github.com/arussellsaw/workbench/runner"
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("workbench")

//APIServer handles retrieval of benchmark results
type APIServer struct {
	Results []runner.ResultSet
	sync.Mutex
}

//Run the api server
func (s *APIServer) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/fetch", s.fetch)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/ui/", 301)
		}
	})
	http.Handle("/", r)

	var staticFileServer http.Handler
	staticFileServer = http.FileServer(http.Dir("/Users/alexrussell-saw/dev/go/src/github.com/arussellsaw/workbench/ui/static"))
	http.Handle("/ui/", http.StripPrefix("/ui/", staticFileServer))
	http.ListenAndServe(":8899", nil)
}

type resultFormat struct {
	Name   string
	Points [][2]float64
}

func (s *APIServer) fetch(w http.ResponseWriter, r *http.Request) {
	var stats = make(map[string][][2]float64)
	for i, set := range s.Results {
		for _, bench := range set.Benchmarks {
			var point [2]float64
			point[0] = float64(i)
			point[1] = float64(bench.OpTime.Nanoseconds())
			stats[bench.Name] = append(stats[bench.Name], point)
		}
	}
	var keys []string
	for k := range stats {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var results []resultFormat
	for _, key := range keys {
		row := resultFormat{
			Name:   key,
			Points: stats[key],
		}
		results = append(results, row)
	}
	output, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		http.Error(w, "failed to parse benchmark results", 500)
		return
	}
	w.Write(output)
}
