

import { callRegisterRestApi } from "./api_calls.js";
import { checkIfLoggedIn, createErrorElement, createSuccessElement } from "./reuseable_functions.js";

const registerForm = document.getElementById('register-form');
const body = document.getElementById("body")


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
            body.prepend(createErrorElement("Passwords don't match"));
            return;
        }

        callRegisterRestApi(userData)
            .then(_data => {
                body.prepend(createSuccessElement("Account created! Redirecting to login..."));
                setTimeout(() => { window.location.href = "/login"; }, 1500);
            })
            .catch(err => {
                body.prepend(createErrorElement(err.message));
            });
    });
}