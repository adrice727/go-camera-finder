package main

import (
	// "errors"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	// "fmt"
	"html/template"
	// "io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	// "net/url"
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
	templates = map[string]*template.Template{
		"brands": template.Must(template.ParseFiles("templates/brands.html", "templates/index.html")),
		"index":  template.Must(template.ParseFiles("templates/index.html")),
	}
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

func getBrandModels(brand string) []byte {
	method := "flickr.cameras.getBrandModels&brand=" + brand + "&api_key="
	request := getRequestUrl(method)
	resp, _ := http.Get(request)
	defer resp.Body.Close()
	models, _ := ioutil.ReadAll(resp.Body)
	return models
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

// type Models struct {
// 	Cameras struct {
// 		Brand  string `json:"brand"`
// 		Camera []struct {
// 			Id   string `json:"id"`
// 			Name struct {
// 				Content string `json:"_content"`
// 			} `json:"name"`
// 			Details struct {
// 				Megapixels struct {
// 					Content int `json:"_content"`
// 				} `json:"megapixels"`
// 				LcdScreenSize struct {
// 					Content int `json:"_content"`
// 				} `json:"lcd_screen_size"`
// 				MemoryType struct {
// 					Content string `json:"_content"`
// 				} `json:"memory_type"`
// 			} `json:"details"`
// 			Images struct {
// 				Small struct {
// 					Content string `json:"_content"`
// 				} `json:"small"`
// 				Large struct {
// 					Content string `json:"_content"`
// 				} `json:"large"`
// 			} `json:"images"`
// 		} `json:"camera"`
// 	} `json:"cameras"`
// 	Stat string `json:"stat"`
// }

type Page struct {
	Title string
	Body  interface{}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("index handler called")
	log.Print(r.URL.Path)
	cameraData := getCameraBrands()
	var c Cameras
	_ = json.Unmarshal(cameraData, &c)
	viewData := c.Brands.Brand
	p := Page{"brands", viewData}
	renderTemplate(w, "brands", p)

}

func brandHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	brand := vars["brand"]
	modelsData := getBrandModels(brand)
	w.Header().Set("Content-Type", "application/json")
	w.Write(modelsData)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates[tmpl].ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/brands/{brand}", brandHandler).Methods("GET")
	http.Handle("/", r)

	flag.Parse()
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
