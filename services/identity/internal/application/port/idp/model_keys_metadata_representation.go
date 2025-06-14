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

// checks if the KeysMetadataRepresentation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &KeysMetadataRepresentation{}

// KeysMetadataRepresentation struct for KeysMetadataRepresentation
type KeysMetadataRepresentation struct {
	Active *map[string]string          `json:"active,omitempty"`
	Keys   []KeyMetadataRepresentation `json:"keys,omitempty"`
}

// NewKeysMetadataRepresentation instantiates a new KeysMetadataRepresentation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewKeysMetadataRepresentation() *KeysMetadataRepresentation {
	this := KeysMetadataRepresentation{}
	return &this
}

// NewKeysMetadataRepresentationWithDefaults instantiates a new KeysMetadataRepresentation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewKeysMetadataRepresentationWithDefaults() *KeysMetadataRepresentation {
	this := KeysMetadataRepresentation{}
	return &this
}

// GetActive returns the Active field value if set, zero value otherwise.
func (o *KeysMetadataRepresentation) GetActive() map[string]string {
	if o == nil || IsNil(o.Active) {
		var ret map[string]string
		return ret
	}
	return *o.Active
}

// GetActiveOk returns a tuple with the Active field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KeysMetadataRepresentation) GetActiveOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Active) {
		return nil, false
	}
	return o.Active, true
}

// HasActive returns a boolean if a field has been set.
func (o *KeysMetadataRepresentation) HasActive() bool {
	if o != nil && !IsNil(o.Active) {
		return true
	}

	return false
}

// SetActive gets a reference to the given map[string]string and assigns it to the Active field.
func (o *KeysMetadataRepresentation) SetActive(v map[string]string) {
	o.Active = &v
}

// GetKeys returns the Keys field value if set, zero value otherwise.
func (o *KeysMetadataRepresentation) GetKeys() []KeyMetadataRepresentation {
	if o == nil || IsNil(o.Keys) {
		var ret []KeyMetadataRepresentation
		return ret
	}
	return o.Keys
}

// GetKeysOk returns a tuple with the Keys field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KeysMetadataRepresentation) GetKeysOk() ([]KeyMetadataRepresentation, bool) {
	if o == nil || IsNil(o.Keys) {
		return nil, false
	}
	return o.Keys, true
}

// HasKeys returns a boolean if a field has been set.
func (o *KeysMetadataRepresentation) HasKeys() bool {
	if o != nil && !IsNil(o.Keys) {
		return true
	}

	return false
}

// SetKeys gets a reference to the given []KeyMetadataRepresentation and assigns it to the Keys field.
func (o *KeysMetadataRepresentation) SetKeys(v []KeyMetadataRepresentation) {
	o.Keys = v
}

func (o KeysMetadataRepresentation) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o KeysMetadataRepresentation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Active) {
		toSerialize["active"] = o.Active
	}
	if !IsNil(o.Keys) {
		toSerialize["keys"] = o.Keys
	}
	return toSerialize, nil
}

type NullableKeysMetadataRepresentation struct {
	value *KeysMetadataRepresentation
	isSet bool
}

func (v NullableKeysMetadataRepresentation) Get() *KeysMetadataRepresentation {
	return v.value
}

func (v *NullableKeysMetadataRepresentation) Set(val *KeysMetadataRepresentation) {
	v.value = val
	v.isSet = true
}

func (v NullableKeysMetadataRepresentation) IsSet() bool {
	return v.isSet
}

func (v *NullableKeysMetadataRepresentation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKeysMetadataRepresentation(val *KeysMetadataRepresentation) *NullableKeysMetadataRepresentation {
	return &NullableKeysMetadataRepresentation{value: val, isSet: true}
}

func (v NullableKeysMetadataRepresentation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKeysMetadataRepresentation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
