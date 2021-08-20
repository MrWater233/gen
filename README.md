# 使用教程

1. 安装Go并通过以下命令安装Gen

   ```shell
   go get -u github.com/MrWater233/gen
   ```

2. 在代码中导入

   ```go
   import "github.com/MrWater233/gen"
   ```

3. 简单的使用

   ```go
   package main
   
   import "github.com/MrWater233/gen"
   
   func main() {
   	r := gen.New()
   	// 加载日志中间件
   	r.Use(gen.Logger())
   	r.GET("/", func(c *gen.Context) {
   		c.HTML(200, "<h1>Hello Gen!</h1>")
   	})
   	// 开启服务，监听8000端口
   	r.Run(":8000")
   }
   ```