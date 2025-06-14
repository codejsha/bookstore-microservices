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

// checks if the FederatedIdentityRepresentation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &FederatedIdentityRepresentation{}

// FederatedIdentityRepresentation struct for FederatedIdentityRepresentation
type FederatedIdentityRepresentation struct {
	IdentityProvider *string `json:"identityProvider,omitempty"`
	UserId           *string `json:"userId,omitempty"`
	UserName         *string `json:"userName,omitempty"`
}

// NewFederatedIdentityRepresentation instantiates a new FederatedIdentityRepresentation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFederatedIdentityRepresentation() *FederatedIdentityRepresentation {
	this := FederatedIdentityRepresentation{}
	return &this
}

// NewFederatedIdentityRepresentationWithDefaults instantiates a new FederatedIdentityRepresentation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFederatedIdentityRepresentationWithDefaults() *FederatedIdentityRepresentation {
	this := FederatedIdentityRepresentation{}
	return &this
}

// GetIdentityProvider returns the IdentityProvider field value if set, zero value otherwise.
func (o *FederatedIdentityRepresentation) GetIdentityProvider() string {
	if o == nil || IsNil(o.IdentityProvider) {
		var ret string
		return ret
	}
	return *o.IdentityProvider
}

// GetIdentityProviderOk returns a tuple with the IdentityProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederatedIdentityRepresentation) GetIdentityProviderOk() (*string, bool) {
	if o == nil || IsNil(o.IdentityProvider) {
		return nil, false
	}
	return o.IdentityProvider, true
}

// HasIdentityProvider returns a boolean if a field has been set.
func (o *FederatedIdentityRepresentation) HasIdentityProvider() bool {
	if o != nil && !IsNil(o.IdentityProvider) {
		return true
	}

	return false
}

// SetIdentityProvider gets a reference to the given string and assigns it to the IdentityProvider field.
func (o *FederatedIdentityRepresentation) SetIdentityProvider(v string) {
	o.IdentityProvider = &v
}

// GetUserId returns the UserId field value if set, zero value otherwise.
func (o *FederatedIdentityRepresentation) GetUserId() string {
	if o == nil || IsNil(o.UserId) {
		var ret string
		return ret
	}
	return *o.UserId
}

// GetUserIdOk returns a tuple with the UserId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederatedIdentityRepresentation) GetUserIdOk() (*string, bool) {
	if o == nil || IsNil(o.UserId) {
		return nil, false
	}
	return o.UserId, true
}

// HasUserId returns a boolean if a field has been set.
func (o *FederatedIdentityRepresentation) HasUserId() bool {
	if o != nil && !IsNil(o.UserId) {
		return true
	}

	return false
}

// SetUserId gets a reference to the given string and assigns it to the UserId field.
func (o *FederatedIdentityRepresentation) SetUserId(v string) {
	o.UserId = &v
}

// GetUserName returns the UserName field value if set, zero value otherwise.
func (o *FederatedIdentityRepresentation) GetUserName() string {
	if o == nil || IsNil(o.UserName) {
		var ret string
		return ret
	}
	return *o.UserName
}

// GetUserNameOk returns a tuple with the UserName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederatedIdentityRepresentation) GetUserNameOk() (*string, bool) {
	if o == nil || IsNil(o.UserName) {
		return nil, false
	}
	return o.UserName, true
}

// HasUserName returns a boolean if a field has been set.
func (o *FederatedIdentityRepresentation) HasUserName() bool {
	if o != nil && !IsNil(o.UserName) {
		return true
	}

	return false
}

// SetUserName gets a reference to the given string and assigns it to the UserName field.
func (o *FederatedIdentityRepresentation) SetUserName(v string) {
	o.UserName = &v
}

func (o FederatedIdentityRepresentation) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o FederatedIdentityRepresentation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.IdentityProvider) {
		toSerialize["identityProvider"] = o.IdentityProvider
	}
	if !IsNil(o.UserId) {
		toSerialize["userId"] = o.UserId
	}
	if !IsNil(o.UserName) {
		toSerialize["userName"] = o.UserName
	}
	return toSerialize, nil
}

type NullableFederatedIdentityRepresentation struct {
	value *FederatedIdentityRepresentation
	isSet bool
}

func (v NullableFederatedIdentityRepresentation) Get() *FederatedIdentityRepresentation {
	return v.value
}

func (v *NullableFederatedIdentityRepresentation) Set(val *FederatedIdentityRepresentation) {
	v.value = val
	v.isSet = true
}

func (v NullableFederatedIdentityRepresentation) IsSet() bool {
	return v.isSet
}

func (v *NullableFederatedIdentityRepresentation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFederatedIdentityRepresentation(val *FederatedIdentityRepresentation) *NullableFederatedIdentityRepresentation {
	return &NullableFederatedIdentityRepresentation{value: val, isSet: true}
}

func (v NullableFederatedIdentityRepresentation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFederatedIdentityRepresentation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
