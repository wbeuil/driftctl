package vsphere

import (
	"github.com/cloudskiff/driftctl/pkg/alerter"
	"github.com/cloudskiff/driftctl/pkg/output"
	"github.com/cloudskiff/driftctl/pkg/remote/cache"
	"github.com/cloudskiff/driftctl/pkg/remote/common"
	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/cloudskiff/driftctl/pkg/terraform"
)

const RemoteVSphereTerraform = "vsphere+tf"

/**
 * Initialize remote (configure credentials, launch tf providers and start gRPC clients)
 * Required to use Scanner
 */

func Init(version string, alerter *alerter.Alerter,
	providerLibrary *terraform.ProviderLibrary,
	supplierLibrary *resource.SupplierLibrary,
	remoteLibrary *common.RemoteLibrary,
	progress output.Progress,
	resourceSchemaRepository *resource.SchemaRepository,
	factory resource.ResourceFactory,
	configDir string) error {
	if version == "" {
		version = "2.0.1"
	}

	provider, err := NewVSphereTerraformProvider(version, progress, configDir)
	if err != nil {
		return err
	}
	err = provider.Init()
	if err != nil {
		return err
	}

	repositoryCache := cache.New(100)

	repository, err := NewVSphereRepository(provider.GetConfig(), repositoryCache)
    if err != nil {
        return err
    }
	deserializer := resource.NewDeserializer(factory)
	providerLibrary.AddProvider(terraform.VSPHERE, provider)

	supplierLibrary.AddSupplier(NewVSphereVirtualMachineSupplier(provider, repository, deserializer))

	err = resourceSchemaRepository.Init(version, provider.Schema())
	if err != nil {
		return err
	}

	return nil
}
