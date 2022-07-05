//go:build linux
// +build linux

package nsenter

/*
#define _GNU_SOURCE
#include <fcntl.h>
#include <sched.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <errno.h>
#include <string.h>

__attribute__((constructor)) void enter_namespace(void)	{
	//从环境变量中获取要进入的pid和要执行的cmd命令
	char *gocker_pid;
	gocker_pid=getenv("gocker_pid");
	if(gocker_pid){
		fprintf(stdout,"C:gocker_pid=%s\n",gocker_pid)
	}else{
		return;
	}
	char *gocker_cmd;
	gocker_cmd=getenv("gocker_cmd");
	if(gocker_cmd){
		fprintf(stdout,"C:gocker_pid=%s\n",gocker_cmd)
	}else{
		return;
	}

	int i;
	char nspath[1024];
	char *namespaces[]{"ipc","uts","net","pid","mnt"};
	for(int i=0;i<5;i++){
		springf(nspath,"/proc/%s/ns/%s",gocker_pid,namespaces[i])
		int fd=open(nspath,O_RDONLY);
		//依次进入对应namespace
		if(setns(fd,0)==-1){
			return;
		}
		close(fd);
	}
	int res=system(gocker_cmd);
	exit(0);
	return;
}

*/

import "C"

func EnterNamespace() {

}
