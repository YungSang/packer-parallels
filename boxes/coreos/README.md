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

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "coreos"
end
```

```
$ vagrant up --provider parallels
```
