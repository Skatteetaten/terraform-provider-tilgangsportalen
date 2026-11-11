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
	msg609 = "609: Object already exists"
)

// newErrorClient returns a Client pointing to a local test server that always
// responds with the given API error message.
func newErrorClient(t *testing.T, message string) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `[{"Message": %q}]`, message) // write the error message as an API-style JSON response body
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

		{name: "group", check: group, apiMessage: msg605},                // Group deleted outside Terraform or not visible to us: not an error
		{name: "group", check: group, apiMessage: msg606},                // Group deleted outside Terraform or not visible to us: not an error
		{name: "group", check: group, apiMessage: msg602, wantErr: true}, // Real failures must not be mistaken for "does not exist"
		{name: "role", check: role, apiMessage: msg605},                  // Role deleted outside Terraform or not visible to us: not an error
		{name: "role", check: role, apiMessage: msg602, wantErr: true},   // Real failures must not be mistaken for "does not exist"
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

// TestTermsOfUseErrorCodes checks that error codes meaning "already in the
// desired state" give success, and that any other code gives an error.
// GetTermsOfUseForRole wraps the API error, so this also checks that the code
// is still found after wrapping.
//   - AssignTermsOfUseToRole: 609 means the terms of use is already assigned.
//   - WaitForTermsOfUseRemoval: 606 means no terms of use is assigned.
//   - WaitForTermsOfUseAssignment: 606 means keep waiting; any other code must fail at once.
func TestTermsOfUseErrorCodes(t *testing.T) {
	assign := func(c *Client) error {
		_, err := c.AssignTermsOfUseToRole(TermsOfUseAssignment{RoleName: "my-role", TermsOfUse: "my-terms"})
		return err
	}
	waitForRemoval := func(c *Client) error {
		return c.WaitForTermsOfUseRemoval("my-role")
	}
	waitForAssignment := func(c *Client) error {
		_, err := c.WaitForTermsOfUseAssignment("my-terms", "my-role")
		return err
	}

	tests := []struct {
		name       string
		call       func(*Client) error
		apiMessage string
		wantErr    bool
	}{
		{name: "assign", call: assign, apiMessage: msg609},                                      // Already assigned is the desired state, so apply must not fail
		{name: "assign", call: assign, apiMessage: msg602, wantErr: true},                       // Other errors must still fail
		{name: "waitForRemoval", call: waitForRemoval, apiMessage: msg606},                      // No terms of use assigned means removal is done; the error is wrapped on the way
		{name: "waitForRemoval", call: waitForRemoval, apiMessage: msg602, wantErr: true},       // Other errors must fail at once instead of waiting until timeout
		{name: "waitForRemoval", call: waitForRemoval, apiMessage: msg605, wantErr: true},       // Role not found or no permission must not be mistaken for "no terms of use assigned"
		{name: "waitForAssignment", call: waitForAssignment, apiMessage: msg602, wantErr: true}, // Other errors must fail at once instead of waiting until timeout
		{name: "waitForAssignment", call: waitForAssignment, apiMessage: msg605, wantErr: true}, // Role not found or no permission must not be mistaken for "not assigned yet"
	}

	for _, tt := range tests {
		t.Run(tt.name+"/"+tt.apiMessage, func(t *testing.T) {
			err := tt.call(newErrorClient(t, tt.apiMessage))
			if (err != nil) != tt.wantErr {
				t.Errorf("err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestErrorCodeFromMessage checks that the error code is read from the start of Message
// only, so object names containing "6xx:" do not give the wrong code.
func TestErrorCodeFromMessage(t *testing.T) {
	tests := []struct {
		message  string
		wantCode string
	}{
		{message: msg605, wantCode: "605"},                                 // Normal format from the API
		{message: "609 : Object already exists", wantCode: "609"},          // Space before the colon
		{message: "609: Role '606: test' already exists", wantCode: "609"}, // Object name containing another code must not affect the result
		{message: "Role '606: test' not found"},                            // Message without a code at the start must not pick up a code from the object name
		{message: "Something failed"},                                      // Message without any colon
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			if got := newAPIError("500", 500, tt.message).Code; got != tt.wantCode {
				t.Errorf("Code = %q, want %q", got, tt.wantCode)
			}
		})
	}
}

// TestErrorMessageFromBody checks that Message is read from the error body, and
// that an empty list or a body that is not JSON does not fail.
func TestErrorMessageFromBody(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{name: "message", body: `[{"Message": "605: Object not found"}]`, want: "605: Object not found"}, // Normal format from the API
		{name: "empty list", body: `[]`, want: ""},                                                       // Used to panic on index out of range
		{name: "not json", body: "<html>Bad Gateway</html>", want: "<html>Bad Gateway</html>"},           // Error page from a gateway or proxy; used to give a JSON parse error instead of the HTTP status
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseErrorMessage([]byte(tt.body)); got != tt.want {
				t.Errorf("parseErrorMessage(%q) = %q, want %q", tt.body, got, tt.want)
			}
		})
	}
}
