package main

import (
	"github.com/firerockets/gator/internal/config"
	"github.com/firerockets/gator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}
