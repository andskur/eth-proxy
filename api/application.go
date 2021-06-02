package api

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-openapi/loads"
	"github.com/justinas/alice"

	"eth-proxy/api/config"
	"eth-proxy/api/handlers"
	"eth-proxy/api/middlewares"
	"eth-proxy/api/server"
	"eth-proxy/api/server/operations"
	"eth-proxy/lib"
	"eth-proxy/pkg/logger"
	version "eth-proxy/pkg/versioner"
	"eth-proxy/utils"
)

// App is
type App struct {
	version     *version.Version
	config      *config.Scheme
	ethProxyAPI lib.IEthProxyAPI
	server      *server.Server
}

// NewApplication create new CTL Application
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
	if err := app.initEthProxyAPI(app.config.Api); err != nil {
		return fmt.Errorf("initialize EthProxy API instance: %w", err)
	}

	if err := app.initServer(app.config.Http); err != nil {
		return fmt.Errorf("initialize application http server: %w", err)
	}
	return nil
}

// initEthProxyAPI is
func (app *App) initEthProxyAPI(cfg *config.Addr) (err error) {
	app.ethProxyAPI, err = lib.New(utils.CreateAddr(cfg.Host, cfg.Port))
	if err != nil {
		return fmt.Errorf("connect to clientAPI: %w", err)
	}

	logger.Log().Infof("Application EthProxy API initialized on %s", utils.CreateAddr(cfg.Host, cfg.Port))
	return nil
}

// initServer initialize Application HTTP server instance
func (app *App) initServer(cfg *config.Addr) error {
	logger.Log().Info("Application API server initializing...")

	api, handler, err := app.initAPI()
	if err != nil {
		return fmt.Errorf("init HTTP server: %w", err)
	}

	app.server = server.NewServer(api)
	app.server.Port = cfg.Port
	app.server.Host = cfg.Host

	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		return fmt.Errorf("init HTTP server: %w", err)
	}

	app.server.ReadTimeout = timeout
	app.server.WriteTimeout = timeout

	app.server.SetHandler(*handler)

	logger.Log().Info("Application API server initialized")
	return nil
}

// initAPI initialize Application API
func (app *App) initAPI() (*operations.ProxyAPIAPI, *http.Handler, error) {
	swaggerSpec, err := loads.Analyzed(server.SwaggerJSON, "")
	if err != nil {
		return nil, nil, fmt.Errorf("initialize API: %w", err)
	}

	api := operations.NewProxyAPIAPI(swaggerSpec)

	api.Logger = logger.Log().Infof

	app.initHandlers(api)

	handler := alice.New(
		middlewares.Recovery(),
		middlewares.Logger(),
		middlewares.Cors,
	).Then(api.Serve(nil))

	return api, &handler, nil
}

// initHandlers initialize Application handlers
func (app *App) initHandlers(api *operations.ProxyAPIAPI) {
	api.BlocksGetBlockHandler = handlers.NewGetBlock(app.ethProxyAPI)
	api.TxsGetTxByHashHandler = handlers.NewGetTxByHash(app.ethProxyAPI)
	api.BlocksGetTxByIndexHandler = handlers.NewGetTxByIndex(app.ethProxyAPI)
}

// Serve start the application, usually running all servers listener
func (app *App) Serve() error {
	go app.server.Serve()

	// Gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	<-quit
	logger.Log().Warning("Shutdown Application...")

	// close(closeChan)
	if err := app.Stop(); err != nil {
		return fmt.Errorf("stop application: %w", err)
	}

	logger.Log().Warning("Application stopped")
	return nil
}

// Stop the application, usually stop all servers listeners
func (app *App) Stop() error {
	if err := app.server.Shutdown(); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	if err := app.ethProxyAPI.Close(); err != nil {
		return fmt.Errorf("close EthProxy API instance: %w", err)
	}
	return nil
}

// Config return application Config Scheme instance
func (app *App) Config() *config.Scheme {
	return app.config
}

// Version return application current version
func (app *App) Version() string {
	return app.version.String()
}
