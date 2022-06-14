package internal

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/meBazil/hr-mvp/internal/acl"
	"github.com/meBazil/hr-mvp/internal/auth"
	"github.com/meBazil/hr-mvp/internal/config"
	"github.com/meBazil/hr-mvp/internal/profile"
	"github.com/meBazil/hr-mvp/internal/rest"
	"github.com/meBazil/hr-mvp/internal/rest/handlers"
	"github.com/meBazil/hr-mvp/internal/vacancy"
	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
	"github.com/meBazil/hr-mvp/pkg/di-container"
	"github.com/meBazil/hr-mvp/pkg/process-manager"
	"github.com/meBazil/hr-mvp/pkg/sqlutil"
)

type Application struct {
	sigChan   <-chan os.Signal
	manager   *process.Manager
	container *di.Container
	cfg       config.App
}

func NewApplication(cfg config.App) (*Application, error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	a := &Application{
		sigChan:   sigChan,
		cfg:       cfg,
		manager:   process.NewManager(),
		container: di.GetDefaultContainer(),
	}

	err := a.bootstrap()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Application) Run() {
	a.manager.StartAll()
	a.registerShutdown()
}

func (a *Application) bootstrap() error {
	a.initServices()

	a.initREST()
	a.initHealthWorker()

	return nil
}

func (a *Application) initServices() {
	connectionManager, err := sqlutil.NewConnectionManager(a.cfg.DB)
	if err != nil {
		logger.Fatal("couldn't initialize connection manager", err, nil)
		return
	}

	a.container.Set(connectionManager)

	err = a.applyMigrations(connectionManager)
	if err != nil {
		logger.Fatal("couldn't initialize migrations", err, nil)
		return
	}

	profileService, err := profile.NewService(profile.NewRepository(connectionManager))
	if err != nil {
		logger.Fatal("couldn't init profile service", err, nil)
		return
	}

	a.container.Set(profileService)

	var authService = auth.NewService(profileService, a.cfg.JWT)
	a.container.Set(authService)

	vacancyService, err := vacancy.NewService(vacancy.NewRepository(connectionManager))
	if err != nil {
		logger.Fatal("couldn't init vacancy service", err, nil)
		return
	}

	a.container.Set(vacancyService)

	aclService, err := acl.NewAclService(a.cfg.ACL)
	if err != nil {
		logger.Fatal("couldn't init acl service", err, nil)
		return
	}

	a.container.Set(aclService)
}

func (a *Application) applyMigrations(conn *sqlutil.ConnectionManager) error {
	master, err := conn.GetWriteConnection().DB()
	if err != nil {
		return err
	}

	migrationDriver, err := postgres.WithInstance(master, &postgres.Config{})
	if err != nil {
		return err
	}

	manager, err := migrate.NewWithDatabaseInstance("file://"+a.cfg.DB.MigrationPath, "taxes", migrationDriver)
	if err != nil {
		return err
	}

	err = manager.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	logger.Info("Migrations were applied", nil)

	return nil
}

func (a *Application) initREST() {
	var authService *auth.Service
	if err := a.container.Load(&authService); err != nil {
		logger.Fatal("couldn't load auth service from container", err)
	}

	var profileService *profile.Service
	if err := a.container.Load(&profileService); err != nil {
		logger.Fatal("couldn't load profile service from container", err)
	}

	var vacancyService *vacancy.Service
	if err := a.container.Load(&vacancyService); err != nil {
		logger.Fatal("couldn't load vacancy service from container", err)
	}

	var aclService *acl.Service
	if err := a.container.Load(&aclService); err != nil {
		logger.Fatal("couldn't load acl service from container", err)
	}

	var apiHandlers = []handlers.APIHandler{
		handlers.NewAuthHandler(authService, profileService),
		handlers.NewProfileHandler(profileService),
		handlers.NewVacancyHandler(vacancyService, authService),
	}

	a.manager.AddWorker(
		process.NewServerWorker(
			"rest",
			rest.NewRestServer(a.cfg.REST, authService, aclService, apiHandlers),
		),
	)
}

func (a *Application) initHealthWorker() {
	server := rest.NewHealthCheckServer(a.cfg.Health.Listen, "/status", rest.DefaultHandler(a.manager))
	a.manager.AddWorker(process.NewServerWorker("health", server))
}

func (a *Application) registerShutdown() {
	go func(manager *process.Manager) {
		<-a.sigChan

		manager.StopAll()
	}(a.manager)

	a.manager.AwaitAll()
}
