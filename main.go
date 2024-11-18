package main

import (
    "bytes"
    "encoding/json"
    //"fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/joho/godotenv"
)

// Константы URL
const (
    EcowittAPIURL    = "https://api.ecowitt.net/api/v3/device/real_time"
    ThingSpeakAPIURL = "https://api.thingspeak.com/update"
)

// Структуры данных
type EcowittResponse struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
    Time string `json:"time"`
    Data struct {
        Outdoor struct {
            Temperature struct {
                Time  string `json:"time"`
                Unit  string `json:"unit"`
                Value string `json:"value"`
            } `json:"temperature"`
            Humidity struct {
                Time  string `json:"time"`
                Unit  string `json:"unit"`
                Value string `json:"value"`
            } `json:"humidity"`
        } `json:"outdoor"`
        Indoor struct {
            Temperature struct {
                Time  string `json:"time"`
                Unit  string `json:"unit"`
                Value string `json:"value"`
            } `json:"temperature"`
            Humidity struct {
                Time  string `json:"time"`
                Unit  string `json:"unit"`
                Value string `json:"value"`
            } `json:"humidity"`
        } `json:"indoor"`
        Wind struct {
            WindSpeed struct {
                Time  string `json:"time"`
                Unit  string `json:"unit"`
                Value string `json:"value"`
            } `json:"wind_speed"`
        } `json:"wind"`
        Rainfall struct {
            Daily struct {
                Time  string `json:"time"`
                Unit  string `json:"unit"`
                Value string `json:"value"`
            } `json:"daily"`
        } `json:"rainfall"`
        Pressure struct {
            Relative struct {
                Time  string `json:"time"`
                Unit  string `json:"unit"`
                Value string `json:"value"`
            } `json:"relative"`
        } `json:"pressure"`
        SolarAndUVI struct {
            Solar struct {
                Time  string `json:"time"`
                Unit  string `json:"unit"`
                Value string `json:"value"`
            } `json:"solar"`
        } `json:"solar_and_uvi"`
    } `json:"data"`
}

type ThingSpeakPayload struct {
    APIKey string `json:"api_key"`
    Field1 string `json:"field1"`
    Field2 string `json:"field2"`
    Field3 string `json:"field3"`
    Field4 string `json:"field4"`
    Field5 string `json:"field5"`
    Field6 string `json:"field6"`
    Field7 string `json:"field7"`
    Field8 string `json:"field8"`
}

func getEcowittData(applicationKey, apiKey, mac string) (*EcowittResponse, error) {
    client := &http.Client{
        Timeout: 10 * time.Second, // Тайм-аут для запроса
    }

    req, err := http.NewRequest("GET", EcowittAPIURL, nil)
    if err != nil {
        return nil, err
    }

    // Добавление параметров запроса
    q := req.URL.Query()
    q.Add("application_key", applicationKey)
    q.Add("api_key", apiKey)
    q.Add("mac", mac)
    q.Add("call_back", "all")
    q.Add("temp_unitid", "1")
    q.Add("pressure_unitid", "5")
    q.Add("wind_speed_unitid", "7")
    q.Add("rainfall_unitid", "12")
    q.Add("solar_irradiance_unitid", "16")
    req.URL.RawQuery = q.Encode()

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var ecowittResp EcowittResponse
    err = json.Unmarshal(body, &ecowittResp)
    if err != nil {
        return nil, err
    }

    return &ecowittResp, nil
}

func sendToThingSpeak(writeAPIKey string, fields ThingSpeakPayload) (string, error) {
    payload := map[string]string{
        "api_key": writeAPIKey,
        "field1":  fields.Field1,
        "field2":  fields.Field2,
        "field3":  fields.Field3,
        "field4":  fields.Field4,
        "field5":  fields.Field5,
        "field6":  fields.Field6,
        "field7":  fields.Field7,
        "field8":  fields.Field8,
    }

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return "", err
    }

    client := &http.Client{
        Timeout: 10 * time.Second, // Тайм-аут для запроса
    }

    resp, err := client.Post(ThingSpeakAPIURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(body), nil
}

func parseData(data *EcowittResponse) ThingSpeakPayload {
    return ThingSpeakPayload{
        Field1: data.Data.Outdoor.Temperature.Value,
        Field2: data.Data.Outdoor.Humidity.Value,
        Field3: data.Data.Indoor.Temperature.Value,
        Field4: data.Data.Indoor.Humidity.Value,
        Field5: data.Data.Wind.WindSpeed.Value,
        Field6: data.Data.Rainfall.Daily.Value,
        Field7: data.Data.Pressure.Relative.Value,
        Field8: data.Data.SolarAndUVI.Solar.Value,
    }
}

func loadEnv() {
    err := godotenv.Load("ecowitt_to_thingspeak.env")
    if err != nil {
        log.Fatalf("Ошибка при загрузке файла .env: %v", err)
    }
}

func main() {
    // Загрузка переменных окружения
    loadEnv()

    // Получение переменных из окружения
    EcowittApplicationKey := os.Getenv("ECOWITT_APPLICATION_KEY")
    EcowittAPIKey := os.Getenv("ECOWITT_API_KEY")
    EcowittMAC := os.Getenv("ECOWITT_MAC")
    ThingSpeakWriteAPIKey := os.Getenv("THINGSPEAK_WRITE_API_KEY")

    // Проверка наличия переменных
    if EcowittApplicationKey == "" || EcowittAPIKey == "" || EcowittMAC == "" || ThingSpeakWriteAPIKey == "" {
        log.Fatal("Необходимые переменные окружения не установлены.")
    }

    ticker := time.NewTicker(1 * time.Minute) // Тикер раз в минуту
    defer ticker.Stop()

    // Немедленный запуск перед первым тикером
    for {
        start := time.Now()
        ecowittData, err := getEcowittData(EcowittApplicationKey, EcowittAPIKey, EcowittMAC)
        if err != nil {
            log.Printf("Ошибка при получении данных с Ecowitt API: %v", err)
        } else {
            if ecowittData.Code != 0 {
                log.Printf("Ошибка от Ecowitt API: %s", ecowittData.Msg)
            } else {
                fields := parseData(ecowittData)
                result, err := sendToThingSpeak(ThingSpeakWriteAPIKey, fields)
                if err != nil {
                    log.Printf("Ошибка при отправке данных на ThingSpeak: %v", err)
                } else {
                    log.Printf("Данные успешно отправлены на ThingSpeak: %s", result)
                }
            }
        }

        // Ждём следующего тика
        elapsed := time.Since(start)
        if elapsed < time.Minute {
            time.Sleep(time.Minute - elapsed)
        }
    }
}
