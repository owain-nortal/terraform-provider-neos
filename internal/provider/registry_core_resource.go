package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	neos "github.com/owain-nortal/neos-client-go"
)

// Ensure the implementation satisfies the expected interfaces.

// New data systemResource is a helper function to simplify the provider implementation.
func NewRegistryCoreResource() resource.Resource {
	return &registryCoreResource{}
}

// registryCoreResource is the resource implementation.
type registryCoreResource struct {
	client *neos.RegistryCoreClient
}

var (
	_ resource.Resource                = &registryCoreResource{}
	_ resource.ResourceWithConfigure   = &registryCoreResource{}
	_ resource.ResourceWithImportState = &registryCoreResource{}
)

// Metadata returns the resource type name.
func (r *registryCoreResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registry_core"
}

// Schema defines the schema for the resource.
func (r *registryCoreResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: "The host which is never passed in",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_key_id": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The access key id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secret_key": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The secret access key ",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identifier": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The identifier key",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"urn": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The URN of the data system which is read only",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "Name of the core",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"account": schema.StringAttribute{
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "Account if not root",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"partition": schema.StringAttribute{
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "The name of the partition",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// registryCoreResourceModel maps the resource schema data.
type registryCoreResourceModel struct {
	Identifier  types.String `tfsdk:"identifier"`
	AccessKeyId types.String `tfsdk:"access_key_id"`
	SecretKey   types.String `tfsdk:"secret_key"`
	URN         types.String `tfsdk:"urn"`
	Name        types.String `tfsdk:"name"`
	Host        types.String `tfsdk:"host"`
	Partition   types.String `tfsdk:"partition"`
	Account     types.String `tfsdk:"account"`
}

// Create a new resource.
func (r *registryCoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "registryCoreResource Create")
	// Retrieve values from plan
	var plan registryCoreResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := neos.RegistryCorePostRequest{
		Name:      plan.Name.String(),
		Partition: plan.Partition.String(),
	}

	result, err := r.client.Post(ctx, item, plan.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating registry entry for core", "Could not create registry entry, unexpected error: "+err.Error())
		return
	}

	plan.Identifier = types.StringValue(result.Identifier)
	plan.AccessKeyId = types.StringValue(result.KeyPair.AccessKeyID)
	plan.SecretKey = types.StringValue(result.KeyPair.SecretAccessKey)
	plan.URN = types.StringValue(result.Urn)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *registryCoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state

	var state registryCoreResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dataSystemList, err := r.client.Get(state.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading NEOS cores from registry", "Could not read NEOS  data system ID "+": "+err.Error())
		return
	}

	for _, ds := range dataSystemList.Cores {
		if ds.Name == state.Name.ValueString() {
			state.Host = types.StringValue(ds.Host)
			state.Name = types.StringValue(ds.Name)
			state.URN = types.StringValue(ds.Urn)
			break
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Info(ctx, "Data system Read Has error")
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *registryCoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan registryCoreResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *registryCoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	tflog.Info(ctx, "registryCoreResource delete")

	var plan registryCoreResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("registryCoreResource delete id %s", plan.Identifier.ValueString()))

	err := r.client.Delete(ctx, plan.Identifier.ValueString(), plan.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting core from registry", "Could not delete core from registry, unexpected error: "+err.Error())
		return
	}

}

func (r *registryCoreResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)

	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *neos.RegistryCoreClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	r.client = &client.RegistryCoreClient
}

func (r *registryCoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("urn"), req, resp)
}
