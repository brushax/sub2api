package service

import (
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Proxy struct {
	ID        int64
	Name      string
	Protocol  string
	Host      string
	Port      int
	Username  string
	Password  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Proxy) IsActive() bool {
	return p.Status == StatusActive
}

func (p *Proxy) URL() string {
	return p.urlWithCredentials(p.Username, p.Password)
}

// URLForAccount renders proxy credentials with account-aware placeholders.
// This keeps existing static proxy behavior unchanged while allowing a single
// proxy definition to expose distinct upstream identities per account.
func (p *Proxy) URLForAccount(account *Account) string {
	return p.urlWithCredentials(
		renderProxyCredentialTemplate(p.Username, account),
		renderProxyCredentialTemplate(p.Password, account),
	)
}

func (p *Proxy) urlWithCredentials(username, password string) string {
	u := &url.URL{
		Scheme: p.Protocol,
		Host:   net.JoinHostPort(p.Host, strconv.Itoa(p.Port)),
	}
	if username != "" && password != "" {
		u.User = url.UserPassword(username, password)
	}
	return u.String()
}

func renderProxyCredentialTemplate(raw string, account *Account) string {
	if raw == "" {
		return ""
	}

	replacements := []string{
		"{ACCOUNT_ID}", "",
		"{ACCOUNT_NAME}", "",
		"{PLATFORM}", "",
		"{TYPE}", "",
		"{RESIN_ACCOUNT}", "",
		"{CHATGPT_ACCOUNT_ID}", "",
		"{PROJECT_ID}", "",
		"{CLAUDE_USER_ID}", "",
	}
	if account != nil {
		replacements = []string{
			"{ACCOUNT_ID}", strconv.FormatInt(account.ID, 10),
			"{ACCOUNT_NAME}", strings.TrimSpace(account.Name),
			"{PLATFORM}", strings.TrimSpace(account.Platform),
			"{TYPE}", strings.TrimSpace(account.Type),
			"{RESIN_ACCOUNT}", account.GetProxyAccountIdentity(),
			"{CHATGPT_ACCOUNT_ID}", strings.TrimSpace(account.GetChatGPTAccountID()),
			"{PROJECT_ID}", strings.TrimSpace(account.GetCredential("project_id")),
			"{CLAUDE_USER_ID}", strings.TrimSpace(account.GetClaudeUserID()),
		}
	}
	return strings.NewReplacer(replacements...).Replace(raw)
}

type ProxyWithAccountCount struct {
	Proxy
	AccountCount   int64
	LatencyMs      *int64
	LatencyStatus  string
	LatencyMessage string
	IPAddress      string
	Country        string
	CountryCode    string
	Region         string
	City           string
	QualityStatus  string
	QualityScore   *int
	QualityGrade   string
	QualitySummary string
	QualityChecked *int64
}

type ProxyAccountSummary struct {
	ID       int64
	Name     string
	Platform string
	Type     string
	Notes    *string
}
