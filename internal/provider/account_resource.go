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
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
// var (
// 	_ resource.Resource = &accountResource{}
// )

// New accountResource is a helper function to simplify the provider implementation.
func NewAccountResource() resource.Resource {
	return &accountResource{}
}

// accountResource is the resource implementation.
type accountResource struct {
	client *neos.AccountClient
}

var (
	_ resource.Resource                = &accountResource{}
	_ resource.ResourceWithConfigure   = &accountResource{}
	_ resource.ResourceWithImportState = &accountResource{}
)

// Metadata returns the resource type name.
func (r *accountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

// Schema defines the schema for the resource.
func (r *accountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The Unique ID of the account",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"urn": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The URN of the account which is read only",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "Name of the account",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "Description of the account",
			},
			"display_name": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "Display name",
			},
			"owner": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "The owner of the account",
			},
			"is_system": schema.BoolAttribute{
				Computed:    true,
				Optional:    false,
				Required:    false,
				Description: "Is system",
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// accountResourceModel maps the resource schema data.
type accountResourceModel struct {
	ID          types.String `tfsdk:"id"`
	URN         types.String `tfsdk:"urn"`
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Owner       types.String `tfsdk:"owner"`
	IsSystem    types.Bool   `tfsdk:"is_system"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// Create a new resource.
func (r *accountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//tflog.Info(ctx, "££ Create Get plan")
	// Retrieve values from plan
	var plan accountResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := neos.AccountPostRequest{
		Name:        plan.Name.String(),
		DisplayName: plan.DisplayName.String(),
		Description: plan.Description.String(),
		Owner:       plan.Owner.String(),
	}

	result, err := r.client.Post(ctx, item)
	if err != nil {
		resp.Diagnostics.AddError("Error creating   account", "Could not create   account, unexpected error: "+err.Error())
		return
	}

	plan.ID = types.StringValue(result.Identifier)
	plan.Name = types.StringValue(result.Name)
	plan.URN = types.StringValue(result.Urn)
	plan.Owner = types.StringValue(result.Owner)
	plan.Description = types.StringValue(result.Description)
	plan.IsSystem = types.BoolValue(result.IsSystem)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.DisplayName = types.StringValue(result.DisplayName)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *accountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var state accountResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	foo := fmt.Sprintf("ID [%s]  Desc [%s]", state.ID.ValueString(), state.Description.ValueString())
	tflog.Info(ctx, foo)

	accountList, err := r.client.Get("")
	if err != nil {
		resp.Diagnostics.AddError("Error Reading NEOS account", "Could not read NEOS  account ID "+state.ID.ValueString()+": "+err.Error())
		return
	}

	for _, ds := range accountList.Accounts {
		if ds.Identifier == state.ID.ValueString() {
			state.ID = types.StringValue(ds.Identifier)
			state.Name = types.StringValue(ds.Name)
			state.Owner = types.StringValue(ds.Owner)
			state.URN = types.StringValue(ds.Urn)
			state.Description = types.StringValue(ds.Description)
			state.Owner = types.StringValue(ds.Owner)
			state.IsSystem = types.BoolValue(ds.IsSystem)
			state.DisplayName = types.StringValue(ds.DisplayName)
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
func (r *accountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan accountResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := neos.AccountPutRequest{
		DisplayName: plan.DisplayName.String(),
		Owner:       plan.Owner.String(),
		Description: plan.Description.String(),
	}

	result, err := r.client.Put(ctx, plan.ID.ValueString(), item)
	if err != nil {
		resp.Diagnostics.AddError("Error updating account", "Could not put account, unexpected error: "+err.Error())
		return
	}

	plan.ID = types.StringValue(result.Identifier)
	plan.Name = types.StringValue(result.Name)
	plan.DisplayName = types.StringValue(result.DisplayName)
	plan.URN = types.StringValue(result.Urn)
	plan.Description = types.StringValue(result.Description)
	plan.IsSystem = types.BoolValue(result.IsSystem)
	plan.Owner = types.StringValue(result.Owner)
	//plan.CreatedAt = types.StringValue(result.CreatedAt.String())
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *accountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var plan accountResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting account", "Could not delete account, unexpected error: "+err.Error())
		return
	}

}

func (r *accountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)

	if !ok {
		resp.Diagnostics.AddError("Unexpected Account Configure Type", fmt.Sprintf("Expected *neos.AccountClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	r.client = &client.AccountClient

}

func (r *accountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
