// SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package avuploader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	coreconfig "github.com/gardener/landscaper-service/pkg/apis/config"
	"github.com/gardener/landscaper-service/pkg/operation"
)

type Controller struct {
	operation.Operation
}

func NewController(log logr.Logger, c client.Client, scheme *runtime.Scheme, config *coreconfig.LandscaperServiceConfiguration) (reconcile.Reconciler, error) {
	ctrl := &Controller{}
	op := operation.NewOperation(log, c, scheme, config)
	ctrl.Operation = *op
	return ctrl, nil
}

func (c *Controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	log := c.Log().WithValues("availabilityCollection", req.NamespacedName.String())
	ctx = logr.NewContext(ctx, log)
	log.V(5).Info("reconcile", "availabilityCollection", req.NamespacedName)

	if c.Config().AvailabilityMonitoring.AvailabilityServiceConfiguration.Url == "" || c.Config().AvailabilityMonitoring.AvailabilityServiceConfiguration.ApiKey == "" {
		log.V(5).Info("av service not configured")
		return reconcile.Result{}, nil
	}

	//get availabilityCollection
	availabilityCollection := &lssv1alpha1.AvailabilityCollection{}
	if err := c.Client().Get(ctx, req.NamespacedName, availabilityCollection); err != nil {
		if apierrors.IsNotFound(err) {
			c.Log().V(5).Info(err.Error())
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	//do not run on spec updates or when status was already uploaded
	if availabilityCollection.ObjectMeta.Generation != availabilityCollection.Status.ObservedGeneration || availabilityCollection.Status.LastRun == availabilityCollection.Status.LastReported {
		return reconcile.Result{}, nil
	}

	request := constructAvsRequest(*availabilityCollection)

	err := doAvsRequest(request, c.Config().AvailabilityMonitoring.AvailabilityServiceConfiguration.Url, c.Config().AvailabilityMonitoring.AvailabilityServiceConfiguration.ApiKey)
	if err != nil {
		return reconcile.Result{}, err
	}

	availabilityCollection.Status.LastReported = availabilityCollection.Status.LastRun

	//write to status
	if err := c.Client().Status().Update(ctx, availabilityCollection); err != nil {
		return reconcile.Result{}, fmt.Errorf("unable to update availability status: %w", err)
	}

	return reconcile.Result{}, nil

}

func constructAvsRequest(availabilityCollection lssv1alpha1.AvailabilityCollection) AvsRequest {
	//Fill failedInstances with all failed. A failed instance will create an instance outage. If this instance is not in the array anymore, instance outage is resolved
	// Overall status is derived if len(failedInstances) > 0
	failedInstances := []AvsInstance{}
	for _, instanceStatus := range availabilityCollection.Status.Instances {
		if instanceStatus.Status == string(lsv1alpha1.LsHealthCheckStatusFailed) {
			avsInstance := AvsInstance{
				InstanceId:   instanceStatus.Name,
				Name:         instanceStatus.Name,
				Status:       AVS_STATUS_DOWN,
				OutageReason: instanceStatus.FailedReason,
			}
			failedInstances = append(failedInstances, avsInstance)
		}
	}
	if availabilityCollection.Status.Self.Status == string(lsv1alpha1.LsHealthCheckStatusFailed) {
		avsInstance := AvsInstance{
			InstanceId:   "Self",
			Name:         "Self",
			Status:       AVS_STATUS_DOWN,
			OutageReason: availabilityCollection.Status.Self.FailedReason,
		}
		failedInstances = append(failedInstances, avsInstance)
	}

	status := AVS_STATUS_UP
	outageReason := ""
	if len(failedInstances) > 0 {
		status = AVS_STATUS_DOWN
		//include self landscaper (--> +1 )
		totalNumberOfMonitoredLandscapers := len(availabilityCollection.Status.Instances) + 1
		outageReason = fmt.Sprintf("%d/%d monitored landscaper down", len(failedInstances), totalNumberOfMonitoredLandscapers)
	}

	request := AvsRequest{
		Timestamp:          availabilityCollection.Status.LastRun.Unix(),
		ResponseTime:       0,
		ResponseStatusCode: 200,
		Instances:          failedInstances,
		Status:             status,
		OutageReason:       outageReason,
	}
	return request
}

func doAvsRequest(request AvsRequest, url string, apiKey string) error {
	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("avs upload post payload json marhsal failed: %w", err)
	}
	client := &http.Client{}
	avsreq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("avs upload request build failed: %w", err)
	}
	avsreq.Header.Add("Api-Key", apiKey)
	avsreq.Header.Add("Content-Type", "application/json")
	avsreq.Header.Add("accept", "application/json")
	resp, err := client.Do(avsreq)
	if err != nil {
		return fmt.Errorf("avs upload client build failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		resBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("avs upload failed read response: %w", err)
		}
		return fmt.Errorf("avs upload request failed with response code %d: %s", resp.StatusCode, resBody)
	}
	return nil
}

// AvsRequest contains the structure for a request to avs to push availability information.
type AvsRequest struct {
	// Timestamp is the timestamp the data was collected. Must be newer than the one in the last request.
	Timestamp int64 `json:"timestamp"`
	// Response time is the time the availability check took.
	ResponseTime int `json:"responseTime"`
	// ResponseStatusCode is the http response code received from the service that is av monitored.
	ResponseStatusCode int `json:"responseStatusCode"`
	// OutageReason is the reason the av monitored service is unavailable.
	OutageReason string `json:"outageReason"`
	// Status is the availability status for the service: 0 = DOWN, 1 = UP.
	Status int `json:"status"`
	// Instances allow to give an av status for multiple instances. If one of the instances is DOWN, the overall status is DOWN.
	Instances []AvsInstance `json:"instances"`
}

// AvsInstance is a status of a single monitored instance.
type AvsInstance struct {
	// InstanceId is the id of the monitored instance.
	InstanceId string `json:"instanceId"`
	// Name is a name for the instance.
	Name string `json:"name"`
	// OutageReason is the reason this instance is unavailable.
	OutageReason string `json:"outageReason"`
	// Status is the availability status for the instance: 0 = DOWN, 1 = UP.
	Status int `json:"status"`
}

const (
	AVS_STATUS_DOWN = 0
	AVS_STATUS_UP   = 1
)
