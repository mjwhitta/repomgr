# RepoMgr

[![Yum](https://img.shields.io/badge/-Buy%20me%20a%20cookie-blue?labelColor=grey&logo=cookiecutter&style=for-the-badge)](https://www.buymeacoffee.com/mjwhitta)

![License](https://img.shields.io/github/license/mjwhitta/repomgr?style=for-the-badge)

## What is this?

Simple git repo manager. At this time it simply ensures your repos are
up-to-date. You can add files and commit or rollback changes. This is
in no way a full git library. It does require that `git` be installed
and in your `PATH`. It currently only clones via SSH. There is a `git`
submodule for lower-level management. Additionally, you can use
`git.Git()` for any `git` functionality that may be missing.

## How to install

Open a terminal and run the following:

```
$ go get -u github.com/mjwhitta/repomgr
```

## How to use

Below is a sample usage to clone multiple repos to a directory.

```
package main

import (
    "os"
    "path/filepath"

    "github.com/mjwhitta/repomgr"
)

func main() {
    var configDir string
    var e error
    var plugins []repomgr.Repository = []repomgr.Repository{
        &repomgr.Repo{URL: "git@github.com:mjwhitta/repomgr.git"},
    }

    if configDir, e = os.UserConfigDir(); e != nil {
        println(e.Error())
        os.Exit(1)
    }

    configDir = filepath.Join(configDir, "toolname", "plugins")

    if e = repomgr.DownloadAll(configDir, plugins); e != nil {
        println(e.Error())
        os.Exit(1)
    }

    // Do stuff with plugins
}
```

## Links

- [Source](https://github.com/mjwhitta/repomgr)
