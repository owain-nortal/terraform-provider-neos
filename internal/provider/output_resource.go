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
// 	_ resource.Resource = &outputResource{}
// )

// New outputResource is a helper function to simplify the provider implementation.
func NewOutputResource() resource.Resource {
	return &outputResource{}
}

// outputResource is the resource implementation.
type outputResource struct {
	client *neos.NeosClient
}

var (
	_ resource.Resource                = &outputResource{}
	_ resource.ResourceWithConfigure   = &outputResource{}
	_ resource.ResourceWithImportState = &outputResource{}
)

// Metadata returns the resource type name.
func (r *outputResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_output"
}

// Schema defines the schema for the resource.
func (r *outputResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The Unique ID of the output",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"urn": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The URN of the output which is read only",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "Name of the output",
			},
			"description": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "Description of the output",
			},
			"label": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "Label for the output",
			},
			"owner": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "The owner of the output",
			},
			"output_type": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "The output type either dashboard or application",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Required:    false,
				Description: "when the output was created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"contact_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    false,
				Optional:    false,
				Required:    true,
				Description: "list of contacts Ids",
			},
			"links": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    false,
				Optional:    false,
				Required:    true,
				Description: "list of links",
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// outputResourceModel maps the resource schema data.
type outputResourceModel struct {
	ID          types.String `tfsdk:"id"`
	URN         types.String `tfsdk:"urn"`
	Name        types.String `tfsdk:"name"`
	Label       types.String `tfsdk:"label"`
	Description types.String `tfsdk:"description"`
	Owner       types.String `tfsdk:"owner"`
	CreatedAt   types.String `tfsdk:"created_at"`
	Links       types.List   `tfsdk:"links"`
	ContactIds  types.List   `tfsdk:"contact_ids"`
	LastUpdated types.String `tfsdk:"last_updated"`
	OutputType  types.String `tfsdk:"output_type"`
}

// Create a new resource.
func (r *outputResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//tflog.Info(ctx, "££ Create Get plan")
	// Retrieve values from plan
	var plan outputResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//tflog.Info(ctx, "££ After Create Get plan")

	linkList, diag := plan.Links.ToListValue(ctx)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	var links []string
	for _, v := range linkList.Elements() {
		links = append(links, v.String())
	}

	contactIDs, diag := plan.ContactIds.ToListValue(ctx)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	var contacts []string
	for _, v := range contactIDs.Elements() {
		contacts = append(contacts, v.String())
	}

	item := neos.OutputPostRequest{
		Entity: neos.OutputPostRequestEntity{
			Name:        plan.Name.String(),
			Label:       plan.Label.String(),
			Description: plan.Description.String(),
			OutputType:  plan.OutputType.String(),
		},
		EntityInfo: neos.OutputPostRequestEntityInfo{
			Owner:      plan.Owner.String(),
			ContactIds: contacts,
			Links:      links,
		},
	}

	//	tflog.Info(ctx, fmt.Sprintf("££ Create Post request [%s] [%s] [%s] [%s]", plan.ID, plan.Name, plan.Label, plan.Description))

	result, err := r.client.OutputPost(ctx, item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating   output",
			"Could not create   output, unexpected error: "+err.Error(),
		)
		return
	}

	//	tflog.Info(ctx, fmt.Sprintf("££ Create Post result [%s] [%s] [%s] [%s] [%s] [%s]", result.Identifier, result.Name, result.Urn, result.Description, result.Label, result.CreatedAt.String()))

	plan.ID = types.StringValue(result.Identifier)
	plan.Name = types.StringValue(result.Name)
	plan.URN = types.StringValue(result.Urn)
	plan.Description = types.StringValue(result.Description)
	plan.Label = types.StringValue(result.Label)
	plan.CreatedAt = types.StringValue(result.CreatedAt.String())
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.OutputType = types.StringValue(result.OutputType)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//	tflog.Info(ctx, fmt.Sprintf("ID [%s] Desc[%s]", plan.ID, plan.Description))

}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *outputResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state

	//	tflog.Info(ctx, "££ READ Get current state")

	var state outputResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	foo := fmt.Sprintf("ID [%s]  Desc [%s]", state.ID.ValueString(), state.Description.ValueString())
	tflog.Info(ctx, foo)

	outputList, err := r.client.OutputGet()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading NEOS output",
			"Could not read NEOS  output ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("££ READ iterate over list looking for: %s", state.ID.ValueString()))
	for _, ds := range outputList.Entities {
		//		tflog.Info(ctx, fmt.Sprintf("££ READ ITEM: [%s] [%s] %v", ds.Identifier, state.ID.ValueString(), (ds.Identifier == state.ID.ValueString())))
		if ds.Identifier == state.ID.ValueString() {
			//			tflog.Info(ctx, fmt.Sprintf("££ READ got one in list [%s]", ds.Identifier))
			state.ID = types.StringValue(ds.Identifier)
			state.Name = types.StringValue(ds.Name)
			state.Label = types.StringValue(ds.Label)
			state.URN = types.StringValue(ds.Urn)
			state.Description = types.StringValue(ds.Description)
			state.Owner = types.StringValue(ds.Owner)
			state.CreatedAt = types.StringValue(ds.CreatedAt.String())
			state.OutputType = types.StringValue(ds.OutputType)
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
func (r *outputResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan outputResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//tflog.Info(ctx, "££ Update After Create Get plan")
	// i, e := plan.ID.ToStringValue(ctx)
	// if e.HasError() {
	// 	tflog.Info(ctx, "Data system update plan get has error")
	// 	return
	// }

	linkList, diag := plan.Links.ToListValue(ctx)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	var links []string
	for _, v := range linkList.Elements() {
		links = append(links, v.String())
	}

	contactIDs, diag := plan.ContactIds.ToListValue(ctx)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	var contacts []string
	for _, v := range contactIDs.Elements() {
		contacts = append(contacts, v.String())
	}

	//tflog.Info(ctx, "££££ update After the ranges")

	item := neos.OutputPutRequest{
		Entity: neos.OutputPutRequestEntity{
			Name:        plan.Name.String(),
			Label:       plan.Label.String(),
			Description: plan.Description.String(),
			OutputType:  plan.OutputType.String(),
		},
	}

	eItem := neos.OutputPutRequestEntityInfo{
		Owner:      plan.Owner.String(),
		ContactIds: contacts,
		Links:      links,
	}

	result, err := r.client.OutputPut(ctx, plan.ID.ValueString(), item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating output",
			"Could not put output, unexpected error: "+err.Error(),
		)
		return
	}
	//tflog.Info(ctx, fmt.Sprintf("££ Create Post result [%s] [%s] [%s] [%s] [%s] [%s]", result.Identifier, result.Name, result.Urn, result.Description, result.Label, result.CreatedAt.String()))

	infoResult, err := r.client.OutputPutInfo(ctx, plan.ID.ValueString(), eItem)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating output",
			"Could not put output, unexpected error: "+err.Error(),
		)
		return
	}

	contactsList, _ := types.ListValueFrom(ctx, types.StringType, infoResult.ContactIds)
	linksList, _ := types.ListValueFrom(ctx, types.StringType, infoResult.Links)

	plan.ID = types.StringValue(result.Identifier)
	plan.Name = types.StringValue(result.Name)
	plan.URN = types.StringValue(result.Urn)
	plan.Description = types.StringValue(result.Description)
	plan.Label = types.StringValue(result.Label)
	plan.CreatedAt = types.StringValue(result.CreatedAt.String())
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.ContactIds = contactsList
	plan.Links = linksList
	plan.Owner = types.StringValue(infoResult.Owner)
	plan.OutputType = types.StringValue(result.OutputType)
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
func (r *outputResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var plan outputResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.OutputDelete(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting output",
			"Could not delete output, unexpected error: "+err.Error(),
		)
		return
	}

}

func (r *outputResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *neos.OutputClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *outputResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
