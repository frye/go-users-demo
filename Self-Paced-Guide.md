# Go Users Demo — Workflow Guide

This guide shows how to accomplish the requested tasks using GitHub Copilot agents in VS Code: Copilot Plan Mode (for structured planning), Copilot Agent Mode (for executing changes), Mission Control (to observe agent runs), and Copilot Chat (to review changes). You'll work in a fork of `go-users-demo`, plan unit tests and delete functionality, create a tracking issue, implement tests via Agent Mode, observe execution, and review changes.

## How Copilot Agents Help

• Copilot Plan Mode: Produces clear, sequenced plans for features and tests before implementation.
• Copilot Agent Mode: Executes changes in your fork (editing files, running builds/tests) following your plans.
• Mission Control: Observes the Agent's actions, tool calls, diffs, and test runs in real time.
• Copilot Chat: Assists with code reviews, summarizes diffs, explains failures, and suggests fixes.

## Quick Start

• Fork and clone: Create a fork on GitHub, then `git clone` your fork; set upstream to original.
• Branch: `git checkout -b feature/tests-and-delete`.
• Build check: `go version`, then `go run main.go` to verify Go and dependencies.

## Repository Summary and Gaps

• Structure:
  ◦ `main.go`: Go entrypoint with Gin setup.
  ◦ `templates/users.html`: loads user data from API.
  ◦ `models/user.go`: in-memory store; `User`, `UserWriteModel`, user not found handling.
  ◦ `controllers/user_controller.go`: REST endpoints.
• Implemented endpoints:
  ◦ `GET /users` → returns `[]User`.
  ◦ `GET /users/:id` → returns `User` or 404 (via error handling).
  ◦ `POST /users` → creates and returns `User` (201 Created).
  ◦ `PUT /users/:id` → updates and returns `User` or 404.
• Constraints:
  ◦ Go (1.21+), Gin; minimal deps — only `gin-gonic/gin` (tests would use standard library `testing` and `net/http/httptest`).
  ◦ In-memory store (slice or map) seeded with demo users.
  ◦ `UserWriteModel` provides defaults when absent; store handles nil/blank as defaults.
• Gaps:
  ◦ No delete method in store or `DELETE` route in controller.
  ◦ No tests or CI.

## Tasks and Agent Usage (with Sample Prompts)

Use these agents and prompts to drive each task efficiently. Run prompts in Copilot Plan Mode for planning and Copilot Agent Mode for execution. Observe runs with Mission Control and review with Copilot Chat.

• Task 1: Plan Unit Test Generation (Plan Mode)

  ◦ Purpose: Produce the unit test strategy before coding.
  ◦ Sample prompt (Plan Mode):
    ■ "Plan unit tests for user store and controller in go-users-demo per project guidance. Include files, coverage points, status codes, and minimal dependencies."

• Task 2: Plan Delete Functionality (Plan Mode)

  ◦ Purpose: Define store and API changes for delete.
  ◦ Sample prompt (Plan Mode):
    ■ "Plan DELETE /users/:id implementation: store delete(id), 204 success, 404 missing, README updates, and test plan."

• Task 3: Create GitHub Issue (Agent Mode + GitHub MCP)

  ◦ Purpose: Track delete implementation in your fork.
  ◦ Sample prompt (Agent Mode, GitHub MCP):
    ■ "Create an issue in my fork your-username/go-users-demo titled 'Add DELETE /users/:id endpoint and store delete()' with requirements, files, and acceptance criteria as specified. Label enhancement, backend; assign to me."

• Task 4: Implement Unit Test Plan (Agent Mode, read upstream issue with MCP)

  ◦ Purpose: Execute the unit test plan using the tracking issue as the exact prompt via GitHub MCP.
  ◦ How to read and use the issue with MCP:
    ■ Ensure GitHub MCP is connected with a PAT (`repo` scope) in VS Code.
    ■ Ask the Agent: "Fetch the latest open issue titled 'Add DELETE /users/:id endpoint and store delete()' from my fork and display its body."
    ■ Copy the full issue body and use it verbatim as the Agent prompt for implementation.
  ◦ Sample prompt (Agent Mode):
    ■ "Using the upstream tracking issue content verbatim, implement the unit test plan: create user_store_test.go and user_controller_test.go with the specified coverage; run go test and fix failures aligned with UserWriteModel defaults. Commit and push a PR referencing the issue."

• Task 5: Observe Agent Execution (Mission Control)

  ◦ Purpose: Monitor the agent's actions and test runs.
  ◦ Sample prompt (Mission Control context):
    ■ "Start observing the unit test implementation run. Mark checkpoints after each test file creation and after successful go test. Surface failures and suggest corrections."

• Task 6: Code Review with Copilot (Copilot Chat)

  ◦ Purpose: Review changes, summarize diffs, and suggest fixes.
  ◦ Sample prompts (Copilot Chat):
    ■ "Summarize changes in user_controller.go, routes.go, and tests. Identify potential issues and suggest minimal fixes."
    ■ "Explain failing tests and propose corrections consistent with project guidance (status codes, defaults, minimal dependencies)."
    ■ "Generate a review checklist covering store/controller paths (create/update/find/delete), 404 mapping, 201/200/204 status codes, and README updates."

## Detailed Steps

### Task 1: Copilot Plan Mode — Unit Test Generation

• Add test files:
  ◦ `models/user_store_test.go`.
  ◦ `controllers/user_controller_test.go`.
• `user_store_test.go` coverage:
  ◦ `findAll()` returns seeded users in insertion order.
  ◦ `findByID(id)` returns present for existing, empty for missing.
  ◦ `create(UserWriteModel)` defaults: generates id, name `"Anonymous"`, emoji wave `"👋"` when absent/blank.
  ◦ `update(id, UserWriteModel)` updates provided fields, retains others; throws not found error if id missing.
  ◦ Seed users exist: "Bramble Fright" 👻, "Sylvie Scream" 🎃, "Eve Eerie" 🧙.
• `user_controller_test.go` with `net/http/httptest`, Gin test context:
  ◦ `GET /users` → 200, JSON array.
  ◦ `GET /users/:id` → 200 for existing, 404 for missing.
  ◦ `POST /users` → 201, JSON payload.
  ◦ `PUT /users/:id` → 200 for existing, 404 for missing.
• Align tests with `UserWriteModel` semantics (nullable fields handled by store). Keep dependencies minimal.

### Task 2: Copilot Plan Mode — Delete Functionality

• Store: add `delete(id)` to remove existing or return not found error.
• Controller: add `DELETE /users/:id` returning `204 No Content` on success; 404 on missing.
• Documentation: update `README.md` for DELETE endpoint behavior and status codes.
• Optional UI: add delete button and `fetch('/users/:id', { method: 'DELETE' })` in `users.html`.
• Tests: plan store and controller delete tests (added in Task 4 execution).

### Task 3: Agent Mode + GitHub MCP — Create Issue

• Configure GitHub MCP with PAT (`repo` scope) in VS Code.
• Start Copilot Agent Mode with GitHub MCP provider.
• Create issue in `your-username/go-users-demo`:
  ◦ Title: `Add DELETE /users/:id endpoint and store delete()`.
  ◦ Body: summary, requirements (store delete, controller mapping, README updates, unit tests), files, acceptance criteria.
  ◦ Labels: `enhancement`, `backend`. Assignee: yourself. Verify on GitHub.

### Task 4: Agent Mode — Implement Unit Test Plan (Read Upstream Issue as Prompt)

• Start by reading the upstream tracking issue content (created in Task 3) and use it verbatim as the Agent Mode prompt. This ensures the agent executes exactly against the acceptance criteria and file targets from the issue.
• Validate Go: `go version`; install deps via `go mod tidy` if needed.
• Implement `user_store_test.go`: seed order, `findByID`, `create` defaults, `update` behavior and exception.
• Implement `user_controller_test.go`: `net/http/httptest`, Gin test context for all endpoints.
• Run tests: `go test ./...`; fix failures aligned with `UserWriteModel` semantics.
• Commit/push: `git add models/user_store_test.go controllers/user_controller_test.go` → commit → push → open PR referencing the issue.

### Task 5: Mission Control — Observe Agent Execution

• Open Mission Control panel; start Agent run for Task 4.
• Observe file edits (test files), tool calls, and `go test`.
• Checkpoints: after each test file creation, after successful test run.
• Intervene with guidance if the agent stalls or assertions fail.
• Export logs/diffs for review.

### Task 6: Code Review — Review with Copilot

• Enable Copilot Chat.
• Summarize diffs in `user_controller.go`, `routes.go`, and tests; ask Copilot for potential issues.
• Validate: `go test ./...`; ask Copilot to explain failures and propose minimal fixes.
• Checklist: coverage (create/update/find/delete), 404 mapping and status codes (201/200/204), minimal deps, README updates.
• Provide PR feedback using Copilot suggestions; request changes or approve when criteria are met.

## Troubleshooting

• Go not found: Install Go or ensure Go is on PATH (`go version`).
• Version issues: Ensure Go 1.21+ is active (`go version`); configure VS Code Go extension accordingly.
• Build fails on emojis: Keep Unicode only where constants exist; otherwise preserve ASCII; ensure source encoding is UTF-8 (Go files are UTF-8 by default).
• Test failures due to `UserWriteModel` semantics: Verify fields are treated as nullable in store and adjust tests; ensure defaults match store logic.
• 404/response mismatches: Confirm not found errors are mapped to 404 in controller; for DELETE, return 204 on success.
• `go test` command hangs or network issues: Run `go clean -testcache` to refresh test cache; check proxy/network settings; run `go test -v ./...` for verbose output.

## Additional Links

- Code
- Issues
- Pull requests
- Actions
- Projects
- Security
