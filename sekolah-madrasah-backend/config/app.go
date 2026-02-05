package config

import "strings"

type DBConnectionTarget string

const (
	MainDB DBConnectionTarget = "MAIN_DB"
)

type Rest struct {
	Origin          string
	Port            int
	SwaggerUser     string
	SwaggerPassword string
	JWTSecret       string
	XAuthCron       string
	RestDebugMode   bool
}

func (r *Rest) GetOrigin() []string {
	noSpace := strings.ToLower(strings.ReplaceAll(r.Origin, " ", ""))
	commaSeparated := strings.Split(noSpace, ",")
	for i := range commaSeparated {
		commaSeparated[i] = strings.TrimSpace(commaSeparated[i])
	}
	return commaSeparated
}

type DBConfig struct {
	Id       DBConnectionTarget
	SslMode  string
	Host     string
	Port     int
	Name     string
	Username string
	Password string
	DbEvent  string
}

type APMConfig struct {
	Enabled          bool
	ServiceName      string
	ServerURL        string
	SecretToken      string
	VerifyServerCert bool
}

// type MsgBrokerConfig struct {
// 	Host     string
// 	Port     string
// 	User     string
// 	Password string
// 	ClientId string
// }

// type BucketConfig struct {
// 	BaseUrl         string
// 	BucketName      string
// 	AccessKeyId     string
// 	SecretAccessKey string
// 	PrefixPath      string
// 	Region          string
// }

type ElasticsearchConfig struct {
	Enabled     bool
	Addresses   []string
	Username    string
	Password    string
	APIKey      string
	CloudID     string
	IndexPrefix string
	CertPath    string
}

type AppConfig struct {
	ServiceName string
	Version     string
	Rest        Rest
	MainDB      DBConfig
	APM         APMConfig
	// Broker        MsgBrokerConfig
	// Bucket        BucketConfig
	Elasticsearch ElasticsearchConfig
	Security      Security
}
type Security struct {
	AesKey string
}

var APP AppConfig
