package main

//TODO: Complete log statements
import (
	"context"
	"time"

	"joomla-backup/internal/config"

	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

var (
	dc  = config.Configuration.DockerClient
	ctx context.Context
)

func setup() {
	ctx = context.Background()

	cleanup()
}

func databaseDump(ctx context.Context, database string) {
	command := []string{"bash", "-c", "/usr/bin/mysqldump -u " + config.Configuration.Database.Credentials.Username + " --password=" + config.Configuration.Database.Credentials.Password + " " + database + " > /dump/" + database + ".sql"}
	log.Debugf("constructed docker exec command: %v", command)
	execConfig := types.ExecConfig{Tty: false, AttachStdout: true, AttachStderr: false, Cmd: command}
	respIdExecCreate, err := dc.ContainerExecCreate(ctx, "mysql", execConfig)
	if err != nil {
		log.Errorf("error creating db dump command: %v", err)
	}
	err = dc.ContainerExecStart(ctx, respIdExecCreate.ID, types.ExecStartCheck{})
	if err != nil {
		log.Errorf("error occured starting db dump: %v", err)
	}

	execStatus, err := dc.ContainerExecInspect(ctx, respIdExecCreate.ID)
	if err != nil {
		log.Errorf("error occured inspecting dump progress: %v", err)
	}
	for execStatus.Running {
		log.Info("waiting for db dump to finish...")
		time.Sleep(2 * time.Second)
	}
}

func main() {
	setup()

	// Dump all databases
	for _, database := range config.Configuration.Database.Databases {
		log.Infof("dumping database %v", database)
		databaseDump(ctx, database) //TODO: Use goroutines to dump databases in parallel and make the backup more efficent
	}

	log.Info("done. exiting.")
}
