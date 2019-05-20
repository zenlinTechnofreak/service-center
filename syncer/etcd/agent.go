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
	"github.com/apache/servicecomb-service-center/pkg/log"
	pb "github.com/apache/servicecomb-service-center/syncer/proto"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/embed"
	"github.com/coreos/etcd/etcdserver/api/v3client"
)

// Agent warps the embed etcd
type Agent struct {
	conf     *Config
	etcd     *embed.Etcd
	client   *clientv3.Client
	syncData *pb.SyncData

	// mapping table for other datacenter instances
	intsMapping map[string]pb.SyncMapping
	//lock        sync.RWMutex
}

// NewAgent new etcd agent with config
func NewAgent(conf *Config) *Agent {
	return &Agent{conf: conf}
}

// Run etcd agent
func (a *Agent) Run() error {
	etcd, err := embed.StartEtcd(a.conf.Config)
	if err != nil {
		return err
	}
	select {
	case <-etcd.Server.ReadyNotify():
		log.Info("ready notify")
	case <-etcd.Server.StopNotify():
		log.Info("stop notify")
	case err = <-etcd.Err():
		log.Error("start etcd failed", err)
		return err
	}
	a.etcd = etcd
	return nil
}

// KVStore Think of ETCD as a key-value store
func (a *Agent) KVStore() clientv3.KV {
	return a.getClient()
}

// Watcher Think of ETCD as a watcher
func (a *Agent) Watcher() clientv3.Watcher {
	return a.getClient()
}

// New a etcd client if not exist
func (a *Agent) getClient() *clientv3.Client {
	if a.client == nil {
		a.client = v3client.New(a.etcd.Server)
	}
	return a.client
}

// Stop etcd agent
func (a *Agent) Stop() {
	a.etcd.Close()
}
