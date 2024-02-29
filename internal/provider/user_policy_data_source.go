package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	neos "github.com/owain-nortal/neos-client-go"
)

func NewUserPolicyDataSource() datasource.DataSource {
	return &userPolicyDataSource{}
}

var (
	_ datasource.DataSource              = &userPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &userPolicyDataSource{}
)

type userPolicyDataSource struct {
	client *neos.PolicyClient
}

func (d *userPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_policy"
}

func (d *userPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	tflog.Info(ctx, "userPolicyDataSource READ")

	var state UserPolicyDataSourceModel

	list, err := d.client.List("", "")
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read User Policy List", err.Error())
		return
	}

	// Map response body to model
	for _, ds := range list.UserPolicies {
		b, err := json.Marshal(ds)
		if err != nil {
			resp.Diagnostics.AddError("Unable to unmarshal policies in List", err.Error())
			return
		}

		policy := string(b)

		userPolicyState := UserPolicyJsonModel{
			ID:     types.StringValue(ds.User),
			Policy: types.StringValue(policy),
		}
		state.UserPolicies = append(state.UserPolicies, userPolicyState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure adds the provider configured client to the data source.
func (d *userPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Info(ctx, "Data source configure")

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected userPolicyDataSource Configure Type", fmt.Sprintf("Expected *neos.PolicyClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	d.client = &client.PolicyClient
}

func (d *userPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"policy": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"policy": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type UserPolicyDataSourceModel struct {
	UserPolicies []UserPolicyJsonModel `tfsdk:"policy"`
}

type UserPolicyJsonModel struct {
	ID     types.String `tfsdk:"id"`
	Policy types.String `tfsdk:"policy"`
}
