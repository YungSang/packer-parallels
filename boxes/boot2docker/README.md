# Pack boot2docker.box for Parallels

```
$ make
$ vagrant box add boot2docker boot2docker.box --provider parallels
$ vagrant init boot2docker
$ vagrant up --provider parallels
```