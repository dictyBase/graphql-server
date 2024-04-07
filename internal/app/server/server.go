package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dictyBase/graphql-server/internal/app/middleware"
	"github.com/dictyBase/graphql-server/internal/authentication"
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
	nreg := registry.NewRegistry()
	if err := setupServices(cltx, nreg); err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	log := getLogger(cltx)
	router := chi.NewRouter()

	// initialize the dataloaders, add middlewares, run the server etc...
	dl := dataloader.NewRetriever()
	s := resolver.NewResolver(nreg, dl, log)
	crs := getCORS(cltx.StringSlice("allowed-origin"))
	router.Use(crs.Handler)
	authMdw, err := middleware.NewJWTAuth(
		cltx.String("jwks-uri"),
		cltx.String("jwt-audience"),
		cltx.String("jwt-issuer"),
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
	hsrv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}
	log.Infof(
		"going to start graphql server with jwt-audience %s and jwt-issuer %s",
		cltx.String("jwt-audience"),
		cltx.String("jwt-issuer"),
	)
	log.Fatal(hsrv.ListenAndServe())

	return nil
}

func setupServices(cltx *cli.Context, nreg registry.Registry) error {
	if err := establishGrpcConnnections(cltx, nreg); err != nil {
		return err
	}
	if err := initRedis(cltx, nreg); err != nil {
		return err
	}
	if err := setupS3Client(cltx, nreg); err != nil {
		return err
	}
	addEndpoints(cltx, nreg)

	return nil
}

func establishGrpcConnnections(
	ctx *cli.Context,
	nreg registry.Registry,
) error {
	for key, val := range nreg.ServiceMap() {
		conn, err := connectToGrpcService(ctx, key)
		if err != nil {
			return fmt.Errorf("error in connecting to grpc service %s", key)
		}
		nreg.AddAPIConnection(val, conn)
	}
	return nil
}

func connectToGrpcService(
	ctx *cli.Context,
	key string,
) (*grpc.ClientConn, error) {
	host := ctx.String(fmt.Sprintf("%s-grpc-host", key))
	port := ctx.String(fmt.Sprintf("%s-grpc-port", key))

	grpcCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(
		grpcCtx,
		fmt.Sprintf("%s:%s", host, port),
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		}...,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to grpc microservice %s", err)
	}
	return conn, nil
}

func addEndpoints(ctx *cli.Context, nreg registry.Registry) {
	nreg.AddAPIEndpoint(registry.PUBLICATION, ctx.String("publication-api"))
	nreg.AddAPIEndpoint(registry.ORGANISM, ctx.String("organism-api"))
	nreg.AddAuthClient(
		registry.AUTH,
		authentication.NewClient(&authentication.LogtoClientParams{
			URL:         ctx.String("auth-api-endpoint"),
			AppID:       ctx.String("app-id"),
			AppSecret:   ctx.String("app-secret"),
			APIResource: ctx.String("api-resource"),
			Key:         "AUTHTOKEN",
			TokenCache:  nreg.GetRedisRepository(registry.REDISREPO),
		}),
	)
}

func setupS3Client(ctx *cli.Context, nreg registry.Registry) error {
	client, err := minio.New(
		fmt.Sprintf(
			"%s:%s",
			ctx.String("s3-server"),
			ctx.String("s3-server-port"),
		),
		&minio.Options{
			Creds: credentials.NewStaticV4(
				ctx.String("access-key"),
				ctx.String("secret-key"),
				"",
			),
			Secure: false,
		})
	if err != nil {
		return fmt.Errorf("error in creating minio client %s", err)
	}
	nreg.AddS3Client(registry.S3CLIENT, client)
	return nil
}

func initRedis(ctx *cli.Context, nreg registry.Registry) error {
	radd := fmt.Sprintf(
		"%s:%s",
		ctx.String("redis-master-service-host"),
		ctx.String("redis-master-service-port"),
	)
	cache, err := redis.NewCache(radd)
	if err != nil {
		return fmt.Errorf("cannot create redis cache: %v", err)
	}
	nreg.AddRepository(registry.REDISREPO, cache)
	return nil
}

func getCORS(origins []string) *cors.Cors {
	origins = append(origins, "http://localhost:*")
	origins = append(origins, "https://dictybase.dev")
	origins = append(origins, "https://dictybase.dev/")
	origins = append(origins, "https://dictybase.dev*")
	return cors.New(cors.Options{
		AllowedOrigins:   origins,
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
