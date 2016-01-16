VAGRANTFILE_API_VERSION = "2"
Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.ssh.forward_agent = true
  config.vm.box = "boxcutter/ubuntu1404"
  config.vm.synced_folder "~/.identity", "/home/vagrant/.identity", create: true
  config.vm.provision "shell", privileged: false, path: "https://s3-us-west-1.amazonaws.com/raptr-us-west-1/bootstrap"

  # box-specific
  config.vm.synced_folder File.join(ENV["GOPATH"],"/src"), "/home/vagrant/src"
  config.vm.provision "shell", inline: "apt-fast install -y golang glock go-packaging"
end
