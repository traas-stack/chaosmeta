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

#include<unistd.h>
#include<stdlib.h>
#include<stdio.h>

int main( int argc, char *argv[])  
{
    int count = atoi(argv[2]);
    int i = 0;
    for(i = 0; i < count; i++) {
        vfork();
    }

    printf("[success]inject success\n");
    for (;;) {
        sleep(86400);
    }

   return 0;
}

/* 
    g++ chaosmeta_cpuload.cpp -o chaosmeta_cpuload
    ./chaosmeta_cpuload 16
    kill -9 -- -PGID
    ps j -A | grep cpuload
 */
// TODO： can consider using golang sc：r1, r2, err := syscall.Syscall(syscall.SYS_VFORK, 0, 0, 0)