# go-binding tool to add the linux Capablities of running process or docker container


## Add capablities to container 
```
./gocapability -cap-add SYS_ADMIN,NET_ADMIN --name {container name}
```

## Add capabilites to the process

```
./gocapability -cap-add SYS_ADMIN,NET_ADMIN --pid {pid}
```