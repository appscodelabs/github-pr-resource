package main

import (
	"encoding/json"
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
	//log.Println("in")
	//log.Println(os.Args[1])

	//takes input from stdin in JSON
	decoder := json.NewDecoder(os.Stdin)
	var inp Input
	err := decoder.Decode(&inp)
	if err != nil {
		log.Fatal(err)
	}

	//log.Println(inp)

	//now it'll fetch the repo
	//and place it in destination $1
	url := "https://github.com/" + inp.Source.Owner + "/" + inp.Source.Repo
	log.Println(url)
	_, err = exec.Command("/git_script.sh", url, os.Args[1], inp.Version.Number, "pull_"+inp.Version.Ref).Output()

	if err != nil {
		log.Fatal(err)
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
