# Lesson Pair Diff Report

- FromBranch: lesson19
- ToBranch: lesson20

## Short Summary

~~~text
 3 files changed, 165 insertions(+), 4 deletions(-)
~~~

## File Stats

~~~text
 internal/order/http.go |  15 +++--
 internal/order/main.go |   1 +
 public/success.html    | 153 +++++++++++++++++++++++++++++++++++++++++++++++++
 3 files changed, 165 insertions(+), 4 deletions(-)
~~~

## Commit Comparison

~~~text
> 2df9032 html
~~~

## Changed Files

~~~text
internal/order/http.go
internal/order/main.go
public/success.html
~~~

## Focus Files (Excluded: go.mod/go.sum, *.pb.go, *.gen.go)

~~~text
internal/order/http.go
internal/order/main.go
public/success.html
~~~

## Full Diff

~~~diff
diff --git a/internal/order/http.go b/internal/order/http.go
index b40adc7..2073a4f 100644
--- a/internal/order/http.go
+++ b/internal/order/http.go
@@ -1,6 +1,7 @@
 package main
 
 import (
+	"fmt"
 	"net/http"
 
 	"github.com/ghost-yu/go_shop_second/common/genproto/orderpb"
@@ -29,9 +30,10 @@ func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID stri
 		return
 	}
 	c.JSON(http.StatusOK, gin.H{
-		"message":     "success",
-		"customer_id": req.CustomerID,
-		"order_id":    r.OrderID,
+		"message":      "success",
+		"customer_id":  req.CustomerID,
+		"order_id":     r.OrderID,
+		"redirect_url": fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID),
 	})
 }
 
@@ -44,5 +46,10 @@ func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerI
 		c.JSON(http.StatusOK, gin.H{"error": err})
 		return
 	}
-	c.JSON(http.StatusOK, gin.H{"message": "success", "data": o})
+	c.JSON(http.StatusOK, gin.H{
+		"message": "success",
+		"data": gin.H{
+			"Order": o,
+		},
+	})
 }
diff --git a/internal/order/main.go b/internal/order/main.go
index 96f2ad3..b35cdc8 100644
--- a/internal/order/main.go
+++ b/internal/order/main.go
@@ -46,6 +46,7 @@ func main() {
 	})
 
 	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
+		router.StaticFile("/success", "../../public/success.html")
 		ports.RegisterHandlersWithOptions(router, HTTPServer{
 			app: application,
 		}, ports.GinServerOptions{
diff --git a/public/success.html b/public/success.html
new file mode 100644
index 0000000..08abb53
--- /dev/null
+++ b/public/success.html
@@ -0,0 +1,153 @@
+<!DOCTYPE html>
+<html lang="en">
+<head>
+    <meta charset="UTF-8">
+    <title>Gorder</title>
+</head>
+<body>
+<section>
+  <p>
+    您已成功下单！
+  </p>
+  <p>
+    订单状态：<span id="orderStatus">等待中...</span>
+  </p>
+  <div class="ready-popup">
+    <p>您的订单正在处理中...</p>
+    <p style="color:burlywood; margin:12px">
+      订单号：<b><span id="orderID"></span></b>
+    </p>
+
+    <button class="close-btn" onclick="document.querySelector('.ready-popup').style.display = 'none'">
+      关闭
+    </button>
+  </div>
+
+  <div class="after-payment-popup">
+    <p>等待支付中...</p>
+    <a id="payment-link" href="#">去支付</a>
+  </div>
+</section>
+</body>
+
+<style>
+  html {
+      margin: 0;
+      padding: 0;
+      background-color: antiquewhite;
+      color: darkblue;
+  }
+
+  section {
+      position: relative;
+      display: flex;
+      flex-direction: column;
+      justify-content: center;
+      align-items: center;
+      height: 100vh;
+  }
+
+  .ready-popup {
+      display: none;
+      flex-direction: column;
+      justify-content: center;
+      align-items: center;
+      position: fixed;
+      top: 50%;
+      left: 50%;
+      transform: translate(-50%, -50%);
+      padding: 20px;
+      background-color: cadetblue;
+      z-index: 1;
+      border: 2px solid black;
+      border-radius: 5px;
+  }
+
+  .ready-popup p {
+      margin: 0;
+  }
+
+  .after-payment-popup {
+      display: none;
+      flex-direction: column;
+      justify-content: center;
+      align-items: center;
+      position: fixed;
+      top: 50%;
+      left: 50%;
+      transform: translate(-50%, -50%);
+      padding: 20px;
+      background-color: cadetblue;
+      z-index: 1;
+      border: 2px solid black;
+      border-radius: 5px;
+  }
+
+  .after-payment-popup p {
+      margin: 0;
+  }
+  .after-payment-popup a {
+      color: white;
+      margin-top: 10px;
+      padding: 5px 10px;
+      background-color: green;
+      border-radius: 5px;
+      text-decoration: none;
+  }
+
+  .close-btn {
+      margin-top: 10px;
+      padding: 5px 10px;
+      background-color: green;
+      border-radius: 5px;
+      border: none;
+      cursor: pointer;
+  }
+</style>
+
+<script>
+  const urlParam = new URLSearchParams(window.location.search);
+  const customerID = urlParam.get('customerID');
+  const orderID = urlParam.get('orderID');
+  const order = {
+      customerID,
+      orderID,
+      status: 'pending'
+  };
+  const getOrder = async() => {
+      const res = await fetch(`/api/customer/${customerID}/orders/${orderID}`);
+      const data = await res.json();
+      console.log("data = ", data)
+
+      /*
+      {
+        "code": 0,
+        "message": "success",
+        "data": {
+          ...
+        }
+      }
+       */
+      if (data.data.Order.Status === 'waiting_for_payment') {
+          order.Status = '等待支付...';
+          document.getElementById('orderStatus').innerText = order.Status;
+          document.querySelector('.after-payment-popup').style.display = 'block';
+          document.getElementById('payment-link').href = data.data.Order.PaymentLink;
+      }
+      if (data.data.Order.Status === 'paid') {
+          order.Status = '已支付成功，请等待...';
+          document.getElementById('orderStatus').innerText = order.Status;
+          setTimeout(getOrder, 5000);
+      } else if (data.data.Order.Status === 'ready') {
+          order.Status = '已完成...';
+          document.querySelector('.after-payment-popup').style.display = 'none';
+          document.querySelector('.ready-popup').style.display = 'block';
+          document.getElementById('orderID').innerText = orderID;
+          document.getElementById('orderStatus').innerText = order.Status;
+      } else {
+          setTimeout(getOrder, 5000);
+      }
+  }
+  getOrder();
+</script>
+</html>
\ No newline at end of file
~~~
