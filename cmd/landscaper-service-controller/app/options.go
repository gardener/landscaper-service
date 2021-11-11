// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	goflag "flag"
	"io/ioutil"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"github.com/gardener/landscaper-service/pkg/apis/config"
	configinstall "github.com/gardener/landscaper-service/pkg/apis/config/install"
	"github.com/gardener/landscaper-service/pkg/apis/config/v1alpha1"

	"github.com/gardener/landscaper/controller-utils/pkg/logger"
	"github.com/go-logr/logr"
	flag "github.com/spf13/pflag"
	ctrl "sigs.k8s.io/controller-runtime"
)

// options holds the landscaper service controller options
type options struct {
	Log        logr.Logger // Log is the logger instance
	ConfigPath string      // ConfigPath is the path to the configuration file

	Config *config.LandscaperServiceConfiguration // Config is the parsed configuration
}

// NewOptions returns a new options instance
func NewOptions() *options {
	return &options{}
}

// AddFlags adds flags passed via command line
func (o *options) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.ConfigPath, "config", "", "Specify the path to the configuration file")
	logger.InitFlags(fs)
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

// Complete initializes the options instance and validates flags
func (o *options) Complete(ctx context.Context) error {
	log, err := logger.New(nil)
	if err != nil {
		return err
	}
	o.Log = log.WithName("setup")
	logger.SetLogger(log)
	ctrl.SetLogger(log)

	o.Config, err = o.parseConfigurationFile(ctx)
	if err != nil {
		return err
	}

	err = o.validate()
	return err
}

func (o *options) parseConfigurationFile(ctx context.Context) (*config.LandscaperServiceConfiguration, error) {
	configScheme := runtime.NewScheme()
	configinstall.Install(configScheme)
	decoder := serializer.NewCodecFactory(configScheme).UniversalDecoder()

	configv1alpha1 := &v1alpha1.LandscaperServiceConfiguration{}

	if len(o.ConfigPath) != 0 {
		data, err := ioutil.ReadFile(o.ConfigPath)
		if err != nil {
			return nil, err
		}

		if _, _, err := decoder.Decode(data, nil, configv1alpha1); err != nil {
			return nil, err
		}
	}

	configScheme.Default(configv1alpha1)

	config := &config.LandscaperServiceConfiguration{}
	err := configScheme.Convert(configv1alpha1, config, ctx)
	if err != nil {
		return nil, err
	}
	configScheme.Default(config)

	return config, nil
}

func (o *options) validate() error {
	return nil
}
