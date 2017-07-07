package generator

import (
	"fmt"

	"github.com/alex-slynko/ondemander/types"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
)

func GenerateOpsfile(originalManifest bosh.BoshManifest) []types.Operation {
	resultOpsFile := []types.Operation{}
	for _, stemcell := range originalManifest.Stemcells {
		if stemcell.OS != "ubuntu-trusty" {
			resultOpsFile = append(resultOpsFile, generateOperationForStemcell(stemcell))
		}
	}
	for _, release := range originalManifest.Releases {
		resultOpsFile = append(resultOpsFile, generateOperationForRelease(release))
	}

	resultOpsFile = append(resultOpsFile, buildServiceDeploymentProperties(originalManifest))
	return resultOpsFile
}

func generateOperationForRelease(release bosh.Release) types.Operation {
	v := make(map[string]interface{})
	v["name"] = release.Name
	v["version"] = fmt.Sprintf("((%s-version))", release.Name)
	op := types.Operation{
		Type:  "replace",
		Path:  "/releases/-",
		Value: v,
	}

	return op
}

func generateOperationForStemcell(stemcell bosh.Stemcell) types.Operation {
	v := make(map[string]interface{})
	v["alias"] = stemcell.Alias
	v["os"] = stemcell.OS
	v["version"] = stemcell.Version
	op := types.Operation{
		Type:  "replace",
		Path:  "/stemcells/-",
		Value: v,
	}

	return op
}

func buildServiceDeploymentProperties(manifest bosh.BoshManifest) types.Operation {
	serviceDeployment := map[string][]string{}
	for _, instanceGroup := range manifest.InstanceGroups {
		for _, job := range instanceGroup.Jobs {
			if serviceDeployment[job.Release] == nil {
				serviceDeployment[job.Release] = []string{}
			}
			if !contains(job.Name, serviceDeployment[job.Release]) {
				serviceDeployment[job.Release] = append(serviceDeployment[job.Release], job.Name)
			}
		}
	}
	return types.Operation{}
}

func contains(job string, jobs []string) bool {
	for _, addedJob := range jobs {
		if addedJob == job {
			return true
		}
	}
	return false
}
