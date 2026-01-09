package chirp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFilterWords(t *testing.T) {
	input := []string{"hello", "kerfuffle", "world", "fornax"}
	got := filterWords(input)
	want := []string{"hello", cleanedString, "world", cleanedString}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("filterWords[%d] = %q; want %q", i, got[i], want[i])
		}
	}
}

func TestValidateChirpHandler_Success(t *testing.T) {
	body := `{"body":"Hello kerfuffle WORLD"}`
	req := httptest.NewRequest("POST", "/validate", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()
	ValidateChirpHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rr.Code, http.StatusOK)
	}
	var resp ValidateChirpResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	want := "hello **** world"
	if resp.CleanedBody != want {
		t.Fatalf("cleaned = %q; want %q", resp.CleanedBody, want)
	}
}

func TestValidateChirpHandler_TooLong(t *testing.T) {
	longBody := make([]byte, 141)
	for i := range longBody {
		longBody[i] = 'a'
	}
	payload := `{"body":"` + string(longBody) + `"}`
	req := httptest.NewRequest("POST", "/validate", bytes.NewBufferString(payload))
	rr := httptest.NewRecorder()
	ValidateChirpHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d; want %d", rr.Code, http.StatusBadRequest)
	}
	var resp ValidateChirpError
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Error != "Chirp is too long" {
		t.Fatalf("error = %q; want %q", resp.Error, "Chirp is too long")
	}
}

func TestValidateChirpHandler_BadJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/validate", bytes.NewBufferString("{bad json"))
	rr := httptest.NewRecorder()
	ValidateChirpHandler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d; want %d", rr.Code, http.StatusInternalServerError)
	}
	var resp ValidateChirpError
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Error != "Something went wrong" {
		t.Fatalf("error = %q; want %q", resp.Error, "Something went wrong")
	}
}
