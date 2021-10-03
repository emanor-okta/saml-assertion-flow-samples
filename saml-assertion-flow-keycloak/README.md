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

### Setup Keycloak IdP and Okta SAML IdP
1. Run `docker run -p 8082:8080 -e KEYCLOAK_USER=admin -e KEYCLOAK_PASSWORD=admin quay.io/keycloak/keycloak:15.0.2`
2. Navigate to `localhost:8082` > **Administration Console** > use **admin**/**admin**
3. In the top left corner where the current Realm is set to *Master* click *Add realm*, name it `okta` and **create**
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
12. Under the **SAML Protocol Settings** enter `http://localhost:8082/auth/realms/okta` for **IdP Issuer URI**. Enter `http://localhost:8082/auth/realms/okta/protocol/saml` for **IdP Single Sign-On URL**
13. For **IdP Signature Certificate** browse to the *.pem* file saved earlier and upload it
14. Click **Show Advanced Settings** and de-select **Sign SAML Authentication Requests**, then click **Add Identity Provider**
15. For the newly create IdP select the drop down arrow on the left side and click **Download metadata**
16. From the Keycloak admin console navigate to **Clients** and click **create**
17. For **import** navigate to the meta-data file saved from Okta and click **save**
  * make note of the **Assertion Consumer Service URL** and **Audience URI** values 
18. De-select **Sign Documents**, **Encrypt Assertions**, and **Client Signature Required**. Click **save** at the bottom
19. Click the **Client Scopes** tab, select **role_list**, then **Remove selected**
20. Select the **Mappers** tab
  * for **firstName**, edit the mapping and enter `firstName` for **User Attribute**, then **save**
  * for **lastName**, edit the mapping and enter `lastName` for **User Attribute**
  * for **email**, edit the mapping and enter `email` for **User Attribute**
21. Navigate back to the main Keycloak admin page, under **Manage** click **Users**, then **Add user**
22. Add a new user by providing **username**, **email**, **firstName**, and **lastName**. Enable **Email Verified** and click **save**
23. Click **Credentials**, enter **Password**, **Password Confirmation**, de-select **Temporary**, click **Set password**
24. In the Okta console create a user in Okta with the same username and **activate** that user.

### Setup the OIDC App to use for the SAML Assertion Token Exchange
Setup an OIDC application in Okta for the SAML Assertion Flow. Instructions can be followed from [here](https://github.com/emanor-okta/saml-assertion-flow-samples/tree/main/self-generated/saml-assertion-flow-self-generated#setup-the-oidc-app-to-use-for-the-saml-assertion-token-exchange). *Step 11* can be ingnored, but the **client_id** and the **client_secret** should be noted. Assign the user created above to the application.

### Setup an Authorization Server for the SAML Assertion Flow
1. The **default** authorization server can be used to test. Navigate to **Security** > **API** > **Authorization Servers**
2. Edit the **default** authorization server and select the **Access Policies** tab.
3. Edit the **Default Policy**, **Default Policy Rule** and enable `SAML 2.0 Assertion` if not already enabled. If there are other policies or rules setup that take priority for this application, then edit the appropriate combo.

### To Run
From the **saml-assertion-flow-keycloak** directory
```
go run main.go
```  

1. Navigate to http://localhost:8080
2. Click **Config**
3. For the **SAML IdP** settings enter,
  * **SAML Request URL**: `http://localhost:8082/auth/realms/okta/protocol/saml`
  * **SAML ACS URL**: value from **step #17** above (ie. `https://DOMAIN.okta.com/sso/saml2/0oa1i54321j7xBCETBd7`)
  * **Entity ID <SAML ISSUER>**: value from **step #17** above (ie. `https://www.okta.com/saml2/service-provider/spggdswuebfgccgopqtr`)
4. Click **Save SAML Settings**
5. For the **Okta OIDC App** settings enter,
  * **Client ID** and **Client Secret** use the values from the [OIDC application](https://github.com/emanor-okta/saml-assertion-flow-samples/tree/main/saml-assertion-flow-keycloak#setup-the-oidc-app-to-use-for-the-saml-assertion-token-exchange) configured.
  * **Token URL**: `https://DOMAIN.okta.com/oauth2/default/v1/token` (unless a different authorization server was setup [above](https://github.com/emanor-okta/saml-assertion-flow-samples/tree/main/saml-assertion-flow-keycloak#setup-an-authorization-server-for-the-saml-assertion-flow))
6. Click **Save OIDC Settings**
7. Click **Get Tokens** to start the flow.   
8. Use the credentials for the Keycloak user created earlier. 
 
### Troubleshooting
Verify the test account is able to login to Okta using Keycloack as the IdP, not Okta. This can be accomplished by setting up a [routing rule](https://help.okta.com/en/prod/Content/Topics/Security/configure-routing-rules.htm) in the Okta Org to redirect to Keycloak, or specifying the Keycloak IdP as part of an [/authorize](https://developer.okta.com/docs/guides/add-an-external-idp/apple/create-authz-url/) request.
    
     
