package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	terraformImagesForkURL = "git@gitlab.com:Deichindianer/terraform-images.git"
)

var (
	terraformVersion = flag.String("version", "", "Specify the new 1.x.x Terraform version")
)

func main() {
	flag.Parse()
	tmpDir, err := os.MkdirTemp("", "tiu-")
	if err != nil {
		panic(err)
	}

	r, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL: terraformImagesForkURL,
	})

	ciFilePath := path.Join(tmpDir, ".gitlab-ci.yml")
	ciYamlBytes, err := os.ReadFile(ciFilePath)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(ciYamlBytes), "\n")
	for idx, line := range lines {
		if strings.Contains(line, "TERRAFORM_BINARY_VERSION: \"1.") {
			lines[idx] = "      - TERRAFORM_BINARY_VERSION: \"" + *terraformVersion + "\""
		}
	}

	updatedCiYaml := strings.Join(lines, "\n")
	err = ioutil.WriteFile(ciFilePath, []byte(updatedCiYaml), 0644)
	if err != nil {
		panic(err)
	}

	w, err := r.Worktree()
	if err != nil {
		panic(err)
	}

	_, err = w.Commit("Updated terraform to "+*terraformVersion, &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  "Philipp Böschen",
			Email: "gitlab@phil.deichindianer.de",
			When:  time.Now(),
		},
	})
	if err != nil {
		panic(err)
	}

	//TODO: open merge request need to test that
}
