import { callChangePasswordApi } from "./api_calls.js";
import { checkIfLoggedIn, showError, showSuccess } from "./reuseable_functions.js";

checkIfLoggedIn();

// Show breach warning if redirected from login with ?breached=1
const params = new URLSearchParams(window.location.search);
if (params.get("breached") === "1") {
    document.getElementById("breach-warning").style.display = "block";
}

const form = document.getElementById("change-password-form");

form.addEventListener("submit", (e) => {
    e.preventDefault();

    const currentPassword = document.getElementById("current-password").value;
    const newPassword = document.getElementById("new-password").value;
    const newPasswordConfirm = document.getElementById("new-password-confirm").value;

    if (newPassword !== newPasswordConfirm) {
        showError("New passwords do not match.");
        return;
    }

    callChangePasswordApi(currentPassword, newPassword, newPasswordConfirm)
        .then(() => {
            document.getElementById("breach-warning").style.display = "none";
            showSuccess("Password updated successfully!");
            form.reset();
        })
        .catch((err) => {
            showError(err.message);
        });
});
