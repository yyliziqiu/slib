package sdb

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/yyliziqiu/slib/slog"
)

type Migration struct {
	Db       *gorm.DB
	Once     []schema.Tabler
	Cron     []schema.Tabler
	Interval time.Duration
}

func Migrates(ctx context.Context, migrations func() []Migration) (err error) {
	for _, migration := range migrations() {
		err = Migrate(ctx, migration)
		if err != nil {
			return err
		}
	}
	return nil
}

func Migrate(ctx context.Context, migration Migration) (err error) {
	db := migration.Db.Set("gorm:table_options", "ENGINE=InnoDB")

	err = migrateTables(db, migration.Once)
	if err != nil {
		return fmt.Errorf("migrate once tables failed [%v]", err)
	}

	if len(migration.Cron) == 0 {
		return nil
	}
	err = migrateTables(db, migration.Cron)
	if err != nil {
		return fmt.Errorf("migrate cron tables failed [%v]", err)
	}
	go runMigrateCronTables(ctx, db, migration.Cron, migration.Interval)

	return nil
}

func migrateTables(db *gorm.DB, tables []schema.Tabler) error {
	for _, table := range tables {
		tab := table.TableName()
		err := db.Table(tab).Migrator().AutoMigrate(&table)
		if err != nil {
			return fmt.Errorf("create table %s failed [%v]", tab, err)
		}
		slog.Infof("Create table succeed, table: %s.", tab)
	}
	return nil
}

func runMigrateCronTables(ctx context.Context, db *gorm.DB, tables []schema.Tabler, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			err := migrateTables(db, tables)
			if err != nil {
				slog.Errorf("Migrate cron tables failed, error: %v.", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
