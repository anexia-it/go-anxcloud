package v1

import (
	"go.anx.io/go-anxcloud/pkg/api/types"
)

func (obj *Location) DeepCopy() types.Object {
	// Initialize arrays

	out := &Location{
		// Primitives
		Identifier:  obj.Identifier,
		Code:        obj.Code,
		Name:        obj.Name,
		CountryCode: obj.CountryCode,
		CityCode:    obj.CityCode,

		// DeepCopyable

		// Arrays

	}

	*out.Latitude = *obj.Latitude
	*out.Longitude = *obj.Longitude

	return out
}

func (obj *Type) DeepCopy() *Type {
	// Initialize arrays

	out := &Type{
		// Primitives
		Identifier: obj.Identifier,
		Name:       obj.Name,

		// DeepCopyable

		// Arrays

	}

	return out
}

func (obj *Resource) DeepCopy() types.Object {
	// Initialize arrays
	copyOfTags := make([]string, 0, len(obj.Tags))
	for _, v := range obj.Tags {
		copyOfTags = append(copyOfTags, v)
	}

	out := &Resource{
		// Primitives
		Identifier: obj.Identifier,
		Name:       obj.Name,
		Type:       obj.Type,
		Attributes: obj.Attributes,

		// DeepCopyable

		// Arrays
		Tags: copyOfTags,
	}

	return out
}

func (obj *ResourceWithTag) DeepCopy() types.Object {
	// Initialize arrays

	out := &ResourceWithTag{
		// Primitives
		Identifier: obj.Identifier,
		Tag:        obj.Tag,

		// DeepCopyable

		// Arrays

	}

	return out
}
