export type UserRole = 'admin' | 'user'

export type User = {
  id: string
  name: string
  email: string
  role: UserRole
  created_at?: string
}

export type LoginResponse = {
  token: string
  user: User
}

export type Notification = {
  id: string
  sender_id: string
  recipient_id: string
  parent_id?: string | null
  title: string
  message: string
  is_read: boolean
  is_completed: boolean
  created_at: string
  read_at?: string | null
  completed_at?: string | null
  sender_name: string
  sender_email: string
  recipient_name: string
  recipient_email: string
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
