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

package filesys

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cmdexec"
	"os"
	"strconv"
	"strings"
)

const (
	FileNotFoundKey = "exit code: 1"
)

func getChmodCmd(path, perm string) string {
	return fmt.Sprintf("chmod %s %s", perm, path)
}

func getCheckFileCmd(file string) string {
	return fmt.Sprintf("test -f %s", file)
}

func getCheckDirCmd(dir string) string {
	return fmt.Sprintf("test -d %s", dir)
}

func getPathExistCmd(path string) string {
	return fmt.Sprintf("test -e %s", path)
}

func getAppendFileCmd(path, content string) string {
	return fmt.Sprintf("echo -e \"%s\" >> %s", content, path)
}

func getOverWriteFileCmd(path, content string) string {
	return fmt.Sprintf("echo -en \"%s\" > %s", content, path)
}

func getDeleteLineByKeyCmd(path, key string) string {
	return fmt.Sprintf("sed -i '/%s/d' %s", key, path)
}

func getMkdirForceCmd(dir string) string {
	return fmt.Sprintf("mkdir -p %s", dir)
}

func getRemoveFile(file string) string {
	return fmt.Sprintf("rm %s", file)
}

func getRemoveRF(path string) string {
	return fmt.Sprintf("rm -rf %s", path)
}

func getMoveFile(src, dst string) string {
	return fmt.Sprintf("mv %s %s", src, dst)
}

func getPermCmd(file string) string {
	return "stat -c '%a' " + file
}

func GetPerm(ctx context.Context, cr, cId string, file string) (string, error) {
	if file == "" {
		return "", fmt.Errorf("\"file\" can not be empty")
	}

	perm, err := cmdexec.ExecCommon(ctx, cr, cId, getPermCmd(file))
	return strings.TrimSpace(perm), err
}

func MoveFile(ctx context.Context, cr, cId string, src, dst string) error {
	if src == "" {
		return fmt.Errorf("\"src\" can not be empty")
	}

	if dst == "" {
		return fmt.Errorf("\"dst\" can not be empty")
	}

	_, err := cmdexec.ExecCommon(ctx, cr, cId, getMoveFile(src, dst))
	return err
}

func RemoveFile(ctx context.Context, cr, cId string, file string) error {
	if file == "" {
		return fmt.Errorf("\"file\" can not be empty")
	}

	_, err := cmdexec.ExecCommon(ctx, cr, cId, getRemoveFile(file))
	return err
}

func RemoveRF(ctx context.Context, cr, cId string, path string) error {
	if path == "" {
		return fmt.Errorf("\"path\" can not be empty")
	}

	_, err := cmdexec.ExecCommon(ctx, cr, cId, getRemoveRF(path))
	return err
}

func OverWriteFile(ctx context.Context, cr, cId string, path, content string) error {
	if path == "" {
		return fmt.Errorf("\"path\" can not be empty")
	}

	if content == "" {
		return fmt.Errorf("\"content\" can not be empty")
	}

	_, err := cmdexec.ExecCommon(ctx, cr, cId, getOverWriteFileCmd(path, content))
	return err
}

func MkdirForce(ctx context.Context, cr, cId string, dir string) error {
	if dir == "" {
		return fmt.Errorf("\"dir\" can not be empty")
	}

	_, err := cmdexec.ExecCommon(ctx, cr, cId, getMkdirForceCmd(dir))
	return err
}

func CheckDir(ctx context.Context, cr, cId string, dir string) (bool, error) {
	if dir == "" {
		return false, fmt.Errorf("\"dir\" can not be empty")
	}

	_, err := cmdexec.ExecCommon(ctx, cr, cId, getCheckDirCmd(dir))
	if err != nil {
		if strings.Index(err.Error(), FileNotFoundKey) >= 0 {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func CheckDirLocal(dir string) error {
	f, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("get file error: %s", err.Error())
	}

	if !f.IsDir() {
		return fmt.Errorf("is not a dir")
	}

	return nil
}

func CheckFile(ctx context.Context, cr, cId string, file string) (bool, error) {
	if file == "" {
		return false, fmt.Errorf("\"file\" can not be empty")
	}

	_, err := cmdexec.ExecCommon(ctx, cr, cId, getCheckFileCmd(file))
	if err != nil {
		if strings.Index(err.Error(), FileNotFoundKey) >= 0 {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// ExistPath in container's namespace
func ExistPath(ctx context.Context, cr, cId string, path string) (bool, error) {
	if path == "" {
		return false, fmt.Errorf("\"path\" can not be empty")
	}

	_, err := cmdexec.ExecCommon(ctx, cr, cId, getPathExistCmd(path))
	if err != nil {
		if strings.Index(err.Error(), FileNotFoundKey) >= 0 {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func ExistPathLocal(file string) (bool, error) {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func Chmod(ctx context.Context, cr, cId string, path, perm string) error {
	if path == "" {
		return fmt.Errorf("\"path\" can not be empty")
	}

	if perm == "" {
		return fmt.Errorf("\"perm\" can not be empty")
	}

	_, err := cmdexec.ExecCommon(ctx, cr, cId, getChmodCmd(path, perm))
	return err
}

// DeleteLineByKey in container's namespace
func DeleteLineByKey(ctx context.Context, cr, cId string, path, key string) error {
	if path == "" {
		return fmt.Errorf("\"path\" can not be empty")
	}

	if key == "" {
		return fmt.Errorf("\"key\" can not be empty")
	}

	_, err := cmdexec.ExecCommon(ctx, cr, cId, getDeleteLineByKeyCmd(path, key))
	return err
}

// AppendFile in container's namespace
func AppendFile(ctx context.Context, cr, cId string, path, content string) error {
	if path == "" {
		return fmt.Errorf("\"path\" can not be empty")
	}

	if content == "" {
		return fmt.Errorf("\"path\" can not be empty")
	}

	_, err := cmdexec.ExecCommon(ctx, cr, cId, getAppendFileCmd(path, content))
	return err
}

func IfPathAbs(ctx context.Context, path string) bool {
	path = strings.TrimSpace(path)
	if path[0] != '/' {
		return false
	}

	return true
}

func MkdirP(ctx context.Context, path string) error {
	return cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("mkdir -p %s", path))
}

func MkdirPInContainer(ctx context.Context, cr, cId, path string) error {
	_, err := cmdexec.ExecContainerRaw(ctx, cr, cId, fmt.Sprintf("mkdir -p %s", path))
	return err
}

func RemoveFileInContainer(ctx context.Context, cr, cId, file string) error {
	if file == "" || file == "/" || strings.Index(file, "*") >= 0 {
		return fmt.Errorf("delete file[%s] not allowed", file)
	}

	_, err := cmdexec.ExecContainerRaw(ctx, cr, cId, fmt.Sprintf("rm %s", file))
	return err
}

func GetPermission(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("get stat of path[%s] error: %s", path, err.Error())
	}

	return fmt.Sprintf("%o", info.Mode().Perm()), nil
}

// ExistFile Must be a file
func ExistFile(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	if f.IsDir() {
		return false, nil
	}

	return true, nil
}

func CheckPermission(permission string) error {
	if len(permission) != 3 {
		return fmt.Errorf("len is not equal 3")
	}

	for _, unit := range permission {
		if unit < '0' || unit > '7' {
			return fmt.Errorf("num is not all in [0,7]")
		}
	}

	return nil
}

func GetProMaxFd(ctx context.Context) (int, error) {
	re, err := cmdexec.RunBashCmdWithOutput(ctx, "ulimit -n")
	if err != nil {
		return -1, fmt.Errorf("cmd exec error: %s", err.Error())
	}

	reStr := strings.TrimSpace(re)
	unitMax, err := strconv.Atoi(reStr)
	if err != nil {
		return -1, fmt.Errorf("%s is not a num: %s", reStr, err.Error())
	}

	return unitMax, nil
}

func GetKernelFdStatus(ctx context.Context) (int, int, error) {
	re, err := cmdexec.RunBashCmdWithOutput(ctx, "cat /proc/sys/fs/file-nr | awk '{print $1,$3}'")
	if err != nil {
		return -1, -1, fmt.Errorf("cmd exec error: %s", err.Error())
	}

	reStr := strings.TrimSpace(re)
	reArr := strings.Split(reStr, " ")
	if len(reArr) != 2 {
		return -1, -1, fmt.Errorf("unexpected output: %s", reStr)
	}

	nowFd, err := strconv.Atoi(reArr[0])
	if err != nil {
		return -1, -1, fmt.Errorf("%s is not a num: %s", reArr[0], err.Error())
	}

	maxFd, err := strconv.Atoi(reArr[1])
	if err != nil {
		return -1, -1, fmt.Errorf("%s is not a num: %s", reArr[1], err.Error())
	}

	return nowFd, maxFd, nil
}

func CreateFdFile(ctx context.Context, dir, filePrefix string, count int) error {
	if err := MkdirP(ctx, dir); err != nil {
		return fmt.Errorf("create dir error: %s", err.Error())
	}

	step := 5000
	if step > count {
		step = count
	}

	start, end := 0, step
	for end <= count {
		if err := cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("cd %s && touch %s{%d..%d}", dir, filePrefix, start, end-1)); err != nil {
			return fmt.Errorf("touch file from[%d] to[%d] error: %s", start, end, err.Error())
		}
		start += step
		end += step
	}

	return nil
}
