package main

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/alex-slynko/ondemander/generator"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
)

func main() {
	var manifest bosh.BoshManifest
	err := yaml.Unmarshal([]byte(os.Args[2]), &manifest)
	if err != nil {
		panic(err)
	}
	newManifest := generator.GenerateOpsfile(manifest)
	output, err := yaml.Marshal(newManifest)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))
}
