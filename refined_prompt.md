Objective: Refactor the `liststrain` GraphQL resolver and filter and the data import pipeline to support bacterial strain filtering, based on new strain annotations and the provided data file.

Context:
1. You are working with three main components:
   - `modware-import` (located in `/home/agent/modware-import`): Responsible for uploading the strain data file.
   - `modware-stock` (located in `/home/agent/modware-stock`): Responsible for managing strain data and executing searches/queries.
   - `@api` folder (located at `/home/agent/workspace/api`): Contains the GraphQL schema definition.
2. The data file to be imported is `strain_characteristics.tsv`.
3. Bacterial strains are now annotated, as discussed in [Migration Issue #83](https://github.com/dictyBase/Migration/issues/83).
4. The `strain_characteristics.tsv` file contains rows detailing strain characteristics (e.g., `DBS0236830	blasticidin resistant`, `DBS0236830	axenic`, etc.), which need to be parsed and stored.

Requirements:
1. Analyze the GraphQL Schema: Locate the `liststrain` query and any related filter input types in the `@api` folder. Extend the schema if necessary to accommodate bacterial strain properties (e.g., specific annotations found in `strain_characteristics.tsv`).
2. Update the Importer: Modify the `modware-import` project to correctly parse and upload the bacterial strain data from `strain_characteristics.tsv`, ensuring that all characteristics (like "drug resistant", "axenic", "null mutant", "GFP marked", etc.) are correctly associated with the respective strain IDs.
3. Refactor the `liststrain` Filter in `modware-stock`: 
   - Locate the `liststrain` resolver/filter logic in the `modware-stock` project.
   - Refactor the code to support filtering based on the new bacterial annotations.
   - Ensure the database query generation is optimized for the new filter arguments.
4. Verify tests pass for the updated `liststrain` resolver. Ensure the returned data matches the expected format defined in the schema.