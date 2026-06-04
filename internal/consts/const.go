package consts

// ルーティング情報
const (
	UserRegisterUrl = "/user_register"
	LoginUrl        = "/login"
	LogoutUrl       = "/logout"
	AddExpenseUrl   = "/user/add_expense"
)

// HTMLファイル情報
const (
	filePath         = "../../web/templates"
	UserRegisterFile = filePath + "/user_register.html"
	LoginFile        = filePath + "/login.html"
	AddExpenseFile   = filePath + "/add_expense.html"
	HeaderFile       = filePath + "/header.html"
)

// ログイン後画面の画面名
const (
	AddExpenseViewName = "add_expense"
)

// 認証情報
const JwtKey = "anhifsjsa32gjjis"
const TokenName = "auth_token"
const CookieValidPath = "/"
