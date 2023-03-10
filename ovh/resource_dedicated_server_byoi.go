package ovh

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceDedicatedServerBringYourOwnImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceDedicatedServerBringYourOwnImageCreate,
		Update: resourceDedicatedServerBringYourOwnImageUpdate,
		Read:   resourceDedicatedServerBringYourOwnImageRead,
		Delete: resourceDedicatedServerBringYourOwnImageDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your dedicated server",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Image URL",
			},
			"checksum": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Image checksum",
			},
			"checksum_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Checksum type",
			},
			"configdrive_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable setting the ConfigDrive",
			},
			"configdrive_hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Set up the server using the provided hostname instead of the default hostname",
			},
			"configdrive_ssh_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "SSH key that should be installed. Password login will be disabled.",
			},
			"configdrive_user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Automatically provision your machine with additional software or settings",
			},
			"configdrive_user_metadatas": {
				Type:        schema.TypeList,
				Description: "Metadata",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Image description",
			},
			"disk_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Disk group id to process install on (only available for some templates)",
			},
			"http_header": {
				Type:        schema.TypeList,
				Description: "HTTP Headers",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Image type",
			},

			//Computed
			"last_checksum": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last checksum",
			},
			"last_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last message",
			},
			"servername": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server name",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "BYOI status",
			},
		},
	}
}

// AttachmentStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an Attachment Task.
func waitForDedicatedServerBringYourOwnImageDone(client *ovh.Client, serviceName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &DedicatedServerBringYourOwnImageResponse{}
		endpoint := fmt.Sprintf("/dedicated/server/%s/bringYourOwnImage", url.PathEscape(serviceName))
		if err := client.Get(endpoint, r); err != nil {
			return r, "", err
		}

		log.Printf("[DEBUG] Pending dedicated server install with bring your own image: %s", r)
		return r, r.Status, nil
	}
}

func resourceDedicatedServerBringYourOwnImageCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf(
		"/dedicated/server/%s/bringYourOwnImage",
		url.PathEscape(serviceName),
	)

	params := (&DedicatedServerBringYourOwnImageCreateOpts{}).FromResource(d)

	r := &DedicatedServerBringYourOwnImageResponse{}

	log.Printf("[DEBUG] Will install dedicated server with bring your own image: %s", params)

	if err := config.OVHClient.Post(endpoint, params, r); err != nil {
		return fmt.Errorf("Error calling POST %s:\n\t %q", endpoint, err)
	}

	log.Printf("[DEBUG] Waiting for dedicated server install with bring your own image %s:", r)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"doing"},
		Target:     []string{"done"},
		Refresh:    waitForDedicatedServerBringYourOwnImageDone(config.OVHClient, serviceName),
		Timeout:    45 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for dedicated server install with bring your own image (%s): %s", params, err)
	}
	log.Printf("[DEBUG] Created dedicated server install with bring your own image %s", r)

	return dedicatedServerBringYourOwnImageRead(d, meta)
}

func dedicatedServerBringYourOwnImageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	r := &DedicatedServerBringYourOwnImageResponse{}

	log.Printf("[DEBUG] Will read dedicated server with bring your own image: %s", serviceName)

	endpoint := fmt.Sprintf("/dedicated/server/%s/bringYourOwnImage", url.PathEscape(serviceName))

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("last_checksum", r.Checksum)
	d.Set("last_message", r.Message)
	d.Set("servername", r.Servername)
	d.Set("status", r.Status)

	log.Printf("[DEBUG] Read dedicated server install with bring your own image %s", r)
	return nil
}

func resourceDedicatedServerBringYourOwnImageUpdate(d *schema.ResourceData, meta interface{}) error {
	// nothing to do on update
	return resourceDedicatedServerBringYourOwnImageRead(d, meta)
}

func resourceDedicatedServerBringYourOwnImageRead(d *schema.ResourceData, meta interface{}) error {
	// Nothing to do on READ
	//
	// IMPORTANT: This resource doesn't represent a real resource
	// but instead a task on a dedicated server. OVH may clean its tasks database after a while
	// so that the API may return a 404 on a task id. If we hit a 404 on a READ, then
	// terraform will understand that it has to recreate the resource, and consequently
	// will trigger new install task on the dedicated server.
	// This is something we must avoid!
	//
	return nil
}

func resourceDedicatedServerBringYourOwnImageDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Will delete dedicated server install with bring your own image: %s", serviceName)

	endpoint := fmt.Sprintf("/dedicated/server/%s/bringYourOwnImage", url.PathEscape(serviceName))

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return fmt.Errorf("calling %s:\n\t %q", endpoint, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"doing"},
		Target:     []string{"done"},
		Refresh:    waitForDedicatedServerBringYourOwnImageDone(config.OVHClient, serviceName),
		Timeout:    45 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("deleting dedicated server install with bring your own image: %s", err)
	}

	log.Printf("[DEBUG] Deleted dedicated server install with bring your own image %s", serviceName)
	return nil
}
