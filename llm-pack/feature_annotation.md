# Feature Annotation API Documentation

This document provides an overview of the DictyBase Feature Annotation API based on the protocol buffer definitions.

## Service Definitions

### FeatureAnnotationService

Service for managing feature annotations.

#### Methods:

- **CreateFeatureAnnotation**: Create a new feature annotation
- **GetFeatureAnnotation**: Retrieve a specific feature annotation by ID
- **UpdateFeatureAnnotation**: Update an existing feature annotation (tags are appended)
- **DeleteFeatureAnnotation**: Delete a feature annotation (soft or hard delete)
- **AddTag**: Add a tag to an existing feature annotation
- **UpdateTag**: Update a tag in a feature annotation
- **RemoveTag**: Remove a tag from a feature annotation
- **ListFeatureAnnotationsByPubmedId**: Retrieve annotations by PubMed ID
- **ListFeatureAnnotationsByDOI**: Retrieve annotations by DOI

### OrganismFeatureService

Service for linking feature annotations to organisms.

#### Methods:

- **LinkFeatureToOrganism**: Link a feature annotation to an organism
- **GetFeatureOrganism**: Get organism for a feature annotation
- **UpdateFeatureOrganism**: Update a feature's organism link
- **RemoveFeatureOrganism**: Remove a feature's organism link

## Main Data Structures

### FeatureAnnotation

Represents a complete feature annotation record with all its metadata.

```protobuf
message FeatureAnnotation {
  string type = 1;
  string id = 2;
  FeatureAnnotationAttributes attributes = 3;
  string created_by = 4;
  string updated_by = 5;
  google.protobuf.Timestamp created_at = 6;
  google.protobuf.Timestamp updated_at = 7;
  bool is_obsolete = 8;
  int64 version = 9 [deprecated = true];
}
```

### FeatureAnnotationAttributes

Core properties and metadata of a feature annotation.

```protobuf
message FeatureAnnotationAttributes {
  string name = 1;
  repeated string synonyms = 2;
  repeated string publications = 3;
  repeated string pubmed = 4;
  repeated Dbxref dbxrefs = 5 [deprecated = true];
  repeated DbLink dblinks = 6;
  repeated TagProperty properties = 7;
}
```

### TagProperty

Key-value pair structure for storing custom attributes.

```protobuf
message TagProperty {
  string tag = 1;
  string value = 2;
  string created_by = 3;
  string updated_by = 4;
  google.protobuf.Timestamp created_at = 5;
  google.protobuf.Timestamp updated_at = 6;
}
```

### DbLink

Reference to an external bioinformatics database entry.

```protobuf
message DbLink {
  string primary_id = 1;
  int64 version = 2;
  string database = 3;
  string linktype = 4;
  string url = 5;
  string label = 6;
}
```

### OrganismFeatureLink

Links a feature to an organism.

```protobuf
message OrganismFeatureLink {
  int64 organism_id = 1;
  string feature_id = 2;
}
```

## Usage Examples

### Creating a Feature Annotation

```go
client := feature_annotation.NewFeatureAnnotationServiceClient(conn)
now := timestamppb.Now()

newAnnotation := &feature_annotation.NewFeatureAnnotation{
    Type: "gene",
    Attributes: &feature_annotation.FeatureAnnotationAttributes{
        Name: "genA",
        Synonyms: []string{"geneA", "GENA"},
        Publications: []string{"10.1234/journal.1234"},
        Pubmed: []string{"12345678"},
        Dblinks: []*feature_annotation.DbLink{
            {
                PrimaryId: "G123",
                Database: "GenBank",
                Linktype: "nucleotide",
                Url: "https://www.ncbi.nlm.nih.gov/gene/G123",
                Label: "GenBank: G123",
            },
        },
        Properties: []*feature_annotation.TagProperty{
            {
                Tag: "location",
                Value: "chromosome 2",
                CreatedBy: "user@example.com",
                CreatedAt: now,
                UpdatedAt: now,
            },
        },
    },
    CreatedBy: "user@example.com",
    CreatedAt: now,
    UpdatedAt: now,
}

response, err := client.CreateFeatureAnnotation(ctx, newAnnotation)
```

### Linking Feature to Organism

```go
orgClient := feature_annotation.NewOrganismFeatureServiceClient(conn)

link := &feature_annotation.OrganismFeatureLink{
    OrganismId: 12345,
    FeatureId: "feat-123",
}

_, err := orgClient.LinkFeatureToOrganism(ctx, link)
```

### Adding a Tag to a Feature

```go
tag := &feature_annotation.TagPropertyCreate{
    Tag: "function",
    Value: "transcription factor",
    CreatedBy: "user@example.com",
}

request := &feature_annotation.AddTagRequest{
    Id: "feat-123",
    Tag: tag,
}

updatedFeature, err := client.AddTag(ctx, request)
```

## Notes

- Email fields are validated to ensure they contain valid email addresses
- The DOI field validates that the format matches a DOI pattern
- Many operations maintain creation and modification timestamps
- Obsolete entities are marked with `is_obsolete` rather than being deleted