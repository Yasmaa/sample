package cmd

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
	"github.com/go-playground/validator/v10"
	"api/internal/delivery/http/router"
	"api/internal/infrastructure/datastore"
	"api/internal/infrastructure/inspector"
	"api/internal/infrastructure/redis"
	"api/internal/infrastructure/schedular"
	"api/internal/registry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	rootCmd.AddCommand(serve)
}

var serve = &cobra.Command{
	Use:   "serve",
	Short: "serve the API",
	Run: func(cmd *cobra.Command, args []string) {

		go http.ListenAndServe(":8080", nil)
		
		ps := datastore.NewPostgreSQL()


		scheduler := schedular.NewScheduler()
		go scheduler.Run()

		inspector.NewInspector()
		go redis.NewRedisClient()

		ze := zap.NewDevelopmentEncoderConfig()
		ze.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger := zap.New(zapcore.NewCore(
			zapcore.NewConsoleEncoder(ze),
			zapcore.AddSync(colorable.NewColorableStdout()),
			zapcore.DebugLevel,
		))

		v := validator.New()
		rg := registry.NewInteractor(ps, v, logger, scheduler)
		h := rg.NewAppHandler()
		
		g := router.NewRouter(logger, h)


		go g.Run(":8008")

		// Create channel for shutdown signals.
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)
		signal.Notify(stop, syscall.SIGTERM)
		//Recieve shutdown signals.
		<-stop

	},
}
