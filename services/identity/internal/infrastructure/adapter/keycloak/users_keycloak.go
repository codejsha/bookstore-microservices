package keycloak

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/client"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/idp"
	"github.com/codejsha/bookstore-microservices/identity/internal/config"
)

var _ idp.UsersAPI = (*usersClient)(nil)

func NewUsersClient(
	cfg *config.Config,
	restyClient *client.RestyClient,
	tokenHelper *AdminTokenHelper,
) idp.UsersAPI {
	return &usersClient{
		cfg:         cfg,
		tokenHelper: tokenHelper,
		restyClient: restyClient.Client,
	}
}

type usersClient struct {
	cfg         *config.Config
	restyClient *resty.Client
	tokenHelper *AdminTokenHelper
}

func (c usersClient) AdminRealmsRealmUsersCountGet(ctx context.Context, param idp.AdminRealmsRealmUsersCountGetParam) (int32, error) {
	// prepare request
	token, _, err := c.tokenHelper.GetTokens()
	if err != nil {
		return 0, fmt.Errorf("failed to get token: %v", err)
	}
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	// execute request
	restyResp, err := c.restyClient.R().
		SetContext(ctx).
		SetDebug(c.cfg.App.Logging.IsDebug).
		SetHeaders(headers).
		Get(c.cfg.Keycloak.Url + fmt.Sprintf("/admin/realms/%s/users/count", param.Realm))
	if err != nil {
		return 0, fmt.Errorf("failed to get user count: %v", err)
	}
	if restyResp.StatusCode() != 200 {
		return 0, fmt.Errorf("failed to get user count: %s", restyResp.String())
	}

	// return response
	var response int32
	err = json.Unmarshal(restyResp.Body(), &response)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal user count response: %v", err)
	}
	return response, nil

}

func (c usersClient) AdminRealmsRealmUsersGet(ctx context.Context, param idp.AdminRealmsRealmUsersGetParam) ([]idp.UserRepresentation, error) {
	// prepare request
	token, _, err := c.tokenHelper.GetTokens()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %v", err)
	}
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	data := map[string]string{
		"client_id":     c.cfg.Keycloak.ClientId,
		"client_secret": c.cfg.Keycloak.ClientSecret,
	}

	// execute request
	request := c.restyClient.R().
		SetContext(ctx).
		SetDebug(c.cfg.App.Logging.IsDebug).
		SetHeaders(headers).
		SetFormData(data)
	if param.Email != "" {
		request.SetQueryParam("username", param.Email)
	}
	restyResp, err := request.
		Get(c.cfg.Keycloak.Url + fmt.Sprintf("/admin/realms/%s/users", param.Realm))
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}
	if restyResp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get users: %s", restyResp.String())
	}

	// return response
	var response []idp.UserRepresentation
	err = json.Unmarshal(restyResp.Body(), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal users response: %v", err)
	}
	return response, nil
}

func (c usersClient) AdminRealmsRealmUsersPost(ctx context.Context, param idp.AdminRealmsRealmUsersPostParam, req idp.UserRepresentation) error {
	// prepare request
	token, _, err := c.tokenHelper.GetTokens()
	if err != nil {
		return fmt.Errorf("failed to get token: %v", err)
	}
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Content-Type":  "application/json",
	}
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}

	// execute request
	restyResp, err := c.restyClient.R().
		SetContext(ctx).
		SetDebug(c.cfg.App.Logging.IsDebug).
		SetHeaders(headers).
		SetBody(body).
		Post(c.cfg.Keycloak.Url + fmt.Sprintf("/admin/realms/%s/users", param.Realm))
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	if restyResp.StatusCode() != 201 {
		return fmt.Errorf("failed to create user: %s", restyResp.String())
	}

	// return response
	return nil
}

func (c usersClient) AdminRealmsRealmUsersProfileGet(ctx context.Context, param idp.AdminRealmsRealmUsersProfileGetParam) (idp.UPConfig, error) {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersProfileMetadataGet(ctx context.Context, param idp.AdminRealmsRealmUsersProfileMetadataGetParam) (idp.UserProfileMetadata, error) {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersProfilePut(ctx context.Context, param idp.AdminRealmsRealmUsersProfilePutParam, req idp.UPConfig) (idp.UPConfig, error) {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdConfiguredUserStorageCredentialTypesGet(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdConfiguredUserStorageCredentialTypesGetParam) ([]string, error) {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdConsentsClientDelete(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdConsentsClientDeleteParam) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdConsentsGet(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdConsentsGetParam) ([]map[string]interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdCredentialsCredentialIdDelete(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdCredentialsCredentialIdDeleteParam) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdCredentialsCredentialIdMoveAfterNewPreviousCredentialIdPost(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdCredentialsCredentialIdMoveAfterNewPreviousCredentialIdPostParam) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdCredentialsCredentialIdMoveToFirstPost(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdCredentialsCredentialIdMoveToFirstPostParam) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdCredentialsCredentialIdUserLabelPut(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdCredentialsCredentialIdUserLabelPutParam, req idp.AdminRealmsRealmUsersUserIdCredentialsCredentialIdUserLabelPutReq) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdCredentialsGet(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdCredentialsGetParam) ([]idp.CredentialRepresentation, error) {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdDelete(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdDeleteParam) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdDisableCredentialTypesPut(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdDisableCredentialTypesPutParam, req idp.AdminRealmsRealmUsersUserIdDisableCredentialTypesPutReq) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdExecuteActionsEmailPut(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdExecuteActionsEmailPutParam, req idp.AdminRealmsRealmUsersUserIdExecuteActionsEmailPutReq) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdFederatedIdentityGet(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdFederatedIdentityGetParam) ([]idp.FederatedIdentityRepresentation, error) {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdFederatedIdentityProviderDelete(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdFederatedIdentityProviderDeleteParam) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdFederatedIdentityProviderPost(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdFederatedIdentityProviderPostParam) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdGet(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdGetParam) (idp.UserRepresentation, error) {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdGroupsCountGet(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdGroupsCountGetParam) (map[string]int64, error) {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdGroupsGet(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdGroupsGetParam) ([]idp.GroupRepresentation, error) {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdGroupsGroupIdDelete(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdGroupsGroupIdDeleteParam) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdGroupsGroupIdPut(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdGroupsGroupIdPutParam) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdImpersonationPost(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdImpersonationPostParam) (map[string]interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdLogoutPost(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdLogoutPostParam) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdOfflineSessionsClientUuidGet(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdOfflineSessionsClientUuidGetParam) ([]idp.UserSessionRepresentation, error) {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdPut(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdPutParam, req idp.UserRepresentation) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdResetPasswordEmailPut(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdResetPasswordEmailPutParam) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdResetPasswordPut(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdResetPasswordPutParam, req idp.CredentialRepresentation) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdSendVerifyEmailPut(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdSendVerifyEmailPutParam) error {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdSessionsGet(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdSessionsGetParam) ([]idp.UserSessionRepresentation, error) {
	// TODO implement me
	panic("implement me")
}

func (c usersClient) AdminRealmsRealmUsersUserIdUnmanagedAttributesGet(ctx context.Context, param idp.AdminRealmsRealmUsersUserIdUnmanagedAttributesGetParam) (map[string][]string, error) {
	// TODO implement me
	panic("implement me")
}
