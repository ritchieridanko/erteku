package di

import (
	"github.com/ritchieridanko/erteku/services/auth/configs"
	"github.com/ritchieridanko/erteku/services/auth/internal/infra"
	"github.com/ritchieridanko/erteku/services/auth/internal/infra/cache"
	"github.com/ritchieridanko/erteku/services/auth/internal/infra/database"
	"github.com/ritchieridanko/erteku/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/erteku/services/auth/internal/infra/publisher"
	"github.com/ritchieridanko/erteku/services/auth/internal/repositories"
	"github.com/ritchieridanko/erteku/services/auth/internal/repositories/caches"
	"github.com/ritchieridanko/erteku/services/auth/internal/repositories/databases"
	"github.com/ritchieridanko/erteku/services/auth/internal/transport/handlers"
	"github.com/ritchieridanko/erteku/services/auth/internal/transport/server"
	"github.com/ritchieridanko/erteku/services/auth/internal/usecases"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/bcrypt"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/jwt"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/validator"
)

type Container struct {
	config     *configs.Config
	cache      *cache.Cache
	database   *database.Database
	transactor *database.Transactor
	logger     *logger.Logger

	acp *publisher.Publisher

	acc caches.AuthCache
	tcc caches.TokenCache

	adb databases.AuthDatabase
	sdb databases.SessionDatabase

	ar repositories.AuthRepository
	sr repositories.SessionRepository
	tr repositories.TokenRepository

	bcrypt    *bcrypt.BCrypt
	jwt       *jwt.JWT
	validator *validator.Validator

	au usecases.AuthUsecase
	su usecases.SessionUsecase

	ah     *handlers.AuthHandler
	server *server.Server
}

func Init(cfg *configs.Config, i *infra.Infra) *Container {
	// Infra
	c := cache.NewCache(&cfg.Cache, i.Cache())
	db := database.NewDatabase(i.Database())
	tx := database.NewTransactor(i.Database())
	l := logger.NewLogger(i.Logger())

	// Publishers
	acp := publisher.NewPublisher(i.PublisherAC())

	// Caches
	acc := caches.NewAuthCache(c)
	tcc := caches.NewTokenCache(&cfg.Auth, c)

	// Databases
	adb := databases.NewAuthDatabase(db)
	sdb := databases.NewSessionDatabase(db)

	// Repositories
	ar := repositories.NewAuthRepository(adb, acc)
	sr := repositories.NewSessionRepository(sdb)
	tr := repositories.NewTokenRepository(tcc)

	// Utils
	b := bcrypt.Init(cfg.Auth.BCrypt.Cost)
	j := jwt.Init(cfg.Auth.JWT.Issuer, cfg.Auth.JWT.Secret, cfg.Auth.JWT.Duration)
	v := validator.Init()

	// Usecases
	au := usecases.NewAuthUsecase(cfg.App.Name, ar, tr, tx, acp, v, b, l)
	su := usecases.NewSessionUsecase(cfg.App.Name, cfg.Auth.Duration.Session, sr, tx, j)

	// Handlers
	ah := handlers.NewAuthHandler(au, su, l)

	// Server
	srv := server.Init(cfg.App.Name, &cfg.Server, l, ah)

	return &Container{
		config:     cfg,
		cache:      c,
		database:   db,
		transactor: tx,
		logger:     l,
		acp:        acp,
		acc:        acc,
		tcc:        tcc,
		adb:        adb,
		sdb:        sdb,
		ar:         ar,
		sr:         sr,
		tr:         tr,
		bcrypt:     b,
		jwt:        j,
		validator:  v,
		au:         au,
		su:         su,
		ah:         ah,
		server:     srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
