package stock

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/hexya/src/models/types"
	"github.com/hexya-erp/pool/h"
)

//import logging
//_logger = logging.getLogger(__name__)
func init() {
	h.ProcurementGroup().DeclareModel()

	h.ProcurementGroup().AddFields(map[string]models.FieldDefinition{
		"PartnerId": models.Many2OneField{
			RelationModel: h.Partner(),
			String:        "Partner",
		},
	})
	h.ProcurementRule().DeclareModel()

	h.ProcurementRule().AddFields(map[string]models.FieldDefinition{
		"LocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Procurement Location",
		},
		"LocationSrcId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Source Location",
			Help:          "Source location is action=move",
		},
		"RouteId": models.Many2OneField{
			RelationModel: h.StockLocationRoute(),
			String:        "Route",
			Help:          "If route_id is False, the rule is global",
		},
		"ProcureMethod": models.SelectionField{
			Selection: types.Selection{
				"make_to_stock": "Take From Stock",
				"make_to_order": "Create Procurement",
			},
			String:   "Move Supply Method",
			Default:  models.DefaultValue("make_to_stock"),
			Required: true,
			Help: "Determines the procurement method of the stock move that" +
				"will be generated: whether it will need to 'take from the" +
				"available stock' in its source location or needs to ignore" +
				"its stock and create a procurement over there.",
		},
		"RouteSequence": models.IntegerField{
			String:  "Route Sequence",
			Related: `RouteId.Sequence`,
			Stored:  true,
		},
		"PickingTypeId": models.Many2OneField{
			RelationModel: h.StockPickingType(),
			String:        "Picking Type",
			Required:      true,
			Help: "Picking Type determines the way the picking should be shown" +
				"in the view, reports, ...",
		},
		"Delay": models.IntegerField{
			String:  "Number of Days",
			Default: models.DefaultValue(0),
		},
		"PartnerAddressId": models.Many2OneField{
			RelationModel: h.Partner(),
			String:        "Partner Address",
		},
		"Propagate": models.BooleanField{
			String:  "Propagate cancel and split",
			Default: models.DefaultValue(true),
			Help: "If checked, when the previous move of the move (which was" +
				"generated by a next procurement) is cancelled or split," +
				"the move generated by this move will too",
		},
		"WarehouseId": models.Many2OneField{
			RelationModel: h.StockWarehouse(),
			String:        "Served Warehouse",
			Help:          "The warehouse this rule is for",
		},
		"PropagateWarehouseId": models.Many2OneField{
			RelationModel: h.StockWarehouse(),
			String:        "Warehouse to Propagate",
			Help: "The warehouse to propagate on the created move/procurement," +
				"which can be different of the warehouse this rule is for" +
				"(e.g for resupplying rules from another warehouse)",
		},
	})
	h.ProcurementRule().Methods().GetAction().DeclareMethod(
		`GetAction`,
		func(rs m.ProcurementRuleSet) {
			//        result = super(ProcurementRule, self)._get_action()
			//        return result + [('move', _('Move From Another Location'))]
		})
	h.ProcurementOrder().DeclareModel()

	h.ProcurementOrder().AddFields(map[string]models.FieldDefinition{
		"LocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Procurement Location",
		},
		"PartnerDestId": models.Many2OneField{
			RelationModel: h.Partner(),
			String:        "Customer Address",
			Help: "In case of dropshipping, we need to know the destination" +
				"address more precisely",
		},
		"MoveIds": models.One2ManyField{
			RelationModel: h.StockMove(),
			ReverseFK:     "",
			String:        "Moves",
			Help:          "Moves created by the procurement",
		},
		"MoveDestId": models.Many2OneField{
			RelationModel: h.StockMove(),
			String:        "Destination Move",
			Help:          "Move which caused (created) the procurement",
		},
		"RouteIds": models.Many2ManyField{
			RelationModel:    h.StockLocationRoute(),
			M2MLinkModelName: "",
			M2MOurField:      "",
			M2MTheirField:    "",
			String:           "Preferred Routes",
			Help: "Preferred route to be followed by the procurement order." +
				"Usually copied from the generating document (SO) but could" +
				"be set up manually.",
		},
		"WarehouseId": models.Many2OneField{
			RelationModel: h.StockWarehouse(),
			String:        "Warehouse",
			Help:          "Warehouse to consider for the route selection",
		},
		"OrderpointId": models.Many2OneField{
			RelationModel: h.StockWarehouseOrderpoint(),
			String:        "Minimum Stock Rule",
		},
	})
	h.ProcurementOrder().Methods().OnchangeWarehouseId().DeclareMethod(
		`OnchangeWarehouseId`,
		func(rs m.ProcurementOrderSet) {
			//        if self.warehouse_id:
			//            self.location_id = self.warehouse_id.lot_stock_id.id
		})
	h.ProcurementOrder().Methods().PropagateCancels().DeclareMethod(
		`PropagateCancels`,
		func(rs m.ProcurementOrderSet) {
			//        cancel_moves = self.with_context(cancel_procurement=True).filtered(
			//            lambda order: order.rule_id.action == 'move').mapped('move_ids')
			//        if cancel_moves:
			//            cancel_moves.action_cancel()
			//        return self.search([('move_dest_id', 'in', cancel_moves.filtered(lambda move: move.propagate).ids)])
		})
	h.ProcurementOrder().Methods().Cancel().DeclareMethod(
		`Cancel`,
		func(rs m.ProcurementOrderSet) {
			//        propagated_procurements = self.filtered(
			//            lambda order: order.state != 'done').propagate_cancels()
			//        if propagated_procurements:
			//            propagated_procurements.cancel()
			//        return super(ProcurementOrder, self).cancel()
		})
	h.ProcurementOrder().Methods().DoViewPickings().DeclareMethod(
		` Return an action to display the pickings belonging to the same
        procurement group of given ids. `,
		func(rs m.ProcurementOrderSet) {
			//        action = self.env.ref('stock.do_view_pickings').read()[0]
			//        action['domain'] = [('group_id', 'in', self.mapped('group_id').ids)]
			//        return action
		})
	h.ProcurementOrder().Methods().FindSuitableRule().DeclareMethod(
		`FindSuitableRule`,
		func(rs m.ProcurementOrderSet) {
			//        rule = super(ProcurementOrder, self)._find_suitable_rule()
			//        if not rule:
			//            # a rule defined on 'Stock' is suitable for a procurement in 'Stock\Bin A'
			//            all_parent_location_ids = self._find_parent_locations()
			//            rule = self._search_suitable_rule(
			//                [('location_id', 'in', all_parent_location_ids.ids)])
			//        return rule
		})
	h.ProcurementOrder().Methods().FindParentLocations().DeclareMethod(
		`FindParentLocations`,
		func(rs m.ProcurementOrderSet) {
			//        parent_locations = self.env['stock.location']
			//        location = self.location_id
			//        while location:
			//            parent_locations |= location
			//            location = location.location_id
			//        return parent_locations
		})
	h.ProcurementOrder().Methods().SearchSuitableRule().DeclareMethod(
		` First find a rule among the ones defined on the procurement order
        group; then try on the routes defined for the product;
finally fallback
        on the default behavior `,
		func(rs m.ProcurementOrderSet, domain interface{}) {
			//        if self.warehouse_id:
			//            domain = expression.AND(
			//                [['|', ('warehouse_id', '=', self.warehouse_id.id), ('warehouse_id', '=', False)], domain])
			//        Pull = self.env['procurement.rule']
			//        res = self.env['procurement.rule']
			//        if self.route_ids:
			//            res = Pull.search(expression.AND(
			//                [[('route_id', 'in', self.route_ids.ids)], domain]), order='route_sequence, sequence', limit=1)
			//        if not res:
			//            product_routes = self.product_id.route_ids | self.product_id.categ_id.total_route_ids
			//            if product_routes:
			//                res = Pull.search(expression.AND(
			//                    [[('route_id', 'in', product_routes.ids)], domain]), order='route_sequence, sequence', limit=1)
			//        if not res:
			//            warehouse_routes = self.warehouse_id.route_ids
			//            if warehouse_routes:
			//                res = Pull.search(expression.AND(
			//                    [[('route_id', 'in', warehouse_routes.ids)], domain]), order='route_sequence, sequence', limit=1)
			//        if not res:
			//            res = Pull.search(expression.AND(
			//                [[('route_id', '=', False)], domain]), order='sequence', limit=1)
			//        return res
		})
	h.ProcurementOrder().Methods().GetStockMoveValues().DeclareMethod(
		` Returns a dictionary of values that will be used to create
a stock move from a procurement.
        This function assumes that the given procurement
has a rule (action == 'move') set on it.

        :param procurement: browse record
        :rtype: dictionary
        `,
		func(rs m.ProcurementOrderSet) {
			//        group_id = False
			//        if self.rule_id.group_propagation_option == 'propagate':
			//            group_id = self.group_id.id
			//        elif self.rule_id.group_propagation_option == 'fixed':
			//            group_id = self.rule_id.group_id.id
			//        date_expected = (datetime.strptime(self.date_planned, DEFAULT_SERVER_DATETIME_FORMAT) -
			//                         relativedelta(days=self.rule_id.delay or 0)).strftime(DEFAULT_SERVER_DATETIME_FORMAT)
			//        qty_done = sum(self.move_ids.filtered(
			//            lambda move: move.state == 'done').mapped('product_uom_qty'))
			//        qty_left = max(self.product_qty - qty_done, 0)
			//        return {
			//            'name': self.name[:2000],
			//            'company_id': self.rule_id.company_id.id or self.rule_id.location_src_id.company_id.id or self.rule_id.location_id.company_id.id or self.company_id.id,
			//            'product_id': self.product_id.id,
			//            'product_uom': self.product_uom.id,
			//            'product_uom_qty': qty_left,
			//            'partner_id': self.rule_id.partner_address_id.id or (self.group_id and self.group_id.partner_id.id) or False,
			//            'location_id': self.rule_id.location_src_id.id,
			//            'location_dest_id': self.location_id.id,
			//            'move_dest_id': self.move_dest_id and self.move_dest_id.id or False,
			//            'procurement_id': self.id,
			//            'rule_id': self.rule_id.id,
			//            'procure_method': self.rule_id.procure_method,
			//            'origin': self.origin,
			//            'picking_type_id': self.rule_id.picking_type_id.id,
			//            'group_id': group_id,
			//            'route_ids': [(4, route.id) for route in self.route_ids],
			//            'warehouse_id': self.rule_id.propagate_warehouse_id.id or self.rule_id.warehouse_id.id,
			//            'date': date_expected,
			//            'date_expected': date_expected,
			//            'propagate': self.rule_id.propagate,
			//            'priority': self.priority,
			//        }
		})
	h.ProcurementOrder().Methods().RunMoveCreate().DeclareMethod(
		`RunMoveCreate`,
		func(rs m.ProcurementOrderSet) {
			//        _logger.warning(
			//            "'_run_move_create' has been renamed into '_get_stock_move_values'... Overrides are ignored")
			//        return self._get_stock_move_values()
		})
	h.ProcurementOrder().Methods().Run().DeclareMethod(
		`Run`,
		func(rs m.ProcurementOrderSet) {
			//        if self.rule_id.action == 'move':
			//            if not self.rule_id.location_src_id:
			//                self.message_post(body=_('No source location defined!'))
			//                return False
			//            # create the move as SUPERUSER because the current user may not have the rights to do it (mto product launched by a sale for example)
			//            self.env['stock.move'].sudo().create(self._get_stock_move_values())
			//            return True
			//        return super(ProcurementOrder, self)._run()
		})
	h.ProcurementOrder().Methods().Run().DeclareMethod(
		`Run`,
		func(rs m.ProcurementOrderSet, autocommit interface{}) {
			//        new_self = self.filtered(lambda order: order.state not in [
			//                                 'running', 'done', 'cancel'])
			//        res = True
			//        if new_self:
			//            res = super(ProcurementOrder, new_self).run(autocommit=autocommit)
			//
			//            # after all the procurements are run, check if some created a draft stock move that needs to be confirmed
			//            # (we do that in batch because it fasts the picking assignation and the picking state computation)
			//            move_ids = new_self.filtered(lambda order: order.state == 'running' and order.rule_id.action == 'move').mapped(
			//                'move_ids').filtered(lambda move: move.state == 'draft')
			//            if move_ids:
			//                move_ids.action_confirm()
			//
			//            # TDE FIXME: action_confirm in stock_move already call run() ... necessary ??
			//            # If procurements created other procurements, run the created in batch
			//            new_procurements = self.search(
			//                [('move_dest_id.procurement_id', 'in', new_self.ids)], order='id')
			//            if new_procurements:
			//                res = new_procurements.run(autocommit=autocommit)
			//        return res
		})
	h.ProcurementOrder().Methods().Check().DeclareMethod(
		` Checking rules of type 'move': satisfied only if all related moves
        are done/cancel and if the requested quantity is moved. `,
		func(rs m.ProcurementOrderSet) {
			//        if self.rule_id.action == 'move':
			//            # In case Phantom BoM splits only into procurements
			//            if not self.move_ids:
			//                return True
			//            move_all_done_or_cancel = all(
			//                move.state in ['done', 'cancel'] for move in self.move_ids)
			//            move_all_cancel = all(
			//                move.state == 'cancel' for move in self.move_ids)
			//            if not move_all_done_or_cancel:
			//                return False
			//            elif move_all_done_or_cancel and not move_all_cancel:
			//                return True
			//            else:
			//                self.message_post(
			//                    body=_('All stock moves have been cancelled for this procurement.'))
			//                # TDE FIXME: strange that a check method actually modified the procurement...
			//                self.write({'state': 'cancel'})
			//                return False
			//        return super(ProcurementOrder, self)._check()
		})
	h.ProcurementOrder().Methods().RunScheduler().DeclareMethod(
		` Call the scheduler in order to check the running procurements
(super method), to check the minimum stock rules
        and the availability of moves. This function is
intended to be run for all the companies at the same time, so
        we run functions as SUPERUSER to avoid intercompanies
and access rights issues. `,
		func(rs m.ProcurementOrderSet, use_new_cursor interface{}, company_id interface{}) {
			//        super(ProcurementOrder, self).run_scheduler(
			//            use_new_cursor=use_new_cursor, company_id=company_id)
			//        try:
			//            if use_new_cursor:
			//                cr = registry(self._cr.dbname).cursor()
			//                self = self.with_env(self.env(cr=cr))  # TDE FIXME
			//
			//            # Minimum stock rules
			//            self.sudo()._procure_orderpoint_confirm(
			//                use_new_cursor=use_new_cursor, company_id=company_id)
			//
			//            # Search all confirmed stock_moves and try to assign them
			//            confirmed_moves = self.env['stock.move'].search([('state', '=', 'confirmed'), (
			//                'product_uom_qty', '!=', 0.0)], limit=None, order='priority desc, date_expected asc')
			//            for x in xrange(0, len(confirmed_moves.ids), 100):
			//                # TDE CLEANME: muf muf
			//                self.env['stock.move'].browse(
			//                    confirmed_moves.ids[x:x + 100]).action_assign()
			//                if use_new_cursor:
			//                    self._cr.commit()
			//            if use_new_cursor:
			//                self._cr.commit()
			//        finally:
			//            if use_new_cursor:
			//                try:
			//                    self._cr.close()
			//                except Exception:
			//                    pass
			//        return {}
		})
	h.ProcurementOrder().Methods().ProcurementFromOrderpointGetOrder().DeclareMethod(
		`ProcurementFromOrderpointGetOrder`,
		func(rs m.ProcurementOrderSet) {
			//        return 'location_id'
		})
	h.ProcurementOrder().Methods().ProcurementFromOrderpointGetGroupingKey().DeclareMethod(
		`ProcurementFromOrderpointGetGroupingKey`,
		func(rs m.ProcurementOrderSet, orderpoint_ids interface{}) {
			//        orderpoints = self.env['stock.warehouse.orderpoint'].browse(
			//            orderpoint_ids)
			//        return orderpoints.location_id.id
		})
	h.ProcurementOrder().Methods().ProcurementFromOrderpointGetGroups().DeclareMethod(
		` Make groups for a given orderpoint; by default schedule
all operations in one without date `,
		func(rs m.ProcurementOrderSet, orderpoint_ids interface{}) {
			//        return [{'to_date': False, 'procurement_values': dict()}]
		})
	h.ProcurementOrder().Methods().ProcurementFromOrderpointPostProcess().DeclareMethod(
		`ProcurementFromOrderpointPostProcess`,
		func(rs m.ProcurementOrderSet, orderpoint_ids interface{}) {
			//        return True
		})
	h.ProcurementOrder().Methods().GetOrderpointDomain().DeclareMethod(
		`GetOrderpointDomain`,
		func(rs m.ProcurementOrderSet, company_id interface{}) {
			//        domain = [('company_id', '=', company_id)] if company_id else []
			//        domain += [('product_id.active', '=', True)]
			//        return domain
		})
	h.ProcurementOrder().Methods().ProcureOrderpointConfirm().DeclareMethod(
		` Create procurements based on orderpoints.
        :param bool use_new_cursor: if set, use a dedicated
cursor and auto-commit after processing
            1000 orderpoints.
            This is appropriate for batch jobs only.
        `,
		func(rs m.ProcurementOrderSet, use_new_cursor interface{}, company_id interface{}) {
			//        if company_id and self.env.user.company_id.id != company_id:
			//            # To ensure that the company_id is taken into account for
			//            # all the processes triggered by this method
			//            # i.e. If a PO is generated by the run of the procurements the
			//            # sequence to use is the one for the specified company not the
			//            # one of the user's company
			//            self = self.with_context(
			//                company_id=company_id, force_company=company_id)
			//        OrderPoint = self.env['stock.warehouse.orderpoint']
			//        domain = self._get_orderpoint_domain(company_id=company_id)
			//        orderpoints_noprefetch = OrderPoint.with_context(prefetch_fields=False).search(domain,
			//                                                                                       order=self._procurement_from_orderpoint_get_order()).ids
			//        while orderpoints_noprefetch:
			//            if use_new_cursor:
			//                cr = registry(self._cr.dbname).cursor()
			//                self = self.with_env(self.env(cr=cr))
			//            OrderPoint = self.env['stock.warehouse.orderpoint']
			//            Procurement = self.env['procurement.order']
			//            ProcurementAutorundefer = Procurement.with_context(
			//                procurement_autorun_defer=True)
			//            procurement_list = []
			//
			//            orderpoints = OrderPoint.browse(orderpoints_noprefetch[:1000])
			//            orderpoints_noprefetch = orderpoints_noprefetch[1000:]
			//
			//            # Calculate groups that can be executed together
			//            location_data = defaultdict(lambda: dict(
			//                products=self.env['product.product'], orderpoints=self.env['stock.warehouse.orderpoint'], groups=list()))
			//            for orderpoint in orderpoints:
			//                key = self._procurement_from_orderpoint_get_grouping_key([
			//                                                                         orderpoint.id])
			//                location_data[key]['products'] += orderpoint.product_id
			//                location_data[key]['orderpoints'] += orderpoint
			//                location_data[key]['groups'] = self._procurement_from_orderpoint_get_groups([
			//                                                                                            orderpoint.id])
			//
			//            for location_id, location_data in location_data.iteritems():
			//                location_orderpoints = location_data['orderpoints']
			//                product_context = dict(
			//                    self._context, location=location_orderpoints[0].location_id.id)
			//                substract_quantity = location_orderpoints.subtract_procurements_from_orderpoints()
			//
			//                for group in location_data['groups']:
			//                    if group.get('from_date'):
			//                        product_context['from_date'] = group['from_date'].strftime(
			//                            DEFAULT_SERVER_DATETIME_FORMAT)
			//                    if group['to_date']:
			//                        product_context['to_date'] = group['to_date'].strftime(
			//                            DEFAULT_SERVER_DATETIME_FORMAT)
			//                    product_quantity = location_data['products'].with_context(
			//                        product_context)._product_available()
			//                    for orderpoint in location_orderpoints:
			//                        try:
			//                            op_product_virtual = product_quantity[orderpoint.product_id.id]['virtual_available']
			//                            if op_product_virtual is None:
			//                                continue
			//                            if float_compare(op_product_virtual, orderpoint.product_min_qty, precision_rounding=orderpoint.product_uom.rounding) <= 0:
			//                                qty = max(
			//                                    orderpoint.product_min_qty, orderpoint.product_max_qty) - op_product_virtual
			//                                remainder = orderpoint.qty_multiple > 0 and qty % orderpoint.qty_multiple or 0.0
			//
			//                                if float_compare(remainder, 0.0, precision_rounding=orderpoint.product_uom.rounding) > 0:
			//                                    qty += orderpoint.qty_multiple - remainder
			//
			//                                if float_compare(qty, 0.0, precision_rounding=orderpoint.product_uom.rounding) < 0:
			//                                    continue
			//
			//                                qty -= substract_quantity[orderpoint.id]
			//                                qty_rounded = float_round(
			//                                    qty, precision_rounding=orderpoint.product_uom.rounding)
			//                                if qty_rounded > 0:
			//                                    new_procurement = ProcurementAutorundefer.create(
			//                                        orderpoint._prepare_procurement_values(qty_rounded, **group['procurement_values']))
			//                                    procurement_list.append(new_procurement)
			//                                    new_procurement.message_post_with_view('mail.message_origin_link',
			//                                                                           values={
			//                                                                               'self': new_procurement, 'origin': orderpoint},
			//                                                                           subtype_id=self.env.ref('mail.mt_note').id)
			//                                    self._procurement_from_orderpoint_post_process(
			//                                        [orderpoint.id])
			//                                if use_new_cursor:
			//                                    cr.commit()
			//
			//                        except OperationalError:
			//                            if use_new_cursor:
			//                                orderpoints_noprefetch += [orderpoint.id]
			//                                cr.rollback()
			//                                continue
			//                            else:
			//                                raise
			//
			//            try:
			//                # TDE CLEANME: use record set ?
			//                procurement_list.reverse()
			//                procurements = self.env['procurement.order']
			//                for p in procurement_list:
			//                    procurements += p
			//                procurements.run()
			//                if use_new_cursor:
			//                    cr.commit()
			//            except OperationalError:
			//                if use_new_cursor:
			//                    cr.rollback()
			//                    continue
			//                else:
			//                    raise
			//
			//            if use_new_cursor:
			//                cr.commit()
			//                cr.close()
			//        return {}
		})
}
