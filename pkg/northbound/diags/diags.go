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

package diags

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/onosproject/onos-config/pkg/manager"
	"github.com/onosproject/onos-config/pkg/northbound"
	"github.com/onosproject/onos-config/pkg/northbound/proto"
	"github.com/onosproject/onos-config/pkg/store"
	"google.golang.org/grpc"
)

// Service is a Service implementation for administration.
type Service struct {
	northbound.Service
}

// Register registers the Service with the gRPC server.
func (s Service) Register(r *grpc.Server) {
	proto.RegisterConfigDiagsServer(r, Server{})
}

// Server implements the gRPC service for diagnostic facilities.
type Server struct {
}

// GetChanges provides a stream of submitted network changes.
func (s Server) GetChanges(r *proto.ChangesRequest, stream proto.ConfigDiags_GetChangesServer) error {
	for _, c := range manager.GetManager().ChangeStore.Store {
		// Build a change message
		msg := &proto.Change{
			Time: &timestamp.Timestamp{Seconds: c.Created.Unix(), Nanos: int32(c.Created.Nanosecond())},
			Id: store.B64(c.ID),
			Desc: c.Description,
		}
		err := stream.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}