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
	_ resources.HashableInputResource = new(RouteOption)
)

func NewRouteOptionHashableResource() resources.HashableResource {
	return new(RouteOption)
}

func NewRouteOption(namespace, name string) *RouteOption {
	routeoption := &RouteOption{}
	routeoption.SetMetadata(&core.Metadata{
		Name:      name,
		Namespace: namespace,
	})
	return routeoption
}

func (r *RouteOption) SetMetadata(meta *core.Metadata) {
	r.Metadata = meta
}

// Deprecated
func (r *RouteOption) SetStatus(status *core.Status) {
	statusutils.SetSingleStatusInNamespacedStatuses(r, status)
}

// Deprecated
func (r *RouteOption) GetStatus() *core.Status {
	if r != nil {
		return statusutils.GetSingleStatusInNamespacedStatuses(r)
	}
	return nil
}

func (r *RouteOption) SetNamespacedStatuses(namespacedStatuses *core.NamespacedStatuses) {
	r.NamespacedStatuses = namespacedStatuses
}

func (r *RouteOption) MustHash() uint64 {
	hashVal, err := r.Hash(nil)
	if err != nil {
		log.Panicf("error while hashing: (%s) this should never happen", err)
	}
	return hashVal
}

func (r *RouteOption) GroupVersionKind() schema.GroupVersionKind {
	return RouteOptionGVK
}

type RouteOptionList []*RouteOption

func (list RouteOptionList) Find(namespace, name string) (*RouteOption, error) {
	for _, routeOption := range list {
		if routeOption.GetMetadata().Name == name && routeOption.GetMetadata().Namespace == namespace {
			return routeOption, nil
		}
	}
	return nil, errors.Errorf("list did not find routeOption %v.%v", namespace, name)
}

func (list RouteOptionList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, routeOption := range list {
		ress = append(ress, routeOption)
	}
	return ress
}

func (list RouteOptionList) AsInputResources() resources.InputResourceList {
	var ress resources.InputResourceList
	for _, routeOption := range list {
		ress = append(ress, routeOption)
	}
	return ress
}

func (list RouteOptionList) Names() []string {
	var names []string
	for _, routeOption := range list {
		names = append(names, routeOption.GetMetadata().Name)
	}
	return names
}

func (list RouteOptionList) NamespacesDotNames() []string {
	var names []string
	for _, routeOption := range list {
		names = append(names, routeOption.GetMetadata().Namespace+"."+routeOption.GetMetadata().Name)
	}
	return names
}

func (list RouteOptionList) Sort() RouteOptionList {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].GetMetadata().Less(list[j].GetMetadata())
	})
	return list
}

func (list RouteOptionList) Clone() RouteOptionList {
	var routeOptionList RouteOptionList
	for _, routeOption := range list {
		routeOptionList = append(routeOptionList, resources.Clone(routeOption).(*RouteOption))
	}
	return routeOptionList
}

func (list RouteOptionList) Each(f func(element *RouteOption)) {
	for _, routeOption := range list {
		f(routeOption)
	}
}

func (list RouteOptionList) EachResource(f func(element resources.Resource)) {
	for _, routeOption := range list {
		f(routeOption)
	}
}

func (list RouteOptionList) AsInterfaces() []interface{} {
	var asInterfaces []interface{}
	list.Each(func(element *RouteOption) {
		asInterfaces = append(asInterfaces, element)
	})
	return asInterfaces
}

// Kubernetes Adapter for RouteOption

func (o *RouteOption) GetObjectKind() schema.ObjectKind {
	t := RouteOptionCrd.TypeMeta()
	return &t
}

func (o *RouteOption) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*RouteOption)
}

func (o *RouteOption) DeepCopyInto(out *RouteOption) {
	clone := resources.Clone(o).(*RouteOption)
	*out = *clone
}

var (
	RouteOptionCrd = crd.NewCrd(
		"routeoptions",
		RouteOptionGVK.Group,
		RouteOptionGVK.Version,
		RouteOptionGVK.Kind,
		"rtopts",
		false,
		&RouteOption{})
)

var (
	RouteOptionGVK = schema.GroupVersionKind{
		Version: "v1",
		Group:   "gateway.solo.io",
		Kind:    "RouteOption",
	}
)
