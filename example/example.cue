package example

import (
	"github.com/cueformproject/cueform/terraform"
	"github.com/cueformproject/cueform/example/providers/hashicorp/aws/5.58.0:aws"
)

infra: (terraform.#makeSchema & {providers: [aws]}).schema

// Providers config
infra: provider: aws: default_tags: [{tags: {managedBy: "cueform"}}]

// Localstack config
infra: provider: aws: {
	access_key: "test"
	secret_key: "test"
	region:     "us-east-1"

	s3_use_path_style:           false
	skip_credentials_validation: true
	skip_metadata_api_check:     "true"
	skip_requesting_account_id:  true

	endpoints: [{
		s3: "http://s3.localhost.localstack.cloud:4566"
	}]
}

// Bucket
infra: resource: aws_s3_bucket: myBucket: {
	bucket: "cueform-example"
}
