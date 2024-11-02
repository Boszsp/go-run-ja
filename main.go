package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)


var inputSet []string = []string{"5555"}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()
	defer func() {
		exec.Command("./clearComplier.sh").Run()
	}()
	//Complier state
	exec.Command("javac", "CodeJv.java").Run()
	cmd := exec.CommandContext(ctx, "java", "CodeJv.java")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer stdin.Close()
		for _,input := range(inputSet){
			io.WriteString(stdin, string(input))
			io.WriteString(stdin, "\n")
		}
	}()

	stderr, err := cmd.StderrPipe()
	var wtf []byte
	if err != nil {
		log.Fatal(err)
	}
	go func ()  {
		defer stderr.Close()
        wtf, _ = io.ReadAll(stderr)
	}()

	stdout, err := cmd.StdoutPipe()
	var out []byte
	if err != nil {
		log.Fatal(err)
	}
	go func ()  {
		defer stdout.Close()
        out, _ = io.ReadAll(stdout)
	}()

	cmd.Run()


	fmt.Print(strings.Trim(string(out), ""),string(wtf))

}
