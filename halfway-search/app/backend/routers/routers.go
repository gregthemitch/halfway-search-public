package main

import (
	"log"
	"net/http"
)

func main() {
	// p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	// p1.save()
	// p2, _ := loadPage("TestPage")
	// fmt.Println(string(p2.Body))

	// http.HandleFunc("/view/", viewHandler)
	// http.HandleFunc("/edit/", editHandler)
	http.Handle("/", http.FileServer(http.Dir("../../frontend/pages")))
	// http.Handle("/yes", fmt.Fprintf(w, "yes"))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func test(w http.ResponseWriter, r *http.Request) {
// 	filepath := "../../frontend/pages"
// 	t, _ := template.ParseFiles(filepath + "home.html")
// 	t.Execute(w, p)
// }

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }

// func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
// 	t, _ := template.ParseFiles(tmpl + ".html")
// 	t.Execute(w, p)
// }

// func viewHandler(w http.ResponseWriter, r *http.Request) {
// 	title := r.URL.Path[len("/view/"):]
// 	p, _ := loadPage(title)
// 	renderTemplate(w, "view", p)
// }

// func editHandler(w http.ResponseWriter, r *http.Request) {
// 	title := r.URL.Path[len("/edit/"):]
// 	p, err := loadPage(title)
// 	if err != nil {
// 		p = &Page{Title: title}
// 	}
// 	renderTemplate(w, "edit", p)
// }

// type Page struct {
// 	Title string
// 	Body  []byte
// }

// func (p *Page) save() error {
// 	// This is a method named save that takes as its receiver p, a pointer to
// 	// Page . It takes no parameters, and returns a value of type error
// 	filename := p.Title + ".txt"
// 	return os.WriteFile(filename, p.Body, 0600)
// }

// func loadPage(title string) (*Page, error) {
// 	filename := title + ".txt"
// 	body, err := os.ReadFile(filename)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Page{Title: title, Body: body}, nil
// }
