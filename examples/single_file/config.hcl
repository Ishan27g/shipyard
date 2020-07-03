network "onprem" {
  subnet = "10.6.0.0/16"
}

container "consul" {
  image   {
    name = "consul:1.6.1"
  }

  command = ["consul", "agent", "-config-file=/config/consul.hcl"]

  volume {
    source      = "./consul_config"
    destination = "/config"
  }

  network   {
    name = "network.onprem"
    ip_address = "10.6.0.200"
  }
  
  network   {
    name = "network.cloud"
  }


  resources {
    # Max CPU to consume, 1024 is one core, default unlimited
    cpu = 2048
    # Pin container to specified CPU cores, default all cores
    cpu_pin = [1,2]
    # max memory in MB to consume, default unlimited
    memory = 1024
  }
}