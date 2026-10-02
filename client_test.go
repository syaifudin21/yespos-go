package sdk

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestSignatureMatchesServerAlgorithm(t *testing.T) {
	c := NewClient("http://example.test", "va-demo", "secret-demo")
	body := []byte(`{"phone":"081234567890","code":"kodelarismanis"}`)
	h := sha256.Sum256(body)
	m := hmac.New(sha256.New, []byte("secret-demo"))
	_, _ = m.Write([]byte("POST:va-demo:" + hex.EncodeToString(h[:]) + ":secret-demo"))
	if got, want := c.signature("POST", body), hex.EncodeToString(m.Sum(nil)); got != want {
		t.Fatalf("signature mismatch: got %s, want %s", got, want)
	}
}

func TestGenerateVoucherCode(t *testing.T) {
	client := NewClient("http://example.test", "va-demo", "secret-demo")
	client.HTTPClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/client/voucher/generate-code" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("va") != "va-demo" || r.Header.Get("signature") == "" {
			t.Fatal("missing auth headers")
		}
		return &http.Response{StatusCode: http.StatusCreated, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"message":"voucher_code_generated","data":{"code":"NCOIFX7F","expired_at":"2026-07-10T17:57:02Z"}}`))}, nil
	})}
	result, err := client.GenerateVoucherCode(context.Background(), VoucherCodeRequest{Phone: "081234567890", Code: "kodelarismanis"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Data.Code != "NCOIFX7F" {
		t.Fatalf("unexpected code: %s", result.Data.Code)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
