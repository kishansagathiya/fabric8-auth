package controller

import (
	"strings"

	"github.com/fabric8-services/fabric8-auth/account/userinfo"
	"github.com/fabric8-services/fabric8-auth/app"
	"github.com/fabric8-services/fabric8-auth/application"
	"github.com/fabric8-services/fabric8-auth/jsonapi"
	"github.com/fabric8-services/fabric8-auth/token"
	"github.com/goadesign/goa"
)

// UserinfoController implements the userinfo resource.
type UserinfoController struct {
	*goa.Controller
	db             application.DB
	tokenManager   token.Manager
	accountService userinfo.AccountService
}

// NewUserinfoController creates a userinfo controller.
func NewUserinfoController(service *goa.Service, accountService userinfo.AccountService, db application.DB, tokenManager token.Manager) *UserinfoController {
	return &UserinfoController{
		Controller:     service.NewController("UserinfoController"),
		accountService: accountService,
		db:             db,
		tokenManager:   tokenManager,
	}
}

// Show runs the show action.
func (c *UserinfoController) Show(ctx *app.ShowUserinfoContext) error {

	user, identity, err := c.accountService.UserInfo(ctx)
	if err != nil {
		return jsonapi.JSONErrorResponse(ctx, err)
	}

	fullName := strings.Split(user.FullName, " ")
	sub := identity.ID.String()
	userInfo := &app.UserInfo{
		Sub:               &sub,
		GivenName:         &fullName[0],
		PreferredUsername: &identity.Username,
		FamilyName:        &fullName[1],
		Email:             &user.Email,
	}

	return ctx.OK(userInfo)
}