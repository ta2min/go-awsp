package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

func main() {

	if len(os.Args) == 1 {
		os.Exit(selectProfile())
	}

	switch os.Args[1] {
	case "init":
		initCmd := `profile="$(go-awsp)"
export AWS_PROFILE=$profile`
		fmt.Println(initCmd)
	default:
		fmt.Fprintln(os.Stderr, "unknown command")
	}
}

func selectProfile() int {
	home, err := os.UserHomeDir()
	bytes, err := ioutil.ReadFile(filepath.Join(home, ".aws/config"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read file: %v\n", err)
		return 1
	}

	awsConfig := string(bytes)
	r := regexp.MustCompile(`\[profile (\S+)\]`)
	result := r.FindAllStringSubmatch(awsConfig, -1)
	var profiles []string
	for _, r := range result {
		profiles = append(profiles, r[1])
	}

	idx, err := fuzzyfinder.Find(profiles, func(i int) string {
		return profiles[i]
	})
	fmt.Println(profiles[idx])
	return 0
}
