// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.
//
// Keycloak Admin REST API
//
// This is a REST API reference for the Keycloak Admin REST API.
//
// API version: 1.0

package idp

// Get the client secret
// (operationId: AdminRealmsRealmClientsClientUuidClientSecretGet)
type AdminRealmsRealmClientsClientUuidClientSecretGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// Generate a new secret for the client
// (operationId: AdminRealmsRealmClientsClientUuidClientSecretPost)
type AdminRealmsRealmClientsClientUuidClientSecretPostParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// Invalidate the rotated secret for the client
// (operationId: AdminRealmsRealmClientsClientUuidClientSecretRotatedDelete)
type AdminRealmsRealmClientsClientUuidClientSecretRotatedDeleteParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// Get the rotated client secret
// (operationId: AdminRealmsRealmClientsClientUuidClientSecretRotatedGet)
type AdminRealmsRealmClientsClientUuidClientSecretRotatedGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// (operationId: AdminRealmsRealmClientsClientUuidDefaultClientScopesClientScopeIdDelete)
type AdminRealmsRealmClientsClientUuidDefaultClientScopesClientScopeIdDeleteParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid    string `uri:"client-uuid"`
	ClientScopeId string `uri:"clientScopeId"`
}

// (operationId: AdminRealmsRealmClientsClientUuidDefaultClientScopesClientScopeIdPut)
type AdminRealmsRealmClientsClientUuidDefaultClientScopesClientScopeIdPutParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid    string `uri:"client-uuid"`
	ClientScopeId string `uri:"clientScopeId"`
}

// Get default client scopes.  Only name and ids are returned.
// (operationId: AdminRealmsRealmClientsClientUuidDefaultClientScopesGet)
type AdminRealmsRealmClientsClientUuidDefaultClientScopesGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// Delete the client
// (operationId: AdminRealmsRealmClientsClientUuidDelete)
type AdminRealmsRealmClientsClientUuidDeleteParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// Create JSON with payload of example access token
// (operationId: AdminRealmsRealmClientsClientUuidEvaluateScopesGenerateExampleAccessTokenGet)
type AdminRealmsRealmClientsClientUuidEvaluateScopesGenerateExampleAccessTokenGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
	Scope      string `form:"scope,omitempty"`
	UserId     string `form:"userId,omitempty"`
}

// Create JSON with payload of example id token
// (operationId: AdminRealmsRealmClientsClientUuidEvaluateScopesGenerateExampleIdTokenGet)
type AdminRealmsRealmClientsClientUuidEvaluateScopesGenerateExampleIdTokenGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
	Scope      string `form:"scope,omitempty"`
	UserId     string `form:"userId,omitempty"`
}

// Create JSON with payload of example user info
// (operationId: AdminRealmsRealmClientsClientUuidEvaluateScopesGenerateExampleUserinfoGet)
type AdminRealmsRealmClientsClientUuidEvaluateScopesGenerateExampleUserinfoGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
	Scope      string `form:"scope,omitempty"`
	UserId     string `form:"userId,omitempty"`
}

// Return list of all protocol mappers, which will be used when generating tokens issued for particular client.
// This means protocol mappers assigned to this client directly and protocol mappers assigned to all client scopes of this client. (operationId: AdminRealmsRealmClientsClientUuidEvaluateScopesProtocolMappersGet)
type AdminRealmsRealmClientsClientUuidEvaluateScopesProtocolMappersGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
	Scope      string `form:"scope,omitempty"`
}

// Get effective scope mapping of all roles of particular role container, which this client is defacto allowed to have in the accessToken issued for him.
// This contains scope mappings, which this client has directly, as well as scope mappings, which are granted to all client scopes, which are linked with this client. (operationId: AdminRealmsRealmClientsClientUuidEvaluateScopesScopeMappingsRoleContainerIdGrantedGet)
type AdminRealmsRealmClientsClientUuidEvaluateScopesScopeMappingsRoleContainerIdGrantedGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
	// either realm name OR client UUID
	RoleContainerId string `uri:"roleContainerId"`
	Scope           string `form:"scope,omitempty"`
}

// Get roles, which this client doesn't have scope for and can't have them in the accessToken issued for him.
// Defacto all the other roles of particular role container, which are not in {@link #getGrantedScopeMappings()} (operationId: AdminRealmsRealmClientsClientUuidEvaluateScopesScopeMappingsRoleContainerIdNotGrantedGet)
type AdminRealmsRealmClientsClientUuidEvaluateScopesScopeMappingsRoleContainerIdNotGrantedGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
	// either realm name OR client UUID
	RoleContainerId string `uri:"roleContainerId"`
	Scope           string `form:"scope,omitempty"`
}

// Get representation of the client
// (operationId: AdminRealmsRealmClientsClientUuidGet)
type AdminRealmsRealmClientsClientUuidGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// (operationId: AdminRealmsRealmClientsClientUuidInstallationProvidersProviderIdGet)
type AdminRealmsRealmClientsClientUuidInstallationProvidersProviderIdGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
	ProviderId string `uri:"providerId"`
}

// Return object stating whether client Authorization permissions have been initialized or not and a reference
// (operationId: AdminRealmsRealmClientsClientUuidManagementPermissionsGet)
type AdminRealmsRealmClientsClientUuidManagementPermissionsGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// Return object stating whether client Authorization permissions have been initialized or not and a reference
// (operationId: AdminRealmsRealmClientsClientUuidManagementPermissionsPut)
type AdminRealmsRealmClientsClientUuidManagementPermissionsPutParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// Unregister a cluster node from the client
// (operationId: AdminRealmsRealmClientsClientUuidNodesNodeDelete)
type AdminRealmsRealmClientsClientUuidNodesNodeDeleteParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
	Node       string `uri:"node"`
}

// Register a cluster node with the client Manually register cluster node to this client - usually it’s not needed to call this directly as adapter should handle by sending registration request to Keycloak
// (operationId: AdminRealmsRealmClientsClientUuidNodesPost)
type AdminRealmsRealmClientsClientUuidNodesPostParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// Get application offline session count Returns a number of offline user sessions associated with this client { \"count\": number }
// (operationId: AdminRealmsRealmClientsClientUuidOfflineSessionCountGet)
type AdminRealmsRealmClientsClientUuidOfflineSessionCountGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// Get offline sessions for client Returns a list of offline user sessions associated with this client
// (operationId: AdminRealmsRealmClientsClientUuidOfflineSessionsGet)
type AdminRealmsRealmClientsClientUuidOfflineSessionsGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
	// Paging offset
	First int32 `form:"first,omitempty"`
	// Maximum results size (defaults to 100)
	Max int32 `form:"max,omitempty"`
}

// (operationId: AdminRealmsRealmClientsClientUuidOptionalClientScopesClientScopeIdDelete)
type AdminRealmsRealmClientsClientUuidOptionalClientScopesClientScopeIdDeleteParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid    string `uri:"client-uuid"`
	ClientScopeId string `uri:"clientScopeId"`
}

// (operationId: AdminRealmsRealmClientsClientUuidOptionalClientScopesClientScopeIdPut)
type AdminRealmsRealmClientsClientUuidOptionalClientScopesClientScopeIdPutParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid    string `uri:"client-uuid"`
	ClientScopeId string `uri:"clientScopeId"`
}

// Get optional client scopes.  Only name and ids are returned.
// (operationId: AdminRealmsRealmClientsClientUuidOptionalClientScopesGet)
type AdminRealmsRealmClientsClientUuidOptionalClientScopesGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// Push the client's revocation policy to its admin URL If the client has an admin URL, push revocation policy to it.
// (operationId: AdminRealmsRealmClientsClientUuidPushRevocationPost)
type AdminRealmsRealmClientsClientUuidPushRevocationPostParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// Update the client
// (operationId: AdminRealmsRealmClientsClientUuidPut)
type AdminRealmsRealmClientsClientUuidPutParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// Generate a new registration access token for the client
// (operationId: AdminRealmsRealmClientsClientUuidRegistrationAccessTokenPost)
type AdminRealmsRealmClientsClientUuidRegistrationAccessTokenPostParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// Get a user dedicated to the service account
// (operationId: AdminRealmsRealmClientsClientUuidServiceAccountUserGet)
type AdminRealmsRealmClientsClientUuidServiceAccountUserGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// Get application session count Returns a number of user sessions associated with this client { \"count\": number }
// (operationId: AdminRealmsRealmClientsClientUuidSessionCountGet)
type AdminRealmsRealmClientsClientUuidSessionCountGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// Test if registered cluster nodes are available Tests availability by sending 'ping' request to all cluster nodes.
// (operationId: AdminRealmsRealmClientsClientUuidTestNodesAvailableGet)
type AdminRealmsRealmClientsClientUuidTestNodesAvailableGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
}

// Get user sessions for client Returns a list of user sessions associated with this client
// (operationId: AdminRealmsRealmClientsClientUuidUserSessionsGet)
type AdminRealmsRealmClientsClientUuidUserSessionsGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// id of client (not client-id!)
	ClientUuid string `uri:"client-uuid"`
	// Paging offset
	First int32 `form:"first,omitempty"`
	// Maximum results size (defaults to 100)
	Max int32 `form:"max,omitempty"`
}

// Get clients belonging to the realm.
// If a client can’t be retrieved from the storage due to a problem with the underlying storage, it is silently removed from the returned list. This ensures that concurrent modifications to the list don’t prevent callers from retrieving this list. (operationId: AdminRealmsRealmClientsGet)
type AdminRealmsRealmClientsGetParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
	// filter by clientId
	ClientId string `form:"clientId,omitempty"`
	// the first result
	First int32 `form:"first,omitempty"`
	// the max results to return
	Max int32  `form:"max,omitempty"`
	Q   string `form:"q,omitempty"`
	// whether this is a search query or a getClientById query
	Search bool `form:"search,omitempty"`
	// filter clients that cannot be viewed in full by admin
	ViewableOnly bool `form:"viewableOnly,omitempty"`
}

// Create a new client Client’s client_id must be unique!
// (operationId: AdminRealmsRealmClientsPost)
type AdminRealmsRealmClientsPostParam struct {
	// realm name (not id!)
	Realm string `uri:"realm"`
}
