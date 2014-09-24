# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.ssh.forward_agent = true
  config.vm.box = "box-cutter/ubuntu1404"
  config.vm.provision "file", source: "~/.ssh/known_hosts", destination: ".ssh/known_hosts"
  config.vm.provision "file", source: "~/.gitconfig", destination: ".gitconfig"
  config.vm.provision "file", source: "~/.profile-aws", destination: ".profile-aws"
  config.vm.provision "shell", inline: "echo 'for f in ~/.profile-*; do . $f; done' >> /home/vagrant/.bashrc"
  config.vm.provision "shell", inline: "apt-get update && apt-get install -y wget curl axel htop vim debhelper git-core"
  config.vm.provision "shell", path: "https://smartystreets-artifacts-us-east-1.s3.amazonaws.com/os-images/vagrant-images/install-scripts/latest-golang.sh"
  config.vm.synced_folder File.join(ENV["GOPATH"],"/src"), "/home/vagrant/src"
  config.vm.provider "virtualbox" do |vb|
    vb.customize ["modifyvm", :id, "--natdnshostresolver1", "on"]
    vb.customize ["modifyvm", :id, "--natdnsproxy1", "on"]
  end
end
