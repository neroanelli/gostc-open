package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
	"server/service/common/node_port"
	"server/service/gost_engine"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(claims jwt.Claims, req DeleteReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		forward, _ := tx.GostClientForward.Where(
			tx.GostClientForward.Code.Eq(req.Code),
			tx.GostClientForward.UserCode.Eq(claims.Code),
		).First()
		if forward == nil {
			return nil
		}
		_, _ = tx.GostNodePort.Where(tx.GostNodePort.Port.Eq(forward.Port), tx.GostNodePort.NodeCode.Eq(forward.NodeCode)).Delete()
		_, _ = tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(forward.Code)).Delete()
		if _, err := tx.GostClientForward.Where(tx.GostClientForward.Code.Eq(forward.Code)).Delete(); err != nil {
			log.Error("删除用户端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		node_port.ReleasePort(forward.NodeCode, forward.Port)
		gost_engine.ClientRemoveForwardConfig(*forward, forward.Node)
		cache.DelTunnelInfo(req.Code)
		return nil
	})
}
