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

// checks if the SocialLinkRepresentation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &SocialLinkRepresentation{}

// SocialLinkRepresentation struct for SocialLinkRepresentation
type SocialLinkRepresentation struct {
	SocialProvider *string `json:"socialProvider,omitempty"`
	SocialUserId   *string `json:"socialUserId,omitempty"`
	SocialUsername *string `json:"socialUsername,omitempty"`
}

// NewSocialLinkRepresentation instantiates a new SocialLinkRepresentation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSocialLinkRepresentation() *SocialLinkRepresentation {
	this := SocialLinkRepresentation{}
	return &this
}

// NewSocialLinkRepresentationWithDefaults instantiates a new SocialLinkRepresentation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSocialLinkRepresentationWithDefaults() *SocialLinkRepresentation {
	this := SocialLinkRepresentation{}
	return &this
}

// GetSocialProvider returns the SocialProvider field value if set, zero value otherwise.
func (o *SocialLinkRepresentation) GetSocialProvider() string {
	if o == nil || IsNil(o.SocialProvider) {
		var ret string
		return ret
	}
	return *o.SocialProvider
}

// GetSocialProviderOk returns a tuple with the SocialProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SocialLinkRepresentation) GetSocialProviderOk() (*string, bool) {
	if o == nil || IsNil(o.SocialProvider) {
		return nil, false
	}
	return o.SocialProvider, true
}

// HasSocialProvider returns a boolean if a field has been set.
func (o *SocialLinkRepresentation) HasSocialProvider() bool {
	if o != nil && !IsNil(o.SocialProvider) {
		return true
	}

	return false
}

// SetSocialProvider gets a reference to the given string and assigns it to the SocialProvider field.
func (o *SocialLinkRepresentation) SetSocialProvider(v string) {
	o.SocialProvider = &v
}

// GetSocialUserId returns the SocialUserId field value if set, zero value otherwise.
func (o *SocialLinkRepresentation) GetSocialUserId() string {
	if o == nil || IsNil(o.SocialUserId) {
		var ret string
		return ret
	}
	return *o.SocialUserId
}

// GetSocialUserIdOk returns a tuple with the SocialUserId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SocialLinkRepresentation) GetSocialUserIdOk() (*string, bool) {
	if o == nil || IsNil(o.SocialUserId) {
		return nil, false
	}
	return o.SocialUserId, true
}

// HasSocialUserId returns a boolean if a field has been set.
func (o *SocialLinkRepresentation) HasSocialUserId() bool {
	if o != nil && !IsNil(o.SocialUserId) {
		return true
	}

	return false
}

// SetSocialUserId gets a reference to the given string and assigns it to the SocialUserId field.
func (o *SocialLinkRepresentation) SetSocialUserId(v string) {
	o.SocialUserId = &v
}

// GetSocialUsername returns the SocialUsername field value if set, zero value otherwise.
func (o *SocialLinkRepresentation) GetSocialUsername() string {
	if o == nil || IsNil(o.SocialUsername) {
		var ret string
		return ret
	}
	return *o.SocialUsername
}

// GetSocialUsernameOk returns a tuple with the SocialUsername field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SocialLinkRepresentation) GetSocialUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.SocialUsername) {
		return nil, false
	}
	return o.SocialUsername, true
}

// HasSocialUsername returns a boolean if a field has been set.
func (o *SocialLinkRepresentation) HasSocialUsername() bool {
	if o != nil && !IsNil(o.SocialUsername) {
		return true
	}

	return false
}

// SetSocialUsername gets a reference to the given string and assigns it to the SocialUsername field.
func (o *SocialLinkRepresentation) SetSocialUsername(v string) {
	o.SocialUsername = &v
}

func (o SocialLinkRepresentation) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o SocialLinkRepresentation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.SocialProvider) {
		toSerialize["socialProvider"] = o.SocialProvider
	}
	if !IsNil(o.SocialUserId) {
		toSerialize["socialUserId"] = o.SocialUserId
	}
	if !IsNil(o.SocialUsername) {
		toSerialize["socialUsername"] = o.SocialUsername
	}
	return toSerialize, nil
}

type NullableSocialLinkRepresentation struct {
	value *SocialLinkRepresentation
	isSet bool
}

func (v NullableSocialLinkRepresentation) Get() *SocialLinkRepresentation {
	return v.value
}

func (v *NullableSocialLinkRepresentation) Set(val *SocialLinkRepresentation) {
	v.value = val
	v.isSet = true
}

func (v NullableSocialLinkRepresentation) IsSet() bool {
	return v.isSet
}

func (v *NullableSocialLinkRepresentation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSocialLinkRepresentation(val *SocialLinkRepresentation) *NullableSocialLinkRepresentation {
	return &NullableSocialLinkRepresentation{value: val, isSet: true}
}

func (v NullableSocialLinkRepresentation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSocialLinkRepresentation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
