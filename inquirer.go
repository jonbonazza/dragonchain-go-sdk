package dragonchain

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"net/http"
	"strings"
	"time"
)

var (
	defaultSigningAlgo = sha256.New()
	defaultEndpoint    = "https://api.dragonchain.com"
)

// Credentials are the credentials used for authenticating requests to the DragonChain API.
type Credentials struct {
	// DragonChainID is the unique ID of the DragonChain that these credentials are for.
	DragonChainID string
	// APIKey the secret API key for to use with the credentials. This is used along with
	// the Client ID to authenticate requests.
	APIKey string
	// ClientID is the unique client id for the credentials. This is used along with APIKey
	// to authenticate requests.
	ClientID string
	// SigningAlogorithm is the hash function used in HMAC signatures.
	SigningAlgorithm hash.Hash
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type inquirer struct {
	httpClient

	Endpoint    string
	Credentials Credentials
	VerifySSL   bool
}

func (iq *inquirer) hmacSign(message string) []byte {
	signingAlgo := iq.Credentials.SigningAlgorithm
	if signingAlgo == nil {
		signingAlgo = defaultSigningAlgo
	}
	h := hmac.New(func() hash.Hash { return signingAlgo }, []byte(iq.Credentials.APIKey))
	h.Write([]byte(message))
	return h.Sum(nil)
}

func (iq *inquirer) Get(resource string, out interface{}) error {
	return iq.doRequest(http.MethodGet, resource, "application/json", nil, out)
}

func (iq *inquirer) Post(resource string, body []byte, out interface{}) error {
	return iq.doRequest(http.MethodPost, resource, "application/json", body, out)
}

func (iq *inquirer) Put(resource string, body []byte, out interface{}) error {
	return iq.doRequest(http.MethodPut, resource, "application/json", body, out)
}

func (iq *inquirer) Delete(resource string) (int, error) {
	return 0, iq.doRequest(http.MethodDelete, resource, "", nil, nil)
}

func (iq *inquirer) doRequest(method, resource, contentType string, body []byte, out interface{}) error {
	if iq.Endpoint == "" {
		iq.Endpoint = defaultEndpoint
	}
	url := iq.Endpoint + resource
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request object: %s", err)
	}
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.999999") + "Z"
	auth := iq.authorizationHeader(method, resource, contentType, timestamp, body)
	req.Header.Set("dragonchain", iq.Credentials.DragonChainID)
	req.Header.Set("timestamp", timestamp)
	req.Header.Set("Authorization", auth)
	if contentType != "" {
		req.Header.Set("Content-type", "application/json")
	}
	resp, err := iq.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %s", err)
	}
	defer resp.Body.Close()
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (iq *inquirer) authorizationHeader(method, resource, contentType string, timestamp string, content []byte) string {
	if content == nil {
		content = []byte("")
	}
	sha := sha256.New()
	sha.Write(content)
	hashedContent := sha.Sum(nil)
	b64Content := base64.StdEncoding.EncodeToString(hashedContent)
	message := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s",
		strings.ToUpper(method), resource,
		iq.Credentials.DragonChainID,
		timestamp,
		contentType,
		b64Content,
	)
	h := hmac.New(func() hash.Hash { return iq.Credentials.SigningAlgorithm }, []byte(iq.Credentials.APIKey))
	h.Write([]byte(message))
	return fmt.Sprintf("DC1-HMAC-SHA256 %s:%s", iq.Credentials.ClientID, string(h.Sum(nil)))
}
