package handlers

import (
	conf "github.com/oleggorj/service-config-data/config-data-util"
	//"gopkg.in/src-d/go-git.v4"
	//log "github.com/oleggorj/service-common-lib/common/logging"
)

type GitHandler struct {
	Environments conf.MappingToEnv
}
