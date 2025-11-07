## 1. 数据结构

### [Configuration](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\config\config.go#L18-L40) 结构体
- **功能**：映射配置文件的主结构体，包含所有应用配置项
- **组成字段**：
    - [Server](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\config\config.go#L19-L21)：服务相关配置
        - `Port string`：服务监听端口
    - [Log](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\config\config.go#L23-L26)：日志相关配置
        - `Path string`：日志文件路径
        - `Level string`：日志级别
    - [MySQL](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\config\config.go#L28-L34)：数据库连接配置
        - `Host string`：数据库主机地址
        - `Port int`：数据库端口号
        - `User string`：数据库用户名
        - `Password string`：数据库密码
        - `Database string`：数据库名称
    - [Auth](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\config\config.go#L36-L39)：认证相关配置
        - `JwtSecret string`：JWT签名密钥
        - `TokenExpiry int`：Token过期时间(秒)

## 2. 全局变量

### [config](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\config\config.go#L13-L13) 变量
- **类型**：`*Configuration`
- **作用**：存储全局配置实例

### [once](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\config\config.go#L14-L14) 变量
- **类型**：`sync.Once`
- **作用**：确保配置只初始化一次，实现单例模式

## 3. 函数接口

### [InitConfig](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\config\config.go#L43-L50) 函数
- **功能**：初始化全局配置（单例模式）
- **参数**：`configPath string` - 配置文件路径
- **行为**：
    - 使用 `sync.Once` 确保只执行一次
    - 调用 [loadConfig](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\config\config.go#L58-L76) 加载配置
    - 失败时记录致命日志并退出程序

### [GetConfig](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\config\config.go#L53-L55) 函数
- **功能**：获取全局配置实例
- **返回值**：`*Configuration` - 全局配置指针
- **特点**：线程安全，支持并发访问

### [loadConfig](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\config\config.go#L58-L76) 函数
- **功能**：加载并解析YAML格式的配置文件
- **参数**：`path string` - 配置文件路径
- **返回值**：`error` - 错误信息
- **处理逻辑**：
    - 打开配置文件
    - 使用 `yaml.Decoder` 解析配置
    - 设置默认值：
        - 默认服务端口为 "8080"
        - 默认MySQL端口为 3306

## 4. 依赖组件

### 外部依赖
- `gopkg.in/yaml.v3`：YAML文件解析库
- 标准库：`fmt`, [log](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\logs\app.log), `os`, `sync`

## 5. 核心特性

### 单例模式实现
- 使用 `sync.Once` 确保配置只初始化一次
- 提供全局访问点 [GetConfig](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\config\config.go#L53-L55)

### 默认值处理
- 自动为关键配置项设置合理默认值
- 提高系统的健壮性

### YAML配置支持
- 支持结构化的YAML配置文件
- 易于维护和扩展配置项