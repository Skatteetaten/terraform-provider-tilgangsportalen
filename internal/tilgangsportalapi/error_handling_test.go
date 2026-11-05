// Unit tests for how the API client handles error codes from Tilgangsportalen.
//
// The API responds to errors with HTTP 500 and the error code as a prefix of the
// message, e.g. [{"Message": "605: Object not found or no permission"}]. Codes are
// listed in the API documentation:
// https://wiki.sits.no/spaces/OIM/pages/1054582648/API+dokumentasjon
//
// These tests only cover what the acceptance tests in internal/provider cannot
// trigger, such as objects deleted outside Terraform or missing API permissions.
//
// To extend: add a msgXXX constant for new error codes and a row in the test
// table. New API functions with the same signature can be added to the table
// directly as a method expression, e.g. (*Client).CheckIfGroupExists.
package tilgangsportalapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// API error messages as returned in the Message field, "<code>: <description>".
const (
	msg602 = "602: Missing permissions to use API"
	msg605 = "605: Object not found or no permission"
	msg606 = "606: Object does not exist"
)

// newErrorClient returns a Client pointing to a local test server that always
// responds with the given API error message.
func newErrorClient(t *testing.T, message string) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `[{"Message": %q}]`, message)
	}))
	t.Cleanup(srv.Close)

	return &Client{baseURL: srv.URL}
}

// TestCheckIfExistsErrorCodes checks that "not found" codes give a non-existing
// object without error, and that any other code gives an error.
//   - GetAzureADGroupV1 can return 500, 600, 602, 605 and 606.
//   - GetRole can return 500, 600, 602 and 605.
func TestCheckIfExistsErrorCodes(t *testing.T) {
	group := (*Client).CheckIfGroupExists
	role := (*Client).CheckIfRoleExists

	tests := []struct {
		name       string
		check      func(*Client, string) (bool, string, error)
		apiMessage string
		wantErr    bool
	}{
		{name: "group", check: group, apiMessage: msg605},
		{name: "group", check: group, apiMessage: msg606},
		{name: "group", check: group, apiMessage: msg602, wantErr: true},
		{name: "role", check: role, apiMessage: msg605},
		{name: "role", check: role, apiMessage: msg602, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/"+tt.apiMessage, func(t *testing.T) {
			exists, name, err := tt.check(newErrorClient(t, tt.apiMessage), "my-object")

			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if exists || name != "" {
				t.Errorf("got exists = %v, name = %q, want false and empty name", exists, name)
			}
		})
	}
}
