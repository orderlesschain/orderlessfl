## Work in Progress

variable "nodes" {
  type    = list(string)
  default = [
    "172.24.33.3",
    "172.24.33.36",
    "172.24.33.63",
    "172.24.33.86",
    "172.24.33.90",
    "172.24.33.91",
    "172.24.33.92",
    "172.24.33.93",
    #    "172.24.33.94",
    #    "172.24.33.95",
    #    "172.24.33.96",
    #    "172.24.33.97",
    #    "172.24.33.98",
    #    "172.24.33.99",
    #    "172.24.33.100",
    #    "172.24.33.101",
  ]
}

variable "clients" {
  type    = list(string)
  default = [
    "172.24.33.134",
    "172.24.33.135",
    "172.24.33.136",
    "172.24.33.133",
    "172.24.33.102",
    "172.24.33.103",
    "172.24.33.104",
    "172.24.33.105",
    "172.24.33.106",
    "172.24.33.107",
    "172.24.33.108",
    "172.24.33.109",
    "172.24.33.110",
    "172.24.33.111",
    "172.24.33.112",
    "172.24.33.113",
    "172.24.33.114",
    "172.24.33.115",
    "172.24.33.116",
    "172.24.33.117",
    "172.24.33.118",
    "172.24.33.119",
    "172.24.33.120",
  ]
}

variable "clis" {
  type    = list(string)
  default = ["172.24.33.167"]
}


#### The Ansible inventory file ############################################
resource "local_file" "ansible_inventory" {
  content = templatefile("./templates/inventory.tmpl",
    {
      nodes-ips   = var.nodes
      clients-ips = var.clients
      cli-ips     = var.clis
    })
  filename = "../ansible/ansible_remote/inventory_dir/inventory"
}

### The Nodes Endpoints
resource "local_file" "nodes_endpoints" {
  content = templatefile("./templates/endpoints.tmpl",
    {
      nodes-ips   = var.nodes
      clients-ips = var.clients
      cli-ips     = var.clis
    })
  filename = "../../configs/endpoints_remote.yml"
}

### The Nodes Endpoints for creating certificates
resource "local_file" "nodes_endpoints_certificate" {
  content = templatefile("./templates/endpoints_certificate.tmpl",
    {
      nodes-ips   = var.nodes
      clients-ips = var.clients
    })
  filename = "../../certificates/endpoints_remote"
}
