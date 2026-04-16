const nav = document.getElementById("nav-bar");
const navRegister = document.getElementById("nav-register");
const navLogin = document.getElementById("nav-login");

const isProfilePage = window.location.pathname.includes("profile");

export function checkIfLoggedIn() {
    const token = localStorage.getItem('token');

    if (!token) {
        // Unauthenticated user trying to access profile — send to login
        if (isProfilePage) {
            window.location.href = "/login";
        }
        return;
    }

    // 1. Hide Login and Register
    if (navLogin) navLogin.style.display = "none";
    if (navRegister) navRegister.style.display = "none";

    // 2. Add Profile link
    const profileLink = document.createElement("a");
    profileLink.href = "/profile";
    profileLink.id = "nav-profile";
    profileLink.textContent = "Profile";
    nav.appendChild(profileLink);

    // 3. Add Logout link
    const logoutLink = document.createElement("a");
    logoutLink.href = "#";
    logoutLink.id = "nav-logout";
    logoutLink.textContent = "Logout";
    logoutLink.addEventListener("click", (e) => {
        e.preventDefault();
        logout();
    });
    nav.appendChild(logoutLink);

    // 4. Breach lockdown — breached users can only use the profile page
    if (sessionStorage.getItem('breachWarning') === '1' && !isProfilePage) {
        window.location.href = "/profile";
    }
}

function logout() {
    localStorage.removeItem('token');
    sessionStorage.removeItem('breachWarning');
    window.location.href = "/";
}

const TOAST_DURATION = 4000;

function getToastContainer() {
    let container = document.getElementById('toast-container');
    if (!container) {
        container = document.createElement('div');
        container.id = 'toast-container';
        document.body.appendChild(container);
    }
    return container;
}

function showToast(message, type) {
    const container = getToastContainer();

    const toast = document.createElement('div');
    toast.className = `toast toast--${type}`;

    const icon = document.createElement('span');
    icon.className = 'toast__icon';
    icon.textContent = type === 'success' ? '✓' : '✕';

    const text = document.createElement('span');
    text.className = 'toast__message';
    text.textContent = message;

    const close = document.createElement('button');
    close.className = 'toast__close';
    close.setAttribute('aria-label', 'Dismiss');
    close.textContent = '×';

    const progress = document.createElement('div');
    progress.className = 'toast__progress';

    toast.appendChild(icon);
    toast.appendChild(text);
    toast.appendChild(close);
    toast.appendChild(progress);
    container.appendChild(toast);

    const dismiss = () => {
        toast.classList.add('toast--out');
        toast.addEventListener('animationend', () => toast.remove(), { once: true });
    };

    const timer = setTimeout(dismiss, TOAST_DURATION);
    close.addEventListener('click', () => { clearTimeout(timer); dismiss(); });
}

export function showError(message) { showToast(message, 'error'); }
export function showSuccess(message) { showToast(message, 'success'); }


