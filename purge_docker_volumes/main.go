package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "path"
  "github.com/samalba/dockerclient"
)

const (
  vfsBase = "/var/lib/docker/vfs/dir"
)

var (
  daemonUrl string 
)

type VolInfo struct {
  inUse bool
}

func markVolumesInUse(volInfos map[string]*VolInfo) []string {
  client, err := dockerclient.NewDockerClient(daemonUrl, nil)
  if err != nil {
    log.Fatal(err)
  }
  
  containers, err := client.ListContainers(true, false, "")
  if err != nil {
    log.Fatal(err)
  }

  var inUse []string

  for _,container := range containers {
    info,err := client.InspectContainer(container.Id)
    if err != nil {
      log.Println("WARN","Could not inspect container", container, err)
    } else {
      for _,vfsDir := range info.Volumes {
        inUse = append(inUse, vfsDir)
        volInfo,exists := volInfos[vfsDir]
        if exists {
          volInfo.inUse = true
        }
      }
    }
  }

  return inUse
}

func getAllKnownVolumes() map[string]*VolInfo {
  result := make(map[string]*VolInfo)
  dir, err := os.Open(vfsBase)
  if err != nil {
    log.Fatal(err)
  }

  contents, err := dir.Readdir(-1)
  if err != nil {
    log.Fatal(err)
  }

  for _,d := range contents {
    joined := path.Join(vfsBase,d.Name())
    result[joined] = new(VolInfo)
  }

  return result
}

func purge(volInfos map[string]*VolInfo) (purged int) {
  for vfsDir,info := range volInfos {
    if !info.inUse {
      fmt.Println("DELETING", vfsDir)
      err := os.RemoveAll(vfsDir)
      if err != nil {
        log.Println("WARN", "Removing", vfsDir, err)
      } else {
        purged++
      }
    }
  }
  return
}

func main() {
  flag.StringVar(&daemonUrl, "H", "unix:///var/run/docker.sock", "The Docker daemon's socket")
  flag.Parse()

  vols := getAllKnownVolumes()

  markVolumesInUse(vols)

  purged := purge(vols)
  if purged == 0 {
    fmt.Println("Congrats, nothing to purge")
  }
}
