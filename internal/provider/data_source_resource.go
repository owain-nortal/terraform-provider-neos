package provider

import (
	"context"
	"fmt"

	"encoding/json"
	"time"

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
// var (
// 	_ resource.Resource = &dataSourceResource{}
// )

// New data systemResource is a helper function to simplify the provider implementation.
func NewDataSourceResource() resource.Resource {
	return &dataSourceResource{}
}

// dataSourceResource is the resource implementation.
type dataSourceResource struct {
	client                 *neos.DataSourceClient
	connectionClient       *neos.DataSourceConnectionClient
	dataSourceSecretClient *neos.DataSourceSecretClient
	secretClient           *neos.SecretClient
}

var (
	_ resource.Resource                = &dataSourceResource{}
	_ resource.ResourceWithConfigure   = &dataSourceResource{}
	_ resource.ResourceWithImportState = &dataSourceResource{}
)

// Metadata returns the resource type name.
func (r *dataSourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_source"
}

// Schema defines the schema for the resource.
func (r *dataSourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The Unique ID of the data system",
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
				Description: "Name of the data system",
			},
			"description": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "Description of the data system",
			},
			"label": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "Label for the data system",
			},
			"owner": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "The owner of the data system",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Required:    false,
				Description: "when the data system was created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_json": schema.StringAttribute{
				Computed:    false,
				Optional:    false,
				Required:    true,
				Description: "connection json",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secret_json": schema.StringAttribute{
				Computed:    false,
				Optional:    false,
				Required:    true,
				Description: "secret json",
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

// dataSourceResourceModel maps the resource schema data.
type dataSourceResourceModel struct {
	ID             types.String `tfsdk:"id"`
	URN            types.String `tfsdk:"urn"`
	Name           types.String `tfsdk:"name"`
	Label          types.String `tfsdk:"label"`
	Description    types.String `tfsdk:"description"`
	Owner          types.String `tfsdk:"owner"`
	CreatedAt      types.String `tfsdk:"created_at"`
	ConnectionJson types.String `tfsdk:"connection_json"`
	SecretJson     types.String `tfsdk:"secret_json"`
	Links          types.List   `tfsdk:"links"`
	ContactIds     types.List   `tfsdk:"contact_ids"`
	LastUpdated    types.String `tfsdk:"last_updated"`
}

// Create a new resource.
func (r *dataSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//tflog.Info(ctx, "££ Create Get plan")
	// Retrieve values from plan
	var plan dataSourceResourceModel
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

	item := neos.DataSourcePostRequest{
		Entity: neos.DataSourcePostRequestEntity{
			Name:        plan.Name.String(),
			Label:       plan.Label.String(),
			Description: plan.Description.String(),
		},
		EntityInfo: neos.DataSourcePostRequestEntityInfo{
			Owner:      plan.Owner.String(),
			ContactIds: contacts,
			Links:      links,
		},
	}

	//	tflog.Info(ctx, fmt.Sprintf("££ Create Post request [%s] [%s] [%s] [%s]", plan.ID, plan.Name, plan.Label, plan.Description))

	result, err := r.client.Post(ctx, item)
	if err != nil {
		resp.Diagnostics.AddError("Error creating data source", "Could not create   data source, unexpected error: "+err.Error())
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
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	connectionResult, err := r.connectionClient.Put(ctx, result.Identifier, plan.ConnectionJson.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating data source connection", "Could not create data source connection, unexpected error: "+err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Connection result %s %s", result.Identifier, connectionResult))

	var secretMap map[string]string
	secretBytes := []byte(plan.SecretJson.ValueString())
	if err := json.Unmarshal(secretBytes, &secretMap); err != nil {
		resp.Diagnostics.AddError("Error unmarshal connection secret", " unexpected error: "+err.Error())
		return
	}

	secretResult, err := r.dataSourceSecretClient.Post(ctx, result.Identifier, secretMap)
	if err != nil {
		resp.Diagnostics.AddError("Error creating data source secret", "Could not create data source secret, unexpected error: "+err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Connection result %s %s", result.Identifier, secretResult))

}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *dataSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state

	//	tflog.Info(ctx, "££ READ Get current state")

	var state dataSourceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dataSourceList, err := r.client.Get()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading NEOS data system", "Could not read NEOS  data system ID "+state.ID.ValueString()+": "+err.Error())
		return
	}

	for _, ds := range dataSourceList.Entities {
		if ds.Identifier == state.ID.ValueString() {
			state.ID = types.StringValue(ds.Identifier)
			state.Name = types.StringValue(ds.Name)
			state.Label = types.StringValue(ds.Label)
			state.URN = types.StringValue(ds.Urn)
			state.Description = types.StringValue(ds.Description)
			state.Owner = types.StringValue(ds.Owner)
			state.CreatedAt = types.StringValue(ds.CreatedAt.String())
			break
		}
	}

	// connection, err := r.connectionClient.Get(state.ID.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError("Error Reading NEOS data source connection", "Could not read NEOS  data source connection ID: "+state.ID.ValueString()+": "+err.Error())
	// 	return
	// }

	connectionJson := state.ConnectionJson

	//state.ConnectionJson = types.StringValue(connection)
	state.ConnectionJson = connectionJson

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Info(ctx, "Data source Read Has error")
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dataSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dataSourceResourceModel
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

	item := neos.DataSourcePutRequest{
		Entity: neos.DataSourcePutRequestEntity{
			Name:        plan.Name.String(),
			Label:       plan.Label.String(),
			Description: plan.Description.String(),
		},
	}

	eItem := neos.DataSourcePutRequestEntityInfo{
		Owner:      plan.Owner.String(),
		ContactIds: contacts,
		Links:      links,
	}

	result, err := r.client.Put(ctx, plan.ID.ValueString(), item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating data system",
			"Could not put data system, unexpected error: "+err.Error(),
		)
		return
	}
	//tflog.Info(ctx, fmt.Sprintf("££ Create Post result [%s] [%s] [%s] [%s] [%s] [%s]", result.Identifier, result.Name, result.Urn, result.Description, result.Label, result.CreatedAt.String()))

	infoResult, err := r.client.PutInfo(ctx, plan.ID.ValueString(), eItem)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating data system",
			"Could not put data system, unexpected error: "+err.Error(),
		)
		return
	}

	connectionResult, err := r.connectionClient.Put(ctx, result.Identifier, plan.ConnectionJson.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating data source connection", "Could not create data source connection, unexpected error: "+err.Error())
		return
	}

	//need to get the secret id from the connection?
	ds, err := r.client.GetById(plan.ID.ValueString())

	tflog.Info(ctx, fmt.Sprintf("Connection result %s %s", result.Identifier, connectionResult))

	var secretMap map[string]string
	secretBytes := []byte(plan.SecretJson.ValueString())
	if err := json.Unmarshal(secretBytes, &secretMap); err != nil {
		resp.Diagnostics.AddError("Error unmarshal connection secret", " unexpected error: "+err.Error())
		return
	}

	secret, err := r.secretClient.GetById(ds.SecretIdentifier)

	spr := neos.SecretPutRequest{
		Name: secret.Name,
		Data: secretMap,
	}

	secretResult, err := r.secretClient.Put(ctx, ds.SecretIdentifier, spr)
	if err != nil {
		resp.Diagnostics.AddError("Error creating data source secret", "Could not create data source secret, unexpected error: "+err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Secret put result %s %s", result.Identifier, secretResult.Identifier))

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
	plan.ConnectionJson = types.StringValue(connectionResult)
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
func (r *dataSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var plan dataSourceResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting data system",
			"Could not delete data system, unexpected error: "+err.Error(),
		)
		return
	}

}

func (r *dataSourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)

	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *neos.DataSourceClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	r.client = &client.DataSourceClient
	r.connectionClient = &client.DataSourceConnectionClient
	r.dataSourceSecretClient = &client.DataSourceSecretClient
	r.secretClient = &client.SecretClient

}

func (r *dataSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
