// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type ManifestMap map[string]Manifest

type ManifestResponse struct {
	Manifests ManifestMap `json:"manifest"`
}

type Manifest struct {
	Tags []string `json:"tag"`
}

func main() {
	if len(os.Args) < 4 {
		panic(fmt.Errorf("usage: %s REGISTRY PATH VERSION", os.Args[0]))
	}

	registry := os.Args[1]
	path := os.Args[2]
	version := os.Args[3]

	res, err := http.Get(fmt.Sprintf("https://%s/v2/%s/tags/list", registry, path))

	if err != nil {
		panic(fmt.Errorf("request failed: %w", err))
	}

	if res.StatusCode != http.StatusOK {
		panic(fmt.Errorf("unexpected response with code: %d", res.StatusCode))
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(fmt.Errorf("failed to read response body: %w", err))
	}

	var manifestsResponse ManifestResponse
	if err := json.Unmarshal(resBody, &manifestsResponse); err != nil {
		panic(fmt.Errorf("respponse body unmarshal failed: %w", err))
	}

	for sha256, manifest := range manifestsResponse.Manifests {
		for _, tag := range manifest.Tags {
			if tag == version {
				fmt.Printf("%s", sha256)
				os.Exit(0)
			}
		}
	}

	fmt.Println("not found")
	os.Exit(1)
}
