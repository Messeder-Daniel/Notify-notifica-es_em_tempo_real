import './style.css'

document.querySelector<HTMLDivElement>('#app')!.innerHTML = `
  <main class="app">
    <section class="hero">
      <p class="eyebrow">Sistema de notificações em tempo real</p>
      <h1>Real-Time Notifications</h1>
      <p class="description">
        Frontend em Vite + TypeScript para consumir uma API em Go,
        autenticar usuários com JWT e receber notificações via WebSocket.
      </p>
    </section>

    <section class="card">
      <h2>Status do frontend</h2>
      <ul class="status-list">
        <li><span class="status-dot success"></span> Vite configurado</li>
        <li><span class="status-dot success"></span> TypeScript ativo</li>
        <li><span class="status-dot pending"></span> Login será integrado na Etapa 14</li>
        <li><span class="status-dot pending"></span> WebSocket será integrado na Etapa 14</li>
      </ul>
    </section>

    <section class="grid">
      <article class="panel">
        <h3>Autenticação</h3>
        <p>
          O frontend usará o endpoint <code>POST /auth/login</code>
          para obter o token JWT.
        </p>
      </article>

      <article class="panel">
        <h3>Notificações</h3>
        <p>
          As notificações serão carregadas por <code>GET /notifications</code>
          e atualizadas em tempo real via WebSocket.
        </p>
      </article>

      <article class="panel">
        <h3>Tempo real</h3>
        <p>
          A conexão WebSocket usará <code>/ws?token=&lt;jwt&gt;</code>
          para receber eventos <code>notification.created</code>.
        </p>
      </article>
    </section>
  </main>
`
