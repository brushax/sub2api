package service

import (
	"net/url"
	"testing"
)

func TestProxyURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		proxy Proxy
		want  string
	}{
		{
			name: "without auth",
			proxy: Proxy{
				Protocol: "http",
				Host:     "proxy.example.com",
				Port:     8080,
			},
			want: "http://proxy.example.com:8080",
		},
		{
			name: "with auth",
			proxy: Proxy{
				Protocol: "socks5",
				Host:     "socks.example.com",
				Port:     1080,
				Username: "user",
				Password: "pass",
			},
			want: "socks5://user:pass@socks.example.com:1080",
		},
		{
			name: "username only keeps no auth for compatibility",
			proxy: Proxy{
				Protocol: "http",
				Host:     "proxy.example.com",
				Port:     8080,
				Username: "user-only",
			},
			want: "http://proxy.example.com:8080",
		},
		{
			name: "with special characters in credentials",
			proxy: Proxy{
				Protocol: "http",
				Host:     "proxy.example.com",
				Port:     3128,
				Username: "first last@corp",
				Password: "p@ ss:#word",
			},
			want: "http://first%20last%40corp:p%40%20ss%3A%23word@proxy.example.com:3128",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.proxy.URL(); got != tc.want {
				t.Fatalf("Proxy.URL() mismatch: got=%q want=%q", got, tc.want)
			}
		})
	}
}

func TestProxyURL_SpecialCharactersRoundTrip(t *testing.T) {
	t.Parallel()

	proxy := Proxy{
		Protocol: "http",
		Host:     "proxy.example.com",
		Port:     3128,
		Username: "first last@corp",
		Password: "p@ ss:#word",
	}

	parsed, err := url.Parse(proxy.URL())
	if err != nil {
		t.Fatalf("parse proxy URL failed: %v", err)
	}
	if got := parsed.User.Username(); got != proxy.Username {
		t.Fatalf("username mismatch after parse: got=%q want=%q", got, proxy.Username)
	}
	pass, ok := parsed.User.Password()
	if !ok {
		t.Fatal("password missing after parse")
	}
	if pass != proxy.Password {
		t.Fatalf("password mismatch after parse: got=%q want=%q", pass, proxy.Password)
	}
}

func TestProxyURLForAccount_TemplatePlaceholders(t *testing.T) {
	t.Parallel()

	proxy := Proxy{
		Protocol: "http",
		Host:     "proxy.example.com",
		Port:     3128,
		Username: "Codex.{RESIN_ACCOUNT}",
		Password: "token-{ACCOUNT_ID}",
	}
	account := &Account{
		ID:       42,
		Name:     "openai-oauth-1",
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
		Credentials: map[string]any{
			"chatgpt_account_id": "acct-123",
		},
	}

	got := proxy.URLForAccount(account)
	want := "http://Codex.acct-123:token-42@proxy.example.com:3128"
	if got != want {
		t.Fatalf("Proxy.URLForAccount() mismatch: got=%q want=%q", got, want)
	}
}

func TestProxyURLForAccount_ProjectIDFallback(t *testing.T) {
	t.Parallel()

	proxy := Proxy{
		Protocol: "http",
		Host:     "proxy.example.com",
		Port:     3128,
		Username: "Gemini.{RESIN_ACCOUNT}",
		Password: "token",
	}
	account := &Account{
		ID:       9,
		Platform: PlatformGemini,
		Type:     AccountTypeOAuth,
		Credentials: map[string]any{
			"project_id": "proj-abc",
		},
	}

	got := proxy.URLForAccount(account)
	want := "http://Gemini.proj-abc:token@proxy.example.com:3128"
	if got != want {
		t.Fatalf("Proxy.URLForAccount() mismatch: got=%q want=%q", got, want)
	}
}

func TestAccountGetProxyURL_FallbacksToAccountID(t *testing.T) {
	t.Parallel()

	account := &Account{
		ID:       7,
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
		Proxy: &Proxy{
			Protocol: "http",
			Host:     "proxy.example.com",
			Port:     8080,
			Username: "Claude.{RESIN_ACCOUNT}",
			Password: "secret",
		},
	}

	got := account.GetProxyURL()
	want := "http://Claude.7:secret@proxy.example.com:8080"
	if got != want {
		t.Fatalf("Account.GetProxyURL() mismatch: got=%q want=%q", got, want)
	}
}
