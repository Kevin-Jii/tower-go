package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type ThirdPartyRouteService struct {
	routeModule *module.ThirdPartyRouteModule
	storeModule *module.StoreModule
	orderModule *module.ThirdPartyOrderModule
}

func NewThirdPartyRouteService(routeModule *module.ThirdPartyRouteModule, storeModule *module.StoreModule, orderModule *module.ThirdPartyOrderModule) *ThirdPartyRouteService {
	return &ThirdPartyRouteService{
		routeModule: routeModule,
		storeModule: storeModule,
		orderModule: orderModule,
	}
}

func (s *ThirdPartyRouteService) List() ([]*model.ThirdPartyRoute, error) {
	return s.routeModule.List()
}

func (s *ThirdPartyRouteService) Create(req *model.UpsertThirdPartyRouteReq) (*model.ThirdPartyRoute, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("路线名称不能为空")
	}
	row := &model.ThirdPartyRoute{
		Name:   name,
		Remark: strings.TrimSpace(req.Remark),
	}
	if err := s.routeModule.Create(row); err != nil {
		return nil, err
	}
	if err := s.validateStores(req.StoreIDs); err != nil {
		return nil, err
	}
	if err := s.routeModule.ReplaceStores(row.ID, req.StoreIDs); err != nil {
		return nil, err
	}
	return s.routeModule.GetByID(row.ID)
}

func (s *ThirdPartyRouteService) Update(id uint, req *model.UpsertThirdPartyRouteReq) error {
	if _, err := s.routeModule.GetByID(id); err != nil {
		return err
	}
	if err := s.validateStores(req.StoreIDs); err != nil {
		return err
	}
	if err := s.routeModule.Update(id, map[string]interface{}{
		"name":   strings.TrimSpace(req.Name),
		"remark": strings.TrimSpace(req.Remark),
	}); err != nil {
		return err
	}
	return s.routeModule.ReplaceStores(id, req.StoreIDs)
}

func (s *ThirdPartyRouteService) Delete(id uint) error {
	return s.routeModule.Delete(id)
}

func (s *ThirdPartyRouteService) ImportByDateRange(routeID uint, startDate, endDate string) ([]model.RouteImportedProductRow, error) {
	start, end, err := normalizeDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}
	startDate = start.Format("2006-01-02")
	endDate = end.Format("2006-01-02")
	route, err := s.routeModule.GetByID(routeID)
	if err != nil {
		return nil, err
	}
	if len(route.Stores) == 0 {
		return []model.RouteImportedProductRow{}, nil
	}

	accountStoreMap := make(map[uint]*model.Store)
	accountIDs := make([]uint, 0)
	for _, rs := range route.Stores {
		if rs.Store == nil || rs.Store.ThirdPartyAccountID == nil {
			continue
		}
		accountID := *rs.Store.ThirdPartyAccountID
		accountStoreMap[accountID] = rs.Store
		accountIDs = append(accountIDs, accountID)
	}
	if len(accountIDs) == 0 {
		return []model.RouteImportedProductRow{}, nil
	}

	orders, err := s.orderModule.ListByAccountsAndDateRange(accountIDs, startDate, endDate)
	if err != nil {
		return nil, err
	}
	productStoreQty := make(map[string]map[uint]float64)
	productTotalQty := make(map[string]float64)
	storeNameByID := make(map[uint]string)
	for _, rs := range route.Stores {
		if rs.Store != nil {
			storeNameByID[rs.Store.ID] = rs.Store.Name
		}
	}
	for _, order := range orders {
		store := accountStoreMap[order.AccountID]
		if store == nil {
			continue
		}
		parsed := parseOrderRawItems(order.RawJSON)
		for _, p := range parsed {
			if strings.TrimSpace(p.Name) == "" {
				continue
			}
			if _, ok := productStoreQty[p.Name]; !ok {
				productStoreQty[p.Name] = make(map[uint]float64)
			}
			productStoreQty[p.Name][store.ID] += p.Quantity
			productTotalQty[p.Name] += p.Quantity
		}
	}
	rows := make([]model.RouteImportedProductRow, 0, len(productStoreQty))
	for productName, storeQtyMap := range productStoreQty {
		storeQty := make([]model.RouteStoreQuantity, 0, len(route.Stores))
		for _, rs := range route.Stores {
			storeQty = append(storeQty, model.RouteStoreQuantity{
				StoreID:   rs.StoreID,
				StoreName: storeNameByID[rs.StoreID],
				Quantity:  storeQtyMap[rs.StoreID],
			})
		}
		rows = append(rows, model.RouteImportedProductRow{
			ProductName: productName,
			TotalQty:    productTotalQty[productName],
			StoreQty:    storeQty,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].TotalQty > rows[j].TotalQty
	})
	return rows, nil
}

func logisticsSheetModelToMap(row *model.ThirdPartyLogisticsSheet) map[string]interface{} {
	headers := make([]string, 0)
	_ = json.Unmarshal([]byte(row.HeadersJSON), &headers)
	matrix := make([][]float64, 0)
	_ = json.Unmarshal([]byte(row.RowsJSON), &matrix)
	products := make([]string, 0)
	_ = json.Unmarshal([]byte(row.ProductsJSON), &products)
	return map[string]interface{}{
		"id":         row.ID,
		"route_id":   row.RouteID,
		"sheet_date": row.SheetDate,
		"start_date": row.StartDate,
		"end_date":   row.EndDate,
		"headers":    headers,
		"rows":       matrix,
		"products":   products,
		"updated_at": row.UpdatedAt,
		"created_at": row.CreatedAt,
	}
}

func (s *ThirdPartyRouteService) SaveLogisticsSheet(routeID uint, req *model.SaveRouteSheetReq) (map[string]interface{}, error) {
	if _, err := s.routeModule.GetByID(routeID); err != nil {
		return nil, err
	}
	if _, _, err := normalizeDateRange(req.StartDate, req.EndDate); err != nil {
		return nil, err
	}
	headersRaw, _ := json.Marshal(req.Headers)
	rowsRaw, _ := json.Marshal(req.Rows)
	productsRaw, _ := json.Marshal(req.Products)
	sheetDate := time.Now().Format("2006-01-02")
	row := &model.ThirdPartyLogisticsSheet{
		RouteID:      routeID,
		SheetDate:    sheetDate,
		StartDate:    strings.TrimSpace(req.StartDate),
		EndDate:      strings.TrimSpace(req.EndDate),
		HeadersJSON:  string(headersRaw),
		RowsJSON:     string(rowsRaw),
		ProductsJSON: string(productsRaw),
	}
	if err := s.routeModule.SaveLogisticsSheet(row); err != nil {
		return nil, err
	}
	saved, err := s.routeModule.GetLogisticsSheet(routeID, row.StartDate, row.EndDate)
	if err != nil {
		return nil, err
	}
	return logisticsSheetModelToMap(saved), nil
}

func (s *ThirdPartyRouteService) ListLogisticsSheets(routeID uint) ([]map[string]interface{}, error) {
	if _, err := s.routeModule.GetByID(routeID); err != nil {
		return nil, err
	}
	rows, err := s.routeModule.ListLogisticsSheets(routeID)
	if err != nil {
		return nil, err
	}
	resp := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		resp = append(resp, logisticsSheetModelToMap(row))
	}
	return resp, nil
}

func (s *ThirdPartyRouteService) validateStores(storeIDs []uint) error {
	for _, storeID := range storeIDs {
		if _, err := s.storeModule.GetByID(storeID); err != nil {
			return fmt.Errorf("门店不存在: %d", storeID)
		}
	}
	return nil
}

type rawItem struct {
	Name     string
	Quantity float64
}

func parseOrderRawItems(raw string) []rawItem {
	text := strings.TrimSpace(raw)
	if text == "" {
		return nil
	}
	var payload struct {
		ItemList []struct {
			ItemName string  `json:"itemName"`
			SkuName  string  `json:"skuName"`
			ItemNum  float64 `json:"itemNum"`
		} `json:"itemList"`
		RowItemTypeInfoList []struct {
			WrapperUnit string `json:"wrapperUnit"`
		} `json:"rowItemTypeInfoList"`
	}
	if err := json.Unmarshal([]byte(text), &payload); err != nil {
		return nil
	}
	items := make([]rawItem, 0, len(payload.ItemList))
	for _, it := range payload.ItemList {
		name := strings.TrimSpace(it.ItemName)
		if name == "" {
			name = strings.TrimSpace(it.SkuName)
		}
		if name == "" {
			name = "未知商品"
		}
		items = append(items, rawItem{
			Name:     name,
			Quantity: it.ItemNum,
		})
	}
	return items
}

func normalizeDateRange(startDate, endDate string) (time.Time, time.Time, error) {
	startDate = strings.TrimSpace(startDate)
	endDate = strings.TrimSpace(endDate)
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("开始日期格式必须为 YYYY-MM-DD")
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("结束日期格式必须为 YYYY-MM-DD")
	}
	if end.Before(start) {
		return time.Time{}, time.Time{}, errors.New("结束日期不能早于开始日期")
	}
	return start, end, nil
}
