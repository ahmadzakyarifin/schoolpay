package parent

import (
	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/ahmadzakyarifin/schoolpay/internal/middleware"
	academichandler "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/delivery"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	academicusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/usecase"
	auditrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	financerepo "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	supporthandler "github.com/ahmadzakyarifin/schoolpay/internal/module/support/delivery"
	supportrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/support/repository"
	supportusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/support/usecase"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/internal/websocket"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"golang.org/x/time/rate"
)

func SetupParentRoutes(g *gin.RouterGroup, db *bun.DB, cfg *config.Config, msg utils.Messenger, redisClient *redis.Client, hub *websocket.Hub) {
	// Repositories
	stuRepo := academicrepo.NewStudentRepo(db)
	notiRepo := notificationrepo.NewNotificationRepo(db)
	userRepo := userauthrepo.NewUserRepo(db)
	authRepo := userauthrepo.NewAuthRepo(db, redisClient)

	majorRepo := academicrepo.NewMajorRepo(db)
	classRepo := academicrepo.NewClassRepo(db)
	ayRepo := academicrepo.NewAcademicYearRepo(db)
	sbRepo := financerepo.NewStudentBillRepo(db)

	auditRepo := auditrepo.NewAuditLogRepo(db)
	auditSvc := auditusecase.NewAuditLogService(auditRepo)
	stuSvc := academicusecase.NewStudentService(db, stuRepo, userRepo, authRepo, msg, notiRepo, ayRepo, majorRepo, classRepo, sbRepo, cfg, auditSvc)
	waSvc := notificationusecase.NewWhatsAppService()
	supportRepo := supportrepo.NewSupportRepo(db)
	supportSvc := supportusecase.NewSupportService(db, supportRepo, userRepo, waSvc, auditSvc)

	// Handlers
	stuHdl := academichandler.NewStudentHandler(stuSvc, cfg)
	supportHdl := supporthandler.NewSupportHandler(supportSvc, hub)
	parentSupportLimit := middleware.RateLimitMiddleware(redisClient, "parent_support", rate.Limit(6.0/60.0), 6)

	students := g.Group("/students")
	{
		students.GET("/me", stuHdl.GetMyStudents)
		students.GET("/:id/class-history", stuHdl.GetMyStudentClassHistory)
	}

	support := g.Group("/support")
	{
		support.GET("/conversation", supportHdl.ParentConversation)
		support.GET("/messages", supportHdl.ParentMessages)
		support.POST("/messages", parentSupportLimit, supportHdl.ParentSendMessage)
	}
}
