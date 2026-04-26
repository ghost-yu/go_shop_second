# 中文深度逐行解读（自动生成）

- 仓库: go_shop_second
- 起始引用: lesson61
- 结束引用: lesson62
- 生成时间: 2026-04-06 18:34:27 +08:00
- 解析范围: 仅 *.go，自动排除 *.pb.go 与 *.gen.go
- 解析模式: 祖先链模式（按提交）

## 提交 1: [10be7fe] 全链路可观测性建设-下(mq异步链路)

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


