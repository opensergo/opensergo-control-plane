#!/bin/sh
##!/usr/bin/env bash
# ======================================================================
# Linux/OSX startup script
# ======================================================================

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
OPENSERGO_CONTROL_PLANE_HOME=$HOME/.opensergo.install/opensergo-control-plane
mkdir -p $OPENSERGO_CONTROL_PLANE_HOME

# Exec Command Args
VERSION=main
IS_LOCAL=false
IS_UNINSTALL=false
IS_HELP=false

# Init Var, which will change by parseArgs()
DOWNLOAD_TARGET_URL_PREFIX=$OPENSERGO_CONTROL_PLANE_HOME/main
DOWNLOAD_SOURCE_URL_PREFIX=$RESOURCE_URL_PREFIX/main

# parse exec command args
allArgs=($*)
for (( i = 1; i <= $#; i++ )); do
  currArg=${allArgs[$i-1]}
  case $currArg in
      "-h" | "-help")
        IS_HELP=true
        ;;
      "-l" | "-local")
        IS_LOCAL=true
        ;;
      "-u" | "-uninstall")
        IS_LOCAL=true
        ;;
    esac

  if [ $i = $# ] ; then
    if [ "-" != ${currArg:0:1} ]; then
      VERSION=$currArg
    fi
  fi
done



function parseVersion() {
  # RESOURCE_URL
  resource_url_prefix=https://raw.githubusercontent.com/opensergo/opensergo-control-plane
  resource_url_version=$VERSION

  # parse VERSION_URL
  if [ "main" != $VERSION ] ; then
    resource_url_version=tree/$VERSION
  fi

  DOWNLOAD_TARGET_URL_PREFIX=$OPENSERGO_CONTROL_PLANE_HOME/$VERSION
  DOWNLOAD_SOURCE_URL_PREFIX=$resource_url_prefix/$resource_url_version
}

# invoke parseArgs()
parseVersion

# Resources in k8s/crd/bases
res_crd[0]=k8s/crd/bases/fault-tolerance.opensergo.io_circuitbreakerstrategies.yaml
res_crd[1]=k8s/crd/bases/fault-tolerance.opensergo.io_concurrencylimitstrategies.yaml
res_crd[2]=k8s/crd/bases/fault-tolerance.opensergo.io_faulttolerancerules.yaml
res_crd[3]=k8s/crd/bases/fault-tolerance.opensergo.io_ratelimitstrategies.yaml
res_crd[4]=k8s/crd/bases/fault-tolerance.opensergo.io_throttlingstrategies.yaml
res_crd[5]=k8s/crd/bases/traffic.opensergo.io_trafficerouters.yaml

function download_CRD() {
  rm -rf $OPENSERGO_CONTROL_PLANE_HOME/$VERSION/k8s/crd/bases/*
  mkdir -p $OPENSERGO_CONTROL_PLANE_HOME/$VERSION/k8s/crd/bases

  for crd in ${res_crd[*]}
  do
    wget --no-check-certificate -O $DOWNLOAD_TARGET_URL_PREFIX/$crd $DOWNLOAD_SOURCE_URL_PREFIX/$crd
  done
}

function install_crd() {
  for crd in ${res_crd[*]}
  do
    kubectl apply -f $DOWNLOAD_TARGET_URL_PREFIX/$crd
  done
}

function uninstall_crd() {
  for crd in ${res_crd[*]}
  do
    kubectl delete -f $DOWNLOAD_TARGET_URL_PREFIX/$crd
  done
}

# Resources in k8s
res_namespace[0]=k8s/namespace.yaml

function download_Namespace() {
  mkdir -p $OPENSERGO_CONTROL_PLANE_HOME/$VERSION/k8s

  for namespace in ${res_namespace[*]}
  do
    wget --no-check-certificate -O $DOWNLOAD_TARGET_URL_PREFIX/$namespace $DOWNLOAD_SOURCE_URL_PREFIX/$namespace
  done
}

function install_namespace() {
  for namespace in ${res_namespace[*]}
  do
    kubectl apply -f $DOWNLOAD_TARGET_URL_PREFIX/$namespace
  done
}

function uninstall_namespace() {
  for namespace in ${res_namespace[*]}
  do
    kubectl delete -f $DOWNLOAD_TARGET_URL_PREFIX/$namespace
  done
}

# Resources in k8s/rbac
res_rbac[0]=k8s/rbac/rbac.yaml

function download_RBAC() {
  rm -rf $OPENSERGO_CONTROL_PLANE_HOME/$VERSION/k8s/rbac/*
  mkdir -p $OPENSERGO_CONTROL_PLANE_HOME/$VERSION/k8s/rbac

  for rbac in ${res_rbac[*]}
  do
    wget --no-check-certificate -O $DOWNLOAD_TARGET_URL_PREFIX/$rbac $DOWNLOAD_SOURCE_URL_PREFIX/$rbac
  done
}

function install_rbac() {
  for rbac in ${res_rbac[*]}
  do
    kubectl apply -f $DOWNLOAD_TARGET_URL_PREFIX/$rbac
  done
}

function uninstall_rbac() {
  for rbac in ${res_rbac[*]}
  do
    kubectl delete -f $DOWNLOAD_TARGET_URL_PREFIX/$rbac
  done
}

function download() {
    download_CRD
    download_Namespace
    download_RBAC
}

function uninstall() {
    uninstall_rbac
    uninstall_namespace
    uninstall_crd
}

function install() {
  # invoke uninstall
  uninstall
  # download or not
  if [ true != $IS_LOCAL ]; then
    download
  fi
  # invoke install
  install_crd
  install_namespace
  install_rbac
}


function usage() {
  echo "Usage:  init.sh [-Option] [version]"
  echo "Resources will download in $OPENSERGO_CONTROL_PLANE_HOME"
  echo ""
  echo "Available Options:"
  echo "  -l, -local         : exec from local resources, which will not download resources"
  echo "  -u, -uninstall     : uninstall CRDs, Namespace, RBAC of opensergo control plane"
  echo "  -h, -help          : helps"
  echo "version:             : default version is main"
}

if [ true = $IS_HELP ] ; then
  usage
elif [ true = $IS_UNINSTALL ] ; then
  uninstall
else
  install
fi