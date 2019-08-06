package utils

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-ps"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

/*
//#include <libproc.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>

void sayHello(){
	printf("proc world");
}
//
//int getPIDPath ( char* pid)
//{
//   pid_t pid; int ret;
//   char pathbuf[PROC_PIDPATHINFO_MAXSIZE];
//
//
//       pid = (pid_t)atoi(pid);
//       ret = proc_pidpath (pid, pathbuf, sizeof(pathbuf));
//       if ( ret <= 0 ) {
//           fprintf(stderr, "PID %d: proc_pidpath ();\n", pid);
//           fprintf(stderr, "    %s\n", strerror(errno));
//       } else {
//           printf("proc %d: %s\n", pid, pathbuf);
//       }
//
//
//   return 0;
//}

// 仅适用于 linux
char* getPIDPath(int pid){
	char file[32];
	char* buf = malloc(sizeof(char)*64);
	pid_t tpid = (pid_t)pid;
	sprintf(file, "/proc/%i/cmdline", tpid);
	FILE *f = fopen(file, "r");
	fgets(buf, 64, f);
	fclose(f);
	//printf("\n=== %s",buf);
	return buf;
}

#cgo CFLAGS: -g

*/
import "C"

func GetProcessPath(pid uint32) (string){

	//i, _ := strconv.Atoi(os.Args[1])
	cs := C.getPIDPath(C.int(int(pid)))
	path :=C.GoString(cs)
	C.free(unsafe.Pointer(cs))

	//fmt.Println(">> ",path)
	return path
}

func GetProcessInfo(pid uint32) (ps.Process,error){
	//NOTE : syscall.Signal is not available in Windows
	if runtime.GOOS != "windows" {
		_, err := getProcessRunningStatus(int(pid))

		if err != nil {
			fmt.Println("Error : ", err)
			//os.Exit(-1)
			return nil, err
		}

	}

	// at this stage the Processes related functions found in Golang's OS package
	// is no longer sufficient, we will use Mitchell Hashimoto's https://github.com/mitchellh/go-ps
	// package to find the application/executable/binary name behind the process ID.

	p, err := ps.FindProcess(int(pid))

	if err != nil {
		fmt.Println("Error : ", err)
		return nil, err
	}

	fmt.Println(p.Executable())
	return p,nil
}

func getProcessRunningStatus(pid int) (*os.Process, error) {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil, err
	}

	//double check if process is running and alive
	//by sending a signal 0
	//NOTE : syscall.Signal is not available in Windows

	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return proc, nil
	}

	if err == syscall.ESRCH {
		return nil, errors.New("process not running")
	}

	// default
	return nil, errors.New("process running but query operation not permitted")
}