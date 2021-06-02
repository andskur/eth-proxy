package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/ethclient"

	"eth-proxy/pkg/logger"
	"eth-proxy/pkg/versioner"
	"eth-proxy/proto"
	"eth-proxy/proxy-service/config"
	"eth-proxy/proxy-service/server"
	"eth-proxy/proxy-service/service"
	"eth-proxy/utils"
)

// App is basic application structure
type App struct {
	version   *version.Version
	config    *config.Scheme
	ethClient *ethclient.Client
	grpc      server.IGrpc
	srv       service.IService
}

// NewApplication create new Hub Application
func NewApplication() (app *App, err error) {
	app = &App{config: &config.Scheme{}}

	app.version, err = version.NewVersion()
	if err != nil {
		return nil, fmt.Errorf("setup app version: %w", err)
	}

	return
}

// Init initialize application and all necessary instances
func (app *App) Init() error {
	if err := app.initEthereum(app.config.Ethereum); err != nil {
		return fmt.Errorf("application Ethereum client initialisation: %w", err)
	}

	if err := app.initService(); err != nil {
		return fmt.Errorf("service layer initialising: %w", err)
	}

	if err := app.initGrpc(app.config.Grpc); err != nil {
		return fmt.Errorf("application GRPC server initialisation: %w", err)
	}

	return nil
}

// initEthereum initialise Application Ethereum client
func (app *App) initEthereum(cfg *config.Ethereum) (err error) {
	app.ethClient, err = ethclient.Dial(cfg.Addr)
	if err != nil {
		return fmt.Errorf("connecting to Ethereum node at %s: %w", cfg.Addr, err)
	}

	logger.Log().Infof("successfully established connection to Ethereum node at %s", cfg.Addr)
	return nil
}

// initService initialize Application service layer instance
func (app *App) initService() (err error) {
	app.srv, err = service.NewService(app.ethClient)
	if err != nil {
		return fmt.Errorf("create new service instance %w", err)
	}

	logger.Log().Info("service layer successfully initialised")
	return nil
}

// initGrpc initialize Application GRPC server instance
func (app *App) initGrpc(cfg *config.Grpc) (err error) {
	app.grpc, err = server.NewServer(cfg)
	if err != nil {
		return fmt.Errorf("create new GRPC server instance: %w", err)
	}

	proto.RegisterEthProxyServiceServer(app.grpc.GetServer(), server.NewEthProxyAPI(app.srv))

	logger.Log().Infof("App GRPC server initialized on %s", utils.CreateAddr(cfg.Host, cfg.Port))
	return nil
}

// Serve start the application
func (app *App) Serve() error {
	go app.grpc.Listen()

	// Gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	<-quit
	logger.Log().Warning("Shutdown Application...")

	app.Stop()

	logger.Log().Warning("Application stopped")
	return nil
}

// Stop the application
func (app *App) Stop() {
	app.grpc.Close()
	app.ethClient.Close()
}

// Config return application Config Scheme instance
func (app *App) Config() *config.Scheme {
	return app.config
}

// Version return application current version
func (app *App) Version() string {
	return app.version.String()
}
