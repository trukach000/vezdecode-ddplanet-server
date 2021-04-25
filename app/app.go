package app

import (
	"ddplanet-server/pkg/config"
	"ddplanet-server/pkg/database"
	"ddplanet-server/pkg/httpext"
	"ddplanet-server/pkg/swagger"
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
		chim.RedirectSlashes,
		database.NewDatabaseMiddleware(db).Attach,
	)

	httpext.ServeFile(r, "/photo", "data/photo.jpg")

	r.Get("/swagger/*", swagger.WrapSwagger)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/support/request", CreateSupportRequest)
	})

	return r
}
