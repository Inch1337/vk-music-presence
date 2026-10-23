# VK Music Presence — Architecture

**Status:** Proposed Architecture  
**Platform:** Windows 10/11  
**Architecture Style:** Modular Monolith  
**Core Language:** Go  
**Desktop Framework:** Wails v3  
**Frontend:** Svelte + TypeScript

> This document describes the intended architecture of VK Music Presence.
> Components marked as planned may not exist in the codebase yet.
> Implementation details will be refined through technical prototypes.

---

## 1. System Overview

VK Music Presence is a lightweight desktop application that detects music playback from supported services and synchronizes the current track with Discord Rich Presence.

The application runs locally on the user's computer.

Its primary responsibilities are:

1. Acquire playback information.
2. Identify the music service.
3. Filter services according to user preferences.
4. Select one active playback source.
5. Synchronize the selected track with Discord.
6. Expose application state and settings through a desktop interface.

The initial supported service is VK Music running in Google Chrome or a Chrome-installed PWA.

Future versions may support additional services without requiring major changes to the application core.

## 2. Architectural Principles

### Separation of Concerns

Each component should have a clearly defined responsibility.

Media acquisition, service identification, source selection, Discord integration, and user interface logic must remain separate.

### Dependency Direction

Core application logic must not depend directly on:

- Wails
- Svelte
- Windows-specific APIs
- Discord RPC implementation details

External integrations should depend on domain models and appropriate abstractions.

### Simplicity

Avoid introducing abstractions without a clear purpose.

Interfaces should primarily exist at meaningful dependency boundaries, not for every struct or function.

### Extensibility

Music service support must be extensible.

Adding a new service should not require rewriting the Discord client or core playback synchronization.

### Testability

Business logic should be testable without a running Discord client, browser, or desktop UI.

---

## 3. High-Level Architecture

```text
                  MUSIC SOURCES
                       |
          +------------+------------+
          |                         |
     Windows GSMTC            Chrome Extension
          |                     (optional)
          |                         |
          +------------+------------+
                       |
                Media Acquisition
                       |
                Service Detection
                       |
                Source Selection
                       |
                Presence Manager
                       |
                Discord Adapter
                       |
                 Discord Desktop


                APPLICATION CORE
                       |
              Desktop Integration
                       |
                  Wails v3
                       |
              Svelte + TypeScript
                       |
                  User Interface
```

The desktop interface controls the application but does not directly interact with media providers or Discord RPC.

## 4. Domain Models

The domain layer defines application concepts independently of external libraries.

### 4.1. Track

Represents a musical composition.

Planned properties:

- Title
- Artist
- Album
- Duration
- Playback position

Optional properties such as artwork and track URLs may be added later.

### 4.2. Playback State

Represents playback status.

Possible values:

- Playing
- Paused
- Stopped
- Unknown

Unknown must be handled explicitly rather than automatically treated as Playing.

### 4.3. Music Service

Identifies the music service associated with playback.

Initial service identifiers:

- `vk`

Potential future identifiers:

- `yandex`
- `youtube_music`

Only implemented and verified services may be enabled in user settings.

### 4.4. Media Observation

Represents information reported by a media provider.

Planned properties:

- Media session identifier
- Source application identifier
- Track metadata
- Playback state
- Playback timeline
- Observation timestamp
- Data source

Media observations are not automatically associated with a specific music service.

Service identity must be established separately.

### 4.5. Identified Playback

Represents playback associated with a verified music service.

Contains:

- Music service identifier
- Source/session identity
- Track information
- Playback state
- Observation time
- Evidence of service identification

The precise Go types will be designed during implementation.

---

## 5. Core Components

### 5.1. Media Provider

**Responsibility:** Acquire media playback information.

The initial provider will use Windows Global System Media Transport Controls (GSMTC) through WinRT.

Expected capabilities:

- Enumerate media sessions.
- Read track metadata.
- Read playback state.
- Read available timeline information.
- Detect playback changes.
- Detect session disappearance.

The provider must not assume that every Chrome media session belongs to VK Music.

The exact WinRT library and event-handling strategy will be selected after prototyping.

### 5.2. Service Detector

**Responsibility:** Identify which supported music service owns a media observation.

The detector should use reliable source information.

Possible evidence includes:

- Distinct application identity.
- Verified browser or PWA metadata.
- Data reported by a trusted service-specific integration.

Track names, artist names, or the fact that Chrome is playing audio are not sufficient proof of service identity.

If identification is ambiguous, the observation must remain unidentified.

The first implementation will focus on VK Music.

### 5.3. Source Selector

**Responsibility:** Select one active playback source from the identified candidates.

Selection rules:

1. Consider only supported and enabled services.
2. Reject unidentified sources.
3. Prefer actively playing sources.
4. Preserve the currently selected source while it remains eligible.
5. Handle multiple simultaneous sources deterministically.
6. Return no active source when no eligible playback exists.

Future versions may support configurable service priorities.

The selector must not depend on Discord-specific data structures.

### 5.4. Presence Manager

**Responsibility:** Coordinate synchronization between selected playback and Discord.

Expected behavior:

- Track current application playback state.
- Detect relevant changes.
- Build a representation of the desired Discord activity.
- Avoid unnecessary duplicate updates.
- Handle pauses and stopped playback.
- Clear stale activities.
- React to configuration changes.
- Coordinate reconnection and shutdown.

The manager should process state transitions consistently.

It must not call Windows-specific APIs directly.

### 5.5. Discord Adapter

**Responsibility:** Manage local Discord Rich Presence communication.

Expected capabilities:

- Connect to Discord Desktop using local IPC.
- Publish an activity.
- Update an activity.
- Clear an activity.
- Detect disconnection.
- Reconnect when possible.
- Handle errors and supported update limits.

A Discord Application ID will be configured for the integration.

The initial prototype will verify support for the `Listening` activity type, timestamps, and relevant metadata.

The implementation must not assume that custom Rich Presence behaves identically to Spotify integration.

The exact Go RPC implementation will be selected after validation.

### 5.6. Configuration Service

**Responsibility:** Manage local application settings.

Planned settings:

- Discord Presence enabled/disabled.
- Enabled music services.
- Playback pause behavior.
- Windows autostart.
- Logging preferences.

Settings will be stored in JSON under an appropriate user-specific Windows directory.

Configuration loading should support validation and safe defaults.

### 5.7. Desktop Integration

**Responsibility:** Connect the application core with the desktop environment.

Technology: Wails v3.

Expected responsibilities:

- Application lifecycle.
- Window management.
- System tray integration.
- Go-to-TypeScript method bindings.
- Application events.
- User actions and settings updates.

Closing or hiding the main window should not terminate the core background processes.

An explicit Exit action should perform graceful shutdown.

### 5.8. Frontend

**Technologies:**

- Svelte
- TypeScript
- Vite
- CSS

The frontend is a presentation layer.

It should display:

- Current track and artist.
- Playback state.
- Discord connection state.
- Enabled music services.
- Basic application settings.

The frontend interacts with Go using Wails bindings and events.

It should obtain an initial application state snapshot and then receive subsequent changes.

The frontend must not communicate directly with Discord RPC or Windows GSMTC.

---

## 6. Multiple Music Service Support

Music service support is a core architectural consideration.

However, only VK Music will be implemented initially.

The application should separate:

**Media acquisition:** How playback information is obtained.

**Service identification:** Which music service is responsible.

**Source selection:** Which identified service controls Discord Presence.

For example, VK Music and Yandex Music may both use Chrome as their playback environment.

A shared media provider may obtain metadata for both, while service-specific detection logic identifies each one.

Additional integrations may provide more reliable source information when required.

Adding a new music service should mainly involve implementing and registering its identification or acquisition integration.

The core Presence Manager and Discord Adapter should remain unchanged.

---

## 7. Optional Chrome Extension

Windows GSMTC may not reliably identify the specific website associated with a Chrome media session.

If this limitation prevents safe VK Music detection, an optional Chrome extension may be introduced.

The extension would:

- Run only with the permissions necessary for supported music sites.
- Identify the active music service.
- Extract or confirm relevant playback information.
- Communicate with the local application through Chrome Native Messaging.

The native application must validate incoming messages.

Extension-reported data and GSMTC observations must not be blindly combined if their session identities cannot be correlated.

The extension is not part of the initial implementation until the GSMTC prototype demonstrates that it is necessary.

---

## 8. Concurrency and Application Lifecycle

The Go core may use concurrent operations for:

- Media monitoring.
- Discord connection management.
- Application state updates.
- Desktop event propagation.

Concurrency should be introduced only when needed.

Use:

- `context.Context` for cancellation.
- Goroutines for independent long-running operations.
- Channels for appropriate event communication.
- Mutexes or single-owner state management where necessary.

Avoid uncontrolled goroutine creation.

Application shutdown should:

1. Stop accepting new work.
2. Cancel background operations.
3. Attempt to clear the Discord activity.
4. Close external connections.
5. Release Windows resources.
6. Exit the application.

Windows-specific threading and WinRT lifecycle requirements must be respected by the media adapter.

---

## 9. Error Handling

Failures in individual integrations should not unnecessarily terminate the application.

Examples:

- Discord is not running.
- Discord connection is lost.
- No media session is available.
- A media session disappears.
- Playback metadata is incomplete.
- Configuration is invalid.
- An enabled service cannot be identified.

Errors should be handled according to their severity.

Use structured logging through Go `log/slog`.

Avoid excessive log output during normal playback.

The user interface should display understandable status information rather than raw internal errors.

---

## 10. Testing Architecture

### Unit Tests

Focus on components that do not require real external integrations:

- Service selection.
- Playback state transitions.
- Track change detection.
- Enabled-service filtering.
- Duplicate update prevention.
- Handling missing or ambiguous metadata.

### Integration Tests

Validate actual behavior involving:

- Windows GSMTC.
- Google Chrome.
- Chrome PWA applications.
- Discord RPC.
- Wails desktop lifecycle.

### Test Doubles

Use simple fake or mock implementations for external boundaries when needed.

Do not introduce a large mocking framework unless it provides clear value.

---

## 11. Security and Privacy

The application should:

- Operate locally.
- Never request VK passwords.
- Avoid requiring VK access tokens.
- Avoid collecting playback history unnecessarily.
- Avoid exposing an unauthenticated local network API.
- Validate all messages from optional external integrations.
- Store configuration only in appropriate local user directories.
- Send only required activity information to Discord.

No dedicated backend server or cloud database is planned.

---

## 12. Repository Structure

Proposed logical structure:

```text
vk-music-presence/
├── cmd/
│   └── media-probe/
│
├── internal/
│   ├── domain/
│   ├── media/
│   ├── services/
│   │   └── vk/
│   ├── selection/
│   ├── presence/
│   ├── discord/
│   ├── config/
│   └── desktop/
│
├── frontend/
│   └── src/
│
├── assets/
│
├── docs/
│   └── ARCHITECTURE.md
│
├── .github/
│   └── workflows/
│
├── AGENTS.md
├── PLAN.md
├── README.md
├── go.mod
└── go.sum
```

This is a target structure, not a requirement to create every directory during initialization.

The official Wails project template may introduce additional files or a different initial layout.

Packages should be added gradually as real implementation needs emerge.

---

## 13. Architectural Decisions

### Accepted

- Modular monolith.
- Go as the main application language.
- Wails v3 with Svelte and TypeScript.
- Windows as the initial target platform.
- No external backend server.
- No database.
- Independent media acquisition and service identification.
- One selected source controlling Discord Presence.
- Multi-service extensibility without premature implementation.
- Core application logic independent of the desktop UI.

### Pending Technical Validation

- Specific Go WinRT library.
- Windows media session identification reliability.
- Need for a Chrome extension.
- Specific Discord RPC implementation.
- Actual Discord Listening activity presentation.
- Wails tray and background behavior in production builds.

These decisions must be finalized using working prototypes.

---

## 14. Implementation Status

**Current state:** Project initialization.

At the time of the initial architecture proposal:

- Git repository has been initialized.
- Initial README and commit have been created.
- Development documentation is being prepared.
- The application core has not been implemented.
- Wails initialization and technical prototypes are pending.

The actual implementation must always be verified against the repository rather than inferred from this document.

---

## 15. Architecture Maintenance

Update this document when:

- A component's responsibility changes.
- A new external integration is introduced.
- A significant dependency is selected or replaced.
- New architectural constraints are discovered.
- A proposed design becomes an implemented design.

Avoid documenting minor implementation details that are better represented by the code itself.

**The architecture should guide development without restricting necessary improvements.**