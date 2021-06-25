package vsphere

import (
	"context"

	"github.com/cloudskiff/driftctl/pkg/remote/cache"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
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
	ctx2 := context.TODO()
	finder := find.NewFinder(r.client.Client, true)
	dc, err := finder.DefaultDatacenter(ctx2)
	if err != nil {
		return nil, err
	}
	finder.SetDatacenter(dc)
	vms, err := finder.VirtualMachineList(ctx2, "*")
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, vm := range vms {
		ids = append(ids, vm.InventoryPath)
	}
	return ids, nil
}
