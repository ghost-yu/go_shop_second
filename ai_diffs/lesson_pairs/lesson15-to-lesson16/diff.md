# Lesson Pair Diff Report

- FromBranch: lesson15
- ToBranch: lesson16

## Short Summary

~~~text
 2 files changed, 41 insertions(+)
~~~

## File Stats

~~~text
 .../payment/infrastructure/consumer/consumer.go    | 38 ++++++++++++++++++++++
 internal/payment/main.go                           |  3 ++
 2 files changed, 41 insertions(+)
~~~

## Commit Comparison

~~~text
> b96b769 payment consumer
~~~

## Changed Files

~~~text
internal/payment/infrastructure/consumer/consumer.go
internal/payment/main.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/payment/infrastructure/consumer/consumer.go
internal/payment/main.go
~~~

## Full Diff

~~~diff
diff --git a/internal/payment/infrastructure/consumer/consumer.go b/internal/payment/infrastructure/consumer/consumer.go
new file mode 100644
index 0000000..b1caba1
--- /dev/null
+++ b/internal/payment/infrastructure/consumer/consumer.go
@@ -0,0 +1,38 @@
+package consumer
+
+import (
+	"github.com/ghost-yu/go_shop_second/common/broker"
+	amqp "github.com/rabbitmq/amqp091-go"
+	"github.com/sirupsen/logrus"
+)
+
+type Consumer struct{}
+
+func NewConsumer() *Consumer {
+	return &Consumer{}
+}
+
+func (c *Consumer) Listen(ch *amqp.Channel) {
+	q, err := ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
+	if err != nil {
+		logrus.Fatal(err)
+	}
+
+	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
+	if err != nil {
+		logrus.Warnf("fail to consume: queue=%s, err=%v", q.Name, err)
+	}
+
+	var forever chan struct{}
+	go func() {
+		for msg := range msgs {
+			c.handleMessage(msg, q, ch)
+		}
+	}()
+	<-forever
+}
+
+func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue, ch *amqp.Channel) {
+	logrus.Infof("Payment receive a message from %s, msg=%v", q.Name, string(msg.Body))
+	_ = msg.Ack(false)
+}
diff --git a/internal/payment/main.go b/internal/payment/main.go
index b8d159b..8bad5a9 100644
--- a/internal/payment/main.go
+++ b/internal/payment/main.go
@@ -5,6 +5,7 @@ import (
 	"github.com/ghost-yu/go_shop_second/common/config"
 	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/ghost-yu/go_shop_second/common/server"
+	"github.com/ghost-yu/go_shop_second/payment/infrastructure/consumer"
 	"github.com/sirupsen/logrus"
 	"github.com/spf13/viper"
 )
@@ -30,6 +31,8 @@ func main() {
 		_ = closeCh()
 	}()
 
+	go consumer.NewConsumer().Listen(ch)
+
 	paymentHandler := NewPaymentHandler()
 	switch serverType {
 	case "http":
~~~
