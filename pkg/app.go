package pkg

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"rdpms25-go-rpc-service/pkg/config"
	"rdpms25-go-rpc-service/pkg/models"
	"rdpms25-go-rpc-service/pkg/util"
	"rdpms25-go-rpc-service/pkg/util/connection"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

func Start(version, buildTime string) {
	util.DurationLogger().Start("app_setup")
	conf, logger, db, err := config.Initialise()
	logger.Info("starting application", "version", version, "build", buildTime)
	if err != nil {
		logger.Error("unable to start application", "err", err)
		os.Exit(1)
	}

	logger.Info("starting application", "version", version, "build", buildTime)

	// initialise repo
	assetRepo, assetTypeRepo, assetViewRepo, ocRepo, ocViewRepo := repository.NewAsset(db), repository.NewAssetType(), repository.NewAssetView(), repository.NewOriginatorClass(), repository.NewOriginatorClassView()
	paramTypeRepo, paramReprRepo, paramReprRangeRepo, paramViewRepo, paramReprRangeView := repository.NewParamType(), repository.NewParamRepr(), repository.NewParamReprRange(), repository.NewParamView(), repository.NewParamReprRangeView()
	telemetryRepo, telemetryViewRepo, historyRepo, historyViewRepo := repository.NewTelemetry(), repository.NewTelemetryView(), repository.NewHistory(), repository.NewHistoryView()
	alarmRepo, alarmViewRepo, alarmTraceRepo, alarmAckHistoryRepo := repository.NewAlarm(), repository.NewAlarmView(), repository.NewAlarmTrace(), repository.NewAckHistoryRepo()
	userProfileRepo, userFcmTokenRepo := repository.NewUserProfile(), repository.NewFcmUserToken()
	userPermissionRepo := repository.NewUserPermission()
	userSessionRepo := repository.NewUserSession()
	alarmMissingRepo, missingAlarmViewRepo := repository.NewAlarmMissingRecord(), repository.NewAlarmMissingRecordView()
	genericRepo, userRepo, iotDeviceRepo, appConRepo := repository.NewGeneric(), repository.NewUser(), repository.NewIotDevice(), repository.NewAppCon()
	extTokenRepo := repository.NewExternalAPIToken()
	assetParamIotDeviceMapRepo := repository.NewAssetParamIotDeviceMap()
	assetParamAvgRepo, assetParamAvgHistoryRepo, assetParamAvgView, assetParamAvgHistoryView := repository.NewAssetParamAverage(), repository.NewAssetParamAverageHistory(), repository.NewAssetParamAverageView(), repository.NewAssetParamAverageHistoryView()
	paramReprAssetMappingRepo := repository.NewParamReprAssetMapping()
	orgRepo, orgViewRepo := repository.NewOrganisation(), repository.NewOrganisationView()
	notRepo, notMediumRepo, notControlRepo, notViewRepo := repository.NewNotification(), repository.NewNotificationMedium(), repository.NewNotificationControl(), repository.NewNotificationView()
	edgeRepo := repository.NewEdgegateway()
	alarmReportRepo := repository.NewAlarmReportRepo(db)
	auditLogRepo, changeLogRepo := repository.NewAuditLog(), repository.NewChangeLog()
	alarmNotificationFreqRepo := repository.NewAlarmNotificationFrequency()
	amRepo, amViewRepo, amStatusRepo, atmRepo := repository.NewAssetMaintenanceModeLog(), repository.NewAssetMaintenanceModeLogView(), repository.NewAssetMaintenanceModeLogStatusView(), repository.NewAssetTypeMaintenanceModeTimeout()
	utilisationRepo := repository.NewUtilisationRepo(db)
	staticRepo := repository.NewStaticOptionRepo()
	assetCoRepo := repository.NewAssetCorrelation()
	ackHistoryViewRepo := repository.NewAlarmAckHistoryView()
	rdpmsHealthReportRepo := repository.NewRDPMSHealthReport(db)
	correlationRepo, CorrelationViewRepo := repository.NewCorrelation(), repository.NewCorrelationAssetView()
	smmsRepo := repository.NewSmmsRepo(db)
	sourceRepo, assetMappingRepo := repository.NewSourceRepo(), repository.NewAssetDataMappingRepo()
	// healthRepo, healthRepoView := repository.NewHealthIssue(), repository.NewHealthIssueViews()
	// initialise services
	assetServ, assetTypeServ, assetViewServ, ocServ, ocViewServ := service.NewAsset(assetRepo), service.NewAssetType(assetTypeRepo), service.NewAssetView(assetViewRepo), service.NewOriginatorClass(ocRepo), service.NewOriginatorClassView(ocViewRepo)
	paramTypeServ, paramReprServ, paramReprRangeServ, paramViewServ, paramReprRangeViewServ := service.NewParamType(paramTypeRepo), service.NewParamRepr(paramReprRepo), service.NewParamReprRange(paramReprRangeRepo), service.NewParamView(paramViewRepo), service.NewParamReprRangeView(paramReprRangeView)
	telemetryServ, telemetryViewServ, historyServ, historyViewServ := service.NewTelemetry(telemetryRepo), service.NewTelemetryView(telemetryViewRepo), service.NewHistory(historyRepo), service.NewHistoryView(historyViewRepo)
	alarmServ, alarmTraceServ, alarmViewServ, alarmAckHistoryServ := service.NewAlarm(alarmRepo, alarmTraceRepo), service.NewAlarmTrace(alarmTraceRepo), service.NewAlarmView(alarmViewRepo), service.NewAckHistoryService(alarmAckHistoryRepo)
	missingAlarmServ, missingAlarmViewServ := service.NewAlarmMissingRecord(alarmMissingRepo), service.NewAlarmMissingRecordView(missingAlarmViewRepo)
	genericServ, userServ, iotDeviceServ, appConServ := service.NewGeneric(genericRepo), service.NewUser(userRepo), service.NewIotDevice(iotDeviceRepo), service.NewAppCon(appConRepo)
	assetParamIotDeviceMapServ := service.NewAssetParamIotDeviceMap(assetParamIotDeviceMapRepo)
	assetParamAvgServ, assetParamAvgHistoryServ, assetParamAvgViewServ, assetParamAvgHistoryViewServ := service.NewAssetParamAverageService(assetParamAvgRepo), service.NewAssetParamAverageHistoryService(assetParamAvgHistoryRepo), service.NewAssetParamAverageView(assetParamAvgView), service.NewAssetParamAverageHistoryView(assetParamAvgHistoryView)
	alarmReportServ := service.NewAlarmReportService(alarmReportRepo)
	paramReprAssetMappingServ := service.NewParamReprAssetMapping(paramReprAssetMappingRepo)
	orgServ, orgViewServ := service.NewOrganisation(orgRepo), service.NewOrganisationView(orgViewRepo)
	notServ, notMediumServ, notControlServ, notViewServ, edgeServ := service.NewNotification(notRepo), service.NewNotificationMedium(notMediumRepo), service.NewNotificationControl(notControlRepo), service.NewNotificationView(notViewRepo), service.NewEdgegateway(edgeRepo)
	userProfileServ, fcmServ := service.NewUserProfile(userProfileRepo), service.NewUserFcmToken(userFcmTokenRepo)
	userPermissionServ := service.NewUserPermission(userPermissionRepo)
	userSessionServ := service.NewUserSession(userSessionRepo)
	auditLogServ, changeLogServ := service.NewAuditLog(auditLogRepo), service.NewChangeLog(changeLogRepo)
	alarmNotificationFreqServ := service.NewAlarmNotificationFrequency(alarmNotificationFreqRepo)
	amServ, amViewServ, amStatusServ, amtServ := service.NewAssetMaintenanceModeLogService(amRepo), service.NewAssetMaintenanceModeLogViewService(amViewRepo), service.NewAssetMaintenanceModeLogStatusViewService(amStatusRepo), service.NewAssetTypeMaintenanceModeTimeoutService(atmRepo)
	extTokenServ := service.NewExternalAPIToken(extTokenRepo)
	smmsServ, utilisationServ, staticServ := service.NewSmmsAssetService(), service.NewUtilisationService(utilisationRepo), service.NewStaticOption(staticRepo)
	smmsDashboardServ := service.NewSmmsService(smmsRepo)
	sourceServ, assetMappingServ := service.NewSourceService(sourceRepo, models.TableNames.Source), service.NewAssetDataMappingService(assetMappingRepo, models.TableNames.AssetDataMapping)
	ackHistoryServView := service.NewAlarmAckHistoryView(ackHistoryViewRepo)
	assetCoServ := service.NewAssetCorrelationService(assetCoRepo)
	rdpmsHealthReportServ := service.NewRDPMSHealthReport(rdpmsHealthReportRepo)
	correlationServ, correlationViewServ := service.NewCorrelation(correlationRepo), service.NewCorrelationAssetView(CorrelationViewRepo)
	// healthServ := service.NewHealthIssueViews(healthRepo)
	// healthServView := service.NewHealthIssueView(healthRepoView)
	//initialise sms provider
	smsProvider := auth.NewSMSProvider()

	ctx := context.Background()

	// workers
	// workers
	var natsConn *nats.Conn
	if conf.Nats.Enable {
		natsConn, err = connection.NewNatsConnection(fmt.Sprintf("rdpms25-cloud-api-%s", conf.App.VendorCloudCode), conf.Nats)
		if err != nil {
			logger.Error("unable to start application", "err", err)
			os.Exit(1)
		}
		crudEmitter := worker.NewCRUDEmitter(natsConn, alarmServ, assetServ, edgeServ, amServ, userServ, fcmServ, orgServ, alarmNotificationFreqServ, notMediumServ, notControlServ, extTokenServ)
		crudEmitter.Subscribe()
	}
	//scheduler
	scheduler := scheduler.NewAssetMaintenanceScheduler(amViewServ, amServ, amtServ)
	go scheduler.BlockingStart(ctx)
	// setup rest apis
	engine := gin.New()

	engine.Use(func(c *gin.Context) {
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), midi=(), sync-xhr=(), microphone=(), camera=(), magnetometer=(), gyroscope=(), fullscreen=(self), payment=()")
		c.Next()
	})

	// rdpms_health
	rdpmsHealthViewServ := service.NewRDPMSHealthView()
	rdpmsHealthGroupbyAssetViewServ := service.NewRDPMSHealthGroupbyAssetView()

	//notification module

	controller.Initialise(
		conf, engine,
		assetServ, assetTypeServ, assetParamAvgServ, assetParamAvgHistoryServ, assetViewServ, assetParamAvgViewServ, assetParamAvgHistoryViewServ, ocServ, ocViewServ,
		paramTypeServ, paramReprServ, paramReprRangeServ, paramViewServ, paramReprRangeViewServ,
		telemetryServ, telemetryViewServ, historyServ, historyViewServ,
		alarmServ, alarmTraceServ, alarmViewServ, alarmAckHistoryServ,
		genericServ, userServ, userProfileServ, userPermissionServ, iotDeviceServ, appConServ,
		orgServ, orgViewServ,
		notServ, notMediumServ, notControlServ, notViewServ,
		edgeServ,
		auditLogServ, changeLogServ,
		alarmNotificationFreqServ,
		paramReprAssetMappingServ,
		alarmReportServ,
		amServ, amViewServ, amStatusServ, amtServ, smmsServ, utilisationServ, staticServ, assetCoServ, smsProvider, natsConn,
		ackHistoryServView,
		rdpmsHealthViewServ, rdpmsHealthGroupbyAssetViewServ,
		rdpmsHealthReportServ,
		assetParamIotDeviceMapServ,
		correlationServ, correlationViewServ,
		fcmServ,
		userSessionServ, missingAlarmServ, missingAlarmViewServ,
		extTokenServ,
		smmsDashboardServ,
		sourceServ, assetMappingServ,
	)

	util.DurationLogger().End("app_setup")

	slog.Info("starting server", "port", conf.App.Port)
	if err := engine.Run(fmt.Sprintf("%s:%d", "0.0.0.0", conf.App.Port)); err != nil {
		slog.Error("unable to start server", "err", err)
	}
}
