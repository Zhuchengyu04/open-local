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

package common

// Configuration stores all the user-defined parameters to the controller
type Configuration struct {
	// Nodename is the kube node name
	Nodename string
	// SysPath is the the mount point of the host sys path
	SysPath string
	// MountPath defines the specified mount path we discover
	MountPath string
	// DisconverInterval is the duration(second) that the agent checks at one time
	DiscoverInterval int
	// LogicalVolumeNamePrefix is the prefix of LogicalVolume Name
	LogicalVolumeNamePrefix string
	// RegExp is used to filter device names
	RegExp string
}

const (
	// DefaultConfigPath is the default configfile path of open-local agent
	DefaultConfigPath string = "/etc/controller/config/"
	// DefaultInterval is the duration(second) that the agent checks at one time
	DefaultInterval int    = 60
	DefaultEndpoint string = "unix://tmp/csi.sock"
)
