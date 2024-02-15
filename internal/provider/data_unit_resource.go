package provider

import (
	"context"
	"fmt"

	"encoding/json"
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
// 	_ resource.Resource = &dataUnitResource{}
// )

// New data unitResource is a helper function to simplify the provider implementation.
func NewDataUnitResource() resource.Resource {
	return &dataUnitResource{}
}

// dataUnitResource is the resource implementation.
type dataUnitResource struct {
	client *neos.DataUnitClient
}

var (
	_ resource.Resource                = &dataUnitResource{}
	_ resource.ResourceWithConfigure   = &dataUnitResource{}
	_ resource.ResourceWithImportState = &dataUnitResource{}
)

// Metadata returns the resource type name.
func (r *dataUnitResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_unit"
}

func (r *dataUnitResource) JSONRemarshal(bytes []byte) ([]byte, error) {
	var ifce interface{}
	err := json.Unmarshal(bytes, &ifce)
	if err != nil {
		return nil, err
	}
	return json.Marshal(ifce)
}

// Schema defines the schema for the resource.
func (r *dataUnitResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The Unique ID of the data unit",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"urn": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The URN of the data unit which is read only",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "Name of the data unit",
			},
			"description": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "Description of thedata unit",
			},
			"label": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "Label for the data unit",
			},
			"owner": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "The owner of the data unit",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Required:    false,
				Description: "when the data unit was created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"config_json": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "json that describes the configuration of the data unit",
			},

			// "config": schema.SingleNestedAttribute{
			// 	Computed:    false,
			// 	Optional:    true,
			// 	Required:    false,
			// 	Description: "configuration of the data unit",
			// 	// NestedObject: schema.NestedAttributeObject{
			// 	Attributes: map[string]schema.Attribute{

			// 		"data_unit_type": schema.StringAttribute{
			// 			Computed:    false,
			// 			Optional:    false,
			// 			Required:    true,
			// 			Description: "the data unit type",
			// 		},

			// 		"query": schema.SingleNestedAttribute{
			// 			Computed:    false,
			// 			Required:    false,
			// 			Optional:    true,
			// 			Description: "query configuration for the data unit",
			// 			//NestedObject: schema.NestedAttributeObject{
			// 			Attributes: map[string]schema.Attribute{
			// 				"query": schema.StringAttribute{
			// 					Computed:    false,
			// 					Optional:    false,
			// 					Required:    true,
			// 					Description: "The query to execute",
			// 				},
			// 				//},
			// 			},
			// 		},

			// 		"parquet": schema.SingleNestedAttribute{
			// 			Computed:    false,
			// 			Required:    false,
			// 			Optional:    true,
			// 			Description: "parquet configuration for the data unit",
			// 			//NestedObject: schema.NestedAttributeObject{
			// 			Attributes: map[string]schema.Attribute{
			// 				// "query": schema.StringAttribute{
			// 				// 	Computed:    false,
			// 				// 	Optional:    false,
			// 				// 	Required:    true,
			// 				// 	Description: "The query to execute",
			// 				// },
			// 				//},
			// 			},
			// 		},
			// 		"table": schema.SingleNestedAttribute{
			// 			Computed:    false,
			// 			Required:    false,
			// 			Optional:    true,
			// 			Description: "query configuration for the data unit",
			// 			//NestedObject: schema.NestedAttributeObject{
			// 			Attributes: map[string]schema.Attribute{
			// 				"table": schema.StringAttribute{
			// 					Computed:    false,
			// 					Optional:    false,
			// 					Required:    true,
			// 					Description: "The table name to use",
			// 				},
			// 				//	},
			// 			},
			// 		},
			// 		"data_product": schema.SingleNestedAttribute{
			// 			Computed:    false,
			// 			Required:    false,
			// 			Optional:    true,
			// 			Description: "data product configuration for a data unit",
			// 			//NestedObject: schema.NestedAttributeObject{
			// 			Attributes: map[string]schema.Attribute{
			// 				"engine": schema.StringAttribute{
			// 					Computed:    false,
			// 					Optional:    false,
			// 					Required:    true,
			// 					Description: "The engine to use",
			// 				},
			// 				"table": schema.StringAttribute{
			// 					Computed:    false,
			// 					Optional:    false,
			// 					Required:    true,
			// 					Description: "The table name to use",
			// 				},
			// 				//	},
			// 			},
			// 		},
			// 		"csv": schema.SingleNestedAttribute{
			// 			Computed:    false,
			// 			Required:    false,
			// 			Optional:    true,
			// 			Description: "csv configuration for a data unit",
			// 			//NestedObject: schema.NestedAttributeObject{
			// 			Attributes: map[string]schema.Attribute{
			// 				"path": schema.StringAttribute{
			// 					Computed:    false,
			// 					Optional:    false,
			// 					Required:    true,
			// 					Description: "The engine to use",
			// 				},
			// 				"has_header": schema.BoolAttribute{
			// 					Computed:    false,
			// 					Optional:    false,
			// 					Required:    true,
			// 					Description: "if the csv has a header",
			// 				},
			// 				"delimiter": schema.StringAttribute{
			// 					Computed:    false,
			// 					Optional:    false,
			// 					Required:    true,
			// 					Description: "The delimiter",
			// 				},
			// 				"quote_char": schema.StringAttribute{
			// 					Computed:    false,
			// 					Optional:    true,
			// 					Required:    false,
			// 					Description: "The quote_char",
			// 				},
			// 				"escape_char": schema.StringAttribute{
			// 					Computed:    false,
			// 					Optional:    true,
			// 					Required:    false,
			// 					Description: "The escape_char",
			// 				},
			// 				//	},
			// 				// },
			// 			},
			// 		},
			// 	},
			// },

			"contact_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "list of contacts Ids",
			},
			"links": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "list of links",
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// dataUnitResourceModel maps the resource schema data.
type dataUnitResourceModel struct {
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
	ConfigJson  types.String `tfsdk:"config_json"`
	//Config      dataUnitConfigModel `tfsdk:"config"`
}

// type dataUnitConfigModel struct {
// 	DataUnitType types.String                    `tfsdk:"data_unit_type"`
// 	Query        *dataUnitConfigQueryModel       `tfsdk:"query"`
// 	Table        *dataUnitConfigTableModel       `tfsdk:"table"`
// 	Csv          *dataUnitConfigCsvModel         `tfsdk:"csv"`
// 	DataProduct  *dataUnitConfigDataProductModel `tfsdk:"data_product"`
// 	Parquet      *dataUnitConfigParquetModel     `tfsdk:"parquet"`
// }

// type dataUnitConfigQueryModel struct {
// 	Query types.String `tfsdk:"query"`
// }

// type dataUnitConfigTableModel struct {
// 	Table types.String `tfsdk:"table"`
// }

// type dataUnitConfigParquetModel struct {
// }

// type dataUnitConfigDataProductModel struct {
// 	Engine types.String `tfsdk:"engine"`
// 	Table  types.String `tfsdk:"table"`
// }

// type dataUnitConfigCsvModel struct {
// 	Path       types.String `tfsdk:"path"`
// 	HasHeader  types.Bool   `tfsdk:"has_header"`
// 	Delimiter  types.String `tfsdk:"delimiter"`
// 	QuoteChar  types.String `tfsdk:"quote_char"`
// 	EscapeChar types.String `tfsdk:"escape_char"`
// }

// Create a new resource.
func (r *dataUnitResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//tflog.Info(ctx, "££ Create Get plan")
	// Retrieve values from plan
	var plan dataUnitResourceModel
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

	item := neos.DataUnitPostRequest{
		Entity: neos.DataUnitPostRequestEntity{
			Name:        plan.Name.String(),
			Label:       plan.Label.String(),
			Description: plan.Description.String(),
		},
		EntityInfo: neos.DataUnitPostRequestEntityInfo{
			Owner:      plan.Owner.String(),
			ContactIds: contacts,
			Links:      links,
		},
	}

	//	tflog.Info(ctx, fmt.Sprintf("££ Create Post request [%s] [%s] [%s] [%s]", plan.ID, plan.Name, plan.Label, plan.Description))

	result, err := r.client.Post(ctx, item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating   data unit",
			"Could not create   data unit, unexpected error: "+err.Error(),
		)
		return
	}
	//	tflog.Info(ctx, fmt.Sprintf("££ Create Post result [%s] [%s] [%s] [%s] [%s] [%s]", result.Identifier, result.Name, result.Urn, result.Description, result.Label, result.CreatedAt.String()))
	id := result.Identifier
	plan.ID = types.StringValue(id)
	plan.Name = types.StringValue(result.Name)
	plan.URN = types.StringValue(result.Urn)
	plan.Description = types.StringValue(result.Description)
	plan.Label = types.StringValue(result.Label)
	plan.CreatedAt = types.StringValue(result.CreatedAt.String())
	// dut := dataUnitConfigTableModel{
	// 	 Table:  types.StringValue("123"),
	// }
	// plan.Config.Table = &dut
	// config

	configJson := plan.ConfigJson.ValueString()
	tflog.Info(ctx, fmt.Sprintf("%s", configJson))

	// ordered, err := r.JSONRemarshal([]byte(configJson))
	// if err != nil {
	// 	resp.Diagnostics.AddError("Error ordering data unit config ", "unexpected error: "+err.Error())
	// 	return
	// }


	

	_, err = r.client.ConfigPut(ctx, result.Identifier, configJson)
	if err != nil {
		resp.Diagnostics.AddError("Error updating data unit config ", "Could not create data unit config, unexpected error: "+err.Error())
		return
	}

	// var res map[string]map[string]interface{}
	// json.Unmarshal(dd, &res)

	plan.ConfigJson = types.StringValue(configJson)
	//tflog.Info(ctx, fmt.Sprintf("%s", dd.Configuration))


	// ordered, err := r.JSONRemarshal([]byte(configJson))
	// if err != nil {
	// 	resp.Diagnostics.AddError("Error ordering data unit config ", "unexpected error: "+err.Error())
	// 	return
	// }
	// dd, err := r.client.ConfigPutBase(ctx, id, ordered)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Error creating data unit config ", "Could not create data unit config, unexpected error: "+err.Error())
	// 	return
	// }

	// plan.ConfigJson = types.StringValue(string(ordered))

	// var res map[string]map[string]interface{}
	// json.Unmarshal(dd, &res)

	// tflog.Info(ctx, fmt.Sprintf("%s", res["configuration"]["data_unit_type"]))

	//dataUnitType := plan.Config.DataUnitType.ValueString()
	//tflog.Info(ctx, dataUnitType)

	//tflog.Info(ctx, fmt.Sprintf("Config %v", plan.Config))

	//tflog.Info(ctx, fmt.Sprintf("Switch [%s] [%s] [%v]", "query", dataUnitType, (dataUnitType == "query")))

	// switch dataUnitType {
	// case "query":
	// 	tflog.Info(ctx, fmt.Sprintf("In query"))
	// 	request := neos.DataUnitConfigurationQueryPutRequest{
	// 		Configuration: neos.DataUnitConfigurationQueryConfigPutRequest{
	// 			DateUnitType: dataUnitType,
	// 			Query:        plan.Config.Query.Query.String(),
	// 		},
	// 	}

	// 	tflog.Info(ctx, fmt.Sprintf("Request %v", request))

	// 	dd, err := r.client.ConfigQueryPut(ctx, id, request)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError("Error creating data unit config - query", "Could not create data unit config - query, unexpected error: "+err.Error())
	// 		return
	// 	}
	// 	plan.Config.Query = &dataUnitConfigQueryModel{
	// 		Query: types.StringValue(dd.Configuration.Query),
	// 	}
	// case "data_product":
	// 	request := neos.DataUnitConfigurationDataProductPutRequest{
	// 		Configuration: neos.DataUnitConfigurationDataProductConfigPutRequest{
	// 			DateUnitType: dataUnitType,
	// 			Engine:       plan.Config.DataProduct.Engine.String(),
	// 			Table:        plan.Config.DataProduct.Table.String(),
	// 		},
	// 	}
	// 	dd, err := r.client.ConfigDataProductPut(ctx, id, request)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError("Error creating data unit config - data_product", "Could not create data unit config - product, unexpected error: "+err.Error())
	// 		return
	// 	}
	// 	plan.Config.DataProduct = &dataUnitConfigDataProductModel{
	// 		Engine: types.StringValue(dd.Configuration.Engine),
	// 		Table:  types.StringValue(dd.Configuration.Table),
	// 	}
	// case "table":
	// 	request := neos.DataUnitConfigurationTablePutRequest{
	// 		Configuration: neos.DataUnitConfigurationTableConfigPutRequest{
	// 			DateUnitType: dataUnitType,
	// 			Table:        plan.Config.Table.Table.String(),
	// 		},
	// 	}
	// 	dd, err := r.client.ConfigTablePut(ctx, id, request)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError("Error creating data unit config - table", "Could not create data unit config - table, unexpected error: "+err.Error())
	// 		return
	// 	}
	// 	plan.Config.Table = &dataUnitConfigTableModel{
	// 		Table: types.StringValue(dd.Configuration.Table),
	// 	}
	// case "parquet":
	// 	request := neos.DataUnitConfigurationParquetPutRequest{
	// 		Configuration: neos.DataUnitConfigurationParquetConfigPutRequest{
	// 			DateUnitType: dataUnitType,
	// 		},
	// 	}
	// 	_, err := r.client.ConfigParquetPut(ctx, id, request)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError("Error creating data unit config - table", "Could not create data unit config - parquet, unexpected error: "+err.Error())
	// 		return
	// 	}
	// 	plan.Config.Parquet = &dataUnitConfigParquetModel{}
	// case "csv":
	// 	request := neos.DataUnitConfigurationCSVPutRequest{
	// 		Configuration: neos.DataUnitConfigurationCSVConfigPutRequest{
	// 			DateUnitType: dataUnitType,
	// 			Delimiter:    plan.Config.Csv.Delimiter.String(),
	// 			Path:         plan.Config.Csv.Path.String(),
	// 			HasHeader:    plan.Config.Csv.HasHeader.ValueBool(),
	// 			EscapeChar:   plan.Config.Csv.EscapeChar.String(),
	// 			QuoteChar:    plan.Config.Csv.QuoteChar.String(),
	// 		},
	// 	}
	// 	dd, err := r.client.ConfigCSVPut(ctx, id, request)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError("Error creating data unit config - csv", "Could not create data unit config - csv, unexpected error: "+err.Error())
	// 		return
	// 	}
	// 	plan.Config.Csv = &dataUnitConfigCsvModel{
	// 		Path:       types.StringValue(dd.Configuration.Path),
	// 		Delimiter:  types.StringValue(dd.Configuration.Delimiter),
	// 		EscapeChar: types.StringValue(dd.Configuration.EscapeChar),
	// 		QuoteChar:  types.StringValue(dd.Configuration.QuoteChar),
	// 		HasHeader:  types.BoolValue(dd.Configuration.HasHeader),
	// 	}
	// }

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
func (r *dataUnitResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state

	//	tflog.Info(ctx, "££ READ Get current state")

	var state dataUnitResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	foo := fmt.Sprintf("ID [%s]  Desc [%s]", state.ID.ValueString(), state.Description.ValueString())
	tflog.Info(ctx, foo)

	dataUnitList, err := r.client.Get()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading NEOS data unit", "Could not read NEOS  data unit ID "+state.ID.ValueString()+": "+err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("££ READ iterate over list looking for: %s", state.ID.ValueString()))
	for _, ds := range dataUnitList.Entities {
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

			// dataUnitConfig, err := r.client.ConfigGetBase(ctx, ds.Identifier)
			// if err != nil {
			// 	resp.Diagnostics.AddError(
			// 		"Error Reading NEOS data unit config",
			// 		"Could not read NEOS  data unit ID "+state.ID.ValueString()+": "+err.Error(),
			// 	)
			// 	return
			// }

			// ordered, err := r.JSONRemarshal([]byte(dataUnitConfig))
			// if err != nil {
			// 	resp.Diagnostics.AddError("Error ordering data unit config ", "unexpected error: "+err.Error())
			// 	return
			// }

			// state.ConfigJson = types.StringValue(string(ordered))

			break
		}
	}

	// state.Config.Table = &dataUnitConfigTableModel{
	// 	Table: types.StringValue("xyz"),
	// }
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
func (r *dataUnitResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dataUnitResourceModel
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

	item := neos.DataUnitPutRequest{
		Entity: neos.DataUnitPutRequestEntity{
			Name:        plan.Name.String(),
			Label:       plan.Label.String(),
			Description: plan.Description.String(),
		},
	}

	eItem := neos.DataUnitPutRequestEntityInfo{
		Owner:      plan.Owner.String(),
		ContactIds: contacts,
		Links:      links,
	}

	result, err := r.client.Put(ctx, plan.ID.ValueString(), item)
	if err != nil {
		resp.Diagnostics.AddError("Error updating data unit", "Could not put data unit, unexpected error: "+err.Error())
		return
	}
	//tflog.Info(ctx, fmt.Sprintf("££ Create Post result [%s] [%s] [%s] [%s] [%s] [%s]", result.Identifier, result.Name, result.Urn, result.Description, result.Label, result.CreatedAt.String()))

	infoResult, err := r.client.PutInfo(ctx, plan.ID.ValueString(), eItem)
	if err != nil {
		resp.Diagnostics.AddError("Error updating data unit", "Could not put data unit, unexpected error: "+err.Error())
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

	configJson := plan.ConfigJson.ValueString()
	tflog.Info(ctx, fmt.Sprintf("%s", configJson))

	// ordered, err := r.JSONRemarshal([]byte(configJson))
	// if err != nil {
	// 	resp.Diagnostics.AddError("Error ordering data unit config ", "unexpected error: "+err.Error())
	// 	return
	// }

	_, err = r.client.ConfigPut(ctx, result.Identifier, configJson)
	if err != nil {
		resp.Diagnostics.AddError("Error updating data unit config ", "Could not update data unit config, unexpected error: "+err.Error())
		return
	}

	// var res map[string]map[string]interface{}
	// json.Unmarshal(dd, &res)

	plan.ConfigJson = types.StringValue(configJson)
	//tflog.Info(ctx, fmt.Sprintf("%s", dd.Configuration))

	// plan.Config.Table = &dataUnitConfigTableModel{
	// 	Table: types.StringValue("upd"),
	// }

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dataUnitResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var plan dataUnitResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting data unit", "Could not delete data unit, unexpected error: "+err.Error())
		return
	}

}

func (r *dataUnitResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)

	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *neos.DataUnitClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	r.client = &client.DataUnitClient

}

func (r *dataUnitResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
