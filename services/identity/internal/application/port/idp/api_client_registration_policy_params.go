// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.
//
// Keycloak Admin REST API
//
// This is a REST API reference for the Keycloak Admin REST API.
//
// API version: 1.0

package idp

// Base path for retrieve providers with the configProperties properly filled
// (operationId: AdminRealmsRealmClientRegistrationPolicyProvidersGet)
type AdminRealmsRealmClientRegistrationPolicyProvidersGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
}
