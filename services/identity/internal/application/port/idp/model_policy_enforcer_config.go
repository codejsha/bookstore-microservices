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

// checks if the PolicyEnforcerConfig type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PolicyEnforcerConfig{}

// PolicyEnforcerConfig struct for PolicyEnforcerConfig
type PolicyEnforcerConfig struct {
	EnforcementMode       *EnforcementMode                   `json:"enforcement-mode,omitempty"`
	Paths                 []PathConfig                       `json:"paths,omitempty"`
	PathCache             *PathCacheConfig                   `json:"path-cache,omitempty"`
	LazyLoadPaths         *bool                              `json:"lazy-load-paths,omitempty"`
	OnDenyRedirectTo      *string                            `json:"on-deny-redirect-to,omitempty"`
	UserManagedAccess     map[string]interface{}             `json:"user-managed-access,omitempty"`
	ClaimInformationPoint *map[string]map[string]interface{} `json:"claim-information-point,omitempty"`
	HttpMethodAsScope     *bool                              `json:"http-method-as-scope,omitempty"`
	Realm                 *string                            `json:"realm,omitempty"`
	AuthServerUrl         *string                            `json:"auth-server-url,omitempty"`
	Credentials           map[string]interface{}             `json:"credentials,omitempty"`
	Resource              *string                            `json:"resource,omitempty"`
}

// NewPolicyEnforcerConfig instantiates a new PolicyEnforcerConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPolicyEnforcerConfig() *PolicyEnforcerConfig {
	this := PolicyEnforcerConfig{}
	return &this
}

// NewPolicyEnforcerConfigWithDefaults instantiates a new PolicyEnforcerConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPolicyEnforcerConfigWithDefaults() *PolicyEnforcerConfig {
	this := PolicyEnforcerConfig{}
	return &this
}

// GetEnforcementMode returns the EnforcementMode field value if set, zero value otherwise.
func (o *PolicyEnforcerConfig) GetEnforcementMode() EnforcementMode {
	if o == nil || IsNil(o.EnforcementMode) {
		var ret EnforcementMode
		return ret
	}
	return *o.EnforcementMode
}

// GetEnforcementModeOk returns a tuple with the EnforcementMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PolicyEnforcerConfig) GetEnforcementModeOk() (*EnforcementMode, bool) {
	if o == nil || IsNil(o.EnforcementMode) {
		return nil, false
	}
	return o.EnforcementMode, true
}

// HasEnforcementMode returns a boolean if a field has been set.
func (o *PolicyEnforcerConfig) HasEnforcementMode() bool {
	if o != nil && !IsNil(o.EnforcementMode) {
		return true
	}

	return false
}

// SetEnforcementMode gets a reference to the given EnforcementMode and assigns it to the EnforcementMode field.
func (o *PolicyEnforcerConfig) SetEnforcementMode(v EnforcementMode) {
	o.EnforcementMode = &v
}

// GetPaths returns the Paths field value if set, zero value otherwise.
func (o *PolicyEnforcerConfig) GetPaths() []PathConfig {
	if o == nil || IsNil(o.Paths) {
		var ret []PathConfig
		return ret
	}
	return o.Paths
}

// GetPathsOk returns a tuple with the Paths field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PolicyEnforcerConfig) GetPathsOk() ([]PathConfig, bool) {
	if o == nil || IsNil(o.Paths) {
		return nil, false
	}
	return o.Paths, true
}

// HasPaths returns a boolean if a field has been set.
func (o *PolicyEnforcerConfig) HasPaths() bool {
	if o != nil && !IsNil(o.Paths) {
		return true
	}

	return false
}

// SetPaths gets a reference to the given []PathConfig and assigns it to the Paths field.
func (o *PolicyEnforcerConfig) SetPaths(v []PathConfig) {
	o.Paths = v
}

// GetPathCache returns the PathCache field value if set, zero value otherwise.
func (o *PolicyEnforcerConfig) GetPathCache() PathCacheConfig {
	if o == nil || IsNil(o.PathCache) {
		var ret PathCacheConfig
		return ret
	}
	return *o.PathCache
}

// GetPathCacheOk returns a tuple with the PathCache field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PolicyEnforcerConfig) GetPathCacheOk() (*PathCacheConfig, bool) {
	if o == nil || IsNil(o.PathCache) {
		return nil, false
	}
	return o.PathCache, true
}

// HasPathCache returns a boolean if a field has been set.
func (o *PolicyEnforcerConfig) HasPathCache() bool {
	if o != nil && !IsNil(o.PathCache) {
		return true
	}

	return false
}

// SetPathCache gets a reference to the given PathCacheConfig and assigns it to the PathCache field.
func (o *PolicyEnforcerConfig) SetPathCache(v PathCacheConfig) {
	o.PathCache = &v
}

// GetLazyLoadPaths returns the LazyLoadPaths field value if set, zero value otherwise.
func (o *PolicyEnforcerConfig) GetLazyLoadPaths() bool {
	if o == nil || IsNil(o.LazyLoadPaths) {
		var ret bool
		return ret
	}
	return *o.LazyLoadPaths
}

// GetLazyLoadPathsOk returns a tuple with the LazyLoadPaths field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PolicyEnforcerConfig) GetLazyLoadPathsOk() (*bool, bool) {
	if o == nil || IsNil(o.LazyLoadPaths) {
		return nil, false
	}
	return o.LazyLoadPaths, true
}

// HasLazyLoadPaths returns a boolean if a field has been set.
func (o *PolicyEnforcerConfig) HasLazyLoadPaths() bool {
	if o != nil && !IsNil(o.LazyLoadPaths) {
		return true
	}

	return false
}

// SetLazyLoadPaths gets a reference to the given bool and assigns it to the LazyLoadPaths field.
func (o *PolicyEnforcerConfig) SetLazyLoadPaths(v bool) {
	o.LazyLoadPaths = &v
}

// GetOnDenyRedirectTo returns the OnDenyRedirectTo field value if set, zero value otherwise.
func (o *PolicyEnforcerConfig) GetOnDenyRedirectTo() string {
	if o == nil || IsNil(o.OnDenyRedirectTo) {
		var ret string
		return ret
	}
	return *o.OnDenyRedirectTo
}

// GetOnDenyRedirectToOk returns a tuple with the OnDenyRedirectTo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PolicyEnforcerConfig) GetOnDenyRedirectToOk() (*string, bool) {
	if o == nil || IsNil(o.OnDenyRedirectTo) {
		return nil, false
	}
	return o.OnDenyRedirectTo, true
}

// HasOnDenyRedirectTo returns a boolean if a field has been set.
func (o *PolicyEnforcerConfig) HasOnDenyRedirectTo() bool {
	if o != nil && !IsNil(o.OnDenyRedirectTo) {
		return true
	}

	return false
}

// SetOnDenyRedirectTo gets a reference to the given string and assigns it to the OnDenyRedirectTo field.
func (o *PolicyEnforcerConfig) SetOnDenyRedirectTo(v string) {
	o.OnDenyRedirectTo = &v
}

// GetUserManagedAccess returns the UserManagedAccess field value if set, zero value otherwise.
func (o *PolicyEnforcerConfig) GetUserManagedAccess() map[string]interface{} {
	if o == nil || IsNil(o.UserManagedAccess) {
		var ret map[string]interface{}
		return ret
	}
	return o.UserManagedAccess
}

// GetUserManagedAccessOk returns a tuple with the UserManagedAccess field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PolicyEnforcerConfig) GetUserManagedAccessOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.UserManagedAccess) {
		return map[string]interface{}{}, false
	}
	return o.UserManagedAccess, true
}

// HasUserManagedAccess returns a boolean if a field has been set.
func (o *PolicyEnforcerConfig) HasUserManagedAccess() bool {
	if o != nil && !IsNil(o.UserManagedAccess) {
		return true
	}

	return false
}

// SetUserManagedAccess gets a reference to the given map[string]interface{} and assigns it to the UserManagedAccess field.
func (o *PolicyEnforcerConfig) SetUserManagedAccess(v map[string]interface{}) {
	o.UserManagedAccess = v
}

// GetClaimInformationPoint returns the ClaimInformationPoint field value if set, zero value otherwise.
func (o *PolicyEnforcerConfig) GetClaimInformationPoint() map[string]map[string]interface{} {
	if o == nil || IsNil(o.ClaimInformationPoint) {
		var ret map[string]map[string]interface{}
		return ret
	}
	return *o.ClaimInformationPoint
}

// GetClaimInformationPointOk returns a tuple with the ClaimInformationPoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PolicyEnforcerConfig) GetClaimInformationPointOk() (*map[string]map[string]interface{}, bool) {
	if o == nil || IsNil(o.ClaimInformationPoint) {
		return nil, false
	}
	return o.ClaimInformationPoint, true
}

// HasClaimInformationPoint returns a boolean if a field has been set.
func (o *PolicyEnforcerConfig) HasClaimInformationPoint() bool {
	if o != nil && !IsNil(o.ClaimInformationPoint) {
		return true
	}

	return false
}

// SetClaimInformationPoint gets a reference to the given map[string]map[string]interface{} and assigns it to the ClaimInformationPoint field.
func (o *PolicyEnforcerConfig) SetClaimInformationPoint(v map[string]map[string]interface{}) {
	o.ClaimInformationPoint = &v
}

// GetHttpMethodAsScope returns the HttpMethodAsScope field value if set, zero value otherwise.
func (o *PolicyEnforcerConfig) GetHttpMethodAsScope() bool {
	if o == nil || IsNil(o.HttpMethodAsScope) {
		var ret bool
		return ret
	}
	return *o.HttpMethodAsScope
}

// GetHttpMethodAsScopeOk returns a tuple with the HttpMethodAsScope field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PolicyEnforcerConfig) GetHttpMethodAsScopeOk() (*bool, bool) {
	if o == nil || IsNil(o.HttpMethodAsScope) {
		return nil, false
	}
	return o.HttpMethodAsScope, true
}

// HasHttpMethodAsScope returns a boolean if a field has been set.
func (o *PolicyEnforcerConfig) HasHttpMethodAsScope() bool {
	if o != nil && !IsNil(o.HttpMethodAsScope) {
		return true
	}

	return false
}

// SetHttpMethodAsScope gets a reference to the given bool and assigns it to the HttpMethodAsScope field.
func (o *PolicyEnforcerConfig) SetHttpMethodAsScope(v bool) {
	o.HttpMethodAsScope = &v
}

// GetRealm returns the Realm field value if set, zero value otherwise.
func (o *PolicyEnforcerConfig) GetRealm() string {
	if o == nil || IsNil(o.Realm) {
		var ret string
		return ret
	}
	return *o.Realm
}

// GetRealmOk returns a tuple with the Realm field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PolicyEnforcerConfig) GetRealmOk() (*string, bool) {
	if o == nil || IsNil(o.Realm) {
		return nil, false
	}
	return o.Realm, true
}

// HasRealm returns a boolean if a field has been set.
func (o *PolicyEnforcerConfig) HasRealm() bool {
	if o != nil && !IsNil(o.Realm) {
		return true
	}

	return false
}

// SetRealm gets a reference to the given string and assigns it to the Realm field.
func (o *PolicyEnforcerConfig) SetRealm(v string) {
	o.Realm = &v
}

// GetAuthServerUrl returns the AuthServerUrl field value if set, zero value otherwise.
func (o *PolicyEnforcerConfig) GetAuthServerUrl() string {
	if o == nil || IsNil(o.AuthServerUrl) {
		var ret string
		return ret
	}
	return *o.AuthServerUrl
}

// GetAuthServerUrlOk returns a tuple with the AuthServerUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PolicyEnforcerConfig) GetAuthServerUrlOk() (*string, bool) {
	if o == nil || IsNil(o.AuthServerUrl) {
		return nil, false
	}
	return o.AuthServerUrl, true
}

// HasAuthServerUrl returns a boolean if a field has been set.
func (o *PolicyEnforcerConfig) HasAuthServerUrl() bool {
	if o != nil && !IsNil(o.AuthServerUrl) {
		return true
	}

	return false
}

// SetAuthServerUrl gets a reference to the given string and assigns it to the AuthServerUrl field.
func (o *PolicyEnforcerConfig) SetAuthServerUrl(v string) {
	o.AuthServerUrl = &v
}

// GetCredentials returns the Credentials field value if set, zero value otherwise.
func (o *PolicyEnforcerConfig) GetCredentials() map[string]interface{} {
	if o == nil || IsNil(o.Credentials) {
		var ret map[string]interface{}
		return ret
	}
	return o.Credentials
}

// GetCredentialsOk returns a tuple with the Credentials field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PolicyEnforcerConfig) GetCredentialsOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.Credentials) {
		return map[string]interface{}{}, false
	}
	return o.Credentials, true
}

// HasCredentials returns a boolean if a field has been set.
func (o *PolicyEnforcerConfig) HasCredentials() bool {
	if o != nil && !IsNil(o.Credentials) {
		return true
	}

	return false
}

// SetCredentials gets a reference to the given map[string]interface{} and assigns it to the Credentials field.
func (o *PolicyEnforcerConfig) SetCredentials(v map[string]interface{}) {
	o.Credentials = v
}

// GetResource returns the Resource field value if set, zero value otherwise.
func (o *PolicyEnforcerConfig) GetResource() string {
	if o == nil || IsNil(o.Resource) {
		var ret string
		return ret
	}
	return *o.Resource
}

// GetResourceOk returns a tuple with the Resource field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PolicyEnforcerConfig) GetResourceOk() (*string, bool) {
	if o == nil || IsNil(o.Resource) {
		return nil, false
	}
	return o.Resource, true
}

// HasResource returns a boolean if a field has been set.
func (o *PolicyEnforcerConfig) HasResource() bool {
	if o != nil && !IsNil(o.Resource) {
		return true
	}

	return false
}

// SetResource gets a reference to the given string and assigns it to the Resource field.
func (o *PolicyEnforcerConfig) SetResource(v string) {
	o.Resource = &v
}

func (o PolicyEnforcerConfig) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PolicyEnforcerConfig) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.EnforcementMode) {
		toSerialize["enforcement-mode"] = o.EnforcementMode
	}
	if !IsNil(o.Paths) {
		toSerialize["paths"] = o.Paths
	}
	if !IsNil(o.PathCache) {
		toSerialize["path-cache"] = o.PathCache
	}
	if !IsNil(o.LazyLoadPaths) {
		toSerialize["lazy-load-paths"] = o.LazyLoadPaths
	}
	if !IsNil(o.OnDenyRedirectTo) {
		toSerialize["on-deny-redirect-to"] = o.OnDenyRedirectTo
	}
	if !IsNil(o.UserManagedAccess) {
		toSerialize["user-managed-access"] = o.UserManagedAccess
	}
	if !IsNil(o.ClaimInformationPoint) {
		toSerialize["claim-information-point"] = o.ClaimInformationPoint
	}
	if !IsNil(o.HttpMethodAsScope) {
		toSerialize["http-method-as-scope"] = o.HttpMethodAsScope
	}
	if !IsNil(o.Realm) {
		toSerialize["realm"] = o.Realm
	}
	if !IsNil(o.AuthServerUrl) {
		toSerialize["auth-server-url"] = o.AuthServerUrl
	}
	if !IsNil(o.Credentials) {
		toSerialize["credentials"] = o.Credentials
	}
	if !IsNil(o.Resource) {
		toSerialize["resource"] = o.Resource
	}
	return toSerialize, nil
}

type NullablePolicyEnforcerConfig struct {
	value *PolicyEnforcerConfig
	isSet bool
}

func (v NullablePolicyEnforcerConfig) Get() *PolicyEnforcerConfig {
	return v.value
}

func (v *NullablePolicyEnforcerConfig) Set(val *PolicyEnforcerConfig) {
	v.value = val
	v.isSet = true
}

func (v NullablePolicyEnforcerConfig) IsSet() bool {
	return v.isSet
}

func (v *NullablePolicyEnforcerConfig) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePolicyEnforcerConfig(val *PolicyEnforcerConfig) *NullablePolicyEnforcerConfig {
	return &NullablePolicyEnforcerConfig{value: val, isSet: true}
}

func (v NullablePolicyEnforcerConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePolicyEnforcerConfig) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
