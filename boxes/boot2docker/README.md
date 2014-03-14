# Pack boot2docker.box for Parallels

```
$ make
$ vagrant box add boot2docker boot2docker.box --provider parallels
$ vagrant init boot2docker
$ vi Vagrantfile
```

```
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "boot2docker"

  config.vm.network "private_network", ip: "192.168.34.10"

  config.vm.synced_folder ".", "/vagrant", type: "nfs"

  config.vm.provision :docker do |d|
    d.pull_images "busybox"
    d.run "busybox",
      cmd: "echo hello"
  end
end
```

```
$ vagrant up --provider parallels
```
