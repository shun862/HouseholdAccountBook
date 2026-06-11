package consts

// ルーティング情報
const (
	UserRegisterUrl  = "/user_register"
	LoginUrl         = "/login"
	LogoutUrl        = "/logout"
	AddExpenseUrl    = "/user/add_expense"
	ExpenseListUrl   = "/user/list"
	ExpenseReportUrl = "/user/report"
)

// HTMLファイル情報
const (
	filePath          = "../../web/templates"
	UserRegisterFile  = filePath + "/user_register.html"
	LoginFile         = filePath + "/login.html"
	AddExpenseFile    = filePath + "/add_expense.html"
	ExpenseListFile   = filePath + "/list.html"
	ExpenseReportFile = filePath + "/report.html"
	HeaderFile        = filePath + "/header.html"
)

// ログイン後画面の画面名
const (
	AddExpenseViewName    = "add_expense"
	ExpenseListViewName   = "list"
	ExpenseReportViewName = "report"
)

// 認証情報
const JwtKey = "anhifsjsa32gjjis"
const TokenName = "auth_token"
const CookieValidPath = "/"
const UserIDKey = "userID"

// フォーマット
const DateFormatH = "2006-01-02"
const DateFormatS = "2006/01/02"
const DateFormatJ = "2006年1月"

const StoreKey = "store_key"
const SessionName = "session_name"
