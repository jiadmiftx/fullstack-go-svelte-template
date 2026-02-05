package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"go.elastic.co/apm/v2"

	"sekolah-madrasah/config"
	"sekolah-madrasah/database/schemas"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDBConnection(cfg config.DBConfig) (*gorm.DB, error) {
	apmLogger := &APMGormLogger{
		logger: logger.Default.LogMode(logger.Warn),
	}
	gcfg := &gorm.Config{
		Logger: apmLogger,
	}

	db, err := openDatabase(cfg, gcfg)
	if err != nil {
		return nil, err
	}

	if err := validateDatabase(db, cfg); err != nil {
		return nil, err
	}

	handleMigration(db, cfg)
	return db, nil
}

func setDefaultPool(sqlDB *sql.DB) {
	if sqlDB == nil {
		return
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
}

func openDatabase(cfg config.DBConfig, gcfg *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.ToURL()), gcfg)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, errors.New("db instance is not ready")
	}
	return db, nil
}

func validateDatabase(db *gorm.DB, cfg config.DBConfig) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	setDefaultPool(sqlDB)
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database connection test failed: %v", err)
	}

	return nil
}

func handleMigration(db *gorm.DB, cfg config.DBConfig) {
	if strings.ToLower(cfg.DbEvent) != "migrate" {
		return
	}

	switch cfg.Id {
	case config.MainDB:
		{
			log.Info("Main Database is migrating")
			if err := db.AutoMigrate(
				// Core modules
				&schemas.User{},
				&schemas.Organization{},
				&schemas.Role{},
				&schemas.Permission{},
				&schemas.RolePermission{},
				&schemas.OrganizationMember{},
				// Unit management
				&schemas.Unit{},
				&schemas.UnitMember{},
				&schemas.UnitSettings{},
				// Posts / Announcements
				&schemas.Post{},
				&schemas.PostComment{},
				&schemas.PostPollOption{},
				&schemas.PostPollVote{},
			); err != nil {
				log.Error("migrate error : ", err)
			}
			log.Info("Main Database migrate successfully")

			runSeeders(db)
			break
		}
	}
}

func runSeeders(db *gorm.DB) {
	log.Info("🌱 Running database seeders...")

	permissions := getCorePermissions()
	for _, p := range permissions {
		var existing schemas.Permission
		if db.Where("name = ?", p.name).First(&existing).Error != nil {
			permission := schemas.Permission{
				Name:        p.name,
				Resource:    p.resource,
				Action:      p.action,
				Description: p.description,
			}
			if err := db.Create(&permission).Error; err == nil {
				log.Infof("✅ Created permission: %s", p.name)
			}
		}
	}

	roles := getDefaultRoles()
	for _, r := range roles {
		var existing schemas.Role
		if db.Where("name = ?", r.name).First(&existing).Error != nil {
			role := schemas.Role{
				Name:        r.name,
				Description: r.description,
			}
			if err := db.Create(&role).Error; err == nil {
				log.Infof("✅ Created role: %s", r.name)

				if r.allPermissions {
					var allPerms []schemas.Permission
					db.Find(&allPerms)
					for _, perm := range allPerms {
						db.Create(&schemas.RolePermission{RoleId: role.Id, PermissionId: perm.Id})
					}
					log.Infof("   → Assigned ALL permissions to %s", r.name)
				} else {
					for _, permName := range r.permissions {
						var perm schemas.Permission
						if db.Where("name = ?", permName).First(&perm).Error == nil {
							db.Create(&schemas.RolePermission{RoleId: role.Id, PermissionId: perm.Id})
						}
					}
					log.Infof("   → Assigned %d permissions to %s", len(r.permissions), r.name)
				}
			}
		}
	}

	var user schemas.User
	if db.Where("email = ?", "superadmin@mail.com").First(&user).Error == nil {
		if !user.IsSuperAdmin {
			user.IsSuperAdmin = true
			db.Save(&user)
			log.Info("✅ User superadmin@mail.com updated to Super Admin")
		}
	}

	log.Info("✅ Database seeding completed!")
}

type permissionData struct {
	name, resource, action, description string
}

func getCorePermissions() []permissionData {
	resources := []struct{ resource, description string }{
		{"users", "User"},
		{"roles", "Role"},
		{"permissions", "Permission"},
		{"organizations", "Organization"},
		{"units", "Unit"},
		{"unit_members", "Unit Member"},
		{"posts", "Post"},
	}
	actions := []struct{ action, description string }{
		{"create", "Create"}, {"read", "Read"}, {"update", "Update"},
		{"delete", "Delete"}, {"list", "List"},
	}

	var perms []permissionData
	for _, r := range resources {
		for _, a := range actions {
			perms = append(perms, permissionData{
				name:        r.resource + "." + a.action,
				resource:    r.resource,
				action:      a.action,
				description: a.description + " " + r.description,
			})
		}
	}
	return perms
}

type roleData struct {
	name, description string
	permissions       []string
	allPermissions    bool
}

func getDefaultRoles() []roleData {
	return []roleData{
		{"Super Admin", "Full system access", nil, true},
		{"Admin", "Organization administrator", []string{
			"organizations.create", "organizations.read", "organizations.update", "organizations.delete", "organizations.list",
			"units.create", "units.read", "units.update", "units.delete", "units.list",
			"unit_members.create", "unit_members.read", "unit_members.update", "unit_members.delete", "unit_members.list",
			"users.read", "users.list", "roles.read", "roles.list", "permissions.read", "permissions.list",
			"posts.create", "posts.read", "posts.update", "posts.delete", "posts.list",
		}, false},
		{"Member", "Basic member with read access", []string{
			"organizations.read",
			"units.read", "units.list",
			"posts.read", "posts.list",
		}, false},
	}
}

type APMGormLogger struct {
	logger logger.Interface
}

func (l *APMGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &APMGormLogger{
		logger: l.logger.LogMode(level),
	}
}

func (l *APMGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Info(ctx, msg, data...)
}

func (l *APMGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Warn(ctx, msg, data...)
}

func (l *APMGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Error(ctx, msg, data...)
}

func (l *APMGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sqlQuery, _ := fc()
	duration := time.Since(begin)

	tx := apm.TransactionFromContext(ctx)
	if tx != nil {
		span := tx.StartSpan("gorm_query", "db", nil)
		defer span.End()

		span.Context.SetDatabase(apm.DatabaseSpanContext{
			Instance:  config.APP.ServiceName + "_db",
			Statement: sqlQuery,
			Type:      "sql",
		})

		span.Duration = duration

		if err != nil {
			span.Outcome = "failure"
			e := apm.DefaultTracer().NewError(err)
			e.Send()
		} else {
			span.Outcome = "success"
		}
	}

	l.logger.Trace(ctx, begin, fc, err)
}
