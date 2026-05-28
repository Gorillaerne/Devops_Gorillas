# Frontend Pages

**Directory:** `Go_Refined_Code/static/html/`

The frontend is five static HTML pages served by Nginx. Each page is a self-contained document that loads the shared stylesheet and one or two JavaScript modules.

All pages share the same navigation bar (`<nav id="nav-bar">`) and the same `style.css`. The `reuseable_functions.js` module dynamically updates the nav bar to show Profile/Logout links when a user is logged in.

---

## index.html — Search Homepage

**URL:** `/`

The main entry point of the application.

**What it contains:**
- A search input field and a search button.
- A results container (`<div id="results">`) where search results are rendered.
- Navigation links to Login and Register.

**Scripts loaded:**
- `search_script.js` — handles search form submission and renders results.
- `reuseable_functions.js` — manages the nav bar and toast notifications.

**Behaviour:**
- Search can be triggered by clicking the button or pressing Enter.
- Results appear below the search bar without a page reload.
- If the user is logged in and breached, they are immediately redirected to `/profile`.

---

## login.html — Login Page

**URL:** `/login`

A form for existing users to log in.

**What it contains:**
- Username and password input fields.
- A submit button.
- A link to the registration page for new users.

**Scripts loaded:**
- `login_page_script.js` — handles form submission and stores the JWT token.
- `reuseable_functions.js`

**Behaviour:**
- On successful login, the JWT token is saved to `localStorage` under the key `"token"`.
- If the server returns `"breached": true`, the flag `"breachWarning=1"` is saved to `sessionStorage` and the user is redirected to `/profile`.
- On failure, an error toast notification appears.

---

## register.html — Registration Page

**URL:** `/register`

A form for new users to create an account.

**What it contains:**
- Username, email, password, and confirm-password fields.
- A submit button.
- A link back to the login page.

**Scripts loaded:**
- `register_page_script.js` — handles form submission and client-side validation.
- `reuseable_functions.js`

**Behaviour:**
- Client-side: checks that both password fields match before sending the request.
- Server-side: the backend also validates this (the client-side check is just for a faster user experience).
- On success, the user is redirected to the login page.
- On failure, an error toast notification appears.

---

## profile.html — User Profile / Change Password

**URL:** `/profile`

A protected page for authenticated users. Allows them to change their password.

**What it contains:**
- A breach warning banner (hidden by default, shown when `breachWarning=1` is in sessionStorage).
- Current password, new password, and confirm new password fields.
- A submit button.

**Scripts loaded:**
- `profile_script.js` — handles the change-password form and reads the breach flag.
- `reuseable_functions.js`

**Behaviour:**
- If no JWT token exists in `localStorage`, the user is immediately redirected to `/login`.
- If the user is flagged as breached (via sessionStorage), the breach warning banner is shown and they cannot navigate away until they change their password.
- On success, the breach flag is cleared from sessionStorage.
- On failure, an error toast appears.

---

## about.html — About Page

**URL:** `/about`

A static informational page about the team.

**What it contains:**
- A brief mission statement about the ¿Who Knows? project.
- A team photo (`teamphoto.png`).

**Scripts loaded:**
- `reuseable_functions.js` (for the nav bar)

**Behaviour:** Static content only — no API calls.
