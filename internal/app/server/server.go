package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dictyBase/graphql-server/internal/app/middleware"
	"github.com/dictyBase/graphql-server/internal/graphql/dataloader"
	"github.com/dictyBase/graphql-server/internal/graphql/generated"
	"github.com/dictyBase/graphql-server/internal/graphql/resolver"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/dictyBase/graphql-server/internal/repository/redis"
	"github.com/go-chi/cors"
	"google.golang.org/grpc"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// RunGraphQLServer starts the GraphQL backend
func RunGraphQLServer(cltx *cli.Context) error {
	log := getLogger(cltx)
	router := chi.NewRouter()
	nreg := registry.NewRegistry()
	for k, v := range nreg.ServiceMap() {
		host := cltx.String(fmt.Sprintf("%s-grpc-host", k))
		port := cltx.String(fmt.Sprintf("%s-grpc-port", k))
		// establish grpc connections
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		conn, err := grpc.DialContext(
			ctx,
			fmt.Sprintf("%s:%s", host, port),
			grpc.WithInsecure(),
			grpc.WithBlock(),
		)
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("cannot connect to grpc microservice %s", err),
				2,
			)
		}
		// add api clients to hashmap
		nreg.AddAPIConnection(v, conn)
	}
	/* endpoints := []string{
		c.String("publication-api") + "/" + "30048658",
		c.String("organism-api"),
	}
	// test all api endpoints
	if err := checkEndpoints(endpoints); err != nil {
		return err
	} */
	// apis came back ok, add to registry
	nreg.AddAPIEndpoint(registry.PUBLICATION, cltx.String("publication-api"))
	nreg.AddAPIEndpoint(registry.ORGANISM, cltx.String("organism-api"))
	// add redis to registry
	radd := fmt.Sprintf(
		"%s:%s",
		cltx.String("redis-master-service-host"),
		cltx.String("redis-master-service-port"),
	)
	cache, err := redis.NewCache(radd)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("cannot create redis cache: %v", err),
			2,
		)
	}
	nreg.AddRepository("redis", cache)
	// initialize the dataloaders
	dl := dataloader.NewRetriever()
	s := resolver.NewResolver(nreg, dl, log)
	crs := getCORS(cltx.StringSlice("allowed-origin"))
	router.Use(crs.Handler)
	authMdw, err := middleware.NewJWTAuth(
		cltx.String("jwks-uri"),
		cltx.String("jwt-audience"),
		cltx.String("issuer"),
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in creating jwt auth middleware %s", err),
			2,
		)
	}
	router.Use(authMdw.JwtHandler)
	router.Use(dataloader.DataloaderMiddleware(nreg))
	execSchema := generated.NewExecutableSchema(generated.Config{Resolvers: s})
	srv := handler.NewDefaultServer(execSchema)
	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	router.Handle("/graphql", srv)
	log.Debugf("connect to port 8080 for GraphQL playground")
	log.Fatal(http.ListenAndServe(":8080", router))
	return nil
}

func checkEndpoints(urls []string) error {
	for _, url := range urls {
		res, err := http.Get(url)
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("cannot reach api endpoint %s", err),
				2,
			)
		}
		if res.StatusCode != http.StatusOK {
			return cli.NewExitError(
				fmt.Sprintf(
					"did not get ok status from api endpoint, got %v instead",
					res.StatusCode,
				),
				2,
			)
		}
	}
	return nil
}

func getCORS(origins []string) *cors.Cors {
	aorg := append(origins, "http://localhost:*")
	aorg = append(aorg, "https://dictybase.dev")
	aorg = append(aorg, "https://dictybase.dev/")
	aorg = append(aorg, "https://dictybase.dev*")
	return cors.New(cors.Options{
		AllowedOrigins:   aorg,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})
}

func getLogger(c *cli.Context) *logrus.Entry {
	log := logrus.New()
	log.Out = os.Stderr
	switch c.GlobalString("log-format") {
	case "text":
		log.Formatter = &logrus.TextFormatter{
			TimestampFormat: "02/Jan/2006:15:04:05",
		}
	case "json":
		log.Formatter = &logrus.JSONFormatter{
			TimestampFormat: "02/Jan/2006:15:04:05",
		}
	}
	lvl := c.GlobalString("log-level")
	switch lvl {
	case "debug":
		log.Level = logrus.DebugLevel
	case "warn":
		log.Level = logrus.WarnLevel
	case "error":
		log.Level = logrus.ErrorLevel
	case "fatal":
		log.Level = logrus.FatalLevel
	case "panic":
		log.Level = logrus.PanicLevel
	}
	return logrus.NewEntry(log)
}
