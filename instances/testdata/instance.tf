data "exoscale_compute_template" "ubuntu" {
  zone = "ch-gva-2"
  name = "Linux Ubuntu 18.04 LTS 64-bit"
}

resource "exoscale_compute" "test" {
  zone         = "ch-gva-2"
  display_name = "mymachine"
  template_id  = data.exoscale_compute_template.ubuntu.id
  size         = "Micro"
  disk_size    = 10
  key_pair     = ""
  state        = "Running"

  affinity_groups = []
  security_groups = ["default"]

}