package main

import (
	"fmt"
	"log"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/mohit/finance-go/internal/auth/hash"
	"github.com/mohit/finance-go/internal/auth/jwt"
	timeservice "github.com/mohit/finance-go/internal/common/time_service"
	"github.com/mohit/finance-go/internal/config"
	"github.com/mohit/finance-go/internal/db"
	"github.com/mohit/finance-go/internal/expense/handler"
	expenserepo "github.com/mohit/finance-go/internal/expense/repository"
	expensecmd "github.com/mohit/finance-go/internal/expense/service/command"
	expensqry "github.com/mohit/finance-go/internal/expense/service/query"
	"github.com/mohit/finance-go/internal/graph"
	"github.com/mohit/finance-go/internal/server/middleware"
	ratelimiter "github.com/mohit/finance-go/internal/server/rate_limiter"
	api "github.com/mohit/finance-go/internal/server/rest"
	"github.com/mohit/finance-go/internal/server/router"
	"github.com/mohit/finance-go/internal/user/handler"
	userrepo "github.com/mohit/finance-go/internal/user/repository"
	registercmd "github.com/mohit/finance-go/internal/user/service/command"
	loginqry "github.com/mohit/finance-go/internal/user/service/query"
	"golang.org/x/time/rate"
)

// Global variables to hold database configuration
var (
	serverPort = config.Envs.ServerPort
	rateLimit  = config.Envs.APIRate
	rateBurst  = config.Envs.RateBurst
	dbUser     = config.Envs.DBUser
	dbPassword = config.Envs.DBPassword
	dbName     = config.Envs.DBName
	dbHost     = config.Envs.DBHost
	dbPort     = config.Envs.DBPort
)

func main() {
	// Connect to the database
	database := db.Connect(db.Config{
		DbUser:     dbUser,
		DbPassword: dbPassword,
		DbName:     dbName,
		DbHost:     dbHost,
		DbPort:     dbPort,
		SSLMode:    config.Envs.DBSSLMode,
	})

	// Run database migrations
	db.RunMigrations(database)

	// Initialize services and repositories
	timeService := timeservice.New()
	userRepository := userrepo.New(database)
	expenseRepository := expenserepo.New(database)
	jwtService := initializeJWTService(timeService)
	hashService := hash.SingletonService()
	ipRateLimiter := ratelimiter.NewIPRateLimiter(rate.Limit(rateLimit), rateBurst, timeService)

	// Initialize middlewares
	authorizationMiddleware := middleware.Authorization(jwtService, true)
	populateClaimsMiddleware := middleware.Authorization(jwtService, false)
	rateLimitingMiddleware := middleware.RateLimitMiddleware(ipRateLimiter)

	// Initialize command and query handlers
	userRegisterCommandHandler := initializeUserRegisterHandler(userRepository, jwtService, hashService, timeService)
	userLoginQueryHandler := initializeUserLoginQueryHandler(userRepository, jwtService, hashService)
	addExpenseHandler := initializeAddExpenseHandler(userRepository, timeService)
	getExpenseHandler := initializeGetExpenseHandler(expenseRepository)
	getExpensesHandler := initializeGetExpensesHandler(expenseRepository)
	patchExpenseHandler := initializePatchExpenseHandler(expenseRepository)
	deleteExpenseHandler := initializeDeleteExpenseHandler(expenseRepository)

	userHandler := user.NewHandler(user.Config{
		UserRepository:  userRepository,
		RegisterHandler: userRegisterCommandHandler,
		LoginHandler:    userLoginQueryHandler,
	})

	// Expense routes
	expenseHandler := expense.NewHandler(expense.Config{
		AddHandler:         addExpenseHandler,
		GetHandler:         getExpenseHandler,
		PatchHandler:       patchExpenseHandler,
		GetMultipleHandler: getExpensesHandler,
		DeleteHandler:      deleteExpenseHandler,
	})

	resolver := graph.NewResolver(graph.ResolverConfig{
		GetExpenseHandler:         getExpenseHandler,
		GetMultipleExpenseHandler: getExpensesHandler,
		AddExpenseHandler:         addExpenseHandler,
		PatchExpenseHandler:       patchExpenseHandler,
	})

	graphHandler := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Create and run the server
	server := router.NewRouter(router.Config{
		Addr:                     fmt.Sprintf(":%s", serverPort),
		RestfullControllers:      []api.IController{userHandler, expenseHandler},
		GraphQlController:        graphHandler,
		AuthorizationMiddleware:  authorizationMiddleware,
		PopulateClaimsMiddleware: populateClaimsMiddleware,
		RateLimitMiddleware:      rateLimitingMiddleware,
	})

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initializePatchExpenseHandler(expenseRepository *expenserepo.Repository) *expensecmd.PatchHandler {
	return expensecmd.NewPatchHandler(expenseRepository)
}

func initializeDeleteExpenseHandler(expenseRepository *expenserepo.Repository) *expensecmd.DeleteHandler {
	return expensecmd.NewDeleteHandler(expenseRepository)
}

func initializeGetExpenseHandler(expenseRepository *expenserepo.Repository) *expensqry.GetHandler {
	return expensqry.NewGetHandler(expenseRepository)
}

func initializeGetExpensesHandler(expenseRepository *expenserepo.Repository) *expensqry.GetMultipleHandler {
	return expensqry.NewGetMultipleHandler(expenseRepository)
}

// initializeJWTService initializes and returns a new JWT service instance.
func initializeJWTService(timeService *timeservice.Service) *jwt.Service {
	return jwt.New(
		jwt.Config{
			SecretKey:   config.Envs.JWTSecret,
			Issuer:      config.Envs.ServerHost,
			ExpTime:     time.Duration(config.Envs.JWTExpirationInSeconds) * time.Second,
			TimeService: timeService,
		})
}

// initializeUserRegisterHandler initializes and returns a new user register command handler.
func initializeUserRegisterHandler(userRepo *userrepo.Repository, jwtService *jwt.Service, hashService *hash.Service, timeService *timeservice.Service) *registercmd.Handler {
	return registercmd.NewHandler(registercmd.Config{
		UserRepo: userRepo,
		JwtSvc:   jwtService,
		HashSvc:  hashService,
		TimeSvc:  timeService,
	})
}

// initializeUserLoginQueryHandler initializes and returns a new user login query handler.
func initializeUserLoginQueryHandler(userRepo *userrepo.Repository, jwtService *jwt.Service, hashService *hash.Service) *loginqry.Handler {
	return loginqry.NewHandler(loginqry.Config{
		UserRepository: userRepo,
		JwtService:     jwtService,
		HashService:    hashService,
	})
}

// initializeAddExpenseHandler initializes and returns a new add expense command handler.
func initializeAddExpenseHandler(userRepo *userrepo.Repository, timeService *timeservice.Service) *expensecmd.AddHandler {
	return expensecmd.NewAddHandler(expensecmd.Config{
		UserRepository: userRepo,
		TimeService:    timeService,
	})
}
