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

// checks if the AbstractPolicyRepresentation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AbstractPolicyRepresentation{}

// AbstractPolicyRepresentation struct for AbstractPolicyRepresentation
type AbstractPolicyRepresentation struct {
	Id               *string                  `json:"id,omitempty"`
	Name             *string                  `json:"name,omitempty"`
	Description      *string                  `json:"description,omitempty"`
	Type             *string                  `json:"type,omitempty"`
	Policies         []string                 `json:"policies,omitempty"`
	Resources        []string                 `json:"resources,omitempty"`
	Scopes           []string                 `json:"scopes,omitempty"`
	Logic            *Logic                   `json:"logic,omitempty"`
	DecisionStrategy *DecisionStrategy        `json:"decisionStrategy,omitempty"`
	Owner            *string                  `json:"owner,omitempty"`
	ResourceType     *string                  `json:"resourceType,omitempty"`
	ResourcesData    []ResourceRepresentation `json:"resourcesData,omitempty"`
	ScopesData       []ScopeRepresentation    `json:"scopesData,omitempty"`
}

// NewAbstractPolicyRepresentation instantiates a new AbstractPolicyRepresentation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAbstractPolicyRepresentation() *AbstractPolicyRepresentation {
	this := AbstractPolicyRepresentation{}
	return &this
}

// NewAbstractPolicyRepresentationWithDefaults instantiates a new AbstractPolicyRepresentation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAbstractPolicyRepresentationWithDefaults() *AbstractPolicyRepresentation {
	this := AbstractPolicyRepresentation{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *AbstractPolicyRepresentation) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AbstractPolicyRepresentation) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *AbstractPolicyRepresentation) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *AbstractPolicyRepresentation) SetId(v string) {
	o.Id = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *AbstractPolicyRepresentation) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AbstractPolicyRepresentation) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *AbstractPolicyRepresentation) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *AbstractPolicyRepresentation) SetName(v string) {
	o.Name = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *AbstractPolicyRepresentation) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AbstractPolicyRepresentation) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *AbstractPolicyRepresentation) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *AbstractPolicyRepresentation) SetDescription(v string) {
	o.Description = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *AbstractPolicyRepresentation) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AbstractPolicyRepresentation) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *AbstractPolicyRepresentation) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *AbstractPolicyRepresentation) SetType(v string) {
	o.Type = &v
}

// GetPolicies returns the Policies field value if set, zero value otherwise.
func (o *AbstractPolicyRepresentation) GetPolicies() []string {
	if o == nil || IsNil(o.Policies) {
		var ret []string
		return ret
	}
	return o.Policies
}

// GetPoliciesOk returns a tuple with the Policies field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AbstractPolicyRepresentation) GetPoliciesOk() ([]string, bool) {
	if o == nil || IsNil(o.Policies) {
		return nil, false
	}
	return o.Policies, true
}

// HasPolicies returns a boolean if a field has been set.
func (o *AbstractPolicyRepresentation) HasPolicies() bool {
	if o != nil && !IsNil(o.Policies) {
		return true
	}

	return false
}

// SetPolicies gets a reference to the given []string and assigns it to the Policies field.
func (o *AbstractPolicyRepresentation) SetPolicies(v []string) {
	o.Policies = v
}

// GetResources returns the Resources field value if set, zero value otherwise.
func (o *AbstractPolicyRepresentation) GetResources() []string {
	if o == nil || IsNil(o.Resources) {
		var ret []string
		return ret
	}
	return o.Resources
}

// GetResourcesOk returns a tuple with the Resources field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AbstractPolicyRepresentation) GetResourcesOk() ([]string, bool) {
	if o == nil || IsNil(o.Resources) {
		return nil, false
	}
	return o.Resources, true
}

// HasResources returns a boolean if a field has been set.
func (o *AbstractPolicyRepresentation) HasResources() bool {
	if o != nil && !IsNil(o.Resources) {
		return true
	}

	return false
}

// SetResources gets a reference to the given []string and assigns it to the Resources field.
func (o *AbstractPolicyRepresentation) SetResources(v []string) {
	o.Resources = v
}

// GetScopes returns the Scopes field value if set, zero value otherwise.
func (o *AbstractPolicyRepresentation) GetScopes() []string {
	if o == nil || IsNil(o.Scopes) {
		var ret []string
		return ret
	}
	return o.Scopes
}

// GetScopesOk returns a tuple with the Scopes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AbstractPolicyRepresentation) GetScopesOk() ([]string, bool) {
	if o == nil || IsNil(o.Scopes) {
		return nil, false
	}
	return o.Scopes, true
}

// HasScopes returns a boolean if a field has been set.
func (o *AbstractPolicyRepresentation) HasScopes() bool {
	if o != nil && !IsNil(o.Scopes) {
		return true
	}

	return false
}

// SetScopes gets a reference to the given []string and assigns it to the Scopes field.
func (o *AbstractPolicyRepresentation) SetScopes(v []string) {
	o.Scopes = v
}

// GetLogic returns the Logic field value if set, zero value otherwise.
func (o *AbstractPolicyRepresentation) GetLogic() Logic {
	if o == nil || IsNil(o.Logic) {
		var ret Logic
		return ret
	}
	return *o.Logic
}

// GetLogicOk returns a tuple with the Logic field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AbstractPolicyRepresentation) GetLogicOk() (*Logic, bool) {
	if o == nil || IsNil(o.Logic) {
		return nil, false
	}
	return o.Logic, true
}

// HasLogic returns a boolean if a field has been set.
func (o *AbstractPolicyRepresentation) HasLogic() bool {
	if o != nil && !IsNil(o.Logic) {
		return true
	}

	return false
}

// SetLogic gets a reference to the given Logic and assigns it to the Logic field.
func (o *AbstractPolicyRepresentation) SetLogic(v Logic) {
	o.Logic = &v
}

// GetDecisionStrategy returns the DecisionStrategy field value if set, zero value otherwise.
func (o *AbstractPolicyRepresentation) GetDecisionStrategy() DecisionStrategy {
	if o == nil || IsNil(o.DecisionStrategy) {
		var ret DecisionStrategy
		return ret
	}
	return *o.DecisionStrategy
}

// GetDecisionStrategyOk returns a tuple with the DecisionStrategy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AbstractPolicyRepresentation) GetDecisionStrategyOk() (*DecisionStrategy, bool) {
	if o == nil || IsNil(o.DecisionStrategy) {
		return nil, false
	}
	return o.DecisionStrategy, true
}

// HasDecisionStrategy returns a boolean if a field has been set.
func (o *AbstractPolicyRepresentation) HasDecisionStrategy() bool {
	if o != nil && !IsNil(o.DecisionStrategy) {
		return true
	}

	return false
}

// SetDecisionStrategy gets a reference to the given DecisionStrategy and assigns it to the DecisionStrategy field.
func (o *AbstractPolicyRepresentation) SetDecisionStrategy(v DecisionStrategy) {
	o.DecisionStrategy = &v
}

// GetOwner returns the Owner field value if set, zero value otherwise.
func (o *AbstractPolicyRepresentation) GetOwner() string {
	if o == nil || IsNil(o.Owner) {
		var ret string
		return ret
	}
	return *o.Owner
}

// GetOwnerOk returns a tuple with the Owner field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AbstractPolicyRepresentation) GetOwnerOk() (*string, bool) {
	if o == nil || IsNil(o.Owner) {
		return nil, false
	}
	return o.Owner, true
}

// HasOwner returns a boolean if a field has been set.
func (o *AbstractPolicyRepresentation) HasOwner() bool {
	if o != nil && !IsNil(o.Owner) {
		return true
	}

	return false
}

// SetOwner gets a reference to the given string and assigns it to the Owner field.
func (o *AbstractPolicyRepresentation) SetOwner(v string) {
	o.Owner = &v
}

// GetResourceType returns the ResourceType field value if set, zero value otherwise.
func (o *AbstractPolicyRepresentation) GetResourceType() string {
	if o == nil || IsNil(o.ResourceType) {
		var ret string
		return ret
	}
	return *o.ResourceType
}

// GetResourceTypeOk returns a tuple with the ResourceType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AbstractPolicyRepresentation) GetResourceTypeOk() (*string, bool) {
	if o == nil || IsNil(o.ResourceType) {
		return nil, false
	}
	return o.ResourceType, true
}

// HasResourceType returns a boolean if a field has been set.
func (o *AbstractPolicyRepresentation) HasResourceType() bool {
	if o != nil && !IsNil(o.ResourceType) {
		return true
	}

	return false
}

// SetResourceType gets a reference to the given string and assigns it to the ResourceType field.
func (o *AbstractPolicyRepresentation) SetResourceType(v string) {
	o.ResourceType = &v
}

// GetResourcesData returns the ResourcesData field value if set, zero value otherwise.
func (o *AbstractPolicyRepresentation) GetResourcesData() []ResourceRepresentation {
	if o == nil || IsNil(o.ResourcesData) {
		var ret []ResourceRepresentation
		return ret
	}
	return o.ResourcesData
}

// GetResourcesDataOk returns a tuple with the ResourcesData field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AbstractPolicyRepresentation) GetResourcesDataOk() ([]ResourceRepresentation, bool) {
	if o == nil || IsNil(o.ResourcesData) {
		return nil, false
	}
	return o.ResourcesData, true
}

// HasResourcesData returns a boolean if a field has been set.
func (o *AbstractPolicyRepresentation) HasResourcesData() bool {
	if o != nil && !IsNil(o.ResourcesData) {
		return true
	}

	return false
}

// SetResourcesData gets a reference to the given []ResourceRepresentation and assigns it to the ResourcesData field.
func (o *AbstractPolicyRepresentation) SetResourcesData(v []ResourceRepresentation) {
	o.ResourcesData = v
}

// GetScopesData returns the ScopesData field value if set, zero value otherwise.
func (o *AbstractPolicyRepresentation) GetScopesData() []ScopeRepresentation {
	if o == nil || IsNil(o.ScopesData) {
		var ret []ScopeRepresentation
		return ret
	}
	return o.ScopesData
}

// GetScopesDataOk returns a tuple with the ScopesData field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AbstractPolicyRepresentation) GetScopesDataOk() ([]ScopeRepresentation, bool) {
	if o == nil || IsNil(o.ScopesData) {
		return nil, false
	}
	return o.ScopesData, true
}

// HasScopesData returns a boolean if a field has been set.
func (o *AbstractPolicyRepresentation) HasScopesData() bool {
	if o != nil && !IsNil(o.ScopesData) {
		return true
	}

	return false
}

// SetScopesData gets a reference to the given []ScopeRepresentation and assigns it to the ScopesData field.
func (o *AbstractPolicyRepresentation) SetScopesData(v []ScopeRepresentation) {
	o.ScopesData = v
}

func (o AbstractPolicyRepresentation) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AbstractPolicyRepresentation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !IsNil(o.Policies) {
		toSerialize["policies"] = o.Policies
	}
	if !IsNil(o.Resources) {
		toSerialize["resources"] = o.Resources
	}
	if !IsNil(o.Scopes) {
		toSerialize["scopes"] = o.Scopes
	}
	if !IsNil(o.Logic) {
		toSerialize["logic"] = o.Logic
	}
	if !IsNil(o.DecisionStrategy) {
		toSerialize["decisionStrategy"] = o.DecisionStrategy
	}
	if !IsNil(o.Owner) {
		toSerialize["owner"] = o.Owner
	}
	if !IsNil(o.ResourceType) {
		toSerialize["resourceType"] = o.ResourceType
	}
	if !IsNil(o.ResourcesData) {
		toSerialize["resourcesData"] = o.ResourcesData
	}
	if !IsNil(o.ScopesData) {
		toSerialize["scopesData"] = o.ScopesData
	}
	return toSerialize, nil
}

type NullableAbstractPolicyRepresentation struct {
	value *AbstractPolicyRepresentation
	isSet bool
}

func (v NullableAbstractPolicyRepresentation) Get() *AbstractPolicyRepresentation {
	return v.value
}

func (v *NullableAbstractPolicyRepresentation) Set(val *AbstractPolicyRepresentation) {
	v.value = val
	v.isSet = true
}

func (v NullableAbstractPolicyRepresentation) IsSet() bool {
	return v.isSet
}

func (v *NullableAbstractPolicyRepresentation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAbstractPolicyRepresentation(val *AbstractPolicyRepresentation) *NullableAbstractPolicyRepresentation {
	return &NullableAbstractPolicyRepresentation{value: val, isSet: true}
}

func (v NullableAbstractPolicyRepresentation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAbstractPolicyRepresentation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
