package main

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
)

type Student struct {
	Name 		string
	Pawprint 	string
	VM 			int
}

type Website struct {
	Title 	string
	URL 	string
}

type Post struct {
	Name	string
	Comment	string
}

var students = make(map[string]Student)
var websites = make(map[string]Website)
var posts = make(map[string]Post)

func studentHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":	getStudentHandler(w, req)
	case "POST": putStudentHandler(w, req)
	default:
		w.WriteHeader(400)
	}
}

func websiteHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":	getWebsiteHandler(w, req)
	case "POST": putWebsiteHandler(w, req)
	default:
		w.WriteHeader(400)
	}
}

func postHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":	getPostHandler(w, req)
	case "POST": putPostHandler(w, req)
	default:
		w.WriteHeader(400)
	}
}

func getStudentHandler(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		text, err := json.Marshal(students)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(text)
	} else {
		student, ok := students[id]
		if !ok {
			w.WriteHeader(404)
			return
		}
		
		text, err := json.Marshal(student)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(text)
	}
}

func putStudentHandler(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		w.WriteHeader(400)
		return
	}
	text, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	student := Student{}
	err = json.Unmarshal(text, &student)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	students[id] = student
}

func getWebsiteHandler(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		text, err := json.Marshal(websites)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(text)
	} else {
		website, ok := websites[id]
		if !ok {
			w.WriteHeader(404)
			return
		}
		
		text, err := json.Marshal(website)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(text)
	}
}

func putWebsiteHandler(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		w.WriteHeader(400)
		return
	}
	text, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	website := Website{}
	err = json.Unmarshal(text, &website)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	websites[id] = website
}

func getPostHandler(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		text, err := json.Marshal(posts)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(text)
	} else {
		post, ok := posts[id]
		if !ok {
			w.WriteHeader(404)
			return
		}
		
		text, err := json.Marshal(post)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(text)
	}
}

func putPostHandler(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		w.WriteHeader(400)
		return
	}
	text, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	post := Post{}
	err = json.Unmarshal(text, &post)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	posts[id] = post
}

//Main function
func main() {
	students["me"] = Student{"Kirk Backus", "kwbr4d", 9}
	posts["0"] = Post{"Kirk", "hello everybody"}
	posts["1"] = Post{"Megan", "Hi Kirk!"}
	websites["portfolio"] = Website{"My Portfolio", "http://babbage.missouri.edu/~kwbr4d/"}

	//Handle each of the REST JSON API calls
	http.HandleFunc("/students", studentHandler)
	http.HandleFunc("/websites", websiteHandler)
	http.HandleFunc("/posts", postHandler)
	
	//Serve files
	http.Handle("/", http.FileServer(http.Dir("./")))

	//Create the http server
	http.ListenAndServe(":8004", nil)
}
