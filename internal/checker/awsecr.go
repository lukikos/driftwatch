package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/driftwatch/driftwatch/internal/config"
)

// checkECRRepositoryPolicy checks whether an ECR repository exists and its
// image tag mutability matches the expected value.
//
// Required fields:
//   - endpoint:    base URL of the ECR-compatible API (e.g. for localstack)
//   - repository:  ECR repository name
//   - expected:    expected imageTagMutability value (MUTABLE or IMMUTABLE)
func checkECRRepositoryPolicy(c config.Check) (bool, string, error) {
	endpoint, ok := c.Params["endpoint"].(string)
	if !ok || strings.TrimSpace(endpoint) == "" {
		return false, "", fmt.Errorf("ecr_repository_policy: missing or empty 'endpoint'")
	}

	repo, ok := c.Params["repository"].(string)
	if !ok || strings.TrimSpace(repo) == "" {
		return false, "", fmt.Errorf("ecr_repository_policy: missing or empty 'repository'")
	}

	expected, ok := c.Params["expected"].(string)
	if !ok || strings.TrimSpace(expected) == "" {
		return false, "", fmt.Errorf("ecr_repository_policy: missing or empty 'expected'")
	}

	url := fmt.Sprintf("%s/repositories/%s", strings.TrimRight(endpoint, "/"), repo)
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return false, "", fmt.Errorf("ecr_repository_policy: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return true, fmt.Sprintf("repository %q not found", repo), nil
	}
	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("ecr_repository_policy: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("ecr_repository_policy: reading body: %w", err)
	}

	var result struct {
		ImageTagMutability string `json:"imageTagMutability"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", fmt.Errorf("ecr_repository_policy: parsing response: %w", err)
	}

	if !strings.EqualFold(result.ImageTagMutability, expected) {
		return true, fmt.Sprintf("imageTagMutability is %q, expected %q", result.ImageTagMutability, expected), nil
	}

	return false, "", nil
}
