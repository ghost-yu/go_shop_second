# Lesson Pair Diff Report

- FromBranch: lesson54
- ToBranch: lesson55

## Short Summary

~~~text
 4 files changed, 128 insertions(+), 23 deletions(-)
~~~

## File Stats

~~~text
 init.sql                                           |  5 +-
 internal/stock/adapters/stock_mysql_repository.go  | 88 ++++++++++++++++------
 .../stock/adapters/stock_mysql_repository_test.go  | 57 ++++++++++++++
 internal/stock/infrastructure/persistent/mysql.go  |  1 +
 4 files changed, 128 insertions(+), 23 deletions(-)
~~~

## Commit Comparison

~~~text
> bdbdfa8 optimistic
~~~

## Changed Files

~~~text
init.sql
internal/stock/adapters/stock_mysql_repository.go
internal/stock/adapters/stock_mysql_repository_test.go
internal/stock/infrastructure/persistent/mysql.go
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
init.sql
internal/stock/adapters/stock_mysql_repository.go
internal/stock/adapters/stock_mysql_repository_test.go
internal/stock/infrastructure/persistent/mysql.go
~~~

## Full Diff

~~~diff
diff --git a/init.sql b/init.sql
index 5585bc0..ce4d6c4 100644
--- a/init.sql
+++ b/init.sql
@@ -7,9 +7,10 @@ CREATE TABLE `o_stock` (
     id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
     product_id VARCHAR(255) NOT NULL,
     quantity INT UNSIGNED NOT NULL DEFAULT 0,
+    version INT NOT NULL DEFAULT 0,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
 ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
 
-INSERT INTO o_stock (product_id, quantity)
-VALUES ('prod_R3g7MikGYsXKzr', 1000), ('prod_R285C3Wb7FDprc', 500);
\ No newline at end of file
+INSERT INTO o_stock (product_id, quantity, version)
+VALUES ('prod_R3g7MikGYsXKzr', 1000, 0), ('prod_R285C3Wb7FDprc', 500, 0);
\ No newline at end of file
diff --git a/internal/stock/adapters/stock_mysql_repository.go b/internal/stock/adapters/stock_mysql_repository.go
index d328331..a4aa7df 100644
--- a/internal/stock/adapters/stock_mysql_repository.go
+++ b/internal/stock/adapters/stock_mysql_repository.go
@@ -54,31 +54,42 @@ func (m MySQLStockRepository) UpdateStock(
 				logrus.Warnf("update stock transaction err=%v", err)
 			}
 		}()
-		var dest []*persistent.StockModel
-		err = tx.Table("o_stock").
-			Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
-			Where("product_id IN ?", getIDFromEntities(data)).
-			Find(&dest).Error
-		if err != nil {
-			return errors.Wrap(err, "failed to find data")
-		}
-		existing := m.unmarshalFromDatabase(dest)
+		err = m.updatePessimistic(ctx, tx, data, updateFn)
+		//err = m.updateOptimistic(ctx, tx, data, updateFn)
+		return err
+	})
+}
+
+func (m MySQLStockRepository) updateOptimistic(
+	ctx context.Context,
+	tx *gorm.DB,
+	data []*entity.ItemWithQuantity,
+	updateFn func(ctx context.Context, existing []*entity.ItemWithQuantity, query []*entity.ItemWithQuantity,
+	) ([]*entity.ItemWithQuantity, error)) error {
+	var dest []*persistent.StockModel
+	if err := tx.Model(&persistent.StockModel{}).
+		Where("product_id IN (?)", getIDFromEntities(data)).
+		Find(&dest).Error; err != nil {
+		return errors.Wrap(err, "failed to find data")
+	}
 
-		updated, err := updateFn(ctx, existing, data)
-		if err != nil {
+	for _, queryData := range data {
+		var newestRecord persistent.StockModel
+		if err := tx.Model(&persistent.StockModel{}).Where("product_id = ?", queryData.ID).
+			First(&newestRecord).Error; err != nil {
 			return err
 		}
-
-		for _, upd := range updated {
-			if err = tx.Table("o_stock").
-				Where("product_id = ?", upd.ID).
-				Update("quantity", upd.Quantity).
-				Error; err != nil {
-				return errors.Wrapf(err, "unable to update %s", upd.ID)
-			}
+		if err := tx.Model(&persistent.StockModel{}).
+			Where("product_id = ? AND version = ? AND quantity - ? >= 0", queryData.ID, newestRecord.Version, queryData.Quantity).
+			Updates(map[string]any{
+				"quantity": gorm.Expr("quantity - ?", queryData.Quantity),
+				"version":  newestRecord.Version + 1,
+			}).Error; err != nil {
+			return err
 		}
-		return nil
-	})
+	}
+
+	return nil
 }
 
 func (m MySQLStockRepository) unmarshalFromDatabase(dest []*persistent.StockModel) []*entity.ItemWithQuantity {
@@ -92,6 +103,41 @@ func (m MySQLStockRepository) unmarshalFromDatabase(dest []*persistent.StockMode
 	return result
 }
 
+func (m MySQLStockRepository) updatePessimistic(
+	ctx context.Context,
+	tx *gorm.DB,
+	data []*entity.ItemWithQuantity,
+	updateFn func(ctx context.Context, existing []*entity.ItemWithQuantity, query []*entity.ItemWithQuantity,
+	) ([]*entity.ItemWithQuantity, error)) error {
+	var dest []*persistent.StockModel
+	if err := tx.Table("o_stock").
+		Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
+		Where("product_id IN ?", getIDFromEntities(data)).
+		Find(&dest).Error; err != nil {
+
+		return errors.Wrap(err, "failed to find data")
+	}
+
+	existing := m.unmarshalFromDatabase(dest)
+	updated, err := updateFn(ctx, existing, data)
+	if err != nil {
+		return err
+	}
+
+	for _, upd := range updated {
+		for _, query := range data {
+			if upd.ID == query.ID {
+				if err = tx.Table("o_stock").Where("product_id = ? AND quantity - ? >= 0", upd.ID, query.Quantity).
+					Update("quantity", gorm.Expr("quantity - ?", query.Quantity)).Error; err != nil {
+					return errors.Wrapf(err, "unable to update %s", upd.ID)
+				}
+			}
+		}
+
+	}
+	return nil
+}
+
 func getIDFromEntities(items []*entity.ItemWithQuantity) []string {
 	var ids []string
 	for _, i := range items {
diff --git a/internal/stock/adapters/stock_mysql_repository_test.go b/internal/stock/adapters/stock_mysql_repository_test.go
index 6244376..9c397be 100644
--- a/internal/stock/adapters/stock_mysql_repository_test.go
+++ b/internal/stock/adapters/stock_mysql_repository_test.go
@@ -5,6 +5,7 @@ import (
 	"fmt"
 	"sync"
 	"testing"
+	"time"
 
 	_ "github.com/ghost-yu/go_shop_second/common/config"
 	"github.com/ghost-yu/go_shop_second/stock/entity"
@@ -101,3 +102,59 @@ func TestMySQLStockRepository_UpdateStock_Race(t *testing.T) {
 	expectedStock := initialStock - int32(concurrentGoroutines)
 	assert.Equal(t, expectedStock, res[0].Quantity)
 }
+
+func TestMySQLStockRepository_UpdateStock_OverSell(t *testing.T) {
+	t.Parallel()
+	ctx := context.Background()
+	db := setupTestDB(t)
+
+	// 准备初始数据
+	var (
+		testItem           = "item-1"
+		initialStock int32 = 5
+	)
+	err := db.Create(ctx, &persistent.StockModel{
+		ProductID: testItem,
+		Quantity:  initialStock,
+	})
+	assert.NoError(t, err)
+
+	repo := NewMySQLStockRepository(db)
+	var wg sync.WaitGroup
+	concurrentGoroutines := 100
+	for i := 0; i < concurrentGoroutines; i++ {
+		wg.Add(1)
+		go func() {
+			defer wg.Done()
+			err := repo.UpdateStock(
+				ctx,
+				[]*entity.ItemWithQuantity{
+					{ID: testItem, Quantity: 1},
+				}, func(ctx context.Context, existing, query []*entity.ItemWithQuantity) ([]*entity.ItemWithQuantity, error) {
+					// 模拟减少库存
+					var newItems []*entity.ItemWithQuantity
+					for _, e := range existing {
+						for _, q := range query {
+							if e.ID == q.ID {
+								newItems = append(newItems, &entity.ItemWithQuantity{
+									ID:       e.ID,
+									Quantity: e.Quantity - q.Quantity,
+								})
+							}
+						}
+					}
+					return newItems, nil
+				},
+			)
+			assert.NoError(t, err)
+		}()
+		time.Sleep(20 * time.Millisecond)
+	}
+
+	wg.Wait()
+	res, err := db.BatchGetStockByID(ctx, []string{testItem})
+	assert.NoError(t, err)
+	assert.NotEmpty(t, res, "res cannot be empty")
+
+	assert.GreaterOrEqual(t, res[0].Quantity, int32(0))
+}
diff --git a/internal/stock/infrastructure/persistent/mysql.go b/internal/stock/infrastructure/persistent/mysql.go
index 08f1d1f..c6eb728 100644
--- a/internal/stock/infrastructure/persistent/mysql.go
+++ b/internal/stock/infrastructure/persistent/mysql.go
@@ -39,6 +39,7 @@ type StockModel struct {
 	ID        int64     `gorm:"column:id"`
 	ProductID string    `gorm:"column:product_id"`
 	Quantity  int32     `gorm:"column:quantity"`
+	Version   int64     `gorm:"column:version"`
 	CreatedAt time.Time `gorm:"column:created_at autoCreateTime"`
 	UpdateAt  time.Time `gorm:"column:updated_at autoUpdateTime"`
 }
~~~
