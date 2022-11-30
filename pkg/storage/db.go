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

package storage

import (
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"path"
)

const storageFile = "chaosmetad.dat"

type dbStorage struct {
	*gorm.DB
}

func newDBStorage() (*dbStorage, error) {
	// TODO: db path can be config
	dsn := path.Join(utils.GetRunPath(), storageFile)

	dsn += "?cache=shared"

	gormDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %s", err.Error())
	}

	tempDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB: %s", err.Error())
	}
	tempDB.SetMaxOpenConns(1)

	db := &dbStorage{
		gormDB,
	}

	return db, nil
}
