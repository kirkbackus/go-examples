package main

import (
	"fmt"
	"net/rpc"
	"log"
)

func main () {
	client, err := rpc.DialHTTP("tcp", "virtual9.cs.missouri.edu:8003")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Synchronous call
	var reply int
	err = client.Call("RPCMethods.RandomInteger", 15, &reply)
	if err != nil {
		log.Fatal("ERROR - RandInt:", err)
	}
	fmt.Println("Random(15) = ", reply)

	var msg string
	err = client.Call("RPCMethods.GetMessage", 0,  &msg)
	if err != nil {
		log.Fatal("ERROR - GetMessage:", err)
	}
	fmt.Println("GetMessage() =", msg) 
}
