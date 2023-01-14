package engine

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"unsafe"

	"github.com/benallfree/pbscript/modules/pbscript/event"

	"github.com/goccy/go-json"

	"github.com/dop251/goja"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

const (
	// App Hooks
	ON_BEFORE_BOOTSTRAP_HOOK = "OnBeforeBootstrap"
	ON_AFTER_BOOTSTRAP_HOOK  = "OnAfterBootstrap"
	ON_BEFORE_SERVE_HOOK     = "OnBeforeServe"
	ON_API_ERROR_HOOK        = "OnBeforeApiError"

	// Db Books
	ON_MODEL_BEFORE_CREATE_HOOK = "OnModelBeforeCreate"
	ON_MODEL_AFTER_CREATE_HOOK  = "OnModelAfterCreate"
	ON_MODEL_BEFORE_UPDATE_HOOK = "OnModelBeforeUpdate"
	ON_MODEL_AFTER_UPDATE_HOOK  = "OnModelAfterUpdate"
	ON_MODEL_BEFORE_DELETE_HOOK = "OnModelBeforeDelete"
	ON_MODEL_AFTER_DELETE_HOOK  = "OnModelAfterDelete"

	// Mailer Hooks
	ON_MAILER_BEFORE_ADMIN_RESET_PASSWORD_SEND_HOOK  = "OnMailerBeforeAdminResetPasswordSend"
	ON_MAILER_AFTER_ADMIN_RESET_PASSWORD_SEND_HOOK   = "OnMailerAfterAdminResetPasswordSend"
	ON_MAILER_BEFORE_RECORD_RESET_PASSWORD_SEND_HOOK = "OnMailerBeforeRecordResetPasswordSend"
	ON_MAILER_AFTER_RECORD_RESET_PASSWORD_SEND_HOOK  = "OnMailerAfterRecordResetPasswordSend"
	ON_MAILER_BEFORE_RECORD_VERIFICATION_SEND_HOOK   = "OnMailerBeforeRecordVerificationSend"
	ON_MAILER_AFTER_RECORD_VERIFICATION_SEND_HOOK    = "OnMailerAfterRecordVerificationSend"
	ON_MAILER_BEFORE_RECORD_CHANGE_EMAIL_SEND_HOOK   = "OnMailerBeforeRecordChangeEmailSend"
	ON_MAILER_AFTER_RECORD_CHANGE_EMAIL_SEND_HOOK    = "OnMailerAferRecordChangeEmailSend"

	// Record API Hooks
	ON_RECORDS_LIST_REQUEST_HOOK                         = "OnRecordsListRequest"
	ON_RECORD_VIEW_REQUEST_HOOK                          = "OnRecordViewRequest"
	ON_RECORD_BEFORE_CREATE_REQUEST_HOOK                 = "OnRecordBeforeCreateRequest"
	ON_RECORD_AFTER_CREATE_REQUEST_HOOK                  = "OnRecordAfterCreateRequest"
	ON_RECORD_BEFORE_UPDATE_REQUEST_HOOK                 = "OnRecordBeforeUpdateRequest"
	ON_RECORD_AFTER_UPDATE_REQUEST_HOOK                  = "OnRecordAfterUpdateRequest"
	ON_RECORD_BEFORE_DELETE_REQUEST_HOOK                 = "OnRecordBeforeDeleteRequest"
	ON_RECORD_AFTER_DELETE_REQUEST_HOOK                  = "OnRecordAfterDeleteRequest"
	ON_RECORD_AUTH_REQUEST_HOOK                          = "OnRecordAuthRequest"
	ON_RECORD_LIST_EXTERNAL_AUTHS_REQUEST_HOOK           = "OnRecordListExternalAuthsRequest"
	ON_RECORD_BEFORE_UNLINK_EXTERNAL_AUTH_REQUEST_HOOK   = "OnRecordBeforeUnlinkExternalAuthRequest"
	ON_RECORD_AFTER_UNLINK_EXTERNAL_AUTH_REQUEST_HOOK    = "OnRecordAfterUnlinkExternalAuthRequest"
	ON_RECORD_BEFORE_REQUEST_VERIFICATION_REQUEST_HOOK   = "OnRecorBeforeRequestVerificationRequest"
	ON_RECORD_AFTER_REQUEST_VERIFICATION_REQUEST_HOOK    = "OnRecordAfterRequestVerificationRequest"
	ON_RECORD_BEFORE_CONFIRM_VERIFICATION_REQUEST_HOOK   = "OnRecordBeforeConfirmVerificationRequest"
	ON_RECORD_AFTER_CONFIRM_VERIFICATION_REQUEST_HOOK    = "OnRecordAfterConfirmVerificationRequest"
	ON_RECORD_BEFORE_REQUEST_PASSWORD_RESET_REQUEST_HOOK = "OnRecordBeforeRequestPasswordResetRequest"
	ON_RECORD_AFTER_REQUEST_PASSWORD_RESET_REQUEST_HOOK  = "OnRecordAfterRequestPasswordResetRequest"
	ON_RECORD_BEFORE_CONFIRM_PASSWORD_RESET_REQUEST_HOOK = "OnRecordBeforeConfirmPasswordResetRequest"
	ON_RECORD_AFTER_CONFIRM_PASSWORD_RESET_REQUEST_HOOK  = "OnRecordAfterConfirmPasswordResetRequest"
	ON_RECORD_BEFORE_REQUEST_EMAIL_CHANGE_REQUEST_HOOK   = "OnRecordBeforeRequestEmailChangeRequest"
	ON_RECORD_AFTER_REQUEST_EMAIL_CHANGE_REQUEST_HOOK    = "OnRecordAfterRequestEmailChangeRequest"
	ON_RECORD_BEFORE_CONFIRM_EMAIL_CHANGE_REQUEST_HOOK   = "OnRecordBeforeConfirmEmailChangeRequest"
	ON_RECORD_AFTER_CONFIRM_EMAIL_CHANGE_REQUEST_HOOK    = "OnRecordAfterConfirmEmailChangeRequest"

	// Realtime API Hooks
	ON_REALTIME_CONNECT_REQUEST_HOOK          = "OnRealtimeConnectRequest"
	ON_REALTIME_DISCONNECT_REQUEST_HOOK       = "OnRealtimeDisconnectRequest"
	ON_REALTIME_BEFORE_SUBSCRIBE_REQUEST_HOOK = "OnRealtimeBeforeSubscribeRequest"
	ON_REALTIME_AFTER_SUBSCRIBE_REQUEST_HOOK  = "OnRealtimeAfterSubscribeRequest"
	ON_REALTIME_BEFORE_MESSAGE_SEND_HOOK      = "OnRealtimeBeforeMessageSend"
	ON_REALTIME_AFTER_MESSAGE_SEND_HOOK       = "OnRealtimeAfterMessageSend"

	// File API Hooks
	ON_FILE_DONWLOAD_REQUEST_HOOK = "OnFileDownloadRequest"

	// Collection API Hooks
	ON_COLLECTION_LIST_REQUEST_HOOK           = "OnCollectionListRequest"
	ON_COLLECTION_VIEW_REQUEST_HOOK           = "OnCollectionViewRequest"
	ON_COLLECTION_BEFORE_CREATE_REQUEST_HOOK  = "OnCollectionBeforeCreateRequest"
	ON_COLLECTION_AFTER_CREATE_REQUEST_HOOK   = "OnCollectionAfterCreateRequest"
	ON_COLLECTION_BEFORE_UPDATE_REQUEST_HOOK  = "OnCollectionBeforeUpdateRequest"
	ON_COLLECTION_AFTER_UPDATE_REQUEST_HOOK   = "OnCollectionAfterUpdateRequest"
	ON_COLLECTION_BEFORE_DELETE_REQUEST_HOOK  = "OnCollectionBeforeDeleteRequest"
	ON_COLLECTION_AFTER_DELETE_REQUEST_HOOK   = "OnCollectionAfterDeleteRequest"
	ON_COLLECTIONS_BEFORE_IMPORT_REQUEST_HOOK = "OnCollectionBeforeImportRequest"
	ON_COLLECTIONS_AFTER_IMPORT_REQUEST_HOOK  = "OnCollectionAfterImportRequest"

	// Settings API Hooks
	ON_SETTINGS_LIST_REQUEST_HOOK          = "OnSettingsListRequest"
	ON_SETTINGS_BEFORE_UPDATE_REQUEST_HOOK = "OnSettingsBeforeUpdateRequest"
	ON_SETTINGS_AFTER_UPDATE_REQUEST_HOOK  = "OnSettingsAfterUpdateRequest"

	// Admin API Hooks
	ON_ADMINS_LIST_REQUEST_HOOK         = "OnAdminsListRequest"
	ON_ADMIN_VIEW_REQUEST_HOOK          = "OnAdminViewRequest"
	ON_ADMIN_BEFORE_CREATE_REQUEST_HOOK = "OnAdminBeforeCreateRequest"
	ON_ADMIN_AFTER_CREATE_REQUEST_HOOK  = "OnAdminAfterCreateRequest"
	ON_ADMIN_BEFORE_UPDATE_REQUEST_HOOK = "OnAdminBeforeUpdateRequest"
	ON_ADMIN_AFTER_UPDATE_REQUEST_HOOK  = "OnAdminAfterUpdateRequest"
	ON_ADMIN_BEFORE_DELETE_REQUEST_HOOK = "OnAdminBeforeDeleteRequest"
	ON_ADMIN_AFTER_DELETE_REQUEST_HOOK  = "OnAdminAfterDeleteRequest"
	ON_ADMIN_AUTH_REQUEST_HOOK          = "OnAdminAuthRequest"
)

var app *pocketbase.PocketBase
var router *echo.Echo
var vm *goja.Runtime
var cleanups = []func(){}
var __go_apis *goja.Object

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func logF(color, format string, args ...any) (n int, err error) {
	s := append(args, string(colorReset))
	fmt.Print(color)
	res, err := fmt.Printf(format, s...)
	fmt.Print(colorReset)
	return res, err
}

func logErrorf(format string, args ...any) (n int, err error) {
	return logF(colorRed, format, args...)
}

func cleanup(msg string, cb func()) {
	fmt.Printf("adding cleanup: %s\n", msg)
	cleanups = append(cleanups, func() {
		fmt.Printf("executing cleanup: %s\n", msg)
		cb()
	})
}

func bootstrapEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.BootstrapEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.BootstrapEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func serveEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.ServeEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.ServeEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func apiErrorEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.ApiErrorEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.ApiErrorEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func modelEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.ModelEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.ModelEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func mailerRecordEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.MailerRecordEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.MailerRecordEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func mailerAdminEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.MailerAdminEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.MailerAdminEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func realtimeConnectEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RealtimeConnectEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RealtimeConnectEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func realtimeDisconnectEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RealtimeDisconnectEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RealtimeDisconnectEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func realtimeSubscribeEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RealtimeSubscribeEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RealtimeSubscribeEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func realtimeMessageEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RealtimeMessageEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RealtimeMessageEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func settingsListEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.SettingsListEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.SettingsListEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func settingsUpdateEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.SettingsUpdateEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.SettingsUpdateEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func recordsListEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RecordsListEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RecordsListEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func recordViewEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RecordViewEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RecordViewEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func recordCreateEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RecordCreateEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RecordCreateEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func recordUpdateEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RecordUpdateEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RecordUpdateEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func recordDeleteEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RecordDeleteEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RecordDeleteEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func recordAuthEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RecordAuthEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RecordAuthEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func recordUnlinkExternalAuthEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RecordUnlinkExternalAuthEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RecordUnlinkExternalAuthEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func recordRequestPasswordResetEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RecordRequestPasswordResetEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RecordRequestPasswordResetEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func recordConfirmPasswordResetEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RecordDeleteEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RecordConfirmPasswordResetEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func recordRequestVerificationEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RecordRequestVerificationDeleteEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RecordRequestVerificationDeleteEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func recordConfirmVerificationEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RecordConfirmVerificationEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RecordConfirmVerificationEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func recordRequestEmailChangeEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RecordRequestEmailChangeEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RecordRequestEmailChangeEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func recordConfirmEmailChangeEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RecordConfirmEmailChangeEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RecordConfirmEmailChangeEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func recordListExternalAuthsEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.RecordListExternalAuthsEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.RecordListExternalAuthsEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func adminsListEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.AdminsListEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.AdminisListEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func adminViewEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.AdminViewEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.AdminViewEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func adminCreateEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.AdminCreateEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.AdminCreateEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func adminUpdateEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.AdminUpdateEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.AdminUpdateEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func adminDeleteEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.AdminDeleteEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.AdminDeleteEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func adminAuthEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.AdminAuthEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.AdminAuthEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func collectionsListEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.CollectionsListEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.CollectionsListEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func collectionViewEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.CollectionViewEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.CollectionViewEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func collectionCreateEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.CollectionCreateEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.CollectionCreateEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func collectionUpdateEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.CollectionUpdateEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.CollectionUpdateEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func collectionDeleteEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.CollectionDeleteEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.CollectionDeleteEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func collectionsImportEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.CollectionsImportEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.CollectiosImportEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func fileDownloadEventCallbackConstructor(eventName string) {
	__go_apis.Set(eventName, func(cb func(e *core.FileDownloadEvent)) {
		fmt.Println("Listening in GO for %s\n", eventName)
		unsub := event.On(event, func(e *event.UnknownPayload) {
			core((*core.FileDownloadEvent)(unsafe.Pointer(e)))
		})
		cleanup(eventName, unsub)
	})
}

func bindApis() {
	__go_apis = vm.NewObject()
	__go_apis.Set("addRoute", func(route echo.Route) {
		method := route.Method
		path := route.Path
		fmt.Printf("Adding route: %s %s\n", method, path)

		router.AddRoute(route)
		cleanup(
			fmt.Sprintf("route %s %s", method, path),
			func() {
				router.Router().Remove(method, path)
			})
	})

	// App Hooks
	bootstrapEventCallbackConstructor(ON_BEFORE_BOOTSTRAP_HOOK)
	bootstrapEventCallbackConstructor(ON_AFTER_BOOTSTRAP_HOOK)
	serveEventCallbackConstructor(ON_BEFORE_SERVE_HOOK)
	apiErrorEventCallbackConstructor(ON_API_ERROR_HOOK)

	// Db Hooks
	modelEventCallbackConstructor(ON_MODEL_BEFORE_CREATE_HOOK)
	modelEventCallbackConstructor(ON_MODEL_AFTER_CREATE_HOOK)
	modelEventCallbackConstructor(ON_MODEL_BEFORE_UPDATE_HOOK)
	modelEventCallbackConstructor(ON_MODEL_AFTER_UPDATE_HOOK)
	modelEventCallbackConstructor(ON_MODEL_BEFORE_DELETE_HOOK)
	modelEventCallbackConstructor(ON_MODEL_AFTER_DELETE_HOOK)

	// Mailer Hooks
	mailerAdminEventCallbackConstructor(ON_MAILER_BEFORE_ADMIN_RESET_PASSWORD_SEND_HOOK)
	mailerAdminEventCallbackConstructor(ON_MAILER_AFTER_ADMIN_RESET_PASSWORD_SEND_HOOK)
	mailerRecordEventCallbackConstructor(ON_MAILER_BEFORE_RECORD_RESET_PASSWORD_SEND_HOOK)
	mailerRecordEventCallbackConstructor(ON_MAILER_AFTER_RECORD_RESET_PASSWORD_SEND_HOOK)
	mailerRecordEventCallbackConstructor(ON_MAILER_BEFORE_RECORD_VERIFICATION_SEND_HOOK)
	mailerRecordEventCallbackConstructor(ON_MAILER_AFTER_RECORD_VERIFICATION_SEND_HOOK)
	mailerRecordEventCallbackConstructor(ON_MAILER_BEFORE_RECORD_CHANGE_EMAIL_SEND_HOOK)
	mailerRecordEventCallbackConstructor(ON_MAILER_AFTER_RECORD_CHANGE_EMAIL_SEND_HOOK)

	// Record API Hooks
	recordsListEventCallbackConstructor(ON_RECORDS_LIST_REQUEST_HOOK)
	recordViewEventCallbackConstructor(ON_RECORD_VIEW_REQUEST_HOOK)
	recordCreateEventCallbackConstructor(ON_RECORD_BEFORE_CREATE_REQUEST_HOOK)
	recordCreateEventCallbackConstructor(ON_RECORD_AFTER_CREATE_REQUEST_HOOK)
	recordUpdateEventCallbackConstructor(ON_RECORD_BEFORE_UPDATE_REQUEST_HOOK)
	recordUpdateEventCallbackConstructor(ON_RECORD_AFTER_UPDATE_REQUEST_HOOK)
	recordDeleteEventCallbackConstructor(ON_RECORD_BEFORE_DELETE_REQUEST_HOOK)
	recordDeleteEventCallbackConstructor(ON_RECORD_AFTER_DELETE_REQUEST_HOOK)
	recordAuthEventCallbackConstructor(ON_RECORD_AUTH_REQUEST_HOOK)
	recordUnlinkExternalAuthEventCallbackConstructor(ON_RECORD_BEFORE_UNLINK_EXTERNAL_AUTH_REQUEST_HOOK)
	recordUnlinkExternalAuthEventCallbackConstructor(ON_RECORD_AFTER_UNLINK_EXTERNAL_AUTH_REQUEST_HOOK)
	recordRequestVerificationEventCallbackConstructor(ON_RECORD_BEFORE_REQUEST_VERIFICATION_REQUEST_HOOK)
	recordRequestVerificationEventCallbackConstructor(ON_RECORD_AFTER_REQUEST_VERIFICATION_REQUEST_HOOK)
	recordConfirmVerificationEventCallbackConstructor(ON_RECORD_BEFORE_CONFIRM_VERIFICATION_REQUEST_HOOK)
	recordConfirmVerificationEventCallbackConstructor(ON_RECORD_AFTER_CONFIRM_VERIFICATION_REQUEST_HOOK)
	recordRequestPasswordResetEventCallbackConstructor(ON_RECORD_BEFORE_REQUEST_PASSWORD_RESET_REQUEST_HOOK)
	recordRequestPasswordResetEventCallbackConstructor(ON_RECORD_AFTER_REQUEST_PASSWORD_RESET_REQUEST_HOOK)
	recordConfirmPasswordResetEventCallbackConstructor(ON_RECORD_BEFORE_CONFIRM_PASSWORD_RESET_REQUEST_HOOK)
	recordConfirmPasswordResetEventCallbackConstructor(ON_RECORD_AFTER_CONFIRM_PASSWORD_RESET_REQUEST_HOOK)
	recordRequestEmailChangeEventCallbackConstructor(ON_RECORD_BEFORE_REQUEST_EMAIL_CHANGE_REQUEST_HOOK)
	recordRequestEmailChangeEventCallbackConstructor(ON_RECORD_AFTER_REQUEST_EMAIL_CHANGE_REQUEST_HOOK)
	recordConfirmEmailChangeEventCallbackConstructor(ON_RECORD_BEFORE_CONFIRM_EMAIL_CHANGE_REQUEST_HOOK)
	recordConfirmEmailChangeEventCallbackConstructor(ON_RECORD_AFTER_CONFIRM_EMAIL_CHANGE_REQUEST_HOOK)

	// Realtime API Hooks
	realtimeConnectEventCallbackConstructor(ON_REALTIME_CONNECT_REQUEST_HOOK)
	realtimeDisconnectEventCallbackConstructor(ON_REALTIME_DISCONNECT_REQUEST_HOOK)
	realtimeSubscribeEventCallbackConstructor(ON_REALTIME_BEFORE_SUBSCRIBE_REQUEST_HOOK)
	realtimeSubscribeEventCallbackConstructor(ON_REALTIME_AFTER_SUBSCRIBE_REQUEST_HOOK)
	realtimeMessageEventCallbackConstructor(ON_REALTIME_BEFORE_MESSAGE_SEND_HOOK)
	realtimeMessageEventCallbackConstructor(ON_REALTIME_AFTER_MESSAGE_SEND_HOOK)

	// File API Hooks
	fileDownloadEventCallbackConstructor(ON_FILE_DONWLOAD_REQUEST_HOOK)

	// Collection API Hooks
	collectionsListEventCallbackConstructor(ON_COLLECTION_LIST_REQUEST_HOOK)
	collectionViewEventCallbackConstructor(ON_COLLECTION_VIEW_REQUEST_HOOK)
	collectionCreateEventCallbackConstructor(ON_COLLECTION_BEFORE_CREATE_REQUEST_HOOK)
	collectionCreateEventCallbackConstructor(ON_COLLECTION_AFTER_CREATE_REQUEST_HOOK)
	collectionUpdateEventCallbackConstructor(ON_COLLECTION_BEFORE_UPDATE_REQUEST_HOOK)
	collectionUpdateEventCallbackConstructor(ON_COLLECTION_AFTER_UPDATE_REQUEST_HOOK)
	collectionDeleteEventCallbackConstructor(ON_COLLECTION_BEFORE_DELETE_REQUEST_HOOK)
	collectionDeleteEventCallbackConstructor(ON_COLLECTION_AFTER_DELETE_REQUEST_HOOK)
	collectionsImportEventCallbackConstructor(ON_COLLECTIONS_BEFORE_IMPORT_REQUEST_HOOK)
	collectionsImportEventCallbackConstructor(ON_COLLECTIONS_AFTER_IMPORT_REQUEST_HOOK)

	// Settings API Hooks
	settingsListEventCallbackConstructor(ON_SETTINGS_LIST_REQUEST_HOOK)
	settingsUpdateEventCallbackConstructor(ON_SETTINGS_BEFORE_UPDATE_REQUEST_HOOK)
	settingsUpdateEventCallbackConstructor(ON_SETTINGS_AFTER_UPDATE_REQUEST_HOOK)

	// Admin API Hooks
	adminsListEventCallbackConstructor(ON_ADMINS_LIST_REQUEST_HOOK)
	adminViewEventCallbackConstructor(ON_ADMIN_VIEW_REQUEST_HOOK)
	adminCreateEventCallbackConstructor(ON_ADMIN_BEFORE_CREATE_REQUEST_HOOK)
	adminCreateEventCallbackConstructor(ON_ADMIN_AFTER_CREATE_REQUEST_HOOK)
	adminUpdateEventCallbackConstructor(ON_ADMIN_BEFORE_UPDATE_REQUEST_HOOK)
	adminUpdateEventCallbackConstructor(ON_ADMIN_AFTER_UPDATE_REQUEST_HOOK)
	adminDeleteEventCallbackConstructor(ON_ADMIN_BEFORE_DELETE_REQUEST_HOOK)
	adminDeleteEventCallbackConstructor(ON_ADMIN_AFTER_DELETE_REQUEST_HOOK)
	adminAuthEventCallbackConstructor(ON_ADMIN_AUTH_REQUEST_HOOK)

	// type TransactionApi struct {
	// 	Execute func(sql string)
	// }
	// __go_apis.Set("withTransaction", func(cb func(e *TransactionApi)) {
	// 	app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
	// 		var api = TransactionApi{
	// 			Execute: func(sql string) error {
	// 				res, err := txDao.DB().Select().NewQuery(sql).Execute()
	// 				if err != nil {
	// 					return err
	// 				}

	// 			}}
	// 	})

	// })

	// Expose API Middlewares
	__go_apis.Set("RequireGuestOnly", apis.RequireGuestOnly)
	__go_apis.Set("RequireRecordAuth", apis.RequireRecordAuth)
	__go_apis.Set("RequireSameContextRecordAuth", apis.RequireSameContextRecordAuth)
	__go_apis.Set("RequireAdminAuth", apis.RequireAdminAuth)
	__go_apis.Set("RequireAdminAuthOnlyIfAny", apis.RequireAdminAuthOnlyIfAny)
	__go_apis.Set("RequireAdminOrRecordAuth", apis.RequireAdminOrRecordAuth)
	__go_apis.Set("RequireAdminOrOwnerAuth", apis.RequireAdminOrOwnerAuth)
	__go_apis.Set("LoadAuthContext", apis.LoadAuthContext)
	__go_apis.Set("LoadCollectionContext", apis.LoadCollectionContext)
	__go_apis.Set("AcitivityLogger", apis.ActivityLogger)

	// Expose API Errors
	__go_apis.Set("NewApiError", apis.NewApiError)
	__go_apis.Set("NewNotFoundError", apis.NewNotFoundError)
	__go_apis.Set("NewBadRequestError", apis.NewBadRequestError)
	__go_apis.Set("NewForbiddenError", apis.NewForbiddenError)
	__go_apis.Set("NewUnauthorizedError", apis.NewUnauthorizedError)

	// Expose App
	__go_apis.Set("app", app)

	// Expose ping for Q&A Check
	__go_apis.Set("ping", func() string {
		return "Hello from Go!"
	})
}

func loadActiveScript() (string, error) {

	collection, err := app.Dao().FindCollectionByNameOrId("pbscript")
	if err != nil {
		return "", err
	}
	recs, err := app.Dao().FindRecordsByExpr(collection, dbx.HashExp{"type": "script", "isActive": true})
	if err != nil {
		return "", err
	}
	if len(recs) > 1 {
		return "", fmt.Errorf("expected one active script record but got %d", len(recs))
	}
	if len(recs) == 0 {
		return "", nil // Empty script
	}
	rec := recs[0]
	jsonData := rec.GetStringDataValue("data")
	type Data struct {
		Source string `json:"source"`
	}
	var json_map Data
	err = json.Unmarshal([]byte(jsonData), &json_map)
	if err != nil {
		return "", err
	}

	script := json_map.Source
	fmt.Printf("Script has been loaded.\n")
	return script, nil

}

func reloadVm() error {
	fmt.Println("Initializing PBScript engine")
	vm = goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())

	// Clean up all handlers
	fmt.Println("Executing cleanups")
	for i := 0; i < len(cleanups); i++ {
		cleanups[i]()
	}
	cleanups = nil

	// Load the main script
	fmt.Println("Loading JS")
	script, err := loadActiveScript()
	if err != nil {
		return err
	}

	// Console proxy
	fmt.Println("Creating console proxy")
	console := vm.NewObject()
	console.Set("log", func(s ...goja.Value) {
		for _, v := range s {
			fmt.Printf("%s ", v.String())
		}
		fmt.Print("\n")
	})
	vm.Set("console", console)

	fmt.Println("Creating apis proxy")
	bindApis()
	vm.Set("__go", __go_apis)

	fmt.Println("Go initialization complete. Running script.")
	source := fmt.Sprintf(`
console.log('Top of PBScript bootstrap')
let __jsfuncs = {ping: ()=>'Hello from PBScript!'}
function registerJsFuncs(funcs) {
__jsfuncs = {__jsfuncs, ...funcs }
}
%s
console.log('Pinging Go')
console.log('Pinging Go succeeded with:', __go.ping())
console.log('Bottom of PBScript bootstrap')
`, script)
	_, err = vm.RunString(source)
	if err != nil {
		return err
	}

	// js api  wireup
	fmt.Println("Wiring up JS API")
	type S struct {
		Ping func() (string, *goja.Exception) `json:"ping"`
	}
	jsFuncs := S{}
	err = vm.ExportTo(vm.Get("__jsfuncs"), &jsFuncs)
	if err != nil {
		return err
	}

	{
		fmt.Println("Pinging JS")
		res, err := jsFuncs.Ping()
		if err != nil {
			return fmt.Errorf("ping() failed with %s", err.Value().Export())
		} else {
			fmt.Printf("Ping succeeded with: %s\n", res)
		}
	}
	return nil
}

func migrate() error {
	fmt.Println("Finding collection")
	_, err := app.Dao().FindCollectionByNameOrId("anything")
	fmt.Println("Finished collection")
	if err != nil {
		err = app.Dao().SaveCollection(&models.Collection{
			Name: "pbscript",
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Type: schema.FieldTypeText,
					Name: "type",
				},
				&schema.SchemaField{
					Type: schema.FieldTypeBool,
					Name: "isActive",
				},
				&schema.SchemaField{
					Type: schema.FieldTypeJson,
					Name: "data",
				},
			),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func watchForScriptChanges() {
	app.OnModelAfterUpdate().Add(func(e *core.ModelEvent) error {
		if e.Model.TableName() == "pbscript" {
			reloadVm()
		}
		return nil
	})

	app.OnModelAfterCreate().Add(func(e *core.ModelEvent) error {
		if e.Model.TableName() == "pbscript" {
			reloadVm()
		}
		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// add new "GET /api/hello" route

		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/pbscript/deploy",
			Handler: func(c echo.Context) error {
				json_map := make(map[string]interface{})
				err := json.NewDecoder(c.Request().Body).Decode(&json_map)
				if err != nil {
					return err
				}
				//json_map has the JSON Payload decoded into a map
				src := json_map["source"]

				err = app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
					fmt.Println("Deactivating active script")
					_, err := txDao.DB().
						NewQuery("UPDATE pbscript SET isActive=false WHERE type='script'").Execute()
					if err != nil {
						return err
					}

					fmt.Println("Packaging new record data")
					bytes, err := json.Marshal(dbx.Params{"source": src})
					if err != nil {
						return err
					}
					_json := string(bytes)

					fmt.Println("Saving new model")
					collection, err := txDao.FindCollectionByNameOrId("pbscript")
					if err != nil {
						return err
					}
					record := models.NewRecord(collection)
					record.SetDataValue("type", "script")
					record.SetDataValue("isActive", "true")
					record.SetDataValue("data", _json)
					err = txDao.SaveRecord(record)
					if err != nil {
						return err
					}
					fmt.Println(("Record saved"))
					// _, err = txDao.DB().
					// 	NewQuery("INSERT INTO pbscript (type,isActive,data) values ('script', true, {data})").Bind(dbx.Params{"data": _json}).Execute()
					// if err != nil {
					// 	return err
					// }
					return nil
				})
				if err != nil {
					return err
				}
				return c.String(http.StatusOK, "ok")

			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireAdminAuth(),
			},
		})

		return nil
	})
}

func initEvents() {

	app.OnModelBeforeCreate().Add(func(e *core.ModelEvent) error {
		fmt.Println("event: OnModelBeforeCreate")
		event.Fire(event.EVT_ON_MODEL_BEFORE_CREATE, (*event.UnknownPayload)(unsafe.Pointer(e)))
		return nil
	})
	app.OnModelAfterCreate().Add(func(e *core.ModelEvent) error {
		fmt.Println("event: OnModelAfterCreate")
		event.Fire(event.EVT_ON_MODEL_AFTER_CREATE, (*event.UnknownPayload)(unsafe.Pointer(e)))
		return nil
	})
}

func StartPBScript(_app *pocketbase.PocketBase) error {
	app = _app

	watchForScriptChanges()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		migrate()
		initEvents()
		router = e.Router
		err := reloadVm()
		if err != nil {
			logErrorf("Error loading VM: %s\n", err)
		}
		return nil
	})

	return nil

}
