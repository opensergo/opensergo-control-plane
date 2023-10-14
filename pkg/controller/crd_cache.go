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

package controller

import (
	"sync"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/opensergo/opensergo-control-plane/pkg/model"
)

type CRDObjectsHolder struct {
	objects []client.Object
	version int64
}

// TODO: need use generic to support two kinds of connection xDS Connection and normal Connection
// CRDCache caches versioned CRD objects in local.
type CRDCache struct {
	kind string
	// crdEntityMap represents a map: (namespace, name) -> unique CRD
	crdEntityMap map[types.NamespacedName]client.Object
	// namespaceAppMap represents a map for CRD group: (namespace, app) -> versionedCRDs
	namespaceAppMap map[model.NamespacedApp]*CRDObjectsHolder

	updateMux sync.RWMutex
}

func NewCRDCache(kind string) *CRDCache {
	return &CRDCache{
		kind:            kind,
		crdEntityMap:    make(map[types.NamespacedName]client.Object),
		namespaceAppMap: make(map[model.NamespacedApp]*CRDObjectsHolder),
	}
}

func (c *CRDCache) GetByNamespacedName(n types.NamespacedName) (client.Object, bool) {
	c.updateMux.RLock()
	defer c.updateMux.RUnlock()

	obj, exists := c.crdEntityMap[n]
	return obj, exists
}

func (c *CRDCache) SetByNamespacedName(n types.NamespacedName, object client.Object) {
	c.updateMux.Lock()
	defer c.updateMux.Unlock()

	c.crdEntityMap[n] = object
}

func (c *CRDCache) DeleteByNamespacedName(n types.NamespacedName) {
	c.updateMux.Lock()
	defer c.updateMux.Unlock()

	delete(c.crdEntityMap, n)
}

func (c *CRDCache) SetByNamespaceApp(n model.NamespacedApp, object client.Object) int {
	c.updateMux.Lock()
	defer c.updateMux.Unlock()

	o, exists := c.namespaceAppMap[n]
	if !exists || o == nil {
		c.namespaceAppMap[n] = &CRDObjectsHolder{
			objects: []client.Object{object},
			version: 1,
		}
		return AddRule
	} else {
		for index := range o.objects {
			// Update object version
			if o.objects[index].GetName() == object.GetName() {
				o.objects[index] = object
				o.version++
				return UpdateRule
			}
		}
		o.objects = append(o.objects, object)
		o.version++
		return AddRule
	}
}

func (c *CRDCache) DeleteByNamespaceApp(n model.NamespacedApp, name string) {
	c.updateMux.Lock()
	defer c.updateMux.Unlock()

	o, exists := c.namespaceAppMap[n]
	if !exists || o == nil {
		return
	}
	for index, obj := range o.objects {
		if obj.GetLabels()["app"] == n.App && obj.GetName() == name {
			// Update object version
			o.objects = append(o.objects[:index], o.objects[index+1:]...)
			o.version++
			return
		}
	}
}

func (c *CRDCache) GetByNamespaceApp(n model.NamespacedApp) ([]client.Object, int64) {
	c.updateMux.RLock()
	defer c.updateMux.RUnlock()

	o, _ := c.namespaceAppMap[n]
	if o == nil {
		return []client.Object{}, 0
	}
	return o.objects, o.version
}

func (c *CRDCache) GetAppByNamespacedName(n types.NamespacedName) (string, bool) {
	c.updateMux.RLock()
	defer c.updateMux.RUnlock()

	crd, exists := c.crdEntityMap[n]
	if exists {
		return crd.GetLabels()["app"], true
	} else {
		return "", false
	}
}
