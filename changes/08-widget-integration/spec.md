# Spec: Python Widget Integration

## Purpose

Enable embedding interactive Python-based widgets (ML models, data visualizations, notebooks) within blog posts via a separate FastAPI service.

## Requirements

### Requirement: The widget service SHALL provide a registry of available widgets

#### Scenario: List available widgets

- **GIVEN:** The widget service is running with registered widgets
- **WHEN:** `GET /widgets` is called
- **THEN:** Response is 200 with a JSON array of widget metadata: `[{"id": "sentiment", "name": "Sentiment Analysis", "description": "..."}]`

#### Scenario: Get widget metadata by ID

- **GIVEN:** A widget with id "sentiment" is registered
- **WHEN:** `GET /widgets/sentiment` is called
- **THEN:** Response is 200 with the widget's metadata

### Requirement: Each widget SHALL serve a self-contained HTML page embeddable via iframe

#### Scenario: Widget renders as standalone HTML

- **GIVEN:** Widget "sentiment" is registered
- **WHEN:** `GET /widgets/sentiment/embed` is called
- **THEN:** Response is HTML with a self-contained UI (form + result area)
- **AND:** The page includes all necessary CSS and JS inline

### Requirement: A sample widget SHALL demonstrate the integration pattern

#### Scenario: Sentiment analysis widget works end-to-end

- **GIVEN:** The widget service is running
- **WHEN:** A user visits the sentiment widget embed and enters text
- **THEN:** The widget processes the input and displays a sentiment result
- **AND:** The interaction happens entirely within the iframe

### Requirement: WidgetEmbed frontend component SHALL render widgets via responsive iframe

#### Scenario: Widget loads in iframe

- **GIVEN:** A blog post contains a widget reference with id "sentiment"
- **WHEN:** The post view renders
- **THEN:** An iframe loads the widget's embed URL from the widget service
- **AND:** A loading indicator shows while the iframe loads

#### Scenario: Widget iframe is responsive

- **GIVEN:** The viewport resizes
- **WHEN:** The WidgetEmbed component renders
- **THEN:** The iframe scales to fit the container width

#### Scenario: Widget service unavailable shows fallback

- **GIVEN:** The widget service is not running
- **WHEN:** The WidgetEmbed component tries to load
- **THEN:** A fallback message is shown instead of a broken iframe

### Requirement: Admin post editor SHALL support the widget post type

#### Scenario: Widget post type references a widget ID

- **GIVEN:** Admin selects "widget" as the post type
- **WHEN:** They fill in the widget ID field
- **THEN:** The post stores the widget reference
- **AND:** The public post view embeds the widget

### Requirement: Widgets route (/widgets/:id) SHALL provide standalone widget viewing

#### Scenario: Direct widget access

- **GIVEN:** Widget "sentiment" exists in the service
- **WHEN:** The user visits `/widgets/sentiment`
- **THEN:** The widget is displayed full-width in the site layout

### Requirement: Widget service SHALL be independently deployable via Docker

#### Scenario: Widget service runs in isolation

- **GIVEN:** The widget Docker image is built
- **WHEN:** It is started as a standalone container
- **THEN:** It serves widgets without depending on the Go backend or frontend
