/*
Copyright 2020 The Kubernetes Authors.

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

package server

import (
	"fmt"

	"github.com/alibaba/open-local/pkg/csi/lib"
	serverhelpers "github.com/google/go-microservice-helpers/server"
	log "github.com/sirupsen/logrus"
)

var (
	lvmdPort string
)

// Start start lvmd
func Start(port string) {
	lvmdPort = port
	address := fmt.Sprintf(":%s", port)
	log.Infof("Lvmd Starting with socket: %s ...", address)

	svr := NewServer()
	serverhelpers.ListenAddress = &address
	grpcServer, _, err := serverhelpers.NewServer()
	if err != nil {
		log.Errorf("failed to init GRPC server: %v", err)
		return
	}

	lib.RegisterLVMServer(grpcServer, &svr)

	err = serverhelpers.ListenAndServe(grpcServer, nil)
	if err != nil {
		log.Errorf("failed to serve: %v", err)
		return
	}
	log.Infof("Lvmd End ...")
}

// GetLvmdPort get lvmd port
func GetLvmdPort() string {
	return lvmdPort
}
