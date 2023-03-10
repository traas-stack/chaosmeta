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

package cmdexec

import (
	"fmt"
	"github.com/containerd/cgroups"
	"github.com/traas-stack/chaosmetad/pkg/utils"
	"os"
)

func addToProCgroup(mPid, cPid int) error {
	cgroup, err := cgroups.Load(hierarchy(utils.RootCgroupPath), pidPath(cPid))
	if err != nil {
		return fmt.Errorf("load cgroup of process[%d] error: %s", cPid, err.Error())
	}

	if err = cgroup.Add(cgroups.Process{Pid: mPid}); err != nil {
		return fmt.Errorf("add process[%d] to cgroup error: %s", mPid, err.Error())
	}

	return nil
}

func pidPath(pid int) cgroups.Path {
	p := fmt.Sprintf("/proc/%d/cgroup", pid)
	paths, err := cgroups.ParseCgroupFile(p)
	if err != nil {
		return func(_ cgroups.Name) (string, error) {
			return "", fmt.Errorf("failed to parse cgroup file %s: %s", p, err.Error())
		}
	}

	return func(name cgroups.Name) (string, error) {
		root, ok := paths[string(name)]
		if !ok {
			if root, ok = paths["name="+string(name)]; !ok {
				return "", fmt.Errorf("controller is not supported")
			}

		}
		return root, nil
	}
}

func hierarchy(root string) func() ([]cgroups.Subsystem, error) {
	return func() ([]cgroups.Subsystem, error) {
		subsystems, err := defaults(root)
		if err != nil {
			return nil, err
		}
		var enabled []cgroups.Subsystem
		for _, s := range pathers(subsystems) {
			// check and remove the default groups that do not exist
			if _, err := os.Lstat(s.Path("/")); err == nil {
				enabled = append(enabled, s)
			}
		}
		return enabled, nil
	}
}

// defaults returns all known groups
func defaults(root string) ([]cgroups.Subsystem, error) {
	h, err := cgroups.NewHugetlb(root)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	s := []cgroups.Subsystem{
		cgroups.NewNamed(root, "systemd"),
		cgroups.NewFreezer(root),
		cgroups.NewPids(root),
		cgroups.NewNetCls(root),
		cgroups.NewNetPrio(root),
		cgroups.NewPerfEvent(root),
		cgroups.NewCpuset(root),
		cgroups.NewCpu(root),
		cgroups.NewCpuacct(root),
		cgroups.NewMemory(root),
		cgroups.NewBlkio(root),
		cgroups.NewRdma(root),
	}
	// only add the devices cgroup if we are not in a user namespace
	// because modifications are not allowed
	if !cgroups.RunningInUserNS() {
		s = append(s, cgroups.NewDevices(root))
	}
	// add the hugetlb cgroup if error wasn't due to missing hugetlb
	// cgroup support on the host
	if err == nil {
		s = append(s, h)
	}
	return s, nil
}

type pather interface {
	cgroups.Subsystem
	Path(path string) string
}

func pathers(subystems []cgroups.Subsystem) []pather {
	var out []pather
	for _, s := range subystems {
		if p, ok := s.(pather); ok {
			out = append(out, p)
		}
	}
	return out
}
