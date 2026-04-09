import { callLoginRestApi } from "./api_calls.js";
import { checkIfLoggedIn, createErrorElement, createSuccessElement } from "./reuseable_functions.js";

checkIfLoggedIn()

const loginForm = document.getElementById('login-form');
const body = document.getElementById("body")


loginForm.addEventListener('submit', (e) => {
    e.preventDefault(); // Stop the browser from reloading the page

    // Grab the values using the IDs we just added
    const userVal = document.getElementById('username').value;
    const passVal = document.getElementById('password').value;

    // Call your function (from the previous step)
    callLoginRestApi(userVal, passVal)
        .then(data => {
            localStorage.setItem("token", data.token);
            body.prepend(createSuccessElement("Login successful! Redirecting..."));
            setTimeout(() => { window.location.href = "/"; }, 1500);
        })
        .catch(err => {
            body.prepend(createErrorElement(err.message));
        });
});