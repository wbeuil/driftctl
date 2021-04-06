// GENERATED, DO NOT EDIT THIS FILE
package aws

import (
	"github.com/zclconf/go-cty/cty"

	"github.com/cloudskiff/driftctl/pkg/dctlcty"
)

const AwsS3BucketMetricResourceType = "aws_s3_bucket_metric"

type AwsS3BucketMetric struct {
	Bucket *string `cty:"bucket"`
	Id     string  `cty:"id" computed:"true"`
	Name   *string `cty:"name"`
	Filter *[]struct {
		Prefix *string           `cty:"prefix"`
		Tags   map[string]string `cty:"tags"`
	} `cty:"filter"`
	CtyVal *cty.Value `diff:"-"`
}

func (r *AwsS3BucketMetric) TerraformId() string {
	return r.Id
}

func (r *AwsS3BucketMetric) TerraformType() string {
	return AwsS3BucketMetricResourceType
}

func (r *AwsS3BucketMetric) CtyValue() *cty.Value {
	return r.CtyVal
}

func awsS3BucketMetricNormalizer(val *dctlcty.CtyAttributes) {
}
