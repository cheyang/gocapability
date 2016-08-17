# go-binding tool to add the linux Capablities of running process or docker container


## Add capablities to container 
```
./gocapability -cap-add SYS_ADMIN,NET_ADMIN --name {container name}
```

## Add capabilites to the process

```
./gocapability -cap-add SYS_ADMIN,NET_ADMIN --pid {pid}
```

Notice: 

The Capabilities of process can only be applied by the current process.


But updating caps by changing cgroup will be supported in kernel level: https://patchwork.kernel.org/patch/9186239/