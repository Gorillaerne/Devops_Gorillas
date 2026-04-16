import { callChangePasswordApi } from "./api_calls.js";
import { checkIfLoggedIn, showError, showSuccess } from "./reuseable_functions.js";

checkIfLoggedIn();

// Show breach warning if set by login redirect (sessionStorage) or ?breached=1 query param
const params = new URLSearchParams(window.location.search);
if (sessionStorage.getItem("breachWarning") === "1" || params.get("breached") === "1") {
    document.getElementById("breach-warning").style.display = "block";
    sessionStorage.removeItem("breachWarning");
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
