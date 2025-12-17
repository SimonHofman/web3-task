| 名称                                                  | 类型   | 可见范围 | 说明                                              |
| ----------------------------------------------------- | ------ | -------- | ------------------------------------------------- |
| Pool                                                  | 结构体 |          | 池信息                                            |
| UnstakeRequest                                        | 结构体 |          | 取消委托信息                                      |
| User                                                  | 结构体 |          | 用户信息                                          |
| initialize(IERC20,uint256,uint256,uint256)            | 函数   | public   | 初始化                                            |
| _authorizeUpgrade(address)                            | 函数   | internal | 授权升级                                          |
| setMetaNode(IERC20)                                   | 函数   | public   | 设置元节点                                        |
| pauseWithdraw                                         | 函数   | public   | 暂停取现                                          |
| unpauseWithdraw                                       | 函数   | public   | 恢复取现                                          |
| pauseClaim                                            | 函数   | public   | 暂停领取                                          |
| unpauseClaim                                          | 函数   | public   | 恢复领取                                          |
| setStartBlock(uint256)                                | 函数   | public   | 设置开始块高                                      |
| setEndBlock(uint256)                                  | 函数   | public   | 设置结束块高                                      |
| setMetaNodePerBlock(uint256)                          | 函数   | public   | 设置元节点每块奖励                                |
| addPool(address,uint256,uint256,uint256,bool)         | 函数   | public   | 添加质押池                                        |
| updatePool(uint256,uint256,uint256)                   | 函数   | public   | 根据pid，更新minDepositAmount和unstakeLockedBlock |
| setPoolWeight(uint256,uint256,bool)                   | 函数   | public   | 根据pid，更新poolWeight                           |
| poolLength                                            | 函数   | external | 获取池长度                                        |
| getMultiplier(uint256,uint256)                        | 函数   | public   | 安全相乘函数                                      |
| pendingMetaNode(uint256,address)                      | 函数   | external | 调用函数pendingMetaNodeByBlockNumber              |
| pendingMetaNodeByBlockNumber(uint256,address,uint256) | 函数   | public   |                                                   |
| stakingBalance(uint256,address)                       | 函数   | external | 查询池中用户抵押的token数量                       |
| withdrawAmount(uint256,address)                       | 函数   | public   | 用户提现                                          |
| updatePool(uint256)                                   | 函数   | public   | 根据pid，更新单个池中的数据                       |
| massUpdatePools                                       | 函数   | public   | 调用updatePool(uint256)，更新所有的池             |
| depositETH()                                          | 函数   | public   | 存入ETH                                           |
| deposit(uint256,uint256)                              | 函数   | public   | 存入ERC20代币                                     |
| unstake(uint256,uint256)                              | 函数   | public   | 根据pid，解除质押amount                           |
| withdraw(uint256)                                     | 函数   | public   | 提现                                              |
| claim(uint256)                                        | 函数   | public   | 根据pid领取奖励                                   |
| _deposit(uint256,uint256)                             | 函数   | internal | 由depositETH()和deposit(uint256,uint256)调用      |
| _safeMetaNodeTransfer(address,uint256)                | 函数   | internal | 安全转账函数                                      |
| _safeETHTransfer(address,uint256)                     | 函数   | internal | 安全转账函数                                      |
|                                                       |        |          |                                                   |
|                                                       |        |          |                                                   |
|                                                       |        |          |                                                   |