package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/api/drive/v3"
)

const mineTypeFolder = "application/vnd.google-apps.folder"

type driveFolder struct {
	ID   string
	Name string
}

func isFolder(f *drive.File) bool {
	return f.MimeType == mineTypeFolder
}

func getFoldersFilesRecursive(service *drive.Service, space string, folder driveFolder, folders, files []*drive.File) ([]*drive.File, []*drive.File) {
	resp, err := service.Files.List().Q(fmt.Sprintf("'%s' in parents", folder.ID)).Do()

	if err != nil {
		fmt.Println(fmt.Errorf("there is an error while fetching from folder %s, error %s", folder.Name, err))
		return folders, files
	}

	if len(resp.Files) == 0 {
		return folders, files
	}

	fmt.Printf("%s- %s\n", space, folder.Name)

	for _, f := range resp.Files {
		if isFolder(f) {
			folders = append(folders, f)
			folders, files = getFoldersFilesRecursive(service, fmt.Sprintf("  %s", space), driveFolder{
				ID:   f.Id,
				Name: f.Name,
			}, folders, files)
		} else {
			files = append(files, f)
			fmt.Printf("  %s* %s\n", space, f.Name)
		}
	}
	return folders, files
}

func GetDriveFolderFilesRecursive(service *drive.Service, folder driveFolder) (folders, files []*drive.File) {
	folders = []*drive.File{}
	files = []*drive.File{}

	resp, err := service.Files.Get(folder.ID).Do()
	if err != nil {
		log.Fatalf("couldn't find root folder %s", folder.Name)
	}

	folders = append(folders, resp)

	folders, files = getFoldersFilesRecursive(service, "", folder, folders, files)

	return folders, files
}

func main() {
	var (
		folderID      string
		folderName    string
		newOwnerEmail string
	)

	service, err := getDriveService()
	if err != nil {
		panic("unable to get drive service")
	}

	fmt.Printf("What folder ID you want to select?\n(i recommend you select most outer folder you need to change owner, i will recursive all inside files and folders)\n")
	if _, err := fmt.Scan(&folderID); err != nil {
		log.Fatalf("Unable to read folder ID %v", err)
	}
	fmt.Printf("What it's name?\n")
	if _, err := fmt.Scan(&folderName); err != nil {
		log.Fatalf("Unable to read folder name %v", err)
	}
	fmt.Printf("Whas's the new owner email?\n")
	if _, err := fmt.Scan(&newOwnerEmail); err != nil {
		log.Fatalf("Unable to read new owner email %v", err)
	}

	targetFolder := driveFolder{
		ID:   folderID,
		Name: folderName,
	}
	targetPermission := &drive.Permission{
		EmailAddress: newOwnerEmail,
		Type:         "user",
		Role:         "owner",
	}

	fmt.Println("Getting list of folders/files need to change owner", targetPermission.EmailAddress, "...")
	folders, files := GetDriveFolderFilesRecursive(service, targetFolder)

	// Need to change owner of folders first or files will be shown in `My Drive`
	// and must not change too fast or else Drive will encounter a bug where sub-folders are both shown in `My Drive` and Parent folder
	fmt.Println("- Folders:")
	for _, f := range folders {
		_, err := service.Permissions.Create(f.Id, targetPermission).TransferOwnership(true).Do()
		if err != nil {
			fmt.Println(fmt.Errorf("unable to change permission of folder `%s`, error: %s", f.Name, err))
			return
		}

		fmt.Printf("owner updated for `%s`\n", f.Name)
		time.Sleep(10 * time.Second)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(files))
	fmt.Println("- Files:")
	for _, f := range files {
		go func(f *drive.File) {
			defer wg.Done()

			_, err := service.Permissions.Create(f.Id, targetPermission).TransferOwnership(true).Do()
			if err != nil {
				fmt.Println(fmt.Errorf("unable to change permission of file `%s`, error: %s", f.Name, err))
				return
			}

			fmt.Printf("owner updated for `%s`\n", f.Name)
		}(f)
	}
	wg.Wait()
}
