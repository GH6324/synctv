package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/synctv-org/synctv/internal/bootstrap"
	"github.com/synctv-org/synctv/internal/email"
	"github.com/synctv-org/synctv/internal/settings"
	"github.com/synctv-org/synctv/server/model"
)

type publicSettings struct {
	PasswordDisableSignup bool `json:"passwordDisableSignup"`

	EmailEnable           bool     `json:"emailEnable"`
	EmailDisableSignup    bool     `json:"emailDisableSignup"`
	EmailWhitelistEnabled bool     `json:"emailWhitelistEnabled"`
	EmailWhitelist        []string `json:"emailWhitelist,omitempty"`

	Oauth2DisableSignup bool `json:"oauth2DisableSignup"`

	GuestEnable bool `json:"guestEnable"`
}

func Settings(ctx *gin.Context) {
	log := ctx.MustGet("log").(*log.Entry)

	oauth2SignupEnabled, err := bootstrap.Oauth2SignupEnabledCache.Get(ctx)
	if err != nil {
		log.Errorf("failed to get oauth2 signup enabled: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, model.NewApiErrorResp(err))
		return
	}
	ctx.JSON(200, model.NewApiDataResp(
		&publicSettings{
			PasswordDisableSignup: settings.DisableUserSignup.Get() || !settings.EnablePasswordSignup.Get(),

			EmailEnable:           email.EnableEmail.Get(),
			EmailDisableSignup:    settings.DisableUserSignup.Get() || email.DisableUserSignup.Get(),
			EmailWhitelistEnabled: email.EmailSignupWhiteListEnable.Get(),
			EmailWhitelist:        strings.Split(email.EmailSignupWhiteList.Get(), ","),

			Oauth2DisableSignup: settings.DisableUserSignup.Get() || len(oauth2SignupEnabled) == 0,

			GuestEnable: settings.EnableGuest.Get(),
		},
	))
}
