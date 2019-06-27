package stock

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/pool/h"
)

func init() {
	h.BarcodeRule().DeclareModel()

	h.BarcodeRule().AddFields(map[string]models.FieldDefinition{
		"Type": models.SelectionField{
			//selection_add=[('weight', _('Weighted Product')),('location', _('Location')),('lot', _('Lot')),('package', _('Package'))]
		},
	})
}
