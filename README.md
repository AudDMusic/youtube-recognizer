# youtube-recognizer
Recognizes music from YouTube video using [AudD music recognition API](https://audd.io/). Also generates .csv with all songs from the video or prints first found song.

Usage:
```bash
./youtube-recognizer
  -api_token string
        AudD API token (default "test")
  -url string
        Link to the YouTube video (default "https://www.youtube.com/watch?v=ANEOD16twxo")
  -csv string
        Path to the .csv which will be created (default "audd.csv")
```

[![Usage demo](https://img.youtube.com/vi/j1ChhoqdlsM/0.jpg)](https://www.youtube.com/watch?v=j1ChhoqdlsM)
