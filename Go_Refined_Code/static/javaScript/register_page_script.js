

import { callRegisterRestApi } from "./api_calls.js";
import { checkIfLoggedIn, createErrorElement } from "./reuseable_functions.js";

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
           return alert("Passwords dont match");
        }

        callRegisterRestApi(userData)
            .then(_data => {
                alert("Account created! Please log in.");
                window.location.href = "/login";
            })
            .catch(err => {
                body.prepend(createErrorElement(err.message))
            });
    });
}