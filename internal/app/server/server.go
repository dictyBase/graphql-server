package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dictyBase/graphql-server/internal/graphql/resolver"
	"google.golang.org/grpc"

	"github.com/dictyBase/graphql-server/internal/registry"

	"github.com/99designs/gqlgen/handler"
	"github.com/dictyBase/graphql-server/internal/graphql/generated"
	"github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v1"
)

// RunGraphQLServer starts the GraphQL backend
func RunGraphQLServer(c *cli.Context) error {
	log := getLogger(c)
	// generate new (empty) hashmap
	nr := registry.NewRegistry()
	for k, v := range registry.ServiceMap {
		host := c.String(fmt.Sprintf("%s-grpc-host", k))
		port := c.String(fmt.Sprintf("%s-grpc-port", k))
		// establish grpc connections
		conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("cannot connect to grpc microservice for %s", err),
				2,
			)
		}
		// add api clients to hashmap
		nr.AddAPIConnection(v, conn)
	}
	// verify if publication api endpoint is running
	// need to use Get method here because Head returns 405 status
	res, err := http.Get(c.String("publication-api") + "/" + "30048658")
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("cannot reach publication api endpoint %s", err),
			2,
		)
	}
	if res.StatusCode != http.StatusOK {
		return cli.NewExitError(
			fmt.Sprintf("did not get ok status from publication api endpoint, got %v instead", res.StatusCode),
			2,
		)
	}
	// publication api status is fine, so add it to registry
	nr.AddAPIEndpoint(registry.PUBLICATION, c.String("publication-api"))
	s := resolver.NewResolver(nr, log)

	http.Handle("/", handler.Playground("GraphQL playground", "/graphql"))
	http.Handle("/graphql", handler.GraphQL(generated.NewExecutableSchema(generated.Config{Resolvers: s})))
	log.Debugf("connect to http://localhost:8080/ for GraphQL playground")
	log.Fatal(http.ListenAndServe(":8080", nil))
	return nil
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
	l := c.GlobalString("log-level")
	switch l {
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
