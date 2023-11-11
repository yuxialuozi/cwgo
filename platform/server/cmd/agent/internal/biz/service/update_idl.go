/*
*
*  * Copyright 2022 CloudWeGo Authors
*  *
*  * Licensed under the Apache License, Version 2.0 (the "License");
*  * you may not use this file except in compliance with the License.
*  * You may obtain a copy of the License at
*  *
*  *     http://www.apache.org/licenses/LICENSE-2.0
*  *
*  * Unless required by applicable law or agreed to in writing, software
*  * distributed under the License is distributed on an "AS IS" BASIS,
*  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  * See the License for the specific language governing permissions and
*  * limitations under the License.
*
 */

package service

import (
	"context"
	"github.com/cloudwego/cwgo/platform/server/cmd/agent/internal/svc"
	agent "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"gorm.io/gorm"
	"net/http"
)

type UpdateIDLService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
} // NewUpdateIDLService new UpdateIDLService
func NewUpdateIDLService(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateIDLService {
	return &UpdateIDLService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Run create note info
func (s *UpdateIDLService) Run(req *agent.UpdateIDLReq) (resp *agent.UpdateIDLRes, err error) {
	err = s.svcCtx.DaoManager.Idl.UpdateIDL(s.ctx, model.IDL{
		Id:              req.Id,
		IdlRepositoryId: req.RepositoryId,
		MainIdlPath:     req.MainIdlPath,
		ServiceName:     req.ServiceName,
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &agent.UpdateIDLRes{
				Code: http.StatusBadRequest,
				Msg:  "idl id not exist",
			}, nil
		}
		return &agent.UpdateIDLRes{
			Code: http.StatusInternalServerError,
			Msg:  "internal err",
		}, nil
	}

	return &agent.UpdateIDLRes{
		Code: 0,
		Msg:  "update idl successfully",
	}, nil
}
