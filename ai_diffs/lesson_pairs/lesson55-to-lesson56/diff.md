# Lesson Pair Diff Report

- FromBranch: lesson55
- ToBranch: lesson56

## Short Summary

~~~text
 2 files changed, 70 insertions(+), 2 deletions(-)
~~~

## File Stats

~~~text
 internal/common/logging/mysql.go                  | 57 +++++++++++++++++++++++
 internal/stock/infrastructure/persistent/mysql.go | 15 +++++-
 2 files changed, 70 insertions(+), 2 deletions(-)
~~~

## Commit Comparison

~~~text
> d49c66a mysql log
~~~

## Changed Files

~~~text
internal/common/logging/mysql.go
internal/stock/infrastructure/persistent/mysql.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/common/logging/mysql.go
internal/stock/infrastructure/persistent/mysql.go
~~~

## Full Diff

~~~diff
diff --git a/internal/common/logging/mysql.go b/internal/common/logging/mysql.go
new file mode 100644
index 0000000..10159f7
--- /dev/null
+++ b/internal/common/logging/mysql.go
@@ -0,0 +1,57 @@
+package logging
+
+import (
+	"context"
+	"encoding/json"
+	"strings"
+	"time"
+
+	"github.com/sirupsen/logrus"
+)
+
+const (
+	Method   = "method"
+	Args     = "args"
+	Cost     = "cost_ms"
+	Response = "response"
+	Error    = "err"
+)
+
+func WhenMySQL(ctx context.Context, method string, args ...any) (logrus.Fields, func(any, *error)) {
+	fields := logrus.Fields{
+		Method: method,
+		Args:   formatMySQLArgs(args),
+	}
+	start := time.Now()
+	return fields, func(resp any, err *error) {
+		level, msg := logrus.InfoLevel, "mysql_success"
+		fields[Cost] = time.Since(start).Milliseconds()
+		fields[Response] = resp
+
+		if err != nil && (*err != nil) {
+			level, msg = logrus.ErrorLevel, "mysql_error"
+			fields[Error] = (*err).Error()
+		}
+
+		logrus.WithContext(ctx).WithFields(fields).Logf(level, "%s", msg)
+	}
+}
+
+func formatMySQLArgs(args []any) string {
+	var item []string
+	for _, arg := range args {
+		item = append(item, formatMySQLArg(arg))
+	}
+	return strings.Join(item, "||")
+}
+
+func formatMySQLArg(arg any) string {
+	switch v := arg.(type) {
+	default:
+		bytes, err := json.Marshal(v)
+		if err != nil {
+			return "unsupported type in formatMySQLArg||err=" + err.Error()
+		}
+		return string(bytes)
+	}
+}
diff --git a/internal/stock/infrastructure/persistent/mysql.go b/internal/stock/infrastructure/persistent/mysql.go
index c6eb728..a6d638d 100644
--- a/internal/stock/infrastructure/persistent/mysql.go
+++ b/internal/stock/infrastructure/persistent/mysql.go
@@ -5,10 +5,12 @@ import (
 	"fmt"
 	"time"
 
+	"github.com/ghost-yu/go_shop_second/common/logging"
 	"github.com/sirupsen/logrus"
 	"github.com/spf13/viper"
 	"gorm.io/driver/mysql"
 	"gorm.io/gorm"
+	"gorm.io/gorm/clause"
 )
 
 type MySQL struct {
@@ -28,6 +30,9 @@ func NewMySQL() *MySQL {
 	if err != nil {
 		logrus.Panicf("connect to mysql failed, err=%v", err)
 	}
+	//db.Callback().Create().Before("gorm:create").Register("set_create_time", func(d *gorm.DB) {
+	//	d.Statement.SetColumn("CreatedAt", time.Now().Format(time.DateTime))
+	//})
 	return &MySQL{db: db}
 }
 
@@ -58,8 +63,10 @@ func (d MySQL) StartTransaction(fc func(tx *gorm.DB) error) error {
 }
 
 func (d MySQL) BatchGetStockByID(ctx context.Context, productIDs []string) ([]StockModel, error) {
+	_, deferLog := logging.WhenMySQL(ctx, "BatchGetStockByID", productIDs)
 	var result []StockModel
-	tx := d.db.WithContext(ctx).Where("product_id IN ?", productIDs).Find(&result)
+	tx := d.db.WithContext(ctx).Clauses(clause.Returning{}).Where("product_id IN ?", productIDs).Find(&result)
+	defer deferLog(result, &tx.Error)
 	if tx.Error != nil {
 		return nil, tx.Error
 	}
@@ -67,5 +74,9 @@ func (d MySQL) BatchGetStockByID(ctx context.Context, productIDs []string) ([]St
 }
 
 func (d MySQL) Create(ctx context.Context, create *StockModel) error {
-	return d.db.WithContext(ctx).Create(create).Error
+	_, deferLog := logging.WhenMySQL(ctx, "Create", create)
+	var returning StockModel
+	err := d.db.WithContext(ctx).Model(&returning).Clauses(clause.Returning{}).Create(create).Error
+	defer deferLog(returning, &err)
+	return err
 }
~~~
