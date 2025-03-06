package service

import (
	"server/pkg/bean"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
	"time"
)

type PageReq struct {
	bean.PageParam
}

type Item struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Ip          string `json:"ip"`
	Online      int    `json:"online"`
	LastTime    string `json:"lastTime"`
	Version     string `json:"version"`
	CreatedAt   string `json:"createdAt"`
	InputBytes  int64  `json:"inputBytes"`
	OutputBytes int64  `json:"outputBytes"`
}

func (service *service) Page(claims jwt.Claims, req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")
	clients, total, _ := db.GostClient.Where(db.GostClient.UserCode.Eq(claims.Code)).Order(db.GostClient.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	for _, client := range clients {
		obsInfo := cache.GetClientObsDateRange(cache.MONTH_DATEONLY_LIST, client.Code)
		list = append(list, Item{
			Code:        client.Code,
			Name:        client.Name,
			Key:         client.Key,
			Ip:          cache.GetClientIp(client.Code),
			Online:      utils.TrinaryOperation(cache.GetClientOnline(client.Code), 1, 2),
			LastTime:    cache.GetClientLastTime(client.Code),
			Version:     cache.GetClientVersion(client.Code),
			CreatedAt:   client.CreatedAt.Format(time.DateTime),
			InputBytes:  obsInfo.InputBytes,
			OutputBytes: obsInfo.OutputBytes,
		})
	}
	return list, total
}
