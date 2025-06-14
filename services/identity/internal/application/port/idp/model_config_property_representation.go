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

// checks if the ConfigPropertyRepresentation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPropertyRepresentation{}

// ConfigPropertyRepresentation struct for ConfigPropertyRepresentation
type ConfigPropertyRepresentation struct {
	Name         *string     `json:"name,omitempty"`
	Label        *string     `json:"label,omitempty"`
	HelpText     *string     `json:"helpText,omitempty"`
	Type         *string     `json:"type,omitempty"`
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	Options      []string    `json:"options,omitempty"`
	Secret       *bool       `json:"secret,omitempty"`
	Required     *bool       `json:"required,omitempty"`
	ReadOnly     *bool       `json:"readOnly,omitempty"`
}

// NewConfigPropertyRepresentation instantiates a new ConfigPropertyRepresentation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPropertyRepresentation() *ConfigPropertyRepresentation {
	this := ConfigPropertyRepresentation{}
	return &this
}

// NewConfigPropertyRepresentationWithDefaults instantiates a new ConfigPropertyRepresentation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPropertyRepresentationWithDefaults() *ConfigPropertyRepresentation {
	this := ConfigPropertyRepresentation{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ConfigPropertyRepresentation) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPropertyRepresentation) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ConfigPropertyRepresentation) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ConfigPropertyRepresentation) SetName(v string) {
	o.Name = &v
}

// GetLabel returns the Label field value if set, zero value otherwise.
func (o *ConfigPropertyRepresentation) GetLabel() string {
	if o == nil || IsNil(o.Label) {
		var ret string
		return ret
	}
	return *o.Label
}

// GetLabelOk returns a tuple with the Label field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPropertyRepresentation) GetLabelOk() (*string, bool) {
	if o == nil || IsNil(o.Label) {
		return nil, false
	}
	return o.Label, true
}

// HasLabel returns a boolean if a field has been set.
func (o *ConfigPropertyRepresentation) HasLabel() bool {
	if o != nil && !IsNil(o.Label) {
		return true
	}

	return false
}

// SetLabel gets a reference to the given string and assigns it to the Label field.
func (o *ConfigPropertyRepresentation) SetLabel(v string) {
	o.Label = &v
}

// GetHelpText returns the HelpText field value if set, zero value otherwise.
func (o *ConfigPropertyRepresentation) GetHelpText() string {
	if o == nil || IsNil(o.HelpText) {
		var ret string
		return ret
	}
	return *o.HelpText
}

// GetHelpTextOk returns a tuple with the HelpText field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPropertyRepresentation) GetHelpTextOk() (*string, bool) {
	if o == nil || IsNil(o.HelpText) {
		return nil, false
	}
	return o.HelpText, true
}

// HasHelpText returns a boolean if a field has been set.
func (o *ConfigPropertyRepresentation) HasHelpText() bool {
	if o != nil && !IsNil(o.HelpText) {
		return true
	}

	return false
}

// SetHelpText gets a reference to the given string and assigns it to the HelpText field.
func (o *ConfigPropertyRepresentation) SetHelpText(v string) {
	o.HelpText = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *ConfigPropertyRepresentation) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPropertyRepresentation) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *ConfigPropertyRepresentation) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *ConfigPropertyRepresentation) SetType(v string) {
	o.Type = &v
}

// GetDefaultValue returns the DefaultValue field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ConfigPropertyRepresentation) GetDefaultValue() interface{} {
	if o == nil {
		var ret interface{}
		return ret
	}
	return o.DefaultValue
}

// GetDefaultValueOk returns a tuple with the DefaultValue field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ConfigPropertyRepresentation) GetDefaultValueOk() (*interface{}, bool) {
	if o == nil || IsNil(o.DefaultValue) {
		return nil, false
	}
	return &o.DefaultValue, true
}

// HasDefaultValue returns a boolean if a field has been set.
func (o *ConfigPropertyRepresentation) HasDefaultValue() bool {
	if o != nil && !IsNil(o.DefaultValue) {
		return true
	}

	return false
}

// SetDefaultValue gets a reference to the given interface{} and assigns it to the DefaultValue field.
func (o *ConfigPropertyRepresentation) SetDefaultValue(v interface{}) {
	o.DefaultValue = v
}

// GetOptions returns the Options field value if set, zero value otherwise.
func (o *ConfigPropertyRepresentation) GetOptions() []string {
	if o == nil || IsNil(o.Options) {
		var ret []string
		return ret
	}
	return o.Options
}

// GetOptionsOk returns a tuple with the Options field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPropertyRepresentation) GetOptionsOk() ([]string, bool) {
	if o == nil || IsNil(o.Options) {
		return nil, false
	}
	return o.Options, true
}

// HasOptions returns a boolean if a field has been set.
func (o *ConfigPropertyRepresentation) HasOptions() bool {
	if o != nil && !IsNil(o.Options) {
		return true
	}

	return false
}

// SetOptions gets a reference to the given []string and assigns it to the Options field.
func (o *ConfigPropertyRepresentation) SetOptions(v []string) {
	o.Options = v
}

// GetSecret returns the Secret field value if set, zero value otherwise.
func (o *ConfigPropertyRepresentation) GetSecret() bool {
	if o == nil || IsNil(o.Secret) {
		var ret bool
		return ret
	}
	return *o.Secret
}

// GetSecretOk returns a tuple with the Secret field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPropertyRepresentation) GetSecretOk() (*bool, bool) {
	if o == nil || IsNil(o.Secret) {
		return nil, false
	}
	return o.Secret, true
}

// HasSecret returns a boolean if a field has been set.
func (o *ConfigPropertyRepresentation) HasSecret() bool {
	if o != nil && !IsNil(o.Secret) {
		return true
	}

	return false
}

// SetSecret gets a reference to the given bool and assigns it to the Secret field.
func (o *ConfigPropertyRepresentation) SetSecret(v bool) {
	o.Secret = &v
}

// GetRequired returns the Required field value if set, zero value otherwise.
func (o *ConfigPropertyRepresentation) GetRequired() bool {
	if o == nil || IsNil(o.Required) {
		var ret bool
		return ret
	}
	return *o.Required
}

// GetRequiredOk returns a tuple with the Required field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPropertyRepresentation) GetRequiredOk() (*bool, bool) {
	if o == nil || IsNil(o.Required) {
		return nil, false
	}
	return o.Required, true
}

// HasRequired returns a boolean if a field has been set.
func (o *ConfigPropertyRepresentation) HasRequired() bool {
	if o != nil && !IsNil(o.Required) {
		return true
	}

	return false
}

// SetRequired gets a reference to the given bool and assigns it to the Required field.
func (o *ConfigPropertyRepresentation) SetRequired(v bool) {
	o.Required = &v
}

// GetReadOnly returns the ReadOnly field value if set, zero value otherwise.
func (o *ConfigPropertyRepresentation) GetReadOnly() bool {
	if o == nil || IsNil(o.ReadOnly) {
		var ret bool
		return ret
	}
	return *o.ReadOnly
}

// GetReadOnlyOk returns a tuple with the ReadOnly field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPropertyRepresentation) GetReadOnlyOk() (*bool, bool) {
	if o == nil || IsNil(o.ReadOnly) {
		return nil, false
	}
	return o.ReadOnly, true
}

// HasReadOnly returns a boolean if a field has been set.
func (o *ConfigPropertyRepresentation) HasReadOnly() bool {
	if o != nil && !IsNil(o.ReadOnly) {
		return true
	}

	return false
}

// SetReadOnly gets a reference to the given bool and assigns it to the ReadOnly field.
func (o *ConfigPropertyRepresentation) SetReadOnly(v bool) {
	o.ReadOnly = &v
}

func (o ConfigPropertyRepresentation) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPropertyRepresentation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Label) {
		toSerialize["label"] = o.Label
	}
	if !IsNil(o.HelpText) {
		toSerialize["helpText"] = o.HelpText
	}
	if !IsNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	if o.DefaultValue != nil {
		toSerialize["defaultValue"] = o.DefaultValue
	}
	if !IsNil(o.Options) {
		toSerialize["options"] = o.Options
	}
	if !IsNil(o.Secret) {
		toSerialize["secret"] = o.Secret
	}
	if !IsNil(o.Required) {
		toSerialize["required"] = o.Required
	}
	if !IsNil(o.ReadOnly) {
		toSerialize["readOnly"] = o.ReadOnly
	}
	return toSerialize, nil
}

type NullableConfigPropertyRepresentation struct {
	value *ConfigPropertyRepresentation
	isSet bool
}

func (v NullableConfigPropertyRepresentation) Get() *ConfigPropertyRepresentation {
	return v.value
}

func (v *NullableConfigPropertyRepresentation) Set(val *ConfigPropertyRepresentation) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPropertyRepresentation) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPropertyRepresentation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPropertyRepresentation(val *ConfigPropertyRepresentation) *NullableConfigPropertyRepresentation {
	return &NullableConfigPropertyRepresentation{value: val, isSet: true}
}

func (v NullableConfigPropertyRepresentation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPropertyRepresentation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
