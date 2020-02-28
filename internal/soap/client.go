package soap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

func populateRequest() *Request {
	req := Request{}
	req.FirstName = "Tony"
	req.MiddleName = ""
	req.LastName = "Blaire"
	req.Dob = "1946-08-08"
	req.AddressLine1 = "866 Atlas Dr"
	req.AddressLine2 = "Apt 999"
	req.City = "London"
	req.State = "England"
	req.ZipCode = "SW15 5PU"
	req.MobilePhone = "9876543210"
	req.Username = "tony1"
	req.Password = "password1"
	return &req
}

// fake function
func Hello() {
	fmt.Println("hello here")
}

func callSOAPClientSteps() {

	req := populateRequest()

	httpReq, err := generateSOAPRequest(req)
	if err != nil {
		fmt.Println("Some problem occurred in request generation")
	}

	response, err := soapCall(httpReq)
	if err != nil {
		fmt.Println("Problem occurred in making a SOAP call")
	}

	fmt.Println(response.SoapBody.Resp.Status)
}

func generateSOAPRequest(req *Request) (*http.Request, error) {
	// Using the var getTemplate to construct request
	template, err := template.New("InputRequest").Parse(getTemplate)
	if err != nil {
		fmt.Printf("Error while marshling object. %s ", err.Error())
		return nil, err
	}

	doc := &bytes.Buffer{}
	// Replacing the doc from template with actual req values
	err = template.Execute(doc, req)
	if err != nil {
		fmt.Printf("template.Execute error. %s ", err.Error())
		return nil, err
	}

	buffer := &bytes.Buffer{}
	encoder := xml.NewEncoder(buffer)
	err = encoder.Encode(doc.String())
	if err != nil {
		fmt.Printf("encoder.Encode error. %s ", err.Error())
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, "https://www.soapurl.com/retreiveIdentity", bytes.NewBuffer([]byte(doc.String())))
	if err != nil {
		fmt.Printf("Error making a request. %s ", err.Error())
		return nil, err
	}

	return r, nil
}
func soapCall(req *http.Request) (*Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := &Response{}
	err = xml.Unmarshal(body, &r)

	if err != nil {
		return nil, err
	}

	if r.SoapBody.Resp.Status != "200" {
		return nil, err
	}

	return r, nil
}

type Response struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	SoapBody struct {
		XMLName xml.Name `xml:"Body"`
		Resp    struct {
			XMLName   xml.Name `xml:"Response"`
			RequestID string   `xml:"RequestID"`
			Response  struct {
				XMLName       xml.Name `xml:"Body"`
				Status        string   `xml:"Status"`
				Salary        string   `xml:"Salary"`
				Designation   string   `xml:"Designation"`
				Manager       string   `xml:"Manager"`
				Company       string   `xml:"Company"`
				EmployedSince string   `xml:"EmployedSince"`
			}
			Status string `xml:"Status"`
		}
		FaultDetails struct {
			XMLName     xml.Name `xml:"Fault"`
			Faultcode   string   `xml:"faultcode"`
			Faultstring string   `xml:"faultstring"`
		}
	}
}

type Request struct {
	//Values are set in below fields as per the request
	FirstName    string
	LastName     string
	MiddleName   string
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	ZipCode      string
	MobilePhone  string
	SSN          string
	Dob          string
	Username     string
	Password     string
}

var getTemplate = `<soapenv:Envelope
 xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
 xmlns:api="http://soapdummies.com/api">
 <soapenv:Header/>
 <soapenv:Body>
  <api:Command
   xmlns="http://soapdummies.com/api">
   <api:Credentials>
    <api:Username>{{.Username}}</api:Username>
    <api:Password>{{.Password}}</api:Password>
   </api:Credentials>
   <api:Body>
    <SOAPDummy schemaVersion="3.0"
     xmlns="http://soapdummies.com/products/request">
     <Identity>      
      <Title/>
      <FirstName>{{.FirstName}}</FirstName>
      <MiddleName>{{.MiddleName}}</MiddleName>
      <LastName>{{.LastName}}</LastName>
      <Suffix/>
      <DOB>{{.Dob}}</DOB>
      <Address>
       <Line1>{{.AddressLine1}}</Line1>
       <Line2>{{.AddressLine2}}</Line2>
      </Address>
      <City>{{.City}}</City>
      <State>{{.State}}</State>
      <Zip>{{.ZipCode}}</Zip>
      <MobilePhone>{{.MobilePhone}}</MobilePhone>
     </Identity>
    </SOAPDummy>
   </api:Body>
  </api:Command>
 </soapenv:Body>
</soapenv:Envelope>`
