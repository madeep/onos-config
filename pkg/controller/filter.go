// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"github.com/onosproject/onos-config/api/types"
	mastershipstore "github.com/onosproject/onos-config/pkg/store/mastership"
	devicetopo "github.com/onosproject/onos-topo/pkg/northbound/device"
	"regexp"
)

// Filter filters individual events for a node
// Each time an event is received from a watcher, the filter has the option to discard the request by
// returning false.
type Filter interface {
	// Accept indicates whether to accept the given object
	Accept(id types.ID) bool
}

// MastershipFilter activates a controller on acquiring mastership
// The MastershipFilter requires a DeviceResolver to extract a device ID from each request. Given a device
// ID, the MastershipFilter rejects any requests for which the local node is not the master for the device.
type MastershipFilter struct {
	Store    mastershipstore.Store
	Resolver DeviceResolver
}

// Accept accepts the given ID if the local node is the master
func (f *MastershipFilter) Accept(id types.ID) bool {
	device, err := f.Resolver.Resolve(id)
	if err != nil {
		return false
	}
	master, err := f.Store.IsMaster(device)
	if err != nil {
		return false
	}
	return master
}

var _ Filter = &MastershipFilter{}

// DeviceResolver resolves a device from a type ID
type DeviceResolver interface {
	// Resolve resolves a device
	Resolve(id types.ID) (devicetopo.ID, error)
}

// RegexpDeviceResolver is a DeviceResolver that reads a device ID from a regexp
type RegexpDeviceResolver struct {
	Regexp regexp.Regexp
}

// Resolve resolves a device ID from the configured regexp
func (r *RegexpDeviceResolver) Resolve(id types.ID) (devicetopo.ID, error) {
	return devicetopo.ID(r.Regexp.FindString(string(id))), nil
}

var _ DeviceResolver = &RegexpDeviceResolver{}
