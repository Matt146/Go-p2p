package main

import (
	"Go-p2p/network"
	"fmt"
	"net/http"
	"sync"

	"github.com/opentracing/opentracing-go/log"
)

var wg sync.WaitGroup

func main() {
	network.SeedRand()
	net := network.MakeNetwork()
	network.InitMSGQueue()
	net.MyIP = "127.0.0.1" + network.Port
	net.MyID = network.GenRandBytes(32)
	http.HandleFunc("/JOIN", net.JoinHandler)
	http.HandleFunc("/PING", net.PingHandler)
	http.HandleFunc("/PONG", net.PongHandler)
	http.HandleFunc("/SendMSG", net.SendMSGHandler)
	http.HandleFunc("/BroadcastMSG", net.BroadcastMSGHandler)
	http.HandleFunc("/BroadcastMSGResponse", net.BroadcastMSGResponseHandler)
	wg.Add(1)
	defer wg.Done()
	go func() {
		log.Error(http.ListenAndServe(network.Port, nil))
	}()

	fmt.Println("Client sending JOIN request")

	// Just echo the messages now
	for {
		for k := range net.Nodes {
			packets := network.HandleMsgQueuePackets(net.Nodes[k].ID)
			for _, v := range packets {
				fmt.Println(string(v.Data))
			}
		}
	}

	// Wait now
	wg.Wait()
}
