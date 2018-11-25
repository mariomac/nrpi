#include <Arduino.h>

#include <ESP8266WiFi.h>

#include <ESP8266HTTPClient.h>

#define ANALOG_PIN 0
#define INTERVAL_uS 60e6
#define CONNECT_LED D1
#define SEND_LED D2

#define JSON_MAX_LENGTH 128

// to be defined by the user
#include "config.h"

int connect();
void disconnect();
double temperature();
void measureAndSend();

void setup()
{
  pinMode(D0, INPUT);           // For wakeup after deep sleep. D0 must be connected with RST pins
  pinMode(CONNECT_LED, OUTPUT); // led indicating the system is workin
  pinMode(SEND_LED, OUTPUT);    // led indicating the system is workin

#if DEBUG
  Serial.begin(9600);
  Serial.setTimeout(2000);
#endif // DEBUG

  connect();

  measureAndSend();
  disconnect();

  ESP.deepSleep(INTERVAL_uS, WAKE_RF_DEFAULT);
}

int connect()
{
  digitalWrite(CONNECT_LED, HIGH);
  WiFi.persistent(false);
  WiFi.mode(WIFI_STA);
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

  int timeout = 0;
  // wait for WiFi connection
  while (WiFi.status() != WL_CONNECTED)
  {
#if DEBUG
    Serial.print(".");
#endif // DEBUG
    delay(1000);
    digitalWrite(CONNECT_LED, LOW);
    if (timeout++ >= 20)
    {
#if DEBUG
      Serial.println("--> ERROR: Wifi commection timeout");
#endif // DEBUG
      return -1;
    }
  }
  return 0;
}

void disconnect()
{
  WiFi.disconnect();
  WiFi.mode(WIFI_OFF);
  digitalWrite(CONNECT_LED, LOW);
  digitalWrite(SEND_LED, LOW);
}

// sends data directly to new relic
void directSubmission(char *json)
{
// TO DO
}
// sends data to the IoT agent
void agentSubmission(char *json)
{
  HTTPClient http;
  // configure traged server and url
  http.begin(NRIOT_AGENT_URL); //HTTP
  int httpCode = http.POST(json);
#if DEBUG
  if (httpCode > 0)
  {
    Serial.printf("Http code: %d\n", httpCode);
    // file found at server
    if (httpCode == HTTP_CODE_OK)
    {
      String payload = http.getString();
      Serial.println(payload);
    }
  }
  else
  {
    Serial.printf("[HTTP] GET... failed, error: %s\n", http.errorToString(httpCode).c_str());
  }
#endif // DEBUG

  http.end();
}

void measureAndSend()
{
  digitalWrite(SEND_LED, HIGH);

  char json[JSON_MAX_LENGTH];

  double temp = temperature();
  sprintf(json, "{\"eventType\":\"TemperatureSample\",\"sensorId\":\"%s\",\"temperature\":%.2lf,\"unit\":\"Celsius\"}", SENSOR_ID, temp);

#if DEBUG
  Serial.printf("Sending: %s\n", json);
#endif // DEBUG

  agentSubmission(json);
}

// Voltage calculation for the TMP36 sensor, connected to the Analog pin
double temperature()
{
  //getting the voltage reading from the temperature sensor
  int reading = analogRead(ANALOG_PIN);

  // converting that reading to voltage, for 3.3v arduino use 3.3
  double voltage = reading * 3.3 / 1024.0;

  // now print out the temperature
  return (voltage - 0.5) * 100; //converting from 10 mv per degree wit 500 mV offset
}
void loop()
{
}