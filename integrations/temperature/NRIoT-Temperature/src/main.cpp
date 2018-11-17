#include <Arduino.h>

#include <ESP8266WiFi.h>
#include <ESP8266WiFiMulti.h>

#include <ESP8266HTTPClient.h>

#define ANALOG_PIN 0
#define DELAY_MS 5000

#define JSON_MAX_LENGTH 256

// to be defined by the user
#define DEBUG true
#define NRIOT_AGENT_URL "http://192.168.1.2:8080/http"

#define SENSOR_ID "livingRoom"
#define WIFI_SSID "MASMOVIL_he9x"
#define WIFI_PASSWORD "********"

double temperature();

ESP8266WiFiMulti WiFiMulti;

char json[JSON_MAX_LENGTH];

void setup() {
  #if DEBUG
  Serial.begin(9600);
  #endif // DEBUG

  WiFi.mode(WIFI_STA);
  WiFiMulti.addAP(WIFI_SSID, WIFI_PASSWORD);
}

void loop() {

  // wait for WiFi connection
  if ((WiFiMulti.run() == WL_CONNECTED)) {

    HTTPClient http;

    // configure traged server and url
    http.begin(NRIOT_AGENT_URL); //HTTP

    double temp = temperature();
    sprintf(json, "{\"eventType\":\"TemperatureSample\",\"sensorId\":\"%s\",\"temperature\":%lf,\"unit\":\"Celsius\"}", SENSOR_ID, temp);
    #if DEBUG
    Serial.printf("Sending: %s\n", json);
    #endif // DEBUG

    int httpCode = http.POST(json);

    #if DEBUG
    if (httpCode > 0) {
      Serial.printf("Http code: %d\n", httpCode);
      // file found at server
      if (httpCode == HTTP_CODE_OK) {
        String payload = http.getString();
        Serial.println(payload);
      }
    } else {
      Serial.printf("[HTTP] GET... failed, error: %s\n", http.errorToString(httpCode).c_str());
    }
    #endif // DEBUG

    http.end();
  }
  delay(DELAY_MS);
}

// Voltage calculation for the TMP36 sensor, connected to the Analog pin
double temperature() {
  //getting the voltage reading from the temperature sensor
  int reading = analogRead(ANALOG_PIN);  

  // converting that reading to voltage, for 3.3v arduino use 3.3
  double voltage = reading * 3.3 / 1024.0;

  // now print out the temperature
  return (voltage - 0.5) * 100 ;  //converting from 10 mv per degree wit 500 mV offset
}
