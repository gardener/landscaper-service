// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package crdmanager

import (
	"embed"

	lsconfig "github.com/gardener/landscaper/apis/config"
	"github.com/gardener/landscaper/controller-utils/pkg/crdmanager"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/gardener/landscaper-service/pkg/apis/config"
)

const (
	embedFSCrdRootDir = "crdresources"
)

//go:embed crdresources/landscaper-service.gardener.cloud*.yaml
var importedCrdFS embed.FS

// NewCrdManager returns a new instance of the CRDManager
func NewCrdManager(log logr.Logger, mgr manager.Manager, config config.CrdManagementConfiguration) (*crdmanager.CRDManager, error) {
	crdConfig := lsconfig.CrdManagementConfiguration{
		DeployCustomResourceDefinitions: config.DeployCustomResourceDefinitions,
		ForceUpdate:                     config.ForceUpdate,
	}
	return crdmanager.NewCrdManager(log, mgr, crdConfig, &importedCrdFS, embedFSCrdRootDir)
}
