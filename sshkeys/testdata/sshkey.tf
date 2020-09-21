
resource "tls_private_key" "test" {
  algorithm = "RSA"
}

resource "exoscale_ssh_keypair" "test" {
  name = "test"
  public_key = replace(tls_private_key.test.public_key_openssh, "\n", "")
}
