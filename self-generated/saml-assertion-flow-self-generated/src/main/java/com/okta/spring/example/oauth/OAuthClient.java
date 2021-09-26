package com.okta.spring.example.oauth;


import okhttp3.*;
import org.json.JSONObject;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.net.URLEncoder;
import java.util.Base64;

@Component
public class OAuthClient {

    @Value("${okta.oidc.client-id}")
    private String clientId;
    @Value("${okta.oidc.client-secret}")
    private String clientSecret;
    @Value("${okta.oidc.issuer}")
    private String issuer;
    @Value("${okta.oidc.scopes}")
    private String scopes;

    private static final String GRANT_TYPE = "urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Asaml2-bearer";
    private static final String FORM_URLENCODED_CONTENT = "application/x-www-form-urlencoded";
    private static final String JSON_CONTENT = "application/json";

    public OAuthClient() {

    }


    public JSONObject sendTokenRequest(String assertion) {
        JSONObject json = null;

        try {
            String content = "grant_type=" + GRANT_TYPE + "&scope=" + scopes.replace(',', ' ') +
                    "&assertion=" + URLEncoder.encode(assertion, "UTF-8");

            byte[] authEncBytes = Base64.getEncoder().encode((clientId+":"+clientSecret).getBytes());
            String authStringEnc = "Basic " + new String(authEncBytes);

            json = sendFormURLEncodedPostRequest(issuer + "/v1/token", authStringEnc, content);
        } catch (IOException e) {
            System.err.println(">> Error sending Token Request: " + e.getLocalizedMessage());
            e.printStackTrace();
            throw e;
        } finally {
            return json;
        }
    }


    public JSONObject sendRestRequest(String url, String authHeader, JSONObject payload) throws IOException {
        return sendPostRequest(url, authHeader, JSON_CONTENT, payload.toString(2));
    }

    private JSONObject sendFormURLEncodedPostRequest(String url, String authHeader,  String content) throws IOException {
        return sendPostRequest(url, authHeader, FORM_URLENCODED_CONTENT, content);
    }

    private JSONObject sendPostRequest(String url, String authHeader, String contentType, String content) throws IOException {
        OkHttpClient client = new OkHttpClient().newBuilder()
                .build();
        MediaType mediaType;

        if (contentType == FORM_URLENCODED_CONTENT) {
            mediaType = MediaType.parse(FORM_URLENCODED_CONTENT);
        } else {
            mediaType = MediaType.parse(JSON_CONTENT);
        }

        RequestBody body = RequestBody.create(mediaType, content);
        Request request = new Request.Builder()
                .url(url)
                .method("POST", body)
                .addHeader("Accept", JSON_CONTENT)
                .addHeader("Authorization", authHeader)
                .addHeader("Content-Type", contentType)
                .build();
        Response response = client.newCall(request).execute();
        JSONObject json = new JSONObject(response.body().string());
        response.close();
        return json;
    }
}
