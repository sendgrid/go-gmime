# -*- mode: ruby -*-
# vi: set ft=ruby :

require "vagrant"
require "berkshelf"


if Vagrant::VERSION < "1.2.1"
  raise "This is only compatible with Vagrant 1.2.1+"
end

project_name = "go-gmime"

Vagrant.configure("2") do |config|
  config.vm.hostname = "#{project_name}-build"

  # Enable the berkshelf-vagrant plugin
  config.berkshelf.enabled = true
  config.berkshelf.berksfile_path = "./Berksfile"

  config.ssh.forward_agent = true
  host_project_path = Dir.getwd
  guest_project_path = "/home/vagrant/sendgrid/#{project_name}"

  config.vm.synced_folder host_project_path, guest_project_path, nfs: true

  # SendGrid CentOs 6 box
  config.vm.define 'centos', primary: true do |c|
    c.berkshelf.berksfile_path = "./Berksfile"
    c.vm.network "private_network", ip: "192.168.100.2"
    c.vm.box = "go-gmime-centos"
    c.vm.box_url = "http://developer.nrel.gov/downloads/vagrant-boxes/CentOS-6.4-x86_64-v20131103.box"
    c.vm.provision "shell", inline: "cd #{guest_project_path}; bin/prepare-environment"
  end

  # Ubuntu box:
  config.vm.define 'ubuntu' do |c|
      c.berkshelf.berksfile_path = "./Berksfile"
      c.vm.network "private_network", ip: "192.168.100.3"
      c.vm.box = "go-gmime-ubuntu-14.04-02222014"
#      c.vm.synced_folder "../../../../", "/vagrant", type: "nfs"
      c.vm.box_url = "http://cloud-images.ubuntu.com/vagrant/trusty/current/trusty-server-cloudimg-amd64-vagrant-disk1.box"
    c.vm.provision "shell", inline: "cd #{guest_project_path}; bin/prepare-environment"
  end

  config.vm.provider :virtualbox do |vb|
    # Give enough horsepower to build without taking all day.
    vb.customize [
      "modifyvm", :id,
      "--memory", "1536",
      "--cpus", "2"
    ]
  end

  config.vm.provision :chef_solo do |chef|
    chef.add_recipe "golang"

    chef.json = {
        "go" => {
            # TODO: parse the package dependency file in bin/install
            "packages" => [
                'github.com/stretchr/testify/assert',
                'github.com/djimenez/iconv-go'
            ]
        }
    }
  end
end
