// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.
//
// Keycloak Admin REST API
//
// This is a REST API reference for the Keycloak Admin REST API.
//
// API version: 1.0

package idp

// Delete the client scope
// (operationId: AdminRealmsRealmClientScopesClientScopeIdDelete)
type AdminRealmsRealmClientScopesClientScopeIdDeleteParam struct {
	// realm name (not id!)
	Realm         string `uri:"realm"`
	ClientScopeId string `uri:"client-scope-id"`
}

// Get representation of the client scope
// (operationId: AdminRealmsRealmClientScopesClientScopeIdGet)
type AdminRealmsRealmClientScopesClientScopeIdGetParam struct {
	// realm name (not id!)
	Realm         string `uri:"realm"`
	ClientScopeId string `uri:"client-scope-id"`
}

// Update the client scope
// (operationId: AdminRealmsRealmClientScopesClientScopeIdPut)
type AdminRealmsRealmClientScopesClientScopeIdPutParam struct {
	// realm name (not id!)
	Realm         string `uri:"realm"`
	ClientScopeId string `uri:"client-scope-id"`
}

// Get client scopes belonging to the realm Returns a list of client scopes belonging to the realm
// (operationId: AdminRealmsRealmClientScopesGet)
type AdminRealmsRealmClientScopesGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
}

// Create a new client scope Client Scope’s name must be unique!
// (operationId: AdminRealmsRealmClientScopesPost)
type AdminRealmsRealmClientScopesPostParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
}

// Delete the client scope
// (operationId: AdminRealmsRealmClientTemplatesClientScopeIdDelete)
type AdminRealmsRealmClientTemplatesClientScopeIdDeleteParam struct {
	// realm name (not id!)
	Realm         string `uri:"realm"`
	ClientScopeId string `uri:"client-scope-id"`
}

// Get representation of the client scope
// (operationId: AdminRealmsRealmClientTemplatesClientScopeIdGet)
type AdminRealmsRealmClientTemplatesClientScopeIdGetParam struct {
	// realm name (not id!)
	Realm         string `uri:"realm"`
	ClientScopeId string `uri:"client-scope-id"`
}

// Update the client scope
// (operationId: AdminRealmsRealmClientTemplatesClientScopeIdPut)
type AdminRealmsRealmClientTemplatesClientScopeIdPutParam struct {
	// realm name (not id!)
	Realm         string `uri:"realm"`
	ClientScopeId string `uri:"client-scope-id"`
}

// Get client scopes belonging to the realm Returns a list of client scopes belonging to the realm
// (operationId: AdminRealmsRealmClientTemplatesGet)
type AdminRealmsRealmClientTemplatesGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
}

// Create a new client scope Client Scope’s name must be unique!
// (operationId: AdminRealmsRealmClientTemplatesPost)
type AdminRealmsRealmClientTemplatesPostParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
}
