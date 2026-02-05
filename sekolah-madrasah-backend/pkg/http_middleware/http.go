package http_middleware

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"sekolah-madrasah/config"
	"sekolah-madrasah/pkg/auth_utils"
	"sekolah-madrasah/pkg/rbac"
	"sekolah-madrasah/pkg/request_utils"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"go.elastic.co/apm/v2"
)

func CORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.APP.Rest.Origin)
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Permissions-Policy", "geolocation=(), camera=(), microphone=()")
	c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	c.Writer.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
	c.Writer.Header().Set("Cross-Origin-Resource-Policy", "same-site")
	c.Writer.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH")
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

func JWTAuthentication(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header provided"})
		c.Abort()
		return
	}
	extractedToken := strings.Split(tokenString, "Bearer ")
	if len(extractedToken) == 2 {
		tokenString = strings.TrimSpace(extractedToken[1])
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect Format of Authorization Token"})
		c.Abort()
		return
	}
	claims, err := auth_utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	ctx := auth_utils.WithAuthClaim(c.Request.Context(), claims)
	c.Request = c.Request.WithContext(ctx)
	c.Set("auth", claims)
	c.Set("token", tokenString)
	c.Set("user_id", claims.UserID) // Set user_id for easy access in controllers
	c.Next()
}

func CronAuthentication(c *gin.Context) {
	cronToken := c.GetHeader("x-auth-cron")
	if cronToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No x-auth-cron header provided"})
		c.Abort()
		return
	}

	if cronToken != config.APP.Rest.XAuthCron {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid cron authentication token"})
		c.Abort()
		return
	}

	// Set a special context marker to identify cron requests
	c.Set("is_cron_request", true)
	c.Next()
}

// getLogLevel determines log level based on HTTP status code
func getLogLevel(statusCode int) string {
	switch {
	case statusCode >= 500:
		return "error"
	case statusCode >= 400:
		return "warning"
	case statusCode >= 300:
		return "info"
	default:
		return "info"
	}
}

func ElasticAPM(c *gin.Context) {
	// Debug: Check if APM is enabled and tracer is available
	tracer := apm.DefaultTracer()
	if tracer == nil {
		log.Warn("APM tracer is nil - skipping transaction")
		c.Next()
		return
	}

	// Start APM transaction for this request
	tx := tracer.StartTransaction(c.Request.Method+" "+c.Request.URL.Path, "request")
	if tx == nil {
		log.Warn("Failed to create APM transaction")
		c.Next()
		return
	}

	defer func() {
		if tx != nil {
			tx.End()
		}
	}()

	// Set transaction context
	if tx != nil {
		tx.Context.SetHTTPRequest(c.Request)
		tx.Context.SetCustom("http.method", c.Request.Method)
		tx.Context.SetCustom("http.url", c.Request.URL.String())
		tx.Context.SetCustom("http.user_agent", c.Request.UserAgent())
		tx.Context.SetCustom("http.remote_addr", c.Request.RemoteAddr)
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				tx.Context.SetCustom("http.body", "invalid body")
			} else {
				bodyString := string(bodyBytes)
				tx.Context.SetCustom("http.body", bodyString)
			}
		}

		// Put transaction in context so handlers can access it
		c.Request = c.Request.WithContext(apm.ContextWithTransaction(c.Request.Context(), tx))
	}

	// Process request
	start := time.Now()
	c.Next()
	duration := time.Since(start)

	// Set response data in transaction
	if tx != nil {
		tx.Context.SetCustom("http.status_code", c.Writer.Status())
		tx.Context.SetCustom("http.duration_ms", duration.Milliseconds())
		tx.Context.SetCustom("request.completed_at", time.Now().Format(time.RFC3339))

		// Add trace and transaction IDs (corrected for APM v2)
		traceID := tx.TraceContext().Trace
		if traceID != (apm.TraceID{}) {
			tx.Context.SetCustom("request.trace_id", traceID.String())
		}
		spanID := tx.TraceContext().Span
		if spanID != (apm.SpanID{}) {
			tx.Context.SetCustom("request.transaction_id", spanID.String())
		}

		// Add request details
		tx.Context.SetCustom("request.method", c.Request.Method)
		tx.Context.SetCustom("request.url", c.Request.URL.Path)
		tx.Context.SetCustom("request.full_url", c.Request.URL.String())
		tx.Context.SetCustom("request.user_agent", c.Request.UserAgent())
		tx.Context.SetCustom("request.remote_addr", c.Request.RemoteAddr)

		// Add query params if exists
		if len(c.Request.URL.RawQuery) > 0 {
			tx.Context.SetCustom("request.query_params", c.Request.URL.RawQuery)
		}

		// Add response size
		if c.Writer.Size() > 0 {
			tx.Context.SetCustom("response.size_bytes", c.Writer.Size())
		}

		// Add user info if available
		if _, exists := c.Get("auth"); exists {
			tx.Context.SetCustom("user.authenticated", "true")
			// Only log basic auth status, avoid exposing sensitive data
			tx.Context.SetCustom("user.has_auth", "true")
		}

		// Add service info
		tx.Context.SetCustom("service.version", "1.0.0")
		tx.Context.SetCustom("service.environment", "development")

		// Add structured log message
		logMessage := fmt.Sprintf("Request completed: %s %s - Status: %d, Duration: %dms",
			c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration.Milliseconds())

		// Log dengan prefix untuk memudahkan searching di Kibana
		tx.Context.SetCustom("log.message", logMessage)
		tx.Context.SetCustom("log.timestamp", time.Now().Format(time.RFC3339))
		tx.Context.SetCustom("log.level", getLogLevel(c.Writer.Status()))
	}

	if tracer != nil && tx != nil {
		// Force flush to ensure transaction is sent immediately
		tracer.Flush(nil)
	}
}

func RBACMiddleware(c *gin.Context) {
	msg := "Access Denied"
	path := c.FullPath()
	role := rbac.Unauthorized
	_, ok := c.Get("claims")
	if ok {
		// auth := helper.GetAuthClaim(c)
		// if strings.ToLower(auth.Role) == "superadmin" {
		// 	role = rbac.SuperAdmin
		// } else if strings.ToLower(auth.Role) == "owner" {
		// 	role = rbac.Owner
		// } else if strings.ToLower(auth.Role) == "member" {
		// 	role = rbac.Member
		// }
	}
	routerDetails := rbac.FindMany(Router, path, rbac.Filter{Method: c.Request.Method})
	if len(routerDetails) > 0 {
		router := rbac.FindOne[rbac.RouterDetail](routerDetails, path, &rbac.Filter{
			Role:   &role,
			Method: c.Request.Method,
		})
		if router != nil {
			if len(router.MustUrlParams) > 0 {
				_, err := request_utils.ParamsToMapOneValue(c.Request.URL.Query()).Must(router.MustUrlParams).Result()
				if err != nil {
					msg = err.Error()
					if router.CustomParamsMessage != "" {
						msg = router.CustomParamsMessage
					}
					c.JSON(http.StatusBadRequest, gin.H{"error": msg})
					c.Abort()
					return
				}
			}
			for _, mid := range router.CustomMiddleware {
				mid(c)
			}
			c.Next()
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": msg})
		c.Abort()
		return
	}
}
