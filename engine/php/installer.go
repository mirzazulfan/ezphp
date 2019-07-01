package php

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/mholt/archiver"
	"github.com/sirupsen/logrus"
)

type Installer struct {
	DownloadUrl string
	Filename    string
	InstallDir  string
}

func (i *Installer) InstallPHP(ioCom IOCom) {

	var err error

	absPath, _ := filepath.Abs(filepath.Dir(i.InstallDir))
	localInstallDir := absPath + string(os.PathSeparator) + i.InstallDir

	ioCom.Outmsg <- NewStdInstall("\nInstalling PHP v7.0.0 in your local directory: " + localInstallDir + "\n")
	ioCom.Outmsg <- NewStdInstall("Downloading PHP from: " + i.DownloadUrl + "/" + i.Filename + "\n")

	_, err = i.download(ioCom)
	if err != nil {
		logrus.Error("Error downloading file " + err.Error())
		NewStderr(err.Error())
		return
	}

	err = i.unzip()
	if err != nil {
		logrus.Error("Error unzipping file " + err.Error())
		NewStderr(err.Error())
		return
	}

}

func (i Installer) download(ioCom IOCom) (*grab.Response, error) {
	logrus.Info("Downloading PHP from " + i.DownloadUrl + "/" + i.Filename)
	client := grab.NewClient()
	req, _ := grab.NewRequest(i.InstallDir+string(os.PathSeparator)+i.Filename, i.DownloadUrl+"/"+i.Filename)
	resp := client.Do(req)
	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			ioCom.Outmsg <- NewStdInstall(fmt.Sprintf("\rDownload in progress: %.2f%%", 100*resp.Progress()))

		case <-resp.Done:
			break Loop
		}

	}

	ioCom.Outmsg <- NewStdInstall(fmt.Sprint("\rDownload in progress: 100%  "))

	if err := resp.Err(); err != nil {
		return nil, resp.Err()
	}

	return resp, nil
}

func (i Installer) unzip() error {
	logrus.Info("Unziping local PHP installation: " + i.InstallDir + string(os.PathSeparator) + i.Filename)
	err := archiver.Unarchive(i.InstallDir+string(os.PathSeparator)+i.Filename, i.InstallDir)
	if err != nil {
		return err
	}

	return nil
}
