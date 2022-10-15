package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"strconv"

	"github.com/goombaio/namegenerator"
	"github.com/gorilla/mux"
)

type Book struct {
	Id         int
	Name       string
	AuthorName string
	AuthorId   int
	PubMonth   string
	PageCount  int
}

type Author struct {
	Id   int
	Name string
}

var bookArray []Book
var authorArray []Author
var authorMap = map[string]int{}

func randomName(mySeed int64) string {
	seed := mySeed
	nameGenerator := namegenerator.NewNameGenerator(seed)
	name := nameGenerator.Generate()
	return name
}

func randomMonth() string {
	months := []string{
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	}

	return months[rand.Intn(len(months))]
}
func populateBookSlice() []Book {

	bookArray := []Book{}
	for i := 0; i < 100; i++ {
		seed := time.Now().UTC().UnixNano()
		tempBook := Book{Id: i, Name: randomName(seed + int64(i)), AuthorName: randomName(int64(i % 7)), AuthorId: i % 7, PubMonth: randomMonth(), PageCount: rand.Intn(300)}
		bookArray = append(bookArray, tempBook)

	}
	return bookArray
}
func populateAuthor() []Author {
	authorArray := []Author{}
	for i := 0; i < 100; i++ {
		_, prs := authorMap[bookArray[i].AuthorName]
		if prs == false {
			tempAuthor := Author{Id: i, Name: bookArray[i].AuthorName}

			authorMap[bookArray[i].AuthorName] = i
			authorArray = append(authorArray, tempAuthor)
		}

	}
	return authorArray
}

func HomeLander(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(bookArray)
}

func formHandler(w http.ResponseWriter, re *http.Request) {

	name := re.FormValue("Name")
	authorName := re.FormValue("AuthorName")
	pageCount, _ := strconv.Atoi(re.FormValue("PageCount"))
	pubMonth := re.FormValue("PubMonth")
	newIndex := len(bookArray)
	if _, prs := authorMap[authorName]; prs == false {
		authorMap[authorName] = len(authorArray)
		authorArray = append(authorArray, Author{Id: authorMap[authorName], Name: authorName})

	}

	bookArray = append(bookArray, Book{Id: newIndex, Name: name, AuthorName: authorName, AuthorId: authorMap[authorName],
		PubMonth:  pubMonth,
		PageCount: pageCount,
	})

}

func deleteBook(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("Name")
	for i := 0; i < len(bookArray); i++ {
		if bookArray[i].Name == name {
			bookArray = append(bookArray[:i], bookArray[i+1:]...)
			break
		}
	}
	w.Header().Set("Content-Type", "application/json")
}

func updateBookTitle(w http.ResponseWriter, req *http.Request) {

	name := req.FormValue("Name")
	newName := req.FormValue("NewName")
	for i := 0; i < len(bookArray); i++ {
		if bookArray[i].Name == name {
			bookArray[i].Name = newName
			break
		}
	}
}

func main() {
	bookArray = populateBookSlice()
	authorArray = populateAuthor()
	r := mux.NewRouter()
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/", fileServer)
	r.HandleFunc("/listBooks", HomeLander)
	r.HandleFunc("/form", formHandler).Methods("POST")
	r.HandleFunc("/delete", deleteBook).Methods("POST")
	r.HandleFunc("/updateTitle", updateBookTitle).Methods("POST")
	http.Handle("/", r)
	err := http.ListenAndServe(":8000", r)
	if err != nil {

		log.Fatal("Error could not connect to server")

	}

	fmt.Println("the server has started")
}
