# Lesson Pair Diff Report

- FromBranch: lesson57
- ToBranch: lesson58

## Short Summary

~~~text
 6 files changed, 112 insertions(+), 67 deletions(-)
~~~

## File Stats

~~~text
 internal/common/logging/mysql.go                   | 24 +++++++---
 internal/common/util/jsonutil.go                   |  8 ++++
 internal/stock/adapters/stock_mysql_repository.go  | 52 +++++++++-------------
 .../stock/adapters/stock_mysql_repository_test.go  |  9 ++--
 .../infrastructure/persistent/builder/stock.go     | 52 ++++++++++++----------
 internal/stock/infrastructure/persistent/mysql.go  | 34 ++++++++++++--
 6 files changed, 112 insertions(+), 67 deletions(-)
~~~

## Commit Comparison

~~~text
> c0504b7 std sql log
~~~

## Changed Files

~~~text
internal/common/logging/mysql.go
internal/common/util/jsonutil.go
internal/stock/adapters/stock_mysql_repository.go
internal/stock/adapters/stock_mysql_repository_test.go
internal/stock/infrastructure/persistent/builder/stock.go
internal/stock/infrastructure/persistent/mysql.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/common/logging/mysql.go
internal/common/util/jsonutil.go
internal/stock/adapters/stock_mysql_repository.go
internal/stock/adapters/stock_mysql_repository_test.go
internal/stock/infrastructure/persistent/builder/stock.go
internal/stock/infrastructure/persistent/mysql.go
~~~

## Full Diff

~~~diff
diff --git a/internal/common/logging/mysql.go b/internal/common/logging/mysql.go
index 10159f7..f2401bf 100644
--- a/internal/common/logging/mysql.go
+++ b/internal/common/logging/mysql.go
@@ -2,10 +2,10 @@ package logging
 
 import (
 	"context"
-	"encoding/json"
 	"strings"
 	"time"
 
+	"github.com/ghost-yu/go_shop_second/common/util"
 	"github.com/sirupsen/logrus"
 )
 
@@ -17,6 +17,10 @@ const (
 	Error    = "err"
 )
 
+type ArgFormatter interface {
+	FormatArg() (string, error)
+}
+
 func WhenMySQL(ctx context.Context, method string, args ...any) (logrus.Fields, func(any, *error)) {
 	fields := logrus.Fields{
 		Method: method,
@@ -46,12 +50,20 @@ func formatMySQLArgs(args []any) string {
 }
 
 func formatMySQLArg(arg any) string {
-	switch v := arg.(type) {
-	default:
-		bytes, err := json.Marshal(v)
+	var (
+		str string
+		err error
+	)
+	defer func() {
 		if err != nil {
-			return "unsupported type in formatMySQLArg||err=" + err.Error()
+			str = "unsupported type in formatMySQLArg||err=" + err.Error()
 		}
-		return string(bytes)
+	}()
+	switch v := arg.(type) {
+	default:
+		str, err = util.MarshalString(v)
+	case ArgFormatter:
+		str, err = v.FormatArg()
 	}
+	return str
 }
diff --git a/internal/common/util/jsonutil.go b/internal/common/util/jsonutil.go
new file mode 100644
index 0000000..440ba98
--- /dev/null
+++ b/internal/common/util/jsonutil.go
@@ -0,0 +1,8 @@
+package util
+
+import "encoding/json"
+
+func MarshalString(v any) (string, error) {
+	bytes, err := json.Marshal(v)
+	return string(bytes), err
+}
diff --git a/internal/stock/adapters/stock_mysql_repository.go b/internal/stock/adapters/stock_mysql_repository.go
index 04108ff..867911b 100644
--- a/internal/stock/adapters/stock_mysql_repository.go
+++ b/internal/stock/adapters/stock_mysql_repository.go
@@ -25,8 +25,7 @@ func (m MySQLStockRepository) GetItems(ctx context.Context, ids []string) ([]*en
 }
 
 func (m MySQLStockRepository) GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error) {
-	query := builder.NewStock().ProductIDs(ids...)
-	data, err := m.db.BatchGetStockByID(ctx, query)
+	data, err := m.db.BatchGetStockByID(ctx, builder.NewStock().ProductIDs(ids...))
 	if err != nil {
 		return nil, errors.Wrap(err, "BatchGetStockByID error")
 	}
@@ -67,25 +66,20 @@ func (m MySQLStockRepository) updateOptimistic(
 	data []*entity.ItemWithQuantity,
 	updateFn func(ctx context.Context, existing []*entity.ItemWithQuantity, query []*entity.ItemWithQuantity,
 	) ([]*entity.ItemWithQuantity, error)) error {
-	var dest []*persistent.StockModel
-
-	if err := builder.NewStock().ProductIDs(getIDFromEntities(data)...).
-		Fill(tx.Model(&persistent.StockModel{})).Find(&dest).Error; err != nil {
-		return errors.Wrap(err, "failed to find data")
-	}
-
 	for _, queryData := range data {
-		var newestRecord persistent.StockModel
-		if err := builder.NewStock().ProductIDs(queryData.ID).
-			Fill(tx.Model(&persistent.StockModel{})).First(&newestRecord).Error; err != nil {
+		var newestRecord *persistent.StockModel
+		newestRecord, err := m.db.GetStockByID(ctx, builder.NewStock().ProductIDs(queryData.ID))
+		if err != nil {
 			return err
 		}
-
-		if err := builder.NewStock().ProductIDs(queryData.ID).Versions(newestRecord.Version).QuantityGT(queryData.Quantity).
-			Fill(tx.Model(&persistent.StockModel{})).Updates(map[string]any{
-			"quantity": gorm.Expr("quantity - ?", queryData.Quantity),
-			"version":  newestRecord.Version + 1,
-		}).Error; err != nil {
+		if err = m.db.Update(
+			ctx,
+			tx,
+			builder.NewStock().ProductIDs(queryData.ID).Versions(newestRecord.Version).QuantityGT(queryData.Quantity),
+			map[string]any{
+				"quantity": gorm.Expr("quantity - ?", queryData.Quantity),
+				"version":  newestRecord.Version + 1,
+			}); err != nil {
 			return err
 		}
 	}
@@ -93,7 +87,7 @@ func (m MySQLStockRepository) updateOptimistic(
 	return nil
 }
 
-func (m MySQLStockRepository) unmarshalFromDatabase(dest []*persistent.StockModel) []*entity.ItemWithQuantity {
+func (m MySQLStockRepository) unmarshalFromDatabase(dest []persistent.StockModel) []*entity.ItemWithQuantity {
 	var result []*entity.ItemWithQuantity
 	for _, i := range dest {
 		result = append(result, &entity.ItemWithQuantity{
@@ -110,10 +104,9 @@ func (m MySQLStockRepository) updatePessimistic(
 	data []*entity.ItemWithQuantity,
 	updateFn func(ctx context.Context, existing []*entity.ItemWithQuantity, query []*entity.ItemWithQuantity,
 	) ([]*entity.ItemWithQuantity, error)) error {
-	var dest []*persistent.StockModel
-	if err := builder.NewStock().ProductIDs(getIDFromEntities(data)...).ForUpdate().
-		Fill(tx.Model(&persistent.StockModel{})).Find(&dest).Error; err != nil {
-
+	var dest []persistent.StockModel
+	dest, err := m.db.BatchGetStockByID(ctx, builder.NewStock().ProductIDs(getIDFromEntities(data)...).ForUpdate())
+	if err != nil {
 		return errors.Wrap(err, "failed to find data")
 	}
 
@@ -125,15 +118,14 @@ func (m MySQLStockRepository) updatePessimistic(
 
 	for _, upd := range updated {
 		for _, query := range data {
-			if upd.ID == query.ID {
-				if err = builder.NewStock().ProductIDs(upd.ID).QuantityGT(query.Quantity).
-					Fill(tx.Model(&persistent.StockModel{})).
-					Update("quantity", gorm.Expr("quantity - ?", query.Quantity)).Error; err != nil {
-					return errors.Wrapf(err, "unable to update %s", upd.ID)
-				}
+			if upd.ID != query.ID {
+				continue
+			}
+			if err = m.db.Update(ctx, tx, builder.NewStock().ProductIDs(upd.ID).QuantityGT(query.Quantity),
+				map[string]any{"quantity": gorm.Expr("quantity - ?", query.Quantity)}); err != nil {
+				return errors.Wrapf(err, "unable to update %s", upd.ID)
 			}
 		}
-
 	}
 	return nil
 }
diff --git a/internal/stock/adapters/stock_mysql_repository_test.go b/internal/stock/adapters/stock_mysql_repository_test.go
index ab5b49c..fa77ed0 100644
--- a/internal/stock/adapters/stock_mysql_repository_test.go
+++ b/internal/stock/adapters/stock_mysql_repository_test.go
@@ -12,6 +12,7 @@ import (
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent"
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent/builder"
 	"github.com/spf13/viper"
+	gormlogger "gorm.io/gorm/logger"
 
 	"github.com/stretchr/testify/assert"
 	"gorm.io/driver/mysql"
@@ -41,7 +42,9 @@ func setupTestDB(t *testing.T) *persistent.MySQL {
 		viper.GetString("mysql.port"),
 		testDB,
 	)
-	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
+	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
+		Logger: gormlogger.Default.LogMode(gormlogger.Info),
+	})
 	assert.NoError(t, err)
 	assert.NoError(t, db.AutoMigrate(&persistent.StockModel{}))
 
@@ -58,7 +61,7 @@ func TestMySQLStockRepository_UpdateStock_Race(t *testing.T) {
 		testItem           = "item-1"
 		initialStock int32 = 100
 	)
-	err := db.Create(ctx, &persistent.StockModel{
+	err := db.Create(ctx, nil, &persistent.StockModel{
 		ProductID: testItem,
 		Quantity:  initialStock,
 	})
@@ -114,7 +117,7 @@ func TestMySQLStockRepository_UpdateStock_OverSell(t *testing.T) {
 		testItem           = "item-1"
 		initialStock int32 = 5
 	)
-	err := db.Create(ctx, &persistent.StockModel{
+	err := db.Create(ctx, nil, &persistent.StockModel{
 		ProductID: testItem,
 		Quantity:  initialStock,
 	})
diff --git a/internal/stock/infrastructure/persistent/builder/stock.go b/internal/stock/infrastructure/persistent/builder/stock.go
index 81a6a72..7a1b418 100644
--- a/internal/stock/infrastructure/persistent/builder/stock.go
+++ b/internal/stock/infrastructure/persistent/builder/stock.go
@@ -1,84 +1,88 @@
 package builder
 
 import (
+	"github.com/ghost-yu/go_shop_second/common/util"
 	"gorm.io/gorm"
 	"gorm.io/gorm/clause"
 )
 
 type Stock struct {
-	id        []int64
-	productID []string
-	quantity  []int32
-	version   []int64
+	ID        []int64  `json:"ID,omitempty"`
+	ProductID []string `json:"product_id,omitempty"`
+	Quantity  []int32  `json:"quantity,omitempty"`
+	Version   []int64  `json:"version,omitempty"`
 
 	// extend fields
-	order     string
-	forUpdate bool
+	OrderBy       string `json:"order_by,omitempty"`
+	ForUpdateLock bool   `json:"for_update,omitempty"`
 }
 
 func NewStock() *Stock {
 	return &Stock{}
 }
 
+func (s *Stock) FormatArg() (string, error) {
+	return util.MarshalString(s)
+}
+
 func (s *Stock) Fill(db *gorm.DB) *gorm.DB {
 	db = s.fillWhere(db)
-	if s.order != "" {
-		db = db.Order(s.order)
+	if s.OrderBy != "" {
+		db = db.Order(s.Order)
 	}
 	return db
 }
 
 func (s *Stock) fillWhere(db *gorm.DB) *gorm.DB {
-	if len(s.id) > 0 {
-		db = db.Where("id in (?)", s.id)
+	if len(s.ID) > 0 {
+		db = db.Where("ID in (?)", s.ID)
 	}
-	if len(s.productID) > 0 {
-		db = db.Where("product_id in (?)", s.productID)
+	if len(s.ProductID) > 0 {
+		db = db.Where("product_id in (?)", s.ProductID)
 	}
-	if len(s.version) > 0 {
-		db = db.Where("version in (?)", s.version)
+	if len(s.Version) > 0 {
+		db = db.Where("Version in (?)", s.Version)
 	}
-	if len(s.quantity) > 0 {
+	if len(s.Quantity) > 0 {
 		db = s.fillQuantityGT(db)
 	}
-
-	if s.forUpdate {
+	if s.ForUpdateLock {
 		db = db.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate})
 	}
 	return db
 }
 
 func (s *Stock) fillQuantityGT(db *gorm.DB) *gorm.DB {
-	db = db.Where("quantity >= ?", s.quantity)
+	db = db.Where("Quantity >= ?", s.Quantity)
 	return db
 }
 
 func (s *Stock) IDs(v ...int64) *Stock {
-	s.id = v
+	s.ID = v
 	return s
 }
 
 func (s *Stock) ProductIDs(v ...string) *Stock {
-	s.productID = v
+	s.ProductID = v
 	return s
 }
 
 func (s *Stock) Order(v string) *Stock {
-	s.order = v
+	s.OrderBy = v
 	return s
 }
 
 func (s *Stock) Versions(v ...int64) *Stock {
-	s.version = v
+	s.Version = v
 	return s
 }
 
 func (s *Stock) QuantityGT(v ...int32) *Stock {
-	s.quantity = v
+	s.Quantity = v
 	return s
 }
 
 func (s *Stock) ForUpdate() *Stock {
-	s.forUpdate = true
+	s.ForUpdateLock = true
 	return s
 }
diff --git a/internal/stock/infrastructure/persistent/mysql.go b/internal/stock/infrastructure/persistent/mysql.go
index 1502d09..dcd4ec6 100644
--- a/internal/stock/infrastructure/persistent/mysql.go
+++ b/internal/stock/infrastructure/persistent/mysql.go
@@ -31,7 +31,7 @@ func NewMySQL() *MySQL {
 	if err != nil {
 		logrus.Panicf("connect to mysql failed, err=%v", err)
 	}
-	//db.Callback().Create().Before("gorm:create").Register("set_create_time", func(d *gorm.DB) {
+	//db.Callback().Create().Before("gorm:create").Register("set_create_time", func(d *gorm.UseTransaction) {
 	//	d.Statement.SetColumn("CreatedAt", time.Now().Format(time.DateTime))
 	//})
 	return &MySQL{db: db}
@@ -59,14 +59,32 @@ func (m *StockModel) BeforeCreate(tx *gorm.DB) (err error) {
 	return nil
 }
 
+func (d *MySQL) UseTransaction(tx *gorm.DB) *gorm.DB {
+	if tx == nil {
+		return d.db
+	}
+	return tx
+}
+
 func (d MySQL) StartTransaction(fc func(tx *gorm.DB) error) error {
 	return d.db.Transaction(fc)
 }
 
+func (d MySQL) GetStockByID(ctx context.Context, query *builder.Stock) (*StockModel, error) {
+	_, deferLog := logging.WhenMySQL(ctx, "GetStockByID", query)
+	var result StockModel
+	tx := query.Fill(d.db.WithContext(ctx)).First(&result)
+	defer deferLog(result, &tx.Error)
+	if tx.Error != nil {
+		return nil, tx.Error
+	}
+	return &result, nil
+}
+
 func (d MySQL) BatchGetStockByID(ctx context.Context, query *builder.Stock) ([]StockModel, error) {
 	_, deferLog := logging.WhenMySQL(ctx, "BatchGetStockByID", query)
 	var result []StockModel
-	tx := query.Fill(d.db.WithContext(ctx).Clauses(clause.Returning{}).Find(&result))
+	tx := query.Fill(d.db.WithContext(ctx)).Find(&result)
 	defer deferLog(result, &tx.Error)
 	if tx.Error != nil {
 		return nil, tx.Error
@@ -74,10 +92,18 @@ func (d MySQL) BatchGetStockByID(ctx context.Context, query *builder.Stock) ([]S
 	return result, nil
 }
 
-func (d MySQL) Create(ctx context.Context, create *StockModel) error {
+func (d MySQL) Update(ctx context.Context, tx *gorm.DB, cond *builder.Stock, update map[string]any) error {
+	_, deferLog := logging.WhenMySQL(ctx, "BatchUpdateStock", cond)
+	var returning StockModel
+	res := cond.Fill(d.UseTransaction(tx).WithContext(ctx).Model(&returning).Clauses(clause.Returning{})).Updates(update)
+	defer deferLog(returning, &res.Error)
+	return res.Error
+}
+
+func (d MySQL) Create(ctx context.Context, tx *gorm.DB, create *StockModel) error {
 	_, deferLog := logging.WhenMySQL(ctx, "Create", create)
 	var returning StockModel
-	err := d.db.WithContext(ctx).Model(&returning).Clauses(clause.Returning{}).Create(create).Error
+	err := d.UseTransaction(tx).WithContext(ctx).Model(&returning).Clauses(clause.Returning{}).Create(create).Error
 	defer deferLog(returning, &err)
 	return err
 }
~~~
