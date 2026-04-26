# 中文深度逐行解读（自动生成）

- 仓库: go_shop_second
- 起始引用: lesson60
- 结束引用: lesson61
- 生成时间: 2026-04-06 18:34:19 +08:00
- 解析范围: 仅 *.go，自动排除 *.pb.go 与 *.gen.go
- 解析模式: 祖先链模式（按提交）

## 提交 1: [74d970e] 全链路可观测性建设“

### 文件: internal/common/handler/redis/client.go

~~~go
   1: package redis
   2: 
   3: import (
   4: 	"context"
   5: 	"errors"
   6: 	"time"
   7: 
   8: 	"github.com/ghost-yu/go_shop_second/common/logging"
   9: 	"github.com/redis/go-redis/v9"
  10: 	"github.com/sirupsen/logrus"
  11: )
  12: 
  13: func SetNX(ctx context.Context, client *redis.Client, key, value string, ttl time.Duration) (err error) {
  14: 	now := time.Now()
  15: 	defer func() {
  16: 		l := logrus.WithContext(ctx).WithFields(logrus.Fields{
  17: 			"start":       now,
  18: 			"key":         key,
  19: 			"value":       value,
  20: 			logging.Error: err,
  21: 			logging.Cost:  time.Since(now).Milliseconds(),
  22: 		})
  23: 		if err == nil {
  24: 			l.Info("_redis_setnx_success")
  25: 		} else {
  26: 			l.Warn("_redis_setnx_error")
  27: 		}
  28: 	}()
  29: 	if client == nil {
  30: 		return errors.New("redis client is nil")
  31: 	}
  32: 	_, err = client.SetNX(ctx, key, value, ttl).Result()
  33: 	return err
  34: }
  35: 
  36: func Del(ctx context.Context, client *redis.Client, key string) (err error) {
  37: 	now := time.Now()
  38: 	defer func() {
  39: 		l := logrus.WithContext(ctx).WithFields(logrus.Fields{
  40: 			"start":       now,
  41: 			"key":         key,
  42: 			logging.Error: err,
  43: 			logging.Cost:  time.Since(now).Milliseconds(),
  44: 		})
  45: 		if err == nil {
  46: 			l.Info("_redis_del_success")
  47: 		} else {
  48: 			l.Warn("_redis_del_error")
  49: 		}
  50: 	}()
  51: 	if client == nil {
  52: 		return errors.New("redis client is nil")
  53: 	}
  54: 	_, err = client.Del(ctx, key).Result()
  55: 	return err
  56: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 语法块结束：关闭 import 或参数列表。 |
| 12 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 13 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 14 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 15 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 16 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 17 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 18 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 19 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 22 | 代码块结束：收束当前函数、分支或类型定义。 |
| 23 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 24 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 25 | 代码块结束：收束当前函数、分支或类型定义。 |
| 26 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 27 | 代码块结束：收束当前函数、分支或类型定义。 |
| 28 | 代码块结束：收束当前函数、分支或类型定义。 |
| 29 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 30 | 返回语句：输出当前结果并结束执行路径。 |
| 31 | 代码块结束：收束当前函数、分支或类型定义。 |
| 32 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 33 | 返回语句：输出当前结果并结束执行路径。 |
| 34 | 代码块结束：收束当前函数、分支或类型定义。 |
| 35 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 36 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 37 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 38 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 39 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 40 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 41 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 42 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 43 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 44 | 代码块结束：收束当前函数、分支或类型定义。 |
| 45 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 46 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 47 | 代码块结束：收束当前函数、分支或类型定义。 |
| 48 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 49 | 代码块结束：收束当前函数、分支或类型定义。 |
| 50 | 代码块结束：收束当前函数、分支或类型定义。 |
| 51 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 52 | 返回语句：输出当前结果并结束执行路径。 |
| 53 | 代码块结束：收束当前函数、分支或类型定义。 |
| 54 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 55 | 返回语句：输出当前结果并结束执行路径。 |
| 56 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/logging/grpc.go

~~~go
   1: package logging
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/sirupsen/logrus"
   7: 	"google.golang.org/grpc"
   8: 	"google.golang.org/grpc/metadata"
   9: )
  10: 
  11: func GRPCUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
  12: 	fields := logrus.Fields{
  13: 		Args: req,
  14: 	}
  15: 	defer func() {
  16: 		fields[Response] = resp
  17: 		if err != nil {
  18: 			fields[Error] = err.Error()
  19: 			logf(ctx, logrus.ErrorLevel, fields, "%s", "_grpc_request_out")
  20: 		}
  21: 	}()
  22: 	md, exist := metadata.FromIncomingContext(ctx)
  23: 	if exist {
  24: 		fields["grpc_metadata"] = md
  25: 	}
  26: 
  27: 	logf(ctx, logrus.InfoLevel, fields, "%s", "_grpc_request_in")
  28: 	return handler(ctx, req)
  29: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 语法块结束：关闭 import 或参数列表。 |
| 10 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 11 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 12 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 13 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 14 | 代码块结束：收束当前函数、分支或类型定义。 |
| 15 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 16 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 17 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 18 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 19 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 20 | 代码块结束：收束当前函数、分支或类型定义。 |
| 21 | 代码块结束：收束当前函数、分支或类型定义。 |
| 22 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 23 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 24 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 25 | 代码块结束：收束当前函数、分支或类型定义。 |
| 26 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | 返回语句：输出当前结果并结束执行路径。 |
| 29 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/logging/logrus.go

~~~go
   1: package logging
   2: 
   3: import (
   4: 	"context"
   5: 	"os"
   6: 	"strconv"
   7: 	"time"
   8: 
   9: 	"github.com/ghost-yu/go_shop_second/common/tracing"
  10: 	"github.com/sirupsen/logrus"
  11: 	prefixed "github.com/x-cray/logrus-prefixed-formatter"
  12: )
  13: 
  14: // 要么用logging.Infof, Warnf...
  15: // 或者直接加hook，用 logrus.Infof...
  16: 
  17: func Init() {
  18: 	SetFormatter(logrus.StandardLogger())
  19: 	logrus.SetLevel(logrus.DebugLevel)
  20: 	logrus.AddHook(&traceHook{})
  21: }
  22: 
  23: func SetFormatter(logger *logrus.Logger) {
  24: 	logger.SetFormatter(&logrus.JSONFormatter{
  25: 		TimestampFormat: time.RFC3339,
  26: 		FieldMap: logrus.FieldMap{
  27: 			logrus.FieldKeyLevel: "severity",
  28: 			logrus.FieldKeyTime:  "time",
  29: 			logrus.FieldKeyMsg:   "message",
  30: 		},
  31: 	})
  32: 	if isLocal, _ := strconv.ParseBool(os.Getenv("LOCAL_ENV")); isLocal {
  33: 		logger.SetFormatter(&prefixed.TextFormatter{
  34: 			ForceColors:     true,
  35: 			ForceFormatting: true,
  36: 			TimestampFormat: time.RFC3339,
  37: 		})
  38: 	}
  39: }
  40: 
  41: func logf(ctx context.Context, level logrus.Level, fields logrus.Fields, format string, args ...any) {
  42: 	logrus.WithContext(ctx).WithFields(fields).Logf(level, format, args...)
  43: }
  44: 
  45: func InfofWithCost(ctx context.Context, fields logrus.Fields, start time.Time, format string, args ...any) {
  46: 	fields[Cost] = time.Since(start).Milliseconds()
  47: 	Infof(ctx, fields, format, args...)
  48: }
  49: 
  50: func Infof(ctx context.Context, fields logrus.Fields, format string, args ...any) {
  51: 	logrus.WithContext(ctx).WithFields(fields).Infof(format, args...)
  52: }
  53: 
  54: func Errorf(ctx context.Context, fields logrus.Fields, format string, args ...any) {
  55: 	logrus.WithContext(ctx).WithFields(fields).Errorf(format, args...)
  56: }
  57: 
  58: func Warnf(ctx context.Context, fields logrus.Fields, format string, args ...any) {
  59: 	logrus.WithContext(ctx).WithFields(fields).Warnf(format, args...)
  60: }
  61: 
  62: func Panicf(ctx context.Context, fields logrus.Fields, format string, args ...any) {
  63: 	logrus.WithContext(ctx).WithFields(fields).Panicf(format, args...)
  64: }
  65: 
  66: type traceHook struct{}
  67: 
  68: func (t traceHook) Levels() []logrus.Level {
  69: 	return logrus.AllLevels
  70: }
  71: 
  72: func (t traceHook) Fire(entry *logrus.Entry) error {
  73: 	if entry.Context != nil {
  74: 		entry.Data["trace"] = tracing.TraceID(entry.Context)
  75: 		entry = entry.WithTime(time.Now())
  76: 	}
  77: 	return nil
  78: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 12 | 语法块结束：关闭 import 或参数列表。 |
| 13 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 14 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 15 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 16 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 17 | init 函数：包加载时自动执行，常用于初始化与注册。 |
| 18 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 19 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 代码块结束：收束当前函数、分支或类型定义。 |
| 22 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 23 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 24 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 25 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 26 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 29 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 30 | 代码块结束：收束当前函数、分支或类型定义。 |
| 31 | 代码块结束：收束当前函数、分支或类型定义。 |
| 32 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 33 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 34 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 35 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 36 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 37 | 代码块结束：收束当前函数、分支或类型定义。 |
| 38 | 代码块结束：收束当前函数、分支或类型定义。 |
| 39 | 代码块结束：收束当前函数、分支或类型定义。 |
| 40 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 41 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 42 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 43 | 代码块结束：收束当前函数、分支或类型定义。 |
| 44 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 45 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 46 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 47 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 48 | 代码块结束：收束当前函数、分支或类型定义。 |
| 49 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 50 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 51 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 52 | 代码块结束：收束当前函数、分支或类型定义。 |
| 53 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 54 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 55 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 56 | 代码块结束：收束当前函数、分支或类型定义。 |
| 57 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 58 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 59 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 60 | 代码块结束：收束当前函数、分支或类型定义。 |
| 61 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 62 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 63 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 64 | 代码块结束：收束当前函数、分支或类型定义。 |
| 65 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 66 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 67 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 68 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 69 | 返回语句：输出当前结果并结束执行路径。 |
| 70 | 代码块结束：收束当前函数、分支或类型定义。 |
| 71 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 72 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 73 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 74 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 75 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 76 | 代码块结束：收束当前函数、分支或类型定义。 |
| 77 | 返回语句：输出当前结果并结束执行路径。 |
| 78 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/logging/mysql.go

~~~go
   1: package logging
   2: 
   3: import (
   4: 	"context"
   5: 	"strings"
   6: 	"time"
   7: 
   8: 	"github.com/ghost-yu/go_shop_second/common/util"
   9: 	"github.com/sirupsen/logrus"
  10: )
  11: 
  12: const (
  13: 	Method   = "method"
  14: 	Args     = "args"
  15: 	Cost     = "cost_ms"
  16: 	Response = "response"
  17: 	Error    = "error"
  18: )
  19: 
  20: type ArgFormatter interface {
  21: 	FormatArg() (string, error)
  22: }
  23: 
  24: func WhenMySQL(ctx context.Context, method string, args ...any) (logrus.Fields, func(any, *error)) {
  25: 	fields := logrus.Fields{
  26: 		Method: method,
  27: 		Args:   formatArgs(args),
  28: 	}
  29: 	start := time.Now()
  30: 	return fields, func(resp any, err *error) {
  31: 		level, msg := logrus.InfoLevel, "mysql_success"
  32: 		fields[Cost] = time.Since(start).Milliseconds()
  33: 		fields[Response] = resp
  34: 
  35: 		if err != nil && (*err != nil) {
  36: 			level, msg = logrus.ErrorLevel, "mysql_error"
  37: 			fields[Error] = (*err).Error()
  38: 		}
  39: 
  40: 		logf(ctx, level, fields, "%s", msg)
  41: 	}
  42: }
  43: 
  44: func formatArgs(args []any) string {
  45: 	var item []string
  46: 	for _, arg := range args {
  47: 		item = append(item, formatArg(arg))
  48: 	}
  49: 	return strings.Join(item, "||")
  50: }
  51: 
  52: func formatArg(arg any) string {
  53: 	var (
  54: 		str string
  55: 		err error
  56: 	)
  57: 	defer func() {
  58: 		if err != nil {
  59: 			str = "unsupported type in formatMySQLArg||err=" + err.Error()
  60: 		}
  61: 	}()
  62: 	switch v := arg.(type) {
  63: 	default:
  64: 		str, err = util.MarshalString(v)
  65: 	case ArgFormatter:
  66: 		str, err = v.FormatArg()
  67: 	}
  68: 	return str
  69: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 语法块结束：关闭 import 或参数列表。 |
| 11 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 12 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 13 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 14 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 15 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 16 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 17 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 18 | 语法块结束：关闭 import 或参数列表。 |
| 19 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 20 | 接口定义：声明能力契约，用于解耦与可替换实现。 |
| 21 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 22 | 代码块结束：收束当前函数、分支或类型定义。 |
| 23 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 24 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 25 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 26 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | 代码块结束：收束当前函数、分支或类型定义。 |
| 29 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 30 | 返回语句：输出当前结果并结束执行路径。 |
| 31 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 32 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 33 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 34 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 35 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 36 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 37 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 38 | 代码块结束：收束当前函数、分支或类型定义。 |
| 39 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 40 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 41 | 代码块结束：收束当前函数、分支或类型定义。 |
| 42 | 代码块结束：收束当前函数、分支或类型定义。 |
| 43 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 44 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 45 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 46 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 47 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 48 | 代码块结束：收束当前函数、分支或类型定义。 |
| 49 | 返回语句：输出当前结果并结束执行路径。 |
| 50 | 代码块结束：收束当前函数、分支或类型定义。 |
| 51 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 52 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 53 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 54 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 55 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 56 | 语法块结束：关闭 import 或参数列表。 |
| 57 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 58 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 59 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 60 | 代码块结束：收束当前函数、分支或类型定义。 |
| 61 | 代码块结束：收束当前函数、分支或类型定义。 |
| 62 | 多分支选择：按状态或类型分流执行路径。 |
| 63 | 默认分支：兜底处理，防止未知输入静默失败。 |
| 64 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 65 | 分支标签：定义 switch 的命中条件。 |
| 66 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 67 | 代码块结束：收束当前函数、分支或类型定义。 |
| 68 | 返回语句：输出当前结果并结束执行路径。 |
| 69 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/logging/when.go

~~~go
   1: package logging
   2: 
   3: import (
   4: 	"context"
   5: 	"time"
   6: 
   7: 	"github.com/sirupsen/logrus"
   8: )
   9: 
  10: func WhenCommandExecute(ctx context.Context, commandName string, cmd any, err error) {
  11: 	fields := logrus.Fields{
  12: 		"cmd": cmd,
  13: 	}
  14: 	if err == nil {
  15: 		logf(ctx, logrus.InfoLevel, fields, "%s_command_success", commandName)
  16: 	} else {
  17: 		logf(ctx, logrus.ErrorLevel, fields, "%s_command_failed", commandName)
  18: 	}
  19: }
  20: 
  21: func WhenRequest(ctx context.Context, method string, args ...any) (logrus.Fields, func(any, *error)) {
  22: 	fields := logrus.Fields{
  23: 		Method: method,
  24: 		Args:   formatArgs(args),
  25: 	}
  26: 	start := time.Now()
  27: 	return fields, func(resp any, err *error) {
  28: 		level, msg := logrus.InfoLevel, "_request_success"
  29: 		fields[Cost] = time.Since(start).Milliseconds()
  30: 		fields[Response] = resp
  31: 
  32: 		if err != nil && (*err != nil) {
  33: 			level, msg = logrus.ErrorLevel, "_request_failed"
  34: 			fields[Error] = (*err).Error()
  35: 		}
  36: 
  37: 		logf(ctx, level, fields, "%s", msg)
  38: 	}
  39: }
  40: 
  41: func WhenEventPublish(ctx context.Context, method string, args ...any) (logrus.Fields, func(any, *error)) {
  42: 	fields := logrus.Fields{
  43: 		Method: method,
  44: 		Args:   formatArgs(args),
  45: 	}
  46: 	start := time.Now()
  47: 	return fields, func(resp any, err *error) {
  48: 		level, msg := logrus.InfoLevel, "_mq_publish_success"
  49: 		fields[Cost] = time.Since(start).Milliseconds()
  50: 		fields[Response] = resp
  51: 
  52: 		if err != nil && (*err != nil) {
  53: 			level, msg = logrus.ErrorLevel, "_mq_publish_failed"
  54: 			fields[Error] = (*err).Error()
  55: 		}
  56: 
  57: 		logf(ctx, level, fields, "%s", msg)
  58: 	}
  59: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 语法块结束：关闭 import 或参数列表。 |
| 9 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 10 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 11 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 12 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 13 | 代码块结束：收束当前函数、分支或类型定义。 |
| 14 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 15 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 16 | 代码块结束：收束当前函数、分支或类型定义。 |
| 17 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 18 | 代码块结束：收束当前函数、分支或类型定义。 |
| 19 | 代码块结束：收束当前函数、分支或类型定义。 |
| 20 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 21 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 22 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 23 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 24 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 25 | 代码块结束：收束当前函数、分支或类型定义。 |
| 26 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 27 | 返回语句：输出当前结果并结束执行路径。 |
| 28 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 29 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 30 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 31 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 32 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 33 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 34 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 35 | 代码块结束：收束当前函数、分支或类型定义。 |
| 36 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 代码块结束：收束当前函数、分支或类型定义。 |
| 39 | 代码块结束：收束当前函数、分支或类型定义。 |
| 40 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 41 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 42 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 43 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 44 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 45 | 代码块结束：收束当前函数、分支或类型定义。 |
| 46 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 47 | 返回语句：输出当前结果并结束执行路径。 |
| 48 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 49 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 50 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 51 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 52 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 53 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 54 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 55 | 代码块结束：收束当前函数、分支或类型定义。 |
| 56 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 57 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 58 | 代码块结束：收束当前函数、分支或类型定义。 |
| 59 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/middleware/request.go

~~~go
   1: package middleware
   2: 
   3: import (
   4: 	"bytes"
   5: 	"encoding/json"
   6: 	"io"
   7: 	"time"
   8: 
   9: 	"github.com/ghost-yu/go_shop_second/common/logging"
  10: 	"github.com/gin-gonic/gin"
  11: 	"github.com/sirupsen/logrus"
  12: )
  13: 
  14: func RequestLog(l *logrus.Entry) gin.HandlerFunc {
  15: 	return func(c *gin.Context) {
  16: 		requestIn(c, l)
  17: 		defer requestOut(c, l)
  18: 		c.Next()
  19: 	}
  20: }
  21: 
  22: func requestOut(c *gin.Context, l *logrus.Entry) {
  23: 	response, _ := c.Get("response")
  24: 	start, _ := c.Get("request_start")
  25: 	startTime := start.(time.Time)
  26: 	l.WithContext(c.Request.Context()).WithFields(logrus.Fields{
  27: 		logging.Cost:     time.Since(startTime).Milliseconds(),
  28: 		logging.Response: response,
  29: 	}).Info("_request_out")
  30: }
  31: 
  32: func requestIn(c *gin.Context, l *logrus.Entry) {
  33: 	c.Set("request_start", time.Now())
  34: 	body := c.Request.Body
  35: 	bodyBytes, _ := io.ReadAll(body)
  36: 	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
  37: 	var compactJson bytes.Buffer
  38: 	_ = json.Compact(&compactJson, bodyBytes)
  39: 	l.WithContext(c.Request.Context()).WithFields(logrus.Fields{
  40: 		"start":      time.Now().Unix(),
  41: 		logging.Args: compactJson.String(),
  42: 		"from":       c.RemoteIP(),
  43: 		"uri":        c.Request.RequestURI,
  44: 	}).Info("_request_in")
  45: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 语法块结束：关闭 import 或参数列表。 |
| 13 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 14 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 15 | 返回语句：输出当前结果并结束执行路径。 |
| 16 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 17 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 18 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 19 | 代码块结束：收束当前函数、分支或类型定义。 |
| 20 | 代码块结束：收束当前函数、分支或类型定义。 |
| 21 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 22 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 23 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 24 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 25 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 26 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 29 | 代码块结束：收束当前函数、分支或类型定义。 |
| 30 | 代码块结束：收束当前函数、分支或类型定义。 |
| 31 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 32 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 33 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 34 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 35 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 36 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 39 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 40 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 41 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 42 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 43 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 44 | 代码块结束：收束当前函数、分支或类型定义。 |
| 45 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/server/gprc.go

~~~go
   1: package server
   2: 
   3: import (
   4: 	"net"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/logging"
   7: 	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
   8: 	grpc_tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
   9: 	"github.com/sirupsen/logrus"
  10: 	"github.com/spf13/viper"
  11: 	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
  12: 	"google.golang.org/grpc"
  13: )
  14: 
  15: func init() {
  16: 	logger := logrus.New()
  17: 	logger.SetLevel(logrus.WarnLevel)
  18: 	grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(logger))
  19: }
  20: 
  21: func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
  22: 	addr := viper.Sub(serviceName).GetString("grpc-addr")
  23: 	if addr == "" {
  24: 		// TODO: Warning log
  25: 		addr = viper.GetString("fallback-grpc-addr")
  26: 	}
  27: 	RunGRPCServerOnAddr(addr, registerServer)
  28: }
  29: 
  30: func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
  31: 	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
  32: 	grpcServer := grpc.NewServer(
  33: 		grpc.StatsHandler(otelgrpc.NewServerHandler()),
  34: 		grpc.ChainUnaryInterceptor(
  35: 			grpc_tags.UnaryServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
  36: 			grpc_logrus.UnaryServerInterceptor(logrusEntry),
  37: 			logging.GRPCUnaryInterceptor,
  38: 		),
  39: 		grpc.ChainStreamInterceptor(
  40: 			grpc_tags.StreamServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
  41: 			grpc_logrus.StreamServerInterceptor(logrusEntry),
  42: 		),
  43: 	)
  44: 	registerServer(grpcServer)
  45: 
  46: 	listen, err := net.Listen("tcp", addr)
  47: 	if err != nil {
  48: 		logrus.Panic(err)
  49: 	}
  50: 	logrus.Infof("Starting gRPC server, Listening: %s", addr)
  51: 	if err := grpcServer.Serve(listen); err != nil {
  52: 		logrus.Panic(err)
  53: 	}
  54: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 8 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 语法块结束：关闭 import 或参数列表。 |
| 14 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 15 | init 函数：包加载时自动执行，常用于初始化与注册。 |
| 16 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 17 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 18 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 19 | 代码块结束：收束当前函数、分支或类型定义。 |
| 20 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 21 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 22 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 23 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 24 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 25 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 26 | 代码块结束：收束当前函数、分支或类型定义。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | 代码块结束：收束当前函数、分支或类型定义。 |
| 29 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 30 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 31 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 32 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 33 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 34 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 35 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 36 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 39 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 40 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 41 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 42 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 43 | 语法块结束：关闭 import 或参数列表。 |
| 44 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 45 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 46 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 47 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 48 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 49 | 代码块结束：收束当前函数、分支或类型定义。 |
| 50 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 51 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 52 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 53 | 代码块结束：收束当前函数、分支或类型定义。 |
| 54 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/kitchen/infrastructure/consumer/consumer.go

~~~go
   1: package consumer
   2: 
   3: import (
   4: 	"context"
   5: 	"encoding/json"
   6: 	"fmt"
   7: 	"time"
   8: 
   9: 	"github.com/ghost-yu/go_shop_second/common/broker"
  10: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
  11: 	"github.com/ghost-yu/go_shop_second/common/logging"
  12: 	"github.com/pkg/errors"
  13: 	amqp "github.com/rabbitmq/amqp091-go"
  14: 	"github.com/sirupsen/logrus"
  15: 	"go.opentelemetry.io/otel"
  16: )
  17: 
  18: type OrderService interface {
  19: 	UpdateOrder(ctx context.Context, request *orderpb.Order) error
  20: }
  21: 
  22: type Consumer struct {
  23: 	orderGRPC OrderService
  24: }
  25: 
  26: type Order struct {
  27: 	ID          string
  28: 	CustomerID  string
  29: 	Status      string
  30: 	PaymentLink string
  31: 	Items       []*orderpb.Item
  32: }
  33: 
  34: func NewConsumer(orderGRPC OrderService) *Consumer {
  35: 	return &Consumer{
  36: 		orderGRPC: orderGRPC,
  37: 	}
  38: }
  39: 
  40: func (c *Consumer) Listen(ch *amqp.Channel) {
  41: 	q, err := ch.QueueDeclare("", true, false, true, false, nil)
  42: 	if err != nil {
  43: 		logrus.Fatal(err)
  44: 	}
  45: 	if err = ch.QueueBind(q.Name, "", broker.EventOrderPaid, false, nil); err != nil {
  46: 		logrus.Fatal(err)
  47: 	}
  48: 	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
  49: 	if err != nil {
  50: 		logrus.Warnf("fail to consume: queue=%s, err=%v", q.Name, err)
  51: 	}
  52: 
  53: 	var forever chan struct{}
  54: 	go func() {
  55: 		for msg := range msgs {
  56: 			c.handleMessage(ch, msg, q)
  57: 		}
  58: 	}()
  59: 	<-forever
  60: }
  61: 
  62: func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
  63: 	tr := otel.Tracer("rabbitmq")
  64: 	ctx, span := tr.Start(broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers), fmt.Sprintf("rabbitmq.%s.consume", q.Name))
  65: 	defer span.End()
  66: 
  67: 	var err error
  68: 	defer func() {
  69: 		if err != nil {
  70: 			logging.Warnf(ctx, nil, "consume failed||from=%s||msg=%+v||err=%v", q.Name, msg, err)
  71: 			_ = msg.Nack(false, false)
  72: 		} else {
  73: 			logging.Infof(ctx, nil, "%s", "consume success")
  74: 			_ = msg.Ack(false)
  75: 		}
  76: 	}()
  77: 
  78: 	o := &Order{}
  79: 	if err = json.Unmarshal(msg.Body, o); err != nil {
  80: 		err = errors.Wrap(err, "error unmarshal msg.body into order")
  81: 		return
  82: 	}
  83: 	if o.Status != "paid" {
  84: 		err = errors.New("order not paid, cannot cook")
  85: 		return
  86: 	}
  87: 	cook(ctx, o)
  88: 	span.AddEvent(fmt.Sprintf("order_cook: %v", o))
  89: 	if err = c.orderGRPC.UpdateOrder(ctx, &orderpb.Order{
  90: 		ID:          o.ID,
  91: 		CustomerID:  o.CustomerID,
  92: 		Status:      "ready",
  93: 		Items:       o.Items,
  94: 		PaymentLink: o.PaymentLink,
  95: 	}); err != nil {
  96: 		logging.Errorf(ctx, nil, "error updating order||orderID=%s||err=%v", o.ID, err)
  97: 		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
  98: 			err = errors.Wrapf(err, "retry_error, error handling retry, messageID=%s||err=%v", msg.MessageId, err)
  99: 		}
 100: 		return
 101: 	}
 102: 	span.AddEvent("kitchen.order.finished.updated")
 103: }
 104: 
 105: func cook(ctx context.Context, o *Order) {
 106: 	logrus.WithContext(ctx).Printf("cooking order: %s", o.ID)
 107: 	time.Sleep(5 * time.Second)
 108: 	logrus.WithContext(ctx).Printf("order %s done!", o.ID)
 109: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 14 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 15 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 16 | 语法块结束：关闭 import 或参数列表。 |
| 17 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 18 | 接口定义：声明能力契约，用于解耦与可替换实现。 |
| 19 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 20 | 代码块结束：收束当前函数、分支或类型定义。 |
| 21 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 22 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 23 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 24 | 代码块结束：收束当前函数、分支或类型定义。 |
| 25 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 26 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 29 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 30 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 31 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 32 | 代码块结束：收束当前函数、分支或类型定义。 |
| 33 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 34 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 35 | 返回语句：输出当前结果并结束执行路径。 |
| 36 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 37 | 代码块结束：收束当前函数、分支或类型定义。 |
| 38 | 代码块结束：收束当前函数、分支或类型定义。 |
| 39 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 40 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 41 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 42 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 43 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 44 | 代码块结束：收束当前函数、分支或类型定义。 |
| 45 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 46 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 47 | 代码块结束：收束当前函数、分支或类型定义。 |
| 48 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 49 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 50 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 51 | 代码块结束：收束当前函数、分支或类型定义。 |
| 52 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 53 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 54 | goroutine 启动：引入并发执行，需关注生命周期管理。 |
| 55 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 56 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 57 | 代码块结束：收束当前函数、分支或类型定义。 |
| 58 | 代码块结束：收束当前函数、分支或类型定义。 |
| 59 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 60 | 代码块结束：收束当前函数、分支或类型定义。 |
| 61 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 62 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 63 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 64 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 65 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 66 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 67 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 68 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 69 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 70 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 71 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 72 | 代码块结束：收束当前函数、分支或类型定义。 |
| 73 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 74 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 75 | 代码块结束：收束当前函数、分支或类型定义。 |
| 76 | 代码块结束：收束当前函数、分支或类型定义。 |
| 77 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 78 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 79 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 80 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 81 | 返回语句：输出当前结果并结束执行路径。 |
| 82 | 代码块结束：收束当前函数、分支或类型定义。 |
| 83 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 84 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 85 | 返回语句：输出当前结果并结束执行路径。 |
| 86 | 代码块结束：收束当前函数、分支或类型定义。 |
| 87 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 88 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 89 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 90 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 91 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 92 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 93 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 94 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 95 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 96 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 97 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 98 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 99 | 代码块结束：收束当前函数、分支或类型定义。 |
| 100 | 返回语句：输出当前结果并结束执行路径。 |
| 101 | 代码块结束：收束当前函数、分支或类型定义。 |
| 102 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 103 | 代码块结束：收束当前函数、分支或类型定义。 |
| 104 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 105 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 106 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 107 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 108 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 109 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/adapters/grpc/stock_grpc.go

~~~go
   1: package grpc
   2: 
   3: import (
   4: 	"context"
   5: 	"errors"
   6: 
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   8: 	"github.com/ghost-yu/go_shop_second/common/genproto/stockpb"
   9: 	"github.com/ghost-yu/go_shop_second/common/logging"
  10: )
  11: 
  12: type StockGRPC struct {
  13: 	client stockpb.StockServiceClient
  14: }
  15: 
  16: func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
  17: 	return &StockGRPC{client: client}
  18: }
  19: 
  20: func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (resp *stockpb.CheckIfItemsInStockResponse, err error) {
  21: 	_, dLog := logging.WhenRequest(ctx, "StockGRPC.CheckIfItemsInStock", items)
  22: 	defer dLog(resp, &err)
  23: 
  24: 	if items == nil {
  25: 		return nil, errors.New("grpc items cannot be nil")
  26: 	}
  27: 	return s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
  28: }
  29: 
  30: func (s StockGRPC) GetItems(ctx context.Context, itemIDs []string) (items []*orderpb.Item, err error) {
  31: 	_, dLog := logging.WhenRequest(ctx, "StockGRPC.GetItems", items)
  32: 	defer dLog(items, &err)
  33: 
  34: 	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{ItemIDs: itemIDs})
  35: 	if err != nil {
  36: 		return nil, err
  37: 	}
  38: 	return resp.Items, nil
  39: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 语法块结束：关闭 import 或参数列表。 |
| 11 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 12 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 13 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 14 | 代码块结束：收束当前函数、分支或类型定义。 |
| 15 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 16 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 17 | 返回语句：输出当前结果并结束执行路径。 |
| 18 | 代码块结束：收束当前函数、分支或类型定义。 |
| 19 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 20 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 21 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 22 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 23 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 24 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 25 | 返回语句：输出当前结果并结束执行路径。 |
| 26 | 代码块结束：收束当前函数、分支或类型定义。 |
| 27 | 返回语句：输出当前结果并结束执行路径。 |
| 28 | 代码块结束：收束当前函数、分支或类型定义。 |
| 29 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 30 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 31 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 32 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 33 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 34 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 35 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 36 | 返回语句：输出当前结果并结束执行路径。 |
| 37 | 代码块结束：收束当前函数、分支或类型定义。 |
| 38 | 返回语句：输出当前结果并结束执行路径。 |
| 39 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/adapters/order_inmem_repository.go

~~~go
   1: package adapters
   2: 
   3: import (
   4: 	"context"
   5: 	"strconv"
   6: 	"sync"
   7: 	"time"
   8: 
   9: 	"github.com/ghost-yu/go_shop_second/common/logging"
  10: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
  11: )
  12: 
  13: type MemoryOrderRepository struct {
  14: 	lock  *sync.RWMutex
  15: 	store []*domain.Order
  16: }
  17: 
  18: func NewMemoryOrderRepository() *MemoryOrderRepository {
  19: 	s := []*domain.Order{
  20: 		{
  21: 			ID:          "fake-ID",
  22: 			CustomerID:  "fake-customer-id",
  23: 			Status:      "fake-status",
  24: 			PaymentLink: "fake-payment-link",
  25: 			Items:       nil,
  26: 		},
  27: 	}
  28: 	return &MemoryOrderRepository{
  29: 		lock:  &sync.RWMutex{},
  30: 		store: s,
  31: 	}
  32: }
  33: 
  34: func (m *MemoryOrderRepository) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
  35: 	_, deferLog := logging.WhenRequest(ctx, "MemoryOrderRepository.Create", map[string]any{"order": order})
  36: 	defer deferLog(created, &err)
  37: 
  38: 	m.lock.Lock()
  39: 	defer m.lock.Unlock()
  40: 	newOrder := &domain.Order{
  41: 		ID:          strconv.FormatInt(time.Now().Unix(), 10),
  42: 		CustomerID:  order.CustomerID,
  43: 		Status:      order.Status,
  44: 		PaymentLink: order.PaymentLink,
  45: 		Items:       order.Items,
  46: 	}
  47: 	return newOrder, nil
  48: }
  49: 
  50: func (m *MemoryOrderRepository) Get(ctx context.Context, id, customerID string) (got *domain.Order, err error) {
  51: 	_, deferLog := logging.WhenRequest(ctx, "MemoryOrderRepository.Get", map[string]any{
  52: 		"id":         id,
  53: 		"customerID": customerID,
  54: 	})
  55: 	defer deferLog(got, &err)
  56: 
  57: 	m.lock.RLock()
  58: 	defer m.lock.RUnlock()
  59: 	for _, o := range m.store {
  60: 		if o.ID == id && o.CustomerID == customerID {
  61: 			return o, nil
  62: 		}
  63: 	}
  64: 	return nil, domain.NotFoundError{OrderID: id}
  65: }
  66: 
  67: func (m *MemoryOrderRepository) Update(
  68: 	ctx context.Context,
  69: 	order *domain.Order,
  70: 	updateFn func(context.Context, *domain.Order) (*domain.Order, error),
  71: ) (err error) {
  72: 	_, deferLog := logging.WhenRequest(ctx, "MemoryOrderRepository.Update", map[string]any{
  73: 		"order": order,
  74: 	})
  75: 	defer deferLog(nil, &err)
  76: 
  77: 	m.lock.Lock()
  78: 	defer m.lock.Unlock()
  79: 	found := false
  80: 	for i, o := range m.store {
  81: 		if o.ID == order.ID && o.CustomerID == order.CustomerID {
  82: 			found = true
  83: 			updatedOrder, err := updateFn(ctx, order)
  84: 			if err != nil {
  85: 				return err
  86: 			}
  87: 			m.store[i] = updatedOrder
  88: 		}
  89: 	}
  90: 	if !found {
  91: 		return domain.NotFoundError{OrderID: order.ID}
  92: 	}
  93: 	return nil
  94: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 11 | 语法块结束：关闭 import 或参数列表。 |
| 12 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 13 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 14 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 15 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 16 | 代码块结束：收束当前函数、分支或类型定义。 |
| 17 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 18 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 19 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 22 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 23 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 24 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 25 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 26 | 代码块结束：收束当前函数、分支或类型定义。 |
| 27 | 代码块结束：收束当前函数、分支或类型定义。 |
| 28 | 返回语句：输出当前结果并结束执行路径。 |
| 29 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 30 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 31 | 代码块结束：收束当前函数、分支或类型定义。 |
| 32 | 代码块结束：收束当前函数、分支或类型定义。 |
| 33 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 34 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 35 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 36 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 37 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 38 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 39 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 40 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 41 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 42 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 43 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 44 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 45 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 46 | 代码块结束：收束当前函数、分支或类型定义。 |
| 47 | 返回语句：输出当前结果并结束执行路径。 |
| 48 | 代码块结束：收束当前函数、分支或类型定义。 |
| 49 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 50 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 51 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 52 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 53 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 54 | 代码块结束：收束当前函数、分支或类型定义。 |
| 55 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 56 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 57 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 58 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 59 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 60 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 61 | 返回语句：输出当前结果并结束执行路径。 |
| 62 | 代码块结束：收束当前函数、分支或类型定义。 |
| 63 | 代码块结束：收束当前函数、分支或类型定义。 |
| 64 | 返回语句：输出当前结果并结束执行路径。 |
| 65 | 代码块结束：收束当前函数、分支或类型定义。 |
| 66 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 67 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 68 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 69 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 70 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 71 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 72 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 73 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 74 | 代码块结束：收束当前函数、分支或类型定义。 |
| 75 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 76 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 77 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 78 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 79 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 80 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 81 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 82 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 83 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 84 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 85 | 返回语句：输出当前结果并结束执行路径。 |
| 86 | 代码块结束：收束当前函数、分支或类型定义。 |
| 87 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 88 | 代码块结束：收束当前函数、分支或类型定义。 |
| 89 | 代码块结束：收束当前函数、分支或类型定义。 |
| 90 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 91 | 返回语句：输出当前结果并结束执行路径。 |
| 92 | 代码块结束：收束当前函数、分支或类型定义。 |
| 93 | 返回语句：输出当前结果并结束执行路径。 |
| 94 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/adapters/order_mongo_repository.go

~~~go
   1: package adapters
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	_ "github.com/ghost-yu/go_shop_second/common/config"
   7: 	"github.com/ghost-yu/go_shop_second/common/logging"
   8: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
   9: 	"github.com/ghost-yu/go_shop_second/order/entity"
  10: 	"github.com/spf13/viper"
  11: 	"go.mongodb.org/mongo-driver/bson"
  12: 	"go.mongodb.org/mongo-driver/bson/primitive"
  13: 	"go.mongodb.org/mongo-driver/mongo"
  14: )
  15: 
  16: var (
  17: 	dbName   = viper.GetString("mongo.db-name")
  18: 	collName = viper.GetString("mongo.coll-name")
  19: )
  20: 
  21: type OrderRepositoryMongo struct {
  22: 	db *mongo.Client
  23: }
  24: 
  25: func NewOrderRepositoryMongo(db *mongo.Client) *OrderRepositoryMongo {
  26: 	return &OrderRepositoryMongo{db: db}
  27: }
  28: 
  29: func (r *OrderRepositoryMongo) collection() *mongo.Collection {
  30: 	return r.db.Database(dbName).Collection(collName)
  31: }
  32: 
  33: type orderModel struct {
  34: 	MongoID     primitive.ObjectID `bson:"_id"`
  35: 	ID          string             `bson:"id"`
  36: 	CustomerID  string             `bson:"customer_id"`
  37: 	Status      string             `bson:"status"`
  38: 	PaymentLink string             `bson:"payment_link"`
  39: 	Items       []*entity.Item     `bson:"items"`
  40: }
  41: 
  42: func (r *OrderRepositoryMongo) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
  43: 	_, deferLog := logging.WhenRequest(ctx, "OrderRepositoryMongo.Create", map[string]any{"order": order})
  44: 	defer deferLog(created, &err)
  45: 
  46: 	write := r.marshalToModel(order)
  47: 	res, err := r.collection().InsertOne(ctx, write)
  48: 	if err != nil {
  49: 		return nil, err
  50: 	}
  51: 	created = order
  52: 	created.ID = res.InsertedID.(primitive.ObjectID).Hex()
  53: 	return created, nil
  54: }
  55: 
  56: func (r *OrderRepositoryMongo) Get(ctx context.Context, id, customerID string) (got *domain.Order, err error) {
  57: 	_, deferLog := logging.WhenRequest(ctx, "OrderRepositoryMongo.Get", map[string]any{
  58: 		"id":         id,
  59: 		"customerID": customerID,
  60: 	})
  61: 	defer deferLog(got, &err)
  62: 
  63: 	read := &orderModel{}
  64: 	mongoID, _ := primitive.ObjectIDFromHex(id)
  65: 	cond := bson.M{"_id": mongoID}
  66: 	if err = r.collection().FindOne(ctx, cond).Decode(read); err != nil {
  67: 		return
  68: 	}
  69: 	if read == nil {
  70: 		return nil, domain.NotFoundError{OrderID: id}
  71: 	}
  72: 	got = r.unmarshal(read)
  73: 	return got, nil
  74: }
  75: 
  76: // Update 先查找对应的order，然后apply updateFn，再写入回去
  77: func (r *OrderRepositoryMongo) Update(
  78: 	ctx context.Context,
  79: 	order *domain.Order,
  80: 	updateFn func(context.Context, *domain.Order) (*domain.Order, error),
  81: ) (err error) {
  82: 	_, deferLog := logging.WhenRequest(ctx, "OrderRepositoryMongo.Update", map[string]any{
  83: 		"order": order,
  84: 	})
  85: 	defer deferLog(nil, &err)
  86: 
  87: 	// 事务
  88: 	session, err := r.db.StartSession()
  89: 	if err != nil {
  90: 		return
  91: 	}
  92: 	defer session.EndSession(ctx)
  93: 
  94: 	if err = session.StartTransaction(); err != nil {
  95: 		return err
  96: 	}
  97: 	defer func() {
  98: 		if err == nil {
  99: 			_ = session.CommitTransaction(ctx)
 100: 		} else {
 101: 			_ = session.AbortTransaction(ctx)
 102: 		}
 103: 	}()
 104: 
 105: 	// inside transaction:
 106: 	oldOrder, err := r.Get(ctx, order.ID, order.CustomerID)
 107: 	if err != nil {
 108: 		return
 109: 	}
 110: 	updated, err := updateFn(ctx, order)
 111: 	if err != nil {
 112: 		return
 113: 	}
 114: 	mongoID, _ := primitive.ObjectIDFromHex(oldOrder.ID)
 115: 	_, err = r.collection().UpdateOne(
 116: 		ctx,
 117: 		bson.M{"_id": mongoID, "customer_id": oldOrder.CustomerID},
 118: 		bson.M{"$set": bson.M{
 119: 			"status":       updated.Status,
 120: 			"payment_link": updated.PaymentLink,
 121: 		}},
 122: 	)
 123: 	if err != nil {
 124: 		return
 125: 	}
 126: 	return
 127: }
 128: 
 129: func (r *OrderRepositoryMongo) marshalToModel(order *domain.Order) *orderModel {
 130: 	return &orderModel{
 131: 		MongoID:     primitive.NewObjectID(),
 132: 		ID:          order.ID,
 133: 		CustomerID:  order.CustomerID,
 134: 		Status:      order.Status,
 135: 		PaymentLink: order.PaymentLink,
 136: 		Items:       order.Items,
 137: 	}
 138: }
 139: 
 140: func (r *OrderRepositoryMongo) unmarshal(m *orderModel) *domain.Order {
 141: 	return &domain.Order{
 142: 		ID:          m.MongoID.Hex(),
 143: 		CustomerID:  m.CustomerID,
 144: 		Status:      m.Status,
 145: 		PaymentLink: m.PaymentLink,
 146: 		Items:       m.Items,
 147: 	}
 148: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 语法块结束：关闭 import 或参数列表。 |
| 15 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 16 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 17 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 18 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 19 | 语法块结束：关闭 import 或参数列表。 |
| 20 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 21 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 22 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 23 | 代码块结束：收束当前函数、分支或类型定义。 |
| 24 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 25 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 26 | 返回语句：输出当前结果并结束执行路径。 |
| 27 | 代码块结束：收束当前函数、分支或类型定义。 |
| 28 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 29 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 30 | 返回语句：输出当前结果并结束执行路径。 |
| 31 | 代码块结束：收束当前函数、分支或类型定义。 |
| 32 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 33 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 34 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 35 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 36 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 39 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 40 | 代码块结束：收束当前函数、分支或类型定义。 |
| 41 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 42 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 43 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 44 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 45 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 46 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 47 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 48 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 49 | 返回语句：输出当前结果并结束执行路径。 |
| 50 | 代码块结束：收束当前函数、分支或类型定义。 |
| 51 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 52 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 53 | 返回语句：输出当前结果并结束执行路径。 |
| 54 | 代码块结束：收束当前函数、分支或类型定义。 |
| 55 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 56 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 57 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 58 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 59 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 60 | 代码块结束：收束当前函数、分支或类型定义。 |
| 61 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 62 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 63 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 64 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 65 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 66 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 67 | 返回语句：输出当前结果并结束执行路径。 |
| 68 | 代码块结束：收束当前函数、分支或类型定义。 |
| 69 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 70 | 返回语句：输出当前结果并结束执行路径。 |
| 71 | 代码块结束：收束当前函数、分支或类型定义。 |
| 72 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 73 | 返回语句：输出当前结果并结束执行路径。 |
| 74 | 代码块结束：收束当前函数、分支或类型定义。 |
| 75 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 76 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 77 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 78 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 79 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 80 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 81 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 82 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 83 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 84 | 代码块结束：收束当前函数、分支或类型定义。 |
| 85 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 86 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 87 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 88 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 89 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 90 | 返回语句：输出当前结果并结束执行路径。 |
| 91 | 代码块结束：收束当前函数、分支或类型定义。 |
| 92 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 93 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 94 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 95 | 返回语句：输出当前结果并结束执行路径。 |
| 96 | 代码块结束：收束当前函数、分支或类型定义。 |
| 97 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 98 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 99 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 100 | 代码块结束：收束当前函数、分支或类型定义。 |
| 101 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 102 | 代码块结束：收束当前函数、分支或类型定义。 |
| 103 | 代码块结束：收束当前函数、分支或类型定义。 |
| 104 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 105 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 106 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 107 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 108 | 返回语句：输出当前结果并结束执行路径。 |
| 109 | 代码块结束：收束当前函数、分支或类型定义。 |
| 110 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 111 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 112 | 返回语句：输出当前结果并结束执行路径。 |
| 113 | 代码块结束：收束当前函数、分支或类型定义。 |
| 114 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 115 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 116 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 117 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 118 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 119 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 120 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 121 | 代码块结束：收束当前函数、分支或类型定义。 |
| 122 | 语法块结束：关闭 import 或参数列表。 |
| 123 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 124 | 返回语句：输出当前结果并结束执行路径。 |
| 125 | 代码块结束：收束当前函数、分支或类型定义。 |
| 126 | 返回语句：输出当前结果并结束执行路径。 |
| 127 | 代码块结束：收束当前函数、分支或类型定义。 |
| 128 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 129 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 130 | 返回语句：输出当前结果并结束执行路径。 |
| 131 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 132 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 133 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 134 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 135 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 136 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 137 | 代码块结束：收束当前函数、分支或类型定义。 |
| 138 | 代码块结束：收束当前函数、分支或类型定义。 |
| 139 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 140 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 141 | 返回语句：输出当前结果并结束执行路径。 |
| 142 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 143 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 144 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 145 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 146 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 147 | 代码块结束：收束当前函数、分支或类型定义。 |
| 148 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/app/command/create_order.go

~~~go
   1: package command
   2: 
   3: import (
   4: 	"context"
   5: 	"encoding/json"
   6: 	"fmt"
   7: 
   8: 	"github.com/ghost-yu/go_shop_second/common/broker"
   9: 	"github.com/ghost-yu/go_shop_second/common/decorator"
  10: 	"github.com/ghost-yu/go_shop_second/common/logging"
  11: 	"github.com/ghost-yu/go_shop_second/order/app/query"
  12: 	"github.com/ghost-yu/go_shop_second/order/convertor"
  13: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
  14: 	"github.com/ghost-yu/go_shop_second/order/entity"
  15: 	"github.com/pkg/errors"
  16: 	amqp "github.com/rabbitmq/amqp091-go"
  17: 	"github.com/sirupsen/logrus"
  18: 	"go.opentelemetry.io/otel"
  19: 	"google.golang.org/grpc/status"
  20: )
  21: 
  22: type CreateOrder struct {
  23: 	CustomerID string
  24: 	Items      []*entity.ItemWithQuantity
  25: }
  26: 
  27: type CreateOrderResult struct {
  28: 	OrderID string
  29: }
  30: 
  31: type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]
  32: 
  33: type createOrderHandler struct {
  34: 	orderRepo domain.Repository
  35: 	stockGRPC query.StockService
  36: 	channel   *amqp.Channel
  37: }
  38: 
  39: func NewCreateOrderHandler(
  40: 	orderRepo domain.Repository,
  41: 	stockGRPC query.StockService,
  42: 	channel *amqp.Channel,
  43: 	logger *logrus.Entry,
  44: 	metricClient decorator.MetricsClient,
  45: ) CreateOrderHandler {
  46: 	if orderRepo == nil {
  47: 		panic("nil orderRepo")
  48: 	}
  49: 	if stockGRPC == nil {
  50: 		panic("nil stockGRPC")
  51: 	}
  52: 	if channel == nil {
  53: 		panic("nil channel ")
  54: 	}
  55: 	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
  56: 		createOrderHandler{
  57: 			orderRepo: orderRepo,
  58: 			stockGRPC: stockGRPC,
  59: 			channel:   channel,
  60: 		},
  61: 		logger,
  62: 		metricClient,
  63: 	)
  64: }
  65: 
  66: func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
  67: 	var err error
  68: 	defer logging.WhenCommandExecute(ctx, "CreateOrderHandler", cmd, err)
  69: 
  70: 	q, err := c.channel.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
  71: 	if err != nil {
  72: 		return nil, err
  73: 	}
  74: 
  75: 	t := otel.Tracer("rabbitmq")
  76: 	ctx, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", q.Name))
  77: 	defer span.End()
  78: 
  79: 	validItems, err := c.validate(ctx, cmd.Items)
  80: 	if err != nil {
  81: 		return nil, err
  82: 	}
  83: 	pendingOrder, err := domain.NewPendingOrder(cmd.CustomerID, validItems)
  84: 	if err != nil {
  85: 		return nil, err
  86: 	}
  87: 	o, err := c.orderRepo.Create(ctx, pendingOrder)
  88: 	if err != nil {
  89: 		return nil, err
  90: 	}
  91: 
  92: 	marshalledOrder, err := json.Marshal(o)
  93: 	if err != nil {
  94: 		return nil, err
  95: 	}
  96: 	header := broker.InjectRabbitMQHeaders(ctx)
  97: 	err = c.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
  98: 		ContentType:  "application/json",
  99: 		DeliveryMode: amqp.Persistent,
 100: 		Body:         marshalledOrder,
 101: 		Headers:      header,
 102: 	})
 103: 	if err != nil {
 104: 		return nil, errors.Wrapf(err, "publish event error q.Name=%s", q.Name)
 105: 	}
 106: 
 107: 	return &CreateOrderResult{OrderID: o.ID}, nil
 108: }
 109: 
 110: func (c createOrderHandler) validate(ctx context.Context, items []*entity.ItemWithQuantity) ([]*entity.Item, error) {
 111: 	if len(items) == 0 {
 112: 		return nil, errors.New("must have at least one item")
 113: 	}
 114: 	items = packItems(items)
 115: 	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, convertor.NewItemWithQuantityConvertor().EntitiesToProtos(items))
 116: 	if err != nil {
 117: 		return nil, status.Convert(err).Err()
 118: 	}
 119: 	return convertor.NewItemConvertor().ProtosToEntities(resp.Items), nil
 120: }
 121: 
 122: func packItems(items []*entity.ItemWithQuantity) []*entity.ItemWithQuantity {
 123: 	merged := make(map[string]int32)
 124: 	for _, item := range items {
 125: 		merged[item.ID] += item.Quantity
 126: 	}
 127: 	var res []*entity.ItemWithQuantity
 128: 	for id, quantity := range merged {
 129: 		res = append(res, &entity.ItemWithQuantity{
 130: 			ID:       id,
 131: 			Quantity: quantity,
 132: 		})
 133: 	}
 134: 	return res
 135: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 14 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 15 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 16 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 17 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 18 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 19 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 20 | 语法块结束：关闭 import 或参数列表。 |
| 21 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 22 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 23 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 24 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 25 | 代码块结束：收束当前函数、分支或类型定义。 |
| 26 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 27 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 28 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 29 | 代码块结束：收束当前函数、分支或类型定义。 |
| 30 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 31 | 类型定义：建立语义模型，影响方法与边界设计。 |
| 32 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 33 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 34 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 35 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 36 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 37 | 代码块结束：收束当前函数、分支或类型定义。 |
| 38 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 39 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 40 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 41 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 42 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 43 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 44 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 45 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 46 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 47 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 48 | 代码块结束：收束当前函数、分支或类型定义。 |
| 49 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 50 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 51 | 代码块结束：收束当前函数、分支或类型定义。 |
| 52 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 53 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 54 | 代码块结束：收束当前函数、分支或类型定义。 |
| 55 | 返回语句：输出当前结果并结束执行路径。 |
| 56 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 57 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 58 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 59 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 60 | 代码块结束：收束当前函数、分支或类型定义。 |
| 61 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 62 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 63 | 语法块结束：关闭 import 或参数列表。 |
| 64 | 代码块结束：收束当前函数、分支或类型定义。 |
| 65 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 66 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 67 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 68 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 69 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 70 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 71 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 72 | 返回语句：输出当前结果并结束执行路径。 |
| 73 | 代码块结束：收束当前函数、分支或类型定义。 |
| 74 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 75 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 76 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 77 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 78 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 79 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 80 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 81 | 返回语句：输出当前结果并结束执行路径。 |
| 82 | 代码块结束：收束当前函数、分支或类型定义。 |
| 83 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 84 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 85 | 返回语句：输出当前结果并结束执行路径。 |
| 86 | 代码块结束：收束当前函数、分支或类型定义。 |
| 87 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 88 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 89 | 返回语句：输出当前结果并结束执行路径。 |
| 90 | 代码块结束：收束当前函数、分支或类型定义。 |
| 91 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 92 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 93 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 94 | 返回语句：输出当前结果并结束执行路径。 |
| 95 | 代码块结束：收束当前函数、分支或类型定义。 |
| 96 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 97 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 98 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 99 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 100 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 101 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 102 | 代码块结束：收束当前函数、分支或类型定义。 |
| 103 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 104 | 返回语句：输出当前结果并结束执行路径。 |
| 105 | 代码块结束：收束当前函数、分支或类型定义。 |
| 106 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 107 | 返回语句：输出当前结果并结束执行路径。 |
| 108 | 代码块结束：收束当前函数、分支或类型定义。 |
| 109 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 110 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 111 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 112 | 返回语句：输出当前结果并结束执行路径。 |
| 113 | 代码块结束：收束当前函数、分支或类型定义。 |
| 114 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 115 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 116 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 117 | 返回语句：输出当前结果并结束执行路径。 |
| 118 | 代码块结束：收束当前函数、分支或类型定义。 |
| 119 | 返回语句：输出当前结果并结束执行路径。 |
| 120 | 代码块结束：收束当前函数、分支或类型定义。 |
| 121 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 122 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 123 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 124 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 125 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 126 | 代码块结束：收束当前函数、分支或类型定义。 |
| 127 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 128 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 129 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 130 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 131 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 132 | 代码块结束：收束当前函数、分支或类型定义。 |
| 133 | 代码块结束：收束当前函数、分支或类型定义。 |
| 134 | 返回语句：输出当前结果并结束执行路径。 |
| 135 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/app/command/update_order.go

~~~go
   1: package command
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/decorator"
   7: 	"github.com/ghost-yu/go_shop_second/common/logging"
   8: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
   9: 	"github.com/sirupsen/logrus"
  10: )
  11: 
  12: type UpdateOrder struct {
  13: 	Order    *domain.Order
  14: 	UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)
  15: }
  16: 
  17: type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]
  18: 
  19: type updateOrderHandler struct {
  20: 	orderRepo domain.Repository
  21: 	//stockGRPC
  22: }
  23: 
  24: func NewUpdateOrderHandler(
  25: 	orderRepo domain.Repository,
  26: 	logger *logrus.Entry,
  27: 	metricClient decorator.MetricsClient,
  28: ) UpdateOrderHandler {
  29: 	if orderRepo == nil {
  30: 		panic("nil orderRepo")
  31: 	}
  32: 	return decorator.ApplyCommandDecorators[UpdateOrder, interface{}](
  33: 		updateOrderHandler{orderRepo: orderRepo},
  34: 		logger,
  35: 		metricClient,
  36: 	)
  37: }
  38: 
  39: func (c updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
  40: 	var err error
  41: 	defer logging.WhenCommandExecute(ctx, "UpdateOrderHandler", cmd, err)
  42: 
  43: 	if cmd.UpdateFn == nil {
  44: 		logrus.Panicf("UpdateOrderHandler got nil order, cmd=%+v", cmd)
  45: 	}
  46: 	if err = c.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn); err != nil {
  47: 		return nil, err
  48: 	}
  49: 	return nil, nil
  50: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 语法块结束：关闭 import 或参数列表。 |
| 11 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 12 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 13 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 14 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 15 | 代码块结束：收束当前函数、分支或类型定义。 |
| 16 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 17 | 类型定义：建立语义模型，影响方法与边界设计。 |
| 18 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 19 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 22 | 代码块结束：收束当前函数、分支或类型定义。 |
| 23 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 24 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 25 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 26 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 29 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 30 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 31 | 代码块结束：收束当前函数、分支或类型定义。 |
| 32 | 返回语句：输出当前结果并结束执行路径。 |
| 33 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 34 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 35 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 36 | 语法块结束：关闭 import 或参数列表。 |
| 37 | 代码块结束：收束当前函数、分支或类型定义。 |
| 38 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 39 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 40 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 41 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 42 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 43 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 44 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 45 | 代码块结束：收束当前函数、分支或类型定义。 |
| 46 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 47 | 返回语句：输出当前结果并结束执行路径。 |
| 48 | 代码块结束：收束当前函数、分支或类型定义。 |
| 49 | 返回语句：输出当前结果并结束执行路径。 |
| 50 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/infrastructure/consumer/consumer.go

~~~go
   1: package consumer
   2: 
   3: import (
   4: 	"context"
   5: 	"encoding/json"
   6: 	"fmt"
   7: 
   8: 	"github.com/ghost-yu/go_shop_second/common/broker"
   9: 	"github.com/ghost-yu/go_shop_second/common/logging"
  10: 	"github.com/ghost-yu/go_shop_second/order/app"
  11: 	"github.com/ghost-yu/go_shop_second/order/app/command"
  12: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
  13: 	"github.com/pkg/errors"
  14: 	amqp "github.com/rabbitmq/amqp091-go"
  15: 	"github.com/sirupsen/logrus"
  16: 	"go.opentelemetry.io/otel"
  17: )
  18: 
  19: type Consumer struct {
  20: 	app app.Application
  21: }
  22: 
  23: func NewConsumer(app app.Application) *Consumer {
  24: 	return &Consumer{
  25: 		app: app,
  26: 	}
  27: }
  28: 
  29: func (c *Consumer) Listen(ch *amqp.Channel) {
  30: 	q, err := ch.QueueDeclare(broker.EventOrderPaid, true, false, true, false, nil)
  31: 	if err != nil {
  32: 		logrus.Fatal(err)
  33: 	}
  34: 	err = ch.QueueBind(q.Name, "", broker.EventOrderPaid, false, nil)
  35: 	if err != nil {
  36: 		logrus.Fatal(err)
  37: 	}
  38: 	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
  39: 	if err != nil {
  40: 		logrus.Fatal(err)
  41: 	}
  42: 	var forever chan struct{}
  43: 	go func() {
  44: 		for msg := range msgs {
  45: 			c.handleMessage(ch, msg, q)
  46: 		}
  47: 	}()
  48: 	<-forever
  49: }
  50: 
  51: func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
  52: 	tr := otel.Tracer("rabbitmq")
  53: 	ctx, span := tr.Start(broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers), fmt.Sprintf("rabbitmq.%s.consume", q.Name))
  54: 	defer span.End()
  55: 
  56: 	var err error
  57: 	defer func() {
  58: 		if err != nil {
  59: 			logging.Warnf(ctx, nil, "consume failed||from=%s||msg=%+v||err=%v", q.Name, msg, err)
  60: 			_ = msg.Nack(false, false)
  61: 		} else {
  62: 			logging.Infof(ctx, nil, "%s", "consume success")
  63: 			_ = msg.Ack(false)
  64: 		}
  65: 	}()
  66: 
  67: 	o := &domain.Order{}
  68: 	if err = json.Unmarshal(msg.Body, o); err != nil {
  69: 		err = errors.Wrap(err, "error unmarshal msg.body into domain.order")
  70: 		return
  71: 	}
  72: 	_, err = c.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
  73: 		Order: o,
  74: 		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
  75: 			if err := order.IsPaid(); err != nil {
  76: 				return nil, err
  77: 			}
  78: 			return order, nil
  79: 		},
  80: 	})
  81: 	if err != nil {
  82: 		logging.Errorf(ctx, nil, "error updating order||orderID=%s||err=%v", o.ID, err)
  83: 		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
  84: 			err = errors.Wrapf(err, "retry_error, error handling retry, messageID=%s||err=%v", msg.MessageId, err)
  85: 		}
  86: 		return
  87: 	}
  88: 
  89: 	span.AddEvent("order.updated")
  90: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 15 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 16 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 17 | 语法块结束：关闭 import 或参数列表。 |
| 18 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 19 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 代码块结束：收束当前函数、分支或类型定义。 |
| 22 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 23 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 24 | 返回语句：输出当前结果并结束执行路径。 |
| 25 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 26 | 代码块结束：收束当前函数、分支或类型定义。 |
| 27 | 代码块结束：收束当前函数、分支或类型定义。 |
| 28 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 29 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 30 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 31 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 32 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 33 | 代码块结束：收束当前函数、分支或类型定义。 |
| 34 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 35 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 36 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 37 | 代码块结束：收束当前函数、分支或类型定义。 |
| 38 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 39 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 40 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 41 | 代码块结束：收束当前函数、分支或类型定义。 |
| 42 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 43 | goroutine 启动：引入并发执行，需关注生命周期管理。 |
| 44 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 45 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 46 | 代码块结束：收束当前函数、分支或类型定义。 |
| 47 | 代码块结束：收束当前函数、分支或类型定义。 |
| 48 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 49 | 代码块结束：收束当前函数、分支或类型定义。 |
| 50 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 51 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 52 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 53 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 54 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 55 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 56 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 57 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 58 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 59 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 60 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 61 | 代码块结束：收束当前函数、分支或类型定义。 |
| 62 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 63 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 64 | 代码块结束：收束当前函数、分支或类型定义。 |
| 65 | 代码块结束：收束当前函数、分支或类型定义。 |
| 66 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 67 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 68 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 69 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 70 | 返回语句：输出当前结果并结束执行路径。 |
| 71 | 代码块结束：收束当前函数、分支或类型定义。 |
| 72 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 73 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 74 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 75 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 76 | 返回语句：输出当前结果并结束执行路径。 |
| 77 | 代码块结束：收束当前函数、分支或类型定义。 |
| 78 | 返回语句：输出当前结果并结束执行路径。 |
| 79 | 代码块结束：收束当前函数、分支或类型定义。 |
| 80 | 代码块结束：收束当前函数、分支或类型定义。 |
| 81 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 82 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 83 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 84 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 85 | 代码块结束：收束当前函数、分支或类型定义。 |
| 86 | 返回语句：输出当前结果并结束执行路径。 |
| 87 | 代码块结束：收束当前函数、分支或类型定义。 |
| 88 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 89 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 90 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/ports/grpc.go

~~~go
   1: package ports
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   7: 	"github.com/ghost-yu/go_shop_second/order/app"
   8: 	"github.com/ghost-yu/go_shop_second/order/app/command"
   9: 	"github.com/ghost-yu/go_shop_second/order/app/query"
  10: 	"github.com/ghost-yu/go_shop_second/order/convertor"
  11: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
  12: 	"github.com/golang/protobuf/ptypes/empty"
  13: 	"google.golang.org/grpc/codes"
  14: 	"google.golang.org/grpc/status"
  15: 	"google.golang.org/protobuf/types/known/emptypb"
  16: )
  17: 
  18: type GRPCServer struct {
  19: 	app app.Application
  20: }
  21: 
  22: func NewGRPCServer(app app.Application) *GRPCServer {
  23: 	return &GRPCServer{app: app}
  24: }
  25: 
  26: func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
  27: 	_, err := G.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
  28: 		CustomerID: request.CustomerID,
  29: 		Items:      convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
  30: 	})
  31: 	if err != nil {
  32: 		return nil, status.Error(codes.Internal, err.Error())
  33: 	}
  34: 	return &empty.Empty{}, nil
  35: }
  36: 
  37: func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
  38: 	o, err := G.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
  39: 		CustomerID: request.CustomerID,
  40: 		OrderID:    request.OrderID,
  41: 	})
  42: 	if err != nil {
  43: 		return nil, status.Error(codes.NotFound, err.Error())
  44: 	}
  45: 	return convertor.NewOrderConvertor().EntityToProto(o), nil
  46: }
  47: 
  48: func (G GRPCServer) UpdateOrder(ctx context.Context, request *orderpb.Order) (_ *emptypb.Empty, err error) {
  49: 	order, err := domain.NewOrder(
  50: 		request.ID,
  51: 		request.CustomerID,
  52: 		request.Status,
  53: 		request.PaymentLink,
  54: 		convertor.NewItemConvertor().ProtosToEntities(request.Items))
  55: 	if err != nil {
  56: 		err = status.Error(codes.Internal, err.Error())
  57: 		return nil, err
  58: 	}
  59: 	_, err = G.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
  60: 		Order: order,
  61: 		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
  62: 			return order, nil
  63: 		},
  64: 	})
  65: 	return nil, err
  66: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 15 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 16 | 语法块结束：关闭 import 或参数列表。 |
| 17 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 18 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 19 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 20 | 代码块结束：收束当前函数、分支或类型定义。 |
| 21 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 22 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 23 | 返回语句：输出当前结果并结束执行路径。 |
| 24 | 代码块结束：收束当前函数、分支或类型定义。 |
| 25 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 26 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 27 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 28 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 29 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 30 | 代码块结束：收束当前函数、分支或类型定义。 |
| 31 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 32 | 返回语句：输出当前结果并结束执行路径。 |
| 33 | 代码块结束：收束当前函数、分支或类型定义。 |
| 34 | 返回语句：输出当前结果并结束执行路径。 |
| 35 | 代码块结束：收束当前函数、分支或类型定义。 |
| 36 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 37 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 38 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 39 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 40 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 41 | 代码块结束：收束当前函数、分支或类型定义。 |
| 42 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 43 | 返回语句：输出当前结果并结束执行路径。 |
| 44 | 代码块结束：收束当前函数、分支或类型定义。 |
| 45 | 返回语句：输出当前结果并结束执行路径。 |
| 46 | 代码块结束：收束当前函数、分支或类型定义。 |
| 47 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 48 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 49 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 50 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 51 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 52 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 53 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 54 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 55 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 56 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 57 | 返回语句：输出当前结果并结束执行路径。 |
| 58 | 代码块结束：收束当前函数、分支或类型定义。 |
| 59 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 60 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 61 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 62 | 返回语句：输出当前结果并结束执行路径。 |
| 63 | 代码块结束：收束当前函数、分支或类型定义。 |
| 64 | 代码块结束：收束当前函数、分支或类型定义。 |
| 65 | 返回语句：输出当前结果并结束执行路径。 |
| 66 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/payment/adapters/order_grpc.go

~~~go
   1: package adapters
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   7: 	"github.com/ghost-yu/go_shop_second/common/tracing"
   8: 	"google.golang.org/grpc/status"
   9: )
  10: 
  11: type OrderGRPC struct {
  12: 	client orderpb.OrderServiceClient
  13: }
  14: 
  15: func NewOrderGRPC(client orderpb.OrderServiceClient) *OrderGRPC {
  16: 	return &OrderGRPC{client: client}
  17: }
  18: 
  19: func (o OrderGRPC) UpdateOrder(ctx context.Context, order *orderpb.Order) (err error) {
  20: 	ctx, span := tracing.Start(ctx, "order_grpc.update_order")
  21: 	defer span.End()
  22: 
  23: 	_, err = o.client.UpdateOrder(ctx, order)
  24: 	return status.Convert(err).Err()
  25: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 语法块结束：关闭 import 或参数列表。 |
| 10 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 11 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 12 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 13 | 代码块结束：收束当前函数、分支或类型定义。 |
| 14 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 15 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 16 | 返回语句：输出当前结果并结束执行路径。 |
| 17 | 代码块结束：收束当前函数、分支或类型定义。 |
| 18 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 19 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 20 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 21 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 22 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 23 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 24 | 返回语句：输出当前结果并结束执行路径。 |
| 25 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/payment/app/command/create_payment.go

~~~go
   1: package command
   2: 
   3: import (
   4: 	"context"
   5: 
   6: 	"github.com/ghost-yu/go_shop_second/common/decorator"
   7: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
   8: 	"github.com/ghost-yu/go_shop_second/common/logging"
   9: 	"github.com/ghost-yu/go_shop_second/payment/domain"
  10: 	"github.com/sirupsen/logrus"
  11: )
  12: 
  13: // TODO: ACL 清理
  14: 
  15: type CreatePayment struct {
  16: 	Order *orderpb.Order
  17: }
  18: 
  19: type CreatePaymentHandler decorator.CommandHandler[CreatePayment, string]
  20: 
  21: type createPaymentHandler struct {
  22: 	processor domain.Processor
  23: 	orderGRPC OrderService
  24: }
  25: 
  26: func (c createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) (string, error) {
  27: 	var err error
  28: 	defer logging.WhenCommandExecute(ctx, "CreatePaymentHandler", cmd, err)
  29: 
  30: 	link, err := c.processor.CreatePaymentLink(ctx, cmd.Order)
  31: 	if err != nil {
  32: 		return "", err
  33: 	}
  34: 	newOrder := &orderpb.Order{
  35: 		ID:          cmd.Order.ID,
  36: 		CustomerID:  cmd.Order.CustomerID,
  37: 		Status:      "waiting_for_payment",
  38: 		Items:       cmd.Order.Items,
  39: 		PaymentLink: link,
  40: 	}
  41: 	err = c.orderGRPC.UpdateOrder(ctx, newOrder)
  42: 	return link, err
  43: }
  44: 
  45: func NewCreatePaymentHandler(
  46: 	processor domain.Processor,
  47: 	orderGRPC OrderService,
  48: 	logger *logrus.Entry,
  49: 	metricClient decorator.MetricsClient,
  50: ) CreatePaymentHandler {
  51: 	return decorator.ApplyCommandDecorators[CreatePayment, string](
  52: 		createPaymentHandler{processor: processor, orderGRPC: orderGRPC},
  53: 		logger,
  54: 		metricClient,
  55: 	)
  56: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 语法块结束：关闭 import 或参数列表。 |
| 12 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 13 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 14 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 15 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 16 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 17 | 代码块结束：收束当前函数、分支或类型定义。 |
| 18 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 19 | 类型定义：建立语义模型，影响方法与边界设计。 |
| 20 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 21 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 22 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 23 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 24 | 代码块结束：收束当前函数、分支或类型定义。 |
| 25 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 26 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 29 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 30 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 31 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 32 | 返回语句：输出当前结果并结束执行路径。 |
| 33 | 代码块结束：收束当前函数、分支或类型定义。 |
| 34 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 35 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 36 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 39 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 40 | 代码块结束：收束当前函数、分支或类型定义。 |
| 41 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 42 | 返回语句：输出当前结果并结束执行路径。 |
| 43 | 代码块结束：收束当前函数、分支或类型定义。 |
| 44 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 45 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 46 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 47 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 48 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 49 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 50 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 51 | 返回语句：输出当前结果并结束执行路径。 |
| 52 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 53 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 54 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 55 | 语法块结束：关闭 import 或参数列表。 |
| 56 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/payment/http.go

~~~go
   1: package main
   2: 
   3: import (
   4: 	"encoding/json"
   5: 	"fmt"
   6: 	"io"
   7: 	"net/http"
   8: 
   9: 	"github.com/ghost-yu/go_shop_second/common/broker"
  10: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
  11: 	"github.com/ghost-yu/go_shop_second/common/logging"
  12: 	"github.com/ghost-yu/go_shop_second/payment/domain"
  13: 	"github.com/gin-gonic/gin"
  14: 	"github.com/pkg/errors"
  15: 	amqp "github.com/rabbitmq/amqp091-go"
  16: 	"github.com/sirupsen/logrus"
  17: 	"github.com/spf13/viper"
  18: 	"github.com/stripe/stripe-go/v79"
  19: 	"github.com/stripe/stripe-go/v79/webhook"
  20: 	"go.opentelemetry.io/otel"
  21: )
  22: 
  23: type PaymentHandler struct {
  24: 	channel *amqp.Channel
  25: }
  26: 
  27: func NewPaymentHandler(ch *amqp.Channel) *PaymentHandler {
  28: 	return &PaymentHandler{channel: ch}
  29: }
  30: 
  31: // stripe listen --forward-to localhost:8284/api/webhook
  32: func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
  33: 	c.POST("/api/webhook", h.handleWebhook)
  34: }
  35: 
  36: func (h *PaymentHandler) handleWebhook(c *gin.Context) {
  37: 	logrus.WithContext(c.Request.Context()).Info("receive webhook from stripe")
  38: 	var err error
  39: 	defer func() {
  40: 		if err != nil {
  41: 			logging.Warnf(c.Request.Context(), nil, "handleWebhook err=%v", err)
  42: 		} else {
  43: 			logging.Infof(c.Request.Context(), nil, "%s", "handleWebhook success")
  44: 		}
  45: 	}()
  46: 
  47: 	const MaxBodyBytes = int64(65536)
  48: 	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
  49: 	payload, err := io.ReadAll(c.Request.Body)
  50: 	if err != nil {
  51: 		err = errors.Wrap(err, "Error reading request body")
  52: 		c.JSON(http.StatusServiceUnavailable, err.Error())
  53: 		return
  54: 	}
  55: 
  56: 	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"),
  57: 		viper.GetString("ENDPOINT_STRIPE_SECRET"))
  58: 
  59: 	if err != nil {
  60: 		err = errors.Wrap(err, "error verifying webhook signature")
  61: 		c.JSON(http.StatusBadRequest, err.Error())
  62: 		return
  63: 	}
  64: 
  65: 	switch event.Type {
  66: 	case stripe.EventTypeCheckoutSessionCompleted:
  67: 		var session stripe.CheckoutSession
  68: 		if err = json.Unmarshal(event.Data.Raw, &session); err != nil {
  69: 			err = errors.Wrap(err, "error unmarshal event.data.raw into session")
  70: 			c.JSON(http.StatusBadRequest, err.Error())
  71: 			return
  72: 		}
  73: 
  74: 		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
  75: 			var items []*orderpb.Item
  76: 			_ = json.Unmarshal([]byte(session.Metadata["items"]), &items)
  77: 
  78: 			marshalledOrder, err := json.Marshal(&domain.Order{
  79: 				ID:          session.Metadata["orderID"],
  80: 				CustomerID:  session.Metadata["customerID"],
  81: 				Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
  82: 				PaymentLink: session.Metadata["paymentLink"],
  83: 				Items:       items,
  84: 			})
  85: 			if err != nil {
  86: 				err = errors.Wrap(err, "error marshal domain.order")
  87: 				c.JSON(http.StatusBadRequest, err.Error())
  88: 				return
  89: 			}
  90: 
  91: 			// TODO: mq logging
  92: 			tr := otel.Tracer("rabbitmq")
  93: 			ctx, span := tr.Start(c.Request.Context(), fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderPaid))
  94: 			defer span.End()
  95: 
  96: 			headers := broker.InjectRabbitMQHeaders(ctx)
  97: 			_ = h.channel.PublishWithContext(ctx, broker.EventOrderPaid, "", false, false, amqp.Publishing{
  98: 				ContentType:  "application/json",
  99: 				DeliveryMode: amqp.Persistent,
 100: 				Body:         marshalledOrder,
 101: 				Headers:      headers,
 102: 			})
 103: 			logrus.WithContext(c).Infof("message published to %s, body: %s", broker.EventOrderPaid, string(marshalledOrder))
 104: 		}
 105: 	}
 106: 	c.JSON(http.StatusOK, nil)
 107: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 15 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 16 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 17 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 18 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 19 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 20 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 21 | 语法块结束：关闭 import 或参数列表。 |
| 22 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 23 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 24 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 25 | 代码块结束：收束当前函数、分支或类型定义。 |
| 26 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 27 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 28 | 返回语句：输出当前结果并结束执行路径。 |
| 29 | 代码块结束：收束当前函数、分支或类型定义。 |
| 30 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 31 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 32 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 33 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 34 | 代码块结束：收束当前函数、分支或类型定义。 |
| 35 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 36 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 39 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 40 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 41 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 42 | 代码块结束：收束当前函数、分支或类型定义。 |
| 43 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 44 | 代码块结束：收束当前函数、分支或类型定义。 |
| 45 | 代码块结束：收束当前函数、分支或类型定义。 |
| 46 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 47 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 48 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 49 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 50 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 51 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 52 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 53 | 返回语句：输出当前结果并结束执行路径。 |
| 54 | 代码块结束：收束当前函数、分支或类型定义。 |
| 55 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 56 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 57 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 58 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 59 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 60 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 61 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 62 | 返回语句：输出当前结果并结束执行路径。 |
| 63 | 代码块结束：收束当前函数、分支或类型定义。 |
| 64 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 65 | 多分支选择：按状态或类型分流执行路径。 |
| 66 | 分支标签：定义 switch 的命中条件。 |
| 67 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 68 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 69 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 70 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 71 | 返回语句：输出当前结果并结束执行路径。 |
| 72 | 代码块结束：收束当前函数、分支或类型定义。 |
| 73 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 74 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 75 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 76 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 77 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 78 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 79 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 80 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 81 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 82 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 83 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 84 | 代码块结束：收束当前函数、分支或类型定义。 |
| 85 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 86 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 87 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 88 | 返回语句：输出当前结果并结束执行路径。 |
| 89 | 代码块结束：收束当前函数、分支或类型定义。 |
| 90 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 91 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 92 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 93 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 94 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 95 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 96 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 97 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 98 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 99 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 100 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 101 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 102 | 代码块结束：收束当前函数、分支或类型定义。 |
| 103 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 104 | 代码块结束：收束当前函数、分支或类型定义。 |
| 105 | 代码块结束：收束当前函数、分支或类型定义。 |
| 106 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 107 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/payment/infrastructure/consumer/consumer.go

~~~go
   1: package consumer
   2: 
   3: import (
   4: 	"context"
   5: 	"encoding/json"
   6: 	"fmt"
   7: 
   8: 	"github.com/ghost-yu/go_shop_second/common/broker"
   9: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
  10: 	"github.com/ghost-yu/go_shop_second/common/logging"
  11: 	"github.com/ghost-yu/go_shop_second/payment/app"
  12: 	"github.com/ghost-yu/go_shop_second/payment/app/command"
  13: 	"github.com/pkg/errors"
  14: 	amqp "github.com/rabbitmq/amqp091-go"
  15: 	"github.com/sirupsen/logrus"
  16: 	"go.opentelemetry.io/otel"
  17: )
  18: 
  19: type Consumer struct {
  20: 	app app.Application
  21: }
  22: 
  23: func NewConsumer(app app.Application) *Consumer {
  24: 	return &Consumer{
  25: 		app: app,
  26: 	}
  27: }
  28: 
  29: func (c *Consumer) Listen(ch *amqp.Channel) {
  30: 	q, err := ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
  31: 	if err != nil {
  32: 		logrus.Fatal(err)
  33: 	}
  34: 
  35: 	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
  36: 	if err != nil {
  37: 		logrus.Warnf("fail to consume: queue=%s, err=%v", q.Name, err)
  38: 	}
  39: 
  40: 	var forever chan struct{}
  41: 	go func() {
  42: 		for msg := range msgs {
  43: 			c.handleMessage(ch, msg, q)
  44: 		}
  45: 	}()
  46: 	<-forever
  47: }
  48: 
  49: func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
  50: 	tr := otel.Tracer("rabbitmq")
  51: 	ctx, span := tr.Start(broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers), fmt.Sprintf("rabbitmq.%s.consume", q.Name))
  52: 	defer span.End()
  53: 
  54: 	logging.Infof(ctx, nil, "Payment receive a message from %s, msg=%v", q.Name, string(msg.Body))
  55: 	var err error
  56: 	defer func() {
  57: 		if err != nil {
  58: 			logging.Warnf(ctx, nil, "consume failed||from=%s||msg=%+v||err=%v", q.Name, msg, err)
  59: 			_ = msg.Nack(false, false)
  60: 		} else {
  61: 			logging.Infof(ctx, nil, "%s", "consume success")
  62: 			_ = msg.Ack(false)
  63: 		}
  64: 	}()
  65: 
  66: 	o := &orderpb.Order{}
  67: 	if err = json.Unmarshal(msg.Body, o); err != nil {
  68: 		err = errors.Wrap(err, "failed to unmarshall msg to order")
  69: 		return
  70: 	}
  71: 	if _, err = c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
  72: 		err = errors.Wrap(err, "failed to create payment")
  73: 		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
  74: 			err = errors.Wrapf(err, "retry_error, error handling retry, messageID=%s, err=%v", msg.MessageId, err)
  75: 		}
  76: 		return
  77: 	}
  78: 
  79: 	span.AddEvent("payment.created")
  80: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 15 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 16 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 17 | 语法块结束：关闭 import 或参数列表。 |
| 18 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 19 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 代码块结束：收束当前函数、分支或类型定义。 |
| 22 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 23 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 24 | 返回语句：输出当前结果并结束执行路径。 |
| 25 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 26 | 代码块结束：收束当前函数、分支或类型定义。 |
| 27 | 代码块结束：收束当前函数、分支或类型定义。 |
| 28 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 29 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 30 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 31 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 32 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 33 | 代码块结束：收束当前函数、分支或类型定义。 |
| 34 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 35 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 36 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 37 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 38 | 代码块结束：收束当前函数、分支或类型定义。 |
| 39 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 40 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 41 | goroutine 启动：引入并发执行，需关注生命周期管理。 |
| 42 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 43 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 44 | 代码块结束：收束当前函数、分支或类型定义。 |
| 45 | 代码块结束：收束当前函数、分支或类型定义。 |
| 46 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 47 | 代码块结束：收束当前函数、分支或类型定义。 |
| 48 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 49 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 50 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 51 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 52 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 53 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 54 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 55 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 56 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 57 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 58 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 59 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 60 | 代码块结束：收束当前函数、分支或类型定义。 |
| 61 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 62 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 63 | 代码块结束：收束当前函数、分支或类型定义。 |
| 64 | 代码块结束：收束当前函数、分支或类型定义。 |
| 65 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 66 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 67 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 68 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 69 | 返回语句：输出当前结果并结束执行路径。 |
| 70 | 代码块结束：收束当前函数、分支或类型定义。 |
| 71 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 72 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 73 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 74 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 75 | 代码块结束：收束当前函数、分支或类型定义。 |
| 76 | 返回语句：输出当前结果并结束执行路径。 |
| 77 | 代码块结束：收束当前函数、分支或类型定义。 |
| 78 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 79 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 80 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/stock/app/query/check_if_items_in_stock.go

~~~go
   1: package query
   2: 
   3: import (
   4: 	"context"
   5: 	"strings"
   6: 	"time"
   7: 
   8: 	"github.com/ghost-yu/go_shop_second/common/decorator"
   9: 	"github.com/ghost-yu/go_shop_second/common/handler/redis"
  10: 	"github.com/ghost-yu/go_shop_second/common/logging"
  11: 	domain "github.com/ghost-yu/go_shop_second/stock/domain/stock"
  12: 	"github.com/ghost-yu/go_shop_second/stock/entity"
  13: 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/integration"
  14: 	"github.com/pkg/errors"
  15: 	"github.com/sirupsen/logrus"
  16: )
  17: 
  18: const (
  19: 	redisLockPrefix = "check_stock_"
  20: )
  21: 
  22: type CheckIfItemsInStock struct {
  23: 	Items []*entity.ItemWithQuantity
  24: }
  25: 
  26: type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*entity.Item]
  27: 
  28: type checkIfItemsInStockHandler struct {
  29: 	stockRepo domain.Repository
  30: 	stripeAPI *integration.StripeAPI
  31: }
  32: 
  33: func NewCheckIfItemsInStockHandler(
  34: 	stockRepo domain.Repository,
  35: 	stripeAPI *integration.StripeAPI,
  36: 	logger *logrus.Entry,
  37: 	metricClient decorator.MetricsClient,
  38: ) CheckIfItemsInStockHandler {
  39: 	if stockRepo == nil {
  40: 		panic("nil stockRepo")
  41: 	}
  42: 	if stripeAPI == nil {
  43: 		panic("nil stripeAPI")
  44: 	}
  45: 	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*entity.Item](
  46: 		checkIfItemsInStockHandler{
  47: 			stockRepo: stockRepo,
  48: 			stripeAPI: stripeAPI,
  49: 		},
  50: 		logger,
  51: 		metricClient,
  52: 	)
  53: }
  54: 
  55: // Deprecated
  56: var stub = map[string]string{
  57: 	"1": "price_1QBYvXRuyMJmUCSsEyQm2oP7",
  58: 	"2": "price_1QBYl4RuyMJmUCSsWt2tgh6d",
  59: }
  60: 
  61: func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
  62: 	if err := lock(ctx, getLockKey(query)); err != nil {
  63: 		return nil, errors.Wrapf(err, "redis lock error: key=%s", getLockKey(query))
  64: 	}
  65: 	defer func() {
  66: 		if err := unlock(ctx, getLockKey(query)); err != nil {
  67: 			logging.Warnf(ctx, nil, "redis unlock fail, err=%v", err)
  68: 		}
  69: 	}()
  70: 
  71: 	var res []*entity.Item
  72: 	for _, i := range query.Items {
  73: 		priceID, err := h.stripeAPI.GetPriceByProductID(ctx, i.ID)
  74: 		if err != nil || priceID == "" {
  75: 			return nil, err
  76: 		}
  77: 		res = append(res, &entity.Item{
  78: 			ID:       i.ID,
  79: 			Quantity: i.Quantity,
  80: 			PriceID:  priceID,
  81: 		})
  82: 	}
  83: 	if err := h.checkStock(ctx, query.Items); err != nil {
  84: 		return nil, err
  85: 	}
  86: 	return res, nil
  87: }
  88: 
  89: func getLockKey(query CheckIfItemsInStock) string {
  90: 	var ids []string
  91: 	for _, i := range query.Items {
  92: 		ids = append(ids, i.ID)
  93: 	}
  94: 	return redisLockPrefix + strings.Join(ids, "_")
  95: }
  96: 
  97: func unlock(ctx context.Context, key string) error {
  98: 	return redis.Del(ctx, redis.LocalClient(), key)
  99: }
 100: 
 101: func lock(ctx context.Context, key string) error {
 102: 	return redis.SetNX(ctx, redis.LocalClient(), key, "1", 5*time.Minute)
 103: }
 104: 
 105: func (h checkIfItemsInStockHandler) checkStock(ctx context.Context, query []*entity.ItemWithQuantity) error {
 106: 	var ids []string
 107: 	for _, i := range query {
 108: 		ids = append(ids, i.ID)
 109: 	}
 110: 	records, err := h.stockRepo.GetStock(ctx, ids)
 111: 	if err != nil {
 112: 		return err
 113: 	}
 114: 	idQuantityMap := make(map[string]int32)
 115: 	for _, r := range records {
 116: 		idQuantityMap[r.ID] += r.Quantity
 117: 	}
 118: 	var (
 119: 		ok       = true
 120: 		failedOn []struct {
 121: 			ID   string
 122: 			Want int32
 123: 			Have int32
 124: 		}
 125: 	)
 126: 	for _, item := range query {
 127: 		if item.Quantity > idQuantityMap[item.ID] {
 128: 			ok = false
 129: 			failedOn = append(failedOn, struct {
 130: 				ID   string
 131: 				Want int32
 132: 				Have int32
 133: 			}{ID: item.ID, Want: item.Quantity, Have: idQuantityMap[item.ID]})
 134: 		}
 135: 	}
 136: 	if ok {
 137: 		return h.stockRepo.UpdateStock(ctx, query, func(
 138: 			ctx context.Context,
 139: 			existing []*entity.ItemWithQuantity,
 140: 			query []*entity.ItemWithQuantity,
 141: 		) ([]*entity.ItemWithQuantity, error) {
 142: 			var newItems []*entity.ItemWithQuantity
 143: 			for _, e := range existing {
 144: 				for _, q := range query {
 145: 					if e.ID == q.ID {
 146: 						newItems = append(newItems, &entity.ItemWithQuantity{
 147: 							ID:       e.ID,
 148: 							Quantity: e.Quantity - q.Quantity,
 149: 						})
 150: 					}
 151: 				}
 152: 			}
 153: 			return newItems, nil
 154: 		})
 155: 	}
 156: 	return domain.ExceedStockError{FailedOn: failedOn}
 157: }
 158: 
 159: func getStubPriceID(id string) string {
 160: 	priceID, ok := stub[id]
 161: 	if !ok {
 162: 		priceID = stub["1"]
 163: 	}
 164: 	return priceID
 165: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 15 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 16 | 语法块结束：关闭 import 或参数列表。 |
| 17 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 18 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 19 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 20 | 语法块结束：关闭 import 或参数列表。 |
| 21 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 22 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 23 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 24 | 代码块结束：收束当前函数、分支或类型定义。 |
| 25 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 26 | 类型定义：建立语义模型，影响方法与边界设计。 |
| 27 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 28 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 29 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 30 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 31 | 代码块结束：收束当前函数、分支或类型定义。 |
| 32 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 33 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 34 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 35 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 36 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 39 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 40 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 41 | 代码块结束：收束当前函数、分支或类型定义。 |
| 42 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 43 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 44 | 代码块结束：收束当前函数、分支或类型定义。 |
| 45 | 返回语句：输出当前结果并结束执行路径。 |
| 46 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 47 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 48 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 49 | 代码块结束：收束当前函数、分支或类型定义。 |
| 50 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 51 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 52 | 语法块结束：关闭 import 或参数列表。 |
| 53 | 代码块结束：收束当前函数、分支或类型定义。 |
| 54 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 55 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 56 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 57 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 58 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 59 | 代码块结束：收束当前函数、分支或类型定义。 |
| 60 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 61 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 62 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 63 | 返回语句：输出当前结果并结束执行路径。 |
| 64 | 代码块结束：收束当前函数、分支或类型定义。 |
| 65 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 66 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 67 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 68 | 代码块结束：收束当前函数、分支或类型定义。 |
| 69 | 代码块结束：收束当前函数、分支或类型定义。 |
| 70 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 71 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 72 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 73 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 74 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 75 | 返回语句：输出当前结果并结束执行路径。 |
| 76 | 代码块结束：收束当前函数、分支或类型定义。 |
| 77 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 78 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 79 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 80 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 81 | 代码块结束：收束当前函数、分支或类型定义。 |
| 82 | 代码块结束：收束当前函数、分支或类型定义。 |
| 83 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 84 | 返回语句：输出当前结果并结束执行路径。 |
| 85 | 代码块结束：收束当前函数、分支或类型定义。 |
| 86 | 返回语句：输出当前结果并结束执行路径。 |
| 87 | 代码块结束：收束当前函数、分支或类型定义。 |
| 88 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 89 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 90 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 91 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 92 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 93 | 代码块结束：收束当前函数、分支或类型定义。 |
| 94 | 返回语句：输出当前结果并结束执行路径。 |
| 95 | 代码块结束：收束当前函数、分支或类型定义。 |
| 96 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 97 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 98 | 返回语句：输出当前结果并结束执行路径。 |
| 99 | 代码块结束：收束当前函数、分支或类型定义。 |
| 100 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 101 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 102 | 返回语句：输出当前结果并结束执行路径。 |
| 103 | 代码块结束：收束当前函数、分支或类型定义。 |
| 104 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 105 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 106 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 107 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 108 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 109 | 代码块结束：收束当前函数、分支或类型定义。 |
| 110 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 111 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 112 | 返回语句：输出当前结果并结束执行路径。 |
| 113 | 代码块结束：收束当前函数、分支或类型定义。 |
| 114 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 115 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 116 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 117 | 代码块结束：收束当前函数、分支或类型定义。 |
| 118 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 119 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 120 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 121 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 122 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 123 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 124 | 代码块结束：收束当前函数、分支或类型定义。 |
| 125 | 语法块结束：关闭 import 或参数列表。 |
| 126 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 127 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 128 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 129 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 130 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 131 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 132 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 133 | 代码块结束：收束当前函数、分支或类型定义。 |
| 134 | 代码块结束：收束当前函数、分支或类型定义。 |
| 135 | 代码块结束：收束当前函数、分支或类型定义。 |
| 136 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 137 | 返回语句：输出当前结果并结束执行路径。 |
| 138 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 139 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 140 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 141 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 142 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 143 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 144 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 145 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 146 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 147 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 148 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 149 | 代码块结束：收束当前函数、分支或类型定义。 |
| 150 | 代码块结束：收束当前函数、分支或类型定义。 |
| 151 | 代码块结束：收束当前函数、分支或类型定义。 |
| 152 | 代码块结束：收束当前函数、分支或类型定义。 |
| 153 | 返回语句：输出当前结果并结束执行路径。 |
| 154 | 代码块结束：收束当前函数、分支或类型定义。 |
| 155 | 代码块结束：收束当前函数、分支或类型定义。 |
| 156 | 返回语句：输出当前结果并结束执行路径。 |
| 157 | 代码块结束：收束当前函数、分支或类型定义。 |
| 158 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 159 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 160 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 161 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 162 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 163 | 代码块结束：收束当前函数、分支或类型定义。 |
| 164 | 返回语句：输出当前结果并结束执行路径。 |
| 165 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/stock/infrastructure/persistent/mysql.go

~~~go
   1: package persistent
   2: 
   3: import (
   4: 	"context"
   5: 	"fmt"
   6: 	"time"
   7: 
   8: 	"github.com/ghost-yu/go_shop_second/common/logging"
   9: 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent/builder"
  10: 	"github.com/sirupsen/logrus"
  11: 	"github.com/spf13/viper"
  12: 	"gorm.io/driver/mysql"
  13: 	"gorm.io/gorm"
  14: 	"gorm.io/gorm/clause"
  15: )
  16: 
  17: type MySQL struct {
  18: 	db *gorm.DB
  19: }
  20: 
  21: func NewMySQL() *MySQL {
  22: 	dsn := fmt.Sprintf(
  23: 		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
  24: 		viper.GetString("mysql.user"),
  25: 		viper.GetString("mysql.password"),
  26: 		viper.GetString("mysql.host"),
  27: 		viper.GetString("mysql.port"),
  28: 		viper.GetString("mysql.dbname"),
  29: 	)
  30: 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
  31: 	if err != nil {
  32: 		logrus.Panicf("connect to mysql failed, err=%v", err)
  33: 	}
  34: 	//db.Callback().Create().Before("gorm:create").Register("set_create_time", func(d *gorm.UseTransaction) {
  35: 	//	d.Statement.SetColumn("CreatedAt", time.Now().Format(time.DateTime))
  36: 	//})
  37: 	return &MySQL{db: db}
  38: }
  39: 
  40: func NewMySQLWithDB(db *gorm.DB) *MySQL {
  41: 	return &MySQL{db: db}
  42: }
  43: 
  44: type StockModel struct {
  45: 	ID        int64     `gorm:"column:id"`
  46: 	ProductID string    `gorm:"column:product_id"`
  47: 	Quantity  int32     `gorm:"column:quantity"`
  48: 	Version   int64     `gorm:"column:version"`
  49: 	CreatedAt time.Time `gorm:"column:created_at autoCreateTime"`
  50: 	UpdateAt  time.Time `gorm:"column:updated_at autoUpdateTime"`
  51: }
  52: 
  53: func (StockModel) TableName() string {
  54: 	return "o_stock"
  55: }
  56: 
  57: func (m *StockModel) BeforeCreate(tx *gorm.DB) (err error) {
  58: 	m.UpdateAt = time.Now()
  59: 	return nil
  60: }
  61: 
  62: func (d *MySQL) UseTransaction(tx *gorm.DB) *gorm.DB {
  63: 	if tx == nil {
  64: 		return d.db
  65: 	}
  66: 	return tx
  67: }
  68: 
  69: func (d MySQL) StartTransaction(fc func(tx *gorm.DB) error) error {
  70: 	return d.db.Transaction(fc)
  71: }
  72: 
  73: func (d MySQL) GetStockByID(ctx context.Context, query *builder.Stock) (result *StockModel, err error) {
  74: 	_, deferLog := logging.WhenMySQL(ctx, "GetStockByID", query)
  75: 	defer deferLog(result, &err)
  76: 
  77: 	err = query.Fill(d.db.WithContext(ctx)).First(&result).Error
  78: 	if err != nil {
  79: 		return nil, err
  80: 	}
  81: 	return result, nil
  82: }
  83: 
  84: func (d MySQL) BatchGetStockByID(ctx context.Context, query *builder.Stock) (result []StockModel, err error) {
  85: 	_, deferLog := logging.WhenMySQL(ctx, "BatchGetStockByID", query)
  86: 	defer deferLog(result, &err)
  87: 
  88: 	err = query.Fill(d.db.WithContext(ctx)).Find(&result).Error
  89: 	if err != nil {
  90: 		return nil, err
  91: 	}
  92: 	return result, nil
  93: }
  94: 
  95: func (d MySQL) Update(ctx context.Context, tx *gorm.DB, cond *builder.Stock, update map[string]any) (err error) {
  96: 	var returning StockModel
  97: 	_, deferLog := logging.WhenMySQL(ctx, "BatchUpdateStock", cond)
  98: 	defer deferLog(returning, &err)
  99: 
 100: 	res := cond.Fill(d.UseTransaction(tx).WithContext(ctx).Model(&returning).Clauses(clause.Returning{})).Updates(update)
 101: 	return res.Error
 102: }
 103: 
 104: func (d MySQL) Create(ctx context.Context, tx *gorm.DB, create *StockModel) (err error) {
 105: 	var returning StockModel
 106: 	_, deferLog := logging.WhenMySQL(ctx, "Create", create)
 107: 	defer deferLog(returning, &err)
 108: 
 109: 	return d.UseTransaction(tx).WithContext(ctx).Model(&returning).Clauses(clause.Returning{}).Create(create).Error
 110: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 15 | 语法块结束：关闭 import 或参数列表。 |
| 16 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 17 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 18 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 19 | 代码块结束：收束当前函数、分支或类型定义。 |
| 20 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 21 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 22 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 23 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 24 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 25 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 26 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 29 | 语法块结束：关闭 import 或参数列表。 |
| 30 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 31 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 32 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 33 | 代码块结束：收束当前函数、分支或类型定义。 |
| 34 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 35 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 36 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 37 | 返回语句：输出当前结果并结束执行路径。 |
| 38 | 代码块结束：收束当前函数、分支或类型定义。 |
| 39 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 40 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 41 | 返回语句：输出当前结果并结束执行路径。 |
| 42 | 代码块结束：收束当前函数、分支或类型定义。 |
| 43 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 44 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 45 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 46 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 47 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 48 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 49 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 50 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 51 | 代码块结束：收束当前函数、分支或类型定义。 |
| 52 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 53 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 54 | 返回语句：输出当前结果并结束执行路径。 |
| 55 | 代码块结束：收束当前函数、分支或类型定义。 |
| 56 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 57 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 58 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 59 | 返回语句：输出当前结果并结束执行路径。 |
| 60 | 代码块结束：收束当前函数、分支或类型定义。 |
| 61 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 62 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 63 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 64 | 返回语句：输出当前结果并结束执行路径。 |
| 65 | 代码块结束：收束当前函数、分支或类型定义。 |
| 66 | 返回语句：输出当前结果并结束执行路径。 |
| 67 | 代码块结束：收束当前函数、分支或类型定义。 |
| 68 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 69 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 70 | 返回语句：输出当前结果并结束执行路径。 |
| 71 | 代码块结束：收束当前函数、分支或类型定义。 |
| 72 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 73 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 74 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 75 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 76 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 77 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 78 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 79 | 返回语句：输出当前结果并结束执行路径。 |
| 80 | 代码块结束：收束当前函数、分支或类型定义。 |
| 81 | 返回语句：输出当前结果并结束执行路径。 |
| 82 | 代码块结束：收束当前函数、分支或类型定义。 |
| 83 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 84 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 85 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 86 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 87 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 88 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 89 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 90 | 返回语句：输出当前结果并结束执行路径。 |
| 91 | 代码块结束：收束当前函数、分支或类型定义。 |
| 92 | 返回语句：输出当前结果并结束执行路径。 |
| 93 | 代码块结束：收束当前函数、分支或类型定义。 |
| 94 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 95 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 96 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 97 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 98 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 99 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 100 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 101 | 返回语句：输出当前结果并结束执行路径。 |
| 102 | 代码块结束：收束当前函数、分支或类型定义。 |
| 103 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 104 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 105 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 106 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 107 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 108 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 109 | 返回语句：输出当前结果并结束执行路径。 |
| 110 | 代码块结束：收束当前函数、分支或类型定义。 |

## 提交 2: [10be7fe] 全链路可观测性建设-下(mq异步链路)

### 文件: internal/common/broker/event.go

~~~go
   1: package broker
   2: 
   3: import (
   4: 	"context"
   5: 	"encoding/json"
   6: 
   7: 	"github.com/ghost-yu/go_shop_second/common/logging"
   8: 	"github.com/pkg/errors"
   9: 	amqp "github.com/rabbitmq/amqp091-go"
  10: 	"github.com/sirupsen/logrus"
  11: )
  12: 
  13: const (
  14: 	EventOrderCreated = "order.created"
  15: 	EventOrderPaid    = "order.paid"
  16: )
  17: 
  18: type RoutingType string
  19: 
  20: const (
  21: 	FanOut RoutingType = "fan-out"
  22: 	Direct RoutingType = "direct"
  23: )
  24: 
  25: type PublishEventReq struct {
  26: 	Channel  *amqp.Channel
  27: 	Routing  RoutingType
  28: 	Queue    string
  29: 	Exchange string
  30: 	Body     any
  31: }
  32: 
  33: func PublishEvent(ctx context.Context, p PublishEventReq) (err error) {
  34: 	_, dLog := logging.WhenEventPublish(ctx, p)
  35: 	defer dLog(nil, &err)
  36: 
  37: 	if err = checkParam(p); err != nil {
  38: 		return err
  39: 	}
  40: 
  41: 	switch p.Routing {
  42: 	default:
  43: 		logrus.WithContext(ctx).Panicf("unsupported routing type: %s", string(p.Routing))
  44: 	case FanOut:
  45: 		return fanOut(ctx, p)
  46: 	case Direct:
  47: 		return directQueue(ctx, p)
  48: 	}
  49: 	return nil
  50: }
  51: 
  52: func checkParam(p PublishEventReq) error {
  53: 	if p.Channel == nil {
  54: 		return errors.New("nil channel")
  55: 	}
  56: 	return nil
  57: }
  58: 
  59: func directQueue(ctx context.Context, p PublishEventReq) (err error) {
  60: 	_, err = p.Channel.QueueDeclare(p.Queue, true, false, false, false, nil)
  61: 	if err != nil {
  62: 		return err
  63: 	}
  64: 	jsonBody, err := json.Marshal(p.Body)
  65: 	if err != nil {
  66: 		return err
  67: 	}
  68: 	return doPublish(ctx, p.Channel, p.Exchange, p.Queue, false, false, amqp.Publishing{
  69: 		ContentType:  "application/json",
  70: 		DeliveryMode: amqp.Persistent,
  71: 		Body:         jsonBody,
  72: 		Headers:      InjectRabbitMQHeaders(ctx),
  73: 	})
  74: }
  75: 
  76: func doPublish(ctx context.Context, ch *amqp.Channel, exchange, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
  77: 	if err := ch.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg); err != nil {
  78: 		logging.Warnf(ctx, nil, "_publish_event_failed||exchange=%s||key=%s||msg=%v", exchange, key, msg)
  79: 		return errors.Wrap(err, "publish event error")
  80: 	}
  81: 	return nil
  82: }
  83: 
  84: func fanOut(ctx context.Context, p PublishEventReq) (err error) {
  85: 	jsonBody, err := json.Marshal(p.Body)
  86: 	if err != nil {
  87: 		return err
  88: 	}
  89: 	return doPublish(ctx, p.Channel, p.Exchange, "", false, false, amqp.Publishing{
  90: 		ContentType:  "application/json",
  91: 		DeliveryMode: amqp.Persistent,
  92: 		Body:         jsonBody,
  93: 		Headers:      InjectRabbitMQHeaders(ctx),
  94: 	})
  95: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 语法块结束：关闭 import 或参数列表。 |
| 12 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 13 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 14 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 15 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 16 | 语法块结束：关闭 import 或参数列表。 |
| 17 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 18 | 类型定义：建立语义模型，影响方法与边界设计。 |
| 19 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 20 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 21 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 22 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 23 | 语法块结束：关闭 import 或参数列表。 |
| 24 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 25 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 26 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 29 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 30 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 31 | 代码块结束：收束当前函数、分支或类型定义。 |
| 32 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 33 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 34 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 35 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 36 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 37 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 38 | 返回语句：输出当前结果并结束执行路径。 |
| 39 | 代码块结束：收束当前函数、分支或类型定义。 |
| 40 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 41 | 多分支选择：按状态或类型分流执行路径。 |
| 42 | 默认分支：兜底处理，防止未知输入静默失败。 |
| 43 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 44 | 分支标签：定义 switch 的命中条件。 |
| 45 | 返回语句：输出当前结果并结束执行路径。 |
| 46 | 分支标签：定义 switch 的命中条件。 |
| 47 | 返回语句：输出当前结果并结束执行路径。 |
| 48 | 代码块结束：收束当前函数、分支或类型定义。 |
| 49 | 返回语句：输出当前结果并结束执行路径。 |
| 50 | 代码块结束：收束当前函数、分支或类型定义。 |
| 51 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 52 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 53 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 54 | 返回语句：输出当前结果并结束执行路径。 |
| 55 | 代码块结束：收束当前函数、分支或类型定义。 |
| 56 | 返回语句：输出当前结果并结束执行路径。 |
| 57 | 代码块结束：收束当前函数、分支或类型定义。 |
| 58 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 59 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 60 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 61 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 62 | 返回语句：输出当前结果并结束执行路径。 |
| 63 | 代码块结束：收束当前函数、分支或类型定义。 |
| 64 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 65 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 66 | 返回语句：输出当前结果并结束执行路径。 |
| 67 | 代码块结束：收束当前函数、分支或类型定义。 |
| 68 | 返回语句：输出当前结果并结束执行路径。 |
| 69 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 70 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 71 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 72 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 73 | 代码块结束：收束当前函数、分支或类型定义。 |
| 74 | 代码块结束：收束当前函数、分支或类型定义。 |
| 75 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 76 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 77 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 78 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 79 | 返回语句：输出当前结果并结束执行路径。 |
| 80 | 代码块结束：收束当前函数、分支或类型定义。 |
| 81 | 返回语句：输出当前结果并结束执行路径。 |
| 82 | 代码块结束：收束当前函数、分支或类型定义。 |
| 83 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 84 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 85 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 86 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 87 | 返回语句：输出当前结果并结束执行路径。 |
| 88 | 代码块结束：收束当前函数、分支或类型定义。 |
| 89 | 返回语句：输出当前结果并结束执行路径。 |
| 90 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 91 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 92 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 93 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 94 | 代码块结束：收束当前函数、分支或类型定义。 |
| 95 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/broker/rabbitmq.go

~~~go
   1: package broker
   2: 
   3: import (
   4: 	"context"
   5: 	"fmt"
   6: 	"time"
   7: 
   8: 	_ "github.com/ghost-yu/go_shop_second/common/config"
   9: 	"github.com/ghost-yu/go_shop_second/common/logging"
  10: 	amqp "github.com/rabbitmq/amqp091-go"
  11: 	"github.com/sirupsen/logrus"
  12: 	"github.com/spf13/viper"
  13: 	"go.opentelemetry.io/otel"
  14: )
  15: 
  16: const (
  17: 	DLX                = "dlx"
  18: 	DLQ                = "dlq"
  19: 	amqpRetryHeaderKey = "x-retry-count"
  20: )
  21: 
  22: var (
  23: 	maxRetryCount = viper.GetInt64("rabbitmq.max-retry")
  24: )
  25: 
  26: func Connect(user, password, host, port string) (*amqp.Channel, func() error) {
  27: 	address := fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port)
  28: 	conn, err := amqp.Dial(address)
  29: 	if err != nil {
  30: 		logrus.Fatal(err)
  31: 	}
  32: 	ch, err := conn.Channel()
  33: 	if err != nil {
  34: 		logrus.Fatal(err)
  35: 	}
  36: 	err = ch.ExchangeDeclare(EventOrderCreated, "direct", true, false, false, false, nil)
  37: 	if err != nil {
  38: 		logrus.Fatal(err)
  39: 	}
  40: 	err = ch.ExchangeDeclare(EventOrderPaid, "fanout", true, false, false, false, nil)
  41: 	if err != nil {
  42: 		logrus.Fatal(err)
  43: 	}
  44: 	if err = createDLX(ch); err != nil {
  45: 		logrus.Fatal(err)
  46: 	}
  47: 	return ch, conn.Close
  48: }
  49: 
  50: func createDLX(ch *amqp.Channel) error {
  51: 	q, err := ch.QueueDeclare("share_queue", true, false, false, false, nil)
  52: 	if err != nil {
  53: 		return err
  54: 	}
  55: 	err = ch.ExchangeDeclare(DLX, "fanout", true, false, false, false, nil)
  56: 	if err != nil {
  57: 		return err
  58: 	}
  59: 	err = ch.QueueBind(q.Name, "", DLX, false, nil)
  60: 	if err != nil {
  61: 		return err
  62: 	}
  63: 	_, err = ch.QueueDeclare(DLQ, true, false, false, false, nil)
  64: 	return err
  65: }
  66: 
  67: func HandleRetry(ctx context.Context, ch *amqp.Channel, d *amqp.Delivery) (err error) {
  68: 	fields, dLog := logging.WhenRequest(ctx, "HandleRetry", map[string]any{
  69: 		"delivery":        d,
  70: 		"max_retry_count": maxRetryCount,
  71: 	})
  72: 	defer dLog(nil, &err)
  73: 
  74: 	if d.Headers == nil {
  75: 		d.Headers = amqp.Table{}
  76: 	}
  77: 	retryCount, ok := d.Headers[amqpRetryHeaderKey].(int64)
  78: 	if !ok {
  79: 		retryCount = 0
  80: 	}
  81: 	retryCount++
  82: 	d.Headers[amqpRetryHeaderKey] = retryCount
  83: 	fields["retry_count"] = retryCount
  84: 
  85: 	if retryCount >= maxRetryCount {
  86: 		logrus.WithContext(ctx).Infof("moving message %s to dlq", d.MessageId)
  87: 		return doPublish(ctx, ch, "", DLQ, false, false, amqp.Publishing{
  88: 			Headers:      d.Headers,
  89: 			ContentType:  "application/json",
  90: 			Body:         d.Body,
  91: 			DeliveryMode: amqp.Persistent,
  92: 		})
  93: 	}
  94: 	logrus.WithContext(ctx).Debugf("retring message %s, count=%d", d.MessageId, retryCount)
  95: 	time.Sleep(time.Second * time.Duration(retryCount))
  96: 	return doPublish(ctx, ch, "", DLQ, false, false, amqp.Publishing{
  97: 		Headers:      d.Headers,
  98: 		ContentType:  "application/json",
  99: 		Body:         d.Body,
 100: 		DeliveryMode: amqp.Persistent,
 101: 	})
 102: }
 103: 
 104: type RabbitMQHeaderCarrier map[string]interface{}
 105: 
 106: func (r RabbitMQHeaderCarrier) Get(key string) string {
 107: 	value, ok := r[key]
 108: 	if !ok {
 109: 		return ""
 110: 	}
 111: 	return value.(string)
 112: }
 113: 
 114: func (r RabbitMQHeaderCarrier) Set(key string, value string) {
 115: 	r[key] = value
 116: }
 117: 
 118: func (r RabbitMQHeaderCarrier) Keys() []string {
 119: 	keys := make([]string, len(r))
 120: 	i := 0
 121: 	for key := range r {
 122: 		keys[i] = key
 123: 		i++
 124: 	}
 125: 	return keys
 126: }
 127: 
 128: func InjectRabbitMQHeaders(ctx context.Context) map[string]interface{} {
 129: 	carrier := make(RabbitMQHeaderCarrier)
 130: 	otel.GetTextMapPropagator().Inject(ctx, carrier)
 131: 	return carrier
 132: }
 133: 
 134: func ExtractRabbitMQHeaders(ctx context.Context, headers map[string]interface{}) context.Context {
 135: 	return otel.GetTextMapPropagator().Extract(ctx, RabbitMQHeaderCarrier(headers))
 136: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 8 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 语法块结束：关闭 import 或参数列表。 |
| 15 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 16 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 17 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 18 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 19 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 20 | 语法块结束：关闭 import 或参数列表。 |
| 21 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 22 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 23 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 24 | 语法块结束：关闭 import 或参数列表。 |
| 25 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 26 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 27 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 28 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 29 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 30 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 31 | 代码块结束：收束当前函数、分支或类型定义。 |
| 32 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 33 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 34 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 35 | 代码块结束：收束当前函数、分支或类型定义。 |
| 36 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 37 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 38 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 39 | 代码块结束：收束当前函数、分支或类型定义。 |
| 40 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 41 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 42 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 43 | 代码块结束：收束当前函数、分支或类型定义。 |
| 44 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 45 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 46 | 代码块结束：收束当前函数、分支或类型定义。 |
| 47 | 返回语句：输出当前结果并结束执行路径。 |
| 48 | 代码块结束：收束当前函数、分支或类型定义。 |
| 49 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 50 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 51 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 52 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 53 | 返回语句：输出当前结果并结束执行路径。 |
| 54 | 代码块结束：收束当前函数、分支或类型定义。 |
| 55 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 56 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 57 | 返回语句：输出当前结果并结束执行路径。 |
| 58 | 代码块结束：收束当前函数、分支或类型定义。 |
| 59 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 60 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 61 | 返回语句：输出当前结果并结束执行路径。 |
| 62 | 代码块结束：收束当前函数、分支或类型定义。 |
| 63 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 64 | 返回语句：输出当前结果并结束执行路径。 |
| 65 | 代码块结束：收束当前函数、分支或类型定义。 |
| 66 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 67 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 68 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 69 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 70 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 71 | 代码块结束：收束当前函数、分支或类型定义。 |
| 72 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 73 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 74 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 75 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 76 | 代码块结束：收束当前函数、分支或类型定义。 |
| 77 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 78 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 79 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 80 | 代码块结束：收束当前函数、分支或类型定义。 |
| 81 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 82 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 83 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 84 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 85 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 86 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 87 | 返回语句：输出当前结果并结束执行路径。 |
| 88 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 89 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 90 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 91 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 92 | 代码块结束：收束当前函数、分支或类型定义。 |
| 93 | 代码块结束：收束当前函数、分支或类型定义。 |
| 94 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 95 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 96 | 返回语句：输出当前结果并结束执行路径。 |
| 97 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 98 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 99 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 100 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 101 | 代码块结束：收束当前函数、分支或类型定义。 |
| 102 | 代码块结束：收束当前函数、分支或类型定义。 |
| 103 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 104 | 类型定义：建立语义模型，影响方法与边界设计。 |
| 105 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 106 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 107 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 108 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 109 | 返回语句：输出当前结果并结束执行路径。 |
| 110 | 代码块结束：收束当前函数、分支或类型定义。 |
| 111 | 返回语句：输出当前结果并结束执行路径。 |
| 112 | 代码块结束：收束当前函数、分支或类型定义。 |
| 113 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 114 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 115 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 116 | 代码块结束：收束当前函数、分支或类型定义。 |
| 117 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 118 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 119 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 120 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 121 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 122 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 123 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 124 | 代码块结束：收束当前函数、分支或类型定义。 |
| 125 | 返回语句：输出当前结果并结束执行路径。 |
| 126 | 代码块结束：收束当前函数、分支或类型定义。 |
| 127 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 128 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 129 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 130 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 131 | 返回语句：输出当前结果并结束执行路径。 |
| 132 | 代码块结束：收束当前函数、分支或类型定义。 |
| 133 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 134 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 135 | 返回语句：输出当前结果并结束执行路径。 |
| 136 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/common/logging/when.go

~~~go
   1: package logging
   2: 
   3: import (
   4: 	"context"
   5: 	"time"
   6: 
   7: 	"github.com/sirupsen/logrus"
   8: )
   9: 
  10: func WhenCommandExecute(ctx context.Context, commandName string, cmd any, err error) {
  11: 	fields := logrus.Fields{
  12: 		"cmd": cmd,
  13: 	}
  14: 	if err == nil {
  15: 		logf(ctx, logrus.InfoLevel, fields, "%s_command_success", commandName)
  16: 	} else {
  17: 		logf(ctx, logrus.ErrorLevel, fields, "%s_command_failed", commandName)
  18: 	}
  19: }
  20: 
  21: func WhenRequest(ctx context.Context, method string, args ...any) (logrus.Fields, func(any, *error)) {
  22: 	fields := logrus.Fields{
  23: 		Method: method,
  24: 		Args:   formatArgs(args),
  25: 	}
  26: 	start := time.Now()
  27: 	return fields, func(resp any, err *error) {
  28: 		level, msg := logrus.InfoLevel, "_request_success"
  29: 		fields[Cost] = time.Since(start).Milliseconds()
  30: 		fields[Response] = resp
  31: 
  32: 		if err != nil && (*err != nil) {
  33: 			level, msg = logrus.ErrorLevel, "_request_failed"
  34: 			fields[Error] = (*err).Error()
  35: 		}
  36: 
  37: 		logf(ctx, level, fields, "%s", msg)
  38: 	}
  39: }
  40: 
  41: func WhenEventPublish(ctx context.Context, args ...any) (logrus.Fields, func(any, *error)) {
  42: 	fields := logrus.Fields{
  43: 		Args: formatArgs(args),
  44: 	}
  45: 	start := time.Now()
  46: 	return fields, func(resp any, err *error) {
  47: 		level, msg := logrus.InfoLevel, "_mq_publish_success"
  48: 		fields[Cost] = time.Since(start).Milliseconds()
  49: 		fields[Response] = resp
  50: 
  51: 		if err != nil && (*err != nil) {
  52: 			level, msg = logrus.ErrorLevel, "_mq_publish_failed"
  53: 			fields[Error] = (*err).Error()
  54: 		}
  55: 
  56: 		logf(ctx, level, fields, "%s", msg)
  57: 	}
  58: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 语法块结束：关闭 import 或参数列表。 |
| 9 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 10 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 11 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 12 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 13 | 代码块结束：收束当前函数、分支或类型定义。 |
| 14 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 15 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 16 | 代码块结束：收束当前函数、分支或类型定义。 |
| 17 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 18 | 代码块结束：收束当前函数、分支或类型定义。 |
| 19 | 代码块结束：收束当前函数、分支或类型定义。 |
| 20 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 21 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 22 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 23 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 24 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 25 | 代码块结束：收束当前函数、分支或类型定义。 |
| 26 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 27 | 返回语句：输出当前结果并结束执行路径。 |
| 28 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 29 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 30 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 31 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 32 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 33 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 34 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 35 | 代码块结束：收束当前函数、分支或类型定义。 |
| 36 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 代码块结束：收束当前函数、分支或类型定义。 |
| 39 | 代码块结束：收束当前函数、分支或类型定义。 |
| 40 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 41 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 42 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 43 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 44 | 代码块结束：收束当前函数、分支或类型定义。 |
| 45 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 46 | 返回语句：输出当前结果并结束执行路径。 |
| 47 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 48 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 49 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 50 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 51 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 52 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 53 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 54 | 代码块结束：收束当前函数、分支或类型定义。 |
| 55 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 56 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 57 | 代码块结束：收束当前函数、分支或类型定义。 |
| 58 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/order/app/command/create_order.go

~~~go
   1: package command
   2: 
   3: import (
   4: 	"context"
   5: 	"fmt"
   6: 
   7: 	"github.com/ghost-yu/go_shop_second/common/broker"
   8: 	"github.com/ghost-yu/go_shop_second/common/decorator"
   9: 	"github.com/ghost-yu/go_shop_second/common/logging"
  10: 	"github.com/ghost-yu/go_shop_second/order/app/query"
  11: 	"github.com/ghost-yu/go_shop_second/order/convertor"
  12: 	domain "github.com/ghost-yu/go_shop_second/order/domain/order"
  13: 	"github.com/ghost-yu/go_shop_second/order/entity"
  14: 	"github.com/pkg/errors"
  15: 	amqp "github.com/rabbitmq/amqp091-go"
  16: 	"github.com/sirupsen/logrus"
  17: 	"go.opentelemetry.io/otel"
  18: 	"google.golang.org/grpc/status"
  19: )
  20: 
  21: type CreateOrder struct {
  22: 	CustomerID string
  23: 	Items      []*entity.ItemWithQuantity
  24: }
  25: 
  26: type CreateOrderResult struct {
  27: 	OrderID string
  28: }
  29: 
  30: type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]
  31: 
  32: type createOrderHandler struct {
  33: 	orderRepo domain.Repository
  34: 	stockGRPC query.StockService
  35: 	channel   *amqp.Channel
  36: }
  37: 
  38: func NewCreateOrderHandler(
  39: 	orderRepo domain.Repository,
  40: 	stockGRPC query.StockService,
  41: 	channel *amqp.Channel,
  42: 	logger *logrus.Entry,
  43: 	metricClient decorator.MetricsClient,
  44: ) CreateOrderHandler {
  45: 	if orderRepo == nil {
  46: 		panic("nil orderRepo")
  47: 	}
  48: 	if stockGRPC == nil {
  49: 		panic("nil stockGRPC")
  50: 	}
  51: 	if channel == nil {
  52: 		panic("nil channel ")
  53: 	}
  54: 	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
  55: 		createOrderHandler{
  56: 			orderRepo: orderRepo,
  57: 			stockGRPC: stockGRPC,
  58: 			channel:   channel,
  59: 		},
  60: 		logger,
  61: 		metricClient,
  62: 	)
  63: }
  64: 
  65: func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
  66: 	var err error
  67: 	defer logging.WhenCommandExecute(ctx, "CreateOrderHandler", cmd, err)
  68: 
  69: 	t := otel.Tracer("rabbitmq")
  70: 	ctx, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderCreated))
  71: 	defer span.End()
  72: 
  73: 	validItems, err := c.validate(ctx, cmd.Items)
  74: 	if err != nil {
  75: 		return nil, err
  76: 	}
  77: 	pendingOrder, err := domain.NewPendingOrder(cmd.CustomerID, validItems)
  78: 	if err != nil {
  79: 		return nil, err
  80: 	}
  81: 	o, err := c.orderRepo.Create(ctx, pendingOrder)
  82: 	if err != nil {
  83: 		return nil, err
  84: 	}
  85: 
  86: 	err = broker.PublishEvent(ctx, broker.PublishEventReq{
  87: 		Channel:  c.channel,
  88: 		Routing:  broker.Direct,
  89: 		Queue:    broker.EventOrderCreated,
  90: 		Exchange: "",
  91: 		Body:     o,
  92: 	})
  93: 	if err != nil {
  94: 		return nil, errors.Wrapf(err, "publish event error q.Name=%s", broker.EventOrderCreated)
  95: 	}
  96: 
  97: 	return &CreateOrderResult{OrderID: o.ID}, nil
  98: }
  99: 
 100: func (c createOrderHandler) validate(ctx context.Context, items []*entity.ItemWithQuantity) ([]*entity.Item, error) {
 101: 	if len(items) == 0 {
 102: 		return nil, errors.New("must have at least one item")
 103: 	}
 104: 	items = packItems(items)
 105: 	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, convertor.NewItemWithQuantityConvertor().EntitiesToProtos(items))
 106: 	if err != nil {
 107: 		return nil, status.Convert(err).Err()
 108: 	}
 109: 	return convertor.NewItemConvertor().ProtosToEntities(resp.Items), nil
 110: }
 111: 
 112: func packItems(items []*entity.ItemWithQuantity) []*entity.ItemWithQuantity {
 113: 	merged := make(map[string]int32)
 114: 	for _, item := range items {
 115: 		merged[item.ID] += item.Quantity
 116: 	}
 117: 	var res []*entity.ItemWithQuantity
 118: 	for id, quantity := range merged {
 119: 		res = append(res, &entity.ItemWithQuantity{
 120: 			ID:       id,
 121: 			Quantity: quantity,
 122: 		})
 123: 	}
 124: 	return res
 125: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 15 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 16 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 17 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 18 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 19 | 语法块结束：关闭 import 或参数列表。 |
| 20 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 21 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 22 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 23 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 24 | 代码块结束：收束当前函数、分支或类型定义。 |
| 25 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 26 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 27 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 28 | 代码块结束：收束当前函数、分支或类型定义。 |
| 29 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 30 | 类型定义：建立语义模型，影响方法与边界设计。 |
| 31 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 32 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 33 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 34 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 35 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 36 | 代码块结束：收束当前函数、分支或类型定义。 |
| 37 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 38 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 39 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 40 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 41 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 42 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 43 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 44 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 45 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 46 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 47 | 代码块结束：收束当前函数、分支或类型定义。 |
| 48 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 49 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 50 | 代码块结束：收束当前函数、分支或类型定义。 |
| 51 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 52 | panic：不可恢复错误路径，常用于关键初始化失败。 |
| 53 | 代码块结束：收束当前函数、分支或类型定义。 |
| 54 | 返回语句：输出当前结果并结束执行路径。 |
| 55 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 56 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 57 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 58 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 59 | 代码块结束：收束当前函数、分支或类型定义。 |
| 60 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 61 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 62 | 语法块结束：关闭 import 或参数列表。 |
| 63 | 代码块结束：收束当前函数、分支或类型定义。 |
| 64 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 65 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 66 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 67 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 68 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 69 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 70 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 71 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 72 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 73 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 74 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 75 | 返回语句：输出当前结果并结束执行路径。 |
| 76 | 代码块结束：收束当前函数、分支或类型定义。 |
| 77 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 78 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 79 | 返回语句：输出当前结果并结束执行路径。 |
| 80 | 代码块结束：收束当前函数、分支或类型定义。 |
| 81 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 82 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 83 | 返回语句：输出当前结果并结束执行路径。 |
| 84 | 代码块结束：收束当前函数、分支或类型定义。 |
| 85 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 86 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 87 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 88 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 89 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 90 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 91 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 92 | 代码块结束：收束当前函数、分支或类型定义。 |
| 93 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 94 | 返回语句：输出当前结果并结束执行路径。 |
| 95 | 代码块结束：收束当前函数、分支或类型定义。 |
| 96 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 97 | 返回语句：输出当前结果并结束执行路径。 |
| 98 | 代码块结束：收束当前函数、分支或类型定义。 |
| 99 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 100 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 101 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 102 | 返回语句：输出当前结果并结束执行路径。 |
| 103 | 代码块结束：收束当前函数、分支或类型定义。 |
| 104 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 105 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 106 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 107 | 返回语句：输出当前结果并结束执行路径。 |
| 108 | 代码块结束：收束当前函数、分支或类型定义。 |
| 109 | 返回语句：输出当前结果并结束执行路径。 |
| 110 | 代码块结束：收束当前函数、分支或类型定义。 |
| 111 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 112 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 113 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 114 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 115 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 116 | 代码块结束：收束当前函数、分支或类型定义。 |
| 117 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 118 | 循环结构：遍历数据或重复执行，需关注终止条件与副作用。 |
| 119 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 120 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 121 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 122 | 代码块结束：收束当前函数、分支或类型定义。 |
| 123 | 代码块结束：收束当前函数、分支或类型定义。 |
| 124 | 返回语句：输出当前结果并结束执行路径。 |
| 125 | 代码块结束：收束当前函数、分支或类型定义。 |

### 文件: internal/payment/http.go

~~~go
   1: package main
   2: 
   3: import (
   4: 	"encoding/json"
   5: 	"fmt"
   6: 	"io"
   7: 	"net/http"
   8: 
   9: 	"github.com/ghost-yu/go_shop_second/common/broker"
  10: 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
  11: 	"github.com/ghost-yu/go_shop_second/common/logging"
  12: 	"github.com/ghost-yu/go_shop_second/payment/domain"
  13: 	"github.com/gin-gonic/gin"
  14: 	"github.com/pkg/errors"
  15: 	amqp "github.com/rabbitmq/amqp091-go"
  16: 	"github.com/sirupsen/logrus"
  17: 	"github.com/spf13/viper"
  18: 	"github.com/stripe/stripe-go/v79"
  19: 	"github.com/stripe/stripe-go/v79/webhook"
  20: 	"go.opentelemetry.io/otel"
  21: )
  22: 
  23: type PaymentHandler struct {
  24: 	channel *amqp.Channel
  25: }
  26: 
  27: func NewPaymentHandler(ch *amqp.Channel) *PaymentHandler {
  28: 	return &PaymentHandler{channel: ch}
  29: }
  30: 
  31: // stripe listen --forward-to localhost:8284/api/webhook
  32: func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
  33: 	c.POST("/api/webhook", h.handleWebhook)
  34: }
  35: 
  36: func (h *PaymentHandler) handleWebhook(c *gin.Context) {
  37: 	logrus.WithContext(c.Request.Context()).Info("receive webhook from stripe")
  38: 	var err error
  39: 	defer func() {
  40: 		if err != nil {
  41: 			logging.Warnf(c.Request.Context(), nil, "handleWebhook err=%v", err)
  42: 		} else {
  43: 			logging.Infof(c.Request.Context(), nil, "%s", "handleWebhook success")
  44: 		}
  45: 	}()
  46: 
  47: 	const MaxBodyBytes = int64(65536)
  48: 	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
  49: 	payload, err := io.ReadAll(c.Request.Body)
  50: 	if err != nil {
  51: 		err = errors.Wrap(err, "Error reading request body")
  52: 		c.JSON(http.StatusServiceUnavailable, err.Error())
  53: 		return
  54: 	}
  55: 
  56: 	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"),
  57: 		viper.GetString("ENDPOINT_STRIPE_SECRET"))
  58: 
  59: 	if err != nil {
  60: 		err = errors.Wrap(err, "error verifying webhook signature")
  61: 		c.JSON(http.StatusBadRequest, err.Error())
  62: 		return
  63: 	}
  64: 
  65: 	switch event.Type {
  66: 	case stripe.EventTypeCheckoutSessionCompleted:
  67: 		var session stripe.CheckoutSession
  68: 		if err = json.Unmarshal(event.Data.Raw, &session); err != nil {
  69: 			err = errors.Wrap(err, "error unmarshal event.data.raw into session")
  70: 			c.JSON(http.StatusBadRequest, err.Error())
  71: 			return
  72: 		}
  73: 
  74: 		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
  75: 			var items []*orderpb.Item
  76: 			_ = json.Unmarshal([]byte(session.Metadata["items"]), &items)
  77: 
  78: 			tr := otel.Tracer("rabbitmq")
  79: 			ctx, span := tr.Start(c.Request.Context(), fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderPaid))
  80: 			defer span.End()
  81: 
  82: 			_ = broker.PublishEvent(ctx, broker.PublishEventReq{
  83: 				Channel:  h.channel,
  84: 				Routing:  broker.FanOut,
  85: 				Queue:    "",
  86: 				Exchange: broker.EventOrderPaid,
  87: 				Body: &domain.Order{
  88: 					ID:          session.Metadata["orderID"],
  89: 					CustomerID:  session.Metadata["customerID"],
  90: 					Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
  91: 					PaymentLink: session.Metadata["paymentLink"],
  92: 					Items:       items,
  93: 				},
  94: 			})
  95: 		}
  96: 	}
  97: 	c.JSON(http.StatusOK, nil)
  98: }
~~~

| 行号 | 中文深度解释 |
| --- | --- |
| 1 | 包声明：定义文件归属命名空间，决定编译与可见性边界。 |
| 2 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 3 | 导入块开始：集中声明依赖，便于快速理解耦合面。 |
| 4 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 5 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 6 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 7 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 8 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 9 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 10 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 11 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 12 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 13 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 14 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 15 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 16 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 17 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 18 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 19 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 20 | 具体依赖项：该包被当前文件使用，体现职责方向。 |
| 21 | 语法块结束：关闭 import 或参数列表。 |
| 22 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 23 | 结构体定义：声明数据载体，承载状态或依赖。 |
| 24 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 25 | 代码块结束：收束当前函数、分支或类型定义。 |
| 26 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 27 | 函数定义：声明逻辑单元，入参与返回值决定职责边界。 |
| 28 | 返回语句：输出当前结果并结束执行路径。 |
| 29 | 代码块结束：收束当前函数、分支或类型定义。 |
| 30 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 31 | 注释：解释意图、风险或待办，帮助理解设计。 |
| 32 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 33 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 34 | 代码块结束：收束当前函数、分支或类型定义。 |
| 35 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 36 | 方法定义：函数绑定接收者类型，体现对象行为。 |
| 37 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 38 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 39 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 40 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 41 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 42 | 代码块结束：收束当前函数、分支或类型定义。 |
| 43 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 44 | 代码块结束：收束当前函数、分支或类型定义。 |
| 45 | 代码块结束：收束当前函数、分支或类型定义。 |
| 46 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 47 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 48 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 49 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 50 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 51 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 52 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 53 | 返回语句：输出当前结果并结束执行路径。 |
| 54 | 代码块结束：收束当前函数、分支或类型定义。 |
| 55 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 56 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 57 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 58 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 59 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 60 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 61 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 62 | 返回语句：输出当前结果并结束执行路径。 |
| 63 | 代码块结束：收束当前函数、分支或类型定义。 |
| 64 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 65 | 多分支选择：按状态或类型分流执行路径。 |
| 66 | 分支标签：定义 switch 的命中条件。 |
| 67 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 68 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 69 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 70 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 71 | 返回语句：输出当前结果并结束执行路径。 |
| 72 | 代码块结束：收束当前函数、分支或类型定义。 |
| 73 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 74 | 条件分支：进行校验、错误拦截或流程分叉。 |
| 75 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 76 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 77 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 78 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 79 | 短变量声明：就地定义并初始化，收窄作用域。 |
| 80 | defer：函数退出前执行，常用于资源释放与收尾。 |
| 81 | 空行：用于分隔逻辑块，提升可读性与维护效率。 |
| 82 | 赋值语句：更新状态或绑定数据，可能影响后续流程。 |
| 83 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 84 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 85 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 86 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 87 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 88 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 89 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 90 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 91 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 92 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 93 | 代码块结束：收束当前函数、分支或类型定义。 |
| 94 | 代码块结束：收束当前函数、分支或类型定义。 |
| 95 | 代码块结束：收束当前函数、分支或类型定义。 |
| 96 | 代码块结束：收束当前函数、分支或类型定义。 |
| 97 | 执行语句：参与当前实现，需结合上下文理解业务意图。 |
| 98 | 代码块结束：收束当前函数、分支或类型定义。 |


