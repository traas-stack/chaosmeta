#include<unistd.h>
#include<stdlib.h>
#include<stdio.h>

int main( int argc, char *argv[])  
{
    int count = atoi(argv[2]);
    for(int i=0; i< count; i++) {
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