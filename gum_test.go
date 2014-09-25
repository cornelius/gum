package main

import (
  "testing"
  "reflect"
)

func doTestParseRpmLine(t *testing.T, line string, expected_file_type string, expected_file_name string, expected_link_target string) {
  file_type, file_name, link_target := parseRpmLine(line)

  if file_type != expected_file_type {
    t.Errorf("parseRpmLine('%s') file type = '%s', want '%s'", line, file_type, expected_file_type)
  }
  if file_name != expected_file_name {
    t.Errorf("parseRpmLine('%s') file name = '%s', want '%s'", line, file_name, expected_file_name)
  }
  if link_target != expected_link_target {
    t.Errorf("parseRpmLine('%s') file type = '%s', want '%s'", line, link_target, expected_link_target)
  }
}

func TestParseRpmLineFile(t *testing.T) {
  line := "-rw-r--r--    1 root    root                 18234080 Mar 31 11:40 /usr/lib64/libruby2.0-static.a"

  expected_file_type := "-"
  expected_file_name := "/usr/lib64/libruby2.0-static.a"
  expected_link_target := ""

  doTestParseRpmLine(t, line, expected_file_type, expected_file_name, expected_link_target)
}

func TestParseRpmLineDir(t *testing.T) {
  line := "drwxr-xr-x    2 root    root                        0 Mar 31 11:45 /usr/include/ruby-2.0.0/x86_64-linux/ruby"

  expected_file_type := "d"
  expected_file_name := "/usr/include/ruby-2.0.0/x86_64-linux/ruby"
  expected_link_target := ""

  doTestParseRpmLine(t, line, expected_file_type, expected_file_name, expected_link_target)
}

func TestParseRpmLineLink(t *testing.T) {
  line := "lrwxrwxrwx    1 root    root                       19 Mar 31 11:45 /usr/lib64/libruby2.0.so -> libruby2.0.so.2.0.0"

  expected_file_type := "l"
  expected_file_name := "/usr/lib64/libruby2.0.so"
  expected_link_target := "libruby2.0.so.2.0.0"

  doTestParseRpmLine(t, line, expected_file_type, expected_file_name, expected_link_target)
}

func TestParseRpmLineFileSpaces(t *testing.T) {
  line := "-rw-r--r--    1 root    root                    61749 Jun 26 01:56 /usr/share/kde4/templates/kipiplugins_photolayoutseditor/data/templates/a4/h/Flipping Tux Black.ple"

  expected_file_type := "-"
  expected_file_name := "/usr/share/kde4/templates/kipiplugins_photolayoutseditor/data/templates/a4/h/Flipping Tux Black.ple"
  expected_link_target := ""

  doTestParseRpmLine(t, line, expected_file_type, expected_file_name, expected_link_target)
}

func TestAddImplicitlyManagedDirs(t *testing.T) {
  files_original := map[string]string{
    "/abc/def/ghf/somefile": "",
    "/zzz": "/abc/def",
  }
  dirs_original := map[string]bool{
    "/abc/def": true,
  }
  dirs_expected := map[string]bool{
    "/abc": false,
    "/abc/def": true,
    "/abc/def/ghf": false,
    "/zzz": false,
  }

  dirs := dirs_original

  addImplicitlyManagedDirs(dirs, files_original)

  if !reflect.DeepEqual(dirs, dirs_expected) {
    t.Errorf("addImplicitlyManagedDirs('%v') = '%v', want '%v'", dirs_original, dirs, dirs_expected)
  }
}
