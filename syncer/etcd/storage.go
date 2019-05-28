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
	"encoding/json"
	"sync"

	"github.com/apache/servicecomb-service-center/pkg/log"
	"github.com/apache/servicecomb-service-center/pkg/util"
	"github.com/apache/servicecomb-service-center/syncer/datacenter"
	pb "github.com/apache/servicecomb-service-center/syncer/proto"
	"github.com/coreos/etcd/clientv3"
	"github.com/gogo/protobuf/proto"
)

var (
	defaultMapping = make(pb.SyncMapping, 0)

	// servicesKey the key of service in etcd
	servicesKey = "/syncer/v1/services"

	// instancesKey the key of instance in etcd
	instancesKey = "/syncer/v1/instances"

	// expansionsKey the key of extended attribute in etcd
	expansionsKey = "/syncer/v1/expansions"

	// mappingsKey the key of instances mapping in etcd
	mappingsKey = "/syncer/v1/mappings"
)

type storage struct {
	client *clientv3.Client
	data   *pb.SyncData
	// mapping table for other datacenter instances
	maps map[string]pb.SyncMapping
	lock sync.RWMutex
}

func NewStorage(client *clientv3.Client) datacenter.Storage {
	storage := &storage{
		client: client,
		data:   &pb.SyncData{},
	}
	storage.maps = storage.getMapsFromEtcd()
	return storage
}

// getMapsFromEtcd load maps from etcd to storage
func (s *storage) getMapsFromEtcd() map[string]pb.SyncMapping {
	maps := make(map[string]pb.SyncMapping)
	s.getPrefixKey(mappingsKey, func(key, val []byte) (next bool) {
		next = true
		entry := &pb.MappingEntry{}
		if err := proto.Unmarshal(val, entry); err != nil {
			log.Errorf(err, "Proto unmarshal '%s' failed: %s", val, err)
			return
		}

		mapping, ok := maps[entry.NodeName]
		if !ok {
			mapping = make(pb.SyncMapping, 0, 10)
		}
		mapping = append(mapping, entry)
		maps[entry.NodeName] = mapping
		return
	})

	data, _ := json.MarshalIndent(&maps, "", "  ")
	log.Infof("%s", data)
	log.Infof("Loaded maps from disk to storage", maps)
	return maps
}

// getPrefixKey Get data from etcd based on the prefix key
func (s *storage) getPrefixKey(prefix string, handler func(key, val []byte) (next bool)) {
	resp, err := s.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		log.Errorf(err, "Get mapping from etcd failed: %s", err)
		return
	}

	for _, kv := range resp.Kvs {
		if !handler(kv.Key, kv.Value) {
			break
		}
	}
}

// UpdateData Update data to storage
func (s *storage) UpdateData(data *pb.SyncData) {
	s.lock.Lock()
	s.data = data
	s.lock.Unlock()
}

// GetData Get data from storage
func (s *storage) GetData() (data *pb.SyncData) {
	s.lock.RLock()
	data = s.data
	s.lock.RUnlock()
	return
}

// cleanExpired clean the expired instances in the maps
func (s *storage) cleanExpired(sources pb.SyncMapping, actives pb.SyncMapping) {
next:
	for _, entry := range sources {
		for _, act := range actives {
			if act.CurInstanceID == entry.CurInstanceID {
				continue next
			}
		}
		key := mappingsKey + "/" + entry.NodeName + "/" + entry.OrgInstanceID
		if _, err := s.client.Delete(context.Background(), key); err != nil {
			log.Errorf(err, "Delete instance nodeName=%s instanceID=%s failed", entry.NodeName, entry.OrgInstanceID)
		}
	}
}

// UpdateMapByNode update map to storage by nodeName of other node
func (s *storage) UpdateMapByNode(nodeName string, mapping pb.SyncMapping) {
	newMaps := make(pb.SyncMapping, 0, len(mapping))
	for _, val := range mapping {
		key := mappingsKey + "/" + nodeName + "/" + val.OrgInstanceID
		data, err := proto.Marshal(val)
		if err != nil {
			log.Errorf(err, "Proto marshal failed: %s", err)
			continue
		}
		_, err = s.client.Put(context.Background(), key, util.BytesToStringWithNoCopy(data))
		if err != nil {
			log.Errorf(err, "Save mapping to etcd failed: %s", err)
		}
		newMaps = append(newMaps, val)
	}
	s.cleanExpired(s.maps[nodeName], newMaps)
	s.lock.Lock()
	s.maps[nodeName] = newMaps
	s.lock.Unlock()
}

// GetMapByNode get map by nodeName of other node
func (s *storage) GetMapByNode(nodeName string) (mapping pb.SyncMapping) {
	s.lock.RLock()
	data, ok := s.maps[nodeName]
	if !ok {
		data = defaultMapping
	}
	s.lock.RUnlock()
	return data
}

// UpdateMaps update all maps to etcd
func (s *storage) UpdateMaps(maps pb.SyncMapping) {
	mappings := make(map[string]pb.SyncMapping)
	for _, val := range maps {
		key := mappingsKey + "/" + val.NodeName + "/" + val.OrgInstanceID
		data, err := proto.Marshal(val)
		if err != nil {
			log.Errorf(err, "Proto marshal failed: %s", err)
			continue
		}
		_, err = s.client.Put(context.Background(), key, util.BytesToStringWithNoCopy(data))
		if err != nil {
			log.Errorf(err, "Save mapping to etcd failed: %s", err)
			continue
		}
		mapping, ok := mappings[val.NodeName]
		if !ok {
			mapping = make(pb.SyncMapping, 0, 10)
		}
		mapping = append(mapping, val)
		mappings[val.NodeName] = mapping
	}

	for node, mapping := range mappings {
		s.cleanExpired(s.maps[node], mapping)
	}
	s.lock.Lock()
	s.maps = mappings
	s.lock.Unlock()
}

// GetMaps Get maps from storage
func (s *storage) GetMaps() (mapping pb.SyncMapping) {
	s.lock.RLock()
	mapping = make(pb.SyncMapping, 0, 10)
	for _, data := range s.maps {
		if data != nil {
			mapping = append(mapping, data...)
		}
	}
	s.lock.RUnlock()
	return
}
