# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.ssh.forward_agent = true
  config.vm.provider "virtualbox" do |vb|
    vb.customize ["modifyvm", :id, "--natdnshostresolver1", "on"]
    vb.customize ["modifyvm", :id, "--natdnsproxy1", "on"]
  end
  config.vm.box = "boxcutter/ubuntu1404"
  config.vm.synced_folder File.join(ENV["GOPATH"],"/src"), "/home/vagrant/src"
  config.vm.synced_folder "~/.identity", "/home/vagrant/.identity"
  config.vm.provision "file", source: "~/.ssh/known_hosts", destination: ".ssh/known_hosts"
  config.vm.provision "file", source: "~/.gitconfig", destination: ".gitconfig"
  config.vm.provision "shell", privileged: false, inline: "sed -i '1i test -d ~/.identity && for f in ~/.identity/*; do . $f; done' ~/.bashrc"
  config.vm.provision "shell", privileged: false, path: "https://s3-us-west-1.amazonaws.com/raptr-us-west-1/bootstrap"
  config.vm.provision "shell", inline: "apt-fast install -y wget curl axel htop vim debhelper git-core mercurial"
  config.vm.provision "shell", inline: "curl --silent -L https://storage.googleapis.com/golang/go1.3.3.linux-amd64.tar.gz | tar xvz --owner root --group root -C /usr/local && ln -s /usr/local/go/bin/* /usr/local/bin"
end
