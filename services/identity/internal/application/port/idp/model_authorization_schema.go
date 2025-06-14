// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.
//
// Keycloak Admin REST API
//
// This is a REST API reference for the Keycloak Admin REST API.
//
// API version: 1.0

package idp

import (
	"encoding/json"
)

// checks if the AuthorizationSchema type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AuthorizationSchema{}

// AuthorizationSchema struct for AuthorizationSchema
type AuthorizationSchema struct {
	ResourceTypes *map[string]ResourceType `json:"resourceTypes,omitempty"`
}

// NewAuthorizationSchema instantiates a new AuthorizationSchema object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuthorizationSchema() *AuthorizationSchema {
	this := AuthorizationSchema{}
	return &this
}

// NewAuthorizationSchemaWithDefaults instantiates a new AuthorizationSchema object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuthorizationSchemaWithDefaults() *AuthorizationSchema {
	this := AuthorizationSchema{}
	return &this
}

// GetResourceTypes returns the ResourceTypes field value if set, zero value otherwise.
func (o *AuthorizationSchema) GetResourceTypes() map[string]ResourceType {
	if o == nil || IsNil(o.ResourceTypes) {
		var ret map[string]ResourceType
		return ret
	}
	return *o.ResourceTypes
}

// GetResourceTypesOk returns a tuple with the ResourceTypes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthorizationSchema) GetResourceTypesOk() (*map[string]ResourceType, bool) {
	if o == nil || IsNil(o.ResourceTypes) {
		return nil, false
	}
	return o.ResourceTypes, true
}

// HasResourceTypes returns a boolean if a field has been set.
func (o *AuthorizationSchema) HasResourceTypes() bool {
	if o != nil && !IsNil(o.ResourceTypes) {
		return true
	}

	return false
}

// SetResourceTypes gets a reference to the given map[string]ResourceType and assigns it to the ResourceTypes field.
func (o *AuthorizationSchema) SetResourceTypes(v map[string]ResourceType) {
	o.ResourceTypes = &v
}

func (o AuthorizationSchema) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AuthorizationSchema) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ResourceTypes) {
		toSerialize["resourceTypes"] = o.ResourceTypes
	}
	return toSerialize, nil
}

type NullableAuthorizationSchema struct {
	value *AuthorizationSchema
	isSet bool
}

func (v NullableAuthorizationSchema) Get() *AuthorizationSchema {
	return v.value
}

func (v *NullableAuthorizationSchema) Set(val *AuthorizationSchema) {
	v.value = val
	v.isSet = true
}

func (v NullableAuthorizationSchema) IsSet() bool {
	return v.isSet
}

func (v *NullableAuthorizationSchema) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAuthorizationSchema(val *AuthorizationSchema) *NullableAuthorizationSchema {
	return &NullableAuthorizationSchema{value: val, isSet: true}
}

func (v NullableAuthorizationSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAuthorizationSchema) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
