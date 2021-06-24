package vsphere

import (
	"context"
	"fmt"

	"github.com/cloudskiff/driftctl/pkg/remote/cache"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
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
	// First version
	ctx := context.TODO()
	manager := view.NewManager(r.client.Client)
	view, err := manager.CreateContainerView(ctx, r.client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return nil, err
	}
	var vms []mo.VirtualMachine
	err = view.Retrieve(ctx, []string{"VirtualMachine"}, []string{"config.uuid"}, &vms)
	if err != nil {
		return nil, err
	}
	fmt.Println("First attempt with ContainerView")
	for _, vm := range vms {
		fmt.Println(vm)
	}

	// Second version
	ctx2 := context.TODO()
	finder := find.NewFinder(r.client.Client, true)
	dc, err := finder.DefaultDatacenter(ctx2)
	if err != nil {
		return nil, err
	}
	finder.SetDatacenter(dc)
	vms2, err := finder.VirtualMachineList(ctx2, "*")
	if err != nil {
		return nil, err
	}
	fmt.Println("Second attempt with Finder")
	for _, vm2 := range vms2 {
		fmt.Println(vm2)
	}

	return []string{}, nil
}
