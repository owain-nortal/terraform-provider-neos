package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	neos "github.com/owain-nortal/neos-client-go"
	"golang.org/x/exp/slices"
	"time"
)

// New groupResource is a helper function to simplify the provider implementation.
func NewGroupResource() resource.Resource {
	return &groupResource{}
}

// groupResource is the resource implementation.
type groupResource struct {
	client *neos.GroupClient
}

var (
	_ resource.Resource                = &groupResource{}
	_ resource.ResourceWithConfigure   = &groupResource{}
	_ resource.ResourceWithImportState = &groupResource{}
)

// Metadata returns the resource type name.
func (r *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

// Schema defines the schema for the resource.
func (r *groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The Unique ID of the group",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "Name of the group",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "Description of the group",
			},
			"principals": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "list of principals",
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
			"account": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "account if not root",
			},
		},
	}
}

// groupResourceModel maps the resource schema data.
type groupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Principals  types.Set    `tfsdk:"principals"`
	IsSystem    types.Bool   `tfsdk:"is_system"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Account     types.String `tfsdk:"account"`
}

func SortListValueIntoStringArray(ctx context.Context, lv basetypes.SetValue) (diag.Diagnostics, []string) {
	rtn := []string{}
	diags := lv.ElementsAs(ctx, &rtn, false)
	return diags, rtn
}

func SortStringArrayToList(input []string) (basetypes.SetValue, diag.Diagnostics) {
	rtn := []attr.Value{}
	for _, v := range input {
		rtn = append(rtn, types.StringValue(v))
	}

	return types.SetValue(types.StringType, rtn)
}

// Create a new resource.
func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan groupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := neos.GroupPostRequest{
		Name:        plan.Name.String(),
		Description: plan.Description.String(),
	}

	result, err := r.client.Post(ctx, item, plan.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating   group", "Could not create   group, unexpected error: "+err.Error())
		return
	}

	plan.ID = types.StringValue(result.Identifier)
	plan.Name = types.StringValue(result.Name)
	plan.Description = types.StringValue(result.Description)
	plan.IsSystem = types.BoolValue(result.IsSystem)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	gppr := neos.GroupPrincipalPostRequest{}

	diags, gppr.Principals = SortListValueIntoStringArray(ctx, plan.Principals)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(gppr.Principals) > 0 {
		_, err := r.client.PrincipalsPost(ctx, plan.ID.ValueString(), gppr, plan.Account.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error creating   group Principals", "Could not create group Principals, unexpected error: "+err.Error())
			return
		}
	}

	plan.Principals, diags = SortStringArrayToList(gppr.Principals)
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

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var state groupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	foo := fmt.Sprintf("ID [%s]  Desc [%s]", state.ID.ValueString(), state.Description.ValueString())
	tflog.Info(ctx, foo)

	groupList, err := r.client.List(state.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading NEOS group", "Could not read NEOS  group ID "+state.ID.ValueString()+": "+err.Error())
		return
	}

	for _, ds := range groupList.Groups {
		if ds.Identifier == state.ID.ValueString() {
			state.ID = types.StringValue(ds.Identifier)
			state.Name = types.StringValue(ds.Name)
			state.Description = types.StringValue(ds.Description)
			state.IsSystem = types.BoolValue(ds.IsSystem)
			state.Principals, diags = SortStringArrayToList(ds.Principals)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				tflog.Info(ctx, "group Read Has error")
				return
			}
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
func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan groupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := neos.GroupPutRequest{
		Name:        plan.Name.String(),
		Description: plan.Description.String(),
	}

	result, err := r.client.Put(ctx, plan.ID.ValueString(), item, plan.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error updating group", "Could not put group, unexpected error: "+err.Error())
		return
	}

	plan.ID = types.StringValue(result.Identifier)
	plan.Name = types.StringValue(result.Name)
	plan.Description = types.StringValue(result.Description)
	plan.IsSystem = types.BoolValue(result.IsSystem)

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	grp, err := r.client.Get(plan.ID.ValueString(),plan.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error gettting group", "Could not get group to workout principals, unexpected error: "+err.Error())
		return
	}

	shouldReturn := r.addUsersToGroup(ctx, plan, resp, grp)
	if shouldReturn {
		return
	}

	shouldReturn = r.deleteUsersFromGroup(ctx, plan, resp, grp)
	if shouldReturn {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupResource) deleteUsersFromGroup(ctx context.Context, plan groupResourceModel, resp *resource.UpdateResponse, grp neos.Group) bool {
	delList := []string{}
	tvl, diags := plan.Principals.ToSetValue(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return true
	}
	diags, planPrincipals := SortListValueIntoStringArray(ctx, tvl)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return true
	}

	for _, gP := range grp.Principals {
		if slices.Contains(planPrincipals, gP) {

		} else {

			delList = append(delList, gP)
		}
	}

	g, err := r.client.PrincipalsDelete(ctx, plan.ID.ValueString(), neos.GroupPrincipalDeleteRequest{Principals: delList}, plan.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting principals", "Could not delete principals, unexpected error: "+err.Error())
		return true
	}

	plan.Principals, diags = basetypes.NewSetValueFrom(ctx, types.StringType, g.Principals)
	resp.Diagnostics.Append(diags...)
	return resp.Diagnostics.HasError()
}

func (r *groupResource) addUsersToGroup(ctx context.Context, plan groupResourceModel, resp *resource.UpdateResponse, grp neos.Group) bool {
	addList := []string{}
	tvl, diags := plan.Principals.ToSetValue(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return true
	}
	diags, planPrincipals := SortListValueIntoStringArray(ctx, tvl)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return true
	}
	for _, planP := range planPrincipals {
		if slices.Contains(grp.Principals, planP) {

		} else {

			addList = append(grp.Principals, planP)
		}
	}

	gg, err := r.client.PrincipalsPost(ctx, plan.ID.ValueString(), neos.GroupPrincipalPostRequest{Principals: addList}, plan.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error posting principals", "Could not post to update principals, unexpected error: "+err.Error())
		return true
	}

	tflog.Info(ctx, fmt.Sprintf("ID: %s ", gg.Identifier))

	return false
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var plan groupResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, plan.ID.ValueString(), plan.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting group", "Could not delete group, unexpected error: "+err.Error())
		return
	}

}

func (r *groupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)

	if !ok {
		resp.Diagnostics.AddError("Unexpected Group Configure Type", fmt.Sprintf("Expected *neos.GroupClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	r.client = &client.GroupClient

}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
