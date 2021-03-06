package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"
)

var (
	failOnError = flag.Bool("f", false, "fail on error")
	parallel    = flag.Bool("p", false, "parallel")
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	targs := make([]*template.Template, len(args))
	for i, arg := range args {
		t, err := template.New(fmt.Sprintf("arg%d", i)).Parse(arg)
		if err != nil {
			log.Fatal(err)
		}
		targs[i] = t
	}

	scan := bufio.NewScanner(os.Stdin)

	for scan.Scan() {
		var v interface{}
		err := json.Unmarshal(scan.Bytes(), &v)
		if err != nil {
			if *failOnError {
				log.Fatal(err)
			}
			continue
		}
		xargs := make([]string, len(targs))
		for i, t := range targs {
			var buf bytes.Buffer
			err = t.Execute(&buf, v)
			if err != nil {
				log.Fatal(err)
			}
			xargs[i] = buf.String()
		}
		cmd := exec.Command(xargs[0], xargs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if *parallel {
			err = cmd.Start()
		} else {
			err = cmd.Run()
		}
		if *failOnError && err != nil {
			log.Fatal(err)
		}
	}

	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
}
