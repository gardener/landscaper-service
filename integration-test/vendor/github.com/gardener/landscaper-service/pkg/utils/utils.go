// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"strconv"

	"sigs.k8s.io/controller-runtime/pkg/client"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// StringToAnyJSON marshals a string as an AnyJSON object.
func StringToAnyJSON(s string) lsv1alpha1.AnyJSON {
	return lsv1alpha1.NewAnyJSON([]byte(fmt.Sprintf("\"%s\"", s)))
}

// BoolToAnyJSON marshals a boolean as an AnyJSON object.
func BoolToAnyJSON(b bool) lsv1alpha1.AnyJSON {
	return lsv1alpha1.NewAnyJSON([]byte(strconv.FormatBool(b)))
}

// IntToAnyJSON marshals an integer as an AnyJSON object.
func IntToAnyJSON(i int) lsv1alpha1.AnyJSON {
	return lsv1alpha1.NewAnyJSON([]byte(strconv.Itoa(i)))
}

// ContainsReference checks whether the object reference list contains the specified object reference.
func ContainsReference(refList []lssv1alpha1.ObjectReference, ref *lssv1alpha1.ObjectReference) bool {
	for _, e := range refList {
		if e.Equals(ref) {
			return true
		}
	}
	return false
}

// RemoveReference removes the given object reference from the object reference list if contained.
func RemoveReference(refList []lssv1alpha1.ObjectReference, ref *lssv1alpha1.ObjectReference) []lssv1alpha1.ObjectReference {
	for i, e := range refList {
		if e.Equals(ref) {
			refList[i] = refList[len(refList)-1]
			refList = refList[:len(refList)-1]
			break
		}
	}
	return refList
}

// HasOperationAnnotation returns true if the object has provided operation annotation set.
func HasOperationAnnotation(object client.Object, operation string) bool {
	annotations := object.GetAnnotations()
	if annotations == nil {
		return false
	}

	operationAnnotation, ok := annotations[lssv1alpha1.LandscaperServiceOperationAnnotation]
	if !ok {
		return false
	}

	return operationAnnotation == operation
}

// SetOperationAnnotation sets the provided operation annotation.
func SetOperationAnnotation(object client.Object, operation string) {
	annotations := object.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
		object.SetAnnotations(annotations)
	}
	annotations[lssv1alpha1.LandscaperServiceOperationAnnotation] = operation
}

// RemoveOperationAnnotation removes the operation annotation if it exists.
func RemoveOperationAnnotation(object client.Object) {
	annotations := object.GetAnnotations()
	if annotations != nil {
		delete(annotations, lssv1alpha1.LandscaperServiceOperationAnnotation)
	}
}
