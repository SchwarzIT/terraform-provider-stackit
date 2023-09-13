# create a load balancer
resource "stackit_load_balancer" "example" {
  project_id = var.project_id
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
