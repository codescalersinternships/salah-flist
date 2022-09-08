# Flist

The flist file format is a general purpose format to store metadata about a (posix) filesystem. It's main goal is keeping a small file with enough information to make a complete filesystem available without the data payload itself, in an efficient way.

`Flist` is an app that provides the ability to run an application in a loosely isolated environment called a container.

## Design

Flist uses a client-server architecture. The Flist client talks to the Flist daemon, which does the heavy lifting of building, running, and handling your Flist containers. The Docker client and daemon communicate over UNIX sockets.

## How to use flist?

1- Build daemon and client:

```bash
go build -o flistd ./daemon/*.go
go build -o flist ./client/*.go
```

2- Run daemon as background process:

```bash
sudo ./flistd
```

3- Run client with available sub-commands

## Avilable sub-command

- To run entrypoint from an flist

```bash
 sudo ./flist run META  ENTRYPOINT
```

- To list running containers

```bash
sudo ./flist ps
```

- To stop container

```bash
sudo ./flist stop CONTAINER_ID
```

- To remove container

```bash
sudo ./flist rm CONTAINER_ID
```
