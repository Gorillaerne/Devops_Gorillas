import { callLoginRestApi } from "./api_calls.js";
import { checkIfLoggedIn, showError, showSuccess } from "./reuseable_functions.js";

checkIfLoggedIn()

const loginForm = document.getElementById('login-form');

loginForm.addEventListener('submit', (e) => {
    e.preventDefault();

    const userVal = document.getElementById('username').value;
    const passVal = document.getElementById('password').value;

    callLoginRestApi(userVal, passVal)
        .then(data => {
            localStorage.setItem("token", data.token);
            if (data.breached) {
                showSuccess("Login successful! Your password was exposed — please change it now.");
                setTimeout(() => { window.location.href = "/profile?breached=1"; }, 2000);
            } else {
                showSuccess("Login successful! Redirecting...");
                setTimeout(() => { window.location.href = "/"; }, 1500);
            }
        })
        .catch(err => {
            showError(err.message);
        });
});