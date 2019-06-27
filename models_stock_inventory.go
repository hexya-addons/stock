package stock

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/hexya/src/models/types/dates"
	"github.com/hexya-erp/pool/h"
	"github.com/hexya-erp/pool/q"
)

func init() {
	h.StockInventory().DeclareModel()

	h.StockInventory().Methods().DefaultLocationId().DeclareMethod(
		`DefaultLocationId`,
		func(rs m.StockInventorySet) {
			//        company_user = self.env.user.company_id
			//        warehouse = self.env['stock.warehouse'].search(
			//            [('company_id', '=', company_user.id)], limit=1)
			//        if warehouse:
			//            return warehouse.lot_stock_id.id
			//        else:
			//            raise UserError(
			//                _('You must define a warehouse for the company: %s.') % (company_user.name))
		})
	h.StockInventory().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{
			String:   "Inventory Reference",
			ReadOnly: true,
			Required: true,
			//states={'draft': [('readonly', False)]}
		},
		"Date": models.DateTimeField{
			String:   "Inventory Date",
			ReadOnly: true,
			Required: true,
			Default:  func(env models.Environment) interface{} { return dates.Now() },
			Help: "The date that will be used for the stock level check of" +
				"the products and the validation of the stock move related" +
				"to this inventory.",
		},
		"LineIds": models.One2ManyField{
			RelationModel: h.StockInventoryLine(),
			ReverseFK:     "",
			String:        "Inventories",
			NoCopy:        false,
			ReadOnly:      false,
			//states={'done': [('readonly', True)]}
		},
		"MoveIds": models.One2ManyField{
			RelationModel: h.StockMove(),
			ReverseFK:     "",
			String:        "Created Moves",
			//states={'done': [('readonly', True)]}
		},
		"State": models.SelectionField{
			String: "Status",
			//selection=[('draft', 'Draft'),('cancel', 'Cancelled'),('confirm', 'In Progress'),('done', 'Validated')]
			NoCopy:   true,
			Index:    true,
			ReadOnly: true,
			Default:  models.DefaultValue("draft"),
		},
		"CompanyId": models.Many2OneField{
			RelationModel: h.Company(),
			String:        "Company",
			ReadOnly:      true,
			Index:         true,
			Required:      true,
			//states={'draft': [('readonly', False)]}
			Default: func(env models.Environment) interface{} { return env["res.company"]._company_default_get() },
		},
		"LocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Inventoried Location",
			ReadOnly:      true,
			Required:      true,
			//states={'draft': [('readonly', False)]}
			Default: models.DefaultValue(_default_location_id),
		},
		"ProductId": models.Many2OneField{
			RelationModel: h.ProductProduct(),
			String:        "Inventoried Product",
			ReadOnly:      true,
			//states={'draft': [('readonly', False)]}
			Help: "Specify Product to focus your inventory on a particular Product.",
		},
		"PackageId": models.Many2OneField{
			RelationModel: h.StockQuantPackage(),
			String:        "Inventoried Pack",
			ReadOnly:      true,
			//states={'draft': [('readonly', False)]}
			Help: "Specify Pack to focus your inventory on a particular Pack.",
		},
		"PartnerId": models.Many2OneField{
			RelationModel: h.Partner(),
			String:        "Inventoried Owner",
			ReadOnly:      true,
			//states={'draft': [('readonly', False)]}
			Help: "Specify Owner to focus your inventory on a particular Owner.",
		},
		"LotId": models.Many2OneField{
			RelationModel: h.StockProductionLot(),
			String:        "Inventoried Lot/Serial Number",
			NoCopy:        true,
			ReadOnly:      true,
			//states={'draft': [('readonly', False)]}
			Help: "Specify Lot/Serial Number to focus your inventory on a" +
				"particular Lot/Serial Number.",
		},
		"Filter": models.SelectionField{
			String: "Inventory of",
			//selection='_selection_filter'
			Required: true,
			Default:  models.DefaultValue("none"),
			Help: "If you do an entire inventory, you can choose 'All Products'" +
				"and it will prefill the inventory with the current stock." +
				" If you only do some products  (e.g. Cycle Counting) you" +
				"can choose 'Manual Selection of Products' and the system" +
				"won't propose anything.  You can also let the system propose" +
				"for a single product / lot /... ",
		},
		"TotalQty": models.FloatField{
			String:  "Total Quantity",
			Compute: h.StockInventory().Methods().ComputeTotalQty(),
		},
		"CategoryId": models.Many2OneField{
			RelationModel: h.ProductCategory(),
			String:        "Inventoried Category",
			ReadOnly:      true,
			//states={'draft': [('readonly', False)]}
			Help: "Specify Product Category to focus your inventory on a particular" +
				"Category.",
		},
		"Exhausted": models.BooleanField{
			String:   "Include Exhausted Products",
			ReadOnly: true,
			//states={'draft': [('readonly', False)]}
		},
	})
	h.StockInventory().Methods().ComputeTotalQty().DeclareMethod(
		` For single product inventory, total quantity of the counted `,
		func(rs h.StockInventorySet) h.StockInventoryData {
			//        if self.product_id:
			//            self.total_qty = sum(self.mapped('line_ids').mapped('product_qty'))
			//        else:
			//            self.total_qty = 0
		})
	h.StockInventory().Methods().Unlink().Extend(
		`Unlink`,
		func(rs m.StockInventorySet) {
			//        for inventory in self:
			//            if inventory.state == 'done':
			//                raise UserError(
			//                    _('You cannot delete a validated inventory adjustement.'))
			//        return super(Inventory, self).unlink()
		})
	h.StockInventory().Methods().SelectionFilter().DeclareMethod(
		` Get the list of filter allowed according to the options checked
        in 'Settings\Warehouse'. `,
		func(rs m.StockInventorySet) {
			//        res_filter = [
			//            ('none', _('All products')),
			//            ('category', _('One product category')),
			//            ('product', _('One product only')),
			//            ('partial', _('Select products manually'))]
			//        if self.user_has_groups('stock.group_tracking_owner'):
			//            res_filter += [('owner', _('One owner only')),
			//                           ('product_owner', _('One product for a specific owner'))]
			//        if self.user_has_groups('stock.group_production_lot'):
			//            res_filter.append(('lot', _('One Lot/Serial Number')))
			//        if self.user_has_groups('stock.group_tracking_lot'):
			//            res_filter.append(('pack', _('A Pack')))
			//        return res_filter
		})
	h.StockInventory().Methods().OnchangeFilter().DeclareMethod(
		`OnchangeFilter`,
		func(rs m.StockInventorySet) {
			//        if self.filter not in ('product', 'product_owner'):
			//            self.product_id = False
			//        if self.filter != 'lot':
			//            self.lot_id = False
			//        if self.filter not in ('owner', 'product_owner'):
			//            self.partner_id = False
			//        if self.filter != 'pack':
			//            self.package_id = False
			//        if self.filter != 'category':
			//            self.category_id = False
			//        if self.filter == 'product':
			//            self.exhausted = True
		})
	h.StockInventory().Methods().OnchangeLocationId().DeclareMethod(
		`OnchangeLocationId`,
		func(rs m.StockInventorySet) {
			//        if self.location_id.company_id:
			//            self.company_id = self.location_id.company_id
		})
	h.StockInventory().Methods().CheckFilterProduct().DeclareMethod(
		`CheckFilterProduct`,
		func(rs m.StockInventorySet) {
			//        if self.filter == 'none' and self.product_id and self.location_id and self.lot_id:
			//            return
			//        if self.filter not in ('product', 'product_owner') and self.product_id:
			//            raise UserError(
			//                _('The selected inventory options are not coherent.'))
			//        if self.filter != 'lot' and self.lot_id:
			//            raise UserError(
			//                _('The selected inventory options are not coherent.'))
			//        if self.filter not in ('owner', 'product_owner') and self.partner_id:
			//            raise UserError(
			//                _('The selected inventory options are not coherent.'))
			//        if self.filter != 'pack' and self.package_id:
			//            raise UserError(
			//                _('The selected inventory options are not coherent.'))
		})
	h.StockInventory().Methods().ActionResetProductQty().DeclareMethod(
		`ActionResetProductQty`,
		func(rs m.StockInventorySet) {
			//        self.mapped('line_ids').write({'product_qty': 0})
			//        return True
		})
	//    reset_real_qty = action_reset_product_qty
	h.StockInventory().Methods().ActionDone().DeclareMethod(
		`ActionDone`,
		func(rs m.StockInventorySet) {
			//        negative = next((line for line in self.mapped(
			//            'line_ids') if line.product_qty < 0 and line.product_qty != line.theoretical_qty), False)
			//        if negative:
			//            raise UserError(_('You cannot set a negative product quantity in an inventory line:\n\t%s - qty: %s') %
			//                            (negative.product_id.name, negative.product_qty))
			//        self.action_check()
			//        self.write({'state': 'done'})
			//        self.post_inventory()
			//        return True
		})
	h.StockInventory().Methods().PostInventory().DeclareMethod(
		`PostInventory`,
		func(rs m.StockInventorySet) {
			//        self.mapped('move_ids').filtered(
			//            lambda move: move.state != 'done').action_done()
		})
	h.StockInventory().Methods().ActionCheck().DeclareMethod(
		` Checks the inventory and computes the stock move to do `,
		func(rs m.StockInventorySet) {
			//        for inventory in self:
			//            # first remove the existing stock moves linked to this inventory
			//            inventory.mapped('move_ids').unlink()
			//            for line in inventory.line_ids:
			//                # compare the checked quantities on inventory lines to the theorical one
			//                stock_move = line._generate_moves()
		})
	h.StockInventory().Methods().ActionCancelDraft().DeclareMethod(
		`ActionCancelDraft`,
		func(rs m.StockInventorySet) {
			//        self.mapped('move_ids').action_cancel()
			//        self.write({
			//            'line_ids': [(5)],
			//            'state': 'draft'
			//        })
		})
	h.StockInventory().Methods().ActionStart().DeclareMethod(
		`ActionStart`,
		func(rs m.StockInventorySet) {
			//        for inventory in self:
			//            vals = {'state': 'confirm', 'date': fields.Datetime.now()}
			//            if (inventory.filter != 'partial') and not inventory.line_ids:
			//                vals.update({'line_ids': [
			//                            (0, 0, line_values) for line_values in inventory._get_inventory_lines_values()]})
			//            inventory.write(vals)
			//        return True
		})
	//    prepare_inventory = action_start
	h.StockInventory().Methods().ActionInventoryLineTree().DeclareMethod(
		`ActionInventoryLineTree`,
		func(rs m.StockInventorySet) {
			//        action = self.env.ref('stock.action_inventory_line_tree').read()[0]
			//        action['context'] = {
			//            'default_location_id': self.location_id.id,
			//            'default_product_id': self.product_id.id,
			//            'default_prod_lot_id': self.lot_id.id,
			//            'default_package_id': self.package_id.id,
			//            'default_partner_id': self.partner_id.id,
			//            'default_inventory_id': self.id,
			//        }
			//        return action
		})
	h.StockInventory().Methods().GetInventoryLinesValues().DeclareMethod(
		`GetInventoryLinesValues`,
		func(rs m.StockInventorySet) {
			//        locations = self.env['stock.location'].search(
			//            [('id', 'child_of', [self.location_id.id])])
			//        domain = ' location_id in %s AND active = TRUE'
			//        args = (tuple(locations.ids))
			//        vals = []
			//        Product = self.env['product.product']
			//        quant_products = self.env['product.product']
			//        products_to_filter = self.env['product.product']
			//        if self.company_id:
			//            domain += ' AND company_id = %s'
			//            args += (self.company_id.id)
			//        if self.partner_id:
			//            domain += ' AND owner_id = %s'
			//            args += (self.partner_id.id)
			//        if self.lot_id:
			//            domain += ' AND lot_id = %s'
			//            args += (self.lot_id.id)
			//        if self.product_id:
			//            domain += ' AND product_id = %s'
			//            args += (self.product_id.id)
			//            products_to_filter |= self.product_id
			//        if self.package_id:
			//            domain += ' AND package_id = %s'
			//            args += (self.package_id.id)
			//        if self.category_id:
			//            categ_products = Product.search(
			//                [('categ_id', '=', self.category_id.id)])
			//            domain += ' AND product_id = ANY (%s)'
			//            args += (categ_products.ids)
			//            products_to_filter |= categ_products
			//        self.env.cr.execute("""SELECT product_id, sum(qty) as product_qty, location_id, lot_id as prod_lot_id, package_id, owner_id as partner_id
			//            FROM stock_quant
			//            LEFT JOIN product_product
			//            ON product_product.id = stock_quant.product_id
			//            WHERE %s
			//            GROUP BY product_id, location_id, lot_id, package_id, partner_id """ % domain, args)
			//        for product_data in self.env.cr.dictfetchall():
			//            # replace the None the dictionary by False, because falsy values are tested later on
			//            for void_field in [item[0] for item in product_data.items() if item[1] is None]:
			//                product_data[void_field] = False
			//            product_data['theoretical_qty'] = product_data['product_qty']
			//            if product_data['product_id']:
			//                product_data['product_uom_id'] = Product.browse(
			//                    product_data['product_id']).uom_id.id
			//                quant_products |= Product.browse(product_data['product_id'])
			//            vals.append(product_data)
			//        if self.exhausted:
			//            exhausted_vals = self._get_exhausted_inventory_line(
			//                products_to_filter, quant_products)
			//            vals.extend(exhausted_vals)
			//        return vals
		})
	h.StockInventory().Methods().GetExhaustedInventoryLine().DeclareMethod(
		`
        This function return inventory lines for exausted products
        :param products: products With Selected Filter.
        :param quant_products: products available in stock_quants
        `,
		func(rs m.StockInventorySet, products interface{}, quant_products interface{}) {
			//        vals = []
			//        exhausted_domain = [
			//            ('type', 'not in', ('service', 'consu', 'digital'))]
			//        if products:
			//            exhausted_products = products - quant_products
			//            exhausted_domain += [('id', 'in', exhausted_products.ids)]
			//        else:
			//            exhausted_domain += [('id', 'not in', quant_products.ids)]
			//        exhausted_products = self.env['product.product'].search(
			//            exhausted_domain)
			//        for product in exhausted_products:
			//            vals.append({
			//                'inventory_id': self.id,
			//                'product_id': product.id,
			//                'location_id': self.location_id.id,
			//            })
			//        return vals
		})
	h.StockInventoryLine().DeclareModel()

	h.StockInventoryLine().AddFields(map[string]models.FieldDefinition{
		"InventoryId": models.Many2OneField{
			RelationModel: h.StockInventory(),
			String:        "Inventory",
			Index:         true,
			OnDelete:      `cascade`,
		},
		"PartnerId": models.Many2OneField{
			RelationModel: h.Partner(),
			String:        "Owner",
		},
		"ProductId": models.Many2OneField{
			RelationModel: h.ProductProduct(),
			String:        "Product",
			Index:         true,
			Required:      true,
		},
		"ProductName": models.CharField{
			String:   "Product Name",
			Related:  `ProductId.Name`,
			Stored:   true,
			ReadOnly: true,
		},
		"ProductCode": models.CharField{
			String:  "Product Code",
			Related: `ProductId.DefaultCode`,
			Stored:  true,
		},
		"ProductUomId": models.Many2OneField{
			RelationModel: h.ProductUom(),
			String:        "Product Unit of Measure",
			Required:      true,
			Default:       func(env models.Environment) interface{} { return env.ref() },
		},
		"ProductQty": models.FloatField{
			String: "Checked Quantity",
			//digits=dp.get_precision('Product Unit of Measure')
			Default: models.DefaultValue(0),
		},
		"LocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Location",
			Index:         true,
			Required:      true,
		},
		"LocationName": models.CharField{
			String:  "Location Name",
			Related: `LocationId.CompleteName`,
			Stored:  true,
		},
		"PackageId": models.Many2OneField{
			RelationModel: h.StockQuantPackage(),
			String:        "Pack",
			Index:         true,
		},
		"ProdLotId": models.Many2OneField{
			RelationModel: h.StockProductionLot(),
			String:        "Lot/Serial Number",
			Filter:        q.ProductId().Equals(product_id),
		},
		"ProdlotName": models.CharField{
			String:   "Serial Number Name",
			Related:  `ProdLotId.Name`,
			Stored:   true,
			ReadOnly: true,
		},
		"CompanyId": models.Many2OneField{
			RelationModel: h.Company(),
			String:        "Company",
			Related:       `InventoryId.CompanyId`,
			Index:         true,
			ReadOnly:      true,
			Stored:        true,
		},
		"State": models.SelectionField{
			Selection: "Status",
			Related:   `InventoryId.State`,
			ReadOnly:  true,
		},
		"TheoreticalQty": models.FloatField{
			String:  "Theoretical Quantity",
			Compute: h.StockInventoryLine().Methods().ComputeTheoreticalQty(),
			//digits=dp.get_precision('Product Unit of Measure')
			ReadOnly: true,
			Stored:   true,
		},
		"InventoryLocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Location",
			Related:       `InventoryId.LocationId`,
			//related_sudo=False
		},
	})
	h.StockInventoryLine().Methods().ComputeTheoreticalQty().DeclareMethod(
		`ComputeTheoreticalQty`,
		func(rs h.StockInventoryLineSet) h.StockInventoryLineData {
			//        if not self.product_id:
			//            self.theoretical_qty = 0
			//            return
			//        theoretical_qty = sum([x.qty for x in self._get_quants()])
			//        if theoretical_qty and self.product_uom_id and self.product_id.uom_id != self.product_uom_id:
			//            theoretical_qty = self.product_id.uom_id._compute_quantity(
			//                theoretical_qty, self.product_uom_id)
			//        self.theoretical_qty = theoretical_qty
		})
	h.StockInventoryLine().Methods().OnchangeProduct().DeclareMethod(
		`OnchangeProduct`,
		func(rs m.StockInventoryLineSet) {
			//        res = {}
			//        if self.product_id:
			//            self.product_uom_id = self.product_id.uom_id
			//            res['domain'] = {'product_uom_id': [
			//                ('category_id', '=', self.product_id.uom_id.category_id.id)]}
			//        return res
		})
	h.StockInventoryLine().Methods().OnchangeQuantityContext().DeclareMethod(
		`OnchangeQuantityContext`,
		func(rs m.StockInventoryLineSet) {
			//        if self.product_id and self.location_id and self.product_id.uom_id.category_id == self.product_uom_id.category_id:
			//            self._compute_theoretical_qty()
			//            self.product_qty = self.theoretical_qty
		})
	h.StockInventoryLine().Methods().Write().Extend(
		`Write`,
		func(rs m.StockInventoryLineSet, values models.RecordData) {
			//        values.pop('product_name', False)
			//        res = super(InventoryLine, self).write(values)
			//        return res
		})
	h.StockInventoryLine().Methods().Create().Extend(
		`Create`,
		func(rs m.StockInventoryLineSet, values models.RecordData) {
			//        values.pop('product_name', False)
			//        if 'product_id' in values and 'product_uom_id' not in values:
			//            values['product_uom_id'] = self.env['product.product'].browse(
			//                values['product_id']).uom_id.id
			//        existings = self.search([
			//            ('product_id', '=', values.get('product_id')),
			//            ('inventory_id.state', '=', 'confirm'),
			//            ('location_id', '=', values.get('location_id')),
			//            ('partner_id', '=', values.get('partner_id')),
			//            ('package_id', '=', values.get('package_id')),
			//            ('prod_lot_id', '=', values.get('prod_lot_id'))])
			//        res = super(InventoryLine, self).create(values)
			//        if existings:
			//            raise UserError(_("You cannot have two inventory adjustements in state 'in Progess' with the same product"
			//                              "(%s), same location(%s), same package, same owner and same lot. Please first validate"
			//                              "the first inventory adjustement with this product before creating another one.") %
			//                            (res.product_id.display_name, res.location_id.display_name))
			//        return res
		})
	h.StockInventoryLine().Methods().GetQuants().DeclareMethod(
		`GetQuants`,
		func(rs m.StockInventoryLineSet) {
			//        return self.env['stock.quant'].search([
			//            ('company_id', '=', self.company_id.id),
			//            ('location_id', '=', self.location_id.id),
			//            ('lot_id', '=', self.prod_lot_id.id),
			//            ('product_id', '=', self.product_id.id),
			//            ('owner_id', '=', self.partner_id.id),
			//            ('package_id', '=', self.package_id.id)])
		})
	h.StockInventoryLine().Methods().GetMoveValues().DeclareMethod(
		`GetMoveValues`,
		func(rs m.StockInventoryLineSet, qty interface{}, location_id interface{}, location_dest_id interface{}) {
			//        self.ensure_one()
			//        return {
			//            'name': _('INV:') + (self.inventory_id.name or ''),
			//            'product_id': self.product_id.id,
			//            'product_uom': self.product_uom_id.id,
			//            'product_uom_qty': qty,
			//            'date': self.inventory_id.date,
			//            'company_id': self.inventory_id.company_id.id,
			//            'inventory_id': self.inventory_id.id,
			//            'state': 'confirmed',
			//            'restrict_lot_id': self.prod_lot_id.id,
			//            'restrict_partner_id': self.partner_id.id,
			//            'location_id': location_id,
			//            'location_dest_id': location_dest_id,
			//        }
		})
	h.StockInventoryLine().Methods().FixupNegativeQuants().DeclareMethod(
		` This will handle the irreconciable quants created by a
force availability followed by a
        return. When generating the moves of an inventory
line, we look for quants of this line's
        product created to compensate a force availability.
If there are some and if the quant
        which it is propagated from is still in the same
location, we move it to the inventory
        adjustment location before getting it back. Getting
the quantity from the inventory
        location will allow the negative quant to be compensated.
        `,
		func(rs m.StockInventoryLineSet) {
			//        self.ensure_one()
			//        for quant in self._get_quants().filtered(lambda q: q.propagated_from_id.location_id.id == self.location_id.id):
			//            # send the quantity to the inventory adjustment location
			//            move_out_vals = self._get_move_values(
			//                quant.qty, self.location_id.id, self.product_id.property_stock_inventory.id)
			//            move_out = self.env['stock.move'].create(move_out_vals)
			//            self.env['stock.quant'].quants_reserve(
			//                [(quant, quant.qty)], move_out)
			//            move_out.action_done()
			//
			//            # get back the quantity from the inventory adjustment location
			//            move_in_vals = self._get_move_values(
			//                quant.qty, self.product_id.property_stock_inventory.id, self.location_id.id)
			//            move_in = self.env['stock.move'].create(move_in_vals)
			//            move_in.action_done()
		})
	h.StockInventoryLine().Methods().GenerateMoves().DeclareMethod(
		`GenerateMoves`,
		func(rs m.StockInventoryLineSet) {
			//        moves = self.env['stock.move']
			//        Quant = self.env['stock.quant']
			//        for line in self:
			//            line._fixup_negative_quants()
			//
			//            if float_utils.float_compare(line.theoretical_qty, line.product_qty, precision_rounding=line.product_id.uom_id.rounding) == 0:
			//                continue
			//            diff = line.theoretical_qty - line.product_qty
			//            if diff < 0:  # found more than expected
			//                vals = line._get_move_values(
			//                    abs(diff), line.product_id.property_stock_inventory.id, line.location_id.id)
			//            else:
			//                vals = line._get_move_values(
			//                    abs(diff), line.location_id.id, line.product_id.property_stock_inventory.id)
			//            move = moves.create(vals)
			//
			//            if diff > 0:
			//                domain = [('qty', '>', 0.0), ('package_id', '=', line.package_id.id), (
			//                    'lot_id', '=', line.prod_lot_id.id), ('location_id', '=', line.location_id.id)]
			//                preferred_domain_list = [[('reservation_id', '=', False)], [(
			//                    'reservation_id.inventory_id', '!=', line.inventory_id.id)]]
			//                quants = Quant.quants_get_preferred_domain(
			//                    move.product_qty, move, domain=domain, preferred_domain_list=preferred_domain_list)
			//                Quant.quants_reserve(quants, move)
			//            elif line.package_id:
			//                move.action_done()
			//                move.quant_ids.write({'package_id': line.package_id.id})
			//                quants = Quant.search([('qty', '<', 0.0), ('product_id', '=', move.product_id.id),
			//                                       ('location_id', '=', move.location_dest_id.id), ('package_id', '!=', False)], limit=1)
			//                if quants:
			//                    for quant in move.quant_ids:
			//                        # To avoid we take a quant that was reconcile already
			//                        if quant.location_id.id == move.location_dest_id.id:
			//                            quant._quant_reconcile_negative(move)
			//        return moves
		})
}
