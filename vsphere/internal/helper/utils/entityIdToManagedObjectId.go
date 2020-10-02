package utils

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"
)

const VM = "VirtualMachine"
const DISTRIBUTED_VIRTUAL_SWITCH = "VmwareDistributedVirtualSwitch"

func GetMoid(client *govmomi.Client, entityType string, id string) (string, error) {

	switch entityType {
	case VM:
		vm, err := virtualmachine.FromUUID(client, id)
		if err != nil {
			return "", fmt.Errorf("error while finiding vm with id %s %s", id, err)
		}
		return vm.Reference().Value, nil
	case DISTRIBUTED_VIRTUAL_SWITCH:
		dvsm := types.ManagedObjectReference{Type: "DistributedVirtualSwitchManager", Value: "DVSManager"}
		req := &types.QueryDvsByUuid{
			This: dvsm,
			Uuid: id,
		}
		resp, err := methods.QueryDvsByUuid(context.TODO(), client, req)
		if err != nil {
			return "", fmt.Errorf("error while finiding distributed virtual switch with id %s %s", id, err)
		}
		return resp.Returnval.Reference().Value, nil
	default:
		return id, nil
	}
}
