package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// DiffResult represents the result of comparing two strings.
type DiffResult struct {
	Text1 string
	Text2 string
	Diff  string
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/diff", diffHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/index.html")
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func diffHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	text1 := r.FormValue("text1")
	text2 := r.FormValue("text2")

	// Perform diffing logic here
	diffResult := computeDiff(text1, text2)

	// Render the HTML template with the diff result
	renderTemplate(w, "templates/diff.html", diffResult)
}

func computeDiff(text1, text2 string) DiffResult {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, false)

	var diffString string
	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			diffString += "<span style=\"color:green\">" + diff.Text + "</span>"
		case diffmatchpatch.DiffDelete:
			diffString += "<span style=\"color:red;text-decoration:line-through\">" + diff.Text + "</span>"
		default:
			diffString += diff.Text
		}
	}

	return DiffResult{
		Text1: text1,
		Text2: text2,
		Diff:  diffString,
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data DiffResult) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
