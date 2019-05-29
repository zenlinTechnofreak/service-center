/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package server

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/servicecomb-service-center/pkg/log"
	"github.com/apache/servicecomb-service-center/syncer/grpc"
	pb "github.com/apache/servicecomb-service-center/syncer/proto"
	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/serf/serf"
)

const (
	EventDiscovered = "discovered"
)

// tickHandler Timed task handler
func (s *Server) tickHandler(ctx context.Context) {
	log.Debugf("Handle Tick")
	// Flush data to the storage of servicecenter
	s.servicecenter.FlushData()

	event, _ := proto.Marshal(&pb.Member{
		NodeName: s.conf.NodeName,
		RPCPort:  int32(s.conf.RPCPort),
		Time:     fmt.Sprintf("%d", time.Now().UTC().Second()),
	})

	// sends a UserEvent on Serf, the event will be broadcast between members
	err := s.agent.UserEvent(EventDiscovered, event, true)
	if err != nil {
		log.Errorf(err, "Syncer send user event failed")
	}
}

// GetData Sync Data to GRPC
func (s *Server) Discovery() *pb.SyncData {
	return s.servicecenter.Discovery()
}

// HandleEvent Handles events from serf
func (s *Server) HandleEvent(event serf.Event) {
	switch event.EventType() {
	case serf.EventUser:
		s.userEvent(event.(serf.UserEvent))
	case serf.EventQuery:
		s.queryEvent(event.(*serf.Query))
	}
}

// userEvent Handles "EventUser" notification events, no response required
func (s *Server) userEvent(event serf.UserEvent) {
	m := &pb.Member{}
	err := proto.Unmarshal(event.Payload, m)
	if err != nil {
		log.Errorf(err, "trigger user event '%s' handler failed", event.EventType())
		return
	}

	// Excludes notifications from self, as the gossip protocol inevitably has redundant notifications
	if s.agent.LocalMember().Name == m.NodeName {
		return
	}

	// Get member information and get synchronized data from it
	member := s.agent.Member(m.NodeName)
	// Get dta from remote member
	endpoint := fmt.Sprintf("%s:%d", member.Addr, m.RPCPort)
	log.Debugf("Going to pull data from %s %s", m.NodeName, endpoint)
	data, err := grpc.Pull(context.Background(), endpoint)
	if err != nil {
		log.Errorf(err, "Pull other serf instances failed, node name is '%s'", m.NodeName)
		return
	}
	// Registry instances to servicecenter and update storage of it
	s.servicecenter.Registry(m.NodeName, data)
}

// queryEvent Handles "EventQuery" query events and respond if conditions are met
func (s *Server) queryEvent(query *serf.Query) {
	// todo: Get instances requested
}
