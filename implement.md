# mychainctl Implementation Plan

## High-Level Phases

1. Project setup and scaffolding
2. Registry client core logic
3. CLI commands and output formatting
4. Error handling and timeout behavior
5. Validation and lightweight testing

## Detailed Tasks

### Phase 1: Project setup and scaffolding

- Initialize Go module and basic project structure (done)
- Create main entry and Cobra root command (done)
- Add Makefile with build/run targets (done)
- Wire global flags (e.g., --output) (done)
- Set version/build metadata placeholders (done)

Implementation notes:

- Global flag `--output/-o` is wired in the root command and stored in `outputFormat`.
- Version/build placeholders live in `cmd/version.go` as `version`, `commit`, and `date` for future `-ldflags` injection.

### Phase 2: Registry client core logic

- Implement registry client constructor (done)
- Implement ListTags for repository tag listing (done)
- Implement InspectImage for image metadata (done)
- Define ImageMetadata DTO (done)
- Ensure client returns errors directly without retries (done)

Implementation notes:

- The registry client lives in `pkg/registry/client.go` with `NewClient`, `ListTags`, and `InspectImage`.
- `ListTags` builds a full repository name using the configured registry and calls `remote.List`.
- `InspectImage` uses `remote.Get`, returning digest/media type and a best-effort platform string.

### Phase 3: CLI commands and output formatting

- Implement version command output (done)
- Implement images list command arg parsing (done)
- Implement images inspect command arg parsing (done)
- Add JSON output formatting option (done)
- Add table output with tabwriter (done)

Implementation notes:

- `version` prints version, commit, and date via a Cobra command in `cmd/version.go`.
- `images list` and `images inspect` live in `cmd/images.go` and honor the global `--output` flag.
- JSON output uses a pretty encoder; table output uses `text/tabwriter` for alignment.

### Phase 4: Error handling and timeout behavior

- Add context.WithTimeout around client calls
- Ensure errors bubble to Cobra RunE
- Ensure clean error output in main

### Phase 5: Validation and lightweight testing

- Verify list command against cgr.dev
- Verify inspect command for a known image tag
- Validate JSON output structure
- Validate table formatting alignment
