# Post-Deploy GraphQL Liveness Workflow — Design

Date: 2026-11-05
Status: approved (approach and shape confirmed by the user)

## Problem

Deployments of graphql-server publish a GraphQL API but nothing verifies the
live operations after the deploy. The `graphql-client` CLI already ships a
`liveness-check` command that probes 27 GraphQL operations and writes a
Markdown report (endpoint health, per-operation health, failure summary).
No CI wiring exists for it.

## Goal

After every deployment of graphql-server (dev/staging and production), run the
`liveness-check` probe against the deployed endpoint and publish the report
where the team can see it. Probe failures must never block the pipeline.

## Decisions (confirmed)

| Decision | Choice |
|---|---|
| Trigger | After both deploys: `ci.yml` (dev/staging) and `tag-build-deploy.yml` (production) |
| Visualization | Actions job summary (Markdown report rendered on the run page) + uploaded artifact |
| Failure policy | Report-only. The workflow stays green even when probes fail |
| CLI source | Build from source: check out `dictyBase/graphql-client` at `develop` and run `go run . liveness-check` |

## Endpoints

| Environment | Endpoint |
|---|---|
| dev/staging (`ci.yml`) | `https://graphql.dictybase.dev/graphql` (verified live 2026-11-05) |
| production (`tag-build-deploy.yml`) | `https://graphql.dictycr.org/graphql` (provided by the user; `/graphql` path assumed to match dev, confirmed by first CI run) |

## Approaches considered

1. **Reusable workflow called from both deploy workflows (chosen).** One
   `workflow_call` workflow holds the probe logic; each deploy workflow adds one
   `needs: call-build-deploy` job. Deploy and probe share one run page, no extra
   permissions, no duplication of probe logic.
2. **Separate `workflow_run`-triggered workflow.** Loses adjacency on the run
   page, needs extra token permissions to read the triggering run, and makes
   the report harder to find.
3. **Manual `workflow_dispatch` only.** Not tied to deployments, so it misses
   exactly the moments it exists to verify.

## Shape

- `graphql-server/.github/workflows/liveness-probe.yml` — reusable workflow,
  one `endpoint` input. Checks out `dictyBase/graphql-client` at `develop`,
  sets up Go 1.26, runs `go run . liveness-check --endpoint <endpoint> --output
  liveness-report.md` with `continue-on-error: true`, appends the report to
  `$GITHUB_STEP_SUMMARY`, uploads the report as an artifact (90-day retention).
- `graphql-server/.github/workflows/ci.yml` — add job `liveness-probe` with
  `needs: call-build-deploy`, endpoint = dev/staging URL.
- `graphql-server/.github/workflows/tag-build-deploy.yml` — add job
  `liveness-probe` with `needs: call-build-deploy`, endpoint = production URL.

## Non-goals

- No change to the `graphql-client` CLI (it already exits 0 and writes the
  report; the report-only policy needs no exit-code work).
- No historical dashboard or GitHub Pages site (job summary + artifact chosen).
- No pipeline gating on probe failures.