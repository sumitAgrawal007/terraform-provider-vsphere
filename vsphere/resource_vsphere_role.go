package vsphere

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/vmware/govmomi/object"
	"log"
	"strconv"
	"strings"
)

func resourceVsphereRole() *schema.Resource {
	sch := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the storage policy.",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of the storage policy.",
		},
		"role_privileges": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "roles",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"label": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The display label of the Role.",
		},
	}

	return &schema.Resource{
		Create: resourceRoleCreate,
		Read:   resourceRoleRead,
		Update: resourceRoleUpdate,
		Delete: resourceRoleDelete,
		Schema: sch,
	}
}

func resourceRoleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Beginning create role %s", d.Get("name").(string))
	client := meta.(*VSphereClient).vimClient

	authorizationManager := object.NewAuthorizationManager(client.Client)

	name := d.Get("name").(string)
	rolePrivileges := structure.SliceInterfacesToStrings(d.Get("role_privileges").([]interface{}))

	roleId, err := authorizationManager.AddRole(context.Background(), name, rolePrivileges)
	if err != nil {
		return fmt.Errorf("error while creating role %s", err)
	}

	d.SetId(strconv.Itoa(int(roleId)))
	return resourceRoleRead(d, meta)
}

func resourceRoleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Reading vm storage policy profile", resourceVSphereVirtualMachineIDString(d))
	client := meta.(*VSphereClient).vimClient
	authorizationManager := object.NewAuthorizationManager(client.Client)
	i, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		panic(err)
	}
	roleId := int32(i)
	rlist, err := authorizationManager.RoleList(context.Background())

	if err != nil {
		return fmt.Errorf("error")
	}
	role := rlist.ById(roleId)
	if role == nil {
		log.Printf("Role %s has been deleted", d.Get("name"))
		d.SetId("")
		return nil
	}

	d.Set("name", role.Name)
	if role.Info != nil && role.Info.GetDescription() != nil {
		d.Set("description", role.Info.GetDescription().Summary)
		d.Set("label", role.Info.GetDescription().Label)
	}

	var arr []string
	for _, str := range role.Privilege {
		if strings.Split(str, ".")[0] != "System" {
			arr = append(arr, str)
		}
	}
	d.Set("role_privileges", arr)
	return nil
}

func resourceRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Beginning update role %s", d.Get("name").(string))
	client := meta.(*VSphereClient).vimClient

	authorizationManager := object.NewAuthorizationManager(client.Client)

	name := d.Get("name").(string)
	i, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		panic(err)
	}
	roleId := int32(i)
	rolePrivileges := structure.SliceInterfacesToStrings(d.Get("role_privileges").([]interface{}))

	err = authorizationManager.UpdateRole(context.Background(), roleId, name, rolePrivileges)
	if err != nil {
		return fmt.Errorf("error while updating role %s", err)
	}
	return resourceRoleRead(d, meta)
}

func resourceRoleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Performing Delete of Role with ID %s", d.Id())
	client := meta.(*VSphereClient).vimClient
	authorizationManager := object.NewAuthorizationManager(client.Client)

	i, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		panic(err)
	}
	roleId := int32(i)
	err = authorizationManager.RemoveRole(context.Background(), roleId, true)
	if err != nil {
		return fmt.Errorf("error while deleting role %s", err)
	}

	d.SetId("")
	log.Printf("[DEBUG] %s: Delete complete", d.Id())
	return nil
}
