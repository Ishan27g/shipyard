Feature: Modules
  In order to test Modules
  I should apply a blueprint

  Scenario: Blueprint containing two modules
    Given the following environment variables are set
      | key            | value                 |
      | CONSUL_VERSION | 1.8.0                 |
      | ENVOY_VERSION  | 1.14.3                |
    And I have a running blueprint
    Then the following resources should be running
      | name                      | type              |
      | cloud                     | network           |
      | onprem                    | network           |
      | consul                    | container         |
      | envoy                     | sidecar           |
      | consul-container-http     | container_ingress |
      | consul-container-http-2   | container_ingress |
      | server.k3s                | k8s_cluster       |
      | docs                      | docs              |
    And a HTTP call to "http://consul.container.shipyard.run:8500/v1/agent/members" should result in status 200
    And a HTTP call to "http://consul-http.ingress.shipyard.run:18500/v1/agent/members" should result in status 200
