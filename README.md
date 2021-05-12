# youtube

## 概要
Youtube Data API の開発練習用。サーバ側はGo言語、フロントエンドはVue.jsを使用しています。bootstrap 5と axios はあまり使ったことがなかったのでこちらも練習として使用してあります。

## 環境
- go 1.14.3
- bootstrap 5.0.0
- vue.js 2
- axios 0.21.1

## 初期対応
Youtube Data API の使用のはAPIキーが必要なので以下のconfig.ini.defaultファイルからconfig.iniを作成してAPIキーを指定してください。

```
cd youtube/config
mv config.ini.default config.ini
```

https://github.com/nero-15/youtube/blob/main/youtube/config/config.ini.default#L2

