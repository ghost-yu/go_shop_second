# Lesson Pair Diff Report

- FromBranch: lesson56
- ToBranch: lesson57

## Short Summary

~~~text
 4 files changed, 109 insertions(+), 23 deletions(-)
~~~

## File Stats

~~~text
 internal/stock/adapters/stock_mysql_repository.go  | 36 +++++-----
 .../stock/adapters/stock_mysql_repository_test.go  |  5 +-
 .../infrastructure/persistent/builder/stock.go     | 84 ++++++++++++++++++++++
 internal/stock/infrastructure/persistent/mysql.go  |  7 +-
 4 files changed, 109 insertions(+), 23 deletions(-)
~~~

## Commit Comparison

~~~text
> 5222389 sql builder
~~~

## Changed Files

~~~text
internal/stock/adapters/stock_mysql_repository.go
internal/stock/adapters/stock_mysql_repository_test.go
internal/stock/infrastructure/persistent/builder/stock.go
internal/stock/infrastructure/persistent/mysql.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/stock/adapters/stock_mysql_repository.go
internal/stock/adapters/stock_mysql_repository_test.go
internal/stock/infrastructure/persistent/builder/stock.go
internal/stock/infrastructure/persistent/mysql.go
~~~

## Full Diff

~~~diff
diff --git a/internal/stock/adapters/stock_mysql_repository.go b/internal/stock/adapters/stock_mysql_repository.go
index a4aa7df..04108ff 100644
--- a/internal/stock/adapters/stock_mysql_repository.go
+++ b/internal/stock/adapters/stock_mysql_repository.go
@@ -5,10 +5,10 @@ import (
 
 	"github.com/ghost-yu/go_shop_second/stock/entity"
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent"
+	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent/builder"
 	"github.com/pkg/errors"
 	"github.com/sirupsen/logrus"
 	"gorm.io/gorm"
-	"gorm.io/gorm/clause"
 )
 
 type MySQLStockRepository struct {
@@ -25,7 +25,8 @@ func (m MySQLStockRepository) GetItems(ctx context.Context, ids []string) ([]*en
 }
 
 func (m MySQLStockRepository) GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error) {
-	data, err := m.db.BatchGetStockByID(ctx, ids)
+	query := builder.NewStock().ProductIDs(ids...)
+	data, err := m.db.BatchGetStockByID(ctx, query)
 	if err != nil {
 		return nil, errors.Wrap(err, "BatchGetStockByID error")
 	}
@@ -67,24 +68,24 @@ func (m MySQLStockRepository) updateOptimistic(
 	updateFn func(ctx context.Context, existing []*entity.ItemWithQuantity, query []*entity.ItemWithQuantity,
 	) ([]*entity.ItemWithQuantity, error)) error {
 	var dest []*persistent.StockModel
-	if err := tx.Model(&persistent.StockModel{}).
-		Where("product_id IN (?)", getIDFromEntities(data)).
-		Find(&dest).Error; err != nil {
+
+	if err := builder.NewStock().ProductIDs(getIDFromEntities(data)...).
+		Fill(tx.Model(&persistent.StockModel{})).Find(&dest).Error; err != nil {
 		return errors.Wrap(err, "failed to find data")
 	}
 
 	for _, queryData := range data {
 		var newestRecord persistent.StockModel
-		if err := tx.Model(&persistent.StockModel{}).Where("product_id = ?", queryData.ID).
-			First(&newestRecord).Error; err != nil {
+		if err := builder.NewStock().ProductIDs(queryData.ID).
+			Fill(tx.Model(&persistent.StockModel{})).First(&newestRecord).Error; err != nil {
 			return err
 		}
-		if err := tx.Model(&persistent.StockModel{}).
-			Where("product_id = ? AND version = ? AND quantity - ? >= 0", queryData.ID, newestRecord.Version, queryData.Quantity).
-			Updates(map[string]any{
-				"quantity": gorm.Expr("quantity - ?", queryData.Quantity),
-				"version":  newestRecord.Version + 1,
-			}).Error; err != nil {
+
+		if err := builder.NewStock().ProductIDs(queryData.ID).Versions(newestRecord.Version).QuantityGT(queryData.Quantity).
+			Fill(tx.Model(&persistent.StockModel{})).Updates(map[string]any{
+			"quantity": gorm.Expr("quantity - ?", queryData.Quantity),
+			"version":  newestRecord.Version + 1,
+		}).Error; err != nil {
 			return err
 		}
 	}
@@ -110,10 +111,8 @@ func (m MySQLStockRepository) updatePessimistic(
 	updateFn func(ctx context.Context, existing []*entity.ItemWithQuantity, query []*entity.ItemWithQuantity,
 	) ([]*entity.ItemWithQuantity, error)) error {
 	var dest []*persistent.StockModel
-	if err := tx.Table("o_stock").
-		Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
-		Where("product_id IN ?", getIDFromEntities(data)).
-		Find(&dest).Error; err != nil {
+	if err := builder.NewStock().ProductIDs(getIDFromEntities(data)...).ForUpdate().
+		Fill(tx.Model(&persistent.StockModel{})).Find(&dest).Error; err != nil {
 
 		return errors.Wrap(err, "failed to find data")
 	}
@@ -127,7 +126,8 @@ func (m MySQLStockRepository) updatePessimistic(
 	for _, upd := range updated {
 		for _, query := range data {
 			if upd.ID == query.ID {
-				if err = tx.Table("o_stock").Where("product_id = ? AND quantity - ? >= 0", upd.ID, query.Quantity).
+				if err = builder.NewStock().ProductIDs(upd.ID).QuantityGT(query.Quantity).
+					Fill(tx.Model(&persistent.StockModel{})).
 					Update("quantity", gorm.Expr("quantity - ?", query.Quantity)).Error; err != nil {
 					return errors.Wrapf(err, "unable to update %s", upd.ID)
 				}
diff --git a/internal/stock/adapters/stock_mysql_repository_test.go b/internal/stock/adapters/stock_mysql_repository_test.go
index 9c397be..ab5b49c 100644
--- a/internal/stock/adapters/stock_mysql_repository_test.go
+++ b/internal/stock/adapters/stock_mysql_repository_test.go
@@ -10,6 +10,7 @@ import (
 	_ "github.com/ghost-yu/go_shop_second/common/config"
 	"github.com/ghost-yu/go_shop_second/stock/entity"
 	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent"
+	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent/builder"
 	"github.com/spf13/viper"
 
 	"github.com/stretchr/testify/assert"
@@ -95,7 +96,7 @@ func TestMySQLStockRepository_UpdateStock_Race(t *testing.T) {
 	}
 
 	wg.Wait()
-	res, err := db.BatchGetStockByID(ctx, []string{testItem})
+	res, err := db.BatchGetStockByID(ctx, builder.NewStock().ProductIDs(testItem))
 	assert.NoError(t, err)
 	assert.NotEmpty(t, res, "res cannot be empty")
 
@@ -152,7 +153,7 @@ func TestMySQLStockRepository_UpdateStock_OverSell(t *testing.T) {
 	}
 
 	wg.Wait()
-	res, err := db.BatchGetStockByID(ctx, []string{testItem})
+	res, err := db.BatchGetStockByID(ctx, builder.NewStock().ProductIDs(testItem))
 	assert.NoError(t, err)
 	assert.NotEmpty(t, res, "res cannot be empty")
 
diff --git a/internal/stock/infrastructure/persistent/builder/stock.go b/internal/stock/infrastructure/persistent/builder/stock.go
new file mode 100644
index 0000000..81a6a72
--- /dev/null
+++ b/internal/stock/infrastructure/persistent/builder/stock.go
@@ -0,0 +1,84 @@
+package builder
+
+import (
+	"gorm.io/gorm"
+	"gorm.io/gorm/clause"
+)
+
+type Stock struct {
+	id        []int64
+	productID []string
+	quantity  []int32
+	version   []int64
+
+	// extend fields
+	order     string
+	forUpdate bool
+}
+
+func NewStock() *Stock {
+	return &Stock{}
+}
+
+func (s *Stock) Fill(db *gorm.DB) *gorm.DB {
+	db = s.fillWhere(db)
+	if s.order != "" {
+		db = db.Order(s.order)
+	}
+	return db
+}
+
+func (s *Stock) fillWhere(db *gorm.DB) *gorm.DB {
+	if len(s.id) > 0 {
+		db = db.Where("id in (?)", s.id)
+	}
+	if len(s.productID) > 0 {
+		db = db.Where("product_id in (?)", s.productID)
+	}
+	if len(s.version) > 0 {
+		db = db.Where("version in (?)", s.version)
+	}
+	if len(s.quantity) > 0 {
+		db = s.fillQuantityGT(db)
+	}
+
+	if s.forUpdate {
+		db = db.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate})
+	}
+	return db
+}
+
+func (s *Stock) fillQuantityGT(db *gorm.DB) *gorm.DB {
+	db = db.Where("quantity >= ?", s.quantity)
+	return db
+}
+
+func (s *Stock) IDs(v ...int64) *Stock {
+	s.id = v
+	return s
+}
+
+func (s *Stock) ProductIDs(v ...string) *Stock {
+	s.productID = v
+	return s
+}
+
+func (s *Stock) Order(v string) *Stock {
+	s.order = v
+	return s
+}
+
+func (s *Stock) Versions(v ...int64) *Stock {
+	s.version = v
+	return s
+}
+
+func (s *Stock) QuantityGT(v ...int32) *Stock {
+	s.quantity = v
+	return s
+}
+
+func (s *Stock) ForUpdate() *Stock {
+	s.forUpdate = true
+	return s
+}
diff --git a/internal/stock/infrastructure/persistent/mysql.go b/internal/stock/infrastructure/persistent/mysql.go
index a6d638d..1502d09 100644
--- a/internal/stock/infrastructure/persistent/mysql.go
+++ b/internal/stock/infrastructure/persistent/mysql.go
@@ -6,6 +6,7 @@ import (
 	"time"
 
 	"github.com/ghost-yu/go_shop_second/common/logging"
+	"github.com/ghost-yu/go_shop_second/stock/infrastructure/persistent/builder"
 	"github.com/sirupsen/logrus"
 	"github.com/spf13/viper"
 	"gorm.io/driver/mysql"
@@ -62,10 +63,10 @@ func (d MySQL) StartTransaction(fc func(tx *gorm.DB) error) error {
 	return d.db.Transaction(fc)
 }
 
-func (d MySQL) BatchGetStockByID(ctx context.Context, productIDs []string) ([]StockModel, error) {
-	_, deferLog := logging.WhenMySQL(ctx, "BatchGetStockByID", productIDs)
+func (d MySQL) BatchGetStockByID(ctx context.Context, query *builder.Stock) ([]StockModel, error) {
+	_, deferLog := logging.WhenMySQL(ctx, "BatchGetStockByID", query)
 	var result []StockModel
-	tx := d.db.WithContext(ctx).Clauses(clause.Returning{}).Where("product_id IN ?", productIDs).Find(&result)
+	tx := query.Fill(d.db.WithContext(ctx).Clauses(clause.Returning{}).Find(&result))
 	defer deferLog(result, &tx.Error)
 	if tx.Error != nil {
 		return nil, tx.Error
~~~
