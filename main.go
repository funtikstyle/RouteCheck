package main

import (
	"RouteCheck/service"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"sync"
)

// type Data struct {
// 	route  string
// 	status string
// 	time   string
// }

func main() {
	file, err := os.Open("routes.csv")
	if err != nil {
		log.Println(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	var wg sync.WaitGroup

	// sum, err := service.GetMD5SumString(file)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(sum)

	r := csv.NewReader(file)

	for {
		records, err := r.Read()
		log.Println(records)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
		}

		for key := range records {
			wg.Add(1)

			go func(key int, records []string){
			route := records[key]
			response := service.SendRequest(route)

			if response > 300 || response < 200 {
				 service.LogFailRoute(route, response)
			}

			timeNow := time.Now().String()
			fmt.Printf("%s %s %s\n", records[key], response, timeNow)
			
			wg.Done()
		}(key, records)
		
	}
	}
	 wg.Wait()
}


