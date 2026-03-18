package config

import "github.com/spf13/viper"

// NewViperConfig 统一加载全局配置文件到内存。
// 这个函数在 init() 里被调用，确保 main 启动前配置已经就位，
// 后续各个函数可以直接用 viper.Get... 获取配置，而不用每次都读文件。
func NewViperConfig() error {
	// SetConfigName 指定配置文件名（不包括扩展名）。
	// 这里是 "global"，对应 global.yaml。
	viper.SetConfigName("global")

	// SetConfigType 告诉 viper 配置格式是什么。
	// viper 支持 json/yaml/toml/hcl 等多种格式。
	viper.SetConfigType("yaml")

	// AddConfigPath 告诉 viper 去哪个目录找配置文件。
	// "../common/config" 是相对当前执行目录的路径，
	// 假设你从项目根目录启动 order 服务，viper 会去 internal/common/config 目录找 global.yaml。
	viper.AddConfigPath("../common/config")

	// AutomaticEnv 让环境变量自动覆盖配置文件中的值。
	// 举例：如果环境里设置了 ORDER_GRPC_ADDR=":6666"，会覆盖 yaml 里的 order.grpc-addr。
	// 这样本地开发、测试、生产环境可以动态调整配置，而不改代码。
	viper.AutomaticEnv()

	// ReadInConfig 真正去磁盘读配置文件，把内容解析到 viper 的内存结构里。
	// 如果文件不存在或格式有问题，这里会返回 error。
	return viper.ReadInConfig()
}
