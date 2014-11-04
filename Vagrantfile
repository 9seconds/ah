# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  ENV['VAGRANT_DEFAULT_PROVIDER'] ||= 'docker'

  ["bash", "zsh", "fish"].each do |shell|
    config.vm.define shell do |vm|
      vm.ssh.username = "root"
      if Gem.win_platform?
        vm.vm.synced_folder ".", "/vagrant", type: "rsync"
      end

      vm.vm.provider "docker" do |docker|
        docker.name = shell
        docker.has_ssh = true
        docker.remains_running = true
        docker.build_dir = "./dockerfiles/#{shell}"
      end
    end
  end

end

