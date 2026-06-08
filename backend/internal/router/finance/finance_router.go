package finance

import (
	"github.com/ahmadzakyarifin/schoolpay/internal/middleware"
	financehandler "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/delivery"
	financeusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupFinanceRoutes(
	api *gin.RouterGroup,
	jwtSecret string,
	userRepo userauthrepo.UserRepo,
	paySvc financeusecase.PaymentService,
	sbSvc financeusecase.StudentBillService,
	redisClient *redis.Client,
) {
	finGroup := api.Group("/finance")
	finGroup.Use(middleware.AuthMiddleware(jwtSecret, userRepo, redisClient))
	finGroup.Use(middleware.RoleMiddleware("admin", "parent"))
	finGroup.Use(middleware.RateLimitPerUser("finance_private", 300))

	paymentLimit := middleware.RateLimitPerUser("finance_payment", 10)

	payHdl := financehandler.NewPaymentHandler(paySvc)
	sbHdl := financehandler.NewStudentBillHandler(sbSvc, paySvc)

	finGroup.POST("/payments", paymentLimit, payHdl.Process)
	finGroup.GET("/my-payments", payHdl.GetHistory)
	finGroup.POST("/payment-intent", paymentLimit, payHdl.CreateIntent)
	finGroup.GET("/my-bills", sbHdl.GetMyBills)
	finGroup.GET("/payments/:id/receipt", payHdl.GetReceipt)
}
