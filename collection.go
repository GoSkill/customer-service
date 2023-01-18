package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

// сортировка SMS
type SortedSMSByCountry []SMSData               //SMS по названию страны
func (c SortedSMSByCountry) Len() int           { return len(c) }
func (c SortedSMSByCountry) Less(i, j int) bool { return c[i].Country < c[j].Country }
func (с SortedSMSByCountry) Swap(i, j int)      { с[i], с[j] = с[j], с[i] }

type SortedSMSByProvider []SMSData               //SMS по названию провайдера
func (c SortedSMSByProvider) Len() int           { return len(c) }
func (c SortedSMSByProvider) Less(i, j int) bool { return c[i].Provider < c[j].Provider }
func (с SortedSMSByProvider) Swap(i, j int)      { с[i], с[j] = с[j], с[i] }

// Этап 2. Сбор данных о системе SMS
func GetSMScollection(fileSMS [][]string) [][]SMSData {
	var (
		SMS       SMSData
		ValidSMS  []SMSData
		sortedSMS [][]SMSData
	)
	for _, sms := range fileSMS {
		alfaCode := CodeISOalpha2()
		if len(sms) != 4 { //проверка наличия 4х полей
			continue
		}
		if !ValidData(alfaCode, sms[0]) { //валидность кода страны
			continue
		}
		if !ValidData(Providers, sms[3]) { //валидность провайдера
			continue
		}
		sms0 := sms[0]
		SMS.Country = NameDecoding(sms0)
		SMS.Bandwidth = sms[1]
		SMS.ResponseTime = sms[2]
		SMS.Provider = sms[3]
		ValidSMS = append(ValidSMS, SMS) //срез правильных SMS
	}
	sort.Stable(SortedSMSByProvider(ValidSMS))
	sms0 := make([]SMSData, 0)
	sms0 = append(sms0, ValidSMS...)

	sort.Stable(SortedSMSByCountry(ValidSMS))
	sms1 := make([]SMSData, 0)
	sms1 = append(sms1, ValidSMS...)
	sortedSMS = append(sortedSMS, sms0, sms1)
	return sortedSMS
}

// сортировка MMS
type SortedMMSByCountry []MMSData               //SMS по названию страны
func (c SortedMMSByCountry) Len() int           { return len(c) }
func (с SortedMMSByCountry) Swap(i, j int)      { с[i], с[j] = с[j], с[i] }
func (c SortedMMSByCountry) Less(i, j int) bool { return c[i].Country < c[j].Country }

type SortedMMSByProvider []MMSData               //SMS по названию провайдера
func (c SortedMMSByProvider) Len() int           { return len(c) }
func (с SortedMMSByProvider) Swap(i, j int)      { с[i], с[j] = с[j], с[i] }
func (c SortedMMSByProvider) Less(i, j int) bool { return c[i].Provider < c[j].Provider }

// Этап 3. Сбор данных о системе MMS
func GetMMScollection(fileMMS []MMSData) [][]MMSData {
	var (
		ValidMMS  []MMSData
		sortedMMS [][]MMSData
	)
	for _, mms := range fileMMS {
		alfaCode := CodeISOalpha2()
		if !ValidData(alfaCode, mms.Country) { //валидность кода страны
			continue
		}
		if !ValidData(Providers, mms.Provider) { //валидность провайдера
			continue
		}
		mms0 := mms.Country
		mms.Country = NameDecoding(mms0)
		ValidMMS = append(ValidMMS, mms)
	}
	sort.Stable(SortedMMSByProvider(ValidMMS))
	mms0 := make([]MMSData, 0)
	mms0 = append(mms0, ValidMMS...)

	sort.Stable(SortedMMSByCountry(ValidMMS))
	mms1 := make([]MMSData, 0)
	mms1 = append(mms1, ValidMMS...)
	sortedMMS = append(sortedMMS, mms0, mms1)
	return sortedMMS
}

// Этап 4. Сбор данных о системе Voice Call
func GetVoiceCollection(fileVoice [][]string) []VoiceCallData {
	var (
		Voice      VoiceCallData
		ValidVoice []VoiceCallData
	)
	alfaCode := CodeISOalpha2()
	for _, voice := range fileVoice {
		if len(voice) != 8 {
			continue
		}
		if !ValidData(alfaCode, voice[0]) {
			continue
		}
		if !ValidData(ProviderVoice, voice[3]) {
			continue
		}
		voice4, err := strconv.ParseFloat(voice[4], 32)
		if err != nil {
			continue
		}
		Voice.TTFB, err = strconv.Atoi(voice[5])
		if err != nil {
			continue
		}
		Voice.VoicePurity, err = strconv.Atoi(voice[6])
		if err != nil {
			continue
		}
		voiceEnd := strings.Trim(voice[7], "\n")
		Voice.MedianOfCallsTime, err = strconv.Atoi(voiceEnd)
		if err != nil {
			continue
		}
		Voice.Country = voice[0]
		Voice.Bandwidth = voice[1]
		Voice.ResponseTime = voice[2]
		Voice.Provider = voice[3]
		Voice.ConnectionStability = float32(voice4)
		ValidVoice = append(ValidVoice, Voice)
	}
	return ValidVoice
}

// Этап 5. Сбор данных о системе Email
func GetEmailCollection(fileEmail [][]string) map[string][][]EmailData {
	var (
		email      EmailData
		ValidEmail []EmailData
	)
	alfaCode := CodeISOalpha2()
	for _, e := range fileEmail {
		if len(e) != 3 {
			continue
		}
		if !ValidData(alfaCode, e[0]) {
			continue
		}
		if ValidData(ProviderEmail, e[2]) {
			continue
		}
		emailEnd := strings.Trim(e[2], "\n")
		email2, err := strconv.Atoi(emailEnd)
		if err != nil {
			continue
		}
		email.Country = e[0]
		email.Provider = e[1]
		email.DeliveryTime = email2
		ValidEmail = append(ValidEmail, email)
	}
	return sortingEmailToMap(ValidEmail)
}

// сортировка структуры "EmailData" по полю "DeliveryTime" с выборкой по коду страны "Country"
// и составление карты map[string][][]EmailData
type SortedEmailByDeliveryTime []EmailData

func (c SortedEmailByDeliveryTime) Len() int           { return len(c) }
func (с SortedEmailByDeliveryTime) Swap(i, j int)      { с[i], с[j] = с[j], с[i] }
func (c SortedEmailByDeliveryTime) Less(i, j int) bool { return c[i].DeliveryTime < c[j].DeliveryTime }

func sortingEmailToMap(std []EmailData) map[string][][]EmailData {
	MapEmail := make(map[string][][]EmailData)
	for _, country := range getCountriesList() {
		var MaxMinTime [][]EmailData
		var timeByCountry []EmailData
		for _, str := range std {
			if str.Country != country {
				continue
			}
			timeByCountry = append(timeByCountry, str)
		}
		if timeByCountry == nil {
			continue
		}
		sort.Stable(SortedEmailByDeliveryTime(timeByCountry))
		maxTime := timeByCountry[(len(timeByCountry) - 3):]
		minTime := timeByCountry[:3]
		for i := range MaxMinTime {
			MaxMinTime[i] = make([]EmailData, 3)
		}
		MaxMinTime = append(MaxMinTime, maxTime, minTime)
		MapEmail[country] = MaxMinTime
	}
	return MapEmail
}

// Этап 6. Сбор данных о системе Billing
func GetBillingCollection(file string) BillingData {
	var (
		Billing BillingData
		maskBit string
	)
	fileBilling, err := os.Open(file)
	if err != nil {
		log.Fatalf("ошибка открытия файла: %s", err)
	}
	defer fileBilling.Close()

	data := make([]byte, 64)
	for {
		n, err := fileBilling.Read(data)
		if err == io.EOF {
			break
		}
		maskBit = string(data[:n])
	}
	billing := strings.Split(maskBit, "")
	var bit [6]bool
	var reverseBit []bool
	for i := 0; len(billing) > 0; i++ {
		n := len(billing) - 1
		if billing[n] == "1" {
			bit[i] = true
		}
		reverseBit = append(reverseBit, bit[i])
		billing = billing[:n] //реверс откусыванием
	}
	Billing.CreateCustomer = reverseBit[0]
	Billing.Purchase = reverseBit[1]
	Billing.Payout = reverseBit[2]
	Billing.Recurring = reverseBit[3]
	Billing.FraudControl = reverseBit[4]
	Billing.CheckoutPage = reverseBit[5]
	return Billing
}

// Этап 7. Сбор данных о системе Support
func GetSupportCollection(fileSupport []SupportData) []int {
	var Support []int
	var sum, Level int
	for _, tic := range fileSupport {
		sum += tic.ActiveTickets
	}
	load := sum / len(fileSupport)
	switch {
	case load < 9:
		Level = 1
	case load >= 9 && load <= 16:
		Level = 2
	case load > 16:
		Level = 3
	}
	time := sum * 60 / 18
	Support = append(Support, Level, time)
	return Support
}

// Этап 8. Сбор данных о системе истории инцидентов
type SortedIncidentByStatus []IncidentData

func (c SortedIncidentByStatus) Len() int           { return len(c) }
func (с SortedIncidentByStatus) Swap(i, j int)      { с[i], с[j] = с[j], с[i] }
func (c SortedIncidentByStatus) Less(i, j int) bool { return c[i].Status < c[j].Status }
func GetIncidentsCollection(url string) []IncidentData {
	Incidents := make([]IncidentData, 0)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		log.Printf("получение ошибки, код состояния: %d \nbody: %s\n", resp.StatusCode, resp.Body)
		return Incidents
	}
	if err = json.NewDecoder(resp.Body).Decode(&Incidents); err != nil {
		log.Println("ERROR: " + err.Error())
		return Incidents
	}
	sort.Stable(SortedIncidentByStatus(Incidents))
	return Incidents
}
