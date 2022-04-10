package acceptance

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// ImportStep returns a Test Step which Imports the Resource
func (td TestData) ImportStep() resource.TestStep {
	return resource.TestStep{
		ResourceName:      td.ResourceName,
		ImportState:       true,
		ImportStateVerify: true,
	}
}
