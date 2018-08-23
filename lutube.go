package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Video struct {
	Id    string
	Title string
}

func loadVideo(id string) (*Video, error) {
	filename := "./videos/" + id + "/videodata.txt"
	videoData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	title := string(videoData)
	return &Video{Id: id, Title: title}, nil
}

func getAvailableVideos() ([]string, error) {
	videoDirectories, err := ioutil.ReadDir("videos")
	if err != nil {
		return nil, err
	}
	availableVideos := make([]string, 0)
	for _, f := range videoDirectories {
		availableVideos = append(availableVideos, f.Name())
	}
	return availableVideos, nil
}

func renderTemplate(writer http.ResponseWriter, templateFile string, data interface{}) {
	templ, _ := template.ParseFiles(templateFile)
	templ.Execute(writer, data)
}

func watchHandler(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Path[len("/watch/"):]
	video, err := loadVideo(id)
	if err != nil {
		http.Redirect(writer, request, "/", http.StatusNotFound)
	}
	renderTemplate(writer, "watch.html", video)
}

func homeHandler(writer http.ResponseWriter, request *http.Request) {
	videos, _ := getAvailableVideos()
	renderTemplate(writer, "home.html", videos)
}

func main() {
	http.HandleFunc("/watch/", watchHandler)
	http.HandleFunc("/", homeHandler)
	http.Handle("/videos", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
