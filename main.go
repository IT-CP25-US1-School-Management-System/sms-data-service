package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	myMiddL "github.com/IT-CP25-US1-School-Management-System/sms-data-service/middleware"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/route/v1"

	helperGRPC "github.com/GodeFvt/go-backend/grpc"
	"github.com/GodeFvt/go-backend/helper"
	helperCookie "github.com/GodeFvt/go-backend/helper/cookie"
	helperMinio "github.com/GodeFvt/go-backend/minio"
	helperRedis "github.com/GodeFvt/go-backend/redis"

	helperMiddl "github.com/GodeFvt/go-backend/helper/middleware"
	helperRoute "github.com/GodeFvt/go-backend/helper/route"
	"github.com/GodeFvt/go-backend/psql"

	_data_handler "github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1/handler"
	_psqldata_repo "github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1/repository"
	_psqldataset_repo "github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1/repository"
	_redis_repo "github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1/repository"
	_data_usecase "github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1/usecase"
	_db_connection_manager_repo "github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1/repository"
	_db_connection_manager_usercase "github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1/usecase"
	_document_repo "github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/document/v1/repository"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echoMiddL "github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cast"
	echoSwagger "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"
)

var (
	grpcMaxReceiveSize = (1024 * 1024) * cast.ToInt(helper.GetENV("GRPC_MAX_RECEIVE_SIZE", "4"))
)

var (
	APP_PORT     = helper.GetENV("APP_PORT", "3200")
	APP_MODE     = helper.GetENV("APP_MODE", "production")
	GRPC_PORT    = helper.GetENV("GRPC_PORT", "4200")
	GRPC_TIMEOUT = cast.ToInt(helper.GetENV("GRPC_TIMEOUT", "120"))
	ALLOW_ORIGIN = func() []string {
		origins := helper.GetENV("ALLOW_ORIGIN", "*")
		origins = strings.Trim(origins, `"`)
		origins = strings.TrimSpace(origins)
		return strings.Split(origins, ",")
	}()
	ALLOW_ORIGIN_HEADER       = strings.Split(helper.GetENV("ALLOW_ORIGIN_HEADER", ""), ",")
	ALLOW_ORIGIN_CREDENTIAL   = cast.ToBool(helper.GetENV("ALLOW_ORIGIN_CREDENTIAL", "true"))
	PSQL_DATABASE_DATASET_URL = helper.GetENV("PSQL_DATABASE_DATASET_URL", "postgres://postgres:postgres@psql_db:5432/app_example?sslmode=disable")

	SENTRY_DSN = helper.GetENV("SENTRY_DSN", "")

	JWKS_URL     = strings.Split(helper.GetENV("JWKS_URL", ""), ",")
	JWT_SECRET   = helper.GetENV("JWT_SECRET", "")
	ALLOWED_ALGS = strings.Split(helper.GetENV("ALLOWED_ALGS", "RS256,HS256"), ",")
	AUDIENCE     = strings.Split(helper.GetENV("AUDIENCE", ""), ",")
	ISSUER       = helper.GetENV("ISSUER", "")

	COOKIE_DOMAIN   = helper.GetENV("COOKIE_DOMAIN", "localhost")
	COOKIE_SECURE   = cast.ToBool(helper.GetENV("COOKIE_SECURE", "false"))
	COOKIE_HTTPONLY = cast.ToBool(helper.GetENV("COOKIE_HTTPONLY", "true"))

	REDIS_ADDRESS = helper.GetENV("REDIS_ADDRESS", "")

	CRYPTO_SECRET = helper.GetENV("CRYPTO_SECRET", "")

	DOCUMENT_SERVICE_GRPC_ADDRESS = helper.GetENV("DOCUMENT_SERVICE_GRPC_ADDRESS", "localhost:3200")
)

func connectPsqlDB(con string) *psql.Client {
	db, err := psql.NewConnection(con, psql.Postgres)
	if err != nil {
		panic(err)
	}
	return db
}

func connectMinio(endpoint, accessKey, secretKey, region string, ssl bool) *helperMinio.Client {
	client, err := helperMinio.NewMinio(endpoint, accessKey, secretKey, ssl, region)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("connecting minio successful!!")
	}

	return client
}

func connectRedis(address string) *helperRedis.Client {
	client, err := helperRedis.NewClient(address)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("connecting redis successful!!")

	return client
}

func startGRPCServer(server *grpc.Server) {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", GRPC_PORT))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	/* serve grpc */
	fmt.Printf("Start grpc Server [::%s]\n", GRPC_PORT)
	if err := server.Serve(listen); err != nil {
		panic(err)
	}
}

func main() {
	// /* init sentry */
	sentryErr := sentry.Init(sentry.ClientOptions{
		Dsn: SENTRY_DSN,
	})

	// /* init psqlClient */
	psqlDatasetClient := connectPsqlDB(PSQL_DATABASE_DATASET_URL)

	defer psqlDatasetClient.GetClient().Close()

	// /* init redisClient */
	redisClient := connectRedis(REDIS_ADDRESS)
	defer redisClient.Close()

	// /* init grpc */
	server := helperGRPC.NewServer(grpc.MaxRecvMsgSize(grpcMaxReceiveSize))
	defer server.GracefulStop()

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.HTTPErrorHandler = helperMiddl.SentryCapture(e)
	helperRoute.RegisterVersion(e)

	e.Use(echoMiddL.Logger())
	e.Use(echoMiddL.Recover())
	e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))

	e.GET("/health-check", func(c echo.Context) error {
		resp := echo.Map{
			"status": "ok",
		}
		return c.JSON(http.StatusOK, resp)
	})

	if APP_MODE == "development" {
		e.Static("/data/api/api-spec/data/v1", "./api-spec/data/v1")
		e.GET("/api-spec/*", echoSwagger.EchoWrapHandler(
			func(c *echoSwagger.Config) {
				c.URLs = []string{"/data/api/api-spec/data/v1/openapi_bundle.json"}
				c.DocExpansion = "list"
				c.DomID = "swagger-ui"
				c.InstanceName = "swagger"
				c.DeepLinking = true
				c.PersistAuthorization = false
				c.SyntaxHighlight = true
			}))
	}

	/* cookie manager */
	cookieManager := helperCookie.NewCookieManager(COOKIE_DOMAIN, COOKIE_SECURE, COOKIE_HTTPONLY)

	/* middleware */
	middL := myMiddL.InitMiddleware(redisClient, cookieManager, JWT_SECRET, JWKS_URL, ALLOWED_ALGS, AUDIENCE, ISSUER)
	e.Use(myMiddL.DefaultCORSConfig(ALLOW_ORIGIN, ALLOW_ORIGIN_HEADER, ALLOW_ORIGIN_CREDENTIAL))
	e.Use(middL.InitContextIfNotExists)
	e.Use(middL.ErrorHandlerMiddleware)
	e.Use(middL.SetTracer)

	/* repository */
	psqlDatasetRepo := _psqldataset_repo.NewPsqlDatasetRepository(psqlDatasetClient)
	redisRepo := _redis_repo.NewRedisRepository(redisClient)

	/* database connection manager */
	dbConnectionManager := _db_connection_manager_repo.NewDBConnectionManagerRepository(psqlDatasetRepo, CRYPTO_SECRET)
	dbConnectionManagerUsecase := _db_connection_manager_usercase.NewDBConnectionManagerUsecase(dbConnectionManager)
	defer dbConnectionManagerUsecase.CloseAll()

	psqlDataRepo := _psqldata_repo.NewPsqlDataRepository(dbConnectionManagerUsecase)
	documentRepo := _document_repo.NewGRPCDocumentRepository(DOCUMENT_SERVICE_GRPC_ADDRESS, GRPC_TIMEOUT)

	/* usecase */
	dataUsecase := _data_usecase.NewDataUsecase(psqlDataRepo, psqlDatasetRepo, documentRepo, redisRepo, CRYPTO_SECRET)

	/* handler */
	dataHandler := _data_handler.NewDataHandler(dataUsecase)

	/* gprc handler */

	/* validate */

	/* inject route */
	router := route.NewRoute(e, middL)
	router.RegisterDataRoute(dataHandler)

	/* inject grpc route */

	/* serve gprc */
	go func() {
		if r := recover(); r != nil {
			fmt.Println("error on start grpc server: ", r.(error))
		}
		startGRPCServer(server)
	}()

	/* serve echo */
	port := fmt.Sprintf(":%s", APP_PORT)
	if sentryErr == nil {
		sentry.CaptureException(e.Start(port))
	} else {
		e.Logger.Fatal(e.Start(port))
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		var errorMessages []string
		for _, e := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf(" Field '%s' failed validation due to: %s", strings.ToLower(e.Field()), e.Tag()+" "+e.Param()))
		}
		return echo.NewHTTPError(http.StatusBadRequest, strings.Join(errorMessages, ", "))
	}
	return nil
}
