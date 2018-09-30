# googledrive-change-owner
Change files, folders owner in batch and recursive from Google Drive folder

### Requirement:
- Go
- Glide

On MacOS, you can install by:

```
$ brew install go
$ brew install glide
```

### Install dependencies:
- Go to the source code directory and run `glide install` to install all Go dependecies.

---

### How to use:
- Go to https://developers.google.com/drive/api/v3/quickstart/go#step_1_turn_on_the to one step enable Drive API for your Google account.
- Then download `credentials.json` in the last screen and put it into the source code directory.
- Simply run `go run main.go client.go` and do step show in the console.

Note: there will be step require you to enter folder ID. Folder ID is the hashed string showed in URL when you enter a folder. E.g https://drive.google.com/drive/u/1/folders/1Py9XWHhlYP0bYWFI2lu2c1R-Emt2Ivjj, then `1Py9XWHhlYP0bYWFI2lu2c1R-Emt2Ivjj` will be your folder ID.
