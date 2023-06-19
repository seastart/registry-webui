package main

import (
	"fmt"

	"github.com/seastart/registry-webui/lib"
)

func main() {
	lib.Init("./config/local.yml")

	repo := &lib.Repo{
		Name: "vcs/log-center",
	}
	tag := &lib.Tag{
		Name: "latest",
	}
	var rspv map[string]any

	_, rspb, err := lib.CallRegistryApi("GET", fmt.Sprintf("%s/tags/list?n=%d&last=%s", repo.Name, 100, ""), nil, nil, &rspv)
	if err != nil {
		fmt.Printf("repo %s tags list error: %s\n", repo.Name, err.Error())
	}
	fmt.Printf("repo %s tags list: %s\n", repo.Name, rspb)

	_, rspb, err = lib.CallRegistryApi("GET", fmt.Sprintf("%s/manifests/%s", repo.Name, tag.Name), map[string]string{"Accept": "application/vnd.docker.distribution.manifest.v2+json,application/vnd.oci.image.index.v1+json"}, nil, &rspv)
	if err != nil {
		fmt.Printf("repo %s tag %s %s error: %s\n", repo.Name, tag.Name, "application/vnd.docker.distribution.manifest.list.v2+json", err.Error())
		_, rspb, err = lib.CallRegistryApi("GET", fmt.Sprintf("%s/manifests/%s", repo.Name, tag.Name), map[string]string{"Accept": "application/vnd.docker.distribution.manifest.list.v2+json,application/vnd.docker.distribution.manifest.v2+json,application/vnd.docker.distribution.manifest.v1+json,application/vnd.oci.image.manifest.v1+json,application/vnd.oci.image.index.v1+json"}, nil, &rspv)
		if err != nil {
			fmt.Printf("repo %s tag %s %s error: %s\n", repo.Name, tag.Name, "application/vnd.docker.distribution.manifest.v1+json", err.Error())
		} else {
			fmt.Printf("repo %s tag %s %s:\n%s\n", repo.Name, tag.Name, "application/vnd.docker.distribution.manifest.v1+json", rspb)
		}
	} else {
		fmt.Printf("repo %s tag %s %s:\n%s\n", repo.Name, tag.Name, "application/vnd.docker.distribution.manifest.list.v2+json", rspb)
	}

	manifestDigest := "sha256:43f9f365da57f247c16cbbd359ee53b1121ddfb04a2baa52098b67b678e648c5"
	_, rspb, err = lib.CallRegistryApi("GET", fmt.Sprintf("%s/manifests/%s", repo.Name, manifestDigest), map[string]string{"Accept": "application/vnd.docker.distribution.manifest.v2+json,application/vnd.oci.image.manifest.v1+json"}, nil, &rspv)
	if err != nil {
		fmt.Printf("repo %s tag %s digest %s %s error: %s\n", repo.Name, tag.Name, manifestDigest, "application/vnd.docker.distribution.manifest.v2+json", err.Error())
	} else {
		fmt.Printf("image %s tag %s digest %s %s:\n%s", repo.Name, tag.Name, manifestDigest, "application/vnd.docker.distribution.manifest.v2+json", rspb)
	}

	imageDigest := "sha256:fc1775f1934e686b7008ee38e88fba1a9d855b876dc8652a9137d25142aa9f70"
	_, rspb, err = lib.CallRegistryApi("GET", fmt.Sprintf("%s/blobs/%s", repo.Name, imageDigest), nil, nil, &rspv)
	if err != nil {
		fmt.Printf("image %s tag %s digest %s container detail error: %s\n", repo.Name, tag.Name, imageDigest, err.Error())
	} else {
		fmt.Printf("image %s tag %s digest %s container detail:\n%s\n", repo.Name, tag.Name, imageDigest, rspb)
	}

}
