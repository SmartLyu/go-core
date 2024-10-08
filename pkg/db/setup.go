package db

import "github.com/yuanJewel/go-core/db/service"

func SetupCmdb() error {
	if err := service.Instance.Setup([]interface{}{
		&Project{},
	}); err != nil {
		return err
	}
	return nil
}
