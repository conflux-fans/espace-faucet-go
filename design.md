# Detailed Design
## 主体思路
以gin框架为主体, 启动服务，在端口进行监听，对前端传回来的信息进行约束，从而提高鲁棒性. 

## 结构体设计
|  名称   | 含义  |
|  ----  | ----  |
| address  | 账户地址 |
| tokenType  | 币种 |
| 合约地址  | erc20合约地址 |

## 接口设计
### sendCFX()
发送cfx代币
#### 输入参数
|  名称   | 含义  |
|  ----  | ----  |
| address  | 账户地址 |

#### 输出参数
交易信息

#### 执行流程
1. gin框架对传至后端的data进行bind
2. 根据账户地址对其领取的时间进行判断
3. 根据账户地址调用rest-api进行转账


### sendERC20()
发送erc20代币
#### 输入参数
|  名称   | 含义  |
|  ----  | ----  |
| address  | 账户地址 |
| tokenType  | 币种 |
| 合约地址  | erc20合约地址 |

#### 输出参数
交易成功信息

#### 执行流程
1. gin框架对传至后端的data进行bind
2. 根据账户地址与合约地址对账户领取的时间进行判断
3. 根据账户地址与合约地址调用rest-api来调用对应的合约方法来实现转账


### 时间约束
1. 一个账户一个小时只能获取一次代币

## 验证码验证
生成一个验证码，传回前端，并进行验证