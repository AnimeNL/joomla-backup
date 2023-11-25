package main

import (
	"joomla-backup/internal/config"
	"os"

	log "github.com/sirupsen/logrus"
)

func cleanup() {
	log.Infof("cleanup")
	cleanWorkdir()
}

func cleanWorkdir() {
	log.Info("remove old dumps")

	log.Debugf("Reading dir %s", config.Configuration.Paths.DatabaseDumps)
	files, err := os.ReadDir(config.Configuration.Paths.DatabaseDumps)
	if err != nil {
		log.Errorf("Error reading dbdump dir: %s", err.Error())
	}

	for _, file := range files {
		log.Infof("Deleting file: %s", file.Name())
		if err := os.Remove(file.Name()); err != nil {
			log.Errorf("error deleting file %s. Error: %s", file.Name(), err.Error())
		}
	}
}
