package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/errorcode"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/usecase"
)

var _ openapi.IdentityAPI = (*identityController)(nil)

func NewIdentityController(
	identityUseCase usecase.IdentityUseCase,
) openapi.IdentityAPI {
	return &identityController{
		identityUseCase: identityUseCase,
	}
}

type identityController struct {
	identityUseCase usecase.IdentityUseCase
}

func (g identityController) ApiV1UsersSigninPost(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var req openapi.SignInWebReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}

	// perform request
	resp, err := g.identityUseCase.SignIn(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return response
	c.SetCookie("refresh_token", resp.RefreshToken, 1800,
		"/api/v1/users/tokens/exchange", "", false, true)
	c.JSON(http.StatusOK, resp.ToSignInWebResp())
}

func (g identityController) ApiV1UsersSignoutPost(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.SignOutWebParam
	if err := c.ShouldBindHeader(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	err := g.identityUseCase.SignOut(ctx, param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return response
	c.Status(http.StatusNoContent)
}

func (g identityController) ApiV1UsersSignupPost(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var req openapi.SignUpWebReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestPayload, err)
		return
	}

	// perform request
	resp, err := g.identityUseCase.SignUp(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return response
	c.JSON(http.StatusCreated, resp.ToSignUpWebResp())
}

func (g identityController) ApiV1UsersTokensExchangePost(c *gin.Context) {
	// prepare request
	ctx := c.Request.Context()
	var param openapi.TokenExchangeWebParam
	if err := c.ShouldBindHeader(&param); err != nil {
		errorcode.RenderResponse(c, errorcode.InvalidRequestParams, err)
		return
	}

	// perform request
	resp, err := g.identityUseCase.ExchangeToken(ctx, param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return response
	c.SetCookie("refresh_token", resp.RefreshToken, 1800,
		"/api/v1/users/tokens/exchange", "", false, true)
	c.JSON(http.StatusOK, resp.ToTokenExchangeWebResp())
}
