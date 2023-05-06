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

package main

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/log"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cmdexec"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/errutil"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/filesys"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	FaultFileAdd    = "add"
	FaultFileAppend = "append"
	FaultFileMv     = "mv"
	FaultFileDelete = "del"
	FaultFileChmod  = "chmod"
)

const (
	BackUpDir = "/tmp/chaosmeta_backup_file"
)

// [func] [fault] [level] [args]
func main() {
	var (
		err                       error
		fName, fault, level, args = os.Args[1], os.Args[2], os.Args[3], os.Args[4:]
		ctx                       = context.Background()
	)
	log.Level = level

	switch fName {
	case utils.MethodValidator:
		err = execValidator(ctx, fault, args)
	case utils.MethodInject:
		err = execInject(ctx, fault, args)
	case utils.MethodRecover:
		err = execRecover(ctx, fault, args)
	default:
		errutil.ExitExpectedErr(fmt.Sprintf("not support method: %s", fName))
	}

	if err != nil {
		errutil.ExitExpectedErr(err.Error())
	}
}

func execValidator(ctx context.Context, fault string, args []string) error {
	switch fault {
	case FaultFileAdd:
		force, err := strconv.ParseBool(args[1])
		if err != nil {
			return fmt.Errorf("\"force\" is not a bool: %s", args[1])
		}

		return validatorAdd(ctx, args[0], force)
	case FaultFileAppend:
		return validatorAppend(ctx, args[0])
	case FaultFileMv:
		return validatorMv(ctx, args[0], args[1])
	case FaultFileChmod:
		return validatorChmod(ctx, args[0])
	case FaultFileDelete:
		return validatorDelete(ctx, args[0])
	default:
		return fmt.Errorf("not support fault: %s", fault)
	}
}

func execInject(ctx context.Context, fault string, args []string) error {
	switch fault {
	case FaultFileAdd:
		return injectAdd(ctx, args[0], args[1], args[2])
	case FaultFileAppend:
		raw, err := strconv.ParseBool(args[2])
		if err != nil {
			return fmt.Errorf("\"raw\" is not a bool: %s", args[1])
		}

		return injectAppend(ctx, args[0], args[1], raw, args[3])
	case FaultFileMv:
		return injectMv(ctx, args[0], args[1])
	case FaultFileChmod:
		return injectChmod(ctx, args[0], args[1])
	case FaultFileDelete:
		return injectDelete(ctx, args[0], args[1])
	default:
		return fmt.Errorf("not support fault: %s", fault)
	}
}

func execRecover(ctx context.Context, fault string, args []string) error {
	switch fault {
	case FaultFileAdd:
		return recoverAdd(ctx, args[0])
	case FaultFileAppend:
		raw, err := strconv.ParseBool(args[2])
		if err != nil {
			return fmt.Errorf("\"raw\" is not a bool: %s", args[1])
		}
		return recoverAppend(ctx, args[0], args[1], raw)
	case FaultFileMv:
		return recoverMv(ctx, args[0], args[1])
	case FaultFileChmod:
		return recoverChmod(ctx, args[0], args[1])
	case FaultFileDelete:
		return recoverDelete(ctx, args[0], args[1])
	default:
		return fmt.Errorf("not support fault: %s", fault)
	}
}

func validatorAdd(ctx context.Context, path string, force bool) error {
	isPathExist, err := filesys.ExistPath(path)
	if err != nil {
		return fmt.Errorf("\"path\"[%s] check exist error: %s", path, err.Error())
	}

	dir := filepath.Dir(path)
	isDirExist, err := filesys.ExistPath(dir)
	if err != nil {
		return fmt.Errorf("check dir[%s] exist error: %s", dir, err.Error())
	}

	if isPathExist {
		isFile, _ := filesys.ExistFile(path)
		if !isFile {
			return fmt.Errorf("\"path\"[%s] is an existed dir", path)
		}
	}

	if !force {
		if isPathExist {
			return fmt.Errorf("file[%s] exist, if want to force to overwrite, please provide [-f] or [--force] args", path)
		}

		if !isDirExist {
			return fmt.Errorf("dir[%s] is not exist, if want to auto create, please provide [-f] or [--force] args", dir)
		}
	}

	return nil
}

func injectAdd(ctx context.Context, path, perm, content string) error {
	logger := log.GetLogger(ctx)

	dir := filepath.Dir(path)
	isDirExist, _ := filesys.ExistPath(dir)

	if !isDirExist {
		if err := filesys.MkdirP(ctx, dir); err != nil {
			return fmt.Errorf("mkdir dir[%s] error: %s", dir, err.Error())
		}
	}

	if err := cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("echo -en \"%s\" > %s", content, path)); err != nil {
		return fmt.Errorf("add content to %s error: %s", path, err.Error())
	}

	if perm != "" {
		if err := filesys.Chmod(ctx, path, perm); err != nil {
			if err := recoverAdd(ctx, path); err != nil {
				logger.Warnf("undo error: %s", err.Error())
			}

			return fmt.Errorf("chmod file[%s] to[%s] error: %s", path, perm, err.Error())
		}
	}

	return nil
}

func recoverAdd(ctx context.Context, path string) error {
	isExist, err := filesys.ExistPath(path)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", path, err.Error())
	}

	if isExist {
		return os.Remove(path)
	}

	return nil
}

func validatorAppend(ctx context.Context, path string) error {
	isFileExist, err := filesys.ExistFile(path)
	if err != nil {
		return fmt.Errorf("\"path\"[%s] check exist error: %s", path, err.Error())
	}

	if !isFileExist {
		return fmt.Errorf("file[%s] is not exist", path)
	}

	return nil
}

func injectAppend(ctx context.Context, path, uid string, raw bool, content string) error {
	flag := getAppendFlag(uid)

	if !raw {
		content = strings.ReplaceAll(content, "\\n", "\n")
		content = fmt.Sprintf("%s%s", strings.ReplaceAll(content, "\n", fmt.Sprintf("%s\n", flag)), flag)
	}

	content = fmt.Sprintf("\n%s", content)
	if err := cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("echo -e \"%s\" >> %s", content, path)); err != nil {
		return fmt.Errorf("append content to %s error: %s", path, err.Error())
	}

	return nil
}

func recoverAppend(ctx context.Context, path, uid string, raw bool) error {
	if raw {
		return nil
	}

	fileExist, err := filesys.ExistFile(path)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", path, err.Error())
	}

	if !fileExist {
		return nil
	}

	flag := getAppendFlag(uid)
	isExist, err := filesys.HasFileLineByKey(ctx, flag, path)
	if err != nil {
		return fmt.Errorf("check file[%s] line exist key[%s] error: %s", path, flag, err.Error())
	}

	if isExist {
		return cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("sed -i '/%s/d' %s", getAppendFlag(uid), path))
	}

	return nil
}

func validatorChmod(ctx context.Context, path string) error {
	isPathExist, err := filesys.ExistFile(path)
	if err != nil {
		return fmt.Errorf("\"path\"[%s] check exist error: %s", path, err.Error())
	}

	if !isPathExist {
		return fmt.Errorf("\"path\"[%s] is not an existed file", path)
	}

	return nil
}

func injectChmod(ctx context.Context, path, permission string) error {
	return filesys.Chmod(ctx, path, permission)
}

func recoverChmod(ctx context.Context, path, perm string) error {
	isExist, err := filesys.ExistFile(path)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", path, err.Error())
	}

	if isExist {
		return filesys.Chmod(ctx, path, perm)
	}

	return nil
}

func validatorDelete(ctx context.Context, path string) error {
	isPathExist, err := filesys.ExistFile(path)
	if err != nil {
		return fmt.Errorf("\"path\"[%s] check exist error: %s", path, err.Error())
	}

	if !isPathExist {
		return fmt.Errorf("\"path\"[%s] is not an existed file", path)
	}

	return nil
}

func injectDelete(ctx context.Context, path, uid string) error {
	backupDir := getBackupDir(uid)
	if err := filesys.MkdirP(ctx, backupDir); err != nil {
		return fmt.Errorf("create backup dir[%s] error: %s", backupDir, err.Error())
	}

	return os.Rename(path, fmt.Sprintf("%s/%s", backupDir, filepath.Base(path)))
}

func recoverDelete(ctx context.Context, path, uid string) error {
	backupDir := getBackupDir(uid)

	isExist, err := filesys.ExistPath(path)
	if err != nil {
		return fmt.Errorf("check path[%s] exist error: %s", path, err.Error())
	}
	if !isExist {
		backupFile := fmt.Sprintf("%s/%s", backupDir, filepath.Base(path))
		if err := os.Rename(backupFile, path); err != nil {
			return fmt.Errorf("mv from[%s] to[%s] error: %s", backupFile, path, err.Error())
		}
	}

	return os.Remove(backupDir)
}

func validatorMv(ctx context.Context, src, dst string) error {
	isPathExist, err := filesys.ExistFile(src)
	if err != nil {
		return fmt.Errorf("\"src\"[%s] check exist error: %s", src, err.Error())
	}

	if !isPathExist {
		return fmt.Errorf("\"src\"[%s] is not an existed file", src)
	}

	isPathExist, err = filesys.ExistPath(dst)
	if err != nil {
		return fmt.Errorf("\"dst\"[%s] check exist error: %s", dst, err.Error())
	}

	if isPathExist {
		return fmt.Errorf("\"dst\"[%s] is existed", dst)
	}

	return nil
}

func injectMv(ctx context.Context, src, dst string) error {
	return os.Rename(src, dst)
}

func recoverMv(ctx context.Context, src, dst string) error {
	isExist, err := filesys.ExistPath(src)
	if err != nil {
		return fmt.Errorf("check src[%s] exist error: %s", src, err.Error())
	}

	if isExist {
		return nil
	}

	return os.Rename(dst, src)
}

func getAppendFlag(uid string) string {
	return fmt.Sprintf(" %s-%s", utils.RootName, uid)
}

func getBackupDir(uid string) string {
	return fmt.Sprintf("%s%s", BackUpDir, uid)
}
