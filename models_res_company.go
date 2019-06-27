package stock

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/pool/h"
)

func init() {

	h.Company().AddFields(map[string]models.FieldDefinition{
		"PropagationMinimumDelta": models.IntegerField{
			String:  "Minimum Delta for Propagation of a Date Change on moves linked together",
			Default: models.DefaultValue(1),
		},
		"InternalTransitLocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Internal Transit Location",
			//on_delete="restrict"
			Help: "Technical field used for resupply routes between warehouses" +
				"that belong to this company",
		},
	})
	h.Company().Methods().CreateTransitLocation().DeclareMethod(
		`Create a transit location with company_id being the given
company_id. This is needed
           in case of resuply routes between warehouses
belonging to the same company, because
           we don't want to create accounting entries at that time.
        `,
		func(rs m.CompanySet) {
			//        parent_location = self.env.ref(
			//            'stock.stock_location_locations', raise_if_not_found=False)
			//        location = self.env['stock.location'].create({
			//            'name': _('%s: Transit Location') % self.name,
			//            'usage': 'transit',
			//            'location_id': parent_location and parent_location.id or False,
			//        })
			//        location.sudo().write({'company_id': self.id})
			//        self.write({'internal_transit_location_id': location.id})
		})
	h.Company().Methods().Create().Extend(
		`Create`,
		func(rs m.CompanySet, vals models.RecordData) {
			//        company = super(Company, self).create(vals)
			//        self.env['stock.warehouse'].check_access_rights('create')
			//        self.env['stock.warehouse'].sudo().create(
			//            {'name': company.name, 'code': company.name[:5], 'company_id': company.id})
			//        company.create_transit_location()
			//        return company
		})
}
