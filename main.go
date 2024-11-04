package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

var inputSet []string = make([]string, 2)

func t() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)

	tempPath, err := os.MkdirTemp("./", "runJava")
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	} else {
		os.Chmod(tempPath, 0o777)
	}
	os.Chdir(tempPath)
	defer os.RemoveAll(path.Join("..",tempPath))
	code := `

	import java.util.Scanner;

	/*
	 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
	 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
	 */
	/**
	 *
	 * @author xxx
	 */
	public class CodeJv {
	
		public static void main(String[] args) {
			Scanner sc = new Scanner(System.in);
			int name = sc.nextInt();
			int b = sc.nextInt();
			CodeJv2.hello();
			System.out.println( b + (name) );
		}
	
	}
	
	`

	os.WriteFile("CodeJv.java", []byte(code), 0o777)
	code = `
/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
 */

/**
 *
 * @author xxx
 */
public class CodeJv2 {
    static public void  hello(){
        System.out.println("Hello! From CodeJv2");
    }
}

	`
	os.WriteFile("CodeJv2.java", []byte(code), 0o777)

	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}

	defer cancel()
	defer func() {
		exec.Command("./clearComplier.sh", ".").Run()
	}()

	//Complier state
	time.Sleep(time.Duration(10) * time.Millisecond)
	comp := exec.Command("javac", "CodeJv.java")
	comp.Run()


	cmd := exec.CommandContext(ctx, "java", "CodeJv.java")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer stdin.Close()
		for _, input := range inputSet {
			io.WriteString(stdin, string(input))
			io.WriteString(stdin, "\n")
		}
	}()

	stderr, err := cmd.StderrPipe()
	var wtf []byte
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer stderr.Close()
		wtf, _ = io.ReadAll(stderr)
	}()

	stdout, err := cmd.StdoutPipe()
	var out []byte
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer stdout.Close()
		out, _ = io.ReadAll(stdout)
	}()

	cmd.Run()

	fmt.Print(strings.Trim(string(out), ""), map[bool]string{false: "", true: "Error = "}[string(wtf) != ""], string(wtf))
	if(ctx.Err() != nil){
		fmt.Println(ctx.Err())
	}
}

func main() {

	for i := 1; i <= 1; i++ {
		inputSet[0] = strconv.Itoa(i*50)
		inputSet[1] = strconv.Itoa(i + 100)

		if i == 3 {
			inputSet[1] = "xxxx"

		}
		fmt.Println("Input: ", inputSet)
		t()
	}
}
