Feature: Docker Container
  In order to test Docker containers
  I should apply a blueprint

  Scenario: Single Container from Local Blueprint
    Given I apply the config
    Then there should be 1 network called "onprem"
    And there should be 1 container running called "consul.container.shipyard.run"
    And there should be 1 container running called "consul.sidecar.shipyard.run"
    And there should be 1 container running called "consul-container-http.ingress.shipyard.run"
    And a call to "http://consul-http.ingress.shipyard.run:8500/v1/agent/members" should result in status 200