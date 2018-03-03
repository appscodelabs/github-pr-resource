package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type P struct {
	Status string `json:"status"`
	Path   string `json:"path"`
}

type Input struct {
	Source struct {
		Owner       string `json:"owner"`
		Repo        string `json:"repo"`
		AccessToken string `json:"access_token"`
		Org         string `json:"org"`
		Label       string `json:"label"`
	} `json:"source"`
	Params P `json:"params"`
}

type Ver struct {
	Ref    string `json:"ref"`
	Number string `json:"number"`
}

type Output struct {
	Version Ver `json:"version"`
}

func main() {
	//log.Println("out")

	//takes input from stdin in JSON
	decoder := json.NewDecoder(os.Stdin)

	var inp Input
	err := decoder.Decode(&inp)

	if err != nil {
		log.Fatal(err)
	}

	//log.Println(inp)
	//log.Println(os.Args[1])

	//create client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: inp.Source.AccessToken})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	//get ref header from directory
	b, err := exec.Command("/find_hash.sh", os.Args[1], inp.Params.Path).Output()

	if err != nil {
		log.Fatal(err)
	}

	//update status of the pr
	newStatus := &github.RepoStatus{
		State: github.String(inp.Params.Status),
	}

	ref := string(b[:len(b)-1])

	_, _, err = client.Repositories.CreateStatus(context.Background(), inp.Source.Owner, inp.Source.Repo, ref, newStatus)
	if err != nil {
		log.Fatal(err)
	}

	//find pr_no. this need to be printed to stdout
	b, err = exec.Command("/fetch_pr.sh", os.Args[1], inp.Params.Path).Output()

	if err != nil {
		log.Fatal(err)
	}

	num := string(b[:len(b)-1])

	id, err := strconv.Atoi(num)

	if err != nil {
		log.Fatal(err)
	}

	//log.Println(id)

	//get pr from api and remove label
	_, err = client.Issues.RemoveLabelForIssue(context.Background(), inp.Source.Owner, inp.Source.Repo, id, inp.Source.Label)

	if err != nil {
		//TODO: thie code removes label from PR, but also shows error 404. IDK why
		log.Println(err.Error())
	}

	//prepare output format
	op := Output{
		Version: Ver{
			Ref:    ref,
			Number: num,
		},
	}

	b, err = json.Marshal(op)

	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stdout.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}
