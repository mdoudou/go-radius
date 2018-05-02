package server

import (
	"bytes"
	"encoding/binary"
	"github.com/hel2o/go-radius/db"
	"github.com/hel2o/go-radius/g"
	"github.com/hel2o/go-radius/http"
	"github.com/hel2o/radius"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type radiusService struct{}

func StartRadius() {
	aut := radius.NewServer(g.Config().GoRadius.AuthListen, g.Config().GoRadius.SharedKey, radiusService{})
	acc := radius.NewServer(g.Config().GoRadius.AcctListen, g.Config().GoRadius.SharedKey, radiusService{})

	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	errChan := make(chan error)
	go func() {
		log.Println("Listening UDP on 0.0.0.0", g.Config().GoRadius.AuthListen)
		err := aut.ListenAndServe()
		if err != nil {
			errChan <- err
		}
	}()
	go func() {
		log.Println("Listening UDP on 0.0.0.0", g.Config().GoRadius.AcctListen)
		err := acc.ListenAndServe()
		if err != nil {
			errChan <- err
		}
	}()
	go http.Start()
	select {
	case <-signalChan:
		aut.Stop()
		log.Println("Stopping authentication server ")

		acc.Stop()
		log.Println("Stopping accounting server")

	case err := <-errChan:
		log.Println(err.Error())
	}
}

func (p radiusService) RadiusHandle(request *radius.Packet) *radius.Packet {

	npac := request.Reply()
	userName := request.GetUsername()
	password := request.GetPassword()
	nasIPAddress := request.GetNasIpAddress().String()
	framedIPAddress := request.GetFramedIPAddress().String()
	nasIdentifier := request.GetNASIdentifier()
	acctSessionId := request.GetAcctSessionId()

	switch request.Code {
	case radius.AccessRequest:
		if db.CheckUserPassword(db.RadiusDb, request.GetUsername(), request.GetPassword()) {
			npac.Code = radius.AccessAccept
			// add Vendor-specific attribute - Vendor Huawei (code 2011) Attribute Huawei-Exec-Privilege (code 29)
			npac.AddVSA(radius.VSA{Vendor: 2011, Type: 29, Value: intToBytes(3)})
			npac.AddAVP(radius.AVP{Type: radius.ServiceType, Value: intToBytes(1)})
			npac.AddAVP(radius.AVP{Type: radius.LoginService, Value: intToBytes(0)})
			npac.AddAVP(radius.AVP{Type: radius.ReplyMessage, Value: []byte("Authentication success - by ISC Team")})
			if acctSessionId != "" {
				db.AuthSuccess(db.RadiusDb, userName, password, nasIPAddress, nasIdentifier, framedIPAddress, acctSessionId)
			} else {
				log.Println(nasIPAddress, "acctSessionId empty")
			}
			log.Println(userName, nasIPAddress, "auth success")

		} else {
			db.AuthFail(db.RadiusDb, userName, password, nasIPAddress, nasIdentifier, framedIPAddress, acctSessionId)
			npac.Code = radius.AccessReject
			npac.AddAVP(radius.AVP{Type: radius.ReplyMessage, Value: []byte("Authentication failed - by ISC Team")})
			log.Println(userName, nasIPAddress, password, "auth failed")
		}
	case radius.AccountingRequest:
		if request.GetAcctStatusType().String() == "Start" {
			db.Login(db.RadiusDb, userName, nasIPAddress, nasIdentifier, framedIPAddress, acctSessionId)
			log.Println(userName, nasIPAddress, "Login")
		}
		if request.GetAcctStatusType().String() == "Stop" {
			db.Logout(db.RadiusDb, framedIPAddress, acctSessionId)
			log.Println(userName, nasIPAddress, "Logout")
		}
		npac.Code = radius.AccountingResponse
	default:
		npac.Code = radius.AccessReject
	}
	return npac
}

//整形转换成字节
func intToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}
