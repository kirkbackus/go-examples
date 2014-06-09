package main

import (
	"fmt"
	"net/rpc"
	"net/http"
	"math/rand"
	"time"
)

type Args struct {
	max int
}

type RPCMethods struct {}

func (m *RPCMethods) RandomInteger(args *Args, result *int) error {
	*result = rand.Intn(args.max);
	return nil
}

func (m *RPCMethods) RandomFloat(args *Args, result *float32) error {
	*result = rand.Float32();
	return nil
}

func main () {
	// Initialize the random number generator
	rand.Seed( time.Now().UTC().UnixNano())

	// expose methods of RPCMethods instance
	rpcMethods := new(RPCMethods)
	rpc.Register(rpcMethods)
	rpc.HandleHTTP()

	//Create the HTTP server to handle RPC requests
	err := http.ListenAndServe(":8002", nil)
	if err != nil {
		fmt.Println(err.Error());
	}
}
