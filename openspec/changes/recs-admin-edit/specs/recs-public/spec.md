## MODIFIED Requirements

### Requirement: Recs page displays recommendations from API
The `/recs` page SHALL fetch published recommendations from `GET /api/v1/recs` and render them grouped by section, visually identical to the previous hardcoded layout. The hardcoded `RECS` constant SHALL be removed.

#### Scenario: Published recs render grouped by section
- **WHEN** a visitor loads `/recs` and the API returns sections with items
- **THEN** each section heading and its rec items are displayed in order, matching the existing card layout (name, optional link, tag pill with bg/text colours, description)

#### Scenario: Loading state is shown while fetching
- **WHEN** a visitor loads `/recs` and the API call is in flight
- **THEN** a loading indicator is shown instead of the rec list

#### Scenario: Empty state when no published recs
- **WHEN** a visitor loads `/recs` and the API returns no sections
- **THEN** a friendly empty state message is shown

#### Scenario: Error state on API failure
- **WHEN** the API call fails
- **THEN** a graceful error message is shown rather than a blank page or console error
