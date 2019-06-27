package stock

import (
	"github.com/hexya-erp/hexya-base/web/webdata"
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/pool/h"
)

//import logging
//_logger = logging.getLogger(__name__)
func init() {
	h.StockQuant().DeclareModel()

	h.StockQuant().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{
			String:  "Identifier",
			Compute: h.StockQuant().Methods().ComputeName(),
		},
		"ProductId": models.Many2OneField{
			RelationModel: h.ProductProduct(),
			String:        "Product",
			Index:         true,
			OnDelete:      `restrict`,
			ReadOnly:      true,
			Required:      true,
		},
		"LocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Location",

			Index:    true,
			OnDelete: `restrict`,
			ReadOnly: true,
			Required: true,
		},
		"Qty": models.FloatField{
			String:   "Quantity",
			Index:    true,
			ReadOnly: true,
			Required: true,
			Help: "Quantity of products in this quant, in the default unit" +
				"of measure of the product",
		},
		"ProductUomId": models.Many2OneField{
			RelationModel: h.ProductUom(),
			String:        "Unit of Measure",
			Related:       `ProductId.UomId`,
			ReadOnly:      true,
		},
		"PackageId": models.Many2OneField{
			RelationModel: h.StockQuantPackage(),
			String:        "Package",
			Index:         true,
			ReadOnly:      true,
			Help:          "The package containing this quant",
		},
		"PackagingTypeId": models.Many2OneField{
			RelationModel: h.ProductPackaging(),
			String:        "Type of packaging",
			Related:       `PackageId.PackagingId`,
			ReadOnly:      true,
			Stored:        true,
		},
		"ReservationId": models.Many2OneField{
			RelationModel: h.StockMove(),
			String:        "Reserved for Move",
			Index:         true,
			ReadOnly:      true,
			Help:          "The move the quant is reserved for",
		},
		"LotId": models.Many2OneField{
			RelationModel: h.StockProductionLot(),
			String:        "Lot/Serial Number",
			Index:         true,
			OnDelete:      `restrict`,
			ReadOnly:      true,
		},
		"Cost": models.FloatField{
			String: "Unit Cost",
			//group_operator='avg'
		},
		"OwnerId": models.Many2OneField{
			RelationModel: h.Partner(),
			String:        "Owner",
			Index:         true,
			ReadOnly:      true,
			Help:          "This is the owner of the quant",
		},
		"InDate": models.DateTimeField{
			String:   "Incoming Date",
			Index:    true,
			ReadOnly: true,
		},
		"HistoryIds": models.Many2ManyField{
			RelationModel:    h.StockMove(),
			M2MLinkModelName: "",
			M2MOurField:      "",
			M2MTheirField:    "",
			String:           "Moves",
			NoCopy:           true,
			Help:             "Moves that operate(d) on this quant",
		},
		"CompanyId": models.Many2OneField{
			RelationModel: h.Company(),
			String:        "Company",
			Index:         true,
			ReadOnly:      true,
			Required:      true,
			Default:       func(env models.Environment) interface{} { return env["res.company"]._company_default_get() },
			Help:          "The company to which the quants belong",
		},
		"InventoryValue": models.FloatField{
			String:   "Inventory Value",
			Compute:  h.StockQuant().Methods().ComputeInventoryValue(),
			ReadOnly: true,
		},
		"PropagatedFromId": models.Many2OneField{
			RelationModel: h.StockQuant(),
			String:        "Linked Quant",
			Index:         true,
			ReadOnly:      true,
			Help:          "The negative quant this is coming from",
		},
		"NegativeMoveId": models.Many2OneField{
			RelationModel: h.StockMove(),
			String:        "Move Negative Quant",
			ReadOnly:      true,
			Help: "If this is a negative quant, this will be the move that" +
				"caused this negative quant.",
		},
		"NegativeDestLocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Negative Destination Location",
			Related:       `NegativeMoveId.LocationDestId`,
			ReadOnly:      true,
			Help: "Technical field used to record the destination location" +
				"of a move that created a negative quant",
		},
	})
	h.StockQuant().Fields().CreateDate().setString("Creation Date")
	h.StockQuant().Fields().CreateDate().setReadOnly(true)
	h.StockQuant().Methods().ComputeName().DeclareMethod(
		` Forms complete name of location from parent location to
child location. `,
		func(rs h.StockQuantSet) h.StockQuantData {
			//        self.name = '%s: %s%s' % (
			//            self.lot_id.name or self.product_id.code or '', self.qty, self.product_id.uom_id.name)
		})
	h.StockQuant().Methods().ComputeInventoryValue().DeclareMethod(
		`ComputeInventoryValue`,
		func(rs h.StockQuantSet) h.StockQuantData {
			//        for quant in self:
			//            if quant.company_id != self.env.user.company_id:
			//                # if the company of the quant is different than the current user company, force the company in the context
			//                # then re-do a browse to read the property fields for the good company.
			//                quant = quant.with_context(force_company=quant.company_id.id)
			//            quant.inventory_value = quant.product_id.standard_price * quant.qty
		})
	h.StockQuant().Methods().Init().DeclareMethod(
		`Init`,
		func(rs m.StockQuantSet) {
			//        self._cr.execute('SELECT indexname FROM pg_indexes WHERE indexname = %s',
			//                         ('stock_quant_product_location_index'))
			//        if not self._cr.fetchone():
			//            self._cr.execute(
			//                'CREATE INDEX stock_quant_product_location_index ON stock_quant (product_id, location_id, company_id, qty, in_date, reservation_id)')
		})
	h.StockQuant().Methods().Unlink().Extend(
		`Unlink`,
		func(rs m.StockQuantSet) {
			//        if not self.env.context.get('force_unlink'):
			//            raise UserError(
			//                _('Under no circumstances should you delete or change quants yourselves!'))
			//        return super(Quant, self).unlink()
		})
	h.StockQuant().Methods().ReadGroup().Extend(
		` Overwrite the read_group in order to sum the function
field 'inventory_value' in group by `,
		func(rs m.StockQuantSet, domain webdata.ReadGroupParams, fields interface{}, groupby interface{}, offset interface{}, limit interface{}, orderby interface{}, lazy interface{}) {
			//        res = super(Quant, self).read_group(domain, fields, groupby,
			//                                            offset=offset, limit=limit, orderby=orderby, lazy=lazy)
			//        if 'inventory_value' in fields:
			//            for line in res:
			//                lines = self.search(line.get('__domain', domain))
			//                inv_value = 0.0
			//                for line2 in lines:
			//                    inv_value += line2.inventory_value
			//                line['inventory_value'] = inv_value
			//        return res
		})
	h.StockQuant().Methods().ActionViewQuantHistory().DeclareMethod(
		` Returns an action that display the history of the quant, which
        mean all the stock moves that lead to this quant
creation with this
        quant quantity. `,
		func(rs m.StockQuantSet) {
			//        action = self.env.ref('stock', 'stock_move_action').read()[0]
			//        action['domain'] = [('id', 'in', self.mapped('history_ids').ids)]
			//        return action
		})
	h.StockQuant().Methods().QuantsReserve().DeclareMethod(
		` This function reserves quants for the given move and optionally
        given link. If the total of quantity reserved is
enough, the move state
        is also set to 'assigned'

        :param quants: list of tuple(quant browse record
or None, qty to reserve). If None is given as first tuple
element, the item will be ignored. Negative quants should
not be received as argument
        :param move: browse record
        :param link: browse record (stock.move.operation.link)
        `,
		func(rs m.StockQuantSet, quants interface{}, move interface{}, link interface{}) {
			//        quants_to_reserve_sudo = self.env['stock.quant'].sudo()
			//        reserved_availability = move.reserved_availability
			//        for quant, qty in quants:
			//            if qty <= 0.0 or (quant and quant.qty <= 0.0):
			//                raise UserError(
			//                    _('You can not reserve a negative quantity or a negative quant.'))
			//            if not quant:
			//                continue
			//            quant._quant_split(qty)
			//            quants_to_reserve_sudo |= quant
			//            reserved_availability += quant.qty
			//        if quants_to_reserve_sudo:
			//            quants_to_reserve_sudo.write({'reservation_id': move.id})
			//        rounding = move.product_id.uom_id.rounding
			//        if float_compare(reserved_availability, move.product_qty, precision_rounding=rounding) == 0 and move.state in ('confirmed', 'waiting'):
			//            move.write({'state': 'assigned'})
			//        elif float_compare(reserved_availability, 0, precision_rounding=rounding) > 0 and not move.partially_available:
			//            move.write({'partially_available': True})
		})
	h.StockQuant().Methods().QuantsMove().DeclareMethod(
		`Moves all given stock.quant in the given destination location.
 Unreserve from current move.
        :param quants: list of tuple(browse record(stock.quant)
or None, quantity to move)
        :param move: browse record (stock.move)
        :param location_to: browse record (stock.location)
depicting where the quants have to be moved
        :param location_from: optional browse record (stock.location)
explaining where the quant has to be taken
                              (may differ from the move
source location in case a removal strategy applied).
                              This parameter is only used
to pass to _quant_create_from_move if a negative quant must be created
        :param lot_id: ID of the lot that must be set on
the quants to move
        :param owner_id: ID of the partner that must own
the quants to move
        :param src_package_id: ID of the package that contains
the quants to move
        :param dest_package_id: ID of the package that
must be set on the moved quant
        `,
		func(rs m.StockQuantSet, quants interface{}, move interface{}, location_to interface{}, location_from interface{}, lot_id interface{}, owner_id interface{}, src_package_id interface{}, dest_package_id interface{}, entire_pack interface{}) {
			//        if location_to.usage == 'view':
			//            raise UserError(
			//                _('You cannot move to a location of type view %s.') % (location_to.name))
			//        quants_reconcile_sudo = self.env['stock.quant'].sudo()
			//        quants_move_sudo = self.env['stock.quant'].sudo()
			//        check_lot = False
			//        for quant, qty in quants:
			//            if not quant:
			//                # If quant is None, we will create a quant to move (and potentially a negative counterpart too)
			//                quant = self._quant_create_from_move(
			//                    qty, move, lot_id=lot_id, owner_id=owner_id, src_package_id=src_package_id, dest_package_id=dest_package_id, force_location_from=location_from, force_location_to=location_to)
			//                check_lot = True
			//            else:
			//                quant._quant_split(qty)
			//                quants_move_sudo |= quant
			//            quants_reconcile_sudo |= quant
			//        if quants_move_sudo:
			//            moves_recompute = quants_move_sudo.filtered(
			//                lambda self: self.reservation_id != move).mapped('reservation_id')
			//            quants_move_sudo._quant_update_from_move(
			//                move, location_to, dest_package_id, lot_id=lot_id, entire_pack=entire_pack)
			//            moves_recompute.recalculate_move_state()
			//        if location_to.usage == 'internal':
			//            # Do manual search for quant to avoid full table scan (order by id)
			//            self._cr.execute("""
			//                SELECT 0 FROM stock_quant, stock_location WHERE product_id = %s AND stock_location.id = stock_quant.location_id AND
			//                ((stock_location.parent_left >= %s AND stock_location.parent_left < %s) OR stock_location.id = %s) AND qty < 0.0 LIMIT 1
			//            """, (move.product_id.id, location_to.parent_left, location_to.parent_right, location_to.id))
			//            if self._cr.fetchone():
			//                quants_reconcile_sudo._quant_reconcile_negative(move)
			//        picking_type = move.picking_id and move.picking_id.picking_type_id or False
			//        if check_lot and lot_id and move.product_id.tracking == 'serial' and (not picking_type or (picking_type.use_create_lots or picking_type.use_existing_lots)):
			//            other_quants = self.search([('product_id', '=', move.product_id.id), ('lot_id', '=', lot_id),
			//                                        ('qty', '>', 0.0), ('location_id.usage', '=', 'internal')])
			//            if other_quants:
			//                # We raise an error if:
			//                # - the total quantity is strictly larger than 1.0
			//                # - there are more than one negative quant, to avoid situations where the user would
			//                #   force the quantity at several steps of the process
			//                if sum(other_quants.mapped('qty')) > 1.0 or len([q for q in other_quants.mapped('qty') if q < 0]) > 1:
			//                    lot_name = self.env['stock.production.lot'].browse(
			//                        lot_id).name
			//                    raise UserError(_('The serial number %s is already in stock.') %
			//                                    lot_name + _("Otherwise make sure the right stock/owner is set."))
		})
	h.StockQuant().Methods().QuantCreateFromMove().DeclareMethod(
		`Create a quant in the destination location and create a negative
        quant in the source location if it's an internal location. `,
		func(rs m.StockQuantSet, qty interface{}, move interface{}, lot_id interface{}, owner_id interface{}, src_package_id interface{}, dest_package_id interface{}, force_location_from interface{}, force_location_to interface{}) {
			//        price_unit = move.get_price_unit()
			//        location = force_location_to or move.location_dest_id
			//        rounding = move.product_id.uom_id.rounding
			//        vals = {
			//            'product_id': move.product_id.id,
			//            'location_id': location.id,
			//            'qty': float_round(qty, precision_rounding=rounding),
			//            'cost': price_unit,
			//            'history_ids': [(4, move.id)],
			//            'in_date': datetime.now().strftime(DEFAULT_SERVER_DATETIME_FORMAT),
			//            'company_id': move.company_id.id,
			//            'lot_id': lot_id,
			//            'owner_id': owner_id,
			//            'package_id': dest_package_id,
			//        }
			//        if move.location_id.usage == 'internal':
			//            # if we were trying to move something from an internal location and reach here (quant creation),
			//            # it means that a negative quant has to be created as well.
			//            negative_vals = vals.copy()
			//            negative_vals['location_id'] = force_location_from and force_location_from.id or move.location_id.id
			//            negative_vals['qty'] = float_round(-qty,
			//                                               precision_rounding=rounding)
			//            negative_vals['cost'] = price_unit
			//            negative_vals['negative_move_id'] = move.id
			//            negative_vals['package_id'] = src_package_id
			//            negative_quant_id = self.sudo().create(negative_vals)
			//            vals.update({'propagated_from_id': negative_quant_id.id})
			//        picking_type = move.picking_id and move.picking_id.picking_type_id or False
			//        if lot_id and move.product_id.tracking == 'serial' and (not picking_type or (picking_type.use_create_lots or picking_type.use_existing_lots)):
			//            if qty != 1.0:
			//                raise UserError(
			//                    _('You should only receive by the piece with the same serial number'))
			//        return self.sudo().create(vals)
		})
	h.StockQuant().Methods().QuantCreate().DeclareMethod(
		`QuantCreate`,
		func(rs m.StockQuantSet, qty interface{}, move interface{}, lot_id interface{}, owner_id interface{}, src_package_id interface{}, dest_package_id interface{}, force_location_from interface{}, force_location_to interface{}) {
			//        _logger.warning(
			//            "'_quant_create' has been renamed into '_quant_create_from_move'... Overrides are ignored")
			//        return self._quant_create_from_move(
			//            qty, move, lot_id=lot_id, owner_id=owner_id,
			//            src_package_id=src_package_id, dest_package_id=dest_package_id,
			//            force_location_from=force_location_from, force_location_to=force_location_to)
		})
	h.StockQuant().Methods().QuantUpdateFromMove().DeclareMethod(
		`QuantUpdateFromMove`,
		func(rs m.StockQuantSet, move interface{}, location_dest_id interface{}, dest_package_id interface{}, lot_id interface{}, entire_pack interface{}) {
			//        vals = {
			//            'location_id': location_dest_id.id,
			//            'history_ids': [(4, move.id)],
			//            'reservation_id': False}
			//        if lot_id and any(quant for quant in self if not quant.lot_id.id):
			//            vals['lot_id'] = lot_id
			//        if not entire_pack:
			//            vals.update({'package_id': dest_package_id})
			//        self.write(vals)
		})
	h.StockQuant().Methods().MoveQuantsWrite().DeclareMethod(
		`MoveQuantsWrite`,
		func(rs m.StockQuantSet, move interface{}, location_dest_id interface{}, dest_package_id interface{}, lot_id interface{}, entire_pack interface{}) {
			//        _logger.warning(
			//            "'move_quants_write' has been renamed into '_quant_update_from_move'... Overrides are ignored")
			//        return self._quant_update_from_move(move, location_dest_id, dest_package_id, lot_id=lot_id, entire_pack=entire_pack)
		})
	h.StockQuant().Methods().QuantReconcileNegative().DeclareMethod(
		`
            When new quant arrive in a location, try to
reconcile it with
            negative quants. If it's possible, apply the cost of the new
            quant to the counterpart of the negative quant.
        `,
		func(rs m.StockQuantSet, move interface{}) {
			//        solving_quant = self
			//        quants = self._search_quants_to_reconcile()
			//        product_uom_rounding = self.product_id.uom_id.rounding
			//        for quant_neg, qty in quants:
			//            if not quant_neg or not solving_quant:
			//                continue
			//            quants_to_solve = self.search(
			//                [('propagated_from_id', '=', quant_neg.id)])
			//            if not quants_to_solve:
			//                continue
			//            solving_qty = qty
			//            solved_quants = self.env['stock.quant'].sudo()
			//            for to_solve_quant in quants_to_solve:
			//                if float_compare(solving_qty, 0, precision_rounding=product_uom_rounding) <= 0:
			//                    continue
			//                solved_quants |= to_solve_quant
			//                to_solve_quant._quant_split(
			//                    min(solving_qty, to_solve_quant.qty))
			//                solving_qty -= min(solving_qty, to_solve_quant.qty)
			//            remaining_solving_quant = solving_quant._quant_split(qty)
			//            remaining_neg_quant = quant_neg._quant_split(-qty)
			//            # if the reconciliation was not complete, we need to link together the remaining parts
			//            if remaining_neg_quant:
			//                remaining_to_solves = self.sudo().search(
			//                    [('propagated_from_id', '=', quant_neg.id), ('id', 'not in', solved_quants.ids)])
			//                if remaining_to_solves:
			//                    remaining_to_solves.write(
			//                        {'propagated_from_id': remaining_neg_quant.id})
			//            if solving_quant.propagated_from_id and solved_quants:
			//                solved_quants.write(
			//                    {'propagated_from_id': solving_quant.propagated_from_id.id})
			//            # delete the reconciled quants, as it is replaced by the solved quants
			//            quant_neg.sudo().with_context(force_unlink=True).unlink()
			//            if solved_quants:
			//                # price update + accounting entries adjustments
			//                solved_quants._price_update(solving_quant.cost)
			//                # merge history (and cost?)
			//                solved_quants.write(solving_quant._prepare_history())
			//            solving_quant.with_context(force_unlink=True).unlink()
			//            solving_quant = remaining_solving_quant
		})
	h.StockQuant().Methods().PrepareHistory().DeclareMethod(
		`PrepareHistory`,
		func(rs m.StockQuantSet) {
			//        return {
			//            'history_ids': [(4, history_move.id) for history_move in self.history_ids],
			//        }
		})
	h.StockQuant().Methods().PriceUpdate().DeclareMethod(
		`PriceUpdate`,
		func(rs m.StockQuantSet, newprice interface{}) {
			//        self.sudo().write({'cost': newprice})
		})
	h.StockQuant().Methods().SearchQuantsToReconcile().DeclareMethod(
		` Searches negative quants to reconcile for where the quant
to reconcile is put `,
		func(rs m.StockQuantSet) {
			//        dom = ['&', '&', '&', '&',
			//               ('qty', '<', 0),
			//               ('location_id', 'child_of', self.location_id.id),
			//               ('product_id', '=', self.product_id.id),
			//               ('owner_id', '=', self.owner_id.id),
			//               # Do not let the quant eat itself, or it will kill its history (e.g. returns / Stock -> Stock)
			//               ('id', '!=', self.propagated_from_id.id)]
			//        if self.package_id.id:
			//            dom = ['&'] + dom + [('package_id', '=', self.package_id.id)]
			//        if self.lot_id:
			//            dom = ['&'] + dom + \
			//                ['|', ('lot_id', '=', False), ('lot_id', '=', self.lot_id.id)]
			//            order = 'lot_id, in_date'
			//        else:
			//            order = 'in_date'
			//        rounding = self.product_id.uom_id.rounding
			//        quants = []
			//        quantity = self.qty
			//        for quant in self.search(dom, order=order):
			//            if float_compare(quantity, abs(quant.qty), precision_rounding=rounding) >= 0:
			//                quants += [(quant, abs(quant.qty))]
			//                quantity -= abs(quant.qty)
			//            elif float_compare(quantity, 0.0, precision_rounding=rounding) != 0:
			//                quants += [(quant, quantity)]
			//                quantity = 0
			//                break
			//        return quants
		})
	h.StockQuant().Methods().QuantsGetPreferredDomain().DeclareMethod(
		` This function tries to find quants for the given domain
and move/ops, by trying to first limit
            the choice on the quants that match the first
item of preferred_domain_list as well. But if the qty requested
is not reached
            it tries to find the remaining quantity by
looping on the preferred_domain_list (tries with the second
item and so on).
            Make sure the quants aren't found twice =>
all the domains of preferred_domain_list should be orthogonal
        `,
		func(rs m.StockQuantSet, qty interface{}, move interface{}, ops interface{}, lot_id interface{}, domain interface{}, preferred_domain_list interface{}) {
			//        return self.quants_get_reservation(
			//            qty, move,
			//            pack_operation_id=ops and ops.id or False,
			//            lot_id=lot_id,
			//            company_id=self.env.context.get('company_id', False),
			//            domain=domain,
			//            preferred_domain_list=preferred_domain_list)
		})
	h.StockQuant().Methods().QuantsGetReservation().DeclareMethod(
		` This function tries to find quants for the given domain
and move/ops, by trying to first limit
            the choice on the quants that match the first
item of preferred_domain_list as well. But if the qty requested
is not reached
            it tries to find the remaining quantity by
looping on the preferred_domain_list (tries with the second
item and so on).
            Make sure the quants aren't found twice =>
all the domains of preferred_domain_list should be orthogonal
        `,
		func(rs m.StockQuantSet, qty interface{}, move interface{}, pack_operation_id interface{}, lot_id interface{}, company_id interface{}, domain interface{}, preferred_domain_list interface{}) {
			//        reservations = [(None, qty)]
			//        pack_operation = self.env['stock.pack.operation'].browse(
			//            pack_operation_id)
			//        location = pack_operation.location_id if pack_operation else move.location_id
			//        if location.usage in ['inventory', 'production', 'supplier']:
			//            return reservations
			//        restrict_lot_id = lot_id if pack_operation else move.restrict_lot_id.id or lot_id
			//        removal_strategy = move.get_removal_strategy()
			//        domain = self._quants_get_reservation_domain(
			//            move,
			//            pack_operation_id=pack_operation_id,
			//            lot_id=lot_id,
			//            company_id=company_id,
			//            initial_domain=domain)
			//        if not restrict_lot_id and not preferred_domain_list:
			//            meta_domains = [[]]
			//        elif restrict_lot_id and not preferred_domain_list:
			//            meta_domains = [[('lot_id', '=', restrict_lot_id)], [
			//                ('lot_id', '=', False)]]
			//        elif restrict_lot_id and preferred_domain_list:
			//            lot_list = []
			//            no_lot_list = []
			//            for inner_domain in preferred_domain_list:
			//                lot_list.append(
			//                    inner_domain + [('lot_id', '=', restrict_lot_id)])
			//                no_lot_list.append(inner_domain + [('lot_id', '=', False)])
			//            meta_domains = lot_list + no_lot_list
			//        else:
			//            meta_domains = preferred_domain_list
			//        res_qty = qty
			//        while (float_compare(res_qty, 0, precision_rounding=move.product_id.uom_id.rounding) and meta_domains):
			//            additional_domain = meta_domains.pop(0)
			//            reservations.pop()
			//            new_reservations = self._quants_get_reservation(
			//                res_qty, move,
			//                ops=pack_operation,
			//                domain=domain + additional_domain,
			//                removal_strategy=removal_strategy)
			//            for quant in new_reservations:
			//                if quant[0]:
			//                    res_qty -= quant[1]
			//            reservations += new_reservations
			//        return reservations
		})
	h.StockQuant().Methods().QuantsGetReservationDomain().DeclareMethod(
		`QuantsGetReservationDomain`,
		func(rs m.StockQuantSet, move interface{}, pack_operation_id interface{}, lot_id interface{}, company_id interface{}, initial_domain interface{}) {
			//        initial_domain = initial_domain if initial_domain is not None else [
			//            ('qty', '>', 0.0)]
			//        domain = initial_domain + [('product_id', '=', move.product_id.id)]
			//        if pack_operation_id:
			//            pack_operation = self.env['stock.pack.operation'].browse(
			//                pack_operation_id)
			//            domain += [('location_id', '=', pack_operation.location_id.id)]
			//            if pack_operation.owner_id:
			//                domain += [('owner_id', '=', pack_operation.owner_id.id)]
			//            if pack_operation.package_id and not pack_operation.product_id:
			//                domain += [('package_id', 'child_of',
			//                            pack_operation.package_id.id)]
			//            elif pack_operation.package_id and pack_operation.product_id:
			//                domain += [('package_id', '=', pack_operation.package_id.id)]
			//            else:
			//                domain += [('package_id', '=', False)]
			//        else:
			//            domain += [('location_id', 'child_of', move.location_id.id)]
			//            if move.restrict_partner_id:
			//                domain += [('owner_id', '=', move.restrict_partner_id.id)]
			//        if company_id:
			//            domain += [('company_id', '=', company_id)]
			//        else:
			//            domain += [('company_id', '=', move.company_id.id)]
			//        return domain
		})
	h.StockQuant().Methods().QuantsRemovalGetOrder().DeclareMethod(
		`QuantsRemovalGetOrder`,
		func(rs m.StockQuantSet, removal_strategy interface{}) {
			//        if removal_strategy == 'fifo':
			//            return 'in_date, id'
			//        elif removal_strategy == 'lifo':
			//            return 'in_date desc, id desc'
			//        raise UserError(_('Removal strategy %s not implemented.') %
			//                        (removal_strategy))
		})
	h.StockQuant().Methods().QuantsGetReservation().DeclareMethod(
		` Implementation of removal strategies.

        :return: a structure containing an ordered list
of tuples: quants and
                 the quantity to remove from them. A tuple (None, qty)
                 represents a qty not possible to reserve.
        `,
		func(rs m.StockQuantSet, quantity interface{}, move interface{}, ops interface{}, domain interface{}, orderby interface{}, removal_strategy interface{}) {
			//        if removal_strategy:
			//            order = self._quants_removal_get_order(removal_strategy)
			//        elif orderby:
			//            order = orderby
			//        else:
			//            order = 'in_date'
			//        rounding = move.product_id.uom_id.rounding
			//        domain = domain if domain is not None else [('qty', '>', 0.0)]
			//        res = []
			//        offset = 0
			//        remaining_quantity = quantity
			//        quants = self.search(domain, order=order, limit=10, offset=offset)
			//        while float_compare(remaining_quantity, 0, precision_rounding=rounding) > 0 and quants:
			//            for quant in quants:
			//                if float_compare(remaining_quantity, abs(quant.qty), precision_rounding=rounding) >= 0:
			//                    # reserved_quants.append(self._ReservedQuant(quant, abs(quant.qty)))
			//                    res += [(quant, abs(quant.qty))]
			//                    remaining_quantity -= abs(quant.qty)
			//                elif float_compare(remaining_quantity, 0.0, precision_rounding=rounding) != 0:
			//                    # reserved_quants.append(self._ReservedQuant(quant, remaining_quantity))
			//                    res += [(quant, remaining_quantity)]
			//                    remaining_quantity = 0
			//            offset += 10
			//            quants = self.search(domain, order=order, limit=10, offset=offset)
			//        if float_compare(remaining_quantity, 0, precision_rounding=rounding) > 0:
			//            res.append((None, remaining_quantity))
			//        return res
		})
	h.StockQuant().Methods().GetTopLevelPackages().DeclareMethod(
		` This method searches for as much possible higher level packages that
        can be moved as a single operation, given a list
of quants to move and
        their suggested destination, and returns the list
of matching packages. `,
		func(rs m.StockQuantSet, product_to_location interface{}) {
			//        top_lvl_packages = self.env['stock.quant.package']
			//        for package in self.mapped('package_id'):
			//            all_in = True
			//            top_package = self.env['stock.quant.package']
			//            while package:
			//                if any(quant not in self for quant in package.get_content()):
			//                    all_in = False
			//                if all_in:
			//                    destinations = set([product_to_location[product]
			//                                        for product in package.get_content().mapped('product_id')])
			//                    if len(destinations) > 1:
			//                        all_in = False
			//                if all_in:
			//                    top_package = package
			//                    package = package.parent_id
			//                else:
			//                    package = False
			//            top_lvl_packages |= top_package
			//        return top_lvl_packages
		})
	h.StockQuant().Methods().GetLatestMove().DeclareMethod(
		`GetLatestMove`,
		func(rs m.StockQuantSet) {
			//        latest_move = self.history_ids[0]
			//        for move in self.history_ids:
			//            if move.date > latest_move.date:
			//                latest_move = move
			//        return latest_move
		})
	h.StockQuant().Methods().QuantSplit().DeclareMethod(
		`QuantSplit`,
		func(rs m.StockQuantSet, qty interface{}) {
			//        self.ensure_one()
			//        rounding = self.product_id.uom_id.rounding
			//        if float_compare(abs(self.qty), abs(qty), precision_rounding=rounding) <= 0:
			//            return False
			//        qty_round = float_round(qty, precision_rounding=rounding)
			//        new_qty_round = float_round(
			//            self.qty - qty, precision_rounding=rounding)
			//        self._cr.execute(
			//            """SELECT move_id FROM stock_quant_move_rel WHERE quant_id = %s""", (self.id))
			//        res = self._cr.fetchall()
			//        new_quant = self.sudo().copy(
			//            default={'qty': new_qty_round, 'history_ids': [(4, x[0]) for x in res]})
			//        self.sudo().write({'qty': qty_round})
			//        return new_quant
		})
	h.StockQuantPackage().DeclareModel()

	h.StockQuantPackage().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{
			String:  "Package Reference",
			NoCopy:  true,
			Index:   true,
			Default: func(env models.Environment) interface{} { return },
		},
		"QuantIds": models.One2ManyField{
			RelationModel: h.StockQuant(),
			ReverseFK:     "",
			String:        "Bulk Content",
			ReadOnly:      true,
		},
		"ParentId": models.Many2OneField{
			RelationModel: h.StockQuantPackage(),
			String:        "Parent Package",
			OnDelete:      `restrict`,
			ReadOnly:      true,
			Help:          "The package containing this item",
		},
		"AncestorIds": models.One2ManyField{
			RelationModel: h.StockQuantPackage(),
			String:        "Ancestors",
			Compute:       h.StockQuantPackage().Methods().ComputeAncestorIds(),
		},
		"ChildrenQuantIds": models.One2ManyField{
			RelationModel: h.StockQuant(),
			String:        "All Bulk Content",
			Compute:       h.StockQuantPackage().Methods().ComputeChildrenQuantIds(),
		},
		"ChildrenIds": models.One2ManyField{
			RelationModel: h.StockQuantPackage(),
			ReverseFK:     "",
			String:        "Contained Packages",
			ReadOnly:      true,
		},
		"ParentLeft": models.IntegerField{
			String: "Left Parent",
			Index:  true,
		},
		"ParentRight": models.IntegerField{
			String: "Right Parent",
			Index:  true,
		},
		"PackagingId": models.Many2OneField{
			RelationModel: h.ProductPackaging(),
			String:        "Package Type",
			Index:         true,
			Help: "This field should be completed only if everything inside" +
				"the package share the same product, otherwise it doesn't" +
				"really makes sense.",
		},
		"LocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Location",
			Compute:       h.StockQuantPackage().Methods().ComputePackageInfo(),
			//search='_search_location'
			Index:    true,
			ReadOnly: true,
		},
		"CompanyId": models.Many2OneField{
			RelationModel: h.Company(),
			String:        "Company",
			Compute:       h.StockQuantPackage().Methods().ComputePackageInfo(),
			//search='_search_company'
			Index:    true,
			ReadOnly: true,
		},
		"OwnerId": models.Many2OneField{
			RelationModel: h.Partner(),
			String:        "Owner",
			Compute:       h.StockQuantPackage().Methods().ComputePackageInfo(),
			//search='_search_owner'
			Index:    true,
			ReadOnly: true,
		},
	})
	h.StockQuantPackage().Methods().ComputeAncestorIds().DeclareMethod(
		`ComputeAncestorIds`,
		func(rs h.StockQuantPackageSet) h.StockQuantPackageData {
			//        self.ancestor_ids = self.env['stock.quant.package'].search(
			//            [('id', 'parent_of', self.id)]).ids
		})
	h.StockQuantPackage().Methods().ComputeChildrenQuantIds().DeclareMethod(
		`ComputeChildrenQuantIds`,
		func(rs h.StockQuantPackageSet) h.StockQuantPackageData {
			//        for package in self:
			//            if package.id:
			//                package.children_quant_ids = self.env['stock.quant'].search(
			//                    [('package_id', 'child_of', package.id)]).ids
		})
	h.StockQuantPackage().Methods().ComputePackageInfo().DeclareMethod(
		`ComputePackageInfo`,
		func(rs h.StockQuantPackageSet) h.StockQuantPackageData {
			//        for package in self:
			//            quants = package.children_quant_ids
			//            if quants:
			//                values = quants[0]
			//            else:
			//                values = {'location_id': False,
			//                          'company_id': self.env.user.company_id.id, 'owner_id': False}
			//            package.location_id = values['location_id']
			//            package.company_id = values['company_id']
			//            package.owner_id = values['owner_id']
		})
	h.StockQuantPackage().Methods().NameGet().Extend(
		`NameGet`,
		func(rs m.StockQuantPackageSet) {
			//        return self._compute_complete_name().items()
		})
	h.StockQuantPackage().Methods().ComputeCompleteName().DeclareMethod(
		` Forms complete name of location from parent location to
child location. `,
		func(rs m.StockQuantPackageSet) {
			//        res = {}
			//        for package in self:
			//            current = package
			//            name = current.name
			//            while current.parent_id:
			//                name = '%s / %s' % (current.parent_id.name, name)
			//                current = current.parent_id
			//            res[package.id] = name
			//        return res
		})
	h.StockQuantPackage().Methods().SearchLocation().DeclareMethod(
		`SearchLocation`,
		func(rs m.StockQuantPackageSet, operator interface{}, value interface{}) {
			//        if value:
			//            packs = self.search([('quant_ids.location_id', operator, value)])
			//        else:
			//            packs = self.search([('quant_ids', operator, value)])
			//        if packs:
			//            return [('id', 'parent_of', packs.ids)]
			//        else:
			//            return [('id', '=', False)]
		})
	h.StockQuantPackage().Methods().SearchCompany().DeclareMethod(
		`SearchCompany`,
		func(rs m.StockQuantPackageSet, operator interface{}, value interface{}) {
			//        if value:
			//            packs = self.search([('quant_ids.company_id', operator, value)])
			//        else:
			//            packs = self.search([('quant_ids', operator, value)])
			//        if packs:
			//            return [('id', 'parent_of', packs.ids)]
			//        else:
			//            return [('id', '=', False)]
		})
	h.StockQuantPackage().Methods().SearchOwner().DeclareMethod(
		`SearchOwner`,
		func(rs m.StockQuantPackageSet, operator interface{}, value interface{}) {
			//        if value:
			//            packs = self.search([('quant_ids.owner_id', operator, value)])
			//        else:
			//            packs = self.search([('quant_ids', operator, value)])
			//        if packs:
			//            return [('id', 'parent_of', packs.ids)]
			//        else:
			//            return [('id', '=', False)]
		})
	h.StockQuantPackage().Methods().CheckLocationConstraint().DeclareMethod(
		`checks that all quants in a package are stored in the same
location. This function cannot be used
           as a constraint because it needs to be checked
on pack operations (they may not call write on the
           package)
        `,
		func(rs m.StockQuantPackageSet) {
			//        for pack in self:
			//            parent = pack
			//            while parent.parent_id:
			//                parent = parent.parent_id
			//            locations = parent.get_content().filtered(
			//                lambda quant: quant.qty > 0.0).mapped('location_id')
			//            if len(locations) != 1:
			//                raise UserError(
			//                    _('Everything inside a package should be in the same location'))
			//        return True
		})
	h.StockQuantPackage().Methods().ActionViewRelatedPicking().DeclareMethod(
		` Returns an action that display the picking related to this
        package (source or destination).
        `,
		func(rs m.StockQuantPackageSet) {
			//        self.ensure_one()
			//        pickings = self.env['stock.picking'].search(
			//            ['|', ('pack_operation_ids.package_id', '=', self.id), ('pack_operation_ids.result_package_id', '=', self.id)])
			//        action = self.env.ref('stock.action_picking_tree_all').read()[0]
			//        action['domain'] = [('id', 'in', pickings.ids)]
			//        return action
		})
	h.StockQuantPackage().Methods().Unpack().DeclareMethod(
		`Unpack`,
		func(rs m.StockQuantPackageSet) {
			//        for package in self:
			//            # TDE FIXME: why superuser ?
			//            package.mapped('quant_ids').sudo().write(
			//                {'package_id': package.parent_id.id})
			//            package.mapped('children_ids').write(
			//                {'parent_id': package.parent_id.id})
			//        return self.env['ir.actions.act_window'].for_xml_id('stock', 'action_package_view')
		})
	h.StockQuantPackage().Methods().ViewContentPackage().DeclareMethod(
		`ViewContentPackage`,
		func(rs m.StockQuantPackageSet) {
			//        action = self.env['ir.actions.act_window'].for_xml_id(
			//            'stock', 'quantsact')
			//        action['domain'] = [('id', 'in', self._get_contained_quants().ids)]
			//        return action
		})
	//    get_content_package = view_content_package
	h.StockQuantPackage().Methods().GetContainedQuants().DeclareMethod(
		`GetContainedQuants`,
		func(rs m.StockQuantPackageSet) {
			//        return self.env['stock.quant'].search([('package_id', 'child_of', self.ids)])
		})
	//    get_content = _get_contained_quants
	h.StockQuantPackage().Methods().GetAllProductsQuantities().DeclareMethod(
		`This function computes the different product quantities
for the given package
        `,
		func(rs m.StockQuantPackageSet) {
			//        res = {}
			//        for quant in self._get_contained_quants():
			//            if quant.product_id not in res:
			//                res[quant.product_id] = 0
			//            res[quant.product_id] += quant.qty
			//        return res
		})
}
