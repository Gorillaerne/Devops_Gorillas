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


