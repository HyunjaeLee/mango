package main

import (
	"net/http"
	"html/template"
	"parse"
	"strconv"
	"github.com/gorilla/mux"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/episode/{seriesId}", episode)
	http.ListenAndServe(":8080", router)
}

func index(w http.ResponseWriter, _ *http.Request) {
	indexPage, _ := template.ParseFiles("templates/index.html")
	indexPage.Execute(w, parse.ParseIndex())
}

func episode(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	seriesId, _ := strconv.Atoi(vars["seriesId"])
	indexPage, _ := template.ParseFiles("templates/episode.html")
	indexPage.Execute(w, parse.ParseEpisode(seriesId, 1))
}
