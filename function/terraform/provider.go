package terraform

import (
	"github.com/gofunct/goscript/api/tf"
	"github.com/hashicorp/terraform/terraform"
)

type GrpcProviderFunc func() tf.ProviderServer
type ProviderFunc func() terraform.ResourceProvider
