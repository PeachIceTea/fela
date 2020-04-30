package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/PeachIceTea/fela/conf"
)

// Define errors
var (
	ErrIDParamMissing = errors.New("id parameter is missing")
	ErrInvalidID      = errors.New("invalid id")
)

// RegisterRoutes registers /v1 routes
func RegisterRoutes(r *gin.RouterGroup, c *conf.Config) {
	v1 := r.Group("/api/v1")
	{
		Login(v1, c)
	}

	protected := v1.Group("/")
	protected.Use(authHeaderRequired(c))
	{

		GetUsers(protected, c)
		GetUser(protected, c)
		UpdateUser(protected, c)
		Register(v1, c) //TODO: Proper setup story
		DeleteUser(protected, c)

		GetAudiobooks(protected, c)
		GetAudiobook(protected, c)
		GetAudiobookFiles(protected, c)
		Upload(protected, c)
		UpdateAudiobook(protected, c)
		DeleteAudiobook(protected, c)

		protected.GET("/token-test", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, conf.M{"msg": "token is valid"})
		})
	}

	ServeFiles(r, c)
}

// The authHeaderRequired guard requires a request to carry a valid JWT in the
// Authorization header using the Bearer schema.
func authHeaderRequired(c *conf.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authStr := ctx.GetHeader("Authorization")
		if authStr == "" {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				conf.M{"err": "auth header missing"},
			)
			return
		}

		header := strings.Split(authStr, "Bearer ")
		if len(header) != 2 {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				conf.M{"err": "invalid auth header"},
			)
			return
		}

		claims, err := parseToken(header[1], c)
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				conf.M{"err": "could not parse auth token"},
			)
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}

// parseToken parses the JWT token and returns its claims.
func parseToken(token string, c *conf.Config) (claims *Claims, err error) {
	claims = &Claims{}
	_, err = jwt.ParseWithClaims(
		token,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(
					"Unexpected signing method: %v",
					token.Header["alg"],
				)
			}

			return c.Secret, nil
		},
	)
	return
}

// getID extracts :id from URL.
func getID(ctx *gin.Context) (int64, error) {
	idStr, ok := ctx.Params.Get("id")
	if !ok {
		return 0, ErrIDParamMissing
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, ErrInvalidID
	}

	return id, err
}

// Claims for the JWT token
// TODO: Refresh token
type Claims struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.StandardClaims
}

// isAdmin returns whether user is an admin
func (c *Claims) isAdmin() bool {
	return c.Role == "admin"
}

// isUploader returns whether user is allowed to upload
func (c *Claims) isUploader() bool {
	return c.Role == "uploader" || c.isAdmin()
}

// getClaims extracts a claims struct from gin.Context.
// This will panic when used in unprotected routes.
func getClaims(ctx *gin.Context) (c *Claims) {
	claims, ok := ctx.Get("claims")
	if !ok {
		panic("could not get claims from context")
	}

	return claims.(*Claims)
}
