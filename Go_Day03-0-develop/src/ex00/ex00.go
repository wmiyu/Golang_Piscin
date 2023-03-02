package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var (
	_     = fmt.Print
	count int
	batch int
)

type RestorantData struct {
	ID     	int    `json:"id"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	Phone    string    `json:"phone"`
	Location GLocation `json:"location"`
}
type GLocation struct {
	Longitude string `json:"lon"`
	Latitude  string `json:"lat"`
}

func createRestorantData(data [][]string) []RestorantData {
	var restorantData []RestorantData

	for i, line := range data {
		if i > 0 {
			var rec RestorantData
			for j, field := range line {
				if j == 0 {
					rec.ID = i
				} else if j == 1 {
					rec.Name = field
				} else if j == 2 {
					rec.Address = field
				} else if j == 3 {
					rec.Phone = field
				} else if j == 4 {
					rec.Location.Longitude = field
				} else if j == 5 {
					rec.Location.Latitude = field
				}
			}
			restorantData = append(restorantData, rec)
		}
	}
	return restorantData
}

func makeRdatafromCSV() []RestorantData {
	f, err := os.Open("../../materials/data.csv")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'
	csvdata, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	restorantData := createRestorantData(csvdata)
	return restorantData
}

func init() {
	flag.IntVar(&count, "count", 15000, "Number of documents to generate")
	flag.IntVar(&batch, "batch", 150, "Number of documents to send in one batch")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
}

func main() {
	log.SetFlags(0)

	type bulkResponse struct {
		Errors bool `json:"errors"`
		Items  []struct {
			Index struct {
				ID     string `json:"_id"`
				Result string `json:"result"`
				Status int    `json:"status"`
				Error  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
					Cause  struct {
						Type   string `json:"type"`
						Reason string `json:"reason"`
					} `json:"caused_by"`
				} `json:"error"`
			} `json:"index"`
		} `json:"items"`
	}

	var (
		buf bytes.Buffer
		res *esapi.Response
		err error
		raw map[string]interface{}
		blk *bulkResponse

		indexName  = "places"
		numItems   int
		numErrors  int
		numIndexed int
		numBatches int
		currBatch  int
	)

	log.Println(strings.Repeat("▁", 65))

	log.Printf("Bulk: documents [%d] batch size [%d]\n",
		count, batch)
	log.Println(strings.Repeat("▁", 65))

	// Create the Elasticsearch client
	//
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	fmt.Print("→ Sending batch ")

	// Re-create the index
	//
	if res, err = es.Indices.Delete([]string{indexName}); err != nil {
		log.Fatalf("Cannot delete index: %s", err)
	}
	mapping := `
    {
      "settings": {
        "number_of_shards": 1
      },
      "mappings": {
		"properties": {
		  "name": {
			  "type":  "text"
		  },
		  "address": {
			  "type":  "text"
		  },
		  "phone": {
			  "type":  "text"
		  },
		  "location": {
			"type": "geo_point"
		  }
		}
	  }
    }`

	res, err = es.Indices.Create(indexName, es.Indices.Create.WithBody(strings.NewReader(mapping)))
	// res, err = es.Indices.Create(indexName)

	if err != nil {
		log.Fatalf("Cannot create index: %s", err)
	}
	if res.IsError() {
		log.Fatalf("Cannot create index: %s", res)
	}

	if count%batch == 0 {
		numBatches = (count / batch)
	} else {
		numBatches = (count / batch) + 1
	}

	start := time.Now().UTC()

	restorantData := makeRdatafromCSV()

	// Loop over the collection
	//

	for i, a := range restorantData {
		numItems++

		currBatch = i / batch
		if i == count-1 {
			currBatch++
		}

		// Prepare the metadata payload
		//
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : %d } }%s`, a.ID, "\n"))
		fmt.Printf("%s", meta) // <-- Uncomment to see the payload

		// Prepare the data payload: encode article to JSON
		//
		data, err := json.Marshal(a)

		if err != nil {
			log.Fatalf("Cannot encode article: %s", err)
		}

		// Append newline to the data payload
		//
		data = append(data, "\n"...) // <-- Comment out to trigger failure for batch

		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)

		// When a threshold is reached, execute the Bulk() request with body from buffer
		//
		if i > 0 && i%batch == 0 || i == count-1 {
			fmt.Printf("[%d/%d] ", currBatch, numBatches)

			res, err = es.Bulk(bytes.NewReader(buf.Bytes()), es.Bulk.WithIndex(indexName))
			if err != nil {
				log.Fatalf("Failure indexing batch %d: %s", currBatch, err)
			}
			// If the whole request failed, print error and mark all documents as failed
			//
			if res.IsError() {
				numErrors += numItems
				if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
					log.Fatalf("Failure to to parse response body: %s", err)
				} else {
					log.Printf("  Error: [%d] %s: %s",
						res.StatusCode,
						raw["error"].(map[string]interface{})["type"],
						raw["error"].(map[string]interface{})["reason"],
					)
				}
				// A successful response might still contain errors for particular documents...
				//
			} else {
				if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
					log.Fatalf("Failure to to parse response body: %s", err)
				} else {
					for _, d := range blk.Items {
						// ... so for any HTTP status above 201 ...
						//
						if d.Index.Status > 201 {
							// ... increment the error counter ...
							//
							numErrors++

							// ... and print the response status and error information ...
							log.Printf("  Error: [%d]: %s: %s: %s: %s",
								d.Index.Status,
								d.Index.Error.Type,
								d.Index.Error.Reason,
								d.Index.Error.Cause.Type,
								d.Index.Error.Cause.Reason,
							)
						} else {
							// ... otherwise increase the success counter.
							//
							numIndexed++
						}
					}
				}
			}

			// Close the response body, to prevent reaching the limit for goroutines or file handles
			//
			res.Body.Close()

			// Reset the buffer and items counter
			//
			buf.Reset()
			numItems = 0
		}
	}

	// Report the results: number of indexed docs, number of errors, duration, indexing rate
	//
	fmt.Print("\n")
	log.Println(strings.Repeat("▔", 65))

	dur := time.Since(start)

	if numErrors > 0 {
		log.Fatalf(
			"Indexed [%d] documents with [%d] errors in %d (%.2f docs/sec)",
			((numIndexed)),
			((numErrors)),
			dur.Truncate(time.Millisecond),
			((1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%d] documents in %d (%.2f docs/sec)",
			((numIndexed)),
			dur.Truncate(time.Millisecond),
			((1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
		)
	}
}
