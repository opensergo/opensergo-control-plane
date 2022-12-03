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
	"github.com/opensergo/opensergo-control-plane/pkg/api/v1alpha1"
	"github.com/opensergo/opensergo-control-plane/pkg/api/v1beta1/networking"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CRDKind = string

// CRDGenerator represents a generator function of an OpenSergo CRD.
type CRDGenerator = func() client.Object

type CRDMetadata struct {
	kind CRDKind

	generator CRDGenerator
}

func (m *CRDMetadata) Kind() CRDKind {
	return m.kind
}

func (m *CRDMetadata) Generator() CRDGenerator {
	return m.generator
}

func NewCRDMetadata(kind CRDKind, generator CRDGenerator) *CRDMetadata {
	return &CRDMetadata{
		kind:      kind,
		generator: generator,
	}
}

const (
	FaultToleranceRuleKind       = "fault-tolerance.opensergo.io/v1alpha1/FaultToleranceRule"
	RateLimitStrategyKind        = "fault-tolerance.opensergo.io/v1alpha1/RateLimitStrategy"
	ThrottlingStrategyKind       = "fault-tolerance.opensergo.io/v1alpha1/ThrottlingStrategy"
	ConcurrencyLimitStrategyKind = "fault-tolerance.opensergo.io/v1alpha1/ConcurrencyLimitStrategy"
	CircuitBreakerStrategyKind   = "fault-tolerance.opensergo.io/v1alpha1/CircuitBreakerStrategy"
	VirtualServiceKind           = "networking.istio.io/v1beta1/VirtualService"
)

var (
	// crdMetadataMap is the universal registry for all OpenSergo CRDs.
	crdMetadataMap = map[CRDKind]*CRDMetadata{
		FaultToleranceRuleKind: NewCRDMetadata(FaultToleranceRuleKind, func() client.Object {
			return &v1alpha1.FaultToleranceRule{}
		}),
		RateLimitStrategyKind: NewCRDMetadata(RateLimitStrategyKind, func() client.Object {
			return &v1alpha1.RateLimitStrategy{}
		}),
		ThrottlingStrategyKind: NewCRDMetadata(ThrottlingStrategyKind, func() client.Object {
			return &v1alpha1.ThrottlingStrategy{}
		}),
		ConcurrencyLimitStrategyKind: NewCRDMetadata(ConcurrencyLimitStrategyKind, func() client.Object {
			return &v1alpha1.ConcurrencyLimitStrategy{}
		}),
		CircuitBreakerStrategyKind: NewCRDMetadata(CircuitBreakerStrategyKind, func() client.Object {
			return &v1alpha1.CircuitBreakerStrategy{}
		}),
		VirtualServiceKind: NewCRDMetadata(VirtualServiceKind, func() client.Object {
			return &networking.VirtualService{}
		}),
	}
)

func GetCrdMetadata(kind CRDKind) (*CRDMetadata, bool) {
	// TODO: should we put a lock here?
	data, exists := crdMetadataMap[kind]
	return data, exists
}
