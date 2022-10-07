// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"log"
	"sort"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/statusutils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// Compile-time assertion
	_ resources.HashableInputResource = new(MatchableHttpGateway)
)

func NewMatchableHttpGatewayHashableResource() resources.HashableResource {
	return new(MatchableHttpGateway)
}

func NewMatchableHttpGateway(namespace, name string) *MatchableHttpGateway {
	matchablehttpgateway := &MatchableHttpGateway{}
	matchablehttpgateway.SetMetadata(&core.Metadata{
		Name:      name,
		Namespace: namespace,
	})
	return matchablehttpgateway
}

func (r *MatchableHttpGateway) SetMetadata(meta *core.Metadata) {
	r.Metadata = meta
}

// Deprecated
func (r *MatchableHttpGateway) SetStatus(status *core.Status) {
	statusutils.SetSingleStatusInNamespacedStatuses(r, status)
}

// Deprecated
func (r *MatchableHttpGateway) GetStatus() *core.Status {
	if r != nil {
		return statusutils.GetSingleStatusInNamespacedStatuses(r)
	}
	return nil
}

func (r *MatchableHttpGateway) SetNamespacedStatuses(namespacedStatuses *core.NamespacedStatuses) {
	r.NamespacedStatuses = namespacedStatuses
}

func (r *MatchableHttpGateway) MustHash() uint64 {
	hashVal, err := r.Hash(nil)
	if err != nil {
		log.Panicf("error while hashing: (%s) this should never happen", err)
	}
	return hashVal
}

func (r *MatchableHttpGateway) GroupVersionKind() schema.GroupVersionKind {
	return MatchableHttpGatewayGVK
}

type MatchableHttpGatewayList []*MatchableHttpGateway

func (list MatchableHttpGatewayList) Find(namespace, name string) (*MatchableHttpGateway, error) {
	for _, matchableHttpGateway := range list {
		if matchableHttpGateway.GetMetadata().Name == name && matchableHttpGateway.GetMetadata().Namespace == namespace {
			return matchableHttpGateway, nil
		}
	}
	return nil, errors.Errorf("list did not find matchableHttpGateway %v.%v", namespace, name)
}

func (list MatchableHttpGatewayList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, matchableHttpGateway := range list {
		ress = append(ress, matchableHttpGateway)
	}
	return ress
}

func (list MatchableHttpGatewayList) AsInputResources() resources.InputResourceList {
	var ress resources.InputResourceList
	for _, matchableHttpGateway := range list {
		ress = append(ress, matchableHttpGateway)
	}
	return ress
}

func (list MatchableHttpGatewayList) Names() []string {
	var names []string
	for _, matchableHttpGateway := range list {
		names = append(names, matchableHttpGateway.GetMetadata().Name)
	}
	return names
}

func (list MatchableHttpGatewayList) NamespacesDotNames() []string {
	var names []string
	for _, matchableHttpGateway := range list {
		names = append(names, matchableHttpGateway.GetMetadata().Namespace+"."+matchableHttpGateway.GetMetadata().Name)
	}
	return names
}

func (list MatchableHttpGatewayList) Sort() MatchableHttpGatewayList {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].GetMetadata().Less(list[j].GetMetadata())
	})
	return list
}

func (list MatchableHttpGatewayList) Clone() MatchableHttpGatewayList {
	var matchableHttpGatewayList MatchableHttpGatewayList
	for _, matchableHttpGateway := range list {
		matchableHttpGatewayList = append(matchableHttpGatewayList, resources.Clone(matchableHttpGateway).(*MatchableHttpGateway))
	}
	return matchableHttpGatewayList
}

func (list MatchableHttpGatewayList) Each(f func(element *MatchableHttpGateway)) {
	for _, matchableHttpGateway := range list {
		f(matchableHttpGateway)
	}
}

func (list MatchableHttpGatewayList) EachResource(f func(element resources.Resource)) {
	for _, matchableHttpGateway := range list {
		f(matchableHttpGateway)
	}
}

func (list MatchableHttpGatewayList) AsInterfaces() []interface{} {
	var asInterfaces []interface{}
	list.Each(func(element *MatchableHttpGateway) {
		asInterfaces = append(asInterfaces, element)
	})
	return asInterfaces
}

// Kubernetes Adapter for MatchableHttpGateway

func (o *MatchableHttpGateway) GetObjectKind() schema.ObjectKind {
	t := MatchableHttpGatewayCrd.TypeMeta()
	return &t
}

func (o *MatchableHttpGateway) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*MatchableHttpGateway)
}

func (o *MatchableHttpGateway) DeepCopyInto(out *MatchableHttpGateway) {
	clone := resources.Clone(o).(*MatchableHttpGateway)
	*out = *clone
}

var (
	MatchableHttpGatewayCrd = crd.NewCrd(
		"httpgateways",
		MatchableHttpGatewayGVK.Group,
		MatchableHttpGatewayGVK.Version,
		MatchableHttpGatewayGVK.Kind,
		"hgw",
		false,
		&MatchableHttpGateway{})
)

var (
	MatchableHttpGatewayGVK = schema.GroupVersionKind{
		Version: "v1",
		Group:   "gateway.solo.io",
		Kind:    "MatchableHttpGateway",
	}
)
