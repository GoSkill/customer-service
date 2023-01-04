package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var resultT ResultT

func handleConnection(w http.ResponseWriter, r *http.Request) {
	rst := getResultData()
	checkingStructure(rst)

	response, _ := json.Marshal(resultT)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func checkingStructure(rst ResultSetT) ResultT {
	//rst.SMS = nil //для проверки ошибки

	if rst.SMS == nil || rst.MMS == nil || rst.VoiceCall == nil ||
		rst.Email == nil || rst.Support == nil || rst.Incidents == nil {
		resultT.Error = "Error on collect data"
		fmt.Printf("ошибка сбора данных:\n %v\n %v\n", resultT.Status, resultT.Error)
		return resultT
	}
	resultT.Status = true
	resultT.Data = rst
	resultT.Error = ""
	return resultT
}

func main() {
	addr := flag.String("addr", ":8282", "Сетевой адрес HTTP")
	flag.Parse()

	r := mux.NewRouter()

	r.HandleFunc("/test", handleConnection)

	srv := &http.Server{
		Addr:         *addr,
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Сервер слушает на 127.0.0.1%s", srv.Addr)
	srv.ListenAndServe()
}

func getResultData() ResultSetT {
	//коллекция SMS - [][]слайс слайсов SMS сообщений
	pathFileSMS := "simulator/sms.data"
	fileSMS := openAndReadCSV(pathFileSMS)
	SMS := GetSMScollection(fileSMS)
	//fmt.Printf("SMS %v\n", SMS)

	//коллекция MMS - [][]слайс слайсов SMS сообщений
	pathURL := "http://127.0.0.1:8383/mms"
	UrlMMS := parsingMMS(pathURL)
	MMS := GetMMScollection(UrlMMS)
	//fmt.Printf("MMS %v\n", MMS)

	//коллекция VoiceCall - []слайс Voice сообщений
	pathFileVoice := "simulator/voice.data"
	fileVoice := openAndReadCSV(pathFileVoice)
	VoiceCall := GetVoiceCollection(fileVoice)
	//fmt.Printf("VoiceCall %v\n", VoiceCall)

	//коллекция EmailData - map[string][][] карта слайс слайсов Email сообщений
	pathFileEmail := "simulator/email.data"
	fileEmail := openAndReadCSV(pathFileEmail)
	Email := GetEmailCollection(fileEmail)
	//fmt.Printf("Email %v\n", Email)

	//коллекция BillingData - структура BillingData
	pathFileBilling := "simulator/billing.data"
	Billing := GetBillingCollection(pathFileBilling)
	//fmt.Printf("Billing %v\n", Billing)

	//коллекция Support - []int слайс нагрузки и времени ответа
	pathURLsupport := "http://127.0.0.1:8383/support"
	UrlSupport := parsingSupport(pathURLsupport)
	Support := GetSupportCollection(UrlSupport)
	//fmt.Printf("Support %v\n", Support)

	//коллекция Incident - []слайс Incident сообщений
	pathURLincident := "http://127.0.0.1:8383/accendent"
	Incidents := GetIncidentsCollection(pathURLincident)
	//fmt.Printf("Incidents %v\n", Incidents)

	rst := ResultSetT{
		SMS:       SMS,
		MMS:       MMS,
		VoiceCall: VoiceCall,
		Email:     Email,
		Billing:   Billing,
		Support:   Support,
		Incidents: Incidents,
	}
	return rst
}
