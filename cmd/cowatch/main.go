package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	codes "github.com/avearmin/stylecodes"
)

func main() {
	colors := []string{codes.ColorGreen, codes.ColorBlue}
	sites := []string{"google.com", "github.com"}

	wg := &sync.WaitGroup{}

	for i, site := range sites {
		wg.Add(1)
		go func(site string, color string) {
			run(site, color)
			wg.Done()
		}(site, colors[i])
	}

	wg.Wait()
}

func run(site string, color string) {
	cmd := exec.Command("ping", site)
	stdout, err := cmd.StdoutPipe()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintln(os.Stdout, color, site, line, codes.ResetColor)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
