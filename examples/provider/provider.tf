# We recommend to use environment variables instead of explicit assignment
provider "tilgangsportalen" {
  hosturl  = var.TILGANGSPORTALEN_URL
  username = var.TILGANGSPORTALEN_USERNAME
  password = var.TILGANGSPORTALEN_PASSWORD
}
