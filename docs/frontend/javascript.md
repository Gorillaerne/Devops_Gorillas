# Frontend JavaScript Modules

**Directory:** `Go_Refined_Code/static/javaScript/`

The frontend uses six ES modules. Each file is focused on a single responsibility. Modules share functionality through `import`/`export` — there is no bundler, so they run natively in the browser.

---

## api_calls.js — API Client

The central place for all `fetch` calls to the backend. Each function returns a Promise.

### `callSearchRestApi(query)`

- **Calls:** `GET /api/search?q=<query>`
- **Returns:** The parsed JSON response `{ data: [...] }`.
- **Throws:** An error if the response status is not OK.

### `callLoginRestApi(username, password)`

- **Calls:** `POST /api/login` with `{ username, password }` as JSON.
- **Returns:** The parsed JSON response including `token` and `breached`.
- **Throws:** An error with the server's error message if the login fails.

### `callRegisterRestApi(userData)`

- **Calls:** `POST /api/register` with the user data as JSON.
- **Returns:** The parsed JSON response.
- **Throws:** An error with the server's error message on failure.

### `callChangePasswordApi(currentPassword, newPassword, newPassword2)`

- **Calls:** `POST /api/change-password`.
- **Reads:** The JWT token from `localStorage` and sends it as `Authorization: Bearer <token>`.
- **Returns:** The parsed JSON response.
- **Throws:** An error with the server's error message on failure.

---

## reuseable_functions.js — Shared Utilities

Imported by every page. Handles authentication state, navigation updates, breach lockdown, and toast notifications.

### `checkIfLoggedIn()`

Called on page load by each page script. Checks for a JWT token in `localStorage`.

**If no token:**
- On the profile page: redirects to `/login`.
- On other pages: does nothing (guest access is allowed).

**If a token exists:**
- Hides the Login and Register nav links.
- Dynamically adds Profile and Logout links to the nav bar.
- Checks `sessionStorage` for `"breachWarning=1"`. If found and the user is not on the profile page, redirects to `/profile`. This enforces the breach lockdown.

**`logout()` (private):**
- Removes the JWT token from `localStorage`.
- Removes the breach warning flag from `sessionStorage`.
- Redirects to the homepage.

### Toast Notification System

A lightweight, non-blocking notification system. Toasts appear in the top-right corner and auto-dismiss after 4 seconds.

**`showError(message)`** — Shows a red error toast.

**`showSuccess(message)`** — Shows a green success toast.

**Internal functions:**
- `getToastContainer()` — Creates or finds the toast container `<div>`.
- `showToast(message, type)` — Creates the toast element, appends it to the container, and sets up the dismiss timer and close button.

---

## search_script.js — Homepage Search

Handles the search form on `index.html`.

**Behaviour:**
- Listens for a button click or Enter key press.
- Calls `callSearchRestApi` with the query from the input field.
- Clears the results container and renders each result as a card with a title, URL link, and content snippet.
- Shows an error toast if the API call fails.
- Calls `checkIfLoggedIn` on page load.

---

## login_page_script.js — Login Form

Handles the login form on `login.html`.

**Behaviour:**
- Listens for the form submit event.
- Calls `callLoginRestApi` with the username and password.
- On success:
  - Saves the JWT token to `localStorage` under `"token"`.
  - If `breached: true` in the response, sets `"breachWarning=1"` in `sessionStorage`.
  - Redirects to the homepage.
- On failure, shows an error toast.
- Calls `checkIfLoggedIn` on page load (redirects already-logged-in users away).

---

## register_page_script.js — Registration Form

Handles the registration form on `register.html`.

**Behaviour:**
- Listens for the form submit event.
- Validates that the two password fields match client-side before making any API call.
- Calls `callRegisterRestApi` with `{ username, email, password, password2 }`.
- On success, redirects to `/login`.
- On failure, shows an error toast.
- Calls `checkIfLoggedIn` on page load.

---

## profile_script.js — Profile / Change Password

Handles the change-password form on `profile.html`.

**Behaviour:**
- On page load:
  - Calls `checkIfLoggedIn` (redirects to `/login` if not authenticated).
  - Checks `sessionStorage` for `"breachWarning=1"` and shows the breach warning banner if found.
- Listens for the form submit event.
- Calls `callChangePasswordApi` with the three password fields.
- On success:
  - Clears `"breachWarning"` from `sessionStorage` (the user has complied).
  - Shows a success toast.
- On failure, shows an error toast.
