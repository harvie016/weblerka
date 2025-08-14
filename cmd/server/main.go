package main

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

type PhotoWithQuote struct {
	Image string
	Quote string
}

type GalleryPageData struct {
	Photos []PhotoWithQuote
}

var quotes = []string{
	"Жизнь прекрасна, когда улыбаешься!",
	"Каждый момент — это новый шанс.",
	"Счастье — это быть с теми, кто тебе дорог.",
	"Верь в себя и всё получится!",
	"Мечтай. Действуй. Побеждай.",
	"Любовь — это всё, что нам нужно.",
	"Сегодня — лучший день для улыбки!",
	"Вдохновляй и будь вдохновлён.",
	"Свети ярко, как солнце!",
	"Делай добро и оно вернётся.",
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("internal/templates/home.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func GalleryHandler(w http.ResponseWriter, r *http.Request) {
	images, err := loadImages("static/img/all") // исправленный путь
	if err != nil {
		http.Error(w, "Ошибка загрузки изображений", http.StatusInternalServerError)
		return
	}
	var photos []PhotoWithQuote
	for i, img := range images {
		quote := quotes[i%len(quotes)]
		photos = append(photos, PhotoWithQuote{Image: img, Quote: quote})
	}
	tmpl, err := template.ParseFiles("internal/templates/gallery.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, GalleryPageData{Photos: photos})
}

func GalleryAPIHandler(w http.ResponseWriter, r *http.Request) {
	images, err := loadImages("static/img")
	if err != nil {
		http.Error(w, "Ошибка загрузки изображений", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}

func loadImages(dir string) ([]string, error) {
	var images []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			ext := filepath.Ext(d.Name())
			switch ext {
			case ".jpg", ".jpeg", ".png", ".gif", ".webp":
				images = append(images, d.Name())
			}
		}
		return nil
	})

	return images, err
}

type IndexPageData struct {
	Photos []PhotoWithQuote
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	images, err := loadImages("static/img/we")
	if err != nil {
		http.Error(w, "Ошибка загрузки изображений", http.StatusInternalServerError)
		return
	}
	var photos []PhotoWithQuote
	for i, img := range images {
		quote := quotes[i%len(quotes)]
		photos = append(photos, PhotoWithQuote{Image: img, Quote: quote})
	}
	tmpl, err := template.ParseFiles("internal/templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, IndexPageData{Photos: photos})
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/gallery", GalleryHandler).Methods("GET")
	r.HandleFunc("/api/gallery", GalleryAPIHandler).Methods("GET")
	r.HandleFunc("/index", IndexHandler).Methods("GET")

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	log.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
