## 1. 全局变量

### [Logger](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\log\log.go#L11-L11)
- **类型**: `*zap.Logger`
- **功能**: 全局日志记录器实例，供整个应用程序使用
- **初始化**: 通过 [InitLogger()](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\log\log.go#L13-L70) 函数进行初始化

## 2. 初始化函数

### [InitLogger()](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\log\log.go#L13-L70)
- **功能**: 初始化全局日志系统
- **返回值**: `error` - 初始化过程中可能发生的错误（当前实现始终返回nil）
- **主要流程**:
    1. 配置日志文件滚动策略
    2. 设置日志级别
    3. 配置日志编码器
    4. 创建复合核心处理器
    5. 初始化全局 [Logger](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\log\log.go#L11-L11) 实例

## 3. 核心组件

### 日志文件滚动配置
- **实现**: 基于 `lumberjack.Logger`
- **配置项**:
    - `Filename`: 日志文件路径，从 [config.GetConfig().Log.Path](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\config\config.go#L53-L55) 获取
    - `MaxSize`: 单个日志文件最大尺寸(128MB)
    - `MaxBackups`: 最多保留备份数(30个)
    - `MaxAge`: 最多保留天数(7天)
    - `Compress`: 是否启用压缩(true)

### 日志级别配置
- **支持级别**: debug, info, warn, error
- **默认级别**: info
- **实现**: 通过 `zapcore.Level` 设置

### 编码器配置
- **基础配置**: `zap.NewProductionEncoderConfig()`
- **时间格式**: `zapcore.ISO8601TimeEncoder` (ISO8601格式)
- **级别格式**: `zapcore.CapitalLevelEncoder` (大写级别)

### 复合核心处理器
- **实现**: `zapcore.NewTee`
- **组成**:
    1. 文件输出核心: JSON格式日志写入文件
    2. 控制台输出核心: 控制台友好格式输出到标准输出

## 4. 依赖组件

### 内部依赖
- `mission04/internal/config`: 配置管理模块

### 外部依赖
- `github.com/natefinch/lumberjack`: 日志文件滚动管理
- `go.uber.org/zap`: 高性能日志库
- `go.uber.org/zap/zapcore`: Zap核心组件
- `os`: 标准输出操作

## 5. 功能特性

### 多输出目标
- 同时输出到日志文件和控制台
- 文件输出采用JSON格式便于解析
- 控制台输出采用可读性更好的格式

### 调用者信息追踪
- 使用 `zap.AddCaller()` 记录调用位置
- 使用 `zap.AddCallerSkip(1)` 确保调用者信息准确

### 日志级别控制
- 支持动态配置日志级别
- 只记录等于或高于配置级别的日志