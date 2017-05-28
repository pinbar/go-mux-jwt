package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"net/url"

	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	responseStatus, responseBody := testHelperProcessNonAuthRequest("/", "", t)
	assert.Equal(t, http.StatusOK, responseStatus, "should return status OK")
	assert.Equal(t, "Welcome to the Matrix!", responseBody, "wrong message")
}
func TestMetacortex(t *testing.T) {
	responseStatus, responseBody := testHelperProcessNonAuthRequest("/metacortex", "", t)
	assert.Equal(t, http.StatusOK, responseStatus, "should return status OK")
	assert.Equal(t, "Mr Anderson's not so secure workplace!", responseBody, "wrong message")
}
func TestAgentName(t *testing.T) {
	responseStatus, responseBody := testHelperProcessNonAuthRequest("/agents/smith", "", t)
	assert.Equal(t, http.StatusOK, responseStatus, "should return status OK")
	assert.Equal(t, "My name is agent smith", string(responseBody), "wrong message")
}
func TestAuthenticateBadCreds(t *testing.T) {
	responseStatus, responseBody := testHelperAuthenticate("neo", "reeves", t)
	assert.Equal(t, http.StatusUnauthorized, responseStatus, "should return status 401")
	assert.Equal(t, "Name and password do not match", string(responseBody), "wrong password")
}
func TestAuthenticateGoodCreds(t *testing.T) {
	responseStatus, responseBody := testHelperAuthenticate("neo", "keanu", t)
	assert.Equal(t, http.StatusOK, responseStatus, "should return status 200")
	assert.Equal(t, "Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibmVvIiwicm9sZSI6InJlZHBpbGwifQ.TS72DgJx5euYy-YXVEPoHt9Pl0Y7YpV4tRecQaxv7Xk", string(responseBody), "token not created")
}
func TestMegacityNoAuth(t *testing.T) {
	responseStatus, responseBody := testHelperProcessNonAuthRequest("/api/megacity", "", t)
	assert.Equal(t, http.StatusUnauthorized, responseStatus, "should return status 401")
	assert.Equal(t, "Missing Authorization Header", string(responseBody), "no token should return unauthorized error")
}
func TestMegacityWithAuth(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibmVvIiwicm9sZSI6InJlZHBpbGwifQ.TS72DgJx5euYy-YXVEPoHt9Pl0Y7YpV4tRecQaxv7Xk"
	responseStatus, responseBody := testHelperProcessNonAuthRequest("/api/megacity", token, t)
	assert.Equal(t, http.StatusOK, responseStatus, "should return status 200")
	assert.Equal(t, "Welcome to the Megacity!", string(responseBody), "request with token failed")
}
func TestMegacityWithAuthBadToken(t *testing.T) {
	token := "some.bad.token"
	responseStatus, responseBody := testHelperProcessNonAuthRequest("/api/megacity", token, t)
	assert.Equal(t, http.StatusUnauthorized, responseStatus, "should return status 401")
	assert.Equal(t, "Error verifying JWT token: invalid character '²' looking for beginning of value", string(responseBody), "request with bad token should return unauthorized error")
}

//test helper functions
func testHelperProcessNonAuthRequest(reqURI string, jwtToken string, t *testing.T) (int, string) {
	r := ConfigureRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	reqURL := ts.URL + reqURI
	req, _ := http.NewRequest("GET", reqURL, nil)
	if len(jwtToken) > 0 {

		req.Header.Set("Authorization", "Bearer "+jwtToken)
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		t.Error("Error making the request")
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return resp.StatusCode, string(responseBody)
}

func testHelperAuthenticate(programName string, programPassword string, t *testing.T) (int, string) {
	r := ConfigureRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	reqURL := ts.URL + "/authenticate"

	var data = make(url.Values)
	data.Set("programName", programName)
	data.Set("programPassword", programPassword)

	resp, err := http.PostForm(reqURL, data)
	if err != nil {
		t.Error("Error making the request")
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return resp.StatusCode, string(responseBody)
}
