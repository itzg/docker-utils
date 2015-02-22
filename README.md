# docker-utils

[![status](https://sourcegraph.com/api/repos/github.com/itzg/docker-utils/.badges/status.svg)](https://sourcegraph.com/github.com/itzg/docker-utils)

Some simple utilities to manage Docker environments.

## purge_docker_volumes

Something I learned after quite a bit of Docker use and running out of space in the `/var/lib/docker` filesystem is that
`docker rm` by default does not remove the associated "vfs" volumes created via the Dockerfile's `VOLUME` declaration. 

For example,

    $ docker inspect -f "{{.Volumes}}" 9a8e6969eced
    map[/conf:/var/lib/docker/vfs/dir/77de7629c684968f98848372db3c490b28a0a22b94227ad7a8fc6cc55c55e16c /data:/var/lib/docker/vfs/dir/dde88b9a0d615ed1b411c142916f92cd5b5b4f829ad97f2e56d786db739bd83c]
    $ sudo ls /var/lib/docker/vfs/dir/77de7629c684968f98848372db3c490b28a0a22b94227ad7a8fc6cc55c55e16c
    elasticsearch.yml  logging.yml
    
Now, let's remove that container
    
    $ docker rm 9a8e6969eced
    9a8e6969eced
    $ sudo ls /var/lib/docker/vfs/dir/77de7629c684968f98848372db3c490b28a0a22b94227ad7a8fc6cc55c55e16c
    elasticsearch.yml  logging.yml

Ohh, it's still there. That's great if I wanted to get at the content later, but not if I have burned through containers without
knowing they had `VOLUME`s.

Safely purge these stale volumes by using:

    $ sudo ~/go/bin/purge_docker_volumes
    DELETING /var/lib/docker/vfs/dir/dde88b9a0d615ed1b411c142916f92cd5b5b4f829ad97f2e56d786db739bd83c
    DELETING /var/lib/docker/vfs/dir/77de7629c684968f98848372db3c490b28a0a22b94227ad7a8fc6cc55c55e16c

And it's safe to re-run it:

    $ sudo ~/go/bin/purge_docker_volumes
    Congrats, nothing to purge

### Building

Get it and install it:

    $ go get github.com/itzg/docker-utils/purge_docker_volumes
    $ go install github.com/itzg/docker-utils/purge_docker_volumes

and it run it from your `$GOPATH`/bin
