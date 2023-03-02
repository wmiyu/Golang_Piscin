/*
 * Team01 Server
 *
 */

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	teamnode "team01/nodeapi"
)

var portNum int
var leaderNode string

func main() {

	flag.IntVar(&portNum, "port", 3333, "port to serv http")
	flag.StringVar(&leaderNode, "leader", "", "[LeaderAddress:port]")
	flag.Parse()

	if leaderNode == "" {
		log.Printf("Server Leader Node started on %d port", portNum)
	} else {
		log.Printf("Server Slave Node started on %d port", portNum)
		nodeInfo, err := json.Marshal(teamnode.NodeInfo{Address: fmt.Sprintf("http://localhost:%d", portNum)})
		if err != nil {
			log.Fatal("Marshall error")
		}
		bodyReq := bytes.NewReader(nodeInfo)
		client := &http.Client{}
		resp, err := client.Post(leaderNode+"/registernode", "application/json; charset=UTF-8", bodyReq)
		if err != nil {
			log.Fatal(err.Error())
		}
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err.Error())
		}
		if err = json.Unmarshal(respBody, &teamnode.Storage); err != nil {
			log.Fatal(err.Error())
		}
	}

	teamnode.PortNum = portNum
	teamnode.LeaderNode = leaderNode
	router := teamnode.NewRouter()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portNum), router))
}
