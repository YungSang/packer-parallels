# Packer Builder Plugin for Parallels

This is a [Packer](http://www.packer.io/) Builder plugin for [Parallels Desktop for Mac](http://www.parallels.com/products/desktop/).

## How to Build the Plugin

	$ make

It will build and install `packer-builder-parallels-iso`.

## How to Use (Pack a box)

- Prepare template.json for Packer with `"type": "parallels-iso"`.

	```
	{
		"builders": [{
			"type": "parallels-iso",
		.
		.
		.
		You can use "provisioners" but "post-processors".
	```

- Use Packer to build a pvm

	`$ packer build template.json`

	It will make `./output-parallels-iso/packer-parallels-iso.pvm`

- Pack a box with the pvm.  
	Packer-Parallels Builder plugin can create a pvm file only but a box directly, because we can not use Packer's post-processor with custom builders. The post-processor seems hard coded, not pluggable.
So you have to pack a box with the pvm, metadata.json and Vagrantfile manually.

	```
	clean up pvm
	$ cd ./output-parallels-iso
	$ rm -rf ./packer-parallels-iso.pvm/Snapshots
	$ rm -f  ./packer-parallels-iso.pvm/*.log
	$ rm -f  ./packer-parallels-iso.pvm/*.backup
	$ rm -f  ./packer-parallels-iso.pvm/harddisk.hdd/*.Backup

	pack a box
	$ cp ../metadata.json .

	If you need a default initialization, you can put Vagrantfile.
	$ vi Vagrantfile

	If you need more files, you can copy them here.

	$ tar zcvf ../<box-name>.box *
	```
