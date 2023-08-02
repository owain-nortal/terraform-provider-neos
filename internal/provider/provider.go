package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/owain-nortal/neos-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &neosProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &neosProvider{
			version: version,
		}
	}
}

type neosProviderModel struct {
	IAMHost      types.String `tfsdk:"iam_host"`
	CoreHost     types.String `tfsdk:"core_host"`
	RegistryHost types.String `tfsdk:"registry_host"`

	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// neosProvider is the provider implementation.
type neosProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *neosProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "neos"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *neosProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"iam_host": schema.StringAttribute{
				Optional: true,
			},
			"registry_host": schema.StringAttribute{
				Optional: true,
			},
			"core_host": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a NEOS API client for data sources and resources.
func (p *neosProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	tflog.Info(ctx, "Configuring NEOS client")

	// Retrieve provider data from configuration
	var config neosProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.IAMHost.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("iam_host"),
			"Unknown NEOS IAM Host",
			"The provider cannot create the NEOS API client as there is an unknown configuration value for the NEOS IAM host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_HOST environment variable.",
		)
	}

	if config.IAMHost.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("registry_host"),
			"Unknown NEOS Registry Host",
			"The provider cannot create the NEOS API client as there is an unknown configuration value for the NEOS Registry host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_HOST environment variable.",
		)
	}

	if config.IAMHost.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("core_host"),
			"Unknown NEOS Core Host",
			"The provider cannot create the NEOS API client as there is an unknown configuration value for the NEOS Core host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown NEOS API Username",
			"The provider cannot create the NEOS API client as there is an unknown configuration value for the NEOS API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown NEOS API Password",
			"The provider cannot create the NEOS API client as there is an unknown configuration value for the NEOS API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	iamhost := os.Getenv("NEOS_IAM_HOST")
	registryhost := os.Getenv("NEOS_REGISTRY_HOST")
	corehost := os.Getenv("NEOS_CORE_HOST")
	username := os.Getenv("NEOS_USERNAME")
	password := os.Getenv("NEOS_PASSWORD")

	if !config.IAMHost.IsNull() {
		iamhost = config.IAMHost.ValueString()
	}

	if !config.CoreHost.IsNull() {
		corehost = config.CoreHost.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if !config.RegistryHost.IsNull() {
		registryhost = config.RegistryHost.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if iamhost == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("iam_host"),
			"Missing NEOS API iam host",
			"The provider cannot create the NEOS API client as there is a missing or empty value for the NEOS iam host. "+
				"Set the host value in the configuration or use the NEOS_IAM_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if corehost == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("core_host"),
			"Missing NEOS API core host",
			"The provider cannot create the NEOS API client as there is a missing or empty value for the NEOS core host. "+
				"Set the host value in the configuration or use the NEOS_CORE_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if registryhost == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("registry_host"),
			"Missing NEOS registry host",
			"The provider cannot create the NEOS API client as there is a missing or empty value for the NEOS registry host. "+
				"Set the host value in the configuration or use the NEOS_REGISTRY_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing NEOS API Username",
			"The provider cannot create the NEOS API client as there is a missing or empty value for the NEOS username. "+
				"Set the username value in the configuration or use the NEOS_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing NEOS API Password",
			"The provider cannot create the NEOS API client as there is a missing or empty value for the NEOS password. "+
				"Set the password value in the configuration or use the NEOS_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "neos_host", iamhost)
	ctx = tflog.SetField(ctx, "neos_username", username)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "neos_password")
	tflog.Info(ctx, "Creating NEOS IAM client")

	iamUrl := fmt.Sprintf("%s/api/iam", iamhost)
	iamClient := neos.NewIAMClient(iamUrl, username, password)
	loginResponse, err := iamClient.Login()

	tflog.Info(ctx, loginResponse.AccessToken)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create NEOS API Client",
			"An unexpected error occurred when creating the NEOS API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"NEOS Client Error: "+err.Error(),
		)
		return
	}

	neos.AccessToken = loginResponse.AccessToken

	client := neos.NewNeosClient(iamhost, registryhost, corehost, "https")

	// Make the NEOS client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = &client
	resp.ResourceData = &client
}

// DataSources defines the data sources implemented in the provider.
func (p *neosProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDataSystemDataSource,
		NewDataProductDataSource,
		NewRegistryCoreDataSource,
		NewDataSourceDataSource,
		NewDataUnitDataSource,
		NewLinksDataSource,
	}
}

func (p *neosProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDataSystemResource,
		NewDataProductResource,
		NewRegistryCoreResource,
		NewDataSourceResource,
		NewDataUnitResource,
		NewLinkDataSourceDataUnitResource,
		NewLinkDataSystemDataSourceResource,
		NewLinkDataUnitDataProductResource,
		NewOutputResource,
		NewLinkDataProductOutputResource,
		NewLinkDataProductDataProductResource,
	}
}
