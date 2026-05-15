package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yourusername/driftwatch/internal/config"
)

// checkOpenSearchDomainStatus checks that an AWS OpenSearch domain has the expected status.
// Required fields: endpoint, domain_name, expected (default: "Active")
func checkOpenSearchDomainStatus(c config.Check) (bool, string, error) {
	endpoint, ok := c.Params["endpoint"].(string)
	if !ok || endpoint == "" {
		return false, "", fmt.Errorf("opensearch_domain_status: missing or empty 'endpoint'")
	}

	domainName, ok := c.Params["domain_name"].(string)
	if !ok || domainName == "" {
		return false, "", fmt.Errorf("opensearch_domain_status: missing or empty 'domain_name'")
	}

	expected, _ := c.Params["expected"].(string)
	if expected == "" {
		expected = "Active"
	}

	url := fmt.Sprintf("%s/2021-01-01/opensearch/domain/%s", endpoint, domainName)
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return false, "", fmt.Errorf("opensearch_domain_status: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("opensearch_domain_status: failed to read response: %w", err)
	}

	var result struct {
		DomainStatus struct {
			Processing bool   `json:"Processing"`
			UpgradeProcessing bool `json:"UpgradeProcessing"`
			Endpoint   string `json:"Endpoint"`
			ClusterConfig struct {
				InstanceType string `json:"InstanceType"`
			} `json:"ClusterConfig"`
			Status string `json:"Status"`
		} `json:"DomainStatus"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("opensearch_domain_status: failed to parse response: %w", err)
	}

	actual := result.DomainStatus.Status
	if actual == "" {
		// Derive status from Processing flag when Status field absent
		if result.DomainStatus.Processing || result.DomainStatus.UpgradeProcessing {
			actual = "Processing"
		} else {
			actual = "Active"
		}
	}

	if actual != expected {
		return true, fmt.Sprintf("domain %q status is %q, expected %q", domainName, actual, expected), nil
	}
	return false, "", nil
}
