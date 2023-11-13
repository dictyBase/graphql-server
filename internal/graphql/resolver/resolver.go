//go:generate go run github.com/99designs/gqlgen generate
package resolver

import (
	"github.com/dictyBase/graphql-server/internal/graphql/dataloader"
	"github.com/dictyBase/graphql-server/internal/graphql/generated"
	"github.com/dictyBase/graphql-server/internal/graphql/resolver/auth"
	"github.com/dictyBase/graphql-server/internal/graphql/resolver/content"
	"github.com/dictyBase/graphql-server/internal/graphql/resolver/gene"
	"github.com/dictyBase/graphql-server/internal/graphql/resolver/order"
	"github.com/dictyBase/graphql-server/internal/graphql/resolver/organism"
	"github.com/dictyBase/graphql-server/internal/graphql/resolver/publication"
	"github.com/dictyBase/graphql-server/internal/graphql/resolver/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/resolver/user"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/sirupsen/logrus"
)

type Resolver struct {
	registry.Registry
	Dataloaders dataloader.Retriever
	Logger      *logrus.Entry
}

type MutationResolver struct {
	registry.Registry
	Logger *logrus.Entry
}

type QueryResolver struct {
	registry.Registry
	Dataloaders dataloader.Retriever
	Logger      *logrus.Entry
}

func NewResolver(
	nr registry.Registry,
	dl dataloader.Retriever,
	l *logrus.Entry,
) *Resolver {
	return &Resolver{Registry: nr, Dataloaders: dl, Logger: l}
}

func (rrs *Resolver) Mutation() generated.MutationResolver {
	return &MutationResolver{
		Registry: rrs.Registry,
		Logger:   rrs.Logger,
	}
}
func (rrs *Resolver) Query() generated.QueryResolver {
	return &QueryResolver{
		Registry:    rrs.Registry,
		Dataloaders: rrs.Dataloaders,
		Logger:      rrs.Logger,
	}
}
func (rrs *Resolver) User() generated.UserResolver {
	return &user.UserResolver{
		Client: rrs.GetAuthClient(registry.AUTH),
		Logger: rrs.Logger,
	}
}
func (rrs *Resolver) Role() generated.RoleResolver {
	return &user.RoleResolver{
		Client: rrs.GetAuthClient(registry.AUTH),
		Logger: rrs.Logger,
	}
}
func (rrs *Resolver) Permission() generated.PermissionResolver {
	return &user.PermissionResolver{
		Client: rrs.GetPermissionClient(registry.PERMISSION),
		Logger: rrs.Logger,
	}
}
func (rrs *Resolver) Publication() generated.PublicationResolver {
	return &publication.PublicationResolver{
		Logger: rrs.Logger,
	}
}
func (rrs *Resolver) Author() generated.AuthorResolver {
	return &publication.AuthorResolver{
		Logger: rrs.Logger,
	}
}
func (rrs *Resolver) Strain() generated.StrainResolver {
	return &stock.StrainResolver{
		Client:           rrs.GetStockClient(registry.STOCK),
		UserClient:       rrs.GetUserClient(registry.USER),
		AnnotationClient: rrs.GetAnnotationClient(registry.ANNOTATION),
		Registry:         rrs.Registry,
		Logger:           rrs.Logger,
	}
}
func (rrs *Resolver) Plasmid() generated.PlasmidResolver {
	return &stock.PlasmidResolver{
		Client:           rrs.GetStockClient(registry.STOCK),
		UserClient:       rrs.GetUserClient(registry.USER),
		AnnotationClient: rrs.GetAnnotationClient(registry.ANNOTATION),
		Registry:         rrs.Registry,
		Logger:           rrs.Logger,
	}
}

func (rrs *Resolver) Order() generated.OrderResolver {
	return &order.OrderResolver{
		Client:      rrs.GetOrderClient(registry.ORDER),
		StockClient: rrs.GetStockClient(registry.STOCK),
		UserClient:  rrs.GetUserClient(registry.USER),
		Logger:      rrs.Logger,
	}
}

func (rrs *Resolver) Content() generated.ContentResolver {
	return &content.ContentResolver{
		Client:     rrs.GetContentClient(registry.CONTENT),
		UserClient: rrs.GetAuthClient(registry.AUTH),
		Logger:     rrs.Logger,
	}
}

func (rrs *Resolver) Auth() generated.AuthResolver {
	return &auth.AuthResolver{}
}

func (rrs *Resolver) Gene() generated.GeneResolver {
	return &gene.GeneResolver{
		Redis:  rrs.GetRedisRepository("redis"),
		Logger: rrs.Logger,
	}
}

func (rrs *Resolver) Organism() generated.OrganismResolver {
	return &organism.OrganismResolver{
		Logger: rrs.Logger,
	}
}
