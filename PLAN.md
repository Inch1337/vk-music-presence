# VK Music Presence — Development Plan

**Project:** VK Music Presence  
**Status:** Planning / Initial Setup  
**Target Platform:** Windows 10/11  
**Initial Supported Service:** VK Music  
**Target Release:** v1.0.0

---

## 1. Project Goal

Build a lightweight, open-source Windows desktop application that detects music playback and automatically displays the current track in Discord Rich Presence.

The application should:

- Work in the background using the system tray.
- Support VK Music in Google Chrome and Chrome-installed PWA applications.
- Display the current track and artist in Discord.
- Synchronize playback changes.
- Handle pauses, stopped playback, and Discord disconnections.
- Require no VK account credentials or VK API tokens.
- Allow users to enable or disable supported music services.
- Support additional music services in future versions.
- Provide a simple and modern desktop interface.

The initial implementation focuses exclusively on VK Music.

The architecture must remain independent of any particular music service.

## 2. Development Principles

1. Build functionality incrementally.
2. Validate technical assumptions before implementing complex architecture.
3. Keep Go core logic independent of the desktop framework.
4. Prefer simple, idiomatic Go solutions.
5. Separate media acquisition, service detection, source selection, and Discord synchronization.
6. Avoid unnecessary abstractions, dependencies, and infrastructure.
7. Add automated tests where they provide meaningful value.
8. Keep documentation synchronized with implemented functionality.
9. Use small, meaningful Git commits.
10. Prioritize understanding and learning over implementation speed.

## 3. Technology Stack

| Component | Technology |
|---|---|
| Core application | Go |
| Desktop framework | Wails v3 |
| Frontend | Svelte + TypeScript |
| Frontend tooling | Vite |
| Styling | CSS |
| Media integration | Windows GSMTC / WinRT |
| Optional media integration | Chrome Extension + Native Messaging |
| Discord integration | Local Discord RPC |
| Configuration | JSON |
| Logging | Go log/slog |
| Go testing | go test |
| Frontend testing | Vitest |
| CI/CD | GitHub Actions |
| Distribution | GitHub Releases |

Specific WinRT and Discord libraries will be selected after technical validation.

Wails v3 is currently a beta dependency. Pin a tested version and verify desktop functionality before relying on it for release builds.

## 4. Development Roadmap

### Phase 0 — Project Initialization

**Goal:** Prepare the repository and development environment.

Tasks:

- [x] Create the project directory.
- [x] Initialize the Git repository.
- [x] Create the initial README.
- [x] Create the initial Git commit.
- [x] Create AGENTS.md.
- [x] Create PLAN.md.
- [x] Create docs/ARCHITECTURE.md.
- [ ] Verify Go installation and version.
- [ ] Verify Node.js and npm.
- [ ] Verify Git configuration.
- [ ] Install the Wails v3 CLI.
- [ ] Initialize the Wails + Svelte + TypeScript project.
- [ ] Verify development mode.
- [ ] Verify the production build.
- [ ] Verify basic Go-to-TypeScript communication.
- [ ] Verify window hiding and system tray functionality.

**Definition of Done:**

The repository contains a working minimal Wails application with Svelte and TypeScript. Development and production builds work on Windows.

### Phase 1 — Windows Media Prototype

**Goal:** Investigate how Windows exposes music playback information.

Tasks:

- [ ] Research Windows GSMTC APIs.
- [ ] Evaluate Go-compatible WinRT libraries.
- [ ] Create a minimal media diagnostic program.
- [ ] Enumerate available media sessions.
- [ ] Read track title and artist.
- [ ] Read playback state.
- [ ] Read playback position and duration where available.
- [ ] Test VK Music in a regular Chrome tab.
- [ ] Test VK Music in a Chrome-installed PWA.
- [ ] Test playback changes and pauses.
- [ ] Test simultaneous media sessions.
- [ ] Investigate reliable VK Music identification.
- [ ] Document limitations and findings.

**Important technical question:**

Can we reliably distinguish VK Music from unrelated media playing in Chrome?

If GSMTC cannot provide sufficient information, investigate Chrome Extension + Native Messaging as an additional source of playback data.

**Definition of Done:**

A Go prototype can obtain playback information and a technically validated approach exists for identifying VK Music without falsely identifying unrelated services.

### Phase 2 — Discord RPC Prototype

**Goal:** Verify local Discord Rich Presence integration.

Tasks:

- [ ] Research the current Discord RPC documentation.
- [ ] Create a Discord application and obtain its Application ID.
- [ ] Evaluate Go RPC libraries or direct IPC implementation.
- [ ] Establish a local connection to Discord Desktop.
- [ ] Publish a test activity.
- [ ] Test the Listening activity type.
- [ ] Test track title and artist display.
- [ ] Investigate timestamps and playback progress.
- [ ] Update an existing activity.
- [ ] Clear an activity.
- [ ] Test Discord restarts and disconnections.
- [ ] Document supported features and limitations.

**Definition of Done:**

A standalone Go prototype successfully publishes and clears Discord activity, with verified behavior for the fields required by the MVP.

### Phase 3 — Domain Model and Core Architecture

**Goal:** Build a testable application core independent of Windows UI and Discord implementation details.

Planned components:

- Domain models.
- Media Provider.
- Service Detector.
- Source Selector.
- Presence Manager.
- Discord Client abstraction.
- Application lifecycle management.

Tasks:

- [ ] Define the Track model.
- [ ] Define playback states.
- [ ] Define music service identifiers.
- [ ] Define media session identifiers and source metadata.
- [ ] Design a minimal Media Provider abstraction.
- [ ] Implement VK Music identification.
- [ ] Implement filtering of enabled music services.
- [ ] Implement selection of one active source.
- [ ] Implement playback state synchronization.
- [ ] Implement track change detection.
- [ ] Prevent redundant Discord updates.
- [ ] Implement cancellation and graceful shutdown.
- [ ] Add unit tests for core behavior.

**Architecture requirements:**

- Media acquisition must not depend on Discord.
- Service detection must be separate from media acquisition.
- Core logic must not depend on Wails or Svelte.
- Unsupported music services must not be exposed as working integrations.
- Only one selected source may control Discord Presence at a time.
- Unknown or ambiguous sources must not be published as a specific music service.

**Definition of Done:**

The core synchronizes confirmed VK Music playback with Discord and can be tested independently of the desktop UI.

### Phase 4 — Desktop Application

**Goal:** Integrate the Go core with Wails and build the user interface.

Tasks:

- [ ] Connect the Go core to Wails.
- [ ] Implement typed Go-to-TypeScript bindings.
- [ ] Implement application state events.
- [ ] Display the current track and artist.
- [ ] Display Discord connection status.
- [ ] Implement enable/disable Presence controls.
- [ ] Implement a basic settings interface.
- [ ] Implement the system tray.
- [ ] Hide the window without terminating background processing.
- [ ] Implement an explicit application exit action.
- [ ] Display supported music services.
- [ ] Implement enabling/disabling supported services.
- [ ] Add basic frontend tests.

**Definition of Done:**

A user can launch VK Music Presence, manage it through the tray, view playback information, and configure the supported functionality without using a terminal.

### Phase 5 — Reliability and User Experience

**Goal:** Prepare the application for everyday use.

Tasks:

- [ ] Handle Discord connection failures.
- [ ] Implement automatic reconnection.
- [ ] Handle disappeared media sessions.
- [ ] Handle pause and resume.
- [ ] Handle track seeking where supported.
- [ ] Improve multiple-session selection.
- [ ] Prevent unnecessary Presence updates.
- [ ] Implement JSON configuration.
- [ ] Store settings in the appropriate Windows user directory.
- [ ] Implement structured logging.
- [ ] Implement Windows autostart.
- [ ] Handle graceful application shutdown.
- [ ] Test prolonged background operation.
- [ ] Measure CPU and memory usage.
- [ ] Add regression tests for reported bugs.

**Definition of Done:**

The application remains stable during normal playback, browser restarts, Discord restarts, and extended background operation.

### Phase 6 — Public Release

**Goal:** Publish a usable open-source Windows application.

Tasks:

- [ ] Finalize the user-facing README.
- [ ] Update the architecture documentation.
- [ ] Document installation and usage.
- [ ] Document known limitations.
- [ ] Choose an open-source license.
- [ ] Add GitHub Actions for automated tests.
- [ ] Add automated Windows builds.
- [ ] Verify production packaging and runtime dependencies.
- [ ] Test on a clean Windows installation.
- [ ] Prepare application icons and screenshots.
- [ ] Publish a release candidate.
- [ ] Fix release-blocking issues.
- [ ] Publish v1.0.0.

**Definition of Done:**

A user can download VK Music Presence from GitHub, run it without installing development tools, and use the supported VK Music integration.

---

## 5. Future Development

The following features are outside the initial MVP.

### Album Artwork

- Investigate artwork metadata availability.
- Investigate Discord Rich Presence image restrictions.
- Implement supported artwork display if technically feasible.

### Additional Music Services

Potential candidates:

- Yandex Music.
- YouTube Music.
- Other browser-based music services.

Each integration must provide reliable service identification before being enabled.

### Advanced Settings

- Source priority configuration.
- Additional playback behavior settings.
- Advanced diagnostic controls.
- More appearance customization.

### Other Platforms

Linux and macOS support may be investigated after the Windows implementation is stable.

Cross-platform support is not a v1.0 requirement.

---

## 6. Technical Risks

### 6.1. Chrome Media Session Identification

Windows may expose browser media sessions without identifying the specific website.

**Mitigation:** Validate GSMTC data experimentally. Introduce an optional Chrome extension only if required for reliable identification.

### 6.2. WinRT Integration with Go

Go bindings may introduce compatibility or API coverage limitations.

**Mitigation:** Evaluate libraries using an isolated prototype. Consider a small platform adapter if necessary.

### 6.3. Discord RPC Restrictions

Custom Rich Presence may not behave exactly like Spotify integration.

**Mitigation:** Test supported activity types, timestamps, and asset behavior before implementing dependent features.

### 6.4. Wails v3 Beta

Wails v3 may introduce breaking changes or platform-specific issues.

**Mitigation:** Pin a tested version, avoid unnecessary dependency updates, and validate tray/background behavior.

### 6.5. Multiple Music Sources

Multiple browser sessions may play simultaneously.

**Mitigation:** Define explicit selection rules and maintain only one active Discord Presence source.

---

## 7. Testing Strategy

### Unit Testing

Test business logic independently of external integrations:

- Playback transitions.
- Track changes.
- Service detection.
- Source selection.
- Enabled service filtering.
- Duplicate update prevention.
- Error handling.

### Integration Testing

Test real Windows and Discord behavior:

- Chrome tab playback.
- Chrome PWA playback.
- Multiple media sessions.
- Pause and resume.
- Track switching.
- Browser closure.
- Discord restart.
- Application shutdown.

### Desktop Testing

Verify:

- Wails window lifecycle.
- System tray behavior.
- Configuration persistence.
- Background operation.
- Production builds.

---

## 8. Versioning

Planned version milestones:

| Version | Focus |
|---|---|
| Prototype | Technical validation |
| 0.1.0 | VK Music to Discord MVP |
| 0.2.0 | Stable synchronization |
| 0.3.0 | Desktop UI and tray |
| 0.4.0 | Usability and reliability improvements |
| 1.0.0 | Public release |

Versions describe functional milestones, not fixed deadlines.

Release contents may be adjusted as implementation progresses.

---

## 9. Documentation Workflow

- `AGENTS.md` defines how Codex assists with development.
- `PLAN.md` tracks development goals and progress.
- `docs/ARCHITECTURE.md` documents components, dependencies, and technical decisions.
- `README.md` explains the application to users and contributors.

Checklist items are marked as complete only after implementation and verification.

The documentation must distinguish planned functionality from functionality that already exists.

Architectural changes should be documented when decisions are made.

---

## 10. Current Next Step

Complete Phase 0.

After project documentation is prepared:

1. Verify the development environment.
2. Install and initialize Wails v3.
3. Start a minimal Svelte application.
4. Verify frontend-to-Go communication.
5. Verify tray and window lifecycle behavior.

Then proceed to Windows Media API prototyping.

**Current priority:** Establish a working development environment without prematurely implementing application features.