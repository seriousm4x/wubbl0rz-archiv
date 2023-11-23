package routes

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/seriousm4x/wubbl0rz-archiv/internal/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOAuthConfig *oauth2.Config
	stateStorage      *StateStorage
	App               *pocketbase.PocketBase
)

func YoutubeRegisterHandler(app *pocketbase.PocketBase) {
	App = app
	stateStorage = NewStateStorage()
}

func YoutubeHandleVerify(c echo.Context) error {
	settings, err := App.Dao().FindFirstRecordByFilter("settings", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return apis.NewApiError(http.StatusInternalServerError, "failed to get settings", nil)
	}

	ytBearerToken := settings.GetString("yt_bearer_token")
	if ytBearerToken == "" || ytBearerToken == "\"\"" {
		logger.Warning.Println("yt_bearer_token is empty")
		return apis.NewApiError(http.StatusInternalServerError, "yt_bearer_token is empty", nil)
	}

	tok := &oauth2.Token{}
	if err := json.Unmarshal([]byte(ytBearerToken), &tok); err != nil {
		logger.Error.Println(err)
		return apis.NewApiError(http.StatusInternalServerError, "failed to unmarshal bearer token", nil)
	}

	return c.NoContent(http.StatusOK)
}

func YoutubeHandleLogin(c echo.Context) error {
	settings, err := App.Dao().FindFirstRecordByFilter("settings", "id != ''")
	if err != nil {
		logger.Error.Println(err)
		return apis.NewApiError(http.StatusInternalServerError, "failed to get settings", nil)
	}

	ytClientSecret := settings.GetString("yt_client_secret")
	if ytClientSecret == "" {
		logger.Error.Println("yt_client_secret is empty")
		return apis.NewApiError(http.StatusInternalServerError, "yt_client_secret is empty", nil)
	}

	scope := "https://www.googleapis.com/auth/youtube.upload"
	googleOAuthConfig, err = google.ConfigFromJSON([]byte(ytClientSecret), scope)
	if err != nil {
		logger.Error.Println(err)
		return apis.NewApiError(http.StatusInternalServerError, "failed to create config from json", nil)
	}
	googleOAuthConfig.RedirectURL = fmt.Sprintf("%s/wubbl0rz/youtube/callback", os.Getenv("PUBLIC_API_URL"))

	state, err := generateRandomState()
	if err != nil {
		return apis.NewApiError(http.StatusInternalServerError, "failed to create random state", nil)
	}

	codeVerifier, codeChallenge, err := generatePKCEChallenge()
	if err != nil {
		return apis.NewApiError(http.StatusInternalServerError, "failed to create pkce challenge", nil)
	}

	authURL := googleOAuthConfig.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.AccessTypeOffline)
	logger.Debug.Println("authURL:", authURL)

	// Store the state and code verifier for later verification
	stateStorage.Store(state, codeVerifier)

	return c.String(http.StatusOK, authURL)
}

func YoutubeHandleCallback(c echo.Context) error {
	state := c.FormValue("state")
	codeVerifier, err := stateStorage.GetCodeVerifier(state)
	if err != nil {
		return apis.NewApiError(http.StatusInternalServerError, "failed to get code verifier", nil)
	}

	// Verify that the state parameter is valid
	if !stateStorage.Verify(state) {
		return apis.NewUnauthorizedError("invalid state parameter", nil)
	}

	code := c.FormValue("code")

	token, err := googleOAuthConfig.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		return apis.NewApiError(http.StatusInternalServerError, "failed to exchange token", nil)
	}

	settings, err := App.Dao().FindFirstRecordByFilter("settings", "id != ''")
	if err != nil {
		return apis.NewApiError(http.StatusInternalServerError, "failed to find settings record", nil)
	}

	settings.Set("yt_bearer_token", token)
	if err := App.Dao().SaveRecord(settings); err != nil {
		return apis.NewApiError(http.StatusInternalServerError, "unable to save record", nil)
	}

	return c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s/admin", os.Getenv("PUBLIC_FRONTEND_URL")))
}

// StateStorage is a simple in-memory storage for state parameters and code verifiers
type StateStorage struct {
	mu            sync.Mutex
	states        map[string]time.Time
	codeVerifiers map[string]string
}

// NewStateStorage creates a new StateStorage instance
func NewStateStorage() *StateStorage {
	return &StateStorage{
		states:        make(map[string]time.Time),
		codeVerifiers: make(map[string]string),
	}
}

// Store stores a state parameter and its creation time
func (s *StateStorage) Store(state, codeVerifier string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[state] = time.Now()
	s.codeVerifiers[state] = codeVerifier
}

// Verify verifies that the state parameter is valid
func (s *StateStorage) Verify(state string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the state parameter exists and is not expired (e.g., valid for 5 minutes)
	creationTime, ok := s.states[state]
	if !ok || time.Since(creationTime) > 5*time.Minute {
		return false
	}

	// Remove the state parameter to prevent replay attacks
	delete(s.states, state)
	return true
}

// GetCodeVerifier retrieves the stored code verifier for a given state
func (s *StateStorage) GetCodeVerifier(state string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	codeVerifier, ok := s.codeVerifiers[state]
	if !ok {
		err := errors.New("verifier not found")
		logger.Error.Println(err)
		return "", err
	}

	// Remove the code verifier to prevent replay attacks
	delete(s.codeVerifiers, state)
	return codeVerifier, nil
}

// Generate a random string (e.g., a UUID)
func generateRandomState() (string, error) {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		logger.Error.Println(err)
		return "", err
	}

	return fmt.Sprintf("%x", buf), nil
}

// Generates a random PKCE string using crypto/rand
func generatePKCEChallenge() (string, string, error) {
	verifier := make([]byte, 32)
	_, err := rand.Read(verifier)
	if err != nil {
		logger.Error.Println(err)
		return "", "", err
	}

	codeVerifier := base64.RawURLEncoding.EncodeToString(verifier)

	sha256Verifier := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(sha256Verifier[:])

	return codeVerifier, codeChallenge, nil
}
