import { FormEvent, ReactNode, useEffect, useMemo, useRef, useState } from "react"
import { toast } from "sonner"
import {
  Bell,
  CheckCircle2,
  CircleAlert,
  Clock3,
  History,
  Inbox,
  LogOut,
  MessageSquareReply,
  Radio,
  RefreshCcw,
  RotateCcw,
  Send,
  Settings,
  ShieldCheck,
  UserCircle,
  UserPlus,
} from "lucide-react"

import type { LoginResponse, Notification, PasswordResponse, User, WebSocketEvent } from "@/types"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Separator } from "@/components/ui/separator"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Textarea } from "@/components/ui/textarea"

const API_URL = "http://localhost:8080"
const WS_URL = "ws://localhost:8080"

type ConnectionStatus = "disconnected" | "connecting" | "connected"
type AppPage = "inbox" | "compose" | "sent" | "account"

function loadStoredUser(): User | null {
  const storedUser = localStorage.getItem("user")

  if (!storedUser) {
    return null
  }

  try {
    return JSON.parse(storedUser) as User
  } catch {
    localStorage.removeItem("user")
    return null
  }
}

function formatDate(value: string): string {
  return new Intl.DateTimeFormat("pt-BR", {
    dateStyle: "short",
    timeStyle: "short",
  }).format(new Date(value))
}

function formatDuration(from: string, to?: string | null): string {
  if (!to) {
    return "Ainda não concluída"
  }

  const start = new Date(from).getTime()
  const end = new Date(to).getTime()
  const diffInMinutes = Math.max(0, Math.round((end - start) / 60000))

  if (diffInMinutes < 1) {
    return "menos de 1 min"
  }

  if (diffInMinutes < 60) {
    return `${diffInMinutes} min`
  }

  const hours = Math.floor(diffInMinutes / 60)
  const minutes = diffInMinutes % 60

  if (minutes === 0) {
    return `${hours}h`
  }

  return `${hours}h ${minutes}min`
}

function getErrorMessage(error: unknown): string {
  if (error instanceof Error) {
    return error.message
  }

  return "Erro inesperado."
}

export function App() {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem("token"))
  const [user, setUser] = useState<User | null>(() => loadStoredUser())
  const [receivedNotifications, setReceivedNotifications] = useState<Notification[]>([])
  const [sentNotifications, setSentNotifications] = useState<Notification[]>([])
  const [connectionStatus, setConnectionStatus] = useState<ConnectionStatus>("disconnected")
  const [currentPage, setCurrentPage] = useState<AppPage>("inbox")
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [successMessage, setSuccessMessage] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  const socketRef = useRef<WebSocket | null>(null)

  const unreadCount = useMemo(
    () => receivedNotifications.filter((notification) => !notification.is_read).length,
    [receivedNotifications],
  )

  const completedCount = useMemo(
    () => receivedNotifications.filter((notification) => notification.is_completed).length,
    [receivedNotifications],
  )

  const roleLabel = user?.role === "admin" ? "Administrador" : "Usuário"

  const saveSession = (newToken: string, newUser: User) => {
    setToken(newToken)
    setUser(newUser)

    localStorage.setItem("token", newToken)
    localStorage.setItem("user", JSON.stringify(newUser))
  }

  const saveUser = (newUser: User) => {
    setUser(newUser)
    localStorage.setItem("user", JSON.stringify(newUser))
  }

  const clearSession = () => {
    setToken(null)
    setUser(null)
    setReceivedNotifications([])
    setSentNotifications([])
    setConnectionStatus("disconnected")
    setCurrentPage("inbox")
    setErrorMessage(null)
    setSuccessMessage(null)

    if (socketRef.current) {
      socketRef.current.close()
      socketRef.current = null
    }

    localStorage.removeItem("token")
    localStorage.removeItem("user")
  }

  async function apiRequest<T>(path: string, options: RequestInit = {}): Promise<T> {
    const headers = new Headers(options.headers)

    if (options.body && !headers.has("Content-Type")) {
      headers.set("Content-Type", "application/json")
    }

    if (token) {
      headers.set("Authorization", `Bearer ${token}`)
    }

    const response = await fetch(`${API_URL}${path}`, {
      ...options,
      headers,
    })

    if (response.status === 401 && token) {
      clearSession()
      throw new Error("Sessão expirada. Faça login novamente.")
    }

    if (!response.ok) {
      let message = "Erro inesperado na comunicação com a API."

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

  async function runWithFeedback(action: () => Promise<void>) {
    setIsLoading(true)
    setErrorMessage(null)
    setSuccessMessage(null)

    try {
      await action()
    } catch (error) {
      setErrorMessage(getErrorMessage(error))
    } finally {
      setIsLoading(false)
    }
  }

  function upsertReceivedNotification(notification: Notification) {
    setReceivedNotifications((currentNotifications) => upsertNotification(currentNotifications, notification))
  }

  function upsertSentNotification(notification: Notification) {
    setSentNotifications((currentNotifications) => upsertNotification(currentNotifications, notification))
  }

  function upsertNotification(currentNotifications: Notification[], notification: Notification) {
    const exists = currentNotifications.some((item) => item.id === notification.id)

    if (exists) {
      return currentNotifications.map((item) => (item.id === notification.id ? notification : item))
    }

    return [notification, ...currentNotifications]
  }

  async function loadReceivedNotifications() {
    const data = await apiRequest<Notification[]>("/notifications")
    setReceivedNotifications(data)
  }

  async function loadSentNotifications() {
    const data = await apiRequest<Notification[]>("/notifications/sent")
    setSentNotifications(data)
  }

  async function refreshAll() {
    await Promise.all([loadReceivedNotifications(), loadSentNotifications()])
  }

  function connectWebSocket(currentToken: string) {
    if (socketRef.current) {
      socketRef.current.close()
    }

    setConnectionStatus("connecting")

    const socket = new WebSocket(`${WS_URL}/ws?token=${encodeURIComponent(currentToken)}`)
    socketRef.current = socket

    socket.addEventListener("open", () => {
      setConnectionStatus("connected")
    })

    socket.addEventListener("message", (event) => {
      try {
        const websocketEvent = JSON.parse(event.data) as WebSocketEvent

        if (websocketEvent.type === "connected") {
          setConnectionStatus("connected")
          return
        }

        if (websocketEvent.type === "notification.created") {
          upsertReceivedNotification(websocketEvent.data)

          toast(websocketEvent.data.title, {
            description: `${websocketEvent.data.sender_name}: ${websocketEvent.data.message}`,
          })

          setSuccessMessage("Nova notificação recebida em tempo real.")
        }
      } catch {
        setErrorMessage("Não foi possível processar uma mensagem WebSocket.")
      }
    })

    socket.addEventListener("error", () => {
      setErrorMessage("Erro na conexão WebSocket.")
    })

    socket.addEventListener("close", () => {
      if (socketRef.current === socket) {
        socketRef.current = null
        setConnectionStatus("disconnected")
      }
    })
  }

  async function handleLogin(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const formData = new FormData(event.currentTarget)
    const email = String(formData.get("email") ?? "").trim()
    const password = String(formData.get("password") ?? "")

    await runWithFeedback(async () => {
      const response = await apiRequest<LoginResponse>("/auth/login", {
        method: "POST",
        body: JSON.stringify({ email, password }),
      })

      saveSession(response.token, response.user)
      setSuccessMessage("Login realizado com sucesso.")
    })
  }

  async function handleRegister(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const formData = new FormData(event.currentTarget)
    const name = String(formData.get("name") ?? "").trim()
    const email = String(formData.get("email") ?? "").trim()
    const password = String(formData.get("password") ?? "")

    await runWithFeedback(async () => {
      const response = await apiRequest<LoginResponse>("/auth/register", {
        method: "POST",
        body: JSON.stringify({ name, email, password }),
      })

      saveSession(response.token, response.user)
      setSuccessMessage("Conta criada com sucesso.")
    })
  }

  async function handleCreateNotification(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const form = event.currentTarget
    const formData = new FormData(form)
    const recipientEmail = String(formData.get("recipient_email") ?? "").trim()
    const title = String(formData.get("title") ?? "").trim()
    const message = String(formData.get("message") ?? "").trim()

    await runWithFeedback(async () => {
      const notification = await apiRequest<Notification>("/notifications", {
        method: "POST",
        body: JSON.stringify({
          recipient_email: recipientEmail,
          title,
          message,
        }),
      })

      upsertSentNotification(notification)

      if (notification.recipient_id === user?.id) {
        upsertReceivedNotification(notification)
      }

      form.reset()
      setSuccessMessage("Notificação enviada com sucesso.")
      setCurrentPage("sent")
    })
  }

  async function handleReply(notification: Notification, message: string) {
    await runWithFeedback(async () => {
      const reply = await apiRequest<Notification>(`/notifications/${notification.id}/replies`, {
        method: "POST",
        body: JSON.stringify({
          title: `Resposta: ${notification.title}`,
          message,
        }),
      })

      upsertSentNotification(reply)

      if (reply.recipient_id === user?.id) {
        upsertReceivedNotification(reply)
      }

      setSuccessMessage("Resposta enviada com sucesso.")
    })
  }

  async function handleUpdateProfile(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const formData = new FormData(event.currentTarget)
    const name = String(formData.get("name") ?? "").trim()
    const email = String(formData.get("email") ?? "").trim()

    await runWithFeedback(async () => {
      const updatedUser = await apiRequest<User>("/auth/me", {
        method: "PATCH",
        body: JSON.stringify({ name, email }),
      })

      saveUser(updatedUser)
      setSuccessMessage("Dados da conta atualizados.")
    })
  }

  async function handleChangePassword(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const form = event.currentTarget
    const formData = new FormData(form)
    const currentPassword = String(formData.get("current_password") ?? "")
    const newPassword = String(formData.get("new_password") ?? "")

    await runWithFeedback(async () => {
      await apiRequest<PasswordResponse>("/auth/password", {
        method: "PATCH",
        body: JSON.stringify({
          current_password: currentPassword,
          new_password: newPassword,
        }),
      })

      form.reset()
      setSuccessMessage("Senha alterada com sucesso.")
    })
  }

  async function handleUpdateReadStatus(notification: Notification) {
    const endpoint = notification.is_read
      ? `/notifications/${notification.id}/unread`
      : `/notifications/${notification.id}/read`

    await runWithFeedback(async () => {
      const updatedNotification = await apiRequest<Notification>(endpoint, {
        method: "PATCH",
      })

      upsertReceivedNotification(updatedNotification)

      setSuccessMessage(
        updatedNotification.is_read
          ? "Notificação marcada como lida."
          : "Notificação marcada como não lida.",
      )
    })
  }

  async function handleUpdateCompletionStatus(notification: Notification) {
    const endpoint = notification.is_completed
      ? `/notifications/${notification.id}/reopen`
      : `/notifications/${notification.id}/complete`

    await runWithFeedback(async () => {
      const updatedNotification = await apiRequest<Notification>(endpoint, {
        method: "PATCH",
      })

      upsertReceivedNotification(updatedNotification)

      setSuccessMessage(
        updatedNotification.is_completed
          ? "Notificação concluída."
          : "Notificação reaberta.",
      )
    })
  }

  useEffect(() => {
    if (!token) {
      return
    }

    let isMounted = true

    async function initializeSession() {
      setIsLoading(true)
      setErrorMessage(null)

      try {
        const currentUser = await apiRequest<User>("/auth/me")

        if (!isMounted) {
          return
        }

        saveUser(currentUser)
        await refreshAll()

        if (token) {
          connectWebSocket(token)
        }
      } catch (error) {
        if (isMounted) {
          setErrorMessage(getErrorMessage(error))
        }
      } finally {
        if (isMounted) {
          setIsLoading(false)
        }
      }
    }

    void initializeSession()

    return () => {
      isMounted = false

      if (socketRef.current) {
        socketRef.current.close()
        socketRef.current = null
      }
    }
  }, [token])

  if (!token) {
    return (
      <main className="min-h-screen bg-slate-50 px-4 py-8 text-slate-950 sm:grid sm:place-items-center">
        <Card className="mx-auto w-full max-w-md border-slate-200 shadow-xl shadow-slate-200/70">
          <CardHeader>
            <div className="mb-3 flex size-10 items-center justify-center rounded-xl bg-blue-600 text-white">
              <Bell className="size-5" />
            </div>
            <p className="text-xs font-bold uppercase tracking-[0.22em] text-blue-600">Notify</p>
            <CardTitle className="text-3xl">Acesse sua central</CardTitle>
            <CardDescription>
              Entre ou crie uma conta para enviar e receber notificações em tempo real.
            </CardDescription>
          </CardHeader>

          <CardContent>
            <Feedback errorMessage={errorMessage} successMessage={successMessage} />

            <Tabs defaultValue="login" className="w-full">
              <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="login">Entrar</TabsTrigger>
                <TabsTrigger value="register">Criar conta</TabsTrigger>
              </TabsList>

              <TabsContent value="login">
                <form className="mt-5 grid gap-4" onSubmit={handleLogin}>
                  <div className="grid gap-2">
                    <Label htmlFor="login-email">E-mail</Label>
                    <Input
                      id="login-email"
                      name="email"
                      type="email"
                      defaultValue="messederdaniel@outlook.com"
                      required
                    />
                  </div>

                  <div className="grid gap-2">
                    <Label htmlFor="login-password">Senha</Label>
                    <Input id="login-password" name="password" type="password" defaultValue="Teste@2026" required />
                  </div>

                  <Button type="submit" disabled={isLoading}>
                    {isLoading ? "Entrando..." : "Entrar"}
                  </Button>
                </form>

                <p className="mt-4 text-xs text-slate-500">
                  Admin: <strong>messederdaniel@outlook.com</strong> / <strong>Teste@2026</strong>
                </p>
              </TabsContent>

              <TabsContent value="register">
                <form className="mt-5 grid gap-4" onSubmit={handleRegister}>
                  <div className="grid gap-2">
                    <Label htmlFor="register-name">Nome</Label>
                    <Input id="register-name" name="name" type="text" placeholder="Seu nome" required />
                  </div>

                  <div className="grid gap-2">
                    <Label htmlFor="register-email">E-mail</Label>
                    <Input id="register-email" name="email" type="email" placeholder="voce@example.com" required />
                  </div>

                  <div className="grid gap-2">
                    <Label htmlFor="register-password">Senha</Label>
                    <Input
                      id="register-password"
                      name="password"
                      type="password"
                      placeholder="Mínimo de 6 caracteres"
                      minLength={6}
                      required
                    />
                  </div>

                  <Button type="submit" disabled={isLoading}>
                    <UserPlus className="mr-2 size-4" />
                    {isLoading ? "Criando..." : "Criar conta"}
                  </Button>
                </form>
              </TabsContent>
            </Tabs>
          </CardContent>
        </Card>
      </main>
    )
  }

  return (
    <main className="min-h-screen bg-slate-50 px-4 py-6 text-slate-950 md:px-8">
      <div className="mx-auto grid max-w-6xl gap-5">
        <header className="flex flex-col gap-4 rounded-2xl border bg-white p-5 shadow-sm xl:flex-row xl:items-center xl:justify-between">
          <div>
            <h1 className="text-4xl font-black tracking-tight text-blue-600 md:text-5xl">Notify</h1>
            <p className="mt-2 text-lg font-semibold text-slate-900">
              {user?.role === "admin" ? "Painel administrador" : "Painel do usuário"}
            </p>
            <p className="mt-1 max-w-2xl text-sm text-slate-600">
              Envie, receba, responda e acompanhe notificações em tempo real via WebSocket.
            </p>
          </div>

          <div className="flex flex-col gap-3 rounded-2xl border bg-slate-50 p-3 text-sm text-slate-600 sm:flex-row sm:items-center">
            <div className="min-w-0 sm:max-w-52">
              <div className="flex items-center gap-2">
                <strong className="block truncate text-slate-950">{user?.name}</strong>
                <Badge variant={user?.role === "admin" ? "default" : "secondary"}>{roleLabel}</Badge>
              </div>
              <span className="block truncate text-xs">{user?.email}</span>
            </div>

            <div className="flex flex-wrap gap-2">
              <NavIconButton
                label="Recebidas"
                active={currentPage === "inbox"}
                onClick={() => setCurrentPage("inbox")}
              >
                <Inbox className="size-4" />
              </NavIconButton>

              <NavIconButton
                label="Enviar"
                active={currentPage === "compose"}
                onClick={() => setCurrentPage("compose")}
              >
                <Send className="size-4" />
              </NavIconButton>

              <NavIconButton
                label="Enviadas"
                active={currentPage === "sent"}
                onClick={() => setCurrentPage("sent")}
              >
                <History className="size-4" />
              </NavIconButton>

              <NavIconButton
                label="Minha conta"
                active={currentPage === "account"}
                onClick={() => setCurrentPage("account")}
              >
                <UserCircle className="size-4" />
              </NavIconButton>

              <Button
                type="button"
                aria-label="Atualizar"
                title="Atualizar"
                variant="outline"
                className="size-9 rounded-full p-0"
                onClick={() => void runWithFeedback(refreshAll)}
              >
                <RefreshCcw className="size-4" />
              </Button>

              <Button
                type="button"
                aria-label="Sair"
                title="Sair"
                variant="outline"
                className="size-9 rounded-full p-0"
                onClick={clearSession}
              >
                <LogOut className="size-4" />
              </Button>
            </div>
          </div>
        </header>

        <Feedback errorMessage={errorMessage} successMessage={successMessage} />

        {currentPage === "inbox" ? (
          <InboxPage
            receivedNotifications={receivedNotifications}
            unreadCount={unreadCount}
            completedCount={completedCount}
            connectionStatus={connectionStatus}
            isLoading={isLoading}
            onUpdateReadStatus={handleUpdateReadStatus}
            onUpdateCompletionStatus={handleUpdateCompletionStatus}
            onReply={handleReply}
          />
        ) : null}

        {currentPage === "compose" ? (
          <ComposePage isLoading={isLoading} onCreateNotification={handleCreateNotification} />
        ) : null}

        {currentPage === "sent" ? <SentPage sentNotifications={sentNotifications} /> : null}

        {currentPage === "account" ? (
          <AccountPage
            user={user}
            isLoading={isLoading}
            onUpdateProfile={handleUpdateProfile}
            onChangePassword={handleChangePassword}
          />
        ) : null}
      </div>
    </main>
  )
}

function NavIconButton({
  label,
  active,
  onClick,
  children,
}: {
  label: string
  active: boolean
  onClick: () => void
  children: ReactNode
}) {
  return (
    <Button
      type="button"
      aria-label={label}
      title={label}
      variant={active ? "default" : "outline"}
      className="size-9 rounded-full p-0"
      onClick={onClick}
    >
      {children}
    </Button>
  )
}

function InboxPage({
  receivedNotifications,
  unreadCount,
  completedCount,
  connectionStatus,
  isLoading,
  onUpdateReadStatus,
  onUpdateCompletionStatus,
  onReply,
}: {
  receivedNotifications: Notification[]
  unreadCount: number
  completedCount: number
  connectionStatus: ConnectionStatus
  isLoading: boolean
  onUpdateReadStatus: (notification: Notification) => Promise<void>
  onUpdateCompletionStatus: (notification: Notification) => Promise<void>
  onReply: (notification: Notification, message: string) => Promise<void>
}) {
  return (
    <>
      <section className="grid gap-4 md:grid-cols-4">
        <SummaryCard
          title="WebSocket"
          value={getConnectionLabel(connectionStatus)}
          icon={<Radio className="size-5" />}
          status={connectionStatus}
        />
        <SummaryCard title="Recebidas" value={String(receivedNotifications.length)} icon={<Bell className="size-5" />} />
        <SummaryCard title="Não lidas" value={String(unreadCount)} icon={<CircleAlert className="size-5" />} />
        <SummaryCard title="Concluídas" value={String(completedCount)} icon={<CheckCircle2 className="size-5" />} />
      </section>

      <Card>
        <CardHeader>
          <CardTitle>Notificações recebidas</CardTitle>
          <CardDescription>
            Mensagens recebidas em tempo real. Ao chegar uma nova notificação, o Sonner aparece na tela.
          </CardDescription>
        </CardHeader>

        <CardContent>
          {receivedNotifications.length === 0 ? (
            <EmptyState title="Nenhuma notificação recebida" description="Quando alguém enviar uma notificação para você, ela aparecerá aqui." />
          ) : (
            <div className="divide-y rounded-xl border">
              {receivedNotifications.map((notification) => (
                <ReceivedNotificationItem
                  key={notification.id}
                  notification={notification}
                  isLoading={isLoading}
                  onUpdateReadStatus={onUpdateReadStatus}
                  onUpdateCompletionStatus={onUpdateCompletionStatus}
                  onReply={onReply}
                />
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </>
  )
}

function ComposePage({
  isLoading,
  onCreateNotification,
}: {
  isLoading: boolean
  onCreateNotification: (event: FormEvent<HTMLFormElement>) => Promise<void>
}) {
  return (
    <section className="grid gap-4 lg:grid-cols-[1.2fr_0.8fr]">
      <Card>
        <CardHeader>
          <CardTitle>Enviar notificação</CardTitle>
          <CardDescription>
            Informe o destinatário, escreva a mensagem e envie em tempo real via WebSocket.
          </CardDescription>
        </CardHeader>

        <CardContent>
          <form className="grid gap-5" onSubmit={onCreateNotification}>
            <div className="grid gap-2">
              <Label htmlFor="recipient-email">E-mail do destinatário</Label>
              <Input
                id="recipient-email"
                name="recipient_email"
                type="email"
                placeholder="barretodaniel11971@hotmail.com"
                required
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="notification-title">Título</Label>
              <Input id="notification-title" name="title" placeholder="Ex: Atualização necessária" required />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="notification-message">Mensagem</Label>
              <Textarea
                id="notification-message"
                name="message"
                rows={10}
                placeholder="Descreva a mensagem, orientação ou solicitação para o destinatário."
                required
              />
            </div>

            <Button type="submit" className="w-full sm:w-fit" disabled={isLoading}>
              <Send className="mr-2 size-4" />
              {isLoading ? "Enviando..." : "Enviar notificação"}
            </Button>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Fluxo em tempo real</CardTitle>
          <CardDescription>Como esta tela se conecta ao requisito do projeto.</CardDescription>
        </CardHeader>

        <CardContent className="grid gap-3 text-sm text-slate-600">
          <p>
            Se o destinatário estiver logado e conectado via WebSocket, ele recebe um Sonner imediatamente.
          </p>
          <p>
            Se estiver offline, a notificação fica persistida no PostgreSQL e aparece quando ele entrar.
          </p>
          <p>
            O remetente consegue acompanhar depois se a notificação foi lida, concluída e quanto tempo levou.
          </p>
        </CardContent>
      </Card>
    </section>
  )
}

function SentPage({ sentNotifications }: { sentNotifications: Notification[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Histórico de enviadas</CardTitle>
        <CardDescription>
          Acompanhe para quem você enviou, se foi lida, se foi concluída e o tempo até conclusão.
        </CardDescription>
      </CardHeader>

      <CardContent>
        {sentNotifications.length === 0 ? (
          <EmptyState title="Nenhuma notificação enviada" description="As notificações que você enviar aparecerão aqui." />
        ) : (
          <div className="divide-y rounded-xl border">
            {sentNotifications.map((notification) => (
              <SentNotificationItem key={notification.id} notification={notification} />
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  )
}

function AccountPage({
  user,
  isLoading,
  onUpdateProfile,
  onChangePassword,
}: {
  user: User | null
  isLoading: boolean
  onUpdateProfile: (event: FormEvent<HTMLFormElement>) => Promise<void>
  onChangePassword: (event: FormEvent<HTMLFormElement>) => Promise<void>
}) {
  return (
    <section className="grid gap-4 lg:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle>Dados da conta</CardTitle>
          <CardDescription>Atualize nome e e-mail usados para acessar o sistema.</CardDescription>
        </CardHeader>

        <CardContent>
          <form className="grid gap-4" onSubmit={onUpdateProfile}>
            <div className="grid gap-2">
              <Label htmlFor="profile-name">Nome</Label>
              <Input id="profile-name" name="name" defaultValue={user?.name ?? ""} required />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="profile-email">E-mail</Label>
              <Input id="profile-email" name="email" type="email" defaultValue={user?.email ?? ""} required />
            </div>

            <Button type="submit" disabled={isLoading}>
              <Settings className="mr-2 size-4" />
              Salvar dados
            </Button>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Segurança</CardTitle>
          <CardDescription>Altere a senha informando a senha atual.</CardDescription>
        </CardHeader>

        <CardContent>
          <form className="grid gap-4" onSubmit={onChangePassword}>
            <div className="grid gap-2">
              <Label htmlFor="current-password">Senha atual</Label>
              <Input id="current-password" name="current_password" type="password" required />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="new-password">Nova senha</Label>
              <Input id="new-password" name="new_password" type="password" minLength={6} required />
            </div>

            <Button type="submit" disabled={isLoading}>
              <ShieldCheck className="mr-2 size-4" />
              Alterar senha
            </Button>
          </form>
        </CardContent>
      </Card>
    </section>
  )
}

function ReceivedNotificationItem({
  notification,
  isLoading,
  onUpdateReadStatus,
  onUpdateCompletionStatus,
  onReply,
}: {
  notification: Notification
  isLoading: boolean
  onUpdateReadStatus: (notification: Notification) => Promise<void>
  onUpdateCompletionStatus: (notification: Notification) => Promise<void>
  onReply: (notification: Notification, message: string) => Promise<void>
}) {
  return (
    <article
      className={`grid gap-4 border-l-4 p-4 ${
        notification.is_completed
          ? "border-emerald-500 bg-white"
          : notification.is_read
            ? "border-blue-500 bg-white"
            : "border-amber-500 bg-amber-50/50"
      }`}
    >
      <div className="grid gap-2">
        <div className="flex flex-wrap items-center gap-2">
          <Badge variant={notification.parent_id ? "outline" : "default"}>
            {notification.parent_id ? "Resposta" : "Notificação"}
          </Badge>
          <Badge variant={notification.is_read ? "secondary" : "destructive"}>
            {notification.is_read ? "Lida" : "Não lida"}
          </Badge>
          <Badge variant={notification.is_completed ? "default" : "outline"}>
            {notification.is_completed ? "Concluída" : "Pendente"}
          </Badge>
          <time className="text-xs text-slate-500">{formatDate(notification.created_at)}</time>
        </div>

        <div>
          <h3 className="font-semibold">{notification.title}</h3>
          <p className="mt-1 text-sm text-slate-600">{notification.message}</p>
        </div>

        <div className="text-xs text-slate-500">
          De: <strong>{notification.sender_name}</strong> &lt;{notification.sender_email}&gt;
        </div>
      </div>

      <div className="flex flex-wrap gap-2">
        <Button variant="outline" disabled={isLoading} onClick={() => void onUpdateReadStatus(notification)}>
          {notification.is_read ? <RotateCcw className="mr-2 size-4" /> : <CheckCircle2 className="mr-2 size-4" />}
          {notification.is_read ? "Marcar como não lida" : "Marcar como lida"}
        </Button>

        <Button variant="outline" disabled={isLoading} onClick={() => void onUpdateCompletionStatus(notification)}>
          {notification.is_completed ? <RotateCcw className="mr-2 size-4" /> : <CheckCircle2 className="mr-2 size-4" />}
          {notification.is_completed ? "Reabrir" : "Concluir"}
        </Button>
      </div>

      <ReplyBox notification={notification} isLoading={isLoading} onReply={onReply} />
    </article>
  )
}

function SentNotificationItem({ notification }: { notification: Notification }) {
  return (
    <article className="grid gap-4 border-l-4 border-blue-500 p-4 sm:grid-cols-[1fr_auto] sm:items-start">
      <div className="grid gap-2">
        <div className="flex flex-wrap items-center gap-2">
          <Badge variant={notification.parent_id ? "outline" : "default"}>
            {notification.parent_id ? "Resposta enviada" : "Enviada"}
          </Badge>
          <Badge variant={notification.is_read ? "secondary" : "outline"}>
            {notification.is_read ? "Lida" : "Não lida"}
          </Badge>
          <Badge variant={notification.is_completed ? "default" : "outline"}>
            {notification.is_completed ? "Concluída" : "Pendente"}
          </Badge>
          <time className="text-xs text-slate-500">{formatDate(notification.created_at)}</time>
        </div>

        <div>
          <h3 className="font-semibold">{notification.title}</h3>
          <p className="mt-1 text-sm text-slate-600">{notification.message}</p>
        </div>

        <div className="text-xs text-slate-500">
          Para: <strong>{notification.recipient_name}</strong> &lt;{notification.recipient_email}&gt;
        </div>
      </div>

      <div className="rounded-xl border bg-slate-50 p-3 text-sm text-slate-600">
        <div className="flex items-center gap-2">
          <Clock3 className="size-4" />
          <span>{formatDuration(notification.created_at, notification.completed_at)}</span>
        </div>
      </div>
    </article>
  )
}

function ReplyBox({
  notification,
  isLoading,
  onReply,
}: {
  notification: Notification
  isLoading: boolean
  onReply: (notification: Notification, message: string) => Promise<void>
}) {
  const [isOpen, setIsOpen] = useState(false)
  const [message, setMessage] = useState("")

  async function submitReply(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const trimmedMessage = message.trim()
    if (!trimmedMessage) {
      return
    }

    await onReply(notification, trimmedMessage)
    setMessage("")
    setIsOpen(false)
  }

  if (!isOpen) {
    return (
      <Button variant="ghost" className="w-fit px-0 text-blue-600" onClick={() => setIsOpen(true)}>
        <MessageSquareReply className="mr-2 size-4" />
        Responder
      </Button>
    )
  }

  return (
    <form className="grid gap-3 rounded-xl border bg-white p-3" onSubmit={submitReply}>
      <Label htmlFor={`reply-${notification.id}`}>Responder notificação</Label>
      <Textarea
        id={`reply-${notification.id}`}
        value={message}
        onChange={(event) => setMessage(event.target.value)}
        placeholder="Escreva sua resposta ou dúvida"
        rows={4}
      />
      <div className="flex flex-wrap gap-2">
        <Button type="submit" disabled={isLoading || !message.trim()}>
          Enviar resposta
        </Button>
        <Button type="button" variant="outline" onClick={() => setIsOpen(false)}>
          Cancelar
        </Button>
      </div>
    </form>
  )
}

function Feedback({
  errorMessage,
  successMessage,
}: {
  errorMessage: string | null
  successMessage: string | null
}) {
  if (errorMessage) {
    return (
      <Alert variant="destructive" className="mb-4">
        <CircleAlert className="size-4" />
        <AlertDescription>{errorMessage}</AlertDescription>
      </Alert>
    )
  }

  if (successMessage) {
    return (
      <Alert className="mb-4 border-emerald-200 bg-emerald-50 text-emerald-800">
        <CheckCircle2 className="size-4" />
        <AlertDescription>{successMessage}</AlertDescription>
      </Alert>
    )
  }

  return null
}

function SummaryCard({
  title,
  value,
  icon,
  status,
}: {
  title: string
  value: string
  icon: ReactNode
  status?: ConnectionStatus
}) {
  return (
    <Card>
      <CardContent className="flex items-center justify-between p-5">
        <div>
          <p className="text-xs font-bold uppercase tracking-wider text-slate-500">{title}</p>
          <strong className="mt-2 flex items-center gap-2 text-2xl">
            {status ? <span className={`size-2 rounded-full ${getConnectionDotClass(status)}`} /> : null}
            {value}
          </strong>
        </div>

        <div className="rounded-xl bg-slate-100 p-3 text-slate-600">{icon}</div>
      </CardContent>
    </Card>
  )
}

function EmptyState({ title, description }: { title: string; description: string }) {
  return (
    <div className="rounded-xl border border-dashed p-8 text-center text-sm text-slate-500">
      <strong className="block text-slate-950">{title}</strong>
      {description}
    </div>
  )
}

function getConnectionLabel(status: ConnectionStatus): string {
  if (status === "connected") {
    return "Conectado"
  }

  if (status === "connecting") {
    return "Conectando"
  }

  return "Desconectado"
}

function getConnectionDotClass(status: ConnectionStatus): string {
  if (status === "connected") {
    return "bg-emerald-500"
  }

  if (status === "connecting") {
    return "bg-amber-500"
  }

  return "bg-slate-400"
}
