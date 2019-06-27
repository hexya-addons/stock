package stock

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/pool/h"
)

func init() {
	h.StockIncoterms().DeclareModel()

	h.StockIncoterms().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{
			String:    "Name",
			Required:  true,
			Translate: true,
			Help: "Incoterms are series of sales terms. They are used to divide" +
				"transaction costs and responsibilities between buyer and" +
				"seller and reflect state-of-the-art transportation practices.",
		},
		"Code": models.CharField{
			String:   "Code",
			Size:     3,
			Required: true,
			Help:     "Incoterm Standard Code",
		},
		"Active": models.BooleanField{
			String:  "Active",
			Default: models.DefaultValue(true),
			Help: "By unchecking the active field, you may hide an INCOTERM" +
				"you will not use.",
		},
	})
}
