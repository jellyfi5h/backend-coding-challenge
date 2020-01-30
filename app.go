package main

import (
	"./server"
	"./settings"
)

func main() {
	usage := settings.GetInfo() //get the imformation ncessair(host, post) for run the server
	if usage != nil {
		//print usage
		return
	}
	err := server.CreateDaemon()
	if err != nil {
		return //something goes wrong with creating the daemon process
	}
	//set handlers(entry points)
	server.WebServer(nil) //start listen from host:port given
}
