docs "docs" {
  path  = "./docs"
  port  = 8080
	open_in_browser = true

  network {
    name = "network.docs"
  }

  index_title = "Test"
  index_pages = ["index", "other"]
}