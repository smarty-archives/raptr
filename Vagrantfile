# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.ssh.forward_agent = true
  config.vm.box = "box-cutter/ubuntu1404"
  config.vm.provision "file", source: "~/.ssh/known_hosts", destination: ".ssh/known_hosts"
  config.vm.provision "file", source: "~/.gitconfig", destination: ".gitconfig"
  config.vm.provision "shell", inline: "apt-get update && apt-get install -y wget curl axel htop vim debhelper git-core"
  config.vm.provision "shell", inline: "curl --silent -L http://golang.org/dl/go1.3.1.linux-amd64.tar.gz | tar xvz --owner root --group root -C /usr/local && ln -s /usr/local/go/bin/* /usr/local/bin"
  config.vm.synced_folder File.join(ENV["GOPATH"],"/src"), "/home/vagrant/src"
  config.vm.synced_folder "~/.identity", "/home/vagrant/.identity"
  config.vm.provision "shell", inline: "echo 'test -d ~/.identity && for f in ~/.identity/*; do . $f; done' >> /home/vagrant/.bashrc"
  config.vm.provider "virtualbox" do |vb|
    vb.customize ["modifyvm", :id, "--natdnshostresolver1", "on"]
    vb.customize ["modifyvm", :id, "--natdnsproxy1", "on"]
  end
end
