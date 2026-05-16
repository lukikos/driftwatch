// Package checker dispatches infrastructure drift checks by type.
package checker

import (
	"fmt"

	"github.com/yourusername/driftwatch/internal/config"
)

// Checker dispatches checks to the appropriate handler.
type Checker struct{}

// New returns a new Checker.
func New() *Checker {
	return &Checker{}
}

// Check runs the appropriate check for the given config.Check and returns
// (drifted bool, message string, err error).
func (c *Checker) Check(check config.Check) (bool, string, error) {
	switch check.Type {
	case "env_var":
		return checkEnvVar(check)
	case "file_hash":
		return checkFileHash(check)
	case "http_status":
		return checkHTTPStatus(check)
	case "process_running":
		return checkProcessRunning(check)
	case "port_open":
		return checkPortOpen(check)
	case "docker_container":
		return checkDockerContainer(check)
	case "sys_command":
		return checkSysCommand(check)
	case "dns_resolve":
		return checkDNSResolve(check)
	case "ssl_expiry":
		return checkSSLExpiry(check)
	case "file_content":
		return checkFileContent(check)
	case "file_exists":
		return checkFileExists(check)
	case "dir_size":
		return checkDirSize(check)
	case "last_cron_run":
		return checkLastCronRun(check)
	case "env_file":
		return checkEnvFile(check)
	case "k8s_pod":
		return checkK8sPod(check)
	case "disk_usage":
		return checkDiskUsage(check)
	case "http_latency":
		return checkHTTPLatency(check)
	case "service_status":
		return checkServiceStatus(check)
	case "cert_pin":
		return checkCertPin(check)
	case "mount_point":
		return checkMountPoint(check)
	case "symlink":
		return checkSymlink(check)
	case "json_field":
		return checkJSONField(check)
	case "s3_bucket_access":
		return checkS3BucketAccess(check)
	case "ec2_instance_metadata":
		return checkEC2InstanceMetadata(check)
	case "lambda_function":
		return checkLambdaFunction(check)
	case "rds_instance_status":
		return checkRDSInstanceStatus(check)
	case "sqs_queue_attributes":
		return checkSQSQueueAttributes(check)
	case "iam_role_policy":
		return checkIAMRolePolicy(check)
	case "security_group_rules":
		return checkSecurityGroupRules(check)
	case "dynamodb_table":
		return checkDynamoDBTable(check)
	case "ecs_service_status":
		return checkECSServiceStatus(check)
	case "sns_topic_attributes":
		return checkSNSTopicAttributes(check)
	case "cloudwatch_alarm":
		return checkCloudWatchAlarm(check)
	case "route53_health_check":
		return checkRoute53HealthCheck(check)
	case "eks_cluster_status":
		return checkEKSClusterStatus(check)
	case "secrets_manager_secret":
		return checkSecretsManagerSecret(check)
	case "elbv2_target_group_health":
		return checkELBv2TargetGroupHealth(check)
	case "kinesis_stream_status":
		return checkKinesisStreamStatus(check)
	case "glue_job_status":
		return checkGlueJobStatus(check)
	case "ssm_parameter":
		return checkSSMParameter(check)
	case "cloudfront_distribution":
		return checkCloudFrontDistribution(check)
	case "redshift_cluster_status":
		return checkRedshiftClusterStatus(check)
	case "opensearch_domain_status":
		return checkOpenSearchDomainStatus(check)
	case "step_functions_state_machine":
		return checkStepFunctionsStateMachine(check)
	case "ecr_repository_policy":
		return checkECRRepositoryPolicy(check)
	case "elasticache_cluster":
		return checkElastiCacheClusterStatus(check)
	default:
		return false, "", fmt.Errorf("unknown check type: %q", check.Type)
	}
}
