package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dictyBase/graphql-server/internal/app/server"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "graphql-server"
	app.Usage = "cli for graphql-server backend"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-format",
			Usage: "format of the logging out, either of json or text.",
			Value: "json",
		},
		cli.StringFlag{
			Name:  "log-level",
			Usage: "log level for the application",
			Value: "error",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "start-server",
			Usage:  "starts the graphql-server backend",
			Action: server.RunGraphQLServer,
			Flags:  getServerFlags(),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error in running command %s", err)
	}
}

func getServerFlags() []cli.Flag {
	var f []cli.Flag
	f = append(f, redisFlags()...)
	f = append(f, dscFlags()...)
	f = append(f, nonGRPCFlags()...)
	f = append(f, allowedOriginFlags()...)
	f = append(f, userFlags()...)
	f = append(f, authFlags()...)
	f = append(f, contentFlags()...)

	return f
}

func userFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     "user-grpc-host",
			EnvVar:   "USER_API_SERVICE_HOST",
			Usage:    "user grpc host",
			Required: true,
		},
		cli.StringFlag{
			Name:     "user-grpc-port",
			EnvVar:   "USER_API_SERVICE_PORT",
			Usage:    "user grpc port",
			Required: true,
		},
		cli.StringFlag{
			Name:     "role-grpc-host",
			EnvVar:   "ROLE_API_SERVICE_HOST",
			Usage:    "role grpc host",
			Required: true,
		},
		cli.StringFlag{
			Name:     "role-grpc-port",
			EnvVar:   "ROLE_API_SERVICE_PORT",
			Usage:    "role grpc port",
			Required: true,
		},
		cli.StringFlag{
			Name:     "permission-grpc-host",
			EnvVar:   "PERMISSION_API_SERVICE_HOST",
			Usage:    "permission grpc host",
			Required: true,
		},
		cli.StringFlag{
			Name:     "permission-grpc-port",
			EnvVar:   "PERMISSION_API_SERVICE_PORT",
			Usage:    "permission grpc port",
			Required: true,
		},
	}
}

func redisFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "redis-master-service-host",
			EnvVar: "REDIS_SERVICE_HOST",
			Usage:  "redis grpc host",
		},
		cli.StringFlag{
			Name:   "redis-master-service-port",
			EnvVar: "REDIS_SERVICE_PORT",
			Usage:  "redis grpc port",
		},
	}
}

func authFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     "auth-api-endpoint",
			Usage:    "base http url of logto authentication api endpoint",
			EnvVar:   "AUTH_ENDPOINT",
			Required: true,
		},
		cli.StringFlag{
			Name:     "app-id",
			Usage:    "api identifier",
			EnvVar:   "APPLICATION_ID",
			Required: true,
		},
		cli.StringFlag{
			Name:  "api-resource",
			Usage: "http url that represents the identity of the resource",
			Value: "https://default.logto.app/api",
		},
		cli.StringFlag{
			Name:     "app-secret",
			Usage:    "secret to access the authentication api",
			EnvVar:   "APPLICATION_SECRET",
			Required: true,
		},
		cli.StringFlag{
			Name:     "jwks-uri",
			Usage:    "url to retrieve JWK public key set",
			EnvVar:   "JWKS_PUBLIC_URI",
			Required: true,
		},
		cli.StringFlag{
			Name:     "jwt-issuer",
			Usage:    "expected jwt issuer of the token",
			EnvVar:   "JWT_ISSUER",
			Required: true,
		},
		cli.StringFlag{
			Name:     "jwt-audience",
			Usage:    "expect jwt audience of the token",
			EnvVar:   "JWT_AUDIENCE",
			Required: true,
		},
	}
}

func dscFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     "stock-grpc-host",
			EnvVar:   "STOCK_API_SERVICE_HOST",
			Usage:    "stock grpc host",
			Required: true,
		},
		cli.StringFlag{
			Name:     "stock-grpc-port",
			EnvVar:   "STOCK_API_SERVICE_PORT",
			Usage:    "stock grpc port",
			Required: true,
		},
		cli.StringFlag{
			Name:     "order-grpc-host",
			EnvVar:   "ORDER_API_SERVICE_HOST",
			Usage:    "order grpc host",
			Required: true,
		},
		cli.StringFlag{
			Name:     "order-grpc-port",
			EnvVar:   "ORDER_API_SERVICE_PORT",
			Usage:    "order grpc port",
			Required: true,
		},
		cli.StringFlag{
			Name:     "annotation-grpc-host",
			EnvVar:   "ANNOTATION_API_SERVICE_HOST",
			Usage:    "annotation grpc host",
			Required: true,
		},
		cli.StringFlag{
			Name:     "annotation-grpc-port",
			EnvVar:   "ANNOTATION_API_SERVICE_PORT",
			Usage:    "annotation grpc port",
			Required: true,
		},
	}
}

func contentFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     "content-grpc-host",
			EnvVar:   "CONTENT_API_SERVICE_HOST",
			Usage:    "content grpc host",
			Required: true,
		},
		cli.StringFlag{
			Name:     "content-grpc-port",
			EnvVar:   "CONTENT_API_SERVICE_PORT",
			Usage:    "content grpc port",
			Required: true,
		},
	}
}

func nonGRPCFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     "publication-api, pub",
			EnvVar:   "PUBLICATION_API_ENDPOINT",
			Usage:    "publication api endpoint",
			Required: true,
		},
		cli.StringFlag{
			Name:   "organism-api, org",
			EnvVar: "ORGANISM_API_ENDPOINT",
			Usage:  "json endpoint for organisms (downloads page)",
			Value: fmt.Sprintf(
				"https://raw.githubusercontent.com/dictyBase/migration-data/%s",
				"master/downloads/organisms-with-citations.staging.json",
			),
		},
	}
}

/*
*

	A list of allowed origins is necessary since server has set
	allow-credentials to true.
	See https://github.com/rs/cors/issues/55
*/
func allowedOriginFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringSliceFlag{
			Name:  "allowed-origin",
			Usage: "allowed origins for CORS",
			Value: &cli.StringSlice{},
		},
	}
}
