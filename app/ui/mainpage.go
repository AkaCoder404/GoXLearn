package ui

import (
	"context"
	"fmt"
	"log"
	"main/app/core"
	"main/app/ui/utils"
	"math"
	"time"

	"github.com/rivo/tview"
)

const (
	offsetRange = 100
	loadDelay   = time.Millisecond * 50
	maxOffset   = 10000
)

// MainPage: This struct contains the grid an the entry table
type MainPage struct {
	Grid                *tview.Grid  // the page grid
	CurrentCoursesTable *tview.Table // the table that contains the list of current courses
	AllCoursesTable     *tview.Table // the table that contains the list of all courses

	CurrentOffset int
	MaxOffset     int

	cWrap *utils.ContextWrapper // For context cancellation
}

// ShowMainPage : Make the app show the main page
func ShowMainPage() {
	// Creating the new main page
	log.Println("Creating new main page...")
	mainPage := newMainPage()

	core.App.TView.SetFocus(mainPage.Grid)
	core.App.PageHolder.AddAndSwitchToPage(utils.MainPageID, mainPage.Grid, true)
}

// newMainPage : Creates a new main page
func newMainPage() *MainPage {
	var dimensions []int
	for i := 0; i < 15; i++ {
		dimensions = append(dimensions, -1)
	}

	grid := utils.NewGrid(dimensions, dimensions)
	// Set grid attributes.
	grid.SetTitleColor(utils.MainPageGridTitleColor).
		SetBorderColor(utils.MainPageGridBorderColor).
		SetBorder(true)

	// Create the base main table
	table := tview.NewTable()
	// Set table attributes
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(utils.MainPageTableBorderColor).
		SetTitleColor(utils.MainPageTableTitleColor).
		SetBorder(true)

	// Create all courses table
	allCoursesTable := tview.NewTable()
	allCoursesTable.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(utils.MainPageTableBorderColor).
		SetTitleColor(utils.MainPageTableTitleColor).
		SetTitle("All Courses").
		SetBorder(true)

	// Add the table to the grid. Table spans the whole page.
	grid.AddItem(table, 0, 0, 5, 15, 0, 0, true).
		AddItem(allCoursesTable, 5, 0, 10, 15, 0, 0, true)
	// AddItem(table, 0, 0, 5, 15, 0, 80, true).
	// AddItem(allCoursesTable, 5, 0, 10, 15, 0, 80, true)

	ctx, cancel := context.WithCancel(context.Background())
	mainPage := &MainPage{
		Grid:                grid,
		CurrentCoursesTable: table,
		AllCoursesTable:     allCoursesTable,
		cWrap: &utils.ContextWrapper{
			Ctx:    ctx,
			Cancel: cancel,
		},
	}

	// Set up the MainPage for user
	mainPage.setPage()
	return mainPage
}

// setPage: Set up the MainPage for a user
func (p *MainPage) setPage() {
	log.Println("Setting up main page...")
	go p.setPageGrid()
	go p.setPageTable()
	go p.setAllCoursesPageTable()
}

// setPageGrid : Show grid title
func (p *MainPage) setPageGrid() {
	// Get User Information
	name, studentTypes, department, err := core.App.Client.User.GetUserInformation()
	if err != nil {
		log.Printf("Got User Info: %s %s %s", name, department, studentTypes)
		log.Printf("Error getting user info: %s", err.Error())
	}

	core.App.TView.QueueUpdateDraw(func() {
		p.Grid.SetTitle(fmt.Sprintf("Welcome to GoXLearn CLI, [lightgreen]%s!", name))
	})
	log.Println("Finished setting grid")
}

// setPageTable : Show table items and title: current semesters courses
func (p *MainPage) setPageTable() {
	ctx, cancel := p.cWrap.ResetContext()

	// Set Handlers
	p.setHandlers(cancel, nil)
	time.Sleep(loadDelay)
	defer cancel()

	// Get this semesters id
	tableTitle := "Classes"
	core.App.TView.QueueUpdateDraw(func() {
		// Clear current entries
		p.CurrentCoursesTable.Clear()

		// Set headers.
		titleHeader := tview.NewTableCell("课程名").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageTitleColor).
			SetSelectable(false)
		englishClassNameHeader := tview.NewTableCell("英文课程名").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageDescColor).
			SetSelectable(false)
		unreadNoticeHeader := tview.NewTableCell("未读公告").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageDescColor).
			SetSelectable(false)
		unreadClassMaterialsHeader := tview.NewTableCell("未浏览课件").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageDescColor).
			SetSelectable(false)
		courseIdHeader := tview.NewTableCell("课程号").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageDescColor).
			SetSelectable(false)
		profHeader := tview.NewTableCell("Professor").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageTagColor).
			SetSelectable(false)
		classTimeHeader := tview.NewTableCell("Class Time").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageTagColor).
			SetSelectable(false)

		p.CurrentCoursesTable.SetCell(0, 0, titleHeader).
			SetCell(0, 1, englishClassNameHeader).
			SetCell(0, 2, courseIdHeader).
			SetCell(0, 3, unreadNoticeHeader).
			SetCell(0, 4, unreadClassMaterialsHeader).
			SetCell(0, 5, profHeader).
			SetCell(0, 6, classTimeHeader).
			SetFixed(1, 0)

		// Set table title.
		// page, first, last := p.calculatePaginationData()
		p.CurrentCoursesTable.SetTitle(fmt.Sprintf("%s [::bu] Loading...", tableTitle))
	})

	// Get list of courses
	if p.cWrap.ToCancel(ctx) {
		return
	}

	semesterID, err := core.App.Client.Class.GetCurrentAndNextSemester()
	if err != nil {
		log.Println("Semester id retreival failed", err.Error())
		return
	}

	list, err := core.App.Client.Class.GetCourseList(semesterID, "student")
	if err != nil {
		log.Println("Course list retreival failed", err.Error())
		return
	}

	// TODO: Limit Total of Entries? maybe when viewing all courses ...
	p.MaxOffset = int(math.Min(float64(len(list.ResultList)), maxOffset))

	//  Update table title with current semester
	core.App.TView.QueueUpdateDraw(func() {
		p.CurrentCoursesTable.SetTitle(fmt.Sprintf("Classes for %s", semesterID))
	})

	// Fill in table details
	for index := 0; index < len(list.ResultList); index++ {
		// cancel in the middle of loading?
		if p.cWrap.ToCancel(ctx) {
			return
		}
		class := list.ResultList[index]

		// 课程名
		classNameCell := tview.NewTableCell(fmt.Sprintf("%-10s", class.Kcm)).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageTitleColor).SetReference(&class)

		// 课程英文名
		classEnglishNameCell := tview.NewTableCell(fmt.Sprintf("%-10s", class.Ywkcm)).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageDescColor).SetReference(&class)
		// 课程号
		classIdCell := tview.NewTableCell(fmt.Sprintf("%-10s", class.Kch+"-"+class.Kxhnumber)).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageDescColor).SetReference(&class)

		// 未读公告
		unreadNoticeCell := tview.NewTableCell(fmt.Sprintf("%-10s", fmt.Sprint(class.Xggs)+"/"+fmt.Sprint(class.Ggzs))).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageDescColor).SetReference(&class)

		// 未读课件
		unreadClassMaterialsCell := tview.NewTableCell(fmt.Sprintf("%-10s", fmt.Sprint(class.Xkjs)+"/"+fmt.Sprint(class.Xskjs))).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageDescColor).SetReference(&class)

		// 教授名称 cell.
		professorCell := tview.NewTableCell(fmt.Sprintf("%-10s", class.Jsm)).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageTagColor).SetReference(&class)

		// 课程教师号 Cell
		professorIdCell := tview.NewTableCell(fmt.Sprintf("%-10s", class.Jsh)).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageTagColor).SetReference(&class)

		p.CurrentCoursesTable.SetCell(index+1, 0, classNameCell).
			SetCell(index+1, 1, classEnglishNameCell).
			SetCell(index+1, 2, classIdCell).
			SetCell(index+1, 3, unreadNoticeCell).
			SetCell(index+1, 4, unreadClassMaterialsCell).
			SetCell(index+1, 5, professorCell).
			SetCell(index+1, 6, professorIdCell)
	}

	core.App.TView.QueueUpdateDraw(func() {
		p.CurrentCoursesTable.Select(1, 0)
		p.CurrentCoursesTable.ScrollToBeginning()
	})

	log.Println("Finish setting current classes page table")
}

// setAllCouresPageTable : Show table of all courses.
func (p *MainPage) setAllCoursesPageTable() {
	ctx, cancel := p.cWrap.ResetContext()

	// Set Handlers
	p.setHandlers(cancel, nil)
	time.Sleep(loadDelay)
	defer cancel()

	// Get this semesters id
	// tableTitle := "Classes"
	core.App.TView.QueueUpdateDraw(func() {
		// Clear current entries
		p.AllCoursesTable.Clear()

		// Set headers.
		titleHeader := tview.NewTableCell("课程名").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageTitleColor).
			SetSelectable(false)
		mainProfessor := tview.NewTableCell("主讲老师").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageDescColor).
			SetSelectable(false)
		semester := tview.NewTableCell("学年学期").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageDescColor).
			SetSelectable(false)
		myHomeworks := tview.NewTableCell("我的作业").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageDescColor).
			SetSelectable(false)
		myQuestions := tview.NewTableCell("我的答疑").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageDescColor).
			SetSelectable(false)
		myDiscussions := tview.NewTableCell("我的讨论").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageDescColor).
			SetSelectable(false)

		p.AllCoursesTable.SetCell(0, 0, titleHeader).
			SetCell(0, 1, mainProfessor).
			SetCell(0, 2, semester).
			SetCell(0, 3, myHomeworks).
			SetCell(0, 4, myQuestions).
			SetCell(0, 5, myDiscussions).
			SetFixed(1, 0)

		// Set table title.
		// page, first, last := p.calculatePaginationData()
		p.AllCoursesTable.SetTitle(fmt.Sprintf("All classes [::bu] Loading..."))
	})

	// Get list of courses
	if p.cWrap.ToCancel(ctx) {
		return
	}

	list, err := core.App.Client.Class.GetAllClasses()
	if err != nil {
		log.Println("Course list retreival failed", err.Error())
		return
	}

	// TODO: Limit Total of Entries? maybe when viewing all courses ...
	// p.MaxOffset = int(math.Min(float64(len(list.ResultList)), maxOffset))

	core.App.TView.QueueUpdateDraw(func() {
		p.AllCoursesTable.SetTitle(fmt.Sprintf("All classes"))
	})

	// Fill in table details
	for index := 0; index < len(list.Object.AaData); index++ {
		// cancel in the middle of loading?
		if p.cWrap.ToCancel(ctx) {
			return
		}
		class := list.Object.AaData[index]

		// 课程名
		classNameCell := tview.NewTableCell(fmt.Sprintf("%-10s", class.Kcm)).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageTitleColor).SetReference(&class)

		// 主讲教师
		classProfessorCell := tview.NewTableCell(fmt.Sprintf("%-10s", class.Jsm)).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageDescColor).SetReference(&class)

		// 学年学期
		semesterCell := tview.NewTableCell(fmt.Sprintf("%-10s", class.Xnxq)).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageDescColor).SetReference(&class)

		// 我的作业
		myHomeworkCell := tview.NewTableCell(fmt.Sprintf("%-10s", fmt.Sprint(class.Zys))).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageDescColor).SetReference(&class)

		// 未读课件
		myQuestionsCell := tview.NewTableCell(fmt.Sprintf("%-10s", fmt.Sprint(class.Xsdys))).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageDescColor).SetReference(&class)

		// 教授名称 cell.
		myDiscussionsCell := tview.NewTableCell(fmt.Sprintf("%-10s", fmt.Sprint(class.Tls))).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageDescColor).SetReference(&class)

		p.AllCoursesTable.SetCell(index+1, 0, classNameCell).
			SetCell(index+1, 1, classProfessorCell).
			SetCell(index+1, 2, semesterCell).
			SetCell(index+1, 3, myHomeworkCell).
			SetCell(index+1, 4, myQuestionsCell).
			SetCell(index+1, 5, myDiscussionsCell)
	}

	core.App.TView.QueueUpdateDraw(func() {
		// p.AllCoursesTable.Select(1, 0)
		p.AllCoursesTable.ScrollToBeginning()
	})

	log.Println("Finish setting all courses table")

}
