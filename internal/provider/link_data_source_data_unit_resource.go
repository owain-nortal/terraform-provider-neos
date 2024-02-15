package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/owain-nortal/neos-client-go"
	"time"
)


func NewLinkDataSourceDataUnitResource() resource.Resource {
	return &linkDataSourceDataUnitResource{}
}

type linkDataSourceDataUnitResource struct {
	client *neos.LinksClient
}

var (
	_ resource.Resource                = &linkDataSourceDataUnitResource{}
	_ resource.ResourceWithConfigure   = &linkDataSourceDataUnitResource{}
	_ resource.ResourceWithImportState = &linkDataSourceDataUnitResource{}
)

// Metadata returns the resource type name.
func (r *linkDataSourceDataUnitResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_link_data_source_data_unit"
}

// Schema defines the schema for the resource.
func (r *linkDataSourceDataUnitResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"parent_identifier": schema.StringAttribute{
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "The data source parent identifier",
			},
			"child_identifier": schema.StringAttribute{
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "The data unit child identifier",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "Last updated time",
			},
		},
	}
}

// linkDataSourceDataUnitResourceModel maps the resource schema data.
type linkDataSourceDataUnitResourceModel struct {
	ParentIdentifier types.String `tfsdk:"parent_identifier"`
	ChildIdentifier  types.String `tfsdk:"child_identifier"`
	LastUpdated      types.String `tfsdk:"last_updated"`
}

// Create a new resource.
func (r *linkDataSourceDataUnitResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "linkDataSourceDataUnitResource Create Get plan")
	// Retrieve values from plan
	var plan linkDataSourceDataUnitResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//tflog.Info(ctx, "££ After Create Get plan")

	tflog.Info(ctx, fmt.Sprintf("linkDataSourceDataUnitResource Create Post request [%s] [%s]", plan.ParentIdentifier.ValueString(), plan.ChildIdentifier.ValueString()))

	result, err := r.client.LinkDataSourceToDataUnit(ctx, plan.ParentIdentifier.ValueString(), plan.ChildIdentifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating link",
			"Could not create link, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("linkDataSourceDataUnitResource Create Post result [%s] [%s] ", result.Parent.Identifier, result.Child.Identifier))

	plan.ParentIdentifier = types.StringValue(result.Parent.Identifier)
	plan.ChildIdentifier = types.StringValue(result.Child.Identifier)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *linkDataSourceDataUnitResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state

	tflog.Info(ctx, "linkDataSourceDataUnitResource READ Get current state")

	var state linkDataSourceDataUnitResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("linkDataSourceDataUnitResource Parent ID [%s]  Desc [%s]", state.ParentIdentifier.ValueString(), state.ChildIdentifier.ValueString()))

	linksList, err := r.client.Get()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading NEOS data system",
			"Could not read Links "+state.ParentIdentifier.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("linkDataSourceDataUnitResource READ iterate over list looking for: %s", state.ParentIdentifier.ValueString()))
	for _, ds := range linksList.Links {
		tflog.Info(ctx, fmt.Sprintf(">>>>> linkDataSourceDataUnitResource READ ITEM: [%s] [%s] ", ds.Parent.Identifier, state.ParentIdentifier.ValueString()))
		if ds.Parent.Identifier == state.ParentIdentifier.ValueString() {
			tflog.Info(ctx, fmt.Sprintf("linkDataSourceDataUnitResource READ got one in list [%s]", ds.Parent.Identifier))
			state.ParentIdentifier = types.StringValue(ds.Parent.Identifier)
			state.ChildIdentifier = types.StringValue(ds.Child.Identifier)
			break
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Info(ctx, "linkDataSourceDataUnitResource Links Read Has error")
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *linkDataSourceDataUnitResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan linkDataSourceDataUnitResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.LinkDataSourceToDataUnit(ctx, plan.ParentIdentifier.ValueString(), plan.ChildIdentifier.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating link",
			"Could not post link , unexpected error: "+err.Error(),
		)
		return
	}

	plan.ParentIdentifier = types.StringValue(result.Parent.Identifier)
	plan.ChildIdentifier = types.StringValue(result.Child.Identifier)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *linkDataSourceDataUnitResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var plan linkDataSourceDataUnitResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteLinkDataSourceToDataUnit(ctx, plan.ParentIdentifier.ValueString(), plan.ChildIdentifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting link", "Could not delete link, unexpected error: "+err.Error())
		return
	}

}

func (r *linkDataSourceDataUnitResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)

	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *neos.LinksClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	r.client = &client.LinksClient
}

func (r *linkDataSourceDataUnitResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
