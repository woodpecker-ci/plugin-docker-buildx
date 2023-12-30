// Source: https://github.com/drone-plugins/drone-docker/tree/939591f01828eceae54f5768dc7ce08ad0ad0bba/cmd/drone-ecr
package plugin

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

const DefaultRegion = "us-east-1"

var repo string
var assumeRole string
var externalID string
var ecr_login Login
var aws_region string

func (p *Plugin) EcrInit() {

	// create a standalone Login object to account for single repo and multi-repo case
	if len(p.settings.Logins) >= 1 {
		for _, login := range p.settings.Logins {
			if strings.Contains(login.Registry, "amazonaws.com") {
				ecr_login = login
				aws_region = login.Aws_region

				// filter repo containing ecr registry
				substrings := make([]string, 0)
				for _, repo := range p.settings.Build.Repo.Value() {
					substrings = append(substrings, strings.Split(repo, ",")...)
				}
				filtered := make([]string, 0)
				for _, s := range substrings {
					if strings.Contains(s, "amazonaws.com") {
						filtered = append(filtered, s)
					}
				}

				// Join the filtered substrings into a comma-separated string
				repo = strings.Join(filtered, ",")

				// set the region
				if aws_region == "" {
					aws_region = DefaultRegion
				}

				os.Setenv("AWS_REGION", aws_region)
				os.Setenv("AWS_ACCESS_KEY_ID", ecr_login.Aws_access_key_id)
				os.Setenv("AWS_SECRET_ACCESS_KEY", ecr_login.Aws_secret_access_key)

			}
		}
	} else {
		ecr_login.Aws_access_key_id = p.settings.AwsAccessKeyId
		ecr_login.Aws_secret_access_key = p.settings.AwsSecretAccessKey
		aws_region = p.settings.AwsRegion
		repo = p.settings.Build.Repo.Value()[0]

		// set the region
		if aws_region == "" {
			aws_region = DefaultRegion
		}

		os.Setenv("AWS_REGION", p.settings.AwsRegion)
		os.Setenv("AWS_ACCESS_KEY_ID", p.settings.AwsAccessKeyId)
		os.Setenv("AWS_SECRET_ACCESS_KEY", p.settings.AwsSecretAccessKey)
	}
	// here the env vars are used for authentication
	sess, err := session.NewSession(&aws.Config{Region: &aws_region})
	if err != nil {
		log.Fatalf("error creating aws session: %v", err)
	}

	svc := getECRClient(sess, assumeRole, externalID)
	username, password, registry, err := getAuthInfo(svc)

	if err != nil {
		log.Fatalf("error getting ECR auth: %v", err)
	}

	if !strings.HasPrefix(repo, registry) {
		repo = fmt.Sprintf("%s/%s", registry, repo)
	}

	if p.settings.EcrCreateRepository {
		err = ensureRepoExists(svc, trimHostname(repo, registry), p.settings.EcrScanOnPush)
		if err != nil {
			log.Fatalf("error creating ECR repo: %v", err)
		}
		err = updateImageScannningConfig(svc, trimHostname(repo, registry), p.settings.EcrScanOnPush)
		if err != nil {
			log.Fatalf("error updating scan on push for ECR repo: %v", err)
		}
	}

	if p.settings.EcrLifecyclePolicy != "" {
		p, err := os.ReadFile(p.settings.EcrLifecyclePolicy)
		if err != nil {
			log.Fatal(err)
		}
		if err := uploadLifeCyclePolicy(svc, string(p), trimHostname(repo, registry)); err != nil {
			log.Fatalf("error uploading ECR lifecycle policy: %v", err)
		}
	}

	if p.settings.EcrRepositoryPolicy != "" {
		p, err := os.ReadFile(p.settings.EcrRepositoryPolicy)
		if err != nil {
			log.Fatal(err)
		}
		if err := uploadRepositoryPolicy(svc, string(p), trimHostname(repo, registry)); err != nil {
			log.Fatalf("error uploading ECR repository policy. %v", err)
		}
	}

	// set Username and Password for all Login which contain an AWS key
	if len(p.settings.Logins) >= 1 {
		for i, login := range p.settings.Logins {
			if login.Aws_secret_access_key != "" && login.Aws_access_key_id != "" {
				p.settings.Logins[i].Username = username
				p.settings.Logins[i].Password = password
				p.settings.Logins[i].Registry = registry
			}
		}
	} else {
		p.settings.DefaultLogin.Username = username
		p.settings.DefaultLogin.Password = password
		p.settings.DefaultLogin.Registry = registry
	}

}

func trimHostname(repo, registry string) string {
	repo = strings.TrimPrefix(repo, registry)
	repo = strings.TrimLeft(repo, "/")
	return repo
}

func ensureRepoExists(svc *ecr.ECR, name string, scanOnPush bool) (err error) {
	input := &ecr.CreateRepositoryInput{}
	input.SetRepositoryName(name)
	input.SetImageScanningConfiguration(&ecr.ImageScanningConfiguration{ScanOnPush: &scanOnPush})
	_, err = svc.CreateRepository(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == ecr.ErrCodeRepositoryAlreadyExistsException {
			// eat it, we skip checking for existing to save two requests
			err = nil
		}
	}

	return
}

func updateImageScannningConfig(svc *ecr.ECR, name string, scanOnPush bool) (err error) {
	input := &ecr.PutImageScanningConfigurationInput{}
	input.SetRepositoryName(name)
	input.SetImageScanningConfiguration(&ecr.ImageScanningConfiguration{ScanOnPush: &scanOnPush})
	_, err = svc.PutImageScanningConfiguration(input)

	return err
}

func uploadLifeCyclePolicy(svc *ecr.ECR, lifecyclePolicy string, name string) (err error) {
	input := &ecr.PutLifecyclePolicyInput{}
	input.SetLifecyclePolicyText(lifecyclePolicy)
	input.SetRepositoryName(name)
	_, err = svc.PutLifecyclePolicy(input)

	return err
}

func uploadRepositoryPolicy(svc *ecr.ECR, repositoryPolicy string, name string) (err error) {
	input := &ecr.SetRepositoryPolicyInput{}
	input.SetPolicyText(repositoryPolicy)
	input.SetRepositoryName(name)
	_, err = svc.SetRepositoryPolicy(input)

	return err
}

func getAuthInfo(svc *ecr.ECR) (username, password, registry string, err error) {
	var result *ecr.GetAuthorizationTokenOutput
	var decoded []byte

	result, err = svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return
	}

	auth := result.AuthorizationData[0]
	token := *auth.AuthorizationToken
	decoded, err = base64.StdEncoding.DecodeString(token)
	if err != nil {
		return
	}

	registry = strings.TrimPrefix(*auth.ProxyEndpoint, "https://")
	creds := strings.Split(string(decoded), ":")
	username = creds[0]
	password = creds[1]
	return
}

func getECRClient(sess *session.Session, role string, externalId string) *ecr.ECR {
	if role == "" {
		return ecr.New(sess)
	}
	if externalId != "" {
		return ecr.New(sess, &aws.Config{
			Credentials: stscreds.NewCredentials(sess, role, func(p *stscreds.AssumeRoleProvider) {
				p.ExternalID = &externalId
			}),
		})
	} else {
		return ecr.New(sess, &aws.Config{
			Credentials: stscreds.NewCredentials(sess, role),
		})
	}
}
