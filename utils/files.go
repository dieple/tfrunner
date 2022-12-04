package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"path/filepath"
	"errors"
	"io"

	opts "github.com/dieple/tfrunner/options"
)

func softLinkingFiles(curr_path string, module_path string, bd *opts.BuildData, args opts.Options){

    var_path := filepath.Join(curr_path, bd.ProviderPath, "/variables")
    rel_path, err := filepath.Rel(filepath.Join(curr_path, module_path), var_path)
    if err != nil {
        // panic(err)
        log.Println(err)
    }

    mp := strings.TrimRight(module_path, "\r\n")

    for _, f := range opts.SoftLinkFiles {
        link_file_path := filepath.Join(curr_path, mp, f)
        os.Chdir(mp)
        if _, err := os.Stat(link_file_path); err == nil {
            // file exists - check if it's symbolic link
            // fmt.Println("file exists...")
            fileInfo, err := os.Lstat(link_file_path)
            if err != nil {
                log.Fatal(err)
            }
            if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
                // fmt.Println("File is a symbolic link ... so remove old linking file")
                err := os.Remove(link_file_path)
                if err != nil {
                    log.Fatal(err)
                }
            }
            err = os.Symlink(filepath.Join(rel_path, f), f)
            if err != nil {
                log.Fatalf("Could not create symlink: %v", err)
            }
        } else if errors.Is(err, os.ErrNotExist) {
            // file does NOT exist
            err = os.Symlink(filepath.Join(rel_path, f), f)
            if err != nil {
                log.Fatalf("Could not create symlink: %v", err)
            }
        }
    }

    os.Chdir(mp)

    // setup backend & provider override files
    backend_override := filepath.Join(curr_path, bd.ProviderPath, "variables/config/backend_override.tf")
    providers_override := filepath.Join(curr_path, bd.ProviderPath, "variables/config/providers_override.tf")

    // remove previous runs -- cd {mpdule_path} && rm -f {plan_output_file} && rm -rf .terraform"
    if stat, err2 := exists(opts.PlanOutputFile); err2 == nil && stat {
        err3 := os.Remove(opts.PlanOutputFile)
        if err3 != nil {
            log.Fatal(err3)
        }
    }

    if stat_dir, err4 := exists(".terraform"); err4 == nil && stat_dir {
        err5 := os.RemoveAll(".terraform")
            if err5 != nil {
                log.Fatal(err5)
           }
    }

    // cp override files
    _, err = copy(backend_override, filepath.Join(curr_path, mp, "backend_override.tf"))
    if err != nil {
        log.Fatal(err)
    }
    _, err = copy(providers_override, filepath.Join(curr_path, mp, "providers_override.tf"))
    if err != nil {
        log.Fatal(err)
    }

    // reset back to executable current path for next module run
    os.Chdir(curr_path)
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func copy(src, dst string) (int64, error) {
    sourceFileStat, err := os.Stat(src)
    if err != nil {
        return 0, err
    }

    if !sourceFileStat.Mode().IsRegular() {
        return 0, fmt.Errorf("%s is not a regular file", src)
    }

    source, err := os.Open(src)
    if err != nil {
        return 0, err
    }
    defer source.Close()

    destination, err := os.Create(dst)
    if err != nil {
        return 0, err
    }
    defer destination.Close()
    nBytes, err := io.Copy(destination, source)
    return nBytes, err
}
