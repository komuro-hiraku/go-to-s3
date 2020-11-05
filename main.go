package main

import (
	"flag"
	"fmt"

	"github.com/komuro-hiraku/go-to-s3/s3"
)

// https://github.com/aws/aws-sdk-go#Quick-Examples
func main() {
	var bucket, key string // s3 params
	var roleArn string
	var timeoutMin int64
	var recordCount int

	// flagをパース
	flag.StringVar(&bucket, "b", "", "Bucket name.")
	flag.StringVar(&key, "k", "", "Object key name.")
	flag.StringVar(&roleArn, "a", "", "role arn")
	flag.Int64Var(&timeoutMin, "t", 30, "Timeout Minutes. Default 30 min")
	flag.IntVar(&recordCount, "r", 1000, "Record Count / 1process. Default 1000 count")
	flag.Parse()

	// Check
	if len(bucket) == 0 {
		panic(fmt.Errorf("You MUST input bucket name. %s", bucket))
	}
	if len(key) == 0 {
		panic(fmt.Errorf("You MUST input key. %s", key))
	}
	if len(roleArn) == 0 {
		panic(fmt.Errorf("You MUST input Assume Role ARN. %s", roleArn))
	}

	// Create Upload
	p := &s3.ParUpload{
		Bucket: bucket,
		Key: key,
		Arn: roleArn,
	}

	// 100並列でUploadを実行
	p.Upload(100, recordCount, timeoutMin)
}
