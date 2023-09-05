# TODO — rebuild GRC from scratch to current state

## Phase 01 — Project initialization

- [x] `go mod init github.com/carlosrabelo/grc`

## Phase 02 — Domain types

- [x] Create `grc/pkg/types/types.go` with structs `Filter`, `Feed`, `Entry`, `Category`, `Property`, `FiltersConfig`, `Author`, `Defaults`
- [x] Define XML namespace constants (`AtomNS`, `AppsNS`), email/domain regexes, `FeedID`, `XMLHeader`
- [x] `go fmt ./...` and `go vet ./...`

## Phase 03 — Test utilities

- [x] Create `grc/internal/testutils/testutils.go` with `CreateTempYAMLFile()`, `CreateTempFile()`, `CleanupFile()`, `BoolPtr()`

## Phase 04 — YAML validation

- [x] Create `grc/internal/validate/validate.go` with `.yaml`/`.yml` extension check, YAML parsing (`KnownFields(true)`), author validation, filter validation, email/domain validation
- [x] Create `grc/internal/validate/validate_test.go` with tests for all validation paths
- [x] `go get gopkg.in/yaml.v3` and `go mod tidy`

## Phase 05 — Config loading

- [x] Create `grc/internal/loader/loader.go` with `LoadConfig()` → validate extension → read file → parse YAML → validate config
- [x] Create `grc/internal/loader/loader_test.go` with tests for loading and error paths

## Phase 06 — Atom XML feed generation

- [ ] Create `grc/internal/generate/generate.go` with `GenerateFeed()`, `createBaseFeed()`, `createFeedEntry()`, `buildFilterProperties()`, `applyDefaults()`
- [ ] Create `grc/internal/generate/generate_test.go` with tests for feed structure, defaults, explicit booleans

## Phase 07 — XML persistence

- [ ] Create `grc/internal/save/save.go` with `SaveXML()`, `validateFileOverwrite()`, `writeXMLFile()`, `ensureXMLExtension()`
- [ ] Create `grc/internal/save/save_test.go` with tests for overwrite protection and force flag

## Phase 08 — Main CLI

- [ ] Create `grc/internal/app/app.go` with `Run()`, `parseCLIArgs()`, `validateRequiredArgs()`, `createLogger()`, `loadConfiguration()`, `generateXMLFeed()`, `resolveOutputPath()`, `persistXMLFile()`, `displayVersion()`, `displayHelp()`, `displaySuccessMessage()`
- [ ] Create `grc/cmd/grc/main.go` with entry point (≤50 lines), `version`/`buildTime` variables via ldflags
- [ ] Create `grc/internal/app/app_test.go` with CLI integration tests (basic execution, custom output, verbose, force, version, help, errors)

## Phase 09 — Example YAML

- [ ] Create example YAML config with sample filters (newsletter, support, internal notices, clients, financial, spam)

## Phase 10 — Documentation

- [ ] Create `docs/GUIDE.md` with complete usage guide
- [ ] Create `docs/GUIDE-PT.md` (Brazilian Portuguese)
- [ ] `LICENSE` MIT

## Phase 11 — Quality and CI

- [ ] Set up `golangci-lint`
- [ ] Push to GitHub and create GitHub Actions workflow for build + test on pushes and tags
