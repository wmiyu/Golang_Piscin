/*
 * Team01 Server
 *
 * API version: 1.0.0
 */
package teamnode01

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

var Storage = make(map[string]string)
var LenUuid = 36
var RepFactor = 2

var NodesPortAddress = make(map[string]bool)

//var NodesPortAddress []string

// var isMaster = true
var PortNum int
var LeaderNode string

// This is for connect & heartbeat checks :
// curl -XGET http://127.0.0.1:3333/connect

type HBResponse struct {
	available   []string
	unavailable []string
}

func GetConnect(w http.ResponseWriter, r *http.Request) {
	if LeaderNode == "" {
		hb := nodeHeartBeat(&w)
		for _, unavailNode := range hb.unavailable {
			msg := fmt.Sprintf("Connection error to %s\n", unavailNode)
			fmt.Fprintf(w, msg)
			fmt.Print(msg)
		}
		for _, availNode := range hb.available {
			msg := fmt.Sprintf("Known nodes: %s\n", availNode)
			if _, err := fmt.Fprintf(w, msg); err != nil {
				fmt.Println("THIS:", err.Error())
			}
			fmt.Print(msg)
		}
		if len(hb.available) < RepFactor {
			msg := fmt.Sprintf("WARNING: cluster size (%d) is smaller than a replication factor (%d)!\n",
				len(hb.available), RepFactor)
			fmt.Fprintf(w, msg)
			fmt.Print(msg)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}

}

// SetBook OK :
// curl -XPOST -H "Content-Type: application/json" -d '{"uuid": "0d5d3807-5fbf-4228-a657-5a091c4e497f","name":"Chapayev_s Mustache comb"}' http://127.0.0.1:3333/setbook
// SetBook Error :
// curl -XPOST -H "Content-Type: application/json" -d '{"uuid": "12345","name":"Chapayev_s Mustache comb"}' http://127.0.0.1:3333/setbook

func nodeHeartBeat(w *http.ResponseWriter) HBResponse {

	var hb HBResponse
	hb.available = append(hb.available, fmt.Sprintf("http://localhost:%d", PortNum)) // добавление своего адресса в доступные

	client := &http.Client{}
	for child := range NodesPortAddress {
		_, err := client.Get(child + "/connect")
		if err != nil {
			if NodesPortAddress[child] != false {
				hb.unavailable = append(hb.unavailable, child)
				NodesPortAddress[child] = false
			}
		} else {
			hb.available = append(hb.available, child)
			if NodesPortAddress[child] == false {
				NodesPortAddress[child] = true
			}
		}
	}
	if len(hb.unavailable) != 0 {
		(*w).WriteHeader(http.StatusRequestTimeout)
	} else {
		(*w).WriteHeader(http.StatusOK)
	}
	return hb
}

func nodeConfirmation(b Book, w *http.ResponseWriter, reqName string) error {
	post, err := json.Marshal(b)
	if err != nil {
		code := http.StatusInternalServerError
		msg := fmt.Sprintf(`Error: json.Marshall get error`)
		(*w).WriteHeader(code)
		fmt.Fprintf(*w, msg)
		fmt.Println(msg)
		return errors.New("json.Marshall get error")
	}

	client := &http.Client{}
	for child := range NodesPortAddress {
		if NodesPortAddress[child] == false {
			continue
		}
		bodyPost := bytes.NewReader(post)
		resp, err := client.Post(child+reqName, "application/json; charset=UTF-8", bodyPost)
		if err != nil || (resp.StatusCode != http.StatusOK) {
			msg := fmt.Sprintf(`Error: the request to the child server(%s) returned an error`, child)
			(*w).WriteHeader(resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				code := http.StatusInternalServerError
				msg := fmt.Sprintf(`Error: ReadAll get error`)
				(*w).WriteHeader(code)
				fmt.Fprintf(*w, msg)
				fmt.Println(msg)
				return errors.New("ReadAll get error")
			}
			fmt.Fprintf(*w, string(respBody))
			fmt.Println(msg)
			return errors.New("the request to the child server(" + child + ") returned an error")
		}
	}
	return nil
}

type NodeInfo struct {
	Address string `json:"address"`
}

func RegisterNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var info NodeInfo

	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	code := http.StatusCreated
	w.WriteHeader(code)

	resp, err := json.Marshal(Storage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	fmt.Fprintf(w, "%s\n", resp)

	msg := fmt.Sprintf(`Added adress node : %s`, info.Address)
	fmt.Println(msg)

	NodesPortAddress[info.Address] = true

}

func SetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var b Book

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !(len(b.Uuid) == LenUuid) {
		code := http.StatusBadRequest
		msg := fmt.Sprintf(`Error: Key %s is not a proper UUID4`, b.Uuid)
		w.WriteHeader(code)
		fmt.Fprintf(w, msg)
		fmt.Println(msg)
		return
	}

	//marugula changes:
	if /* is Lieder && */ nodeConfirmation(b, &w, "/setbook") != nil {
		return
	}
	//

	Storage[b.Uuid] = b.Name
	code := http.StatusOK
	msg := fmt.Sprintf(`Created (%d replicas)`, countAvailNodes())
	w.WriteHeader(code)
	fmt.Fprintf(w, msg)
	fmt.Println(msg)
}

func countAvailNodes() int {
	count := 1
	for _, value := range NodesPortAddress {
		fmt.Println("THIS1")

		if value {
			fmt.Println("THIS2")
			count++
		}
	}
	fmt.Println("THIS3")

	return count
}

// GetBook OK :
// curl -XPOST -H "Content-Type: application/json" -d '{"uuid": "0d5d3807-5fbf-4228-a657-5a091c4e497f"}' http://127.0.0.1:3333/getbook
// GetBook Error :
// curl -XPOST -H "Content-Type: application/json" -d '{"uuid": "0d5d3807-5fbf-4228-a657-5a091c4e497a"}' http://127.0.0.1:3333/getbook

func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var b Book

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !(len(b.Uuid) == LenUuid) {
		code := http.StatusBadRequest
		msg := fmt.Sprintf(`Error: Key %s is not a proper UUID4`, b.Uuid)
		w.WriteHeader(code)
		fmt.Fprintf(w, msg)
		fmt.Println(msg)
		return
	}

	if Storage[b.Uuid] == "" {
		code := http.StatusNotFound
		w.WriteHeader(code)
		msg := fmt.Sprintf(`Not found %s`, b.Uuid)
		fmt.Fprintf(w, msg)
		fmt.Println(msg)
		fmt.Println(b.Uuid)
		return
	}

	//marugula changes:
	if /* is Lieder && */ nodeConfirmation(b, &w, "/getbook") != nil {
		return
	}
	//

	code := http.StatusOK
	msg := fmt.Sprintf(`{"name": "%s"}`, Storage[b.Uuid])
	w.WriteHeader(code)
	fmt.Fprintf(w, msg)
	fmt.Println(msg)
}

// DeleteBook OK :
// curl -XPOST -H "Content-Type: application/json" -d '{"uuid": "0d5d3807-5fbf-4228-a657-5a091c4e497f","name":"Chapayev_s Mustache comb"}' http://127.0.0.1:3333/deletebook
// DeleteBook Error :
// curl -XPOST -H "Content-Type: application/json" -d '{"uuid": "12345","name":"Chapayev_s Mustache comb"}' http://127.0.0.1:3333/deletebook

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var b Book

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !(len(b.Uuid) == LenUuid) {
		code := http.StatusBadRequest
		msg := fmt.Sprintf(`Error: Key %s is not a proper UUID4`, b.Uuid)
		w.WriteHeader(code)
		fmt.Fprintf(w, msg)
		fmt.Println(msg)
		return
	}
	if Storage[b.Uuid] == "" {
		code := http.StatusNotFound
		w.WriteHeader(code)
		msg := fmt.Sprintf(`Not found %s`, b.Uuid)
		fmt.Fprintf(w, msg)
		fmt.Println(msg)
		return
	}

	//marugula changes:
	if /* is Lieder && */ nodeConfirmation(b, &w, "/deletebook") != nil {
		return
	}
	//

	delete(Storage, b.Uuid)
	code := http.StatusOK
	msg := fmt.Sprintf(`Deleted (%d replicas)`, countAvailNodes())
	w.WriteHeader(code)
	fmt.Fprintf(w, msg)
	fmt.Println(msg)
}
