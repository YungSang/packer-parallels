# Pack coreos.box for Parallels

```
$ make
```

```
$ vagrant plugin install vagrant-parallels
$ vagrant box add coreos coreos.box --provider parallels
$ vagrant init coreos
$ vi Vagrantfile
```

```
VAGRANTFILE_API_VERSION = "2"

Vagrant.require_version ">= 1.5.0"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "coreos"

  config.vm.network "forwarded_port", guest: 4243, host: 4243

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
