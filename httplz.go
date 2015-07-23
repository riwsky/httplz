package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type stringConverter func(string) (string, error)

func (f stringConverter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Fatalln("bad method")
	}
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	ret, err := f(string(bs))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ret = err.Error()
	}
	w.Write([]byte(ret))
}

// args does not include any arguments to httplz itself
func makeStringConverter(args []string) (stringConverter, error) {
	if len(args) == 0 {
		log.Fatalln("Didn't receive any args to run!")
	}

	path, err := exec.LookPath(args[0])
	if err != nil {
		log.Println("Got error looking for command to run")
		log.Println(err)
		return nil, err
	}

	log.Printf("Constructing stringConverter for args %v\n", args)

	return func(in string) (string, error) {
		cmd := exec.Command(path, args[1:]...)
		cmd.Stdin = strings.NewReader(in)
		var errbuf bytes.Buffer
		cmd.Stderr = &errbuf
		out, err := cmd.Output()
		outs := string(out)
		if err != nil {
			log.Printf("Hit problems running command %v:\n", args)
			log.Println(err)
			log.Printf("Std out was:\n%v", errbuf.String())
			return "", err
		}
		return outs, nil
	}, nil
}

var DEFAULT_PORT = 8080

func printUsage() {
	fmt.Printf("Usage: httplz [--port %d] CMD\n", DEFAULT_PORT)
	fmt.Println("Will pass POST bodies on PORT to CMD, responding with stdout")
	os.Exit(1)
}

func getOpts(args []string) (int, []string) {
	port := DEFAULT_PORT
	var err error
	if len(args) == 1 {
		fmt.Println("Error: didn't pass any command to execute!")
		printUsage()
	}
	if args[1] == "--port" || args[1] == "-port" {
		port, err = strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Error: couldn't parse port number!")
			printUsage()
		}
		args = args[2:]
	}
	return port, args[1:]
}

func main() {
	port, subCommandArgs := getOpts(os.Args)
	conv, err := makeStringConverter(subCommandArgs)
	if err != nil {
		log.Fatalln(err)
	}
	http.Handle("/", conv)
	connStr := fmt.Sprintf(":%d", port)
	log.Fatalln(http.ListenAndServe(connStr, nil))
}
