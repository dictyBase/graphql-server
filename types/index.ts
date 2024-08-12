import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
const defaultOptions = {} as const;
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


export const LoginDocument = gql`
    mutation Login($input: LoginInput!) {
  login(input: $input) {
    token
    user {
      id
      email
      first_name
      last_name
      roles {
        role
        permissions {
          permission
          resource
        }
      }
    }
    identity {
      provider
    }
  }
}
    `;
export type LoginMutationFn = Apollo.MutationFunction<LoginMutation, LoginMutationVariables>;

/**
 * __useLoginMutation__
 *
 * To run a mutation, you first call `useLoginMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useLoginMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [loginMutation, { data, loading, error }] = useLoginMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useLoginMutation(baseOptions?: Apollo.MutationHookOptions<LoginMutation, LoginMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<LoginMutation, LoginMutationVariables>(LoginDocument, options);
      }
export type LoginMutationHookResult = ReturnType<typeof useLoginMutation>;
export type LoginMutationResult = Apollo.MutationResult<LoginMutation>;
export type LoginMutationOptions = Apollo.BaseMutationOptions<LoginMutation, LoginMutationVariables>;
export const LogoutDocument = gql`
    mutation Logout {
  logout {
    success
  }
}
    `;
export type LogoutMutationFn = Apollo.MutationFunction<LogoutMutation, LogoutMutationVariables>;

/**
 * __useLogoutMutation__
 *
 * To run a mutation, you first call `useLogoutMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useLogoutMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [logoutMutation, { data, loading, error }] = useLogoutMutation({
 *   variables: {
 *   },
 * });
 */
export function useLogoutMutation(baseOptions?: Apollo.MutationHookOptions<LogoutMutation, LogoutMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<LogoutMutation, LogoutMutationVariables>(LogoutDocument, options);
      }
export type LogoutMutationHookResult = ReturnType<typeof useLogoutMutation>;
export type LogoutMutationResult = Apollo.MutationResult<LogoutMutation>;
export type LogoutMutationOptions = Apollo.BaseMutationOptions<LogoutMutation, LogoutMutationVariables>;
export const CreateContentDocument = gql`
    mutation CreateContent($input: CreateContentInput!) {
  createContent(input: $input) {
    name
    created_by {
      id
    }
    content
    namespace
  }
}
    `;
export type CreateContentMutationFn = Apollo.MutationFunction<CreateContentMutation, CreateContentMutationVariables>;

/**
 * __useCreateContentMutation__
 *
 * To run a mutation, you first call `useCreateContentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateContentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createContentMutation, { data, loading, error }] = useCreateContentMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useCreateContentMutation(baseOptions?: Apollo.MutationHookOptions<CreateContentMutation, CreateContentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateContentMutation, CreateContentMutationVariables>(CreateContentDocument, options);
      }
export type CreateContentMutationHookResult = ReturnType<typeof useCreateContentMutation>;
export type CreateContentMutationResult = Apollo.MutationResult<CreateContentMutation>;
export type CreateContentMutationOptions = Apollo.BaseMutationOptions<CreateContentMutation, CreateContentMutationVariables>;
export const UpdateContentDocument = gql`
    mutation UpdateContent($input: UpdateContentInput!) {
  updateContent(input: $input) {
    id
    updated_by {
      id
    }
    content
  }
}
    `;
export type UpdateContentMutationFn = Apollo.MutationFunction<UpdateContentMutation, UpdateContentMutationVariables>;

/**
 * __useUpdateContentMutation__
 *
 * To run a mutation, you first call `useUpdateContentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateContentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateContentMutation, { data, loading, error }] = useUpdateContentMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateContentMutation(baseOptions?: Apollo.MutationHookOptions<UpdateContentMutation, UpdateContentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateContentMutation, UpdateContentMutationVariables>(UpdateContentDocument, options);
      }
export type UpdateContentMutationHookResult = ReturnType<typeof useUpdateContentMutation>;
export type UpdateContentMutationResult = Apollo.MutationResult<UpdateContentMutation>;
export type UpdateContentMutationOptions = Apollo.BaseMutationOptions<UpdateContentMutation, UpdateContentMutationVariables>;
export const DeleteContentDocument = gql`
    mutation DeleteContent($id: ID!) {
  deleteContent(id: $id) {
    success
  }
}
    `;
export type DeleteContentMutationFn = Apollo.MutationFunction<DeleteContentMutation, DeleteContentMutationVariables>;

/**
 * __useDeleteContentMutation__
 *
 * To run a mutation, you first call `useDeleteContentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteContentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteContentMutation, { data, loading, error }] = useDeleteContentMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useDeleteContentMutation(baseOptions?: Apollo.MutationHookOptions<DeleteContentMutation, DeleteContentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteContentMutation, DeleteContentMutationVariables>(DeleteContentDocument, options);
      }
export type DeleteContentMutationHookResult = ReturnType<typeof useDeleteContentMutation>;
export type DeleteContentMutationResult = Apollo.MutationResult<DeleteContentMutation>;
export type DeleteContentMutationOptions = Apollo.BaseMutationOptions<DeleteContentMutation, DeleteContentMutationVariables>;
export const CreateOrderDocument = gql`
    mutation CreateOrder($input: CreateOrderInput!) {
  createOrder(input: $input) {
    id
  }
}
    `;
export type CreateOrderMutationFn = Apollo.MutationFunction<CreateOrderMutation, CreateOrderMutationVariables>;

/**
 * __useCreateOrderMutation__
 *
 * To run a mutation, you first call `useCreateOrderMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateOrderMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createOrderMutation, { data, loading, error }] = useCreateOrderMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useCreateOrderMutation(baseOptions?: Apollo.MutationHookOptions<CreateOrderMutation, CreateOrderMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateOrderMutation, CreateOrderMutationVariables>(CreateOrderDocument, options);
      }
export type CreateOrderMutationHookResult = ReturnType<typeof useCreateOrderMutation>;
export type CreateOrderMutationResult = Apollo.MutationResult<CreateOrderMutation>;
export type CreateOrderMutationOptions = Apollo.BaseMutationOptions<CreateOrderMutation, CreateOrderMutationVariables>;
export const UploadFileDocument = gql`
    mutation UploadFile($file: Upload!) {
  uploadFile(file: $file) {
    url
  }
}
    `;
export type UploadFileMutationFn = Apollo.MutationFunction<UploadFileMutation, UploadFileMutationVariables>;

/**
 * __useUploadFileMutation__
 *
 * To run a mutation, you first call `useUploadFileMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUploadFileMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [uploadFileMutation, { data, loading, error }] = useUploadFileMutation({
 *   variables: {
 *      file: // value for 'file'
 *   },
 * });
 */
export function useUploadFileMutation(baseOptions?: Apollo.MutationHookOptions<UploadFileMutation, UploadFileMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UploadFileMutation, UploadFileMutationVariables>(UploadFileDocument, options);
      }
export type UploadFileMutationHookResult = ReturnType<typeof useUploadFileMutation>;
export type UploadFileMutationResult = Apollo.MutationResult<UploadFileMutation>;
export type UploadFileMutationOptions = Apollo.BaseMutationOptions<UploadFileMutation, UploadFileMutationVariables>;
export const CreateUserDocument = gql`
    mutation CreateUser($input: CreateUserInput!) {
  createUser(input: $input) {
    id
  }
}
    `;
export type CreateUserMutationFn = Apollo.MutationFunction<CreateUserMutation, CreateUserMutationVariables>;

/**
 * __useCreateUserMutation__
 *
 * To run a mutation, you first call `useCreateUserMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateUserMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createUserMutation, { data, loading, error }] = useCreateUserMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useCreateUserMutation(baseOptions?: Apollo.MutationHookOptions<CreateUserMutation, CreateUserMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateUserMutation, CreateUserMutationVariables>(CreateUserDocument, options);
      }
export type CreateUserMutationHookResult = ReturnType<typeof useCreateUserMutation>;
export type CreateUserMutationResult = Apollo.MutationResult<CreateUserMutation>;
export type CreateUserMutationOptions = Apollo.BaseMutationOptions<CreateUserMutation, CreateUserMutationVariables>;
export const UpdateUserDocument = gql`
    mutation UpdateUser($id: ID!, $input: UpdateUserInput!) {
  updateUser(id: $id, input: $input) {
    id
  }
}
    `;
export type UpdateUserMutationFn = Apollo.MutationFunction<UpdateUserMutation, UpdateUserMutationVariables>;

/**
 * __useUpdateUserMutation__
 *
 * To run a mutation, you first call `useUpdateUserMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateUserMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateUserMutation, { data, loading, error }] = useUpdateUserMutation({
 *   variables: {
 *      id: // value for 'id'
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateUserMutation(baseOptions?: Apollo.MutationHookOptions<UpdateUserMutation, UpdateUserMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateUserMutation, UpdateUserMutationVariables>(UpdateUserDocument, options);
      }
export type UpdateUserMutationHookResult = ReturnType<typeof useUpdateUserMutation>;
export type UpdateUserMutationResult = Apollo.MutationResult<UpdateUserMutation>;
export type UpdateUserMutationOptions = Apollo.BaseMutationOptions<UpdateUserMutation, UpdateUserMutationVariables>;
export const ListContentByNamespaceDocument = gql`
    query ListContentByNamespace($namespace: String!) {
  listContentByNamespace(namespace: $namespace) {
    id
    content
    name
    slug
    created_at
    updated_at
    created_by {
      id
      email
      first_name
      last_name
    }
    updated_by {
      id
      email
      first_name
      last_name
    }
  }
}
    `;

/**
 * __useListContentByNamespaceQuery__
 *
 * To run a query within a React component, call `useListContentByNamespaceQuery` and pass it any options that fit your needs.
 * When your component renders, `useListContentByNamespaceQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useListContentByNamespaceQuery({
 *   variables: {
 *      namespace: // value for 'namespace'
 *   },
 * });
 */
export function useListContentByNamespaceQuery(baseOptions: Apollo.QueryHookOptions<ListContentByNamespaceQuery, ListContentByNamespaceQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ListContentByNamespaceQuery, ListContentByNamespaceQueryVariables>(ListContentByNamespaceDocument, options);
      }
export function useListContentByNamespaceLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ListContentByNamespaceQuery, ListContentByNamespaceQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ListContentByNamespaceQuery, ListContentByNamespaceQueryVariables>(ListContentByNamespaceDocument, options);
        }
export type ListContentByNamespaceQueryHookResult = ReturnType<typeof useListContentByNamespaceQuery>;
export type ListContentByNamespaceLazyQueryHookResult = ReturnType<typeof useListContentByNamespaceLazyQuery>;
export type ListContentByNamespaceQueryResult = Apollo.QueryResult<ListContentByNamespaceQuery, ListContentByNamespaceQueryVariables>;
export const ContentBySlugDocument = gql`
    query ContentBySlug($slug: String!) {
  contentBySlug(slug: $slug) {
    id
    content
    name
    slug
    created_at
    updated_at
    created_by {
      id
      email
      first_name
      last_name
    }
    updated_by {
      id
      email
      first_name
      last_name
    }
  }
}
    `;

/**
 * __useContentBySlugQuery__
 *
 * To run a query within a React component, call `useContentBySlugQuery` and pass it any options that fit your needs.
 * When your component renders, `useContentBySlugQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useContentBySlugQuery({
 *   variables: {
 *      slug: // value for 'slug'
 *   },
 * });
 */
export function useContentBySlugQuery(baseOptions: Apollo.QueryHookOptions<ContentBySlugQuery, ContentBySlugQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ContentBySlugQuery, ContentBySlugQueryVariables>(ContentBySlugDocument, options);
      }
export function useContentBySlugLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ContentBySlugQuery, ContentBySlugQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ContentBySlugQuery, ContentBySlugQueryVariables>(ContentBySlugDocument, options);
        }
export type ContentBySlugQueryHookResult = ReturnType<typeof useContentBySlugQuery>;
export type ContentBySlugLazyQueryHookResult = ReturnType<typeof useContentBySlugLazyQuery>;
export type ContentBySlugQueryResult = Apollo.QueryResult<ContentBySlugQuery, ContentBySlugQueryVariables>;
export const ContentDocument = gql`
    query Content($id: ID!) {
  content(id: $id) {
    id
    content
    name
    slug
    namespace
    created_at
    updated_at
    created_by {
      id
      email
      first_name
      last_name
    }
    updated_by {
      id
      email
      first_name
      last_name
    }
  }
}
    `;

/**
 * __useContentQuery__
 *
 * To run a query within a React component, call `useContentQuery` and pass it any options that fit your needs.
 * When your component renders, `useContentQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useContentQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useContentQuery(baseOptions: Apollo.QueryHookOptions<ContentQuery, ContentQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ContentQuery, ContentQueryVariables>(ContentDocument, options);
      }
export function useContentLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ContentQuery, ContentQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ContentQuery, ContentQueryVariables>(ContentDocument, options);
        }
export type ContentQueryHookResult = ReturnType<typeof useContentQuery>;
export type ContentLazyQueryHookResult = ReturnType<typeof useContentLazyQuery>;
export type ContentQueryResult = Apollo.QueryResult<ContentQuery, ContentQueryVariables>;
export const ListOrganismsDocument = gql`
    query ListOrganisms {
  listOrganisms {
    taxon_id
    scientific_name
    citations {
      title
      authors
      pubmed_id
      journal
    }
    downloads {
      title
      items {
        title
        url
      }
    }
  }
}
    `;

/**
 * __useListOrganismsQuery__
 *
 * To run a query within a React component, call `useListOrganismsQuery` and pass it any options that fit your needs.
 * When your component renders, `useListOrganismsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useListOrganismsQuery({
 *   variables: {
 *   },
 * });
 */
export function useListOrganismsQuery(baseOptions?: Apollo.QueryHookOptions<ListOrganismsQuery, ListOrganismsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ListOrganismsQuery, ListOrganismsQueryVariables>(ListOrganismsDocument, options);
      }
export function useListOrganismsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ListOrganismsQuery, ListOrganismsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ListOrganismsQuery, ListOrganismsQueryVariables>(ListOrganismsDocument, options);
        }
export type ListOrganismsQueryHookResult = ReturnType<typeof useListOrganismsQuery>;
export type ListOrganismsLazyQueryHookResult = ReturnType<typeof useListOrganismsLazyQuery>;
export type ListOrganismsQueryResult = Apollo.QueryResult<ListOrganismsQuery, ListOrganismsQueryVariables>;
export const GeneSummaryDocument = gql`
    query GeneSummary($gene: String!) {
  geneGeneralInformation(gene: $gene) {
    id
    name_description
    gene_product
    synonyms
    description
  }
  geneOntologyAnnotation(gene: $gene) {
    id
    type
    date
    go_term
    evidence_code
    with {
      id
      db
      name
    }
    extensions {
      id
      db
      relation
      name
    }
  }
  listPublicationsWithGene(gene: $gene) {
    id
    title
    journal
    pages
    issue
    authors {
      last_name
    }
  }
}
    `;

/**
 * __useGeneSummaryQuery__
 *
 * To run a query within a React component, call `useGeneSummaryQuery` and pass it any options that fit your needs.
 * When your component renders, `useGeneSummaryQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGeneSummaryQuery({
 *   variables: {
 *      gene: // value for 'gene'
 *   },
 * });
 */
export function useGeneSummaryQuery(baseOptions: Apollo.QueryHookOptions<GeneSummaryQuery, GeneSummaryQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GeneSummaryQuery, GeneSummaryQueryVariables>(GeneSummaryDocument, options);
      }
export function useGeneSummaryLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GeneSummaryQuery, GeneSummaryQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GeneSummaryQuery, GeneSummaryQueryVariables>(GeneSummaryDocument, options);
        }
export type GeneSummaryQueryHookResult = ReturnType<typeof useGeneSummaryQuery>;
export type GeneSummaryLazyQueryHookResult = ReturnType<typeof useGeneSummaryLazyQuery>;
export type GeneSummaryQueryResult = Apollo.QueryResult<GeneSummaryQuery, GeneSummaryQueryVariables>;
export const GeneOntologyAnnotationDocument = gql`
    query GeneOntologyAnnotation($gene: String!) {
  geneOntologyAnnotation(gene: $gene) {
    id
    type
    date
    go_term
    evidence_code
    qualifier
    publication
    assigned_by
    with {
      id
      db
      name
    }
    extensions {
      id
      db
      relation
      name
    }
  }
}
    `;

/**
 * __useGeneOntologyAnnotationQuery__
 *
 * To run a query within a React component, call `useGeneOntologyAnnotationQuery` and pass it any options that fit your needs.
 * When your component renders, `useGeneOntologyAnnotationQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGeneOntologyAnnotationQuery({
 *   variables: {
 *      gene: // value for 'gene'
 *   },
 * });
 */
export function useGeneOntologyAnnotationQuery(baseOptions: Apollo.QueryHookOptions<GeneOntologyAnnotationQuery, GeneOntologyAnnotationQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GeneOntologyAnnotationQuery, GeneOntologyAnnotationQueryVariables>(GeneOntologyAnnotationDocument, options);
      }
export function useGeneOntologyAnnotationLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GeneOntologyAnnotationQuery, GeneOntologyAnnotationQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GeneOntologyAnnotationQuery, GeneOntologyAnnotationQueryVariables>(GeneOntologyAnnotationDocument, options);
        }
export type GeneOntologyAnnotationQueryHookResult = ReturnType<typeof useGeneOntologyAnnotationQuery>;
export type GeneOntologyAnnotationLazyQueryHookResult = ReturnType<typeof useGeneOntologyAnnotationLazyQuery>;
export type GeneOntologyAnnotationQueryResult = Apollo.QueryResult<GeneOntologyAnnotationQuery, GeneOntologyAnnotationQueryVariables>;
export const ListStrainsWithGeneDocument = gql`
    query ListStrainsWithGene($gene: String!) {
  listStrainsWithGene(gene: $gene) {
    id
    label
    characteristics
    in_stock
    phenotypes {
      phenotype
      publication {
        id
        title
        journal
        pages
        volume
        pub_date
        authors {
          last_name
          rank
        }
      }
    }
  }
}
    `;

/**
 * __useListStrainsWithGeneQuery__
 *
 * To run a query within a React component, call `useListStrainsWithGeneQuery` and pass it any options that fit your needs.
 * When your component renders, `useListStrainsWithGeneQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useListStrainsWithGeneQuery({
 *   variables: {
 *      gene: // value for 'gene'
 *   },
 * });
 */
export function useListStrainsWithGeneQuery(baseOptions: Apollo.QueryHookOptions<ListStrainsWithGeneQuery, ListStrainsWithGeneQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ListStrainsWithGeneQuery, ListStrainsWithGeneQueryVariables>(ListStrainsWithGeneDocument, options);
      }
export function useListStrainsWithGeneLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ListStrainsWithGeneQuery, ListStrainsWithGeneQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ListStrainsWithGeneQuery, ListStrainsWithGeneQueryVariables>(ListStrainsWithGeneDocument, options);
        }
export type ListStrainsWithGeneQueryHookResult = ReturnType<typeof useListStrainsWithGeneQuery>;
export type ListStrainsWithGeneLazyQueryHookResult = ReturnType<typeof useListStrainsWithGeneLazyQuery>;
export type ListStrainsWithGeneQueryResult = Apollo.QueryResult<ListStrainsWithGeneQuery, ListStrainsWithGeneQueryVariables>;
export const PublicationDocument = gql`
    query Publication($id: ID!) {
  publication(id: $id) {
    id
    doi
    title
    abstract
    journal
    pub_date
    pages
    issue
    volume
    authors {
      initials
      last_name
    }
  }
}
    `;

/**
 * __usePublicationQuery__
 *
 * To run a query within a React component, call `usePublicationQuery` and pass it any options that fit your needs.
 * When your component renders, `usePublicationQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = usePublicationQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function usePublicationQuery(baseOptions: Apollo.QueryHookOptions<PublicationQuery, PublicationQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<PublicationQuery, PublicationQueryVariables>(PublicationDocument, options);
      }
export function usePublicationLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<PublicationQuery, PublicationQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<PublicationQuery, PublicationQueryVariables>(PublicationDocument, options);
        }
export type PublicationQueryHookResult = ReturnType<typeof usePublicationQuery>;
export type PublicationLazyQueryHookResult = ReturnType<typeof usePublicationLazyQuery>;
export type PublicationQueryResult = Apollo.QueryResult<PublicationQuery, PublicationQueryVariables>;
export const ListRecentPublicationsDocument = gql`
    query ListRecentPublications($limit: Int! = 4) {
  listRecentPublications(limit: $limit) {
    id
    doi
    title
    abstract
    journal
    pub_date
    pages
    issue
    volume
    authors {
      initials
      last_name
    }
  }
}
    `;

/**
 * __useListRecentPublicationsQuery__
 *
 * To run a query within a React component, call `useListRecentPublicationsQuery` and pass it any options that fit your needs.
 * When your component renders, `useListRecentPublicationsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useListRecentPublicationsQuery({
 *   variables: {
 *      limit: // value for 'limit'
 *   },
 * });
 */
export function useListRecentPublicationsQuery(baseOptions?: Apollo.QueryHookOptions<ListRecentPublicationsQuery, ListRecentPublicationsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ListRecentPublicationsQuery, ListRecentPublicationsQueryVariables>(ListRecentPublicationsDocument, options);
      }
export function useListRecentPublicationsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ListRecentPublicationsQuery, ListRecentPublicationsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ListRecentPublicationsQuery, ListRecentPublicationsQueryVariables>(ListRecentPublicationsDocument, options);
        }
export type ListRecentPublicationsQueryHookResult = ReturnType<typeof useListRecentPublicationsQuery>;
export type ListRecentPublicationsLazyQueryHookResult = ReturnType<typeof useListRecentPublicationsLazyQuery>;
export type ListRecentPublicationsQueryResult = Apollo.QueryResult<ListRecentPublicationsQuery, ListRecentPublicationsQueryVariables>;
export const ListPublicationsWithGeneDocument = gql`
    query ListPublicationsWithGene($gene: String!) {
  listPublicationsWithGene(gene: $gene) {
    related_genes {
      id
      name
    }
    id
    doi
    title
    journal
    pub_date
    volume
    pages
    pub_type
    source
    issue
    authors {
      last_name
      rank
    }
  }
}
    `;

/**
 * __useListPublicationsWithGeneQuery__
 *
 * To run a query within a React component, call `useListPublicationsWithGeneQuery` and pass it any options that fit your needs.
 * When your component renders, `useListPublicationsWithGeneQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useListPublicationsWithGeneQuery({
 *   variables: {
 *      gene: // value for 'gene'
 *   },
 * });
 */
export function useListPublicationsWithGeneQuery(baseOptions: Apollo.QueryHookOptions<ListPublicationsWithGeneQuery, ListPublicationsWithGeneQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ListPublicationsWithGeneQuery, ListPublicationsWithGeneQueryVariables>(ListPublicationsWithGeneDocument, options);
      }
export function useListPublicationsWithGeneLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ListPublicationsWithGeneQuery, ListPublicationsWithGeneQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ListPublicationsWithGeneQuery, ListPublicationsWithGeneQueryVariables>(ListPublicationsWithGeneDocument, options);
        }
export type ListPublicationsWithGeneQueryHookResult = ReturnType<typeof useListPublicationsWithGeneQuery>;
export type ListPublicationsWithGeneLazyQueryHookResult = ReturnType<typeof useListPublicationsWithGeneLazyQuery>;
export type ListPublicationsWithGeneQueryResult = Apollo.QueryResult<ListPublicationsWithGeneQuery, ListPublicationsWithGeneQueryVariables>;
export const StrainListDocument = gql`
    query StrainList($cursor: Int!, $limit: Int!, $filter: StrainListFilter) {
  listStrains(cursor: $cursor, limit: $limit, filter: $filter) {
    nextCursor
    totalCount
    strains {
      id
      label
      summary
      in_stock
    }
  }
}
    `;

/**
 * __useStrainListQuery__
 *
 * To run a query within a React component, call `useStrainListQuery` and pass it any options that fit your needs.
 * When your component renders, `useStrainListQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useStrainListQuery({
 *   variables: {
 *      cursor: // value for 'cursor'
 *      limit: // value for 'limit'
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useStrainListQuery(baseOptions: Apollo.QueryHookOptions<StrainListQuery, StrainListQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<StrainListQuery, StrainListQueryVariables>(StrainListDocument, options);
      }
export function useStrainListLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<StrainListQuery, StrainListQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<StrainListQuery, StrainListQueryVariables>(StrainListDocument, options);
        }
export type StrainListQueryHookResult = ReturnType<typeof useStrainListQuery>;
export type StrainListLazyQueryHookResult = ReturnType<typeof useStrainListLazyQuery>;
export type StrainListQueryResult = Apollo.QueryResult<StrainListQuery, StrainListQueryVariables>;
export const ListStrainsWithPhenotypeDocument = gql`
    query ListStrainsWithPhenotype($cursor: Int!, $limit: Int!, $type: String!, $annotation: String!) {
  listStrainsWithAnnotation(
    cursor: $cursor
    limit: $limit
    type: $type
    annotation: $annotation
  ) {
    totalCount
    nextCursor
    strains {
      id
      label
      genes {
        name
      }
      publications {
        id
        pub_date
        title
        journal
        volume
        pages
        authors {
          last_name
        }
      }
    }
  }
}
    `;

/**
 * __useListStrainsWithPhenotypeQuery__
 *
 * To run a query within a React component, call `useListStrainsWithPhenotypeQuery` and pass it any options that fit your needs.
 * When your component renders, `useListStrainsWithPhenotypeQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useListStrainsWithPhenotypeQuery({
 *   variables: {
 *      cursor: // value for 'cursor'
 *      limit: // value for 'limit'
 *      type: // value for 'type'
 *      annotation: // value for 'annotation'
 *   },
 * });
 */
export function useListStrainsWithPhenotypeQuery(baseOptions: Apollo.QueryHookOptions<ListStrainsWithPhenotypeQuery, ListStrainsWithPhenotypeQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ListStrainsWithPhenotypeQuery, ListStrainsWithPhenotypeQueryVariables>(ListStrainsWithPhenotypeDocument, options);
      }
export function useListStrainsWithPhenotypeLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ListStrainsWithPhenotypeQuery, ListStrainsWithPhenotypeQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ListStrainsWithPhenotypeQuery, ListStrainsWithPhenotypeQueryVariables>(ListStrainsWithPhenotypeDocument, options);
        }
export type ListStrainsWithPhenotypeQueryHookResult = ReturnType<typeof useListStrainsWithPhenotypeQuery>;
export type ListStrainsWithPhenotypeLazyQueryHookResult = ReturnType<typeof useListStrainsWithPhenotypeLazyQuery>;
export type ListStrainsWithPhenotypeQueryResult = Apollo.QueryResult<ListStrainsWithPhenotypeQuery, ListStrainsWithPhenotypeQueryVariables>;
export const ListBacterialStrainsDocument = gql`
    query ListBacterialStrains {
  bacterialFoodSource: listStrainsWithAnnotation(
    cursor: 0
    limit: 100
    type: "characteristic"
    annotation: "bacterial food source"
  ) {
    totalCount
    nextCursor
    strains {
      id
      label
      summary
      in_stock
    }
  }
  symbioticFarmerBacterium: listStrainsWithAnnotation(
    cursor: 0
    limit: 100
    type: "characteristic"
    annotation: "symbiotic farmer bacterium"
  ) {
    totalCount
    nextCursor
    strains {
      id
      label
      summary
      in_stock
    }
  }
}
    `;

/**
 * __useListBacterialStrainsQuery__
 *
 * To run a query within a React component, call `useListBacterialStrainsQuery` and pass it any options that fit your needs.
 * When your component renders, `useListBacterialStrainsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useListBacterialStrainsQuery({
 *   variables: {
 *   },
 * });
 */
export function useListBacterialStrainsQuery(baseOptions?: Apollo.QueryHookOptions<ListBacterialStrainsQuery, ListBacterialStrainsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ListBacterialStrainsQuery, ListBacterialStrainsQueryVariables>(ListBacterialStrainsDocument, options);
      }
export function useListBacterialStrainsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ListBacterialStrainsQuery, ListBacterialStrainsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ListBacterialStrainsQuery, ListBacterialStrainsQueryVariables>(ListBacterialStrainsDocument, options);
        }
export type ListBacterialStrainsQueryHookResult = ReturnType<typeof useListBacterialStrainsQuery>;
export type ListBacterialStrainsLazyQueryHookResult = ReturnType<typeof useListBacterialStrainsLazyQuery>;
export type ListBacterialStrainsQueryResult = Apollo.QueryResult<ListBacterialStrainsQuery, ListBacterialStrainsQueryVariables>;
export const ListStrainsInventoryDocument = gql`
    query ListStrainsInventory($cursor: Int!, $limit: Int!) {
  listStrainsWithAnnotation(
    cursor: $cursor
    limit: $limit
    type: "strain_inventory"
    annotation: "strain_inventory"
  ) {
    totalCount
    nextCursor
    strains {
      id
      label
      summary
      in_stock
    }
  }
}
    `;

/**
 * __useListStrainsInventoryQuery__
 *
 * To run a query within a React component, call `useListStrainsInventoryQuery` and pass it any options that fit your needs.
 * When your component renders, `useListStrainsInventoryQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useListStrainsInventoryQuery({
 *   variables: {
 *      cursor: // value for 'cursor'
 *      limit: // value for 'limit'
 *   },
 * });
 */
export function useListStrainsInventoryQuery(baseOptions: Apollo.QueryHookOptions<ListStrainsInventoryQuery, ListStrainsInventoryQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ListStrainsInventoryQuery, ListStrainsInventoryQueryVariables>(ListStrainsInventoryDocument, options);
      }
export function useListStrainsInventoryLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ListStrainsInventoryQuery, ListStrainsInventoryQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ListStrainsInventoryQuery, ListStrainsInventoryQueryVariables>(ListStrainsInventoryDocument, options);
        }
export type ListStrainsInventoryQueryHookResult = ReturnType<typeof useListStrainsInventoryQuery>;
export type ListStrainsInventoryLazyQueryHookResult = ReturnType<typeof useListStrainsInventoryLazyQuery>;
export type ListStrainsInventoryQueryResult = Apollo.QueryResult<ListStrainsInventoryQuery, ListStrainsInventoryQueryVariables>;
export const ListPlasmidsInventoryDocument = gql`
    query ListPlasmidsInventory($cursor: Int!, $limit: Int!) {
  listPlasmidsWithAnnotation(
    cursor: $cursor
    limit: $limit
    type: "plasmid_inventory"
    annotation: "plasmid inventory"
  ) {
    totalCount
    nextCursor
    plasmids {
      id
      name
      summary
      in_stock
    }
  }
}
    `;

/**
 * __useListPlasmidsInventoryQuery__
 *
 * To run a query within a React component, call `useListPlasmidsInventoryQuery` and pass it any options that fit your needs.
 * When your component renders, `useListPlasmidsInventoryQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useListPlasmidsInventoryQuery({
 *   variables: {
 *      cursor: // value for 'cursor'
 *      limit: // value for 'limit'
 *   },
 * });
 */
export function useListPlasmidsInventoryQuery(baseOptions: Apollo.QueryHookOptions<ListPlasmidsInventoryQuery, ListPlasmidsInventoryQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ListPlasmidsInventoryQuery, ListPlasmidsInventoryQueryVariables>(ListPlasmidsInventoryDocument, options);
      }
export function useListPlasmidsInventoryLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ListPlasmidsInventoryQuery, ListPlasmidsInventoryQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ListPlasmidsInventoryQuery, ListPlasmidsInventoryQueryVariables>(ListPlasmidsInventoryDocument, options);
        }
export type ListPlasmidsInventoryQueryHookResult = ReturnType<typeof useListPlasmidsInventoryQuery>;
export type ListPlasmidsInventoryLazyQueryHookResult = ReturnType<typeof useListPlasmidsInventoryLazyQuery>;
export type ListPlasmidsInventoryQueryResult = Apollo.QueryResult<ListPlasmidsInventoryQuery, ListPlasmidsInventoryQueryVariables>;
export const PlasmidListFilterDocument = gql`
    query PlasmidListFilter($cursor: Int!, $limit: Int!, $filter: String!) {
  listPlasmids(cursor: $cursor, limit: $limit, filter: $filter) {
    nextCursor
    totalCount
    plasmids {
      id
      name
      summary
      in_stock
    }
  }
}
    `;

/**
 * __usePlasmidListFilterQuery__
 *
 * To run a query within a React component, call `usePlasmidListFilterQuery` and pass it any options that fit your needs.
 * When your component renders, `usePlasmidListFilterQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = usePlasmidListFilterQuery({
 *   variables: {
 *      cursor: // value for 'cursor'
 *      limit: // value for 'limit'
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function usePlasmidListFilterQuery(baseOptions: Apollo.QueryHookOptions<PlasmidListFilterQuery, PlasmidListFilterQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<PlasmidListFilterQuery, PlasmidListFilterQueryVariables>(PlasmidListFilterDocument, options);
      }
export function usePlasmidListFilterLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<PlasmidListFilterQuery, PlasmidListFilterQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<PlasmidListFilterQuery, PlasmidListFilterQueryVariables>(PlasmidListFilterDocument, options);
        }
export type PlasmidListFilterQueryHookResult = ReturnType<typeof usePlasmidListFilterQuery>;
export type PlasmidListFilterLazyQueryHookResult = ReturnType<typeof usePlasmidListFilterLazyQuery>;
export type PlasmidListFilterQueryResult = Apollo.QueryResult<PlasmidListFilterQuery, PlasmidListFilterQueryVariables>;
export const PlasmidDocument = gql`
    query Plasmid($id: ID!) {
  plasmid(id: $id) {
    id
    name
    summary
    depositor {
      first_name
      last_name
    }
    publications {
      id
      pub_date
      title
      journal
      volume
      pages
      doi
      authors {
        last_name
      }
    }
    dbxrefs
    genes {
      name
    }
    image_map
    sequence
    keywords
    genbank_accession
    in_stock
  }
}
    `;

/**
 * __usePlasmidQuery__
 *
 * To run a query within a React component, call `usePlasmidQuery` and pass it any options that fit your needs.
 * When your component renders, `usePlasmidQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = usePlasmidQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function usePlasmidQuery(baseOptions: Apollo.QueryHookOptions<PlasmidQuery, PlasmidQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<PlasmidQuery, PlasmidQueryVariables>(PlasmidDocument, options);
      }
export function usePlasmidLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<PlasmidQuery, PlasmidQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<PlasmidQuery, PlasmidQueryVariables>(PlasmidDocument, options);
        }
export type PlasmidQueryHookResult = ReturnType<typeof usePlasmidQuery>;
export type PlasmidLazyQueryHookResult = ReturnType<typeof usePlasmidLazyQuery>;
export type PlasmidQueryResult = Apollo.QueryResult<PlasmidQuery, PlasmidQueryVariables>;
export const StrainDocument = gql`
    query Strain($id: ID!) {
  strain(id: $id) {
    id
    label
    summary
    species
    parent {
      id
      label
    }
    depositor {
      first_name
      last_name
    }
    plasmid
    dbxrefs
    publications {
      id
      pub_date
      title
      journal
      volume
      pages
      doi
      authors {
        last_name
      }
    }
    genes {
      name
    }
    in_stock
    systematic_name
    genotypes
    mutagenesis_method
    genetic_modification
    names
    characteristics
    phenotypes {
      phenotype
      note
      assay
      environment
      publication {
        id
        pub_date
        title
        journal
        volume
        pages
        authors {
          last_name
        }
      }
    }
  }
}
    `;

/**
 * __useStrainQuery__
 *
 * To run a query within a React component, call `useStrainQuery` and pass it any options that fit your needs.
 * When your component renders, `useStrainQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useStrainQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useStrainQuery(baseOptions: Apollo.QueryHookOptions<StrainQuery, StrainQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<StrainQuery, StrainQueryVariables>(StrainDocument, options);
      }
export function useStrainLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<StrainQuery, StrainQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<StrainQuery, StrainQueryVariables>(StrainDocument, options);
        }
export type StrainQueryHookResult = ReturnType<typeof useStrainQuery>;
export type StrainLazyQueryHookResult = ReturnType<typeof useStrainLazyQuery>;
export type StrainQueryResult = Apollo.QueryResult<StrainQuery, StrainQueryVariables>;
export const ListRecentPlasmidsDocument = gql`
    query ListRecentPlasmids($limit: Int! = 4) {
  listRecentPlasmids(limit: $limit) {
    id
    created_at
    name
  }
}
    `;

/**
 * __useListRecentPlasmidsQuery__
 *
 * To run a query within a React component, call `useListRecentPlasmidsQuery` and pass it any options that fit your needs.
 * When your component renders, `useListRecentPlasmidsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useListRecentPlasmidsQuery({
 *   variables: {
 *      limit: // value for 'limit'
 *   },
 * });
 */
export function useListRecentPlasmidsQuery(baseOptions?: Apollo.QueryHookOptions<ListRecentPlasmidsQuery, ListRecentPlasmidsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ListRecentPlasmidsQuery, ListRecentPlasmidsQueryVariables>(ListRecentPlasmidsDocument, options);
      }
export function useListRecentPlasmidsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ListRecentPlasmidsQuery, ListRecentPlasmidsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ListRecentPlasmidsQuery, ListRecentPlasmidsQueryVariables>(ListRecentPlasmidsDocument, options);
        }
export type ListRecentPlasmidsQueryHookResult = ReturnType<typeof useListRecentPlasmidsQuery>;
export type ListRecentPlasmidsLazyQueryHookResult = ReturnType<typeof useListRecentPlasmidsLazyQuery>;
export type ListRecentPlasmidsQueryResult = Apollo.QueryResult<ListRecentPlasmidsQuery, ListRecentPlasmidsQueryVariables>;
export const ListRecentStrainsDocument = gql`
    query ListRecentStrains($limit: Int! = 4) {
  listRecentStrains(limit: $limit) {
    id
    created_at
    systematic_name
  }
}
    `;

/**
 * __useListRecentStrainsQuery__
 *
 * To run a query within a React component, call `useListRecentStrainsQuery` and pass it any options that fit your needs.
 * When your component renders, `useListRecentStrainsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useListRecentStrainsQuery({
 *   variables: {
 *      limit: // value for 'limit'
 *   },
 * });
 */
export function useListRecentStrainsQuery(baseOptions?: Apollo.QueryHookOptions<ListRecentStrainsQuery, ListRecentStrainsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ListRecentStrainsQuery, ListRecentStrainsQueryVariables>(ListRecentStrainsDocument, options);
      }
export function useListRecentStrainsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ListRecentStrainsQuery, ListRecentStrainsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ListRecentStrainsQuery, ListRecentStrainsQueryVariables>(ListRecentStrainsDocument, options);
        }
export type ListRecentStrainsQueryHookResult = ReturnType<typeof useListRecentStrainsQuery>;
export type ListRecentStrainsLazyQueryHookResult = ReturnType<typeof useListRecentStrainsLazyQuery>;
export type ListRecentStrainsQueryResult = Apollo.QueryResult<ListRecentStrainsQuery, ListRecentStrainsQueryVariables>;
export const UserByEmailDocument = gql`
    query UserByEmail($email: String!) {
  userByEmail(email: $email) {
    id
  }
}
    `;

/**
 * __useUserByEmailQuery__
 *
 * To run a query within a React component, call `useUserByEmailQuery` and pass it any options that fit your needs.
 * When your component renders, `useUserByEmailQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useUserByEmailQuery({
 *   variables: {
 *      email: // value for 'email'
 *   },
 * });
 */
export function useUserByEmailQuery(baseOptions: Apollo.QueryHookOptions<UserByEmailQuery, UserByEmailQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<UserByEmailQuery, UserByEmailQueryVariables>(UserByEmailDocument, options);
      }
export function useUserByEmailLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<UserByEmailQuery, UserByEmailQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<UserByEmailQuery, UserByEmailQueryVariables>(UserByEmailDocument, options);
        }
export type UserByEmailQueryHookResult = ReturnType<typeof useUserByEmailQuery>;
export type UserByEmailLazyQueryHookResult = ReturnType<typeof useUserByEmailLazyQuery>;
export type UserByEmailQueryResult = Apollo.QueryResult<UserByEmailQuery, UserByEmailQueryVariables>;