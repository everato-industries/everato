package event

import (
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/jackc/pgx/v5"
)

// This function will handle the creation of an event.
func CreateEvent(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	// TODO
}
