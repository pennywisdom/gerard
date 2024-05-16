package cmd

import (
	"log"
	"path/filepath"
)

type input struct {
	envFile   string
	workDir   string
	vars      map[string]string
	productId string
}

type svcCatCreateRepoInput struct {
	input
	repoType     string
	product      string
	businessUnit string
	division     string
	project      string
}

func (i *input) resolve(path string) string {
	basedir, err := filepath.Abs(i.workDir)
	if err != nil {
		log.Fatal(err)
	}
	if path == "" {
		return path
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(basedir, path)
	}
	return path
}

func (i *svcCatCreateRepoInput) Envfile() string {
	return i.resolve(i.envFile)
}

func (i *svcCatCreateRepoInput) WorkDir() string {
	return i.resolve("*")
}

func (i *svcCatCreateRepoInput) Vars() map[string]string {
	return i.vars
}

func (i *svcCatCreateRepoInput) ProductId() string {
	return i.productId
}
