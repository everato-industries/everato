package user

import (
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/jackc/pgx/v5"
)

func VerifyUser(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {}
