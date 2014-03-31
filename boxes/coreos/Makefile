coreos.box: coreos.iso template.json vagrantfile.tpl \
	tmp/insecure_private_key tmp/override-plugin.rb \
	oem/coreos-install oem/cloud-config.yml
	packer build template.json

coreos.iso:
	curl -LO https://github.com/YungSang/coreos-live-iso/releases/download/iso%2Fv0.1.0/coreos.iso

tmp/insecure_private_key:
	mkdir -p tmp
	curl -L https://raw.github.com/mitchellh/vagrant/master/keys/vagrant -o tmp/insecure_private_key

oem/coreos-install:
	mkdir -p oem
	curl -L https://raw.github.com/coreos/init/master/bin/coreos-install -o oem/coreos-install
	sed -e "s/amd64-generic/amd64-usr/g" -i "" oem/coreos-install
	sed -e "s/partprobe/partprobe \|\| true/g" -i "" oem/coreos-install

oem/cloud-config.yml:
	mkdir -p oem
	curl -L https://raw.github.com/coreos/coreos-overlay/master/coreos-base/oem-vagrant/files/cloud-config.yml -o oem/cloud-config.yml

tmp/override-plugin.rb:
	mkdir -p tmp
	cd tmp; \
	curl -LO https://raw.github.com/coreos/coreos-vagrant/master/override-plugin.rb

clean:
	rm -f coreos.iso
	rm -f coreos.box
	rm -rf oem/
	rm -rf tmp/
	rm -rf output-*/

.PHONY: clean
