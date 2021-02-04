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