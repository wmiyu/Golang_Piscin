package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	// "os"
	"context"
	"strings"
	"time"
)

type Book struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

type Inputs struct {
	n   int
	err error
}

func heartBeat(client *http.Client, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			makeConnect(client, false)
		}
		time.Sleep(time.Minute / 10)
	}
}

func makeConnect(client *http.Client, firstConnect bool) {
	resp, err := client.Get(fmt.Sprintf("http://%s:%d/connect", Host, Port))
	if err != nil {
		fmt.Printf("Client error: %v\n", "client.Get")
		return
	}
	defer resp.Body.Close()
	if firstConnect == true || resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Client error: %v\n", "ioutil.ReadAll")
			return
		}
		fmt.Printf("Status: %s  Body: %s\n", resp.Status, string(body))
	}
}

var Port int
var Host string

func main() {
	flag.IntVar(&Port, "P", 3333, "port to serv http")
	flag.StringVar(&Host, "H", "localhost", "[LeaderAddress]")
	flag.Parse()

	client := &http.Client{}
	ctx, cancel := context.WithCancel(context.Background())
	makeConnect(client, true)
	go heartBeat(client, ctx)

	var tmp, uuid, name string
	var inpt Inputs
	for inpt.n, inpt.err = fmt.Scanln(&tmp, &uuid, &name); inpt.err != io.EOF; inpt.n, inpt.err = fmt.Scanln(&tmp, &uuid, &name) {
		tmp = strings.ToUpper(tmp)
		if strings.Compare(tmp, "QUIT") == 0 {
			cancel()
			break
		} else if strings.Compare(tmp, "SET") == 0 {
			makeSet(client, uuid, name)
		} else if strings.Compare(tmp, "GET") == 0 {
			makeGet(client, uuid, name, fmt.Sprintf("http://%s:%d/getbook", Host, Port))
		} else if strings.Compare(tmp, "DELETE") == 0 {
			makeDelete(client, uuid, name)
		} else {
			fmt.Printf("Command (%s) not suported\n", tmp)
		}
		uuid = ""
		name = ""
	}
}

func makeDelete(client *http.Client, uuid, name string) {
	makeGet(client, uuid, name, fmt.Sprintf("http://%s:%d/deletebook", Host, Port))
}

func makeSet(client *http.Client, uuid, name string) {
	makeGet(client, uuid, name, fmt.Sprintf("http://%s:%d/setbook", Host, Port))
}

func makeGet(client *http.Client, uuid, name, url string) {
	var book Book

	if strings.Compare(url, fmt.Sprintf("http://%s:%d/setbook", Host, Port)) == 0 {
		if len(name) == 0 {
			fmt.Printf("Client error: %v\n", "empty name")
			return
		}
		book.Name = name
	}
	book.Uuid = uuid
	clientJson, err := json.Marshal(book)
	if err != nil {
		fmt.Printf("Client error: %v\n", "json.Marshal")
		return
	}
	bodyPost := bytes.NewReader(clientJson)
	resp, err := client.Post(url, "application/json; charset=UTF-8", bodyPost)
	if err != nil {
		fmt.Printf("Client error: %v\n", "client.Post")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Client error: %v\n", "ioutil.ReadAll")
		return
	}
	fmt.Printf("Status: %s  Body: %s\n", resp.Status, string(body))
}
