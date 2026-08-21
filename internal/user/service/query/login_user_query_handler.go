// Package loginqry provides functionality for handling login queries.
// It processes login requests by verifying user credentials and generating authentication results.
package loginqry

import (
	"errors"
	"fmt"

	ijwt "github.com/mohit/finance-go/internal/auth/jwt_interface"
	appError "github.com/mohit/finance-go/internal/common/app_error"
	iquery "github.com/mohit/finance-go/internal/common/cqrs/query"
	errdmn "github.com/mohit/finance-go/internal/common/error"
	"github.com/mohit/finance-go/internal/common/hash"
	irepository "github.com/mohit/finance-go/internal/common/repository"
	auth "github.com/mohit/finance-go/internal/user/service/common"
)

// Handler processes login queries by validating user credentials and generating authentication results.
type Handler struct {
	userRepo irepository.IUserRepository
	jwtSvc   ijwt.IService
	hashSvc  hash.IService
}

// Ensure Handler implements the iquery.IHandler interface for Query type and auth.Result type.
var _ iquery.IHandler[*Query, *auth.Result] = &Handler{}

// Config holds the dependencies needed to create a new Handler.
// It includes the user repository, JWT service, and hash service.
type Config struct {
	UserRepository irepository.IUserRepository
	JwtService     ijwt.IService
	HashService    hash.IService
}

// NewHandler creates a new Handler with the provided configuration.
// It initializes the Handler with the necessary services for processing login queries.
func NewHandler(config Config) *Handler {
	return &Handler{
		userRepo: config.UserRepository,
		jwtSvc:   config.JwtService,
		hashSvc:  config.HashService,
	}
}

// Handle processes a login query and returns an authentication result if successful.
// Returns:
// - *auth.Result: A pointer to the authentication result containing user ID, username, and token.
// - error: An error if the login fails. Possible errors include:
//   - InvalidCredential: If the username is not found or the password is incorrect.
//   - Unexpected: For unexpected errors during user retrieval or password validation.
func (h *Handler) Handle(query *Query) (*auth.Result, error) {
	user, err := h.userRepo.ByUsername(query.Username)
	if err != nil {
		var domainErr *errdmn.Error
		if errors.As(err, &domainErr) {
			return nil, appError.InvalidCredential(domainErr.Message)
		}
		return nil, errdmn.NewUnexpected(fmt.Sprintf("failed to retrieve user, %v", err))
	}

	isPasswordCorrect, err := h.hashSvc.Match(user.PasswordHash(), query.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to validate user password, %w", err)
	}

	if !isPasswordCorrect {
		return nil, appError.InvalidCredential("incorrect password")
	}

	token, err := h.jwtSvc.Generate(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT for user, %w", err)
	}

	return auth.NewResult(user.ID(), user.Username(), token), nil
}
