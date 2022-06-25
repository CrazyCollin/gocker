package command

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"mini-docker/src/cgroups/subsystem"
	"mini-docker/src/container"
	"mini-docker/src/run"
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
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "mem",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			fmt.Errorf("missing container command args")
		}
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		tty := context.Bool("ti")
		resConfig := &subsystem.ResourceConfig{
			MemoryLimit: context.String("mem"),
			CpuSet:      context.String("cpuShare"),
			CpuShare:    context.String("cpuSet"),
		}
		run.Run(tty, cmdArray, resConfig)
		return nil
	},
}
