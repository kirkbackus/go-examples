//---------------------------------------
// MyHTTPServer
//   Written By Kirk Backus
//---------------------------------------

//The main package
package main

//Imports
import (
	//"fmt"	
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

//Structure to handle the options
type Options struct {
	Path string
	Port string
}

//Declare the options to be global
var op = &Options{Path: "./", Port: "8001"}

//Helper function to print the complex message
func PrintLogMessage(n int, r *http.Request) {
	log.Println(r.Method+" "+r.Proto+" - "+r.RemoteAddr+" [", n, "] \""+r.URL.Path+"\"")
}

//HTTP Request Handler
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	//Get adjusted filepath which will ensure proper directory
	cleanpath := op.Path+filepath.Clean(r.URL.Path);

	//Check whether the file we are looking for exists
	file, err := os.Stat(cleanpath)
	if file == nil && err != nil {
		http.Error(w, "404: Resource Not Found", 404);
		PrintLogMessage(404, r)
	} else {
		http.ServeFile(w, r, cleanpath)
		PrintLogMessage(200, r)
	}
}

//Main function
func main() {
	//Read config.json for configuration options
	data, _ := ioutil.ReadFile("./config.json")

	//Parse config file, store results in "op"
	json.Unmarshal(data, op)

	//Set log output to the logfile
	logfile, _ := os.OpenFile("log.txt", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 644)
	log.SetOutput(logfile)

	//Create the HTTP server
	http.HandleFunc("/", HandleRequest); 
	err := http.ListenAndServe(":" + op.Port, nil)
	
	//Check for any errors and log them
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
