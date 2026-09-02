"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UserByEmail = exports.ListPhenotypeAssays = exports.ListPhenotypeEnvironments = exports.ListPhenotypes = exports.ListRecentStrains = exports.ListRecentPlasmids = exports.Strain = exports.Plasmid = exports.PlasmidListFilter = exports.ListPlasmidsInventory = exports.ListStrainsInventory = exports.ListBacterialStrains = exports.ListStrainsWithPhenotype = exports.StrainList = exports.ListPublicationsWithGene = exports.ListRecentPublications = exports.Publication = exports.ListStrainsWithGene = exports.GeneOntologyAnnotation = exports.ListPublicationsWithGeneSummary = exports.GeneOntologyAnnotationSummary = exports.GeneGeneralInformationSummary = exports.ListOrganisms = exports.Content = exports.ContentBySlug = exports.ListContentByNamespace = exports.UpdateUser = exports.CreateUser = exports.UploadFile = exports.UpdateStrainPhenotype = exports.AddStrainPhenotype = exports.CreateOrder = exports.UpdateGeneGeneralInfo = exports.CreateGeneGeneralInfo = exports.DeleteContent = exports.UpdateContent = exports.CreateContent = exports.Logout = exports.Login = void 0;
const graphql_tag_1 = require("graphql-tag");
exports.Login = (0, graphql_tag_1.default) `
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
exports.Logout = (0, graphql_tag_1.default) `
    mutation Logout {
  logout {
    success
  }
}
    `;
exports.CreateContent = (0, graphql_tag_1.default) `
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
exports.UpdateContent = (0, graphql_tag_1.default) `
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
exports.DeleteContent = (0, graphql_tag_1.default) `
    mutation DeleteContent($id: ID!) {
  deleteContent(id: $id) {
    success
  }
}
    `;
exports.CreateGeneGeneralInfo = (0, graphql_tag_1.default) `
    mutation CreateGeneGeneralInfo($id: ID!, $input: CreateGeneGeneralInfoInput!) {
  createGeneGeneralInfo(id: $id, input: $input) {
    id
    created_by {
      id
    }
  }
}
    `;
exports.UpdateGeneGeneralInfo = (0, graphql_tag_1.default) `
    mutation UpdateGeneGeneralInfo($id: ID!, $input: UpdateGeneGeneralInfoInput!) {
  updateGeneGeneralInfo(id: $id, input: $input) {
    id
    updated_by {
      id
    }
  }
}
    `;
exports.CreateOrder = (0, graphql_tag_1.default) `
    mutation CreateOrder($input: CreateOrderInput!) {
  createOrder(input: $input) {
    id
  }
}
    `;
exports.AddStrainPhenotype = (0, graphql_tag_1.default) `
    mutation AddStrainPhenotype($strainId: ID!, $input: AddStrainPhenotypeInput!) {
  addStrainPhenotype(strainId: $strainId, input: $input) {
    id
    label
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
exports.UpdateStrainPhenotype = (0, graphql_tag_1.default) `
    mutation UpdateStrainPhenotype($strainId: ID!, $target: UpdateStrainPhenotypeTargetInput!, $payload: UpdateStrainPhenotypePayloadInput!) {
  updateStrainPhenotype(strainId: $strainId, target: $target, payload: $payload) {
    id
    label
    phenotypes {
      phenotype
      environment
      assay
      note
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
exports.UploadFile = (0, graphql_tag_1.default) `
    mutation UploadFile($file: Upload!) {
  uploadFile(file: $file) {
    url
  }
}
    `;
exports.CreateUser = (0, graphql_tag_1.default) `
    mutation CreateUser($input: CreateUserInput!) {
  createUser(input: $input) {
    id
  }
}
    `;
exports.UpdateUser = (0, graphql_tag_1.default) `
    mutation UpdateUser($id: ID!, $input: UpdateUserInput!) {
  updateUser(id: $id, input: $input) {
    id
  }
}
    `;
exports.ListContentByNamespace = (0, graphql_tag_1.default) `
    query ListContentByNamespace($namespace: String!, $limit: Int) {
  listContentByNamespace(namespace: $namespace, limit: $limit) {
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
exports.ContentBySlug = (0, graphql_tag_1.default) `
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
exports.Content = (0, graphql_tag_1.default) `
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
exports.ListOrganisms = (0, graphql_tag_1.default) `
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
exports.GeneGeneralInformationSummary = (0, graphql_tag_1.default) `
    query GeneGeneralInformationSummary($gene: String!) {
  geneGeneralInformation(gene: $gene) {
    id
    name_description
    gene_product
    synonyms
    description
  }
}
    `;
exports.GeneOntologyAnnotationSummary = (0, graphql_tag_1.default) `
    query GeneOntologyAnnotationSummary($gene: String!) {
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
}
    `;
exports.ListPublicationsWithGeneSummary = (0, graphql_tag_1.default) `
    query ListPublicationsWithGeneSummary($gene: String!) {
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
exports.GeneOntologyAnnotation = (0, graphql_tag_1.default) `
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
exports.ListStrainsWithGene = (0, graphql_tag_1.default) `
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
exports.Publication = (0, graphql_tag_1.default) `
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
exports.ListRecentPublications = (0, graphql_tag_1.default) `
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
exports.ListPublicationsWithGene = (0, graphql_tag_1.default) `
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
exports.StrainList = (0, graphql_tag_1.default) `
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
exports.ListStrainsWithPhenotype = (0, graphql_tag_1.default) `
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
exports.ListBacterialStrains = (0, graphql_tag_1.default) `
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
exports.ListStrainsInventory = (0, graphql_tag_1.default) `
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
exports.ListPlasmidsInventory = (0, graphql_tag_1.default) `
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
exports.PlasmidListFilter = (0, graphql_tag_1.default) `
    query PlasmidListFilter($cursor: Int!, $limit: Int!, $filter: PlasmidListFilter) {
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
exports.Plasmid = (0, graphql_tag_1.default) `
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
exports.Strain = (0, graphql_tag_1.default) `
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
exports.ListRecentPlasmids = (0, graphql_tag_1.default) `
    query ListRecentPlasmids($limit: Int! = 4) {
  listRecentPlasmids(limit: $limit) {
    id
    created_at
    name
  }
}
    `;
exports.ListRecentStrains = (0, graphql_tag_1.default) `
    query ListRecentStrains($limit: Int! = 4) {
  listRecentStrains(limit: $limit) {
    id
    created_at
    systematic_name
  }
}
    `;
exports.ListPhenotypes = (0, graphql_tag_1.default) `
    query ListPhenotypes($search: String!) {
  listPhenotypes(search: $search)
}
    `;
exports.ListPhenotypeEnvironments = (0, graphql_tag_1.default) `
    query ListPhenotypeEnvironments($search: String!) {
  listPhenotypeEnvironments(search: $search)
}
    `;
exports.ListPhenotypeAssays = (0, graphql_tag_1.default) `
    query ListPhenotypeAssays($search: String!) {
  listPhenotypeAssays(search: $search)
}
    `;
exports.UserByEmail = (0, graphql_tag_1.default) `
    query UserByEmail($email: String!) {
  userByEmail(email: $email) {
    id
  }
}
    `;
