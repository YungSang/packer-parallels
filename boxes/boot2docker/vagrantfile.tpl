docker_daemon_running = Vagrant.source_root.join("plugins/provisioners/docker/cap/linux/docker_daemon_running.rb")
if File.exist?(docker_daemon_running)
  require docker_daemon_running

  module VagrantPlugins
    module Docker
      module Cap
        module Linux
          module DockerDaemonRunning
            def self.docker_daemon_running(machine)
              machine.communicate.test("test -f /var/run/docker.pid")
            end
          end
        end
      end
    end
  end
end

Vagrant.configure("2") do |config|
  config.ssh.shell = "sh -l"
  config.ssh.username = "docker"

  # Attach the b2d ISO so that it can boot
  config.vm.provider "parallels" do |p|
    p.check_guest_tools = false
    p.customize "pre-boot", [
      "set", :id,
      "--device-set", "cdrom0",
      "--image", File.expand_path("../boot2docker.iso", __FILE__),
      "--enable", "--connect"
    ]
  end
end
