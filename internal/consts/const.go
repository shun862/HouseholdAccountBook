package consts

// ルーティング情報
const (
	UserRegisterUrl = "/user_register"
	LoginUrl        = "/login"
	AddExpenseUrl   = "/user/add_expense"
)

// HTMLファイル情報
const (
	UserRegisterFile = "../../web/templates/user_register.html"
	LoginFile        = "../../web/templates/login.html"
	AddExpenseFile   = "../../web/templates/add_expense.html"
)

// 認証情報
const JwtKey = "anhifsjsa32gjjis"
const TokenName = "auth_token"
const CookiePath = "/user"
