package stock

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/pool/h"
)

func init() {
	h.StockPackOperation().DeclareModel()

	h.StockPackOperation().Methods().GetDefaultFromLoc().DeclareMethod(
		`GetDefaultFromLoc`,
		func(rs m.StockPackOperationSet) {
			//        default_loc = self.env.context.get('default_location_id')
			//        if default_loc:
			//            return self.env['stock.location'].browse(default_loc).name
		})
	h.StockPackOperation().Methods().GetDefaultToLoc().DeclareMethod(
		`GetDefaultToLoc`,
		func(rs m.StockPackOperationSet) {
			//        default_loc = self.env.context.get('default_location_dest_id')
			//        if default_loc:
			//            return self.env['stock.location'].browse(default_loc).name
		})
	h.StockPackOperation().AddFields(map[string]models.FieldDefinition{
		"PickingId": models.Many2OneField{
			RelationModel: h.StockPicking(),
			String:        "Stock Picking",
			Required:      true,
			Help:          "The stock operation where the packing has been made",
		},
		"ProductId": models.Many2OneField{
			RelationModel: h.ProductProduct(),
			String:        "Product",
			OnDelete:      `cascade`,
		},
		"ProductUomId": models.Many2OneField{
			RelationModel: h.ProductUom(),
			String:        "Unit of Measure",
		},
		"ProductQty": models.FloatField{
			String:  "To Do",
			Default: models.DefaultValue(0),
			//digits=dp.get_precision('Product Unit of Measure')
			Required: true,
		},
		"OrderedQty": models.FloatField{
			String: "Ordered Quantity",
			//digits=dp.get_precision('Product Unit of Measure')
		},
		"QtyDone": models.FloatField{
			String:  "Done",
			Default: models.DefaultValue(0),
			//digits=dp.get_precision('Product Unit of Measure')
		},
		"QtyDoneUomOrdered": models.FloatField{
			String: "Quantity Done",
			//digits=dp.get_precision('Product Unit of Measure')
			Compute: h.StockPackOperation().Methods().ComputeQtyDoneUomOrdered(),
			Help:    "Quantity done in UOM ordered",
		},
		"IsDone": models.BooleanField{
			Compute:  h.StockPackOperation().Methods().ComputeIsDone(),
			String:   "Done",
			ReadOnly: false,
			//oldname='processed_boolean'
		},
		"PackageId": models.Many2OneField{
			RelationModel: h.StockQuantPackage(),
			String:        "Source Package",
		},
		"PackLotIds": models.One2ManyField{
			RelationModel: h.StockPackOperationLot(),
			ReverseFK:     "",
			String:        "Lots/Serial Numbers Used",
		},
		"ResultPackageId": models.Many2OneField{
			RelationModel: h.StockQuantPackage(),
			String:        "Destination Package",
			OnDelete:      `cascade`,
			Required:      false,
			Help:          "If set, the operations are packed into this package",
		},
		"Date": models.DateTimeField{
			String:   "Date",
			Default:  func(env models.Environment) interface{} { return odoo.fields.Date.context_today },
			Required: true,
		},
		"OwnerId": models.Many2OneField{
			RelationModel: h.Partner(),
			String:        "Owner",
			Help:          "Owner of the quants",
		},
		"LinkedMoveOperationIds": models.One2ManyField{
			RelationModel: h.StockMoveOperationLink(),
			ReverseFK:     "",
			String:        "Linked Moves",
			ReadOnly:      true,
			Help: "Moves impacted by this operation for the computation of" +
				"the remaining quantities",
		},
		"RemainingQty": models.FloatField{
			Compute: h.StockPackOperation().Methods().GetRemainingQty(),
			String:  "Remaining Qty",
			//digits=0
			Help: "Remaining quantity in default UoM according to moves matched" +
				"with this operation.",
		},
		"LocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Source Location",
			Required:      true,
		},
		"LocationDestId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Destination Location",
			Required:      true,
		},
		"PickingSourceLocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			Related:       `PickingId.LocationId`,
		},
		"PickingDestinationLocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			Related:       `PickingId.LocationDestId`,
		},
		"FromLoc": models.CharField{
			Compute: h.StockPackOperation().Methods().ComputeLocationDescription(),
			Default: models.DefaultValue(_get_default_from_loc),
			String:  "From",
		},
		"ToLoc": models.CharField{
			Compute: h.StockPackOperation().Methods().ComputeLocationDescription(),
			Default: models.DefaultValue(_get_default_to_loc),
			String:  "To",
		},
		"FreshRecord": models.BooleanField{
			String:  "Newly created pack operation",
			Default: models.DefaultValue(true),
		},
		"LotsVisible": models.BooleanField{
			Compute: h.StockPackOperation().Methods().ComputeLotsVisible(),
		},
		"State": models.SelectionField{
			//selection=[('draft', 'Draft'),('cancel', 'Cancelled'),('waiting', 'Waiting Another Operation'),('confirmed', 'Waiting Availability'),('partially_available', 'Partially Available'),('assigned', 'Available'),('done', 'Done')]
			Related: `PickingId.State`,
		},
	})
	h.StockPackOperation().Methods().ComputeIsDone().DeclareMethod(
		`ComputeIsDone`,
		func(rs h.StockPackOperationSet) h.StockPackOperationData {
			//        self.is_done = self.qty_done > 0.0
		})
	h.StockPackOperation().Methods().OnChangeIsDone().DeclareMethod(
		`OnChangeIsDone`,
		func(rs m.StockPackOperationSet) {
			//        if not self.product_id:
			//            if self.is_done and self.qty_done == 0:
			//                self.qty_done = 1.0
			//            if not self.is_done and self.qty_done != 0:
			//                self.qty_done = 0.0
		})
	h.StockPackOperation().Methods().GetRemainingProdQuantities().DeclareMethod(
		`Get the remaining quantities per product on an operation
with a package. This function returns a dictionary`,
		func(rs m.StockPackOperationSet) {
			//        if not self.package_id or self.product_id:
			//            return {self.product_id: self.remaining_qty}
			//        res = self.package_id._get_all_products_quantities()
			//        for record in self.linked_move_operation_ids:
			//            if record.move_id.product_id not in res:
			//                res[record.move_id.product_id] = 0
			//            res[record.move_id.product_id] -= record.qty
			//        return res
		})
	h.StockPackOperation().Methods().GetRemainingQty().DeclareMethod(
		`GetRemainingQty`,
		func(rs h.StockPackOperationSet) h.StockPackOperationData {
			//        if self.package_id and not self.product_id:
			//            # dont try to compute the remaining quantity for packages because it's not relevant (a package could include different products).
			//            # should use _get_remaining_prod_quantities instead
			//            # TDE FIXME: actually resolve the comment hereabove
			//            self.remaining_qty = 0
			//        else:
			//            qty = self.product_qty
			//            if self.product_uom_id:
			//                qty = self.product_uom_id._compute_quantity(
			//                    self.product_qty, self.product_id.uom_id)
			//            for record in self.linked_move_operation_ids:
			//                qty -= record.qty
			//            self.remaining_qty = float_round(
			//                qty, precision_rounding=self.product_id.uom_id.rounding)
		})
	h.StockPackOperation().Methods().ComputeLocationDescription().DeclareMethod(
		`ComputeLocationDescription`,
		func(rs h.StockPackOperationSet) h.StockPackOperationData {
			//        for operation, operation_sudo in zip(self, self.sudo()):
			//            operation.from_loc = '%s%s' % (
			//                operation_sudo.location_id.name, operation.product_id and operation_sudo.package_id.name or '')
			//            operation.to_loc = '%s%s' % (
			//                operation_sudo.location_dest_id.name, operation_sudo.result_package_id.name or '')
		})
	h.StockPackOperation().Methods().ComputeLotsVisible().DeclareMethod(
		`ComputeLotsVisible`,
		func(rs h.StockPackOperationSet) h.StockPackOperationData {
			//        if self.pack_lot_ids:
			//            self.lots_visible = True
			//        # TDE FIXME: not sure correctly migrated
			//        elif self.picking_id.picking_type_id and self.product_id.tracking != 'none':
			//            picking = self.picking_id
			//            self.lots_visible = picking.picking_type_id.use_existing_lots or picking.picking_type_id.use_create_lots
			//        else:
			//            self.lots_visible = self.product_id.tracking != 'none'
		})
	h.StockPackOperation().Methods().ComputeQtyDoneUomOrdered().DeclareMethod(
		`ComputeQtyDoneUomOrdered`,
		func(rs h.StockPackOperationSet) h.StockPackOperationData {
			//        for pack in self:
			//            if pack.product_uom_id and pack.linked_move_operation_ids:
			//                pack.qty_done_uom_ordered = pack.product_uom_id._compute_quantity(
			//                    pack.qty_done, pack.linked_move_operation_ids[0].move_id.product_uom)
			//            else:
			//                pack.qty_done_uom_ordered = pack.qty_done
		})
	h.StockPackOperation().Methods().OnchangePacklots().DeclareMethod(
		`OnchangePacklots`,
		func(rs m.StockPackOperationSet) {
			//        self.qty_done = sum([x.qty for x in self.pack_lot_ids])
		})
	h.StockPackOperation().Methods().OnchangeProductId().DeclareMethod(
		`OnchangeProductId`,
		func(rs m.StockPackOperationSet) {
			//        if self.product_id:
			//            self.lots_visible = self.product_id.tracking != 'none'
			//            if not self.product_uom_id or self.product_uom_id.category_id != self.product_id.uom_id.category_id:
			//                self.product_uom_id = self.product_id.uom_id.id
			//            res = {'domain': {'product_uom_id': [
			//                ('category_id', '=', self.product_uom_id.category_id.id)]}}
			//        else:
			//            res = {'domain': {'product_uom_id': []}}
			//        return res
		})
	h.StockPackOperation().Methods().Create().Extend(
		`Create`,
		func(rs m.StockPackOperationSet, vals models.RecordData) {
			//        vals['ordered_qty'] = vals.get('product_qty')
			//        return super(PackOperation, self).create(vals)
		})
	h.StockPackOperation().Methods().Write().Extend(
		`Write`,
		func(rs m.StockPackOperationSet, values models.RecordData) {
			//        values['fresh_record'] = False
			//        return super(PackOperation, self).write(values)
		})
	h.StockPackOperation().Methods().Unlink().Extend(
		`Unlink`,
		func(rs m.StockPackOperationSet) {
			//        if any([operation.state in ('done', 'cancel') for operation in self]):
			//            raise UserError(
			//                _('You can not delete pack operations of a done picking'))
			//        return super(PackOperation, self).unlink()
		})
	h.StockPackOperation().Methods().SplitQuantities().DeclareMethod(
		`SplitQuantities`,
		func(rs m.StockPackOperationSet) {
			//        for operation in self:
			//            if float_compare(operation.product_qty, operation.qty_done, precision_rounding=operation.product_uom_id.rounding) == 1:
			//                cpy = operation.copy(default={
			//                                     'qty_done': 0.0, 'product_qty': operation.product_qty - operation.qty_done})
			//                operation.write({'product_qty': operation.qty_done})
			//                operation._copy_remaining_pack_lot_ids(cpy)
			//            else:
			//                raise UserError(
			//                    _('The quantity to split should be smaller than the quantity To Do.  '))
			//        return True
		})
	h.StockPackOperation().Methods().Save().DeclareMethod(
		`Save`,
		func(rs m.StockPackOperationSet) {
			//        for pack in self:
			//            if pack.product_id.tracking != 'none':
			//                pack.write({'qty_done': sum(pack.pack_lot_ids.mapped('qty'))})
			//        return {'type': 'ir.actions.act_window_close'}
		})
	h.StockPackOperation().Methods().ActionSplitLots().DeclareMethod(
		`ActionSplitLots`,
		func(rs m.StockPackOperationSet) {
			//        action_ctx = dict(self.env.context)
			//        returned_move = self.linked_move_operation_ids.mapped(
			//            'move_id').mapped('origin_returned_move_id')
			//        picking_type = self.picking_id.picking_type_id
			//        action_ctx.update({
			//            'serial': self.product_id.tracking == 'serial',
			//            'only_create': picking_type.use_create_lots and not picking_type.use_existing_lots and not returned_move,
			//            'create_lots': picking_type.use_create_lots,
			//            'state_done': self.picking_id.state == 'done',
			//            'show_reserved': any([lot for lot in self.pack_lot_ids if lot.qty_todo > 0.0])})
			//        view_id = self.env.ref('stock.view_pack_operation_lot_form').id
			//        return {
			//            'name': _('Lot/Serial Number Details'),
			//            'type': 'ir.actions.act_window',
			//            'view_type': 'form',
			//            'view_mode': 'form',
			//            'res_model': 'stock.pack.operation',
			//            'views': [(view_id, 'form')],
			//            'view_id': view_id,
			//            'target': 'new',
			//            'res_id': self.ids[0],
			//            'context': action_ctx}
		})
	//    split_lot = action_split_lots
	h.StockPackOperation().Methods().ShowDetails().DeclareMethod(
		`ShowDetails`,
		func(rs m.StockPackOperationSet) {
			//        view_id = self.env.ref(
			//            'stock.view_pack_operation_details_form_save').id
			//        return {
			//            'name': _('Operation Details'),
			//            'type': 'ir.actions.act_window',
			//            'view_type': 'form',
			//            'view_mode': 'form',
			//            'res_model': 'stock.pack.operation',
			//            'views': [(view_id, 'form')],
			//            'view_id': view_id,
			//            'target': 'new',
			//            'res_id': self.ids[0],
			//            'context': self.env.context}
		})
	h.StockPackOperation().Methods().CheckSerialNumber().DeclareMethod(
		`CheckSerialNumber`,
		func(rs m.StockPackOperationSet) {
			//        for operation in self:
			//            if operation.picking_id and \
			//                    (operation.picking_id.picking_type_id.use_existing_lots or operation.picking_id.picking_type_id.use_create_lots) and \
			//                    operation.product_id and operation.product_id.tracking != 'none' and \
			//                    operation.qty_done > 0.0:
			//                if not operation.pack_lot_ids:
			//                    raise UserError(
			//                        _('You need to provide a Lot/Serial Number for product %s') % operation.product_id.name)
			//                if operation.product_id.tracking == 'serial':
			//                    for opslot in operation.pack_lot_ids:
			//                        if opslot.qty not in (1.0, 0.0):
			//                            raise UserError(
			//                                _('You should provide a different serial number for each piece'))
		})
	//    check_tracking = _check_serial_number
	h.StockPackOperation().Methods().CopyRemainingPackLotIds().DeclareMethod(
		`CopyRemainingPackLotIds`,
		func(rs m.StockPackOperationSet, new_operation interface{}) {
			//        for op in self:
			//            for lot in op.pack_lot_ids:
			//                new_qty_todo = lot.qty_todo - lot.qty
			//
			//                if float_compare(new_qty_todo, 0, precision_rounding=op.product_uom_id.rounding) > 0:
			//                    lot.copy({
			//                        'operation_id': new_operation.id,
			//                        'qty_todo': new_qty_todo,
			//                        'qty': 0,
			//                    })
		})
	h.StockPackOperationLot().DeclareModel()
	h.StockPackOperationLot().AddSQLConstraint("qty", "CHECK(qty >= 0.0)", "Quantity must be greater than or equal to 0.0!")
	h.StockPackOperationLot().AddSQLConstraint("uniq_lot_id", "unique(operation_id, lot_id)", "You have already mentioned this lot in another line")
	h.StockPackOperationLot().AddSQLConstraint("uniq_lot_name", "unique(operation_id, lot_name)", "You have already mentioned this lot name in another line")

	h.StockPackOperationLot().AddFields(map[string]models.FieldDefinition{
		"OperationId": models.Many2OneField{
			RelationModel: h.StockPackOperation(),
		},
		"Qty": models.FloatField{
			String:  "Done",
			Default: models.DefaultValue(1),
			//digits=dp.get_precision('Product Unit of Measure')
		},
		"LotId": models.Many2OneField{
			RelationModel: h.StockProductionLot(),
			String:        "Lot/Serial Number",
		},
		"LotName": models.CharField{
			String: "Lot/Serial Number",
		},
		"QtyTodo": models.FloatField{
			String:  "To Do",
			Default: models.DefaultValue(0),
			//digits=dp.get_precision('Product Unit of Measure')
		},
		"PlusVisible": models.BooleanField{
			Compute: h.StockPackOperationLot().Methods().ComputePlusVisible(),
			Default: models.DefaultValue(true),
		},
	})
	h.StockPackOperationLot().Methods().ComputePlusVisible().DeclareMethod(
		`ComputePlusVisible`,
		func(rs h.StockPackOperationLotSet) h.StockPackOperationLotData {
			//        if self.operation_id.product_id.tracking == 'serial':
			//            self.plus_visible = (self.qty == 0.0)
			//        else:
			//            self.plus_visible = (self.qty_todo == 0.0) or (
			//                self.qty < self.qty_todo)
		})
	h.StockPackOperationLot().Methods().CheckLot().DeclareMethod(
		`CheckLot`,
		func(rs m.StockPackOperationLotSet) {
			//        if any(not lot.lot_name and not lot.lot_id for lot in self):
			//            raise ValidationError(_('Lot/Serial Number required'))
			//        return True
		})
	h.StockPackOperationLot().Methods().ActionAddQuantity().DeclareMethod(
		`ActionAddQuantity`,
		func(rs m.StockPackOperationLotSet, quantity interface{}) {
			//        for lot in self:
			//            lot.write({'qty': lot.qty + quantity})
			//            lot.operation_id.write({'qty_done': sum(
			//                operation_lot.qty for operation_lot in lot.operation_id.pack_lot_ids)})
			//        return self.mapped('operation_id').action_split_lots()
		})
	h.StockPackOperationLot().Methods().DoPlus().DeclareMethod(
		`DoPlus`,
		func(rs m.StockPackOperationLotSet) {
			//        return self.action_add_quantity(1)
		})
	h.StockPackOperationLot().Methods().DoMinus().DeclareMethod(
		`DoMinus`,
		func(rs m.StockPackOperationLotSet) {
			//        return self.action_add_quantity(-1)
		})
	h.StockMoveOperationLink().DeclareModel()

	h.StockMoveOperationLink().AddFields(map[string]models.FieldDefinition{
		"Qty": models.FloatField{
			String: "Quantity",
			Help: "Quantity of products to consider when talking about the" +
				"contribution of this pack operation towards the remaining" +
				"quantity of the move (and inverse). Given in the product main uom.",
		},
		"OperationId": models.Many2OneField{
			RelationModel: h.StockPackOperation(),
			String:        "Operation",
			OnDelete:      `cascade`,
			Required:      true,
		},
		"MoveId": models.Many2OneField{
			RelationModel: h.StockMove(),
			String:        "Move",
			OnDelete:      `cascade`,
			Required:      true,
		},
		"ReservedQuantId": models.Many2OneField{
			RelationModel: h.StockQuant(),
			String:        "Reserved Quant",
			Help: "Technical field containing the quant that created this" +
				"link between an operation and a stock move. Used at the" +
				"stock_move_obj.action_done() time to avoid seeking a matching" +
				"quant again",
		},
	})
}
