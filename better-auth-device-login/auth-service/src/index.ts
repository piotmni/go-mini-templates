import { Hono } from "hono";
import { cors } from "hono/cors";
import { auth } from "./auth";

const app = new Hono();

app.use(
  "/api/auth/*",
  cors({
    origin: "http://localhost:3000",
    allowHeaders: ["Content-Type", "Authorization"],
    allowMethods: ["POST", "GET", "OPTIONS"],
    credentials: true,
  })
);

app.on(["POST", "GET"], "/api/auth/*", (c) => {
  return auth.handler(c.req.raw);
});

app.get("/", (c) => {
  return c.json({ message: "Auth service is running" });
});

app.get("/device", (c) => {
  const userCode = c.req.query("user_code") || "";
  const html = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Device Authorization</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f5f5f5;
    }
    .container {
      background: white;
      padding: 2rem;
      border-radius: 8px;
      box-shadow: 0 2px 10px rgba(0,0,0,0.1);
      text-align: center;
      max-width: 400px;
      width: 90%;
    }
    h1 { margin-bottom: 0.5rem; color: #333; }
    .subtitle { color: #666; margin-bottom: 1.5rem; }
    .code-input {
      width: 100%;
      padding: 1rem;
      font-size: 1.5rem;
      text-align: center;
      border: 2px solid #ddd;
      border-radius: 8px;
      margin-bottom: 1rem;
      text-transform: uppercase;
      letter-spacing: 0.2em;
    }
    .code-input:focus {
      outline: none;
      border-color: #24292e;
    }
    .submit-btn {
      width: 100%;
      padding: 0.75rem 1.5rem;
      background: #24292e;
      color: white;
      border: none;
      border-radius: 6px;
      font-size: 1rem;
      cursor: pointer;
      transition: background 0.2s;
    }
    .submit-btn:hover { background: #1a1e22; }
    .submit-btn:disabled { background: #ccc; cursor: not-allowed; }
    .message { margin-top: 1rem; padding: 0.75rem; border-radius: 6px; }
    .message.success { background: #d4edda; color: #155724; }
    .message.error { background: #f8d7da; color: #721c24; }
    .message.info { background: #d1ecf1; color: #0c5460; }
    .login-prompt { margin-top: 1.5rem; padding-top: 1.5rem; border-top: 1px solid #eee; }
    .login-btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.5rem 1rem;
      background: #24292e;
      color: white;
      border: none;
      border-radius: 6px;
      font-size: 0.9rem;
      cursor: pointer;
      text-decoration: none;
    }
    .login-btn:hover { background: #1a1e22; }
    .login-btn svg { width: 18px; height: 18px; fill: white; }
  </style>
</head>
<body>
  <div class="container">
    <h1>Device Authorization</h1>
    <p class="subtitle">Enter the code shown on your device</p>
    <div id="session-status"></div>
    <form id="device-form">
      <input 
        type="text" 
        id="user-code" 
        class="code-input" 
        placeholder="ABCD-1234"
        value="${userCode}"
        maxlength="9"
        autocomplete="off"
        required
      >
      <button type="submit" class="submit-btn" id="submit-btn">Authorize Device</button>
    </form>
    <div id="message"></div>
    <div id="login-prompt" class="login-prompt" style="display: none;">
      <p style="margin-bottom: 0.75rem; color: #666;">You need to sign in first</p>
      <button class="login-btn" onclick="loginWithGithub()">
        <svg viewBox="0 0 24 24"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>
        Sign in with GitHub
      </button>
    </div>
  </div>
  <script type="module">
    import { createAuthClient } from "https://esm.sh/better-auth@1.4.10/client";
    import { deviceAuthorizationClient } from "https://esm.sh/better-auth@1.4.10/client/plugins";
    
    const authClient = createAuthClient({
      baseURL: "http://localhost:3000",
      plugins: [deviceAuthorizationClient()]
    });

    // Check if user is logged in
    async function checkSession() {
      try {
        const session = await authClient.getSession();
        if (session?.data?.user) {
          document.getElementById('session-status').innerHTML = 
            '<div class="message info">Signed in as ' + session.data.user.email + '</div>';
          document.getElementById('login-prompt').style.display = 'none';
          document.getElementById('submit-btn').disabled = false;
          return true;
        }
      } catch (e) {
        console.log('Not signed in');
      }
      document.getElementById('login-prompt').style.display = 'block';
      document.getElementById('submit-btn').disabled = true;
      return false;
    }

    window.loginWithGithub = async () => {
      const currentUrl = window.location.href;
      await authClient.signIn.social({
        provider: "github",
        callbackURL: currentUrl
      });
    };

    document.getElementById('device-form').addEventListener('submit', async (e) => {
      e.preventDefault();
      const userCode = document.getElementById('user-code').value.trim().replace(/-/g, '').toUpperCase();
      const messageEl = document.getElementById('message');
      const submitBtn = document.getElementById('submit-btn');
      
      if (!userCode) {
        messageEl.innerHTML = '<div class="message error">Please enter a code</div>';
        return;
      }

      submitBtn.disabled = true;
      submitBtn.textContent = 'Verifying...';

      try {
        // Approve the device
        const result = await authClient.device.approve({
          userCode: userCode
        });
        
        console.log('Approve result:', result);
        
        if (result.error) {
          throw new Error(result.error.message || 'Authorization failed');
        }
        
        messageEl.innerHTML = '<div class="message success">Device authorized! You can close this window and return to your CLI.</div>';
        submitBtn.textContent = 'Authorized';
      } catch (error) {
        console.error('Approve error:', error);
        messageEl.innerHTML = '<div class="message error">' + (error.message || 'Authorization failed. Please try again.') + '</div>';
        submitBtn.disabled = false;
        submitBtn.textContent = 'Authorize Device';
      }
    });

    // Format code input as user types
    document.getElementById('user-code').addEventListener('input', (e) => {
      let value = e.target.value.toUpperCase().replace(/[^A-Z0-9]/g, '');
      if (value.length > 4) {
        value = value.slice(0, 4) + '-' + value.slice(4, 8);
      }
      e.target.value = value;
    });

    // Check session on load
    checkSession();
  </script>
</body>
</html>
  `;
  return c.html(html);
});

app.get("/login", (c) => {
  const html = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Login</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f5f5f5;
    }
    .container {
      background: white;
      padding: 2rem;
      border-radius: 8px;
      box-shadow: 0 2px 10px rgba(0,0,0,0.1);
      text-align: center;
    }
    h1 { margin-bottom: 1.5rem; color: #333; }
    .github-btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.75rem 1.5rem;
      background: #24292e;
      color: white;
      border: none;
      border-radius: 6px;
      font-size: 1rem;
      cursor: pointer;
      transition: background 0.2s;
    }
    .github-btn:hover { background: #1a1e22; }
    .github-btn svg { width: 20px; height: 20px; fill: white; }
  </style>
</head>
<body>
  <div class="container">
    <h1>Welcome</h1>
    <button class="github-btn" onclick="loginWithGithub()">
      <svg viewBox="0 0 24 24"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>
      Sign in with GitHub
    </button>
  </div>
  <script type="module">
    import { createAuthClient } from "https://esm.sh/better-auth@1.2.5/client";
    
    const authClient = createAuthClient({
      baseURL: "http://localhost:3000"
    });

    window.loginWithGithub = async () => {
      await authClient.signIn.social({
        provider: "github",
        callbackURL: "/"
      });
    };
  </script>
</body>
</html>
  `;
  return c.html(html);
});

export default {
  port: 3000,
  fetch: app.fetch,
};
