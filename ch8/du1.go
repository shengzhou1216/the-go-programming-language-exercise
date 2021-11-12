package main

import (
	"io/ioutil"
	"log"
	"path"
	"time"
)

func walkdir_recursive(dir string) (size int64) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("Error", err)
		return 0
	}
	for _, e := range entries {
		if e.IsDir() {
			size += walkdir_recursive(path.Join(dir, e.Name()))
			continue
		}
		size += e.Size()
	}
	return
}

func walkdir_concurrent(dir string, sizeChan chan int64) {
	log.Println("walkdir_concurrent", dir)
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("Error", err)
		return
	}
	for _, e := range entries {
		if e.IsDir() {
			go walkdir_concurrent(path.Join(dir, e.Name()), sizeChan)
			continue
		}
		sizeChan <- e.Size()
	}
}

func main() {
	start := time.Now()
	dir := "/home/sheng"
	var size int64

	// size = walkdir_recursive("/home/sheng")
	// log.Println(fmt.Sprintf("%dG", walkdir_recursive("/home/sheng")/1024/1024/1024))
	sizeChan := make(chan int64)
	go walkdir_concurrent(dir, sizeChan)
	for s,ok :=  <-sizeChan; ok; {
		size += s
	}
	log.Printf("%dG\n", size/1024/1024/1024)
	end := time.Now()
	log.Println("time:", end.Sub(start).Seconds())

}
