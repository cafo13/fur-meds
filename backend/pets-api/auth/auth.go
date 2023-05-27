package auth

import (
	"net/http"
	"reflect"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type AuthMiddleware interface {
	Middleware() gin.HandlerFunc
	GetUserUidByMail(ctx *gin.Context, userMail string) (string, error)
}

type FirebaseAuthMiddleware struct {
	AuthClient *auth.Client
}

func NewFirebaseAuthMiddleware(authClient *auth.Client) AuthMiddleware {
	return FirebaseAuthMiddleware{AuthClient: authClient}
}

func (a FirebaseAuthMiddleware) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := a.tokenFromHeader(ctx.Request)
		if bearerToken == "" {
			log.Error(errors.New("empty bearer token"))
			ctx.String(http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}

		token, err := a.AuthClient.VerifyIDToken(ctx, bearerToken)
		if err != nil {
			log.Error(errors.Wrap(err, "unable to verify jwt"))
			ctx.String(http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}

		user := User{
			UID:   token.UID,
			Email: token.Claims["email"].(string),
		}
		ctx.Set("user", user)

		log.Infof("setting user %+v to context value 'user'", user)

		ctx.Next()
	}
}

func (a FirebaseAuthMiddleware) GetUserUidByMail(ctx *gin.Context, userMail string) (string, error) {
	user, err := a.AuthClient.GetUserByEmail(ctx, userMail)
	if err != nil {
		return "", err
	}

	return user.UID, nil
}

func (a FirebaseAuthMiddleware) tokenFromHeader(r *http.Request) string {
	headerValue := r.Header.Get("Authorization")

	if len(headerValue) > 7 && strings.ToLower(headerValue[0:6]) == "bearer" {
		return headerValue[7:]
	}

	return ""
}

type User struct {
	UID   string
	Email string
}

type ctxKey int

const (
	userContextKey ctxKey = iota
)

var (
	ErrNoUserInContext = errors.New("auth error: no user in context")
)

func UserFromCtx(ctx *gin.Context) (User, error) {
	user, ok := ctx.Get("user")
	log.Infof("got user %+v from context value 'user'", user)

	if ok {
		reflectUser := reflect.ValueOf(user)
		if reflectUser.IsZero() {
			return User{}, errors.New("auth error: user in context was zero value")
		}
		userUid := reflectUser.FieldByName("UID")
		userMail := reflectUser.FieldByName("Email")

		if userUid.IsZero() {
			return User{}, errors.New("auth error: user UID in context was zero value")
		}

		return User{
			UID:   userUid.String(),
			Email: userMail.String(),
		}, nil
	}

	return User{}, ErrNoUserInContext
}
