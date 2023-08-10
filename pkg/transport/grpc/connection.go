// Copyright 2022, OpenSergo Authors
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

package grpc

import (
	"sync"

	"github.com/opensergo/opensergo-control-plane/pkg/model"
	pb "github.com/opensergo/opensergo-control-plane/pkg/proto/transport/v1"
	"github.com/pkg/errors"
)

type OpenSergoTransportStream = pb.OpenSergoUniversalTransportService_SubscribeConfigServer

type Connection struct {
	identifier model.ClientIdentifier
	stream     OpenSergoTransportStream

	valid bool
}

func (c *Connection) Identifier() model.ClientIdentifier {
	return c.identifier
}

func (c *Connection) Stream() OpenSergoTransportStream {
	return c.stream
}

func (c *Connection) IsValid() bool {
	return c.stream != nil && c.valid
}

type mapValue interface {
	interface {
		*Connection | *model.XDsConnection
	}
}
type ConnectionMap[T mapValue] map[model.ClientIdentifier]T

func NewConnection(identifier model.ClientIdentifier, stream OpenSergoTransportStream) *Connection {
	return &Connection{
		identifier: identifier,
		stream:     stream,
		valid:      true,
	}
}

type ConnectionManager[T mapValue] struct {
	// connectionMap is used to save the connections which subscribed to the same namespace, app and kind.
	// (namespace+app, (kind, connections...))
	connectionMap map[model.NamespacedApp]map[string]ConnectionMap[T]
	// identifier: NamespaceApp: kinds
	// The identifier is used to distinguish the requested process instance and remove stream when disconnected
	identifierMap map[model.ClientIdentifier]map[model.NamespacedApp][]string

	updateMux sync.RWMutex
}

func (c *ConnectionManager[T]) Add(namespace, app, kind string, connection T, identifier model.ClientIdentifier) error {
	if connection == nil {
		return errors.New("nil connection")
	}

	c.updateMux.Lock()
	defer c.updateMux.Unlock()

	nsa := model.NamespacedApp{
		Namespace: namespace,
		App:       app,
	}
	if c.connectionMap[nsa] == nil {
		c.connectionMap[nsa] = make(map[string]ConnectionMap[T])
	}
	connectionMap := c.connectionMap[nsa][kind]
	if connectionMap == nil {
		connectionMap = make(ConnectionMap[T])
		c.connectionMap[nsa][kind] = connectionMap
	}
	if connectionMap[identifier] == nil {
		connectionMap[identifier] = connection
	}

	// TODO: legacy logic, rearrange it later
	if c.identifierMap[identifier] == nil {
		c.identifierMap[identifier] = make(map[model.NamespacedApp][]string)
	}
	c.identifierMap[identifier][nsa] = append(c.identifierMap[identifier][nsa], kind)

	return nil
}

func (c *ConnectionManager[T]) Get(namespace, app, kind string) ([]T, bool) {
	c.updateMux.RLock()
	defer c.updateMux.RUnlock()

	kindMap, exists := c.connectionMap[model.NamespacedApp{
		Namespace: namespace,
		App:       app,
	}]
	if !exists || kindMap == nil {
		return nil, false
	}
	connectionMap, exists := kindMap[kind]
	if !exists || connectionMap == nil {
		return nil, false
	}

	connectionList := make([]T, len(connectionMap))
	for _, conn := range connectionMap {
		connectionList = append(connectionList, conn)
	}
	return connectionList, true
}

func (c *ConnectionManager[mapValue]) removeInternal(n model.NamespacedApp, kind string, identifier model.ClientIdentifier) error {
	// Guarded in the outer function, if a lock is added here, it will deadlock
	kindMap, exists := c.connectionMap[n]

	// TODO: handle error
	if !exists || kindMap == nil {
		return nil
	}
	streams, exists := kindMap[kind]
	if !exists || streams == nil {
		return nil
	}
	delete(streams, identifier)
	return nil
}

func (c *ConnectionManager[mapValue]) RemoveByIdentifier(identifier model.ClientIdentifier) error {
	c.updateMux.Lock()
	defer c.updateMux.Unlock()

	NamespaceAppKinds, exists := c.identifierMap[identifier]
	if !exists {
		return nil
	}
	for n, kinds := range NamespaceAppKinds {
		for _, kind := range kinds {
			err := c.removeInternal(n, kind, identifier)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewConnectionManager[T mapValue]() *ConnectionManager[T] {
	return &ConnectionManager[T]{
		connectionMap: make(map[model.NamespacedApp]map[string]ConnectionMap[T]),
		identifierMap: make(map[model.ClientIdentifier]map[model.NamespacedApp][]string),
	}
}
