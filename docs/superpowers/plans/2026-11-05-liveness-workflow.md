# Post-Deploy GraphQL Liveness Workflow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Run the `graphql-client` liveness probe after every graphql-server deployment and publish the Markdown report to the Actions job summary and as an artifact.

**Architecture:** One reusable `workflow_call` workflow holds the probe logic. Both deploy workflows (`ci.yml`, `tag-build-deploy.yml`) add one job each that calls it after `call-build-deploy`. The internal job carries `continue-on-error: true` (caller jobs cannot use that key), so step failures and probe failures do not fail the deploy run. Workflow-level validation errors and cancellation are not covered by this key.

**Tech Stack:** GitHub Actions (reusable workflows, job summary, upload-artifact), Go 1.26, `graphql-client` CLI `liveness-check` subcommand.

**Spec:** `docs/superpowers/plans/2026-11-05-liveness-workflow-design.md`

## Global Constraints

- Work happens in `~/Projects/devenv/golang/graphql-server`, branch `develop`.
- Report-only policy: no job or step failure caused by probe results may turn the workflow run red.
- Dev/staging endpoint: `https://graphql.dictybase.dev/graphql` (verified live 2026-11-05).
- Production endpoint: `https://graphql.dictycr.org/graphql` (user-provided; `/graphql` path assumed to match dev).
- CLI source: `dictyBase/graphql-client` at `develop`, run via `go run . liveness-check`.
- Action versions: `actions/checkout@v7`, `actions/setup-go@v7`, `actions/upload-artifact@v7` (latest majors, 2026-11-05).
- Artifact retention: 90 days.
- The sandbox cannot push. Commit locally and hand the push to the user.

## Review Focus

1. **Report-only guarantee.** A probe failure, a CLI crash, or a missing report file must not fail the run. The internal job carries `continue-on-error: true` (caller jobs cannot use that key; verified with actionlint). A stale `liveness-report.md` must never publish: delete it before the probe runs.
2. **Report file missing at summary time.** If the CLI dies before writing `liveness-report.md`, the summary step must not fail on `cat`. Guard with a file-existence check.
3. **Reusable workflow call syntax.** A job that calls a local reusable workflow uses `uses: ./.github/workflows/<file>` with `with:` inputs and cannot also have `runs-on` or `steps`.
4. **`needs` chaining.** The liveness job must start only after `call-build-deploy` succeeds, in both workflows.
5. **Production endpoint path.** `https://graphql.dictycr.org/graphql` is an assumption. The first tag-build run is the test; a wrong path shows as failed probes in the report, not a workflow failure.

## File Structure

- Create `.github/workflows/liveness-probe.yml` — reusable workflow: one job, four steps (checkout graphql-client, setup Go, run probe, publish report).
- Modify `.github/workflows/ci.yml` — add the `liveness-probe` caller job.
- Modify `.github/workflows/tag-build-deploy.yml` — add the `liveness-probe` caller job.

---

### Task 1: Reusable liveness probe workflow

**Files:**
- Create: `.github/workflows/liveness-probe.yml`

**Interfaces:**
- Produces: reusable workflow `liveness-probe.yml` with `workflow_call` input `endpoint` (string, required) — the GraphQL API URL to probe. Tasks 2 and 3 call it.

- [ ] **Step 1: Create `.github/workflows/liveness-probe.yml`**

```yaml
name: GraphQL Liveness Probe
on:
  workflow_call:
    inputs:
      endpoint:
        description: GraphQL API endpoint URL to probe
        required: true
        type: string
jobs:
  liveness-probe:
    name: Liveness probe
    runs-on: ubuntu-latest
    timeout-minutes: 10
    continue-on-error: true
    steps:
      - name: check out graphql-client
        uses: actions/checkout@v7
        with:
          repository: dictyBase/graphql-client
          ref: develop
      - name: set up golang
        uses: actions/setup-go@v7
        with:
          go-version: 1.26
      - name: run liveness check
        run: |
          rm -f liveness-report.md
          go run . liveness-check --endpoint "$GRAPHQL_ENDPOINT" --output liveness-report.md
        env:
          GOPROXY: https://proxy.golang.org
          GRAPHQL_ENDPOINT: ${{ inputs.endpoint }}
      - name: publish report to job summary
        if: always()
        run: |
          if [ -f liveness-report.md ]; then
            cat liveness-report.md >> "$GITHUB_STEP_SUMMARY"
          else
            echo "No liveness report was produced." >> "$GITHUB_STEP_SUMMARY"
          fi
      - name: upload report artifact
        if: always()
        uses: actions/upload-artifact@v7
        with:
          name: liveness-report
          path: liveness-report.md
          retention-days: 90
          if-no-files-found: ignore
```

- [ ] **Step 2: Lint the workflow**

Run: `cd ~/Projects/devenv/golang/graphql-server && go run github.com/rhysd/actionlint/cmd/actionlint@latest`
Expected: only the two pre-existing baseline errors (`ci.yml:25`, `testcov.yml:25`, the `ports[6379]` map access — present at the parent commit, out of scope). No error in `liveness-probe.yml`.

- [ ] **Step 3: Commit**

```bash
cd ~/Projects/devenv/golang/graphql-server
git add .github/workflows/liveness-probe.yml
git commit -m "ci: add reusable GraphQL liveness probe workflow"
```

---

### Task 2: Call the probe after the dev/staging deploy

**Files:**
- Modify: `.github/workflows/ci.yml`

**Interfaces:**
- Consumes: `liveness-probe.yml` from Task 1 (input `endpoint`).

- [ ] **Step 1: Add the `liveness-probe` job to `ci.yml`**

Append this job after `call-build-deploy`:

```yaml
  liveness-probe:
    name: GraphQL liveness probe
    needs: call-build-deploy
    uses: ./.github/workflows/liveness-probe.yml
    with:
      endpoint: https://graphql.dictybase.dev/graphql
```

- [ ] **Step 2: Lint the workflow**

Run: `cd ~/Projects/devenv/golang/graphql-server && go run github.com/rhysd/actionlint/cmd/actionlint@latest`
Expected: only the two pre-existing baseline errors (`ci.yml:25`, `testcov.yml:25`). No error in `ci.yml`'s `liveness-probe` job.

- [ ] **Step 3: Commit**

```bash
cd ~/Projects/devenv/golang/graphql-server
git add .github/workflows/ci.yml
git commit -m "ci: run liveness probe after dev/staging deploy"
```

---

### Task 3: Call the probe after the production deploy

**Files:**
- Modify: `.github/workflows/tag-build-deploy.yml`

**Interfaces:**
- Consumes: `liveness-probe.yml` from Task 1 (input `endpoint`).

- [ ] **Step 1: Add the `liveness-probe` job to `tag-build-deploy.yml`**

Append this job after `call-build-deploy`:

```yaml
  liveness-probe:
    name: GraphQL liveness probe
    needs: call-build-deploy
    uses: ./.github/workflows/liveness-probe.yml
    with:
      endpoint: https://graphql.dictycr.org/graphql
```

- [ ] **Step 2: Lint all workflows**

Run: `cd ~/Projects/devenv/golang/graphql-server && go run github.com/rhysd/actionlint/cmd/actionlint@latest`
Expected: only the two pre-existing baseline errors (`ci.yml:25`, `testcov.yml:25`). No error in `tag-build-deploy.yml`.

- [ ] **Step 3: Commit**

```bash
cd ~/Projects/devenv/golang/graphql-server
git add .github/workflows/tag-build-deploy.yml
git commit -m "ci: run liveness probe after production deploy"
```

---

### Task 4: Design and plan documents

**Files:**
- Create: `docs/superpowers/plans/2026-11-05-liveness-workflow-design.md` (already written)
- Create: `docs/superpowers/plans/2026-11-05-liveness-workflow.md` (this file)

- [ ] **Step 1: Commit the documents**

```bash
cd ~/Projects/devenv/golang/graphql-server
git add docs/superpowers/plans/
git commit -m "docs: add liveness workflow design and implementation plan"
```

- [ ] **Step 2: Hand the push to the user**

Report the branch (`develop`), the commit list, and remind the user to push. Do not push from the sandbox.