const switchToRegister = document.getElementById('switch-to-register');
const switchToLogin = document.getElementById('switch-to-login');
const loginContainer = document.getElementById('login-container');
const registerContainer = document.getElementById('register-container');

switchToRegister.addEventListener('click', () => {
  loginContainer.style.display = 'none';
  registerContainer.style.display = 'block';
});

switchToLogin.addEventListener('click', () => {
  registerContainer.style.display = 'none';
  loginContainer.style.display = 'block';
});

document.getElementById('login-form').addEventListener('submit', async function(event) {
  event.preventDefault();
  const username = document.getElementById('login-username').value;
  const password = document.getElementById('login-password').value;

  const response = await fetch('/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ username, password })
  });

  const result = await response.json();

  if (result.success) {
    window.location.href = '/welcome'; // Redirect to welcome page on success
  } else {
    alert('Login failed: ' + result.message);
  }
});

document.getElementById('register-form').addEventListener('submit', async function(event) {
  event.preventDefault();
  const username = document.getElementById('register-username').value;
  const password = document.getElementById('register-password').value;

  const response = await fetch('/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ username, password })
  });

  const result = await response.json();

  if (result.success) {
    alert('Registration successful! Please log in.');
    registerContainer.style.display = 'none';
    loginContainer.style.display = 'block';
  } else {
    alert('Registration failed: ' + result.message);
  }
});
