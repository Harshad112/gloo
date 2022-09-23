// Code generated by protoc-gen-ext. DO NOT EDIT.
// source: github.com/solo-io/gloo/projects/gloo/api/v1/options/grpc_json/grpc_json.proto

package grpc_json

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"hash/fnv"

	safe_hasher "github.com/solo-io/protoc-gen-ext/pkg/hasher"
	"github.com/solo-io/protoc-gen-ext/pkg/hasher/hashstructure"
)

// ensure the imports are used
var (
	_ = errors.New("")
	_ = fmt.Print
	_ = binary.LittleEndian
	_ = new(hash.Hash64)
	_ = fnv.New64
	_ = hashstructure.Hash
	_ = new(safe_hasher.SafeHasher)
)

// Hash function
func (m *GrpcJsonTranscoder) Hash(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("grpc_json.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/grpc_json.GrpcJsonTranscoder")); err != nil {
		return 0, err
	}

	for _, v := range m.GetServices() {

		if _, err = hasher.Write([]byte(v)); err != nil {
			return 0, err
		}

	}

	if h, ok := interface{}(m.GetPrintOptions()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("PrintOptions")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetPrintOptions(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("PrintOptions")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	err = binary.Write(hasher, binary.LittleEndian, m.GetMatchIncomingRequestRoute())
	if err != nil {
		return 0, err
	}

	for _, v := range m.GetIgnoredQueryParameters() {

		if _, err = hasher.Write([]byte(v)); err != nil {
			return 0, err
		}

	}

	err = binary.Write(hasher, binary.LittleEndian, m.GetAutoMapping())
	if err != nil {
		return 0, err
	}

	err = binary.Write(hasher, binary.LittleEndian, m.GetIgnoreUnknownQueryParameters())
	if err != nil {
		return 0, err
	}

	err = binary.Write(hasher, binary.LittleEndian, m.GetConvertGrpcStatus())
	if err != nil {
		return 0, err
	}

	switch m.DescriptorSet.(type) {

	case *GrpcJsonTranscoder_ProtoDescriptor:

		if _, err = hasher.Write([]byte(m.GetProtoDescriptor())); err != nil {
			return 0, err
		}

	case *GrpcJsonTranscoder_ProtoDescriptorBin:

		if _, err = hasher.Write(m.GetProtoDescriptorBin()); err != nil {
			return 0, err
		}

	}

	return hasher.Sum64(), nil
}

// Hash function
func (m *GrpcJsonTranscoder_PrintOptions) Hash(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("grpc_json.options.gloo.solo.io.github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/grpc_json.GrpcJsonTranscoder_PrintOptions")); err != nil {
		return 0, err
	}

	err = binary.Write(hasher, binary.LittleEndian, m.GetAddWhitespace())
	if err != nil {
		return 0, err
	}

	err = binary.Write(hasher, binary.LittleEndian, m.GetAlwaysPrintPrimitiveFields())
	if err != nil {
		return 0, err
	}

	err = binary.Write(hasher, binary.LittleEndian, m.GetAlwaysPrintEnumsAsInts())
	if err != nil {
		return 0, err
	}

	err = binary.Write(hasher, binary.LittleEndian, m.GetPreserveProtoFieldNames())
	if err != nil {
		return 0, err
	}

	return hasher.Sum64(), nil
}
