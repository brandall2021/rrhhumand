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
	"github.com/rrhhumand/api/internal/compensation"
	"github.com/rrhhumand/api/internal/config"
	"github.com/rrhhumand/api/internal/departments"
	"github.com/rrhhumand/api/internal/document_categories"
	"github.com/rrhhumand/api/internal/documents"
	"github.com/rrhhumand/api/internal/employees"
	"github.com/rrhhumand/api/internal/events"
	"github.com/rrhhumand/api/internal/feed"
	"github.com/rrhhumand/api/internal/handlers"
	"github.com/rrhhumand/api/internal/leave"
	"github.com/rrhhumand/api/internal/notifications"
	"github.com/rrhhumand/api/internal/onboarding"
	"github.com/rrhhumand/api/internal/organization"
	"github.com/rrhhumand/api/internal/overtime"
	"github.com/rrhhumand/api/internal/payroll"
	featapp "github.com/rrhhumand/api/internal/payroll/features/application"
	featengine "github.com/rrhhumand/api/internal/payroll/features/engine"
	feathttp "github.com/rrhhumand/api/internal/payroll/features/http"
	featrepo "github.com/rrhhumand/api/internal/payroll/features/repository"
	"github.com/rrhhumand/api/internal/payroll/features"
	benefitsapp "github.com/rrhhumand/api/internal/benefits/application"
	benefitshttp "github.com/rrhhumand/api/internal/benefits/http"
	benefitsrepo "github.com/rrhhumand/api/internal/benefits/repository"
	"github.com/rrhhumand/api/internal/benefits"
	expensesapp "github.com/rrhhumand/api/internal/expenses/application"
	expenseshttp "github.com/rrhhumand/api/internal/expenses/http"
	expensesrepo "github.com/rrhhumand/api/internal/expenses/repository"
	"github.com/rrhhumand/api/internal/expenses"
	"github.com/rrhhumand/api/internal/performance"
	"github.com/rrhhumand/api/internal/positions"
	"github.com/rrhhumand/api/internal/profile"
	recrapp "github.com/rrhhumand/api/internal/recruitment/application"
	recrengine "github.com/rrhhumand/api/internal/recruitment/engine"
	recrhttp "github.com/rrhhumand/api/internal/recruitment/http"
	recrrepo "github.com/rrhhumand/api/internal/recruitment/repository"
	"github.com/rrhhumand/api/internal/recruitment"
	"github.com/rrhhumand/api/internal/roles"
	"github.com/rrhhumand/api/internal/scheduling"
	"github.com/rrhhumand/api/internal/server"
	"github.com/rrhhumand/api/internal/surveys"
	"github.com/rrhhumand/api/internal/training"
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

	payRepo := payroll.NewRepository(pool, logger.Get())
	payService := payroll.NewService(payRepo, logger.Get())
	payHandler := payroll.NewHandler(payService)

	perfRepo := performance.NewRepository(pool)
	perfService := performance.NewService(perfRepo)
	perfHandler := performance.NewHandler(perfService)

	recReqRepo := recrrepo.NewRequisitionRepo(pool)
	recPosRepo := recrrepo.NewPositionRepo(pool)
	recPostRepo := recrrepo.NewPostingRepo(pool)
	recCandRepo := recrrepo.NewCandidateRepo(pool)
	recAppRepo := recrrepo.NewApplicationRepo(pool)
	recIntRepo := recrrepo.NewInterviewRepo(pool)
	recAssessRepo := recrrepo.NewAssessmentRepo(pool)
	recOfferRepo := recrrepo.NewOfferRepo(pool)
	recHiringRepo := recrrepo.NewHiringProcessRepo(pool)
	recWfRepo := recrrepo.NewWorkflowRepo(pool)
	recScoreRepo := recrrepo.NewScoringRepo(pool)
	recEmailRepo := recrrepo.NewEmailRepo(pool)
	recDashRepo := recrrepo.NewDashboardRepo(pool)
	recSettingsRepo := recrrepo.NewSettingsRepo(pool)

	recScoreEngine := recrengine.NewScoringEngine()
	recMatchEngine := recrengine.NewMatchingEngine(recScoreEngine, recCandRepo, recPosRepo, recScoreRepo)

	recReqSvc := recrapp.NewRequisitionService(recReqRepo, recPosRepo)
	recPosSvc := recrapp.NewPositionService(recPosRepo, recReqRepo)
	recPostSvc := recrapp.NewPostingService(recPostRepo, recPosRepo)
	recCandSvc := recrapp.NewCandidateService(recCandRepo)
	recAppSvc := recrapp.NewApplicationService(recAppRepo, recCandRepo, recPostRepo)
	recIntSvc := recrapp.NewInterviewService(recIntRepo, recAppRepo)
	recAssessSvc := recrapp.NewAssessmentService(recAssessRepo, recAppRepo)
	recOfferSvc := recrapp.NewOfferService(recOfferRepo, recAppRepo)
	recHiringSvc := recrapp.NewHiringService(recHiringRepo, recOfferRepo, recAppRepo, recCandRepo)
	recWfSvc := recrapp.NewWorkflowService(recWfRepo)
	recScoreSvc := recrapp.NewScoringService(recScoreRepo, recScoreEngine, recMatchEngine, recCandRepo, recPosRepo)
	recEmailSvc := recrapp.NewEmailService(recEmailRepo, recEmailRepo)
	recDashSvc := recrapp.NewDashboardService(recDashRepo)
	recSettingsSvc := recrapp.NewSettingsService(recSettingsRepo)

	recHandler := recrhttp.NewHandler(
		recReqSvc, recPosSvc, recPostSvc, recCandSvc, recAppSvc,
		recIntSvc, recAssessSvc, recOfferSvc, recHiringSvc, recWfSvc,
		recScoreSvc, recEmailSvc, recDashSvc, recSettingsSvc,
	)

	onboardingRepo := onboarding.NewRepository(pool)
	onboardingService := onboarding.NewService(onboardingRepo)
	onboardingHandler := onboarding.NewHandler(onboardingService)

	trainingRepo := training.NewRepository(pool, logger.Get())
	eventsSvc := events.NewService()
	notifSvc := notifications.NewService()
	trainingService := training.NewService(trainingRepo, eventsSvc, notifSvc, nil, logger.Get())
	trainingHandler := training.NewHandler(trainingService)

	compRepo := compensation.NewRepository(pool, logger.Get())
	compService := compensation.NewService(compRepo, logger.Get())
	compHandler := compensation.NewHandler(compService)

	// FASE 19B — Payroll Features
	featReceiptRepo := featrepo.NewReceiptRepo(pool)
	featArcaRepo := featrepo.NewArcaRepo(pool)
	featBookRepo := featrepo.NewBookRepo(pool)
	featBankRepo := featrepo.NewBankRepo(pool)
	featAccountingRepo := featrepo.NewAccountingRepo(pool)
	featReportRepo := featrepo.NewReportRepo(pool)

	featReceiptSvc := featapp.NewReceiptService(featReceiptRepo)
	featArcaSvc := featapp.NewArcaService(featArcaRepo)
	featBookSvc := featapp.NewBookService(featBookRepo)
	featBankSvc := featapp.NewBankService(featBankRepo)
	featAccountingSvc := featapp.NewAccountingService(featAccountingRepo)
	featReportSvc := featapp.NewReportService(featReportRepo)

	// Engine generators
	featReceiptGen := featengine.NewReceiptGenerator()
	featArcaGen := featengine.NewArcaGenerator()
	featBankGen := featengine.NewBankGenerator()
	featAccGen := featengine.NewAccountingGenerator()
	featReportGen := featengine.NewReportGenerator()
	_ = featReceiptGen
	_ = featArcaGen
	_ = featBankGen
	_ = featAccGen
	_ = featReportGen

	featuresHandler := feathttp.NewHandler(featReceiptSvc, featArcaSvc, featBookSvc, featBankSvc, featAccountingSvc, featReportSvc)

	// FASE 20 — Benefits & Total Rewards
	benCatalogRepo := benefitsrepo.NewCatalogRepo(pool)
	benBenefitRepo := benefitsrepo.NewBenefitRepo(pool)
	benPlanRepo := benefitsrepo.NewPlanRepo(pool)
	benEligibilityRepo := benefitsrepo.NewEligibilityRepo(pool)
	benWorkflowRepo := benefitsrepo.NewWorkflowRepo(pool)
	benAssignmentRepo := benefitsrepo.NewAssignmentRepo(pool)
	benWalletRepo := benefitsrepo.NewWalletRepo(pool)
	benReimbursementRepo := benefitsrepo.NewReimbursementRepo(pool)
	benFlexibleRepo := benefitsrepo.NewFlexibleRepo(pool)
	benCostRepo := benefitsrepo.NewCostRepo(pool)
	benBonusRepo := benefitsrepo.NewBonusRepo(pool)
	benRewardsRepo := benefitsrepo.NewRewardsRepo(pool)

	benCatalogSvc := benefitsapp.NewCatalogService(benCatalogRepo)
	benBenefitSvc := benefitsapp.NewBenefitService(benBenefitRepo, benCatalogRepo, benPlanRepo)
	benEligibilitySvc := benefitsapp.NewEligibilityService(benEligibilityRepo, benBenefitRepo, nil)
	benWorkflowSvc := benefitsapp.NewWorkflowService(benWorkflowRepo)
	benAssignmentSvc := benefitsapp.NewAssignmentService(benAssignmentRepo, benBenefitRepo, benWalletRepo)
	benWalletSvc := benefitsapp.NewWalletService(benWalletRepo, benFlexibleRepo)
	benReimbursementSvc := benefitsapp.NewReimbursementService(benReimbursementRepo, benWalletRepo)
	benBonusSvc := benefitsapp.NewBonusService(benBonusRepo)
	benRewardsSvc := benefitsapp.NewRewardsService(benRewardsRepo, benBonusRepo)
	benCostSvc := benefitsapp.NewCostService(benCostRepo)

	benefitsHandler := benefitshttp.NewHandler(
		benCatalogSvc, benBenefitSvc, benEligibilitySvc, benWorkflowSvc,
		benAssignmentSvc, benWalletSvc, benReimbursementSvc, benBonusSvc,
		benRewardsSvc, benCostSvc,
	)

	// FASE 21 — Expenses & Travel
	expCatalogRepo := expensesrepo.NewCatalogRepo(pool)
	expExpenseRepo := expensesrepo.NewExpenseRepo(pool)
	expTravelRepo := expensesrepo.NewTravelRepo(pool)
	expReportRepo := expensesrepo.NewReportRepo(pool)
	expAdvanceRepo := expensesrepo.NewAdvanceRepo(pool)
	expReimbursementRepo := expensesrepo.NewReimbursementRepo(pool)
	expReceiptRepo := expensesrepo.NewReceiptRepo(pool)
	expPolicyRepo := expensesrepo.NewPolicyRepo(pool)
	expApprovalRepo := expensesrepo.NewApprovalRepo(pool)
	expWorkflowRepo := expensesrepo.NewWorkflowRepo(pool)
	expBudgetRepo := expensesrepo.NewBudgetRepo(pool)
	expExchangeRepo := expensesrepo.NewExchangeRepo(pool)
	expAllowanceRepo := expensesrepo.NewAllowanceRepo(pool)
	expCardRepo := expensesrepo.NewCardRepo(pool)
	expAuditRepo := expensesrepo.NewAuditRepo(pool)
	expDuplicateRepo := expensesrepo.NewDuplicateRepo(pool)

	expCatalogSvc := expensesapp.NewCatalogService(expCatalogRepo, expAuditRepo)
	expExpenseSvc := expensesapp.NewExpenseService(expExpenseRepo, expAuditRepo, expReceiptRepo, expDuplicateRepo, nil)
	expTravelSvc := expensesapp.NewTravelService(expTravelRepo, expAuditRepo)
	expReportSvc := expensesapp.NewReportService(expReportRepo, expExpenseRepo, expAdvanceRepo, expAuditRepo, nil)
	expAdvanceSvc := expensesapp.NewAdvanceService(expAdvanceRepo, expAuditRepo)
	expReimbursementSvc := expensesapp.NewReimbursementService(expReimbursementRepo, expAuditRepo)
	expPolicySvc := expensesapp.NewPolicyService(expPolicyRepo, expExpenseRepo)
	expApprovalSvc := expensesapp.NewApprovalService(expApprovalRepo, expWorkflowRepo, expExpenseRepo, expTravelRepo, expReportRepo, expAdvanceRepo)
	expBudgetSvc := expensesapp.NewBudgetService(expBudgetRepo)
	expExchangeSvc := expensesapp.NewExchangeService(expExchangeRepo)
	expAllowanceSvc := expensesapp.NewAllowanceService(expAllowanceRepo)

	expensesHandler := expenseshttp.NewHandler(
		expCatalogSvc, expExpenseSvc, expTravelSvc, expReportSvc,
		expAdvanceSvc, expReimbursementSvc, expPolicySvc, expApprovalSvc,
		expBudgetSvc, expExchangeSvc, expAllowanceSvc,
	)

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
		onboardingHandler,
		trainingHandler,
		compHandler,
		featuresHandler,
		benefitsHandler,
		expensesHandler,
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

	onboardingWorker := onboarding.NewWorker(onboardingRepo)
	go onboardingWorker.Run(ctx)

	trainingWorker := training.NewWorker(trainingService, logger.Get())
	go trainingWorker.Start(ctx, "", 5*time.Minute)

	compWorker := compensation.NewWorker(compService, compRepo)
	compWorker.Start(30 * time.Minute)

	payWorker := payroll.NewWorker(payService, payRepo)
	payWorker.Start(15 * time.Minute)

	// FASE 19B workers
	receiptWorker := features.NewReceiptWorker()
	exportWorker := features.NewExportWorker()
	bankWorker := features.NewBankWorker()
	accountingWorker := features.NewAccountingWorker()
	reportWorker := features.NewReportWorker()
	receiptWorker.Start(5 * time.Minute)
	exportWorker.Start(10 * time.Minute)
	bankWorker.Start(10 * time.Minute)
	accountingWorker.Start(15 * time.Minute)
	reportWorker.Start(15 * time.Minute)

	// FASE 20 — Benefits workers
	benefits.StartBenefitWorkers(pool)

	// FASE 21 — Expense workers
	expenses.StartExpenseWorkers(pool)

	// ATS workers
	recruitment.StartRecruitmentWorkers(pool)

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
