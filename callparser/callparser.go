package callparser

import (
    //"fmt"
    "strings"
    "strconv"
    "log"
    "os"
    "bufio"
    "regexp"
)

type Station struct {
    Valid                  bool
    Call                   string
    Prefix                 string
    PrimaryPrefix          string
    Homecall               string
    Country                string
    Latitude               float32
    Longitude              float32
    Cqz                    int
    Ituz                   int
    Continent              string
    Offset                 float32
    Mm                     bool
    Am                     bool
    Beacon                 bool
    CallArea               string
}

type CountryInfo struct {
    Country                string
    Cqz                    int
    Ituz                   int
    Continent              string
    Latitude               float32
    Longitude              float32
    Offset                 float32
    PrimaryPrefix          string
	CountryNum			   int
}

type PrefixAlias struct {
    Prefix                 string
    Cqz                    int
    Ituz                   int
    Parent                 *CountryInfo
}

var countries, countriesByNo, prefixes = loadCtyMap("./callparser")

var contintents = map[string]int{
  "NA": 1,
  "SA": 2,
  "EU": 3,
  "AF": 4,
  "AS": 5,
  "OC": 6,
  "AN": 7,
}

var reEndPrefix *regexp.Regexp = regexp.MustCompile("[([]")
var reGetCQZ *regexp.Regexp = regexp.MustCompile("(?:[(])([0-9]+)(?:[)])")
var reGetITUZ *regexp.Regexp = regexp.MustCompile("(?:[[])([0-9]+)(?:[]])")
var reHas3Char = *regexp.MustCompile("(?i)[/A-Z0-9\\-]{3,15}")  // Make sure the call has at least 3 characters
var reLeadingNumber = *regexp.MustCompile("(?i)^[0-9]{1}[A-Z]{1,2}?([0-9]{1})[A-Z]+$")
var reLeadingAlpha = *regexp.MustCompile("(?i)^[A-Z]{1,2}?([0-9]{1,4})[A-Z]+$")
var reRemoveDashSuffix = *regexp.MustCompile("(?i)[-]{1}[0-9#-]{1,4}$")

func NewStation(input string) *Station {
    s := new(Station)
    s.Valid = false
    s.Call = strings.ToUpper(strings.TrimSpace(input))
    s.parseCall(s.Call)
    if !s.Valid {
        log.Printf("Busted Homecall: '%s' of %s could not be decoded", s.Homecall, s.Call)
    } else {
        if !s.Mm && !s.Am {
            if s.Prefix == ""  {
                log.Printf("Busted Prefix: '%s' of %s could not be decoded", s.Prefix, s.Call)
            } else if ctyInfo, ok := prefixes[s.Prefix]; !ok {
                s.Valid = false
                log.Printf("Warning Busted: No country info found for '%s'", s.Call)
            } else {
                s.Country = ctyInfo.Parent.Country
                s.Latitude = ctyInfo.Parent.Latitude
                s.Longitude = ctyInfo.Parent.Longitude
                s.PrimaryPrefix = ctyInfo.Parent.PrimaryPrefix
                s.Cqz = ctyInfo.Cqz
                s.Ituz = ctyInfo.Ituz
                s.Continent = ctyInfo.Parent.Continent
                s.Offset = ctyInfo.Parent.Offset
            }
        }
    }
    return s
}


func (st *Station) parseCall(call string) {
    if a := reHas3Char.FindStringSubmatch(call); a != nil {
        call = reRemoveDashSuffix.ReplaceAllString(call, "")
        switch segments := strings.Split(call, "/"); len(segments) {
            case 1:    // simple call no appendices
                st.checkCall(call, call)
            case 2:   // One appendix
                hasDesig := st.checkDesig(segments[1])
                if _, err := strconv.Atoi(segments[1]); !hasDesig && len(segments[1]) == 1 && err == nil {  // a call area?
                    st.CallArea = segments[1]
                }
                if hasDesig || st.CallArea != "" {  // First segment must be a call
                    st.checkCall(segments[0], segments[0])
                    if st.Mm || st.Am {
                        st.Prefix = ""
                    }
                } else { 
                    if okc1, okp1 := st.checkCall(segments[0], segments[1]); okc1 && okp1 {
                        if okc2, okp2 := st.checkCall(segments[1], segments[0]); okc2 && okp2 {
                            if prefix, ok := iteratePrefix(segments[1]); ok {
                               // Handle situation where homecall is also a prefix (i.e. VP2E)
                               if st.Homecall == prefix {
                                  st.Prefix = segments[1]
                                  st.Homecall = segments[0]
                               }
                            }
                        } else {
                            st.Homecall = segments[0]
                            if prefix, ok := iteratePrefix(segments[1]); ok {
                                st.Prefix = prefix
                            }
                        }
                    } else {
                        st.checkCall(segments[1], segments[0])
                    }
                }
            case 3:   // Two appendices last one must be designator
                st.checkDesig(segments[2])
                var callArea string
                if segments[2] == "P" {
                    if _, err := strconv.Atoi(segments[1]); err == nil {  // a call area?
                        callArea = segments[1]
                    }
                }
                // Portable lighthouse combos
                if (segments[2] == "P" && segments[1] == "LH") || (segments[1] == "P" && segments[2] == "LH") {
                    st.checkCall(segments[0], segments[0])
                    break
                }
                if !st.Mm && !st.Am {
                    if okc, okp := st.checkCall(segments[1], segments[0]); okc && okp {
                        break;
                    } 
                    if okc, okp := st.checkCall(segments[0], segments[1]); okc && okp {
                        break;
                    } else if okc && !okp && callArea != "" && segments[2] == "P" {   // Special case N7ZG/1/P
                        st.CallArea = callArea
                        st.Valid = true
                        st.Homecall = segments[0]
                    }
                } else {  // This is not a valid scenario
                    st.Valid = false
                }
            default:
                log.Printf("Nothing passed in should not be here '%s'", call)
        }
    } else {
        log.Printf("A valid call sign must be at least 3 characters '%s'", call)
    }
}


func (st *Station) checkCall(call string, prefix string) (validCall bool, validPrefix bool) {
    var callArea string
    validCall = false
    if a := reLeadingAlpha.FindStringSubmatch(call); a != nil {
        if len(a[1]) <= 2 {
            callArea = a[1][len(a[1]) - 1:]
        }
        validCall = true
    } else if a := reLeadingNumber.FindStringSubmatch(call); a != nil {
         callArea = a[1][len(a[1]) - 1:]
         validCall = true
    } else {
         validCall = false
    }
    prefix, valid := iteratePrefix(prefix)
    validPrefix = valid
    if validCall {
        st.Homecall = call
    }
    if validPrefix {
       st.Prefix = prefix
    }
    if validCall && validPrefix {
       st.Valid = true
       if st.CallArea == "" {
           st.CallArea = callArea
       }
    }
    return
}  


func (st *Station) checkDesig(s string) bool {
    hasDesig := false
    if s == "MM" {
        hasDesig = true
        st.Mm = true
    } else if s == "AM" {
        hasDesig = true
        st.Am = true
    } else if s == "BCN" || s == "B"  {
        hasDesig = true
        st.Beacon = true
    } else if s == "LH" {
        hasDesig = true
    } else if s == "M" {
        hasDesig = true
    } else if s == "P" {
        hasDesig = true
    } else if s == "QRP" {
        hasDesig = true
    } else if s == "QRPP" {
        hasDesig = true
    } else {
        hasDesig = false
    }
    return hasDesig
}

func LookupCountry(prefix string) (*CountryInfo, bool) {

	if c, ok := countries[prefix]; ok {
		return &c, ok
	}
	return nil, false
} 

func LookupCountryByNo(id int) (*CountryInfo, bool) {

	if c, ok := countriesByNo[id]; ok {
		return &c, ok
	}
	return nil, false
} 

func loadCtyMap(path string) (countries map[string]CountryInfo, countriesByNo map[int]CountryInfo, aliases map[string]PrefixAlias) {

    file, err := os.Open(path + "/cty.dat")
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
    defer file.Close()

    countries = make(map[string]CountryInfo, 0)
    countriesByNo = make(map[int]CountryInfo, 0)
    aliases = make(map[string]PrefixAlias, 0)

    //scan := bufio.NewScanner(result.Body)
    scan := bufio.NewScanner(file)
    
    var fields []string
    var prefixes [] string
    var c *CountryInfo

	countryNum := 0
    for scan.Scan() {
        line := scan.Text()
        fields = strings.Split(line, ":")

        var isFieldLine bool = false
        var isPrefixLine bool = false
        var isLastPrefixLine bool = false

        if len(fields) == 9 {
            isFieldLine = true
        } else {
            isPrefixLine = true
        }
           
        if last := len(line) - 1; isPrefixLine && last >= 0 && line[last] == ';' {
           isLastPrefixLine = true
           line = line[:last]
       
        }

        if isPrefixLine {
            prefixes = strings.Split(line, ",")
        }

        if isFieldLine {
            c = new(CountryInfo)
            c.Country = fields[0]
            if cqz, err := strconv.ParseInt(strings.TrimSpace(fields[1]), 10, 32); err == nil {
                c.Cqz = int(cqz)
            }
            if ituz, err := strconv.ParseInt(strings.TrimSpace(fields[2]), 10, 32); err == nil {
                c.Ituz = int(ituz)
            }
            c.Continent = strings.TrimSpace(fields[3])
            if latitude, err := strconv.ParseFloat(strings.TrimSpace(fields[4]), 32); err == nil {
                c.Latitude = float32(latitude)
            }
            if longitude, err := strconv.ParseFloat(strings.TrimSpace(fields[5]), 32); err == nil {
                c.Longitude = float32(longitude)
            }
            if offset, err := strconv.ParseFloat(strings.TrimSpace(fields[6]), 32); err == nil {
                c.Offset = float32(offset)
            }
            c.PrimaryPrefix = strings.TrimSpace(fields[7])
        } else {
            for _, v := range prefixes {
                a := new(PrefixAlias)
                a.Parent = c
                a.Cqz = c.Cqz
                a.Ituz = c.Ituz
                s := strings.TrimSpace(v)
                i := reEndPrefix.FindStringIndex(s)
                prefix := s
                if len(i) > 0 {
                    prefix = s[:i[0]]
                    if z := reGetCQZ.FindStringSubmatch(s); z != nil {
                        if cqz, err := strconv.ParseInt(z[1], 10, 32); err == nil {
                            a.Cqz = int(cqz)
                        }
                    }
                    if y := reGetITUZ.FindStringSubmatch(s); y != nil {
                        if ituz, err := strconv.ParseInt(y[1], 10, 32); err == nil {
                            a.Ituz = int(ituz)
                        }
                    }
                }
                a.Prefix = prefix
                aliases[prefix] = *a
            }
        }

        if isLastPrefixLine {
			countryNum++
			c.CountryNum = countryNum
            countries[c.PrimaryPrefix] = *c
            countriesByNo[c.CountryNum] = *c
        }
    }
    return
}


// Truncate call until it corresponds to a Prefix in the database
func iteratePrefix(call string) (prefix string, ok bool) {

    ok = true
    prefix = call
    for len(prefix) > 0 {
        if _, found := prefixes[prefix]; found {
            return
        }
        prefix = strings.Replace(prefix, " ", "", -1)
        prefix = prefix[:len(prefix) - 1]
    }
    ok = false
    return
}


func Use(vals ...interface{}) {
    for _, val := range vals {
        _ = val
    }
}


/*
func main() {
    s := NewStation(os.Args[1])
    //fmt.Printf("%v\n", prefixes)
    fmt.Printf("Homecall = %s\n", s.Homecall)
    fmt.Printf("%v\n", s)
    //fmt.Println(iteratePrefix(os.Args[1]))
}
*/


