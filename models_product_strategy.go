package stock

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/pool/h"
)

func init() {
	h.ProductRemoval().DeclareModel()

	h.ProductRemoval().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{
			String:   "Name",
			Required: true,
		},
		"Method": models.CharField{
			String:   "Method",
			Required: true,
			Help:     "FIFO, LIFO...",
		},
	})
	h.ProductPutaway().DeclareModel()

	h.ProductPutaway().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{
			String:   "Name",
			Required: true,
		},
		"Method": models.SelectionField{
			Selection: "_get_putaway_options",
			String:    "Method",
			Default:   models.DefaultValue("fixed"),
			Required:  true,
		},
		"FixedLocationIds": models.One2ManyField{
			RelationModel: h.StockFixedPutawayStrat(),
			ReverseFK:     "",
			String:        "Fixed Locations Per Product Category",
			NoCopy:        false,
			Help: "When the method is fixed, this location will be used to" +
				"store the products",
		},
	})
	h.ProductPutaway().Methods().GetPutawayOptions().DeclareMethod(
		`GetPutawayOptions`,
		func(rs m.ProductPutawaySet) {
			//        return [('fixed', 'Fixed Location')]
		})
	h.ProductPutaway().Methods().PutawayApplyFixed().DeclareMethod(
		`PutawayApplyFixed`,
		func(rs m.ProductPutawaySet, product interface{}) {
			//        for strat in self.fixed_location_ids:
			//            categ = product.categ_id
			//            while categ:
			//                if strat.category_id.id == categ.id:
			//                    return strat.fixed_location_id.id
			//                categ = categ.parent_id
			//        return self.env['stock.location']
		})
	h.ProductPutaway().Methods().PutawayApply().DeclareMethod(
		`PutawayApply`,
		func(rs m.ProductPutawaySet, product interface{}) {
			//        if hasattr(self, '_putaway_apply_%s' % (self.method)):
			//            return getattr(self, '_putaway_apply_%s' % (self.method))(product)
			//        return self.env['stock.location']
		})
	h.StockFixedPutawayStrat().DeclareModel()

	h.StockFixedPutawayStrat().AddFields(map[string]models.FieldDefinition{
		"PutawayId": models.Many2OneField{
			RelationModel: h.ProductPutaway(),
			String:        "Put Away Method",
			Required:      true,
		},
		"CategoryId": models.Many2OneField{
			RelationModel: h.ProductCategory(),
			String:        "Product Category",
			Required:      true,
		},
		"FixedLocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Location",
			Required:      true,
		},
		"Sequence": models.IntegerField{
			String: "Priority",
			Help: "Give to the more specialized category, a higher priority" +
				"to have them in top of the list.",
		},
	})
}
