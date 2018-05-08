# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "alpine/alpine64"
  config.vm.network "forwarded_port", guest: 80, host: 8080, host_ip: "127.0.0.1"
end
