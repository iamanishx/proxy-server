# README.md Fix - Implementation Plan

- [x] **Task 1: Re-encode readme.md from UTF-16 LE to UTF-8**
  - The file is currently UTF-16 Little Endian (detected via `file` command) with a BOM (`fffe`).
  - Convert it to standard UTF-8 encoding using `iconv` or `recode`.

- [x] **Task 2: Fix garbled/truncated content at the end**
  - Remove the corrupted bytes at the end of the file (current content shows `=؀�` unrecoverable garbage).
  - The sentence should end cleanly: `> This mf can proxy upto 100k requests per sec`

- [x] **Task 3: Replace profanity / informal language with professional wording**
  - Change `This mf can proxy` to professional wording, e.g., `This project can proxy` or simply `A high-performance HTTP reverse proxy`.
  - Ensure the overall tone is appropriate for an open-source README.

- [x] **Task 4: Add a proper project title and description**
  - Keep the heading `# proxy-server`.
  - Expand the tagline/blockquote to describe what the project does professionally.
  - Add a concise description paragraph explaining this is an HTTP reverse proxy server in Go.

- [x] **Task 5: Add "Features" section**
  - Document key features: reverse proxying, configurable via `.env`, connection pooling (MaxIdleConns, MaxIdleConnsPerHost), high throughput.

- [x] **Task 6: Add "Prerequisites" section**
  - Document Go version requirement (1.23.3 as per `go.mod`).
  - Note that `make`, `golangci-lint` are needed for development tooling.

- [x] **Task 7: Add "Getting Started / Quick Start" section**
  - Show how to clone, configure, build, and run the server.
  - Include a minimal `.env` example.

- [x] **Task 8: Add "Configuration" section**
  - Document environment variables: `LISTEN_ADDR`, `BACKEND_URL`, `MAX_IDLE_CONNS`, `MAX_IDLE_CONNS_PER_HOST`.
  - Explain each variable's purpose and default values.

- [x] **Task 9: Add "Makefile Targets" section**
  - Document available make commands: `fmt`, `vet`, `tidy`, `lint`, `build`, `run`, `test`, `check`, `clean`.
  - Explain what each target does.

- [x] **Task 10: Add "Project Structure" section**
  - Briefly describe the code layout: `cmd/main.go` (entry point), `.env` (configuration), `makefile` (build automation), etc.

- [x] **Task 11: Final review and verification**
  - Verify the file is valid UTF-8 markdown with `file` command.
  - Check no profanity or garbled content remains.
  - Ensure all sections are properly formatted.
