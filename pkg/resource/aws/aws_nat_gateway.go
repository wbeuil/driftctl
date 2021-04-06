// GENERATED, DO NOT EDIT THIS FILE
package aws

import (
	"github.com/zclconf/go-cty/cty"

	"github.com/cloudskiff/driftctl/pkg/dctlcty"
)

const AwsNatGatewayResourceType = "aws_nat_gateway"

type AwsNatGateway struct {
	AllocationId       *string           `cty:"allocation_id"`
	Id                 string            `cty:"id" computed:"true"`
	NetworkInterfaceId *string           `cty:"network_interface_id" computed:"true"`
	PrivateIp          *string           `cty:"private_ip" computed:"true"`
	PublicIp           *string           `cty:"public_ip" computed:"true"`
	SubnetId           *string           `cty:"subnet_id"`
	Tags               map[string]string `cty:"tags"`
	CtyVal             *cty.Value        `diff:"-"`
}

func (r *AwsNatGateway) TerraformId() string {
	return r.Id
}

func (r *AwsNatGateway) TerraformType() string {
	return AwsNatGatewayResourceType
}

func (r *AwsNatGateway) CtyValue() *cty.Value {
	return r.CtyVal
}

func awsNatGatewayNormalizer(val *dctlcty.CtyAttributes) {
}
