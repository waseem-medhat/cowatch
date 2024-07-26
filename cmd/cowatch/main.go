package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	codes "github.com/avearmin/stylecodes"
)

var colors = []string{
	codes.ColorGreen,
	codes.ColorBlue,
	codes.ColorYellow,
	codes.ColorCyan,
	codes.ColorMagenta,
}

func main() {
	sites := []string{"google.com", "github.com", "duck.com"}
	wg := &sync.WaitGroup{}
	sigintChan := make(chan os.Signal, 1)
	signal.Notify(sigintChan, os.Interrupt)

	for i, site := range sites {
		wg.Add(1)
		go func(site string, color string) {
			run(site, color, sigintChan)
			wg.Done()
		}(site, colors[i%len(colors)])
	}

	wg.Wait()
}

func run(site string, color string, sigintChan chan os.Signal) {
	cmd := exec.Command("ping", site)
	stdout, err := cmd.StdoutPipe()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintf(os.Stdout, "%v[%v]: %v%v\n", color, site, line, codes.ResetColor)
	}

	go func() {
		<-sigintChan
		if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
			fmt.Println(err)
		}
	}()

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
