// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

//go:generate ../../hack/generate-code.sh
//go:generate go run -mod=vendor ../../hack/generate-schemes --schema-dir ./.schemes --crd-dir ../crdmanager/crdresources

package apis
