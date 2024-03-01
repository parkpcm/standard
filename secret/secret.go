package secret

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"fmt"
	"os"
)

// Get retrieves the secret data for the specified path using the Google Secret Manager API.
// It returns the secret data as a byte array and an error if the retrieval fails.
// The function relies on the getSecret function to access the secret version data.
// The getSecret function creates a SecretManager client and sends a request to access the specified secret version.
// If an error occurs during the process, the function returns nil for the secret data and an error describing the issue.
func Get(ctx context.Context, path string) ([]byte, error) {
	return getSecret(ctx, path)
}

func GetFromVolume(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// getSecret retrieves the secret data for the specified path using the Google Secret Manager API.
// It takes a context and a path string as input.
// It returns the secret data as a byte array and an error if the retrieval fails.
// The function creates a SecretManager client using the provided context.
// If an error occurs during the creation of the client, the function returns nil for the secret data and an error describing the issue.
// It then constructs an AccessSecretVersionRequest with the provided path and sends the request to access the specified secret version.
// If an error occurs during the process of accessing the secret version, the function returns nil for the secret data and an error describing the issue.
// Finally, it returns the payload data from the result and any error that occurred during the process.
func getSecret(ctx context.Context, path string) ([]byte, error) {

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: path,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %v", err)
	}

	return result.Payload.Data, err

}
