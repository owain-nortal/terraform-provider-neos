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
	"github.com/owain-nortal/neos-client-go"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
// var (
// 	_ resource.Resource = &registryCoreResource{}
// )

// New data systemResource is a helper function to simplify the provider implementation.
func NewRegistryCoreResource() resource.Resource {
	return &registryCoreResource{}
}

// registryCoreResource is the resource implementation.
type registryCoreResource struct {
	client *neos.NeosClient
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
			"access_key": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The access key",
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
			},
			"partition": schema.StringAttribute{
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "The name of the partition",
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// registryCoreResourceModel maps the resource schema data.
type registryCoreResourceModel struct {
	Identifier  types.String `tfsdk:"identifier"`
	AccessKey   types.String `tfsdk:"access_key"`
	URN         types.String `tfsdk:"urn"`
	Name        types.String `tfsdk:"name"`
	Host        types.String `tfsdk:"host"`
	Partition   types.String `tfsdk:"partition"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// Create a new resource.
func (r *registryCoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//tflog.Info(ctx, "££ Create Get plan")
	// Retrieve values from plan
	var plan registryCoreResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//tflog.Info(ctx, "££ After Create Get plan")

	item := neos.RegistryCorePostRequest{
		Name:      plan.Name.String(),
		Partition: plan.Partition.String(),
	}

	//	tflog.Info(ctx, fmt.Sprintf("££ Create Post request [%s] [%s] [%s] [%s]", plan.ID, plan.Name, plan.Label, plan.Description))

	result, err := r.client.RegistryCorePost(ctx, item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating registry entry for core",
			"Could not create registry entry, unexpected error: "+err.Error(),
		)
		return
	}

	//	tflog.Info(ctx, fmt.Sprintf("££ Create Post result [%s] [%s] [%s] [%s] [%s] [%s]", result.Identifier, result.Name, result.Urn, result.Description, result.Label, result.CreatedAt.String()))

	plan.Identifier = types.StringValue(result.Identifier)
	plan.AccessKey = types.StringValue(result.AccessKey)
	plan.URN = types.StringValue(result.Urn)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//	tflog.Info(ctx, fmt.Sprintf("ID [%s] Desc[%s]", plan.ID, plan.Description))

}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *registryCoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state

	//	tflog.Info(ctx, "££ READ Get current state")

	var state registryCoreResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dataSystemList, err := r.client.RegistryCoreGet()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading NEOS cores from registry",
			"Could not read NEOS  data system ID "+": "+err.Error(),
		)
		return
	}

	for _, ds := range dataSystemList.Cores {
		//		tflog.Info(ctx, fmt.Sprintf("££ READ ITEM: [%s] [%s] %v", ds.Identifier, state.ID.ValueString(), (ds.Identifier == state.ID.ValueString())))
		if ds.Name == state.Name.ValueString() {
			//			tflog.Info(ctx, fmt.Sprintf("££ READ got one in list [%s]", ds.Identifier))
			state.Host = types.StringValue(ds.Host)
			state.Name = types.StringValue(ds.Name)
			state.URN = types.StringValue(ds.Urn)
			break
		}
	}

	//	tsv, _ := state.ID.ToStringValue(ctx)
	// Set refreshed state
	//	tflog.Info(ctx, "££ READ iterate over list")
	//	tflog.Info(ctx, tsv.String())
	//	tflog.Info(ctx, state.ID.ValueString())

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

	// Can't do update

	//tflog.Info(ctx, "££££ update After the ranges")

	// item := neos.RegistryCorePutRequest{
	// 	Entity: neos.RegistryCorePutRequestEntity{
	// 		Name:        plan.Name.String(),
	// 		Label:       plan.Label.String(),
	// 		Description: plan.Description.String(),
	// 	},
	// }

	// eItem := neos.RegistryCorePutRequestEntityInfo{
	// 	Owner:      plan.Owner.String(),
	// 	ContactIds: contacts,
	// 	Links:      links,
	// }

	// result, err := r.client.RegistryCorePut(ctx, plan.ID.ValueString(), item)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error updating data system",
	// 		"Could not put data system, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }
	// //tflog.Info(ctx, fmt.Sprintf("££ Create Post result [%s] [%s] [%s] [%s] [%s] [%s]", result.Identifier, result.Name, result.Urn, result.Description, result.Label, result.CreatedAt.String()))

	// infoResult, err := r.client.RegistryCorePutInfo(ctx, plan.ID.ValueString(), eItem)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error updating data system",
	// 		"Could not put data system, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }

	// contactsList, _ := types.ListValueFrom(ctx, types.StringType, infoResult.ContactIds)
	// linksList, _ := types.ListValueFrom(ctx, types.StringType, infoResult.Links)

	// plan.ID = types.StringValue(result.Identifier)
	// plan.Name = types.StringValue(result.Name)
	// plan.URN = types.StringValue(result.Urn)
	// plan.Description = types.StringValue(result.Description)
	// plan.Label = types.StringValue(result.Label)
	// plan.CreatedAt = types.StringValue(result.CreatedAt.String())
	// plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	// plan.ContactIds = contactsList
	// plan.Links = linksList
	// plan.Owner = types.StringValue(infoResult.Owner)
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
func (r *registryCoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var plan registryCoreResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rcdr := neos.RegistryCoreDeleteRequest{
		Urn: plan.URN.ValueString(),
	}

	err := r.client.RegistryCoreDelete(ctx, rcdr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting core from registry",
			"Could not delete core from registry, unexpected error: "+err.Error(),
		)
		return
	}

}

func (r *registryCoreResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *neos.RegistryCoreClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *registryCoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("urn"), req, resp)
}
