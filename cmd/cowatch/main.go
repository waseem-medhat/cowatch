package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
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
	fmt.Println(codes.ColorMagenta, logo, codes.ResetColor)

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
	cmdFields := strings.Fields(c.Cmd)

	args := []string{}
	if len(cmdFields) > 1 {
		args = cmdFields[1:]
	}

	cmd := exec.Command(cmdFields[0], args...)

	cmd.Stdin = os.Stdin
	stdoutPipe, err := cmd.StdoutPipe()
	stderrPipe, err := cmd.StderrPipe()

	err = cmd.Start()
	if err != nil {
		log.Fatal("couldn't start command", c.Name, err)
	}

	go func() {
		stdoutScanner := bufio.NewScanner(stdoutPipe)
		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			fmt.Fprintf(os.Stdout, "%v[%v > stdout]: %v%v\n", color, c.Name, line, codes.ResetColor)
		}
	}()

	go func() {
		stderrScanner := bufio.NewScanner(stderrPipe)
		for stderrScanner.Scan() {
			line := stderrScanner.Text()
			fmt.Fprintf(os.Stdout, "%v[%v > stderr]: %v%v\n", color, c.Name, line, codes.ResetColor)
		}
	}()

	go func() {
		<-sigintChan
		cmd.Process.Signal(syscall.SIGINT)
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v[%v error]: %v%v\n", color, c.Name, err, codes.ResetColor)
	}
}
