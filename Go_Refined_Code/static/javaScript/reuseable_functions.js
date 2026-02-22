const nav = document.getElementById("nav-bar");
const navRegister = document.getElementById("nav-register");
const navLogin = document.getElementById("nav-login");

export function checkIfLoggedIn() {

    const token = localStorage.getItem('token');

    if (token) {
        // 1. Hide Login and Register
        if (navLogin) navLogin.style.display = "none";
        if (navRegister) navRegister.style.display = "none";

        // 2. Create and Append Logout link
        const logoutLink = document.createElement("a");
        logoutLink.href = "#";
        logoutLink.id = "nav-logout";
        logoutLink.textContent = "Logout";

        logoutLink.addEventListener("click", (e) => {
            e.preventDefault();
            logout();

        });

        nav.appendChild(logoutLink);
    }
}

function logout() {
    // Simply remove the token and refresh
    localStorage.removeItem('token');
    window.location.href = "/";
}

export function createErrorElement(message) {
    const div = document.createElement('div');
    div.className = 'error';
    const strong = document.createElement('strong');
    strong.textContent = message; 
    div.appendChild(strong);
    
    return div;
}


