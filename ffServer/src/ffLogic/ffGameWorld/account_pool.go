package ffGameWorld

import (
	"ffCommon/pool"

	"fmt"
)

type accountPool struct {
	pool *pool.Pool
}

func (ap *accountPool) apply() *account {
	account, _ := ap.pool.Apply().(*account)
	return account
}

func (ap *accountPool) back(account *account) {
	ap.pool.Back(account)
}

func (ap *accountPool) String() string {
	return ap.pool.String()
}

func (ap *accountPool) init(initCount int) error {
	if initCount < 1 {
		return fmt.Errorf("accountPool.Init: invalid initCount[%v]", initCount)
	}

	ap.pool = pool.New("accountPool", false, newAccount, initCount, 50)
	return nil
}
