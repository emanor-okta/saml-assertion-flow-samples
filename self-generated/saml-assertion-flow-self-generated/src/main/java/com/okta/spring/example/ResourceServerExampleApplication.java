package com.okta.spring.example;

import com.okta.spring.example.oauth.OAuthClient;
import com.okta.spring.example.saml.SAMLAssertion;
import org.json.JSONArray;
import org.json.JSONObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.security.config.annotation.method.configuration.EnableGlobalMethodSecurity;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.WebSecurityConfigurerAdapter;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;
import org.springframework.security.oauth2.server.resource.authentication.JwtIssuerAuthenticationManagerResolver;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.*;

@SpringBootApplication
@EnableGlobalMethodSecurity(prePostEnabled = true, securedEnabled = true)
public class ResourceServerExampleApplication {

    public static void main(String[] args) {
        ApplicationContext applicationContext = SpringApplication.run(ResourceServerExampleApplication.class, args);
    }


    @Configuration
    static class OktaOAuth2WebSecurityConfigurerAdapter extends WebSecurityConfigurerAdapter {
        @Autowired
        SAMLAssertion samlAssertion;
        @Autowired
        OAuthClient oAuthClient;

        @Value("${okta.authorization.server1}")
        private String server1;
        @Value("${okta.authorization.server2}")
        private String server2;

        @Override
        protected void configure(HttpSecurity http) throws Exception {
            JwtIssuerAuthenticationManagerResolver authenticationManagerResolver = new JwtIssuerAuthenticationManagerResolver
                    (server1, server2);

            http.authorizeRequests(authorize -> authorize
                .anyRequest().authenticated()
            )
            .oauth2ResourceServer(oauth2 -> oauth2
                .authenticationManagerResolver(authenticationManagerResolver)
            );

            // process CORS annotations
            http.cors();
        }
    }


    @RestController
    @CrossOrigin(origins = "http://localhost:8080")
    public class MessageOfTheDayController {
        @Autowired
        SAMLAssertion samlAssertion;
        @Autowired
        OAuthClient oAuthClient;

        @GetMapping("/api/messages")
        @PreAuthorize("hasAuthority('SCOPE_email')")
        public Map<String, Object> messages(JwtAuthenticationToken authentication) {
            Map<String, Object> extendedResult = getExtendedMessages(authentication.getName());

            Map<String, Object> result = new HashMap<>();
            List<Message> messages = new ArrayList<>();
            messages.add(new Message("Hello, world!"));
            messages.add(new Message("I am a robot."));

            JSONArray extendedMessages = (JSONArray)extendedResult.get("messages");
            if (extendedMessages != null) {
                for (Object o : extendedMessages.toList()) {
                    messages.add(new Message((String) ((Map) o).get("text")));
                }
            }

            result.put("messages", messages);
            return result;
        }


        @PostMapping("/api/extendedMessages")
        @PreAuthorize("hasAuthority('SCOPE_saml_flow')")
        public Map<String, Object> extendedMessages(JwtAuthenticationToken authentication) {
            Map<String, Object> result = new HashMap<>();
            result.put("messages", Arrays.asList(
                    new Message("SAML Assertion"),
                    new Message("Flow Messages")
            ));

            return result;
        }


        private Map<String, Object> getExtendedMessages(String subject) {
            Map<String, Object> result = new HashMap<>();

            try {
                String assertion = samlAssertion.generateSamlAssertion(subject);
                JSONObject json = oAuthClient.sendTokenRequest(assertion);
                String accessToken = (String) json.remove("access_token");
                System.out.println(json);
                json = oAuthClient.sendRestRequest("http://localhost:8000/api/extendedMessages",
                        "Bearer " + accessToken, json);
                result.put("messages", json.get("messages"));
            } catch (Exception e) {
                e.printStackTrace();
                return result;
            }

            return result;
        }


        @GetMapping("/api/samlAssertionFlow")
        @PreAuthorize("hasAuthority('SCOPE_email')")
        public String runSamlAssertionFlow(JwtAuthenticationToken authentication) {

            try {
                String assertion = samlAssertion.generateSamlAssertion(authentication.getName());
                JSONObject json = oAuthClient.sendTokenRequest(assertion);
                System.out.println(json);
                String accessToken = (String) json.remove("access_token");
                json = oAuthClient.sendRestRequest("https://httpbin.org/post", "Bearer " + accessToken, json);
                json.remove("data");
                json.remove("args");
                json.remove("form");
                json.remove("files");
                json.remove("origin");
                System.out.println(json);
                return json.toString(2);
            }catch (Exception e) {
                e.printStackTrace();
                return null;
            }
        }
    }

    class Message {
        public Date date = new Date();
        public String text;

        Message(String text) {
            this.text = text;
        }

    }
}
