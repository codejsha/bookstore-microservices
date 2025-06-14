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

// checks if the RequiredActionProviderRepresentation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &RequiredActionProviderRepresentation{}

// RequiredActionProviderRepresentation struct for RequiredActionProviderRepresentation
type RequiredActionProviderRepresentation struct {
	Alias         *string            `json:"alias,omitempty"`
	Name          *string            `json:"name,omitempty"`
	ProviderId    *string            `json:"providerId,omitempty"`
	Enabled       *bool              `json:"enabled,omitempty"`
	DefaultAction *bool              `json:"defaultAction,omitempty"`
	Priority      *int32             `json:"priority,omitempty"`
	Config        *map[string]string `json:"config,omitempty"`
}

// NewRequiredActionProviderRepresentation instantiates a new RequiredActionProviderRepresentation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRequiredActionProviderRepresentation() *RequiredActionProviderRepresentation {
	this := RequiredActionProviderRepresentation{}
	return &this
}

// NewRequiredActionProviderRepresentationWithDefaults instantiates a new RequiredActionProviderRepresentation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRequiredActionProviderRepresentationWithDefaults() *RequiredActionProviderRepresentation {
	this := RequiredActionProviderRepresentation{}
	return &this
}

// GetAlias returns the Alias field value if set, zero value otherwise.
func (o *RequiredActionProviderRepresentation) GetAlias() string {
	if o == nil || IsNil(o.Alias) {
		var ret string
		return ret
	}
	return *o.Alias
}

// GetAliasOk returns a tuple with the Alias field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequiredActionProviderRepresentation) GetAliasOk() (*string, bool) {
	if o == nil || IsNil(o.Alias) {
		return nil, false
	}
	return o.Alias, true
}

// HasAlias returns a boolean if a field has been set.
func (o *RequiredActionProviderRepresentation) HasAlias() bool {
	if o != nil && !IsNil(o.Alias) {
		return true
	}

	return false
}

// SetAlias gets a reference to the given string and assigns it to the Alias field.
func (o *RequiredActionProviderRepresentation) SetAlias(v string) {
	o.Alias = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *RequiredActionProviderRepresentation) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequiredActionProviderRepresentation) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *RequiredActionProviderRepresentation) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *RequiredActionProviderRepresentation) SetName(v string) {
	o.Name = &v
}

// GetProviderId returns the ProviderId field value if set, zero value otherwise.
func (o *RequiredActionProviderRepresentation) GetProviderId() string {
	if o == nil || IsNil(o.ProviderId) {
		var ret string
		return ret
	}
	return *o.ProviderId
}

// GetProviderIdOk returns a tuple with the ProviderId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequiredActionProviderRepresentation) GetProviderIdOk() (*string, bool) {
	if o == nil || IsNil(o.ProviderId) {
		return nil, false
	}
	return o.ProviderId, true
}

// HasProviderId returns a boolean if a field has been set.
func (o *RequiredActionProviderRepresentation) HasProviderId() bool {
	if o != nil && !IsNil(o.ProviderId) {
		return true
	}

	return false
}

// SetProviderId gets a reference to the given string and assigns it to the ProviderId field.
func (o *RequiredActionProviderRepresentation) SetProviderId(v string) {
	o.ProviderId = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *RequiredActionProviderRepresentation) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled) {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequiredActionProviderRepresentation) GetEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.Enabled) {
		return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *RequiredActionProviderRepresentation) HasEnabled() bool {
	if o != nil && !IsNil(o.Enabled) {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *RequiredActionProviderRepresentation) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetDefaultAction returns the DefaultAction field value if set, zero value otherwise.
func (o *RequiredActionProviderRepresentation) GetDefaultAction() bool {
	if o == nil || IsNil(o.DefaultAction) {
		var ret bool
		return ret
	}
	return *o.DefaultAction
}

// GetDefaultActionOk returns a tuple with the DefaultAction field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequiredActionProviderRepresentation) GetDefaultActionOk() (*bool, bool) {
	if o == nil || IsNil(o.DefaultAction) {
		return nil, false
	}
	return o.DefaultAction, true
}

// HasDefaultAction returns a boolean if a field has been set.
func (o *RequiredActionProviderRepresentation) HasDefaultAction() bool {
	if o != nil && !IsNil(o.DefaultAction) {
		return true
	}

	return false
}

// SetDefaultAction gets a reference to the given bool and assigns it to the DefaultAction field.
func (o *RequiredActionProviderRepresentation) SetDefaultAction(v bool) {
	o.DefaultAction = &v
}

// GetPriority returns the Priority field value if set, zero value otherwise.
func (o *RequiredActionProviderRepresentation) GetPriority() int32 {
	if o == nil || IsNil(o.Priority) {
		var ret int32
		return ret
	}
	return *o.Priority
}

// GetPriorityOk returns a tuple with the Priority field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequiredActionProviderRepresentation) GetPriorityOk() (*int32, bool) {
	if o == nil || IsNil(o.Priority) {
		return nil, false
	}
	return o.Priority, true
}

// HasPriority returns a boolean if a field has been set.
func (o *RequiredActionProviderRepresentation) HasPriority() bool {
	if o != nil && !IsNil(o.Priority) {
		return true
	}

	return false
}

// SetPriority gets a reference to the given int32 and assigns it to the Priority field.
func (o *RequiredActionProviderRepresentation) SetPriority(v int32) {
	o.Priority = &v
}

// GetConfig returns the Config field value if set, zero value otherwise.
func (o *RequiredActionProviderRepresentation) GetConfig() map[string]string {
	if o == nil || IsNil(o.Config) {
		var ret map[string]string
		return ret
	}
	return *o.Config
}

// GetConfigOk returns a tuple with the Config field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequiredActionProviderRepresentation) GetConfigOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Config) {
		return nil, false
	}
	return o.Config, true
}

// HasConfig returns a boolean if a field has been set.
func (o *RequiredActionProviderRepresentation) HasConfig() bool {
	if o != nil && !IsNil(o.Config) {
		return true
	}

	return false
}

// SetConfig gets a reference to the given map[string]string and assigns it to the Config field.
func (o *RequiredActionProviderRepresentation) SetConfig(v map[string]string) {
	o.Config = &v
}

func (o RequiredActionProviderRepresentation) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o RequiredActionProviderRepresentation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Alias) {
		toSerialize["alias"] = o.Alias
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.ProviderId) {
		toSerialize["providerId"] = o.ProviderId
	}
	if !IsNil(o.Enabled) {
		toSerialize["enabled"] = o.Enabled
	}
	if !IsNil(o.DefaultAction) {
		toSerialize["defaultAction"] = o.DefaultAction
	}
	if !IsNil(o.Priority) {
		toSerialize["priority"] = o.Priority
	}
	if !IsNil(o.Config) {
		toSerialize["config"] = o.Config
	}
	return toSerialize, nil
}

type NullableRequiredActionProviderRepresentation struct {
	value *RequiredActionProviderRepresentation
	isSet bool
}

func (v NullableRequiredActionProviderRepresentation) Get() *RequiredActionProviderRepresentation {
	return v.value
}

func (v *NullableRequiredActionProviderRepresentation) Set(val *RequiredActionProviderRepresentation) {
	v.value = val
	v.isSet = true
}

func (v NullableRequiredActionProviderRepresentation) IsSet() bool {
	return v.isSet
}

func (v *NullableRequiredActionProviderRepresentation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRequiredActionProviderRepresentation(val *RequiredActionProviderRepresentation) *NullableRequiredActionProviderRepresentation {
	return &NullableRequiredActionProviderRepresentation{value: val, isSet: true}
}

func (v NullableRequiredActionProviderRepresentation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRequiredActionProviderRepresentation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
