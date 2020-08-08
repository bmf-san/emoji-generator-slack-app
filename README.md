# emoji-generator-slack-app
This is an emoji generator app for slack.

![download](https://user-images.githubusercontent.com/13291041/89861807-9577c380-dbe1-11ea-8510-3288b3767ef3.png)
![download (1)](https://user-images.githubusercontent.com/13291041/89861804-94df2d00-dbe1-11ea-8864-e162a0dffe06.png)

# Example
`@botname [color] [bgColor] [line1] [line2(optional)]`

![Screen Shot 2020-08-11 at 14 49 46](https://user-images.githubusercontent.com/13291041/89861979-f3a4a680-dbe1-11ea-8c93-7c118c89e813.png)
![Screen Shot 2020-08-11 at 14 49 40](https://user-images.githubusercontent.com/13291041/89861975-f1dae300-dbe1-11ea-8e59-10ef38800cce.png)


# Get Started
You need to be prepared to set environment variables.
```
SLACK_SIGNIN_SECRET=
SLACK_BOT_TOKEN=
```

Also, you need to have these settings.

### Singing Secret
<img src="https://user-images.githubusercontent.com/13291041/89857888-3a8d9e80-dbd8-11ea-88af-70fb2031bef8.png" width="300px">

### Bot Token
<img src="https://user-images.githubusercontent.com/13291041/89857887-39f50800-dbd8-11ea-8edf-c392866cdef6.png" width="300px">

### Required Permission - Bot Token Scopes
<img src="https://user-images.githubusercontent.com/13291041/89857886-395c7180-dbd8-11ea-845d-1072a5d7da0f.png" width="300px">

### Event subscriptions - Request URL
<img src="https://user-images.githubusercontent.com/13291041/89857879-36618100-dbd8-11ea-8e70-ccf18590f7bf.png" width="300px">


If you want to try it out easily in a local environment, try using ngork.
```
go run main.go
ngrok http 9999
```

Please read a caution before using in production environment if you want this app in production environment.

# API
## `/generator`
- API for generatoring image. 
- The following items are prepared as query parameters.
    - color
      - text color.
    - bgColor
      - Background color.
    - line1
      - First line text. 
    - line2
      - Second line text.
      - Optional
- ex. http://localhost:9999/generator?color=red&bgColor=green&line1=foo&line2=bar

## `/slack/events`
- API for slack events subscription.

# Caution
There is a bug where images are posted multiple times.

# References
- [note.com - Goでheadless browserを用いた動的画像生成](https://note.com/timakin/n/n55d483d11b22)
- [qiita.com - Go で Slack Bot を作る (2020年3月版)](https://qiita.com/frozenbonito/items/cf75dadce12ef9a048e9)
- [qiita.com - Go で Interactive な Slack Bot を作る (2020年5月版)](https://qiita.com/frozenbonito/items/1df9bb685e6173160991#%E3%81%BE%E3%81%A8%E3%82%81)
- [dev.to - Slackで送った文字を画像で返すbot作った](https://dev.to/amotarao/slackbot-376)
- [lab.syncer.jp - 複数行のテキストを描く方法](https://lab.syncer.jp/Web/JavaScript/Canvas/8)
- [Stackoverflow - Size to fit font on a canvas](https://stackoverflow.com/questions/20551534/size-to-fit-font-on-a-canvas)

