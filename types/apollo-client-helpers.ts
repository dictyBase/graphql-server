import { FieldPolicy, FieldReadFunction, TypePolicies, TypePolicy } from '@apollo/client/cache';
export type AuthKeySpecifier = ('identity' | 'token' | 'user' | AuthKeySpecifier)[];
export type AuthFieldPolicy = {
	identity?: FieldPolicy<any> | FieldReadFunction<any>,
	token?: FieldPolicy<any> | FieldReadFunction<any>,
	user?: FieldPolicy<any> | FieldReadFunction<any>
};
export type AuthorKeySpecifier = ('first_name' | 'initials' | 'last_name' | 'rank' | AuthorKeySpecifier)[];
export type AuthorFieldPolicy = {
	first_name?: FieldPolicy<any> | FieldReadFunction<any>,
	initials?: FieldPolicy<any> | FieldReadFunction<any>,
	last_name?: FieldPolicy<any> | FieldReadFunction<any>,
	rank?: FieldPolicy<any> | FieldReadFunction<any>
};
export type BasePublicationKeySpecifier = ('abstract' | 'authors' | 'doi' | 'id' | 'issn' | 'issue' | 'journal' | 'pages' | 'pub_date' | 'pub_type' | 'source' | 'status' | 'title' | 'volume' | BasePublicationKeySpecifier)[];
export type BasePublicationFieldPolicy = {
	abstract?: FieldPolicy<any> | FieldReadFunction<any>,
	authors?: FieldPolicy<any> | FieldReadFunction<any>,
	doi?: FieldPolicy<any> | FieldReadFunction<any>,
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	issn?: FieldPolicy<any> | FieldReadFunction<any>,
	issue?: FieldPolicy<any> | FieldReadFunction<any>,
	journal?: FieldPolicy<any> | FieldReadFunction<any>,
	pages?: FieldPolicy<any> | FieldReadFunction<any>,
	pub_date?: FieldPolicy<any> | FieldReadFunction<any>,
	pub_type?: FieldPolicy<any> | FieldReadFunction<any>,
	source?: FieldPolicy<any> | FieldReadFunction<any>,
	status?: FieldPolicy<any> | FieldReadFunction<any>,
	title?: FieldPolicy<any> | FieldReadFunction<any>,
	volume?: FieldPolicy<any> | FieldReadFunction<any>
};
export type CitationKeySpecifier = ('authors' | 'journal' | 'pubmed_id' | 'title' | CitationKeySpecifier)[];
export type CitationFieldPolicy = {
	authors?: FieldPolicy<any> | FieldReadFunction<any>,
	journal?: FieldPolicy<any> | FieldReadFunction<any>,
	pubmed_id?: FieldPolicy<any> | FieldReadFunction<any>,
	title?: FieldPolicy<any> | FieldReadFunction<any>
};
export type ContentKeySpecifier = ('content' | 'created_at' | 'created_by' | 'id' | 'name' | 'namespace' | 'slug' | 'updated_at' | 'updated_by' | ContentKeySpecifier)[];
export type ContentFieldPolicy = {
	content?: FieldPolicy<any> | FieldReadFunction<any>,
	created_at?: FieldPolicy<any> | FieldReadFunction<any>,
	created_by?: FieldPolicy<any> | FieldReadFunction<any>,
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	name?: FieldPolicy<any> | FieldReadFunction<any>,
	namespace?: FieldPolicy<any> | FieldReadFunction<any>,
	slug?: FieldPolicy<any> | FieldReadFunction<any>,
	updated_at?: FieldPolicy<any> | FieldReadFunction<any>,
	updated_by?: FieldPolicy<any> | FieldReadFunction<any>
};
export type DeleteContentKeySpecifier = ('success' | DeleteContentKeySpecifier)[];
export type DeleteContentFieldPolicy = {
	success?: FieldPolicy<any> | FieldReadFunction<any>
};
export type DeletePermissionKeySpecifier = ('success' | DeletePermissionKeySpecifier)[];
export type DeletePermissionFieldPolicy = {
	success?: FieldPolicy<any> | FieldReadFunction<any>
};
export type DeleteRoleKeySpecifier = ('success' | DeleteRoleKeySpecifier)[];
export type DeleteRoleFieldPolicy = {
	success?: FieldPolicy<any> | FieldReadFunction<any>
};
export type DeleteStockKeySpecifier = ('success' | DeleteStockKeySpecifier)[];
export type DeleteStockFieldPolicy = {
	success?: FieldPolicy<any> | FieldReadFunction<any>
};
export type DeleteUserKeySpecifier = ('success' | DeleteUserKeySpecifier)[];
export type DeleteUserFieldPolicy = {
	success?: FieldPolicy<any> | FieldReadFunction<any>
};
export type DownloadKeySpecifier = ('items' | 'title' | DownloadKeySpecifier)[];
export type DownloadFieldPolicy = {
	items?: FieldPolicy<any> | FieldReadFunction<any>,
	title?: FieldPolicy<any> | FieldReadFunction<any>
};
export type DownloadItemKeySpecifier = ('title' | 'url' | DownloadItemKeySpecifier)[];
export type DownloadItemFieldPolicy = {
	title?: FieldPolicy<any> | FieldReadFunction<any>,
	url?: FieldPolicy<any> | FieldReadFunction<any>
};
export type ExtensionKeySpecifier = ('db' | 'id' | 'name' | 'relation' | ExtensionKeySpecifier)[];
export type ExtensionFieldPolicy = {
	db?: FieldPolicy<any> | FieldReadFunction<any>,
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	name?: FieldPolicy<any> | FieldReadFunction<any>,
	relation?: FieldPolicy<any> | FieldReadFunction<any>
};
export type GOAnnotationKeySpecifier = ('assigned_by' | 'date' | 'evidence_code' | 'extensions' | 'go_term' | 'id' | 'publication' | 'qualifier' | 'type' | 'with' | GOAnnotationKeySpecifier)[];
export type GOAnnotationFieldPolicy = {
	assigned_by?: FieldPolicy<any> | FieldReadFunction<any>,
	date?: FieldPolicy<any> | FieldReadFunction<any>,
	evidence_code?: FieldPolicy<any> | FieldReadFunction<any>,
	extensions?: FieldPolicy<any> | FieldReadFunction<any>,
	go_term?: FieldPolicy<any> | FieldReadFunction<any>,
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	publication?: FieldPolicy<any> | FieldReadFunction<any>,
	qualifier?: FieldPolicy<any> | FieldReadFunction<any>,
	type?: FieldPolicy<any> | FieldReadFunction<any>,
	with?: FieldPolicy<any> | FieldReadFunction<any>
};
export type GeneKeySpecifier = ('id' | 'name' | GeneKeySpecifier)[];
export type GeneFieldPolicy = {
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	name?: FieldPolicy<any> | FieldReadFunction<any>
};
export type GeneGeneralInfoKeySpecifier = ('description' | 'gene_product' | 'id' | 'name_description' | 'synonyms' | GeneGeneralInfoKeySpecifier)[];
export type GeneGeneralInfoFieldPolicy = {
	description?: FieldPolicy<any> | FieldReadFunction<any>,
	gene_product?: FieldPolicy<any> | FieldReadFunction<any>,
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	name_description?: FieldPolicy<any> | FieldReadFunction<any>,
	synonyms?: FieldPolicy<any> | FieldReadFunction<any>
};
export type IdentityKeySpecifier = ('created_at' | 'id' | 'identifier' | 'provider' | 'updated_at' | 'user_id' | IdentityKeySpecifier)[];
export type IdentityFieldPolicy = {
	created_at?: FieldPolicy<any> | FieldReadFunction<any>,
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	identifier?: FieldPolicy<any> | FieldReadFunction<any>,
	provider?: FieldPolicy<any> | FieldReadFunction<any>,
	updated_at?: FieldPolicy<any> | FieldReadFunction<any>,
	user_id?: FieldPolicy<any> | FieldReadFunction<any>
};
export type ImageFileKeySpecifier = ('url' | ImageFileKeySpecifier)[];
export type ImageFileFieldPolicy = {
	url?: FieldPolicy<any> | FieldReadFunction<any>
};
export type LogoutKeySpecifier = ('success' | LogoutKeySpecifier)[];
export type LogoutFieldPolicy = {
	success?: FieldPolicy<any> | FieldReadFunction<any>
};
export type MutationKeySpecifier = ('createContent' | 'createOrder' | 'createPermission' | 'createPlasmid' | 'createRole' | 'createRolePermissionRelationship' | 'createStrain' | 'createUser' | 'createUserRoleRelationship' | 'deleteContent' | 'deletePermission' | 'deleteRole' | 'deleteStock' | 'deleteUser' | 'login' | 'logout' | 'updateContent' | 'updateOrder' | 'updatePermission' | 'updatePlasmid' | 'updateRole' | 'updateStrain' | 'updateUser' | 'uploadFile' | MutationKeySpecifier)[];
export type MutationFieldPolicy = {
	createContent?: FieldPolicy<any> | FieldReadFunction<any>,
	createOrder?: FieldPolicy<any> | FieldReadFunction<any>,
	createPermission?: FieldPolicy<any> | FieldReadFunction<any>,
	createPlasmid?: FieldPolicy<any> | FieldReadFunction<any>,
	createRole?: FieldPolicy<any> | FieldReadFunction<any>,
	createRolePermissionRelationship?: FieldPolicy<any> | FieldReadFunction<any>,
	createStrain?: FieldPolicy<any> | FieldReadFunction<any>,
	createUser?: FieldPolicy<any> | FieldReadFunction<any>,
	createUserRoleRelationship?: FieldPolicy<any> | FieldReadFunction<any>,
	deleteContent?: FieldPolicy<any> | FieldReadFunction<any>,
	deletePermission?: FieldPolicy<any> | FieldReadFunction<any>,
	deleteRole?: FieldPolicy<any> | FieldReadFunction<any>,
	deleteStock?: FieldPolicy<any> | FieldReadFunction<any>,
	deleteUser?: FieldPolicy<any> | FieldReadFunction<any>,
	login?: FieldPolicy<any> | FieldReadFunction<any>,
	logout?: FieldPolicy<any> | FieldReadFunction<any>,
	updateContent?: FieldPolicy<any> | FieldReadFunction<any>,
	updateOrder?: FieldPolicy<any> | FieldReadFunction<any>,
	updatePermission?: FieldPolicy<any> | FieldReadFunction<any>,
	updatePlasmid?: FieldPolicy<any> | FieldReadFunction<any>,
	updateRole?: FieldPolicy<any> | FieldReadFunction<any>,
	updateStrain?: FieldPolicy<any> | FieldReadFunction<any>,
	updateUser?: FieldPolicy<any> | FieldReadFunction<any>,
	uploadFile?: FieldPolicy<any> | FieldReadFunction<any>
};
export type NumberOfPublicationsWithGeneKeySpecifier = ('num_pubs' | 'publications' | NumberOfPublicationsWithGeneKeySpecifier)[];
export type NumberOfPublicationsWithGeneFieldPolicy = {
	num_pubs?: FieldPolicy<any> | FieldReadFunction<any>,
	publications?: FieldPolicy<any> | FieldReadFunction<any>
};
export type OrderKeySpecifier = ('comments' | 'consumer' | 'courier' | 'courier_account' | 'created_at' | 'id' | 'items' | 'payer' | 'payment' | 'purchase_order_num' | 'purchaser' | 'status' | 'updated_at' | OrderKeySpecifier)[];
export type OrderFieldPolicy = {
	comments?: FieldPolicy<any> | FieldReadFunction<any>,
	consumer?: FieldPolicy<any> | FieldReadFunction<any>,
	courier?: FieldPolicy<any> | FieldReadFunction<any>,
	courier_account?: FieldPolicy<any> | FieldReadFunction<any>,
	created_at?: FieldPolicy<any> | FieldReadFunction<any>,
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	items?: FieldPolicy<any> | FieldReadFunction<any>,
	payer?: FieldPolicy<any> | FieldReadFunction<any>,
	payment?: FieldPolicy<any> | FieldReadFunction<any>,
	purchase_order_num?: FieldPolicy<any> | FieldReadFunction<any>,
	purchaser?: FieldPolicy<any> | FieldReadFunction<any>,
	status?: FieldPolicy<any> | FieldReadFunction<any>,
	updated_at?: FieldPolicy<any> | FieldReadFunction<any>
};
export type OrderListWithCursorKeySpecifier = ('limit' | 'nextCursor' | 'orders' | 'previousCursor' | 'totalCount' | OrderListWithCursorKeySpecifier)[];
export type OrderListWithCursorFieldPolicy = {
	limit?: FieldPolicy<any> | FieldReadFunction<any>,
	nextCursor?: FieldPolicy<any> | FieldReadFunction<any>,
	orders?: FieldPolicy<any> | FieldReadFunction<any>,
	previousCursor?: FieldPolicy<any> | FieldReadFunction<any>,
	totalCount?: FieldPolicy<any> | FieldReadFunction<any>
};
export type OrganismKeySpecifier = ('citations' | 'downloads' | 'scientific_name' | 'taxon_id' | OrganismKeySpecifier)[];
export type OrganismFieldPolicy = {
	citations?: FieldPolicy<any> | FieldReadFunction<any>,
	downloads?: FieldPolicy<any> | FieldReadFunction<any>,
	scientific_name?: FieldPolicy<any> | FieldReadFunction<any>,
	taxon_id?: FieldPolicy<any> | FieldReadFunction<any>
};
export type PermissionKeySpecifier = ('created_at' | 'description' | 'id' | 'permission' | 'resource' | 'updated_at' | PermissionKeySpecifier)[];
export type PermissionFieldPolicy = {
	created_at?: FieldPolicy<any> | FieldReadFunction<any>,
	description?: FieldPolicy<any> | FieldReadFunction<any>,
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	permission?: FieldPolicy<any> | FieldReadFunction<any>,
	resource?: FieldPolicy<any> | FieldReadFunction<any>,
	updated_at?: FieldPolicy<any> | FieldReadFunction<any>
};
export type PhenotypeKeySpecifier = ('assay' | 'environment' | 'note' | 'phenotype' | 'publication' | PhenotypeKeySpecifier)[];
export type PhenotypeFieldPolicy = {
	assay?: FieldPolicy<any> | FieldReadFunction<any>,
	environment?: FieldPolicy<any> | FieldReadFunction<any>,
	note?: FieldPolicy<any> | FieldReadFunction<any>,
	phenotype?: FieldPolicy<any> | FieldReadFunction<any>,
	publication?: FieldPolicy<any> | FieldReadFunction<any>
};
export type PlasmidKeySpecifier = ('created_at' | 'created_by' | 'dbxrefs' | 'depositor' | 'editable_summary' | 'genbank_accession' | 'genes' | 'id' | 'image_map' | 'in_stock' | 'keywords' | 'name' | 'publications' | 'sequence' | 'summary' | 'updated_at' | 'updated_by' | PlasmidKeySpecifier)[];
export type PlasmidFieldPolicy = {
	created_at?: FieldPolicy<any> | FieldReadFunction<any>,
	created_by?: FieldPolicy<any> | FieldReadFunction<any>,
	dbxrefs?: FieldPolicy<any> | FieldReadFunction<any>,
	depositor?: FieldPolicy<any> | FieldReadFunction<any>,
	editable_summary?: FieldPolicy<any> | FieldReadFunction<any>,
	genbank_accession?: FieldPolicy<any> | FieldReadFunction<any>,
	genes?: FieldPolicy<any> | FieldReadFunction<any>,
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	image_map?: FieldPolicy<any> | FieldReadFunction<any>,
	in_stock?: FieldPolicy<any> | FieldReadFunction<any>,
	keywords?: FieldPolicy<any> | FieldReadFunction<any>,
	name?: FieldPolicy<any> | FieldReadFunction<any>,
	publications?: FieldPolicy<any> | FieldReadFunction<any>,
	sequence?: FieldPolicy<any> | FieldReadFunction<any>,
	summary?: FieldPolicy<any> | FieldReadFunction<any>,
	updated_at?: FieldPolicy<any> | FieldReadFunction<any>,
	updated_by?: FieldPolicy<any> | FieldReadFunction<any>
};
export type PlasmidListWithCursorKeySpecifier = ('limit' | 'nextCursor' | 'plasmids' | 'previousCursor' | 'totalCount' | PlasmidListWithCursorKeySpecifier)[];
export type PlasmidListWithCursorFieldPolicy = {
	limit?: FieldPolicy<any> | FieldReadFunction<any>,
	nextCursor?: FieldPolicy<any> | FieldReadFunction<any>,
	plasmids?: FieldPolicy<any> | FieldReadFunction<any>,
	previousCursor?: FieldPolicy<any> | FieldReadFunction<any>,
	totalCount?: FieldPolicy<any> | FieldReadFunction<any>
};
export type PublicationKeySpecifier = ('abstract' | 'authors' | 'doi' | 'id' | 'issn' | 'issue' | 'journal' | 'pages' | 'pub_date' | 'pub_type' | 'source' | 'status' | 'title' | 'volume' | PublicationKeySpecifier)[];
export type PublicationFieldPolicy = {
	abstract?: FieldPolicy<any> | FieldReadFunction<any>,
	authors?: FieldPolicy<any> | FieldReadFunction<any>,
	doi?: FieldPolicy<any> | FieldReadFunction<any>,
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	issn?: FieldPolicy<any> | FieldReadFunction<any>,
	issue?: FieldPolicy<any> | FieldReadFunction<any>,
	journal?: FieldPolicy<any> | FieldReadFunction<any>,
	pages?: FieldPolicy<any> | FieldReadFunction<any>,
	pub_date?: FieldPolicy<any> | FieldReadFunction<any>,
	pub_type?: FieldPolicy<any> | FieldReadFunction<any>,
	source?: FieldPolicy<any> | FieldReadFunction<any>,
	status?: FieldPolicy<any> | FieldReadFunction<any>,
	title?: FieldPolicy<any> | FieldReadFunction<any>,
	volume?: FieldPolicy<any> | FieldReadFunction<any>
};
export type PublicationWithGeneKeySpecifier = ('abstract' | 'authors' | 'doi' | 'id' | 'issn' | 'issue' | 'journal' | 'pages' | 'pub_date' | 'pub_type' | 'related_genes' | 'source' | 'status' | 'title' | 'volume' | PublicationWithGeneKeySpecifier)[];
export type PublicationWithGeneFieldPolicy = {
	abstract?: FieldPolicy<any> | FieldReadFunction<any>,
	authors?: FieldPolicy<any> | FieldReadFunction<any>,
	doi?: FieldPolicy<any> | FieldReadFunction<any>,
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	issn?: FieldPolicy<any> | FieldReadFunction<any>,
	issue?: FieldPolicy<any> | FieldReadFunction<any>,
	journal?: FieldPolicy<any> | FieldReadFunction<any>,
	pages?: FieldPolicy<any> | FieldReadFunction<any>,
	pub_date?: FieldPolicy<any> | FieldReadFunction<any>,
	pub_type?: FieldPolicy<any> | FieldReadFunction<any>,
	related_genes?: FieldPolicy<any> | FieldReadFunction<any>,
	source?: FieldPolicy<any> | FieldReadFunction<any>,
	status?: FieldPolicy<any> | FieldReadFunction<any>,
	title?: FieldPolicy<any> | FieldReadFunction<any>,
	volume?: FieldPolicy<any> | FieldReadFunction<any>
};
export type QueryKeySpecifier = ('content' | 'contentBySlug' | 'geneGeneralInformation' | 'geneOntologyAnnotation' | 'listContentByNamespace' | 'listOrders' | 'listOrganisms' | 'listPermissions' | 'listPlasmids' | 'listPlasmidsWithAnnotation' | 'listPublicationsWithGene' | 'listRecentPlasmids' | 'listRecentPublications' | 'listRecentStrains' | 'listRoles' | 'listStrains' | 'listStrainsWithAnnotation' | 'listStrainsWithGene' | 'listUsers' | 'order' | 'organism' | 'permission' | 'plasmid' | 'publication' | 'role' | 'strain' | 'user' | 'userByEmail' | QueryKeySpecifier)[];
export type QueryFieldPolicy = {
	content?: FieldPolicy<any> | FieldReadFunction<any>,
	contentBySlug?: FieldPolicy<any> | FieldReadFunction<any>,
	geneGeneralInformation?: FieldPolicy<any> | FieldReadFunction<any>,
	geneOntologyAnnotation?: FieldPolicy<any> | FieldReadFunction<any>,
	listContentByNamespace?: FieldPolicy<any> | FieldReadFunction<any>,
	listOrders?: FieldPolicy<any> | FieldReadFunction<any>,
	listOrganisms?: FieldPolicy<any> | FieldReadFunction<any>,
	listPermissions?: FieldPolicy<any> | FieldReadFunction<any>,
	listPlasmids?: FieldPolicy<any> | FieldReadFunction<any>,
	listPlasmidsWithAnnotation?: FieldPolicy<any> | FieldReadFunction<any>,
	listPublicationsWithGene?: FieldPolicy<any> | FieldReadFunction<any>,
	listRecentPlasmids?: FieldPolicy<any> | FieldReadFunction<any>,
	listRecentPublications?: FieldPolicy<any> | FieldReadFunction<any>,
	listRecentStrains?: FieldPolicy<any> | FieldReadFunction<any>,
	listRoles?: FieldPolicy<any> | FieldReadFunction<any>,
	listStrains?: FieldPolicy<any> | FieldReadFunction<any>,
	listStrainsWithAnnotation?: FieldPolicy<any> | FieldReadFunction<any>,
	listStrainsWithGene?: FieldPolicy<any> | FieldReadFunction<any>,
	listUsers?: FieldPolicy<any> | FieldReadFunction<any>,
	order?: FieldPolicy<any> | FieldReadFunction<any>,
	organism?: FieldPolicy<any> | FieldReadFunction<any>,
	permission?: FieldPolicy<any> | FieldReadFunction<any>,
	plasmid?: FieldPolicy<any> | FieldReadFunction<any>,
	publication?: FieldPolicy<any> | FieldReadFunction<any>,
	role?: FieldPolicy<any> | FieldReadFunction<any>,
	strain?: FieldPolicy<any> | FieldReadFunction<any>,
	user?: FieldPolicy<any> | FieldReadFunction<any>,
	userByEmail?: FieldPolicy<any> | FieldReadFunction<any>
};
export type RoleKeySpecifier = ('created_at' | 'description' | 'id' | 'permissions' | 'role' | 'updated_at' | RoleKeySpecifier)[];
export type RoleFieldPolicy = {
	created_at?: FieldPolicy<any> | FieldReadFunction<any>,
	description?: FieldPolicy<any> | FieldReadFunction<any>,
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	permissions?: FieldPolicy<any> | FieldReadFunction<any>,
	role?: FieldPolicy<any> | FieldReadFunction<any>,
	updated_at?: FieldPolicy<any> | FieldReadFunction<any>
};
export type StockKeySpecifier = ('created_at' | 'created_by' | 'dbxrefs' | 'depositor' | 'editable_summary' | 'genes' | 'id' | 'in_stock' | 'publications' | 'summary' | 'updated_at' | 'updated_by' | StockKeySpecifier)[];
export type StockFieldPolicy = {
	created_at?: FieldPolicy<any> | FieldReadFunction<any>,
	created_by?: FieldPolicy<any> | FieldReadFunction<any>,
	dbxrefs?: FieldPolicy<any> | FieldReadFunction<any>,
	depositor?: FieldPolicy<any> | FieldReadFunction<any>,
	editable_summary?: FieldPolicy<any> | FieldReadFunction<any>,
	genes?: FieldPolicy<any> | FieldReadFunction<any>,
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	in_stock?: FieldPolicy<any> | FieldReadFunction<any>,
	publications?: FieldPolicy<any> | FieldReadFunction<any>,
	summary?: FieldPolicy<any> | FieldReadFunction<any>,
	updated_at?: FieldPolicy<any> | FieldReadFunction<any>,
	updated_by?: FieldPolicy<any> | FieldReadFunction<any>
};
export type StrainKeySpecifier = ('characteristics' | 'created_at' | 'created_by' | 'dbxrefs' | 'depositor' | 'editable_summary' | 'genes' | 'genetic_modification' | 'genotypes' | 'id' | 'in_stock' | 'label' | 'mutagenesis_method' | 'names' | 'parent' | 'phenotypes' | 'plasmid' | 'publications' | 'species' | 'summary' | 'systematic_name' | 'updated_at' | 'updated_by' | StrainKeySpecifier)[];
export type StrainFieldPolicy = {
	characteristics?: FieldPolicy<any> | FieldReadFunction<any>,
	created_at?: FieldPolicy<any> | FieldReadFunction<any>,
	created_by?: FieldPolicy<any> | FieldReadFunction<any>,
	dbxrefs?: FieldPolicy<any> | FieldReadFunction<any>,
	depositor?: FieldPolicy<any> | FieldReadFunction<any>,
	editable_summary?: FieldPolicy<any> | FieldReadFunction<any>,
	genes?: FieldPolicy<any> | FieldReadFunction<any>,
	genetic_modification?: FieldPolicy<any> | FieldReadFunction<any>,
	genotypes?: FieldPolicy<any> | FieldReadFunction<any>,
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	in_stock?: FieldPolicy<any> | FieldReadFunction<any>,
	label?: FieldPolicy<any> | FieldReadFunction<any>,
	mutagenesis_method?: FieldPolicy<any> | FieldReadFunction<any>,
	names?: FieldPolicy<any> | FieldReadFunction<any>,
	parent?: FieldPolicy<any> | FieldReadFunction<any>,
	phenotypes?: FieldPolicy<any> | FieldReadFunction<any>,
	plasmid?: FieldPolicy<any> | FieldReadFunction<any>,
	publications?: FieldPolicy<any> | FieldReadFunction<any>,
	species?: FieldPolicy<any> | FieldReadFunction<any>,
	summary?: FieldPolicy<any> | FieldReadFunction<any>,
	systematic_name?: FieldPolicy<any> | FieldReadFunction<any>,
	updated_at?: FieldPolicy<any> | FieldReadFunction<any>,
	updated_by?: FieldPolicy<any> | FieldReadFunction<any>
};
export type StrainListWithCursorKeySpecifier = ('limit' | 'nextCursor' | 'previousCursor' | 'strains' | 'totalCount' | StrainListWithCursorKeySpecifier)[];
export type StrainListWithCursorFieldPolicy = {
	limit?: FieldPolicy<any> | FieldReadFunction<any>,
	nextCursor?: FieldPolicy<any> | FieldReadFunction<any>,
	previousCursor?: FieldPolicy<any> | FieldReadFunction<any>,
	strains?: FieldPolicy<any> | FieldReadFunction<any>,
	totalCount?: FieldPolicy<any> | FieldReadFunction<any>
};
export type UserKeySpecifier = ('city' | 'country' | 'created_at' | 'email' | 'first_address' | 'first_name' | 'group_name' | 'id' | 'is_active' | 'last_name' | 'organization' | 'phone' | 'roles' | 'second_address' | 'state' | 'updated_at' | 'zipcode' | UserKeySpecifier)[];
export type UserFieldPolicy = {
	city?: FieldPolicy<any> | FieldReadFunction<any>,
	country?: FieldPolicy<any> | FieldReadFunction<any>,
	created_at?: FieldPolicy<any> | FieldReadFunction<any>,
	email?: FieldPolicy<any> | FieldReadFunction<any>,
	first_address?: FieldPolicy<any> | FieldReadFunction<any>,
	first_name?: FieldPolicy<any> | FieldReadFunction<any>,
	group_name?: FieldPolicy<any> | FieldReadFunction<any>,
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	is_active?: FieldPolicy<any> | FieldReadFunction<any>,
	last_name?: FieldPolicy<any> | FieldReadFunction<any>,
	organization?: FieldPolicy<any> | FieldReadFunction<any>,
	phone?: FieldPolicy<any> | FieldReadFunction<any>,
	roles?: FieldPolicy<any> | FieldReadFunction<any>,
	second_address?: FieldPolicy<any> | FieldReadFunction<any>,
	state?: FieldPolicy<any> | FieldReadFunction<any>,
	updated_at?: FieldPolicy<any> | FieldReadFunction<any>,
	zipcode?: FieldPolicy<any> | FieldReadFunction<any>
};
export type UserListKeySpecifier = ('pageNum' | 'pageSize' | 'totalCount' | 'users' | UserListKeySpecifier)[];
export type UserListFieldPolicy = {
	pageNum?: FieldPolicy<any> | FieldReadFunction<any>,
	pageSize?: FieldPolicy<any> | FieldReadFunction<any>,
	totalCount?: FieldPolicy<any> | FieldReadFunction<any>,
	users?: FieldPolicy<any> | FieldReadFunction<any>
};
export type WithKeySpecifier = ('db' | 'id' | 'name' | WithKeySpecifier)[];
export type WithFieldPolicy = {
	db?: FieldPolicy<any> | FieldReadFunction<any>,
	id?: FieldPolicy<any> | FieldReadFunction<any>,
	name?: FieldPolicy<any> | FieldReadFunction<any>
};
export type StrictTypedTypePolicies = {
	Auth?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | AuthKeySpecifier | (() => undefined | AuthKeySpecifier),
		fields?: AuthFieldPolicy,
	},
	Author?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | AuthorKeySpecifier | (() => undefined | AuthorKeySpecifier),
		fields?: AuthorFieldPolicy,
	},
	BasePublication?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | BasePublicationKeySpecifier | (() => undefined | BasePublicationKeySpecifier),
		fields?: BasePublicationFieldPolicy,
	},
	Citation?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | CitationKeySpecifier | (() => undefined | CitationKeySpecifier),
		fields?: CitationFieldPolicy,
	},
	Content?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | ContentKeySpecifier | (() => undefined | ContentKeySpecifier),
		fields?: ContentFieldPolicy,
	},
	DeleteContent?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | DeleteContentKeySpecifier | (() => undefined | DeleteContentKeySpecifier),
		fields?: DeleteContentFieldPolicy,
	},
	DeletePermission?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | DeletePermissionKeySpecifier | (() => undefined | DeletePermissionKeySpecifier),
		fields?: DeletePermissionFieldPolicy,
	},
	DeleteRole?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | DeleteRoleKeySpecifier | (() => undefined | DeleteRoleKeySpecifier),
		fields?: DeleteRoleFieldPolicy,
	},
	DeleteStock?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | DeleteStockKeySpecifier | (() => undefined | DeleteStockKeySpecifier),
		fields?: DeleteStockFieldPolicy,
	},
	DeleteUser?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | DeleteUserKeySpecifier | (() => undefined | DeleteUserKeySpecifier),
		fields?: DeleteUserFieldPolicy,
	},
	Download?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | DownloadKeySpecifier | (() => undefined | DownloadKeySpecifier),
		fields?: DownloadFieldPolicy,
	},
	DownloadItem?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | DownloadItemKeySpecifier | (() => undefined | DownloadItemKeySpecifier),
		fields?: DownloadItemFieldPolicy,
	},
	Extension?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | ExtensionKeySpecifier | (() => undefined | ExtensionKeySpecifier),
		fields?: ExtensionFieldPolicy,
	},
	GOAnnotation?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | GOAnnotationKeySpecifier | (() => undefined | GOAnnotationKeySpecifier),
		fields?: GOAnnotationFieldPolicy,
	},
	Gene?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | GeneKeySpecifier | (() => undefined | GeneKeySpecifier),
		fields?: GeneFieldPolicy,
	},
	GeneGeneralInfo?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | GeneGeneralInfoKeySpecifier | (() => undefined | GeneGeneralInfoKeySpecifier),
		fields?: GeneGeneralInfoFieldPolicy,
	},
	Identity?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | IdentityKeySpecifier | (() => undefined | IdentityKeySpecifier),
		fields?: IdentityFieldPolicy,
	},
	ImageFile?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | ImageFileKeySpecifier | (() => undefined | ImageFileKeySpecifier),
		fields?: ImageFileFieldPolicy,
	},
	Logout?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | LogoutKeySpecifier | (() => undefined | LogoutKeySpecifier),
		fields?: LogoutFieldPolicy,
	},
	Mutation?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | MutationKeySpecifier | (() => undefined | MutationKeySpecifier),
		fields?: MutationFieldPolicy,
	},
	NumberOfPublicationsWithGene?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | NumberOfPublicationsWithGeneKeySpecifier | (() => undefined | NumberOfPublicationsWithGeneKeySpecifier),
		fields?: NumberOfPublicationsWithGeneFieldPolicy,
	},
	Order?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | OrderKeySpecifier | (() => undefined | OrderKeySpecifier),
		fields?: OrderFieldPolicy,
	},
	OrderListWithCursor?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | OrderListWithCursorKeySpecifier | (() => undefined | OrderListWithCursorKeySpecifier),
		fields?: OrderListWithCursorFieldPolicy,
	},
	Organism?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | OrganismKeySpecifier | (() => undefined | OrganismKeySpecifier),
		fields?: OrganismFieldPolicy,
	},
	Permission?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | PermissionKeySpecifier | (() => undefined | PermissionKeySpecifier),
		fields?: PermissionFieldPolicy,
	},
	Phenotype?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | PhenotypeKeySpecifier | (() => undefined | PhenotypeKeySpecifier),
		fields?: PhenotypeFieldPolicy,
	},
	Plasmid?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | PlasmidKeySpecifier | (() => undefined | PlasmidKeySpecifier),
		fields?: PlasmidFieldPolicy,
	},
	PlasmidListWithCursor?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | PlasmidListWithCursorKeySpecifier | (() => undefined | PlasmidListWithCursorKeySpecifier),
		fields?: PlasmidListWithCursorFieldPolicy,
	},
	Publication?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | PublicationKeySpecifier | (() => undefined | PublicationKeySpecifier),
		fields?: PublicationFieldPolicy,
	},
	PublicationWithGene?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | PublicationWithGeneKeySpecifier | (() => undefined | PublicationWithGeneKeySpecifier),
		fields?: PublicationWithGeneFieldPolicy,
	},
	Query?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | QueryKeySpecifier | (() => undefined | QueryKeySpecifier),
		fields?: QueryFieldPolicy,
	},
	Role?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | RoleKeySpecifier | (() => undefined | RoleKeySpecifier),
		fields?: RoleFieldPolicy,
	},
	Stock?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | StockKeySpecifier | (() => undefined | StockKeySpecifier),
		fields?: StockFieldPolicy,
	},
	Strain?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | StrainKeySpecifier | (() => undefined | StrainKeySpecifier),
		fields?: StrainFieldPolicy,
	},
	StrainListWithCursor?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | StrainListWithCursorKeySpecifier | (() => undefined | StrainListWithCursorKeySpecifier),
		fields?: StrainListWithCursorFieldPolicy,
	},
	User?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | UserKeySpecifier | (() => undefined | UserKeySpecifier),
		fields?: UserFieldPolicy,
	},
	UserList?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | UserListKeySpecifier | (() => undefined | UserListKeySpecifier),
		fields?: UserListFieldPolicy,
	},
	With?: Omit<TypePolicy, "fields" | "keyFields"> & {
		keyFields?: false | WithKeySpecifier | (() => undefined | WithKeySpecifier),
		fields?: WithFieldPolicy,
	}
};
export type TypedTypePolicies = StrictTypedTypePolicies & TypePolicies;