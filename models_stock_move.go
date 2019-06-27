package stock

	import (
		"net/http"

		"github.com/hexya-erp/hexya/src/controllers"
		"github.com/hexya-erp/hexya/src/models"
		"github.com/hexya-erp/hexya/src/models/types"
		"github.com/hexya-erp/hexya/src/models/types/dates"
		"github.com/hexya-erp/pool/h"
		"github.com/hexya-erp/pool/q"
	)
	
//import time
func init() {
h.StockMove().DeclareModel()



h.StockMove().Methods().DefaultGroupId().DeclareMethod(
`DefaultGroupId`,
func(rs m.StockMoveSet)  {
//        if self.env.context.get('default_picking_id'):
//            return self.env['stock.picking'].browse(self.env.context['default_picking_id']).group_id.id
//        return False
})
h.StockMove().AddFields(map[string]models.FieldDefinition{
"Name": models.CharField{
String: "Description",
Index: true,
Required: true,
},
"Sequence": models.IntegerField{
String: "Sequence",
Default: models.DefaultValue(10),
},
"Priority": models.SelectionField{
Selection: procurement.PROCUREMENT_PRIORITIES,
String: "Priority",
Default: models.DefaultValue("1"),
},
"Date": models.DateTimeField{
String: "Date",
Default: func (env models.Environment) interface{} { return dates.Now() },
Index: true,
Required: true,
//states={'done': [('readonly', True)]}
Help: "Move date: scheduled date until move is done, then date" + 
"of actual move processing",
},
"CompanyId": models.Many2OneField{
RelationModel: h.Company(),
String: "Company",
Default: func (env models.Environment) interface{} { return env["res.company"]._company_default_get() },
Index: true,
Required: true,
},
"DateExpected": models.DateTimeField{
String: "Expected Date",
Default: func (env models.Environment) interface{} { return dates.Now() },
Index: true,
Required: true,
//states={'done': [('readonly', True)]}
Help: "Scheduled date for the processing of this move",
},
"ProductId": models.Many2OneField{
RelationModel: h.ProductProduct(),
String: "Product",
Filter: q.Type().In(%!s(<nil>)),
Index: true,
Required: true,
//states={'done': [('readonly', True)]}
},
"OrderedQty": models.FloatField{
String: "Ordered Quantity",
//digits=dp.get_precision('Product Unit of Measure')
},
"ProductQty": models.FloatField{
String: "Real Quantity",
Compute: h.StockMove().Methods().ComputeProductQty(),
//inverse='_set_product_qty'
//digits=0
Stored: true,
Help: "Quantity in the default UoM of the product",
},
"ProductUomQty": models.FloatField{
String: "Quantity",
//digits=dp.get_precision('Product Unit of Measure')
Default: models.DefaultValue(1),
Required: true,
//states={'done': [('readonly', True)]}
Help: "This is the quantity of products from an inventory point" + 
"of view. For moves in the state 'done', this is the quantity" + 
"of products that were actually moved. For other moves," + 
"this is the quantity of product that is planned to be moved." + 
"Lowering this quantity does not generate a backorder. Changing" + 
"this quantity on assigned moves affects the product reservation," + 
"and should be done with care.",
},
"ProductUom": models.Many2OneField{
RelationModel: h.ProductUom(),
String: "Unit of Measure",
Required: true,
//states={'done': [('readonly', True)]}
},
"ProductTmplId": models.Many2OneField{
RelationModel: h.ProductTemplate(),
String: "Product Template",
Related: `ProductId.ProductTmplId`,
Help: "Technical: used in views",
},
"ProductPackaging": models.Many2OneField{
RelationModel: h.ProductPackaging(),
String: "Preferred Packaging",
Help: "It specifies attributes of packaging like type, quantity" + 
"of packaging,etc.",
},
"LocationId": models.Many2OneField{
RelationModel: h.StockLocation(),
String: "Source Location",

Index: true,
Required: true,
//states={'done': [('readonly', True)]}
Help: "Sets a location if you produce at a fixed location. This" + 
"can be a partner location if you subcontract the manufacturing" + 
"operations.",
},
"LocationDestId": models.Many2OneField{
RelationModel: h.StockLocation(),
String: "Destination Location",

Index: true,
Required: true,
//states={'done': [('readonly', True)]}
Help: "Location where the system will stock the finished products.",
},
"PartnerId": models.Many2OneField{
RelationModel: h.Partner(),
String: "Destination Address ",
//states={'done': [('readonly', True)]}
Help: "Optional address where goods are to be delivered, specifically" + 
"used for allotment",
},
"MoveDestId": models.Many2OneField{
RelationModel: h.StockMove(),
String: "Destination Move",
NoCopy: true,
Index: true,
Help: "Optional: next stock move when chaining them",
},
"MoveOrigIds": models.One2ManyField{
RelationModel: h.StockMove(),
ReverseFK: "",
String: "Original Move",
Help: "Optional: previous stock move when chaining them",
},
"PickingId": models.Many2OneField{
RelationModel: h.StockPicking(),
String: "Transfer Reference",
Index: true,
//states={'done': [('readonly', True)]}
},
"PickingPartnerId": models.Many2OneField{
RelationModel: h.Partner(),
String: "Transfer Destination Address",
Related: `PickingId.PartnerId`,
},
"Note": models.TextField{
String: "Notes",
},
"State": models.SelectionField{
Selection: types.Selection{
"draft": "New",
"cancel": "Cancelled",
"waiting": "Waiting Another Move",
"confirmed": "Waiting Availability",
"assigned": "Available",
"done": "Done",
},
String: "Status",
NoCopy: true,
Default: models.DefaultValue("draft"),
Index: true,
ReadOnly: true,
Help: "* New: When the stock move is created and not yet confirmed." + 
"* Waiting Another Move: This state can be seen when a move" + 
"is waiting for another one, for example in a chained flow." + 
"* Waiting Availability: This state is reached when the" + 
"procurement resolution is not straight forward. It may" + 
"need the scheduler to run, a component to be manufactured..." + 
"* Available: When products are reserved, it is set to 'Available'." + 
"* Done: When the shipment is processed, the state is 'Done'.",
},
"PartiallyAvailable": models.BooleanField{
String: "Partially Available",
NoCopy: true,
ReadOnly: true,
Help: "Checks if the move has some stock reserved",
},
"PriceUnit": models.FloatField{
String: "Unit Price",
Help: "Technical field used to record the product cost set by" + 
"the user during a picking confirmation (when costing method" + 
"used is 'average price' or 'real'). Value given in company" + 
"currency and in product uom.",
},
"SplitFrom": models.Many2OneField{
RelationModel: h.StockMove(),
String: "Move Split From",
NoCopy: true,
Help: "Technical field used to track the origin of a split move," + 
"which can be useful in case of debug",
},
"BackorderId": models.Many2OneField{
RelationModel: h.StockPicking(),
String: "Back Order of",
Related: `PickingId.BackorderId`,
Index: true,
},
"Origin": models.CharField{
String: "Source Document",
},
"ProcureMethod": models.SelectionField{
Selection: types.Selection{
"make_to_stock": "Default: Take From Stock",
"make_to_order": "Advanced: Apply Procurement Rules",
},
String: "Supply Method",
Default: models.DefaultValue("make_to_stock"),
Required: true,
Help: "By default, the system will take from the stock in the" + 
"source location and passively wait for availability.The" + 
"other possibility allows you to directly create a procurement" + 
"on the source location (and thus ignore its current stock)" + 
"to gather products. If we want to chain moves and have" + 
"this one to wait for the previous,this second option should be chosen.",
},
"Scrapped": models.BooleanField{
String: "Scrapped",
Related: `LocationDestId.ScrapLocation`,
ReadOnly: true,
Stored: true,
},
"QuantIds": models.Many2ManyField{
RelationModel: h.StockQuant(),
M2MLinkModelName: "",
M2MOurField: "",
M2MTheirField: "",
String: "Moved Quants",
NoCopy: true,
},
"ReservedQuantIds": models.One2ManyField{
RelationModel: h.StockQuant(),
ReverseFK: "",
String: "Reserved quants",
},
"LinkedMoveOperationIds": models.One2ManyField{
RelationModel: h.StockMoveOperationLink(),
ReverseFK: "",
String: "Linked Operations",
ReadOnly: true,
Help: "Operations that impact this move for the computation of" + 
"the remaining quantities",
},
"RemainingQty": models.FloatField{
String: "Remaining Quantity",
Compute: h.StockMove().Methods().GetRemainingQty(),
//digits=0
//states={'done': [('readonly', True)]}
Help: "Remaining Quantity in default UoM according to operations" + 
"matched with this move",
},
"ProcurementId": models.Many2OneField{
RelationModel: h.ProcurementOrder(),
String: "Procurement",
},
"GroupId": models.Many2OneField{
RelationModel: h.ProcurementGroup(),
String: "Procurement Group",
Default: models.DefaultValue(_default_group_id),
},
"RuleId": models.Many2OneField{
RelationModel: h.ProcurementRule(),
String: "Procurement Rule",
Help: "The procurement rule that created this stock move",
},
"PushRuleId": models.Many2OneField{
RelationModel: h.StockLocationPath(),
String: "Push Rule",
Help: "The push rule that created this stock move",
},
"Propagate": models.BooleanField{
String: "Propagate cancel and split",
Default: models.DefaultValue(true),
Help: "If checked, when this move is cancelled, cancel the linked move too",
},
"PickingTypeId": models.Many2OneField{
RelationModel: h.StockPickingType(),
String: "Picking Type",
},
"InventoryId": models.Many2OneField{
RelationModel: h.StockInventory(),
String: "Inventory",
},
"LotIds": models.Many2ManyField{
RelationModel: h.StockProductionLot(),
String: "Lots/Serial Numbers",
Compute: h.StockMove().Methods().ComputeLotIds(),
},
"OriginReturnedMoveId": models.Many2OneField{
RelationModel: h.StockMove(),
String: "Origin return move",
NoCopy: true,
Help: "Move that created the return move",
},
"ReturnedMoveIds": models.One2ManyField{
RelationModel: h.StockMove(),
ReverseFK: "",
String: "All returned moves",
Help: "Optional: all returned moves created from this move",
},
"ReservedAvailability": models.FloatField{
String: "Quantity Reserved",
Compute: h.StockMove().Methods().ComputeReservedAvailability(),
ReadOnly: true,
Help: "Quantity that has already been reserved for this move",
},
"Availability": models.FloatField{
String: "Forecasted Quantity",
Compute: h.StockMove().Methods().ComputeProductAvailability(),
ReadOnly: true,
Help: "Quantity in stock that can still be reserved for this move",
},
"StringAvailabilityInfo": models.TextField{
String: "Availability",
Compute: h.StockMove().Methods().ComputeStringQtyInformation(),
ReadOnly: true,
Help: "Show various information on stock availability for this move",
},
"RestrictLotId": models.Many2OneField{
RelationModel: h.StockProductionLot(),
String: "Lot/Serial Number",
Help: "Technical field used to depict a restriction on the lot/serial" + 
"number of quants to consider when marking this move as 'done'",
},
"RestrictPartnerId": models.Many2OneField{
RelationModel: h.Partner(),
String: "Owner ",
Help: "Technical field used to depict a restriction on the ownership" + 
"of quants to consider when marking this move as 'done'",
},
"RouteIds": models.Many2ManyField{
RelationModel: h.StockLocationRoute(),
M2MLinkModelName: "",
M2MOurField: "",
M2MTheirField: "",
String: "Destination route",
Help: "Preferred route to be followed by the procurement order",
},
"WarehouseId": models.Many2OneField{
RelationModel: h.StockWarehouse(),
String: "Warehouse",
Help: "Technical field depicting the warehouse to consider for" + 
"the route selection on the next procurement (if any).",
},
})
h.StockMove().Fields().CreateDate().setString( "Creation Date",)
h.StockMove().Fields().CreateDate().setIndex( true,)
h.StockMove().Fields().CreateDate().setReadOnly( true,)
h.StockMove().Methods().ComputeProductQty().DeclareMethod(
`ComputeProductQty`,
func(rs h.StockMoveSet) h.StockMoveData {
//        if self.product_uom:
//            rounding_method = self._context.get('rounding_method', 'UP')
//            self.product_qty = self.product_uom._compute_quantity(
//                self.product_uom_qty, self.product_id.uom_id, rounding_method=rounding_method)
})
h.StockMove().Methods().SetProductQty().DeclareMethod(
` The meaning of product_qty field changed lately and is
now a functional field computing the quantity
        in the default product UoM. This code has been
added to raise an error if a write is made given a value
        for `product_qty`, where the same write should
set the `product_uom_qty` field instead, in order to
        detect errors. `,
func(rs m.StockMoveSet)  {
//        raise UserError(
//            _('The requested operation cannot be processed because of a programming error setting the `product_qty` field instead of the `product_uom_qty`.'))
})
h.StockMove().Methods().GetRemainingQty().DeclareMethod(
`GetRemainingQty`,
func(rs h.StockMoveSet) h.StockMoveData {
//        self.remaining_qty = float_round(self.product_qty - sum(self.mapped(
//            'linked_move_operation_ids').mapped('qty')), precision_rounding=self.product_id.uom_id.rounding)
})
h.StockMove().Methods().ComputeLotIds().DeclareMethod(
`ComputeLotIds`,
func(rs h.StockMoveSet) h.StockMoveData {
//        if self.state == 'done':
//            self.lot_ids = self.mapped('quant_ids').mapped('lot_id').ids
//        else:
//            self.lot_ids = self.mapped(
//                'reserved_quant_ids').mapped('lot_id').ids
})
h.StockMove().Methods().ComputeReservedAvailability().DeclareMethod(
`ComputeReservedAvailability`,
func(rs h.StockMoveSet) h.StockMoveData {
//        result = {data['reservation_id'][0]: data['qty'] for data in
//                  self.env['stock.quant'].read_group([('reservation_id', 'in', self.ids)], ['reservation_id', 'qty'], ['reservation_id'])}
//        for rec in self:
//            rec.reserved_availability = result.get(rec.id, 0.0)
})
h.StockMove().Methods().ComputeProductAvailability().DeclareMethod(
`ComputeProductAvailability`,
func(rs h.StockMoveSet) h.StockMoveData {
//        if self.state == 'done':
//            self.availability = self.product_qty
//        else:
//            qty_tot = self.env['stock.quant'].read_group([('location_id', 'child_of', self.location_id.id), (
//                'product_id', '=', self.product_id.id), ('reservation_id', '=', False)], ['qty'], [])[0]
//            self.availability = min(self.product_qty, qty_tot['qty'] or 0)
})
h.StockMove().Methods().ComputeStringQtyInformation().DeclareMethod(
`ComputeStringQtyInformation`,
func(rs h.StockMoveSet) h.StockMoveData {
//        precision = self.env['decimal.precision'].precision_get(
//            'Product Unit of Measure')
//        void_moves = self.filtered(lambda move: move.state in (
//            'draft', 'done', 'cancel') or move.location_id.usage != 'internal')
//        other_moves = self - void_moves
//        for move in void_moves:
//            move.string_availability_info = ''  # 'not applicable' or 'n/a' could work too
//        for move in other_moves:
//            total_available = min(
//                move.product_qty, move.reserved_availability + move.availability)
//            total_available = move.product_id.uom_id._compute_quantity(
//                total_available, move.product_uom, round=False)
//            total_available = float_round(
//                total_available, precision_digits=precision)
//            info = str(total_available)
//            if self.user_has_groups('product.group_uom'):
//                info += ' ' + move.product_uom.name
//            if move.reserved_availability:
//                if move.reserved_availability != total_available:
//                    # some of the available quantity is assigned and some are available but not reserved
//                    reserved_available = move.product_id.uom_id._compute_quantity(
//                        move.reserved_availability, move.product_uom, round=False)
//                    reserved_available = float_round(
//                        reserved_available, precision_digits=precision)
//                    info += _(' (%s reserved)') % str(reserved_available)
//                else:
//                    # all available quantity is assigned
//                    info += _(' (reserved)')
//            move.string_availability_info = info
})
h.StockMove().Methods().CheckUom().DeclareMethod(
`CheckUom`,
func(rs m.StockMoveSet)  {
//        moves_error = self.filtered(
//            lambda move: move.product_id.uom_id.category_id != move.product_uom.category_id)
//        if moves_error:
//            user_warning = _(
//                'You try to move a product using a UoM that is not compatible with the UoM of the product moved. Please use an UoM in the same UoM category.')
//            for move in moves_error:
//                user_warning += _('\n\n%s --> Product UoM is %s (%s) - Move UoM is %s (%s)') % (move.product_id.display_name,
//                                                                                                move.product_id.uom_id.name, move.product_id.uom_id.category_id.name, move.product_uom.name, move.product_uom.category_id.name)
//            user_warning += _('\n\nBlocking: %s') % ' ,'.join(moves_error.mapped('name'))
//            raise UserError(user_warning)
})
h.StockMove().Methods().Init().DeclareMethod(
`Init`,
func(rs m.StockMoveSet)  {
//        self._cr.execute('SELECT indexname FROM pg_indexes WHERE indexname = %s',
//                         ('stock_move_product_location_index'))
//        if not self._cr.fetchone():
//            self._cr.execute(
//                'CREATE INDEX stock_move_product_location_index ON stock_move (product_id, location_id, location_dest_id, company_id, state)')
})
h.StockMove().Methods().NameGet().Extend(
`NameGet`,
func(rs m.StockMoveSet)  {
//        res = []
//        for move in self:
//            res.append((move.id, '%s%s%s>%s' % (
//                move.picking_id.origin and '%s/' % move.picking_id.origin or '',
//                move.product_id.code and '%s: ' % move.product_id.code or '',
//                move.location_id.name, move.location_dest_id.name)))
//        return res
})
h.StockMove().Methods().Create().Extend(
`Create`,
func(rs m.StockMoveSet, vals models.RecordData)  {
//        perform_tracking = not self.env.context.get(
//            'mail_notrack') and vals.get('picking_id')
//        if perform_tracking:
//            picking = self.env['stock.picking'].browse(vals['picking_id'])
//            initial_values = {picking.id: {'state': picking.state}}
//        vals['ordered_qty'] = vals.get('product_uom_qty')
//        res = super(StockMove, self).create(vals)
//        if perform_tracking:
//            picking.message_track(
//                picking.fields_get(['state']), initial_values)
//        return res
})
h.StockMove().Methods().Write().Extend(
`Write`,
func(rs m.StockMoveSet, vals models.RecordData)  {
//        Picking = self.env['stock.picking']
//        frozen_fields = ['product_qty', 'product_uom',
//                         'location_id', 'location_dest_id', 'product_id']
//        if any(fname in frozen_fields for fname in vals.keys()) and any(move.state == 'done' for move in self):
//            raise UserError(
//                _('Quantities, Units of Measure, Products and Locations cannot be modified on stock moves that have already been processed (except by the Administrator).'))
//        propagated_changes_dict = {}
//        propagated_date_field = False
//        if vals.get('date_expected'):
//            # propagate any manual change of the expected date
//            propagated_date_field = 'date_expected'
//        elif (vals.get('state', '') == 'done' and vals.get('date')):
//            # propagate also any delta observed when setting the move as done
//            propagated_date_field = 'date'
//        if not self._context.get('do_not_propagate', False) and (propagated_date_field or propagated_changes_dict):
//            # any propagation is (maybe) needed
//            for move in self:
//                if move.move_dest_id and move.propagate:
//                    if 'date_expected' in propagated_changes_dict:
//                        propagated_changes_dict.pop('date_expected')
//                    if propagated_date_field:
//                        current_date = datetime.strptime(
//                            move.date_expected, DEFAULT_SERVER_DATETIME_FORMAT)
//                        new_date = datetime.strptime(
//                            vals.get(propagated_date_field), DEFAULT_SERVER_DATETIME_FORMAT)
//                        delta_days = (
//                            new_date - current_date).total_seconds() / 86400
//                        if abs(delta_days) >= move.company_id.propagation_minimum_delta:
//                            old_move_date = datetime.strptime(
//                                move.move_dest_id.date_expected, DEFAULT_SERVER_DATETIME_FORMAT)
//                            new_move_date = (old_move_date + relativedelta.relativedelta(
//                                days=delta_days or 0)).strftime(DEFAULT_SERVER_DATETIME_FORMAT)
//                            propagated_changes_dict['date_expected'] = new_move_date
//                    # For pushed moves as well as for pulled moves, propagate by recursive call of write().
//                    # Note that, for pulled moves we intentionally don't propagate on the procurement.
//                    if propagated_changes_dict:
//                        move.move_dest_id.write(propagated_changes_dict)
//        track_pickings = not self._context.get('mail_notrack') and any(
//            field in vals for field in ['state', 'picking_id', 'partially_available'])
//        if track_pickings:
//            to_track_picking_ids = set(
//                [move.picking_id.id for move in self if move.picking_id])
//            if vals.get('picking_id'):
//                to_track_picking_ids.add(vals['picking_id'])
//            to_track_picking_ids = list(to_track_picking_ids)
//            pickings = Picking.browse(to_track_picking_ids)
//            initial_values = dict(
//                (picking.id, {'state': picking.state}) for picking in pickings)
//        res = super(StockMove, self).write(vals)
//        if track_pickings:
//            pickings.message_track(
//                pickings.fields_get(['state']), initial_values)
//        return res
})
h.StockMove().Methods().GetPriceUnit().DeclareMethod(
` Returns the unit price to store on the quant `,
func(rs m.StockMoveSet)  {
//        return self.price_unit or self.product_id.standard_price
})
h.StockMove().Methods().GetRemovalStrategy().DeclareMethod(
` Returns the removal strategy to consider for the given move/ops `,
func(rs m.StockMoveSet)  {
//        if self.product_id.categ_id.removal_strategy_id:
//            return self.product_id.categ_id.removal_strategy_id.method
//        loc = self.location_id
//        while loc:
//            if loc.removal_strategy_id:
//                return loc.removal_strategy_id.method
//            loc = loc.location_id
//        return 'fifo'
})
h.StockMove().Methods().GetAncestors().DeclareMethod(
`Find the first level ancestors of given move `,
func(rs m.StockMoveSet)  {
//        ancestors = self.env['stock.move']
//        move = self
//        while move:
//            ancestors |= move.move_orig_ids
//            move = not move.move_orig_ids and move.split_from or False
//        return ancestors
})
//    find_move_ancestors = get_ancestors
h.StockMove().Methods().FilterClosedMoves().DeclareMethod(
` Helper methods when having to avoid working on moves that are
        already done or canceled. In a lot of cases you
may handle a batch
        of stock moves, some being already done / canceled,
other being still
        under computation. Instead of having to use filtered
everywhere and
        forgot some of them, use this tool instead. `,
func(rs m.StockMoveSet)  {
//        return self.filtered(lambda move: move.state not in ('done', 'cancel'))
})
h.StockMove().Methods().DoUnreserve().DeclareMethod(
`DoUnreserve`,
func(rs m.StockMoveSet)  {
//        if any(move.state in ('done', 'cancel') for move in self):
//            raise UserError(_('Cannot unreserve a done move'))
//        self.quants_unreserve()
//        if not self.env.context.get('no_state_change'):
//            waiting = self.filtered(
//                lambda move: move.procure_method == 'make_to_order' or move.get_ancestors())
//            waiting.write({'state': 'waiting'})
//            (self - waiting).write({'state': 'confirmed'})
})
h.StockMove().Methods().PushApply().DeclareMethod(
`PushApply`,
func(rs m.StockMoveSet)  {
//        Push = self.env['stock.location.path']
//        for move in self:
//            # if the move is already chained, there is no need to check push rules
//            if move.move_dest_id:
//                continue
//            # if the move is a returned move, we don't want to check push rules, as returning a returned move is the only decent way
//            # to receive goods without triggering the push rules again (which would duplicate chained operations)
//            domain = [('location_from_id', '=', move.location_dest_id.id)]
//            # priority goes to the route defined on the product and product category
//            routes = move.product_id.route_ids | move.product_id.categ_id.total_route_ids
//            rules = Push.search(
//                domain + [('route_id', 'in', routes.ids)], order='route_sequence, sequence', limit=1)
//            if not rules:
//                # TDE FIXME/ should those really be in a if / elif ??
//                # then we search on the warehouse if a rule can apply
//                if move.warehouse_id:
//                    rules = Push.search(
//                        domain + [('route_id', 'in', move.warehouse_id.route_ids.ids)], order='route_sequence, sequence', limit=1)
//                elif move.picking_id.picking_type_id.warehouse_id:
//                    rules = Push.search(
//                        domain + [('route_id', 'in', move.picking_id.picking_type_id.warehouse_id.route_ids.ids)], order='route_sequence, sequence', limit=1)
//            if not rules:
//                # if no specialized push rule has been found yet, we try to find a general one (without route)
//                rules = Push.search(
//                    domain + [('route_id', '=', False)], order='sequence', limit=1)
//            # Make sure it is not returning the return
//            if rules and (not move.origin_returned_move_id or move.origin_returned_move_id.location_dest_id.id != rules.location_dest_id.id):
//                rules._apply(move)
//        return True
})
h.StockMove().Methods().OnchangeQuantity().DeclareMethod(
`OnchangeQuantity`,
func(rs m.StockMoveSet)  {
//        if not self.product_id or self.product_qty < 0.0:
//            self.product_qty = 0.0
//        if self.product_qty < self._origin.product_qty:
//            warning_mess = {
//                'title': _('Quantity decreased!'),
//                'message': _("By changing this quantity here, you accept the "
//                             "new quantity as complete: Odoo will not "
//                             "automatically generate a back order."),
//            }
//            return {'warning': warning_mess}
})
h.StockMove().Methods().OnchangeProductId().DeclareMethod(
`OnchangeProductId`,
func(rs m.StockMoveSet)  {
//        product = self.product_id.with_context(
//            lang=self.partner_id.lang or self.env.user.lang)
//        self.name = product.partner_ref
//        self.product_uom = product.uom_id.id
//        self.product_uom_qty = 1.0
//        return {'domain': {'product_uom': [('category_id', '=', product.uom_id.category_id.id)]}}
})
h.StockMove().Methods().OnchangeDate().DeclareMethod(
`OnchangeDate`,
func(rs m.StockMoveSet)  {
//        if self.date_expected:
//            self.date = self.date_expected
})
h.StockMove().Methods().AssignPicking().DeclareMethod(
` Try to assign the moves to an existing picking that has not been
        reserved yet and has the same procurement group,
locations and picking
        type (moves should already have them identical).
Otherwise, create a new
        picking to assign them to. `,
func(rs m.StockMoveSet)  {
//        Picking = self.env['stock.picking']
//        for move in self:
//            recompute = False
//            picking = Picking.search([
//                ('group_id', '=', move.group_id.id),
//                ('location_id', '=', move.location_id.id),
//                ('location_dest_id', '=', move.location_dest_id.id),
//                ('picking_type_id', '=', move.picking_type_id.id),
//                ('printed', '=', False),
//                ('state', 'in', ['draft', 'confirmed', 'waiting', 'partially_available', 'assigned'])], limit=1)
//            if not picking:
//                recompute = True
//                picking = Picking.create(move._get_new_picking_values())
//            move.write({'picking_id': picking.id})
//
//            # If this method is called in batch by a write on a one2many and
//            # at some point had to create a picking, some next iterations could
//            # try to find back the created picking. As we look for it by searching
//            # on some computed fields, we have to force a recompute, else the
//            # record won't be found.
//            if recompute:
//                move.recompute()
//        return True
})
//    _picking_assign = assign_picking
h.StockMove().Methods().GetNewPickingValues().DeclareMethod(
` Prepares a new picking for this move as it could not be assigned to
        another picking. This method is designed to be inherited. `,
func(rs m.StockMoveSet)  {
//        return {
//            'origin': self.origin,
//            'company_id': self.company_id.id,
//            'move_type': self.group_id and self.group_id.move_type or 'direct',
//            'partner_id': self.partner_id.id,
//            'picking_type_id': self.picking_type_id.id,
//            'location_id': self.location_id.id,
//            'location_dest_id': self.location_dest_id.id,
//        }
})
//    _prepare_picking_assign = _get_new_picking_values
h.StockMove().Methods().ActionConfirm().DeclareMethod(
` Confirms stock move or put it in waiting if it's linked
to another move. `,
func(rs m.StockMoveSet)  {
//        move_create_proc = self.env['stock.move']
//        move_to_confirm = self.env['stock.move']
//        move_waiting = self.env['stock.move']
//        to_assign = {}
//        self.set_default_price_unit_from_product()
//        for move in self:
//            # if the move is preceeded, then it's waiting (if preceeding move is done, then action_assign has been called already and its state is already available)
//            if move.move_orig_ids:
//                move_waiting |= move
//            # if the move is split and some of the ancestor was preceeded, then it's waiting as well
//            else:
//                inner_move = move.split_from
//                while inner_move:
//                    if inner_move.move_orig_ids:
//                        move_waiting |= move
//                        break
//                    inner_move = inner_move.split_from
//                else:
//                    if move.procure_method == 'make_to_order':
//                        move_create_proc |= move
//                    else:
//                        move_to_confirm |= move
//
//            if not move.picking_id and move.picking_type_id:
//                key = (move.group_id.id, move.location_id.id,
//                       move.location_dest_id.id)
//                if key not in to_assign:
//                    to_assign[key] = self.env['stock.move']
//                to_assign[key] |= move
//        procurements = self.env['procurement.order']
//        for move in move_create_proc:
//            procurements |= procurements.create(
//                move._prepare_procurement_from_move())
//        if procurements:
//            procurements.run()
//        move_to_confirm.write({'state': 'confirmed'})
//        (move_waiting | move_create_proc).write({'state': 'waiting'})
//        for key, moves in to_assign.items():
//            moves.assign_picking()
//        self._push_apply()
//        return self
})
h.StockMove().Methods().SetDefaultPriceMoves().DeclareMethod(
`SetDefaultPriceMoves`,
func(rs m.StockMoveSet)  {
//        return self.filtered(lambda move: not move.price_unit)
})
h.StockMove().Methods().SetDefaultPriceUnitFromProduct().DeclareMethod(
` Set price to move, important in inter-company moves or
receipts with only one partner `,
func(rs m.StockMoveSet)  {
//        for move in self._set_default_price_moves():
//            move.write({'price_unit': move.product_id.standard_price})
})
//    attribute_price = set_default_price_unit_from_product
h.StockMove().Methods().PrepareProcurementFromMove().DeclareMethod(
`PrepareProcurementFromMove`,
func(rs m.StockMoveSet)  {
//        origin = (self.group_id and (self.group_id.name + ":") or "") + \
//            (self.rule_id and self.rule_id.name or self.origin or self.picking_id.name or "/")
//        group_id = self.group_id and self.group_id.id or False
//        if self.rule_id:
//            if self.rule_id.group_propagation_option == 'fixed' and self.rule_id.group_id:
//                group_id = self.rule_id.group_id.id
//            elif self.rule_id.group_propagation_option == 'none':
//                group_id = False
//        return {
//            'name': self.rule_id and self.rule_id.name or "/",
//            'origin': origin,
//            'company_id': self.company_id.id,
//            'date_planned': self.date_expected,
//            'product_id': self.product_id.id,
//            'product_qty': self.product_uom_qty,
//            'product_uom': self.product_uom.id,
//            'location_id': self.location_id.id,
//            'move_dest_id': self.id,
//            'group_id': group_id,
//            'route_ids': [(4, x.id) for x in self.route_ids],
//            'warehouse_id': self.warehouse_id.id or (self.picking_type_id and self.picking_type_id.warehouse_id.id or False),
//            'priority': self.priority,
//        }
})
h.StockMove().Methods().ForceAssign().DeclareMethod(
`ForceAssign`,
func(rs m.StockMoveSet)  {
//        self.write({'state': 'assigned'})
//        self.check_recompute_pack_op()
})
h.StockMove().Methods().CheckRecomputePackOp().DeclareMethod(
`CheckRecomputePackOp`,
func(rs m.StockMoveSet)  {
//        pickings = self.mapped('picking_id').filtered(
//            lambda picking: picking.state not in ('waiting', 'confirmed'))
//        pickings_partial = pickings.filtered(lambda picking: not any(
//            operation.qty_done for operation in picking.pack_operation_ids))
//        pickings_partial.do_prepare_partial()
//        (pickings - pickings_partial).write({'recompute_pack_op': True})
})
h.StockMove().Methods().CheckTracking().DeclareMethod(
` Checks if serial number is assigned to stock move or not
and raise an error if it had to. `,
func(rs m.StockMoveSet, pack_operation interface{})  {
//        for move in self:
//            if move.picking_id and \
//                    (move.picking_id.picking_type_id.use_existing_lots or move.picking_id.picking_type_id.use_create_lots) and \
//                    move.product_id.tracking != 'none' and \
//                    not (move.restrict_lot_id or (pack_operation and (pack_operation.product_id and pack_operation.pack_lot_ids)) or (pack_operation and not pack_operation.product_id)):
//                raise UserError(_('You need to provide a Lot/Serial Number for product %s') %
//                                ("%s (%s)" % (move.product_id.name, move.picking_id.name)))
})
h.StockMove().Methods().ActionAssign().DeclareMethod(
` Checks the product type and accordingly writes the state. `,
func(rs m.StockMoveSet, no_prepare interface{})  {
//        main_domain = {}
//        Quant = self.env['stock.quant']
//        Uom = self.env['product.uom']
//        moves_to_assign = self.env['stock.move']
//        moves_to_do = self.env['stock.move']
//        operations = self.env['stock.pack.operation']
//        ancestors_list = {}
//        moves = self.filtered(lambda move: move.state in [
//                              'confirmed', 'waiting', 'assigned'])
//        moves.filtered(lambda move: move.reserved_quant_ids).do_unreserve()
//        for move in moves:
//            if move.location_id.usage in ('supplier', 'inventory', 'production'):
//                moves_to_assign |= move
//                # TDE FIXME: what ?
//                # in case the move is returned, we want to try to find quants before forcing the assignment
//                if not move.origin_returned_move_id:
//                    continue
//            # if the move is preceeded, restrict the choice of quants in the ones moved previously in original move
//            ancestors = move.find_move_ancestors()
//            if move.product_id.type == 'consu' and not ancestors:
//                moves_to_assign |= move
//                continue
//            else:
//                moves_to_do |= move
//
//                # we always search for yet unassigned quants
//                main_domain[move.id] = [
//                    ('reservation_id', '=', False), ('qty', '>', 0)]
//
//                ancestors_list[move.id] = True if ancestors else False
//                if move.state == 'waiting' and not ancestors:
//                    # if the waiting move hasn't yet any ancestor (PO/MO not confirmed yet), don't find any quant available in stock
//                    main_domain[move.id] += [('id', '=', False)]
//                elif ancestors:
//                    main_domain[move.id] += [('history_ids',
//                                              'in', ancestors.ids)]
//
//                # if the move is returned from another, restrict the choice of quants to the ones that follow the returned move
//                if move.origin_returned_move_id:
//                    main_domain[move.id] += [('history_ids',
//                                              'in', move.origin_returned_move_id.id)]
//                for link in move.linked_move_operation_ids:
//                    operations |= link.operation_id
//        operations = operations.sorted(key=lambda x: (
//            (x.package_id and not x.product_id) and -4 or 0) + (x.package_id and -2 or 0) + (x.pack_lot_ids and -1 or 0))
//        for ops in operations:
//            # TDE FIXME: this code seems to be in action_done, isn't it ?
//            # first try to find quants based on specific domains given by linked operations for the case where we want to rereserve according to existing pack operations
//            if not (ops.product_id and ops.pack_lot_ids):
//                for record in ops.linked_move_operation_ids:
//                    move = record.move_id
//                    if move.id in main_domain:
//                        qty = record.qty
//                        domain = main_domain[move.id]
//                        if qty:
//                            quants = Quant.quants_get_preferred_domain(
//                                qty, move, ops=ops, domain=domain, preferred_domain_list=[])
//                            Quant.quants_reserve(quants, move, record)
//            else:
//                lot_qty = {}
//                rounding = ops.product_id.uom_id.rounding
//                for pack_lot in ops.pack_lot_ids:
//                    lot_qty[pack_lot.lot_id.id] = ops.product_uom_id._compute_quantity(
//                        pack_lot.qty, ops.product_id.uom_id)
//                for record in ops.linked_move_operation_ids:
//                    move_qty = record.qty
//                    move = record.move_id
//                    domain = main_domain[move.id]
//                    for lot in lot_qty:
//                        if float_compare(lot_qty[lot], 0, precision_rounding=rounding) > 0 and float_compare(move_qty, 0, precision_rounding=rounding) > 0:
//                            qty = min(lot_qty[lot], move_qty)
//                            quants = Quant.quants_get_preferred_domain(
//                                qty, move, ops=ops, lot_id=lot, domain=domain, preferred_domain_list=[])
//                            Quant.quants_reserve(quants, move, record)
//                            lot_qty[lot] -= qty
//                            move_qty -= qty
//        for move in sorted(moves_to_do, key=lambda x: -1 if ancestors_list.get(x.id) else 0):
//            # then if the move isn't totally assigned, try to find quants without any specific domain
//            if move.state != 'assigned' and not self.env.context.get('reserve_only_ops'):
//                qty_already_assigned = move.reserved_availability
//                qty = move.product_qty - qty_already_assigned
//
//                quants = Quant.quants_get_preferred_domain(
//                    qty, move, domain=main_domain[move.id], preferred_domain_list=[])
//                Quant.quants_reserve(quants, move)
//        if moves_to_assign:
//            moves_to_assign.write({'state': 'assigned'})
//        if not no_prepare:
//            self.check_recompute_pack_op()
})
h.StockMove().Methods().PropagateCancel().DeclareMethod(
`PropagateCancel`,
func(rs m.StockMoveSet)  {
//        self.ensure_one()
//        if self.move_dest_id:
//            if self.propagate:
//                if self.move_dest_id.state not in ('done', 'cancel'):
//                    self.move_dest_id.action_cancel()
//            elif self.move_dest_id.state == 'waiting':
//                # If waiting, the chain will be broken and we are not sure if we can still wait for it (=> could take from stock instead)
//                self.move_dest_id.write({'state': 'confirmed'})
})
h.StockMove().Methods().ActionCancel().DeclareMethod(
` Cancels the moves and if all moves are cancelled it cancels
the picking. `,
func(rs m.StockMoveSet)  {
//        if any(move.state == 'done' for move in self):
//            raise UserError(
//                _('You cannot cancel a stock move that has been set to \'Done\'.'))
//        procurements = self.env['procurement.order']
//        for move in self:
//            if move.reserved_quant_ids:
//                move.quants_unreserve()
//            if self.env.context.get('cancel_procurement'):
//                if move.propagate:
//                    pass
//                    # procurements.search([('move_dest_id', '=', move.id)]).cancel()
//            else:
//                move._propagate_cancel()
//                if move.procurement_id:
//                    procurements |= move.procurement_id
//        self.write({'state': 'cancel', 'move_dest_id': False})
//        if procurements:
//            procurements.check()
//        return True
})
h.StockMove().Methods().RecalculateMoveState().DeclareMethod(
`Recompute the state of moves given because their reserved
quants were used to fulfill another operation`,
func(rs m.StockMoveSet)  {
//        for move in self:
//            vals = {}
//            reserved_quant_ids = move.reserved_quant_ids
//            if len(reserved_quant_ids) > 0 and not move.partially_available:
//                vals['partially_available'] = True
//            if len(reserved_quant_ids) == 0 and move.partially_available:
//                vals['partially_available'] = False
//            if move.state == 'assigned':
//                if move.procure_method == 'make_to_order' or move.find_move_ancestors():
//                    vals['state'] = 'waiting'
//                else:
//                    vals['state'] = 'confirmed'
//            if vals:
//                move.write(vals)
})
h.StockMove().Methods().MoveQuantsByLot().DeclareMethod(
`
        This function is used to process all the pack operation
lots of a pack operation
        For every move:
            First, we check the quants with lot already
reserved (and those are already subtracted from the lots to do)
            Then go through all the lots to process:
                Add reserved false lots lot by lot
                Check if there are not reserved quants
or reserved elsewhere with that lot or without lot (with
the traditional method)
        `,
func(rs m.StockMoveSet, ops interface{}, lot_qty interface{}, quants_taken interface{}, false_quants interface{}, lot_move_qty interface{}, quant_dest_package_id interface{})  {
//        return self.browse(lot_move_qty.keys())._move_quants_by_lot_v10(quants_taken, false_quants, ops, lot_qty, lot_move_qty, quant_dest_package_id)
})
h.StockMove().Methods().MoveQuantsByLotV10().DeclareMethod(
`MoveQuantsByLotV10`,
func(rs m.StockMoveSet, quants_taken interface{}, false_quants interface{}, pack_operation interface{}, lot_quantities interface{}, lot_move_quantities interface{}, dest_package_id interface{})  {
//        Quant = self.env['stock.quant']
//        rounding = pack_operation.product_id.uom_id.rounding
//        preferred_domain_list = [[('reservation_id', '=', False)], [
//            '&', ('reservation_id', 'not in', self.ids), ('reservation_id', '!=', False)]]
//        for move_rec_updateme in self:
//            from collections import defaultdict
//            lot_to_quants = defaultdict(list)
//
//            # Assign quants already reserved with lot to the correct
//            for quant in quants_taken:
//                if quant[0] <= move_rec_updateme.reserved_quant_ids:
//                    lot_to_quants[quant[0].lot_id.id].append(quant)
//
//            false_quants_move = [
//                x for x in false_quants if x[0].reservation_id.id == move_rec_updateme.id]
//            for lot_id in lot_quantities.keys():
//                redo_false_quants = False
//
//                # Take remaining reserved quants with  no lot first
//                # (This will be used mainly when incoming had no lot and you do outgoing with)
//                while false_quants_move and float_compare(lot_quantities[lot_id], 0, precision_rounding=rounding) > 0 and float_compare(lot_move_quantities[move_rec_updateme.id], 0, precision_rounding=rounding) > 0:
//                    qty_min = min(
//                        lot_quantities[lot_id], lot_move_quantities[move_rec_updateme.id])
//                    if false_quants_move[0].qty > qty_min:
//                        lot_to_quants[lot_id] += [
//                            (false_quants_move[0], qty_min)]
//                        qty = qty_min
//                        redo_false_quants = True
//                    else:
//                        qty = false_quants_move[0].qty
//                        lot_to_quants[lot_id] += [(false_quants_move[0], qty)]
//                        false_quants_move.pop(0)
//                    lot_quantities[lot_id] -= qty
//                    lot_move_quantities[move_rec_updateme.id] -= qty
//
//                # Search other with first matching lots and then without lots
//                if float_compare(lot_move_quantities[move_rec_updateme.id], 0, precision_rounding=rounding) > 0 and float_compare(lot_quantities[lot_id], 0, precision_rounding=rounding) > 0:
//                    # Search if we can find quants with that lot
//                    qty = min(
//                        lot_quantities[lot_id], lot_move_quantities[move_rec_updateme.id])
//                    quants = Quant.quants_get_preferred_domain(
//                        qty, move_rec_updateme, ops=pack_operation, lot_id=lot_id, domain=[
//                            ('qty', '>', 0)],
//                        preferred_domain_list=preferred_domain_list)
//                    lot_to_quants[lot_id] += quants
//                    lot_quantities[lot_id] -= qty
//                    lot_move_quantities[move_rec_updateme.id] -= qty
//
//                # Move all the quants related to that lot/move
//                if lot_to_quants[lot_id]:
//                    Quant.quants_move(
//                        lot_to_quants[lot_id], move_rec_updateme, pack_operation.location_dest_id,
//                        location_from=pack_operation.location_id, lot_id=lot_id,
//                        owner_id=pack_operation.owner_id.id, src_package_id=pack_operation.package_id.id,
//                        dest_package_id=dest_package_id)
//                    if redo_false_quants:
//                        false_quants_move = [x for x in move_rec_updateme.reserved_quant_ids if (not x.lot_id) and (x.owner_id.id == pack_operation.owner_id.id) and
//                                             (x.location_id.id == pack_operation.location_id.id) and (x.package_id.id == pack_operation.package_id.id)]
//        return True
})
h.StockMove().Methods().ActionDone().DeclareMethod(
` Process completely the moves given and if all moves are
done, it will finish the picking. `,
func(rs m.StockMoveSet)  {
//        self.filtered(lambda move: move.state == 'draft').action_confirm()
//        Uom = self.env['product.uom']
//        Quant = self.env['stock.quant']
//        pickings = self.env['stock.picking']
//        procurements = self.env['procurement.order']
//        operations = self.env['stock.pack.operation']
//        remaining_move_qty = {}
//        for move in self:
//            if move.picking_id:
//                pickings |= move.picking_id
//            remaining_move_qty[move.id] = move.product_qty
//            for link in move.linked_move_operation_ids:
//                operations |= link.operation_id
//                pickings |= link.operation_id.picking_id
//        operations = operations.sorted(key=lambda x: (
//            (x.package_id and not x.product_id) and -4 or 0) + (x.package_id and -2 or 0) + (x.pack_lot_ids and -1 or 0))
//        for operation in operations:
//
//            # product given: result put immediately in the result package (if False: without package)
//            # but if pack moved entirely, quants should not be written anything for the destination package
//            quant_dest_package_id = operation.product_id and operation.result_package_id.id or False
//            entire_pack = not operation.product_id and True or False
//
//            # compute quantities for each lot + check quantities match
//            lot_quantities = dict((pack_lot.lot_id.id, operation.product_uom_id._compute_quantity(pack_lot.qty, operation.product_id.uom_id)
//                                   ) for pack_lot in operation.pack_lot_ids)
//
//            qty = operation.product_qty
//            if operation.product_uom_id and operation.product_uom_id != operation.product_id.uom_id:
//                qty = operation.product_uom_id._compute_quantity(
//                    qty, operation.product_id.uom_id)
//            if operation.pack_lot_ids and float_compare(sum(lot_quantities.values()), qty, precision_rounding=operation.product_id.uom_id.rounding) != 0.0:
//                raise UserError(
//                    _('You have a difference between the quantity on the operation and the quantities specified for the lots. '))
//
//            quants_taken = []
//            false_quants = []
//            lot_move_qty = {}
//
//            prout_move_qty = {}
//            for link in operation.linked_move_operation_ids:
//                prout_move_qty[link.move_id] = prout_move_qty.get(
//                    link.move_id, 0.0) + link.qty
//
//            # Process every move only once for every pack operation
//            for move in prout_move_qty.keys():
//                # TDE FIXME: do in batch ?
//                move.check_tracking(operation)
//
//                # TDE FIXME: I bet the message error is wrong
//                if not remaining_move_qty.get(move.id):
//                    raise UserError(_("The roundings of your unit of measure %s on the move vs. %s on the product don't allow to do these operations or you are not transferring the picking at once. ") % (
//                        move.product_uom.name, move.product_id.uom_id.name))
//
//                if not operation.pack_lot_ids:
//                    preferred_domain_list = [[('reservation_id', '=', move.id)], [('reservation_id', '=', False)], [
//                        '&', ('reservation_id', '!=', move.id), ('reservation_id', '!=', False)]]
//                    quants = Quant.quants_get_preferred_domain(
//                        prout_move_qty[move], move, ops=operation, domain=[
//                            ('qty', '>', 0)],
//                        preferred_domain_list=preferred_domain_list)
//                    Quant.quants_move(quants, move, operation.location_dest_id, location_from=operation.location_id,
//                                      lot_id=False, owner_id=operation.owner_id.id, src_package_id=operation.package_id.id,
//                                      dest_package_id=quant_dest_package_id, entire_pack=entire_pack)
//                else:
//                    # Check what you can do with reserved quants already
//                    qty_on_link = prout_move_qty[move]
//                    rounding = operation.product_id.uom_id.rounding
//                    for reserved_quant in move.reserved_quant_ids:
//                        if (reserved_quant.owner_id.id != operation.owner_id.id) or (reserved_quant.location_id.id != operation.location_id.id) or \
//                                (reserved_quant.package_id.id != operation.package_id.id):
//                            continue
//                        if not reserved_quant.lot_id:
//                            false_quants += [reserved_quant]
//                        elif float_compare(lot_quantities.get(reserved_quant.lot_id.id, 0), 0, precision_rounding=rounding) > 0:
//                            if float_compare(lot_quantities[reserved_quant.lot_id.id], reserved_quant.qty, precision_rounding=rounding) >= 0:
//                                qty_taken = min(
//                                    reserved_quant.qty, qty_on_link)
//                                lot_quantities[reserved_quant.lot_id.id] -= qty_taken
//                                quants_taken += [(reserved_quant, qty_taken)]
//                                qty_on_link -= qty_taken
//                            else:
//                                qty_taken = min(
//                                    qty_on_link, lot_quantities[reserved_quant.lot_id.id])
//                                quants_taken += [(reserved_quant, qty_taken)]
//                                lot_quantities[reserved_quant.lot_id.id] -= qty_taken
//                                qty_on_link -= qty_taken
//                    lot_move_qty[move.id] = qty_on_link
//
//                remaining_move_qty[move.id] -= prout_move_qty[move]
//
//            # Handle lots separately
//            if operation.pack_lot_ids:
//                # TDE FIXME: fix call to move_quants_by_lot to ease understanding
//                self._move_quants_by_lot(
//                    operation, lot_quantities, quants_taken, false_quants, lot_move_qty, quant_dest_package_id)
//
//            # Handle pack in pack
//            if not operation.product_id and operation.package_id and operation.result_package_id.id != operation.package_id.parent_id.id:
//                operation.package_id.sudo().write(
//                    {'parent_id': operation.result_package_id.id})
//        move_dest_ids = set()
//        for move in self:
//            # In case no pack operations in picking
//            if float_compare(remaining_move_qty[move.id], 0, precision_rounding=move.product_id.uom_id.rounding) > 0:
//                # TDE: do in batch ? redone ? check this
//                move.check_tracking(False)
//
//                preferred_domain_list = [[('reservation_id', '=', move.id)], [('reservation_id', '=', False)], [
//                    '&', ('reservation_id', '!=', move.id), ('reservation_id', '!=', False)]]
//                quants = Quant.quants_get_preferred_domain(
//                    remaining_move_qty[move.id], move, domain=[
//                        ('qty', '>', 0)],
//                    preferred_domain_list=preferred_domain_list)
//                Quant.quants_move(
//                    quants, move, move.location_dest_id,
//                    lot_id=move.restrict_lot_id.id, owner_id=move.restrict_partner_id.id)
//
//            # If the move has a destination, add it to the list to reserve
//            if move.move_dest_id and move.move_dest_id.state in ('waiting', 'confirmed'):
//                move_dest_ids.add(move.move_dest_id.id)
//
//            if move.procurement_id:
//                procurements |= move.procurement_id
//
//            # unreserve the quants and make them available for other operations/moves
//            move.quants_unreserve()
//        self.mapped('quant_ids').filtered(lambda quant: quant.package_id and quant.qty > 0).mapped(
//            'package_id')._check_location_constraint()
//        self.write({'state': 'done', 'date': time.strftime(
//            DEFAULT_SERVER_DATETIME_FORMAT)})
//        procurements.check()
//        if move_dest_ids:
//            # TDE FIXME: record setise me
//            self.browse(list(move_dest_ids)).action_assign()
//        pickings.filtered(lambda picking: picking.state == 'done' and not picking.date_done).write(
//            {'date_done': time.strftime(DEFAULT_SERVER_DATETIME_FORMAT)})
//        return True
})
h.StockMove().Methods().Unlink().Extend(
`Unlink`,
func(rs m.StockMoveSet)  {
//        if any(move.state not in ('draft', 'cancel') for move in self):
//            raise UserError(_('You can only delete draft moves.'))
//        return super(StockMove, self).unlink()
})
h.StockMove().Methods().PropagateSplit().DeclareMethod(
`PropagateSplit`,
func(rs m.StockMoveSet, new_move interface{}, qty interface{})  {
//        if self.move_dest_id and self.propagate and self.move_dest_id.state not in ('done', 'cancel'):
//            new_move_prop = self.move_dest_id.split(qty)
//            new_move.write({'move_dest_id': new_move_prop})
})
h.StockMove().Methods().PrepareMoveSplitVals().DeclareMethod(
`PrepareMoveSplitVals`,
func(rs m.StockMoveSet, defaults interface{})  {
//        return defaults
})
h.StockMove().Methods().Split().DeclareMethod(
` Splits qty from move move into a new move

        :param qty: float. quantity to split (given in product UoM)
        :param restrict_lot_id: optional production lot
that can be given in order to force the new move to restrict
its choice of quants to this lot.
        :param restrict_partner_id: optional partner that
can be given in order to force the new move to restrict
its choice of quants to the ones belonging to this partner.
        :param context: dictionay. can contains the special
key 'source_location_id' in order to force the source location
when copying the move
        :returns: id of the backorder move created `,
func(rs m.StockMoveSet, qty interface{}, restrict_lot_id interface{}, restrict_partner_id interface{})  {
//        self = self.with_prefetch(
//        )  # This makes the ORM only look for one record and not 300 at a time, which improves performance
//        if self.state in ('done', 'cancel'):
//            raise UserError(_('You cannot split a move done'))
//        elif self.state == 'draft':
//            # we restrict the split of a draft move because if not confirmed yet, it may be replaced by several other moves in
//            # case of phantom bom (with mrp module). And we don't want to deal with this complexity by copying the product that will explode.
//            raise UserError(
//                _('You cannot split a draft move. It needs to be confirmed first.'))
//        if float_is_zero(qty, precision_rounding=self.product_id.uom_id.rounding) or self.product_qty <= qty:
//            return self.id
//        uom_qty = self.product_id.uom_id._compute_quantity(
//            qty, self.product_uom, rounding_method='HALF-UP')
//        defaults = {
//            'product_uom_qty': uom_qty,
//            'procure_method': 'make_to_stock',
//            'restrict_lot_id': restrict_lot_id,
//            'split_from': self.id,
//            'procurement_id': self.procurement_id.id,
//            'move_dest_id': self.move_dest_id.id,
//            'origin_returned_move_id': self.origin_returned_move_id.id,
//        }
//        defaults = self._prepare_move_split_vals(defaults)
//        if restrict_partner_id:
//            defaults['restrict_partner_id'] = restrict_partner_id
//        if self.env.context.get('source_location_id'):
//            defaults['location_id'] = self.env.context['source_location_id']
//        new_move = self.with_context(rounding_method='HALF-UP').copy(defaults)
//        self.with_context(do_not_propagate=True, rounding_method='HALF-UP').write(
//            {'product_uom_qty': self.product_uom_qty - uom_qty})
//        self._propagate_split(new_move, qty)
//        new_move = new_move.action_confirm()
//        return new_move.id
})
h.StockMove().Methods().ActionShowPicking().DeclareMethod(
`ActionShowPicking`,
func(rs m.StockMoveSet)  {
//        view = self.env.ref('stock.view_picking_form')
//        return {
//            'name': _('Transfer'),
//            'type': 'ir.actions.act_window',
//            'view_type': 'form',
//            'view_mode': 'form',
//            'res_model': 'stock.picking',
//            'views': [(view.id, 'form')],
//            'view_id': view.id,
//            'target': 'new',
//            'res_id': self.picking_id.id}
})
//    show_picking = action_show_picking
h.StockMove().Methods().QuantsUnreserve().DeclareMethod(
`QuantsUnreserve`,
func(rs m.StockMoveSet)  {
//        self.filtered(lambda x: x.partially_available).write(
//            {'partially_available': False})
//        self.mapped('reserved_quant_ids').sudo().write(
//            {'reservation_id': False})
})
}