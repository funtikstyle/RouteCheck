package main

import (
	"RouteCheck/routeService"
	"encoding/csv"
	"io"
	"log"
	"os"
)



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

	// var wg sync.WaitGroup

	// sum, err := service.GetMD5SumFile(file)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(sum)

	r := csv.NewReader(file)
	// records, _ := r.ReadAll()
	// fmt.Println(records)

	// reply := service.GetListDomainWithSelectionRoutes(records)
	// fmt.Println(reply)

	for {
		record, err := r.Read()

		log.Println(record)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
		}
		routeService.AddRouteToMap(record)
		// reply := routeService.AddRouteToMap(record)
		// fmt.Println(reply)

	}

	routeService.RouteCheckingWG()
	
	// fmt.Println("-------------------------------------------------------------------------------")
	// routeService.RouteChecking()

	// for key := range record {
	// 	wg.Add(1)

	// 	go func(key int, record []string) {
	// 		route := record[key]
	// 		response := service.SendRequest(route)

	// 		if response > 300 || response < 200 {
	// 			service.LogFailRoute(route, response)
	// 		}

	// 		timeNow := time.Now().String()

	// 		fmt.Printf("%s %s %s\n", record[key], response, timeNow)

	// 		wg.Done()
	// 	}(key, record)

	// }
	// wg.Wait()
}
