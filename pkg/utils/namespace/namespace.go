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

package namespace

//func JoinProcNs(pid int, nsType string) error {
//	filePath := fmt.Sprintf("/proc/%d/ns/%s", pid, nsType)
//	f, err := os.Open(filePath)
//	if err != nil {
//		return fmt.Errorf("open ns file[%s] error: %s", filePath, err.Error())
//	}
//
//	return unix.Setns(int(f.Fd()), 0)
//}

func GetNsOption(namespaces []string) string {
	var nsOptionStr string
	for _, unitNs := range namespaces {
		switch unitNs {
		case MNT:
			nsOptionStr += " -m"
		case PID:
			nsOptionStr += " -p"
		case UTS:
			nsOptionStr += " -u"
		case NET:
			nsOptionStr += " -n"
		case IPC:
			nsOptionStr += " -i"
		}
	}

	return nsOptionStr
}