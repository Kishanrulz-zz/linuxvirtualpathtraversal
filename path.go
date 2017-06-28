package main

import (
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"bufio"
)

//Folder struct is used to hold one folder
type Folder struct {
	Name    string
	Folders map[string]*Folder
	Parent  *Folder

}

//reset clears the session and initializes working directory from root
func reset () *Folder{
	rootFolder := &Folder{"/", make(map[string]*Folder), nil}
	return rootFolder
}

// newFolder returns a newly created folder
func newFolder(name string, parentFolder *Folder) *Folder {
	return &Folder{name, make(map[string]*Folder), parentFolder}
}

// addFolder adds a folder to the current directory/folder
func addFolder(folder *Folder, path []string) *Folder {
	head := folder	//head stores the head reference of the folder
	currFolder := folder	// currFolder is used to traverse the folder to create relative folder to the head folder
	i := 0
	pathLen := len(path)
	for {
		// No directory can be created with the directory name "/" as "/" is used for root directory
		if path[i] == "/" {
			fmt.Println("SORRY CANNOT INITIALIZE ROOT DIRECTORY")
			break
		}
		if pathLen == 1 {	//a folder is created when the segment is the element in the path array
			// returning if folder already exists with a proper message
			if _, ok := currFolder.Folders[path[i]]; ok {
				fmt.Println("DIRECTORY ALREADY EXIST")
				break
			}
			// adding a directory
			newDirectory := newFolder(path[i], currFolder)
			currFolder.Folders[path[i]] = newDirectory
			pathLen--
			i++
			fmt.Println("SUCC: CREATED")
			break
		}
		// validating whether an intermediate directory already exists
		if _, ok := currFolder.Folders[path[i]]; ok {
			currFolder = currFolder.Folders[path[i]]
			pathLen--
			i++
			continue
		} else {	// returning as no intermediate directory exists
			fmt.Println("INVALID PATH")
			break
		}
	}
	return head	//returning back the head pointer
}

// getWorkingDirectory returns the current working directory
func getWorkingDirectory(folder *Folder) {
	s := "/"
	for {
		if folder.Parent == nil {
			fmt.Println("PWD: ",s)
			break
		}
		s =  "/" + folder.Name + s
		folder = folder.Parent
	}
}

// changeDirectory is used to change the directory
func changeDirectory(folder *Folder, path [] string) *Folder {
	pathLen := len(path)
	i := 0
	for {
		if pathLen == 0 {	//if pathLen equals 0 we have reached the final destination
			fmt.Println("SUCC: REACHED")
			break
		}
		if path[i] == ".." { //handling basic linux path traversal for getting one directory ahead of the current one
			folder = folder.Parent
			pathLen--
			i++
			continue
		}
		//if directory doesn't exists, returning with a proper message
		if _, ok := folder.Folders[path[i]]; !ok {
			fmt.Println("ERR: INVALID DIRECTORY")
			break
		}
		//traversing down the folder
		folder = folder.Folders[path[i]]
		pathLen--
		i++
	}
	fmt.Println(folder.Name)
	return folder	//returning the current directory after the traversal
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("SIMPLE SHELL")
	fmt.Println("---------------------")
	folder := reset()
	fmt.Println("WELCOME TO HOME DIRECTORY")
	fmt.Println(folder.Name)

	for {

		command := make([]string, 2)
		text, _ := reader.ReadString('\n')

		if strings.Contains(text, " ") {
			command = strings.Split(text, " ")
		} else {
			command[0] = text
		}

		cmd := strings.Trim(command[0], "\n")
		action := strings.Trim(command[1], "\n")

		switch cmd {
		case "clear":
			folder = reset()
			fmt.Println("WELCOME TO HOME DIRECTORY")
			fmt.Println(folder.Name)

		case "mkdir":
			pathArr := strings.Split(action, string(filepath.Separator))
			folder = addFolder(folder,pathArr)

		case "pwd":
				getWorkingDirectory(folder)

		case "cd":
			pathArr := strings.Split(action, string(filepath.Separator))
			folder = changeDirectory(folder, pathArr)

		case "ls":
			var dir string
			for key, _ := range folder.Folders {
				dir = key + " " + dir
			}
			fmt.Println("DIR: ", dir)
		}
	}
}