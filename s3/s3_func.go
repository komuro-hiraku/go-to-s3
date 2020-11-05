package s3

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

var lettersRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type ParUpload struct {
	Bucket string
	Key string
	Arn string
}

func (p *ParUpload) Upload(paralells, uploadRecordCount int, timeoutMin int64) {

	// Assume roleしないとダメ
	creds := credentials.NewSharedCredentials("", "default")	// defaultを取得
	config := aws.Config{Region: aws.String("ap-northeast-1"), Credentials: creds}
	defaultSession := session.New(&config)

	// Assume role
	var sCreds *credentials.Credentials

	if len(p.Arn) != 0 {
		sCreds = stscreds.NewCredentials(defaultSession, p.Arn)
	} else {
		log.Info("set default Credentials")
		sCreds = aws.NewConfig().WithRegion("ap-northeast-1").Credentials
	}
	sConfig := aws.Config{
		Region: aws.String("ap-northeast-1"), 
		Credentials: sCreds, 
		S3ForcePathStyle: aws.Bool(true),
	}
	assumeSess := session.New(&sConfig)
	uploader := s3manager.NewUploader(assumeSess)
	uploader.Concurrency = 1000
			
	message := make(chan string)
	for i := 0; i<paralells; i++ {
		go func() {
			_, e := execUpload(message, p.Bucket, p.Key, uploadRecordCount, uploader)
			if e != nil {
				fmt.Printf("error: %v\n", e)
				return
			}
		}()
	}
	defer close(message)

	timeout := time.After(time.Duration(timeoutMin) * time.Minute)
	resultCount := 0
	for {
		// TimeoutもChannelらしい
		select {
		case _, ok := <-message:
			if !ok {
				fmt.Printf("%T\n", ok)
				return	
			}
			// 1000件で1 point追加
			if resultCount % 1000 == 0 {
				fmt.Printf(".")
			}
			resultCount++
		case <-timeout:
			fmt.Printf("Finished: Add %d\n", resultCount)
			return
		}
	}
}

func execUpload(ch chan string, bucket, key string, count int, uploader *s3manager.Uploader) (int, error) {
	resultCount := 0
	for i := 0; i<count; i++ {
		b := randomStringRunes(1024)
		u, err := uuid.NewRandom()
		if err != nil {
			return resultCount, fmt.Errorf("failed to uuid, %v", err)
		}
		k := u.String()
	
		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key: aws.String(fmt.Sprintf("%s/%s", key, k)),
			Body: strings.NewReader(b),
		})
	
		if err != nil {
			return resultCount, fmt.Errorf("failed to upload file, %v", err)
		}
		msg := fmt.Sprintf("file uploaded to, %s\n", aws.StringValue(&result.Location))
		ch <- msg
		time.Sleep(time.Millisecond * 10)	// 10ms Sleep
		resultCount++
	}
	return resultCount, nil
}

// 指定した数のランダム文字列を生成
func randomStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano()) // randを初期化
	b := make([]rune, n)             // 最大長のsliceを定義
	for i := range b {
		b[i] = lettersRunes[rand.Intn(len(lettersRunes))]
	}
	return string(b)
}
