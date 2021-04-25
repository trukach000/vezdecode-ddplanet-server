package app

import (
	"ddplanet-server/pkg/config"
	"ddplanet-server/pkg/cors"
	"ddplanet-server/pkg/database"
	"ddplanet-server/pkg/httpext"
	"ddplanet-server/pkg/swagger"
	"net/http"
	"os"

	chilogrus "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi"
	chim "github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

func Setup() *chi.Mux {

	r := chi.NewRouter()
	log := logrus.New()

	cfg := config.Get()

	db, err := database.InitDatabase(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
	)

	if err != nil {
		logrus.Errorf("can't connect to database: %s", err)
		os.Exit(1)
	}

	r.Use(
		chilogrus.Logger("logger", log),
		chim.Recoverer,
		chim.NoCache,
		database.NewDatabaseMiddleware(db).Attach,
		cors.CORS(),
	)

	r.Get("/swagger/*", swagger.WrapSwagger)

	// frontend files
	httpext.ServeFile(r, "/site", "./site/index.html")
	httpext.ServeDir(r, "/site/*", http.Dir("./site"))

	// API
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/support/request", CreateSupportRequest)
		r.Get("/support/requests", GetSupportRequests)
	})

	return r
}
