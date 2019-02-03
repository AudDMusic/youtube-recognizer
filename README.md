# youtube-recognizer
Recognize music from YouTube video using [AudD music recognition API](https://audd.io/).
Also generates .csv with all songs from the video or print first found song.

Needs ffmpeg to be [installed](https://github.com/AudDMusic/youtube-recognozer/wiki/Installing-FFmpeg).

Usage:
```bash
./youtube-recognizer
-api_token string
        AudD API token (default "test")
  -csv string
        Path to the .csv which will be created (default "audd.csv")
  -first
        Send requests only until first result
  -s int
        Seconds per audio file (default 9)
  -url string
        Link to the YouTube video (default "https://www.youtube.com/watch?v=ANEOD16twxo")
```
