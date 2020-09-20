data "exoscale_compute_template" "ubuntu" {
  zone = "ch-gva-2"
  name = "Linux Ubuntu 18.04 LTS 64-bit"
}

resource "exoscale_instance_pool" "test" {
  zone = "at-vie-1"
  name = "test"
  template_id = data.exoscale_compute_template.ubuntu.id
  size = 1
  service_offering = "micro"
  disk_size = 10
  key_pair = ""
}