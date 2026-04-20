export function callSearchRestApi(query) {
    return fetch(`/api/search?q=${encodeURIComponent(query)}`)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.json();
        })
        .catch(error => {
            console.error("Search API error:", error);
            throw error;
        });
}


export function callLoginRestApi(username, password) {
    return fetch("/api/login", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            username: username,
            password: password
        })
    })
        .then(response => {
            if (!response.ok) {
                return response.json().then(errData => {
                    throw new Error(errData.message || `Login failed (${response.status})`);
                });
            }
            return response.json();
        })
        .catch(error => {
            console.error("Login error:", error);
            throw error;
        });
}

export function callChangePasswordApi(currentPassword, newPassword, newPassword2) {
    const token = localStorage.getItem("token");
    return fetch("/api/change-password", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`,
        },
        body: JSON.stringify({
            current_password: currentPassword,
            new_password: newPassword,
            new_password2: newPassword2,
        }),
    })
        .then(response => {
            if (!response.ok) {
                return response.json().then(errData => {
                    throw new Error(errData.message || `Request failed (${response.status})`);
                });
            }
            return response.json();
        });
}

export function callRegisterRestApi(userData) {
    return fetch(`/api/register`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(userData)
    })
        .then(response => {
            if (!response.ok) {
                return response.json().then(errData => {
                    throw new Error(errData.message || `Registration failed (${response.status})`);
                });
            }
            return response.json();
        });
}