package cloudinit

import (
	"fmt"
	"github.com/danielerez/openshift-appliance/pkg/genisoimage"
)

type CloudInit struct {
	SSHKeyCustomer string
	SSHKeyHost     string
}

func (c *CloudInit) GenerateCloudInit() string {
	return fmt.Sprintf(`#cloud-config
		ssh_authorized_keys:
		`)
}

func genISOImage(s string) string {
	var genIsoImage genisoimage.GenIsoImage = genisoimage.NewGenIsoImage()
	genIsoImage.GenerateDataImage("imagePath", "dataPath")

	return s
}

func (c *CloudInit) GenerateCloudInitISO() string {
	return genISOImage(c.GenerateCloudInit())
}
