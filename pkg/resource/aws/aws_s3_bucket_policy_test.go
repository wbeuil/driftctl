package aws_test

import (
	"fmt"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/cloudskiff/driftctl/test/acceptance"
	"github.com/cloudskiff/driftctl/test/acceptance/awsutils"
)

func TestAcc_AwsS3BucketPolicy_WithFilter(t *testing.T) {
	var bucket string
	region := "us-east-1"
	acceptance.Run(t, acceptance.AccTestCase{
		Paths: []string{"./testdata/acc/aws_s3_bucket_policy_with_filter"},
		Args:  []string{"scan"},
		Checks: []acceptance.AccCheck{
			{
				Env: map[string]string{
					"AWS_REGION": region,
				},
				Args: func() []string {
					return []string{"--filter", "Type=='aws_s3_bucket_policy'"}
				},
				Check: func(result *acceptance.ScanResult, stdout string, err error) {
					if err != nil {
						t.Fatal(err)
					}
					result.AssertManagedCount(1)
					bucket = result.Analysis.Managed()[0].TerraformId()
					result.AssertResourceHasNoDrift(bucket, "aws_s3_bucket_policy")
				},
			},
			{
				Env: map[string]string{
					"AWS_REGION": region,
				},
				Args: func() []string {
					return []string{"--filter", fmt.Sprintf("Type=='aws_s3_bucket_policy' && Attr.bucket!='%s'", bucket)}
				},
				PreExec: func() {
					// Edit the policy to trigger a change on the aws_s3_bucket_policy that will be created through the middleware
					// Then we filter to exlude this resource and ensure there will not be any changes in the final output
					client := s3.New(awsutils.Session())
					policy := fmt.Sprintf("{\"Version\":\"2012-10-17\",\"Id\":\"MYBUCKETPOLICY\",\"Statement\":[{\"Sid\":\"IPAllow\",\"Effect\":\"Deny\",\"Principal\":\"*\",\"Action\":\"s3:*\",\"Resource\":[\"arn:aws:s3:::%s\",\"arn:aws:s3:::%s/*\"],\"Condition\":{\"IpAddress\":{\"aws:SourceIp\":\"1.1.1.1/32\"}}}]}", bucket, bucket)
					_, err := client.PutBucketPolicy(&s3.PutBucketPolicyInput{
						Bucket: awssdk.String(bucket),
						Policy: awssdk.String(policy),
					})
					if err != nil {
						t.Fatal(err)
					}
				},
				Check: func(result *acceptance.ScanResult, stdout string, err error) {
					if err != nil {
						t.Fatal(err)
					}
					result.AssertManagedCount(0)
					result.AssertResourceHasNoDrift(bucket, "aws_s3_bucket_policy")
				},
			},
		},
	})
}
