## ADDED Requirements

### Requirement: Admin dashboard has a Music section
The admin dashboard SHALL include a Music section listing all recs (published and draft) with edit and delete actions, and a button to add a new rec.

#### Scenario: Music section shows all recs
- **WHEN** the admin loads the dashboard
- **THEN** the Music section lists every rec with its title, artist, and published status

#### Scenario: Empty state in Music section
- **WHEN** no recs exist
- **THEN** the Music section shows an empty state prompt to add the first rec

---

### Requirement: Admin can create a new music rec
The admin SHALL be able to open an inline form and submit a new rec with title, artist, optional album, optional cover URL, optional Spotify URL, optional Apple Music URL, optional note, and a published toggle.

#### Scenario: New rec is created and appears in the list
- **WHEN** the admin fills in at least title and artist and submits the form
- **THEN** the new rec appears in the Music section list and the form closes

#### Scenario: Submission blocked when required fields are empty
- **WHEN** the admin submits the form with title or artist empty
- **THEN** the form shows a validation error and does not submit

---

### Requirement: Admin can edit an existing music rec
Each rec in the Music section SHALL have an Edit action that opens the inline form pre-filled with the rec's current values.

#### Scenario: Edit form is pre-filled
- **WHEN** the admin clicks Edit on a rec
- **THEN** the inline form opens with all fields populated from the rec's current data

#### Scenario: Saving edits updates the rec
- **WHEN** the admin changes one or more fields and saves
- **THEN** the rec in the list reflects the updated values

---

### Requirement: Admin can toggle published status
Each rec SHALL have a quick-toggle for published/unpublished without opening the full edit form.

#### Scenario: Toggling published status
- **WHEN** the admin clicks the publish toggle on a rec
- **THEN** the rec's published status flips and the badge updates immediately

---

### Requirement: Admin can delete a music rec
Each rec in the Music section SHALL have a Delete action protected by a confirmation dialog.

#### Scenario: Delete requires confirmation
- **WHEN** the admin clicks Delete on a rec
- **THEN** a confirmation dialog appears before the rec is removed

#### Scenario: Confirmed delete removes the rec
- **WHEN** the admin confirms the delete
- **THEN** the rec is removed from the list and from the public page
