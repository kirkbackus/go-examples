//---------------------------------------
// Search
//   Written By Kirk Backus
//---------------------------------------

//The main package
package main

//"

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"bufio"
	"strings"
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
var op = &Options{Path: "./", Port: "8008"}

//Helper function to print the complex message
func PrintLogMessage(n int, r *http.Request) {
	log.Println(r.Method+" "+r.Proto+" - "+r.RemoteAddr+" [", n, "] \""+r.URL.Path+"\"")
}

//Get SplitPath
func WebPathSplit(path string) []string {
	//Get the path and do some trimming
	path = strings.TrimSpace(path)
	splitPath := strings.Split(filepath.Clean(path), "/")
	
	//Remove empty items at the beginning
	if len(strings.TrimSpace(splitPath[0])) == 0 {
		splitPath = splitPath[1:]
	}
	
	//Remove empty items at the end
	if len(strings.TrimSpace(splitPath[len(splitPath)-1])) == 0 {
		splitPath = splitPath[:len(splitPath)-1]
	}
	
	//Return the split path
	return(splitPath)
}

//Append string to a string array
func Append(slice []string, data string) []string {
    slen := len(slice)
    if slen + 1 > cap(slice) {  // reallocate
        // Allocate double what's needed, for future growth.
        newSlice := make([]string, slen*2)
        // The copy function is predeclared and works for any slice type.
        copy(newSlice, slice)
        slice = newSlice
    }
    slice = slice[0:slen+1]
    slice[slen] = data
    return slice
}

//HTTP Request Handler
func HandleRequest(rw http.ResponseWriter, req *http.Request) {
	//Get adjusted filepath which will ensure proper directory
	cleanpath := op.Path+filepath.Clean(req.URL.Path);
	s_path := WebPathSplit(cleanpath)
	
	if len(s_path) != 1 || s_path[0] == "." || s_path[0] == "favicon.ico" {
		rw.WriteHeader(404)
		return
	}
	
	getreq := fmt.Sprintf("GET /%s HTTP/1.1\r\n\r\n", s_path[0])
	//fmt.Printf(getreq)
	
	search_arr := make([]string, 0, 1)
	
	for i:=1; i<=40; i++ {
		strurl := "virtual" + fmt.Sprintf("%d", i) + ".cs.missouri.edu:8006";
		conn, err := net.Dial("tcp", strurl)
		
		if err != nil {
			//fmt.Printf("[%d] - Could not connect\n", i)
			continue
		}

		fmt.Fprintf(conn, "%s", getreq)
		bread := bufio.NewReader(conn)
		_, err = bread.ReadString('\n')
		//fmt.Printf("[%d] - Status: %s\n", i, status)
		
		for ;; {
			header, err := bread.ReadString('\n')
			//fmt.Printf("(%d) - %s\n", len(header), header)
			if len(header) == 2 || err != nil {
				break
			}
		}
		
		body_length_str, err := bread.ReadString('\n')
		
		if (err != nil) {
			continue
		}
		
		//Trim dat
		body_length_str = strings.TrimSpace(body_length_str)
		
		if (len(body_length_str) == 0) {
			continue
		}
		
		//Get the body length and parse the integer as a 
		body_length, err := strconv.ParseInt(body_length_str, 16, 32)
		
		if err != nil {
			//fmt.Printf("ERR Regarding String Conversion: \"%s\"\n", body_length_str)
			continue
		}
		
		//Read the bytes
		content := make([]byte, body_length)
		bread.Read(content)
		
		//Parse the content as json
		int_str := make([]string, 1024)
		err = json.Unmarshal(content, &int_str)
		if err != nil {
			//fmt.Printf("[%d] - FAILED TO UNMARSHAL: %s\n\tCONTENT: \"%s\"\n", i, err.Error(), content)
			continue
		}
		
		for i:=0; i<len(int_str); i++ {
			//fmt.Printf("ADDED: %s\n", int_str[i])
			search_arr = Append(search_arr, int_str[i])
		}
	}
	
	text, err := json.Marshal(search_arr)
	//fmt.Printf("MARSHALED: %s\n", text)
	
	if err != nil {
		rw.WriteHeader(500)
		return
	}

	rw.Write(text)
}

//Main function
func main() {
	//Read config.json for configuration options
	data, _ := ioutil.ReadFile("./config.json")

	//Parse config file, store results in "op"
	json.Unmarshal(data, op)

	//Set log output to the logfile
	logfile, _ := os.OpenFile("search_log.txt", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 644)
	log.SetOutput(logfile)

	//Create the HTTP server
	http.HandleFunc("/", HandleRequest); 
	err := http.ListenAndServe(":" + op.Port, nil)
	
	//Check for any errors and log them
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
