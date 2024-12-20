package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/supermicah/dionysus-admin/pkg/logging"
	"github.com/supermicah/dionysus-admin/pkg/util"
)

type AuthConfig struct {
	AllowedPathPrefixes []string
	SkippedPathPrefixes []string
	RootID              int64
	Skipper             func(c *gin.Context) bool
	ParseUserID         func(c *gin.Context) (int64, error)
}

func AuthWithConfig(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !AllowedPathPrefixes(c, config.AllowedPathPrefixes...) ||
			SkippedPathPrefixes(c, config.SkippedPathPrefixes...) ||
			(config.Skipper != nil && config.Skipper(c)) {
			c.Next()
			return
		}

		userID, err := config.ParseUserID(c)
		if err != nil {
			util.ResError(c, err)
			return
		}

		ctx := util.NewUserID(c.Request.Context(), userID)
		ctx = logging.NewUserID(ctx, userID)
		if userID == config.RootID {
			ctx = util.NewIsRootUser(ctx)
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
