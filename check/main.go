package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Input struct {
	Source struct {
		Owner       string `json:"owner"`
		Repo        string `json:"repo"`
		AccessToken string `json:"access_token"`
		Org         string `json:"org"`
	} `json:"source"`
	Version struct {
		Ref string `json:"ref"`
	} `json:"version"`
}

type PullReq struct {
	Number int `json:"number"`
	User   struct {
		Login string `json:"login"`
	} `json:"user"`

	Head struct {
		Sha string `json:"sha"`
	} `json:"head"`
	UpdatedAt string `json:"updated_at"`
}

type Output struct {
	Number string `json:"number"`
	Ref    string `json:"ref"`
}

func main() {
	//takes input from stdin in JSON
	decoder := json.NewDecoder(os.Stdin)
	var inp Input
	err := decoder.Decode(&inp)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(inp)

	//get prs from github api
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: inp.Source.AccessToken})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	list, _, err := client.PullRequests.List(context.Background(), inp.Source.Owner, inp.Source.Repo, nil)

	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(list)

	var pullReq []PullReq

	err = json.Unmarshal(b, &pullReq)
	if err != nil {
		log.Fatal(err)
	}

	for _, pr := range pullReq {
		log.Println(pr)
	}

	//sort by update date
	sort.Slice(pullReq, func(i, j int) bool {
		return pullReq[i].UpdatedAt < pullReq[j].UpdatedAt
	})

	log.Println("--------")
	for _, pr := range pullReq {
		log.Println(pr)
	}

	//check which index matches with current version
	index := 0
	for i, pr := range pullReq {
		if pr.Head.Sha == inp.Version.Ref {
			index = i
			break
		}
	}
	log.Println(index)

	//from index to rest, go through the rest of the array and creat json output

	var output []Output
	for i := index; i < len(pullReq); i++ {
		if inp.Source.Org == "" {
			//if source org = nil, no need to check org
			output = append(output, Output{strconv.Itoa(pullReq[i].Number), pullReq[i].Head.Sha})
		} else {
			//now check user org
			flag, _, err := client.Organizations.IsMember(context.Background(), inp.Source.Org, pullReq[i].User.Login)
			if err != nil {
				log.Fatal(err)
			}
			if flag == true {
				output = append(output, Output{strconv.Itoa(pullReq[i].Number), pullReq[i].Head.Sha})
			}
		}
	}

	b, err = json.Marshal(output)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("------final-----")

	_, err = os.Stdout.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}
