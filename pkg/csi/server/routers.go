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
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server lvm grpc server
type Server struct {
	lib.UnimplementedLVMServer
}

// NewServer new server
func NewServer() Server {
	return Server{}
}

// ListLV list lvm volume
func (s Server) ListLV(ctx context.Context, in *lib.ListLVRequest) (*lib.ListLVReply, error) {
	log.Debugf("List LVM for vg: %s", in.VolumeGroup)
	lvs, err := ListLV(in.VolumeGroup)
	if err != nil {
		log.Errorf("List LVM with error: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to list LVs: %v", err)
	}

	pblvs := make([]*lib.LogicalVolume, len(lvs))
	for i, v := range lvs {
		pblvs[i] = v.ToProto()
	}
	log.Debugf("List LVM Successful with result: %+v", pblvs)
	return &lib.ListLVReply{Volumes: pblvs}, nil
}

// CreateLV create lvm volume
func (s Server) CreateLV(ctx context.Context, in *lib.CreateLVRequest) (*lib.CreateLVReply, error) {
	log.Debugf("Create LVM with: %+v", in)
	out, err := CreateLV(ctx, in.VolumeGroup, in.Name, in.Size, in.Mirrors, in.Tags, in.Striping)
	if err != nil {
		log.Errorf("Create LVM with error: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to create lv: %v", err)
	}
	log.Debugf("Create LVM Successful with result: %+v", out)
	return &lib.CreateLVReply{CommandOutput: out}, nil
}

// RemoveLV remove lvm volume
func (s Server) RemoveLV(ctx context.Context, in *lib.RemoveLVRequest) (*lib.RemoveLVReply, error) {
	log.Debugf("Remove LVM with: %+v", in)
	out, err := RemoveLV(ctx, in.VolumeGroup, in.Name)
	if err != nil {
		log.Errorf("Remove LVM with error: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to remove lv: %v", err)
	}
	log.Debugf("Remove LVM Successful with result: %+v", out)
	return &lib.RemoveLVReply{CommandOutput: out}, nil
}

// CloneLV clone lvm volume
func (s Server) CloneLV(ctx context.Context, in *lib.CloneLVRequest) (*lib.CloneLVReply, error) {
	out, err := CloneLV(ctx, in.SourceName, in.DestName)
	if err != nil {
		log.Errorf("Clone LVM with error: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to clone lv: %v", err)
	}
	log.Debugf("Clone LVM with result: %+v", out)
	return &lib.CloneLVReply{CommandOutput: out}, nil
}

// ExpandLV expand lvm volume
func (s Server) ExpandLV(ctx context.Context, in *lib.ExpandLVRequest) (*lib.ExpandLVReply, error) {
	out, err := ExpandLV(ctx, in.VolumeGroup, in.Name, in.Size)
	if err != nil {
		log.Errorf("Expand LVM with error: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to expand lv: %v", err)
	}
	log.Debugf("Expand LVM with result: %+v", out)
	return &lib.ExpandLVReply{CommandOutput: out}, nil
}

// ListVG list volume group
func (s Server) ListVG(ctx context.Context, in *lib.ListVGRequest) (*lib.ListVGReply, error) {
	vgs, err := ListVG()
	if err != nil {
		log.Errorf("List VG with error: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to list VGs: %v", err)
	}

	pbvgs := make([]*lib.VolumeGroup, len(vgs))
	for i, v := range vgs {
		pbvgs[i] = v.ToProto()
	}
	log.Debugf("List VG with result: %+v", pbvgs)
	return &lib.ListVGReply{VolumeGroups: pbvgs}, nil
}

// CreateSnapshot create lvm snapshot
func (s Server) CreateSnapshot(ctx context.Context, in *lib.CreateSnapshotRequest) (*lib.CreateSnapshotReply, error) {
	log.Debugf("Create LVM Snapshot with: %+v", in)
	out, err := CreateSnapshot(ctx, in.VolumeGroup, in.SnapName, in.LvName, in.Size)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateSnapshot: create snapshot with error: %s", err.Error())
	}
	log.Debugf("Create LVM Snapshot Successful with result: %+v", out)
	return &lib.CreateSnapshotReply{CommandOutput: out}, nil
}

// RemoveSnapshot remove lvm snapshot
func (s Server) RemoveSnapshot(ctx context.Context, in *lib.RemoveSnapshotRequest) (*lib.RemoveSnapshotReply, error) {
	log.Debugf("Remove LVM Snapshot with: %+v", in)
	out, err := RemoveSnapshot(ctx, in.VolumeGroup, in.SnapName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "RemoveSnapshot: remove snapshot with error: %s", err.Error())
	}
	log.Debugf("Remove LVM Snapshot Successful with result: %+v", out)
	return &lib.RemoveSnapshotReply{CommandOutput: out}, nil
}

// CreateVG create volume group
func (s Server) CreateVG(ctx context.Context, in *lib.CreateVGRequest) (*lib.CreateVGReply, error) {
	out, err := CreateVG(ctx, in.Name, in.PhysicalVolume, in.Tags)
	if err != nil {
		log.Errorf("Create VG with error: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to create vg: %v", err)
	}
	log.Debugf("Create VG with result: %+v", out)
	return &lib.CreateVGReply{CommandOutput: out}, nil
}

// RemoveVG remove volume group
func (s Server) RemoveVG(ctx context.Context, in *lib.CreateVGRequest) (*lib.RemoveVGReply, error) {
	out, err := RemoveVG(ctx, in.Name)
	if err != nil {
		log.Errorf("Remove VG with error: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to remove vg: %v", err)
	}
	log.Debugf("Remove VG with result: %+v", out)
	return &lib.RemoveVGReply{CommandOutput: out}, nil
}

// CleanPath remove file under path
func (s Server) CleanPath(ctx context.Context, in *lib.CleanPathRequest) (*lib.CleanPathReply, error) {
	err := CleanPath(ctx, in.Path)
	if err != nil {
		log.Errorf("CleanPath with error: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to remove vg: %v", err)
	}
	log.Debugf("CleanPath with result Successful")
	return &lib.CleanPathReply{CommandOutput: "Successful remove path: " + in.Path}, nil
}

// CleanDevice wipefs
func (s Server) CleanDevice(ctx context.Context, in *lib.CleanDeviceRequest) (*lib.CleanDeviceReply, error) {
	out, err := CleanDevice(ctx, in.Device)
	if err != nil {
		log.Errorf("failed to clean device %s: %s", in.Device, err.Error())
		return nil, status.Errorf(codes.Internal, "failed to clean device %s: %v", in.Device, err)
	}
	log.Debugf("clean device %s successfully", in.Device)
	return &lib.CleanDeviceReply{CommandOutput: fmt.Sprintf("clean device %s successfully with output: %s", in.Device, out)}, nil
}

// AddTagLV add tag
func (s Server) AddTagLV(ctx context.Context, in *lib.AddTagLVRequest) (*lib.AddTagLVReply, error) {
	log, err := AddTagLV(ctx, in.VolumeGroup, in.Name, in.Tags)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add tags to lv: %v", err)
	}
	return &lib.AddTagLVReply{CommandOutput: log}, nil
}

// RemoveTagLV remove tag
func (s Server) RemoveTagLV(ctx context.Context, in *lib.RemoveTagLVRequest) (*lib.RemoveTagLVReply, error) {
	log, err := RemoveTagLV(ctx, in.VolumeGroup, in.Name, in.Tags)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove tags from lv: %v", err)
	}
	return &lib.RemoveTagLVReply{CommandOutput: log}, nil
}
