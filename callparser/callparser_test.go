package callparser

import "testing"

func TestAllPropertiesWithValidCall(t *testing.T) {
    assertEqual(t, NewStation("HC2/DH1TW/P").Prefix, "HC")
    assertEqual(t, NewStation("HC2/DH1TW/P").Valid, true)
    assertEqual(t, NewStation("HC2/DH1TW/P").Call, "HC2/DH1TW/P")
    assertEqual(t, NewStation("HC2/DH1TW/P").Homecall, "DH1TW")
    assertEqual(t, NewStation("HC2/DH1TW/P").Country, "Ecuador")
    assertEqual(t, NewStation("HC2/DH1TW/P").Latitude, float32(-1.4))
    assertEqual(t, NewStation("HC2/DH1TW/P").Longitude, float32(78.4))
    assertEqual(t, NewStation("HC2/DH1TW/P").Cqz, 10)
    assertEqual(t, NewStation("HC2/DH1TW/P").Ituz, 12)
    assertEqual(t, NewStation("HC2/DH1TW/P").Continent, "SA")
    assertEqual(t, NewStation("HC2/DH1TW/P").Mm, false)
    assertEqual(t, NewStation("HC2/DH1TW/P").Beacon, false)
    assertEqual(t, NewStation("HC2/DH1TW/P").Am, false)
}

func TestValidCalls(t *testing.T) {  

    assertEqual(t, NewStation("DH1T").Prefix, "DH")
    assertEqual(t, NewStation("DH1TW/P").Prefix, "DH")
    assertEqual(t, NewStation("DH1TW/MM").Prefix, "")
    assertEqual(t, NewStation("FT5WQ/MM").Prefix, "")
    assertEqual(t, NewStation("DH1TW/AM").Prefix, "")
    assertEqual(t, NewStation("DH1TW/VP5").Prefix, "VP5")
    assertEqual(t, NewStation("VP5/DH1TW").Prefix, "VP5")
    assertEqual(t, NewStation("VP5/DH1TW/P").Prefix, "VP5")
    assertEqual(t, NewStation("MM/DH1TW/P").Prefix, "MM")
    assertEqual(t, NewStation("DH1TW/QRP").Prefix, "DH")
    assertEqual(t, NewStation("DH1TW/QRPP").Prefix, "DH")        
    assertEqual(t, NewStation("MM/DH1TW/QRP").Prefix, "MM")
    assertEqual(t, NewStation("MM/DH1TW/QRPP").Prefix, "MM")
    assertEqual(t, NewStation("MM/DH1TW/B").Prefix, "MM")
    assertEqual(t, NewStation("MM/DH1TW/BCN").Prefix, "MM")
    assertEqual(t, NewStation("EA1/DH1TW").Prefix, "EA")
    assertEqual(t, NewStation("EA1/DH1TW/P").Prefix, "EA")
    assertEqual(t, NewStation("DH1TW/EA1").Prefix, "EA")
    assertEqual(t, NewStation("DH1TW/EA").Prefix, "EA")
    assertEqual(t, NewStation("VP2E/AL1O/P").Prefix, "VP2E")
    assertEqual(t, NewStation("VP2E/DL2001IRTA/P").Prefix, "VP2E")
    assertEqual(t, NewStation("DH1TW/EA8/QRP").Prefix, "EA8")
    assertEqual(t, NewStation("W0ERE/B").Prefix, "W")
    assertEqual(t, NewStation("W0ERE/B").Valid, true)
    assertEqual(t, NewStation("ER/KL1A").Prefix, "ER")
    assertEqual(t, NewStation("DL4SDW/HI3").Prefix, "HI")
    assertEqual(t, NewStation("SV9/M1PAH/HH").Prefix, "SV9")
    assertEqual(t, NewStation("8J3XVIII").Prefix, "8J")
    assertEqual(t, NewStation("3DA0TM").Prefix, "3DA")
    assertEqual(t, NewStation("DL4SDW/HI3").Prefix, "HI")
    assertEqual(t, NewStation("9A2HQ").Prefix, "9A")
    assertEqual(t, NewStation("RU27TT").Prefix, "R")
    assertEqual(t, NewStation("UE90K").Prefix, "UE9")
    assertEqual(t, NewStation("DL2000ALMK").Prefix, "DL")
    assertEqual(t, NewStation("HF450NS").Prefix, "HF")    
    assertEqual(t, NewStation("GB558VUL").Prefix, "G")    
    assertEqual(t, NewStation("F/ON5OF").Prefix, "F")
    assertEqual(t, NewStation("OX1A/OZ1ABC").Prefix, "OX")
    assertEqual(t, NewStation("OX1A/OZ").Prefix, "OZ")
    assertEqual(t, NewStation("OZ5V").Prefix, "OZ")
    assertEqual(t, NewStation("OV9DV").Prefix, "OV")
    assertEqual(t, NewStation("CQ59HQ").Prefix, "CQ")
    assertEqual(t, NewStation("RW3DQC/1/P").Prefix, "R")
    assertEqual(t, NewStation("RW3DQC/1/P").Homecall, "RW3DQC")
    assertEqual(t, NewStation("DB0SUE-10").Prefix, "DB")
    assertEqual(t, NewStation("DK0WYC-2").Prefix, "DK")
    assertEqual(t, NewStation("DK0WYC-2").Valid, true)
    assertEqual(t, NewStation("G0KTD/P").Prefix, "G")
    assertEqual(t, NewStation("GW8IZR-#").Prefix, "GW")
}

func TestInvalidCalls(t *testing.T) {
    sbInvalid(t, NewStation("DH"))
    sbInvalid(t, NewStation("DH1"))
    sbInvalid(t, NewStation("DH1TW/012"))
    sbInvalid(t, NewStation("01A/DH1TW"))
    sbInvalid(t, NewStation("01A/DH1TW/P"))
    sbInvalid(t, NewStation("01A/DH1TW/MM"))
    sbInvalid(t, NewStation("QSL"))
    sbInvalid(t, NewStation("QRV"))
    sbInvalid(t, NewStation("T0NTO"))
    sbInvalid(t, NewStation("T0ALL"))
    sbInvalid(t, NewStation("H1GHMUF"))
    sbInvalid(t, NewStation("C1BBI"))
    sbInvalid(t, NewStation("PU1MHZ/QAP"))
    sbInvalid(t, NewStation("DU7/PA0"))
    sbInvalid(t, NewStation("DIPLOMA"))
    sbInvalid(t, NewStation("CQAS"))
    sbInvalid(t, NewStation("IK2SAV/P1"))
    sbInvalid(t, NewStation("IKOFTA"))
    sbInvalid(t, NewStation("SP2/SP3"))
    sbInvalid(t, NewStation("CQ"))
    sbInvalid(t, NewStation("RADAR"))
    sbInvalid(t, NewStation("MUF/INFO"))
    sbInvalid(t, NewStation("RAVIDEO"))
    sbInvalid(t, NewStation("PIRATE"))
    sbInvalid(t, NewStation("XE1/H"))
    sbInvalid(t, NewStation("Z125VZ"))
    assertEqual(t, NewStation("ZD6DYA").Prefix, "")
    sbInvalid(t, NewStation("ZD6DYA"))
    sbInvalid(t, NewStation("F5BUU1"))
    sbInvalid(t, NewStation("0"))
    sbInvalid(t, NewStation("0123456789"))
    sbInvalid(t, NewStation("CD43000"))
    sbInvalid(t, NewStation("GN"))
    assertEqual(t, NewStation("GN").Homecall, "")
    assertEqual(t, NewStation("ARABS").Homecall, "")
    sbInvalid(t, NewStation("2320900"))
    sbInvalid(t, NewStation("ITT9APL"))
    sbInvalid(t, NewStation("MUF"))
}


func TestLighthouse(t *testing.T) {
    sbValid(t, NewStation("DH1TW/LH"))
    assertEqual(t, NewStation("DH1TW/LH").Prefix, "DH")
    sbValid(t, NewStation("UR7GO/P/LH"))
    assertEqual(t, NewStation("UR7GO/P/LH").Prefix, "UR")
}


func TestPortable(t *testing.T) {
    sbValid(t, NewStation("MM/DH1TW/P"))
    assertEqual(t, NewStation("MM/DH1TW/P").Prefix, "MM")
}


func TestMobile(t *testing.T) {
    sbValid(t, NewStation("VK3/DH1TW/M"))
    assertEqual(t, NewStation("VK3/DH1TW/M").Prefix, "VK")
}


func TestNumberAppendix(t *testing.T) {
    assertEqual(t, NewStation("DH1TW/EA3").PrimaryPrefix, "EA")
    assertEqual(t, NewStation("YB9IR/3").PrimaryPrefix, "YB")
    assertEqual(t, NewStation("UA9MAT/1").PrimaryPrefix, "UA9")		
    assertEqual(t, NewStation("W3LPL/5").PrimaryPrefix, "K")
    assertEqual(t, NewStation("UA9KRM/3").PrimaryPrefix, "UA9")
    assertEqual(t, NewStation("UR900CC/4").PrimaryPrefix, "UR")
}


func TestInvalidCallsWithSpecialCharacters(t *testing.T) {
    sbInvalid(t, NewStation("DK()DK"))
    sbInvalid(t, NewStation("DK/DK"))
    sbInvalid(t, NewStation("'!$&/()@"))
    sbInvalid(t, NewStation(""))
}


func TestBeaconFlag(t *testing.T) {
    assertEqual(t, NewStation("DH1TW/BCN").Beacon, true)
    sbValid(t, NewStation("DH1TW/BCN"))
    assertEqual(t, NewStation("DH1TW/B").Beacon, true)
    sbValid(t, NewStation("DH1TW/B"))
    sbValid(t, NewStation("VP2M/DH1TW/BCN"))
    assertEqual(t, NewStation("VP2M/DH1TW/BCN").Beacon, true)
    assertEqual(t, NewStation("VP2M/DH1TW").Beacon, false)
}


func TestAeronauticalMobile(t *testing.T) {
    assertEqual(t, NewStation("DH1TW/AM").Am, true)
    sbValid(t, NewStation("DH1TW/AM"))
    assertEqual(t, NewStation("VP2M/DH1TW/AM").Am, true)
    sbInvalid(t, NewStation("VP2M/DH1TW/AM"))
    sbValid(t, NewStation("VP2M/DH1TW"))
    assertEqual(t, NewStation("VP2M/DH1TW").Am, false)
}


func TestMaritimeMobile(t *testing.T) {
    assertEqual(t, NewStation("DH1TW/MM").Mm, true)
    sbValid(t, NewStation("DH1TW/MM"))
    assertEqual(t, NewStation("DH1TW/MM").Prefix, "")
    assertEqual(t, NewStation("VP2M/DH1TW/MM").Mm, true)
    sbInvalid(t, NewStation("VP2M/DH1TW/MM"))
    sbValid(t, NewStation("VP2M/DH1TW"))
    assertEqual(t, NewStation("VP2M/DH1TW").Mm, false)
    sbValid(t, NewStation("R7GA/MM"))
    assertEqual(t, NewStation("R7GA/MM").Prefix, "")
    assertEqual(t, NewStation("R7GA/MM").Mm, true)
}


func assertEqual(t *testing.T, s1 interface{}, s2 interface{}) {
    if s1 != s2 {
        t.Errorf("Got '%v' but wanted '%v'\n", s1, s2)
    }
}

func sbInvalid(t *testing.T, station *Station) {
    if station.Valid == true {
        t.Errorf("Call '%s' should be invalid but Valid = true", station.Call)
    }
}


func sbValid(t *testing.T, station *Station) {
    if station.Valid == false {
        t.Errorf("Call '%s' should be valid but Valid = false", station.Call)
    }
}


