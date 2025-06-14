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

// checks if the ScopeRepresentation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ScopeRepresentation{}

// ScopeRepresentation struct for ScopeRepresentation
type ScopeRepresentation struct {
	Id          *string                  `json:"id,omitempty"`
	Name        *string                  `json:"name,omitempty"`
	IconUri     *string                  `json:"iconUri,omitempty"`
	Policies    []PolicyRepresentation   `json:"policies,omitempty"`
	Resources   []ResourceRepresentation `json:"resources,omitempty"`
	DisplayName *string                  `json:"displayName,omitempty"`
}

// NewScopeRepresentation instantiates a new ScopeRepresentation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewScopeRepresentation() *ScopeRepresentation {
	this := ScopeRepresentation{}
	return &this
}

// NewScopeRepresentationWithDefaults instantiates a new ScopeRepresentation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewScopeRepresentationWithDefaults() *ScopeRepresentation {
	this := ScopeRepresentation{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *ScopeRepresentation) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ScopeRepresentation) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ScopeRepresentation) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ScopeRepresentation) SetId(v string) {
	o.Id = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ScopeRepresentation) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ScopeRepresentation) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ScopeRepresentation) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ScopeRepresentation) SetName(v string) {
	o.Name = &v
}

// GetIconUri returns the IconUri field value if set, zero value otherwise.
func (o *ScopeRepresentation) GetIconUri() string {
	if o == nil || IsNil(o.IconUri) {
		var ret string
		return ret
	}
	return *o.IconUri
}

// GetIconUriOk returns a tuple with the IconUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ScopeRepresentation) GetIconUriOk() (*string, bool) {
	if o == nil || IsNil(o.IconUri) {
		return nil, false
	}
	return o.IconUri, true
}

// HasIconUri returns a boolean if a field has been set.
func (o *ScopeRepresentation) HasIconUri() bool {
	if o != nil && !IsNil(o.IconUri) {
		return true
	}

	return false
}

// SetIconUri gets a reference to the given string and assigns it to the IconUri field.
func (o *ScopeRepresentation) SetIconUri(v string) {
	o.IconUri = &v
}

// GetPolicies returns the Policies field value if set, zero value otherwise.
func (o *ScopeRepresentation) GetPolicies() []PolicyRepresentation {
	if o == nil || IsNil(o.Policies) {
		var ret []PolicyRepresentation
		return ret
	}
	return o.Policies
}

// GetPoliciesOk returns a tuple with the Policies field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ScopeRepresentation) GetPoliciesOk() ([]PolicyRepresentation, bool) {
	if o == nil || IsNil(o.Policies) {
		return nil, false
	}
	return o.Policies, true
}

// HasPolicies returns a boolean if a field has been set.
func (o *ScopeRepresentation) HasPolicies() bool {
	if o != nil && !IsNil(o.Policies) {
		return true
	}

	return false
}

// SetPolicies gets a reference to the given []PolicyRepresentation and assigns it to the Policies field.
func (o *ScopeRepresentation) SetPolicies(v []PolicyRepresentation) {
	o.Policies = v
}

// GetResources returns the Resources field value if set, zero value otherwise.
func (o *ScopeRepresentation) GetResources() []ResourceRepresentation {
	if o == nil || IsNil(o.Resources) {
		var ret []ResourceRepresentation
		return ret
	}
	return o.Resources
}

// GetResourcesOk returns a tuple with the Resources field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ScopeRepresentation) GetResourcesOk() ([]ResourceRepresentation, bool) {
	if o == nil || IsNil(o.Resources) {
		return nil, false
	}
	return o.Resources, true
}

// HasResources returns a boolean if a field has been set.
func (o *ScopeRepresentation) HasResources() bool {
	if o != nil && !IsNil(o.Resources) {
		return true
	}

	return false
}

// SetResources gets a reference to the given []ResourceRepresentation and assigns it to the Resources field.
func (o *ScopeRepresentation) SetResources(v []ResourceRepresentation) {
	o.Resources = v
}

// GetDisplayName returns the DisplayName field value if set, zero value otherwise.
func (o *ScopeRepresentation) GetDisplayName() string {
	if o == nil || IsNil(o.DisplayName) {
		var ret string
		return ret
	}
	return *o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ScopeRepresentation) GetDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.DisplayName) {
		return nil, false
	}
	return o.DisplayName, true
}

// HasDisplayName returns a boolean if a field has been set.
func (o *ScopeRepresentation) HasDisplayName() bool {
	if o != nil && !IsNil(o.DisplayName) {
		return true
	}

	return false
}

// SetDisplayName gets a reference to the given string and assigns it to the DisplayName field.
func (o *ScopeRepresentation) SetDisplayName(v string) {
	o.DisplayName = &v
}

func (o ScopeRepresentation) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ScopeRepresentation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.IconUri) {
		toSerialize["iconUri"] = o.IconUri
	}
	if !IsNil(o.Policies) {
		toSerialize["policies"] = o.Policies
	}
	if !IsNil(o.Resources) {
		toSerialize["resources"] = o.Resources
	}
	if !IsNil(o.DisplayName) {
		toSerialize["displayName"] = o.DisplayName
	}
	return toSerialize, nil
}

type NullableScopeRepresentation struct {
	value *ScopeRepresentation
	isSet bool
}

func (v NullableScopeRepresentation) Get() *ScopeRepresentation {
	return v.value
}

func (v *NullableScopeRepresentation) Set(val *ScopeRepresentation) {
	v.value = val
	v.isSet = true
}

func (v NullableScopeRepresentation) IsSet() bool {
	return v.isSet
}

func (v *NullableScopeRepresentation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableScopeRepresentation(val *ScopeRepresentation) *NullableScopeRepresentation {
	return &NullableScopeRepresentation{value: val, isSet: true}
}

func (v NullableScopeRepresentation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableScopeRepresentation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
