import { StrictMode } from "react"
import { createRoot } from "react-dom/client"

import { App } from "@/App"
import { Toaster } from "@/components/ui/sonner"
import "./style.css"

createRoot(document.getElementById("app")!).render(
  <StrictMode>
    <App />
    <Toaster richColors position="top-right" />
  </StrictMode>,
)
