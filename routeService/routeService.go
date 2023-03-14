package routeService

import (
	"crypto/md5"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// var wg sync.WaitGroup

var domainList = make(map[string]*Domain)

type Route struct {
	url       string
	NextCheck int64
}

type Domain struct {
	wg        WaitGroupCount
	NextCheck int64
	routes    []*Route
}

// type Domain struct {
// 	DomainUrl  string
// 	NexCheckDT string
// }

func SendRequest(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}

	status := resp.StatusCode

	return status
}

func GetMD5SumFile(f *os.File) (string, error) {
	file1Sum := md5.New()
	_, err := io.Copy(file1Sum, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", file1Sum.Sum(nil)), nil
}

func LogFailRoute(url string, status int) {
	file, err := os.OpenFile("routesOut.csv", os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		log.Println(err)
	}

	w := csv.NewWriter(file)

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	timeNow := time.Now()

	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		timeNow.Year(), timeNow.Month(), timeNow.Day(),
		timeNow.Hour(), timeNow.Minute(), timeNow.Second())
	// log.Println(status)

	row := []string{url, strconv.Itoa(status), formatted}
	err = w.Write(row)

	if err != nil {
		log.Println(err)
	}
	w.Flush()

}

func GetMD5SumString() {

}

func AddRouteToMap(route []string) {
	// listDomain := map[string]*Data{}

	// for _, val := range listRoute {
	domain, err := url.Parse(route[0])

	if err != nil {
		log.Println(err)
	}

	// ListDomain[domain.Hostname()] = &Data{
	// 	Routes: append(ListDomain[domain.Hostname()].Routes, &Route{url: val}),
	// }

	if domainList[domain.Hostname()] == nil {
		domainList[domain.Hostname()] = &Domain{}
	}

	domainList[domain.Hostname()].routes = append(domainList[domain.Hostname()].routes, &Route{url: route[0]})
	// }

}

func RouteChecking() {

	guard := make(chan int, 3)
	// count := 0
	// var wg sync.WaitGroup
	// wg.Add(1)
	for {
		//
		// guard <- 1

		for _, domain := range domainList {
			// wg.Wait()
			timeNow := time.Now().Unix()

			if domain.NextCheck < timeNow {
				if domain.wg.GetCount() == 0 {
					domain.wg.Add(1)

					go func(domain *Domain) {

						// defer func() {
						// wg.Done()
						// <-guard
						// }()

						route := domain.routes[0]

						if route.NextCheck < timeNow {
							status := SendRequest(route.url)

							fmt.Println("Текущее время:", timeNow, "URL:", route.url)

							fmt.Println(domain.NextCheck, route.NextCheck)
							timeNow := time.Now().Unix()

							domain.NextCheck = timeNow + 10
							route.NextCheck = timeNow + 20

							fmt.Println(domain.NextCheck, route.NextCheck)

							domain.routes = domain.routes[1:len(domain.routes)]
							domain.routes = append(domain.routes, route)

							if status > 300 || status < 200 {
								LogFailRoute(route.url, status)
							}
						}
						<-guard
						domain.wg.Done()
					}(domain)
				}
			}
		}
		// wg.Wait()
		//
		// for _, domain := range domainList {
		// 	domain.wg.Wait()
		// }
	}
}

type WaitGroupCount struct {
	sync.WaitGroup
	count int64
}

func (wg *WaitGroupCount) Add(delta int) {
	atomic.AddInt64(&wg.count, int64(delta))
	wg.WaitGroup.Add(delta)
}

func (wg *WaitGroupCount) Done() {
	atomic.AddInt64(&wg.count, -1)
	wg.WaitGroup.Done()
}

func (wg *WaitGroupCount) GetCount() int {
	return int(atomic.LoadInt64(&wg.count))
}

func RouteCheckingWG() {
	var wg WaitGroupCount

	for {
		for key, domain := range domainList {
			if wg.GetCount() < 2 {
				timeNow := time.Now().Unix()

				if domain.NextCheck < timeNow {
					if domain.routes[0].NextCheck < timeNow {
						if domain.wg.GetCount() == 0 {

							wg.Add(1)
							fmt.Println("Потоки:", wg.GetCount(), key)
							domain.wg.Add(1)

							go func(dom *Domain) {
								route := dom.routes[0]

									status := SendRequest(route.url)

									fmt.Println("Текущее время:", timeNow, "URL:", route.url)

									fmt.Println(dom.NextCheck, route.NextCheck)
									timeNow := time.Now().Unix()

									dom.NextCheck = timeNow + 10
									route.NextCheck = timeNow + 20

									fmt.Println(dom.NextCheck, route.NextCheck)

									dom.routes = dom.routes[1:len(dom.routes)]
									dom.routes = append(dom.routes, route)

									if status > 300 || status < 200 {
										LogFailRoute(route.url, status)
									}

								dom.wg.Done()
								wg.Done()
							}(domain)
						}
					}
				}
			}
		}
	}
}
