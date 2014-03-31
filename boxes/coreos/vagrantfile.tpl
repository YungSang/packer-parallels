require_relative "override-plugin.rb"

Vagrant.configure("2") do |config|
  # SSH in as the default 'core' user, it has the vagrant ssh key.
  config.ssh.username = "core"

  # Disable the base shared folder, Guest Additions are unavailable.
  config.vm.synced_folder ".", "/vagrant", disabled: true

  # Attach the b2d ISO so that it can boot
  config.vm.provider :parallels do |p|
    # Guest Tools are unavailable.
    p.check_guest_tools = false
  end
end
