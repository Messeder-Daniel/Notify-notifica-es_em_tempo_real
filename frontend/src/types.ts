export type User = {
  id: string
  name: string
  email: string
  created_at?: string
}

export type LoginResponse = {
  token: string
  user: User
}

export type Notification = {
  id: string
  user_id: string
  title: string
  message: string
  is_read: boolean
  created_at: string
  read_at?: string | null
}

export type PasswordResponse = {
  message: string
}

export type WebSocketEvent =
  | {
      type: 'connected'
      user_id: string
      message: string
    }
  | {
      type: 'notification.created'
      data: Notification
    }
