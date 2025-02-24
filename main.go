package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "html/template"
    "io/ioutil"
    "net/http"
)

type HuggingFaceRequest struct {
    Inputs string `json:"inputs"`
}

type HuggingFaceResponse []struct {
    GeneratedText string `json:"generated_text"`
}

func main() {
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(template.ParseFiles("templates/index.html"))
        tmpl.Execute(w, nil)
    })

    http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
            return
        }

        err := r.ParseForm()
        if err != nil {
            http.Error(w, "Failed to parse form", http.StatusBadRequest)
            return
        }

        userInput := r.Form.Get("message")
        if userInput == "" {
            http.Error(w, "Message cannot be empty", http.StatusBadRequest)
            return
        }

        apiURL := "https://api-inference.huggingface.co/models/gpt2"
        apiToken := "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"

        response, err := sendRequestToHuggingFace(apiURL, apiToken, userInput)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
            return
        }

        if len(*response) > 0 {
            w.Write([]byte((*response)[0].GeneratedText))
        } else {
            w.Write([]byte("No response generated."))
        }
    })

    fmt.Println("Starting server at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}

func sendRequestToHuggingFace(apiURL, apiToken, userInput string) (*HuggingFaceResponse, error) {
    requestBody := HuggingFaceRequest{
        Inputs: userInput,
    }

    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request body: %v", err)
    }

    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }

    req.Header.Set("Authorization", "Bearer "+apiToken)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
    }

    var huggingFaceResponse HuggingFaceResponse
    err = json.Unmarshal(body, &huggingFaceResponse)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
    }

    return &huggingFaceResponse, nil
}