package main

import (
  "fmt"
  "os/exec"
  "bytes"
  "log"
  "strings"
  "io/ioutil"
  "encoding/json"
  "os"
  "sort"
)

func getRpms() []string {
  cmd := exec.Command("rpm", "-qlav")
  var out bytes.Buffer
  cmd.Stdout = &out
  err := cmd.Run()
  if err != nil {
    log.Fatal(err)
  }

  f := func(c rune) bool {
    return c == '\n'
  }
  packages := strings.FieldsFunc(out.String(), f)
  return packages
}

func showRpms() {
  for _, pkg := range getRpms() {
    fmt.Printf("PACKAGE: %s\n", pkg)
  }
}

func showDir(dir string) {
  files, _ := ioutil.ReadDir("/home")
  for _, f := range files {
    fmt.Println(f.Name())
  }
}

func parseRpmLine(line string) (file_type string, file_name string, link_target string) {
  file_type = line[0:1]

  index := strings.Index(line, "/")
  if index < 0 {
    panic(line)
  }
  file := line[index:]
  fields := strings.Split(file, " -> ")
  if len(fields) == 2 {
    file_name = fields[0]
    link_target = fields[1]
  } else {
    file_name = file
  }
  return
}

func addImplicitlyManagedDirs(dirs map[string]bool, files map[string]string) {
  for file, target := range files {
    for i:= 1; i < len(file); i++ {
      if file[i] == '/' {
        topdir := file[:i]
        if _, ok := dirs[topdir]; !ok {
          dirs[topdir] = false
        }
      }
    }

    if target != "" {
      if _, ok := dirs[target]; ok {
        dirs[file] = false
      }
    }
  }
  return
}

func getManagedFiles() (map[string]string, map[string]bool) {
  files := make(map[string]string)
  dirs := make(map[string]bool)

  for _, pkg := range getRpms() {
    if pkg != "(contains no files)" {
      file_type, file_name, link_target := parseRpmLine(pkg)
      switch file_type {
        case "-":
          files[file_name] = ""
        case "d":
          dirs[file_name] = true
        case "l":
          files[file_name] = link_target
      }
    }
  }

  addImplicitlyManagedDirs(dirs, files)

  return files, dirs
}

func printJsonString(jsonMap map[string]string, file_name string) {
  b, err := json.Marshal(jsonMap)

  if err != nil {
    log.Fatal("JSON conversion failes")
  }

  f, err := os.Create(file_name)
  if err != nil {
    panic(err)
  }

  defer f.Close()

  var out bytes.Buffer
  json.Indent(&out, b, "", "  ")
  out.WriteTo(f)
}

func printJsonArray(jsonMap []map[string]string, file_name string) {
  b, err := json.Marshal(jsonMap)

  if err != nil {
    log.Fatal("JSON conversion failes")
  }

  f, err := os.Create(file_name)
  if err != nil {
    panic(err)
  }

  defer f.Close()

  var out bytes.Buffer
  json.Indent(&out, b, "", "  ")
  out.WriteTo(f)
}

func printJsonBool(jsonMap map[string]bool, file_name string) {
  b, err := json.Marshal(jsonMap)

  if err != nil {
    log.Fatal("JSON conversion failes")
  }

  f, err := os.Create(file_name)
  if err != nil {
    panic(err)
  }

  defer f.Close()

  var out bytes.Buffer
  json.Indent(&out, b, "", "  ")
  out.WriteTo(f)
}

func findUnmanagedFiles(dir string, rpm_files map[string]string, rpm_dirs map[string]bool, unmanaged_files map[string]string) {
  ignore_list := map[string]bool{
    "/etc/group": true,
    "/etc/passwd": true,
    "/etc/shadow": true,
    "/etc/init.d/boot.d": true,
    "/etc/init.d/rc0.d": true,
    "/etc/init.d/rc1.d": true,
    "/etc/init.d/rc2.d": true,
    "/etc/init.d/rc3.d": true,
    "/etc/init.d/rc4.d": true,
    "/etc/init.d/rc5.d": true,
    "/etc/init.d/rc6.d": true,
    "/etc/init.d/rcS.d": true,
    "/dev": true,
    "/proc": true,
    "/tmp": true,
    "/run": true,
    "/sys": true,
    "/var/tmp": true,
    "/lost+found": true,
    "/var/run": true,
    "/var/lib/rpm": true,
  }

  files, _ := ioutil.ReadDir(dir)
  for _, f := range files {
    file_name := dir + f.Name()
    if _, ok := ignore_list[file_name]; !ok {
      if f.IsDir() {
        if _, ok := rpm_dirs[file_name]; ok {
          findUnmanagedFiles(file_name + "/", rpm_files, rpm_dirs, unmanaged_files)
        } else {
          unmanaged_files[file_name + "/"] = "dir"
        }
      } else {
        if _, ok := rpm_files[file_name]; !ok {
          if f.Mode() & os.ModeSymlink == os.ModeSymlink {
            unmanaged_files[file_name] = "link"
          } else {
            unmanaged_files[file_name] = "file"
          }
        }
      }
    }
  }
}

func main() {
//  showRpms()
//  showDir("/home")

  rpm_files, rpm_dirs := getManagedFiles()

  var unmanaged_files map[string]string
  unmanaged_files = make(map[string]string)

  printJsonString(rpm_files, "RPM_FILES")
  printJsonBool(rpm_dirs, "RPM_DIRS")

  findUnmanagedFiles("/", rpm_files, rpm_dirs, unmanaged_files)

  files := make([]string, len(unmanaged_files))
  i := 0
  for k, _ := range unmanaged_files {
    files[i] = k
    i++
  }
  sort.Strings(files)

  unmanaged_files_json := make([]map[string]string, len(unmanaged_files))
  for j := range files {
    entry := make(map[string]string)
    entry["name"] = files[j]
    entry["type"] = unmanaged_files[files[j]]
    unmanaged_files_json[j] = entry
  }

  printJsonArray(unmanaged_files_json, "UNMANAGED_FILES")
}