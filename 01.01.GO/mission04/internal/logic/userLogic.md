## 1. 数据结构

### [userLogic](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\logic\userLogic.go#L10-L10) 结构体
- **功能**：实现用户相关业务逻辑的结构体
- **特点**：空结构体，仅用于方法绑定

## 2. 全局变量

### [UserLogic](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\logic\userLogic.go#L12-L12) 变量
- **类型**：`*userLogic`
- **功能**：用户业务逻辑的单例实例
- **初始化**：使用 `new(userLogic)` 创建

## 3. 业务逻辑方法

### [Page](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\db\pagination.go#L6-L6) 方法
- **功能**：分页查询用户列表
- **参数**：`req *model.UserPageReq` - 分页查询请求参数
- **返回值**：
    - `*db.PagedResult` - 分页结果
    - `error` - 查询错误
- **特点**：查询时省略密码字段 `password`

### [Register](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\userHandler.go#L26-L37) 方法
- **功能**：用户注册业务逻辑
- **参数**：`req *model.User` - 用户注册信息
- **主要流程**：
    1. 使用 `bcrypt` 对用户密码进行加密
    2. 将加密后的密码存储到 `req.Password`
    3. 将用户信息保存到数据库
- **返回值**：`error` - 注册过程中的错误

### [Login](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\userHandler.go#L40-L53) 方法
- **功能**：用户登录业务逻辑
- **参数**：`req *model.UserLoginReq` - 用户登录请求信息
- **主要流程**：
    1. 根据用户名查询用户信息
    2. 使用 `bcrypt` 验证密码正确性
    3. 调用 [auth.GenerateToken](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\auth\jwt.go#L36-L66) 生成JWT令牌
    4. 返回包含令牌的登录响应
- **返回值**：
    - `*model.UserLoginResp` - 登录响应（包含token）
    - `error` - 登录过程中的错误

## 4. 依赖组件

### 内部模块
- `mission04/internal/model`：数据模型定义
- `mission04/pkg/auth`：认证相关功能
- `mission04/pkg/db`：数据库访问层和分页工具

### 外部库
- `golang.org/x/crypto/bcrypt`：密码加密和验证

## 5. 安全特性

### 密码安全
- 使用 `bcrypt` 算法对用户密码进行加密存储
- 登录时使用 `bcrypt.CompareHashAndPassword` 进行密码验证

### 数据保护
- 分页查询用户列表时使用 `Omit("password")` 排除密码字段
- 防止敏感信息泄露

### 认证集成
- 登录成功后调用 [auth.GenerateToken](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\auth\jwt.go#L36-L66) 生成JWT令牌
- 实现基于token的无状态认证机制