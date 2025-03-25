package sboot

import (
	"context"
	"fmt"

	"github.com/yyliziqiu/slib/sdb"
	"github.com/yyliziqiu/slib/ses"
	"github.com/yyliziqiu/slib/skafka"
	"github.com/yyliziqiu/slib/slog"
	"github.com/yyliziqiu/slib/sredis"
	"github.com/yyliziqiu/slib/sreflect"
)

func BaseInit(config any) InitFunc {
	return func() (err error) {
		// db
		if val, ok := sreflect.FieldValue(config, "Db"); ok {
			c, ok2 := val.(sdb.Config)
			if ok2 && c.Dsn != "" {
				slog.Info("Init database.")
				err = sdb.Init(c)
				if err != nil {
					return fmt.Errorf("init database failed [%v]", err)
				}
			}
		}
		if val, ok := sreflect.FieldValue(config, "DbList"); ok {
			c, ok2 := val.([]sdb.Config)
			if ok2 && len(c) > 0 {
				slog.Info("Init database list.")
				err = sdb.Init(c...)
				if err != nil {
					return fmt.Errorf("init database list failed [%v]", err)
				}
			}
		}

		// es
		if val, ok := sreflect.FieldValue(config, "Es"); ok {
			c, ok2 := val.(ses.Config)
			if ok2 && len(c.Hosts) > 0 {
				slog.Info("Init es.")
				err = ses.Init(c)
				if err != nil {
					return fmt.Errorf("init es failed [%v]", err)
				}
			}
		}
		if val, ok := sreflect.FieldValue(config, "EsList"); ok {
			c, ok2 := val.([]ses.Config)
			if ok2 && len(c) > 0 {
				slog.Info("Init es list.")
				err = ses.Init(c...)
				if err != nil {
					return fmt.Errorf("init es list failed [%v]", err)
				}
			}
		}

		// redis
		if val, ok := sreflect.FieldValue(config, "Redis"); ok {
			c, ok2 := val.(sredis.Config)
			if ok2 && (c.Addr != "" || len(c.Addrs) > 0) {
				slog.Info("Init redis.")
				err = sredis.Init(c)
				if err != nil {
					return fmt.Errorf("init redis failed [%v]", err)
				}
			}
		}
		if val, ok := sreflect.FieldValue(config, "RedisList"); ok {
			c, ok2 := val.([]sredis.Config)
			if ok2 && len(c) > 0 {
				slog.Info("Init redis list.")
				err = sredis.Init(c...)
				if err != nil {
					return fmt.Errorf("init redis list failed [%v]", err)
				}
			}
		}

		// kafka
		if val, ok := sreflect.FieldValue(config, "Kafka"); ok {
			c, ok2 := val.(skafka.Config)
			if ok2 && c.Server.BootstrapServers != "" {
				slog.Info("Init kafka.")
				err = skafka.Init(c)
				if err != nil {
					return fmt.Errorf("init kafka failed [%v]", err)
				}
			}
		}
		if val, ok := sreflect.FieldValue(config, "KafkaList"); ok {
			c, ok2 := val.([]skafka.Config)
			if ok2 && len(c) > 0 {
				slog.Info("Init kafka list.")
				err = skafka.Init(c...)
				if err != nil {
					return fmt.Errorf("init kafka list failed [%v]", err)
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
			ses.Finally()
			sredis.Finally()
			skafka.Finally()
		}()
		return nil
	}
}
