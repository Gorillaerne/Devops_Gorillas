

import { callRegisterRestApi } from "./api_calls.js";
import { checkIfLoggedIn, showError, showSuccess } from "./reuseable_functions.js";

const registerForm = document.getElementById('register-form');

checkIfLoggedIn()

if (registerForm) {
    registerForm.addEventListener('submit', (e) => {
        e.preventDefault();

        const userData = {
            username: document.getElementById('reg-username').value,
            email: document.getElementById('reg-email').value,
            password: document.getElementById('reg-password').value,
            password2: document.getElementById('reg-password-confirm').value
        };

        if (userData.password !== userData.password2) {
            showError("Passwords don't match");
            return;
        }

        callRegisterRestApi(userData)
            .then(_data => {
                showSuccess("Account created! Redirecting to login...");
                setTimeout(() => { window.location.href = "/login"; }, 1500);
            })
            .catch(err => {
                showError(err.message);
            });
    });
}