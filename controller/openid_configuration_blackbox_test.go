package controller_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/fabric8-services/fabric8-auth/app"
	"github.com/fabric8-services/fabric8-auth/app/test"
	"github.com/fabric8-services/fabric8-auth/client"
	config "github.com/fabric8-services/fabric8-auth/configuration"
	. "github.com/fabric8-services/fabric8-auth/controller"
	"github.com/fabric8-services/fabric8-auth/resource"
	"github.com/fabric8-services/fabric8-auth/rest"

	"github.com/goadesign/goa"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestOpenIDConfigurationREST struct {
	suite.Suite
	configuration *config.ConfigurationData
}

func TestRunOpenIDConfigurationREST(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	configuration, err := config.GetConfigurationData()
	if err != nil {
		log.Panic(nil, map[string]interface{}{
			"err": err,
		}, "failed to setup the configuration")
	}
	suite.Run(t, &TestOpenIDConfigurationREST{configuration: configuration})
}

func (s *TestOpenIDConfigurationREST) UnSecuredController() (*goa.Service, *OpenidConfigurationController) {
	svc := goa.New("Status-Service")
	return svc, NewOpenidConfigurationController(svc, s.configuration)
}

func (s *TestOpenIDConfigurationREST) TestShowOpenIDConfiguration() {
	t := s.T()
	svc, ctrl := s.UnSecuredController()

	_, openIDConfiguration := test.ShowOpenidConfigurationOK(t, svc.Context, svc, ctrl)

	u := &url.URL{
		Path: fmt.Sprintf(client.ShowOpenidConfigurationPath()),
	}
	prms := url.Values{}
	req, err := http.NewRequest("GET", u.String(), nil)

	ctx := context.Background()
	rw := httptest.NewRecorder()
	goaCtx := goa.NewContext(goa.WithAction(ctx, "OpenIDConfigurationTest"), rw, req, prms)
	openIDConfigurationCtx, err := app.NewShowOpenidConfigurationContext(goaCtx, req, goa.New("LoginService"))
	require.Nil(t, err)

	issuer := rest.AbsoluteURL(openIDConfigurationCtx.RequestData, "")
	authorizationEndpoint := rest.AbsoluteURL(openIDConfigurationCtx.RequestData, client.AuthorizeAuthorizePath())
	tokenEndpoint := rest.AbsoluteURL(openIDConfigurationCtx.RequestData, client.ExchangeTokenPath())
	logoutEndpoint := rest.AbsoluteURL(openIDConfigurationCtx.RequestData, client.LogoutLogoutPath())
	jwksURI := s.configuration.GetKeycloakEndpointCerts()

	expectedOpenIDConfiguration := &app.OpenIDConfiguration{
		Issuer:                           &issuer,
		AuthorizationEndpoint:            &authorizationEndpoint,
		TokenEndpoint:                    &tokenEndpoint,
		EndSessionEndpoint:               &logoutEndpoint,
		ResponseTypesSupported:           []string{"code"},
		JwksURI:                          &jwksURI,
		GrantTypesSupported:              []string{"authorization_code", "refresh_token", "client_credentials"},
		SubjectTypesSupported:            []string{},
		IDTokenSigningAlgValuesSupported: []string{"RS256"},
	}

	require.Equal(t, openIDConfiguration, expectedOpenIDConfiguration)
}