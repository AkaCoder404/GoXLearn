# XLearnCLI

This is a CLI terminal application for Tsinghua Universities 网络学堂 portal. It supports [learn2018](https://learn.tsinghua.edu.cn). It uses the go library [gothulearn](https://github.com/AkaCoder404/gothulearn)

*This was built for learning purposes, to become more familiar with go programming language and building a terminal application*

**Keywords:** golang, tview, command-line application, THU, Tsinghua Unversity

## Installation

## Usage

### Login Page/Logout
If remember, saves credentials, and tries to load credentials next time. 

![](https://s2.loli.net/2023/02/15/xFCEGZvgaO9QnXi.png)

In order to logout, use the shortcut CTRL+L

![](https://s2.loli.net/2023/02/15/FpbKduVSgWiLZT4.png)

### Main Page
This is the first page after login, it lists all courses for this semester and previous semesters
![](https://s2.loli.net/2023/02/15/aVOc79mU5AzQj6S.png)


### Class Page
It incorporates a text view for class information, and three tables for class notices, homeworks, and class materials

![](https://s2.loli.net/2023/02/15/GrEenfuApkZQcWP.png)

In order to return to main page, press ESC
 
### Files Page
*in progress*


### Notice Page
*in progress*

### Homework Page
*in progress*

## TODO
1. Keep session alive (learn timeout after a certain time, makes cookie unusable and csrf unusable)
2. Download files
3. Handle all homework views (unsubmitted, submmited, graded)
4. Support more shortcuts to switch between pages

