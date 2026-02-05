package config

import (
	"fmt"
)

func (p *DBConfig) ToURL() string {
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		p.Host, p.Port, p.Username, p.Password, p.Name,
	)
	if p.SslMode != "" {
		connString = connString + fmt.Sprintf(" sslmode=%s", p.SslMode)
	}
	return connString
}
