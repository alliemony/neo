## ADDED Requirements

### Requirement: Music tab is accessible from site navigation
The site SHALL include a Music link in the main navigation that routes to `/music`.

#### Scenario: Navigating to music tab
- **WHEN** a visitor clicks the Music nav link
- **THEN** they are taken to `/music` and the page renders

---

### Requirement: Published recs are displayed as cards
The `/music` page SHALL display all published music recs as a list of cards, each showing title, artist, and optional album cover.

#### Scenario: Recs are shown when published recs exist
- **WHEN** a visitor loads `/music` and there is at least one published rec
- **THEN** each published rec is rendered as a card with title and artist visible

#### Scenario: Empty state when no recs are published
- **WHEN** a visitor loads `/music` and no recs are published
- **THEN** a friendly empty state message is shown instead of an empty list

---

### Requirement: Streaming service pills link to external services
Each rec card SHALL display pill-style buttons for any streaming URLs that are set (Spotify, Apple Music). Pills SHALL only appear if the corresponding URL is non-empty.

#### Scenario: Spotify pill shown when URL is set
- **WHEN** a rec has a non-empty `spotify_url`
- **THEN** a Spotify pill is rendered that opens `spotify_url` in a new tab

#### Scenario: Apple Music pill shown when URL is set
- **WHEN** a rec has a non-empty `apple_url`
- **THEN** an Apple Music pill is rendered that opens `apple_url` in a new tab

#### Scenario: No pill shown when URL is empty
- **WHEN** a rec has an empty or absent `spotify_url` or `apple_url`
- **THEN** the corresponding pill is not rendered

---

### Requirement: Lyrics section expands on demand
Each rec card SHALL have a collapsible section that fetches and displays lyrics when opened.

#### Scenario: Lyrics load on first expand
- **WHEN** a visitor clicks to expand the lyrics section of a rec
- **THEN** a loading indicator is shown, then the lyrics text is rendered

#### Scenario: Graceful error when lyrics unavailable
- **WHEN** the lyrics API returns an error or the song is not found
- **THEN** a "Lyrics unavailable" message is shown instead of an empty section

#### Scenario: Lyrics section collapses on second click
- **WHEN** the lyrics section is open and the visitor clicks the toggle again
- **THEN** the section collapses and hides the lyrics

#### Scenario: Lyrics are cached for the session
- **WHEN** a visitor expands lyrics, collapses, then expands again
- **THEN** lyrics are shown immediately without a second network request

---

### Requirement: Optional personal note is displayed
If a rec has a non-empty `note` field, the card SHALL display it as a short personal annotation.

#### Scenario: Note shown when present
- **WHEN** a rec has a non-empty `note`
- **THEN** the note text is visible on the card

#### Scenario: Note omitted when absent
- **WHEN** a rec has no `note`
- **THEN** no note element is rendered on the card
