# Self Generate SAML Assertion Example

This sample application demonstrates a resource server receiving a request containing a Bearer Access Token. The value of the subject in this token is used to generate a SAML Assertion and exchange it for a new access token to be used to call another API with a different authorization server and added scope.
The first Access Token is received from one of the Okta Front-end sample SPA applications referenced below.

## Prerequisites

Before running this sample, you will need the following:

* An Okta Developer Account, you can sign up for one at https://developer.okta.com/signup/.
* Configure one of the below Okta Front-end sample SPA applications.
  * [Okta Angular Sample Apps](https://github.com/okta/samples-js-angular)
  * [Okta Vue Sample Apps](https://github.com/okta/samples-js-vue)
  * [Okta React Sample Apps](https://github.com/okta/samples-js-react)


## Setup This Example
#### Setup the SAML IdP in Okta to be used by the native OIDC application for the SAML Assertion Flow..
1. Generate a a self signed certificate to be used for testing with the following openSSL command, `openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes -subj 'CN={YOUR_ORG_NAME}'`. Keep these .pem files in an accessible location.
2. Login to your Org and navigate to **security** > **identity providers**
3. Click **Add Identity Provider** and select **Add SAML 2.0 IdP**.
4. Enter a Name for your IdP.
5. For **IdP Username** select **idpuser.subjectNameId**.
6. Under SAML Protocol Settings for **Idp Issuer URI** enter `http://localhost:8080/self/generated`. For the **Idp Single Sign-On URL** enter `http://localhost:8080/self/generated/saml`.
7. For IdP Signature Certificate browse to where **cert.pem** was saved and upload it.
8. Click **Add Identity Provider**.
9. After the IdP is created, from the Identity Providers screen click the drop down arrow on the left hand side of the IdP just created. Copy down the values of the **Assertion Consumer Service URL** and **Audience URI**.
10. Edit `/src/main/resources/application.yml` and modify the SAML section with the values from above.
```yaml
  saml2:
    certificate: {PATH_TO_CERT_PEM}
    private-key: {PATH_TO_KEY_PEM}
    issuer: {SAML_ISSUER_URI}
    acs: {SAML_ACS_URL}
    audience: {SAML_AUDIENCE}
```

#### Setup the OIDC App to use for the SAML Assertion Token Exchange
1. Navigate to **applications** > **applications** > **Create App Integration**
2. Select **OIDC - OpenID Connect** > **Native Application** > **Next**
3. Give it a meaningful name
4. Select Grant Types **Authorization Code**, **Refresh Token**, and **SAML 2.0 Assertion**
5. For **Assignments** select **Skip group assignment for now**
6. Keep the rest of the defaults, **save**
7. From the **General** tab under **Client Credentials** click **edit**.
8. Select **Use Client Authentication** and save.
9. Make note of the **Client ID** and **Client secret** values.
10. Click **Assignments** and add a test user.
11. Edit `/src/main/resources/application.yml` and modify the SAML section with the values from above.

#### Setup an Authorization Server for the SAML Assertion Flow
1. Navigate to **Security** > **API** > **Authorization servers** > **Add Authorization Server**.
2. Copy down the value for **Issuer**
3. Naviagte to **Scopes** > **Add Scope**
4. For **Name**, **Display Name**, and **Description** enter `saml_flow` then **save**.
5. Navigate to **Access Policies** > **Add Policy**
6. Enter `SAML Flow` for **Name** and **Description**
7. Select **The following clients** and add the above OIDC application just created, then **Create Policy**.
8. Select **Add Rule**
9. Provide a name and only allow Grant Types **Authorization Code** and **SAML 2.0 Assertion**, then **Create Rule**
10. Edit `/src/main/resources/application.yml` and modify the OIDC section with the **issuer**, **Client ID**, and **Client secret** values. 

```yaml
  oidc:
    issuer: {OKTA_ISSUER_WITH_SAML_ASSERTION_ENABLED}
    client-id: {OIDC_APP_ID_SAML_ASSERTION_ENABLED}
    client-secret: {OIDC_APP_SECRET}
    scopes: openid,profile,email,offline_access,saml_flow
  authorization:
    server1: https://{OKTA_ORG}/oauth2/{AS_SAML_ASSERTION_ENABLED} #same as issuer under oidc
    server2: https://{OKTA_ORG}/oauth2/default # if the Okta front-end app does not use the default auth server, modify this value to be the same as the front-end app.
```

**backend:**
```bash
./mvnw
```
> **NOTE:** The above command starts the resource server on port 8000. You can browse to `http://localhost:8000` to ensure it has started. If you get the message "401 Unauthorized", it indicates that the resource server is up. You will need to pass an access token to access the resource, which will be done by the front-end below.

**front-end:**

Follow the instructions from one of the Okta Sample front-end applications used.
  * [Okta Angular Sample Apps](https://github.com/okta/samples-js-angular)
  * [Okta Vue Sample Apps](https://github.com/okta/samples-js-vue)
  * [Okta React Sample Apps](https://github.com/okta/samples-js-react)


