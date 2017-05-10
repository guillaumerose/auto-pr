package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"flag"

	"strings"

	"os"

	"github.com/ryanleary/jsonparser"
)

func readUser(file string) user {
	user := user{}
	dat, _ := ioutil.ReadFile(file)
	json.Unmarshal(dat, &user)
	return user
}

func main() {
	var userPath = flag.String("user", "./config.json", "GitHub configuration")
	var organization = flag.String("organization", "", "Organization")
	var repository = flag.String("repository", "", "Repository")
	var filename = flag.String("file", "", "File to change")
	var key = flag.String("key", "", "Key to change")
	var value = flag.String("value", "", "New value")
	var branch = flag.String("branch", "", "Branch name")
	var message = flag.String("message", "", "Commit message")
	var dry = flag.Bool("dry", false, "Dry run")
	var pr = flag.Bool("pr", false, "Create a PR on upstream repository")

	flag.Parse()

	if *filename == "" || *key == "" || *value == "" || *repository == "" || *branch == "" || *message == "" {
		flag.Usage()
		os.Exit(2)
	}

	ctx := context.Background()
	user := readUser(*userPath)
	git := newGit(ctx, user, *organization, *repository)

	json, sha, err := git.getContents(*filename)
	if err != nil {
		log.Fatalln("Cannot find " + *filename + " on github")
	}

	updated, err := jsonparser.Set(json, []byte("\""+*value+"\""), strings.Split(*key, ".")...)
	if err != nil {
		log.Fatalln("Cannot read/update json", err)
	}

	if *dry {
		fmt.Println(string(updated))
	} else {
		err = git.createBranch(*branch)
		if err != nil {
			log.Fatalln("Cannot create branch", err)
		}
		log.Println("Branch created")

		url, err := git.createFile(*message, *branch, *filename, sha, updated)
		if err != nil {
			log.Fatalln("Cannot push json to github", err)
		}
		log.Println("Updates pushed, see changes: " + url)

		if *pr {
			pr, err := git.createPR(*message, *branch)
			if err != nil {
				log.Fatalln("Cannot create PR", err)
			}
			log.Println("PR created: " + pr)
		}
	}
}
