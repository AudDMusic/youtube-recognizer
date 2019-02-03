# youtube-recognozer
Recognize music from YouTube video using AudD music recognition API.
Also generates .csv with all songs from the video or print first found song.

## WILL NOT WORK UNTIL NEW AUDD API VERSION RELEASE ON FEBRUARY 6 

Needs ffmpeg to be [installed](https://github.com/adaptlearning/adapt_authoring/wiki/Installing-FFmpeg).

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
