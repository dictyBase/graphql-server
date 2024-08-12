import { graphql } from 'msw'
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Timestamp: { input: any; output: any; }
  /** The `Upload` scalar type represents a multipart file upload. */
  Upload: { input: any; output: any; }
};

export type Auth = {
  __typename?: 'Auth';
  identity: Identity;
  token: Scalars['String']['output'];
  user: User;
};

export type Author = {
  __typename?: 'Author';
  first_name?: Maybe<Scalars['String']['output']>;
  initials?: Maybe<Scalars['String']['output']>;
  last_name: Scalars['String']['output'];
  rank?: Maybe<Scalars['String']['output']>;
};

export type BasePublication = {
  abstract: Scalars['String']['output'];
  authors: Array<Author>;
  doi?: Maybe<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  issn?: Maybe<Scalars['String']['output']>;
  issue?: Maybe<Scalars['String']['output']>;
  journal: Scalars['String']['output'];
  pages?: Maybe<Scalars['String']['output']>;
  pub_date?: Maybe<Scalars['Timestamp']['output']>;
  pub_type: Scalars['String']['output'];
  source: Scalars['String']['output'];
  status?: Maybe<Scalars['String']['output']>;
  title: Scalars['String']['output'];
  volume?: Maybe<Scalars['String']['output']>;
};

export type Citation = {
  __typename?: 'Citation';
  authors: Scalars['String']['output'];
  journal: Scalars['String']['output'];
  pubmed_id: Scalars['String']['output'];
  title: Scalars['String']['output'];
};

export type Content = {
  __typename?: 'Content';
  content: Scalars['String']['output'];
  created_at: Scalars['Timestamp']['output'];
  created_by: User;
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  namespace: Scalars['String']['output'];
  slug: Scalars['String']['output'];
  updated_at: Scalars['Timestamp']['output'];
  updated_by: User;
};

export type CreateContentInput = {
  content: Scalars['String']['input'];
  created_by: Scalars['String']['input'];
  name: Scalars['String']['input'];
  namespace: Scalars['String']['input'];
};

export type CreateOrderInput = {
  comments?: InputMaybe<Scalars['String']['input']>;
  consumer: Scalars['String']['input'];
  courier: Scalars['String']['input'];
  courier_account: Scalars['String']['input'];
  items: Array<Scalars['String']['input']>;
  payer: Scalars['String']['input'];
  payment: Scalars['String']['input'];
  purchase_order_num?: InputMaybe<Scalars['String']['input']>;
  purchaser: Scalars['String']['input'];
  status: StatusEnum;
};

export type CreatePermissionInput = {
  description: Scalars['String']['input'];
  permission: Scalars['String']['input'];
  resource: Scalars['String']['input'];
};

export type CreatePlasmidInput = {
  created_by: Scalars['String']['input'];
  dbxrefs?: InputMaybe<Array<Scalars['String']['input']>>;
  depositor?: InputMaybe<Scalars['String']['input']>;
  editable_summary?: InputMaybe<Scalars['String']['input']>;
  genbank_accession?: InputMaybe<Scalars['String']['input']>;
  genes?: InputMaybe<Array<Scalars['String']['input']>>;
  image_map?: InputMaybe<Scalars['String']['input']>;
  in_stock: Scalars['Boolean']['input'];
  keywords?: InputMaybe<Array<Scalars['String']['input']>>;
  name: Scalars['String']['input'];
  publications?: InputMaybe<Array<Scalars['String']['input']>>;
  sequence?: InputMaybe<Scalars['String']['input']>;
  summary?: InputMaybe<Scalars['String']['input']>;
  updated_by: Scalars['String']['input'];
};

export type CreateRoleInput = {
  description: Scalars['String']['input'];
  role: Scalars['String']['input'];
};

export type CreateStrainInput = {
  characteristics?: InputMaybe<Array<Scalars['String']['input']>>;
  created_by: Scalars['String']['input'];
  dbxrefs?: InputMaybe<Array<Scalars['String']['input']>>;
  depositor?: InputMaybe<Scalars['String']['input']>;
  editable_summary?: InputMaybe<Scalars['String']['input']>;
  genes?: InputMaybe<Array<Scalars['String']['input']>>;
  genetic_modification?: InputMaybe<Scalars['String']['input']>;
  genotypes?: InputMaybe<Array<Scalars['String']['input']>>;
  in_stock: Scalars['Boolean']['input'];
  label: Scalars['String']['input'];
  mutagenesis_method?: InputMaybe<Scalars['String']['input']>;
  names?: InputMaybe<Array<Scalars['String']['input']>>;
  parent?: InputMaybe<Scalars['String']['input']>;
  phenotypes?: InputMaybe<Array<Scalars['String']['input']>>;
  plasmid?: InputMaybe<Scalars['String']['input']>;
  publications?: InputMaybe<Array<Scalars['String']['input']>>;
  species: Scalars['String']['input'];
  summary?: InputMaybe<Scalars['String']['input']>;
  systematic_name: Scalars['String']['input'];
  updated_by: Scalars['String']['input'];
};

export type CreateUserInput = {
  city?: InputMaybe<Scalars['String']['input']>;
  country?: InputMaybe<Scalars['String']['input']>;
  email: Scalars['String']['input'];
  first_address?: InputMaybe<Scalars['String']['input']>;
  first_name: Scalars['String']['input'];
  group_name?: InputMaybe<Scalars['String']['input']>;
  is_active: Scalars['Boolean']['input'];
  last_name: Scalars['String']['input'];
  organization?: InputMaybe<Scalars['String']['input']>;
  phone?: InputMaybe<Scalars['String']['input']>;
  second_address?: InputMaybe<Scalars['String']['input']>;
  state?: InputMaybe<Scalars['String']['input']>;
  zipcode?: InputMaybe<Scalars['String']['input']>;
};

export type DeleteContent = {
  __typename?: 'DeleteContent';
  success: Scalars['Boolean']['output'];
};

export type DeletePermission = {
  __typename?: 'DeletePermission';
  success: Scalars['Boolean']['output'];
};

export type DeleteRole = {
  __typename?: 'DeleteRole';
  success: Scalars['Boolean']['output'];
};

export type DeleteStock = {
  __typename?: 'DeleteStock';
  success: Scalars['Boolean']['output'];
};

export type DeleteUser = {
  __typename?: 'DeleteUser';
  success: Scalars['Boolean']['output'];
};

export type Download = {
  __typename?: 'Download';
  items: Array<DownloadItem>;
  title: Scalars['String']['output'];
};

export type DownloadItem = {
  __typename?: 'DownloadItem';
  title: Scalars['String']['output'];
  url: Scalars['String']['output'];
};

export type Extension = {
  __typename?: 'Extension';
  db: Scalars['String']['output'];
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
  relation: Scalars['String']['output'];
};

/** The `UploadFile` type, represents the request for uploading a image file with a certain payload. */
export type FileToUpload = {
  file: Scalars['Upload']['input'];
  id: Scalars['Int']['input'];
};

export type GoAnnotation = {
  __typename?: 'GOAnnotation';
  assigned_by: Scalars['String']['output'];
  date: Scalars['String']['output'];
  evidence_code: Scalars['String']['output'];
  extensions?: Maybe<Array<Extension>>;
  go_term: Scalars['String']['output'];
  id: Scalars['String']['output'];
  publication: Scalars['String']['output'];
  qualifier: Scalars['String']['output'];
  type: Scalars['String']['output'];
  with?: Maybe<Array<With>>;
};

export type Gene = {
  __typename?: 'Gene';
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type GeneGeneralInfo = {
  __typename?: 'GeneGeneralInfo';
  description?: Maybe<Scalars['String']['output']>;
  gene_product?: Maybe<Scalars['String']['output']>;
  id: Scalars['String']['output'];
  name_description: Array<Maybe<Scalars['String']['output']>>;
  synonyms: Array<Maybe<Scalars['String']['output']>>;
};

export type Identity = {
  __typename?: 'Identity';
  created_at: Scalars['Timestamp']['output'];
  id: Scalars['ID']['output'];
  identifier: Scalars['String']['output'];
  provider: Scalars['String']['output'];
  updated_at: Scalars['Timestamp']['output'];
  user_id: Scalars['ID']['output'];
};

/** The `ImageFile` type, represents the response of uploading an image file. */
export type ImageFile = {
  __typename?: 'ImageFile';
  url: Scalars['String']['output'];
};

export type LoginInput = {
  client_id: Scalars['String']['input'];
  code: Scalars['String']['input'];
  provider: Scalars['String']['input'];
  redirect_url: Scalars['String']['input'];
  scopes: Scalars['String']['input'];
  state: Scalars['String']['input'];
};

export type Logout = {
  __typename?: 'Logout';
  success: Scalars['Boolean']['output'];
};

export type Mutation = {
  __typename?: 'Mutation';
  createContent?: Maybe<Content>;
  createOrder?: Maybe<Order>;
  createPermission?: Maybe<Permission>;
  createPlasmid?: Maybe<Plasmid>;
  createRole?: Maybe<Role>;
  createRolePermissionRelationship?: Maybe<Role>;
  createStrain?: Maybe<Strain>;
  createUser?: Maybe<User>;
  createUserRoleRelationship?: Maybe<User>;
  deleteContent?: Maybe<DeleteContent>;
  deletePermission?: Maybe<DeletePermission>;
  deleteRole?: Maybe<DeleteRole>;
  deleteStock?: Maybe<DeleteStock>;
  deleteUser?: Maybe<DeleteUser>;
  login?: Maybe<Auth>;
  logout?: Maybe<Logout>;
  updateContent?: Maybe<Content>;
  updateOrder?: Maybe<Order>;
  updatePermission?: Maybe<Permission>;
  updatePlasmid?: Maybe<Plasmid>;
  updateRole?: Maybe<Role>;
  updateStrain?: Maybe<Strain>;
  updateUser?: Maybe<User>;
  uploadFile: ImageFile;
};


export type MutationCreateContentArgs = {
  input?: InputMaybe<CreateContentInput>;
};


export type MutationCreateOrderArgs = {
  input?: InputMaybe<CreateOrderInput>;
};


export type MutationCreatePermissionArgs = {
  input?: InputMaybe<CreatePermissionInput>;
};


export type MutationCreatePlasmidArgs = {
  input?: InputMaybe<CreatePlasmidInput>;
};


export type MutationCreateRoleArgs = {
  input?: InputMaybe<CreateRoleInput>;
};


export type MutationCreateRolePermissionRelationshipArgs = {
  permissionId: Scalars['ID']['input'];
  roleId: Scalars['ID']['input'];
};


export type MutationCreateStrainArgs = {
  input?: InputMaybe<CreateStrainInput>;
};


export type MutationCreateUserArgs = {
  input?: InputMaybe<CreateUserInput>;
};


export type MutationCreateUserRoleRelationshipArgs = {
  roleId: Scalars['ID']['input'];
  userId: Scalars['ID']['input'];
};


export type MutationDeleteContentArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeletePermissionArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeleteRoleArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeleteStockArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeleteUserArgs = {
  id: Scalars['ID']['input'];
};


export type MutationLoginArgs = {
  input?: InputMaybe<LoginInput>;
};


export type MutationUpdateContentArgs = {
  input?: InputMaybe<UpdateContentInput>;
};


export type MutationUpdateOrderArgs = {
  id: Scalars['ID']['input'];
  input?: InputMaybe<UpdateOrderInput>;
};


export type MutationUpdatePermissionArgs = {
  id: Scalars['ID']['input'];
  input?: InputMaybe<UpdatePermissionInput>;
};


export type MutationUpdatePlasmidArgs = {
  id: Scalars['ID']['input'];
  input?: InputMaybe<UpdatePlasmidInput>;
};


export type MutationUpdateRoleArgs = {
  id: Scalars['ID']['input'];
  input?: InputMaybe<UpdateRoleInput>;
};


export type MutationUpdateStrainArgs = {
  id: Scalars['ID']['input'];
  input?: InputMaybe<UpdateStrainInput>;
};


export type MutationUpdateUserArgs = {
  id: Scalars['ID']['input'];
  input?: InputMaybe<UpdateUserInput>;
};


export type MutationUploadFileArgs = {
  file: Scalars['Upload']['input'];
};

export type NumberOfPublicationsWithGene = {
  __typename?: 'NumberOfPublicationsWithGene';
  num_pubs: Scalars['Int']['output'];
  publications: Array<PublicationWithGene>;
};

export type Order = {
  __typename?: 'Order';
  comments?: Maybe<Scalars['String']['output']>;
  consumer?: Maybe<User>;
  courier?: Maybe<Scalars['String']['output']>;
  courier_account?: Maybe<Scalars['String']['output']>;
  created_at: Scalars['Timestamp']['output'];
  id: Scalars['ID']['output'];
  items?: Maybe<Array<Stock>>;
  payer?: Maybe<User>;
  payment?: Maybe<Scalars['String']['output']>;
  purchase_order_num?: Maybe<Scalars['String']['output']>;
  purchaser?: Maybe<User>;
  status?: Maybe<StatusEnum>;
  updated_at: Scalars['Timestamp']['output'];
};

export type OrderListWithCursor = {
  __typename?: 'OrderListWithCursor';
  limit?: Maybe<Scalars['Int']['output']>;
  nextCursor: Scalars['Int']['output'];
  orders: Array<Order>;
  previousCursor: Scalars['Int']['output'];
  totalCount: Scalars['Int']['output'];
};

export type Organism = {
  __typename?: 'Organism';
  citations: Array<Citation>;
  downloads: Array<Download>;
  scientific_name: Scalars['String']['output'];
  taxon_id: Scalars['String']['output'];
};

export type Permission = {
  __typename?: 'Permission';
  created_at: Scalars['Timestamp']['output'];
  description: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  permission: Scalars['String']['output'];
  resource?: Maybe<Scalars['String']['output']>;
  updated_at: Scalars['Timestamp']['output'];
};

export type Phenotype = {
  __typename?: 'Phenotype';
  assay?: Maybe<Scalars['String']['output']>;
  environment?: Maybe<Scalars['String']['output']>;
  note?: Maybe<Scalars['String']['output']>;
  phenotype: Scalars['String']['output'];
  publication?: Maybe<Publication>;
};

export type Plasmid = Stock & {
  __typename?: 'Plasmid';
  created_at: Scalars['Timestamp']['output'];
  created_by: User;
  dbxrefs?: Maybe<Array<Scalars['String']['output']>>;
  depositor: User;
  editable_summary?: Maybe<Scalars['String']['output']>;
  genbank_accession?: Maybe<Scalars['String']['output']>;
  genes?: Maybe<Array<Gene>>;
  id: Scalars['ID']['output'];
  image_map?: Maybe<Scalars['String']['output']>;
  in_stock: Scalars['Boolean']['output'];
  keywords?: Maybe<Array<Scalars['String']['output']>>;
  name: Scalars['String']['output'];
  publications?: Maybe<Array<Publication>>;
  sequence?: Maybe<Scalars['String']['output']>;
  summary?: Maybe<Scalars['String']['output']>;
  updated_at: Scalars['Timestamp']['output'];
  updated_by: User;
};

export type PlasmidListWithCursor = {
  __typename?: 'PlasmidListWithCursor';
  limit?: Maybe<Scalars['Int']['output']>;
  nextCursor: Scalars['Int']['output'];
  plasmids: Array<Plasmid>;
  previousCursor: Scalars['Int']['output'];
  totalCount: Scalars['Int']['output'];
};

export type Publication = BasePublication & {
  __typename?: 'Publication';
  abstract: Scalars['String']['output'];
  authors: Array<Author>;
  doi?: Maybe<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  issn?: Maybe<Scalars['String']['output']>;
  issue?: Maybe<Scalars['String']['output']>;
  journal: Scalars['String']['output'];
  pages?: Maybe<Scalars['String']['output']>;
  pub_date?: Maybe<Scalars['Timestamp']['output']>;
  pub_type: Scalars['String']['output'];
  source: Scalars['String']['output'];
  status?: Maybe<Scalars['String']['output']>;
  title: Scalars['String']['output'];
  volume?: Maybe<Scalars['String']['output']>;
};

export type PublicationWithGene = BasePublication & {
  __typename?: 'PublicationWithGene';
  abstract: Scalars['String']['output'];
  authors: Array<Author>;
  doi?: Maybe<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  issn?: Maybe<Scalars['String']['output']>;
  issue?: Maybe<Scalars['String']['output']>;
  journal: Scalars['String']['output'];
  pages?: Maybe<Scalars['String']['output']>;
  pub_date?: Maybe<Scalars['Timestamp']['output']>;
  pub_type: Scalars['String']['output'];
  related_genes: Array<Gene>;
  source: Scalars['String']['output'];
  status?: Maybe<Scalars['String']['output']>;
  title: Scalars['String']['output'];
  volume?: Maybe<Scalars['String']['output']>;
};

export type Query = {
  __typename?: 'Query';
  content?: Maybe<Content>;
  contentBySlug?: Maybe<Content>;
  geneGeneralInformation?: Maybe<GeneGeneralInfo>;
  geneOntologyAnnotation?: Maybe<Array<GoAnnotation>>;
  listContentByNamespace: Array<Content>;
  listOrders?: Maybe<OrderListWithCursor>;
  listOrganisms?: Maybe<Array<Organism>>;
  listPermissions?: Maybe<Array<Permission>>;
  listPlasmids?: Maybe<PlasmidListWithCursor>;
  listPlasmidsWithAnnotation?: Maybe<PlasmidListWithCursor>;
  listPublicationsWithGene: Array<PublicationWithGene>;
  listRecentPlasmids?: Maybe<Array<Plasmid>>;
  listRecentPublications?: Maybe<Array<Publication>>;
  listRecentStrains?: Maybe<Array<Strain>>;
  listRoles?: Maybe<Array<Role>>;
  listStrains?: Maybe<StrainListWithCursor>;
  listStrainsWithAnnotation?: Maybe<StrainListWithCursor>;
  listStrainsWithGene?: Maybe<Array<Strain>>;
  listUsers?: Maybe<UserList>;
  order?: Maybe<Order>;
  organism?: Maybe<Organism>;
  permission?: Maybe<Permission>;
  plasmid?: Maybe<Plasmid>;
  publication?: Maybe<Publication>;
  role?: Maybe<Role>;
  strain?: Maybe<Strain>;
  user?: Maybe<User>;
  userByEmail?: Maybe<User>;
};


export type QueryContentArgs = {
  id: Scalars['ID']['input'];
};


export type QueryContentBySlugArgs = {
  slug: Scalars['String']['input'];
};


export type QueryGeneGeneralInformationArgs = {
  gene: Scalars['String']['input'];
};


export type QueryGeneOntologyAnnotationArgs = {
  gene: Scalars['String']['input'];
};


export type QueryListContentByNamespaceArgs = {
  namespace: Scalars['String']['input'];
};


export type QueryListOrdersArgs = {
  cursor?: InputMaybe<Scalars['Int']['input']>;
  filter?: InputMaybe<Scalars['String']['input']>;
  limit?: InputMaybe<Scalars['Int']['input']>;
};


export type QueryListPlasmidsArgs = {
  cursor?: InputMaybe<Scalars['Int']['input']>;
  filter?: InputMaybe<Scalars['String']['input']>;
  limit?: InputMaybe<Scalars['Int']['input']>;
};


export type QueryListPlasmidsWithAnnotationArgs = {
  annotation: Scalars['String']['input'];
  cursor?: InputMaybe<Scalars['Int']['input']>;
  limit?: InputMaybe<Scalars['Int']['input']>;
  type: Scalars['String']['input'];
};


export type QueryListPublicationsWithGeneArgs = {
  gene: Scalars['String']['input'];
};


export type QueryListRecentPlasmidsArgs = {
  limit: Scalars['Int']['input'];
};


export type QueryListRecentPublicationsArgs = {
  limit: Scalars['Int']['input'];
};


export type QueryListRecentStrainsArgs = {
  limit: Scalars['Int']['input'];
};


export type QueryListStrainsArgs = {
  cursor?: InputMaybe<Scalars['Int']['input']>;
  filter?: InputMaybe<StrainListFilter>;
  limit?: InputMaybe<Scalars['Int']['input']>;
};


export type QueryListStrainsWithAnnotationArgs = {
  annotation: Scalars['String']['input'];
  cursor?: InputMaybe<Scalars['Int']['input']>;
  limit?: InputMaybe<Scalars['Int']['input']>;
  type: Scalars['String']['input'];
};


export type QueryListStrainsWithGeneArgs = {
  gene: Scalars['String']['input'];
};


export type QueryListUsersArgs = {
  filter: Scalars['String']['input'];
  pagenum: Scalars['String']['input'];
  pagesize: Scalars['String']['input'];
};


export type QueryOrderArgs = {
  id: Scalars['ID']['input'];
};


export type QueryOrganismArgs = {
  taxon_id: Scalars['String']['input'];
};


export type QueryPermissionArgs = {
  id: Scalars['ID']['input'];
};


export type QueryPlasmidArgs = {
  id: Scalars['ID']['input'];
};


export type QueryPublicationArgs = {
  id: Scalars['ID']['input'];
};


export type QueryRoleArgs = {
  id: Scalars['ID']['input'];
};


export type QueryStrainArgs = {
  id: Scalars['ID']['input'];
};


export type QueryUserArgs = {
  id: Scalars['ID']['input'];
};


export type QueryUserByEmailArgs = {
  email: Scalars['String']['input'];
};

export type Role = {
  __typename?: 'Role';
  created_at: Scalars['Timestamp']['output'];
  description: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  permissions?: Maybe<Array<Permission>>;
  role: Scalars['String']['output'];
  updated_at: Scalars['Timestamp']['output'];
};

export enum StatusEnum {
  Cancelled = 'CANCELLED',
  Growing = 'GROWING',
  InPreparation = 'IN_PREPARATION',
  Shipped = 'SHIPPED'
}

export type Stock = {
  created_at: Scalars['Timestamp']['output'];
  created_by: User;
  dbxrefs?: Maybe<Array<Scalars['String']['output']>>;
  depositor: User;
  editable_summary?: Maybe<Scalars['String']['output']>;
  genes?: Maybe<Array<Gene>>;
  id: Scalars['ID']['output'];
  in_stock: Scalars['Boolean']['output'];
  publications?: Maybe<Array<Publication>>;
  summary?: Maybe<Scalars['String']['output']>;
  updated_at: Scalars['Timestamp']['output'];
  updated_by: User;
};

export type Strain = Stock & {
  __typename?: 'Strain';
  characteristics?: Maybe<Array<Scalars['String']['output']>>;
  created_at: Scalars['Timestamp']['output'];
  created_by: User;
  dbxrefs?: Maybe<Array<Scalars['String']['output']>>;
  depositor: User;
  editable_summary?: Maybe<Scalars['String']['output']>;
  genes?: Maybe<Array<Gene>>;
  genetic_modification?: Maybe<Scalars['String']['output']>;
  genotypes?: Maybe<Array<Scalars['String']['output']>>;
  id: Scalars['ID']['output'];
  in_stock: Scalars['Boolean']['output'];
  label: Scalars['String']['output'];
  mutagenesis_method?: Maybe<Scalars['String']['output']>;
  names?: Maybe<Array<Scalars['String']['output']>>;
  parent?: Maybe<Strain>;
  phenotypes?: Maybe<Array<Phenotype>>;
  plasmid?: Maybe<Scalars['String']['output']>;
  publications?: Maybe<Array<Publication>>;
  species: Scalars['String']['output'];
  summary?: Maybe<Scalars['String']['output']>;
  systematic_name: Scalars['String']['output'];
  updated_at: Scalars['Timestamp']['output'];
  updated_by: User;
};

export type StrainListFilter = {
  id?: InputMaybe<Scalars['ID']['input']>;
  in_stock?: InputMaybe<Scalars['Boolean']['input']>;
  label?: InputMaybe<Scalars['String']['input']>;
  strain_type: StrainType;
  summary?: InputMaybe<Scalars['String']['input']>;
};

export type StrainListWithCursor = {
  __typename?: 'StrainListWithCursor';
  limit?: Maybe<Scalars['Int']['output']>;
  nextCursor: Scalars['Int']['output'];
  previousCursor: Scalars['Int']['output'];
  strains: Array<Strain>;
  totalCount: Scalars['Int']['output'];
};

export enum StrainType {
  All = 'ALL',
  Bacterial = 'BACTERIAL',
  Gwdi = 'GWDI',
  Regular = 'REGULAR'
}

export type UpdateContentInput = {
  content: Scalars['String']['input'];
  id: Scalars['ID']['input'];
  updated_by: Scalars['String']['input'];
};

export type UpdateOrderInput = {
  comments?: InputMaybe<Scalars['String']['input']>;
  courier?: InputMaybe<Scalars['String']['input']>;
  courier_account?: InputMaybe<Scalars['String']['input']>;
  items?: InputMaybe<Array<Scalars['String']['input']>>;
  payment?: InputMaybe<Scalars['String']['input']>;
  purchase_order_num?: InputMaybe<Scalars['String']['input']>;
  status?: InputMaybe<StatusEnum>;
};

export type UpdatePermissionInput = {
  description: Scalars['String']['input'];
  permission: Scalars['String']['input'];
  resource: Scalars['String']['input'];
};

export type UpdatePlasmidInput = {
  dbxrefs?: InputMaybe<Array<Scalars['String']['input']>>;
  depositor?: InputMaybe<Scalars['String']['input']>;
  editable_summary?: InputMaybe<Scalars['String']['input']>;
  genbank_accession?: InputMaybe<Scalars['String']['input']>;
  genes?: InputMaybe<Array<Scalars['String']['input']>>;
  image_map?: InputMaybe<Scalars['String']['input']>;
  in_stock?: InputMaybe<Scalars['Boolean']['input']>;
  keywords?: InputMaybe<Array<Scalars['String']['input']>>;
  name?: InputMaybe<Scalars['String']['input']>;
  publications?: InputMaybe<Array<Scalars['String']['input']>>;
  sequence?: InputMaybe<Scalars['String']['input']>;
  summary?: InputMaybe<Scalars['String']['input']>;
  updated_by: Scalars['String']['input'];
};

export type UpdateRoleInput = {
  description: Scalars['String']['input'];
  role: Scalars['String']['input'];
};

export type UpdateStrainInput = {
  characteristics?: InputMaybe<Array<Scalars['String']['input']>>;
  dbxrefs?: InputMaybe<Array<Scalars['String']['input']>>;
  depositor?: InputMaybe<Scalars['String']['input']>;
  editable_summary?: InputMaybe<Scalars['String']['input']>;
  genes?: InputMaybe<Array<Scalars['String']['input']>>;
  genetic_modification?: InputMaybe<Scalars['String']['input']>;
  genotypes?: InputMaybe<Array<Scalars['String']['input']>>;
  in_stock?: InputMaybe<Scalars['Boolean']['input']>;
  label?: InputMaybe<Scalars['String']['input']>;
  mutagenesis_method?: InputMaybe<Scalars['String']['input']>;
  names?: InputMaybe<Array<Scalars['String']['input']>>;
  parent?: InputMaybe<Scalars['String']['input']>;
  phenotypes?: InputMaybe<Array<Scalars['String']['input']>>;
  plasmid?: InputMaybe<Scalars['String']['input']>;
  publications?: InputMaybe<Array<Scalars['String']['input']>>;
  species?: InputMaybe<Scalars['String']['input']>;
  summary?: InputMaybe<Scalars['String']['input']>;
  systematic_name?: InputMaybe<Scalars['String']['input']>;
  updated_by: Scalars['String']['input'];
};

export type UpdateUserInput = {
  city?: InputMaybe<Scalars['String']['input']>;
  country?: InputMaybe<Scalars['String']['input']>;
  first_address?: InputMaybe<Scalars['String']['input']>;
  first_name?: InputMaybe<Scalars['String']['input']>;
  group_name?: InputMaybe<Scalars['String']['input']>;
  is_active?: InputMaybe<Scalars['Boolean']['input']>;
  last_name?: InputMaybe<Scalars['String']['input']>;
  organization?: InputMaybe<Scalars['String']['input']>;
  phone?: InputMaybe<Scalars['String']['input']>;
  second_address?: InputMaybe<Scalars['String']['input']>;
  state?: InputMaybe<Scalars['String']['input']>;
  zipcode?: InputMaybe<Scalars['String']['input']>;
};

export type User = {
  __typename?: 'User';
  city?: Maybe<Scalars['String']['output']>;
  country?: Maybe<Scalars['String']['output']>;
  created_at: Scalars['Timestamp']['output'];
  email: Scalars['String']['output'];
  first_address?: Maybe<Scalars['String']['output']>;
  first_name: Scalars['String']['output'];
  group_name?: Maybe<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  is_active: Scalars['Boolean']['output'];
  last_name: Scalars['String']['output'];
  organization?: Maybe<Scalars['String']['output']>;
  phone?: Maybe<Scalars['String']['output']>;
  roles?: Maybe<Array<Role>>;
  second_address?: Maybe<Scalars['String']['output']>;
  state?: Maybe<Scalars['String']['output']>;
  updated_at: Scalars['Timestamp']['output'];
  zipcode?: Maybe<Scalars['String']['output']>;
};

export type UserList = {
  __typename?: 'UserList';
  pageNum?: Maybe<Scalars['String']['output']>;
  pageSize?: Maybe<Scalars['String']['output']>;
  totalCount: Scalars['Int']['output'];
  users: Array<User>;
};

export type With = {
  __typename?: 'With';
  db: Scalars['String']['output'];
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type LoginMutationVariables = Exact<{
  input: LoginInput;
}>;


export type LoginMutation = { __typename?: 'Mutation', login?: { __typename?: 'Auth', token: string, user: { __typename?: 'User', id: string, email: string, first_name: string, last_name: string, roles?: Array<{ __typename?: 'Role', role: string, permissions?: Array<{ __typename?: 'Permission', permission: string, resource?: string | null }> | null }> | null }, identity: { __typename?: 'Identity', provider: string } } | null };

export type LogoutMutationVariables = Exact<{ [key: string]: never; }>;


export type LogoutMutation = { __typename?: 'Mutation', logout?: { __typename?: 'Logout', success: boolean } | null };

export type CreateContentMutationVariables = Exact<{
  input: CreateContentInput;
}>;


export type CreateContentMutation = { __typename?: 'Mutation', createContent?: { __typename?: 'Content', name: string, content: string, namespace: string, created_by: { __typename?: 'User', id: string } } | null };

export type UpdateContentMutationVariables = Exact<{
  input: UpdateContentInput;
}>;


export type UpdateContentMutation = { __typename?: 'Mutation', updateContent?: { __typename?: 'Content', id: string, content: string, updated_by: { __typename?: 'User', id: string } } | null };

export type DeleteContentMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteContentMutation = { __typename?: 'Mutation', deleteContent?: { __typename?: 'DeleteContent', success: boolean } | null };

export type CreateOrderMutationVariables = Exact<{
  input: CreateOrderInput;
}>;


export type CreateOrderMutation = { __typename?: 'Mutation', createOrder?: { __typename?: 'Order', id: string } | null };

export type UploadFileMutationVariables = Exact<{
  file: Scalars['Upload']['input'];
}>;


export type UploadFileMutation = { __typename?: 'Mutation', uploadFile: { __typename?: 'ImageFile', url: string } };

export type CreateUserMutationVariables = Exact<{
  input: CreateUserInput;
}>;


export type CreateUserMutation = { __typename?: 'Mutation', createUser?: { __typename?: 'User', id: string } | null };

export type UpdateUserMutationVariables = Exact<{
  id: Scalars['ID']['input'];
  input: UpdateUserInput;
}>;


export type UpdateUserMutation = { __typename?: 'Mutation', updateUser?: { __typename?: 'User', id: string } | null };

export type ListContentByNamespaceQueryVariables = Exact<{
  namespace: Scalars['String']['input'];
}>;


export type ListContentByNamespaceQuery = { __typename?: 'Query', listContentByNamespace: Array<{ __typename?: 'Content', id: string, content: string, name: string, slug: string, created_at: any, updated_at: any, created_by: { __typename?: 'User', id: string, email: string, first_name: string, last_name: string }, updated_by: { __typename?: 'User', id: string, email: string, first_name: string, last_name: string } }> };

export type ContentBySlugQueryVariables = Exact<{
  slug: Scalars['String']['input'];
}>;


export type ContentBySlugQuery = { __typename?: 'Query', contentBySlug?: { __typename?: 'Content', id: string, content: string, name: string, slug: string, created_at: any, updated_at: any, created_by: { __typename?: 'User', id: string, email: string, first_name: string, last_name: string }, updated_by: { __typename?: 'User', id: string, email: string, first_name: string, last_name: string } } | null };

export type ContentQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type ContentQuery = { __typename?: 'Query', content?: { __typename?: 'Content', id: string, content: string, name: string, slug: string, namespace: string, created_at: any, updated_at: any, created_by: { __typename?: 'User', id: string, email: string, first_name: string, last_name: string }, updated_by: { __typename?: 'User', id: string, email: string, first_name: string, last_name: string } } | null };

export type ListOrganismsQueryVariables = Exact<{ [key: string]: never; }>;


export type ListOrganismsQuery = { __typename?: 'Query', listOrganisms?: Array<{ __typename?: 'Organism', taxon_id: string, scientific_name: string, citations: Array<{ __typename?: 'Citation', title: string, authors: string, pubmed_id: string, journal: string }>, downloads: Array<{ __typename?: 'Download', title: string, items: Array<{ __typename?: 'DownloadItem', title: string, url: string }> }> }> | null };

export type GeneSummaryQueryVariables = Exact<{
  gene: Scalars['String']['input'];
}>;


export type GeneSummaryQuery = { __typename?: 'Query', geneGeneralInformation?: { __typename?: 'GeneGeneralInfo', id: string, name_description: Array<string | null>, gene_product?: string | null, synonyms: Array<string | null>, description?: string | null } | null, geneOntologyAnnotation?: Array<{ __typename?: 'GOAnnotation', id: string, type: string, date: string, go_term: string, evidence_code: string, with?: Array<{ __typename?: 'With', id: string, db: string, name: string }> | null, extensions?: Array<{ __typename?: 'Extension', id: string, db: string, relation: string, name: string }> | null }> | null, listPublicationsWithGene: Array<{ __typename?: 'PublicationWithGene', id: string, title: string, journal: string, pages?: string | null, issue?: string | null, authors: Array<{ __typename?: 'Author', last_name: string }> }> };

export type GeneOntologyAnnotationQueryVariables = Exact<{
  gene: Scalars['String']['input'];
}>;


export type GeneOntologyAnnotationQuery = { __typename?: 'Query', geneOntologyAnnotation?: Array<{ __typename?: 'GOAnnotation', id: string, type: string, date: string, go_term: string, evidence_code: string, qualifier: string, publication: string, assigned_by: string, with?: Array<{ __typename?: 'With', id: string, db: string, name: string }> | null, extensions?: Array<{ __typename?: 'Extension', id: string, db: string, relation: string, name: string }> | null }> | null };

export type ListStrainsWithGeneQueryVariables = Exact<{
  gene: Scalars['String']['input'];
}>;


export type ListStrainsWithGeneQuery = { __typename?: 'Query', listStrainsWithGene?: Array<{ __typename?: 'Strain', id: string, label: string, characteristics?: Array<string> | null, in_stock: boolean, phenotypes?: Array<{ __typename?: 'Phenotype', phenotype: string, publication?: { __typename?: 'Publication', id: string, title: string, journal: string, pages?: string | null, volume?: string | null, pub_date?: any | null, authors: Array<{ __typename?: 'Author', last_name: string, rank?: string | null }> } | null }> | null }> | null };

export type PublicationQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type PublicationQuery = { __typename?: 'Query', publication?: { __typename?: 'Publication', id: string, doi?: string | null, title: string, abstract: string, journal: string, pub_date?: any | null, pages?: string | null, issue?: string | null, volume?: string | null, authors: Array<{ __typename?: 'Author', initials?: string | null, last_name: string }> } | null };

export type ListRecentPublicationsQueryVariables = Exact<{
  limit?: Scalars['Int']['input'];
}>;


export type ListRecentPublicationsQuery = { __typename?: 'Query', listRecentPublications?: Array<{ __typename?: 'Publication', id: string, doi?: string | null, title: string, abstract: string, journal: string, pub_date?: any | null, pages?: string | null, issue?: string | null, volume?: string | null, authors: Array<{ __typename?: 'Author', initials?: string | null, last_name: string }> }> | null };

export type ListPublicationsWithGeneQueryVariables = Exact<{
  gene: Scalars['String']['input'];
}>;


export type ListPublicationsWithGeneQuery = { __typename?: 'Query', listPublicationsWithGene: Array<{ __typename?: 'PublicationWithGene', id: string, doi?: string | null, title: string, journal: string, pub_date?: any | null, volume?: string | null, pages?: string | null, pub_type: string, source: string, issue?: string | null, related_genes: Array<{ __typename?: 'Gene', id: string, name: string }>, authors: Array<{ __typename?: 'Author', last_name: string, rank?: string | null }> }> };

export type StrainListQueryVariables = Exact<{
  cursor: Scalars['Int']['input'];
  limit: Scalars['Int']['input'];
  filter?: InputMaybe<StrainListFilter>;
}>;


export type StrainListQuery = { __typename?: 'Query', listStrains?: { __typename?: 'StrainListWithCursor', nextCursor: number, totalCount: number, strains: Array<{ __typename?: 'Strain', id: string, label: string, summary?: string | null, in_stock: boolean }> } | null };

export type ListStrainsWithPhenotypeQueryVariables = Exact<{
  cursor: Scalars['Int']['input'];
  limit: Scalars['Int']['input'];
  type: Scalars['String']['input'];
  annotation: Scalars['String']['input'];
}>;


export type ListStrainsWithPhenotypeQuery = { __typename?: 'Query', listStrainsWithAnnotation?: { __typename?: 'StrainListWithCursor', totalCount: number, nextCursor: number, strains: Array<{ __typename?: 'Strain', id: string, label: string, genes?: Array<{ __typename?: 'Gene', name: string }> | null, publications?: Array<{ __typename?: 'Publication', id: string, pub_date?: any | null, title: string, journal: string, volume?: string | null, pages?: string | null, authors: Array<{ __typename?: 'Author', last_name: string }> }> | null }> } | null };

export type ListBacterialStrainsQueryVariables = Exact<{ [key: string]: never; }>;


export type ListBacterialStrainsQuery = { __typename?: 'Query', bacterialFoodSource?: { __typename?: 'StrainListWithCursor', totalCount: number, nextCursor: number, strains: Array<{ __typename?: 'Strain', id: string, label: string, summary?: string | null, in_stock: boolean }> } | null, symbioticFarmerBacterium?: { __typename?: 'StrainListWithCursor', totalCount: number, nextCursor: number, strains: Array<{ __typename?: 'Strain', id: string, label: string, summary?: string | null, in_stock: boolean }> } | null };

export type ListStrainsInventoryQueryVariables = Exact<{
  cursor: Scalars['Int']['input'];
  limit: Scalars['Int']['input'];
}>;


export type ListStrainsInventoryQuery = { __typename?: 'Query', listStrainsWithAnnotation?: { __typename?: 'StrainListWithCursor', totalCount: number, nextCursor: number, strains: Array<{ __typename?: 'Strain', id: string, label: string, summary?: string | null, in_stock: boolean }> } | null };

export type ListPlasmidsInventoryQueryVariables = Exact<{
  cursor: Scalars['Int']['input'];
  limit: Scalars['Int']['input'];
}>;


export type ListPlasmidsInventoryQuery = { __typename?: 'Query', listPlasmidsWithAnnotation?: { __typename?: 'PlasmidListWithCursor', totalCount: number, nextCursor: number, plasmids: Array<{ __typename?: 'Plasmid', id: string, name: string, summary?: string | null, in_stock: boolean }> } | null };

export type PlasmidListFilterQueryVariables = Exact<{
  cursor: Scalars['Int']['input'];
  limit: Scalars['Int']['input'];
  filter: Scalars['String']['input'];
}>;


export type PlasmidListFilterQuery = { __typename?: 'Query', listPlasmids?: { __typename?: 'PlasmidListWithCursor', nextCursor: number, totalCount: number, plasmids: Array<{ __typename?: 'Plasmid', id: string, name: string, summary?: string | null, in_stock: boolean }> } | null };

export type PlasmidQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type PlasmidQuery = { __typename?: 'Query', plasmid?: { __typename?: 'Plasmid', id: string, name: string, summary?: string | null, dbxrefs?: Array<string> | null, image_map?: string | null, sequence?: string | null, keywords?: Array<string> | null, genbank_accession?: string | null, in_stock: boolean, depositor: { __typename?: 'User', first_name: string, last_name: string }, publications?: Array<{ __typename?: 'Publication', id: string, pub_date?: any | null, title: string, journal: string, volume?: string | null, pages?: string | null, doi?: string | null, authors: Array<{ __typename?: 'Author', last_name: string }> }> | null, genes?: Array<{ __typename?: 'Gene', name: string }> | null } | null };

export type StrainQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type StrainQuery = { __typename?: 'Query', strain?: { __typename?: 'Strain', id: string, label: string, summary?: string | null, species: string, plasmid?: string | null, dbxrefs?: Array<string> | null, in_stock: boolean, systematic_name: string, genotypes?: Array<string> | null, mutagenesis_method?: string | null, genetic_modification?: string | null, names?: Array<string> | null, characteristics?: Array<string> | null, parent?: { __typename?: 'Strain', id: string, label: string } | null, depositor: { __typename?: 'User', first_name: string, last_name: string }, publications?: Array<{ __typename?: 'Publication', id: string, pub_date?: any | null, title: string, journal: string, volume?: string | null, pages?: string | null, doi?: string | null, authors: Array<{ __typename?: 'Author', last_name: string }> }> | null, genes?: Array<{ __typename?: 'Gene', name: string }> | null, phenotypes?: Array<{ __typename?: 'Phenotype', phenotype: string, note?: string | null, assay?: string | null, environment?: string | null, publication?: { __typename?: 'Publication', id: string, pub_date?: any | null, title: string, journal: string, volume?: string | null, pages?: string | null, authors: Array<{ __typename?: 'Author', last_name: string }> } | null }> | null } | null };

export type ListRecentPlasmidsQueryVariables = Exact<{
  limit?: Scalars['Int']['input'];
}>;


export type ListRecentPlasmidsQuery = { __typename?: 'Query', listRecentPlasmids?: Array<{ __typename?: 'Plasmid', id: string, created_at: any, name: string }> | null };

export type ListRecentStrainsQueryVariables = Exact<{
  limit?: Scalars['Int']['input'];
}>;


export type ListRecentStrainsQuery = { __typename?: 'Query', listRecentStrains?: Array<{ __typename?: 'Strain', id: string, created_at: any, systematic_name: string }> | null };

export type UserByEmailQueryVariables = Exact<{
  email: Scalars['String']['input'];
}>;


export type UserByEmailQuery = { __typename?: 'Query', userByEmail?: { __typename?: 'User', id: string } | null };


/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockLoginMutation((req, res, ctx) => {
 *   const { input } = req.variables;
 *   return res(
 *     ctx.data({ login })
 *   )
 * })
 */
export const mockLoginMutation = (resolver: Parameters<typeof graphql.mutation<LoginMutation, LoginMutationVariables>>[1]) =>
  graphql.mutation<LoginMutation, LoginMutationVariables>(
    'Login',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockLogoutMutation((req, res, ctx) => {
 *   return res(
 *     ctx.data({ logout })
 *   )
 * })
 */
export const mockLogoutMutation = (resolver: Parameters<typeof graphql.mutation<LogoutMutation, LogoutMutationVariables>>[1]) =>
  graphql.mutation<LogoutMutation, LogoutMutationVariables>(
    'Logout',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockCreateContentMutation((req, res, ctx) => {
 *   const { input } = req.variables;
 *   return res(
 *     ctx.data({ createContent })
 *   )
 * })
 */
export const mockCreateContentMutation = (resolver: Parameters<typeof graphql.mutation<CreateContentMutation, CreateContentMutationVariables>>[1]) =>
  graphql.mutation<CreateContentMutation, CreateContentMutationVariables>(
    'CreateContent',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockUpdateContentMutation((req, res, ctx) => {
 *   const { input } = req.variables;
 *   return res(
 *     ctx.data({ updateContent })
 *   )
 * })
 */
export const mockUpdateContentMutation = (resolver: Parameters<typeof graphql.mutation<UpdateContentMutation, UpdateContentMutationVariables>>[1]) =>
  graphql.mutation<UpdateContentMutation, UpdateContentMutationVariables>(
    'UpdateContent',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockDeleteContentMutation((req, res, ctx) => {
 *   const { id } = req.variables;
 *   return res(
 *     ctx.data({ deleteContent })
 *   )
 * })
 */
export const mockDeleteContentMutation = (resolver: Parameters<typeof graphql.mutation<DeleteContentMutation, DeleteContentMutationVariables>>[1]) =>
  graphql.mutation<DeleteContentMutation, DeleteContentMutationVariables>(
    'DeleteContent',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockCreateOrderMutation((req, res, ctx) => {
 *   const { input } = req.variables;
 *   return res(
 *     ctx.data({ createOrder })
 *   )
 * })
 */
export const mockCreateOrderMutation = (resolver: Parameters<typeof graphql.mutation<CreateOrderMutation, CreateOrderMutationVariables>>[1]) =>
  graphql.mutation<CreateOrderMutation, CreateOrderMutationVariables>(
    'CreateOrder',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockUploadFileMutation((req, res, ctx) => {
 *   const { file } = req.variables;
 *   return res(
 *     ctx.data({ uploadFile })
 *   )
 * })
 */
export const mockUploadFileMutation = (resolver: Parameters<typeof graphql.mutation<UploadFileMutation, UploadFileMutationVariables>>[1]) =>
  graphql.mutation<UploadFileMutation, UploadFileMutationVariables>(
    'UploadFile',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockCreateUserMutation((req, res, ctx) => {
 *   const { input } = req.variables;
 *   return res(
 *     ctx.data({ createUser })
 *   )
 * })
 */
export const mockCreateUserMutation = (resolver: Parameters<typeof graphql.mutation<CreateUserMutation, CreateUserMutationVariables>>[1]) =>
  graphql.mutation<CreateUserMutation, CreateUserMutationVariables>(
    'CreateUser',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockUpdateUserMutation((req, res, ctx) => {
 *   const { id, input } = req.variables;
 *   return res(
 *     ctx.data({ updateUser })
 *   )
 * })
 */
export const mockUpdateUserMutation = (resolver: Parameters<typeof graphql.mutation<UpdateUserMutation, UpdateUserMutationVariables>>[1]) =>
  graphql.mutation<UpdateUserMutation, UpdateUserMutationVariables>(
    'UpdateUser',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListContentByNamespaceQuery((req, res, ctx) => {
 *   const { namespace } = req.variables;
 *   return res(
 *     ctx.data({ listContentByNamespace })
 *   )
 * })
 */
export const mockListContentByNamespaceQuery = (resolver: Parameters<typeof graphql.query<ListContentByNamespaceQuery, ListContentByNamespaceQueryVariables>>[1]) =>
  graphql.query<ListContentByNamespaceQuery, ListContentByNamespaceQueryVariables>(
    'ListContentByNamespace',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockContentBySlugQuery((req, res, ctx) => {
 *   const { slug } = req.variables;
 *   return res(
 *     ctx.data({ contentBySlug })
 *   )
 * })
 */
export const mockContentBySlugQuery = (resolver: Parameters<typeof graphql.query<ContentBySlugQuery, ContentBySlugQueryVariables>>[1]) =>
  graphql.query<ContentBySlugQuery, ContentBySlugQueryVariables>(
    'ContentBySlug',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockContentQuery((req, res, ctx) => {
 *   const { id } = req.variables;
 *   return res(
 *     ctx.data({ content })
 *   )
 * })
 */
export const mockContentQuery = (resolver: Parameters<typeof graphql.query<ContentQuery, ContentQueryVariables>>[1]) =>
  graphql.query<ContentQuery, ContentQueryVariables>(
    'Content',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListOrganismsQuery((req, res, ctx) => {
 *   return res(
 *     ctx.data({ listOrganisms })
 *   )
 * })
 */
export const mockListOrganismsQuery = (resolver: Parameters<typeof graphql.query<ListOrganismsQuery, ListOrganismsQueryVariables>>[1]) =>
  graphql.query<ListOrganismsQuery, ListOrganismsQueryVariables>(
    'ListOrganisms',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockGeneSummaryQuery((req, res, ctx) => {
 *   const { gene } = req.variables;
 *   return res(
 *     ctx.data({ geneGeneralInformation, geneOntologyAnnotation, listPublicationsWithGene })
 *   )
 * })
 */
export const mockGeneSummaryQuery = (resolver: Parameters<typeof graphql.query<GeneSummaryQuery, GeneSummaryQueryVariables>>[1]) =>
  graphql.query<GeneSummaryQuery, GeneSummaryQueryVariables>(
    'GeneSummary',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockGeneOntologyAnnotationQuery((req, res, ctx) => {
 *   const { gene } = req.variables;
 *   return res(
 *     ctx.data({ geneOntologyAnnotation })
 *   )
 * })
 */
export const mockGeneOntologyAnnotationQuery = (resolver: Parameters<typeof graphql.query<GeneOntologyAnnotationQuery, GeneOntologyAnnotationQueryVariables>>[1]) =>
  graphql.query<GeneOntologyAnnotationQuery, GeneOntologyAnnotationQueryVariables>(
    'GeneOntologyAnnotation',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListStrainsWithGeneQuery((req, res, ctx) => {
 *   const { gene } = req.variables;
 *   return res(
 *     ctx.data({ listStrainsWithGene })
 *   )
 * })
 */
export const mockListStrainsWithGeneQuery = (resolver: Parameters<typeof graphql.query<ListStrainsWithGeneQuery, ListStrainsWithGeneQueryVariables>>[1]) =>
  graphql.query<ListStrainsWithGeneQuery, ListStrainsWithGeneQueryVariables>(
    'ListStrainsWithGene',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockPublicationQuery((req, res, ctx) => {
 *   const { id } = req.variables;
 *   return res(
 *     ctx.data({ publication })
 *   )
 * })
 */
export const mockPublicationQuery = (resolver: Parameters<typeof graphql.query<PublicationQuery, PublicationQueryVariables>>[1]) =>
  graphql.query<PublicationQuery, PublicationQueryVariables>(
    'Publication',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListRecentPublicationsQuery((req, res, ctx) => {
 *   const { limit } = req.variables;
 *   return res(
 *     ctx.data({ listRecentPublications })
 *   )
 * })
 */
export const mockListRecentPublicationsQuery = (resolver: Parameters<typeof graphql.query<ListRecentPublicationsQuery, ListRecentPublicationsQueryVariables>>[1]) =>
  graphql.query<ListRecentPublicationsQuery, ListRecentPublicationsQueryVariables>(
    'ListRecentPublications',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListPublicationsWithGeneQuery((req, res, ctx) => {
 *   const { gene } = req.variables;
 *   return res(
 *     ctx.data({ listPublicationsWithGene })
 *   )
 * })
 */
export const mockListPublicationsWithGeneQuery = (resolver: Parameters<typeof graphql.query<ListPublicationsWithGeneQuery, ListPublicationsWithGeneQueryVariables>>[1]) =>
  graphql.query<ListPublicationsWithGeneQuery, ListPublicationsWithGeneQueryVariables>(
    'ListPublicationsWithGene',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockStrainListQuery((req, res, ctx) => {
 *   const { cursor, limit, filter } = req.variables;
 *   return res(
 *     ctx.data({ listStrains })
 *   )
 * })
 */
export const mockStrainListQuery = (resolver: Parameters<typeof graphql.query<StrainListQuery, StrainListQueryVariables>>[1]) =>
  graphql.query<StrainListQuery, StrainListQueryVariables>(
    'StrainList',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListStrainsWithPhenotypeQuery((req, res, ctx) => {
 *   const { cursor, limit, type, annotation } = req.variables;
 *   return res(
 *     ctx.data({ listStrainsWithAnnotation })
 *   )
 * })
 */
export const mockListStrainsWithPhenotypeQuery = (resolver: Parameters<typeof graphql.query<ListStrainsWithPhenotypeQuery, ListStrainsWithPhenotypeQueryVariables>>[1]) =>
  graphql.query<ListStrainsWithPhenotypeQuery, ListStrainsWithPhenotypeQueryVariables>(
    'ListStrainsWithPhenotype',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListBacterialStrainsQuery((req, res, ctx) => {
 *   return res(
 *     ctx.data({ listStrainsWithAnnotation, listStrainsWithAnnotation })
 *   )
 * })
 */
export const mockListBacterialStrainsQuery = (resolver: Parameters<typeof graphql.query<ListBacterialStrainsQuery, ListBacterialStrainsQueryVariables>>[1]) =>
  graphql.query<ListBacterialStrainsQuery, ListBacterialStrainsQueryVariables>(
    'ListBacterialStrains',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListStrainsInventoryQuery((req, res, ctx) => {
 *   const { cursor, limit } = req.variables;
 *   return res(
 *     ctx.data({ listStrainsWithAnnotation })
 *   )
 * })
 */
export const mockListStrainsInventoryQuery = (resolver: Parameters<typeof graphql.query<ListStrainsInventoryQuery, ListStrainsInventoryQueryVariables>>[1]) =>
  graphql.query<ListStrainsInventoryQuery, ListStrainsInventoryQueryVariables>(
    'ListStrainsInventory',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListPlasmidsInventoryQuery((req, res, ctx) => {
 *   const { cursor, limit } = req.variables;
 *   return res(
 *     ctx.data({ listPlasmidsWithAnnotation })
 *   )
 * })
 */
export const mockListPlasmidsInventoryQuery = (resolver: Parameters<typeof graphql.query<ListPlasmidsInventoryQuery, ListPlasmidsInventoryQueryVariables>>[1]) =>
  graphql.query<ListPlasmidsInventoryQuery, ListPlasmidsInventoryQueryVariables>(
    'ListPlasmidsInventory',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockPlasmidListFilterQuery((req, res, ctx) => {
 *   const { cursor, limit, filter } = req.variables;
 *   return res(
 *     ctx.data({ listPlasmids })
 *   )
 * })
 */
export const mockPlasmidListFilterQuery = (resolver: Parameters<typeof graphql.query<PlasmidListFilterQuery, PlasmidListFilterQueryVariables>>[1]) =>
  graphql.query<PlasmidListFilterQuery, PlasmidListFilterQueryVariables>(
    'PlasmidListFilter',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockPlasmidQuery((req, res, ctx) => {
 *   const { id } = req.variables;
 *   return res(
 *     ctx.data({ plasmid })
 *   )
 * })
 */
export const mockPlasmidQuery = (resolver: Parameters<typeof graphql.query<PlasmidQuery, PlasmidQueryVariables>>[1]) =>
  graphql.query<PlasmidQuery, PlasmidQueryVariables>(
    'Plasmid',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockStrainQuery((req, res, ctx) => {
 *   const { id } = req.variables;
 *   return res(
 *     ctx.data({ strain })
 *   )
 * })
 */
export const mockStrainQuery = (resolver: Parameters<typeof graphql.query<StrainQuery, StrainQueryVariables>>[1]) =>
  graphql.query<StrainQuery, StrainQueryVariables>(
    'Strain',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListRecentPlasmidsQuery((req, res, ctx) => {
 *   const { limit } = req.variables;
 *   return res(
 *     ctx.data({ listRecentPlasmids })
 *   )
 * })
 */
export const mockListRecentPlasmidsQuery = (resolver: Parameters<typeof graphql.query<ListRecentPlasmidsQuery, ListRecentPlasmidsQueryVariables>>[1]) =>
  graphql.query<ListRecentPlasmidsQuery, ListRecentPlasmidsQueryVariables>(
    'ListRecentPlasmids',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListRecentStrainsQuery((req, res, ctx) => {
 *   const { limit } = req.variables;
 *   return res(
 *     ctx.data({ listRecentStrains })
 *   )
 * })
 */
export const mockListRecentStrainsQuery = (resolver: Parameters<typeof graphql.query<ListRecentStrainsQuery, ListRecentStrainsQueryVariables>>[1]) =>
  graphql.query<ListRecentStrainsQuery, ListRecentStrainsQueryVariables>(
    'ListRecentStrains',
    resolver
  )

/**
 * @param resolver a function that accepts a captured request and may return a mocked response.
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockUserByEmailQuery((req, res, ctx) => {
 *   const { email } = req.variables;
 *   return res(
 *     ctx.data({ userByEmail })
 *   )
 * })
 */
export const mockUserByEmailQuery = (resolver: Parameters<typeof graphql.query<UserByEmailQuery, UserByEmailQueryVariables>>[1]) =>
  graphql.query<UserByEmailQuery, UserByEmailQueryVariables>(
    'UserByEmail',
    resolver
  )
