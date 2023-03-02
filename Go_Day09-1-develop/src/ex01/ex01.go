package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func getBodyByUrl2(src string) string{
	
	tmp, err := http.Get(src)
	if err != nil {
		fmt.Printf(" (!) some error with http.Get: %v\n", err)
		return err.Error()
	   }
	   body, err := ioutil.ReadAll(tmp.Body)
	   if err != nil {
		fmt.Printf(" (!) some error with ioutil.ReadAll: %v\n", err)
		return err.Error()
	   }
   
	   res := string(body)
	   time.Sleep(time.Second * time.Duration(rand.Intn(3) + 1))
	   return (res)
}

func createUrlChannel(urlSlice []string) <-chan string {
	urlStringChannel := make(chan string)
	var wg = &sync.WaitGroup{}

	for i := 0; i < len(urlSlice); i++ {
		wg.Add(1)
		go func(urlString string) {
			urlStringChannel <- urlString
			wg.Done()
		}(urlSlice[i])
	}
	go func() {
		wg.Wait()
		close(urlStringChannel)
	}()
	return urlStringChannel
}

func crawlWeb(input <-chan string, ctx context.Context) <-chan string{
	outBodyChannel := make(chan string)
	var wg = &sync.WaitGroup{}

	for {
		msg, open := <-input 
		if !open {
			break
		}
		wg.Add(1)
		go func(msgStr string) {
			select {
			case <- ctx.Done():
				return
			default:
				outBodyChannel <- getBodyByUrl2(msgStr)
			}
			wg.Done()
		}(msg)
	}
	go func() {
		wg.Wait()
		close(outBodyChannel)
	}()
	return outBodyChannel
}

func writeHtmlFile(bodyStr string, t time.Time) {
	 filename := fmt.Sprintf("%d%d%d_%.2f_body.html", time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(),  time.Since(t).Seconds())
	 body := []byte(bodyStr)
   	_ = ioutil.WriteFile(filename, body, 0644)
}

func main() {
	var writeToFile bool
	flag.BoolVar(&writeToFile, "w", false, "write HTML body answer to a file")
    flag.Parse()

	t := time.Now()
	fmt.Printf(" === Starting CrawlWeb ==== \n") 
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)
	go func() {
		<-sigChannel
		fmt.Println("\n ---- GOT CTRL+C SIGNAL. Interrupting ----\n")
		stop()
		os.Exit(1)
	}()
	urlTestSlice := []string{"http://ya.ru","http://wmiyu.ru","http://vk.com","http://nic.ru",
	"http://mail.ru", "http://ok.ru",
	"http://ya.ru","http://wmiyu.ru","http://vk.com","http://nic.ru",
	"http://mail.ru", "http://ok.ru"}
	urlChannel := createUrlChannel(urlTestSlice)
	outChannel := crawlWeb(urlChannel, ctx)
	for {
		bodyStr, open := <-outChannel
		if !open{
			break
		}
		fmt.Printf(" == GET complete. Body len - %d Kb. Time elapsed: %.2f sec\n", len(bodyStr) / 1000, time.Since(t).Seconds())
		if writeToFile {
			writeHtmlFile(bodyStr, t)
		}
	}
	fmt.Printf("All Done. Time elapsed: %.2f sec\n", time.Since(t).Seconds())
}