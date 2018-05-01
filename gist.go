package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	fileName            string
	public              bool
	token               string // GITHUB_TOKEN_FOR_GIST
	enterpriseBaseUrl   string // GIST_ENTERPRISE_BASE_URL
	enterpriseUploadUrl string // GIST_ENTERPRISE_UPLOAD_URL
	isEnterprise        bool
)

func init() {
	flag.StringVar(&fileName, "f", "", "gist file name")
	flag.BoolVar(&public, "p", false, "make gist public")
}

type TokenSource oauth2.Token

func (t *TokenSource) Token() (*oauth2.Token, error) {
	return (*oauth2.Token)(t), nil
}

type GistFiles map[github.GistFilename]github.GistFile

func getFilesFromStdin() (GistFiles, error) {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	content := string(data)
	return GistFiles{
		github.GistFilename(fileName): github.GistFile{
			Content: &content,
		},
	}, nil
}

func readFile(fname string) (string, error) {
	f, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func getFilesFromArgs() (GistFiles, error) {
	files := make(GistFiles)
	for _, arg := range flag.Args() {
		content, err := readFile(arg)
		if err != nil {
			return nil, err
		}
		name := github.GistFilename(path.Base(arg))
		files[name] = github.GistFile{
			Content: &content,
		}
	}
	return files, nil
}

func getFiles() (GistFiles, error) {
	if flag.NArg() > 0 {
		return getFilesFromArgs()
	} else {
		return getFilesFromStdin()
	}
}

func getValuesFromEnv() {
	token = os.Getenv("GITHUB_TOKEN_FOR_GIST")
	if len(token) < 1 {
		log.Fatal("must set GITHUB_TOKEN_FOR_GIST")
	}

	baseUrl := os.Getenv("GIST_ENTERPRISE_BASE_URL")
	if len(baseUrl) < 1 {
		return
	}

	uploadUrl := os.Getenv("GIST_ENTERPRISE_UPLOAD_URL")
	if len(uploadUrl) < 1 {
		return
	}

	enterpriseBaseUrl = baseUrl
	enterpriseUploadUrl = uploadUrl
	isEnterprise = true
}

func main() {
	flag.Parse()

	getValuesFromEnv()

	files, err := getFiles()
	if err != nil {
		log.Fatal(err)
	}
	ts := TokenSource{AccessToken: token}

	var client *github.Client
	if isEnterprise {
		client, err = github.NewEnterpriseClient(
			enterpriseBaseUrl,
			enterpriseUploadUrl,
			oauth2.NewClient(oauth2.NoContext, &ts),
		)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		client = github.NewClient(
			oauth2.NewClient(oauth2.NoContext, &ts),
		)
	}

	gist, _, err := client.Gists.Create(context.Background(), &github.Gist{
		Files:  files,
		Public: &public,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*gist.HTMLURL)
}
