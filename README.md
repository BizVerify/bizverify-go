# BizVerify Go SDK

Official Go SDK for the [BizVerify](https://bizverify.co) business entity verification API.

## Installation

```bash
go get github.com/bizverify/bizverify-go
```

## Quick Start

### Authentication (Passwordless)

```go
package main

import (
    "context"
    "fmt"
    "log"

    bv "github.com/bizverify/bizverify-go"
)

func main() {
    client := bv.New()
    ctx := context.Background()

    // Step 1: Request an OTP code via email
    _, err := client.Auth.RequestAccess(ctx, "user@example.com", true)
    if err != nil {
        log.Fatal(err)
    }

    // Step 2: Verify the code — client is auto-configured with the returned API key
    resp, err := client.Auth.VerifyAccess(ctx, "user@example.com", "123456", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("API Key:", resp.APIKey) // Store this for future use

    // All authenticated endpoints now work automatically
    account, err := client.Account.Get(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(account.CreditBalance)
}
```

### Using an Existing API Key

```go
client := bv.New(bv.WithAPIKey("bv_live_..."))
```

### Verify a Business Entity

```go
client := bv.New(bv.WithAPIKey("bv_live_..."))
ctx := context.Background()

// Synchronous verification (cached result)
resp, err := client.Verification.Verify(ctx, bv.VerifyParams{
    EntityName:   "Acme Inc",
    Jurisdiction: "us-fl",
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(resp.Status, string(resp.Data))

// Verify and wait for async job to complete
job, err := client.Verification.VerifyAndWait(ctx, bv.VerifyParams{
    EntityName:   "Acme Inc",
    Jurisdiction: "us-fl",
}, nil)
if err != nil {
    log.Fatal(err)
}
fmt.Println(job.Status, string(job.Result))
```

### Search for Entities

```go
// Single page
resp, err := client.Search.Find(ctx, bv.SearchParams{
    EntityName:   "Acme",
    Jurisdiction: "us-fl",
})
for _, r := range resp.Results {
    fmt.Println(r.EntityName, r.Confidence)
}

// Auto-paginate through all results
iter := client.Search.FindAll(ctx, bv.SearchParams{EntityName: "Acme"})
for iter.Next() {
    r := iter.Value()
    fmt.Println(r.EntityName)
}
if err := iter.Err(); err != nil {
    log.Fatal(err)
}
```

### Response Metadata (Credits & Rate Limits)

After any API call, inspect credits and rate limit info from response headers:

```go
resp, err := client.Verification.Verify(ctx, bv.VerifyParams{
    EntityName:   "Acme Inc",
    Jurisdiction: "us-fl",
})
if err != nil {
    log.Fatal(err)
}

meta := client.LastResponseMeta()
if meta.CreditsRemaining != nil {
    fmt.Printf("Credits remaining: %d\n", *meta.CreditsRemaining)
}
if meta.CreditsCharged != nil {
    fmt.Printf("Credits charged: %d\n", *meta.CreditsCharged)
}
```

### Configuration & Jurisdictions

```go
// Get full API configuration (no auth required)
config, err := client.Config.Get(ctx)

// List available jurisdictions (no auth required)
jurisdictions, err := client.Config.Jurisdictions(ctx)
for _, j := range jurisdictions.Jurisdictions {
    fmt.Printf("%s: %s\n", j.Code, j.Name)
}
```

### Error Handling

```go
import "errors"

entity, err := client.Entities.Get(ctx, "ent_nonexistent")
if err != nil {
    var notFound *bv.NotFoundError
    var noCredits *bv.InsufficientCreditsError
    var rateLimit *bv.RateLimitError

    switch {
    case errors.As(err, &notFound):
        fmt.Printf("Not found: %s (code=%s)\n", notFound.Message, notFound.Code)
    case errors.As(err, &noCredits):
        fmt.Println("Need more credits")
    case errors.As(err, &rateLimit):
        fmt.Printf("Rate limited, retry after %ds\n", rateLimit.RetryAfter)
    default:
        log.Fatal(err)
    }
}
```

## API Reference

### Services

| Service | Methods |
|---------|---------|
| `client.Auth` | `RequestAccess()`, `VerifyAccess()` |
| `client.Config` | `Get()`, `Jurisdictions()` |
| `client.Verification` | `Verify()`, `VerifyAndWait()`, `GetStatus()` |
| `client.Entities` | `Get()`, `History()` |
| `client.Search` | `Find()`, `FindAll()` |
| `client.Account` | `Get()`, `Usage()`, `DataExport()`, `UpdateEmail()`, `CreateKey()`, `RevokeKey()` |
| `client.Billing` | `Get()`, `Purchase()` |
| `client.Checker` | `Check()` |

### Client Options

```go
client := bv.New(
    bv.WithAPIKey("bv_live_..."),          // API key authentication
    bv.WithBaseURL("https://..."),         // Custom base URL
    bv.WithMaxRetries(2),                  // Retry on 5xx (default: 2)
    bv.WithTimeout(30 * time.Second),      // Request timeout (default: 30s)
    bv.WithHTTPClient(customHTTPClient),   // Custom http.Client
)
```

## Requirements

- Go >= 1.21
- Zero external dependencies (stdlib only)

## License

MIT
