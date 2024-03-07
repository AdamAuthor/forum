document.addEventListener("DOMContentLoaded", function () {
    const wrapper = document.querySelector('.wrapper');
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    const errorMessage = document.getElementById('error-message');
    const successMessage = document.getElementById('success-message');
    const btnOut = document.querySelector('.btnLogout-popup');
    const btnPopup = document.querySelector('.btnLogin-popup');

    loginForm.addEventListener('submit', function (event) {
        event.preventDefault();
        const formData = new FormData(loginForm);
        const signInData = {
            username: formData.get('login-username'),
            password: formData.get('login-password')
        };

        fetch('http://localhost:8080/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(signInData)
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Login failed');
                }
                return response.text();
            })
            .then(data => {
                wrapper.classList.remove('active-popup');
                successMessage.textContent = 'Login successful!';
                successMessage.style.display = 'block';
                btnPopup.style.display = 'none'
                btnOut.style.display = 'inline-block'
                setTimeout(function () {
                    successMessage.style.display = 'none';
                }, 2000);

            })
            .catch(error => {
                errorMessage.textContent = 'Invalid username or password. Please try again.';
                errorMessage.style.display = 'block';
                setTimeout(function () {
                    errorMessage.style.display = 'none';
                }, 5000);
            });
    });


    registerForm.addEventListener('submit', function (event) {
        event.preventDefault();
        const formData = new FormData(registerForm);
        const registerData = {
            username: formData.get('username'),
            email: formData.get('email'),
            password_hash: formData.get('password')
        };

        fetch('http://localhost:8080/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(registerData)
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Registration failed');
                }
                return response.text();
            })
            .then(data => {
                wrapper.classList.remove('active');
                successMessage.textContent = 'Registration successful!';
                successMessage.style.display = 'block';
                setTimeout(function () {
                    successMessage.style.display = 'none';
                }, 5000);
            })
            .catch(error => {
                errorMessage.textContent = 'Username or email has already use. Please try again.';
                errorMessage.style.display = 'block';
                setTimeout(function () {
                    errorMessage.style.display = 'none';
                }, 5000);
            });
    });

    btnOut.addEventListener('click', () => {
        fetch('http://localhost:8080/logout', {
            method: 'GET',
            credentials: 'include'
        })
            .then(response => {
                if (response.ok) {
                    document.cookie = 'session_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
                    btnPopup.style.display = 'inline-block';
                    btnOut.style.display = 'none';
                    window.location.reload();
                } else {
                    console.error('Failed to logout');
                }
            })
            .catch(error => {
                console.error('Error during logout:', error);
            });
    });
});