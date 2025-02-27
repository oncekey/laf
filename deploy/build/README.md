
## **WARNING!**
> **The image in this directory is temporarily unavailable because SealOS has some inadequate support. Stay here for the time being, and then use it when the time is mature.**

## Intro

This directory contains sealos cluster image building scripts for laf.

## Build

> This would be build automatically in github action. see [build_cluster_images.yml](../.github/workflows/build_cluster_images.yml)

## Usage Images

```bash
# Install laf controllers
sealos run lafyun/laf-controllers:dev

# Install casdoor & service auth in sealos namespace (because service-auth has hard-coding `sealos` namespace in code)
export NODE_IP=$(kubectl get nodes -o jsonpath="{.items[0].status.addresses[0].address}")
export CASDOOR_NODE_PORT=30080
export CASDOOR_ENDPOINT=http://$NODE_IP:$CASDOOR_NODE_PORT
export CASDOOR_CALLBACK_URL=http://localhost:8080/login/callback

kubectl create namespace laf
sealos run --env NAMESPACE=laf --env DATABASE=casdoor lafyun/laf-postgresql:dev
sealos run --env NAMESPACE=laf --env NODE_PORT=${CASDOOR_NODE_PORT} lafyun/laf-casdoor:dev 


# Install mongodb
kubectl create namespace laf
sealos run --env NAMESPACE=laf lafyun/laf-mongodb:dev
```

