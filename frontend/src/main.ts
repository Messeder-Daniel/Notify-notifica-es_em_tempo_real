import './style.css'
import type { LoginResponse, Notification, User, WebSocketEvent } from './types'

const API_URL = 'http://localhost:8080'
const WS_URL = 'ws://localhost:8080'

type ConnectionStatus = 'disconnected' | 'connecting' | 'connected'

type AppState = {
  token: string | null
  user: User | null
  notifications: Notification[]
  connectionStatus: ConnectionStatus
  errorMessage: string | null
  successMessage: string | null
  isLoading: boolean
  socket: WebSocket | null
}

const state: AppState = {
  token: localStorage.getItem('token'),
  user: loadUserFromStorage(),
  notifications: [],
  connectionStatus: 'disconnected',
  errorMessage: null,
  successMessage: null,
  isLoading: false,
  socket: null,
}

function loadUserFromStorage(): User | null {
  const storedUser = localStorage.getItem('user')

  if (!storedUser) {
    return null
  }

  try {
    return JSON.parse(storedUser) as User
  } catch {
    localStorage.removeItem('user')
    return null
  }
}

function saveSession(token: string, user: User): void {
  state.token = token
  state.user = user

  localStorage.setItem('token', token)
  localStorage.setItem('user', JSON.stringify(user))
}

function clearSession(): void {
  state.token = null
  state.user = null
  state.notifications = []
  state.connectionStatus = 'disconnected'

  if (state.socket) {
    state.socket.close()
    state.socket = null
  }

  localStorage.removeItem('token')
  localStorage.removeItem('user')
}

async function apiRequest<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers = new Headers(options.headers)

  if (options.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json')
  }

  if (state.token) {
    headers.set('Authorization', `Bearer ${state.token}`)
  }

  const response = await fetch(`${API_URL}${path}`, {
    ...options,
    headers,
  })

  if (response.status === 401) {
    clearSession()
    render()
    throw new Error('Sessão expirada. Faça login novamente.')
  }

  if (!response.ok) {
    let message = 'Erro inesperado na comunicação com a API.'

    try {
      const body = (await response.json()) as { error?: string }
      message = body.error ?? message
    } catch {
      // Mantém a mensagem padrão.
    }

    throw new Error(message)
  }

  return response.json() as Promise<T>
}

function render(): void {
  const app = document.querySelector<HTMLDivElement>('#app')

  if (!app) {
    return
  }

  app.innerHTML = state.token ? renderDashboard() : renderLogin()

  if (state.token) {
    bindDashboardEvents()
  } else {
    bindLoginEvents()
  }
}

function renderLogin(): string {
  return `
    <main class="auth-page">
      <section class="auth-card">
        <p class="eyebrow">Sistema de notificações em tempo real</p>
        <h1>Entrar</h1>
        <p class="description">
          Faça login para acessar o painel de notificações e receber eventos em tempo real via WebSocket.
        </p>

        ${renderFeedback()}

        <form id="login-form" class="form">
          <label>
            E-mail
            <input
              type="email"
              name="email"
              value="daniel@example.com"
              autocomplete="email"
              required
            />
          </label>

          <label>
            Senha
            <input
              type="password"
              name="password"
              value="password"
              autocomplete="current-password"
              required
            />
          </label>

          <button type="submit" ${state.isLoading ? 'disabled' : ''}>
            ${state.isLoading ? 'Entrando...' : 'Entrar no sistema'}
          </button>
        </form>

        <p class="hint">
          Usuário de teste: <strong>daniel@example.com</strong> / <strong>password</strong>
        </p>
      </section>
    </main>
  `
}

function renderDashboard(): string {
  const unreadCount = state.notifications.filter((notification) => !notification.is_read).length
  const totalCount = state.notifications.length

  return `
    <main class="dashboard">
      <header class="dashboard-header">
        <div>
          <p class="eyebrow">Dashboard</p>
          <h1>Real-Time Notifications</h1>
          <p class="description">
            Painel do usuário autenticado com integração HTTP + WebSocket.
          </p>
        </div>

        <div class="user-box">
          <span>${escapeHtml(state.user?.email ?? 'Usuário autenticado')}</span>
          <button id="logout-button" class="secondary-button" type="button">Sair</button>
        </div>
      </header>

      ${renderFeedback()}

      <section class="summary-grid">
        <article class="summary-card">
          <span class="summary-label">WebSocket</span>
          <strong class="connection ${state.connectionStatus}">
            <span class="connection-dot"></span>
            ${getConnectionLabel()}
          </strong>
        </article>

        <article class="summary-card">
          <span class="summary-label">Total</span>
          <strong>${totalCount}</strong>
        </article>

        <article class="summary-card">
          <span class="summary-label">Não lidas</span>
          <strong>${unreadCount}</strong>
        </article>
      </section>

      <section class="content-grid">
        <article class="panel">
          <h2>Criar notificação</h2>
          <p>
            Esta ação chama <code>POST /notifications</code>. Depois que o backend salva no banco,
            ele envia o evento <code>notification.created</code> pelo WebSocket.
          </p>

          <form id="notification-form" class="form">
            <label>
              Título
              <input
                type="text"
                name="title"
                placeholder="Ex: Nova mensagem"
                required
              />
            </label>

            <label>
              Mensagem
              <textarea
                name="message"
                rows="5"
                placeholder="Escreva a mensagem da notificação"
                required
              ></textarea>
            </label>

            <button type="submit" ${state.isLoading ? 'disabled' : ''}>
              ${state.isLoading ? 'Enviando...' : 'Criar notificação'}
            </button>
          </form>
        </article>

        <article class="panel notifications-panel">
          <div class="panel-header">
            <div>
              <h2>Notificações</h2>
              <p>Lista carregada do PostgreSQL e atualizada em tempo real.</p>
            </div>

            <button id="refresh-button" class="secondary-button" type="button">
              Atualizar
            </button>
          </div>

          <div class="notifications-list">
            ${renderNotifications()}
          </div>
        </article>
      </section>
    </main>
  `
}

function renderFeedback(): string {
  if (state.errorMessage) {
    return `<div class="alert error">${escapeHtml(state.errorMessage)}</div>`
  }

  if (state.successMessage) {
    return `<div class="alert success">${escapeHtml(state.successMessage)}</div>`
  }

  return ''
}

function renderNotifications(): string {
  if (state.notifications.length === 0) {
    return `
      <div class="empty-state">
        <strong>Nenhuma notificação ainda.</strong>
        <span>Crie uma notificação para testar o envio em tempo real.</span>
      </div>
    `
  }

  return state.notifications.map(renderNotification).join('')
}

function renderNotification(notification: Notification): string {
  const readClass = notification.is_read ? 'read' : 'unread'
  const readLabel = notification.is_read ? 'Lida' : 'Nova'
  const readAt = notification.read_at ? formatDate(notification.read_at) : null

  return `
    <article class="notification-card ${readClass}">
      <div class="notification-content">
        <div class="notification-topline">
          <span class="badge ${readClass}">${readLabel}</span>
          <time>${formatDate(notification.created_at)}</time>
        </div>

        <h3>${escapeHtml(notification.title)}</h3>
        <p>${escapeHtml(notification.message)}</p>

        ${readAt ? `<small>Lida em ${readAt}</small>` : ''}
      </div>

      <button
        type="button"
        class="secondary-button"
        data-action="mark-as-read"
        data-id="${notification.id}"
        ${notification.is_read ? 'disabled' : ''}
      >
        Marcar como lida
      </button>
    </article>
  `
}

function bindLoginEvents(): void {
  const form = document.querySelector<HTMLFormElement>('#login-form')

  form?.addEventListener('submit', async (event) => {
    event.preventDefault()

    const formData = new FormData(form)
    const email = String(formData.get('email') ?? '').trim()
    const password = String(formData.get('password') ?? '').trim()

    state.isLoading = true
    state.errorMessage = null
    state.successMessage = null
    render()

    try {
      const response = await apiRequest<LoginResponse>('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
      })

      const user: User = response.user ?? {
        id: '',
        name: 'Usuário',
        email,
      }

      saveSession(response.token, user)

      await loadNotifications()
      connectWebSocket()

      state.successMessage = 'Login realizado com sucesso.'
    } catch (error) {
      state.errorMessage = getErrorMessage(error)
    } finally {
      state.isLoading = false
      render()
    }
  })
}

function bindDashboardEvents(): void {
  const logoutButton = document.querySelector<HTMLButtonElement>('#logout-button')
  const refreshButton = document.querySelector<HTMLButtonElement>('#refresh-button')
  const notificationForm = document.querySelector<HTMLFormElement>('#notification-form')
  const markAsReadButtons = document.querySelectorAll<HTMLButtonElement>('[data-action="mark-as-read"]')

  logoutButton?.addEventListener('click', () => {
    clearSession()
    state.successMessage = null
    state.errorMessage = null
    render()
  })

  refreshButton?.addEventListener('click', async () => {
    await runWithFeedback(async () => {
      await loadNotifications()
      state.successMessage = 'Notificações atualizadas.'
    })
  })

  notificationForm?.addEventListener('submit', async (event) => {
    event.preventDefault()

    const formData = new FormData(notificationForm)
    const title = String(formData.get('title') ?? '').trim()
    const message = String(formData.get('message') ?? '').trim()

    await runWithFeedback(async () => {
      const notification = await apiRequest<Notification>('/notifications', {
        method: 'POST',
        body: JSON.stringify({ title, message }),
      })

      upsertNotification(notification)
      notificationForm.reset()
      state.successMessage = 'Notificação criada com sucesso.'
    })
  })

  markAsReadButtons.forEach((button) => {
    button.addEventListener('click', async () => {
      const id = button.dataset.id

      if (!id) {
        return
      }

      await runWithFeedback(async () => {
        const notification = await apiRequest<Notification>(`/notifications/${id}/read`, {
          method: 'PATCH',
        })

        upsertNotification(notification)
        state.successMessage = 'Notificação marcada como lida.'
      })
    })
  })
}

async function runWithFeedback(action: () => Promise<void>): Promise<void> {
  state.isLoading = true
  state.errorMessage = null
  state.successMessage = null
  render()

  try {
    await action()
  } catch (error) {
    state.errorMessage = getErrorMessage(error)
  } finally {
    state.isLoading = false
    render()
  }
}

async function loadNotifications(): Promise<void> {
  state.notifications = await apiRequest<Notification[]>('/notifications')
}

function connectWebSocket(): void {
  if (!state.token) {
    return
  }

  if (state.socket) {
    state.socket.close()
  }

  state.connectionStatus = 'connecting'

  const socket = new WebSocket(`${WS_URL}/ws?token=${encodeURIComponent(state.token)}`)
  state.socket = socket

  socket.addEventListener('open', () => {
    state.connectionStatus = 'connected'
    render()
  })

  socket.addEventListener('message', (event) => {
    try {
      const websocketEvent = JSON.parse(event.data) as WebSocketEvent

      if (websocketEvent.type === 'connected') {
        state.connectionStatus = 'connected'
        render()
        return
      }

      if (websocketEvent.type === 'notification.created') {
        upsertNotification(websocketEvent.data)
        state.successMessage = 'Nova notificação recebida em tempo real.'
        render()
      }
    } catch {
      state.errorMessage = 'Não foi possível processar uma mensagem WebSocket.'
      render()
    }
  })

  socket.addEventListener('error', () => {
    state.errorMessage = 'Erro na conexão WebSocket.'
    render()
  })

  socket.addEventListener('close', () => {
    if (state.socket === socket) {
      state.socket = null
      state.connectionStatus = 'disconnected'
      render()
    }
  })
}

function upsertNotification(notification: Notification): void {
  const index = state.notifications.findIndex((item) => item.id === notification.id)

  if (index >= 0) {
    state.notifications[index] = notification
  } else {
    state.notifications = [notification, ...state.notifications]
  }
}

function getConnectionLabel(): string {
  if (state.connectionStatus === 'connected') {
    return 'Conectado'
  }

  if (state.connectionStatus === 'connecting') {
    return 'Conectando'
  }

  return 'Desconectado'
}

function getErrorMessage(error: unknown): string {
  if (error instanceof Error) {
    return error.message
  }

  return 'Erro inesperado.'
}

function formatDate(value: string): string {
  return new Intl.DateTimeFormat('pt-BR', {
    dateStyle: 'short',
    timeStyle: 'short',
  }).format(new Date(value))
}

function escapeHtml(value: string): string {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;')
}

async function initialize(): Promise<void> {
  if (!state.token) {
    render()
    return
  }

  state.isLoading = true
  render()

  try {
    await loadNotifications()
    connectWebSocket()
  } catch (error) {
    state.errorMessage = getErrorMessage(error)
  } finally {
    state.isLoading = false
    render()
  }
}

void initialize()
