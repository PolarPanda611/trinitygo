package errutil

import "errors"

var (

	//Real Error

	// ErrTokenExpired return token expired err
	ErrTokenExpired = errors.New("app.err.TokenExpired")
	// ErrTokenWrongIssuer return wrong issuer err
	ErrTokenWrongIssuer = errors.New("app.err.TokenWrongIssuer")
	// ErrTokenWrongHeaderPrefix return wrong header prefix
	ErrTokenWrongHeaderPrefix = errors.New("app.err.TokenWrongHeaderPrefix")
	// ErrTokenWrongAuthorization return wrong authorization
	ErrTokenWrongAuthorization = errors.New("app.err.TokenWrongAuthorization")

	//User Error

	// ErrUnverifiedToken unverified token
	ErrUnverifiedToken = errors.New("app.err.UnverifiedToken")
	// ErrGeneratedToken GenerateTokenFailed
	ErrGeneratedToken = errors.New("app.err.GenerateTokenFailed")
	// ErrGetUserAuth wrong get user auth func
	ErrGetUserAuth = errors.New("app.err.WrongGetUserAuthFunc")
	// ErrAccessAuthCheckFailed fail to pass access auth check
	ErrAccessAuthCheckFailed = errors.New("app.error.AccessAuthorizationCheckFailed")
	// ErrLoadDataFailed app.error.LoadDataFailed
	ErrLoadDataFailed = errors.New("app.error.LoadDataFailed")
	// ErrResolveDataFailed app.error.ResolveDataFailed
	ErrResolveDataFailed = errors.New("app.error.ResolveDataFailed")
	// ErrCreateDataFailed app.error.CreateDataFailed
	ErrCreateDataFailed = errors.New("app.error.CreateDataFailed")
	// ErrUpdateDataFailed app.error.UpdateDataFailed
	ErrUpdateDataFailed = errors.New("app.error.UpdateDataFailed")
	// ErrDeleteDataFailed app.error.DeleteDataFailed
	ErrDeleteDataFailed = errors.New("app.error.DeleteDataFailed")
	// ErrUnknownService  app.error.UnknownService
	ErrUnknownService = errors.New("app.error.UnknownService")
	// ErrDataAssertion  app.error.DataAssetion
	ErrDataAssertion = errors.New("app.error.DataAssetion")
	// ErrUpdateZeroAffected  app.error.DataAssetion
	ErrUpdateZeroAffected = errors.New("app.error.ErrUpdateZeroAffected")
)
