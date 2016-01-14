# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/trusty64"

  config.vm.provision "shell", inline: <<-SHELL
    curl -sf -o /tmp/go1.5.3.linux-amd64.tar.gz -L https://storage.googleapis.com/golang/go1.5.3.linux-amd64.tar.gz
    sudo mkdir -p /opt && cd /opt && sudo tar xfz /tmp/go1.5.3.linux-amd64.tar.gz && rm -f /tmp/go1.5.3.linux-amd64.tar.gz
    curl -s https://packagecloud.io/install/repositories/darron/consul/script.deb.sh | sudo bash
    sudo apt-get install -y consul git make graphviz
    sudo cat > /etc/profile.d/go.sh << EOF
export GOROOT="/opt/go"
export GOPATH="/home/vagrant/gocode"
export PATH="/opt/go/bin://home/vagrant/gocode/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
export KVEXPRESS_DEBUG=1
EOF
  SHELL
end
