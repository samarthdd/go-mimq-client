<h1 align="center">MqDemo</h1>


# MqDemo

This app mimic  the icap service  int the diagram below 



###  build
```
git clone https://github.com/k8-proxy/go-mimq-client.git
cd mq-demo
go build
```

### build docker image for other services
info : for the moment use "git checkout develop" before docker build
```
git clone https://github.com/k8-proxy/go-k8s-process.git
cd go-k8s-process
docker build -t go-k8s-process .
```
```
git clone https://github.com/k8-proxy/go-k8s-srv1.git
cd go-k8s-srv1
docker build -t go-k8s-srv1 .
```
```
git clone https://github.com/k8-proxy/go-k8s-srv2.git
cd go-k8s-srv2
docker build -t go-k8s-srv2 .
```
- edit the docker-compose file to replace /path/to/shared/dir with one in your localmachine 

# Testing steps

in the mq-demo directory run 
  ```
-  docker-compose up  # wait 30 seconds for mq to initialize
-  ./mq-demo  file.pdf   # the file.pdf should be in  path/to/shared/dir
```
- the processed file will be saved in path/to/shared/dir/rebuild-file.pdf

# Rebuild flow to implement

![new-rebuild-flow-v2](https://github.com/k8-proxy/go-k8s-infra/raw/main/diagram/go-k8s-infra.png)

