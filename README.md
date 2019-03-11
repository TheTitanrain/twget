##### Program for download audio from twitch with parameters

**Configuration:**

**config.yaml:**  
baseUrl: "https://api.twitch.tv/helix/"  
clientID: ""  # twitch developer id  
keywords: "Сериалоги; Кинологи" #  filter video title. Split with ";"  
period: "week" #  period to retrive videos  
first: "15" # fist # videos to retrive  
argsFilename: "/tools/youtube-dl.exe -f bestaudio --get-filename -o %(title)s.%(ext)s "  
argsEncode: "/tools/youtube-dl.exe -f bestaudio --external-downloader aria2c --keep-video -o %(title)s.%(ext)s "  
argsMp3gain : "/tools/mp3gain.exe /r"  
userName: "stopgameru"  # twitch channel username  