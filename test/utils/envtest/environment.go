// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package envtest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"

	"html/template"

	kutil "github.com/gardener/landscaper/controller-utils/pkg/kubernetes"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	lssv1alpha1 "github.com/gardener/landscaper-service/pkg/apis/core/v1alpha1"
)

// Environment is the test environment
type Environment struct {
	// Env is the kubernetes envtest test environment.
	Env *envtest.Environment
	// Client is the kubernetes client
	Client client.Client
}

// NewEnvironment creates a new test environment.
func NewEnvironment(projectRoot string) (*Environment, error) {
	projectRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return nil, err
	}
	testBinPath := filepath.Join(projectRoot, "tmp", "test", "bin")

	// if the test binary path doesn't exist, the envtest default values will be used.
	if _, err := os.Stat(testBinPath); err == nil {
		if err := os.Setenv("TEST_ASSET_KUBE_APISERVER", filepath.Join(testBinPath, "kube-apiserver")); err != nil {
			return nil, err
		}
		if err := os.Setenv("TEST_ASSET_ETCD", filepath.Join(testBinPath, "etcd")); err != nil {
			return nil, err
		}
		if err := os.Setenv("TEST_ASSET_KUBECTL", filepath.Join(testBinPath, "kubectl")); err != nil {
			return nil, err
		}
	}

	return &Environment{
		Env: &envtest.Environment{
			CRDDirectoryPaths: []string{
				filepath.Join(projectRoot, "pkg", "crdmanager", "crdresources"),
				filepath.Join(projectRoot, "tmp", "landscapercrd"),
			},
		},
	}, nil
}

// Start starts the test environment.
func (e *Environment) Start() (client.Client, error) {
	restConfig, err := e.Env.Start()
	if err != nil {
		return nil, err
	}

	fakeClient, err := client.New(restConfig, client.Options{Scheme: LandscaperServiceScheme})
	if err != nil {
		return nil, err
	}

	e.Client = fakeClient
	return fakeClient, nil
}

// Stop stops the test environment
func (e *Environment) Stop() error {
	return e.Env.Stop()
}

// InitResources initializes the test environment with the resources stored in the given path.
func (e *Environment) InitResources(ctx context.Context, resourcesPath string) (*State, error) {
	var state *State
	ns := &corev1.Namespace{}
	ns.GenerateName = "tests-"
	if err := e.Client.Create(ctx, ns); err != nil {
		return nil, err
	}
	state = NewState(ns.GetName())

	resources, err := e.parseResources(resourcesPath, state)
	if err != nil {
		return nil, err
	}

	resourcesChan := make(chan client.Object, len(resources))

	for _, obj := range resources {
		select {
		case resourcesChan <- obj:
		default:
		}
	}

	injectOwnerUUIDs := func(obj client.Object) error {
		refs := obj.GetOwnerReferences()
		for i, ownerRef := range obj.GetOwnerReferences() {
			uObj := &unstructured.Unstructured{}
			uObj.SetAPIVersion(ownerRef.APIVersion)
			uObj.SetKind(ownerRef.Kind)
			uObj.SetName(ownerRef.Name)
			uObj.SetNamespace(obj.GetNamespace())
			if err := e.Client.Get(ctx, kutil.ObjectKeyFromObject(uObj), uObj); err != nil {
				return fmt.Errorf("no owner found for %s\n", kutil.ObjectKeyFromObject(obj).String())
			}
			refs[i].UID = uObj.GetUID()
		}
		obj.SetOwnerReferences(refs)
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	for obj := range resourcesChan {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("context canceled; check resources as there might be a cyclic dependency")
		}
		objName := kutil.ObjectKeyFromObject(obj).String()
		// create namespaces if not exist before
		if len(obj.GetNamespace()) != 0 {
			ns := &corev1.Namespace{}
			ns.Name = obj.GetNamespace()
			if _, err := controllerutil.CreateOrUpdate(ctx, e.Client, ns, func() error {
				return nil
			}); err != nil {
				return nil, err
			}
		}
		// inject real uuids if possible
		if len(obj.GetOwnerReferences()) != 0 {
			if err := injectOwnerUUIDs(obj); err != nil {
				// try to requeue
				// todo: somehow detect cyclic dependencies (maybe just use a context with an timeout)
				resourcesChan <- obj
				continue
			}
		}
		if err := e.createObject(ctx, obj, state); err != nil {
			return nil, fmt.Errorf("unable to create %s %s: %w", objName, obj.GetObjectKind().GroupVersionKind().String(), err)
		}
		if len(resourcesChan) == 0 {
			close(resourcesChan)
		}
	}

	return state, nil
}

// CleanupResources removes all resources that have been created by function InitResources.
func (e *Environment) CleanupResources(ctx context.Context, state *State) error {
	for _, obj := range state.Deployments {
		if err := e.deleteObject(ctx, obj); err != nil {
			return err
		}
	}
	for _, obj := range state.Instances {
		if err := e.deleteObject(ctx, obj); err != nil {
			return err
		}
	}
	for _, obj := range state.Configs {
		if err := e.deleteObject(ctx, obj); err != nil {
			return err
		}
	}
	for _, obj := range state.Secrets {
		if err := e.deleteObject(ctx, obj); err != nil {
			return err
		}
	}
	return nil
}

// WaitForObjectToBeDeleted waits for the given object to be deleted or the given timout duration to end.
func (e *Environment) WaitForObjectToBeDeleted(ctx context.Context, obj client.Object, timeout time.Duration) error {
	var (
		lastErr error
		uObj    client.Object
	)
	err := wait.PollImmediate(2*time.Second, timeout, func() (done bool, err error) {
		uObj = obj.DeepCopyObject().(client.Object)
		if err := e.Client.Get(ctx, client.ObjectKeyFromObject(obj), uObj); err != nil {
			if apierrors.IsNotFound(err) {
				return true, nil
			}
			lastErr = err
			return false, nil
		}
		return false, nil
	})
	if err != nil {
		if lastErr != nil {
			return lastErr
		}
		// try to print the whole object to debug
		d, err2 := json.Marshal(uObj)
		if err2 != nil {
			return err
		}
		return fmt.Errorf("deletion timeout: %s", string(d))
	}
	return nil
}
func (e *Environment) createObject(ctx context.Context, obj client.Object, state *State) error {
	tmp := obj.DeepCopyObject().(client.Object)
	if err := e.Client.Create(ctx, obj); err != nil {
		return err
	}

	tmp.SetName(obj.GetName())
	tmp.SetNamespace(obj.GetNamespace())
	tmp.SetResourceVersion(obj.GetResourceVersion())
	tmp.SetGeneration(obj.GetGeneration())
	tmp.SetUID(obj.GetUID())
	tmp.SetCreationTimestamp(obj.GetCreationTimestamp())

	if err := e.Client.Status().Update(ctx, tmp); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
	}

	state.AddObject(tmp)
	return nil
}

func (e *Environment) deleteObject(ctx context.Context, obj client.Object) error {
	if err := e.Client.Get(ctx, client.ObjectKeyFromObject(obj), obj); err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	if obj.GetDeletionTimestamp().IsZero() {
		if err := e.Client.Delete(ctx, obj); err != nil {
			if apierrors.IsNotFound(err) {
				return nil
			}
			return err
		}
	}
	if err := e.WaitForObjectToBeDeleted(ctx, obj, 5*time.Second); err != nil {
		if err := e.removeFinalizer(ctx, obj); err != nil {
			return err
		}
	}

	return nil
}

func (e *Environment) removeFinalizer(ctx context.Context, obj client.Object) error {
	if len(obj.GetFinalizers()) == 0 {
		return nil
	}
	if err := e.Client.Get(ctx, kutil.ObjectKey(obj.GetName(), obj.GetNamespace()), obj); err != nil {
		return err
	}
	currObj := obj.DeepCopyObject().(client.Object)

	obj.SetFinalizers([]string{})
	patch := client.MergeFrom(currObj)
	if err := e.Client.Patch(ctx, obj, patch); err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("unable to remove finalizer from object: %w", err)
	}
	return nil
}

func (e *Environment) parseResources(path string, state *State) ([]client.Object, error) {
	objects := make([]client.Object, 0)
	errOuter := filepath.Walk(path, func(path string, info os.FileInfo, walkerr error) error {
		if walkerr != nil {
			return walkerr
		}
		if info.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "unable to read file %s", path)
		}

		// template files
		tmpl, err := template.New("init").Parse(string(data))
		if err != nil {
			return err
		}
		buf := bytes.NewBuffer([]byte{})
		if err := tmpl.Execute(buf, map[string]string{"Namespace": state.Namespace}); err != nil {
			return err
		}

		var (
			decoder    = yaml.NewYAMLOrJSONDecoder(buf, 1024)
			decodedObj json.RawMessage
		)

		for {
			if err := decoder.Decode(&decodedObj); err != nil {
				if err == io.EOF {
					return nil
				}
				continue
			}

			if decodedObj != nil {
				objects, err = e.decodeAndAppendLSSObject(decodedObj, objects)
				if err != nil {
					return errors.Wrapf(err, "unable to decode file %s", path)
				}
			}
		}
	})
	if errOuter != nil {
		return nil, errOuter
	}

	return objects, nil
}

func (e *Environment) decodeAndAppendLSSObject(data []byte, objects []client.Object) ([]client.Object, error) {
	decoder := serializer.NewCodecFactory(LandscaperServiceScheme).UniversalDecoder()

	_, gvk, err := decoder.Decode(data, nil, &unstructured.Unstructured{})
	if err != nil {
		return nil, fmt.Errorf("unable to decode object into unstructured: %w", err)
	}

	switch gvk.Kind {
	case DeploymentGVK.Kind:
		deployment := &lssv1alpha1.LandscaperDeployment{}
		if _, _, err := decoder.Decode(data, nil, deployment); err != nil {
			return nil, fmt.Errorf("unable to decode file as landscaper deployment: %w", err)
		}
		return append(objects, deployment), nil
	case InstanceGVK.Kind:
		instance := &lssv1alpha1.Instance{}
		if _, _, err := decoder.Decode(data, nil, instance); err != nil {
			return nil, fmt.Errorf("unable to decode file as instance: %w", err)
		}
		return append(objects, instance), nil
	case ConfigGVK.Kind:
		config := &lssv1alpha1.ServiceTargetConfig{}
		if _, _, err := decoder.Decode(data, nil, config); err != nil {
			return nil, fmt.Errorf("unable to decode file as service target config: %w", err)
		}
		return append(objects, config), nil
	case SecretGVK.Kind:
		secret := &corev1.Secret{}
		if _, _, err := decoder.Decode(data, nil, secret); err != nil {
			return nil, fmt.Errorf("unable to decode file as secret: %w", err)
		}
		return append(objects, secret), nil
	case ConfigMapGVK.Kind:
		configMap := &corev1.ConfigMap{}
		if _, _, err := decoder.Decode(data, nil, configMap); err != nil {
			return nil, fmt.Errorf("unable to decode file as secret: %w", err)
		}
		return append(objects, configMap), nil
	case InstallationGVK.Kind:
		installation := &lsv1alpha1.Installation{}
		if _, _, err := decoder.Decode(data, nil, installation); err != nil {
			return nil, fmt.Errorf("unable to decode file as secret: %w", err)
		}
		return append(objects, installation), nil
	case TargetGVK.Kind:
		target := &lsv1alpha1.Target{}
		if _, _, err := decoder.Decode(data, nil, target); err != nil {
			return nil, fmt.Errorf("unable to decode file as secret: %w", err)
		}
		return append(objects, target), nil
	case ContextGVK.Kind:
		context := &lsv1alpha1.Context{}
		if _, _, err := decoder.Decode(data, nil, context); err != nil {
			return nil, fmt.Errorf("unable to decode file as context: %w", err)
		}
		return append(objects, context), nil
	case AvailabilityCollectionGVK.Kind:
		availabilityCollection := &lssv1alpha1.AvailabilityCollection{}
		if _, _, err := decoder.Decode(data, nil, availabilityCollection); err != nil {
			return nil, fmt.Errorf("unable to decode file as availability collection: %w", err)
		}
		return append(objects, availabilityCollection), nil
	case LsHealthCheckGVK.Kind:
		lshealthcheck := &lsv1alpha1.LsHealthCheck{}
		if _, _, err := decoder.Decode(data, nil, lshealthcheck); err != nil {
			return nil, fmt.Errorf("unable to decode file as  lshealthcheck: %w", err)
		}
		return append(objects, lshealthcheck), nil

	default:
		return objects, nil
	}
}
