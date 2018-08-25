package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Video struct {
	Id    string
	Title string
}

type HomePage struct {
	Videos []string
	Error  bool
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

func requestHandler(
	writer http.ResponseWriter,
	request *http.Request,
	handler func(request *http.Request) (string, interface{}, error)) {
	template, data, err := handler(request)
	if err != nil {
		http.Redirect(writer, request, "/?error=true", http.StatusSeeOther)
	}
	renderTemplate(writer, template, data)
}

func watchHandler(request *http.Request) (string, interface{}, error) {
	id := request.URL.Path[len("/watch/"):]
	video, err := loadVideo(id)
	return "watch.html", video, err
}

func getPreviousError(request *http.Request) bool {
	errParamArr := request.URL.Query()["error"]
	errParam := strings.Join(errParamArr, "")
	if errParam == "true" {
		return true
	}
	return false
}

func homeHandler(request *http.Request) (string, interface{}, error) {
	previousErr := getPreviousError(request)
	videos, err := getAvailableVideos()
	showErrorMessage := err != nil || previousErr
	homePage := HomePage{Videos: videos, Error: showErrorMessage}
	return "home.html", homePage, nil
}

func main() {
	http.HandleFunc("/watch/", func(writer http.ResponseWriter, request *http.Request) {
		requestHandler(writer, request, watchHandler)
	})
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		requestHandler(writer, request, homeHandler)
	})
	http.Handle("/videos", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
