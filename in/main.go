package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Ver struct {
	Number string `json:"number"`
	Ref    string `json:"ref"`
}

type Input struct {
	Source struct {
		Owner       string `json:"owner"`
		Repo        string `json:"repo"`
		AccessToken string `json:"access_token"`
		Org         string `json:"org"`
	} `json:"source"`
	Version Ver `json:"version"`
}

type Output struct {
	Version Ver `json:"version"`
}

func main() {
	//takes JSON input from stdin
	decoder := json.NewDecoder(os.Stdin)
	var inp Input
	err := decoder.Decode(&inp)
	if err != nil {
		log.Fatal(err)
	}

	//now it'll fetch the repo
	//and place it in destination $1
	url := "https://github.com/" + inp.Source.Owner + "/" + inp.Source.Repo
	log.Println(url)

	cmd := exec.Command("/git_script.sh", url, os.Args[1], inp.Version.Number)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()

	if err != nil {
		log.Fatal(fmt.Sprint(err) + " : " + stderr.String())
	}

	//print output
	b, err := json.Marshal(Output{inp.Version})

	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stdout.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}
