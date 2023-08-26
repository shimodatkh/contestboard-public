# Contestboard

# requirement

- GCS Bucket

## Server

- Go

## Client

- Go
- GCP auth
- alp
- pt-query-digest

# How to use

## GCS 設定

鍵ペア作成
https://cloud.google.com/docs/authentication/production?hl=ja#manually

```shell
gcloud auth activate-service-account --key-file ~/service_account.json
gcloud auth list

# 動作確認
gsutil ls -l gs://yourbucket
gsutil cp pt.log gs://yourbucket
gsutil ls -l gs://yourbucket
gsutil rm gs://yourbucket/pt.log
```

## Server

```shell
go run main.go
```

- ブラウザから http://localhost:8080 にアクセス
- 認証は main.go の user auth を参照（デフォルトは user/pass）

## Client

- alp をダウンロードして`bin/alp`に配置
- pt-query-digest をインストール
- `client/.env`を設定する
  - SERVER_HOST: ダッシュボードを起動しているサーバのホスト名
  - ALP_AGGR_COND: alp の集計条件

```shell
go run client/client.go
```

# サービス定義

```shell
sudo touch /etc/systemd/system/contestboard.service

sudo systemctl daemon-reload
sudo systemctl start contestboard.service
sudo systemctl status contestboard.service
```

sudo systemctl restart contestboard.service

```conf
[Unit]
Description = SimpleHttp.service daemon

[Service]
ExecStart=/bin/bash -c 'go run main.go'
WorkingDirectory=/home/isucon/contestboard/
Restart=always
Type=simple
User=isucon

[Install]
WantedBy = multi-user.target
```
