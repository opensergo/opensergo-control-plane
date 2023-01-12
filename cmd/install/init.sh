#!/bin/sh

# Copyright 2022, OpenSergo Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# Prepare Env
# OPENSERGO_CONTROL_PLANE_HOME
OPENSERGO_CONTROL_PLANE_HOME=$HOME/opensergo/opensergo-control-plane
mkdir -p $OPENSERGO_CONTROL_PLANE_HOME


# Uninstall CRDs
kubectl delete -f $OPENSERGO_CONTROL_PLANE_HOME/k8s/crd/bases
# Download CRDs ./k8s/crd/bases
mkdir -p $OPENSERGO_CONTROL_PLANE_HOME/k8s/crd/bases
wget --no-check-certificate -O $OPENSERGO_CONTROL_PLANE_HOME/k8s/crd/bases/fault-tolerance.opensergo.io_circuitbreakerstrategies.yaml   https://raw.githubusercontent.com/opensergo/opensergo-control-plane/main/k8s/crd/bases/fault-tolerance.opensergo.io_circuitbreakerstrategies.yaml
wget --no-check-certificate -O $OPENSERGO_CONTROL_PLANE_HOME/k8s/crd/bases/fault-tolerance.opensergo.io_concurrencylimitstrategies.yaml https://raw.githubusercontent.com/opensergo/opensergo-control-plane/main/k8s/crd/bases/fault-tolerance.opensergo.io_concurrencylimitstrategies.yaml
wget --no-check-certificate -O $OPENSERGO_CONTROL_PLANE_HOME/k8s/crd/bases/fault-tolerance.opensergo.io_faulttolerancerules.yaml        https://raw.githubusercontent.com/opensergo/opensergo-control-plane/main/k8s/crd/bases/fault-tolerance.opensergo.io_faulttolerancerules.yaml
wget --no-check-certificate -O $OPENSERGO_CONTROL_PLANE_HOME/k8s/crd/bases/fault-tolerance.opensergo.io_ratelimitstrategies.yaml        https://raw.githubusercontent.com/opensergo/opensergo-control-plane/main/k8s/crd/bases/fault-tolerance.opensergo.io_ratelimitstrategies.yaml
wget --no-check-certificate -O $OPENSERGO_CONTROL_PLANE_HOME/k8s/crd/bases/fault-tolerance.opensergo.io_throttlingstrategies.yaml       https://raw.githubusercontent.com/opensergo/opensergo-control-plane/main/k8s/crd/bases/fault-tolerance.opensergo.io_throttlingstrategies.yaml
wget --no-check-certificate -O $OPENSERGO_CONTROL_PLANE_HOME/k8s/crd/bases/traffic.opensergo.io_trafficerouters.yaml                    https://raw.githubusercontent.com/opensergo/opensergo-control-plane/main/k8s/crd/bases/traffic.opensergo.io_trafficerouters.yaml
# Install CRDs
kubectl apply -f $OPENSERGO_CONTROL_PLANE_HOME/k8s/crd/bases


# Uninstall Namespace
kubectl delete -f $OPENSERGO_CONTROL_PLANE_HOME/cmd/install/k8s/namespace.yaml
# Download Namespace
mkdir -p $OPENSERGO_CONTROL_PLANE_HOME/k8s
wget --no-check-certificate -O $OPENSERGO_CONTROL_PLANE_HOME/k8s/namespace.yaml  https://raw.githubusercontent.com/opensergo/opensergo-control-plane/main/k8s/namespace.yaml
# Install Namespace
kubectl apply -f $OPENSERGO_CONTROL_PLANE_HOME/k8s/namespace.yaml


# Uninstall rbac.yaml
kubectl delete -f $OPENSERGO_CONTROL_PLANE_HOME/cmd/install/k8s/rbac/rbac.yaml
# Download rbac.yaml
mkdir -p $OPENSERGO_CONTROL_PLANE_HOME/cmd/install/k8s/rbac
wget --no-check-certificate -O $OPENSERGO_CONTROL_PLANE_HOME/cmd/install/k8s/rbac/rbac.yaml https://raw.githubusercontent.com/opensergo/opensergo-control-plane/main/cmd/install/k8s/rbac/rbac.yaml
# Install rbac.yaml
kubectl apply -f $OPENSERGO_CONTROL_PLANE_HOME/cmd/install/k8s/rbac/rbac.yaml