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
		Label       string `json:"label"`
	} `json:"source"`
	Version struct {
		Ref string `json:"ref"`
	} `json:"version"`
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

	//log.Println(inp)

	//get list of prs from github api
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: inp.Source.AccessToken})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	pullReq, _, err := client.PullRequests.List(context.Background(), inp.Source.Owner, inp.Source.Repo, nil)

	//	log.Println("list")
	//	for i, pr := range list {
	//		log.Println(i, pr)
	//	}

	//if err != nil {
	//	log.Fatal(err)
	//}

	//sort by update date
	sort.Slice(pullReq, func(i, j int) bool {
		x := *pullReq[i].UpdatedAt
		y := *pullReq[j].UpdatedAt
		return x.Before(y)
	})

	//check which index matches with current version
	//from this index to last, all refs need to be written at stdout
	index := 0
	for i, pr := range pullReq {
		if *pr.Head.SHA == inp.Version.Ref {
			index = i
			break
		}
	}
	//log.Println(index)

	var output []Output

	//from index to rest, go through the rest of the array to filter correct prs
	for i := index; i < len(pullReq); i++ {
		//log.Println(i, *pullReq[i].User.Login)
		flag := false
		//if both is undefined, add all prs
		if inp.Source.Org == "" && inp.Source.Label == "" {
			flag = true
		} else if inp.Source.Org != "" {
			//if org is defined, first check - if user.org == inp.org
			//if yes then add pr
			//if no then check label
			flag, _, err = client.Organizations.IsMember(context.Background(), inp.Source.Org, *pullReq[i].User.Login)
			//log.Println(flag)
			if err != nil {
				log.Fatal(err)
			}
		}
		if flag == false && inp.Source.Label != "" {
			//only if label is defined
			//if label is defined then check label
			for _, lab := range pullReq[i].Labels {
				if *lab.Name == inp.Source.Label {
					flag = true
					break
				}
			}
		}

		//add to output
		if flag == true {
			output = append(output, Output{strconv.Itoa(*pullReq[i].Number), *pullReq[i].Head.SHA})
		}
	}

	//log.Println(output)
	//	log.Println("PR")
	//	for i, pr := range pullReq {
	//		log.Println(i, pr)
	//	}

	b, err := json.Marshal(output)

	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stdout.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}
