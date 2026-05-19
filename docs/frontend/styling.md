# Styling — style.css

**File:** `Go_Refined_Code/static/style.css`

A single global stylesheet (~668 lines) shared by all five HTML pages. There is no CSS framework — everything is written in plain CSS.

---

## Structure

### CSS Variables (Custom Properties)

Colours, spacing, and font sizes are defined as CSS custom properties at the `:root` level so they can be reused and changed in one place.

### Layout

- Pages use a centered column layout with max-width constraints to stay readable on large screens.
- The navigation bar is a horizontal flex container that stays at the top of every page.

### Navigation Bar

The nav bar (`#nav-bar`) is styled to show links horizontally. The JavaScript in `reuseable_functions.js` dynamically adds or removes links from it depending on authentication state. The CSS ensures new links inherit the same appearance without needing extra classes.

### Forms

Login, registration, and profile forms share consistent styles:
- Stacked input labels and fields.
- Full-width inputs on smaller screens.
- A prominent submit button.

### Search Results

Each search result is displayed as a card with:
- A title (linked to the original URL).
- A content snippet.
- Hover effects for interactivity feedback.

### Toast Notifications

The toast system has its own CSS section:
- `#toast-container` — fixed position in the top-right corner, stacks multiple toasts vertically.
- `.toast` — base toast style with padding, shadow, and border-radius.
- `.toast--success` — green left border and icon.
- `.toast--error` — red left border and icon.
- `.toast--out` — animation class added when the toast is dismissed; triggers a fade-out/slide-out CSS animation.
- `.toast__progress` — a thin progress bar at the bottom of each toast that depletes over 4 seconds, giving a visual countdown.

### Breach Warning Banner

The breach warning banner on `profile.html` is styled as a prominent red alert box. It is `display: none` by default and shown by JavaScript when the breach flag is detected in `sessionStorage`.

### Responsive Design

Media queries adjust layout at common breakpoints to ensure the application is usable on mobile devices. Inputs become full-width and the nav bar collapses appropriately on small screens.
