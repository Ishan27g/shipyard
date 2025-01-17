resource "k8s_cluster" "dev" {
  // use a random port for the cluster
  api_port = random_number(10000, 20000)

  network {
    id = variable.network_id
  }
}

resource "helm" "consul" {
  cluster = resource.k8s_cluster.dev.id
  chart   = "github.com/hashicorp/consul-k8s?ref=${variable.consul_helm_version}//charts/consul"
  values  = "./helm/consul-values.yaml"

  health_check {
    timeout = "240s"
    pods    = ["release=consul"]
  }
}

resource "ingress" "consul_http" {
  port = variable.consul_port

  target {
    id   = resource.k8s_cluster.dev.id
    port = 8500

    config = {
      service   = "consul-consul-server"
      namespace = "default"
    }
  }
}