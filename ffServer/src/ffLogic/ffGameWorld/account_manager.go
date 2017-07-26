package ffGameWorld

import (
	"ffAutoGen/ffError"
	"ffCommon/log/log"
	"ffCommon/uuid"
)

type accountManager struct {
	mapAccountAgent    map[uuid.UUID]uuid.UUID // key: account uuid; value: agent uuid
	mapOnlineAccounts  map[uuid.UUID]*account  // 在线用户. key: agent uuid
	mapOfflineAccounts map[uuid.UUID]*account  // 缓存的不在线用户. key: account uuid
}

// onEventOnline 上线事件
// 执行到此处, 外界已确保没有重复帐号在线!
func (am *accountManager) onEventOnline(account *account) {
	log.RunLogger.Printf("accountManager.onEventOnline: uuidAgent[%x]-uuidAcount[%x]", account.uuidAgent, account.uuidAccount)

	// 记录
	am.mapOnlineAccounts[account.uuidAgent] = account
	am.mapAccountAgent[account.uuidAccount] = account.uuidAgent
}

// onEventOffline 离线事件
func (am *accountManager) onEventOffline(account *account) {
	// 从在线列表内移除
	_, okOnlineAccount := am.mapOnlineAccounts[account.uuidAgent]
	if okOnlineAccount {
		log.RunLogger.Printf("accountManager.onEventOffline: uuidAgent[%x]-uuidAcount[%x] offline", account.uuidAgent, account.uuidAccount)

		delete(am.mapOnlineAccounts, account.uuidAgent)
		delete(am.mapAccountAgent, account.uuidAccount)
	} else {
		// 错误: 在线列表内没有即将离线的account
		log.FatalLogger.Printf("accountManager.onEventOffline: uuidAgent[%x]-uuidAcount[%x] offline, but not in online list", account.uuidAgent, account.uuidAccount)
	}

	// 检查离线列表
	_, okOfflineAccount := am.mapOfflineAccounts[account.uuidAccount]
	if okOfflineAccount {
		// 错误: 离线列表内居然也存在account, 则记录错误, 并且直接回收掉
		log.FatalLogger.Printf("accountManager.onEventOffline: uuidAgent[%x]-uuidAcount[%x] offline mutil times", account.uuidAgent, account.uuidAccount)

		account.Clear()
		poolOfAccount.back(account)

		return
	}

	// 添加到离线列表
	am.mapOfflineAccounts[account.uuidAccount] = account
}

// onEventKickAll 所有人离线
func (am *accountManager) onEventKickAll(kickReason ffError.Error) {
	log.RunLogger.Printf("accountManager.onEventKickAll: kickReason[%v]", kickReason)

	for uuidAgent, account := range am.mapOnlineAccounts {
		log.RunLogger.Printf("accountManager.onEventKickAll: uuidAgent[%x]-uuidAcount[%x] kickReason[%v]", account.uuidAgent, account.uuidAccount, kickReason)

		// 处理踢出
		account.kicked(kickReason, false)

		// 记录到离线列表
		managerOfAccount.mapOfflineAccounts[uuidAgent] = account
	}

	// 重置在线列表
	managerOfAccount.mapOnlineAccounts = make(map[uuid.UUID]*account, worldFrame.DefaultOnlineCount())
	managerOfAccount.mapAccountAgent = make(map[uuid.UUID]uuid.UUID, worldFrame.DefaultOnlineCount())
}

// applyAccount
func (am *accountManager) applyAccount(uuidAgent uuid.UUID, uuidAccount uuid.UUID) *account {
	account, ok := am.mapOfflineAccounts[uuidAccount]
	if !ok {
		account = poolOfAccount.apply()

		log.RunLogger.Printf("accountManager.applyAccount: uuidAgent[%x]-uuidAcount[%x] fresh", uuidAgent, uuidAccount)
	} else {
		delete(am.mapOfflineAccounts, uuidAccount)

		log.RunLogger.Printf("accountManager.applyAccount: uuidAgent[%x]-uuidAcount[%x] relive", uuidAgent, uuidAccount)
	}

	account.uuidAccount, account.uuidAgent = uuidAccount, uuidAgent

	return account
}

func (am *accountManager) init() {
	managerOfAccount.mapOnlineAccounts = make(map[uuid.UUID]*account, worldFrame.DefaultOnlineCount())
	managerOfAccount.mapAccountAgent = make(map[uuid.UUID]uuid.UUID, worldFrame.DefaultOnlineCount())
	managerOfAccount.mapOfflineAccounts = make(map[uuid.UUID]*account, worldFrame.DefaultOnlineCount()*50/100)
}
