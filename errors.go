package sentrytemporal

import (
	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	"go.temporal.io/sdk/workflow"
)

func isContinueAsNewError(err error) bool {
	var continueAsNewErr *workflow.ContinueAsNewError
	return errors.As(err, &continueAsNewErr)
}

func prepareSentryReport(err error) *sentry.Event {
	event, extraDetails := errors.BuildSentryReport(err)

	for extraKey, extraValue := range extraDetails {
		event.Extra[extraKey] = extraValue
	}

	// Avoid leaking the machine's hostname by injecting the literal "<redacted>".
	// Otherwise, sentry.Client.Capture will see an empty ServerName field and
	// automatically fill in the machine's hostname.
	event.ServerName = "<redacted>"

	tags := map[string]string{
		"report_type": "error",
	}
	for key, value := range tags {
		event.Tags[key] = value
	}
	return event
}
