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
	"github.com/IDEA/SERVER/pkg/util"
	"github.com/morikuni/failure"

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
	config := newOAuthConfig()
	client, err := newClient(config, cache)
	if err != nil {
		// 通知処理を入れたいが層が複雑になるので悩み所
		log.Printf("Failed to init google oauth2 client: %v", err)
	}
	return &googleOAuthGateway{client: client, cache: cache}
}

func newOAuthConfig() *oauth2.Config {
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
	return oauthConf
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

func newClient(config *oauth2.Config, cache repository.CacheRepo) (*http.Client, error) {
	var token *oauth2.Token
	strToken, err := cache.Get(oauthTokenKey)
	if err != nil || strToken == "" {
		token ,err= getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		jsonToken, err := json.Marshal(token)
		if err != nil {
			return nil, err
		}
		if cache.Set(oauthTokenKey, string(jsonToken), time.Hour*24*30); err != nil {
			return nil, err
		}
	} else {
		token = &oauth2.Token{}
		if err = json.Unmarshal([]byte(strToken), token); err != nil {
			return nil, err
		}
		token,err=refleshOAuthToken(token,config,cache)
		if err!=nil{
			return nil,failure.Wrap(err)
		}
	}
	return config.Client(context.Background(), token), nil
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token,error ){
	authRedirectURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authRedirectURL)

	// 一度使用したら使用不可になるTokenなのでハードコードしても大丈夫
	// (一ヶ月申請がこないとRefreshTokenが機能しなくなるので、ローカルでauthTokenを取得しなおしてTokenを生成する必要がある)
	authCode := "4/0AfJohXkjN8uqzAbVCMAhxL7jS-PWbil9h-sV2j5Fa8MKFcc0edSN9mbtaj1NpKraSpEvjQ"
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil,failure.Wrap(err)
	}
	return tok,nil
}


func refleshOAuthToken(token *oauth2.Token, conf *oauth2.Config, cache repository.CacheRepo)(*oauth2.Token,error){
	if token.Valid(){
		return token,nil
	}
	reuseToken:=oauth2.ReuseTokenSource(token,conf.TokenSource(context.Background(),token))
	if token.Expiry.Unix() < time.Now().Unix() {
		return nil, fmt.Errorf("refresh token has expired at %v",token.Expiry)
	}

	tkn,err:=reuseToken.Token()
	if err!=nil{
		return nil,failure.Wrap(err)
	}
	strToken,err:=util.Serialize(tkn)
	log.Println(strToken)
	if err:=cache.Set(oauthTokenKey,strToken,24*30*time.Hour);err!=nil{
		return nil,err
	}
	return tkn,nil
}
