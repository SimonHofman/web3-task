## 1. 中间件函数

### [LoggerMiddleware](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\middleware\loggerMiddleware.go#L11-L28) 函数
- **功能**：HTTP请求日志记录中间件，记录每个请求的基本信息和处理耗时
- **返回值**：`gin.HandlerFunc` - Gin处理器函数
- **主要流程**：
    1. 在请求开始时记录时间戳
    2. 执行后续处理器链
    3. 计算请求处理耗时
    4. 记录请求日志信息

## 2. 日志记录内容

### 请求信息
- `method`：HTTP请求方法（GET、POST等）
- `path`：请求路径
- `status`：HTTP响应状态码
- `latency`：请求处理耗时

## 3. 依赖组件

### 内部模块
- `mission04/pkg/log`：日志系统，提供 [log.Logger](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\log\log.go#L11-L11) 实例

### 外部库
- `github.com/gin-gonic/gin`：Gin Web框架
- `go.uber.org/zap`：高性能日志库
- `time`：时间处理

## 4. 技术实现

### 时间测量
- 使用 `time.Now()` 记录请求开始时间
- 使用 `time.Since(start)` 计算处理耗时

### 日志输出
- 使用 `log.Logger.Info` 输出结构化日志
- 采用Zap的字段化日志格式，便于日志分析和检索

## 5. 设计特点

### 非阻塞设计
- 日志记录在请求处理完成后进行，不影响请求处理性能

### 结构化日志
- 使用结构化字段记录请求信息，便于后续日志分析和监控

### 性能考虑
- 使用Zap日志库保证高性能日志记录
- 仅记录必要信息，避免日志冗余