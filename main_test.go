package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	responseStatus, responseBody := testHelperProcessRequest("/", "", t)
	assert.Equal(t, http.StatusOK, responseStatus, "should return status OK")
	assert.Equal(t, "Welcome to the Matrix!", responseBody, "wrong message")
}
func TestMetacortex(t *testing.T) {
	responseStatus, responseBody := testHelperProcessRequest("/metacortex", "", t)
	assert.Equal(t, http.StatusOK, responseStatus, "should return status OK")
	assert.Equal(t, "Mr Anderson's not so secure workplace!", responseBody, "wrong message")
}
func TestAgentName(t *testing.T) {
	responseStatus, responseBody := testHelperProcessRequest("/agents/smith", "", t)
	assert.Equal(t, http.StatusOK, responseStatus, "should return status OK")
	assert.Equal(t, "My name is agent smith", string(responseBody), "wrong message")
}
func TestAuthenticateBadCreds(t *testing.T) {
	responseStatus, responseBody := testHelperAuthenticate("neo", "lawrence", t)
	assert.Equal(t, http.StatusUnauthorized, responseStatus, "should return status 401")
	assert.Equal(t, "Name and password do not match", string(responseBody), "wrong password")
}
func TestAuthenticateNeoCreds(t *testing.T) {
	responseStatus, responseBody := testHelperAuthenticate("neo", "keanu", t)
	assert.Equal(t, http.StatusOK, responseStatus, "should return status 200")
	assert.Equal(t, "Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibmVvIiwicm9sZSI6InJlZHBpbGwifQ.TS72DgJx5euYy-YXVEPoHt9Pl0Y7YpV4tRecQaxv7Xk", string(responseBody), "token not created")
}
func TestAuthenticateMorpheusCreds(t *testing.T) {
	responseStatus, responseBody := testHelperAuthenticate("morpheus", "lawrence", t)
	assert.Equal(t, http.StatusOK, responseStatus, "should return status 200")
	assert.Equal(t, "Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibW9ycGhldXMiLCJyb2xlIjoicmVkcGlsbCJ9.zNV3twzLyIv0GeIx4BhcVM0dtlC73izClB0XIlZoxz4", string(responseBody), "token not created")
}
func TestMegacityNoAuth(t *testing.T) {
	responseStatus, responseBody := testHelperProcessRequest("/api/megacity", "", t)
	assert.Equal(t, http.StatusUnauthorized, responseStatus, "should return status 401")
	assert.Equal(t, "Missing Authorization Header", string(responseBody), "no token should return unauthorized error")
}
func TestMegacityWithAuth(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibmVvIiwicm9sZSI6InJlZHBpbGwifQ.TS72DgJx5euYy-YXVEPoHt9Pl0Y7YpV4tRecQaxv7Xk"
	responseStatus, responseBody := testHelperProcessRequest("/api/megacity", token, t)
	assert.Equal(t, http.StatusOK, responseStatus, "should return status 200")
	assert.Equal(t, "Welcome to the Megacity!", string(responseBody), "request with token failed")
}
func TestMegacityWithAuthBadToken(t *testing.T) {
	token := "some.bad.token"
	responseStatus, responseBody := testHelperProcessRequest("/api/megacity", token, t)
	assert.Equal(t, http.StatusUnauthorized, responseStatus, "should return status 401")
	assert.Equal(t, "Error verifying JWT token: invalid character '²' looking for beginning of value", string(responseBody), "request with bad token should return unauthorized error")
}
func TestLeVraiWithNeoAuth(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibmVvIiwicm9sZSI6InJlZHBpbGwifQ.TS72DgJx5euYy-YXVEPoHt9Pl0Y7YpV4tRecQaxv7Xk"
	responseStatus, responseBody := testHelperProcessRequest("/api/levrai", token, t)
	assert.Equal(t, http.StatusOK, responseStatus, "should return status 200")
	assert.Equal(t, "Welcome to the LeVrai!", string(responseBody), "request should have neo token")
}
func TestLeVraiWithMorpheusAuth(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibW9ycGhldXMiLCJyb2xlIjoicmVkcGlsbCJ9.zNV3twzLyIv0GeIx4BhcVM0dtlC73izClB0XIlZoxz4"
	responseStatus, responseBody := testHelperProcessRequest("/api/levrai", token, t)
	assert.Equal(t, http.StatusUnauthorized, responseStatus, "should return status 401")
	assert.Equal(t, "Only Neo can enter the Merovingian's restaurant!", string(responseBody), "request should have neo token")
}

//test helper functions
func testHelperProcessRequest(reqURI string, jwtToken string, t *testing.T) (int, string) {
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
