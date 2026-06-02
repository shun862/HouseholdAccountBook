function validateAuth(formId) {
  document.getElementById(formId).addEventListener("submit", function (e) {
    let isValid = true;
    const username = document.getElementById("username");
    const password = document.getElementById("password");
    const usernameError = document.getElementById("usernameError");
    const passwordError = document.getElementById("passwordError");

    usernameError.textContent = "";
    passwordError.textContent = "";

    // ユーザー名チェック
    if (username.value.trim() === "") {
      usernameError.textContent = "ユーザー名を入力してください";
      isValid = false;
    } else if (username.value.length < 4 || username.value.length > 10) {
      usernameError.textContent = "ユーザー名は4～10文字で入力してください";
      isValid = false;
    }

    // パスワードチェック
    if (password.value.trim() === "") {
      passwordError.textContent = "パスワードを入力してください";

      isValid = false;
    } else if (password.value.length < 8 || password.value.length > 12) {
      passwordError.textContent = "パスワードは8～12文字で入力してください";
      isValid = false;
    }

    if (!isValid) {
      e.preventDefault();
    }
  });
}

function togglePassword() {
  const password = document.getElementById("password");
  const eye = document.getElementById("eyeIcon");

  if (password.type === "password") {
    password.type = "text";
    eye.classList.remove("bi-eye");
    eye.classList.add("bi-eye-slash");
  } else {
    password.type = "password";
    eye.classList.remove("bi-eye-slash");
    eye.classList.add("bi-eye");
  }
}
