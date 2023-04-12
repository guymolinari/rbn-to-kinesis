package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	sessionKey    string
	notFoundCache map[string]struct{}
)

const (
	httppostUrl = "https://xmldata.qrz.com/xml/current/"
	username    = "N7ZG"
	password    = "tempest"
)

func GetCallFromQRZ(call string) (*QRZDatabase, error) {

	var qrz *QRZDatabase
	var err error
	if notFoundCache == nil {
		notFoundCache = make(map[string]struct{})
	}

	if sessionKey == "" {
		login()
	}

	if _, found := notFoundCache[call]; found {
		return nil, fmt.Errorf("Ignoring, %s in not found cache.", call)
	}

	qrz, err = tryCall(call)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Not found") {
			notFoundCache[call] = struct{}{}
			return nil, err
		}
		if err.Error() != "Session Timeout" {
			return nil, err
		}
		log.Printf("Session timed out, logging in again.")
	}
	if qrz == nil || qrz.Key == "" {
		login()
		qrz, err = tryCall(call)
		if err == nil {
			return qrz, nil
		}
	} else {
		return qrz, nil
	}
	return nil, err
}

// lookup call via QRZ API
func tryCall(call string) (*QRZDatabase, error) {

	params := make(map[string]string)
	params["s"] = sessionKey
	params["callsign"] = call
	return QRZAPI(params)
}

func login() {

	// Log into QRZ
	params := make(map[string]string)
	params["username"] = username
	params["password"] = password
	qrz, err := QRZAPI(params)
	if err != nil {
		log.Fatal(err)
	}
	if qrz.Key == "" {
		log.Fatal(fmt.Errorf("NO KEY QRZ ERR: %v", qrz.Error))
	}
	sessionKey = qrz.Key
	return
}

func QRZAPI(params map[string]string) (*QRZDatabase, error) {

	request, err := http.NewRequest("GET", httppostUrl, nil)
	if err != nil {
		return nil, err
	}

	q := request.URL.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	request.URL.RawQuery = q.Encode()
	//fmt.Println(request.URL.String())

	client := &http.Client{Timeout: time.Second * 2}
	var response *http.Response
	response, err = client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	//fmt.Println("response Status:", response.Status)
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("%v", response.Status)
	}
	//fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	//fmt.Println("response Body:", string(body))
	qrz := &QRZDatabase{}
	err = xml.Unmarshal(body, &qrz)
	if err != nil {
		return nil, err
	}
	//log.Printf("RESP = %#v", qrz)
	if qrz != nil && qrz.Error != "" {
		return nil, fmt.Errorf(qrz.Error)
	}
	if qrz != nil && qrz.Expdate == "0000-00-00" {
		qrz.Expdate = ""
	}
	if qrz != nil && qrz.Efdate == "0000-00-00" {
		qrz.Efdate = ""
	}
	return qrz, nil

}

type QRZDatabase struct {
	// Session
	Key     string `xml:"Session>Key,omitempty"`
	SubExp  string `xml:"Session>SubExp,omitempty"`
	GMTime  string `xml:"Session>GMTime,omitempty"`
	Count   int    `xml:"Session>Count,omitempty"`
	Message string `xml:"Session>Message,omitempty"`
	Remark  string `xml:"Session>Remark,omitempty"`
	Error   string `xml:"Session>Error,omitempty"`

	//Callsign
	Call    string `xml:"Callsign>call,omitempty" sql:"call"`
	Aliases string `xml:"Callsign>aliases,omitempty" sql:"aliases"`
	Dxcc    string `xml:"Callsign>dxcc,omitempty" sql:"dxcc_id"`
	Fname   string `xml:"Callsign>fname,omitempty" sql:"fname"`
	Name    string `xml:"Callsign>name,omitempty" sql:"lname"`
	Addr1   string `xml:"Callsign>addr1,omitempty" sql:"addr1"`
	Addr2   string `xml:"Callsign>addr2,omitempty" sql:"addr2"`
	State   string `xml:"Callsign>state,omitempty" sql:"state"`
	Zip     string `xml:"Callsign>zip,omitempty" sql:"zip"`
	Country string `xml:"Callsign>country,omitempty" sql:"mail_country"`
	Ccode   string `xml:"Callsign>ccode,omitempty" sql:"country_code"`
	//Lat       float64 `xml:"Callsign>lat,omitempty" sql:"lat"`
	//Lon       float64 `xml:"Callsign>lon,omitempty" sql:"lon"`
	Lat     string `xml:"Callsign>lat,omitempty" sql:"lat"`
	Lon     string `xml:"Callsign>lon,omitempty" sql:"lon"`
	Grid    string `xml:"Callsign>grid,omitempty" sql:"grid"`
	County  string `xml:"Callsign>county,omitempty" sql:"county"`
	Fips    string `xml:"Callsign>fips,omitempty" sql:"fips"`
	Land    string `xml:"Callsign>land,omitempty" sql:"dxcc_country"`
	Efdate  string `xml:"Callsign>efdate,omitempty" sql:"license_issue_date"`
	Expdate string `xml:"Callsign>expdate,omitempty" sql:"license_exp_date"`
	P_call  string `xml:"Callsign>p_call,omitempty" sql:"prev_call"`
	Class   string `xml:"Callsign>class,omitempty" sql:"class"`
	Codes   string `xml:"Callsign>codes,omitempty" sql:"codes"`
	Qslmgr  string `xml:"Callsign>qslmgr,omitempty" sql:"qslmgr"`
	Email   string `xml:"Callsign>email,omitempty" sql:"email"`
	Url     string `xml:"Callsign>url,omitempty"`
	//U_views   int     `xml:"Callsign>u_views,omitempty" sql:"u_views"`
	U_views   string `xml:"Callsign>u_views,omitempty" sql:"u_views"`
	Bio       string `xml:"Callsign>bio,omitempty"`
	Biodate   string `xml:"Callsign>biodate,omitempty"`
	Image     string `xml:"Callsign>image,omitempty"`
	Serial    string `xml:"Callsign>serial,omitempty"`
	Moddate   string `xml:"Callsign>moddate,omitempty" sql:"mod_date"`
	MSA       string `xml:"Callsign>MSA,omitempty" sql:"msa"`
	AreaCode  string `xml:"Callsign>AreaCode,omitempty" sql:"area_code"`
	TimeZone  string `xml:"Callsign>TimeZone,omitempty" sql:"time_zone"`
	GMTOffset string `xml:"Callsign>GMTOffset,omitempty" sql:"gmt_offset"`
	DST       string `xml:"Callsign>DST,omitempty" sql:"dst"`
	Eqsl      string `xml:"Callsign>eqsl,omitempty" sql:"eqsl"`
	Mqsl      string `xml:"Callsign>mqsl,omitempty" sql:"mqsl"`
	Lotw      string `xml:"Callsign>lotw,omitempty" sql:"lotw"`
	Cqzone    string `xml:"Callsign>cqzone,omitempty" sql:"cq_zone"`
	Ituzone   string `xml:"Callsign>ituzone,omitempty" sql:"itu_zone"`
	Geoloc    string `xml:"Callsign>geoloc,omitempty" sql:"geoloc"`
	Attn      string `xml:"Callsign>attn,omitempty" sql:"attn"`
	Nickname  string `xml:"Callsign>nickname,omitempty" sql:"nickname"`
	Name_fmt  string `xml:"Callsign>name_fmt,omitempty" sql:"lname_fmt"`
	//Born      int     `xml:"Callsign>born,omitempty" sql:"born"`
	Born string `xml:"Callsign>born,omitempty" sql:"born"`
}
