# youtube-recognizer
Recognizes music from a YouTube video using the [Music Recognition API](https://audd.io/). Generates a .csv file with all the songs from the video.

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
