import { callLoginRestApi } from "./api_calls.js";
import { checkIfLoggedIn, createErrorElement } from "./reuseable_functions.js";

checkIfLoggedIn()

const loginForm = document.getElementById('login-form');
const body = document.getElementById("body")


loginForm.addEventListener('submit', function(e) {
    e.preventDefault(); // Stop the browser from reloading the page

    // Grab the values using the IDs we just added
    const userVal = document.getElementById('username').value;
    const passVal = document.getElementById('password').value;

    // Call your function (from the previous step)
    callLoginRestApi(userVal, passVal)
        .then(data => {
            localStorage.setItem("token",data.token)

            alert("Login Successful!");
            window.location.href = "/"; // Redirect manually after success
        })
        .catch(err => {
            body.prepend(createErrorElement(err.message))
        });
});