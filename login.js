document.getElementById('login-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    const response = await fetch('http://localhost:8080/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
    });

    if (response.ok) {
        const userData = await response.json();
        localStorage.setItem('user', JSON.stringify(userData)); // Save user data in local storage
        window.location.href = 'index.html'; // Redirect to home page
    } else {
        alert('Login failed. Please check your credentials and try again.');
    }
});
