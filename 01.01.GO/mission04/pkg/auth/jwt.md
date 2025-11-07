## 1. 数据结构

### [Claims](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\auth\jwt.go#L10-L33) 结构体
- **功能**：定义 JWT token 的声明结构
- **组成**：
    - `UserID uint`：用户唯一标识，JSON标签为 `user_id`
    - `Role string`：用户角色，JSON标签为 `role`
    - `jwt.RegisteredClaims`：JWT 标准声明字段集合，包含 RFC 7519 规范定义的预注册声明

### 标准声明字段 (`jwt.RegisteredClaims`)
- 包含以下标准 JWT 声明：
    - `Issuer (iss)`：签发者标识
    - `Subject (sub)`：主题标识
    - `Audience (aud)`：接收者标识
    - `Expiration Time (exp)`：过期时间戳
    - `Not Before (nbf)`：生效时间戳
    - `Issued At (iat)`：签发时间戳
    - `JWT ID (jti)`：唯一标识符

## 2. 函数接口

### [GenerateToken](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\auth\jwt.go#L36-L66) 函数
- **功能**：为指定用户生成 JWT token
- **输入参数**：`user model.User` - 用户信息对象
- **输出**：
    - `string`：签名后的 JWT token 字符串
    - `error`：可能的错误信息
- **主要流程**：
    1. 计算 token 过期时间
    2. 构造包含用户信息和过期时间的声明
    3. 使用 HS256 算法创建并签名 token

### [ParseToken](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\auth\jwt.go#L69-L83) 函数
- **功能**：解析并验证 JWT token
- **输入参数**：`tokenString string` - JWT token 字符串
- **输出**：
    - `*Claims`：解析后的声明信息指针
    - `error`：解析或验证错误
- **验证机制**：使用配置中的 JWT 密钥验证签名

## 3. 依赖组件

### 外部依赖
- `github.com/golang-jwt/jwt/v5`：JWT 实现库
- `mission04/internal/config`：配置管理模块
- `mission04/internal/model`：数据模型定义

### 配置项引用
- [config.GetConfig().Auth.TokenExpiry](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\config\config.go#L53-L55)：token 有效期（秒）
- [config.GetConfig().Auth.JwtSecret](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\config\config.go#L53-L55)：JWT 签名密钥

## 4. 核心流程

### Token 生成流程
1. 获取 token 有效期配置
2. 计算过期时间 = 当前时间 + 有效期
3. 创建 Claims 对象，填充用户 ID、角色和过期时间
4. 使用 HS256 算法和配置密钥生成签名 token

### Token 解析流程
1. 使用配置密钥解析 token 字符串
2. 验证 token 有效性和签名
3. 提取并返回声明信息

这些元素共同构成了完整的 JWT 认证机制，实现了 token 的生成、签名、解析和验证功能。