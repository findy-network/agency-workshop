package agent

import (
	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
)

// Default implementation for agency listener
type DefaultListener struct{}

func (d *DefaultListener) HandleNewConnection(
	*agency.Notification,
	*agency.ProtocolStatus_DIDExchangeStatus,
) {
}

func (d *DefaultListener) HandleBasicMesssageDone(
	*agency.Notification,
	*agency.ProtocolStatus_BasicMessageStatus,
) {
}

func (d *DefaultListener) HandleIssueCredentialDone(
	*agency.Notification,
	*agency.ProtocolStatus_IssueCredentialStatus,
) {
}

func (d *DefaultListener) HandlePresentProofPaused(
	*agency.Notification,
	*agency.ProtocolStatus_PresentProofStatus,
) {
}

func (d *DefaultListener) HandlePresentProofDone(
	*agency.Notification,
	*agency.ProtocolStatus_PresentProofStatus,
) {
}
