package command

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gocker/src/cgroups/subsystem"
	"gocker/src/container"
	"gocker/src/nsenter"
	"gocker/src/record"
	"gocker/src/run"
	"os"
	"strings"
)

var (
	tty bool
)

var InitGockerCMD = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		cmd := context.Args().Get(0)
		log.Infof("command %s", cmd)
		return container.RunContainerInitProcess()
	},
}

var RunGockerCMD = cli.Command{
	Name:  "run",
	Usage: `Create a container with namespace and cgroups limit mydocker run -ti [command]`,
	Flags: []cli.Flag{
		//交互式
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		//内存限制
		cli.StringFlag{
			Name:  "mem",
			Usage: "memory limit",
		},
		//cpu资源限制
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
		//cpu共享限制
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		//后台运行
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
		},
		//挂载数据卷
		cli.StringFlag{
			Name:  "v",
			Usage: "mount volume",
		},
		//设置容器名称
		cli.StringFlag{
			Name:  "name",
			Usage: "setup a name to container",
		},
		//设置容器网络
		cli.StringFlag{
			Name:  "net",
			Usage: "setup container's internet",
		},
		cli.StringSliceFlag{
			Name:  "e",
			Usage: "setup container environment",
		},
	},
	//  run命令执行的函数
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container command args")
		}
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}

		tty := context.Bool("ti")
		detach := context.Bool("d")
		if tty && detach {
			return fmt.Errorf("can not enable tty and detach at the same time")
		}
		containerName := context.String("name")
		resConfig := &subsystem.ResourceConfig{
			MemoryLimit: context.String("mem"),
			CpuSet:      context.String("cpuShare"),
			CpuShare:    context.String("cpuSet"),
		}
		containerID := container.RandContainerIDGenerator(10)
		volume := context.String("v")
		envSlice := context.StringSlice("e")
		portMappings := []string{""}
		networkName := context.String("net")
		run.Run(tty, cmdArray, resConfig, "gocker", volume, containerName, "./busybox.tar", containerID, envSlice, portMappings, networkName)
		return nil
	},
}

var ExecContainerCMD = cli.Command{
	Name:  "exec",
	Usage: "exec container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 2 {
			return fmt.Errorf("missing exec container id or command")
		}
		//第二次调用直接进入对应容器namespace
		//todo 待支持判断不同exec command
		if os.Getenv(record.EnvExecPid) != "" {
			log.Infof("enter container namespace")
			nsenter.EnterNamespace()
		}
		cID, cmdArray := context.Args().Get(0), strings.Split(context.Args().Get(1), " ")
		container.ExecContainer(cID, cmdArray)
		return nil
	},
}

var ListContainerInfoCMD = cli.Command{
	Name:  "ps",
	Usage: "list all containers info",
	Action: func(context *cli.Context) error {
		container.ListContainerInfo()
		return nil
	},
}

var CommitContainerCMD = cli.Command{
	Name:  "commit",
	Usage: "commit container into image",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing commit container name of a command")
		}
		cID, imageName := context.Args().Get(0), context.Args().Get(1)
		container.CommitContainer(cID, imageName)
		return nil
	},
}

var CheckLogCMD = cli.Command{
	Name:  "logs",
	Usage: "check the logs of a container with cID",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing cID of container when you want to check contaienr logs")
		}
		cID := context.Args().Get(0)
		container.CheckLogsOfContainer(cID)
		return nil
	},
}

var StopContainerCMD = cli.Command{
	Name:  "stop",
	Usage: "stop the running container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing the cID of container when you want to stop a contaienr")
		}
		cID := context.Args().Get(0)
		container.StopContainer(cID)
		return nil
	},
}

var RemoveContainerCMD = cli.Command{
	Name:  "rm",
	Usage: "remove container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing the cID of container when you want to remove a contaienr")
		}
		cID := context.Args().Get(0)
		container.RemoveContainer(cID)
		return nil
	},
}
