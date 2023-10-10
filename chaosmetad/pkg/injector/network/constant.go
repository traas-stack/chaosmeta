/*
 * Copyright 2022-2023 Chaos Meta Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package network

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/log"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/net"
)

const (
	TargetNetwork = "network"
	DirectionOut  = "out"

	FaultOccupy = "occupy"
	OccupyKey   = "chaosmeta_occupy"

	FaultLimit = "limit"

	FaultDelay = "delay"

	FaultLoss = "loss"

	FaultDuplicate = "duplicate"

	FaultCorrupt = "corrupt"

	FaultReorder   = "reorder"
	DefaultGap     = 3
	DefaultLatency = "1s"

	//NetworkExec = "chaosmeta_network"
)

func undoTcWithErr(ctx context.Context, cr, cId string, netInterface, msg string) error {
	if err := execRecover(ctx, cr, cId, netInterface); err != nil {
		log.GetLogger(ctx).Warnf("undo tc rule error: %s", err.Error())
	}

	return fmt.Errorf(msg)
}

func execRecover(ctx context.Context, cr, cId, netInterface string) error {
	isTcExist, err := net.ExistTCRootQdisc(ctx, cr, cId, netInterface)
	if err != nil {
		return fmt.Errorf("check tc rule exist error: %s", err.Error())
	}

	if isTcExist {
		return net.ClearTcRule(ctx, cr, cId, netInterface)
	}

	return nil
}
