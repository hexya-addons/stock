package stock

import (
	"github.com/hexya-addons/base"
	"github.com/hexya-erp/hexya/src/models/security"
	"github.com/hexya-erp/pool/h"
)

//vars

var (
	//Manage multiple stock_locations
	GroupStockMultiLocations *security.Group
	//Manage multiple warehouses
	GroupStockMultiWarehouses *security.Group
	//User
	GroupStockUser *security.Group
	//Manager
	GroupStockManager *security.Group
	//Manage Lots / Serial Numbers
	GroupProductionLot *security.Group
	//Manage Packages
	GroupTrackingLot *security.Group
	//Manage Push and Pull inventory flows
	GroupAdvLocation *security.Group
	//Manage Different Stock Owners
	GroupTrackingOwner *security.Group
	//A warning can be set on a partner (Stock)
	GroupWarningStock *security.Group
)


//rights
func init() {
	h.StockIncoterms().Methods().Load().AllowGroup(security.GroupEveryone)
	h.StockIncoterms().Methods().AllowAllToGroup(GroupStockManager)
	h.StockWarehouse().Methods().AllowAllToGroup(GroupStockManager)
	h.StockWarehouse().Methods().Load().AllowGroup(base.GroupUser)
	h.StockLocation().Methods().Load().AllowGroup(base.GroupPartnerManager)
	h.StockLocation().Methods().AllowAllToGroup(GroupStockManager)
	h.StockLocation().Methods().Load().AllowGroup(base.GroupUser)
	h.StockPicking().Methods().AllowAllToGroup(GroupStockUser)
	h.StockPicking().Methods().AllowAllToGroup(GroupStockManager)
	h.StockPickingType().Methods().Load().AllowGroup(base.GroupUser)
	h.StockPickingType().Methods().Load().AllowGroup(GroupStockUser)
	h.StockPickingType().Methods().AllowAllToGroup(GroupStockManager)
	h.StockProductionLot().Methods().Load().AllowGroup(GroupStockManager)
	h.StockProductionLot().Methods().AllowAllToGroup(GroupStockUser)
	h.StockMove().Methods().AllowAllToGroup(GroupStockManager)
	h.StockMove().Methods().Load().AllowGroup(GroupStockUser)
	h.StockMove().Methods().Write().AllowGroup(GroupStockUser)
	h.StockMove().Methods().Create().AllowGroup(GroupStockUser)
	h.StockInventory().Methods().Load().AllowGroup(GroupStockUser)
	h.StockInventory().Methods().Write().AllowGroup(GroupStockUser)
	h.StockInventory().Methods().Create().AllowGroup(GroupStockUser)
	h.StockInventory().Methods().AllowAllToGroup(GroupStockManager)
	h.StockInventoryLine().Methods().Load().AllowGroup(GroupStockUser)
	h.StockInventoryLine().Methods().Write().AllowGroup(GroupStockUser)
	h.StockInventoryLine().Methods().Create().AllowGroup(GroupStockUser)
	h.StockInventoryLine().Methods().AllowAllToGroup(GroupStockManager)
	h.StockLocation().Methods().Load().AllowGroup(GroupStockManager)
	h.ReportStockLinesDate().Methods().Load().AllowGroup(GroupStockUser)
	h.ReportStockLinesDate().Methods().AllowAllToGroup(GroupStockManager)
	h.Product.ModelProductProduct().Methods().Load().AllowGroup(GroupStockUser)
	h.Product.ModelProductProduct().Methods().Write().AllowGroup(GroupStockUser)
	h.Product.ModelProductProduct().Methods().Create().AllowGroup(GroupStockUser)
	h.Product.ModelProductTemplate().Methods().Load().AllowGroup(GroupStockUser)
	h.Product.ModelProductTemplate().Methods().Write().AllowGroup(GroupStockUser)
	h.Product.ModelProductTemplate().Methods().Create().AllowGroup(GroupStockUser)
	h.Product.ModelProductUomCateg().Methods().AllowAllToGroup(GroupStockManager)
	h.Product.ModelProductUom().Methods().AllowAllToGroup(GroupStockManager)
	h.Product.ModelProductCategory().Methods().AllowAllToGroup(GroupStockManager)
	h.Product.ModelProductTemplate().Methods().AllowAllToGroup(GroupStockManager)
	h.Product.ModelProductProduct().Methods().AllowAllToGroup(GroupStockManager)
	h.Product.ModelProductPackaging().Methods().AllowAllToGroup(GroupStockManager)
	h.Product.ModelProductSupplierinfo().Methods().AllowAllToGroup(GroupStockManager)
	h.Product.ModelProductPricelist().Methods().AllowAllToGroup(GroupStockManager)
	h.Base.ModelIrProperty().Methods().AllowAllToGroup(GroupStockManager)
	h.Base.ModelResPartner().Methods().Load().AllowGroup(GroupStockManager)
	h.Base.ModelResPartner().Methods().Write().AllowGroup(GroupStockManager)
	h.Base.ModelResPartner().Methods().Create().AllowGroup(GroupStockManager)
	h.Product.ModelProductPricelistItem().Methods().AllowAllToGroup(GroupStockManager)
	h.StockWarehouseOrderpoint().Methods().Load().AllowGroup(GroupStockUser)
	h.StockWarehouseOrderpoint().Methods().AllowAllToGroup(GroupStockManager)
	h.StockQuant().Methods().Load().AllowGroup(GroupStockManager)
	h.StockQuant().Methods().Load().AllowGroup(GroupStockUser)
	h.StockQuant().Methods().Load().AllowGroup(base.GroupUser)
	h.StockQuantPackage().Methods().Load().AllowGroup(base.GroupUser)
	h.StockQuantPackage().Methods().AllowAllToGroup(GroupStockManager)
	h.StockQuantPackage().Methods().AllowAllToGroup(GroupStockUser)
	h.ProcurementRule().Methods().Load().AllowGroup(GroupStockUser)
	h.ProcurementRule().Methods().AllowAllToGroup(GroupStockManager)
	h.ProcurementRule().Methods().AllowAllToGroup(GroupStockManager)
	h.StockLocationPath().Methods().Load().AllowGroup(GroupStockUser)
	h.StockLocationPath().Methods().Load().AllowGroup(base.GroupUser)
	h.StockLocationPath().Methods().AllowAllToGroup(GroupStockUser)
	h.StockLocationRoute().Methods().AllowAllToGroup(GroupStockManager)
	h.StockLocationRoute().Methods().Load().AllowGroup(base.GroupUser)
	h.ProcurementRule().Methods().Load().AllowGroup(base.GroupUser)
	h.StockPackOperation().Methods().AllowAllToGroup(GroupStockManager)
	h.StockPackOperation().Methods().AllowAllToGroup(GroupStockUser)
	h.StockPackOperation().Methods().AllowAllToGroup(base.GroupUser)
	h.ProductPutaway().Methods().Load().AllowGroup(base.GroupUser)
	h.ProductPutaway().Methods().AllowAllToGroup(GroupStockManager)
	h.ProductRemoval().Methods().Load().AllowGroup(base.GroupUser)
	h.StockFixedPutawayStrat().Methods().AllowAllToGroup(GroupStockManager)
	h.StockFixedPutawayStrat().Methods().Load().AllowGroup(GroupStockUser)
	h.StockMoveOperationLink().Methods().AllowAllToGroup(GroupStockManager)
	h.StockMoveOperationLink().Methods().AllowAllToGroup(GroupStockUser)
	h.StockMoveOperationLink().Methods().AllowAllToGroup(base.GroupUser)
	h.StockPackOperationLot().Methods().AllowAllToGroup(GroupStockManager)
	h.StockPackOperationLot().Methods().AllowAllToGroup(GroupStockUser)
	h.StockPackOperationLot().Methods().AllowAllToGroup(base.GroupUser)
	h.Product.ModelProductPriceHistory().Methods().Load().AllowGroup(GroupStockUser)
	h.Product.ModelProductPriceHistory().Methods().Create().AllowGroup(GroupStockUser)
	h.Product.ModelProductPriceHistory().Methods().AllowAllToGroup(GroupStockManager)
	h.Barcodes.ModelBarcodeNomenclature().Methods().Load().AllowGroup(GroupStockUser)
	h.Barcodes.ModelBarcodeNomenclature().Methods().AllowAllToGroup(GroupStockManager)
	h.Barcodes.ModelBarcodeRule().Methods().Load().AllowGroup(GroupStockUser)
	h.Barcodes.ModelBarcodeRule().Methods().AllowAllToGroup(GroupStockManager)
	h.ReportStockForecast().Methods().Load().AllowGroup(GroupStockUser)
	h.ReportStockForecast().Methods().AllowAllToGroup(GroupStockManager)
	h.StockScrap().Methods().Load().AllowGroup(GroupStockUser)
	h.StockScrap().Methods().Write().AllowGroup(GroupStockUser)
	h.StockScrap().Methods().Create().AllowGroup(GroupStockUser)
	h.Product.ModelProductAttribute().Methods().AllowAllToGroup(GroupStockManager)
	h.Product.ModelProductAttributeValue().Methods().AllowAllToGroup(GroupStockManager)
	h.Product.ModelProductAttributePrice().Methods().AllowAllToGroup(GroupStockManager)
	h.Product.ModelProductAttributeLine().Methods().AllowAllToGroup(GroupStockManager)
}
