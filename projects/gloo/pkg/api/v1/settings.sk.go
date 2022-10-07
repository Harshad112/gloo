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
	_ resources.HashableInputResource = new(Settings)
)

func NewSettingsHashableResource() resources.HashableResource {
	return new(Settings)
}

func NewSettings(namespace, name string) *Settings {
	settings := &Settings{}
	settings.SetMetadata(&core.Metadata{
		Name:      name,
		Namespace: namespace,
	})
	return settings
}

func (r *Settings) SetMetadata(meta *core.Metadata) {
	r.Metadata = meta
}

// Deprecated
func (r *Settings) SetStatus(status *core.Status) {
	statusutils.SetSingleStatusInNamespacedStatuses(r, status)
}

// Deprecated
func (r *Settings) GetStatus() *core.Status {
	if r != nil {
		return statusutils.GetSingleStatusInNamespacedStatuses(r)
	}
	return nil
}

func (r *Settings) SetNamespacedStatuses(namespacedStatuses *core.NamespacedStatuses) {
	r.NamespacedStatuses = namespacedStatuses
}

func (r *Settings) MustHash() uint64 {
	hashVal, err := r.Hash(nil)
	if err != nil {
		log.Panicf("error while hashing: (%s) this should never happen", err)
	}
	return hashVal
}

func (r *Settings) GroupVersionKind() schema.GroupVersionKind {
	return SettingsGVK
}

type SettingsList []*Settings

func (list SettingsList) Find(namespace, name string) (*Settings, error) {
	for _, settings := range list {
		if settings.GetMetadata().Name == name && settings.GetMetadata().Namespace == namespace {
			return settings, nil
		}
	}
	return nil, errors.Errorf("list did not find settings %v.%v", namespace, name)
}

func (list SettingsList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, settings := range list {
		ress = append(ress, settings)
	}
	return ress
}

func (list SettingsList) AsInputResources() resources.InputResourceList {
	var ress resources.InputResourceList
	for _, settings := range list {
		ress = append(ress, settings)
	}
	return ress
}

func (list SettingsList) Names() []string {
	var names []string
	for _, settings := range list {
		names = append(names, settings.GetMetadata().Name)
	}
	return names
}

func (list SettingsList) NamespacesDotNames() []string {
	var names []string
	for _, settings := range list {
		names = append(names, settings.GetMetadata().Namespace+"."+settings.GetMetadata().Name)
	}
	return names
}

func (list SettingsList) Sort() SettingsList {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].GetMetadata().Less(list[j].GetMetadata())
	})
	return list
}

func (list SettingsList) Clone() SettingsList {
	var settingsList SettingsList
	for _, settings := range list {
		settingsList = append(settingsList, resources.Clone(settings).(*Settings))
	}
	return settingsList
}

func (list SettingsList) Each(f func(element *Settings)) {
	for _, settings := range list {
		f(settings)
	}
}

func (list SettingsList) EachResource(f func(element resources.Resource)) {
	for _, settings := range list {
		f(settings)
	}
}

func (list SettingsList) AsInterfaces() []interface{} {
	var asInterfaces []interface{}
	list.Each(func(element *Settings) {
		asInterfaces = append(asInterfaces, element)
	})
	return asInterfaces
}

// Kubernetes Adapter for Settings

func (o *Settings) GetObjectKind() schema.ObjectKind {
	t := SettingsCrd.TypeMeta()
	return &t
}

func (o *Settings) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*Settings)
}

func (o *Settings) DeepCopyInto(out *Settings) {
	clone := resources.Clone(o).(*Settings)
	*out = *clone
}

var (
	SettingsCrd = crd.NewCrd(
		"settings",
		SettingsGVK.Group,
		SettingsGVK.Version,
		SettingsGVK.Kind,
		"st",
		false,
		&Settings{})
)

var (
	SettingsGVK = schema.GroupVersionKind{
		Version: "v1",
		Group:   "gloo.solo.io",
		Kind:    "Settings",
	}
)
