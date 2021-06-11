package pkg

import (
	"context"

	"github.com/cloudskiff/driftctl/pkg/alerter"
	"github.com/cloudskiff/driftctl/pkg/parallel"
	"github.com/cloudskiff/driftctl/pkg/remote"
	"github.com/cloudskiff/driftctl/pkg/remote/common"
	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ScannerOptions struct {
	Deep bool
}

type Scanner struct {
	resourceSuppliers []resource.Supplier
	runner            *parallel.ParallelRunner
	remoteLibrary     *common.RemoteLibrary
	alerter           *alerter.Alerter
	options           ScannerOptions
}

func NewScanner(resourceSuppliers []resource.Supplier, remoteLibrary *common.RemoteLibrary, alerter *alerter.Alerter, options ScannerOptions) *Scanner {
	return &Scanner{
		resourceSuppliers: resourceSuppliers,
		runner:            parallel.NewParallelRunner(context.TODO(), 10),
		remoteLibrary:     remoteLibrary,
		alerter:           alerter,
		options:           options,
	}
}

func (s *Scanner) retrieveRunnerResults() ([]resource.Resource, error) {
	results := make([]resource.Resource, 0)
loop:
	for {
		select {
		case resources, ok := <-s.runner.Read():
			if !ok || resources == nil {
				break loop
			}

			for _, res := range resources.([]resource.Resource) {
				if res != nil {
					results = append(results, res)
				}
			}
		case <-s.runner.DoneChan():
			break loop
		}
	}
	return results, s.runner.Err()
}

func (s *Scanner) legacyScan() ([]resource.Resource, error) {
	for _, resourceProvider := range s.resourceSuppliers {
		supplier := resourceProvider
		s.runner.Run(func() (interface{}, error) {
			res, err := supplier.Resources()
			if err != nil {
				err := remote.HandleResourceEnumerationError(err, s.alerter)
				if err == nil {
					return []resource.Resource{}, nil
				}
				return nil, err
			}
			for _, resource := range res {
				logrus.WithFields(logrus.Fields{
					"id":   resource.TerraformId(),
					"type": resource.TerraformType(),
				}).Debug("[DEPRECATED] Found cloud resource")
			}
			return res, nil
		})
	}

	return s.retrieveRunnerResults()
}

func (s *Scanner) scan() ([]resource.Resource, error) {
	for _, enumerator := range s.remoteLibrary.Enumerators() {
		enumerator := enumerator
		s.runner.Run(func() (interface{}, error) {
			resources, err := enumerator.Enumerate()
			if err != nil {
				return nil, err
			}
			for _, resource := range resources {
				if resource == nil {
					continue
				}
				logrus.WithFields(logrus.Fields{
					"id":   resource.TerraformId(),
					"type": resource.TerraformType(),
				}).Debug("Found cloud resource")
			}
			return resources, nil
		})
	}

	enumerationResult, err := s.retrieveRunnerResults()
	if err != nil {
		return nil, err
	}

	s.runner = parallel.NewParallelRunner(context.TODO(), 10)

	for _, res := range enumerationResult {
		res := res
		s.runner.Run(func() (interface{}, error) {
			fetcher := s.remoteLibrary.GetDetailFetcher(resource.ResourceType(res.TerraformType()))
			if fetcher != nil {
				// If we are in deep mode, retrieve resource details
				if s.options.Deep {
					resourceWithDetails, err := fetcher.ReadDetails(res)
					if err != nil {
						return nil, err
					}
					return []resource.Resource{resourceWithDetails}, nil
				}
			}
			return []resource.Resource{res}, nil
		})
	}

	return s.retrieveRunnerResults()
}

func (s *Scanner) Resources() ([]resource.Resource, error) {

	resources, err := s.legacyScan()
	if err != nil {
		return nil, err
	}

	s.runner = parallel.NewParallelRunner(context.TODO(), 10)

	enumerationResult, err := s.scan()
	if err != nil {
		return nil, err
	}
	resources = append(resources, enumerationResult...)

	return resources, err
}

func (s *Scanner) Stop() {
	logrus.Debug("Stopping scanner")
	s.runner.Stop(errors.New("interrupted"))
}
