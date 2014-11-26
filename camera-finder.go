package main

import (
	// "errors"
	"encoding/json"
	"flag"
	// "fmt"
	"html/template"
	// "io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	// "reflect"
	// "regexp"
)

const (
	endpoint = "https://api.flickr.com/services/rest/?method="
	apiKey   = "034084e694de5e0315e418b599615241"
	format   = "&format=json&nojsoncallback=1"
)

var (
	addr      = flag.Bool("addr", false, "find open address and print to final-port.txt")
	templates = template.Must(template.ParseFiles("brands.html"))
)

func getRequestUrl(method string) string {
	return endpoint + method + apiKey + format
}

func getCameraBrands() []byte {
	method := "flickr.cameras.getBrands&api_key="
	request := getRequestUrl(method)
	resp, _ := http.Get(request)
	defer resp.Body.Close()
	brands, _ := ioutil.ReadAll(resp.Body)
	return brands
}

// http://mholt.github.io/json-to-go/
type Cameras struct {
	Brands struct {
		Brand []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"brand"`
	} `json:"brands"`
}

type Page struct {
	Title string
	Body  interface{}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	cameraData := getCameraBrands()
	var c Cameras
	_ = json.Unmarshal(cameraData, &c)
	viewData := c.Brands.Brand
	p := Page{"brands", viewData}
	renderTemplate(w, "brands", p)

}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	flag.Parse()
	http.HandleFunc("/", indexHandler)
	if *addr {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			log.Fatal(err)
		}
		s := &http.Server{}
		s.Serve(l)
		return
	}

	port := ":8080"
	log.Print("Now listening on port " + port[1:])
	http.ListenAndServe(port, nil)
}
