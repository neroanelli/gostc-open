package service

import (
	"server/pkg/jwt"
	"server/repository"
	"server/service/common/cache"
)

type InfoResp struct {
	Client      int64 `json:"client"`
	Host        int64 `json:"host"`
	Forward     int64 `json:"forward"`
	Tunnel      int64 `json:"tunnel"`
	InputBytes  int64 `json:"inputBytes"`
	OutputBytes int64 `json:"outputBytes"`
}

func (service *service) Info(claims jwt.Claims) (result InfoResp) {
	db, _, _ := repository.Get("")
	result.Client, _ = db.GostClient.Where(db.GostClient.UserCode.Eq(claims.Code)).Count()
	result.Host, _ = db.GostClientHost.Where(db.GostClientHost.UserCode.Eq(claims.Code)).Count()
	result.Forward, _ = db.GostClientForward.Where(db.GostClientForward.UserCode.Eq(claims.Code)).Count()
	result.Tunnel, _ = db.GostClientTunnel.Where(db.GostClientTunnel.UserCode.Eq(claims.Code)).Count()
	obsInfo := cache.GetUserObsDateRange(cache.MONTH_DATEONLY_LIST, claims.Code)
	result.InputBytes = obsInfo.InputBytes
	result.OutputBytes = obsInfo.OutputBytes
	return result
}
