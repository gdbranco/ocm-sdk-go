package testing

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/onsi/gomega/ghttp"
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	. "github.com/onsi/ginkgo/v2/dsl/core" // nolint
	. "github.com/onsi/gomega"             // nolint
)

var _ = Describe("RespondWithOCMObjectMarshal", func() {
	var server *ghttp.Server

	BeforeEach(func() {
		server = MakeTCPServer()
	})

	AfterEach(func() {
		server.Close()
	})

	It("should respond with JSON encoded object using provided marshal function", func() {
		// Create a mock cluster object (we'll use a simple struct for testing)
		cluster, err := v1.NewCluster().
			ID("test-cluster").
			Name("Test Cluster").
			Build()
		Expect(err).ToNot(HaveOccurred())

		// Configure the server to respond with the cluster
		server.AppendHandlers(
			RespondWithOCMObjectMarshal(200, cluster, v1.MarshalCluster),
		)

		// Make a request
		resp, err := http.Get(server.URL())
		Expect(err).ToNot(HaveOccurred())
		defer resp.Body.Close()

		// Verify response
		Expect(resp.StatusCode).To(Equal(200))
		Expect(resp.Header.Get("Content-Type")).To(Equal("application/json"))

		// Read and verify the response body contains JSON
		body, err := io.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(body)).To(ContainSubstring("test-cluster"))
		Expect(string(body)).To(ContainSubstring("Test Cluster"))
	})

	It("should handle nil objects gracefully", func() {
		server.AppendHandlers(
			RespondWithOCMObjectMarshal(200, nil, v1.MarshalCluster),
		)

		resp, err := http.Get(server.URL())
		Expect(err).ToNot(HaveOccurred())
		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(500))
	})

	It("should handle invalid marshal function", func() {
		cluster := &v1.Cluster{}
		server.AppendHandlers(
			RespondWithOCMObjectMarshal(200, cluster, "not-a-function"),
		)

		resp, err := http.Get(server.URL())
		Expect(err).ToNot(HaveOccurred())
		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(500))
	})
})

// TestRespondWithOCMObjectMarshal provides unit tests for the marshal function
func TestRespondWithOCMObjectMarshal(t *testing.T) {
	// Create a test cluster
	cluster, err := v1.NewCluster().
		ID("test-cluster-123").
		Name("Test Cluster for Unit Test").
		Build()
	if err != nil {
		t.Fatalf("Failed to create test cluster: %v", err)
	}

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Create the handler
	handler := RespondWithOCMObjectMarshal(201, cluster, v1.MarshalCluster)

	// Call the handler
	handler(rr, req)

	// Check the status code
	if status := rr.Code; status != 201 {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, 201)
	}

	// Check the content type
	expected := "application/json"
	if ct := rr.Header().Get("Content-Type"); ct != expected {
		t.Errorf("Handler returned wrong content type: got %v want %v", ct, expected)
	}

	// Check that the response body contains expected data
	body := rr.Body.String()
	if !strings.Contains(body, "test-cluster-123") {
		t.Errorf("Response body should contain cluster ID: %v", body)
	}
	if !strings.Contains(body, "Test Cluster for Unit Test") {
		t.Errorf("Response body should contain cluster name: %v", body)
	}
}

// TestMarshalFunction tests that we can call marshal functions directly
func TestMarshalFunction(t *testing.T) {
	// Create a test cluster
	cluster, err := v1.NewCluster().
		ID("direct-test").
		Name("Direct Marshal Test").
		Build()
	if err != nil {
		t.Fatalf("Failed to create test cluster: %v", err)
	}

	// Test marshaling directly
	var buf bytes.Buffer
	err = v1.MarshalCluster(cluster, &buf)
	if err != nil {
		t.Fatalf("Failed to marshal cluster: %v", err)
	}

	// Verify the output
	result := buf.String()
	if !strings.Contains(result, "direct-test") {
		t.Errorf("Marshaled result should contain cluster ID: %v", result)
	}
	if !strings.Contains(result, "Direct Marshal Test") {
		t.Errorf("Marshaled result should contain cluster name: %v", result)
	}
}
