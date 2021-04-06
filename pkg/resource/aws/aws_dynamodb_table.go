// GENERATED, DO NOT EDIT THIS FILE
package aws

import (
	"github.com/zclconf/go-cty/cty"

	"github.com/cloudskiff/driftctl/pkg/dctlcty"
)

const AwsDynamodbTableResourceType = "aws_dynamodb_table"

type AwsDynamodbTable struct {
	Arn            *string           `cty:"arn" computed:"true"`
	BillingMode    *string           `cty:"billing_mode"`
	HashKey        *string           `cty:"hash_key"`
	Id             string            `cty:"id" computed:"true"`
	Name           *string           `cty:"name"`
	RangeKey       *string           `cty:"range_key"`
	ReadCapacity   *int              `cty:"read_capacity"`
	StreamArn      *string           `cty:"stream_arn" computed:"true"`
	StreamEnabled  *bool             `cty:"stream_enabled"`
	StreamLabel    *string           `cty:"stream_label" computed:"true"`
	StreamViewType *string           `cty:"stream_view_type" computed:"true"`
	Tags           map[string]string `cty:"tags"`
	WriteCapacity  *int              `cty:"write_capacity"`
	Attribute      *[]struct {
		Name *string `cty:"name"`
		Type *string `cty:"type"`
	} `cty:"attribute"`
	GlobalSecondaryIndex *[]struct {
		HashKey          *string  `cty:"hash_key"`
		Name             *string  `cty:"name"`
		NonKeyAttributes []string `cty:"non_key_attributes"`
		ProjectionType   *string  `cty:"projection_type"`
		RangeKey         *string  `cty:"range_key"`
		ReadCapacity     *int     `cty:"read_capacity"`
		WriteCapacity    *int     `cty:"write_capacity"`
	} `cty:"global_secondary_index"`
	LocalSecondaryIndex *[]struct {
		Name             *string  `cty:"name"`
		NonKeyAttributes []string `cty:"non_key_attributes"`
		ProjectionType   *string  `cty:"projection_type"`
		RangeKey         *string  `cty:"range_key"`
	} `cty:"local_secondary_index"`
	PointInTimeRecovery *[]struct {
		Enabled *bool `cty:"enabled"`
	} `cty:"point_in_time_recovery"`
	Replica *[]struct {
		RegionName *string `cty:"region_name"`
	} `cty:"replica"`
	ServerSideEncryption *[]struct {
		Enabled   *bool   `cty:"enabled"`
		KmsKeyArn *string `cty:"kms_key_arn" computed:"true"`
	} `cty:"server_side_encryption"`
	Timeouts *struct {
		Create *string `cty:"create"`
		Delete *string `cty:"delete"`
		Update *string `cty:"update"`
	} `cty:"timeouts" diff:"-"`
	Ttl *[]struct {
		AttributeName *string `cty:"attribute_name"`
		Enabled       *bool   `cty:"enabled"`
	} `cty:"ttl"`
	CtyVal *cty.Value `diff:"-"`
}

func (r *AwsDynamodbTable) TerraformId() string {
	return r.Id
}

func (r *AwsDynamodbTable) TerraformType() string {
	return AwsDynamodbTableResourceType
}

func (r *AwsDynamodbTable) CtyValue() *cty.Value {
	return r.CtyVal
}

func awsDynamodbTableNormalizer(val *dctlcty.CtyAttributes) {
	val.SafeDelete([]string{"timeouts"})
}
