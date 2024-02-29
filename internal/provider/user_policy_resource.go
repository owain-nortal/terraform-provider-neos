package provider

import (
	"context"
	"fmt"

	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	jt "github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/owain-nortal/neos-client-go"
	//"strings"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
// var (
// 	_ resource.Resource = &userPolicyResource{}
// )

// New userPolicyResource is a helper function to simplify the provider implementation.
func NewUserPolicyResource() resource.Resource {
	return &userPolicyResource{}
}

// userPolicyResource is the resource implementation.
type userPolicyResource struct {
	client *neos.PolicyClient
}

var (
	_ resource.Resource                = &userPolicyResource{}
	_ resource.ResourceWithConfigure   = &userPolicyResource{}
	_ resource.ResourceWithImportState = &userPolicyResource{}
)

// Metadata returns the resource type name.
func (r *userPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_policy"
}

func (r *userPolicyResource) JSONRemarshal(bytes []byte) ([]byte, error) {
	var ifce interface{}
	err := json.Unmarshal(bytes, &ifce)
	if err != nil {
		return nil, err
	}
	return json.Marshal(ifce)
}

// Schema defines the schema for the resource.
func (r *userPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "The Unique ID of the policy",
			},
			"policy_json": schema.StringAttribute{
				CustomType:  jt.NormalizedType{},
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "the Policy",
			},
			// "account": schema.StringAttribute{
			// 	Computed:    true,
			// 	Optional:    false,
			// 	Required:    false,
			// 	Description: "account",
			// },
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// userPolicyResourceModel maps the resource schema data.
type userPolicyResourceModel struct {
	ID          types.String  `tfsdk:"id"`
	Policy      jt.Normalized `tfsdk:"policy_json"`
	LastUpdated types.String  `tfsdk:"last_updated"`
	//Account     types.String `tfsdk:"account"`
}

// Create a new resource.
func (r *userPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve values from plan
	var plan userPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ss := plan.Policy.ValueString()

	item := neos.PolicyPostRequest{
		Policy: ss,
	}

	_, err := r.client.Post(ctx, item, r.client.Account)
	if err != nil {
		resp.Diagnostics.AddError("Error creating policy", "Could not create policy, unexpected error: "+err.Error())
		return
	}

	normailisedPolicy, err := r.client.NormalizeJson(plan.Policy.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error normailisedPolicy", "Could not normailise Policy, unexpected error: "+err.Error())
		return
	}

	// id := result.
	//plan.ID = types.StringValue(id)
	plan.Policy = jt.NewNormalizedValue(normailisedPolicy)
	//plan.Policy = types.StringValue(normailisedPolicy)
	//plan.Account = types.StringValue(r.client.Account)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *userPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state

	var state userPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: need to not hard code the paritition to ksa
	//nrn := fmt.Sprintf("nrn:ksa:iam::%s:user:%s", state.Account.ValueString(), )

	userPolicy, err := r.client.Get(state.ID.ValueString(), r.client.Account)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading NEOS userPolicy", "Could not read NEOS  userPolicy ID "+state.ID.ValueString()+": "+err.Error())
		return
	}

	//	x , diag :=  jt.NormalizedType.ValueFromString(jt.NormalizedType{},ctx,)

	getPolicy := jt.NewNormalizedValue(userPolicy.Policy)
	eq, _ := getPolicy.StringSemanticEquals(ctx, state.Policy)
	if !eq {
		state.Policy = getPolicy
	}

	//state.Policy = types.StringValue(userPolicy.Policy)

	// if dd.HasError() {
	// 	return
	// }

	//state.Account = types.StringValue(r.client.Account)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Info(ctx, "Data system Read Has error")
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *userPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ss := plan.Policy.ValueString()

	dspr := neos.PolicyPutRequest{
		Policy: ss,
	}

	_, err := r.client.Put(ctx, plan.ID.ValueString(), dspr, r.client.Account)
	if err != nil {
		resp.Diagnostics.AddError("Error updating userPolicy", "Could not put userPolicy, unexpected error: "+err.Error())
		return
	}

	plan.Policy = jt.NewNormalizedValue(ss) // types.StringValue(result.Policy)
	// _, dd := plan.Policy.ValueFromString(ctx, types.StringValue(result.Policy))
	// if dd.HasError() {
	// 	return
	// }

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var plan userPolicyResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, plan.ID.ValueString(), r.client.Account)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting userPolicy", "Could not delete userPolicy, unexpected error: "+err.Error())
		return
	}
}

func (r *userPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*neos.NeosClient)

	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *neos.PolicyClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	r.client = &client.PolicyClient

}

func (r *userPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
