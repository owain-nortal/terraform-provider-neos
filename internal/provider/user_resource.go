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
// 	_ resource.Resource = &userResource{}
// )

// New userResource is a helper function to simplify the provider implementation.
func NewUserResource() resource.Resource {
	return &userResource{}
}

// userResource is the resource implementation.
type userResource struct {
	client *neos.UserClient
}

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

// Metadata returns the resource type name.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) JSONRemarshal(bytes []byte) ([]byte, error) {
	var ifce interface{}
	err := json.Unmarshal(bytes, &ifce)
	if err != nil {
		return nil, err
	}
	return json.Marshal(ifce)
}

// "identifier": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
// "urn": "string",
// "first_name": "string",
// "last_name": "string",
// "email": "string",
// "username": "string",
// "is_system": true,
// "enabled": true

// Schema defines the schema for the resource.
func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The Unique ID of the user",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"urn": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The URN of the user which is read only",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"first_name": schema.StringAttribute{
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "First name of the user",
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.RequiresReplace(),
				// },
			},
			"last_name": schema.StringAttribute{
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "Last name of the user",
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.RequiresReplace(),
				// },
			},
			"username": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "Username of the user",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"email": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "Email of the user",
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.RequiresReplace(),
				// },
			},
			"enabled": schema.BoolAttribute{
				Computed:    false,
				Optional:    false,
				Required:    true,
				Description: "user enabled",
			},
			"account": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "account if not root",
			},
			"is_system": schema.BoolAttribute{
				Computed:    true,
				Optional:    false,
				Required:    false,
				Description: "The owner of the user",
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// userResourceModel maps the resource schema data.
type userResourceModel struct {
	ID          types.String `tfsdk:"id"`
	URN         types.String `tfsdk:"urn"`
	FirstName   types.String `tfsdk:"first_name"`
	LastName    types.String `tfsdk:"last_name"`
	Email       types.String `tfsdk:"email"`
	Username    types.String `tfsdk:"username"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	IsSystem    types.Bool   `tfsdk:"is_system"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Account     types.String `tfsdk:"account"`
}

// Create a new resource.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve values from plan
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//tflog.Info(ctx, "££ After Create Get plan")

	item := neos.UserPostRequest{
		Username:  plan.Username.String(),
		LastName:  plan.LastName.String(),
		FirstName: plan.FirstName.String(),
		Email:     plan.Email.String(),
	}

	result, err := r.client.Post(ctx, item, plan.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", "Could not create user, unexpected error: "+err.Error())
		return
	}
	id := result.Identifier
	plan.ID = types.StringValue(id)
	plan.LastName = types.StringValue(result.LastName)
	plan.URN = types.StringValue(result.Urn)
	plan.FirstName = types.StringValue(result.FirstName)
	plan.Email = types.StringValue(result.Email)
	plan.Username = types.StringValue(result.Username)
	plan.Enabled = types.BoolValue(result.Enabled)
	plan.IsSystem = types.BoolValue(result.IsSystem)
	plan.Account = types.StringValue(result.Account)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state

	//	tflog.Info(ctx, "££ READ Get current state")

	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userList, err := r.client.List("", "", state.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading NEOS user", "Could not read NEOS  user ID "+state.ID.ValueString()+": "+err.Error())
		return
	}

	for _, ds := range userList.Users {
		if ds.Identifier == state.ID.ValueString() {
			//bits := strings.Split(ds.Urn, ":")
			//account := bits[4]
			state.ID = types.StringValue(ds.Identifier)
			state.FirstName = types.StringValue(ds.FirstName)
			state.LastName = types.StringValue(ds.LastName)
			state.Email = types.StringValue(ds.Email)
			state.Username = types.StringValue(ds.Username)
			state.Enabled = types.BoolValue(ds.Enabled)
			state.IsSystem = types.BoolValue(ds.IsSystem)
			state.URN = types.StringValue(ds.Urn)
			state.Account = types.StringValue(state.Account.ValueString())
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
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dspr := neos.UserPostRequest{
		Username:  plan.Username.String(),
		LastName:  plan.LastName.String(),
		FirstName: plan.FirstName.String(),
		Email:     plan.Email.String(),
	}

	result, err := r.client.Post(ctx, dspr, plan.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error updating user", "Could not put user, unexpected error: "+err.Error())
		return
	}

	plan.ID = types.StringValue(result.Identifier)
	plan.FirstName = types.StringValue(result.FirstName)
	plan.URN = types.StringValue(result.Urn)
	plan.LastName = types.StringValue(result.LastName)
	plan.Email = types.StringValue(result.Email)
	plan.Account = types.StringValue(result.Account)
	plan.Username = types.StringValue(result.Username)
	plan.Enabled = types.BoolValue(result.Enabled)
	plan.IsSystem = types.BoolValue(result.IsSystem)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var plan userResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, plan.ID.ValueString(), plan.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting user", "Could not delete user, unexpected error: "+err.Error())
		return
	}

}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)

	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *neos.UserClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	r.client = &client.UserClient

}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
