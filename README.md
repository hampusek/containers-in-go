From [Liz Rice talk at GOTO 2018](https://www.youtube.com/watch?v=8fi7uSYlOdc).

## Running the container

```bash
go run main.go run /bin/bash
```

require root privileges to run.


## Getting the "base image"

```bash
docker run -d --rm --name ubuntufs ubuntu:20.04 sleep 1000
docker export ubuntufs -o ubuntufs.tar
docker stop ubuntufs
mkdir -p "$HOME"/ubuntufs
tar xf ubuntufs.tar -C "$HOME"/ubuntufs/
```
