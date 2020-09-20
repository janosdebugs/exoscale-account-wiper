

resource "exoscale_security_group" "test" {
  name = "test"
}

resource "exoscale_security_group_rules" "test" {
  security_group = exoscale_security_group.test.name

  ingress {
    ports = ["1"]
    protocol = "tcp"
    description = "test"
    user_security_group_list = ["default"]
  }
}
