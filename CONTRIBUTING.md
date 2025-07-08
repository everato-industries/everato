# Contributing to Everato

Thank you for considering contributing to Everato! This document provides guidelines and instructions to help you contribute effectively to this project.

## Table of Contents

1. [Development Environment Setup](#development-environment-setup)
2. [Code Style and Standards](#code-style-and-standards)
3. [Branch Naming Conventions](#branch-naming-conventions)
4. [Commit Message Guidelines](#commit-message-guidelines)
5. [Pull Request Process](#pull-request-process)
6. [Reporting Bugs](#reporting-bugs)
7. [Feature Requests](#feature-requests)
8. [Testing](#testing)
9. [Project Structure](#project-structure)
10. [Communication](#communication)

## Development Environment Setup

Follow these steps to set up your development environment for Everato:

1. **Fork the repository**: Create a personal copy of the Everato repository on GitHub.
2. **Clone your fork**: Clone the repository to your local machine.
3. **Set up environment variables**: Copy the `.env.example` file to `.env` and configure it with your settings.
4. **Install development tools**: Use `make install` to install necessary dependencies and tools.
5. **Start the database**: Use `make db` to start the PostgreSQL database.
6. **Run migrations**: Apply database migrations using `make migrate-up`.
7. **Run the application**: Start the Everato application in development mode with `make dev`.
8. **Access the application**: Open your web browser and navigate to `http://localhost:8080` to see the running application.

### Prerequisites

- Go 1.24+
- PostgreSQL 15+
- Docker and Docker Compose (for development environment)
- Make (for running development commands)
- Node.js and npm/pnpm (for TailwindCSS compilation)

### Getting Started

1. Fork the repository and clone your fork:

    ```bash
    git clone https://github.com/YOUR_USERNAME/everato.git
    cd everato
    ```

2. Add the original repository as upstream:

    ```bash
    git remote add upstream https://github.com/dtg-lucifer/everato.git
    ```

3. Set up environment variables:

    ```bash
    cp .env.example .env
    # Edit .env file with your configuration
    ```

4. Install development tools:

    ```bash
    make install
    ```

5. Start the database:

    ```bash
    make db
    ```

6. Run migrations:

    ```bash
    make migrate-up
    ```

7. Run the application in development mode:

    ```bash
    make dev
    ```

## Code Style and Standards

Everato follows strict code style guidelines to maintain code quality and consistency:

### Go Code

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Format your code using `gofmt` or `go fmt`
- Run `golangci-lint` before submitting code
- Document all exported functions, types, and methods
- Organize imports alphabetically with standard library imports first

### HTML/CSS/JavaScript

- Use 4 spaces for indentation in HTML and CSS
- Format HTML templates with appropriate indentation
- Follow BEM naming conventions for CSS classes

### SQL

- Use uppercase for SQL keywords
- Format queries with appropriate indentation and line breaks
- Add comments for complex queries

### General Guidelines

- Write clear, descriptive comments for complex logic
- Keep functions focused on a single responsibility
- Avoid deep nesting of control structures
- Prefer explicit error handling over implicit one
- Include appropriate logging at different severity levels

## Branch Naming Conventions

Follow these naming conventions for branches:

- `feature/short-description` - For new features
- `bugfix/issue-number-short-description` - For bug fixes
- `refactor/component-name` - For code refactoring
- `docs/what-changed` - For documentation updates
- `test/what-tested` - For adding or updating tests

## Commit Message Guidelines

Write clear, meaningful commit messages:

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests after the first line
- Consider using the following format:

```
type(scope): Short description

Longer description if needed, explaining the context and motivation.

Fixes #123
```

Where `type` can be:

- feat: (new feature)
- fix: (bug fix)
- docs: (documentation changes)
- style: (formatting, missing semi-colons, etc; no code change)
- refactor: (refactoring production code)
- test: (adding missing tests, refactoring tests)
- chore: (updating build tasks, package manager configs, etc)

## Pull Request Process

1. Ensure your branch is up to date with the main branch
2. Run all tests and ensure they pass
3. Format your code according to the style guidelines
4. Create a pull request with a clear title and description
5. Reference any relevant issues
6. Wait for code review and address any feedback

### Pull Request Template

When creating a pull request, include:

- A description of the changes
- The motivation behind the changes
- Any breaking changes
- Screenshots (if applicable)
- Steps to test the changes

## Reporting Bugs

When reporting bugs, include:

1. A clear, descriptive title
2. Steps to reproduce the issue
3. Expected behavior
4. Actual behavior
5. Environment information (OS, browser, Go version, etc.)
6. Any relevant logs or screenshots

## Feature Requests

When requesting features, include:

1. A clear, descriptive title
2. A detailed description of the proposed feature
3. The motivation behind the feature
4. Any alternatives you've considered
5. Example use cases

## Testing

- Write tests for all new features and bug fixes
- Ensure all tests pass before submitting a pull request
- Include unit tests, integration tests, and end-to-end tests as appropriate
- Run `make test` to execute the test suite

### Testing Guidelines

- Use table-driven tests where appropriate
- Mock external dependencies in unit tests
- Write clear test descriptions
- Test edge cases and error conditions

## Project Structure

Everato follows a well-organized directory structure:

```
everato/
├── assets/                # Project assets like architecture diagrams
├── components/            # UI components for templ rendering
├── config/                # Configuration management
├── docker/                # Docker-related files for development
├── internal/              # Private application code
│   ├── db/                # Database-related code
│   ├── handlers/          # HTTP request handlers
│   ├── middlewares/       # HTTP middleware components
│   ├── services/          # Business logic
│   └── utils/             # Utility functions
├── pages/                 # Page templates (templ)
├── pkg/                   # Shared public libraries
├── public/                # Static assets (served directly)
├── scripts/               # Utility scripts
├── styles/                # Source CSS files (TailwindCSS)
└── templates/             # HTML templates
```

Follow this structure when adding new code to the project.

## Communication

- Use GitHub Issues for bug reports and feature requests
- Use Pull Requests for code contributions and reviews
- Follow a respectful and inclusive communication style
- Be patient with responses and feedback

### Code of Conduct

By participating in this project, you agree to abide by the following principles:

- Be respectful and inclusive
- Exercise empathy and kindness
- Be open to constructive feedback
- Focus on what's best for the community
- Show courtesy and respect in all interactions

## License

By contributing to Everato, you agree that your contributions will be licensed under the project's [license](LICENSE).

---

Thank you for contributing to Everato! Your efforts help make this project better for everyone.
