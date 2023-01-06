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

var (
	addr      *string
	mms       *string
	support   *string
	accendent *string
)

func init() {
	addr = flag.String("addr", ":8282", "Сетевой адрес HTTP")
	mms = flag.String("mms", "http://127.0.0.1:8383/mms", "путь к данным MMS")
	support = flag.String("support", "http://127.0.0.1:8383/support", "путь к данным Support")
	accendent = flag.String("accendent", "http://127.0.0.1:8383/accendent", "путь к данным Incidents")
	flag.Parse()
}

func main() {

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

var resultT ResultT

func handleConnection(w http.ResponseWriter, r *http.Request) {
	rst := getResultData()
	checkingStructure(rst)

	response, _ := json.Marshal(resultT)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func checkingStructure(rst ResultSetT) ResultT {
	//rst.SMS = nil //для проверки ошибки записи в ResultT

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

func getResultData() ResultSetT {
	//коллекция SMS - [][]слайс SMS сообщений
	pathFileSMS := "simulator/sms.data"
	fileSMS := openAndReadCSV(pathFileSMS)
	SMS := GetSMScollection(fileSMS)

	//коллекция MMS - [][]слайс MMS сообщений
	pathURL := mms
	UrlMMS := parsingMMS(*pathURL)
	MMS := GetMMScollection(UrlMMS)

	//коллекция VoiceCall - []слайс Voice сообщений
	pathFileVoice := "simulator/voice.data"
	fileVoice := openAndReadCSV(pathFileVoice)
	VoiceCall := GetVoiceCollection(fileVoice)

	//коллекция EmailData - map[string][][] карта Email сообщений
	pathFileEmail := "simulator/email.data"
	fileEmail := openAndReadCSV(pathFileEmail)
	Email := GetEmailCollection(fileEmail)

	//коллекция BillingData - структура BillingData
	pathFileBilling := "simulator/billing.data"
	Billing := GetBillingCollection(pathFileBilling)

	//коллекция Support - []int слайс нагрузки и времени ответа
	pathURLsupport := support
	UrlSupport := parsingSupport(*pathURLsupport)
	Support := GetSupportCollection(UrlSupport)

	//коллекция Incident - []слайс Incident сообщений
	pathURLincident := accendent
	Incidents := GetIncidentsCollection(*pathURLincident)

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
