package stock

import (
	"github.com/hexya-erp/hexya-base/web/webdata"
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/hexya/src/models/types"
	"github.com/hexya-erp/hexya/src/models/types/dates"
	"github.com/hexya-erp/pool/h"
	"github.com/hexya-erp/pool/q"
)

//import json
//import time
func init() {
	h.StockPickingType().DeclareModel()

	h.StockPickingType().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{
			String:    "Picking Type Name",
			Required:  true,
			Translate: true,
		},
		"Color": models.IntegerField{
			String: "Color",
		},
		"Sequence": models.IntegerField{
			String: "Sequence",
			Help:   "Used to order the 'All Operations' kanban view",
		},
		"SequenceId": models.Many2OneField{
			RelationModel: h.IrSequence(),
			String:        "Reference Sequence",
			Required:      true,
		},
		"DefaultLocationSrcId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Default Source Location",
			Help: "This is the default source location when you create a picking" +
				"manually with this picking type. It is possible however" +
				"to change it or that the routes put another location. If" +
				"it is empty, it will check for the supplier location on the partner. ",
		},
		"DefaultLocationDestId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Default Destination Location",
			Help: "This is the default destination location when you create" +
				"a picking manually with this picking type. It is possible" +
				"however to change it or that the routes put another location." +
				"If it is empty, it will check for the customer location on the partner. ",
		},
		"Code": models.SelectionField{
			Selection: types.Selection{
				"incoming": "Vendors",
				"outgoing": "Customers",
				"internal": "Internal",
			},
			String:   "Type of Operation",
			Required: true,
		},
		"ReturnPickingTypeId": models.Many2OneField{
			RelationModel: h.StockPickingType(),
			String:        "Picking Type for Returns",
		},
		"ShowEntirePacks": models.BooleanField{
			String: "Allow moving packs",
			Help: "If checked, this shows the packs to be moved as a whole" +
				"in the Operations tab all the time, even if there was no" +
				"entire pack reserved.",
		},
		"WarehouseId": models.Many2OneField{
			RelationModel: h.StockWarehouse(),
			String:        "Warehouse",
			OnDelete:      `cascade`,
			Default:       func(env models.Environment) interface{} { return env["stock.warehouse"].search() },
		},
		"Active": models.BooleanField{
			String:  "Active",
			Default: models.DefaultValue(true),
		},
		"UseCreateLots": models.BooleanField{
			String:  "Create New Lots/Serial Numbers",
			Default: models.DefaultValue(true),
			Help: "If this is checked only, it will suppose you want to create" +
				"new Lots/Serial Numbers, so you can provide them in a text field. ",
		},
		"UseExistingLots": models.BooleanField{
			String:  "Use Existing Lots/Serial Numbers",
			Default: models.DefaultValue(true),
			Help: "If this is checked, you will be able to choose the Lots/Serial" +
				"Numbers. You can also decide to not put lots in this picking" +
				"type.  This means it will create stock with no lot or not" +
				"put a restriction on the lot taken. ",
		},
		"LastDonePicking": models.CharField{
			String:  "Last 10 Done Pickings",
			Compute: h.StockPickingType().Methods().ComputeLastDonePicking(),
		},
		"CountPickingDraft": models.IntegerField{
			Compute: h.StockPickingType().Methods().ComputePickingCount(),
		},
		"CountPickingReady": models.IntegerField{
			Compute: h.StockPickingType().Methods().ComputePickingCount(),
		},
		"CountPicking": models.IntegerField{
			Compute: h.StockPickingType().Methods().ComputePickingCount(),
		},
		"CountPickingWaiting": models.IntegerField{
			Compute: h.StockPickingType().Methods().ComputePickingCount(),
		},
		"CountPickingLate": models.IntegerField{
			Compute: h.StockPickingType().Methods().ComputePickingCount(),
		},
		"CountPickingBackorders": models.IntegerField{
			Compute: h.StockPickingType().Methods().ComputePickingCount(),
		},
		"RatePickingLate": models.IntegerField{
			Compute: h.StockPickingType().Methods().ComputePickingCount(),
		},
		"RatePickingBackorders": models.IntegerField{
			Compute: h.StockPickingType().Methods().ComputePickingCount(),
		},
		"BarcodeNomenclatureId": models.Many2OneField{
			RelationModel: h.BarcodeNomenclature(),
			String:        "Barcode Nomenclature",
		},
	})
	h.StockPickingType().Methods().ComputeLastDonePicking().DeclareMethod(
		`ComputeLastDonePicking`,
		func(rs h.StockPickingTypeSet) h.StockPickingTypeData {
			//        tristates = []
			//        for picking in self.env['stock.picking'].search([('picking_type_id', '=', self.id), ('state', '=', 'done')], order='date_done desc', limit=10):
			//            if picking.date_done > picking.date:
			//                tristates.insert(
			//                    0, {'tooltip': picking.name or '' + ": " + _('Late'), 'value': -1})
			//            elif picking.backorder_id:
			//                tristates.insert(
			//                    0, {'tooltip': picking.name or '' + ": " + _('Backorder exists'), 'value': 0})
			//            else:
			//                tristates.insert(
			//                    0, {'tooltip': picking.name or '' + ": " + _('OK'), 'value': 1})
			//        self.last_done_picking = json.dumps(tristates)
		})
	h.StockPickingType().Methods().ComputePickingCount().DeclareMethod(
		`ComputePickingCount`,
		func(rs h.StockPickingTypeSet) h.StockPickingTypeData {
			//        domains = {
			//            'count_picking_draft': [('state', '=', 'draft')],
			//            'count_picking_waiting': [('state', 'in', ('confirmed', 'waiting'))],
			//            'count_picking_ready': [('state', 'in', ('assigned', 'partially_available'))],
			//            'count_picking': [('state', 'in', ('assigned', 'waiting', 'confirmed', 'partially_available'))],
			//            'count_picking_late': [('min_date', '<', time.strftime(DEFAULT_SERVER_DATETIME_FORMAT)), ('state', 'in', ('assigned', 'waiting', 'confirmed', 'partially_available'))],
			//            'count_picking_backorders': [('backorder_id', '!=', False), ('state', 'in', ('confirmed', 'assigned', 'waiting', 'partially_available'))],
			//        }
			//        for field in domains:
			//            data = self.env['stock.picking'].read_group(domains[field] +
			//                                                        [('state', 'not in', ('done', 'cancel')),
			//                                                         ('picking_type_id', 'in', self.ids)],
			//                                                        ['picking_type_id'], ['picking_type_id'])
			//            count = dict(map(lambda x: (
			//                x['picking_type_id'] and x['picking_type_id'][0], x['picking_type_id_count']), data))
			//            for record in self:
			//                record[field] = count.get(record.id, 0)
			//        for record in self:
			//            record.rate_picking_late = record.count_picking and record.count_picking_late * \
			//                100 / record.count_picking or 0
			//            record.rate_picking_backorders = record.count_picking and record.count_picking_backorders * \
			//                100 / record.count_picking or 0
		})
	h.StockPickingType().Methods().NameGet().Extend(
		` Display 'Warehouse_name: PickingType_name' `,
		func(rs m.StockPickingTypeSet) {
			//        res = []
			//        for picking_type in self:
			//            if self.env.context.get('special_shortened_wh_name'):
			//                if picking_type.warehouse_id:
			//                    name = picking_type.warehouse_id.name
			//                else:
			//                    name = _('Customer') + ' (' + picking_type.name + ')'
			//            elif picking_type.warehouse_id:
			//                name = picking_type.warehouse_id.name + ': ' + picking_type.name
			//            else:
			//                name = picking_type.name
			//            res.append((picking_type.id, name))
			//        return res
		})
	h.StockPickingType().Methods().NameSearch().Extend(
		`NameSearch`,
		func(rs m.StockPickingTypeSet, name webdata.NameSearchParams, args interface{}, operator interface{}, limit interface{}) {
			//        args = args or []
			//        domain = []
			//        if name:
			//            domain = ['|', ('name', operator, name),
			//                      ('warehouse_id.name', operator, name)]
			//        picks = self.search(domain + args, limit=limit)
			//        return picks.name_get()
		})
	h.StockPickingType().Methods().OnchangePickingCode().DeclareMethod(
		`OnchangePickingCode`,
		func(rs m.StockPickingTypeSet) {
			//        if self.code == 'incoming':
			//            self.default_location_src_id = self.env.ref(
			//                'stock.stock_location_suppliers').id
			//            self.default_location_dest_id = self.env.ref(
			//                'stock.stock_location_stock').id
			//        elif self.code == 'outgoing':
			//            self.default_location_src_id = self.env.ref(
			//                'stock.stock_location_stock').id
			//            self.default_location_dest_id = self.env.ref(
			//                'stock.stock_location_customers').id
		})
	h.StockPickingType().Methods().GetAction().DeclareMethod(
		`GetAction`,
		func(rs m.StockPickingTypeSet, action_xmlid interface{}) {
			//        action = self.env.ref(action_xmlid).read()[0]
			//        if self:
			//            action['display_name'] = self.display_name
			//        return action
		})
	h.StockPickingType().Methods().GetActionPickingTreeLate().DeclareMethod(
		`GetActionPickingTreeLate`,
		func(rs m.StockPickingTypeSet) {
			//        return self._get_action('stock.action_picking_tree_late')
		})
	h.StockPickingType().Methods().GetActionPickingTreeBackorder().DeclareMethod(
		`GetActionPickingTreeBackorder`,
		func(rs m.StockPickingTypeSet) {
			//        return self._get_action('stock.action_picking_tree_backorder')
		})
	h.StockPickingType().Methods().GetActionPickingTreeWaiting().DeclareMethod(
		`GetActionPickingTreeWaiting`,
		func(rs m.StockPickingTypeSet) {
			//        return self._get_action('stock.action_picking_tree_waiting')
		})
	h.StockPickingType().Methods().GetActionPickingTreeReady().DeclareMethod(
		`GetActionPickingTreeReady`,
		func(rs m.StockPickingTypeSet) {
			//        return self._get_action('stock.action_picking_tree_ready')
		})
	h.StockPickingType().Methods().GetStockPickingActionPickingType().DeclareMethod(
		`GetStockPickingActionPickingType`,
		func(rs m.StockPickingTypeSet) {
			//        return self._get_action('stock.stock_picking_action_picking_type')
		})
	h.StockPicking().DeclareModel()
	h.StockPicking().AddSQLConstraint("name_uniq", "unique(name, company_id)", "Reference must be unique per company!")

	h.StockPicking().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{
			String:  "Reference",
			Default: models.DefaultValue("/"),
			NoCopy:  true,
			Index:   true,
			//states={'done': [('readonly', True)], 'cancel': [('readonly', True)]}
		},
		"Origin": models.CharField{
			String: "Source Document",
			Index:  true,
			//states={'done': [('readonly', True)], 'cancel': [('readonly', True)]}
			Help: "Reference of the document",
		},
		"Note": models.TextField{
			String: "Notes",
		},
		"BackorderId": models.Many2OneField{
			RelationModel: h.StockPicking(),
			String:        "Back Order of",
			NoCopy:        true,
			Index:         true,
			//states={'done': [('readonly', True)], 'cancel': [('readonly', True)]}
			Help: "If this shipment was split, then this field links to the" +
				"shipment which contains the already processed part.",
		},
		"MoveType": models.SelectionField{
			Selection: types.Selection{
				"direct": "Partial",
				"one":    "All at once",
			},
			String:   "Delivery Type",
			Default:  models.DefaultValue("direct"),
			Required: true,
			//states={'done': [('readonly', True)], 'cancel': [('readonly', True)]}
			Help: "It specifies goods to be deliver partially or all at once",
		},
		"State": models.SelectionField{
			Selection: types.Selection{
				"draft":               "Draft",
				"cancel":              "Cancelled",
				"waiting":             "Waiting Another Operation",
				"confirmed":           "Waiting Availability",
				"partially_available": "Partially Available",
				"assigned":            "Available",
				"done":                "Done",
			},
			String:   "Status",
			Compute:  h.StockPicking().Methods().ComputeState(),
			NoCopy:   true,
			Index:    true,
			ReadOnly: true,
			Stored:   true,
			//track_visibility='onchange'
			Help: " * Draft: not confirmed yet and will not be scheduled until confirmed" +
				" * Waiting Another Operation: waiting for another move" +
				"to proceed before it becomes automatically available (e.g." +
				"in Make-To-Order flows)" +
				" * Waiting Availability: still waiting for the availability of products" +
				" * Partially Available: some products are available and reserved" +
				" * Ready to Transfer: products reserved, simply waiting" +
				"for confirmation." +
				" * Transferred: has been processed, can't be modified or" +
				"cancelled anymore" +
				" * Cancelled: has been cancelled, can't be confirmed anymore",
		},
		"GroupId": models.Many2OneField{
			RelationModel: h.ProcurementGroup(),
			String:        "Procurement Group",
			ReadOnly:      true,
			Related:       `MoveLines.GroupId`,
			Stored:        true,
		},
		"Priority": models.SelectionField{
			Selection: procurement.PROCUREMENT_PRIORITIES,
			String:    "Priority",
			Compute:   h.StockPicking().Methods().ComputePriority(),
			//inverse='_set_priority'
			Stored: true,
			Index:  true,
			//track_visibility='onchange'
			//states={'done': [('readonly', True)], 'cancel': [('readonly', True)]}
			Help: "Priority for this picking. Setting manually a value here" +
				"would set it as priority for all the moves",
		},
		"MinDate": models.DateTimeField{
			String:  "Scheduled Date",
			Compute: h.StockPicking().Methods().ComputeDates(),
			//inverse='_set_min_date'
			Stored: true,
			Index:  true,
			//track_visibility='onchange'
			//states={'done': [('readonly', True)], 'cancel': [('readonly', True)]}
			Help: "Scheduled time for the first part of the shipment to be" +
				"processed. Setting manually a value here would set it as" +
				"expected date for all the stock moves.",
		},
		"MaxDate": models.DateTimeField{
			String:  "Max. Expected Date",
			Compute: h.StockPicking().Methods().ComputeDates(),
			Stored:  true,
			Index:   true,
			Help:    "Scheduled time for the last part of the shipment to be processed",
		},
		"Date": models.DateTimeField{
			String:  "Creation Date",
			Default: func(env models.Environment) interface{} { return dates.Now() },
			Index:   true,
			//track_visibility='onchange'
			//states={'done': [('readonly', True)], 'cancel': [('readonly', True)]}
			Help: "Creation Date, usually the time of the order",
		},
		"DateDone": models.DateTimeField{
			String:   "Date of Transfer",
			NoCopy:   true,
			ReadOnly: true,
			Help:     "Completion Date of Transfer",
		},
		"LocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Source Location Zone",
			Default: func(env models.Environment) interface{} {
				return env["stock.picking.type"].browse().default_location_src_id
			},
			ReadOnly: true,
			Required: true,
			//states={'draft': [('readonly', False)]}
		},
		"LocationDestId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Destination Location Zone",
			Default: func(env models.Environment) interface{} {
				return env["stock.picking.type"].browse().default_location_dest_id
			},
			ReadOnly: true,
			Required: true,
			//states={'draft': [('readonly', False)]}
		},
		"MoveLines": models.One2ManyField{
			RelationModel: h.StockMove(),
			ReverseFK:     "",
			String:        "Stock Moves",
			NoCopy:        false,
		},
		"HasScrapMove": models.BooleanField{
			String:  "Has Scrap Moves",
			Compute: h.StockPicking().Methods().HasScrapMove(),
		},
		"PickingTypeId": models.Many2OneField{
			RelationModel: h.StockPickingType(),
			String:        "Picking Type",
			Required:      true,
			ReadOnly:      true,
			//states={'draft': [('readonly', False)]}
		},
		"PickingTypeCode": models.SelectionField{
			Selection: types.Selection{
				"incoming": "Vendors",
				"outgoing": "Customers",
				"internal": "Internal",
			},
			Related:  `PickingTypeId.Code`,
			ReadOnly: true,
		},
		"PickingTypeEntirePacks": models.BooleanField{
			Related:  `PickingTypeId.ShowEntirePacks`,
			ReadOnly: true,
		},
		"QuantReservedExist": models.BooleanField{
			String:  "Has quants already reserved",
			Compute: h.StockPicking().Methods().ComputeQuantReservedExist(),
			Help:    "Check the existance of quants linked to this picking",
		},
		"PartnerId": models.Many2OneField{
			RelationModel: h.Partner(),
			String:        "Partner",
			//states={'done': [('readonly', True)], 'cancel': [('readonly', True)]}
		},
		"CompanyId": models.Many2OneField{
			RelationModel: h.Company(),
			String:        "Company",
			Default:       func(env models.Environment) interface{} { return env["res.company"]._company_default_get() },
			Index:         true,
			Required:      true,
			//states={'done': [('readonly', True)], 'cancel': [('readonly', True)]}
		},
		"PackOperationIds": models.One2ManyField{
			RelationModel: h.StockPackOperation(),
			ReverseFK:     "",
			String:        "Related Packing Operations",
			//states={'done': [('readonly', True)], 'cancel': [('readonly', True)]}
		},
		"PackOperationProductIds": models.One2ManyField{
			RelationModel: h.StockPackOperation(),
			ReverseFK:     "",
			String:        "Non pack",
			Filter:        q.ProductId().NotEquals(False),
			//states={'done': [('readonly', True)], 'cancel': [('readonly', True)]}
		},
		"PackOperationPackIds": models.One2ManyField{
			RelationModel: h.StockPackOperation(),
			ReverseFK:     "",
			String:        "Pack",
			Filter:        q.ProductId().Equals(False),
			//states={'done': [('readonly', True)], 'cancel': [('readonly', True)]}
		},
		"PackOperationExist": models.BooleanField{
			String:  "Has Pack Operations",
			Compute: h.StockPicking().Methods().ComputePackOperationExist(),
			Help:    "Check the existence of pack operation on the picking",
		},
		"OwnerId": models.Many2OneField{
			RelationModel: h.Partner(),
			String:        "Owner",
			//states={'done': [('readonly', True)], 'cancel': [('readonly', True)]}
			Help: "Default Owner",
		},
		"Printed": models.BooleanField{
			String: "Printed",
		},
		"ProductId": models.Many2OneField{
			RelationModel: h.ProductProduct(),
			String:        "Product",
			Related:       `MoveLines.ProductId`,
		},
		"RecomputePackOp": models.BooleanField{
			String: "Recompute pack operation?",
			NoCopy: true,
			Help: "True if reserved quants changed, which mean we might need" +
				"to recompute the package operations",
		},
		"LaunchPackOperations": models.BooleanField{
			String: "Launch Pack Operations",
			NoCopy: true,
		},
	})
	h.StockPicking().Methods().ComputeState().DeclareMethod(
		` State of a picking depends on the state of its related stock.move
         - no moves: draft or assigned (launch_pack_operations)
         - all moves canceled: cancel
         - all moves done (including possible canceled): done
         - All at once picking: least of confirmed / waiting / assigned
         - Partial picking
          - all moves assigned: assigned
          - one of the move is assigned or partially available:
partially available
          - otherwise in waiting or confirmed state
        `,
		func(rs h.StockPickingSet) h.StockPickingData {
			//        if not self.move_lines and self.launch_pack_operations:
			//            self.state = 'assigned'
			//        elif not self.move_lines:
			//            self.state = 'draft'
			//        elif any(move.state == 'draft' for move in self.move_lines):  # TDE FIXME: should be all ?
			//            self.state = 'draft'
			//        elif all(move.state == 'cancel' for move in self.move_lines):
			//            self.state = 'cancel'
			//        elif all(move.state in ['cancel', 'done'] for move in self.move_lines):
			//            self.state = 'done'
			//        else:
			//            # We sort our moves by importance of state: "confirmed" should be first, then we'll have
			//            # "waiting" and finally "assigned" at the end.
			//            moves_todo = self.move_lines\
			//                .filtered(lambda move: move.state not in ['cancel', 'done'])\
			//                .sorted(key=lambda move: (move.state == 'assigned' and 2) or (move.state == 'waiting' and 1) or 0)
			//            if self.move_type == 'one':
			//                self.state = moves_todo[0].state or 'draft'
			//            elif moves_todo[0].state != 'assigned' and any(x.partially_available or x.state == 'assigned' for x in moves_todo):
			//                self.state = 'partially_available'
			//            else:
			//                self.state = moves_todo[-1].state or 'draft'
		})
	h.StockPicking().Methods().ComputePriority().DeclareMethod(
		`ComputePriority`,
		func(rs h.StockPickingSet) h.StockPickingData {
			//        self.priority = self.mapped('move_lines') and max(
			//            self.mapped('move_lines').mapped('priority')) or '1'
		})
	h.StockPicking().Methods().SetPriority().DeclareMethod(
		`SetPriority`,
		func(rs m.StockPickingSet) {
			//        self.move_lines.write({'priority': self.priority})
		})
	h.StockPicking().Methods().ComputeDates().DeclareMethod(
		`ComputeDates`,
		func(rs h.StockPickingSet) h.StockPickingData {
			//        self.min_date = min(self.move_lines.mapped('date_expected') or [False])
			//        self.max_date = max(self.move_lines.mapped('date_expected') or [False])
		})
	h.StockPicking().Methods().SetMinDate().DeclareMethod(
		`SetMinDate`,
		func(rs m.StockPickingSet) {
			//        self.move_lines.write({'date_expected': self.min_date})
		})
	h.StockPicking().Methods().HasScrapMove().DeclareMethod(
		`HasScrapMove`,
		func(rs h.StockPickingSet) h.StockPickingData {
			//        self.has_scrap_move = bool(self.env['stock.move'].search_count(
			//            [('picking_id', '=', self.id), ('scrapped', '=', True)]))
		})
	h.StockPicking().Methods().ComputeQuantReservedExist().DeclareMethod(
		`ComputeQuantReservedExist`,
		func(rs h.StockPickingSet) h.StockPickingData {
			//        self.quant_reserved_exist = any(
			//            move.reserved_quant_ids for move in self.mapped('move_lines'))
		})
	h.StockPicking().Methods().ComputePackOperationExist().DeclareMethod(
		`ComputePackOperationExist`,
		func(rs h.StockPickingSet) h.StockPickingData {
			//        self.pack_operation_exist = bool(self.pack_operation_ids)
		})
	h.StockPicking().Methods().OnchangePickingType().DeclareMethod(
		`OnchangePickingType`,
		func(rs m.StockPickingSet) {
			//        if self.picking_type_id:
			//            if self.picking_type_id.default_location_src_id:
			//                location_id = self.picking_type_id.default_location_src_id.id
			//            elif self.partner_id:
			//                location_id = self.partner_id.property_stock_supplier.id
			//            else:
			//                customerloc, location_id = self.env['stock.warehouse']._get_partner_locations(
			//                )
			//
			//            if self.picking_type_id.default_location_dest_id:
			//                location_dest_id = self.picking_type_id.default_location_dest_id.id
			//            elif self.partner_id:
			//                location_dest_id = self.partner_id.property_stock_customer.id
			//            else:
			//                location_dest_id, supplierloc = self.env['stock.warehouse']._get_partner_locations(
			//                )
			//
			//            self.location_id = location_id
			//            self.location_dest_id = location_dest_id
			//        if self.partner_id:
			//            if self.partner_id.picking_warn == 'no-message' and self.partner_id.parent_id:
			//                partner = self.partner_id.parent_id
			//            elif self.partner_id.picking_warn not in ('no-message', 'block') and self.partner_id.parent_id.picking_warn == 'block':
			//                partner = self.partner_id.parent_id
			//            else:
			//                partner = self.partner_id
			//            if partner.picking_warn != 'no-message':
			//                if partner.picking_warn == 'block':
			//                    self.partner_id = False
			//                return {'warning': {
			//                    'title': ("Warning for %s") % partner.name,
			//                    'message': partner.picking_warn_msg
			//                }}
		})
	h.StockPicking().Methods().Create().Extend(
		`Create`,
		func(rs m.StockPickingSet, vals models.RecordData) {
			//        defaults = self.default_get(['name', 'picking_type_id'])
			//        if vals.get('name', '/') == '/' and defaults.get('name', '/') == '/' and vals.get('picking_type_id', defaults.get('picking_type_id')):
			//            vals['name'] = self.env['stock.picking.type'].browse(vals.get(
			//                'picking_type_id', defaults.get('picking_type_id'))).sequence_id.next_by_id()
			//        if vals.get('move_lines') and vals.get('location_id') and vals.get('location_dest_id'):
			//            for move in vals['move_lines']:
			//                if len(move) == 3:
			//                    move[2]['location_id'] = vals['location_id']
			//                    move[2]['location_dest_id'] = vals['location_dest_id']
			//        return super(Picking, self).create(vals)
		})
	h.StockPicking().Methods().Write().Extend(
		`Write`,
		func(rs m.StockPickingSet, vals models.RecordData) {
			//        res = super(Picking, self).write(vals)
			//        after_vals = {}
			//        if vals.get('location_id'):
			//            after_vals['location_id'] = vals['location_id']
			//        if vals.get('location_dest_id'):
			//            after_vals['location_dest_id'] = vals['location_dest_id']
			//        if after_vals:
			//            self.mapped('move_lines').filtered(
			//                lambda move: not move.scrapped).write(after_vals)
			//        return res
		})
	h.StockPicking().Methods().Unlink().Extend(
		`Unlink`,
		func(rs m.StockPickingSet) {
			//        self.mapped('move_lines').action_cancel()
			//        self.mapped('move_lines').unlink()  # Checks if moves are not done
			//        return super(Picking, self).unlink()
		})
	h.StockPicking().Methods().ActionAssignOwner().DeclareMethod(
		`ActionAssignOwner`,
		func(rs m.StockPickingSet) {
			//        self.pack_operation_ids.write({'owner_id': self.owner_id.id})
		})
	h.StockPicking().Methods().DoPrintPicking().DeclareMethod(
		`DoPrintPicking`,
		func(rs m.StockPickingSet) {
			//        self.write({'printed': True})
			//        return self.env["report"].get_action(self, 'stock.report_picking')
		})
	h.StockPicking().Methods().ActionConfirm().DeclareMethod(
		`ActionConfirm`,
		func(rs m.StockPickingSet) {
			//        self.filtered(lambda picking: not picking.move_lines).write(
			//            {'launch_pack_operations': True})
			//        self.mapped('move_lines').filtered(
			//            lambda move: move.state == 'draft').action_confirm()
			//        self.filtered(lambda picking: picking.location_id.usage in (
			//            'supplier', 'inventory', 'production')).force_assign()
			//        return True
		})
	h.StockPicking().Methods().ActionAssign().DeclareMethod(
		` Check availability of picking moves.
        This has the effect of changing the state and reserve
quants on available moves, and may
        also impact the state of the picking as it is computed
based on move's states.
        @return: True
        `,
		func(rs m.StockPickingSet) {
			//        self.filtered(lambda picking: picking.state ==
			//                      'draft').action_confirm()
			//        moves = self.mapped('move_lines').filtered(
			//            lambda move: move.state not in ('draft', 'cancel', 'done'))
			//        if not moves:
			//            raise UserError(_('Nothing to check the availability for.'))
			//        moves.action_assign()
			//        return True
		})
	h.StockPicking().Methods().ForceAssign().DeclareMethod(
		` Changes state of picking to available if moves are confirmed
or waiting.
        @return: True
        `,
		func(rs m.StockPickingSet) {
			//        self.mapped('move_lines').filtered(lambda move: move.state in [
			//            'confirmed', 'waiting']).force_assign()
			//        return True
		})
	h.StockPicking().Methods().ActionCancel().DeclareMethod(
		`ActionCancel`,
		func(rs m.StockPickingSet) {
			//        self.mapped('move_lines').action_cancel()
			//        return True
		})
	h.StockPicking().Methods().ActionDone().DeclareMethod(
		`Changes picking state to done by processing the Stock Moves
of the Picking

        Normally that happens when the button "Done" is
pressed on a Picking view.
        @return: True
        `,
		func(rs m.StockPickingSet) {
			//        draft_moves = self.mapped('move_lines').filtered(
			//            lambda self: self.state == 'draft')
			//        todo_moves = self.mapped('move_lines').filtered(
			//            lambda self: self.state in ['draft', 'assigned', 'confirmed'])
			//        draft_moves.action_confirm()
			//        todo_moves.action_done()
			//        return True
		})
	h.StockPicking().Methods().RecheckAvailability().DeclareMethod(
		`RecheckAvailability`,
		func(rs m.StockPickingSet) {
			//        self.action_assign()
			//        self.do_prepare_partial()
		})
	h.StockPicking().Methods().PreparePackOps().DeclareMethod(
		` Prepare pack_operations, returns a list of dict to give at create `,
		func(rs m.StockPickingSet, quants interface{}, forced_qties interface{}) {
			//        valid_quants = quants.filtered(lambda quant: quant.qty > 0)
			//        _Mapping = namedtuple(
			//            'Mapping', ('product', 'package', 'owner', 'location', 'location_dst_id'))
			//        all_products = valid_quants.mapped('product_id') | self.env['product.product'].browse(
			//            p.id for p in forced_qties.keys()) | self.move_lines.mapped('product_id')
			//        computed_putaway_locations = dict(
			//            (product, self.location_dest_id.get_putaway_strategy(product) or self.location_dest_id.id) for product in all_products)
			//        product_to_uom = dict((product.id, product.uom_id)
			//                              for product in all_products)
			//        picking_moves = self.move_lines.filtered(
			//            lambda move: move.state not in ('done', 'cancel'))
			//        for move in picking_moves:
			//            # If we encounter an UoM that is smaller than the default UoM or the one already chosen, use the new one instead.
			//            if move.product_uom != product_to_uom[move.product_id.id] and move.product_uom.factor > product_to_uom[move.product_id.id].factor:
			//                product_to_uom[move.product_id.id] = move.product_uom
			//        if len(picking_moves.mapped('location_id')) > 1:
			//            raise UserError(
			//                _('The source location must be the same for all the moves of the picking.'))
			//        if len(picking_moves.mapped('location_dest_id')) > 1:
			//            raise UserError(
			//                _('The destination location must be the same for all the moves of the picking.'))
			//        pack_operation_values = []
			//        top_lvl_packages = valid_quants._get_top_level_packages(
			//            computed_putaway_locations)
			//        for pack in top_lvl_packages:
			//            pack_quants = pack.get_content()
			//            pack_operation_values.append({
			//                'picking_id': self.id,
			//                'package_id': pack.id,
			//                'product_qty': 1.0,
			//                'location_id': pack.location_id.id,
			//                'location_dest_id': computed_putaway_locations[pack_quants[0].product_id],
			//                'owner_id': pack.owner_id.id,
			//            })
			//            valid_quants -= pack_quants
			//        qtys_grouped = {}
			//        lots_grouped = {}
			//        for quant in valid_quants:
			//            key = _Mapping(quant.product_id, quant.package_id, quant.owner_id,
			//                           quant.location_id, computed_putaway_locations[quant.product_id])
			//            qtys_grouped.setdefault(key, 0.0)
			//            qtys_grouped[key] += quant.qty
			//            if quant.product_id.tracking != 'none' and quant.lot_id:
			//                lots_grouped.setdefault(key, dict()).setdefault(
			//                    quant.lot_id.id, 0.0)
			//                lots_grouped[key][quant.lot_id.id] += quant.qty
			//        for product, qty in forced_qties.items():
			//            if qty <= 0.0:
			//                continue
			//            key = _Mapping(product, self.env['stock.quant.package'], self.owner_id,
			//                           self.location_id, computed_putaway_locations[product])
			//            qtys_grouped.setdefault(key, 0.0)
			//            qtys_grouped[key] += qty
			//        Uom = self.env['product.uom']
			//        product_id_to_vals = {}
			//        for mapping, qty in qtys_grouped.items():
			//            uom = product_to_uom[mapping.product.id]
			//            val_dict = {
			//                'picking_id': self.id,
			//                'product_qty': mapping.product.uom_id._compute_quantity(qty, uom),
			//                'product_id': mapping.product.id,
			//                'package_id': mapping.package.id,
			//                'owner_id': mapping.owner.id,
			//                'location_id': mapping.location.id,
			//                'location_dest_id': mapping.location_dst_id,
			//                'product_uom_id': uom.id,
			//                'pack_lot_ids': [
			//                    (0, 0, {
			//                        'lot_id': lot,
			//                        'qty': 0.0,
			//                        'qty_todo': mapping.product.uom_id._compute_quantity(lots_grouped[mapping][lot], uom)
			//                    }) for lot in lots_grouped.get(mapping, {}).keys()],
			//            }
			//            product_id_to_vals.setdefault(
			//                mapping.product.id, list()).append(val_dict)
			//        for move in self.move_lines.filtered(lambda move: move.state not in ('done', 'cancel')):
			//            values = product_id_to_vals.pop(move.product_id.id, [])
			//            pack_operation_values += values
			//        return pack_operation_values
		})
	h.StockPicking().Methods().DoPreparePartial().DeclareMethod(
		`DoPreparePartial`,
		func(rs m.StockPickingSet) {
			//        PackOperation = self.env['stock.pack.operation']
			//        existing_packages = PackOperation.search(
			//            [('picking_id', 'in', self.ids)])  # TDE FIXME: o2m / m2o ?
			//        if existing_packages:
			//            existing_packages.unlink()
			//        for picking in self:
			//            forced_qties = {}  # Quantity remaining after calculating reserved quants
			//            picking_quants = self.env['stock.quant']
			//            # Calculate packages, reserved quants, qtys of this picking's moves
			//            for move in picking.move_lines:
			//                if move.state not in ('assigned', 'confirmed', 'waiting'):
			//                    continue
			//                move_quants = move.reserved_quant_ids
			//                picking_quants += move_quants
			//                forced_qty = 0.0
			//                if move.state == 'assigned':
			//                    qty = move.product_uom._compute_quantity(
			//                        move.product_uom_qty, move.product_id.uom_id, round=False)
			//                    forced_qty = qty - sum([x.qty for x in move_quants])
			//                # if we used force_assign() on the move, or if the move is incoming, forced_qty > 0
			//                if float_compare(forced_qty, 0, precision_rounding=move.product_id.uom_id.rounding) > 0:
			//                    if forced_qties.get(move.product_id):
			//                        forced_qties[move.product_id] += forced_qty
			//                    else:
			//                        forced_qties[move.product_id] = forced_qty
			//            for vals in picking._prepare_pack_ops(picking_quants, forced_qties):
			//                vals['fresh_record'] = False
			//                PackOperation |= PackOperation.create(vals)
			//        self.do_recompute_remaining_quantities()
			//        for pack in PackOperation:
			//            pack.ordered_qty = sum(
			//                pack.mapped('linked_move_operation_ids').mapped('move_id').filtered(
			//                    lambda r: r.state != 'cancel').mapped('ordered_qty')
			//            )
			//        self.write({'recompute_pack_op': False})
		})
	h.StockPicking().Methods().DoUnreserve().DeclareMethod(
		`
          Will remove all quants for picking in picking_ids
        `,
		func(rs m.StockPickingSet) {
			//        moves_to_unreserve = self.mapped('move_lines').filtered(
			//            lambda move: move.state not in ('done', 'cancel'))
			//        pack_line_to_unreserve = self.mapped('pack_operation_ids')
			//        if moves_to_unreserve:
			//            if pack_line_to_unreserve:
			//                pack_line_to_unreserve.unlink()
			//            moves_to_unreserve.do_unreserve()
		})
	h.StockPicking().Methods().RecomputeRemainingQty().DeclareMethod(
		`RecomputeRemainingQty`,
		func(rs m.StockPickingSet, done_qtys interface{}) {
			//        def _create_link_for_index(operation_id, index, product_id, qty_to_assign, quant_id=False):
			//            move_dict = prod2move_ids[product_id][index]
			//            qty_on_link = min(move_dict['remaining_qty'], qty_to_assign)
			//            self.env['stock.move.operation.link'].create(
			//                {'move_id': move_dict['move'].id, 'operation_id': operation_id, 'qty': qty_on_link, 'reserved_quant_id': quant_id})
			//            if float_compare(move_dict['remaining_qty'], qty_on_link, precision_rounding=move_dict['move'].product_uom.rounding) == 0:
			//                prod2move_ids[product_id].pop(index)
			//            else:
			//                move_dict['remaining_qty'] -= qty_on_link
			//            return qty_on_link
			//        def _create_link_for_quant(operation_id, quant, qty):
			//            """create a link for given operation and reserved move of given quant, for the max quantity possible, and returns this quantity"""
			//            if not quant.reservation_id.id:
			//                return _create_link_for_product(operation_id, quant.product_id.id, qty)
			//            qty_on_link = 0
			//            for i in range(0, len(prod2move_ids[quant.product_id.id])):
			//                if prod2move_ids[quant.product_id.id][i]['move'].id != quant.reservation_id.id:
			//                    continue
			//                qty_on_link = _create_link_for_index(
			//                    operation_id, i, quant.product_id.id, qty, quant_id=quant.id)
			//                break
			//            return qty_on_link
			//        def _create_link_for_product(operation_id, product_id, qty):
			//            '''method that creates the link between a given operation and move(s) of given product, for the given quantity.
			//            Returns True if it was possible to create links for the requested quantity (False if there was not enough quantity on stock moves)'''
			//            qty_to_assign = qty
			//            Product = self.env["product.product"]
			//            product = Product.browse(product_id)
			//            rounding = product.uom_id.rounding
			//            qtyassign_cmp = float_compare(
			//                qty_to_assign, 0.0, precision_rounding=rounding)
			//            if prod2move_ids.get(product_id):
			//                while prod2move_ids[product_id] and qtyassign_cmp > 0:
			//                    qty_on_link = _create_link_for_index(
			//                        operation_id, 0, product_id, qty_to_assign, quant_id=False)
			//                    qty_to_assign -= qty_on_link
			//                    qtyassign_cmp = float_compare(
			//                        qty_to_assign, 0.0, precision_rounding=rounding)
			//            return qtyassign_cmp == 0
			//        Uom = self.env['product.uom']
			//        QuantPackage = self.env['stock.quant.package']
			//        OperationLink = self.env['stock.move.operation.link']
			//        quants_in_package_done = set()
			//        prod2move_ids = {}
			//        still_to_do = []
			//        moves = sorted([x for x in self.move_lines if x.state not in ('done', 'cancel')], key=lambda x: (
			//            ((x.state == 'assigned') and -2 or 0) + (x.partially_available and -1 or 0)))
			//        for move in moves:
			//            if not prod2move_ids.get(move.product_id.id):
			//                prod2move_ids[move.product_id.id] = [
			//                    {'move': move, 'remaining_qty': move.product_qty}]
			//            else:
			//                prod2move_ids[move.product_id.id].append(
			//                    {'move': move, 'remaining_qty': move.product_qty})
			//        need_rereserve = False
			//        operations = self.pack_operation_ids
			//        operations = sorted(operations, key=lambda x: ((x.package_id and not x.product_id)
			//                                                       and -4 or 0) + (x.package_id and -2 or 0) + (x.pack_lot_ids and -1 or 0))
			//        links = OperationLink.search(
			//            [('operation_id', 'in', [x.id for x in operations])])
			//        if links:
			//            links.unlink()
			//        for ops in operations:
			//            lot_qty = {}
			//            for packlot in ops.pack_lot_ids:
			//                lot_qty[packlot.lot_id.id] = ops.product_uom_id._compute_quantity(
			//                    packlot.qty, ops.product_id.uom_id)
			//            # for each operation, create the links with the stock move by seeking on the matching reserved quants,
			//            # and deffer the operation if there is some ambiguity on the move to select
			//            if ops.package_id and not ops.product_id and (not done_qtys or ops.qty_done):
			//                # entire package
			//                for quant in ops.package_id.get_content():
			//                    remaining_qty_on_quant = quant.qty
			//                    if quant.reservation_id:
			//                        # avoid quants being counted twice
			//                        quants_in_package_done.add(quant.id)
			//                        qty_on_link = _create_link_for_quant(
			//                            ops.id, quant, quant.qty)
			//                        remaining_qty_on_quant -= qty_on_link
			//                    if remaining_qty_on_quant:
			//                        still_to_do.append(
			//                            (ops, quant.product_id.id, remaining_qty_on_quant))
			//                        need_rereserve = True
			//            elif ops.product_id.id:
			//                # Check moves with same product
			//                product_qty = ops.qty_done if done_qtys else ops.product_qty
			//                qty_to_assign = ops.product_uom_id._compute_quantity(
			//                    product_qty, ops.product_id.uom_id)
			//                precision_rounding = ops.product_id.uom_id.rounding
			//                for move_dict in prod2move_ids.get(ops.product_id.id, []):
			//                    move = move_dict['move']
			//                    for quant in move.reserved_quant_ids:
			//                        if float_compare(qty_to_assign, 0, precision_rounding=precision_rounding) != 1:
			//                            break
			//                        if quant.id in quants_in_package_done:
			//                            continue
			//
			//                        # check if the quant is matching the operation details
			//                        if ops.package_id:
			//                            flag = quant.package_id == ops.package_id
			//                        else:
			//                            flag = not quant.package_id.id
			//                        flag = flag and (ops.owner_id.id == quant.owner_id.id) and (
			//                            ops.location_id.id == quant.location_id.id)
			//                        if flag:
			//                            if not lot_qty:
			//                                max_qty_on_link = min(quant.qty, qty_to_assign)
			//                                qty_on_link = _create_link_for_quant(
			//                                    ops.id, quant, max_qty_on_link)
			//                                qty_to_assign -= qty_on_link
			//                            else:
			//                                # if there is still some qty left
			//                                if lot_qty.get(quant.lot_id.id):
			//                                    max_qty_on_link = min(
			//                                        quant.qty, qty_to_assign, lot_qty[quant.lot_id.id])
			//                                    qty_on_link = _create_link_for_quant(
			//                                        ops.id, quant, max_qty_on_link)
			//                                    qty_to_assign -= qty_on_link
			//                                    lot_qty[quant.lot_id.id] -= qty_on_link
			//
			//                qty_assign_cmp = float_compare(
			//                    qty_to_assign, 0, precision_rounding=precision_rounding)
			//                if qty_assign_cmp > 0:
			//                    # qty reserved is less than qty put in operations. We need to create a link but it's deferred after we processed
			//                    # all the quants (because they leave no choice on their related move and needs to be processed with higher priority)
			//                    still_to_do += [(ops, ops.product_id.id, qty_to_assign)]
			//                    need_rereserve = True
			//        all_op_processed = True
			//        for ops, product_id, remaining_qty in still_to_do:
			//            all_op_processed = _create_link_for_product(
			//                ops.id, product_id, remaining_qty) and all_op_processed
			//        return (need_rereserve, all_op_processed)
		})
	h.StockPicking().Methods().PickingRecomputeRemainingQuantities().DeclareMethod(
		`PickingRecomputeRemainingQuantities`,
		func(rs m.StockPickingSet, done_qtys interface{}) {
			//        need_rereserve = False
			//        all_op_processed = True
			//        if self.pack_operation_ids:
			//            need_rereserve, all_op_processed = self.recompute_remaining_qty(
			//                done_qtys=done_qtys)
			//        return need_rereserve, all_op_processed
		})
	h.StockPicking().Methods().DoRecomputeRemainingQuantities().DeclareMethod(
		`DoRecomputeRemainingQuantities`,
		func(rs m.StockPickingSet, done_qtys interface{}) {
			//        tmp = self.filtered(lambda picking: picking.pack_operation_ids)
			//        if tmp:
			//            for pick in tmp:
			//                pick.recompute_remaining_qty(done_qtys=done_qtys)
		})
	h.StockPicking().Methods().RereserveQuants().DeclareMethod(
		` Unreserve quants then try to reassign quants.`,
		func(rs m.StockPickingSet, move_ids interface{}) {
			//        if not move_ids:
			//            self.do_unreserve()
			//            self.action_assign()
			//        else:
			//            moves = self.env['stock.move'].browse(move_ids)
			//            if self.env.context.get('no_state_change'):
			//                moves.filtered(lambda m: m.reserved_quant_ids).do_unreserve()
			//            else:
			//                moves.do_unreserve()
			//            moves.action_assign(no_prepare=True)
		})
	h.StockPicking().Methods().DoNewTransfer().DeclareMethod(
		`DoNewTransfer`,
		func(rs m.StockPickingSet) {
			//        for pick in self:
			//            if pick.state == 'done':
			//                raise UserError(_('The pick is already validated'))
			//            pack_operations_delete = self.env['stock.pack.operation']
			//            if not pick.move_lines and not pick.pack_operation_ids:
			//                raise UserError(
			//                    _('Please create some Initial Demand or Mark as Todo and create some Operations. '))
			//            # In draft or with no pack operations edited yet, ask if we can just do everything
			//            if pick.state == 'draft' or all([x.qty_done == 0.0 for x in pick.pack_operation_ids]):
			//                # If no lots when needed, raise error
			//                picking_type = pick.picking_type_id
			//                if (picking_type.use_create_lots or picking_type.use_existing_lots):
			//                    for pack in pick.pack_operation_ids:
			//                        if pack.product_id and pack.product_id.tracking != 'none':
			//                            raise UserError(
			//                                _('Some products require lots/serial numbers, so you need to specify those first!'))
			//                view = self.env.ref('stock.view_immediate_transfer')
			//                wiz = self.env['stock.immediate.transfer'].create(
			//                    {'pick_id': pick.id})
			//                # TDE FIXME: a return in a loop, what a good idea. Really.
			//                return {
			//                    'name': _('Immediate Transfer?'),
			//                    'type': 'ir.actions.act_window',
			//                    'view_type': 'form',
			//                    'view_mode': 'form',
			//                    'res_model': 'stock.immediate.transfer',
			//                    'views': [(view.id, 'form')],
			//                    'view_id': view.id,
			//                    'target': 'new',
			//                    'res_id': wiz.id,
			//                    'context': self.env.context,
			//                }
			//
			//            # Check backorder should check for other barcodes
			//            if pick.check_backorder():
			//                view = self.env.ref('stock.view_backorder_confirmation')
			//                wiz = self.env['stock.backorder.confirmation'].create(
			//                    {'pick_id': pick.id})
			//                # TDE FIXME: same reamrk as above actually
			//                return {
			//                    'name': _('Create Backorder?'),
			//                    'type': 'ir.actions.act_window',
			//                    'view_type': 'form',
			//                    'view_mode': 'form',
			//                    'res_model': 'stock.backorder.confirmation',
			//                    'views': [(view.id, 'form')],
			//                    'view_id': view.id,
			//                    'target': 'new',
			//                    'res_id': wiz.id,
			//                    'context': self.env.context,
			//                }
			//            for operation in pick.pack_operation_ids:
			//                if operation.qty_done < 0:
			//                    raise UserError(_('No negative quantities allowed'))
			//                if operation.qty_done > 0:
			//                    operation.write({'product_qty': operation.qty_done})
			//                else:
			//                    pack_operations_delete |= operation
			//            if pack_operations_delete:
			//                pack_operations_delete.unlink()
			//        self.do_transfer()
			//        return
		})
	h.StockPicking().Methods().CheckBackorder().DeclareMethod(
		`CheckBackorder`,
		func(rs m.StockPickingSet) {
			//        need_rereserve, all_op_processed = self.picking_recompute_remaining_quantities(
			//            done_qtys=True)
			//        for move in self.move_lines:
			//            if float_compare(move.remaining_qty, 0, precision_rounding=move.product_id.uom_id.rounding) != 0:
			//                return True
			//        return False
		})
	h.StockPicking().Methods().DoTransfer().DeclareMethod(
		` If no pack operation, we do simple action_done of the picking.
        Otherwise, do the pack operations. `,
		func(rs m.StockPickingSet) {
			//        self._create_lots_for_picking()
			//        no_pack_op_pickings = self.filtered(
			//            lambda picking: not picking.pack_operation_ids)
			//        no_pack_op_pickings.action_done()
			//        other_pickings = self - no_pack_op_pickings
			//        for picking in other_pickings:
			//            need_rereserve, all_op_processed = picking.picking_recompute_remaining_quantities()
			//            todo_moves = self.env['stock.move']
			//            toassign_moves = self.env['stock.move']
			//
			//            # create extra moves in the picking (unexpected product moves coming from pack operations)
			//            if not all_op_processed:
			//                todo_moves |= picking._create_extra_moves()
			//
			//            if need_rereserve or not all_op_processed:
			//                moves_reassign = any(
			//                    x.origin_returned_move_id or x.move_orig_ids for x in picking.move_lines if x.state not in ['done', 'cancel'])
			//                if moves_reassign and picking.location_id.usage not in ("supplier", "production", "inventory"):
			//                    # unnecessary to assign other quants than those involved with pack operations as they will be unreserved anyways.
			//                    picking.with_context(reserve_only_ops=True, no_state_change=True).rereserve_quants(
			//                        move_ids=picking.move_lines.ids)
			//                picking.do_recompute_remaining_quantities()
			//
			//            # split move lines if needed
			//            for move in picking.move_lines:
			//                rounding = move.product_id.uom_id.rounding
			//                remaining_qty = move.remaining_qty
			//                if move.state in ('done', 'cancel'):
			//                    # ignore stock moves cancelled or already done
			//                    continue
			//                elif move.state == 'draft':
			//                    toassign_moves |= move
			//                if float_compare(remaining_qty, 0,  precision_rounding=rounding) == 0:
			//                    if move.state in ('draft', 'assigned', 'confirmed'):
			//                        todo_moves |= move
			//                elif float_compare(remaining_qty, 0, precision_rounding=rounding) > 0 and float_compare(remaining_qty, move.product_qty, precision_rounding=rounding) < 0:
			//                    # TDE FIXME: shoudl probably return a move - check for no track key, by the way
			//                    new_move_id = move.split(remaining_qty)
			//                    new_move = self.env['stock.move'].with_context(
			//                        mail_notrack=True).browse(new_move_id)
			//                    todo_moves |= move
			//                    # Assign move as it was assigned before
			//                    toassign_moves |= new_move
			//
			//            # TDE FIXME: do_only_split does not seem used anymore
			//            if todo_moves and not self.env.context.get('do_only_split'):
			//                todo_moves.action_done()
			//            elif self.env.context.get('do_only_split'):
			//                picking = picking.with_context(split=todo_moves.ids)
			//
			//            picking._create_backorder()
			//        return True
		})
	h.StockPicking().Methods().CreateLotsForPicking().DeclareMethod(
		`CreateLotsForPicking`,
		func(rs m.StockPickingSet) {
			//        Lot = self.env['stock.production.lot']
			//        for pack_op_lot in self.mapped('pack_operation_ids').mapped('pack_lot_ids'):
			//            if not pack_op_lot.lot_id:
			//                lot = Lot.create(
			//                    {'name': pack_op_lot.lot_name, 'product_id': pack_op_lot.operation_id.product_id.id})
			//                pack_op_lot.write({'lot_id': lot.id})
			//        self.mapped('pack_operation_ids').mapped('pack_lot_ids').filtered(
			//            lambda op_lot: op_lot.qty == 0.0).unlink()
		})
	//    create_lots_for_picking = _create_lots_for_picking
	h.StockPicking().Methods().CreateExtraMoves().DeclareMethod(
		`This function creates move lines on a picking, at the time
of do_transfer, based on
        unexpected product transfers (or exceeding quantities)
found in the pack operations.
        `,
		func(rs m.StockPickingSet) {
			//        self.ensure_one()
			//        moves = self.env['stock.move']
			//        for pack_operation in self.pack_operation_ids:
			//            for product, remaining_qty in pack_operation._get_remaining_prod_quantities().items():
			//                if float_compare(remaining_qty, 0, precision_rounding=product.uom_id.rounding) > 0:
			//                    vals = self._prepare_values_extra_move(
			//                        pack_operation, product, remaining_qty)
			//                    moves |= moves.create(vals)
			//        if moves:
			//            moves.with_context(skip_check=True).action_confirm()
			//        return moves
		})
	h.StockPicking().Methods().PrepareValuesExtraMove().DeclareMethod(
		`
        Creates an extra move when there is no corresponding
original move to be copied
        `,
		func(rs m.StockPickingSet, op interface{}, product interface{}, remaining_qty interface{}) {
			//        Uom = self.env["product.uom"]
			//        uom_id = product.uom_id.id
			//        qty = remaining_qty
			//        if op.product_id and op.product_uom_id and op.product_uom_id.id != product.uom_id.id:
			//            if op.product_uom_id.factor > product.uom_id.factor:  # If the pack operation's is a smaller unit
			//                uom_id = op.product_uom_id.id
			//                # HALF-UP rounding as only rounding errors will be because of propagation of error from default UoM
			//                qty = product.uom_id._compute_quantity(
			//                    remaining_qty, op.product_uom_id, rounding_method='HALF-UP')
			//        picking = op.picking_id
			//        ref = product.default_code
			//        name = '[' + ref + ']' + ' ' + product.name if ref else product.name
			//        proc_id = False
			//        for m in op.linked_move_operation_ids:
			//            if m.move_id.procurement_id:
			//                proc_id = m.move_id.procurement_id.id
			//                break
			//        return {
			//            'picking_id': picking.id,
			//            'location_id': picking.location_id.id,
			//            'location_dest_id': picking.location_dest_id.id,
			//            'product_id': product.id,
			//            'procurement_id': proc_id,
			//            'product_uom': uom_id,
			//            'product_uom_qty': qty,
			//            'name': _('Extra Move: ') + name,
			//            'state': 'draft',
			//            'restrict_partner_id': op.owner_id.id,
			//            'group_id': picking.group_id.id,
			//        }
		})
	h.StockPicking().Methods().CreateBackorder().DeclareMethod(
		` Move all non-done lines into a new backorder picking.
If the key 'do_only_split' is given in the context, then
move all lines not in context.get('split', []) instead
of all non-done lines.
        `,
		func(rs m.StockPickingSet, backorder_moves interface{}) {
			//        backorders = self.env['stock.picking']
			//        for picking in self:
			//            backorder_moves = backorder_moves or picking.move_lines
			//            if self._context.get('do_only_split'):
			//                not_done_bo_moves = backorder_moves.filtered(
			//                    lambda move: move.id not in self._context.get('split', []))
			//            else:
			//                not_done_bo_moves = backorder_moves.filtered(
			//                    lambda move: move.state not in ('done', 'cancel'))
			//            if not not_done_bo_moves:
			//                continue
			//            backorder_picking = picking.copy({
			//                'name': '/',
			//                'move_lines': [],
			//                'pack_operation_ids': [],
			//                'backorder_id': picking.id
			//            })
			//            picking.message_post(
			//                body=_("Back order <em>%s</em> <b>created</b>.") % (backorder_picking.name))
			//            not_done_bo_moves.write({'picking_id': backorder_picking.id})
			//            if not picking.date_done:
			//                picking.write({'date_done': time.strftime(
			//                    DEFAULT_SERVER_DATETIME_FORMAT)})
			//            backorder_picking.action_confirm()
			//            backorder_picking.action_assign()
			//            backorders |= backorder_picking
			//        return backorders
		})
	h.StockPicking().Methods().PutInPack().DeclareMethod(
		`PutInPack`,
		func(rs m.StockPickingSet) {
			//        QuantPackage = self.env["stock.quant.package"]
			//        package = False
			//        for pick in self:
			//            operations = [x for x in pick.pack_operation_ids if x.qty_done > 0 and (
			//                not x.result_package_id)]
			//            pack_operation_ids = self.env['stock.pack.operation']
			//            for operation in operations:
			//                # If we haven't done all qty in operation, we have to split into 2 operation
			//                op = operation
			//                if operation.qty_done < operation.product_qty:
			//                    new_operation = operation.copy(
			//                        {'product_qty': operation.qty_done, 'qty_done': operation.qty_done})
			//
			//                    operation.write(
			//                        {'product_qty': operation.product_qty - operation.qty_done, 'qty_done': 0})
			//                    if operation.pack_lot_ids:
			//                        packlots_transfer = [(4, x.id)
			//                                             for x in operation.pack_lot_ids]
			//                        new_operation.write(
			//                            {'pack_lot_ids': packlots_transfer})
			//
			//                        # the stock.pack.operation.lot records now belong to the new, packaged stock.pack.operation
			//                        # we have to create new ones with new quantities for our original, unfinished stock.pack.operation
			//                        new_operation._copy_remaining_pack_lot_ids(operation)
			//
			//                    op = new_operation
			//                pack_operation_ids |= op
			//            if operations:
			//                pack_operation_ids.check_tracking()
			//                package = QuantPackage.create({})
			//                pack_operation_ids.write({'result_package_id': package.id})
			//            else:
			//                raise UserError(
			//                    _('Please process some quantities to put in the pack first!'))
			//        return package
		})
	h.StockPicking().Methods().ButtonScrap().DeclareMethod(
		`ButtonScrap`,
		func(rs m.StockPickingSet) {
			//        self.ensure_one()
			//        return {
			//            'name': _('Scrap'),
			//            'view_type': 'form',
			//            'view_mode': 'form',
			//            'res_model': 'stock.scrap',
			//            'view_id': self.env.ref('stock.stock_scrap_form_view2').id,
			//            'type': 'ir.actions.act_window',
			//            'context': {'default_picking_id': self.id, 'product_ids': self.pack_operation_product_ids.mapped('product_id').ids},
			//            'target': 'new',
			//        }
		})
	h.StockPicking().Methods().ActionSeeMoveScrap().DeclareMethod(
		`ActionSeeMoveScrap`,
		func(rs m.StockPickingSet) {
			//        self.ensure_one()
			//        action = self.env.ref('stock.action_stock_scrap').read()[0]
			//        scraps = self.env['stock.scrap'].search([('picking_id', '=', self.id)])
			//        action['domain'] = [('id', 'in', scraps.ids)]
			//        return action
		})
}
