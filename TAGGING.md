# Tagging and Versioning Guide

## Repository Structure

**Important:** `shared-password-hasher` is an **independent GitHub repository** (`github.com/Bengo-Hub/shared-password-hasher`) in the Bengo-Hub organization. Each BengoBox service is also an independent repository. The `BengoBox` folder is just a local root directory where developers clone repositories.

## Semantic Versioning

Follow semantic versioning (MAJOR.MINOR.PATCH):

- **MAJOR**: Breaking changes (e.g., changing hash format, parameters)
- **MINOR**: New features (e.g., adding new methods)
- **PATCH**: Bug fixes, performance improvements

## Tagging the Library

### Step 1: Tag the Repository

```bash
# In the shared-password-hasher repository
cd shared/password-hasher/

# Create and push tag
git tag v0.1.0 -m "Initial release: Argon2id password hasher"
git push origin v0.1.0
```

### Step 2: Update Service go.mod Files

Each Go service should import the library:

```go
require (
    github.com/Bengo-Hub/shared-password-hasher v0.1.0
)
```

### Step 3: Update Auth Service

Update auth-service to use the shared library:

```bash
cd auth-service/auth-api/
go get github.com/Bengo-Hub/shared-password-hasher@v0.1.0
go mod tidy
```

Replace internal password package with shared library:

```go
// Before
import "github.com/bengobox/auth-service/internal/password"

// After
import passwordhasher "github.com/Bengo-Hub/shared-password-hasher"
```

### Step 4: Update Other Services

```bash
# TruLoad Go services
cd TruLoad/truload-sync-service/
go get github.com/Bengo-Hub/shared-password-hasher@v0.1.0

# Ordering service
cd ordering-service/ordering-backend/
go get github.com/Bengo-Hub/shared-password-hasher@v0.1.0

# Notifications service
cd notifications-service/
go get github.com/Bengo-Hub/shared-password-hasher@v0.1.0
```

## Local Development Setup

When developing locally, use `go.work` at the `BengoBox` root:

```bash
# Create parent directory
mkdir -p BengoBox/shared
cd BengoBox/

# Clone repositories
git clone https://github.com/Bengo-Hub/shared-password-hasher.git shared/password-hasher
git clone https://github.com/Bengo-Hub/auth-service.git auth-service
# ... clone other services

# Create go.work file
go work init
go work use ./shared/password-hasher
go work use ./auth-service/auth-api
go work use ./ordering-service/ordering-backend
go work use ./notifications-service
# ... add other services
```

Example `go.work` file:

```go
go 1.24.0

use (
    ./shared/password-hasher
    ./shared/auth-client
    ./shared/events
    ./auth-service/auth-api
    ./ordering-service/ordering-backend
    ./TruLoad/truload-backend
    ./notifications-service
)
```

## Version History

### v0.1.0 (Initial Release)
- Argon2id password hashing with default parameters
- Compatible with auth-service implementation
- PHC string format support
- Constant-time verification
- Comprehensive test suite

## Upgrading

When a new version is released:

```bash
# Update to latest version
go get github.com/Bengo-Hub/shared-password-hasher@latest
go mod tidy

# Or specific version
go get github.com/Bengo-Hub/shared-password-hasher@v0.2.0
go mod tidy
```

## Breaking Changes

If making breaking changes (MAJOR version bump):

1. Document migration guide in CHANGELOG
2. Provide backward compatibility period
3. Update all services gradually
4. Announce in team channels

Example breaking change:
```
v1.0.0 -> v2.0.0: Changed hash format from PHC to custom format
Migration: Use v1.x for 6 months, then migrate to v2.x
```

## CI/CD Integration

### GitHub Actions

Create `.github/workflows/release.yml`:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      - name: Run tests
        run: go test -v ./...
      
      - name: Create Release
        uses: actions/create-release@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          tag_name: ${{ github.ref }}
          release_name: Release ${{ github.ref }}
          draft: false
          prerelease: false
```

## Verification

After tagging and updating services:

```bash
# Verify version
go list -m github.com/Bengo-Hub/shared-password-hasher

# Test import
go run -mod=readonly ./cmd/api

# Run tests
go test ./...
```

## Best Practices

1. **Always test before tagging**: Run full test suite
2. **Document changes**: Update CHANGELOG.md
3. **Announce updates**: Notify team of new versions
4. **Gradual rollout**: Update dev → staging → production
5. **Monitor errors**: Watch for hash verification failures after updates
