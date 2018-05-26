package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/anovakovic01/tap-backend/auth/google"
	"github.com/anovakovic01/tap-backend/events"
	eventshttp "github.com/anovakovic01/tap-backend/events/http"
	eventsdb "github.com/anovakovic01/tap-backend/events/postgres"
	"github.com/anovakovic01/tap-backend/news"
	newshttp "github.com/anovakovic01/tap-backend/news/http"
	newsdb "github.com/anovakovic01/tap-backend/news/postgres"
	"github.com/anovakovic01/tap-backend/news/rss"
	"github.com/go-zoo/bone"
	"github.com/mmcdole/gofeed"
)

const (
	envGoogleClientID = "TAP_GOOGLE_CLIENT_ID"
	envGoogleSecret   = "TAP_GOOGLE_SECRET"
	envNewsDBHost     = "TAP_NEWS_DB_HOST"
	envNewsDBPort     = "TAP_NEWS_DB_PORT"
	envNewsDBName     = "TAP_NEWS_DB_NAME"
	envNewsDBUser     = "TAP_NEWS_DB_USER"
	envNewsDBPass     = "TAP_NEWS_DB_PASS"
	envEventsDBHost   = "TAP_EVENTS_DB_HOST"
	envEventsDBPort   = "TAP_EVENTS_DB_PORT"
	envEventsDBName   = "TAP_EVENTS_DB_NAME"
	envEventsDBUser   = "TAP_EVENTS_DB_USER"
	envEventsDBPass   = "TAP_EVENTS_DB_PASS"
	envPort           = "PORT"

	defNewsDBHost   = "localhost"
	defNewsDBPort   = "5432"
	defNewsDBName   = "news"
	defNewsDBUser   = "postgres"
	defNewsDBPass   = "postgres"
	defEventsDBHost = "localhost"
	defEventsDBPort = "5432"
	defEventsDBName = "events"
	defEventsDBUser = "postgres"
	defEventsDBPass = "postgres"
	defPort         = "8000"
)

type config struct {
	clientID     string
	secret       string
	newsDBHost   string
	newsDBPort   string
	newsDBName   string
	newsDBUser   string
	newsDBPass   string
	eventsDBHost string
	eventsDBPort string
	eventsDBName string
	eventsDBUser string
	eventsDBPass string
	port         string
}

func load() config {
	return config{
		clientID:     getenv(envGoogleClientID, ""),
		secret:       getenv(envGoogleSecret, ""),
		newsDBHost:   getenv(envNewsDBHost, defNewsDBHost),
		newsDBPort:   getenv(envNewsDBPort, defNewsDBPort),
		newsDBName:   getenv(envNewsDBName, defNewsDBName),
		newsDBUser:   getenv(envNewsDBUser, defNewsDBUser),
		newsDBPass:   getenv(envNewsDBPass, defNewsDBPass),
		eventsDBHost: getenv(envEventsDBHost, defEventsDBHost),
		eventsDBPort: getenv(envEventsDBPort, defEventsDBPort),
		eventsDBName: getenv(envEventsDBName, defEventsDBName),
		eventsDBUser: getenv(envEventsDBUser, defEventsDBUser),
		eventsDBPass: getenv(envEventsDBPass, defEventsDBPass),
		port:         getenv(envPort, defPort),
	}
}

func main() {
	cfg := load()
	newsDB, err := newsdb.Connect(cfg.newsDBHost, cfg.newsDBPort, cfg.newsDBName, cfg.newsDBUser, cfg.newsDBPass)
	if err != nil {
		fmt.Printf("Failed to connect to postgres: %s\n", err)
		os.Exit(1)
	}

	eventsDB, err := eventsdb.Connect(cfg.eventsDBHost, cfg.eventsDBPort, cfg.eventsDBName, cfg.eventsDBUser, cfg.eventsDBPass)
	if err != nil {
		fmt.Printf("Failed to connect to postgres: %s\n", err)
		os.Exit(1)
	}

	authSvc := google.NewService(cfg.clientID)

	parser := gofeed.NewParser()
	collector := rss.NewCollector(parser)
	newsRepo := newsdb.NewNewsRepositry(newsDB)
	newsSvc := news.NewService(newsRepo, collector)

	eventsRepo := eventsdb.NewEventsRepository(eventsDB)
	eventsSvc := events.NewService(eventsRepo)

	go newsSvc.Collect("https://auto.economictimes.indiatimes.com/rss/topstories")

	mux := bone.New()
	newshttp.UpdateHandler(mux, newsSvc, authSvc)
	eventshttp.UpdateHandler(mux, eventsSvc, authSvc)
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%s", cfg.port), mux))
}

func getenv(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		val = def
	}
	return val
}
