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

// checks if the UserRepresentation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UserRepresentation{}

// UserRepresentation struct for UserRepresentation
type UserRepresentation struct {
	Id                         *string                           `json:"id,omitempty"`
	Username                   *string                           `json:"username,omitempty"`
	FirstName                  *string                           `json:"firstName,omitempty"`
	LastName                   *string                           `json:"lastName,omitempty"`
	Email                      *string                           `json:"email,omitempty"`
	EmailVerified              *bool                             `json:"emailVerified,omitempty"`
	Attributes                 *map[string][]string              `json:"attributes,omitempty"`
	UserProfileMetadata        *UserProfileMetadata              `json:"userProfileMetadata,omitempty"`
	Self                       *string                           `json:"self,omitempty"`
	Origin                     *string                           `json:"origin,omitempty"`
	CreatedTimestamp           *int64                            `json:"createdTimestamp,omitempty"`
	Enabled                    *bool                             `json:"enabled,omitempty"`
	Totp                       *bool                             `json:"totp,omitempty"`
	FederationLink             *string                           `json:"federationLink,omitempty"`
	ServiceAccountClientId     *string                           `json:"serviceAccountClientId,omitempty"`
	Credentials                []CredentialRepresentation        `json:"credentials,omitempty"`
	DisableableCredentialTypes []string                          `json:"disableableCredentialTypes,omitempty"`
	RequiredActions            []string                          `json:"requiredActions,omitempty"`
	FederatedIdentities        []FederatedIdentityRepresentation `json:"federatedIdentities,omitempty"`
	RealmRoles                 []string                          `json:"realmRoles,omitempty"`
	ClientRoles                *map[string][]string              `json:"clientRoles,omitempty"`
	ClientConsents             []UserConsentRepresentation       `json:"clientConsents,omitempty"`
	NotBefore                  *int32                            `json:"notBefore,omitempty"`
	// Deprecated
	ApplicationRoles *map[string][]string `json:"applicationRoles,omitempty"`
	// Deprecated
	SocialLinks []SocialLinkRepresentation `json:"socialLinks,omitempty"`
	Groups      []string                   `json:"groups,omitempty"`
	Access      *map[string]bool           `json:"access,omitempty"`
}

// NewUserRepresentation instantiates a new UserRepresentation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserRepresentation() *UserRepresentation {
	this := UserRepresentation{}
	return &this
}

// NewUserRepresentationWithDefaults instantiates a new UserRepresentation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserRepresentationWithDefaults() *UserRepresentation {
	this := UserRepresentation{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *UserRepresentation) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *UserRepresentation) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *UserRepresentation) SetId(v string) {
	o.Id = &v
}

// GetUsername returns the Username field value if set, zero value otherwise.
func (o *UserRepresentation) GetUsername() string {
	if o == nil || IsNil(o.Username) {
		var ret string
		return ret
	}
	return *o.Username
}

// GetUsernameOk returns a tuple with the Username field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.Username) {
		return nil, false
	}
	return o.Username, true
}

// HasUsername returns a boolean if a field has been set.
func (o *UserRepresentation) HasUsername() bool {
	if o != nil && !IsNil(o.Username) {
		return true
	}

	return false
}

// SetUsername gets a reference to the given string and assigns it to the Username field.
func (o *UserRepresentation) SetUsername(v string) {
	o.Username = &v
}

// GetFirstName returns the FirstName field value if set, zero value otherwise.
func (o *UserRepresentation) GetFirstName() string {
	if o == nil || IsNil(o.FirstName) {
		var ret string
		return ret
	}
	return *o.FirstName
}

// GetFirstNameOk returns a tuple with the FirstName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetFirstNameOk() (*string, bool) {
	if o == nil || IsNil(o.FirstName) {
		return nil, false
	}
	return o.FirstName, true
}

// HasFirstName returns a boolean if a field has been set.
func (o *UserRepresentation) HasFirstName() bool {
	if o != nil && !IsNil(o.FirstName) {
		return true
	}

	return false
}

// SetFirstName gets a reference to the given string and assigns it to the FirstName field.
func (o *UserRepresentation) SetFirstName(v string) {
	o.FirstName = &v
}

// GetLastName returns the LastName field value if set, zero value otherwise.
func (o *UserRepresentation) GetLastName() string {
	if o == nil || IsNil(o.LastName) {
		var ret string
		return ret
	}
	return *o.LastName
}

// GetLastNameOk returns a tuple with the LastName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetLastNameOk() (*string, bool) {
	if o == nil || IsNil(o.LastName) {
		return nil, false
	}
	return o.LastName, true
}

// HasLastName returns a boolean if a field has been set.
func (o *UserRepresentation) HasLastName() bool {
	if o != nil && !IsNil(o.LastName) {
		return true
	}

	return false
}

// SetLastName gets a reference to the given string and assigns it to the LastName field.
func (o *UserRepresentation) SetLastName(v string) {
	o.LastName = &v
}

// GetEmail returns the Email field value if set, zero value otherwise.
func (o *UserRepresentation) GetEmail() string {
	if o == nil || IsNil(o.Email) {
		var ret string
		return ret
	}
	return *o.Email
}

// GetEmailOk returns a tuple with the Email field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetEmailOk() (*string, bool) {
	if o == nil || IsNil(o.Email) {
		return nil, false
	}
	return o.Email, true
}

// HasEmail returns a boolean if a field has been set.
func (o *UserRepresentation) HasEmail() bool {
	if o != nil && !IsNil(o.Email) {
		return true
	}

	return false
}

// SetEmail gets a reference to the given string and assigns it to the Email field.
func (o *UserRepresentation) SetEmail(v string) {
	o.Email = &v
}

// GetEmailVerified returns the EmailVerified field value if set, zero value otherwise.
func (o *UserRepresentation) GetEmailVerified() bool {
	if o == nil || IsNil(o.EmailVerified) {
		var ret bool
		return ret
	}
	return *o.EmailVerified
}

// GetEmailVerifiedOk returns a tuple with the EmailVerified field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetEmailVerifiedOk() (*bool, bool) {
	if o == nil || IsNil(o.EmailVerified) {
		return nil, false
	}
	return o.EmailVerified, true
}

// HasEmailVerified returns a boolean if a field has been set.
func (o *UserRepresentation) HasEmailVerified() bool {
	if o != nil && !IsNil(o.EmailVerified) {
		return true
	}

	return false
}

// SetEmailVerified gets a reference to the given bool and assigns it to the EmailVerified field.
func (o *UserRepresentation) SetEmailVerified(v bool) {
	o.EmailVerified = &v
}

// GetAttributes returns the Attributes field value if set, zero value otherwise.
func (o *UserRepresentation) GetAttributes() map[string][]string {
	if o == nil || IsNil(o.Attributes) {
		var ret map[string][]string
		return ret
	}
	return *o.Attributes
}

// GetAttributesOk returns a tuple with the Attributes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetAttributesOk() (*map[string][]string, bool) {
	if o == nil || IsNil(o.Attributes) {
		return nil, false
	}
	return o.Attributes, true
}

// HasAttributes returns a boolean if a field has been set.
func (o *UserRepresentation) HasAttributes() bool {
	if o != nil && !IsNil(o.Attributes) {
		return true
	}

	return false
}

// SetAttributes gets a reference to the given map[string][]string and assigns it to the Attributes field.
func (o *UserRepresentation) SetAttributes(v map[string][]string) {
	o.Attributes = &v
}

// GetUserProfileMetadata returns the UserProfileMetadata field value if set, zero value otherwise.
func (o *UserRepresentation) GetUserProfileMetadata() UserProfileMetadata {
	if o == nil || IsNil(o.UserProfileMetadata) {
		var ret UserProfileMetadata
		return ret
	}
	return *o.UserProfileMetadata
}

// GetUserProfileMetadataOk returns a tuple with the UserProfileMetadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetUserProfileMetadataOk() (*UserProfileMetadata, bool) {
	if o == nil || IsNil(o.UserProfileMetadata) {
		return nil, false
	}
	return o.UserProfileMetadata, true
}

// HasUserProfileMetadata returns a boolean if a field has been set.
func (o *UserRepresentation) HasUserProfileMetadata() bool {
	if o != nil && !IsNil(o.UserProfileMetadata) {
		return true
	}

	return false
}

// SetUserProfileMetadata gets a reference to the given UserProfileMetadata and assigns it to the UserProfileMetadata field.
func (o *UserRepresentation) SetUserProfileMetadata(v UserProfileMetadata) {
	o.UserProfileMetadata = &v
}

// GetSelf returns the Self field value if set, zero value otherwise.
func (o *UserRepresentation) GetSelf() string {
	if o == nil || IsNil(o.Self) {
		var ret string
		return ret
	}
	return *o.Self
}

// GetSelfOk returns a tuple with the Self field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetSelfOk() (*string, bool) {
	if o == nil || IsNil(o.Self) {
		return nil, false
	}
	return o.Self, true
}

// HasSelf returns a boolean if a field has been set.
func (o *UserRepresentation) HasSelf() bool {
	if o != nil && !IsNil(o.Self) {
		return true
	}

	return false
}

// SetSelf gets a reference to the given string and assigns it to the Self field.
func (o *UserRepresentation) SetSelf(v string) {
	o.Self = &v
}

// GetOrigin returns the Origin field value if set, zero value otherwise.
func (o *UserRepresentation) GetOrigin() string {
	if o == nil || IsNil(o.Origin) {
		var ret string
		return ret
	}
	return *o.Origin
}

// GetOriginOk returns a tuple with the Origin field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetOriginOk() (*string, bool) {
	if o == nil || IsNil(o.Origin) {
		return nil, false
	}
	return o.Origin, true
}

// HasOrigin returns a boolean if a field has been set.
func (o *UserRepresentation) HasOrigin() bool {
	if o != nil && !IsNil(o.Origin) {
		return true
	}

	return false
}

// SetOrigin gets a reference to the given string and assigns it to the Origin field.
func (o *UserRepresentation) SetOrigin(v string) {
	o.Origin = &v
}

// GetCreatedTimestamp returns the CreatedTimestamp field value if set, zero value otherwise.
func (o *UserRepresentation) GetCreatedTimestamp() int64 {
	if o == nil || IsNil(o.CreatedTimestamp) {
		var ret int64
		return ret
	}
	return *o.CreatedTimestamp
}

// GetCreatedTimestampOk returns a tuple with the CreatedTimestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetCreatedTimestampOk() (*int64, bool) {
	if o == nil || IsNil(o.CreatedTimestamp) {
		return nil, false
	}
	return o.CreatedTimestamp, true
}

// HasCreatedTimestamp returns a boolean if a field has been set.
func (o *UserRepresentation) HasCreatedTimestamp() bool {
	if o != nil && !IsNil(o.CreatedTimestamp) {
		return true
	}

	return false
}

// SetCreatedTimestamp gets a reference to the given int64 and assigns it to the CreatedTimestamp field.
func (o *UserRepresentation) SetCreatedTimestamp(v int64) {
	o.CreatedTimestamp = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *UserRepresentation) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled) {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.Enabled) {
		return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *UserRepresentation) HasEnabled() bool {
	if o != nil && !IsNil(o.Enabled) {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *UserRepresentation) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetTotp returns the Totp field value if set, zero value otherwise.
func (o *UserRepresentation) GetTotp() bool {
	if o == nil || IsNil(o.Totp) {
		var ret bool
		return ret
	}
	return *o.Totp
}

// GetTotpOk returns a tuple with the Totp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetTotpOk() (*bool, bool) {
	if o == nil || IsNil(o.Totp) {
		return nil, false
	}
	return o.Totp, true
}

// HasTotp returns a boolean if a field has been set.
func (o *UserRepresentation) HasTotp() bool {
	if o != nil && !IsNil(o.Totp) {
		return true
	}

	return false
}

// SetTotp gets a reference to the given bool and assigns it to the Totp field.
func (o *UserRepresentation) SetTotp(v bool) {
	o.Totp = &v
}

// GetFederationLink returns the FederationLink field value if set, zero value otherwise.
func (o *UserRepresentation) GetFederationLink() string {
	if o == nil || IsNil(o.FederationLink) {
		var ret string
		return ret
	}
	return *o.FederationLink
}

// GetFederationLinkOk returns a tuple with the FederationLink field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetFederationLinkOk() (*string, bool) {
	if o == nil || IsNil(o.FederationLink) {
		return nil, false
	}
	return o.FederationLink, true
}

// HasFederationLink returns a boolean if a field has been set.
func (o *UserRepresentation) HasFederationLink() bool {
	if o != nil && !IsNil(o.FederationLink) {
		return true
	}

	return false
}

// SetFederationLink gets a reference to the given string and assigns it to the FederationLink field.
func (o *UserRepresentation) SetFederationLink(v string) {
	o.FederationLink = &v
}

// GetServiceAccountClientId returns the ServiceAccountClientId field value if set, zero value otherwise.
func (o *UserRepresentation) GetServiceAccountClientId() string {
	if o == nil || IsNil(o.ServiceAccountClientId) {
		var ret string
		return ret
	}
	return *o.ServiceAccountClientId
}

// GetServiceAccountClientIdOk returns a tuple with the ServiceAccountClientId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetServiceAccountClientIdOk() (*string, bool) {
	if o == nil || IsNil(o.ServiceAccountClientId) {
		return nil, false
	}
	return o.ServiceAccountClientId, true
}

// HasServiceAccountClientId returns a boolean if a field has been set.
func (o *UserRepresentation) HasServiceAccountClientId() bool {
	if o != nil && !IsNil(o.ServiceAccountClientId) {
		return true
	}

	return false
}

// SetServiceAccountClientId gets a reference to the given string and assigns it to the ServiceAccountClientId field.
func (o *UserRepresentation) SetServiceAccountClientId(v string) {
	o.ServiceAccountClientId = &v
}

// GetCredentials returns the Credentials field value if set, zero value otherwise.
func (o *UserRepresentation) GetCredentials() []CredentialRepresentation {
	if o == nil || IsNil(o.Credentials) {
		var ret []CredentialRepresentation
		return ret
	}
	return o.Credentials
}

// GetCredentialsOk returns a tuple with the Credentials field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetCredentialsOk() ([]CredentialRepresentation, bool) {
	if o == nil || IsNil(o.Credentials) {
		return nil, false
	}
	return o.Credentials, true
}

// HasCredentials returns a boolean if a field has been set.
func (o *UserRepresentation) HasCredentials() bool {
	if o != nil && !IsNil(o.Credentials) {
		return true
	}

	return false
}

// SetCredentials gets a reference to the given []CredentialRepresentation and assigns it to the Credentials field.
func (o *UserRepresentation) SetCredentials(v []CredentialRepresentation) {
	o.Credentials = v
}

// GetDisableableCredentialTypes returns the DisableableCredentialTypes field value if set, zero value otherwise.
func (o *UserRepresentation) GetDisableableCredentialTypes() []string {
	if o == nil || IsNil(o.DisableableCredentialTypes) {
		var ret []string
		return ret
	}
	return o.DisableableCredentialTypes
}

// GetDisableableCredentialTypesOk returns a tuple with the DisableableCredentialTypes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetDisableableCredentialTypesOk() ([]string, bool) {
	if o == nil || IsNil(o.DisableableCredentialTypes) {
		return nil, false
	}
	return o.DisableableCredentialTypes, true
}

// HasDisableableCredentialTypes returns a boolean if a field has been set.
func (o *UserRepresentation) HasDisableableCredentialTypes() bool {
	if o != nil && !IsNil(o.DisableableCredentialTypes) {
		return true
	}

	return false
}

// SetDisableableCredentialTypes gets a reference to the given []string and assigns it to the DisableableCredentialTypes field.
func (o *UserRepresentation) SetDisableableCredentialTypes(v []string) {
	o.DisableableCredentialTypes = v
}

// GetRequiredActions returns the RequiredActions field value if set, zero value otherwise.
func (o *UserRepresentation) GetRequiredActions() []string {
	if o == nil || IsNil(o.RequiredActions) {
		var ret []string
		return ret
	}
	return o.RequiredActions
}

// GetRequiredActionsOk returns a tuple with the RequiredActions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetRequiredActionsOk() ([]string, bool) {
	if o == nil || IsNil(o.RequiredActions) {
		return nil, false
	}
	return o.RequiredActions, true
}

// HasRequiredActions returns a boolean if a field has been set.
func (o *UserRepresentation) HasRequiredActions() bool {
	if o != nil && !IsNil(o.RequiredActions) {
		return true
	}

	return false
}

// SetRequiredActions gets a reference to the given []string and assigns it to the RequiredActions field.
func (o *UserRepresentation) SetRequiredActions(v []string) {
	o.RequiredActions = v
}

// GetFederatedIdentities returns the FederatedIdentities field value if set, zero value otherwise.
func (o *UserRepresentation) GetFederatedIdentities() []FederatedIdentityRepresentation {
	if o == nil || IsNil(o.FederatedIdentities) {
		var ret []FederatedIdentityRepresentation
		return ret
	}
	return o.FederatedIdentities
}

// GetFederatedIdentitiesOk returns a tuple with the FederatedIdentities field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetFederatedIdentitiesOk() ([]FederatedIdentityRepresentation, bool) {
	if o == nil || IsNil(o.FederatedIdentities) {
		return nil, false
	}
	return o.FederatedIdentities, true
}

// HasFederatedIdentities returns a boolean if a field has been set.
func (o *UserRepresentation) HasFederatedIdentities() bool {
	if o != nil && !IsNil(o.FederatedIdentities) {
		return true
	}

	return false
}

// SetFederatedIdentities gets a reference to the given []FederatedIdentityRepresentation and assigns it to the FederatedIdentities field.
func (o *UserRepresentation) SetFederatedIdentities(v []FederatedIdentityRepresentation) {
	o.FederatedIdentities = v
}

// GetRealmRoles returns the RealmRoles field value if set, zero value otherwise.
func (o *UserRepresentation) GetRealmRoles() []string {
	if o == nil || IsNil(o.RealmRoles) {
		var ret []string
		return ret
	}
	return o.RealmRoles
}

// GetRealmRolesOk returns a tuple with the RealmRoles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetRealmRolesOk() ([]string, bool) {
	if o == nil || IsNil(o.RealmRoles) {
		return nil, false
	}
	return o.RealmRoles, true
}

// HasRealmRoles returns a boolean if a field has been set.
func (o *UserRepresentation) HasRealmRoles() bool {
	if o != nil && !IsNil(o.RealmRoles) {
		return true
	}

	return false
}

// SetRealmRoles gets a reference to the given []string and assigns it to the RealmRoles field.
func (o *UserRepresentation) SetRealmRoles(v []string) {
	o.RealmRoles = v
}

// GetClientRoles returns the ClientRoles field value if set, zero value otherwise.
func (o *UserRepresentation) GetClientRoles() map[string][]string {
	if o == nil || IsNil(o.ClientRoles) {
		var ret map[string][]string
		return ret
	}
	return *o.ClientRoles
}

// GetClientRolesOk returns a tuple with the ClientRoles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetClientRolesOk() (*map[string][]string, bool) {
	if o == nil || IsNil(o.ClientRoles) {
		return nil, false
	}
	return o.ClientRoles, true
}

// HasClientRoles returns a boolean if a field has been set.
func (o *UserRepresentation) HasClientRoles() bool {
	if o != nil && !IsNil(o.ClientRoles) {
		return true
	}

	return false
}

// SetClientRoles gets a reference to the given map[string][]string and assigns it to the ClientRoles field.
func (o *UserRepresentation) SetClientRoles(v map[string][]string) {
	o.ClientRoles = &v
}

// GetClientConsents returns the ClientConsents field value if set, zero value otherwise.
func (o *UserRepresentation) GetClientConsents() []UserConsentRepresentation {
	if o == nil || IsNil(o.ClientConsents) {
		var ret []UserConsentRepresentation
		return ret
	}
	return o.ClientConsents
}

// GetClientConsentsOk returns a tuple with the ClientConsents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetClientConsentsOk() ([]UserConsentRepresentation, bool) {
	if o == nil || IsNil(o.ClientConsents) {
		return nil, false
	}
	return o.ClientConsents, true
}

// HasClientConsents returns a boolean if a field has been set.
func (o *UserRepresentation) HasClientConsents() bool {
	if o != nil && !IsNil(o.ClientConsents) {
		return true
	}

	return false
}

// SetClientConsents gets a reference to the given []UserConsentRepresentation and assigns it to the ClientConsents field.
func (o *UserRepresentation) SetClientConsents(v []UserConsentRepresentation) {
	o.ClientConsents = v
}

// GetNotBefore returns the NotBefore field value if set, zero value otherwise.
func (o *UserRepresentation) GetNotBefore() int32 {
	if o == nil || IsNil(o.NotBefore) {
		var ret int32
		return ret
	}
	return *o.NotBefore
}

// GetNotBeforeOk returns a tuple with the NotBefore field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetNotBeforeOk() (*int32, bool) {
	if o == nil || IsNil(o.NotBefore) {
		return nil, false
	}
	return o.NotBefore, true
}

// HasNotBefore returns a boolean if a field has been set.
func (o *UserRepresentation) HasNotBefore() bool {
	if o != nil && !IsNil(o.NotBefore) {
		return true
	}

	return false
}

// SetNotBefore gets a reference to the given int32 and assigns it to the NotBefore field.
func (o *UserRepresentation) SetNotBefore(v int32) {
	o.NotBefore = &v
}

// GetApplicationRoles returns the ApplicationRoles field value if set, zero value otherwise.
// Deprecated
func (o *UserRepresentation) GetApplicationRoles() map[string][]string {
	if o == nil || IsNil(o.ApplicationRoles) {
		var ret map[string][]string
		return ret
	}
	return *o.ApplicationRoles
}

// GetApplicationRolesOk returns a tuple with the ApplicationRoles field value if set, nil otherwise
// and a boolean to check if the value has been set.
// Deprecated
func (o *UserRepresentation) GetApplicationRolesOk() (*map[string][]string, bool) {
	if o == nil || IsNil(o.ApplicationRoles) {
		return nil, false
	}
	return o.ApplicationRoles, true
}

// HasApplicationRoles returns a boolean if a field has been set.
func (o *UserRepresentation) HasApplicationRoles() bool {
	if o != nil && !IsNil(o.ApplicationRoles) {
		return true
	}

	return false
}

// SetApplicationRoles gets a reference to the given map[string][]string and assigns it to the ApplicationRoles field.
// Deprecated
func (o *UserRepresentation) SetApplicationRoles(v map[string][]string) {
	o.ApplicationRoles = &v
}

// GetSocialLinks returns the SocialLinks field value if set, zero value otherwise.
// Deprecated
func (o *UserRepresentation) GetSocialLinks() []SocialLinkRepresentation {
	if o == nil || IsNil(o.SocialLinks) {
		var ret []SocialLinkRepresentation
		return ret
	}
	return o.SocialLinks
}

// GetSocialLinksOk returns a tuple with the SocialLinks field value if set, nil otherwise
// and a boolean to check if the value has been set.
// Deprecated
func (o *UserRepresentation) GetSocialLinksOk() ([]SocialLinkRepresentation, bool) {
	if o == nil || IsNil(o.SocialLinks) {
		return nil, false
	}
	return o.SocialLinks, true
}

// HasSocialLinks returns a boolean if a field has been set.
func (o *UserRepresentation) HasSocialLinks() bool {
	if o != nil && !IsNil(o.SocialLinks) {
		return true
	}

	return false
}

// SetSocialLinks gets a reference to the given []SocialLinkRepresentation and assigns it to the SocialLinks field.
// Deprecated
func (o *UserRepresentation) SetSocialLinks(v []SocialLinkRepresentation) {
	o.SocialLinks = v
}

// GetGroups returns the Groups field value if set, zero value otherwise.
func (o *UserRepresentation) GetGroups() []string {
	if o == nil || IsNil(o.Groups) {
		var ret []string
		return ret
	}
	return o.Groups
}

// GetGroupsOk returns a tuple with the Groups field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetGroupsOk() ([]string, bool) {
	if o == nil || IsNil(o.Groups) {
		return nil, false
	}
	return o.Groups, true
}

// HasGroups returns a boolean if a field has been set.
func (o *UserRepresentation) HasGroups() bool {
	if o != nil && !IsNil(o.Groups) {
		return true
	}

	return false
}

// SetGroups gets a reference to the given []string and assigns it to the Groups field.
func (o *UserRepresentation) SetGroups(v []string) {
	o.Groups = v
}

// GetAccess returns the Access field value if set, zero value otherwise.
func (o *UserRepresentation) GetAccess() map[string]bool {
	if o == nil || IsNil(o.Access) {
		var ret map[string]bool
		return ret
	}
	return *o.Access
}

// GetAccessOk returns a tuple with the Access field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRepresentation) GetAccessOk() (*map[string]bool, bool) {
	if o == nil || IsNil(o.Access) {
		return nil, false
	}
	return o.Access, true
}

// HasAccess returns a boolean if a field has been set.
func (o *UserRepresentation) HasAccess() bool {
	if o != nil && !IsNil(o.Access) {
		return true
	}

	return false
}

// SetAccess gets a reference to the given map[string]bool and assigns it to the Access field.
func (o *UserRepresentation) SetAccess(v map[string]bool) {
	o.Access = &v
}

func (o UserRepresentation) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UserRepresentation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Username) {
		toSerialize["username"] = o.Username
	}
	if !IsNil(o.FirstName) {
		toSerialize["firstName"] = o.FirstName
	}
	if !IsNil(o.LastName) {
		toSerialize["lastName"] = o.LastName
	}
	if !IsNil(o.Email) {
		toSerialize["email"] = o.Email
	}
	if !IsNil(o.EmailVerified) {
		toSerialize["emailVerified"] = o.EmailVerified
	}
	if !IsNil(o.Attributes) {
		toSerialize["attributes"] = o.Attributes
	}
	if !IsNil(o.UserProfileMetadata) {
		toSerialize["userProfileMetadata"] = o.UserProfileMetadata
	}
	if !IsNil(o.Self) {
		toSerialize["self"] = o.Self
	}
	if !IsNil(o.Origin) {
		toSerialize["origin"] = o.Origin
	}
	if !IsNil(o.CreatedTimestamp) {
		toSerialize["createdTimestamp"] = o.CreatedTimestamp
	}
	if !IsNil(o.Enabled) {
		toSerialize["enabled"] = o.Enabled
	}
	if !IsNil(o.Totp) {
		toSerialize["totp"] = o.Totp
	}
	if !IsNil(o.FederationLink) {
		toSerialize["federationLink"] = o.FederationLink
	}
	if !IsNil(o.ServiceAccountClientId) {
		toSerialize["serviceAccountClientId"] = o.ServiceAccountClientId
	}
	if !IsNil(o.Credentials) {
		toSerialize["credentials"] = o.Credentials
	}
	if !IsNil(o.DisableableCredentialTypes) {
		toSerialize["disableableCredentialTypes"] = o.DisableableCredentialTypes
	}
	if !IsNil(o.RequiredActions) {
		toSerialize["requiredActions"] = o.RequiredActions
	}
	if !IsNil(o.FederatedIdentities) {
		toSerialize["federatedIdentities"] = o.FederatedIdentities
	}
	if !IsNil(o.RealmRoles) {
		toSerialize["realmRoles"] = o.RealmRoles
	}
	if !IsNil(o.ClientRoles) {
		toSerialize["clientRoles"] = o.ClientRoles
	}
	if !IsNil(o.ClientConsents) {
		toSerialize["clientConsents"] = o.ClientConsents
	}
	if !IsNil(o.NotBefore) {
		toSerialize["notBefore"] = o.NotBefore
	}
	if !IsNil(o.ApplicationRoles) {
		toSerialize["applicationRoles"] = o.ApplicationRoles
	}
	if !IsNil(o.SocialLinks) {
		toSerialize["socialLinks"] = o.SocialLinks
	}
	if !IsNil(o.Groups) {
		toSerialize["groups"] = o.Groups
	}
	if !IsNil(o.Access) {
		toSerialize["access"] = o.Access
	}
	return toSerialize, nil
}

type NullableUserRepresentation struct {
	value *UserRepresentation
	isSet bool
}

func (v NullableUserRepresentation) Get() *UserRepresentation {
	return v.value
}

func (v *NullableUserRepresentation) Set(val *UserRepresentation) {
	v.value = val
	v.isSet = true
}

func (v NullableUserRepresentation) IsSet() bool {
	return v.isSet
}

func (v *NullableUserRepresentation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserRepresentation(val *UserRepresentation) *NullableUserRepresentation {
	return &NullableUserRepresentation{value: val, isSet: true}
}

func (v NullableUserRepresentation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserRepresentation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
