package data

import (
	"encoding/json"
	"go-stock/backend/db"
	"go-stock/backend/logger"
	"testing"
)

// @Author spark
// @Date 2025/4/23 17:58
// @Desc
//-----------------------------------------------------------------------------------

func TestGetSinaNews(t *testing.T) {
	db.Init("../../data/stock.db")
	NewMarketNewsApi().GetSinaNews(30)
	//NewMarketNewsApi().GetNewTelegraph(30)

}

func TestGlobalStockIndexes(t *testing.T) {
	resp := NewMarketNewsApi().GlobalStockIndexes(30)
	bytes, err := json.Marshal(resp)
	if err != nil {
		return
	}
	logger.SugaredLogger.Debugf("resp: %+v", string(bytes))
}
