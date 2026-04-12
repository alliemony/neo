# Proposal: CI/CD with GitHub Actions

## Problem
No automated testing, building, or deployment pipeline exists. All testing and deployment is manual, increasing the risk of regressions and slowing down the development workflow.

## Solution
Add GitHub Actions workflows for:
1. **CI (Test)** — Run all tests (frontend, backend, widgets) on every push and PR
2. **Build** — Build the production Docker image on pushes to `main`
3. **Deploy** — Deploy to Fly.io on pushes to `main` (will fail until deployment secrets are configured)

## Scope
- Three workflow files in `.github/workflows/`
- No changes to application code
- Deploy workflow designed to fail gracefully without secrets

## Non-Goals
- CDN / asset pipeline
- Staging environments
- Multi-region deployment
