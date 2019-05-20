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
package datacenter

import (
	"context"

	"github.com/apache/servicecomb-service-center/pkg/log"
	pb "github.com/apache/servicecomb-service-center/syncer/proto"
)

// sync synchronize information from other datacenters
func (s *store) sync(data *pb.SyncData, curMappings []*pb.SyncMapping) ([]*pb.SyncMapping, error) {
	orgMappings, curMappings := s.syncServiceInstances(data.Services, curMappings)
	return s.deleteInstances(orgMappings, curMappings), nil
}

// syncServiceInstances register instances from other datacenter
func (s *store) syncServiceInstances(services []*pb.SyncService, curMappings []*pb.SyncMapping) ([]*pb.SyncMapping, []*pb.SyncMapping) {
	var err error
	ctx := context.Background()
	orgMappings := make([]*pb.SyncMapping, 0, 10)

	for _, svc := range services {
		curServiceID := ""
		for _, inst := range svc.Instances {

			// Send an instance heartbeat if the instance has already been registered
			skip := false
			for _, val := range curMappings {
				if val.OrgInstanceID == inst.InstanceId {
					curServiceID = val.CurServiceID
					err = s.datacenter.Heartbeat(ctx, val.DomainProject, val.CurServiceID, val.CurInstanceID)
					if err != nil {
						log.Errorf(err, "Syncer heartbeat instance failed")
					} else {
						orgMappings = append(orgMappings, val)
					}
					skip = true
					break
				}
			}
			if skip {
				continue
			}

			// Create microservice if the service to which the instance belongs does not exist
			if curServiceID == "" {
				if curServiceID, _ = s.datacenter.ServiceExistence(ctx, svc.DomainProject, svc.Service); curServiceID == "" {
					svc.Service.ServiceId = ""
					curServiceID, err = s.datacenter.CreateService(ctx, svc.DomainProject, svc.Service)
					if err != nil {
						log.Errorf(err, "Syncer create service failed")
						continue
					}
				}
			}
			// Register instance information when the instance does not exist
			item := &pb.SyncMapping{
				DomainProject: svc.DomainProject,
				OrgServiceID:  inst.ServiceId,
				OrgInstanceID: inst.InstanceId,
				CurServiceID:  curServiceID,
			}
			inst.ServiceId = curServiceID
			inst.InstanceId = ""
			curInstanceID, err := s.datacenter.RegisterInstance(ctx, svc.DomainProject, curServiceID, inst)
			if err != nil {
				log.Errorf(err, "Syncer create service failed")
				continue
			}

			item.CurInstanceID = curInstanceID
			orgMappings = append(orgMappings, item)
			curMappings = append(curMappings, item)
		}
	}
	return orgMappings, curMappings
}

// deleteInstances Unregister instances of mapping table that has been unregistered from other datacenter
func (s *store) deleteInstances(orgMappings, curMappings []*pb.SyncMapping) []*pb.SyncMapping {
	l := 0
	ol := len(orgMappings)
	cl := len(curMappings)
	if ol < cl {
		l = cl - ol
	} else {
		l = ol - cl
	}
	if l == 0 {
		return curMappings
	}

	ctx := context.Background()
	nm := make([]*pb.SyncMapping, 0, l)
	for _, val := range curMappings {
		skip := false
		for _, org := range orgMappings {
			if val.CurInstanceID == org.CurInstanceID {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		err := s.datacenter.UnregisterInstance(ctx, val.DomainProject, val.CurServiceID, val.CurInstanceID)
		if err != nil {
			log.Errorf(err, "Syncer delete service failed")
		}
	}
	return nm
}
