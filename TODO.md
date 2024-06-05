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

- [x] Create `grc/internal/generate/generate.go` with `GenerateFeed()`, `createBaseFeed()`, `createFeedEntry()`, `buildFilterProperties()`, `applyDefaults()`
- [x] Create `grc/internal/generate/generate_test.go` with tests for feed structure, defaults, explicit booleans

## Phase 07 — XML persistence

- [x] Create `grc/internal/save/save.go` with `SaveXML()`, `validateFileOverwrite()`, `writeXMLFile()`, `ensureXMLExtension()`
- [x] Create `grc/internal/save/save_test.go` with tests for overwrite protection and force flag

## Phase 08 — Main CLI

- [x] Create `grc/internal/app/app.go` with `Run()`, `parseCLIArgs()`, `validateRequiredArgs()`, `createLogger()`, `loadConfiguration()`, `generateXMLFeed()`, `resolveOutputPath()`, `persistXMLFile()`, `displayVersion()`, `displayHelp()`, `displaySuccessMessage()`
- [x] Create `grc/cmd/grc/main.go` with entry point (≤50 lines), `version`/`buildTime` variables via ldflags
- [x] Create `grc/internal/app/app_test.go` with CLI integration tests (basic execution, custom output, verbose, force, version, help, errors)

## Phase 09 — Example YAML

- [x] Create example YAML config with sample filters (newsletter, support, internal notices, clients, financial, spam)

## Phase 10 — Documentation

- [ ] Create `docs/GUIDE.md` with complete usage guide
- [ ] Create `docs/GUIDE-PT.md` (Brazilian Portuguese)
- [x] `LICENSE` GPL-3

## Phase 11 — Quality

- [ ] `make fmt`, `make vet`, and `make quality` (fmt + vet + test)

## Phase 12 — Gmail filter types

- [ ] Extend `Filter` with `Criteria *GmailCriteria` and `Action *GmailAction`
- [ ] Define structs `GmailCriteria` (From, To, Subject, HasTheWord, etc.) and `GmailAction` (Archive, MarkRead, Star, Delete, NeverSpam, AlwaysSpam, Label, Forward, SmartLabel, Category)
- [ ] Keep backward compatibility: existing `expression` field continues to work

## Phase 13 — Expression parser

- [ ] Create `grc/internal/parser/parser.go` that parses Gmail search expressions into structured `Criteria` + `Action`
- [ ] Support operators: `AND`, `OR`, `-` (negation), `from:`, `to:`, `subject:`, `has:`, `filename:`, `list:`, `after:`, `before:`, `older_than:`, `newer_than:`, `size:`, `larger:`, `smaller:`, `in:`, `is:`, `label:`, `cc:`, `bcc:`, `deliveredto:`
- [ ] Handle quoted values and parenthesized groups
- [ ] Create `grc/internal/parser/parser_test.go` with table-driven tests for all operators

## Phase 14 — Gmail XML generation

- [ ] Create `grc/internal/gmail/gmail.go` with `GenerateGmailXML()` that produces Gmail import XML (`apps:property` entries)
- [ ] Map `GmailCriteria` fields to `apps:property` name/value pairs (from, to, subject, hasTheWord, hasAttachment, etc.)
- [ ] Map `GmailAction` fields to `apps:property` actions (shouldArchive, shouldMarkAsRead, shouldStar, shouldDelete, etc.)
- [ ] Create `grc/internal/gmail/gmail_test.go` with tests for all criteria/action combinations

## Phase 15 — CLI `--format` flag

- [ ] Add `--format atom|gmail` flag to CLI (default: `atom` for backward compatibility)
- [ ] Wire `--format gmail` to `gmail.GenerateGmailXML()` pipeline
- [ ] Update `resolveOutputPath()` to use `.gmail.xml` extension for gmail format
- [ ] Update help text and tests
- [ ] Create `docs/example-gmail.yaml` with Gmail-specific criteria/action filters

## Phase 16 — Gmail filter import (reverse)

- [ ] Create `grc/internal/importer/importer.go` that reads Gmail export XML (Google Takeout) and produces `[]Filter`
- [ ] Map `apps:property` criteria entries back to `GmailCriteria`
- [ ] Map `apps:property` action entries back to `GmailAction`
- [ ] Add `--import` flag or `grc import <gmail.xml>` subcommand
- [ ] Create `grc/internal/importer/importer_test.go` with roundtrip tests

## Phase 17 — Search query builder

- [ ] Create `grc/internal/query/query.go` with `BuildSearchURL()` that generates `https://mail.google.com/mail/u/0/#search/...` links
- [ ] Reconstruct Gmail search string from `GmailCriteria`
- [ ] Create `grc/internal/query/query_test.go`

## Phase 18 — YAML config with Gmail examples

- [ ] Create `docs/example-gmail.yaml` with filters using `criteria` + `action` (newsletter, support, invoices, etc.)
- [ ] Demonstrate archive, label, star, markRead, forward, delete actions

## Phase 19 — Gmail integration tests

- [ ] End-to-end tests: YAML → parser → gmail XML → validate XML schema
- [ ] Roundtrip tests: YAML → gmail XML → importer → YAML
- [ ] Test `grc --format gmail` CLI integration
