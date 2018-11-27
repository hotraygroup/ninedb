package controller

import (
	"fmt"
	"log"
	"ninedb/engine"
	"ninedb/models"
)

//用户转账
func Transfer(fromID int, toID int, asset string, amount string) error {
	from := &models.User{UID: fromID}
	to := &models.User{UID: toID}

	//先检查转账的两个账号是否存在
	if engine.Get(from) == nil || engine.Get(to) == nil {
		log.Printf("user %d or %d not found", fromID, toID)
		return fmt.Errorf("user %d or %d not found", fromID, toID)
	}

	//先扣钱
	if err := engine.UpdateField(from, asset, "DESC", amount); err != nil {
		log.Printf("user %d asset[%s] DESC failed", fromID, asset)
		return err
	}

	//再发钱
	if err := engine.UpdateField(to, asset, "INC", amount); err != nil {
		log.Printf("user %d asset[%s] INC failed", toID, asset)
		return err
	}

	return nil
}
