# encryptah
This is a two-tier application written in Go (backend) and Python (frontend). It takes in a user's input string and encrypts it with an AES-256 key and returns back the ciphertext.

![screenshot](images/screenshot.png)

## Build Docker Images Manually
```bash
# build the docker image for the backend
$ docker build --tag jacobmammoliti/encryptah-be:0.1 src/backend/

# push it to dockerhub
$ docker push jacobmammoliti/encryptah-be:0.1

# build the docker image for the frontend
$ docker build --tag jacobmammoliti/encryptah-fe:0.1 src/frontend/

# push it to dockerhub
$ docker push jacobmammoliti/encryptah-fe:0.1
```

## Build Multi-Arch Docker Images
```bash
# Create a new builder instance
$ docker buildx create --use --name=mult-arch
mult-arch

# Inspect the builder instance and bootstrap it
$ docker buildx inspect --bootstrap
[+] Building 13.3s (1/1) FINISHED                                                                                                   
 => [internal] booting buildkit                                                                                               13.3s
 => => pulling image moby/buildkit:buildx-stable-1                                                                            10.8s
 => => creating container buildx_buildkit_test0                                                                                2.5s
Name:   mult-arch
Driver: docker-container

Nodes:
Name:      mult-arch0
Endpoint:  unix:///var/run/docker.sock
Status:    running
Platforms: linux/amd64, linux/arm64, linux/riscv64, linux/ppc64le, linux/s390x, linux/386, linux/arm/v7, linux/arm/v6

# Build the frontend image
$ docker buildx build --push -t jacobmammoliti/encryptah-frontend:1.0 --platform linux/amd64,linux/arm64 src/frontend/
[+] Building 18.3s (14/14) FINISHED
...
 => [auth] jacobmammoliti/frontend:pull,push token for registry-1.docker.io

# Build the backend image
$  docker buildx build --push -t jacobmammoliti/encryptah-backend:1.0 --platform linux/amd64,linux/arm64 src/backend/
[+] Building 27.4s (19/19) FINISHED
...
 => [auth] jacobmammoliti/backend:pull,push token for registry-1.docker.io     
```

## Deploy to Kubernetes
```bash
# deploy the backend
$ kubectl apply -f kubernetes/encryptah-be.yaml
serviceaccount/encryptah-fe-sa created
pod/encryptah-be created
service/encryptah-be created

# deploy the frontend
$ kubectl apply -f kubernetes/encryptah-fe.yaml
serviceaccount/encryptah-fe-sa unchanged
pod/encryptah-fe created
service/encryptah-fe created
```