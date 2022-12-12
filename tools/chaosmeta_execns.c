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

#include <sys/types.h>
#include <stdlib.h>
#include <unistd.h>
#include <stdio.h>
#include <errno.h>
#include <sys/stat.h>
#include <sys/syscall.h>
#include <fcntl.h>
#include <getopt.h>

int enter_ns(int pid, const char* type) {
    char path[64], selfpath[64];
    snprintf(path, sizeof(path), "/proc/%d/ns/%s", pid, type);
    snprintf(selfpath, sizeof(selfpath), "/proc/self/ns/%s", type);

    struct stat oldns_stat, newns_stat;
    int oldre = stat(selfpath, &oldns_stat);
    int newre = stat(path, &newns_stat);
    if (oldre != 0) {
        fprintf(stderr, "stat self namespace file[%s] error\n", selfpath);
        return oldre;
    }

    if (newre != 0) {
        fprintf(stderr, "stat target namespace file[%s] error\n", path);
        return oldre;
    }

    if (oldns_stat.st_ino != newns_stat.st_ino) {
        int newns = open(path, O_RDONLY);
        if (newns < 0) {
            fprintf(stderr, "open target file[%s] error\n", path);
            return newns;
        }

        int result = syscall(__NR_setns, newns, 0);
        close(newns);
        if (result < 0) {
            fprintf(stderr, "setns error\n");
            return result;
        }
    }

    return 0;
}

int main(int argc, char *argv[]) {
    int opt;
    char *cmd;
    int target = 0;
    int ipcns = 0;
    int utsns = 0;
    int netns = 0;
    int pidns = 0;
    int mntns = 0;
    char *string = "c:t:mpuni";

    while((opt =getopt(argc, argv, string))!= -1) {
        switch (opt) {
            case 'c':
                cmd = optarg;
                break;
            case 't':
                target = atoi(optarg);
                break;
            case 'm':
                mntns = 1;
                break;
            case 'p':
                pidns = 1;
                break;
            case 'u':
                utsns = 1;
                break;
            case 'n':
                netns = 1;
                break;
            case 'i':
                ipcns = 1;
                break;
            default:
                break;
        }
    }

    if (target <= 0) {
        fprintf(stderr, "%s is not a valid process ID\n", target);
        return 1;
    }

    if (!cmd) {
        fprintf(stderr, "cmd args is empty\n");
        return 1;
    }

    if(ipcns) {
        int re = enter_ns(target, "ipc");
        if (re != 0) {
            return re;
        }
    }

    if(utsns) {
        int re = enter_ns(target, "uts");
        if (re != 0) {
            return re;
        }
    }

    if(netns) {
        int re = enter_ns(target, "net");
        if (re != 0) {
            return re;
        }
    }

    if(pidns) {
        int re = enter_ns(target, "pid");
        if (re != 0) {
            return re;
        }
    }

    if(mntns) {
        int re = enter_ns(target, "mnt");
        if (re != 0) {
            return re;
        }
    }

    int re = system(cmd);
    if (re < 0) {
        fprintf(stderr, "cmd exec error\n");
        return re;
    }

    return 0;
}
