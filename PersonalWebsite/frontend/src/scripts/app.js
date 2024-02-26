const loginForm = document.getElementById('loginForm');
const errorMessage = document.querySelector('.error-message');

loginForm.addEventListener('submit', (event) => {
    event.preventDefault(); // Prevent default form submission

    const password = document.querySelector('input[name="password"]').value;

    // Send data to backend (using Fetch API here, adjust if needed)
    fetch('/login', {  
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded' 
        },
        body: new URLSearchParams({
            'password': password 
        })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Login failed'); 
        }
        return response.text(); // Or response.json() if your backend sends JSON
    })
    .then(data => {
        // Successful login
        window.location.href = data;  // Redirect to the received URL 
    })
    .catch(error => {
        errorMessage.textContent = 'Invalid password. Please try again.'; 
    });
});