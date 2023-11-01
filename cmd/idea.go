/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"strings"

	"ungitz/util"

	"github.com/spf13/cobra"
)

// codeCmd represents the code command
var ideaCmd = &cobra.Command{
	Use:   "idea",
	Short: "It will open the directory in Intellij Idea",
	Long: `This command will help to open the unzipped folder to Intellij Idea.
	In order for this command to work, Intellij Idea should be installed in your system`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(File_flag) < 1 && link_flag == "" && len(args) < 1 {
			return errors.New("accept(s) 1 argument")
		}
		return nil
	},
	Example: `ungitz code -f <filename>,<repo name>,<branch name>
ungitz code -l <URL>`,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		var fileName string
		var err error
		var argument string
		var link_arg string
		var name_arg string
		var branch_arg string
		var repo_name string
		var err1 error
		var fileName1 string
		var branch_name string
		var repo_name1 string

		// flag check
		if len(File_flag) != 0 {
			argument = File_flag[0]
			repo_name1 = File_flag[1]
			branch_name = File_flag[2]
			repo_name = repo_name1 + "-" + branch_name

		} else if link_flag != "" {
			link_arg = link_flag
			name_arg = util.RegexFilter(link_arg, fname_pattern)
			branch_arg = strings.TrimSuffix(util.RegexFilter(link_arg, branch_pattern), ".zip")

			// wait period implementation using go-routines for download function
			wg.Add(1)
			go func(name_arg, link_arg string) {
				util.Download(name_arg, link_arg)
				wg.Done()
			}(name_arg+".zip", link_arg)
			wg.Wait()
			argument = name_arg + ".zip"
		}

		// file exist check
		FileExists, err := util.FileExists(argument)
		if err != nil {
			fmt.Println(err.Error())
		}
		if FileExists {
			fileName, err = filepath.Abs(argument)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Printf("File %v does not exist", argument)
			return
		}

		// initialisation of working directory
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
		}

		util.Unzip(fileName, wd)

		if link_flag != "" {
			var testname = util.FilenameWithoutExtension(fileName)
			os.Chdir(testname + "-" + branch_arg)
		} else if len(File_flag) != 0 {
			fileName1, err1 = filepath.Abs(repo_name + ".zip")
			if err1 != nil {
				fmt.Println(err1.Error())
			}
			var testname1 = util.FilenameWithoutExtension(fileName1)
			os.Chdir(testname1)
		}

		// updation of working directory
		wd, err = os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
		}

		commandCode := exec.Command("idea", wd)
		err = commandCode.Run()

		if err != nil {
			fmt.Println("Intellij Idea executable not found in %PATH%")
		} else {
			fmt.Println("Unzipping and opening file.")
		}

	},
}

func init() {
	rootCmd.AddCommand(ideaCmd)
	ideaCmd.PersistentFlags().StringSliceVarP(&File_flag, "file", "f", []string{}, "Arguments:<file name>,<repo name>,<branch name>")
	ideaCmd.PersistentFlags().StringVarP(&link_flag, "link", "l", "", "Argument:<URL>")
}
