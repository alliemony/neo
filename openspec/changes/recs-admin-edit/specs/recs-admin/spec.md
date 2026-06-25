## ADDED Requirements

### Requirement: Admin dashboard has a Recs section
The admin dashboard SHALL include a Recs section listing all recs (published and draft) with their name, section label, published status badge, and Edit / Delete actions, plus a button to add a new rec.

#### Scenario: All recs are listed regardless of published status
- **WHEN** the admin loads the dashboard
- **THEN** the Recs section shows every rec with name, section, and published badge

#### Scenario: Empty state when no recs exist
- **WHEN** no recs exist in the database
- **THEN** the Recs section shows an empty-state prompt

---

### Requirement: Admin can create a new rec
The admin SHALL be able to open an inline form and submit a new rec with: name*, section* (text input with datalist of existing sections), tag*, tag_bg*, tag_color*, description*, href (optional), published toggle, sort_order.

#### Scenario: New rec appears in list after creation
- **WHEN** the admin fills all required fields and submits
- **THEN** the new rec appears in the Recs section list and the form closes

#### Scenario: Form blocks submission when required fields are empty
- **WHEN** the admin submits with any required field empty
- **THEN** a validation error is shown and the form does not submit

---

### Requirement: Admin can edit an existing rec
Each rec SHALL have an Edit action that opens the inline form pre-filled with current values.

#### Scenario: Edit form is pre-filled
- **WHEN** the admin clicks Edit on a rec
- **THEN** the inline form opens with all fields populated from the rec's data

#### Scenario: Saving edits updates the rec in the list
- **WHEN** the admin changes one or more fields and saves
- **THEN** the list reflects the updated values immediately

---

### Requirement: Admin can toggle published status
Each rec SHALL have a quick-toggle for published/unpublished without opening the full edit form.

#### Scenario: Published badge toggles on click
- **WHEN** the admin clicks the published badge/toggle on a rec
- **THEN** the rec's published status flips and the badge updates immediately

---

### Requirement: Admin can delete a rec
Each rec SHALL have a Delete action protected by a confirmation dialog.

#### Scenario: Delete requires confirmation
- **WHEN** the admin clicks Delete on a rec
- **THEN** a confirmation dialog appears before deletion

#### Scenario: Confirmed delete removes the rec
- **WHEN** the admin confirms the delete
- **THEN** the rec is removed from the list and from the public page
