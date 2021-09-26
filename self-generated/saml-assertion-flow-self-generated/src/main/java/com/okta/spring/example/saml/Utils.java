package com.okta.spring.example.saml;

import org.w3c.dom.Document;
import org.w3c.dom.Node;
import org.w3c.dom.NodeList;

import javax.xml.xpath.*;
import java.io.*;
import java.util.Random;

public class Utils {
    private static Random random = new Random();
    private static final long MIN = 1000000000000000000L;
    private static final long MAX = 9223372036854775807L;

    public static String loadPemAsString(String location) {
        StringBuilder sb = new StringBuilder();
        BufferedReader br = null;

        try {
            br = new BufferedReader(new FileReader(location));
            char buffer[] = new char[1024];
            for (int read = br.read(buffer); read != -1; read = br.read(buffer)) {
                sb.append(buffer, 0, read);
            }
        } catch (Exception e) {
            System.err.println("Unable to load PEM file: " + location);
        } finally {
            if (br != null) {
                try {
                    br.close();
                } catch (Exception e) {}
            }
        }

        return sb.toString();
    }

    public static String generateSamlID() {
        return "id" + random.longs(MIN,MAX).findFirst().getAsLong();
    }

    public static Node getAssertionNode(Document document, String id) throws XPathExpressionException {
        XPathFactory xPathFactory;
        try {
            xPathFactory = XPathFactory.newInstance("http://java.sun.com/jaxp/xpath/dom",
                                         "com.sun.org.apache.xpath.internal.jaxp.XPathFactoryImpl",
                                                         ClassLoader.getSystemClassLoader());
        } catch (XPathFactoryConfigurationException e) {
            xPathFactory = XPathFactory.newInstance();
        }
        XPath xpath = xPathFactory.newXPath();
        XPathExpression expr = xpath.compile("//*[@ID='" + id + "']");
        NodeList nodeList = (NodeList)expr.evaluate(document, XPathConstants.NODESET);
        return (Node)expr.evaluate(document, XPathConstants.NODE);
    }
}
