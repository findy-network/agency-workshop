package agent

import (
	"log"
	"os"
	"strconv"

	"github.com/findy-network/findy-agent-auth/acator/authn"
	"github.com/findy-network/findy-common-go/agency/client"
	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"google.golang.org/grpc"
)

const (
	subCmdLogin    = "login"
	subCmdRegister = "register"
)

type AgencyClient struct {
	Conn           client.Conn
	AgentClient    agency.AgentServiceClient
	ProtocolClient agency.ProtocolServiceClient
}

func execAuthCmd(cmd string) (res authn.Result, err error) {
	defer err2.Handle(&err)

	myCmd := authn.Cmd{
		SubCmd:   subCmdLogin,
		UserName: os.Getenv("FCLI_USER"),
		Url:      os.Getenv("FCLI_URL"),
		AAGUID:   "12c85a48-4baf-47bd-b51f-f192871a1511",
		Key:      os.Getenv("FCLI_KEY"),
		Counter:  0,
		Token:    "",
		Origin:   os.Getenv("FCLI_ORIGIN"),
	}

	myCmd.SubCmd = cmd

	try.To(myCmd.Validate())

	return myCmd.Exec(os.Stdout)
}

func LoginAgent() (
	agencyClient *AgencyClient,
	err error,
) {
	defer err2.Handle(&err)

	// first try to login
	res, firstTryErr := execAuthCmd(subCmdLogin)
	if firstTryErr != nil {
		// if login fails, try to register and relogin
		_ = try.To1(execAuthCmd(subCmdRegister))
		res = try.To1(execAuthCmd(subCmdLogin))
	}

	log.Println("Agent login succeeded")

	token := res.Token
	// set up API connection
	conf := client.BuildClientConnBase(
		os.Getenv("FCLI_TLS_PATH"),
		os.Getenv("AGENCY_API_SERVER"),
		try.To1(strconv.Atoi(os.Getenv("AGENCY_API_SERVER_PORT"))),
		[]grpc.DialOption{},
	)

	conn := client.TryAuthOpen(token, conf)

	return &AgencyClient{
		Conn:           conn,
		AgentClient:    agency.NewAgentServiceClient(conn),
		ProtocolClient: agency.NewProtocolServiceClient(conn),
	}, nil
}
