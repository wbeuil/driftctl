package vsphere

import (
	remoteerror "github.com/cloudskiff/driftctl/pkg/remote/error"
	"github.com/cloudskiff/driftctl/pkg/resource"
	resourcevsphere "github.com/cloudskiff/driftctl/pkg/resource/vsphere"
	"github.com/cloudskiff/driftctl/pkg/terraform"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
)

type VSphereVirtualMachineSupplier struct {
	reader       terraform.ResourceReader
	deserializer *resource.Deserializer
	repository   VSphereRepository
	runner       *terraform.ParallelResourceReader
}

func NewVSphereVirtualMachineSupplier(provider *VSphereTerraformProvider, repository VSphereRepository, deserializer *resource.Deserializer) *VSphereVirtualMachineSupplier {
	return &VSphereVirtualMachineSupplier{
		provider,
		deserializer,
		repository,
		terraform.NewParallelResourceReader(provider.Runner().SubRunner()),
	}
}

func (s VSphereVirtualMachineSupplier) Resources() ([]resource.Resource, error) {
	resourceList, err := s.repository.ListVirtualMachines()
	if err != nil {
		return nil, remoteerror.NewResourceEnumerationError(err, resourcevsphere.VSphereVirtualMachineResourceType)
	}

	for _, id := range resourceList {
		id := id
		s.runner.Run(func() (cty.Value, error) {
			completeResource, err := s.reader.ReadResource(terraform.ReadResourceArgs{
				Ty: resourcevsphere.VSphereVirtualMachineResourceType,
				ID: id,
			})
			if err != nil {
				logrus.Warnf("Error reading %s[%s]: %+v", id, resourcevsphere.VSphereVirtualMachineResourceType, err)
				return cty.NilVal, err
			}
			return *completeResource, nil
		})
	}

	results, err := s.runner.Wait()
	if err != nil {
		return nil, err
	}

	return s.deserializer.Deserialize(resourcevsphere.VSphereVirtualMachineResourceType, results)
}
