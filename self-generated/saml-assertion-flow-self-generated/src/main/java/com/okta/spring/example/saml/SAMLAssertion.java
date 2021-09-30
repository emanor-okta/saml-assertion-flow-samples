package com.okta.spring.example.saml;

import com.onelogin.saml2.util.Constants;
import com.onelogin.saml2.util.Util;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.w3c.dom.Document;
import org.w3c.dom.Node;

import javax.annotation.PostConstruct;
import java.security.PrivateKey;
import java.security.cert.X509Certificate;
import java.time.Instant;
import java.util.Base64;

@Component
public class SAMLAssertion {

    private X509Certificate certificate;
    private PrivateKey privateKey;

    @Value("${okta.saml2.certificate}")
    private String certificate_;
    @Value("${okta.saml2.private-key}")
    private String privateKey_;
    @Value("${okta.saml2.issuer}")
    private String issuer;
    @Value("${okta.saml2.acs}")
    private String acs;
    @Value("${okta.saml2.audience}")
    private String audience;


    public SAMLAssertion() {

    }

    @PostConstruct
    private void init() {
        try {
            String cert = Utils.loadPemAsString(certificate_);
//            System.out.println(cert);
            certificate = Util.loadCert(cert);
//            System.out.println(certificate);
            String privKey = Utils.loadPemAsString(privateKey_);
//            System.out.println(privKey);
            privateKey = Util.loadPrivateKey(privKey);
//            System.out.println(privateKey);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }


    public String generateSamlAssertion(String email) throws Exception {
        String assertion = new String(samlResponse);
        String assertionId = Utils.generateSamlID();
        assertion = assertion.replaceAll("\\{ISSUER\\}", issuer)
                .replaceAll("\\{ACS_URL\\}", acs)
                .replaceAll("\\{AUDIENCE\\}", audience)
                .replaceAll("\\{USER_EMAIL\\}", email)
                .replaceAll("\\{AUTH_INSTANT\\}", Instant.now().toString())
                .replaceAll("\\{ISSUE_INSTANT\\}", Instant.now().toString())
                .replaceAll("\\{NOT_BEFORE\\}", Instant.now().minusSeconds(5*60L).toString())
                .replaceAll("\\{NOT_ON_OR_AFTER\\}", Instant.now().plusSeconds(5*60L).toString())
                .replaceAll("\\{RESPONSE_ID\\}", Utils.generateSamlID())
                .replaceAll("\\{ASSERTION_ID\\}", assertionId);;

//        System.out.println(assertion);

        Document document = Util.loadXML(assertion);
        Node assertionNode = Utils.getAssertionNode(document, assertionId);
        assertion = Util.addSign(assertionNode, privateKey, certificate, Constants.RSA_SHA256, Constants.SHA1);
//        System.out.println(assertion);
        assertion = Base64.getEncoder().encodeToString(assertion.getBytes());
//        System.out.println(assertion);
        return assertion;
    }


    private final String samlResponse = "<saml2p:Response Destination=\"{ACS_URL}\"" +
            " ID=\"{RESPONSE_ID}\" IssueInstant=\"{ISSUE_INSTANT}\" Version=\"2.0\"" +
            " xmlns:saml2p=\"urn:oasis:names:tc:SAML:2.0:protocol\"><saml2:Issuer Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:entity\"" +
            " xmlns:saml2=\"urn:oasis:names:tc:SAML:2.0:assertion\">{ISSUER}</saml2:Issuer>" +
            "<saml2p:Status xmlns:saml2p=\"urn:oasis:names:tc:SAML:2.0:protocol\"><saml2p:StatusCode Value=\"" +
            "urn:oasis:names:tc:SAML:2.0:status:Success\"/></saml2p:Status><saml2:Assertion ID=\"{ASSERTION_ID}\"" +
            " IssueInstant=\"{ISSUE_INSTANT}\" Version=\"2.0\" xmlns:saml2=\"urn:oasis:names:tc:SAML:2.0:assertion\">" +
            "<saml2:Issuer Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:entity\" xmlns:saml2=\"urn:oasis:names:tc:SAML:2.0:assertion\">" +
            "{ISSUER}</saml2:Issuer><saml2:Subject xmlns:saml2=\"urn:oasis:names:tc:SAML:2.0:assertion\">" +
            "<saml2:NameID Format=\"urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified\">{USER_EMAIL}</saml2:NameID>" +
            "<saml2:SubjectConfirmation Method=\"urn:oasis:names:tc:SAML:2.0:cm:bearer\"><saml2:SubjectConfirmationData" +
            " NotOnOrAfter=\"{NOT_ON_OR_AFTER}\" Recipient=\"{ACS_URL}\"/>" +
            "</saml2:SubjectConfirmation></saml2:Subject><saml2:Conditions NotBefore=\"{NOT_BEFORE}\"" +
            " NotOnOrAfter=\"{NOT_ON_OR_AFTER}\" xmlns:saml2=\"urn:oasis:names:tc:SAML:2.0:assertion\"><saml2:AudienceRestriction>" +
            "<saml2:Audience>{AUDIENCE}</saml2:Audience></saml2:AudienceRestriction>" +
            "</saml2:Conditions><saml2:AuthnStatement AuthnInstant=\"{AUTH_INSTANT}\" SessionIndex=\"id1632266117710.137242044\"" +
            " xmlns:saml2=\"urn:oasis:names:tc:SAML:2.0:assertion\"><saml2:AuthnContext>" +
            "<saml2:AuthnContextClassRef>urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport</saml2:AuthnContextClassRef>" +
            "</saml2:AuthnContext></saml2:AuthnStatement></saml2:Assertion></saml2p:Response>";
}
