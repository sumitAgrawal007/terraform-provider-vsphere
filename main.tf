provider "vsphere" {
  user           = "administrator@vsphere.local"
  password       = "Admin!23"
  vsphere_server = "sc2-10-186-4-243"

  # If you have a self-signed cert
  allow_unverified_ssl = true
}

data "vsphere_datacenter" "dc" {
  name = "Sample_DC_2"
}

data "vsphere_host" h2 {
  name = "10.186.15.23"
  datacenter_id = data.vsphere_datacenter.dc.id

}

data "vsphere_role" "role1" {
  label = "Datastore consumer (sample)"
}


data "vsphere_role" "role2" {
  label = "Virtual machine user (sample)"
}

data "vsphere_virtual_machine" vm1 {
  name = "Sample_Exhaustive_VM_for_Simple_Testbed"
  datacenter_id = data.vsphere_datacenter.dc.id
}

resource "vsphere_entity_permission" p2 {
  entity_id = data.vsphere_virtual_machine.vm1.id
  entity_type = "VirtualMachine"
  permissions {
    user_or_group = "vsphere.local\\ExternalIDPUsers"
    propagate = true
    is_group = true
    role_id = data.vsphere_role.role1.id
  }
  permissions {
    user_or_group = "vsphere.local\\DCClients"
    propagate = true
    is_group = true
    role_id = data.vsphere_role.role2.id
  }
}
