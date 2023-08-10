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

# Resources in k8s/workload
res_deploy[0]=k8s/workload/opensergo-control-plane.yaml

function download_workload() {
  rm -rf $OPENSERGO_CONTROL_PLANE_HOME/$VERSION/k8s/workload/*
  mkdir -p $OPENSERGO_CONTROL_PLANE_HOME/$VERSION/k8s/workload

  for crd in ${res_deploy[*]}
  do
    wget --no-check-certificate -O $DOWNLOAD_TARGET_URL_PREFIX/$crd $DOWNLOAD_SOURCE_URL_PREFIX/$crd
  done
}

function install_workload() {
  for crd in ${res_deploy[*]}
  do
    kubectl apply -f $DOWNLOAD_TARGET_URL_PREFIX/$crd
  done
}

function uninstall_workload() {
  for crd in ${res_deploy[*]}
  do
    kubectl delete -f $DOWNLOAD_TARGET_URL_PREFIX/$crd
  done
}

function install() {
  # invoke uninstall
  uninstall_workload
  # download or not
  if [ true != $IS_LOCAL ]; then
    download_workload
  fi
  # invoke install
  install_workload
}


function usage() {
  echo "Usage:  deploy.sh [-Option] [version]"
  echo "Resources will download in $OPENSERGO_CONTROL_PLANE_HOME"
  echo ""
  echo "Available Options:"
  echo "  -l, -local         : exec from local resources, which will not download resources"
  echo "  -u, -uninstall     : uninstall Workload of opensergo control plane"
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