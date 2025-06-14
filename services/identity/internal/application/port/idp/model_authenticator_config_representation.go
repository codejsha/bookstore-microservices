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

// checks if the AuthenticatorConfigRepresentation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AuthenticatorConfigRepresentation{}

// AuthenticatorConfigRepresentation struct for AuthenticatorConfigRepresentation
type AuthenticatorConfigRepresentation struct {
	Id     *string            `json:"id,omitempty"`
	Alias  *string            `json:"alias,omitempty"`
	Config *map[string]string `json:"config,omitempty"`
}

// NewAuthenticatorConfigRepresentation instantiates a new AuthenticatorConfigRepresentation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuthenticatorConfigRepresentation() *AuthenticatorConfigRepresentation {
	this := AuthenticatorConfigRepresentation{}
	return &this
}

// NewAuthenticatorConfigRepresentationWithDefaults instantiates a new AuthenticatorConfigRepresentation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuthenticatorConfigRepresentationWithDefaults() *AuthenticatorConfigRepresentation {
	this := AuthenticatorConfigRepresentation{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *AuthenticatorConfigRepresentation) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthenticatorConfigRepresentation) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *AuthenticatorConfigRepresentation) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *AuthenticatorConfigRepresentation) SetId(v string) {
	o.Id = &v
}

// GetAlias returns the Alias field value if set, zero value otherwise.
func (o *AuthenticatorConfigRepresentation) GetAlias() string {
	if o == nil || IsNil(o.Alias) {
		var ret string
		return ret
	}
	return *o.Alias
}

// GetAliasOk returns a tuple with the Alias field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthenticatorConfigRepresentation) GetAliasOk() (*string, bool) {
	if o == nil || IsNil(o.Alias) {
		return nil, false
	}
	return o.Alias, true
}

// HasAlias returns a boolean if a field has been set.
func (o *AuthenticatorConfigRepresentation) HasAlias() bool {
	if o != nil && !IsNil(o.Alias) {
		return true
	}

	return false
}

// SetAlias gets a reference to the given string and assigns it to the Alias field.
func (o *AuthenticatorConfigRepresentation) SetAlias(v string) {
	o.Alias = &v
}

// GetConfig returns the Config field value if set, zero value otherwise.
func (o *AuthenticatorConfigRepresentation) GetConfig() map[string]string {
	if o == nil || IsNil(o.Config) {
		var ret map[string]string
		return ret
	}
	return *o.Config
}

// GetConfigOk returns a tuple with the Config field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthenticatorConfigRepresentation) GetConfigOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Config) {
		return nil, false
	}
	return o.Config, true
}

// HasConfig returns a boolean if a field has been set.
func (o *AuthenticatorConfigRepresentation) HasConfig() bool {
	if o != nil && !IsNil(o.Config) {
		return true
	}

	return false
}

// SetConfig gets a reference to the given map[string]string and assigns it to the Config field.
func (o *AuthenticatorConfigRepresentation) SetConfig(v map[string]string) {
	o.Config = &v
}

func (o AuthenticatorConfigRepresentation) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AuthenticatorConfigRepresentation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Alias) {
		toSerialize["alias"] = o.Alias
	}
	if !IsNil(o.Config) {
		toSerialize["config"] = o.Config
	}
	return toSerialize, nil
}

type NullableAuthenticatorConfigRepresentation struct {
	value *AuthenticatorConfigRepresentation
	isSet bool
}

func (v NullableAuthenticatorConfigRepresentation) Get() *AuthenticatorConfigRepresentation {
	return v.value
}

func (v *NullableAuthenticatorConfigRepresentation) Set(val *AuthenticatorConfigRepresentation) {
	v.value = val
	v.isSet = true
}

func (v NullableAuthenticatorConfigRepresentation) IsSet() bool {
	return v.isSet
}

func (v *NullableAuthenticatorConfigRepresentation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAuthenticatorConfigRepresentation(val *AuthenticatorConfigRepresentation) *NullableAuthenticatorConfigRepresentation {
	return &NullableAuthenticatorConfigRepresentation{value: val, isSet: true}
}

func (v NullableAuthenticatorConfigRepresentation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAuthenticatorConfigRepresentation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
