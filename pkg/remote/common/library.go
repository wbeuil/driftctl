package common

import (
	"github.com/cloudskiff/driftctl/pkg/resource"
)

type Enumerator interface {
	SupportedType() resource.ResourceType
	Enumerate() ([]resource.Resource, error)
}

type DetailFetcher interface {
	ReadDetails(resource.Resource) (resource.Resource, error)
}

type RemoteLibrary struct {
	enumerators    []Enumerator
	detailFetchers map[resource.ResourceType]DetailFetcher
}

func NewRemoteLibrary() *RemoteLibrary {
	return &RemoteLibrary{
		make([]Enumerator, 0),
		make(map[resource.ResourceType]DetailFetcher),
	}
}

func (r *RemoteLibrary) AddEnumerator(enumerator Enumerator) {
	r.enumerators = append(r.enumerators, enumerator)
}

func (r *RemoteLibrary) Enumerators() []Enumerator {
	return r.enumerators
}

func (r *RemoteLibrary) AddDetailFetcher(ty resource.ResourceType, detailFetcher DetailFetcher) {
	r.detailFetchers[ty] = detailFetcher
}

func (r *RemoteLibrary) GetDetailFetcher(ty resource.ResourceType) DetailFetcher {
	return r.detailFetchers[ty]
}
