import gql from 'graphql-tag';

export const Login = gql`
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
export const Logout = gql`
    mutation Logout {
  logout {
    success
  }
}
    `;
export const CreateContent = gql`
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
export const UpdateContent = gql`
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
export const DeleteContent = gql`
    mutation DeleteContent($id: ID!) {
  deleteContent(id: $id) {
    success
  }
}
    `;
export const CreateOrder = gql`
    mutation CreateOrder($input: CreateOrderInput!) {
  createOrder(input: $input) {
    id
  }
}
    `;
export const UploadFile = gql`
    mutation UploadFile($file: Upload!) {
  uploadFile(file: $file) {
    url
  }
}
    `;
export const CreateUser = gql`
    mutation CreateUser($input: CreateUserInput!) {
  createUser(input: $input) {
    id
  }
}
    `;
export const UpdateUser = gql`
    mutation UpdateUser($id: ID!, $input: UpdateUserInput!) {
  updateUser(id: $id, input: $input) {
    id
  }
}
    `;
export const ListContentByNamespace = gql`
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
export const ContentBySlug = gql`
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
export const Content = gql`
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
export const ListOrganisms = gql`
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
export const GeneSummary = gql`
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
export const GeneOntologyAnnotation = gql`
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
export const ListStrainsWithGene = gql`
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
export const Publication = gql`
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
export const ListRecentPublications = gql`
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
export const ListPublicationsWithGene = gql`
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
export const StrainList = gql`
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
export const ListStrainsWithPhenotype = gql`
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
export const ListBacterialStrains = gql`
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
export const ListStrainsInventory = gql`
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
export const ListPlasmidsInventory = gql`
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
export const PlasmidListFilter = gql`
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
export const Plasmid = gql`
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
export const Strain = gql`
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
export const ListRecentPlasmids = gql`
    query ListRecentPlasmids($limit: Int! = 4) {
  listRecentPlasmids(limit: $limit) {
    id
    created_at
    name
  }
}
    `;
export const ListRecentStrains = gql`
    query ListRecentStrains($limit: Int! = 4) {
  listRecentStrains(limit: $limit) {
    id
    created_at
    systematic_name
  }
}
    `;
export const UserByEmail = gql`
    query UserByEmail($email: String!) {
  userByEmail(email: $email) {
    id
  }
}
    `;