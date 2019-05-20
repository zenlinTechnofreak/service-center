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

package etcd

import (
	"context"

	"github.com/apache/servicecomb-service-center/pkg/log"
	pb "github.com/apache/servicecomb-service-center/syncer/proto"
	"github.com/coreos/etcd/clientv3"
	"github.com/gogo/protobuf/proto"
)

var (
	defaultMapping = make([]*pb.SyncMapping, 0)
	// servicesKey the key of service in etcd
	servicesKey = "/syncer/v1/services"

	// instancesKey the key of instance in etcd
	instancesKey = "/syncer/v1/instances"

	// expansionsKey the key of extended attribute in etcd
	expansionsKey = "/syncer/v1/expansions"

	// mappingsKey the key of instances mapping in etcd
	mappingsKey = "/syncer/v1/mappings"
)

// SaveSyncData Save self sync data
func (a *Agent) SaveSyncData(data *pb.SyncData) {
	a.syncData = data
}

// GetSyncData Get self sync data
func (a *Agent) GetSyncData() (data *pb.SyncData) {
	data = &pb.SyncData{Services: a.syncData.Services[:]}
	return
}

// SaveSyncMappings Save mappings for other datacenter instances
func (a *Agent) SaveSyncMappings(nodeName string, mappings []*pb.SyncMapping) error {
	for _, val := range mappings {
		key := mappingsKey + "/" + nodeName + "/" + val.OrgInstanceID
		data, err := proto.Marshal(val)
		if err != nil {
			log.Errorf(err, "Proto marshal failed: %s", err)
			return err
		}
		_, err = a.getClient().Put(context.Background(), key, string(data))
		if err != nil {
			log.Errorf(err, "Save mapping to etcd failed: %s", err)
			return err
		}
	}
	return nil
}

// GetSyncMapping Get mappings for other datacenter instances
func (a *Agent) GetSyncMappings(nodeName string) []*pb.SyncMapping {
	return a.getMappingsByPrefixKey(mappingsKey + "/" + nodeName)
}

// GetAllMappings Get all mappings for other datacenters instances
func (a *Agent) GetAllMappings() []*pb.SyncMapping {
	return a.getMappingsByPrefixKey(mappingsKey)
}

// getMappingsByPrefixKey Get mappings from etcd based on the prefix
func (a *Agent) getMappingsByPrefixKey(prefix string) []*pb.SyncMapping {
	resp, err := a.getClient().Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		log.Errorf(err, "Get mapping from etcd failed: %s", err)
		return defaultMapping
	}

	if len(resp.Kvs) == 0 {
		return defaultMapping
	}

	mappings := make([]*pb.SyncMapping, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		mapping := &pb.SyncMapping{}
		err = proto.Unmarshal(kv.Value, mapping)
		if err != nil {
			log.Errorf(err, "Proto unmarshal '%s' failed: %s", string(kv.Value), err)
			continue
		}
		mappings = append(mappings, mapping)
	}
	return mappings
}
