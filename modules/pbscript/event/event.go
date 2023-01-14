package event

import "fmt"

type UnknownPayload any

var inc = 0
var events = map[string]map[int]func(payload *UnknownPayload){}

var HOOKS = [69]string{
	// App Hooks
	"OnBeforeBootstrap",
	"OnAfterBootstrap",
	"OnBeforeServe",
	"OnBeforeApiError",
	// Db Books
	"OnModelBeforeCreate",
	"OnModelAfterCreate",
	"OnModelBeforeUpdate",
	"OnModelAfterUpdate",
	"OnModelBeforeDelete",
	"OnModelAfterDelete",
	// Mailer Hooks
	"OnMailerBeforeAdminResetPasswordSend",
	"OnMailerAfterAdminResetPasswordSend",
	"OnMailerBeforeRecordResetPasswordSend",
	"OnMailerAfterRecordResetPasswordSend",
	"OnMailerBeforeRecordVerificationSend",
	"OnMailerBeforeRecordChangeEmailSend",
	"OnMailerAferRecordChangeEmailSend",
	// Record API Hooks
	"OnRecordsListRequest",
	"OnRecordViewRequest",
	"OnRecordBeforeCreateRequest",
	"OnRecordAfterCreateRequest",
	"OnRecordBeforeUpdateRequest",
	"OnRecordBeforeDeleteRequest",
	"OnRecordAfterDeleteRequest",
	"OnRecordAuthRequest",
	"OnRecordListExternalAuthsRequest",
	"OnRecordBeforeUnlinkExternalAuthRequest",
	"OnRecordAfterUnlinkExternalAuthRequest",
	"OnRecorBeforeRequestVerificationRequest",
	"OnRecordAfterRequestVerificationRequest",
	"OnRecordBeforeConfirmVerificationRequest",
	"OnRecordAfterConfirmVerificationRequest",
	"OnRecordBeforeRequestPasswordResetRequest",
	"OnRecordAfterRequestPasswordResetRequest",
	"OnRecordBeforeConfirmPasswordResetRequest",
	"OnRecordAfterConfirmPasswordResetRequest",
	"OnRecordBeforeRequestEmailChangeRequest",
	"OnRecordAfterRquestEmailChangeRequest",
	"OnRecordBeforeConfirmEmailChangeRequest",
	"OnRecordAfterConfirmEmailChangeRequest",
	// Realtime API Hooks
	"OnRealtimeConnectRequest",
	"OnRealtimeDisconnectRequest",
	"OnRealtimeBeforeSubscribeRequest",
	"OnRealtimeAfterSubscribeRequest",
	"OnRealtimeBeforeMessageSend",
	"OnRealtimeAfterMessageSend",
	// File API Hooks
	"OnFileDownloadRequest",
	// Collection API Hooks
	"OnCollectionListRequest",
	"OnCollectionViewRequest",
	"OnCollectionBeforeCreateRequest",
	"OnCollectionAfterCreateRequest",
	"OnCollectionBeforeUpdateRequest",
	"OnCollectionAfterUpdateRequest",
	"OnCollectionBeforeDeleteRequest",
	"OnCollectionAfterDeleteRequest",
	"OnCollectionBeforeImportRequest",
	"OnCollectionAfterImportRequest",
	// Settings API Hooks
	"OnSettingsListRequest",
	"OnSettingsBeforeUpdateRequest",
	"OnSettingsAfterUpdateRequest",
	// Admin API Hooks
	"OnAdminsListRequest",
	"OnAdminViewRequest",
	"OnAdminBeforeCreateRequest",
	"OnAdminAfterCreateRequest",
	"OnAdminBeforeUpdateRequest",
	"OnAdminAfterUpdateRequest",
	"OnAdminBeforeDeleteRequest",
	"OnAdminAfterDeleteRequest",
	"OnAdminAuthRequest",
}

func isValid(eventName string) bool {
	for _, element := range HOOKS {
		if element == eventName {
			return true
		}
	}
	return false
}

func ensureEvent(eventName string) {
	if !isValid((eventName)) {
		panic(fmt.Sprintf("%s is not a valid event name", eventName))
	}
	if _, ok := events[eventName]; !ok {
		fmt.Printf("Creating collection for %s\n", eventName)
		events[eventName] = make(map[int]func(payload *UnknownPayload))
	}
}

func On(eventName string, cb func(payload *UnknownPayload)) func() {
	ensureEvent(eventName)

	inc++
	idx := inc
	events[eventName][idx] = cb
	fmt.Printf("Adding %d to %s\n", idx, eventName)
	return func() {
		delete(events[eventName], idx)
	}
}

func Fire(eventName string, payload *UnknownPayload) {
	ensureEvent(eventName)

	fmt.Printf("Firing %s\n", eventName)
	for fnId, v := range events[eventName] {
		fmt.Printf("Dispatching %s to %d\n", eventName, fnId)
		v(payload)
	}
}
