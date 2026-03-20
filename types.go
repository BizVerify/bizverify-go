package bizverify

import (
	"encoding/json"
	"time"
)

type Address struct {
	Line1      string  `json:"line1"`
	Line2      *string `json:"line2"`
	City       string  `json:"city"`
	State      *string `json:"state"`
	PostalCode *string `json:"postal_code"`
	Country    string  `json:"country"`
}

type RegisteredAgent struct {
	Name    string   `json:"name"`
	Address *Address `json:"address"`
}

type Officer struct {
	Name    string   `json:"name"`
	Title   string   `json:"title"`
	Address *Address `json:"address"`
}

type FilingSummary struct {
	Date        string  `json:"date"`
	Type        string  `json:"type"`
	Description *string `json:"description"`
}

type User struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Plan          string `json:"plan"`
	CreditBalance int    `json:"credit_balance"`
	CreatedAt     string `json:"created_at"`
}

type RequestAccessParams struct {
	Email       string `json:"email"`
	AcceptTerms bool   `json:"accept_terms"`
}

type RequestAccessResponse struct {
	Message string `json:"message"`
}

type VerifyAccessParams struct {
	Email string  `json:"email"`
	Code  string  `json:"code"`
	Label *string `json:"label,omitempty"`
}

type VerifyAccessResponse struct {
	APIKey string `json:"api_key"`
	KeyID  string `json:"key_id"`
	Label  string `json:"label"`
}

type ResponseMeta struct {
	CreditsRemaining   *int `json:"credits_remaining"`
	CreditsCharged     *int `json:"credits_charged"`
	RateLimitLimit     *int `json:"rate_limit_limit"`
	RateLimitRemaining *int `json:"rate_limit_remaining"`
	RateLimitReset     *int `json:"rate_limit_reset"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ConfigResponse struct {
	Jurisdictions interface{} `json:"jurisdictions"`
	Checker       interface{} `json:"checker"`
	Pricing       interface{} `json:"pricing"`
	Features      interface{} `json:"features"`
	RateLimits    interface{} `json:"rateLimits"`
	Status        interface{} `json:"status"`
	Legal         interface{} `json:"legal"`
	Docs          interface{} `json:"docs"`
}

type JurisdictionInfo struct {
	Code     string          `json:"code"`
	Name     string          `json:"name"`
	Features map[string]bool `json:"features"`
}

type JurisdictionsResponse struct {
	Jurisdictions []JurisdictionInfo `json:"jurisdictions"`
}

type VerifyParams struct {
	EntityName        string `json:"entity_name"`
	Jurisdiction      string `json:"jurisdiction"`
	EntityType        string `json:"entity_type,omitempty"`
	VerificationLevel string `json:"verification_level,omitempty"`
	ForceRefresh      *bool  `json:"force_refresh,omitempty"`
	WebhookURL        string `json:"webhook_url,omitempty"`
}

type VerifyResponse struct {
	Status         string          `json:"status"`
	Data           json.RawMessage `json:"data,omitempty"`
	JobID          *string         `json:"job_id,omitempty"`
	EntityID       *string         `json:"entity_id,omitempty"`
	Cached         bool            `json:"cached"`
	CreditsCharged int             `json:"credits_charged"`
}

type JobStatusResponse struct {
	ID                string          `json:"id"`
	Status            string          `json:"status"`
	Jurisdiction      string          `json:"jurisdiction"`
	Query             string          `json:"query"`
	VerificationLevel string          `json:"verification_level"`
	CreditsCharged    int             `json:"credits_charged"`
	Result            json.RawMessage `json:"result,omitempty"`
	Error             *string         `json:"error,omitempty"`
	CreatedAt         string          `json:"created_at"`
	CompletedAt       *string         `json:"completed_at"`
}

type PollOptions struct {
	PollInterval   time.Duration
	Timeout        time.Duration
	OnStatusChange func(JobStatusResponse)
}

type Entity struct {
	ID                   string           `json:"id"`
	EntityName           string           `json:"entity_name"`
	Jurisdiction         string           `json:"jurisdiction"`
	EntityType           string           `json:"entity_type"`
	Status               string           `json:"status"`
	JurisdictionID       *string          `json:"jurisdiction_id"`
	GoodStanding         *bool            `json:"good_standing"`
	FormationDate        *string          `json:"formation_date"`
	RegisteredAgent      *RegisteredAgent `json:"registered_agent"`
	Officers             []Officer        `json:"officers"`
	PrincipalAddress     *Address         `json:"principal_address"`
	FilingHistorySummary []FilingSummary   `json:"filing_history_summary"`
	CreatedAt            string           `json:"created_at"`
	UpdatedAt            string           `json:"updated_at"`
}

type HistoryParams struct {
	Limit  *int
	Offset *int
}

type PaginatedSnapshots struct {
	Snapshots []json.RawMessage `json:"snapshots"`
	Total     int               `json:"total"`
	Limit     int               `json:"limit"`
	Offset    int               `json:"offset"`
}

type SearchParams struct {
	EntityName   string `json:"entity_name"`
	Jurisdiction string `json:"jurisdiction,omitempty"`
	EntityType   string `json:"entity_type,omitempty"`
	Limit        *int   `json:"limit,omitempty"`
	Offset       *int   `json:"offset,omitempty"`
}

type SearchResult struct {
	EntityName     string  `json:"entity_name"`
	Jurisdiction   string  `json:"jurisdiction"`
	EntityType     string  `json:"entity_type"`
	Status         string  `json:"status"`
	JurisdictionID *string `json:"jurisdiction_id"`
	Confidence     float64 `json:"confidence"`
}

type SearchResponse struct {
	Results               []SearchResult `json:"results"`
	Total                 int            `json:"total"`
	Limit                 int            `json:"limit"`
	Offset                int            `json:"offset"`
	JurisdictionsSearched []string       `json:"jurisdictions_searched"`
	JurisdictionsFailed   []string       `json:"jurisdictions_failed"`
	CreditsCharged        int            `json:"credits_charged"`
}

type ApiKeyInfo struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Prefix    string `json:"prefix"`
	RateLimit int    `json:"rate_limit"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

type Account struct {
	ID            string       `json:"id"`
	Email         string       `json:"email"`
	EmailVerified bool         `json:"email_verified"`
	Plan          string       `json:"plan"`
	CreditBalance int          `json:"credit_balance"`
	APIKeys       []ApiKeyInfo `json:"api_keys"`
	CreatedAt     string       `json:"created_at"`
}

type UsageEntry struct {
	Date         string  `json:"date"`
	Endpoint     string  `json:"endpoint"`
	Jurisdiction *string `json:"jurisdiction"`
	RequestCount int     `json:"request_count"`
	CreditsUsed  int     `json:"credits_used"`
}

type EndpointSummary struct {
	Endpoint      string `json:"endpoint"`
	TotalRequests int    `json:"total_requests"`
	TotalCredits  int    `json:"total_credits"`
}

type JurisdictionSummary struct {
	Jurisdiction  string `json:"jurisdiction"`
	TotalRequests int    `json:"total_requests"`
	TotalCredits  int    `json:"total_credits"`
}

type UsageStats struct {
	PeriodDays     int                   `json:"period_days"`
	Daily          []UsageEntry          `json:"daily"`
	ByEndpoint     []EndpointSummary     `json:"by_endpoint"`
	ByJurisdiction []JurisdictionSummary `json:"by_jurisdiction"`
}

type DataExportProfile struct {
	ID              string  `json:"id"`
	Email           string  `json:"email"`
	EmailVerified   bool    `json:"email_verified"`
	Plan            string  `json:"plan"`
	CreditBalance   int     `json:"credit_balance"`
	TermsAcceptedAt *string `json:"terms_accepted_at"`
	TermsVersion    *string `json:"terms_version"`
	CreatedAt       string  `json:"created_at"`
}

type DataExportApiKey struct {
	ID        string  `json:"id"`
	Label     string  `json:"label"`
	Prefix    string  `json:"prefix"`
	RateLimit int     `json:"rate_limit"`
	IsActive  bool    `json:"is_active"`
	CreatedAt string  `json:"created_at"`
	RevokedAt *string `json:"revoked_at"`
}

type DataExportTransaction struct {
	ID           string  `json:"id"`
	Amount       int     `json:"amount"`
	BalanceAfter int     `json:"balance_after"`
	Type         string  `json:"type"`
	Description  string  `json:"description"`
	ReferenceID  *string `json:"reference_id"`
	CreatedAt    string  `json:"created_at"`
}

type DataExportJob struct {
	ID             string  `json:"id"`
	Jurisdiction   string  `json:"jurisdiction"`
	Query          string  `json:"query"`
	Status         string  `json:"status"`
	CreditsCharged int     `json:"credits_charged"`
	CreatedAt      string  `json:"created_at"`
	CompletedAt    *string `json:"completed_at"`
}

type DataExport struct {
	Profile            DataExportProfile       `json:"profile"`
	APIKeys            []DataExportApiKey       `json:"api_keys"`
	CreditTransactions []DataExportTransaction  `json:"credit_transactions"`
	VerificationJobs   []DataExportJob          `json:"verification_jobs"`
	UsageStats         []UsageEntry             `json:"usage_stats"`
}

type CreateKeyResponse struct {
	ID      string `json:"id"`
	Key     string `json:"key"`
	Prefix  string `json:"prefix"`
	Label   string `json:"label"`
	Message string `json:"message"`
}

type BillingParams struct {
	Limit  *int
	Offset *int
}

type BillingInfo struct {
	Balance      int             `json:"balance"`
	Packages     json.RawMessage `json:"packages"`
	Transactions json.RawMessage `json:"transactions"`
}

type PurchaseResponse struct {
	SessionID string `json:"session_id"`
	URL       string `json:"url"`
}

type CheckerResult struct {
	EntityName   string  `json:"entity_name"`
	EntityType   string  `json:"entity_type"`
	Status       string  `json:"status"`
	Jurisdiction string  `json:"jurisdiction"`
	Confidence   float64 `json:"confidence"`
}

type CheckerResponse struct {
	Results      []CheckerResult `json:"results"`
	Query        string          `json:"query"`
	Jurisdiction string          `json:"jurisdiction"`
	Total        int             `json:"total"`
}
