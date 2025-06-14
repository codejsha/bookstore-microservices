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

// checks if the ClaimRepresentation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ClaimRepresentation{}

// ClaimRepresentation struct for ClaimRepresentation
type ClaimRepresentation struct {
	Name     *bool `json:"name,omitempty"`
	Username *bool `json:"username,omitempty"`
	Profile  *bool `json:"profile,omitempty"`
	Picture  *bool `json:"picture,omitempty"`
	Website  *bool `json:"website,omitempty"`
	Email    *bool `json:"email,omitempty"`
	Gender   *bool `json:"gender,omitempty"`
	Locale   *bool `json:"locale,omitempty"`
	Address  *bool `json:"address,omitempty"`
	Phone    *bool `json:"phone,omitempty"`
}

// NewClaimRepresentation instantiates a new ClaimRepresentation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClaimRepresentation() *ClaimRepresentation {
	this := ClaimRepresentation{}
	return &this
}

// NewClaimRepresentationWithDefaults instantiates a new ClaimRepresentation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClaimRepresentationWithDefaults() *ClaimRepresentation {
	this := ClaimRepresentation{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ClaimRepresentation) GetName() bool {
	if o == nil || IsNil(o.Name) {
		var ret bool
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClaimRepresentation) GetNameOk() (*bool, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ClaimRepresentation) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given bool and assigns it to the Name field.
func (o *ClaimRepresentation) SetName(v bool) {
	o.Name = &v
}

// GetUsername returns the Username field value if set, zero value otherwise.
func (o *ClaimRepresentation) GetUsername() bool {
	if o == nil || IsNil(o.Username) {
		var ret bool
		return ret
	}
	return *o.Username
}

// GetUsernameOk returns a tuple with the Username field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClaimRepresentation) GetUsernameOk() (*bool, bool) {
	if o == nil || IsNil(o.Username) {
		return nil, false
	}
	return o.Username, true
}

// HasUsername returns a boolean if a field has been set.
func (o *ClaimRepresentation) HasUsername() bool {
	if o != nil && !IsNil(o.Username) {
		return true
	}

	return false
}

// SetUsername gets a reference to the given bool and assigns it to the Username field.
func (o *ClaimRepresentation) SetUsername(v bool) {
	o.Username = &v
}

// GetProfile returns the Profile field value if set, zero value otherwise.
func (o *ClaimRepresentation) GetProfile() bool {
	if o == nil || IsNil(o.Profile) {
		var ret bool
		return ret
	}
	return *o.Profile
}

// GetProfileOk returns a tuple with the Profile field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClaimRepresentation) GetProfileOk() (*bool, bool) {
	if o == nil || IsNil(o.Profile) {
		return nil, false
	}
	return o.Profile, true
}

// HasProfile returns a boolean if a field has been set.
func (o *ClaimRepresentation) HasProfile() bool {
	if o != nil && !IsNil(o.Profile) {
		return true
	}

	return false
}

// SetProfile gets a reference to the given bool and assigns it to the Profile field.
func (o *ClaimRepresentation) SetProfile(v bool) {
	o.Profile = &v
}

// GetPicture returns the Picture field value if set, zero value otherwise.
func (o *ClaimRepresentation) GetPicture() bool {
	if o == nil || IsNil(o.Picture) {
		var ret bool
		return ret
	}
	return *o.Picture
}

// GetPictureOk returns a tuple with the Picture field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClaimRepresentation) GetPictureOk() (*bool, bool) {
	if o == nil || IsNil(o.Picture) {
		return nil, false
	}
	return o.Picture, true
}

// HasPicture returns a boolean if a field has been set.
func (o *ClaimRepresentation) HasPicture() bool {
	if o != nil && !IsNil(o.Picture) {
		return true
	}

	return false
}

// SetPicture gets a reference to the given bool and assigns it to the Picture field.
func (o *ClaimRepresentation) SetPicture(v bool) {
	o.Picture = &v
}

// GetWebsite returns the Website field value if set, zero value otherwise.
func (o *ClaimRepresentation) GetWebsite() bool {
	if o == nil || IsNil(o.Website) {
		var ret bool
		return ret
	}
	return *o.Website
}

// GetWebsiteOk returns a tuple with the Website field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClaimRepresentation) GetWebsiteOk() (*bool, bool) {
	if o == nil || IsNil(o.Website) {
		return nil, false
	}
	return o.Website, true
}

// HasWebsite returns a boolean if a field has been set.
func (o *ClaimRepresentation) HasWebsite() bool {
	if o != nil && !IsNil(o.Website) {
		return true
	}

	return false
}

// SetWebsite gets a reference to the given bool and assigns it to the Website field.
func (o *ClaimRepresentation) SetWebsite(v bool) {
	o.Website = &v
}

// GetEmail returns the Email field value if set, zero value otherwise.
func (o *ClaimRepresentation) GetEmail() bool {
	if o == nil || IsNil(o.Email) {
		var ret bool
		return ret
	}
	return *o.Email
}

// GetEmailOk returns a tuple with the Email field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClaimRepresentation) GetEmailOk() (*bool, bool) {
	if o == nil || IsNil(o.Email) {
		return nil, false
	}
	return o.Email, true
}

// HasEmail returns a boolean if a field has been set.
func (o *ClaimRepresentation) HasEmail() bool {
	if o != nil && !IsNil(o.Email) {
		return true
	}

	return false
}

// SetEmail gets a reference to the given bool and assigns it to the Email field.
func (o *ClaimRepresentation) SetEmail(v bool) {
	o.Email = &v
}

// GetGender returns the Gender field value if set, zero value otherwise.
func (o *ClaimRepresentation) GetGender() bool {
	if o == nil || IsNil(o.Gender) {
		var ret bool
		return ret
	}
	return *o.Gender
}

// GetGenderOk returns a tuple with the Gender field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClaimRepresentation) GetGenderOk() (*bool, bool) {
	if o == nil || IsNil(o.Gender) {
		return nil, false
	}
	return o.Gender, true
}

// HasGender returns a boolean if a field has been set.
func (o *ClaimRepresentation) HasGender() bool {
	if o != nil && !IsNil(o.Gender) {
		return true
	}

	return false
}

// SetGender gets a reference to the given bool and assigns it to the Gender field.
func (o *ClaimRepresentation) SetGender(v bool) {
	o.Gender = &v
}

// GetLocale returns the Locale field value if set, zero value otherwise.
func (o *ClaimRepresentation) GetLocale() bool {
	if o == nil || IsNil(o.Locale) {
		var ret bool
		return ret
	}
	return *o.Locale
}

// GetLocaleOk returns a tuple with the Locale field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClaimRepresentation) GetLocaleOk() (*bool, bool) {
	if o == nil || IsNil(o.Locale) {
		return nil, false
	}
	return o.Locale, true
}

// HasLocale returns a boolean if a field has been set.
func (o *ClaimRepresentation) HasLocale() bool {
	if o != nil && !IsNil(o.Locale) {
		return true
	}

	return false
}

// SetLocale gets a reference to the given bool and assigns it to the Locale field.
func (o *ClaimRepresentation) SetLocale(v bool) {
	o.Locale = &v
}

// GetAddress returns the Address field value if set, zero value otherwise.
func (o *ClaimRepresentation) GetAddress() bool {
	if o == nil || IsNil(o.Address) {
		var ret bool
		return ret
	}
	return *o.Address
}

// GetAddressOk returns a tuple with the Address field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClaimRepresentation) GetAddressOk() (*bool, bool) {
	if o == nil || IsNil(o.Address) {
		return nil, false
	}
	return o.Address, true
}

// HasAddress returns a boolean if a field has been set.
func (o *ClaimRepresentation) HasAddress() bool {
	if o != nil && !IsNil(o.Address) {
		return true
	}

	return false
}

// SetAddress gets a reference to the given bool and assigns it to the Address field.
func (o *ClaimRepresentation) SetAddress(v bool) {
	o.Address = &v
}

// GetPhone returns the Phone field value if set, zero value otherwise.
func (o *ClaimRepresentation) GetPhone() bool {
	if o == nil || IsNil(o.Phone) {
		var ret bool
		return ret
	}
	return *o.Phone
}

// GetPhoneOk returns a tuple with the Phone field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClaimRepresentation) GetPhoneOk() (*bool, bool) {
	if o == nil || IsNil(o.Phone) {
		return nil, false
	}
	return o.Phone, true
}

// HasPhone returns a boolean if a field has been set.
func (o *ClaimRepresentation) HasPhone() bool {
	if o != nil && !IsNil(o.Phone) {
		return true
	}

	return false
}

// SetPhone gets a reference to the given bool and assigns it to the Phone field.
func (o *ClaimRepresentation) SetPhone(v bool) {
	o.Phone = &v
}

func (o ClaimRepresentation) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ClaimRepresentation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Username) {
		toSerialize["username"] = o.Username
	}
	if !IsNil(o.Profile) {
		toSerialize["profile"] = o.Profile
	}
	if !IsNil(o.Picture) {
		toSerialize["picture"] = o.Picture
	}
	if !IsNil(o.Website) {
		toSerialize["website"] = o.Website
	}
	if !IsNil(o.Email) {
		toSerialize["email"] = o.Email
	}
	if !IsNil(o.Gender) {
		toSerialize["gender"] = o.Gender
	}
	if !IsNil(o.Locale) {
		toSerialize["locale"] = o.Locale
	}
	if !IsNil(o.Address) {
		toSerialize["address"] = o.Address
	}
	if !IsNil(o.Phone) {
		toSerialize["phone"] = o.Phone
	}
	return toSerialize, nil
}

type NullableClaimRepresentation struct {
	value *ClaimRepresentation
	isSet bool
}

func (v NullableClaimRepresentation) Get() *ClaimRepresentation {
	return v.value
}

func (v *NullableClaimRepresentation) Set(val *ClaimRepresentation) {
	v.value = val
	v.isSet = true
}

func (v NullableClaimRepresentation) IsSet() bool {
	return v.isSet
}

func (v *NullableClaimRepresentation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClaimRepresentation(val *ClaimRepresentation) *NullableClaimRepresentation {
	return &NullableClaimRepresentation{value: val, isSet: true}
}

func (v NullableClaimRepresentation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClaimRepresentation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
