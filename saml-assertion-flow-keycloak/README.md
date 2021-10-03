# Sample SAML Assertion Flow using Keycloak as a SAML IdP

This example shows how to use the [Okta SAML Assertion Flow](https://developer.okta.com/docs/guides/implement-saml2/overview/) to exchange an assertion for tokens.
The application does the following:
1. Programmatically sends a SAML Request to a Keycload IdP 
2. Handles client authentication with Keycloak if needed
3. Receives the SAML response *intended for Okta*, and strips out the SAML Assertion
4. Makes a SAML Assertion Flow Token Exchange call to Okta with the Assertion
5. Receives Tokens and displays the stack of calls made 

## Prerequisites

Before running this sample, you will need the following:

* An Okta Developer Account, you can sign up for one at https://developer.okta.com/signup/.
* A running [Keycloak](https://www.keycloak.org/) instance.
  * instructions for setting up and running a dockerized instance below
* [Go](https://golang.org/) installed, *1.16+* 


### To Install
```
git clone https://github.com/emanor-okta/saml-assertion-flow-samples.git
cd saml-assertion-flow-keycloak
go mod tidy
```

### Setup Keycloak IdP
1. Run `docker run -p 8082:8080 -e KEYCLOAK_USER=admin -e KEYCLOAK_PASSWORD=admin quay.io/keycloak/keycloak:15.0.2`
2. Navigate to `localhost:8082` > **Administration Console** > use **admin**/**admin**
3. In the top left corner where the current Realm is set to *Master* click *Add realm*, name if `okta` and **create**
4. With **okta** set as the current realm select **Indentity Providers** > **SAML v2.0**
5. For **Single Sign-On Service URL** enter **http://localhost:8082/auth/realms/okta/protocol/saml**
6. Set **NameID Policy Format** to **Unspecified**, then **save**
7. Pull up the meta-data from URL `http://localhost:8082/auth/realms/okta/protocol/saml/descriptor`
8. Copy the Value of the x509 certificate found between the `<ds:X509Certificate>` elements. Create a new unformatted plain text file and create the Begin/End lines with this value between like below. Save it with a *.pem* extension
```
-----BEGIN CERTIFICATE-----
{x509_VALUE}
-----END CERTIFICATE-----
```
9. In your Okta Org navigate to **Security** > **Identitiy Providers** > **Add Identity Provider** > **SAML 2.0**
10. Give the IdP a meaningful name. For **IdP Username** select `idpuser.subjectNameId`
11. For **if no match is found** select `Redirect to Okta sign-in page`
12. Under the **SAML Protocol Settings** enter `http://localhost:8082/auth/realms/okta` for **IdP Issuer URI**. Enter `http://localhost:8082/auth/realms/okta/protocol/saml` for **IdP Single Sign-On URL **
13. For **IdP Signature Certificate** browse to the *.pem* file saved earlier and upload it. 

### To Run
```
go run main.go
```  

1. Navigate to http://localhost:8080   
2. Click 'Get Tokens' to start the flow.   
3. Use credentials `read.only`/`Th1sPassword` to login. 
    
     
