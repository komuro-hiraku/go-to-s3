# Go-To-S3

指定されたS3の Bucket に 1kb のダミーデータをひたすら突っ込むだけ

## つかいかた

```bash
$ go run main.go -b <YOUR_BUCKET_NAME> -k <KEY_NAME> -a <YOUR_ASSUME_ROLE_ARN> -t <TIMEOUT> -r <RECORD_COUNT>
```

- `YOUR_BUCKET_NAME`: Put対象のBucket名
- `KEY_NAME`: Bucket 内にダミーデータを格納するためのディレクトリ名(`at` と入力すると `bucketName/at/` 配下にダミーデータが大量に投入される)
- `YOUR_ASSUME_ROLE_ARN`: Assume Roleによるアクセスを前提にしてるので `default` プロファイルをソースにSwitchできる Assume Role のARNを記述
- `TIMEOUT`: タイムアウト（単位は分）
- `RECORD_COUNT`: 1プロセスにつきどれくらいレコードを突っ込むか。100並列で実行されるので、 `1` を入力すると合計100件データが投入されます

