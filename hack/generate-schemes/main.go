// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gardener/landscaper/apis/hack/generate-schemes/app"
	lsschema "github.com/gardener/landscaper/apis/schema"

	corelsv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"

	"github.com/gardener/landscaper-service/pkg/apis/openapi"
)

var Exports = []string{
	"LandscaperDeployment",
	"Instance",
	"ServiceTargetConfig",
	"NamespaceRegistration",
	"SubjectList",
}

var CRDs = []lsschema.CustomResourceDefinitions{
	corelsv1alpha1.ResourceDefinition,
}

var (
	schemaDir string
	crdDir    string
)

func init() {
	flag.StringVar(&schemaDir, "schema-dir", "", "output directory for jsonschemas")
	flag.StringVar(&crdDir, "crd-dir", "", "output directory for crds")
}

func main() {
	flag.Parse()
	if len(schemaDir) == 0 {
		log.Fatalln("expected --schema-dir to be set")
	}
	schemaGenerator := app.NewSchemaGenerator(Exports, CRDs, openapi.GetOpenAPIDefinitions)
	if err := schemaGenerator.Run(schemaDir, crdDir); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
