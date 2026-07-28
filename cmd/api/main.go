package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rrhhumand/api/internal/attendance"
	"github.com/rrhhumand/api/internal/auth"
	"github.com/rrhhumand/api/internal/branches"
	"github.com/rrhhumand/api/internal/companies"
	"github.com/rrhhumand/api/internal/config"
	"github.com/rrhhumand/api/internal/departments"
	"github.com/rrhhumand/api/internal/document_categories"
	"github.com/rrhhumand/api/internal/documents"
	"github.com/rrhhumand/api/internal/employees"
	"github.com/rrhhumand/api/internal/feed"
	"github.com/rrhhumand/api/internal/handlers"
	"github.com/rrhhumand/api/internal/leave"
	"github.com/rrhhumand/api/internal/organization"
	"github.com/rrhhumand/api/internal/overtime"
	"github.com/rrhhumand/api/internal/payroll"
	"github.com/rrhhumand/api/internal/performance"
	"github.com/rrhhumand/api/internal/positions"
	"github.com/rrhhumand/api/internal/profile"
	"github.com/rrhhumand/api/internal/recruitment"
	"github.com/rrhhumand/api/internal/roles"
	"github.com/rrhhumand/api/internal/scheduling"
	"github.com/rrhhumand/api/internal/server"
	"github.com/rrhhumand/api/internal/surveys"
	"github.com/rrhhumand/api/internal/users"
	"github.com/rrhhumand/api/pkg/database"
	"github.com/rrhhumand/api/pkg/logger"
)

func main() {
	cfg := config.Load()

	if err := logger.Init("debug", "console"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Get().Sync()

	logger.Info("Starting application",
		logger.String("app", cfg.AppName),
		logger.String("env", cfg.AppEnv),
		logger.String("port", cfg.AppPort),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbCfg := database.NewConfig()
	pool, err := database.Connect(ctx, dbCfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", logger.Err(err))
	}
	defer pool.Close()
	logger.Info("Database connected successfully")

	healthHandler := handlers.NewHealthHandler(pool)

	userRepo := users.NewUserRepository(pool)
	roleRepo := roles.NewRoleRepository(pool)
	refreshTokenRepo := auth.NewRefreshTokenRepository(pool)
	jwtService := auth.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.Expiration,
		cfg.JWT.RefreshExpiration,
	)
	authService := auth.NewAuthService(userRepo, refreshTokenRepo, jwtService, roleRepo)
	authHandler := auth.NewAuthHandler(authService)

	companyRepo := companies.NewCompanyRepository(pool)
	companyService := companies.NewCompanyService(companyRepo)
	companyHandler := companies.NewCompanyHandler(companyService)

	branchRepo := branches.NewBranchRepository(pool)
	branchService := branches.NewBranchService(branchRepo)
	branchHandler := branches.NewBranchHandler(branchService)

	departmentRepo := departments.NewDepartmentRepository(pool)
	departmentService := departments.NewDepartmentService(departmentRepo)
	departmentHandler := departments.NewDepartmentHandler(departmentService)

	positionRepo := positions.NewPositionRepository(pool)
	positionService := positions.NewPositionService(positionRepo)
	positionHandler := positions.NewPositionHandler(positionService)

	employeeRepo := employees.NewEmployeeRepository(pool)
	orgRepo := organization.NewOrgRepository(pool)
	employeeService := employees.NewEmployeeService(employeeRepo, orgRepo)
	employeeHandler := employees.NewEmployeeHandler(employeeService)

	orgHandler := organization.NewOrgHandler(orgRepo)

	profileRepo := profile.NewProfileRepository(pool)
	profileService := profile.NewProfileService(profileRepo)
	profileHandler := profile.NewProfileHandler(profileService)

	feedRepo := feed.NewFeedRepository(pool)
	feedService := feed.NewFeedService(feedRepo)
	feedHandler := feed.NewFeedHandler(feedService)

	surveyRepo := surveys.NewSurveyRepository(pool)
	surveyService := surveys.NewSurveyService(surveyRepo)
	surveyStatsService := surveys.NewSurveyStatsService()
	surveyHandler := surveys.NewSurveyHandler(surveyService, surveyStatsService)

	docRepo := documents.NewDocumentRepository(pool)
	docStorage, err := documents.NewStorageService(cfg.MinIO)
	if err != nil {
		logger.Warn("MinIO not available, document storage disabled", logger.Err(err))
	}
	docService := documents.NewDocumentService(docRepo, docStorage)
	docHandler := documents.NewDocumentHandler(docService)

	categoryRepo := document_categories.NewCategoryRepository(pool)
	categoryService := document_categories.NewCategoryService(categoryRepo)
	categoryHandler := document_categories.NewCategoryHandler(categoryService)

	leaveRepo := leave.NewRepository(pool)
	leaveCalc := leave.NewDayCalculator(leaveRepo)
	leaveService := leave.NewService(leaveRepo, leaveCalc)
	leaveHandler := leave.NewHandler(leaveService)

	attRepo := attendance.NewRepository(pool)
	attPunches := attendance.NewPunches(pool)
	attCalc := attendance.NewCalculator()
	attGeofence := attendance.NewGeoFence()
	attService := attendance.NewService(attRepo, attPunches, attCalc, attGeofence)
	attHandler := attendance.NewHandler(attService)

	schedRepo := scheduling.NewRepository(pool)
	schedService := scheduling.NewService(schedRepo)
	schedHandler := scheduling.NewHandler(schedService)

	otRepo := overtime.NewRepository(pool)
	otService := overtime.NewService(otRepo)
	otHandler := overtime.NewHandler(otService)

	payRepo := payroll.NewRepository(pool)
	payService := payroll.NewService(payRepo)
	payHandler := payroll.NewHandler(payService)

	perfRepo := performance.NewRepository(pool)
	perfService := performance.NewService(perfRepo)
	perfHandler := performance.NewHandler(perfService)

	recRepo := recruitment.NewRepository(pool)
	recService := recruitment.NewService(recRepo)
	recHandler := recruitment.NewHandler(recService)

	router := server.NewRouter(
		healthHandler,
		authHandler,
		jwtService,
		companyHandler,
		branchHandler,
		departmentHandler,
		positionHandler,
		employeeHandler,
		orgHandler,
		profileHandler,
		feedHandler,
		surveyHandler,
		docHandler,
		categoryHandler,
		leaveHandler,
		attHandler,
		schedHandler,
		otHandler,
		payHandler,
		perfHandler,
		recHandler,
		pool,
	)

	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("HTTP server listening", logger.String("port", cfg.AppPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", logger.Err(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Server forced to shutdown", logger.Err(err))
	}

	logger.Info("Server exited gracefully")
}
