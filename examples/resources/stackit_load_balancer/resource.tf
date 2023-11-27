

# Create a network
resource "openstack_networking_network_v2" "example" {
  name = "example_network"
}

# Create a subnet
resource "openstack_networking_subnet_v2" "example" {
  name            = "example_subnet"
  cidr            = "192.168.0.0/25"
  dns_nameservers = ["8.8.8.8"] // DNS is needed to reach the control plane
  network_id      = openstack_networking_network_v2.example.id
}


data "openstack_networking_network_v2" "public" {
  name = "floating-net"
}

# Create a floating IP
resource "openstack_networking_floatingip_v2" "example_ip" {
  pool = data.openstack_networking_network_v2.public.name
}

# Get flavor for instance
data "openstack_compute_flavor_v2" "example" {
  name = "g1.1"
}

# Create instance
resource "openstack_compute_instance_v2" "example" {
  depends_on      = [openstack_networking_subnet_v2.example]
  name            = "example_instance"
  flavor_id       = data.openstack_compute_flavor_v2.example.id
  admin_pass      = "example"
  security_groups = ["default"]

  block_device {
    uuid                  = "4364cdb2-dacd-429b-803e-f0f7cfde1c24" // Ubuntu 22.04
    volume_size           = 32
    source_type           = "image"
    destination_type      = "volume"
    delete_on_termination = true
  }

  network {
    name = openstack_networking_network_v2.example.name
  }
}

# Create a router and attach it to the public network
resource "openstack_networking_router_v2" "example_router" {
  name                = "example_router"
  admin_state_up      = "true"
  external_network_id = data.openstack_networking_network_v2.public.id
}

# Attach the subnet to the router
resource "openstack_networking_router_interface_v2" "example_interface" {
  router_id = openstack_networking_router_v2.example_router.id
  subnet_id = openstack_networking_subnet_v2.example.id
}

# create a load balancer
resource "stackit_load_balancer" "example" {
  project_id = ""
  name       = "example"
  target_pools = [{
    name        = "example-target-pool"
    target_port = 80
    targets = [{
      display_name = "example-target"
      ip_address   = openstack_compute_instance_v2.example.network.0.fixed_ip_v4
    }]
  }]
  listeners = [{
    display_name = "example-listener"
    port         = 80
    protocol     = "PROTOCOL_TCP"
    target_pool  = "example-target-pool"
  }]
  networks = [
    { network_id = openstack_networking_network_v2.example.id }
  ]
  external_address     = openstack_networking_floatingip_v2.example_ip.address
  private_network_only = false
}
