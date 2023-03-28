#!/bin/bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

# Check if buildah is present in host os
if [ -z "$(command -v podman)" ]; then
  echo "ERR: podman is not installed"
  echo "RUN: sudo apt install podman"
  exit 1
fi

NET_NAME="cluster"

# Creating net network
podman network create $NET_NAME

# Creating first node
podman run --network=$NET_NAME --name scylla-node1 -d scylladb/scylla --overprovisioned 1 --smp 1 --memory 256M

# Extracting IP of first node
SEED=$(podman inspect scylla-node1 --format "{{.NetworkSettings.Networks.$NET_NAME.IPAddress}}")

# Creating second and third node
podman run --network=$NET_NAME --name scylla-node2 -d scylladb/scylla --seed=$SEED --overprovisioned 1 --smp 1 --memory 256M

podman run --network=$NET_NAME --name scylla-node3 -d scylladb/scylla --seed=$SEED --overprovisioned 1 --smp 1 --memory 256M
