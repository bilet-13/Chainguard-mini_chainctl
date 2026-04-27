# mychainctl Implementation Plan

## High-Level Phases

1. Project setup and scaffolding
2. Registry client core logic
3. CLI commands and output formatting
4. Error handling and timeout behavior
5. Validation and lightweight testing

## Detailed Tasks

### Phase 1: Project setup and scaffolding

- Initialize Go module and basic project structure
- Create main entry and Cobra root command
- Add Makefile with build/run targets
- Wire global flags (e.g., --output)
- Set version/build metadata placeholders

### Phase 2: Registry client core logic

- Implement registry client constructor
- Implement ListTags for repository tag listing
- Implement InspectImage for image metadata
- Define ImageMetadata DTO
- Ensure client returns errors directly without retries

### Phase 3: CLI commands and output formatting

- Implement version command output
- Implement images list command arg parsing
- Implement images inspect command arg parsing
- Add JSON output formatting option
- Add table output with tabwriter

### Phase 4: Error handling and timeout behavior

- Add context.WithTimeout around client calls
- Ensure errors bubble to Cobra RunE
- Ensure clean error output in main

### Phase 5: Validation and lightweight testing

- Verify list command against cgr.dev
- Verify inspect command for a known image tag
- Validate JSON output structure
- Validate table formatting alignment
