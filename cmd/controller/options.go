/*
Copyright © 2021 Alibaba Group Holding Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"strings"

	"github.com/alibaba/open-local/pkg/controller"
	"github.com/spf13/pflag"
	cliflag "k8s.io/component-base/cli/flag"
)

type controllerOption struct {
	Master       string
	Kubeconfig   string
	InitConfig   string
	FeatureGates map[string]bool
}

func (option *controllerOption) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&option.Kubeconfig, "kubeconfig", option.Kubeconfig, "Path to the kubeconfig file to use.")
	fs.StringVar(&option.Master, "master", option.Master, "URL/IP for master.")
	fs.StringVar(&option.InitConfig, "initconfig", "open-local", "initconfig is NodeLocalStorageInitConfig(CRD) for controller to create NodeLocalStorage")
	fs.Var(cliflag.NewMapStringBool(&option.FeatureGates), "feature-gates", "A set of key=value pairs that describe feature gates for alpha/experimental features. "+
		"Options are:\n"+strings.Join(controller.DefaultFeatureGate.KnownFeatures(), "\n"))
}
