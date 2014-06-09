//---------------------------------------
// Keystore
//   Written By Kirk Backus
//---------------------------------------
package main

import (
	"strings"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"path/filepath"
)

//(bucketObjects[bucket])[object])
//var bucketObjects = make(map[string](map[string] (map[string] Dynamic)))
var bucketObjects = make(map[string](map[string] interface{}))
var bucketObjectList = make(map[string]([]string))
var buckets = make([]string, 0, 1)

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

//Retreives all of the buckets
func getBuckets(w http.ResponseWriter, req *http.Request) {
	text, err := json.Marshal(buckets)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Write(text)
}

//Retreive all of the bucket objects
func getBucketObjects(w http.ResponseWriter, req *http.Request, reqObject string) {
	_, hasObject := bucketObjectList[reqObject]
	if !hasObject {
		w.WriteHeader(404)
		//w.Write([]byte("null"))
		return
	}
	
	text, err := json.Marshal(bucketObjectList[reqObject])
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Write(text)
}

//Retreive an object
func getObject(w http.ResponseWriter, req *http.Request, reqBucket string, reqObject string) {
	//First, check if we have the bucket, if not we need to add it
	_, hasBucket := bucketObjectList[reqBucket]
	if !hasBucket {
		w.WriteHeader(404)
		//w.Write([]byte("null"))
		return
	}
	
	//Second, check if we have the object in the bucket
	_, hasObject := bucketObjects[reqBucket][reqObject]
	if !hasObject {
		w.WriteHeader(404)
		//w.Write([]byte("null"))
		return
	}
	
	//Convert the structure to json
	text, err := json.Marshal(bucketObjects[reqBucket][reqObject])
	if err != nil {
		w.WriteHeader(500)
		return
	}
	
	//Write the result
	w.Write(text)
}

//Put object
func putObject(w http.ResponseWriter, req *http.Request, reqBucket string, reqObject string) {
	//Make sure that the bucket and object exist
	if reqBucket == "" || reqObject == "" {
		w.WriteHeader(500)
		return
	}
	
	//Read the body text
	text, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	
	//Unmarshall to make sure we have valid json
	var i interface{}
	err = json.Unmarshal(text, &i)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	
	//First, check if we have the bucket, if not we need to add it
	_, hasBucket := bucketObjectList[reqBucket]
	if !hasBucket {
		addBucket(reqBucket)
	}
	
	//Second, check if we have the object in the bucket
	_, hasObject := bucketObjects[reqBucket][reqObject]
	if !hasObject {
		addObjectToBucket(reqBucket, reqObject)
		w.WriteHeader(201)
	}

	//Add the raw json text to the bucket
	bucketObjects[reqBucket][reqObject] = i
	w.Write(text)
}

//Key Store Handler
func keyStoreHandler(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	s_path := WebPathSplit(path)

	//Switch between each of the request methods
	switch req.Method {
		case "GET": 
			switch len(s_path) {
				case 0:
					getBuckets(w, req)
				case 1:
					getBucketObjects(w, req, s_path[0])
				case 2:
					getObject(w, req, s_path[0], s_path[1])
				default:
					w.WriteHeader(500)
			}
		case "PUT":
			switch len(s_path) {
				case 2:
					putObject(w, req, s_path[0], s_path[1])
				default:
					w.WriteHeader(500)
			}
		default:
			w.WriteHeader(400)
	}
}

func addBucket(bucket string) {
	//buckets = buckets[:len(buckets)+1]
	//buckets[len(buckets)-1] = bucket
	buckets = Append(buckets, bucket)
	bucketObjectList[bucket] = make([]string, 0, 1)
	bucketObjects[bucket] = make(map[string] interface{})
}

func addObjectToBucket(bucket string, object string) {
	bucketObjectList[bucket] = Append(bucketObjectList[bucket], object)
}

//Main function
func main() {
	//Handle each of the REST JSON API calls
	http.HandleFunc("/", keyStoreHandler);

	//Create the http server
	http.ListenAndServe(":8006", nil)
}
