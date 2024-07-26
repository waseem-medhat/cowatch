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
	"github.com/pelletier/go-toml"
)

type cfg struct {
	Commands []command
}

type command struct {
	Name string
	Cmd  string
}

var colors = []string{
	codes.ColorGreen,
	codes.ColorBlue,
	codes.ColorYellow,
	codes.ColorCyan,
	codes.ColorMagenta,
}

func main() {
	configBytes, err := os.ReadFile("cowatch.toml")
	if err != nil {
		log.Fatal(err)
	}

	var config cfg
	toml.Unmarshal(configBytes, &config)

	wg := &sync.WaitGroup{}
	sigintChan := make(chan os.Signal, 1)
	signal.Notify(sigintChan, os.Interrupt)

	for i, c := range config.Commands {
		wg.Add(1)
		go func(c command, color string) {
			run(c, color, sigintChan)
			wg.Done()
		}(c, colors[i%len(colors)])
	}

	wg.Wait()
}

func run(c command, color string, sigintChan chan os.Signal) {
	cmd := exec.Command("bash", "-c", c.Cmd)

	stdoutPipe, err := cmd.StdoutPipe()
	stderrPipe, err := cmd.StderrPipe()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	stdoutScanner := bufio.NewScanner(stdoutPipe)
	for stdoutScanner.Scan() {
		line := stdoutScanner.Text()
		fmt.Fprintf(os.Stdout, "%v[%v > stdout]: %v%v\n", color, c.Name, line, codes.ResetColor)
	}

	stderrScanner := bufio.NewScanner(stderrPipe)
	for stderrScanner.Scan() {
		line := stderrScanner.Text()
		fmt.Fprintf(os.Stdout, "%v[%v > stderr]: %v%v\n", color, c.Name, line, codes.ResetColor)
	}

	go func() {
		<-sigintChan
		if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
			fmt.Println(err)
		}
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v[%v error]: %v%v\n", color, c.Name, err, codes.ResetColor)
	}
}
