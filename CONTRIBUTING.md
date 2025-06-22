
# Contributing to TPC Portal Backend

First of all, thank you for your interest in contributing to the Official TPC Portal of CCET. This document outlines the contribution process to help maintainers and contributors collaborate effectively.

## Table of Contents
- [How to Contribute](#how-to-contribute)
- [Getting Started](#getting-started)
- [Code Guidelines](#code-guidelines)
- [Submitting Changes](#submitting-changes)
- [Reporting Issues](#reporting-issues)

## How to Contribute
1. Fork the repository and clone your fork.
2. Create a new branch for your feature or bug fix:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. Make your changes, ensuring that the code adheres to the project's style guide.
4. Test your changes thoroughly.
5. Push your changes to your fork and submit a pull request.

## Getting Started
To get started with contributing:
1. Ensure you have Go, Docker, and Nginx installed.
2. Clone the repository:
   ```bash
   git clone https://github.com/DevSoc-exe/placement-portal-backend
   cd placement-portal-backend
   ```
3. Run the backend server:
   ```bash
   go run main.go
   ```

## Code Guidelines
- **Go Version:** Ensure your Go version matches the one specified in the `go.mod` file.
- **Linting:** Run `golangci-lint run` before submitting a pull request.
- **Commit Messages:** Write clear and concise commit messages that explain your changes.

## Submitting Changes
1. Ensure that your changes pass all existing tests.
2. Ensure that your code is well-documented.
3. Open a pull request with a clear title and description of your changes.
4. The project maintainers will review your changes and provide feedback if necessary.

## Reporting Issues
If you encounter any issues:
1. Search the issue tracker to see if the problem has already been reported.
2. If not found, create a new issue and provide a clear description of the problem.
   - Mention the environment (Go version, OS, etc.).
   - Provide steps to reproduce the issue if applicable.

## Thank You
We appreciate your time and effort in contributing to this project. Happy coding!
