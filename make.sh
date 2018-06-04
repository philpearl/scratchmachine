#!/bin/bash

# Build the go binary for linux
GOOS=linux go build

# Creating the vagrant box runs the vagrant provisioner, which creates the ISO.
vagrant destroy -f
vagrant up

# This is only useful if you've configured a VM "fred" that uses the ISO
VBoxHeadless -s fred