package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	auth "github.com/abbot/go-http-auth"
)

type Bird struct {
	Species     string `json:"species"`
	Descritpion string `json:"description"`
	Count		int		`json:"count,omitempty"`
}

var birds []Bird


func handlerHello(w http.ResponseWriter, r *http.Request) {
	// For this case, we will always pipe "Hello World" into the response writer
	fmt.Fprintf(w, `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>My Wikipedia</title>
	   </head>
	   
	   <body>
		<img src="/assets/lo-go.png" alt="lo-go">
		 <h1>My Wikipedia</h1>
		 <h2>Welcome onboard!</h2>`)
	fmt.Fprintf(w, `<H3><LI><a href="/assets/">Show contents<a></LI></H3>`)
	fmt.Fprintf(w, `<H3><LI><a href="/admin">Admin panel<a></LI></H3>`)
	fmt.Fprintf(w, "</body></html>")
	fmt.Printf(r.RequestURI)
}

func adminHandle(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	fmt.Fprintf(w, "<html><body><h1>Hello, %s!</h1></body></html>", r.Username)
}


func getBirdHandlers(w http.ResponseWriter, r *http.Request) {

	birds, err := store.GetBirds()

	for idx, bird := range birds {
		bird.Count = idx + 1
	}

	birdListBytes, err := json.MarshalIndent(birds, "\t", "    ")

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(birdListBytes)
	// fmt.Printf("Got N birds %d ", len(birds))
}

func createBirdHandler(w http.ResponseWriter, r *http.Request) {
	bird := Bird{}
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bird.Species = r.Form.Get("species")
	bird.Descritpion = r.Form.Get("description")

	// birds = append(birds, bird)
	err = store.CreateBird(&bird)
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/assets/", http.StatusFound)

	fmt.Println("Added bird", bird.Species, bird.Descritpion)
}


func getDocHandler(w http.ResponseWriter, r *http.Request) {

	birds, err := store.GetBirds()

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		return
	}

	fmt.Fprintf(w, "<html><body>")
	for idx, bird := range birds {
		bird.Count = idx + 1
		fmt.Fprintf(w, `<LI><a href=".">%d %s %s<a></LI>`, bird.Count, bird.Species, bird.Descritpion)
	}

	fmt.Fprintf(w, "</body></html>")

}
