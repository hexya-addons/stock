package stock

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/hexya/src/models/types"
	"github.com/hexya-erp/pool/h"
)

func init() {
	h.StockLocation().DeclareModel()
	h.StockLocation().AddSQLConstraint("barcode_company_uniq", "unique (barcode,company_id)", "The barcode for a location must be unique per company !")

	h.StockLocation().Methods().DefaultGet().Extend(
		`DefaultGet`,
		func(rs m.StockLocationSet, fields interface{}) {
			//        res = super(Location, self).default_get(fields)
			//        if 'barcode' in fields and 'barcode' not in res and res.get('complete_name'):
			//            res['barcode'] = res['complete_name']
			//        return res
		})
	h.StockLocation().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{
			String:    "Location Name",
			Required:  true,
			Translate: true,
		},
		"CompleteName": models.CharField{
			String:  "Full Location Name",
			Compute: h.StockLocation().Methods().ComputeCompleteName(),
			Stored:  true,
		},
		"Active": models.BooleanField{
			String:  "Active",
			Default: models.DefaultValue(true),
			Help: "By unchecking the active field, you may hide a location" +
				"without deleting it.",
		},
		"Usage": models.SelectionField{
			Selection: types.Selection{
				"supplier":    "Vendor Location",
				"view":        "View",
				"internal":    "Internal Location",
				"customer":    "Customer Location",
				"inventory":   "Inventory Loss",
				"procurement": "Procurement",
				"production":  "Production",
				"transit":     "Transit Location",
			},
			String:   "Location Type",
			Default:  models.DefaultValue("internal"),
			Index:    true,
			Required: true,
			Help: "* Vendor Location: Virtual location representing the source" +
				"location for products coming from your vendors" +
				"* View: Virtual location used to create a hierarchical" +
				"structures for your warehouse, aggregating its child locations" +
				"; can't directly contain products" +
				"* Internal Location: Physical locations inside your own warehouses," +
				"* Customer Location: Virtual location representing the" +
				"destination location for products sent to your customers" +
				"* Inventory Loss: Virtual location serving as counterpart" +
				"for inventory operations used to correct stock levels (Physical" +
				"inventories)" +
				"* Procurement: Virtual location serving as temporary counterpart" +
				"for procurement operations when the source (vendor or production)" +
				"is not known yet. This location should be empty when the" +
				"procurement scheduler has finished running." +
				"* Production: Virtual counterpart location for production" +
				"operations: this location consumes the raw material and" +
				"produces finished products" +
				"* Transit Location: Counterpart location that should be" +
				"used in inter-companies or inter-warehouses operations",
		},
		"LocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Parent Location",
			Index:         true,
			OnDelete:      `cascade`,
			Help: "The parent location that includes this location. Example" +
				": The 'Dispatch Zone' is the 'Gate 1' parent location.",
		},
		"ChildIds": models.One2ManyField{
			RelationModel: h.StockLocation(),
			ReverseFK:     "",
			String:        "Contains",
		},
		"PartnerId": models.Many2OneField{
			RelationModel: h.Partner(),
			String:        "Owner",
			Help:          "Owner of the location if not internal",
		},
		"Comment": models.TextField{
			String: "Additional Information",
		},
		"Posx": models.IntegerField{
			String:  "Corridor (X)",
			Default: models.DefaultValue(0),
			Help:    "Optional localization details, for information purpose only",
		},
		"Posy": models.IntegerField{
			String:  "Shelves (Y)",
			Default: models.DefaultValue(0),
			Help:    "Optional localization details, for information purpose only",
		},
		"Posz": models.IntegerField{
			String:  "Height (Z)",
			Default: models.DefaultValue(0),
			Help:    "Optional localization details, for information purpose only",
		},
		"ParentLeft": models.IntegerField{
			String: "Left Parent",
			Index:  true,
		},
		"ParentRight": models.IntegerField{
			String: "Right Parent",
			Index:  true,
		},
		"CompanyId": models.Many2OneField{
			RelationModel: h.Company(),
			String:        "Company",
			Default:       func(env models.Environment) interface{} { return env["res.company"]._company_default_get() },
			Index:         true,
			Help:          "Let this field empty if this location is shared between companies",
		},
		"ScrapLocation": models.BooleanField{
			String:  "Is a Scrap Location?",
			Default: models.DefaultValue(false),
			Help: "Check this box to allow using this location to put scrapped/damaged" +
				"goods.",
		},
		"ReturnLocation": models.BooleanField{
			String: "Is a Return Location?",
			Help:   "Check this box to allow using this location as a return location.",
		},
		"RemovalStrategyId": models.Many2OneField{
			RelationModel: h.ProductRemoval(),
			String:        "Removal Strategy",
			Help: "Defines the default method used for suggesting the exact" +
				"location (shelf) where to take the products from, which" +
				"lot etc. for this location. This method can be enforced" +
				"at the product category level, and a fallback is made on" +
				"the parent locations if none is set here.",
		},
		"PutawayStrategyId": models.Many2OneField{
			RelationModel: h.ProductPutaway(),
			String:        "Put Away Strategy",
			Help: "Defines the default method used for suggesting the exact" +
				"location (shelf) where to store the products. This method" +
				"can be enforced at the product category level, and a fallback" +
				"is made on the parent locations if none is set here.",
		},
		"Barcode": models.CharField{
			String: "Barcode",
			NoCopy: true,
			//oldname='loc_barcode'
		},
	})
	h.StockLocation().Methods().ComputeCompleteName().DeclareMethod(
		` Forms complete name of location from parent location to
child location. `,
		func(rs h.StockLocationSet) h.StockLocationData {
			//        if self.location_id.complete_name:
			//            self.complete_name = '%s/%s' % (
			//                self.location_id.complete_name, self.name)
			//        else:
			//            self.complete_name = self.name
		})
	h.StockLocation().Methods().NameGet().Extend(
		`NameGet`,
		func(rs m.StockLocationSet) {
			//        ret_list = []
			//        for location in self:
			//            orig_location = location
			//            name = location.name
			//            while location.location_id and location.usage != 'view':
			//                location = location.location_id
			//                name = location.name + "/" + name
			//            ret_list.append((orig_location.id, name))
			//        return ret_list
		})
	h.StockLocation().Methods().GetPutawayStrategy().DeclareMethod(
		` Returns the location where the product has to be put,
if any compliant putaway strategy is found. Otherwise returns None.`,
		func(rs m.StockLocationSet, product interface{}) {
			//        current_location = self
			//        putaway_location = self.env['stock.location']
			//        while current_location and not putaway_location:
			//            if current_location.putaway_strategy_id:
			//                putaway_location = current_location.putaway_strategy_id.putaway_apply(
			//                    product)
			//            current_location = current_location.location_id
			//        return putaway_location
		})
	h.StockLocation().Methods().GetWarehouse().DeclareMethod(
		` Returns warehouse id of warehouse that contains location `,
		func(rs m.StockLocationSet) {
			//        return self.env['stock.warehouse'].search([
			//            ('view_location_id.parent_left', '<=', self.parent_left),
			//            ('view_location_id.parent_right', '>=', self.parent_left)], limit=1)
		})
	h.StockLocationRoute().DeclareModel()

	h.StockLocationRoute().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{
			String:    "Route Name",
			Required:  true,
			Translate: true,
		},
		"Active": models.BooleanField{
			String:  "Active",
			Default: models.DefaultValue(true),
			Help: "If the active field is set to False, it will allow you" +
				"to hide the route without removing it.",
		},
		"Sequence": models.IntegerField{
			String:  "Sequence",
			Default: models.DefaultValue(0),
		},
		"PullIds": models.One2ManyField{
			RelationModel: h.ProcurementRule(),
			ReverseFK:     "",
			String:        "Procurement Rules",
			NoCopy:        false,
		},
		"PushIds": models.One2ManyField{
			RelationModel: h.StockLocationPath(),
			ReverseFK:     "",
			String:        "Push Rules",
			NoCopy:        false,
		},
		"ProductSelectable": models.BooleanField{
			String:  "Applicable on Product",
			Default: models.DefaultValue(true),
			Help: "When checked, the route will be selectable in the Inventory" +
				"tab of the Product form.  It will take priority over the" +
				"Warehouse route. ",
		},
		"ProductCategSelectable": models.BooleanField{
			String: "Applicable on Product Category",
			Help: "When checked, the route will be selectable on the Product" +
				"Category.  It will take priority over the Warehouse route. ",
		},
		"WarehouseSelectable": models.BooleanField{
			String: "Applicable on Warehouse",
			Help: "When a warehouse is selected for this route, this route" +
				"should be seen as the default route when products pass" +
				"through this warehouse.  This behaviour can be overridden" +
				"by the routes on the Product/Product Categories or by the" +
				"Preferred Routes on the Procurement",
		},
		"SuppliedWhId": models.Many2OneField{
			RelationModel: h.StockWarehouse(),
			String:        "Supplied Warehouse",
		},
		"SupplierWhId": models.Many2OneField{
			RelationModel: h.StockWarehouse(),
			String:        "Supplying Warehouse",
		},
		"CompanyId": models.Many2OneField{
			RelationModel: h.Company(),
			String:        "Company",
			Default:       func(env models.Environment) interface{} { return env["res.company"]._company_default_get() },
			Index:         true,
			Help:          "Leave this field empty if this route is shared between all companies",
		},
		"ProductIds": models.Many2ManyField{
			RelationModel:    h.ProductTemplate(),
			M2MLinkModelName: "",
			M2MOurField:      "",
			M2MTheirField:    "",
			String:           "Products",
		},
		"CategIds": models.Many2ManyField{
			RelationModel:    h.ProductCategory(),
			M2MLinkModelName: "",
			M2MOurField:      "",
			M2MTheirField:    "",
			String:           "Product Categories",
		},
		"WarehouseIds": models.Many2ManyField{
			RelationModel:    h.StockWarehouse(),
			M2MLinkModelName: "",
			M2MOurField:      "",
			M2MTheirField:    "",
			String:           "Warehouses",
		},
	})
	h.StockLocationRoute().Methods().Write().Extend(
		`when a route is deactivated, deactivate also its pull and push rules`,
		func(rs m.StockLocationRouteSet, values models.RecordData) {
			//        res = super(Route, self).write(values)
			//        if 'active' in values:
			//            self.mapped('push_ids').filtered(lambda path: path.active !=
			//                                             values['active']).write({'active': values['active']})
			//            self.mapped('pull_ids').filtered(lambda rule: rule.active !=
			//                                             values['active']).write({'active': values['active']})
			//        return res
		})
	h.StockLocationRoute().Methods().ViewProductIds().DeclareMethod(
		`ViewProductIds`,
		func(rs m.StockLocationRouteSet) {
			//        return {
			//            'name': _('Products'),
			//            'view_type': 'form',
			//            'view_mode': 'tree,form',
			//            'res_model': 'product.template',
			//            'type': 'ir.actions.act_window',
			//            'domain': [('route_ids', 'in', self.ids)],
			//        }
		})
	h.StockLocationRoute().Methods().ViewCategIds().DeclareMethod(
		`ViewCategIds`,
		func(rs m.StockLocationRouteSet) {
			//        return {
			//            'name': _('Product Categories'),
			//            'view_type': 'form',
			//            'view_mode': 'tree,form',
			//            'res_model': 'product.category',
			//            'type': 'ir.actions.act_window',
			//            'domain': [('route_ids', 'in', self.ids)],
			//        }
		})
	h.StockLocationPath().DeclareModel()

	h.StockLocationPath().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{
			String:   "Operation Name",
			Required: true,
		},
		"CompanyId": models.Many2OneField{
			RelationModel: h.Company(),
			String:        "Company",
			Default:       func(env models.Environment) interface{} { return env["res.company"]._company_default_get() },
			Index:         true,
		},
		"RouteId": models.Many2OneField{
			RelationModel: h.StockLocationRoute(),
			String:        "Route",
		},
		"LocationFromId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Source Location",
			Index:         true,
			OnDelete:      `cascade`,
			Required:      true,
			Help: "This rule can be applied when a move is confirmed that" +
				"has this location as destination location",
		},
		"LocationDestId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Destination Location",
			Index:         true,
			OnDelete:      `cascade`,
			Required:      true,
			Help:          "The new location where the goods need to go",
		},
		"Delay": models.IntegerField{
			String:  "Delay (days)",
			Default: models.DefaultValue(0),
			Help:    "Number of days needed to transfer the goods",
		},
		"PickingTypeId": models.Many2OneField{
			RelationModel: h.StockPickingType(),
			String:        "Picking Type",
			Required:      true,
			Help:          "This is the picking type that will be put on the stock moves",
		},
		"Auto": models.SelectionField{
			Selection: types.Selection{
				"manual":      "Manual Operation",
				"transparent": "Automatic No Step Added",
			},
			String:   "Automatic Move",
			Default:  models.DefaultValue("manual"),
			Index:    true,
			Required: true,
			Help: "The 'Manual Operation' value will create a stock move after" +
				"the current one.With 'Automatic No Step Added', the location" +
				"is replaced in the original move.",
		},
		"Propagate": models.BooleanField{
			String:  "Propagate cancel and split",
			Default: models.DefaultValue(true),
			Help: "If checked, when the previous move is cancelled or split," +
				"the move generated by this move will too",
		},
		"Active": models.BooleanField{
			String:  "Active",
			Default: models.DefaultValue(true),
		},
		"WarehouseId": models.Many2OneField{
			RelationModel: h.StockWarehouse(),
			String:        "Warehouse",
		},
		"RouteSequence": models.IntegerField{
			String:  "Route Sequence",
			Related: `RouteId.Sequence`,
			Stored:  true,
		},
		"Sequence": models.IntegerField{
			String: "Sequence",
		},
	})
	h.StockLocationPath().Methods().Apply().DeclareMethod(
		`Apply`,
		func(rs m.StockLocationPathSet, move interface{}) {
			//        new_date = (datetime.strptime(move.date_expected, DEFAULT_SERVER_DATETIME_FORMAT) +
			//                    relativedelta.relativedelta(days=self.delay)).strftime(DEFAULT_SERVER_DATETIME_FORMAT)
			//        if self.auto == 'transparent':
			//            move.write({
			//                'date': new_date,
			//                'date_expected': new_date,
			//                'location_dest_id': self.location_dest_id.id})
			//            # avoid looping if a push rule is not well configured; otherwise call again push_apply to see if a next step is defined
			//            if self.location_dest_id != move.location_dest_id:
			//                # TDE FIXME: should probably be done in the move model IMO
			//                move._push_apply()
			//        else:
			//            new_move_vals = self._prepare_move_copy_values(move, new_date)
			//            new_move = move.copy(new_move_vals)
			//            move.write({'move_dest_id': new_move.id})
			//            new_move.action_confirm()
		})
	h.StockLocationPath().Methods().PrepareMoveCopyValues().DeclareMethod(
		`PrepareMoveCopyValues`,
		func(rs m.StockLocationPathSet, move_to_copy interface{}, new_date interface{}) {
			//        new_move_vals = {
			//            'origin': move_to_copy.origin or move_to_copy.picking_id.name or "/",
			//            'location_id': move_to_copy.location_dest_id.id,
			//            'location_dest_id': self.location_dest_id.id,
			//            'date': new_date,
			//            'date_expected': new_date,
			//            'company_id': self.company_id.id,
			//            'picking_id': False,
			//            'picking_type_id': self.picking_type_id.id,
			//            'propagate': self.propagate,
			//            'push_rule_id': self.id,
			//            'warehouse_id': self.warehouse_id.id,
			//            'procurement_id': False,
			//        }
			//        return new_move_vals
		})
}
