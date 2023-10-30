package conf

import (
	"flag"
	"strings"

	"github.com/mangohow/cloud-ide/pkg/conf"
	"github.com/spf13/viper"
)

var (
	ServerConfig conf.ServerConf
	MysqlConfig  conf.MysqlConf
	RedisConfig  conf.RedisConf
	LoggerConfig conf.LoggerConf
	GrpcConfig   conf.GrpcConf
	EmailConfig  conf.EmailConf
)

func LoadConf() error {
	viper.SetConfigName("webserver")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	initServerConf()
	initMysqlConf()
	initRedisConf()
	initLogConf()
	initGrpcConf()
	initEmailConf()

	parseFlags()

	return nil
}

func initServerConf() {
	ServerConfig = conf.ServerConf{
		Host: viper.GetString("server.host"),
		Port: viper.GetInt("server.port"),
		Name: viper.GetString("server.name"),
		Mode: viper.GetString("server.mode"),
	}
}

func initMysqlConf() {
	MysqlConfig = conf.MysqlConf{
		DataSourceName: viper.GetString("mysql.dataSourceName"),
		MaxOpenConns:   viper.GetInt("mysql.maxOpenConns"),
		MaxIdleConns:   viper.GetInt("mysql.maxIdleConns"),
	}
}

func initLogConf() {
	LoggerConfig = conf.LoggerConf{
		Level:       viper.GetString("logger.level"),
		FilePath:    viper.GetString("logger.filePath"),
		FileName:    viper.GetString("logger.fileName"),
		MaxFileSize: viper.GetUint64("logger.maxFileSize"),
		ToFile:      viper.GetBool("logger.toFile"),
	}
}

func initRedisConf() {
	RedisConfig = conf.RedisConf{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetUint32("redis.db"),
		PoolSize:     viper.GetUint32("redis.poolSize"),
		MinIdleConns: viper.GetUint32("redis.minIdleConns"),
	}
}

func initGrpcConf() {
	GrpcConfig = conf.GrpcConf{Addr: viper.GetString("grpc.addr")}
}

func initEmailConf() {
	EmailConfig = conf.EmailConf{
		Enabled:     viper.GetBool("email.enabled"),
		Host:        viper.GetString("email.host"),
		Port:        viper.GetUint32("email.port"),
		SenderEmail: viper.GetString("email.senderEmail"),
		AuthCode:    viper.GetString("email.authCode"),
	}
}

// 解析命令行参数
func parseFlags() {
	var (
		mode           string
		port           int
		dataSourceName string
		logLevel       string
		email          string
		emailHost      string
		emailPort      int
		senderEmail    string
		authCode       string
		grpcAddr       string
	)

	flag.StringVar(&mode, "mode", "", "specify server running mode [dev, release]")
	flag.IntVar(&port, "port", -1, "specify server listen port")
	flag.StringVar(&dataSourceName, "mysql-datasource", "", "specify mysql datasource eg. user:password@(svc:port)/database?charset=utf8mb4&parseTime=true&loc=Local")
	flag.StringVar(&logLevel, "log-level", "", "specify log level [debug, info, warn, error]")
	flag.StringVar(&email, "email-enabled", "", "enable email register [enabled, disabled]")
	flag.StringVar(&emailHost, "email-host", "", "specify email host if email is enabled")
	flag.IntVar(&emailPort, "email-port", -1, "specify email port if email is enabled")
	flag.StringVar(&senderEmail, "email-sender", "", "specify sender email if email is enabled")
	flag.StringVar(&authCode, "email-authcode", "", "specify email auth code if email is enabled")
	flag.StringVar(&grpcAddr, "grpc-addr", "", "specify control plane grpc addr eg:cloud-ide-control-plane-svc:6387")
	flag.Parse()

	setString(&ServerConfig.Mode, &mode)
	setString(&EmailConfig.SenderEmail, &senderEmail)
	setString(&EmailConfig.AuthCode, &authCode)
	setString(&GrpcConfig.Addr, &grpcAddr)
	setString(&MysqlConfig.DataSourceName, &dataSourceName)
	setString(&LoggerConfig.Level, &logLevel)
	setString(&EmailConfig.Host, &emailHost)
	if port != -1 {
		ServerConfig.Port = port
	}
	if emailPort != -1 {
		EmailConfig.Port = uint32(emailPort)
	}
	switch strings.ToLower(email) {
	case "enabled":
		EmailConfig.Enabled = true
	case "disabled":
		EmailConfig.Enabled = false
	}
}

func setString(dst *string, src *string) {
	if *src != "" {
		*dst = *src
	}
}
