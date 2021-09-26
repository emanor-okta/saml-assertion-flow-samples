package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/hokaccha/go-prettyjson"
	"mellium.im/xmlstream"
)

func GetAssertionFromKCResponse(resp string) string {
	samlRespDec, _ := base64.StdEncoding.DecodeString(resp)
	// pull out SAML assertion / KC sets xlmns:saml at response level so add it hack
	assertion := strings.Split(string(samlRespDec), "<saml:Assertion xmlns")
	assertion = strings.Split(assertion[1], "</saml:Assertion>")

	return (base64.URLEncoding.EncodeToString([]byte("<saml:Assertion xmlns:saml" + assertion[0] + "</saml:Assertion>")))
}

func GetAssertionFromOktaResponse(resp string) string {
	samlRespDec, _ := base64.StdEncoding.DecodeString(resp)
	assertion := strings.Split(string(samlRespDec), "<saml2:Assertion")
	assertion = strings.Split(assertion[1], "</saml2:Assertion>")

	return (base64.URLEncoding.EncodeToString([]byte("<saml2:Assertion" + assertion[0] + "</saml2:Assertion>")))
}

func FormatXML(s string) string {
	tokenizer := xmlstream.Fmt(xml.NewDecoder(strings.NewReader(s)), xmlstream.Indent("   "))
	buf := new(bytes.Buffer)
	e := xml.NewEncoder(buf)
	for t, err := tokenizer.Token(); err == nil; t, err = tokenizer.Token() {
		e.EncodeToken(t)
	}
	e.Flush()

	return buf.String()
}

func FormatJSON(s string) string {
	fmt.Printf("\nV\n%v\n\n", s)
	formatter := prettyjson.NewFormatter()
	formatter.DisabledColor = true
	formatter.Indent = 3
	bytes, _ := formatter.Format([]byte(s))
	return string(bytes)
}

func DecodeB64(s string) string {
	decoded, _ := base64.StdEncoding.DecodeString(s)
	return string(decoded)
}

func URLDecodeB64(s string) string {
	decoded, _ := base64.URLEncoding.DecodeString(s)
	return string(decoded)
}

func RawDecodeB64(s string) string {
	decoded, _ := base64.RawStdEncoding.DecodeString(s)
	return string(decoded)
}

func EncodeB64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func RawEncodeB64(s string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(s))
}
