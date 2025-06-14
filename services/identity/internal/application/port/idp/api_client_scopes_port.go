// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.
//
// Keycloak Admin REST API
//
// This is a REST API reference for the Keycloak Admin REST API.
//
// API version: 1.0

package idp

import (
	"context"
)

type ClientScopesAPI interface {

	// Delete the client scope
	// (operationId: AdminRealmsRealmClientScopesClientScopeIdDelete)
	// DELETE /admin/realms/{realm}/client-scopes/{client-scope-id}
	AdminRealmsRealmClientScopesClientScopeIdDelete(ctx context.Context, param AdminRealmsRealmClientScopesClientScopeIdDeleteParam) error

	// Get representation of the client scope
	// (operationId: AdminRealmsRealmClientScopesClientScopeIdGet)
	// GET /admin/realms/{realm}/client-scopes/{client-scope-id}
	AdminRealmsRealmClientScopesClientScopeIdGet(ctx context.Context, param AdminRealmsRealmClientScopesClientScopeIdGetParam) (ClientScopeRepresentation, error)

	// Update the client scope
	// (operationId: AdminRealmsRealmClientScopesClientScopeIdPut)
	// PUT /admin/realms/{realm}/client-scopes/{client-scope-id}
	AdminRealmsRealmClientScopesClientScopeIdPut(ctx context.Context, param AdminRealmsRealmClientScopesClientScopeIdPutParam, req ClientScopeRepresentation) error

	// Get client scopes belonging to the realm Returns a list of client scopes belonging to the realm
	// (operationId: AdminRealmsRealmClientScopesGet)
	// GET /admin/realms/{realm}/client-scopes
	AdminRealmsRealmClientScopesGet(ctx context.Context, param AdminRealmsRealmClientScopesGetParam) ([]ClientScopeRepresentation, error)

	// Create a new client scope Client Scope’s name must be unique!
	// (operationId: AdminRealmsRealmClientScopesPost)
	// POST /admin/realms/{realm}/client-scopes
	AdminRealmsRealmClientScopesPost(ctx context.Context, param AdminRealmsRealmClientScopesPostParam, req ClientScopeRepresentation) error

	// Delete the client scope
	// (operationId: AdminRealmsRealmClientTemplatesClientScopeIdDelete)
	// DELETE /admin/realms/{realm}/client-templates/{client-scope-id}
	AdminRealmsRealmClientTemplatesClientScopeIdDelete(ctx context.Context, param AdminRealmsRealmClientTemplatesClientScopeIdDeleteParam) error

	// Get representation of the client scope
	// (operationId: AdminRealmsRealmClientTemplatesClientScopeIdGet)
	// GET /admin/realms/{realm}/client-templates/{client-scope-id}
	AdminRealmsRealmClientTemplatesClientScopeIdGet(ctx context.Context, param AdminRealmsRealmClientTemplatesClientScopeIdGetParam) (ClientScopeRepresentation, error)

	// Update the client scope
	// (operationId: AdminRealmsRealmClientTemplatesClientScopeIdPut)
	// PUT /admin/realms/{realm}/client-templates/{client-scope-id}
	AdminRealmsRealmClientTemplatesClientScopeIdPut(ctx context.Context, param AdminRealmsRealmClientTemplatesClientScopeIdPutParam, req ClientScopeRepresentation) error

	// Get client scopes belonging to the realm Returns a list of client scopes belonging to the realm
	// (operationId: AdminRealmsRealmClientTemplatesGet)
	// GET /admin/realms/{realm}/client-templates
	AdminRealmsRealmClientTemplatesGet(ctx context.Context, param AdminRealmsRealmClientTemplatesGetParam) ([]ClientScopeRepresentation, error)

	// Create a new client scope Client Scope’s name must be unique!
	// (operationId: AdminRealmsRealmClientTemplatesPost)
	// POST /admin/realms/{realm}/client-templates
	AdminRealmsRealmClientTemplatesPost(ctx context.Context, param AdminRealmsRealmClientTemplatesPostParam, req ClientScopeRepresentation) error
}
