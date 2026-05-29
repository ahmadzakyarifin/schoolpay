package parent

import (
	"github.com/ahmadzakyarifin/schoolpay/config"
	academichandler "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/delivery"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	academicusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/usecase"
	auditrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	financerepo "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
)

func SetupParentRoutes(g *gin.RouterGroup, db *bun.DB, cfg *config.Config, msg utils.Messenger, redisClient *redis.Client) {
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

	// Handlers
	stuHdl := academichandler.NewStudentHandler(stuSvc, cfg)

	students := g.Group("/students")
	{
		students.GET("/me", stuHdl.GetMyStudents)
	}
}
