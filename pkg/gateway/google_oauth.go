package gateway

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"time"

	"github.com/IDEA/SERVER/conf"
	"github.com/IDEA/SERVER/pkg/dto"
	"github.com/IDEA/SERVER/pkg/repository"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	readSpreadSheetsScope = "https://www.googleapis.com/auth/spreadsheets.readonly"
	spreadSheetsAllScope  = "https://www.googleapis.com/auth/spreadsheets"
	sendGmailScope        = "https://www.googleapis.com/auth/gmail.send"
)

const (
	oauthTokenKey = "idea.oauth.token"
)

const ()

type GoogleOAuthGateway interface {
	SendEmailByGmail(payload *dto.EmailPayload) error
	GetSpreadSheetValues(sheetID, range_ string) ([][]interface{}, error)
	UpdateSpreadSheetValues(sheetID, range_ string, values [][]interface{}) error
}

type googleOAuthGateway struct {
	client *http.Client
	cache  repository.CacheRepo
}

func NewGoogleOAuthGateway(cache repository.CacheRepo) GoogleOAuthGateway {
	credentials, err := conf.TokenData.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Failed to read credentials file: %v", err)
	}
	oauthScopes := []string{spreadSheetsAllScope, sendGmailScope}
	if err != nil {
		log.Fatalf("Failed to marshal credentials: %v", err)
	}
	oauthConf, err := google.ConfigFromJSON(credentials, oauthScopes...)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := NewClient(oauthConf, cache)
	return &googleOAuthGateway{client: client, cache: cache}
}

func (g *googleOAuthGateway) SendEmailByGmail(payload *dto.EmailPayload) error {
	ctx := context.Background()
	service, err := gmail.NewService(ctx, option.WithHTTPClient(g.client))
	if err != nil {
		return err
	}
	fromName := mime.QEncoding.Encode("utf-8", "学生団体IDEA")
	fromAddress := "idea.unei.drive@gmail.com"
	rawEmail := "To: " + payload.ToAddress + "\r\n" +
		"From: " + fromName + " <" + fromAddress + ">" + "\r\n" +
		"Subject: " + mime.QEncoding.Encode("utf-8", payload.Subject) + "\r\n" + // MIMEエンコーディング
		"Content-Type: text/plain; charset=utf-8" + "\r\n" + // コンテンツタイプを指定
		"\r\n" +
		payload.Content

	message := &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(rawEmail)),
	}
	_, err = service.Users.Messages.Send("me", message).Do()
	if err != nil {
		return err
	}
	return nil
}

func (g *googleOAuthGateway) GetSpreadSheetValues(sheetID, range_ string) ([][]interface{}, error) {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(g.client))
	if err != nil {
		return nil, err
	}
	resp, err := srv.Spreadsheets.Values.Get(sheetID, range_).Do()
	if err != nil {
		return nil, err
	}
	return resp.Values, nil
}

func (g *googleOAuthGateway) UpdateSpreadSheetValues(sheetID, range_ string, values [][]interface{}) error {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(g.client))
	if err != nil {
		return err
	}
	_, err = srv.Spreadsheets.Values.Update(sheetID, range_, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("RAW").Do()
	return err
}

func NewClient(config *oauth2.Config, cache repository.CacheRepo) *http.Client {
	var token *oauth2.Token
	strToken, err := cache.Get(oauthTokenKey)
	if err != nil || strToken == "" {
		token = getTokenFromWeb(config)
		fmt.Println(token)
		jsonToken, err := json.Marshal(token)
		if err != nil {
			log.Fatal(err.Error())
		}
		if cache.Set(oauthTokenKey, string(jsonToken), time.Hour*24*30); err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println("Cache token")
	} else {
		token = &oauth2.Token{}
		if err = json.Unmarshal([]byte(strToken), token); err != nil {
			log.Fatal(err.Error())
		}
	}
	return config.Client(context.Background(), token)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	// 一度使用したら使用不可になるTokenなのでハードコードしても大丈夫
	// (一ヶ月申請がこないとRefreshTokenが機能しなくなるので、ローカルでauthTokenを取得しなおしてTokenを生成する必要がある)
	authCode := "4/0Adeu5BV6xuPOtOTJ-6YIuIWQ8CgSUcurfgRhwGOoWK5tXCbmVYXi1YvP6q7BwBpnuDb-Hg"
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}
