package registry

import (
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-genproto/dictybaseapis/auth"
	"github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/dictyBase/go-genproto/dictybaseapis/identity"
	"github.com/dictyBase/go-genproto/dictybaseapis/order"
	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/repository"
	"github.com/emirpasic/gods/maps/hashmap"
	"google.golang.org/grpc"
)

// constants used with modware-annotation
const (
	PhenoOntology           = "Dicty Phenotypes"
	EnvOntology             = "Dicty Environment"
	AssayOntology           = "Dictyostelium Assay"
	MutagenesisOntology     = "Dd Mutagenesis Method"
	DictyAnnoOntology       = "dicty_annotation"
	StrainCharOnto          = "strain_characteristics"
	StrainInvOnto           = "strain_inventory"
	PlasmidInvOnto          = "plasmid_inventory"
	StrainInvTag            = "strain_inventory"
	PlasmidInvTag           = "plasmid inventory"
	DictyStrainPropOntology = "dicty_strain_property"
	InvLocationTag          = "location"
	GwdiStrainTag           = "REMI-seq"
	GeneralStrainTag        = "general strain"
	BacterialStrainTag      = "bacterial strain"
	LiteratureTag           = "pubmed id"
	NoteTag                 = "public note"
	SysnameTag              = "systematic name"
	MutmethodTag            = "mutagenesis method"
	MuttypeTag              = "mutant type"
	GenoTag                 = "genotype"
	SynTag                  = "synonym"
	EmptyValue              = "novalue"
)

const (
	USER        = "user"
	ROLE        = "role"
	PERMISSION  = "permission"
	PUBLICATION = "publication"
	STOCK       = "stock"
	ORDER       = "order"
	CONTENT     = "content"
	ANNOTATION  = "annotation"
	AUTH        = "auth"
	IDENTITY    = "identity"
	ORGANISM    = "organism"
)

type collection struct {
	connMap *hashmap.Map
}

type Registry interface {
	ServiceMap() map[string]string
	AddAPIEndpoint(key, endpoint string)
	AddAPIConnection(key string, conn *grpc.ClientConn)
	AddRepository(key string, st repository.Repository)
	GetAPIConnection(key string) (conn *grpc.ClientConn)
	GetAPIEndpoint(key string) string
	GetUserClient(key string) user.UserServiceClient
	GetRoleClient(key string) user.RoleServiceClient
	GetPermissionClient(key string) user.PermissionServiceClient
	GetStockClient(key string) stock.StockServiceClient
	GetOrderClient(key string) order.OrderServiceClient
	GetContentClient(key string) content.ContentServiceClient
	GetAnnotationClient(key string) annotation.TaggedAnnotationServiceClient
	GetAuthClient(key string) auth.AuthServiceClient
	GetIdentityClient(key string) identity.IdentityServiceClient
	GetRedisRepository(key string) repository.Repository
}

// NewRegistry constructs a hashmap for our grpc clients
func NewRegistry() Registry {
	m := hashmap.New()
	return &collection{connMap: m}
}

func (coll *collection) ServiceMap() map[string]string {
	return map[string]string{
		"stock":      STOCK,
		"order":      ORDER,
		"annotation": ANNOTATION,
		"role":       ROLE,
		"permission": PERMISSION,
		"user":       USER,
		"content":    CONTENT,
		/*"auth":       AUTH,
		"identity":   IDENTITY, */
	}
}

// AddAPIEndpoint adds a new REST endpoint to the hashmap
func (coll *collection) AddAPIEndpoint(key, endpoint string) {
	coll.connMap.Put(key, endpoint)
}

// AddAPIClient adds a new gRPC client to the hashmap
func (coll *collection) AddAPIConnection(key string, conn *grpc.ClientConn) {
	coll.connMap.Put(key, conn)
}

// AddRepository adds a new repository client to the hashmap
func (coll *collection) AddRepository(key string, st repository.Repository) {
	coll.connMap.Put(key, st)
}

// GetAPIClient looks up a client in the hashmap
func (coll *collection) GetAPIConnection(key string) (conn *grpc.ClientConn) {
	v, ok := coll.connMap.Get(key)
	if !ok {
		panic("could not get grpc client connection")
	}
	conn, _ = v.(*grpc.ClientConn)
	return conn
}

func (coll *collection) GetUserClient(key string) user.UserServiceClient {
	return user.NewUserServiceClient(coll.GetAPIConnection(key))
}

func (coll *collection) GetRoleClient(key string) user.RoleServiceClient {
	return user.NewRoleServiceClient(coll.GetAPIConnection(key))
}

func (coll *collection) GetPermissionClient(
	key string,
) user.PermissionServiceClient {
	return user.NewPermissionServiceClient(coll.GetAPIConnection(key))
}

func (coll *collection) GetStockClient(key string) stock.StockServiceClient {
	return stock.NewStockServiceClient(coll.GetAPIConnection(key))
}

func (coll *collection) GetOrderClient(key string) order.OrderServiceClient {
	return order.NewOrderServiceClient(coll.GetAPIConnection(key))
}

func (coll *collection) GetContentClient(key string) content.ContentServiceClient {
	return content.NewContentServiceClient(coll.GetAPIConnection(key))
}

func (coll *collection) GetAnnotationClient(
	key string,
) annotation.TaggedAnnotationServiceClient {
	return annotation.NewTaggedAnnotationServiceClient(coll.GetAPIConnection(key))
}

func (coll *collection) GetAuthClient(key string) auth.AuthServiceClient {
	return auth.NewAuthServiceClient(coll.GetAPIConnection(key))
}

func (coll *collection) GetIdentityClient(
	key string,
) identity.IdentityServiceClient {
	return identity.NewIdentityServiceClient(coll.GetAPIConnection(key))
}

func (coll *collection) GetAPIEndpoint(key string) string {
	v, _ := coll.connMap.Get(key)
	endpoint, _ := v.(string)
	return endpoint
}

func (coll *collection) GetRedisRepository(key string) repository.Repository {
	v, _ := coll.connMap.Get(key)
	st, _ := v.(repository.Repository)
	return st
}
