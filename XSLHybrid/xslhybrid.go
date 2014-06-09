package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/http"
	"math/rand"
	"time"
	"log"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type IntArgs struct {
	color int
}

type StringArgs struct {
	msg string
}

//Structure to handle the options
type Options struct {
        Path string
        Port string
}

//Declare the options to be global
var op = &Options{Path: "./", Port: "8012"}

//Helper function to print the complex message
func PrintLogMessage(n int, r *http.Request) {
        log.Println(r.Method+" "+r.Proto+" - "+r.RemoteAddr+" [", n, "] \""+r.URL.Path+"\"")
}

type RPCMethods struct {}

var _color int = 0x000000
var _message string = "Default Message"

func (m *RPCMethods) RandomInteger(randInt int, result *int) error {
	*result = rand.Intn(randInt);
	return nil
}

func (m *RPCMethods) RandomFloat(args int, result *float32) error {
	*result = rand.Float32();
	return nil
}

func (m *RPCMethods) SetColor(args *IntArgs, result *int) error {
	_color = args.color
	*result = 0
	return nil
}

func (m *RPCMethods) GetColor(args *IntArgs, result *int) error {
	*result = _color;
	return nil
}

func (m *RPCMethods) SetMessage(str string, result *string) error {
	_message = str
	fmt.Println("Message Changed: ", _message) 
	return nil
}

func (m *RPCMethods) GetMessage(none int, result *string) error {
	*result = _message;
	return nil
}

//HTTP Request Handler
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	//Get adjusted filepath which will ensure proper directory
	cleanpath := op.Path+filepath.Clean(r.URL.Path);
	
	//PROCESS GET REQUEST
	if (r.Method == "GET") {
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
}

func main () {
	//Read config.json for configuration options
	data, _ := ioutil.ReadFile("./config.json")

	//Parse config file, store results in "op"
	json.Unmarshal(data, op)

	//Set log output to the logfile
	logfile, _ := os.OpenFile("log.txt", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 644)
	log.SetOutput(logfile)

	//Initialize the random number generator
	rand.Seed( time.Now().UTC().UnixNano())

	//Expose methods of RPCMethods instance
	rpcMethods := new(RPCMethods)
	rpc.Register(rpcMethods)
	rpc.HandleHTTP()
	
	l, e := net.Listen("tcp", ":8012")
	if e != nil {
		fmt.Println("ERR");
	}
	
	http.HandleFunc("/", HandleRequest)
	http.Serve(l, nil)
}
