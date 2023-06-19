package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/hashicorp/go-version"
)

// public call for cmd debug
func CallRegistryApi(method string, api string, headers map[string]string, body map[string]any, rspv any) (http.Header, []byte, error) {
	return call(method, api, headers, body, rspv)
}

// public get repo detail for cmd debug
func GetRegistryRepoDetail(name string) (*Repo, error) {
	return getRepoDetail(name)
}

// https://docs.docker.com/registry/spec/api/
// call registry httpapi, scan response to rspv
func call(method string, api string, headers map[string]string, body map[string]any, rspv any) (http.Header, []byte, error) {
	// registry url
	registryURL := config.GetString("registry.url")

	// basic auth
	username := config.GetString("registry.username")
	password := config.GetString("registry.password")

	// request
	var reqbody io.Reader
	if body != nil {
		jsonbody, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		reqbody = bytes.NewBuffer(jsonbody)
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%s/v2/%s", registryURL, api), reqbody)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if username != "" {
		req.SetBasicAuth(username, password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	logger.Debugf("response header:%v\n", resp.Header)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("http status code: %d, body: %s", resp.StatusCode, string(respBody))
	}
	if rspv != nil {
		err = json.Unmarshal(respBody, rspv)
		if err != nil {
			return nil, nil, err
		}
	}
	return resp.Header, respBody, nil
}

// get all repo list
func getAllRepos() ([]*Repo, error) {
	repoes := make([]*Repo, 0, 20)
	repoNames := make([]string, 0)
	pageSize := 100
	last := ""
	var rspv map[string][]string
	for {
		// loop page get repo list
		rsh, rspb, err := call("GET", fmt.Sprintf("_catalog?n=%d&last=%s", pageSize, last), nil, nil, &rspv)
		if err != nil {
			logger.Errorf("repo list error: %s", err.Error())
			return nil, err
		}
		logger.Debugf("repo list: %s", rspb)
		repoNames = append(repoNames, rspv["repositories"]...)
		// has next page
		if rsh.Get("Link") != "" {
			last = repoNames[len(repoNames)-1]
		} else {
			break
		}
	}
	// loop get repo detail
	for _, name := range repoNames {
		repo, err := getRepoDetail(name)
		if err != nil {
			return nil, err
		}
		repoes = append(repoes, repo)
	}
	return repoes, nil
}

// get repo detail by name
func getRepoDetail(name string) (*Repo, error) {
	repo := &Repo{
		Name: name,
	}
	// repo taglist
	pageSize := 100
	last := ""
	var rspv map[string]any
	for {
		rsh, rspb, err := call("GET", fmt.Sprintf("%s/tags/list?n=%d&last=%s", repo.Name, pageSize, last), nil, nil, &rspv)
		if err != nil {
			logger.Errorf("repo %s tags list error: %s", repo.Name, err.Error())
			return nil, err
		}
		logger.Debugf("repo %s tags list: %s", repo.Name, rspb)
		for _, t := range rspv["tags"].([]any) {
			tag := &Tag{
				Name: t.(string),
			}
			repo.Tags = append(repo.Tags, tag)
		}
		// has next page
		if rsh.Get("Link") != "" {
			last = repo.Tags[len(repo.Tags)-1].Name
		} else {
			break
		}
	}

	for _, tag := range repo.Tags {
		// https://docs.docker.com/registry/spec/manifest-v2-2/#media-types
		// https://nova.moe/docker-attestation/ 会遇到MANIFEST_UNKNOWN错误，好像是buildx的问题，需加上application/vnd.oci.image.index.v1+json并过滤unknown architecture
		// get every tag manifest list: fat manifest
		_, rspb, err := call("GET", fmt.Sprintf("%s/manifests/%s", repo.Name, tag.Name), map[string]string{"Accept": "application/vnd.docker.distribution.manifest.list.v2+json,application/vnd.oci.image.index.v1+json"}, nil, &rspv)
		if err != nil {
			logger.Errorf("repo %s tag %s %s error: %s", repo.Name, tag.Name, "application/vnd.docker.distribution.manifest.list.v2+json", err.Error())
			continue
		}
		logger.Debugf("repo %s tag %s %s:\n%s", repo.Name, tag.Name, "application/vnd.docker.distribution.manifest.list.v2+json", rspb)

		// loop manifest digest
		var manifestDigests []string

		if rspv["schemaVersion"].(float64) == 1 {
			// if no manifest list, use tagname
			manifestDigests = append(manifestDigests, tag.Name)
		} else {
			for _, m := range rspv["manifests"].([]any) {
				if m.(map[string]any)["platform"].(map[string]any)["architecture"].(string) == "unknown" {
					continue
				}
				manifestDigests = append(manifestDigests, m.(map[string]any)["digest"].(string))
			}
		}
		// loop manifest detail
		for _, manifestDigest := range manifestDigests {
			_, rspb, err = call("GET", fmt.Sprintf("%s/manifests/%s", repo.Name, manifestDigest), map[string]string{"Accept": "application/vnd.docker.distribution.manifest.v2+json,application/vnd.oci.image.manifest.v1+json"}, nil, &rspv)
			if err != nil {
				logger.Errorf("repo %s tag %s digest %s %s error: %s", repo.Name, tag.Name, manifestDigest, "application/vnd.docker.distribution.manifest.v2+json", err.Error())
				continue
			}
			logger.Debugf("image %s tag %s digest %s %s:\n%s", repo.Name, tag.Name, manifestDigest, "application/vnd.docker.distribution.manifest.v2+json", rspb)
			imageDigest := rspv["config"].(map[string]any)["digest"].(string)
			// image size
			var layerSizes []float64
			var size float64
			for _, layer := range rspv["layers"].([]any) {
				layerSize := layer.(map[string]any)["size"].(float64)
				layerSizes = append(layerSizes, layerSize)
				size += layerSize
			}
			// get "application/vnd.docker.container.image.v1+json" blob detail
			_, rspb, err = call("GET", fmt.Sprintf("%s/blobs/%s", repo.Name, imageDigest), nil, nil, &rspv)
			if err != nil {
				logger.Errorf("image %s tag %s digest %s container detail error: %s", repo.Name, tag.Name, imageDigest, err.Error())
				continue
			}
			logger.Debugf("image %s tag %s digest %s container detail:\n%s", repo.Name, tag.Name, imageDigest, rspb)
			// compose image
			image := &Image{
				Digest: imageDigest,
				Size:   int64(size),
				Arch:   rspv["architecture"].(string),
				Os:     rspv["os"].(string),
			}
			// layers
			i := -1
			for _, history := range rspv["history"].([]any) {
				layer := &Layer{
					Script: history.(map[string]any)["created_by"].(string),
					Size:   0,
				}
				// not empty layer, get size from manifest
				if history.(map[string]any)["empty_layer"] != true {
					i++
					layer.Size = int64(layerSizes[i])
				}
				image.Layers = append(image.Layers, layer)
			}
			tag.Images = append(tag.Images, image)
			// created 秒数后小数点位数可能不同
			// 2021-09-26T04:15:53.296355Z
			// 2020-12-22T03:27:46.7090378Z
			// 2023-02-23T08:32:23.210335293Z
			created, _ := time.Parse(time.RFC3339Nano, rspv["created"].(string))
			tag.Created = created.Unix()
			if repo.LastUpdate < tag.Created {
				repo.LastUpdate = tag.Created
			}
			// image description
			// in Dockerfile: LABEL description="xxx"
			// image changelog
			// in Dockerfile: LABEL changelog="xxx"
			config := rspv["config"].(map[string]any)["Labels"].(map[string]any)
			if config["description"] != nil {
				repo.Desc = config["description"].(string)
			}
			if config["changelog"] != nil {
				tag.ChangeLog = config["changelog"].(string)
			}
		}
	}
	// sort tags by time
	sort.Slice(repo.Tags, func(i, j int) bool {
		// time same, latest first and sort by version name(tag name usally be version number)
		if repo.Tags[i].Created == repo.Tags[j].Created {
			if repo.Tags[i].Name == "latest" {
				return true
			} else if repo.Tags[j].Name == "latest" {
				return false
			} else {
				vi, erri := version.NewVersion(repo.Tags[i].Name)
				vj, errj := version.NewVersion(repo.Tags[j].Name)
				if erri == nil && errj == nil {
					return vi.GreaterThan(vj)
				}
				return false
			}
		}
		return repo.Tags[i].Created > repo.Tags[j].Created
	})
	js, _ := json.Marshal(repo)
	logger.Debugf("repo detail: %s", js)
	return repo, nil
}
