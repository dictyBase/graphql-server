"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.mockUserByEmailQuery = exports.mockListPhenotypeAssaysQuery = exports.mockListPhenotypeEnvironmentsQuery = exports.mockListPhenotypesQuery = exports.mockListRecentStrainsQuery = exports.mockListRecentPlasmidsQuery = exports.mockStrainQuery = exports.mockPlasmidQuery = exports.mockPlasmidListFilterQuery = exports.mockListPlasmidsInventoryQuery = exports.mockListStrainsInventoryQuery = exports.mockListBacterialStrainsQuery = exports.mockListStrainsWithPhenotypeQuery = exports.mockStrainListQuery = exports.mockListPublicationsWithGeneQuery = exports.mockListRecentPublicationsQuery = exports.mockPublicationQuery = exports.mockListStrainsWithGeneQuery = exports.mockGeneOntologyAnnotationQuery = exports.mockListPublicationsWithGeneSummaryQuery = exports.mockGeneOntologyAnnotationSummaryQuery = exports.mockGeneGeneralInformationSummaryQuery = exports.mockListOrganismsQuery = exports.mockContentQuery = exports.mockContentBySlugQuery = exports.mockListContentByNamespaceQuery = exports.mockUpdateUserMutation = exports.mockCreateUserMutation = exports.mockUploadFileMutation = exports.mockUpdateStrainPhenotypeMutation = exports.mockAddStrainPhenotypeMutation = exports.mockCreateOrderMutation = exports.mockUpdateGeneGeneralInfoMutation = exports.mockCreateGeneGeneralInfoMutation = exports.mockDeleteContentMutation = exports.mockUpdateContentMutation = exports.mockCreateContentMutation = exports.mockLogoutMutation = exports.mockLoginMutation = exports.StrainType = exports.StatusEnum = exports.PlasmidType = void 0;
const msw_1 = require("msw");
var PlasmidType;
(function (PlasmidType) {
    PlasmidType["All"] = "ALL";
    PlasmidType["GoldenBraid"] = "GOLDEN_BRAID";
    PlasmidType["Regular"] = "REGULAR";
})(PlasmidType || (exports.PlasmidType = PlasmidType = {}));
var StatusEnum;
(function (StatusEnum) {
    StatusEnum["Cancelled"] = "CANCELLED";
    StatusEnum["Growing"] = "GROWING";
    StatusEnum["InPreparation"] = "IN_PREPARATION";
    StatusEnum["Shipped"] = "SHIPPED";
})(StatusEnum || (exports.StatusEnum = StatusEnum = {}));
var StrainType;
(function (StrainType) {
    StrainType["All"] = "ALL";
    StrainType["Bacterial"] = "BACTERIAL";
    StrainType["Gwdi"] = "GWDI";
    StrainType["Regular"] = "REGULAR";
})(StrainType || (exports.StrainType = StrainType = {}));
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockLoginMutation(
 *   ({ query, variables }) => {
 *     const { input } = variables;
 *     return HttpResponse.json({
 *       data: { login }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockLoginMutation = (resolver, options) => msw_1.graphql.mutation('Login', resolver, options);
exports.mockLoginMutation = mockLoginMutation;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockLogoutMutation(
 *   ({ query, variables }) => {
 *     return HttpResponse.json({
 *       data: { logout }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockLogoutMutation = (resolver, options) => msw_1.graphql.mutation('Logout', resolver, options);
exports.mockLogoutMutation = mockLogoutMutation;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockCreateContentMutation(
 *   ({ query, variables }) => {
 *     const { input } = variables;
 *     return HttpResponse.json({
 *       data: { createContent }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockCreateContentMutation = (resolver, options) => msw_1.graphql.mutation('CreateContent', resolver, options);
exports.mockCreateContentMutation = mockCreateContentMutation;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockUpdateContentMutation(
 *   ({ query, variables }) => {
 *     const { input } = variables;
 *     return HttpResponse.json({
 *       data: { updateContent }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockUpdateContentMutation = (resolver, options) => msw_1.graphql.mutation('UpdateContent', resolver, options);
exports.mockUpdateContentMutation = mockUpdateContentMutation;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockDeleteContentMutation(
 *   ({ query, variables }) => {
 *     const { id } = variables;
 *     return HttpResponse.json({
 *       data: { deleteContent }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockDeleteContentMutation = (resolver, options) => msw_1.graphql.mutation('DeleteContent', resolver, options);
exports.mockDeleteContentMutation = mockDeleteContentMutation;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockCreateGeneGeneralInfoMutation(
 *   ({ query, variables }) => {
 *     const { id, input } = variables;
 *     return HttpResponse.json({
 *       data: { createGeneGeneralInfo }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockCreateGeneGeneralInfoMutation = (resolver, options) => msw_1.graphql.mutation('CreateGeneGeneralInfo', resolver, options);
exports.mockCreateGeneGeneralInfoMutation = mockCreateGeneGeneralInfoMutation;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockUpdateGeneGeneralInfoMutation(
 *   ({ query, variables }) => {
 *     const { id, input } = variables;
 *     return HttpResponse.json({
 *       data: { updateGeneGeneralInfo }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockUpdateGeneGeneralInfoMutation = (resolver, options) => msw_1.graphql.mutation('UpdateGeneGeneralInfo', resolver, options);
exports.mockUpdateGeneGeneralInfoMutation = mockUpdateGeneGeneralInfoMutation;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockCreateOrderMutation(
 *   ({ query, variables }) => {
 *     const { input } = variables;
 *     return HttpResponse.json({
 *       data: { createOrder }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockCreateOrderMutation = (resolver, options) => msw_1.graphql.mutation('CreateOrder', resolver, options);
exports.mockCreateOrderMutation = mockCreateOrderMutation;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockAddStrainPhenotypeMutation(
 *   ({ query, variables }) => {
 *     const { strainId, input } = variables;
 *     return HttpResponse.json({
 *       data: { addStrainPhenotype }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockAddStrainPhenotypeMutation = (resolver, options) => msw_1.graphql.mutation('AddStrainPhenotype', resolver, options);
exports.mockAddStrainPhenotypeMutation = mockAddStrainPhenotypeMutation;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockUpdateStrainPhenotypeMutation(
 *   ({ query, variables }) => {
 *     const { strainId, target, payload } = variables;
 *     return HttpResponse.json({
 *       data: { updateStrainPhenotype }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockUpdateStrainPhenotypeMutation = (resolver, options) => msw_1.graphql.mutation('UpdateStrainPhenotype', resolver, options);
exports.mockUpdateStrainPhenotypeMutation = mockUpdateStrainPhenotypeMutation;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockUploadFileMutation(
 *   ({ query, variables }) => {
 *     const { file } = variables;
 *     return HttpResponse.json({
 *       data: { uploadFile }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockUploadFileMutation = (resolver, options) => msw_1.graphql.mutation('UploadFile', resolver, options);
exports.mockUploadFileMutation = mockUploadFileMutation;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockCreateUserMutation(
 *   ({ query, variables }) => {
 *     const { input } = variables;
 *     return HttpResponse.json({
 *       data: { createUser }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockCreateUserMutation = (resolver, options) => msw_1.graphql.mutation('CreateUser', resolver, options);
exports.mockCreateUserMutation = mockCreateUserMutation;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockUpdateUserMutation(
 *   ({ query, variables }) => {
 *     const { id, input } = variables;
 *     return HttpResponse.json({
 *       data: { updateUser }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockUpdateUserMutation = (resolver, options) => msw_1.graphql.mutation('UpdateUser', resolver, options);
exports.mockUpdateUserMutation = mockUpdateUserMutation;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListContentByNamespaceQuery(
 *   ({ query, variables }) => {
 *     const { namespace, limit } = variables;
 *     return HttpResponse.json({
 *       data: { listContentByNamespace }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockListContentByNamespaceQuery = (resolver, options) => msw_1.graphql.query('ListContentByNamespace', resolver, options);
exports.mockListContentByNamespaceQuery = mockListContentByNamespaceQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockContentBySlugQuery(
 *   ({ query, variables }) => {
 *     const { slug } = variables;
 *     return HttpResponse.json({
 *       data: { contentBySlug }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockContentBySlugQuery = (resolver, options) => msw_1.graphql.query('ContentBySlug', resolver, options);
exports.mockContentBySlugQuery = mockContentBySlugQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockContentQuery(
 *   ({ query, variables }) => {
 *     const { id } = variables;
 *     return HttpResponse.json({
 *       data: { content }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockContentQuery = (resolver, options) => msw_1.graphql.query('Content', resolver, options);
exports.mockContentQuery = mockContentQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListOrganismsQuery(
 *   ({ query, variables }) => {
 *     return HttpResponse.json({
 *       data: { listOrganisms }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockListOrganismsQuery = (resolver, options) => msw_1.graphql.query('ListOrganisms', resolver, options);
exports.mockListOrganismsQuery = mockListOrganismsQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockGeneGeneralInformationSummaryQuery(
 *   ({ query, variables }) => {
 *     const { gene } = variables;
 *     return HttpResponse.json({
 *       data: { geneGeneralInformation }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockGeneGeneralInformationSummaryQuery = (resolver, options) => msw_1.graphql.query('GeneGeneralInformationSummary', resolver, options);
exports.mockGeneGeneralInformationSummaryQuery = mockGeneGeneralInformationSummaryQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockGeneOntologyAnnotationSummaryQuery(
 *   ({ query, variables }) => {
 *     const { gene } = variables;
 *     return HttpResponse.json({
 *       data: { geneOntologyAnnotation }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockGeneOntologyAnnotationSummaryQuery = (resolver, options) => msw_1.graphql.query('GeneOntologyAnnotationSummary', resolver, options);
exports.mockGeneOntologyAnnotationSummaryQuery = mockGeneOntologyAnnotationSummaryQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListPublicationsWithGeneSummaryQuery(
 *   ({ query, variables }) => {
 *     const { gene } = variables;
 *     return HttpResponse.json({
 *       data: { listPublicationsWithGene }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockListPublicationsWithGeneSummaryQuery = (resolver, options) => msw_1.graphql.query('ListPublicationsWithGeneSummary', resolver, options);
exports.mockListPublicationsWithGeneSummaryQuery = mockListPublicationsWithGeneSummaryQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockGeneOntologyAnnotationQuery(
 *   ({ query, variables }) => {
 *     const { gene } = variables;
 *     return HttpResponse.json({
 *       data: { geneOntologyAnnotation }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockGeneOntologyAnnotationQuery = (resolver, options) => msw_1.graphql.query('GeneOntologyAnnotation', resolver, options);
exports.mockGeneOntologyAnnotationQuery = mockGeneOntologyAnnotationQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListStrainsWithGeneQuery(
 *   ({ query, variables }) => {
 *     const { gene } = variables;
 *     return HttpResponse.json({
 *       data: { listStrainsWithGene }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockListStrainsWithGeneQuery = (resolver, options) => msw_1.graphql.query('ListStrainsWithGene', resolver, options);
exports.mockListStrainsWithGeneQuery = mockListStrainsWithGeneQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockPublicationQuery(
 *   ({ query, variables }) => {
 *     const { id } = variables;
 *     return HttpResponse.json({
 *       data: { publication }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockPublicationQuery = (resolver, options) => msw_1.graphql.query('Publication', resolver, options);
exports.mockPublicationQuery = mockPublicationQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListRecentPublicationsQuery(
 *   ({ query, variables }) => {
 *     const { limit } = variables;
 *     return HttpResponse.json({
 *       data: { listRecentPublications }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockListRecentPublicationsQuery = (resolver, options) => msw_1.graphql.query('ListRecentPublications', resolver, options);
exports.mockListRecentPublicationsQuery = mockListRecentPublicationsQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListPublicationsWithGeneQuery(
 *   ({ query, variables }) => {
 *     const { gene } = variables;
 *     return HttpResponse.json({
 *       data: { listPublicationsWithGene }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockListPublicationsWithGeneQuery = (resolver, options) => msw_1.graphql.query('ListPublicationsWithGene', resolver, options);
exports.mockListPublicationsWithGeneQuery = mockListPublicationsWithGeneQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockStrainListQuery(
 *   ({ query, variables }) => {
 *     const { cursor, limit, filter } = variables;
 *     return HttpResponse.json({
 *       data: { listStrains }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockStrainListQuery = (resolver, options) => msw_1.graphql.query('StrainList', resolver, options);
exports.mockStrainListQuery = mockStrainListQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListStrainsWithPhenotypeQuery(
 *   ({ query, variables }) => {
 *     const { cursor, limit, type, annotation } = variables;
 *     return HttpResponse.json({
 *       data: { listStrainsWithAnnotation }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockListStrainsWithPhenotypeQuery = (resolver, options) => msw_1.graphql.query('ListStrainsWithPhenotype', resolver, options);
exports.mockListStrainsWithPhenotypeQuery = mockListStrainsWithPhenotypeQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListBacterialStrainsQuery(
 *   ({ query, variables }) => {
 *     return HttpResponse.json({
 *       data: { listStrainsWithAnnotation, listStrainsWithAnnotation }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockListBacterialStrainsQuery = (resolver, options) => msw_1.graphql.query('ListBacterialStrains', resolver, options);
exports.mockListBacterialStrainsQuery = mockListBacterialStrainsQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListStrainsInventoryQuery(
 *   ({ query, variables }) => {
 *     const { cursor, limit } = variables;
 *     return HttpResponse.json({
 *       data: { listStrainsWithAnnotation }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockListStrainsInventoryQuery = (resolver, options) => msw_1.graphql.query('ListStrainsInventory', resolver, options);
exports.mockListStrainsInventoryQuery = mockListStrainsInventoryQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListPlasmidsInventoryQuery(
 *   ({ query, variables }) => {
 *     const { cursor, limit } = variables;
 *     return HttpResponse.json({
 *       data: { listPlasmidsWithAnnotation }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockListPlasmidsInventoryQuery = (resolver, options) => msw_1.graphql.query('ListPlasmidsInventory', resolver, options);
exports.mockListPlasmidsInventoryQuery = mockListPlasmidsInventoryQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockPlasmidListFilterQuery(
 *   ({ query, variables }) => {
 *     const { cursor, limit, filter } = variables;
 *     return HttpResponse.json({
 *       data: { listPlasmids }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockPlasmidListFilterQuery = (resolver, options) => msw_1.graphql.query('PlasmidListFilter', resolver, options);
exports.mockPlasmidListFilterQuery = mockPlasmidListFilterQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockPlasmidQuery(
 *   ({ query, variables }) => {
 *     const { id } = variables;
 *     return HttpResponse.json({
 *       data: { plasmid }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockPlasmidQuery = (resolver, options) => msw_1.graphql.query('Plasmid', resolver, options);
exports.mockPlasmidQuery = mockPlasmidQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockStrainQuery(
 *   ({ query, variables }) => {
 *     const { id } = variables;
 *     return HttpResponse.json({
 *       data: { strain }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockStrainQuery = (resolver, options) => msw_1.graphql.query('Strain', resolver, options);
exports.mockStrainQuery = mockStrainQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListRecentPlasmidsQuery(
 *   ({ query, variables }) => {
 *     const { limit } = variables;
 *     return HttpResponse.json({
 *       data: { listRecentPlasmids }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockListRecentPlasmidsQuery = (resolver, options) => msw_1.graphql.query('ListRecentPlasmids', resolver, options);
exports.mockListRecentPlasmidsQuery = mockListRecentPlasmidsQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListRecentStrainsQuery(
 *   ({ query, variables }) => {
 *     const { limit } = variables;
 *     return HttpResponse.json({
 *       data: { listRecentStrains }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockListRecentStrainsQuery = (resolver, options) => msw_1.graphql.query('ListRecentStrains', resolver, options);
exports.mockListRecentStrainsQuery = mockListRecentStrainsQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListPhenotypesQuery(
 *   ({ query, variables }) => {
 *     const { search } = variables;
 *     return HttpResponse.json({
 *       data: { listPhenotypes }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockListPhenotypesQuery = (resolver, options) => msw_1.graphql.query('ListPhenotypes', resolver, options);
exports.mockListPhenotypesQuery = mockListPhenotypesQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListPhenotypeEnvironmentsQuery(
 *   ({ query, variables }) => {
 *     const { search } = variables;
 *     return HttpResponse.json({
 *       data: { listPhenotypeEnvironments }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockListPhenotypeEnvironmentsQuery = (resolver, options) => msw_1.graphql.query('ListPhenotypeEnvironments', resolver, options);
exports.mockListPhenotypeEnvironmentsQuery = mockListPhenotypeEnvironmentsQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockListPhenotypeAssaysQuery(
 *   ({ query, variables }) => {
 *     const { search } = variables;
 *     return HttpResponse.json({
 *       data: { listPhenotypeAssays }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockListPhenotypeAssaysQuery = (resolver, options) => msw_1.graphql.query('ListPhenotypeAssays', resolver, options);
exports.mockListPhenotypeAssaysQuery = mockListPhenotypeAssaysQuery;
/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockUserByEmailQuery(
 *   ({ query, variables }) => {
 *     const { email } = variables;
 *     return HttpResponse.json({
 *       data: { userByEmail }
 *     })
 *   },
 *   requestOptions
 * )
 */
const mockUserByEmailQuery = (resolver, options) => msw_1.graphql.query('UserByEmail', resolver, options);
exports.mockUserByEmailQuery = mockUserByEmailQuery;
