package main

import (
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/disney/quanta/shared"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hamba/avro"
	"github.com/reiver/go-telnet"
	"gitlab.disney.com/guys-workspace/rbn-to-kinesis/callparser"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Variables to identify the build
var (
	Version string
	Build   string
)

// Exit Codes
const (
	Success = 0
	Select  = "select * from callsign where call = ?"
	Alias   = "select * from callsign where aliases like ?"
	Insert  = "insert into callsign (%s) values (%s)"
)

// Main strct defines command line arguments variables and various global meta-data associated with record loads.
type Main struct {
	RBNHost    string
	RBNPort    int
	Stream     string
	Region     string
	DBHostPort string
	DBUser     string
	DBSchema   string
	SelectStmt *sql.Stmt
	InsertStmt *sql.Stmt
	AliasStmt  *sql.Stmt
}

// NewMain allocates a new pointer to Main struct with empty record counter
func NewMain() *Main {
	return &Main{}
}

func main() {

	app := kingpin.New(os.Args[0], "RBN to Kinesis Bridge").DefaultEnvars()
	app.Version("Version: " + Version + "\nBuild: " + Build)

	stream := app.Arg("stream", "Kinesis stream name.").Required().String()
	dbHostPort := app.Arg("db-host-port", "Quanta host:port").Required().String()
	dbUser := app.Arg("db-user", "Quanta user").Required().String()
	rbnHost := app.Arg("rbn-host", "Host for RBN endpoint.").Default("telnet.reversebeacon.net").String()
	rbnPort := app.Arg("rbn-port", "Port number for service").Default("7000").Int32()
	rbnClientCall := app.Arg("rbn-client-call", "RBN login call").Default("N7ZG").String()
	region := app.Arg("region", "AWS region").Default("us-east-1").String()
	dbSchema := app.Arg("db-schema", "Quanta database").Default("quanta").String()

	splitex := regexp.MustCompile("[[:space:]]+")
	kingpin.MustParse(app.Parse(os.Args[1:]))

	main := NewMain()
	main.RBNHost = *rbnHost
	main.RBNPort = int(*rbnPort)
	main.Region = *region
	main.Stream = *stream
	main.DBHostPort = *dbHostPort
	main.DBUser = *dbUser
	main.DBSchema = *dbSchema

	log.Printf("RBN host %v.\n", main.RBNHost)
	log.Printf("RBN port %d.\n", main.RBNPort)
	log.Printf("AWS region %s.\n", main.Region)
	log.Printf("Kinesis stream %s.\n", main.Stream)
	log.Printf("DB host:port %s.\n", main.DBHostPort)
	log.Printf("DB user %s.\n", main.DBUser)
	log.Printf("DB schema %s.\n", main.DBSchema)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(main.Region),
	})

	if err != nil {
		log.Fatal(err)
	}

	kc := kinesis.New(sess)
	streamName := aws.String(main.Stream)
	_, err = kc.DescribeStream(&kinesis.DescribeStreamInput{StreamName: streamName})

	//if no stream name in AWS
	if err != nil {
		log.Fatal(err)
	}

	schema, err := avro.Parse(`{
	    "type": "record",
	    "name": "spot_events",
	    "namespace": "quanta",
	    "fields" : [
	        {"name": "band", "type": "string"},
	        {"name": "callsign", "type": "string"},
	        {"name": "de_cont", "type": "string"},
	        {"name": "de_pfx", "type": "string"},
	        {"name": "dx", "type": "string"},
	        {"name": "dx_cont", "type": "string"},
	        {"name": "dx_pfx", "type": "string"},
	        {"name": "freq", "type": "double"},
	        {"name": "mode", "type": "string"},
	        {"name": "tx_mode", "type": "string"},
	        {"name": "db", "type": "int"},
	        {"name": "speed", "type": "int"},
	        {"name": "date", "type": "long"}
	    ]
	}`)

	if err != nil {
		log.Fatal(err)
	}

	conn, err := telnet.DialTo("telnet.reversebeacon.net:7000")
	if nil != err {
		log.Fatal(err)
	}

	log.Printf("Connected.")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:@tcp(%s)/%s", main.DBUser, main.DBHostPort, main.DBSchema))
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	main.SelectStmt, err = db.Prepare(Select)
	if err != nil {
		log.Fatal(err)
	}
	defer main.SelectStmt.Close()

	main.AliasStmt, err = db.Prepare(Alias)
	if err != nil {
		log.Fatal(err)
	}
	defer main.AliasStmt.Close()

	main.InsertStmt, err = db.Prepare(shared.GenerateSQLInsert("callsign", &QRZDatabase{}))
	if err != nil {
		log.Fatal(err)
	}
	defer main.InsertStmt.Close()

	log.Print(ReaderTelnet(conn, "Please enter your call:"))
	WriterTelnet(conn, *rbnClientCall)
	for {
		record := make(map[string]interface{})
		str := ReaderTelnet(conn, "\r\n")
		s := splitex.Split(str, 13)
		if len(s) < 10 || s[2] == "de" {
			continue
		}
		c := strings.Split(s[2], "-")
		call := strings.TrimRight(c[0], "0123456789")

		record["callsign"] = call
		if freq, err := strconv.ParseFloat(s[3], 64); err != nil {
			continue
		} else {
			f := Round(freq, .1)
			i := fmt.Sprintf("%.2f", f)
			record["freq"], _ = strconv.ParseFloat(i, 64)
		}
		record["dx"] = s[4]
		record["mode"] = s[5]
		if strength, err2 := strconv.ParseInt(s[6], 10, 64); err2 != nil {
			log.Printf("STRENGTH ERR %v\n", s[6])
			continue
		} else {
			record["db"] = int(strength)
		}
		if speed, err3 := strconv.ParseInt(s[8], 10, 32); err3 != nil {
			log.Printf("SPEED ERR %v\n", s[8])
			continue
		} else {
			record["speed"] = int(speed)
		}
		record["tx_mode"] = s[10]
		timeStr := s[11]
		if timeStr == "B" {
			timeStr = strings.TrimSpace(s[12])
		}
		if !strings.HasSuffix(timeStr, "Z") || len(timeStr) != 5 {
			log.Printf("TIME ERR [%v] - %v\n", s[11], str)
			continue
		}
		hr := timeStr[0:2]
		mn := timeStr[2:4]
		hri, _ := strconv.ParseInt(hr, 10, 64)
		mni, _ := strconv.ParseInt(mn, 10, 64)
		x := time.Now().UTC()
		y := time.Date(x.Year(), x.Month(), x.Day(), int(hri), int(mni), 0, 0, time.UTC)
		// Adjust for clock skew
		if x.Hour() == 0 && hri == 23 {
			y.AddDate(0, 0, -1)
		}
		if x.Hour() == 23 && hri == 0 {
			y.AddDate(0, 0, 1)
		}
		record["date"] = y.Unix() * 1000
		//log.Printf(">%v", y)

		deRow, err := main.getAndInsertRowForCall(record["callsign"].(string))
		if err != nil {
			log.Fatal(err)
		}

		dxRow, errx := main.getAndInsertRowForCall(record["dx"].(string))
		if errx != nil {
			log.Fatal(errx)
		}
		_ = deRow
		_ = dxRow

		err = Decorate(record)
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		data, err := avro.Marshal(schema, record)
		if err != nil {
			log.Fatal(err)
		}

		// put data to stream
		putOutput, err := kc.PutRecord(&kinesis.PutRecordInput{
			Data:         data,
			StreamName:   aws.String(main.Stream),
			PartitionKey: aws.String(y.Format(time.RFC3339)),
		})

		if err != nil {
			log.Fatal(err)
		}
		_ = *putOutput
		//log.Printf("%v\n", *putOutput)
		//log.Printf("%#v", record)
	}
}

// Thin function reads from Telnet session. "expect" is a string I use as signal to stop reading
func ReaderTelnet(conn *telnet.Conn, expect string) (out string) {
	var buffer [1]byte
	recvData := buffer[:]
	var n int
	var err error

	for {
		n, err = conn.Read(recvData)
		//fmt.Println("Bytes: ", n, "Data: ", recvData, string(recvData))
		if n <= 0 || err != nil || strings.Contains(out, expect) {
			break
		} else {
			out += string(recvData)
		}
	}
	return out
}

// convert a command to bytes, and send to Telnet connection followed by '\r\n'
func WriterTelnet(conn *telnet.Conn, command string) {
	var commandBuffer []byte
	for _, char := range command {
		commandBuffer = append(commandBuffer, byte(char))
	}

	var crlfBuffer [2]byte = [2]byte{'\r', '\n'}
	crlf := crlfBuffer[:]

	//fmt.Println(commandBuffer)

	conn.Write(commandBuffer)
	conn.Write(crlf)
}

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func Decorate(record map[string]interface{}) error {

	de := callparser.NewStation(record["callsign"].(string))
	if de.Valid {
		record["de_pfx"] = de.PrimaryPrefix
		record["de_cont"] = de.Continent
	} else {
		return fmt.Errorf("PrefixMapper: cannot locate prefix for '%s'.", record["callsign"])
	}
	dx := callparser.NewStation(record["dx"].(string))
	if dx.Valid {
		record["dx_pfx"] = dx.PrimaryPrefix
		record["dx_cont"] = dx.Continent
	} else {
		return fmt.Errorf("PrefixMapper: cannot locate prefix for '%s'.", record["dx"])
	}

	freq := record["freq"].(float64)

	if freq >= 1800.0 && freq <= 2000.0 {
		record["band"] = "160m"
	} else if freq >= 3500.0 && freq <= 4000.0 {
		record["band"] = "80m"
	} else if freq >= 5300.0 && freq <= 5500.0 {
		record["band"] = "60m"
	} else if freq >= 7000.0 && freq <= 7300.0 {
		record["band"] = "40m"
	} else if freq >= 10100.0 && freq <= 10150.0 {
		record["band"] = "30m"
	} else if freq >= 14000.0 && freq <= 14300.0 {
		record["band"] = "20m"
	} else if freq >= 18068.0 && freq <= 18168.0 {
		record["band"] = "17m"
	} else if freq >= 21000.0 && freq <= 21450.0 {
		record["band"] = "15m"
	} else if freq >= 24890.0 && freq <= 24990.0 {
		record["band"] = "12m"
	} else if freq >= 28000.0 && freq <= 30000.0 {
		record["band"] = "10m"
	} else if freq >= 50000.0 && freq <= 54000.0 {
		record["band"] = "6m"
	} else if freq >= 69900.0 && freq <= 70500.0 {
		record["band"] = "4m"
	} else if freq >= 472 && freq <= 479 {
		record["band"] = "600m"
	} else if freq >= 135.7 && freq <= 137.8 {
		record["band"] = "2200m"
	} else if freq >= 144000.0 && freq <= 148000.0 {
		record["band"] = "2m"
	} else if freq >= 219000.0 && freq <= 225000.0 {
		record["band"] = "1.25m"
	} else if freq >= 420000.0 && freq <= 450000.0 {
		record["band"] = "70cm"
	} else if freq >= 902000.0 && freq <= 928000.0 {
		record["band"] = "33cm"
	} else if freq >= 1240000.0 && freq <= 1300000.0 {
		record["band"] = "23cm"
	} else {
		return fmt.Errorf("Cannot acertain band for '%.1f'.", freq)
	}
	return nil
}

func (m *Main) getRowByCall(call string) (map[string]interface{}, error) {

	rows, err := m.SelectStmt.Query(strings.TrimSpace(call))
	if rows != nil {
		defer rows.Close()
	}
	if err != nil && err != sql.ErrNoRows {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	ret, errx := shared.GetAllRows(rows)
	if errx != nil {
		return nil, errx
	}
	if len(ret) > 0 {
		return ret[0], nil
	}

	aliasRows, err := m.AliasStmt.Query(strings.TrimSpace(call))
	if aliasRows != nil {
		defer aliasRows.Close()
	}
	if err != nil && err != sql.ErrNoRows {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	ret2, err := shared.GetAllRows(aliasRows)
	for _, v := range ret2 {
		aliases := v["aliases"]
		if aliases != nil {
			if strings.Contains(aliases.(string), call) {
				//realCall := v["call"].(string)
				//log.Printf("Found %s in local DB as an alias for %s", realCall, call)
				return v, nil
			}
		}
	}
	return nil, nil
}

func (m *Main) getAndInsertRowForCall(call string) (map[string]interface{}, error) {

	s := strings.Split(call, "/")
	call = s[0]
	if len(call) < 3 {
		return nil, nil
	}

	row, err := m.getRowByCall(call)
	if err != nil {
		log.Fatal(err)
	}

	if row == nil {
		// lookup call via QRZ API
		qrz, qerr := GetCallFromQRZ(call)
		if qerr != nil {
			if !strings.HasPrefix(qerr.Error(), "Ignoring") {
				log.Println(qerr)
			}
			return nil, nil
			// What next?
		} else {
			if strings.Contains(qrz.Aliases, call) && row != nil {
				log.Printf("Call [%s] is an alias for [%s]", call, qrz.Call)
			} else {
				log.Printf("Call [%s] not found, inserting. [%s]", call, qrz.Call)
				// insert into callsign table
				if _, err := m.InsertStmt.Exec(shared.BindParams(qrz)...); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	return row, err
}
