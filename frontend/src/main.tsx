import { StrictMode } from "react"
import { createRoot } from "react-dom/client"
import "./style.css"

function App() {
  return (
    <main className="min-h-screen bg-slate-50 p-6 text-slate-950">
      <div className="mx-auto max-w-md rounded-xl border bg-white p-6 shadow-sm">
        <p className="text-xs font-bold uppercase tracking-wider text-blue-600">Notify</p>
        <h1 className="mt-2 text-3xl font-bold">Frontend com shadcn pronto</h1>
        <p className="mt-3 text-sm text-slate-600">
          Configuração base concluída. Na próxima parte vamos reconstruir login,
          cadastro, dashboard, notificações e conta usando componentes shadcn.
        </p>
      </div>
    </main>
  )
}

createRoot(document.getElementById("app")!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
