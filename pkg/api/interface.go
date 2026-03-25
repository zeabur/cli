package api

import (
	"context"
	"time"

	"github.com/zeabur/cli/pkg/model"
)

// Client is the interface of the Zeabur API client.
type Client interface {
	UserAPI
	ProjectAPI
	ServiceAPI
	EnvironmentAPI
	DeploymentAPI
	LogAPI
	GitAPI
	TemplateAPI
	DomainAPI
	VariableAPI
	ServerAPI
	AIHubAPI
	ZSendAPI
	RegisteredDomainAPI
}

type (
	UserAPI interface {
		GetUserInfo(ctx context.Context) (*model.User, error)
	}

	ProjectAPI interface {
		ListProjects(ctx context.Context, skip, limit int) (*model.ProjectConnection, error)
		ListAllProjects(ctx context.Context) (model.Projects, error)
		GetProject(ctx context.Context, id string, ownerName string, name string) (*model.Project, error)
		CreateProject(ctx context.Context, region string, name *string) (*model.Project, error)
		DeleteProject(ctx context.Context, id string) error
		ExportProject(ctx context.Context, id string, environmentID string) (*model.ExportedTemplate, error)

		GetRegions(ctx context.Context) ([]model.Region, error)
		GetServers(ctx context.Context) ([]model.Server, error)
		GetGenericRegions(ctx context.Context) ([]model.GenericRegion, error)

		CloneProject(ctx context.Context, projectID, environmentID, targetRegion string, suspendOldProject bool) (*model.CloneProjectResult, error)
		CloneProjectStatus(ctx context.Context, newProjectID string) (*model.CloneProjectStatusResult, error)
	}

	ServerAPI interface {
		ListServers(ctx context.Context) (model.ServerListItems, error)
		GetServer(ctx context.Context, id string) (*model.ServerDetail, error)
		RebootServer(ctx context.Context, id string) error
		ListDedicatedServerProviders(ctx context.Context) ([]model.CloudProvider, error)
		ListDedicatedServerRegions(ctx context.Context, provider string) ([]model.DedicatedServerRegion, error)
		ListDedicatedServerPlans(ctx context.Context, provider, region string) (model.DedicatedServerPlans, error)
		RentServer(ctx context.Context, provider, region, plan string) (string, error)
		RevealServerPassword(ctx context.Context, serverID string) (string, error)
	}

	EnvironmentAPI interface {
		ListEnvironments(ctx context.Context, projectID string) (model.Environments, error)
		GetEnvironment(ctx context.Context, id string) (*model.Environment, error)
	}

	ServiceAPI interface {
		ListServices(ctx context.Context, projectID string, skip, limit int) (*model.Connection[model.Service], error)
		ListAllServices(ctx context.Context, projectID string) (model.Services, error)
		ListServicesDetailByEnvironment(ctx context.Context, projectID, environmentID string, skip, limit int) (*model.Connection[model.ServiceDetail], error)
		ListAllServicesDetailByEnvironment(ctx context.Context, projectID, environmentID string) (model.ServiceDetails, error)
		GetService(ctx context.Context, id, ownerName, projectName, name string) (*model.Service, error)
		GetServiceDetailByEnvironment(ctx context.Context, id, ownerName, projectName, name, environmentID string) (*model.ServiceDetail, error)
		ServiceMetric(ctx context.Context, id, projectID, environmentID, metricType string, startTime, endTime time.Time) (*model.ServiceMetric, error)
		ServiceInstructions(ctx context.Context, id, environmentID string) ([]model.ServiceInstruction, error)
		GetPrebuiltItems(ctx context.Context) ([]model.PrebuiltItem, error)
		SearchGitRepositories(ctx context.Context, keyword *string) ([]model.GitRepo, error)

		RestartService(ctx context.Context, id string, environmentID string) error
		RedeployService(ctx context.Context, id string, environmentID string) error
		SuspendService(ctx context.Context, id string, environmentID string) error
		CreatePrebuiltService(ctx context.Context, projectID string, marketplaceCode string) (*model.Service, error)
		CreatePrebuiltServiceCustom(ctx context.Context, projectID string, schema model.ServiceSpecSchemaInput) (*model.Service, error)
		CreatePrebuiltServiceRaw(ctx context.Context, projectID string, rawSchema string) (*model.Service, error)
		CreateService(ctx context.Context, projectID string, name string, repoID int, branchName string) (*model.Service, error)
		CreateEmptyService(ctx context.Context, projectID string, name string) (*model.Service, error)
		UploadZipToService(ctx context.Context, projectID string, serviceID string, environmentID string, zipBytes []byte) (*model.Service, error)
		GetDNSName(ctx context.Context, serviceID string) (string, error)
		GetPortForwardingMode(ctx context.Context, serviceID string, environmentID string) (model.PortForwardingMode, error)
		UpdatePortForwardingMode(ctx context.Context, serviceID string, environmentID string, mode model.PortForwardingMode) error
		GetServicePorts(ctx context.Context, serviceID string, environmentID string) ([]model.ServicePort, error)
		GetPortForwardedHost(ctx context.Context, serviceID string) (string, error)
		UpdateImageTag(ctx context.Context, serviceID string, environmentID string, tag string) error
		DeleteService(ctx context.Context, id string) error
		ExecuteCommand(ctx context.Context, serviceID string, environmentID string, command []string) (*model.CommandResult, error)
	}

	VariableAPI interface {
		ListVariables(ctx context.Context, serviceID string, environmentID string) (model.Variables, model.Variables, error)
		UpdateVariables(ctx context.Context, serviceID string, environmentID string, data map[string]string) (bool, error)
	}

	DomainAPI interface {
		AddDomain(ctx context.Context, serviceID string, environmentID string, isGenerated bool, domain string, options ...string) (*string, error)
		ListDomains(ctx context.Context, serviceID string, environmentID string) (model.Domains, error)
		RemoveDomain(ctx context.Context, domain string) (bool, error)
		CheckDomainAvailable(ctx context.Context, domain string, isGenerated bool, region string) (bool, string, error)
	}

	DeploymentAPI interface {
		ListDeployments(ctx context.Context, serviceID string, environmentID string, perPage int) (*model.DeploymentConnection, error)
		ListAllDeployments(ctx context.Context, serviceID string, environmentID string) (model.Deployments, error)
		GetDeployment(ctx context.Context, id string) (*model.Deployment, error)
		GetLatestDeployment(ctx context.Context, serviceID string, environmentID string) (*model.Deployment, bool, error)
	}

	LogAPI interface {
		GetRuntimeLogs(ctx context.Context, serviceID, environmentID, deploymentID string) (model.Logs, error)
		GetBuildLogs(ctx context.Context, deploymentID string) (model.Logs, error)

		WatchRuntimeLogs(ctx context.Context, projectID, serviceID, environmentID, deploymentID string) (<-chan model.Log, error)
		WatchBuildLogs(ctx context.Context, projectID, deploymentID string) (<-chan model.Log, error)
	}

	GitAPI interface {
		GetRepoBranches(ctx context.Context, repoOwner string, repoName string) ([]string, error)
		GetRepoID(repoOwner string, repoName string) (int, error)
		GetRepoInfo() (string, string, error)
		GetRepoBranchesByRepoID(repoID int) ([]string, error)
	}

	AIHubAPI interface {
		GetAIHubTenant(ctx context.Context) (*model.AIHubTenant, error)
		AddAIHubBalance(ctx context.Context, amount int, provider *string) (*model.AddAIHubBalanceResult, error)
		CreateAIHubKey(ctx context.Context, alias *string) (*model.CreateAIHubKeyResult, error)
		DeleteAIHubKey(ctx context.Context, keyID string) error
		UpdateAIHubAutoRechargeSettings(ctx context.Context, threshold, amount int) (*model.UpdateAIHubAutoRechargeSettingsResult, error)
		GetAIHubSpendLogs(ctx context.Context, startDate, endDate *time.Time) ([]model.AIHubSpendLog, error)
		GetAIHubMonthlyUsage(ctx context.Context, month *string) (*model.AIHubMonthlyUsage, error)
	}

	TemplateAPI interface {
		ListTemplates(ctx context.Context, skip, limit int) (*model.TemplateConnection, error)
		ListAllTemplates(ctx context.Context) (model.Templates, error)
		GetTemplate(ctx context.Context, code string) (*model.Template, error)

		DeployTemplate(
			ctx context.Context,
			rawSpecYaml string,
			variables model.Map,
			repoConfigs model.RepoConfigs,
			projectID string,
		) (*model.Project, error)
		DeleteTemplate(ctx context.Context, code string) error

		CreateTemplateFromFile(ctx context.Context, rawSpecYaml string) (*model.Template, error)
		UpdateTemplateFromFile(ctx context.Context, code, rawSpecYaml string) (bool, error)
	}

	RegisteredDomainAPI interface {
		CheckDomainRegistrationAvailability(ctx context.Context, domain string) (*model.DomainSearchResult, error)
		PurchaseDomain(ctx context.Context, domain, registrantProfileID string) (*model.PurchaseDomainResult, error)
		ListRegisteredDomains(ctx context.Context) (model.RegisteredDomains, error)
		GetRegisteredDomain(ctx context.Context, id string) (*model.RegisteredDomain, error)
		RenewDomain(ctx context.Context, id string) (*model.RegisteredDomain, error)
		SetDomainAutoRenew(ctx context.Context, id string, autoRenew bool) (*model.RegisteredDomain, error)

		ListDNSRecords(ctx context.Context, registeredDomainID string) (model.DNSRecords, error)
		CreateDNSRecord(ctx context.Context, registeredDomainID string, input model.CreateDNSRecordInput) (*model.DNSRecord, error)
		UpdateDNSRecord(ctx context.Context, registeredDomainID, recordID string, input model.UpdateDNSRecordInput) (*model.DNSRecord, error)
		DeleteDNSRecord(ctx context.Context, registeredDomainID, recordID string) error

		ListRegistrantProfiles(ctx context.Context) (model.RegistrantProfiles, error)
		CreateRegistrantProfile(ctx context.Context, input model.CreateRegistrantProfileInput) (*model.RegistrantProfile, error)
		UpdateRegistrantProfile(ctx context.Context, id string, input model.UpdateRegistrantProfileInput) (*model.RegistrantProfile, error)
		DeleteRegistrantProfile(ctx context.Context, id string) error

		ResendRegistrantVerificationEmail(ctx context.Context, registeredDomainID string) error
		UpdateRegistrantContact(ctx context.Context, registeredDomainID string, input model.UpdateRegistrantContactInput) error
	}

	ZSendAPI interface {
		GetZSendOnboardingStatus(ctx context.Context) (*model.ZSendOnboardingStatus, error)
		GetZSendUserStatus(ctx context.Context) (*model.ZSendUserStatus, error)
		OnboardZSend(ctx context.Context) (*model.ZSendOnboardingStatus, error)

		ListZSendDomains(ctx context.Context, page, pageSize *int) (*model.ListZSendDomainsReply, error)
		GetZSendDomain(ctx context.Context, id string) (*model.ZSendDomain, error)
		CreateZSendDomain(ctx context.Context, domain, region string) (*model.ZSendDomain, error)
		VerifyZSendDomain(ctx context.Context, id string) (*model.ZSendDomain, error)
		DeleteZSendDomain(ctx context.Context, id string) error

		ListZSendAPIKeys(ctx context.Context, page, pageSize *int) (*model.ListZSendAPIKeysReply, error)
		GetZSendAPIKey(ctx context.Context, id string) (*model.ZSendAPIKey, error)
		CreateZSendAPIKey(ctx context.Context, input model.CreateZSendAPIKeyInput) (*model.CreateZSendAPIKeyReply, error)
		DeleteZSendAPIKey(ctx context.Context, id string) error

		ListZSendWebhooks(ctx context.Context, page, pageSize *int) (*model.ListZSendWebhooksReply, error)
		GetZSendWebhook(ctx context.Context, id string) (*model.ZSendWebhook, error)
		CreateZSendWebhook(ctx context.Context, input model.CreateZSendWebhookInput) (*model.CreateZSendWebhookReply, error)
		DeleteZSendWebhook(ctx context.Context, id string) error
		VerifyZSendWebhook(ctx context.Context, id string) (*model.VerifyZSendWebhookReply, error)

		// Email record queries (GraphQL)
		ListZSendEmails(ctx context.Context, page, pageSize *int, status, jobType, jobID *string) (*model.ListZSendEmailsReply, error)
		GetZSendEmail(ctx context.Context, id string) (*model.ZSendEmail, error)

		// Send (REST, requires Z-Send API key)
		SendZSendEmail(ctx context.Context, apiKey string, req model.ZSendSendEmailRequest) (*model.ZSendSendEmailReply, error)
		ScheduleZSendEmail(ctx context.Context, apiKey string, req model.ZSendScheduleEmailRequest) (*model.ZSendScheduleEmailReply, error)
		SendZSendBatchEmail(ctx context.Context, apiKey string, req model.ZSendBatchEmailRequest) (*model.ZSendBatchEmailReply, error)

		// Scheduled email management (REST, requires Z-Send API key)
		ListZSendScheduledEmails(ctx context.Context, apiKey string, page, pageSize *int, status *string) (*model.ZSendListScheduledEmailsReply, error)
		GetZSendScheduledEmail(ctx context.Context, apiKey string, id string) (*model.ZSendScheduledEmail, error)
		CancelZSendScheduledEmail(ctx context.Context, apiKey string, id string) error

		// Batch job management (REST, requires Z-Send API key)
		ListZSendBatchEmailJobs(ctx context.Context, apiKey string, page, pageSize *int, status *string) (*model.ZSendListBatchJobsReply, error)
		GetZSendBatchEmailJob(ctx context.Context, apiKey string, id string) (*model.ZSendBatchJob, error)
	}
)
