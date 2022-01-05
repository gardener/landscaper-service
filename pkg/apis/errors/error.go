// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package errors

import (
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// Error is a wrapper around the landscaper service crd error
// that implements the go error interface.
type Error struct {
	lssErr lssv1alpha1.Error
	err    error
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return fmt.Sprintf("Operation: %q Reaseon: %q - Message: %q", e.lssErr.Operation, e.lssErr.Reason, e.lssErr.Message)
}

// LandscaperServiceError returns the wrapped landscaper error.
func (e *Error) LandscaperServiceError() *lssv1alpha1.Error {
	return e.lssErr.DeepCopy()
}

// Unwrap implements the unwrap interface
func (e *Error) Unwrap() error {
	return e.err
}

// NewError creates a new landscaper service internal error
func NewError(operation, reason, message string) *Error {
	return &Error{
		lssErr: lssv1alpha1.Error{
			Operation:          operation,
			Reason:             reason,
			Message:            message,
			LastTransitionTime: metav1.Now(),
			LastUpdateTime:     metav1.Now(),
		},
	}
}

// NewWrappedError creates a new landscaper service internal error that wraps another error
func NewWrappedError(err error, operation, reason, message string) *Error {
	return &Error{
		lssErr: lssv1alpha1.Error{
			Operation:          operation,
			Reason:             reason,
			Message:            message,
			LastTransitionTime: metav1.Now(),
			LastUpdateTime:     metav1.Now(),
		},
		err: err,
	}
}

// IsError returns the innermost landscaper service error if the given error is one.
// If the err does not contain a landscaper service error, nil is returned.
func IsError(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	switch e := err.(type) {
	case *Error:
		return e, true
	default:
		uErr := errors.Unwrap(err)
		if uErr == nil {
			return nil, false
		}
		return IsError(uErr)
	}
}

// TryUpdateError tries to update the properties of the last error if the err is an internal landscaper services error.
func TryUpdateError(lastErr *lssv1alpha1.Error, err error) *lssv1alpha1.Error {
	if err == nil {
		return nil
	}
	if intErr, ok := IsError(err); ok {
		return intErr.UpdatedError(lastErr)
	}
	return nil
}

// UpdatedError updates the properties of an error.
func UpdatedError(lastError *lssv1alpha1.Error, operation, reason, message string) *lssv1alpha1.Error {
	newError := &lssv1alpha1.Error{
		Operation:          operation,
		Reason:             reason,
		Message:            message,
		LastTransitionTime: metav1.Now(),
		LastUpdateTime:     metav1.Now(),
	}

	if lastError != nil && lastError.Operation == operation {
		newError.LastTransitionTime = lastError.LastTransitionTime
	}
	return newError
}

// UpdatedError updates the properties of an existing error.
func (e Error) UpdatedError(lastError *lssv1alpha1.Error) *lssv1alpha1.Error {
	return UpdatedError(lastError, e.lssErr.Operation, e.lssErr.Reason, e.lssErr.Message)
}
