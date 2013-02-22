package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

var (
	prefix    *string
	org       *string
	repoCount *uint
	dir       *string
)

func init() {
	prefix = flag.String("p", "https://github.com", "prefix of repositories")
	org = flag.String("o", "golang-samples", "organization name")
	repoCount = flag.Uint("n", 200, "number of repositories")
	dir = flag.String("f", "", "dest directory")
}

type repo struct {
	Name string `json: "name"`
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case error:
				err := r.(error)
				fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			default:
				fmt.Fprintf(os.Stderr, "Error: %v\n", r)
			}
		}
	}()

	flag.Parse()
	res, err := http.Get(fmt.Sprintf("https://api.github.com/orgs/%s/repos?per_page=%d", *org, *repoCount))
	if err != nil {
		panic(err)
	}

	var repos []repo
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&repos); err != nil {
		panic(err)
	}

	for _, r := range repos {
		url := fmt.Sprintf("%s/%s/%s", *prefix, *org, r.Name)
		var git *exec.Cmd
		if *dir == "" {
			git = exec.Command("git", "clone", url)
		} else {
			dest := fmt.Sprintf("%s/%s", *dir, r.Name)
			git = exec.Command("git", "clone", url, dest)
		}

		if stdout, err := git.StderrPipe(); err == nil {
			out, _ := git.Output()
			if string(out) != "" {
				fmt.Printf("%s", out)
			}
			io.Copy(os.Stdout, stdout)
		} else {
			panic(err)
		}
	}
}
