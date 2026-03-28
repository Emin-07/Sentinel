package main

import "github.com/Emin-07/Sentinel/internal/models"

type templateData struct {
	Process   *models.Process
	Processes []*models.Process
}
