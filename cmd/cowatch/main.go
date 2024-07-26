package main

import (
	"log"
	"os"
	"os/exec"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}

	sites := []string{"google.com", "github.com"}
	wg.Add(len(sites))

	for _, site := range sites {
		go func(site string) {
			cmd := exec.Command("ping", site)

			cmd.Stdout = os.Stdout
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(site)
	}

	wg.Wait()
}
