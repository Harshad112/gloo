// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"fmt"
	"hash"
	"hash/fnv"
	"log"

	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/hashutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type SetupSnapshot struct {
	Settings SettingsList
}

func (s SetupSnapshot) Clone() SetupSnapshot {
	return SetupSnapshot{
		Settings: s.Settings.Clone(),
	}
}

func (s SetupSnapshot) Hash(hasher hash.Hash64) (uint64, error) {
	if hasher == nil {
		hasher = fnv.New64()
	}
	if _, err := s.hashSettings(hasher); err != nil {
		return 0, err
	}
	return hasher.Sum64(), nil
}

func (s SetupSnapshot) hashSettings(hasher hash.Hash64) (uint64, error) {
	return hashutils.HashAllSafe(hasher, s.Settings.AsInterfaces()...)
}

func (s SetupSnapshot) HashFields() []zap.Field {
	var fields []zap.Field
	hasher := fnv.New64()
	SettingsHash, err := s.hashSettings(hasher)
	if err != nil {
		log.Println(eris.Wrapf(err, "error hashing, this should never happen"))
	}
	fields = append(fields, zap.Uint64("settings", SettingsHash))
	snapshotHash, err := s.Hash(hasher)
	if err != nil {
		log.Println(eris.Wrapf(err, "error hashing, this should never happen"))
	}
	return append(fields, zap.Uint64("snapshotHash", snapshotHash))
}

func (s *SetupSnapshot) GetResourcesList(resource resources.Resource) (resources.ResourceList, error) {
	switch resource.(type) {
	case *Settings:
		return s.Settings.AsResources(), nil
	default:
		return resources.ResourceList{}, eris.New("did not contain the input resource type returning empty list")
	}
}

func (s *SetupSnapshot) AddToResourceList(resource resources.Resource) error {
	switch typed := resource.(type) {
	case *Settings:
		s.Settings = append(s.Settings, typed)
		s.Settings.Sort()
		return nil
	default:
		return eris.New("did not add the input resource type because it does not exist")
	}
}

func (s *SetupSnapshot) ReplaceResource(i int, resource resources.Resource) error {
	switch typed := resource.(type) {
	case *Settings:
		s.Settings[i] = typed
	default:
		return eris.Wrapf(eris.New("did not contain the input resource type"), "did not replace the resource at index %d", i)
	}
	return nil
}

func (s *SetupSnapshot) GetInputResourcesList(resource resources.InputResource) (resources.InputResourceList, error) {
	switch resource.(type) {
	case *Settings:
		return s.Settings.AsInputResources(), nil
	default:
		return resources.InputResourceList{}, eris.New("did not contain the input resource type returning empty list")
	}
}

func (s *SetupSnapshot) AddToInputResourceList(resource resources.InputResource) error {
	switch typed := resource.(type) {
	case *Settings:
		s.Settings = append(s.Settings, typed)
		s.Settings.Sort()
		return nil
	default:
		return eris.New("did not add the input resource type because it does not exist")
	}
}

func (s *SetupSnapshot) ReplaceInputResource(i int, resource resources.InputResource) error {
	switch typed := resource.(type) {
	case *Settings:
		s.Settings[i] = typed
	default:
		return eris.Wrapf(eris.New("did not contain the input resource type"), "did not replace the resource at index %d", i)
	}
	return nil
}

type SetupSnapshotStringer struct {
	Version  uint64
	Settings []string
}

func (ss SetupSnapshotStringer) String() string {
	s := fmt.Sprintf("SetupSnapshot %v\n", ss.Version)

	s += fmt.Sprintf("  Settings %v\n", len(ss.Settings))
	for _, name := range ss.Settings {
		s += fmt.Sprintf("    %v\n", name)
	}

	return s
}

func (s SetupSnapshot) Stringer() SetupSnapshotStringer {
	snapshotHash, err := s.Hash(nil)
	if err != nil {
		log.Println(eris.Wrapf(err, "error hashing, this should never happen"))
	}
	return SetupSnapshotStringer{
		Version:  snapshotHash,
		Settings: s.Settings.NamespacesDotNames(),
	}
}

var SetupGvkToHashableResource = map[schema.GroupVersionKind]func() resources.HashableResource{
	SettingsGVK: NewSettingsHashableResource,
}
