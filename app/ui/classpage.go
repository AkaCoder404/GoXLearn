package ui

import (
	"context"
	"fmt"
	"log"
	"main/app/core"
	"main/app/ui/utils"
	"time"

	"github.com/AkaCoder404/gothulearn"
	"github.com/rivo/tview"
)

// ClassPage : This struct contains the required primitives for class page
type ClassPage struct {
	Class     *gothulearn.Class
	Grid      *tview.Grid
	GgTable   *tview.Table // 公告
	HwTable   *tview.Table // 作业
	FileTable *tview.Table // 文件

	Info *tview.TextView

	sWrap *utils.SelectorWrapper
	cWrap *utils.ContextWrapper // For context cancellation
}

// ShowClassPage : Make the app show the class page
func ShowClassPage(class *gothulearn.Class) {
	classPage := newClassPage(class)
	core.App.TView.SetFocus(classPage.Grid)
	core.App.PageHolder.AddAndSwitchToPage(utils.ClassPageID, classPage.Grid, true)
}

// newClassPage : Creates a new class page
func newClassPage(class *gothulearn.Class) *ClassPage {
	// Size of grid
	var dimensions []int
	for i := 0; i < 15; i++ {
		dimensions = append(dimensions, -1)
	}

	grid := utils.NewGrid(dimensions, dimensions)
	// Set grid attributes
	grid.SetTitleColor(utils.ClassPageTitleColor).
		SetBorderColor(utils.ClassPageGridBorderColor).
		SetTitle("Class Information").
		SetBorder(true)

	// Use a Textview for 课程信息
	info := tview.NewTextView()
	// Set textview attributes
	info.SetWrap(true).SetWordWrap(true).
		SetBorderColor(utils.ClassPageInfoViewBorderColor).
		SetTitleColor(utils.ClassPageInfoViewTitleColor).
		SetTitle("About").
		SetBorder(true)

	// TODO: Use tables to show 公告 information
	table := tview.NewTable()
	noticeHeader := tview.NewTableCell("公告").
		SetTextColor(utils.ClassPageTitleColor).
		SetSelectable(false)

	readStatus := tview.NewTableCell("阅读").
		SetTextColor(utils.ClassPageChapNumColor).
		SetSelectable(false)

	publisherHeader := tview.NewTableCell("发布人").
		SetTextColor(utils.ClassPageChapNumColor).
		SetSelectable(false)

	noticeDateHeader := tview.NewTableCell("发布时间").
		SetTextColor(utils.ClassPageChapNumColor).
		SetSelectable(false)

	table.SetCell(0, 0, noticeHeader).
		SetCell(0, 1, readStatus).
		SetCell(0, 2, publisherHeader).
		SetCell(0, 3, noticeDateHeader).
		SetFixed(1, 0)

	// Set table attributes
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(utils.ClassPageTableBorderColor).
		SetTitle("公告").
		SetTitleColor(utils.ClassPageTableTitleColor).
		SetBorder(true)

	// TODO set homework table
	homeworkTable := tview.NewTable()
	// 作业题目
	homeworkHeader := tview.NewTableCell("作业题目").
		SetTextColor(utils.ClassPageTitleColor).
		SetSelectable(false)

	// 完成方式
	completionMethodHeader := tview.NewTableCell("完成方式").
		SetTextColor(utils.ClassPageChapNumColor).
		SetSelectable(false)

	// 截止时间
	dueDateHeader := tview.NewTableCell("截止时间").
		SetTextColor(utils.ClassPageChapNumColor).
		SetSelectable(false)

	homeworkTable.SetCell(0, 0, homeworkHeader).
		SetCell(0, 1, completionMethodHeader).
		SetCell(0, 2, dueDateHeader).
		SetFixed(1, 0)

	homeworkTable.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(utils.ClassPageTableBorderColor).
		SetTitle("我的作业").
		SetTitleColor(utils.ClassPageTableTitleColor).
		SetBorder(true)

	// TODO file table
	fileTable := tview.NewTable()

	fileHeader := tview.NewTableCell("文件").
		SetTextColor(utils.ClassPageTitleColor).
		SetSelectable(false)

	fileSize := tview.NewTableCell("文件大小").
		SetTextColor(utils.ClassPageChapNumColor).
		SetSelectable(false)

	fileTable.SetCell(0, 0, fileHeader).
		SetCell(0, 1, fileSize).
		SetFixed(1, 0)

	fileTable.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(utils.ClassPageTableBorderColor).
		SetTitle("课程文件").
		SetTitleColor(utils.ClassPageTableTitleColor).
		SetBorder(true)

	// Add info and table to the grid. Set the focus to the chapter table.
	// grid.AddItem(info, 0, 0, 5, 15, 0, 0, false).
	// 	AddItem(table, 5, 0, 10, 15, 0, 0, true).
	// 	AddItem(info, 0, 0, 15, 5, 0, 80, false).
	// 	AddItem(table, 0, 5, 15, 10, 0, 80, true)

	grid.AddItem(info, 0, 0, 5, 15, 0, 0, false).
		AddItem(table, 5, 0, 10, 15, 0, 0, true).
		AddItem(homeworkTable, 5, 0, 10, 15, 0, 0, false).
		AddItem(fileTable, 5, 0, 10, 15, 0, 0, false).
		AddItem(info, 0, 0, 15, 5, 0, 80, false).
		AddItem(table, 0, 5, 5, 10, 0, 80, true).
		AddItem(homeworkTable, 5, 5, 5, 10, 0, 80, false).
		AddItem(fileTable, 10, 5, 5, 10, 0, 80, false)

	// Construct page
	ctx, cancel := context.WithCancel(context.Background())

	classPage := &ClassPage{
		Class:     class,
		Grid:      grid,
		Info:      info,
		GgTable:   table,
		HwTable:   homeworkTable,
		FileTable: fileTable,
		// TODO: Comment
		sWrap: &utils.SelectorWrapper{
			Selection: map[int]struct{}{},
		},
		// TODO: Comment
		cWrap: &utils.ContextWrapper{
			Ctx:    ctx,
			Cancel: cancel,
		},
	}

	// Set up values
	go classPage.setClassInfo()
	go classPage.setClassNotice()
	go classPage.setClassHomework()
	go classPage.setClassMaterial()
	go classPage.setFrontPageInfo()
	return classPage
}

// setClassInfo : Set up 课程信息
func (p *ClassPage) setClassInfo() {
	log.Println("Setting up class information...")
	// Get information on selected class
	classInfo, err := core.App.Client.Class.GetClassInformation(p.Class.Wlkcid, "student")
	if err != nil {
		log.Println("Failed to setup class information")
		core.App.TView.QueueUpdateDraw(func() {
			p.Info.SetText(err.Error())
		})
		return
	}

	// Set up information text
	className := classInfo.Title
	mainProfessor := classInfo.MainProfessor
	hekaiProfessor := classInfo.HekaiProfessor
	classCredits := classInfo.ClassCredits
	classHours := classInfo.ClassHours
	classScope := classInfo.ClassScope
	classMaterialsCount := classInfo.ClassMaterialsCount
	classHomeworkCount := classInfo.ClassHomeworkCount
	classDiscussionCount := classInfo.ClassDiscussionCount

	classDescription := classInfo.ClassDescription
	classEnglishDescription := classInfo.ClassEnglishDescription
	classSchedule := classInfo.ClassSchedule
	classAssesmentMethod := classInfo.ClassAssesmentMethod
	classReferenceMaterials := classInfo.ClassReferenceMaterials
	classProfessor := classInfo.ClassProfessor
	classSelectionGuidance := classInfo.ClassSelectionGuidance
	classPrerequisites := classInfo.ClassPrerequisites
	classOpenOfficeHour := classInfo.ClassOpenOfficeHour
	classGradingStandard := classInfo.ClassGradingStandard
	teacherTeachingCharacteristics := classInfo.TeacherTeachingCharacteristics

	// TODO add color to the text
	infoText := fmt.Sprintf(
		"基本信息\n课程名: %s \n\n主讲教师: %s\n合开教师: %s\n课程学分: %s\n课程学时: %s\n课程开放范围: %s\n课程文件数: %s\n部署作业数: %s\n讨论贴数: %s\n\n课程简介\n%s\n\n课程英文简介\n%s\n\n课程时间安排\n%s\n\n课程考核方式\n%s\n\n教材参考资料\n%s\n\n教师教学特点\n%s\n\n教师简介\n%s\n\n课程选课指南\n%s\n\n课程先修要求\n%s\n\n课程开放办公时间\n%s\n\n课程评分标准\n%s\n\n",
		className, mainProfessor, hekaiProfessor, classCredits, classHours, classScope, classMaterialsCount, classHomeworkCount, classDiscussionCount,
		classDescription, classEnglishDescription, classSchedule, classAssesmentMethod, classReferenceMaterials, teacherTeachingCharacteristics,
		classProfessor, classSelectionGuidance, classPrerequisites, classOpenOfficeHour, classGradingStandard)

	core.App.TView.QueueUpdateDraw(func() {
		p.Info.SetText(infoText)
	})
}

// setClassNotice : Set up 公告 information
func (p *ClassPage) setClassNotice() {
	log.Println("Setting up class notice...")

	// Get 公告 information on class
	noticeInfo, err := core.App.Client.Notice.GetNoticeList(p.Class.Wlkcid)
	if err != nil {
		log.Println("Failed to setup class notice")
		return
	}

	// Set up 公告 table
	for index, notice := range noticeInfo.Object.AaData {
		// 公告标题
		noticeNameCell := tview.NewTableCell(fmt.Sprintf("%-10s", notice.Bt)).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageTitleColor).SetReference(&notice)

		// 公告阅读状态
		noticeReadCell := tview.NewTableCell(fmt.Sprintf("%-10s", notice.Sfyd)).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageDescColor).SetReference(&notice)

		// 公告发布人
		noticePublisherCell := tview.NewTableCell(fmt.Sprintf("%-10s", notice.Fbr)).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageDescColor).SetReference(&notice)

		// 公告发布时间
		noticeDateCell := tview.NewTableCell(fmt.Sprintf("%-10s", notice.Fbsj)).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageDescColor).SetReference(&notice)

		p.GgTable.SetCell(index+1, 0, noticeNameCell).
			SetCell(index+1, 1, noticeReadCell).
			SetCell(index+1, 2, noticePublisherCell).
			SetCell(index+1, 3, noticeDateCell)
	}

	core.App.TView.QueueUpdateDraw(func() {
		p.GgTable.Select(1, 0)
		p.GgTable.ScrollToBeginning()
	})

	log.Println("Finish setting class 公告 table")
}

// setClassHomework : Set up 作业 information
func (p *ClassPage) setClassHomework() {
	log.Println("Setting up class homework...")

	// Get 作业 information on class
	// TODO Handle different homeworks, not submitted, submitted, graded,
	// only showing not submitted as of now
	unsubmittedHomeworks, err := core.App.Client.Homework.GetUnsubmittedHomeworks(p.Class.Wlkcid)
	if err != nil {
		log.Println("Failed to setup class homework")
		return
	}

	// Set up 作业 table
	for index, homework := range unsubmittedHomeworks.Object.AaData {
		// 作业标题
		homeworkNameCell := tview.NewTableCell(fmt.Sprintf("%-10s", homework.Bt)).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageTitleColor).SetReference(&homework)

		// 完成方式
		// TODO handle 作业完成方式
		homeworkCompleteTypeCell := tview.NewTableCell(fmt.Sprintf("%-10s", "个人")).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageDescColor).SetReference(&homework)

		// 作业截止时间
		homeworkDeadlineCell := tview.NewTableCell(fmt.Sprintf("%-10s", homework.JzsjStr)).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageDescColor).SetReference(&homework)

		p.HwTable.SetCell(index+1, 0, homeworkNameCell).
			SetCell(index+1, 1, homeworkCompleteTypeCell).
			SetCell(index+1, 2, homeworkDeadlineCell)

	}

}

// setClassMaterial : Set up 课程资料 information
func (p *ClassPage) setClassMaterial() {
	log.Println("Setting up class material...")

	// Get 课程资料 information on class
	// TODO Handle all files of different tabs
	pageFile, err := core.App.Client.File.GetFilePageList(p.Class.Wlkcid, "student")
	if err != nil {
		log.Println("Failed to get file page list")
		return
	}
	materialInfo, err := core.App.Client.File.GetFileList(p.Class.Wlkcid, pageFile.Object.Rows[0].Kjflid, "student")
	if err != nil {
		log.Println("Failed to get file list")
		return
	}

	// Set up 课程资料 table
	for index, material := range materialInfo.Object {
		// 课程资料标题
		materialNameCell := tview.NewTableCell(fmt.Sprintf("%-10s", material.Filename)).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageTitleColor).SetReference(&material)

		// 课程资料大小
		materialSizeCell := tview.NewTableCell(fmt.Sprintf("%-10s", fmt.Sprint(material.FileSize))).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageDescColor).SetReference(&material)

		p.FileTable.SetCell(index+1, 0, materialNameCell).
			SetCell(index+1, 1, materialSizeCell)
	}

	core.App.TView.QueueUpdateDraw(func() {
		p.FileTable.Select(1, 0)
		p.FileTable.ScrollToBeginning()
	})

	log.Println("Finish setting class 课程资料 table")
}

// setFrontPageInfo : Set up 课程首页 information
func (p *ClassPage) setFrontPageInfo() {
	ctx, cancel := p.cWrap.ResetContext()

	// Set handlers
	p.setHandlers(cancel)
	time.Sleep(loadDelay)
	defer cancel()

	// Show loading status
	// core.App.TView.QueueUpdateDraw(func() {
	// 	loadingCell := tview.NewTableCell("Loading...").SetSelectable(false)
	// 	p.GgTable.SetCell(1, 1, loadingCell)
	// })

	if p.cWrap.ToCancel(ctx) {
		return
	}
}
