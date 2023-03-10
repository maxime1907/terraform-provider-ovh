package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type DedicatedServerBringYourOwnImageConfigDrive struct {
	Enable        bool     `json:"enable"`
	Hostname      string   `json:"hostname,omitempty"`
	SshKey        string   `json:"sshKey,omitempty"`
	UserData      string   `json:"userData,omitempty"`
	UserMetadatas []string `json:"userMetadatas,omitempty"`
}

func (p *DedicatedServerBringYourOwnImageConfigDrive) String() string {
	return fmt.Sprintf("enable: %t, hostname:%s, sshKey:%s, userData:%s, userMetadatas:%v", p.Enable, p.Hostname, p.SshKey, p.UserData, p.UserMetadatas)
}

type DedicatedServerBringYourOwnImageCreateOpts struct {
	URL          string                                       `json:"URL"`
	CheckSum     string                                       `json:"checkSum,omitempty"`
	CheckSumType string                                       `json:"checkSumType,omitempty"`
	Configdrive  *DedicatedServerBringYourOwnImageConfigDrive `json:"configdrive,omitempty"`
	Description  string                                       `json:"description,omitempty"`
	DiskGroupId  int                                          `json:"diskGroupId,omitempty"`
	HttpHeader   []string                                     `json:"httpHeader,omitempty"`
	Type         string                                       `json:"type"`
}

func (p *DedicatedServerBringYourOwnImageCreateOpts) String() string {
	return fmt.Sprintf("URL:%s, checkSum:%s, checkSumType:%s, configdrive:%s, description:%s, diskGroupId:%d, httpHeader:%v, type:%s", p.URL, p.CheckSum, p.CheckSumType, p.Configdrive, p.Description, p.DiskGroupId, p.HttpHeader, p.Type)
}

func (p *DedicatedServerBringYourOwnImageCreateOpts) FromResource(d *schema.ResourceData) *DedicatedServerBringYourOwnImageCreateOpts {
	params := &DedicatedServerBringYourOwnImageCreateOpts{
		URL:  d.Get("url").(string),
		Type: d.Get("type").(string),
	}
	if v, ok := d.GetOk("checksum"); ok && v.(string) != "" {
		params.CheckSum = v.(string)
	}
	if v, ok := d.GetOk("checksum_type"); ok && v.(string) != "" {
		params.CheckSumType = v.(string)
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		params.Description = v.(string)
	}
	if v, ok := d.GetOk("disk_group_id"); ok {
		params.DiskGroupId = v.(int)
	}

	httpHeader, _ := helpers.StringsFromSchema(d, "http_header")
	if httpHeader != nil {
		params.HttpHeader = httpHeader
	}

	configDriveEnable := d.Get("configdrive_enable").(bool)

	if configDriveEnable {
		configDrive := DedicatedServerBringYourOwnImageConfigDrive{
			Enable: configDriveEnable,
		}

		if v, ok := d.GetOk("configdrive_hostname"); ok && v.(string) != "" {
			configDrive.Hostname = v.(string)
		}
		if v, ok := d.GetOk("configdrive_ssh_key"); ok && v.(string) != "" {
			configDrive.SshKey = v.(string)
		}
		if v, ok := d.GetOk("configdrive_user_data"); ok && v.(string) != "" {
			configDrive.UserData = v.(string)
		}

		userMetadatas, _ := helpers.StringsFromSchema(d, "configdrive_user_metadatas")
		if userMetadatas != nil {
			configDrive.UserMetadatas = userMetadatas
		}

		params.Configdrive = &configDrive
	}

	return params
}

type DedicatedServerBringYourOwnImageResponse struct {
	Checksum   string `json:"checksum"`
	Message    string `json:"message"`
	Servername string `json:"servername"`
	Status     string `json:"status"`
}

func (p *DedicatedServerBringYourOwnImageResponse) String() string {
	return fmt.Sprintf("checksum: %s, message:%s, servername:%s, status:%s", p.Checksum, p.Message, p.Servername, p.Status)
}
