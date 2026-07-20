# STATUS: DONE

## RUN_ID
20260720-123954

## COMMIT
2e55b7f

## JOB_ID
data-tooling

## SKILL
tdd

## WORKER
data-tooling (agent)

## SUMMARY
Created versioned JSON catalog contract (version 1) shared by development seed and production import. Implemented idempotent development seed-on-empty using the embedded contract. Added `cmd/migrate` and `cmd/importcatalog` commands with process-level tests.

The JSON contract (`CatalogFile`/`CatalogProduct`/`GiaJSON`) validates version, required fields, known enum values, RFC3339 timestamps, and price consistency. The embedded `seed_data.json` represents all 16 products from the existing `SeedData()` (12 approved, 4 non-approved).

`SeedSQLite` now checks for existing products before seeding — seeds only when empty, skips when populated. `ImportCatalogJSON` validates the full input before beginning the write transaction, rejects duplicate IDs by default (including duplicates inside the input and conflicts with existing rows), and imports all-or-nothing in one transaction.

`cmd/migrate` runs goose migrations explicitly using `OpenSQLite`. `cmd/importcatalog` uses `OpenSQLiteProd` (new shared function that applies PRAGMAs and verifies schema without auto-migration). `main.go` updated to use `catalog.OpenSQLiteProd` instead of the private `openProdDB`.

## ARTIFACTS

| File | Description |
|------|-------------|
| `src/api/catalog/catalog_json.go` | Versioned JSON contract: `CatalogFile`, `CatalogProduct`, `GiaJSON` types; `ValidateCatalogJSON()`, `ToSanPhamSo()`, `SanPhamSoToCatalogProduct()` |
| `src/api/catalog/catalog_json_test.go` | 14 test cases covering valid validation, all rejection paths, round-trip conversion, and invalid timestamps |
| `src/api/catalog/seed_data.json` | Embedded versioned JSON with 16 products matching `SeedData()` |
| `src/api/catalog/import.go` | `ImportCatalogJSON()`: validate first, all-or-nothing transaction, duplicate rejection |
| `src/api/catalog/import_test.go` | 5 test cases: success, reject input duplicates, reject DB conflicts, allow-duplicates mode, invalid JSON rejection |
| `src/api/catalog/seed_test.go` | 5 test cases: JSON parsing correctness, product data fidelity, empty-only idempotence, queryability, skip without transaction |
| `src/api/cmd/migrate/main.go` + `main_test.go` | Explicit migrate command; test for success and missing APP_ENV failure |
| `src/api/cmd/importcatalog/main.go` + `main_test.go` | Production import command with `-path` and `-allow-duplicates` flags; tests for success, duplicate rejection, and invalid version |

## CHANGED_FILES

| File | Status | Change |
|------|--------|--------|
| `src/api/catalog/catalog_json.go` | new | JSON contract types, validation, and conversion |
| `src/api/catalog/catalog_json_test.go` | new | Validation test suite |
| `src/api/catalog/seed_data.json` | new | Embedded versioned seed data |
| `src/api/catalog/import.go` | new | `ImportCatalogJSON` all-or-nothing import |
| `src/api/catalog/import_test.go` | new | Import test suite |
| `src/api/catalog/seed_test.go` | new | Seed behavior test suite |
| `src/api/cmd/migrate/main.go` | new | Explicit migration command |
| `src/api/cmd/migrate/main_test.go` | new | Process-level migrate tests |
| `src/api/cmd/importcatalog/main.go` | new | Production import command |
| `src/api/cmd/importcatalog/main_test.go` | new | Process-level import tests |
| `src/api/catalog/seed.go` | modified | Embedded JSON, empty-only `SeedSQLite`, `SeedFromJSON` |
| `src/api/catalog/sqlite.go` | modified | Added `OpenSQLiteProd()` production DB opener |
| `src/api/main.go` | modified | Replaced private `openProdDB` with `catalog.OpenSQLiteProd()`; removed unused `fmt` |

Unchanged: `src/web/package-lock.json` (never touched), all existing handler/sqlite tests, API routes/response shapes.

## VERIFY

All 50+ tests pass across 4 packages:

```
ok  ecommerce/api                   0.969s
ok  ecommerce/api/catalog           0.127s
ok  ecommerce/api/cmd/importcatalog  1.332s
ok  ecommerce/api/cmd/migrate        0.452s
```

| Test | Result |
|------|--------|
| TestValidateCatalogJSON (10 subtests) | PASS |
| TestSanPhamSoToCatalogProduct_RoundTrip | PASS |
| TestToSanPhamSo_RoundTrip | PASS |
| TestToSanPhamSo_RejectsInvalidTimestamp | PASS |
| TestSeedFromJSON (2 subtests) | PASS |
| TestSeedSQLite (3 subtests) | PASS |
| TestImportCatalogJSON (5 subtests) | PASS |
| TestMigrate_Success | PASS |
| TestMigrate_FailsWithoutAppEnv | PASS |
| TestImportCatalog_Success | PASS |
| TestImportCatalog_RejectsDuplicateIDs | PASS |
| TestImportCatalog_RejectsInvalidVersion | PASS |
| All existing catalog tests (50+) | PASS |

Build & lint:
- `go vet ./...` — clean
- `go build ./...` — clean
- `go fmt ./...` — clean
- `git diff --check` — clean

## NEXT_SKILL_HINT
code-review

## NEXT_INPUTS
- Issue #20: Migrate catalog runtime persistence from in-memory to SQLite
- All artifact files listed above
- Base commit: 07f95e9

## BLOCKERS
None.

## RISKS
- `pressly/goose/v3 v3.20.0` is the latest goose release compatible with Go 1.22.5. Future goose releases may require Go ≥ 1.25; the go.mod `go` directive must not be upgraded.
- Embedded `seed_data.json` must match `SeedData()` in data and count. Regenerate by running `SeedData()` through `SanPhamSoToCatalogProduct()` if products are added or removed.
- `ImportCatalogJSON` double-checks existing DB rows before the write transaction. On large databases with many products, this could be slow; production datasets are expected to be small at MVP scale.

## NOTES_FOR_ORCHESTRATOR
- Dev seed (`SeedSQLite`) now uses the embedded JSON and seeds only when empty. Existing tests using `SeedData()` directly (handler tests, sqlite_test.go helpers) are unaffected.
- `OpenSQLiteProd` is the production-aware DB opener, shared by `main.go` and `cmd/importcatalog`.
- The goose version pin (v3.20.0) and Go pin (1.22.5) are preserved.
- `src/web/package-lock.json` was never modified (verified via `git status`).
