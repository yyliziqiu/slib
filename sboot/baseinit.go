package sboot

import (
	"context"
	"fmt"

	"github.com/yyliziqiu/slib/sdb"
	"github.com/yyliziqiu/slib/selastic"
	"github.com/yyliziqiu/slib/skafka"
	"github.com/yyliziqiu/slib/slog"
	"github.com/yyliziqiu/slib/sredis"
)

func BaseInit(config any) InitFunc {
	return func() (err error) {
		// db
		if val, ok := fieldValue(config, "Db"); ok {
			c, ok2 := val.(sdb.Config)
			if ok2 && c.Dsn != "" {
				slog.Info("Init DB.")
				err = sdb.Init(c)
				if err != nil {
					return fmt.Errorf("init DB failed [%v]", err)
				}
			}
		}
		if val, ok := fieldValue(config, "Dbs"); ok {
			c, ok2 := val.([]sdb.Config)
			if ok2 && len(c) > 0 {
				slog.Info("Init DB.")
				err = sdb.Init(c...)
				if err != nil {
					return fmt.Errorf("init DB failed [%v]", err)
				}
			}
		}

		// redis
		if val, ok := fieldValue(config, "Redis"); ok {
			c, ok2 := val.(sredis.Config)
			if ok2 && (c.Addr != "" || len(c.Addrs) > 0) {
				slog.Info("Init redis.")
				err = sredis.Init(c)
				if err != nil {
					return fmt.Errorf("init redis failed [%v]", err)
				}
			}
		}
		if val, ok := fieldValue(config, "Redises"); ok {
			c, ok2 := val.([]sredis.Config)
			if ok2 && len(c) > 0 {
				slog.Info("Init redis.")
				err = sredis.Init(c...)
				if err != nil {
					return fmt.Errorf("init redis failed [%v]", err)
				}
			}
		}

		// elastic
		if val, ok := fieldValue(config, "Elastic"); ok {
			c, ok2 := val.(selastic.Config)
			if ok2 && len(c.Hosts) > 0 {
				slog.Info("Init elastic.")
				err = selastic.Init(c)
				if err != nil {
					return fmt.Errorf("init elastic failed [%v]", err)
				}
			}
		}
		if val, ok := fieldValue(config, "Elastics"); ok {
			c, ok2 := val.([]selastic.Config)
			if ok2 && len(c) > 0 {
				slog.Info("Init elastic.")
				err = selastic.Init(c...)
				if err != nil {
					return fmt.Errorf("init elastic failed [%v]", err)
				}
			}
		}

		// kafka
		if val, ok := fieldValue(config, "Kafka"); ok {
			c, ok2 := val.(skafka.Config)
			if ok2 && c.Server.BootstrapServers != "" {
				slog.Info("Init kafka.")
				err = skafka.Init(c)
				if err != nil {
					return fmt.Errorf("init kafka failed [%v]", err)
				}
			}
		}
		if val, ok := fieldValue(config, "Kafkas"); ok {
			c, ok2 := val.([]skafka.Config)
			if ok2 && len(c) > 0 {
				slog.Info("Init kafka.")
				err = skafka.Init(c...)
				if err != nil {
					return fmt.Errorf("init kafka failed [%v]", err)
				}
			}
		}

		return nil
	}
}

func BaseBoot() BootFunc {
	return func(ctx context.Context) error {
		go func() {
			<-ctx.Done()
			sdb.Finally()
			sredis.Finally()
			skafka.Finally()
			selastic.Finally()
		}()
		return nil
	}
}
