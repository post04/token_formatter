package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	in = bufio.NewReader(os.Stdin)
)

func getInput(toPrint string) string {
	fmt.Print(toPrint)
	value, _ := in.ReadString('\n')
	return strings.TrimSpace(value)
}

type tokenStruct struct {
	Token         string
	Password      string
	Email         string
	Phone         string
	Username      string
	Discriminator string
	ID            string
	Flags         string
	Verified      string
}

type meReq struct {
	ID             string      `json:"id"`
	Username       string      `json:"username"`
	Avatar         string      `json:"avatar"`
	Discriminator  string      `json:"discriminator"`
	PublicFlags    int         `json:"public_flags"`
	Flags          int         `json:"flags"`
	PurchasedFlags int         `json:"purchased_flags"`
	Banner         interface{} `json:"banner"`
	BannerColor    interface{} `json:"banner_color"`
	AccentColor    interface{} `json:"accent_color"`
	Bio            string      `json:"bio"`
	Locale         string      `json:"locale"`
	NsfwAllowed    bool        `json:"nsfw_allowed"`
	MfaEnabled     bool        `json:"mfa_enabled"`
	Email          string      `json:"email"`
	Verified       bool        `json:"verified"`
	Phone          string      `json:"phone"`
}

func (t *tokenStruct) CheckToken() {
	req, err := http.NewRequest("GET", "https://discord.com/api/v9/users/@me", nil)
	if err != nil {
		return
	}
	for _, header := range strings.Split(`accept: */*|accept-language: en-US|cache-control: no-cache|pragma: no-cache|referer: https://discord.com/login|sec-fetch-dest: empty|sec-fetch-mode: cors|sec-fetch-site: same-origin|user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36|x-debug-options: bugReporterEnabled|x-super-properties: eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiQ2hyb21lIiwiZGV2aWNlIjoiIiwic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiYnJvd3Nlcl91c2VyX2FnZW50IjoiTW96aWxsYS81LjAgKFdpbmRvd3MgTlQgMTAuMDsgV2luNjQ7IHg2NCkgQXBwbGVXZWJLaXQvNTM3LjM2IChLSFRNTCwgbGlrZSBHZWNrbykgQ2hyb21lLzkyLjAuNDUxNS4xMzEgU2FmYXJpLzUzNy4zNiIsImJyb3dzZXJfdmVyc2lvbiI6IjkyLjAuNDUxNS4xMzEiLCJvc192ZXJzaW9uIjoiMTAiLCJyZWZlcnJlciI6IiIsInJlZmVycmluZ19kb21haW4iOiIiLCJyZWZlcnJlcl9jdXJyZW50IjoiIiwicmVmZXJyaW5nX2RvbWFpbl9jdXJyZW50IjoiIiwicmVsZWFzZV9jaGFubmVsIjoic3RhYmxlIiwiY2xpZW50X2J1aWxkX251bWJlciI6OTYzNTUsImNsaWVudF9ldmVudF9zb3VyY2UiOm51bGx9`, "|") {
		parts := strings.Split(header, ": ")
		req.Header.Set(parts[0], parts[1])
	}
	req.Header.Set("authorization", t.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	f := &meReq{}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &f)
	if err != nil {
		return
	}
	if f.Phone == "" {
		f.Phone = "null"
	}
	t.Phone = f.Phone
	if f.Email == "" {
		f.Email = "null"
	}
	t.Email = f.Email
	t.Username = f.Username
	t.Discriminator = f.Discriminator
	t.ID = f.ID
	t.Flags = fmt.Sprint(f.Flags)
	t.Verified = fmt.Sprintf("%t", f.Verified)
	return
}

func main() {

	splitKey := getInput("Split Key (typically colon): ")
	tokenFile := getInput("Token file: ")
	currFormat := getInput("Current format: ")
	currFormat = strings.ToLower(currFormat)
	endFormat := getInput("Ending format: ")
	endFormat = strings.ToLower(endFormat)
	checkTokens := getInput("Check tokens (y/n): ")
	checkTokens = strings.ToLower(checkTokens)

	currKeys := strings.Split(currFormat, splitKey)
	endKeys := strings.Split(endFormat, splitKey)
	tokens := []string{}
	tokenStr, err := os.ReadFile(tokenFile)
	if err != nil {
		fmt.Println(err)
		time.Sleep(5 * time.Second)
		return
	}
	tokens = strings.Split(string(tokenStr), "\r\n")
	fmt.Println("Loaded", len(tokens), "tokens!")
	structedTokens := []*tokenStruct{}
	endTokens := []string{}
	for _, token := range tokens {
		parts := strings.Split(token, splitKey)
		t := &tokenStruct{}
		for i, part := range parts {
			switch currKeys[i] {
			case "password":
				t.Password = part
			case "email":
				t.Email = part
			case "token":
				t.Token = part
			case "phone":
				t.Token = part
			case "username":
				t.Username = part
			case "discriminator":
				t.Discriminator = part
			case "id":
				t.ID = part
			case "flags":
				t.Flags = part
			case "verified":
				t.Verified = part
			default:
				fmt.Println("Unsupported key", currKeys[i], "found!")
			}
		}
		if checkTokens == "y" || checkTokens == "yes" {
			t.CheckToken()
		}
		structedTokens = append(structedTokens, t)
	}
	for _, token := range structedTokens {
		t := ""
		for _, part := range endKeys {
			switch part {
			case "password":
				t += token.Password + splitKey
			case "email":
				t += token.Email + splitKey
			case "token":
				t += token.Token + splitKey
			case "phone":
				t += token.Phone + splitKey
			case "username":
				t += token.Username + splitKey
			case "discriminator":
				t += token.Discriminator + splitKey
			case "id":
				t += token.ID + splitKey
			case "flags":
				t += token.Flags + splitKey
			case "verified":
				t += token.Verified + splitKey
			default:
				fmt.Println("Unsupported key", part, "found!")
			}
		}
		endTokens = append(endTokens, t[:len(t)-1])
	}
	endFile := getInput("Output file: ")
	os.WriteFile(endFile, []byte(strings.Join(endTokens, "\n")), 0064)
}
