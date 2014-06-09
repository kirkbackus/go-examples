//---------------------------------------
// Multiplexer
//   Written By Kirk Backus
//---------------------------------------

//The main package
package main

//Imports
import (
	"fmt"	
	"net/http"
	"net/url"
	"log"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"net/http/httputil"
	"hash/adler32"
)

//Structure to handle the options
type Options struct {
	Path string
	Port string
}

//Declare the options to be global
var op = &Options{Path: "./", Port: "8007"}

//Helper function to print the complex message
func PrintLogMessage(n int, r *http.Request) {
	log.Println(r.Method+" "+r.Proto+" - "+r.RemoteAddr+" [", n, "] \""+r.URL.Path+"\"")
}

//HTTP Request Handler
func HandleRequest(rw http.ResponseWriter, req *http.Request) {
	//Get adjusted filepath which will ensure proper directory
	cleanpath := op.Path+filepath.Clean(req.URL.Path);
	
	//Get the checksum of the path
	n := adler32.Checksum([]byte(cleanpath)) % 40 + 1
	//n = 9
	
	strurl := "http://virtual" + fmt.Sprintf("%d", n) + ".cs.missouri.edu:8006";
	
	fmt.Println(strurl);
	
	u, err := url.Parse(strurl)
	if err != nil {
		log.Fatal(err)
	}
	
	//Create a reverse proxy
	reverse_proxy := httputil.NewSingleHostReverseProxy(u)
	reverse_proxy.ServeHTTP(rw, req)
}

//Main function
func main() {
	//Read config.json for configuration options
	data, _ := ioutil.ReadFile("./config.json")

	//Parse config file, store results in "op"
	json.Unmarshal(data, op)

	//Set log output to the logfile
	logfile, _ := os.OpenFile("mux_log.txt", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 644)
	log.SetOutput(logfile)

	//Create the HTTP server
	http.HandleFunc("/", HandleRequest); 
	err := http.ListenAndServe(":" + op.Port, nil)
	
	//Check for any errors and log them
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
