package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type problemInfo struct {
	ID    int    `json:"problem_id"`
	Owner int    `json:"owner,string"`
	Hash  string `json:"problem_spec_hash"`
}

type snapshot struct {
	Problems []problemInfo `json:"problems"`
}

const myID = 101
const apiKey = "101-ef09387a07b469087372e29dca268d27"

func main() {
	var filename string
	flag.StringVar(&filename, "f", "snapshot.json", "snapshot file")
	flag.Parse()

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	dec := json.NewDecoder(file)

	var s snapshot
	dec.Decode(&s)

	// b, _ := json.MarshalIndent(&s, "", "  ")
	// fmt.Println(string(b))

	for _, p := range s.Problems {
		if p.Owner == myID {
			continue
		}

		pname := fmt.Sprintf("problem-%d.txt", p.ID)

		if _, err := os.Stat(pname); os.IsNotExist(err) {
			log.Println("fetching", pname)

			time.Sleep(time.Second)

			client := &http.Client{}

			urlstr := fmt.Sprintf("http://2016sv.icfpcontest.org/api/blob/%s", p.Hash)
			req, err := http.NewRequest("GET", urlstr, nil)
			if err != nil {
				log.Println(err)
				continue
			}
			// ...
			req.Header.Add("X-API-Key", apiKey)
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Println(resp.Status)
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				continue
			}

			// log.Println(string(body))
			writeProblemFile(pname, body)

		} else {
			log.Printf("%s exists", pname)
		}
	}
}

func writeProblemFile(filename string, data []byte) {
	pfile, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer pfile.Close()

	_, err = pfile.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}
