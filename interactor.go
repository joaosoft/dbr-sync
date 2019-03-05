package session

import (
	"auth-types/wst"
	"strings"
	"time"

	"github.com/joaosoft/auth-types/jwt"
	"github.com/satori/go.uuid"
)

type IStorageDB interface {
	GetUserByEmailAndPassword(email, password string) (*User, error)
	GetUserByIdUserAndRefreshToken(idUser, refreshToken string) (*User, error)
	UpdateUserRefreshToken(idUser, refreshToken string) error
}

type Interactor struct {
	config  *SessionConfig
	storage IStorageDB
}

func NewInteractor(config *SessionConfig, storageDB IStorageDB) *Interactor {
	return &Interactor{
		config:  config,
		storage: storageDB,
	}
}

func (i *Interactor) newToken(user *User) (string, error) {
	expirateAt := time.Now().Add(time.Minute * time.Duration(i.config.ExpirationMinutes)).Unix()

	claims := wst.Claims{
		wst.ClaimsExpireAtKey: expirateAt,
		wst.ClaimsAudienceKey: "session",
		wst.ClaimsSubjectKey:  "get-token",
		claimsIdUser:          user.IdUser,
	}
	return wst.New(wst.SignatureHS384, wst.EncodeAscii85, wst.EncodeBase64, wst.EncodeHexadecimal).Generate(claims, i.config.TokenKey)
}

func (i *Interactor) newRefreshToken(user *User) (string, error) {
	jwtId, _ := uuid.NewV4()

	claims := wst.Claims{
		wst.ClaimsAudienceKey: "session",
		wst.ClaimsSubjectKey:  "refresh-token",
		claimsIdUser:          user.IdUser,
		wst.CLaimsJwtId:       jwtId,
	}

	return wst.New(wst.SignatureHS384, wst.EncodeAscii85, wst.EncodeBase64, wst.EncodeHexadecimal).Generate(claims, i.config.TokenKey)
}

func (i *Interactor) GetSession(request *GetSessionRequest) (*SessionResponse, error) {
	log.WithFields(map[string]interface{}{"method": "GetSession"})
	log.Infof("getting user session [email: %s]", request.Email)
	user, err := i.storage.GetUserByEmailAndPassword(request.Email, request.Password)
	if err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error getting user session [email: %s] %s", request.Email, err).ToError()
		return nil, err
	}

	// token
	token, err := i.newToken(user)
	if err != nil {
		return nil, err
	}

	// refresh token
	refreshToken, err := i.newRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// set user refresh token
	if err := i.storage.UpdateUserRefreshToken(user.IdUser, refreshToken); err != nil {
		return nil, err
	}

	return &SessionResponse{
		TokenType:    tokenTypeBearer,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (i *Interactor) loadUserFromRefreshToken(request *RefreshSessionRequest) (*User, error) {
	tokenString := strings.Replace(request.Authorization, "Bearer ", "", 1)

	keyFunc := func(*wst.Token) (interface{}, error) {
		return i.config.TokenKey, nil
	}

	checkFunc := func(wst.Claims) (bool, error) {
		// validate the jti
		return true, nil
	}

	claims := wst.Claims{}
	ok, err := wst.New(wst.SignatureHS384, wst.EncodeAscii85, wst.EncodeBase64, wst.EncodeHexadecimal).Check(tokenString, keyFunc, checkFunc, claims, true)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, wst.ErrorInvalidAuthorization
	}

	if idUser, ok := claims[claimsIdUser]; ok {
		user, err := i.storage.GetUserByIdUserAndRefreshToken(idUser.(string), tokenString)
		return user, err
	}

	return nil, wst.ErrorInvalidAuthorization
}

func (i *Interactor) RefreshToken(request *RefreshSessionRequest) (*SessionResponse, error) {
	log.WithFields(map[string]interface{}{"method": "RefreshToken"})

	// load refresh token
	user, err := i.loadUserFromRefreshToken(request)
	if err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error loading refresh token")
		return nil, wst.ErrorInvalidAuthorization
	}

	if user == nil {
		log.Error("error loading user")
		return nil, jwt.ErrorInvalidAuthorization
	}

	log.Infof("refreshing user session [email: %s]", user.Email)

	// token
	newToken, err := i.newToken(user)
	if err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error generating token")
		return nil, wst.ErrorInvalidAuthorization
	}

	// refresh token
	newRefreshToken, err := i.newRefreshToken(user)
	if err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error generating refresh token")
		return nil, wst.ErrorInvalidAuthorization
	}

	if err := i.storage.UpdateUserRefreshToken(user.IdUser, newRefreshToken); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error updating refresh token of user %s on storage database %s", user.IdUser, err).ToError()
		return nil, err
	}

	return &SessionResponse{
		TokenType:    tokenTypeBearer,
		Token:        newToken,
		RefreshToken: newRefreshToken,
	}, nil
}
