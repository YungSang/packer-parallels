# Packer Builder Plugin for Parallels

This is a [Packer](http://www.packer.io/) Builder plugin for [Parallels Desktop for Mac](http://www.parallels.com/products/desktop/).

## How to Build the Plugin

	$ make

It will build and install `packer-builder-parallels-iso` and `packer-post-processor-vagrant` to make a box for Vagrant.

## How to Use (Pack a box)

- Prepare template.json for Packer with `"type": "parallels-iso"`.

	```
	{
		"builders": [{
			"type": "parallels-iso",
		.
		.
		.
	```

- Use Packer to build a pvm or a box

	`$ packer build template.json`

	Cf.) [A sample project for boot2docker](https://github.com/YungSang/packer-parallels/tree/boot2docker/boot2docker)