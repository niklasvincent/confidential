package aws

import (
	"fmt"
	"os"
	"runtime"
	"path/filepath"

	//"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// UserHomeDir returns the home directory for the user the process is running under
func UserHomeDir() string {
	if runtime.GOOS == "windows" { // Windows
		return os.Getenv("USERPROFILE")
	}

	// *nix
	return os.Getenv("HOME")
}

// Get AWS credentials for a specific profile in ~/.aws/credentials
func GetAwsCredentialsForProfile(profile string) (*credentials.Credentials, error) {
	awsCredentialsFilename := filepath.Join(UserHomeDir(), ".aws", "credentials")
	if _, err := os.Stat(awsCredentialsFilename); !os.IsNotExist(err) {
		credentials := credentials.NewSharedCredentials(awsCredentialsFilename, profile)
		return credentials, nil
	}
	return nil, fmt.Errorf("could not read AWS profiles from %s", awsCredentialsFilename)
}
