package vsphere

import (
	"context"

	"github.com/cloudskiff/driftctl/pkg/remote/cache"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type VSphereRepository interface {
	ListVirtualMachines() ([]string, error)
}

type vsphereRepository struct {
	client *govmomi.Client
	ctx    context.Context
	config vsphereConfig
	cache  cache.Cache
}

func NewVSphereRepository(config vsphereConfig, c cache.Cache) (*vsphereRepository, error) {
	url, err := config.url()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	client, err := govmomi.NewClient(ctx, url, false)
	if err != nil {
		return nil, err
	}

	repo := &vsphereRepository{
		client: client,
		ctx:    context.Background(),
		config: config,
		cache:  c,
	}

	return repo, nil
}

func (r *vsphereRepository) ListVirtualMachines() ([]string, error) {
	ctx := context.TODO()
	finder := find.NewFinder(r.client.Client, true)
	dc, err := finder.DefaultDatacenter(ctx)
	if err != nil {
		return nil, err
	}
	finder.SetDatacenter(dc)
	vms, err := finder.VirtualMachineList(ctx, "*")
	if err != nil {
		return nil, err
	}
	collector := property.DefaultCollector(r.client.Client)
	var refs []types.ManagedObjectReference
	for _, vm := range vms {
		refs = append(refs, vm.Reference())
	}
	var vmt []mo.VirtualMachine
	err = collector.Retrieve(ctx, refs, []string{"config.uuid"}, &vmt)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, vm := range vmt {
		ids = append(ids, vm.Config.Uuid)
	}
	return ids, nil
}
